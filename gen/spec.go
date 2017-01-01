package gen

import (
	"fmt"
	"io"
	"math"

	"github.com/tinylib/msgp/internal/log"
	"github.com/tinylib/msgp/msgp"
)

const (
	errcheck    = "\nif err != nil { return }"
	lenAsUint32 = "uint32(len(%s))"
	literalFmt  = "%s"
	intFmt      = "%d"
	quotedFmt   = `"%s"`
	mapHeader   = "MapHeader"
	arrayHeader = "ArrayHeader"
	mapKey      = "MapKey"
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

func strtoMeth(s string) Method {
	switch s {
	case "encode":
		return Encode
	case "decode":
		return Decode
	case "marshal":
		return Marshal
	case "unmarshal":
		return Unmarshal
	case "size":
		return Size
	case "test":
		return Test
	default:
		return 0
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
	gens []generator
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
	for _, g := range p.gens {
		err := g.Execute(e)
		if err != nil {
			return err
		}
	}
	return nil
}

// generator is the interface through
// which code is generated.
type generator interface {
	Method() Method
	Add(p TransformPass)
	Execute(Elem) error // execute writes the method for the provided object.
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

type fields struct {
	cache map[string]bool
}

func (f *fields) declareOnce(p printer, name, typ string) {
	key := name + "." + typ
	if f.cache == nil {
		f.cache = make(map[string]bool)
	} else if f.cache[key] {
		return
	}

	p.printf("\nvar %s %s;", name, typ)
	f.cache[key] = true
}

func (f *fields) drop() {
	for k := range f.cache {
		delete(f.cache, k)
	}
}

type traversal interface {
	gMap(*Map)
	gSlice(*Slice)
	gArray(*Array)
	gPtr(*Ptr)
	gBase(*BaseElem)
	gStruct(*Struct)
}

type declarer interface {
	declareOnce(p printer, name, typ string)
}

type assigner interface {
	assignAndCheck(name, base string)
	nextTypeAndCheck(name string)
	skipAndCheck()
}

type fuser interface {
	Fuse([]byte)
}

type traversalAssigner interface {
	traversal
	assigner
	declarer
}

type traversalFuser interface {
	traversal
	fuser
}

func genStructFieldsSerializer(t traversalFuser, p printer, fields []StructField) {
	data := msgp.AppendMapHeader(nil, uint32(len(fields)))
	p.printf("\n// map header, size %d", len(fields))
	t.Fuse(data)
	for _, f := range fields {
		if !p.ok() {
			return
		}

		data, p.err = msgp.AppendIntf(nil, f.FieldTag)
		p.printf(
			"\n// [field %q] write label `%v` as msgp.%s",
			f.FieldName, f.FieldTag, msgp.NextType(data),
		)
		t.Fuse(data)

		next(t, f.FieldElem)
	}
}

func genStructFieldsParser(t traversalAssigner, p printer, fields []StructField) {
	const (
		fieldBytes = "fieldBytes"
		fieldInt   = "fieldInt"
		fieldUint  = "fieldUint"
		typ        = "typ"
	)

	groups := groupFieldsByType(fields)
	hasUint := len(groups[msgp.UintType]) > 0
	hasInt := len(groups[msgp.IntType]) > 0
	hasStr := len(groups[msgp.StrType]) > 0
	singleType := len(groups) == 1

	if hasStr {
		t.declareOnce(p, fieldBytes, "[]byte")
	}
	if hasUint {
		t.declareOnce(p, fieldUint, "uint64")

		// Append to int fields also uint fields that do not overflow int64.
		// This is necessary because uint field with value <= (1<<7)-1 could be serialized as fixint.
		// and become a msgp.IntType. This is done for best compatibility with other libraries
		// (and with other languages). That is, some endpoint could serialize uint16 key with value <= (1<<7)-1 as
		// real msgpack uint16, but also could serialize it like fixint.
		for _, f := range groups[msgp.UintType] {
			v := f.FieldTag.(uint64)
			if v <= math.MaxInt64 {
				groups[msgp.IntType] = append(groups[msgp.IntType], f)
				hasInt = true
			}
		}
	}
	if hasInt {
		t.declareOnce(p, fieldInt, "int64")
	}
	if !singleType || hasUint {
		t.declareOnce(p, typ, "msgp.Type")
	}

	sz := randIdent()
	p.declare(sz, u32)
	t.assignAndCheck(sz, mapHeader)
	p.printf("\nfor %s > 0 {\n%s--", sz, sz)
	switch {
	case singleType && hasStr:
		t.assignAndCheck(fieldBytes, mapKey)
		switchFieldKeysStr(t, p, fields, fieldBytes)

	case singleType && !hasUint && hasInt:
		t.assignAndCheck(fieldInt, "Int64")
		switchFieldKeys(t, p, fields, fieldInt)

	default:
		// switch on inferred type of next field
		t.nextTypeAndCheck(typ)
		p.printf("\nswitch %s {", typ)
		if hasUint {
			p.print("\ncase msgp.UintType:")
			t.assignAndCheck(fieldUint, "Uint64")
			switchFieldKeys(t, p, groups[msgp.UintType], fieldUint)
		}
		if hasInt {
			p.print("\ncase msgp.IntType:")
			t.assignAndCheck(fieldInt, "Int64")
			switchFieldKeys(t, p, groups[msgp.IntType], fieldInt)
		}
		if hasStr {
			// double case is done for backward compatibility with previous implementation
			p.print("\ncase msgp.StrType, msgp.BinType:")
			t.assignAndCheck(fieldBytes, mapKey)
			switchFieldKeysStr(t, p, groups[msgp.StrType], fieldBytes)
		}
		p.print("\ndefault:")
		t.skipAndCheck()
		p.closeblock() // close switch
	}
	p.closeblock() // close loop
}

func groupFieldsByType(fields []StructField) map[msgp.Type][]StructField {
	groups := make(map[msgp.Type][]StructField, len(fields))
	for _, f := range fields {
		var t msgp.Type
		switch f.FieldTag.(type) {
		case int, int8, int16, int32, int64:
			t = msgp.IntType
		case uint, uint8, uint16, uint32, uint64:
			t = msgp.UintType
		case string:
			t = msgp.StrType
		default:
			log.Fatalf(
				"could not generate code to work with field's %q label: has unknown type %T",
				f.FieldName, f.FieldTag,
			)
		}
		groups[t] = append(groups[t], f)
	}
	return groups
}

func switchFieldKeysStr(t traversalAssigner, p printer, fields []StructField, label string) {
	p.printf("\nswitch msgp.UnsafeString(%s) {", label)
	for _, f := range fields {
		p.printf("\ncase \"%s\":", f.FieldTag)
		next(t, f.FieldElem)
		if !p.ok() {
			return
		}
	}
	p.print("\ndefault:")
	t.skipAndCheck()
	p.closeblock() // close switch
}

func switchFieldKeys(t traversalAssigner, p printer, fields []StructField, label string) {
	p.printf("\nswitch %s {", label)
	for _, f := range fields {
		p.printf("\ncase %v:", f.FieldTag)
		next(t, f.FieldElem)
		if !p.ok() {
			return
		}
	}
	p.print("\ndefault:")
	t.skipAndCheck()
	p.closeblock() // close switch
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
	switch e := p.(type) {
	case *Struct:
		// TODO(HACK): actually do real math here.
		if len(e.Fields) <= 3 {
			for i := range e.Fields {
				if be, ok := e.Fields[i].FieldElem.(*BaseElem); !ok || (be.Value == IDENT || be.Value == Bytes) {
					goto nope
				}
			}
			return p.TypeName()
		}
	nope:
		return "*" + p.TypeName()

	// gets dereferenced automatically
	case *Array:
		return "*" + p.TypeName()

	// everything else can be
	// by-value.
	default:
		return p.TypeName()
	}
}

// if necessary, wraps a type
// so that its method receiver
// is of the write type.
func methodReceiver(p Elem) string {
	switch p.(type) {

	// structs and arrays are
	// dereferenced automatically,
	// so no need to alter varname
	case *Struct, *Array:
		return "*" + p.TypeName()
	// set variable name to
	// *varname
	default:
		p.SetVarname("(*" + p.Varname() + ")")
		return "*" + p.TypeName()
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
// if m != nil && size > 0 {
//     m = make(type, size)
// } else if len(m) > 0 {
//     for key, _ := range m { delete(m, key) }
// }
//
func (p *printer) resizeMap(size string, m *Map) {
	vn := m.Varname()
	if !p.ok() {
		return
	}
	p.printf("\nif %s == nil && %s > 0 {", vn, size)
	p.printf("\n%s = make(%s, %s)", vn, m.TypeName(), size)
	p.printf("\n} else if len(%s) > 0 {", vn)
	p.clearMap(vn)
	p.closeblock()
}

// assign key to value based on varnames
func (p *printer) mapAssign(m *Map) {
	if !p.ok() {
		return
	}
	p.printf("\n%s[%s] = %s", m.Varname(), m.Keyidx, m.Validx)
}

// clear map keys
func (p *printer) clearMap(name string) {
	p.printf("\nfor key, _ := range %[1]s { delete(%[1]s, key) }", name)
}

func (p *printer) resizeSlice(size string, s *Slice) {
	p.printf("\nif cap(%[1]s) >= int(%[2]s) { %[1]s = (%[1]s)[:%[2]s] } else { %[1]s = make(%[3]s, %[2]s) }", s.Varname(), size, s.TypeName())
}

func (p *printer) arrayCheck(want string, got string) {
	p.printf("\nif %[1]s != %[2]s { err = msgp.ArrayError{Wanted: %[2]s, Got: %[1]s}; return }", got, want)
}

func (p *printer) closeblock() { p.print("\n}") }

// does:
//
// for idx := range iter {
//     {{generate inner}}
// }
//
func (p *printer) rangeBlock(idx string, iter string, t traversal, inner Elem) {
	p.printf("\n for %s := range %s {", idx, iter)
	next(t, inner)
	p.closeblock()
}

func (p *printer) nakedReturn() {
	if p.ok() {
		p.print("\nreturn\n}\n")
	}
}

func (p *printer) comment(s string) {
	p.print("\n// " + s)
}

func (p *printer) printf(format string, args ...interface{}) {
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
