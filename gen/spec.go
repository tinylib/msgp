package gen

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"regexp"
	"strings"
)

const (
	lenAsUint32 = "uint32(len(%s))"
	literalFmt  = "%s"
	intFmt      = "%d"
	quotedFmt   = `"%s"`
	mapHeader   = "MapHeader"
	arrayHeader = "ArrayHeader"
	mapKey      = "MapKeyPtr"
	stringTyp   = "String"
	u32         = "uint32"
)

// Method is a bitfield representing something that the
// generator knows how to print.
type Method uint8

// are the bits in 'f' set in 'm'?
func (m Method) isset(f Method) bool { return (m&f == f) }

// String implements fmt.Stringer
func (m Method) String() string {
	switch m {
	case 0, invalidmeth:
		return "<invalid method>"
	case Decode:
		return "decode"
	case Encode:
		return "encode"
	case Marshal:
		return "marshal"
	case Unmarshal:
		return "unmarshal"
	case Size:
		return "size"
	case Test:
		return "test"
	default:
		// return e.g. "decode+encode+test"
		modes := [...]Method{Decode, Encode, Marshal, Unmarshal, Size, Test}
		any := false
		nm := ""
		for _, mm := range modes {
			if m.isset(mm) {
				if any {
					nm += "+" + mm.String()
				} else {
					nm += mm.String()
					any = true
				}
			}
		}
		return nm

	}
}

const (
	Decode      Method                       = 1 << iota // msgp.Decodable
	Encode                                               // msgp.Encodable
	Marshal                                              // msgp.Marshaler
	Unmarshal                                            // msgp.Unmarshaler
	Size                                                 // msgp.Sizer
	Test                                                 // generate tests
	invalidmeth                                          // this isn't a method
	encodetest  = Encode | Decode | Test                 // tests for Encodable and Decodable
	marshaltest = Marshal | Unmarshal | Test             // tests for Marshaler and Unmarshaler
)

type Printer struct {
	gens          []generator
	CompactFloats bool
	ClearOmitted  bool
	NewTime       bool
	AsUTC         bool
	ArrayLimit    uint32
	MapLimit      uint32
	MarshalLimits bool
	LimitPrefix   string
}

func NewPrinter(m Method, out io.Writer, tests io.Writer) *Printer {
	if m.isset(Test) && tests == nil {
		panic("cannot print tests with 'nil' tests argument!")
	}
	gens := make([]generator, 0, 7)
	if m.isset(Decode) {
		gens = append(gens, decode(out))
	}
	if m.isset(Encode) {
		gens = append(gens, encode(out))
	}
	if m.isset(Marshal) {
		gens = append(gens, marshal(out))
	}
	if m.isset(Unmarshal) {
		gens = append(gens, unmarshal(out))
	}
	if m.isset(Size) {
		gens = append(gens, sizes(out))
	}
	if m.isset(marshaltest) {
		gens = append(gens, mtest(tests))
	}
	if m.isset(encodetest) {
		gens = append(gens, etest(tests))
	}
	if len(gens) == 0 {
		panic("NewPrinter called with invalid method flags")
	}
	return &Printer{gens: gens}
}

// TransformPass is a pass that transforms individual
// elements. (Note that if the returned is different from
// the argument, it should not point to the same objects.)
type TransformPass func(Elem) Elem

// IgnoreTypename is a pass that just ignores
// types of a given name.
func IgnoreTypename(name string) TransformPass {
	if name, ok := strings.CutPrefix(name, "regex:"); ok {
		name = strings.TrimPrefix(name, "regex:")
		rx, err := regexp.Compile(name)
		if err != nil {
			panic(fmt.Sprintf("Error compiling ignore regex %q: %v", name, err))
		}
		return func(e Elem) Elem {
			if rx.MatchString(e.TypeName()) {
				return nil
			}
			return e
		}
	}
	return func(e Elem) Elem {
		if e.TypeName() == name {
			return nil
		}
		return e
	}
}

// ApplyDirective applies a directive to a named pass
// and all of its dependents.
func (p *Printer) ApplyDirective(pass Method, t TransformPass) {
	for _, g := range p.gens {
		if g.Method().isset(pass) {
			g.Add(t)
		}
	}
}

// Print prints an Elem.
func (p *Printer) Print(e Elem) error {
	e.SetIsAllowNil(false)
	for _, g := range p.gens {
		// Elem.SetVarname() is called before the Print() step in parse.FileSet.PrintTo().
		// Elem.SetVarname() generates identifiers as it walks the Elem. This can cause
		// collisions between idents created during SetVarname and idents created during Print,
		// hence the separate prefixes.
		resetIdent("zb")
		err := g.Execute(e, Context{
			compFloats:             p.CompactFloats,
			clearOmitted:           p.ClearOmitted,
			newTime:                p.NewTime,
			asUTC:                  p.AsUTC,
			arrayLimit:             p.ArrayLimit,
			mapLimit:               p.MapLimit,
			marshalLimits:          p.MarshalLimits,
			limitPrefix:            p.LimitPrefix,
			currentFieldArrayLimit: math.MaxUint32, // Initialize to "no field limit"
			currentFieldMapLimit:   math.MaxUint32, // Initialize to "no field limit"
		})
		resetIdent("za")

		if err != nil {
			return err
		}
	}
	return nil
}

type contextItem interface {
	Arg() string
}

type contextString string

func (c contextString) Arg() string {
	return fmt.Sprintf("%q", c)
}

type contextVar string

func (c contextVar) Arg() string {
	return string(c)
}

type Context struct {
	path                   []contextItem
	compFloats             bool
	clearOmitted           bool
	newTime                bool
	asUTC                  bool
	arrayLimit             uint32
	mapLimit               uint32
	marshalLimits          bool
	limitPrefix            string
	currentFieldArrayLimit uint32 // Current field's array limit (0 = no field-level limit)
	currentFieldMapLimit   uint32 // Current field's map limit (0 = no field-level limit)
}

func (c *Context) PushString(s string) {
	c.path = append(c.path, contextString(s))
}

func (c *Context) PushVar(s string) {
	c.path = append(c.path, contextVar(s))
}

func (c *Context) Pop() {
	c.path = c.path[:len(c.path)-1]
}

// SetFieldLimits sets the current field-specific limits for the context
func (c *Context) SetFieldLimits(arrayLimit, mapLimit uint32) {
	c.currentFieldArrayLimit = arrayLimit
	c.currentFieldMapLimit = mapLimit
}

// ClearFieldLimits clears the current field-specific limits (use file limits)
func (c *Context) ClearFieldLimits() {
	c.currentFieldArrayLimit = math.MaxUint32
	c.currentFieldMapLimit = math.MaxUint32
}

func (c *Context) ArgsStr() string {
	var out string
	for idx, p := range c.path {
		if idx > 0 {
			out += ", "
		}
		out += p.Arg()
	}
	return out
}

// generator is the interface through
// which code is generated.
type generator interface {
	Method() Method
	Add(p TransformPass)
	Execute(Elem, Context) error // execute writes the method for the provided object.
}

type passes []TransformPass

func (p *passes) Add(t TransformPass) {
	*p = append(*p, t)
}

func (p *passes) applyall(e Elem) Elem {
	for _, t := range *p {
		e = t(e)
		if e == nil {
			return nil
		}
	}
	return e
}

type traversal interface {
	gMap(*Map)
	gSlice(*Slice)
	gArray(*Array)
	gPtr(*Ptr)
	gBase(*BaseElem)
	gStruct(*Struct)
}

// type-switch dispatch to the correct
// method given the type of 'e'
func next(t traversal, e Elem) {
	switch e := e.(type) {
	case *Map:
		t.gMap(e)
	case *Struct:
		t.gStruct(e)
	case *Slice:
		t.gSlice(e)
	case *Array:
		t.gArray(e)
	case *Ptr:
		t.gPtr(e)
	case *BaseElem:
		t.gBase(e)
	default:
		panic("bad element type")
	}
}

// possibly-immutable method receiver
func imutMethodReceiver(p Elem) string {
	typeName := p.BaseTypeName()
	typeParams := p.TypeParams()

	switch e := p.(type) {
	case *Struct:
		// TODO(HACK): actually do real math here.
		if len(e.Fields) <= 3 {
			for i := range e.Fields {
				if be, ok := e.Fields[i].FieldElem.(*BaseElem); !ok || (be.Value == IDENT || be.Value == Bytes) {
					goto nope
				}
			}
			return typeName + typeParams.TypeParams
		}
	nope:
		return "*" + typeName + typeParams.TypeParams

	// gets dereferenced automatically
	case *Array:
		return "*" + typeName + typeParams.TypeParams

	// everything else can be
	// by-value.
	default:
		return typeName + typeParams.TypeParams
	}
}

// if necessary, wraps a type
// so that its method receiver
// is of the write type.
func methodReceiver(p Elem) string {
	typeName := p.BaseTypeName()
	typeParams := p.TypeParams()

	switch p.(type) {

	// structs and arrays are
	// dereferenced automatically,
	// so no need to alter varname
	case *Struct, *Array:
		return "*" + typeName + typeParams.TypeParams
	// set variable name to
	// *varname
	default:
		p.SetVarname("(*" + p.Varname() + ")")
		return "*" + typeName + typeParams.TypeParams
	}
}

func unsetReceiver(p Elem) {
	switch p.(type) {
	case *Struct, *Array:
	default:
		p.SetVarname("z")
	}
}

// shared utility for generators
type printer struct {
	w   io.Writer
	err error
}

// writes "var {{name}} {{typ}};"
func (p *printer) declare(name string, typ string) {
	p.printf("\nvar %s %s", name, typ)
}

// does:
//
//	if m == nil {
//	    m = make(type, size)
//	} else if len(m) > 0 {
//
//	    for key := range m { delete(m, key) }
//	}
func (p *printer) resizeMap(size string, m *Map) {
	vn := m.Varname()
	if !p.ok() {
		return
	}
	p.printf("\nif %s == nil {", vn)
	p.printf("\n%s = make(%s, %s)", vn, m.TypeName(), size)
	p.printf("\n} else if len(%s) > 0 {", vn)
	p.clearMap(vn)
	p.closeblock()
}

// CanAutoShim contains the primitives that can be auto-shimmed.
var CanAutoShim = map[Primitive]bool{Uint: true, Uint8: true, Uint16: true, Uint32: true, Uint64: true, Int: true, Int8: true, Int16: true, Int32: true, Int64: true, Bool: true, Float32: true, Float64: true, Byte: true}

// assign key to value based on varnames
func (p *printer) mapAssign(m *Map) {
	if !p.ok() {
		return
	}

	if key, ok := m.Key.(*BaseElem); ok {
		fromBase := key.FromBase()
		shimErr := key.ShimErrs
		if m.AutoMapShims && CanAutoShim[key.Value] {
			fromBase = "msgp.AutoShim{}.Parse" + key.Value.String()
			shimErr = true
		}
		if !m.AllowBinMaps && fromBase != "" {
			if key.Value == String && key.ShimFromBase != "" {
				p.printf("\nvar %sTmp %s", m.Keyidx, key.TypeName())
				if key.ShimErrs {
					p.printf("\n%sTmp, err = %s(%s)", m.Keyidx, fromBase, m.Keyidx)
					p.wrapErrCheck("\"shim: " + fromBase + "\"")
				} else {
					p.printf("\n%sTmp = %s(%s)", m.Keyidx, fromBase, m.Keyidx)
				}
				p.printf("\n%s[%sTmp] = %s", m.Varname(), m.Keyidx, m.Validx)
				return
			} else if key.Value == IDENT {
				p.printf("\n%s[%s(%s)] = %s", m.Varname(), fromBase, m.Keyidx, m.Validx)
				return
			} else {
				if shimErr {
					p.printf("\nvar %sTmp %s", m.Keyidx, strings.ToLower(key.Value.String()))
					p.printf("\n%sTmp, err = %s(%s)", m.Keyidx, fromBase, m.Keyidx)
					p.wrapErrCheck("\"shim: " + m.Varname() + "\"")
					p.printf("\n%s[%s(%sTmp)] = %s", m.Varname(), key.FromBase(), m.Keyidx, m.Validx)
				} else {
					p.printf("\n%s[%s(%s)] = %s", m.Varname(), fromBase, m.Keyidx, m.Validx)
				}
				return
			}
		}
	}
	p.printf("\n%s[%s] = %s", m.Varname(), m.Keyidx, m.Validx)
}

// clear map keys
func (p *printer) clearMap(name string) {
	p.printf("\nclear(%[1]s)", name)
}

func (p *printer) wrapErrCheck(ctx string) {
	p.print("\nif err != nil {")
	if ctx != "" {
		p.printf("\nerr = msgp.WrapError(err, %s)", ctx)
	} else {
		p.print("\nerr = msgp.WrapError(err)")
	}
	p.printf("\nreturn")
	p.print("\n}")
}

func (p *printer) resizeSlice(size string, s *Slice) {
	p.printf("\nif cap(%[1]s) >= int(%[2]s) { %[1]s = (%[1]s)[:%[2]s] } else { %[1]s = make(%[3]s, %[2]s) }", s.Varname(), size, s.TypeName())
}

// resizeSliceNoNil will resize a slice and will not allow nil slices.
func (p *printer) resizeSliceNoNil(size string, s *Slice) {
	p.printf("\nif %[1]s != nil && cap(%[1]s) >= int(%[2]s) {", s.Varname(), size)
	p.printf("\n%[1]s = (%[1]s)[:%[2]s]", s.Varname(), size)
	p.printf("\n} else { %[1]s = make(%[3]s, %[2]s) }", s.Varname(), size, s.TypeName())
}

func (p *printer) arrayCheck(want string, got string) {
	p.printf("\nif %[1]s != %[2]s { err = msgp.ArrayError{Wanted: %[2]s, Got: %[1]s}; return }", got, want)
}

func (p *printer) closeblock() { p.print("\n}") }

// does:
//
//	for idx := range iter {
//	    {{generate inner}}
//	}
func (p *printer) rangeBlock(ctx *Context, idx string, iter string, t traversal, inner Elem) {
	ctx.PushVar(idx)
	// Tags on slices do not extend to the elements, so we always disable allownil on elements.
	// If we want this to happen in the future, it should be a unique tag.
	inner.SetIsAllowNil(false)
	p.printf("\n for %s := range %s {", idx, iter)
	next(t, inner)
	p.closeblock()
	ctx.Pop()
}

func (p *printer) nakedReturn() {
	if p.ok() {
		p.print("\nreturn\n}\n")
	}
}

func (p *printer) comment(s string) {
	p.print("\n// " + s)
}

func (p *printer) printf(format string, args ...any) {
	if p.err == nil {
		_, p.err = fmt.Fprintf(p.w, format, args...)
	}
}

func (p *printer) print(format string) {
	if p.err == nil {
		_, p.err = io.WriteString(p.w, format)
	}
}

func (p *printer) initPtr(pt *Ptr) {
	if pt.Needsinit() {
		vname := pt.Varname()
		p.printf("\nif %s == nil { %s = new(%s); }", vname, vname, pt.Value.TypeName())
	}
}

func (p *printer) ok() bool { return p.err == nil }

func tobaseConvert(b *BaseElem) string {
	return b.ToBase() + "(" + b.Varname() + ")"
}

func (p *printer) varWriteMapHeader(receiver string, sizeVarname string, maxSize int) {
	if maxSize <= 15 {
		p.printf("\nerr = %s.Append(0x80 | uint8(%s))", receiver, sizeVarname)
	} else {
		p.printf("\nerr = %s.WriteMapHeader(%s)", receiver, sizeVarname)
	}
}

func (p *printer) varAppendMapHeader(sliceVarname string, sizeVarname string, maxSize int) {
	if maxSize <= 15 {
		p.printf("\n%s = append(%s, 0x80 | uint8(%s))", sliceVarname, sliceVarname, sizeVarname)
	} else {
		p.printf("\n%s = msgp.AppendMapHeader(%s, %s)", sliceVarname, sliceVarname, sizeVarname)
	}
}

// bmask is a bitmask of a the specified number of bits
type bmask struct {
	bitlen  int
	varname string
}

// typeDecl returns the variable declaration as a var statement
func (b *bmask) typeDecl() string {
	return fmt.Sprintf("var %s %s /* %d bits */", b.varname, b.typeName(), b.bitlen)
}

// typeName returns the type, e.g. "uint8" or "[2]uint64"
func (b *bmask) typeName() string {
	if b.bitlen <= 8 {
		return "uint8"
	}
	if b.bitlen <= 16 {
		return "uint16"
	}
	if b.bitlen <= 32 {
		return "uint32"
	}
	if b.bitlen <= 64 {
		return "uint64"
	}

	return fmt.Sprintf("[%d]uint64", (b.bitlen+64-1)/64)
}

// readExpr returns the expression to read from a position in the bitmask.
// Compare ==0 for false or !=0 for true.
func (b *bmask) readExpr(bitoffset int) string {
	if bitoffset < 0 || bitoffset >= b.bitlen {
		panic(fmt.Errorf("bitoffset %d out of range for bitlen %d", bitoffset, b.bitlen))
	}

	var buf bytes.Buffer
	buf.Grow(len(b.varname) + 16)
	buf.WriteByte('(')
	buf.WriteString(b.varname)
	if b.bitlen > 64 {
		fmt.Fprintf(&buf, "[%d]", (bitoffset / 64))
	}
	buf.WriteByte('&')
	fmt.Fprintf(&buf, "0x%X", (uint64(1) << (uint64(bitoffset) % 64)))
	buf.WriteByte(')')

	return buf.String()
}

// setStmt returns the statement to set the specified bit in the bitmask.
func (b *bmask) setStmt(bitoffset int) string {
	var buf bytes.Buffer
	buf.Grow(len(b.varname) + 16)
	buf.WriteString(b.varname)
	if b.bitlen > 64 {
		fmt.Fprintf(&buf, "[%d]", (bitoffset / 64))
	}
	fmt.Fprintf(&buf, " |= 0x%X", (uint64(1) << (uint64(bitoffset) % 64)))

	return buf.String()
}

// notAllSet returns a check against all fields having been set in set.
func (b *bmask) notAllSet() string {
	var buf bytes.Buffer
	buf.Grow(len(b.varname) + 16)
	buf.WriteString(b.varname)
	if b.bitlen > 64 {
		var bytes []string
		remain := b.bitlen
		for remain >= 8 {
			bytes = append(bytes, "0xff")
		}
		if remain > 0 {
			bytes = append(bytes, fmt.Sprintf("0x%X", remain))
		}
		fmt.Fprintf(&buf, " != [%d]byte{%s}\n", (b.bitlen+63)/64, strings.Join(bytes, ","))
	}
	fmt.Fprintf(&buf, " != 0x%x", uint64(1<<b.bitlen)-1)

	return buf.String()
}
