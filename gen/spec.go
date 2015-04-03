package gen

import (
	"fmt"
	"io"
)

var (
	NopGen Generator = nopgen{}
)

// no-op Generator
type nopgen struct{}

func (n nopgen) Execute(_ Elem) error { return nil }

const (
	errcheck    = "\nif err != nil { return }"
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

// Method is a method that the generator
// knows how to generate
type Method uint8

const (
	Decode      Method = iota // msgp.Decodable
	Encode                    // msgp.Encodable
	Marshal                   // msgp.Marshaler
	Unmarshal                 // msgp.Unmarshaler
	Size                      // msgp.Sizer
	MarshalTest               // test Marshaler and Unmarshaler
	EncodeTest                // test Encodable and Decodable
)

// New returns a new Generator that generates
// a particular method and writes its body
// to the provided io.Writer. The output
// may not be in canonical Go format.
func New(m Method, out io.Writer) Generator {
	switch m {
	case Decode:
		return &decodeGen{printer{w: out}, false}
	case Encode:
		return &encodeGen{printer{w: out}}
	case Marshal:
		return &marshalGen{printer{w: out}}
	case Unmarshal:
		return &unmarshalGen{printer{w: out}, false}
	case Size:
		return &sizeGen{printer{w: out}, assign}
	case MarshalTest:
		return mtestGen{out}
	case EncodeTest:
		return etestGen{out}
	}
	panic("bad Method type")
}

// Generator is the interface through
// which code is generated.
type Generator interface {
	// Execute writes the method for the provided object.
	Execute(Elem) error
}

// every recursive type generator
// implements this
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

// if necessary, wraps a type
// so that its method receiver
// is of the write type.
func methodReceiver(p Elem, encodeFunc bool) string {
	star := "*"
	if encodeFunc && p.EncodeValueReceivers() {
		star = ""
	}

	switch p.(type) {

	// structs and arrays are
	// dereferenced automatically,
	// so no need to alter varname
	case *Struct, *Array:
		return star + p.TypeName()

	// set variable name to
	// *varname
	default:
		p.SetVarname("(" + star + p.Varname() + ")")
		return star + p.TypeName()
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
	p.printf("\nif cap(%[1]s) >= int(%[2]s) { %[1]s = %[1]s[:%[2]s] } else { %[1]s = make(%[3]s, %[2]s) }", s.Varname(), size, s.TypeName())
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
