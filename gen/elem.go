package gen

import (
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
)

type ElemType int

const (
	InvalidType ElemType = iota
	PtrType
	SliceType
	StructType
	BaseType
)

// Elem is Base, Slice, Stuct, or Ptr
type Elem interface {
	Type() ElemType // Returns element type
	Ptr() *Ptr
	Slice() *Slice
	Struct() *Struct
	Base() *BaseElem
	TypeName() string
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
func (s *Slice) TypeName() string { return s.Els.TypeName() }

type Ptr struct {
	Varname string
	Value   Elem
}

func (s *Ptr) Type() ElemType   { return PtrType }
func (s *Ptr) Ptr() *Ptr        { return s }
func (s *Ptr) Slice() *Slice    { return nil }
func (s *Ptr) Struct() *Struct  { return nil }
func (s *Ptr) Base() *BaseElem  { return nil }
func (s *Ptr) TypeName() string { return s.Value.TypeName() }

type Struct struct {
	Name   string
	Fields []StructField
}

func (s *Struct) Type() ElemType   { return StructType }
func (s *Struct) Ptr() *Ptr        { return nil }
func (s *Struct) Slice() *Slice    { return nil }
func (s *Struct) Struct() *Struct  { return s }
func (s *Struct) Base() *BaseElem  { return nil }
func (s *Struct) TypeName() string { return s.Name }

type StructField struct {
	FieldTag  string
	FieldElem Elem
}

type BaseElem struct {
	Varname string
	Value   Base
}

func (s *BaseElem) Type() ElemType   { return BaseType }
func (s *BaseElem) Ptr() *Ptr        { return nil }
func (s *BaseElem) Slice() *Slice    { return nil }
func (s *BaseElem) Struct() *Struct  { return nil }
func (s *BaseElem) Base() *BaseElem  { return s }
func (s *BaseElem) String() string   { return s.Value.String() }
func (s *BaseElem) TypeName() string { return strings.ToLower(s.String()) }

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
	default:
		return "INVALID"
	}
}
