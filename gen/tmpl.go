package gen

import (
	"text/template"
)

// Templates
var (
	WriteMsgTempl *template.Template
)

// Family is the family
// that a type belongs to.
// (This is distinct from reflect.Kind,
// because Bytes has its own "kind".)
type Kind byte

const (
	Bytes Kind = iota
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
	Nil
	Struct
)

func (k Kind) String() string {
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
	case Nil:
		return "Nil"
	case Struct:
		return "Struct"
	default:
		return "INVALID"
	}
}

type StructDef struct {
	StructName string
	Fields     []StructField
}

type StructField struct {
	FieldName string // struct name + field name, e.g. "b.Val"
	FieldTag  string // field string name (`msg:"tag"`)
	IsSlice   bool   // is slice (or array)
	IsPtr     bool   // is ptr?
	IsStruct  bool   // is struct? (Kind == Struct)
	Kind      Kind
	KindName  string // Kind.String()

	// IsSlice and IsPtr can both be true if
	// the slice is a slice-of-pointers. Not sure
	// how we'll handle pointer-to-slice or pointer-
	// to-slice-of-pointers yet, among other things.
}

func init() {
	var err error
	WriteMsgTempl, err = template.ParseGlob(WriteMsgTxt)
	if err != nil {
		panic(err)
	}
}

// template for WriteMsg(w)(n, err)
const WriteMsgTxt = `
func (z *{{.StructName}}) WriteMsg(w io.Writer) (n int, err error) {
	var nn int

	nn, err = enc.WriteMapHeader(w, {{len .Fields}})
	n += nn
	if err != nil {
		return
	}

	{{range .Fields}}
	nn, err = enc.WriteString(w, "{{.FieldTag}}")
	n += nn
	if err != nil {
		return
	}

	{{if .IsPtr}}
	if {{.FieldName}} == nil {
		nn, err = enc.WriteNil(w)
	} else {
		{{if .IsSlice}}
		nn, err = enc.WriteArrayHeader(w, len({{.FieldName}}))
		{{else}}
		nn, err = enc.Write{{.KindName}}(w, {{.FieldName}})
		{{end}}
	}
	n += nn
	if err != nil {
		return
	}

	{{if .IsSlice}}{{/*slice of pointers*/}}
	for _, v := range {{.FieldName}}[:] {
		{{if .IsStruct}}
		nn, err = v.MarshalTo(w)
		{{else}}
		nn, err = enc.Write{{.KindName}}(w, *v)
		{{end}}
		n += nn
		if err != nil {
			return
		}
	}
	{{end}}
	{{else if .IsSlice}}{{/*slice to non-pointer type*/}}
	nn, err = enc.WriteArrayHeader(w, uint32(len({{.FieldName}})))
	n += nn
	if err != nil {
		return
	}
	
	for _, v := range {{.FieldName}}[:] {
		{{if .IsStruct}}
		nn, err = v.MarshalTo(w)
		{{else}}
		nn, err = enc.Write{{.KindName}}(w, v)
		{{end}}
		n += nn
		if err != nil {
			return
		}
	}
	
	{{else}}
	nn, err = enc.Write{{.KindName}}(w, {{.FieldName}})
	n += nn
	if err != nil {
		return
	}
	{{end}}
	{{end}}
	return
}
`
