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

// This code defines the type declaration tree.
//
// Consider the following:
//
// type Marshaler struct {
// 	  Thing1 *float64 `msg:"thing1"`
// 	  Body   []byte   `msg:"body"`
// }
//
// A parser using this generator as a backend
// should parse the above into:
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
type Primitive uint8

// this is effectively the
// list of currently available
// ReadXxxx / WriteXxxx methods.
const (
	Invalid Primitive = iota
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

// Elem is a go type capable of being
// serialized into MessagePack. It is
// implemented by *Ptr, *Struct, *Array,
// *Slice, *Map, and *BaseElem
type Elem interface {
	// SetVarname sets this nodes
	// variable name and recursively
	// sets the names of all its children.
	// In general, this should only be
	// called on the parent of the tree.
	SetVarname(s string)

	// Varname returns the variable
	// name of the element.
	Varname() string

	// TypeName is the canonical
	// go type name of the node
	// e.g. "string", "int", "map[string]float64"
	TypeName() string

	hidden()
}

type Array struct {
	name  string // Varname
	Index string // index variable name
	Size  string // array size
	Els   Elem   // child
}

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
func (a *Array) Varname() string { return a.name }

func (a *Array) TypeName() string { return fmt.Sprintf("[%s]%s", a.Size, a.Els.TypeName()) }

func (a *Array) hidden() {}

// Map is a map[string]Elem
type Map struct {
	name   string // varname
	Keyidx string // key variable name
	Validx string // value variable name
	Value  Elem   // value element
}

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

func (m *Map) Varname() string { return m.name }

func (m *Map) TypeName() string { return fmt.Sprintf("map[string]%s", m.Value.TypeName()) }

func (m *Map) hidden() {}

type Slice struct {
	name  string
	Index string
	Els   Elem // The type of each element
}

func (s *Slice) SetVarname(a string) {
	s.name = a
	s.Index = randIdx()
	s.Els.SetVarname(fmt.Sprintf("%s[%s]", s.name, s.Index))
}

func (s *Slice) Varname() string { return s.name }

func (s *Slice) TypeName() string { return "[]" + s.Els.TypeName() }

func (s *Slice) hidden() {}

type Ptr struct {
	name  string
	Value Elem
}

func (s *Ptr) SetVarname(a string) {
	s.name = a

	// struct fields are dereferenced
	// automatically...
	switch x := s.Value.(type) {
	case *Struct:
		// struct fields are automatically dereferenced
		x.SetVarname(a)
		return

	case *BaseElem:
		// identities and extensions have pointer receivers
		if x.Value == IDENT {
			x.SetVarname(a)
		} else {
			x.SetVarname("*" + a)
		}
		return

	default:
		s.Value.SetVarname("*" + a)
		return
	}
}

func (s *Ptr) Varname() string { return s.name }

func (s *Ptr) TypeName() string { return "*" + s.Value.TypeName() }

func (s *Ptr) hidden() {}

type Struct struct {
	vname   string        // varname
	Name    string        // struct type name
	Fields  []StructField // field list
	AsTuple bool          // write as an array instead of a map
}

func (s *Struct) Varname() string { return s.vname }

func (s *Struct) SetVarname(a string) {
	s.vname = a
	writeStructFields(s.Fields, a)
}

func (s *Struct) TypeName() string { return s.Name }

func (s *Struct) hidden() {}

type StructField struct {
	FieldTag  string
	FieldName string
	FieldElem Elem
}

// BaseElem is an element that
// can be represented by a primitive
// MessagePack type.
type BaseElem struct {
	name         string    // varname
	Ident        string    // IDENT name if unresolved, or empty
	ShimToBase   string    // shim to base type, or empty
	ShimFromBase string    // shim from base type, or empty
	Value        Primitive // Type of element
	Convert      bool      // should we do an explicit conversion?
}

func (s *BaseElem) hidden() {}

func (s *BaseElem) Varname() string { return s.name }

func (s *BaseElem) SetVarname(a string) {
	// extensions whose parents
	// are not pointers need to
	// be explicitly referenced
	if s.Value == Ext {
		if strings.HasPrefix(a, "*") {
			s.name = a[1:]
		} else {
			s.name = "&" + a
		}
		return
	}

	s.name = a
}

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

func (k Primitive) String() string {
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
