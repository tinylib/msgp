package gen

import (
	//"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	mssType reflect.Type
	msiType reflect.Type
	btsType reflect.Type
)

func init() {
	var mss map[string]string
	mssType = reflect.TypeOf(mss)
	var msi map[string]interface{}
	msiType = reflect.TypeOf(msi)
	var bts []byte
	btsType = reflect.TypeOf(bts)
}

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
	MapStrStr  // map[string]string
	MapStrIntf // map[string]interface{}
	Intf       // interface{}

	IDENT
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
	String() string
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
func (s *Ptr) TypeName() string { return s.Value.TypeName() }
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
	Ident   string
}

func (s *BaseElem) Type() ElemType  { return BaseType }
func (s *BaseElem) Ptr() *Ptr       { return nil }
func (s *BaseElem) Slice() *Slice   { return nil }
func (s *BaseElem) Struct() *Struct { return nil }
func (s *BaseElem) Base() *BaseElem { return s }
func (s *BaseElem) String() string  { return fmt.Sprintf("(%s - %s)", s.BaseName(), s.Varname) }

// TypeName returns the syntactically correct Go
// type name for the base element.
func (s *BaseElem) TypeName() string {
	switch s.Value {
	case IDENT:
		return s.Ident
	case MapStrIntf:
		return "map[string]interface{}"
	case MapStrStr:
		return "map[string]string"
	default:
		return strings.ToLower(s.BaseName())
	}
}

// BaseName returns the string form of the
// base type (e.g. Float64, Ident, MapStrStr, etc)
func (s *BaseElem) BaseName() string { return s.Value.String() }

// is this a map?
func (s *BaseElem) IsMap() bool { return (s.Value == MapStrStr || s.Value == MapStrIntf) }

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
	case MapStrStr:
		return "MapStrStr"
	case MapStrIntf:
		return "MapStrIntf"
	case Intf:
		return "Intf"
	case IDENT:
		return "Ident"
	default:
		return "INVALID"
	}
}

// construct make the type tree
// out of a reflect.Type. it doesn't try
// to insert variable names, in most cases. additionally,
// some fields may be nil.
func construct(t reflect.Type) Elem {
	switch t.Kind() {
	case reflect.Struct:
		out := &Struct{
			Name:   t.Name(),
			Fields: make([]StructField, t.NumField()),
		}
		for i := range out.Fields {
			field := t.Field(i)
			tag := field.Tag.Get("msg")
			if tag == "" {
				tag = field.Tag.Get("json")
				if tag == "" {
					tag = field.Name
				}
			}
			out.Fields[i] = StructField{
				FieldName: field.Name,
				FieldTag:  tag,
				FieldElem: construct(field.Type),
			}
		}
		return out

	case reflect.Map:
		if t.Key().Kind() != reflect.String {
			return nil
		}
		switch t.Elem().Kind() {
		case reflect.String:
			return &BaseElem{Value: MapStrStr}
		case reflect.Interface:
			return &BaseElem{Value: MapStrIntf}
		default:
			return nil
		}

	case reflect.String:
		return &BaseElem{Value: String}
	case reflect.Slice, reflect.Array:
		if t.AssignableTo(btsType) {
			return &BaseElem{Value: Bytes}
		} else {
			return &Slice{Els: construct(t.Elem())}
		}
	case reflect.Float32:
		return &BaseElem{Value: Float32}
	case reflect.Float64:
		return &BaseElem{Value: Float64}
	case reflect.Int:
		return &BaseElem{Value: Int}
	case reflect.Int8:
		return &BaseElem{Value: Int8}
	case reflect.Int16:
		return &BaseElem{Value: Int16}
	case reflect.Int32:
		return &BaseElem{Value: Int32}
	case reflect.Int64:
		return &BaseElem{Value: Int64}
	case reflect.Uint:
		return &BaseElem{Value: Uint}
	case reflect.Uint8:
		return &BaseElem{Value: Uint8}
	case reflect.Uint16:
		return &BaseElem{Value: Uint16}
	case reflect.Uint32:
		return &BaseElem{Value: Uint32}
	case reflect.Uint64:
		return &BaseElem{Value: Uint64}
	case reflect.Complex64:
		return &BaseElem{Value: Complex64}
	case reflect.Complex128:
		return &BaseElem{Value: Complex128}
	case reflect.Ptr:
		return &Ptr{Value: construct(t.Elem())}
	default:
		return nil
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
		}
	}
}

// writeBase recursively writes variable names
// on pointers, slices, and base types.
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
		case StructType:
			// struct fields are dereferenced automatically
			writeStructFields(b.Struct().Fields, parent.Ptr().Varname)
		}
	}
}

func GenerateTree(v reflect.Type) (Elem, error) {
	var err error
	el := construct(v)
	Propogate(el.Ptr(), "z")
	return el, err
}
