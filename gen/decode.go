package gen

import (
	"text/template"
)

var (
	DecodeMsgTmpl *template.Template
)

const DecodeMsgTxt = `
func (z *{{.StructName}}) DecodeMsg(r io.Reader) (n int, err error) {
	var nn int
	var sz uint32 

	sz, nn, err = enc.ReadMapHeader(r)
	n += nn
	if err != nil {
		return
	}

	var scratch [64]byte
	var field []byte
	for field, nn, err = enc.ReadStringAsBytes(r, scratch[:]); err == nil; {
		n += nn
		switch enc.UnsafeString(field) {
		{{range .Fields}}
		case {{.FieldTag}}:
			// DECODE - TODO
		{{end}}
		default:
			// unrecognized field - do nothing
		}
	}
	n += nn
	if err == io.EOF {
		err = nil
	}
}
`
