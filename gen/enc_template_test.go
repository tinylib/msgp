package gen

import (
	"bufio"
	"os"
	"testing"
)

func TestEncTemplate(t *testing.T) {

	// type Marshaler struct {
	// 	Thing1 *float64
	// 	Body   []byte
	//	Things []complex128
	//  M      map[string]string
	// }
	var val *Ptr = &Ptr{
		Varname: "z",
		Value: &Struct{
			Name: "Marshaler",
			Fields: []StructField{
				{
					FieldTag: "thing1",
					FieldElem: &Ptr{
						Varname: "z.Thing1",
						Value: &BaseElem{
							Varname: "*z.Thing1",
							Value:   Float64,
						},
					},
				},
				{
					FieldTag: "body",
					FieldElem: &BaseElem{
						Varname: "z.Body",
						Value:   Bytes,
					},
				},
				{
					FieldTag: "things",
					FieldElem: &Slice{
						Varname: "z.Things",
						Els: &BaseElem{
							Varname: "z.Things[i]",
							Value:   Complex128,
						},
					},
				},
				{
					FieldTag: "a_map",
					FieldElem: &Map{
						Varname: "z.M",
						Value: &BaseElem{
							Varname: "val",
							Value:   String,
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

	err = WriteDecoderMethod(wr, val)
	if err != nil {
		t.Fatal(err)
	}
	err = wr.Flush()
	if err != nil {
		t.Error(err)
	}
}
