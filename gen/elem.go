package gen

import (
	"fmt"
	"math/rand"
	"strings"
)

const (
	idxChars = "abcdefghijlkmnopqrstuvwxyz"
	idxLen   = 3
)

// generate a random index variable name
func randIdx() string {
	bts := make([]byte, idxLen)
	for i := range bts {
		bts[i] = idxChars[rand.Intn(len(idxChars))]
	}
	return string(bts)
}

// This code defines the template
// syntax tree. If the input were:
//
// type Marshaler struct {
// 	  Thing1 *float64 `msg:"thing1"`
// 	  Body   []byte   `msg:"body"`
// }
//
// then the AST parsing of *Marshaler should produce:
//
// var val Elem = &Ptr{
// 	name: "z",
// 	Value: &Struct{
// 		Name: "Marshaler",
// 		Fields: []StructField{
// 			{
// 				FieldTag: "thing1",
// 				FieldElem: &Ptr{
// 					name: "z.thing1",
// 					Value: &BaseElem{
// 						name:    "*z.thing1",
// 						Value:   Float64,
//						Convert: false,
// 					},
// 				},
// 			},
// 			{
// 				FieldTag: "",
// 				FieldElem: &BaseElem{
// 					name:    "z.body",
// 					Value:   Bytes,
// 					Convert: False,
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
	Time // time.Time
	Ext  // extension

	IDENT // IDENT means an unrecognized identifier
)

// ElemType is one of Ptr, Map
// Slice, Struct, or Base
type ElemType int

const (
	InvalidType ElemType = iota
	PtrType              // pointer-to-object
	SliceType            // slice-of-object
	StructType           // struct-of-objects
	BaseType             // object
	MapType              // map[string]object
	ArrayType            // [Size]object
)

// Elem is Base, Slice, Stuct, Map, or Ptr
type Elem interface {
	Type() ElemType // Returns element type

	// We need the following
	// because objects cannot
	// be type-asserted in templates.
	// Yes, this is gross.
	Ptr() *Ptr
	Slice() *Slice
	Struct() *Struct
	Base() *BaseElem
	Map() *Map
	Array() *Array

	// SetVarname sets this nodes
	// variable name and recursively
	// sets the names of all its children
	SetVarname(s string)

	// Varname is the variable
	// name of the node
	Varname() string

	// TypeName is the canonical
	// go type name of the node
	// e.g. "string", "int", "map[string]float64"
	TypeName() string

	// fmt.Stringer for debugging
	String() string
}

type Array struct {
	name  string // Varname
	Index string // index variable name
	Size  string // array size
	Els   Elem   // child
}

func (a *Array) Type() ElemType  { return ArrayType }
func (a *Array) Ptr() *Ptr       { return nil }
func (a *Array) Slice() *Slice   { return nil }
func (a *Array) Struct() *Struct { return nil }
func (a *Array) Base() *BaseElem { return nil }
func (a *Array) Map() *Map       { return nil }
func (a *Array) Array() *Array   { return a }
func (a *Array) SetVarname(s string) {
	a.name = s
ridx:
	a.Index = randIdx()

	// try to avoid using the same
	// index as a parent slice
	if strings.Contains(a.name, a.Index) {
		goto ridx
	}

	a.Els.SetVarname(fmt.Sprintf("%s[%s]", a.name, a.Index))
}
func (a *Array) Varname() string  { return a.name }
func (a *Array) TypeName() string { return fmt.Sprintf("[%s]%s", a.Size, a.Els.TypeName()) }
func (a *Array) String() string {
	return fmt.Sprintf("Array[%s]Of(%s - %s)", a.Size, a.Els.String(), a.Varname())
}

// Map is a map[string]Elem
type Map struct {
	name   string
	Keyidx string // key variable name
	Validx string // value variable name
	Value  Elem
}

func (m *Map) Type() ElemType  { return MapType }
func (m *Map) Ptr() *Ptr       { return nil }
func (m *Map) Slice() *Slice   { return nil }
func (m *Map) Struct() *Struct { return nil }
func (m *Map) Base() *BaseElem { return nil }
func (m *Map) Map() *Map       { return m }
func (m *Map) Array() *Array   { return nil }
func (m *Map) SetVarname(s string) {
	m.name = s
ridx:
	m.Keyidx = randIdx()
	m.Validx = randIdx()

	// just in case
	if m.Keyidx == m.Validx {
		goto ridx
	}

	m.Value.SetVarname(m.Validx)
}
func (m *Map) Varname() string  { return m.name }
func (m *Map) TypeName() string { return fmt.Sprintf("map[string]%s", m.Value.TypeName()) }
func (m *Map) String() string {
	return fmt.Sprintf("MapOf([string]%s - %s)", m.Value.String(), m.Varname())
}

type Slice struct {
	name  string
	Index string
	Els   Elem // The type of each element
}

func (s *Slice) Type() ElemType  { return SliceType }
func (s *Slice) Ptr() *Ptr       { return nil }
func (s *Slice) Slice() *Slice   { return s }
func (s *Slice) Struct() *Struct { return nil }
func (s *Slice) Base() *BaseElem { return nil }
func (s *Slice) Map() *Map       { return nil }
func (s *Slice) Array() *Array   { return nil }
func (s *Slice) SetVarname(a string) {
	s.name = a
	s.Index = randIdx()
	s.Els.SetVarname(fmt.Sprintf("%s[%s]", s.name, s.Index))
}
func (s *Slice) Varname() string  { return s.name }
func (s *Slice) TypeName() string { return "[]" + s.Els.TypeName() }
func (s *Slice) String() string {
	return fmt.Sprintf("SliceOf(%s - %s)", s.Els.String(), s.Varname())
}

type Ptr struct {
	name  string
	Value Elem
}

func (s *Ptr) Type() ElemType  { return PtrType }
func (s *Ptr) Ptr() *Ptr       { return s }
func (s *Ptr) Slice() *Slice   { return nil }
func (s *Ptr) Struct() *Struct { return nil }
func (s *Ptr) Base() *BaseElem { return nil }
func (s *Ptr) Map() *Map       { return nil }
func (s *Ptr) Array() *Array   { return nil }
func (s *Ptr) SetVarname(a string) {
	s.name = a

	// struct fields are dereferenced
	// automatically...
	switch s.Value.Type() {
	case StructType:
		// struct fields are automatically dereferenced
		s.Value.SetVarname(a)
		return

	case BaseType:
		// identities and extensions have pointer receivers
		if s.Value.Base().IsIdent() {
			s.Value.SetVarname(a)
			return
		}

		fallthrough
	default:
		s.Value.SetVarname("*" + a)
		return
	}
}
func (s *Ptr) Varname() string  { return s.name }
func (s *Ptr) TypeName() string { return "*" + s.Value.TypeName() }
func (s *Ptr) String() string {
	return fmt.Sprintf("PointerTo(%s - %s)", s.Value.String(), s.Varname())
}

type Struct struct {
	Name   string
	Fields []StructField
}

func (s *Struct) Type() ElemType  { return StructType }
func (s *Struct) Ptr() *Ptr       { return nil }
func (s *Struct) Slice() *Slice   { return nil }
func (s *Struct) Struct() *Struct { return s }
func (s *Struct) Base() *BaseElem { return nil }
func (s *Struct) Map() *Map       { return nil }
func (s *Struct) Array() *Array   { return nil }
func (s *Struct) Varname() string { return "" } // structs are special
func (s *Struct) SetVarname(a string) {
	writeStructFields(s.Fields, a)
}
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
	name         string
	Value        Base
	Ident        string // IDENT name if unresolved
	Convert      bool   // should we do an explicit conversion?
	ShimToBase   string // shim to base type
	ShimFromBase string // shim from base type
}

func (s *BaseElem) Type() ElemType  { return BaseType }
func (s *BaseElem) Ptr() *Ptr       { return nil }
func (s *BaseElem) Slice() *Slice   { return nil }
func (s *BaseElem) Struct() *Struct { return nil }
func (s *BaseElem) Map() *Map       { return nil }
func (s *BaseElem) Base() *BaseElem { return s }
func (s *BaseElem) Array() *Array   { return nil }
func (s *BaseElem) Varname() string { return s.name }
func (s *BaseElem) SetVarname(a string) {

	// extensions are assumed
	// to have pointer receivers, so
	// we need to *not* dereference it
	// (if it's a pointer) OR we need
	// to take a reference
	if s.Value == Ext {
		if strings.HasPrefix(a, "*") {
			s.name = strings.TrimPrefix(a, "*")
		} else {
			s.name = "&" + a
		}
		return

		// if we're using a shim
	}

	s.name = a
}

func (s *BaseElem) String() string { return fmt.Sprintf("(%s - %s)", s.BaseName(), s.Varname()) }

// TypeName returns the syntactically correct Go
// type name for the base element.
func (s *BaseElem) TypeName() string {
	if s.Ident != "" {
		return s.Ident
	}
	return s.BaseType()
}

// ToBase, used if Convert==true, is used as tmp = {{ToBase}}({{Varname}})
func (s *BaseElem) ToBase() string {
	if s.ShimToBase != "" {
		return s.ShimToBase
	}
	return s.BaseType()
}

// FromBase, used if Convert==true, is used as {{Varname}} = {{FromBase}}(tmp)
func (s *BaseElem) FromBase() string {
	if s.ShimFromBase != "" {
		return s.ShimFromBase
	}
	return s.Ident
}

// BaseName returns the string form of the
// base type (e.g. Float64, Ident, etc)
func (s *BaseElem) BaseName() string {
	// time is a special case;
	// we strip the package prefix
	if s.Value == Time {
		return "Time"
	}
	return s.Value.String()
}

func (s *BaseElem) BaseType() string {
	switch s.Value {
	case IDENT:
		return s.Ident

	// exceptions to the naming/capitalization
	// rule:
	case Intf:
		return "interface{}"
	case Bytes:
		return "[]byte"
	case Time:
		return "time.Time"
	case Ext:
		return "msgp.Extension"

	// everything else is base.String() with
	// the first letter as lowercase
	default:
		return strings.ToLower(s.BaseName())
	}
}

// is this an interface{} ?
func (s *BaseElem) IsIntf() bool { return s.Value == Intf }

// is this an extension?
func (s *BaseElem) IsExt() bool { return s.Value == Ext }

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
	case Time:
		return "time.Time"
	case Ext:
		return "Extension"
	case IDENT:
		return "Ident"
	default:
		return "INVALID"
	}
}

// writeStructFields is a trampoline for writeBase for
// all of the fields in a struct
func writeStructFields(s []StructField, name string) {
	for i := range s {
		s[i].FieldElem.SetVarname(fmt.Sprintf("%s.%s", name, s[i].FieldName))
	}
}
