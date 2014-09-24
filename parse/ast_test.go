package parse

import (
	"github.com/philhofer/msgp/gen"
	"reflect"
	"testing"
)

func TestAST(t *testing.T) {
	f, err := GetElems("./_to_parse.go")
	if err != nil {
		t.Error(err)
	}

	if len(f) != 1 {
		t.Fatal("Got %d elements; expected %d", 1, len(f))
	}

	if !reflect.DeepEqual(f[0], want) {
		t.Logf("Got: %s", f[0])
		t.Logf("Want: %s", want)
		t.Error("Values not equal.")
	}

}

var want gen.Elem = &gen.Ptr{
	Value: &gen.Struct{
		Name: "TestType",
		Fields: []gen.StructField{
			{
				FieldTag:  "map",
				FieldName: "M",
				FieldElem: &gen.Map{
					Value: &gen.BaseElem{
						Value: gen.String,
					},
				},
			},
			{
				FieldTag:  "float",
				FieldName: "V",
				FieldElem: &gen.Ptr{
					Value: &gen.BaseElem{
						Value: gen.Float64,
					},
				},
			},
			{
				FieldTag:  "thing",
				FieldName: "Thing",
				FieldElem: &gen.Struct{
					Fields: []gen.StructField{
						{
							FieldTag:  "a",
							FieldName: "A",
							FieldElem: &gen.BaseElem{
								Value: gen.Bytes,
							},
						},
						{
							FieldTag:  "b",
							FieldName: "B",
							FieldElem: &gen.BaseElem{
								Value: gen.Bytes,
							},
						},
					},
				},
			},
		},
	},
}
