package gen

import (
	"bufio"
	"os"
	"testing"
	"text/template"
)

func TestEncTemplate(t *testing.T) {
	tmp, err := template.ParseFiles("encode.tmpl", "elem_enc.tmpl")
	if err != nil {
		t.Fatal(err)
	}

	// type Marshaler struct {
	// 	thing1 *float64
	// 	body   []byte
	//	things []complex128
	// }
	var val *Ptr = &Ptr{
		Varname: "z",
		Value: &Struct{
			Name: "Marshaler",
			Fields: []StructField{
				{
					FieldTag: "thing1",
					FieldElem: &Ptr{
						Varname: "z.thing1",
						Value: &BaseElem{
							Varname: "*z.thing1",
							Value:   Float64,
						},
					},
				},
				{
					FieldTag: "body",
					FieldElem: &BaseElem{
						Varname: "z.body",
						Value:   Bytes,
					},
				},
				{
					FieldTag: "things",
					FieldElem: &Slice{
						Varname: "z.things",
						Els: &BaseElem{
							Varname: "z.things[i]",
							Value:   Complex128,
						},
					},
				},
			},
		},
	}

	out, err := os.Create("_test-gen-enc.go")
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	wr := bufio.NewWriter(out)

	_, err = wr.WriteString(`package gen
import(
	"github.com/philhofer/msgp/enc"
	"io"
)

type Marshaler struct {
	thing1 *float64
	body []byte
	things []complex128
}
`)

	if err != nil {
		t.Fatal(err)
	}

	err = tmp.Execute(wr, val)
	if err != nil {
		t.Error(err)
	}

	err = wr.Flush()
	if err != nil {
		t.Error(err)
	}
}
