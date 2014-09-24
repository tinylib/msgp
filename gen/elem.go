package gen

import (
	"fmt"
	"strings"
)

// This code defines the template
// syntax tree. If the input were:
//
// type Marshaler struct {
// 	thing1 *float64
// 	body   []byte
// }
//
// then the AST parsing of *Marshaler should produce:
//
// var val Elem = &Ptr{
// 	Varname: "z",
// 	Value: &Struct{
// 		Name: "Marshaler",
// 		Fields: []StructField{
// 			{
// 				FieldTag: "thing1",
// 				FieldElem: &Ptr{
// 					Varname: "z.thing1",
// 					Value: &BaseElem{
// 						Varname: "*z.thing1",
// 						Value:   Float64,
// 					},
// 				},
// 			},
// 			{
// 				FieldTag: "",
// 				FieldElem: &BaseElem{
// 					Varname: "z.body",
// 					Value:   Bytes,
// 				},
// 			},
// 		},
// 	},
// }

// Base is one of the
// base types
type Base int

// this is effectively the
// list of currently available
// ReadXxxx / WriteXxxx methods.
const (
	Invalid Base = iota
	Bytes
	String
	Float32
	Float64
	Complex64
	Complex128
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Byte
	Int
	Int8
	Int16
	Int32
	Int64
	Bool
	Intf // interface{}

	IDENT // IDENT means an unrecognized identifier
)

// ElemType is one of Ptr, Map
// Slice, Struct, or Base
type ElemType int

const (
	InvalidType ElemType = iota
	PtrType
	SliceType
	StructType
	BaseType
	MapType
)

// Elem is Base, Slice, Stuct, Map, or Ptr
type Elem interface {
	Type() ElemType // Returns element type
	Ptr() *Ptr
	Slice() *Slice
	Struct() *Struct
	Base() *BaseElem
	Map() *Map
	TypeName() string
	String() string
}

// Map is a map[string]Elem
type Map struct {
	Varname string
	Value   Elem
}

func (m *Map) Type() ElemType   { return MapType }
func (m *Map) Ptr() *Ptr        { return nil }
func (m *Map) Slice() *Slice    { return nil }
func (m *Map) Struct() *Struct  { return nil }
func (m *Map) Base() *BaseElem  { return nil }
func (m *Map) Map() *Map        { return m }
func (m *Map) TypeName() string { return fmt.Sprintf("map[string]%s", m.Value.TypeName()) }
func (m *Map) String() string {
	return fmt.Sprintf("MapOf([string]%s - %s)", m.Value.String(), m.Varname)
}

type Slice struct {
	Varname string
	Els     Elem // The type of each element
}

func (s *Slice) Type() ElemType   { return SliceType }
func (s *Slice) Ptr() *Ptr        { return nil }
func (s *Slice) Slice() *Slice    { return s }
func (s *Slice) Struct() *Struct  { return nil }
func (s *Slice) Base() *BaseElem  { return nil }
func (s *Slice) Map() *Map        { return nil }
func (s *Slice) TypeName() string { return "[]" + s.Els.TypeName() }
func (s *Slice) String() string {
	return fmt.Sprintf("SliceOf(%s - %s)", s.Els.String(), s.Varname)
}

type Ptr struct {
	Varname string
	Value   Elem
}

func (s *Ptr) Type() ElemType   { return PtrType }
func (s *Ptr) Ptr() *Ptr        { return s }
func (s *Ptr) Slice() *Slice    { return nil }
func (s *Ptr) Struct() *Struct  { return nil }
func (s *Ptr) Base() *BaseElem  { return nil }
func (s *Ptr) Map() *Map        { return nil }
func (s *Ptr) TypeName() string { return "*" + s.Value.TypeName() }
func (s *Ptr) String() string {
	return fmt.Sprintf("PointerTo(%s - %s)", s.Value.String(), s.Varname)
}

type Struct struct {
	Name   string
	Fields []StructField
}

func (s *Struct) Type() ElemType   { return StructType }
func (s *Struct) Ptr() *Ptr        { return nil }
func (s *Struct) Slice() *Slice    { return nil }
func (s *Struct) Struct() *Struct  { return s }
func (s *Struct) Base() *BaseElem  { return nil }
func (s *Struct) Map() *Map        { return nil }
func (s *Struct) TypeName() string { return s.Name }
func (s *Struct) String() string {
	return fmt.Sprintf("%s{%s}", s.Name, s.Fields)
}

type StructField struct {
	FieldTag  string
	FieldName string
	FieldElem Elem
}

func (s StructField) String() string {
	return fmt.Sprintf("\n\t%s: %s %q, ", s.FieldName, s.FieldElem, s.FieldTag)
}

type BaseElem struct {
	Varname string
	Value   Base
	Ident   string // IDENT name if unresolved
	Convert bool   // should we do an explicit conversion?
}

func (s *BaseElem) Type() ElemType  { return BaseType }
func (s *BaseElem) Ptr() *Ptr       { return nil }
func (s *BaseElem) Slice() *Slice   { return nil }
func (s *BaseElem) Struct() *Struct { return nil }
func (s *BaseElem) Map() *Map       { return nil }
func (s *BaseElem) Base() *BaseElem { return s }
func (s *BaseElem) String() string  { return fmt.Sprintf("(%s - %s)", s.BaseName(), s.Varname) }

// TypeName returns the syntactically correct Go
// type name for the base element.
func (s *BaseElem) TypeName() string {
	if s.Ident != "" {
		return s.Ident
	}
	return s.BaseType()
}

// BaseName returns the string form of the
// base type (e.g. Float64, Ident, etc)
func (s *BaseElem) BaseName() string { return s.Value.String() }
func (s *BaseElem) BaseType() string {
	switch s.Value {
	case IDENT:
		return s.Ident
	case Intf:
		return "interface{}"
	case Bytes:
		return "[]byte"
	default:
		return strings.ToLower(s.BaseName())
	}
}

// is this an interface{}
func (s *BaseElem) IsIntf() bool { return s.Value == Intf }

// is this an external identity?
func (s *BaseElem) IsIdent() bool { return s.Value == IDENT }

func (k Base) String() string {
	switch k {
	case String:
		return "String"
	case Bytes:
		return "Bytes"
	case Float32:
		return "Float32"
	case Float64:
		return "Float64"
	case Complex64:
		return "Complex64"
	case Complex128:
		return "Complex128"
	case Uint:
		return "Uint"
	case Uint8:
		return "Uint8"
	case Uint16:
		return "Uint16"
	case Uint32:
		return "Uint32"
	case Uint64:
		return "Uint64"
	case Byte:
		return "Byte"
	case Int:
		return "Int"
	case Int8:
		return "Int8"
	case Int16:
		return "Int16"
	case Int32:
		return "Int32"
	case Int64:
		return "Int64"
	case Bool:
		return "Bool"
	case Intf:
		return "Intf"
	case IDENT:
		return "Ident"
	default:
		return "INVALID"
	}
}

// propNames propogates names through a *Ptr-to-*Struct
func Propogate(p *Ptr, name string) {
	if p == nil || p.Value == nil || p.Value.Struct() == nil {
		panic("propogate called on non-Ptr-to-Struct")
	}
	p.Varname = name
	writeStructFields(p.Value.Struct().Fields, name)
}

// writeStructFields is a trampoline for writeBase for
// all of the fields in a struct
func writeStructFields(s []StructField, name string) {
	prefix := name + "."
	for i := range s {
		switch s[i].FieldElem.Type() {

		case BaseType:
			s[i].FieldElem.Base().Varname = prefix + s[i].FieldName

		case StructType:
			writeStructFields(s[i].FieldElem.Struct().Fields, prefix+s[i].FieldName)

		case SliceType:
			s[i].FieldElem.Slice().Varname = prefix + s[i].FieldName
			writeBase(s[i].FieldElem.Slice().Els, s[i].FieldElem.Slice())

		case PtrType:
			s[i].FieldElem.Ptr().Varname = prefix + s[i].FieldName
			writeBase(s[i].FieldElem.Ptr().Value, s[i].FieldElem.Ptr())

		case MapType:
			s[i].FieldElem.Map().Varname = prefix + s[i].FieldName
			writeBase(s[i].FieldElem.Map().Value, s[i].FieldElem.Map())

		default:
			panic("unrecognized type!")
		}
	}
}

// writeBase recursively writes variable names
// on pointers, slices, and base types, using the name
// of the parent to name the child
func writeBase(b Elem, parent Elem) {
	switch parent.Type() {
	case SliceType:
		switch b.Type() {
		case BaseType:
			b.Base().Varname = parent.Slice().Varname + "[i]"
		case SliceType:
			b.Slice().Varname = parent.Slice().Varname + "[i]"
			writeBase(b.Slice().Els, b.Slice())
		case PtrType:
			b.Ptr().Varname = parent.Slice().Varname + "[i]"
			writeBase(b.Ptr().Value, b.Ptr())
		case StructType:
			writeStructFields(b.Struct().Fields, parent.Slice().Varname+"[i]")
		case MapType:
			b.Map().Varname = parent.Slice().Varname + "[i]"
			writeBase(b.Map().Value, b.Map())
		}
	case PtrType:
		switch b.Type() {
		case BaseType:
			if b.Base().IsIdent() {
				// SPECIAL CASE
				// WE ASSUME POINTER IDENTS
				// IMPLEMENT THE DECODER/ENCODER INTERFACES.
				b.Base().Varname = parent.Ptr().Varname
			} else {
				b.Base().Varname = "*" + parent.Ptr().Varname
			}
		case SliceType:
			b.Slice().Varname = "*" + parent.Ptr().Varname
			writeBase(b.Slice().Els, b.Slice())
		case PtrType:
			b.Ptr().Varname = "*" + parent.Ptr().Varname
			writeBase(b.Ptr().Value, b.Ptr())
		case MapType:
			b.Map().Varname = "*" + parent.Ptr().Varname
			writeBase(b.Map().Value, b.Map())
		case StructType:
			// struct fields are dereferenced automatically
			writeStructFields(b.Struct().Fields, parent.Ptr().Varname)
		}

	case MapType:
		switch b.Type() {
		case BaseType:
			b.Base().Varname = "val"
		case PtrType:
			b.Ptr().Varname = "val"
			writeBase(b.Ptr().Value, b.Ptr())
		case SliceType:
			b.Slice().Varname = "val"
			writeBase(b.Slice().Els, b.Slice())
		case MapType:
			b.Map().Varname = "val"
			writeBase(b.Map().Value, b.Map())
		case StructType:
			writeStructFields(b.Struct().Fields, "val")
		}
	}
}
