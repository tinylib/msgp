package gen

import (
//"reflect"
//"testing"
)

type TestType struct {
	Thing  *float64          `msg:"blah"`
	Things []*string         `msg:"things"`
	Imap   map[string]string `msg:"imap"`
	Simple int               `msg:"simple"`
}

/*
func TestGenerateTree(t *testing.T) {
	var want *Ptr = &Ptr{
		Varname: "z",
		Value: &Struct{
			Name: "TestType",
			Fields: []StructField{
				{
					FieldName: "Thing",
					FieldTag:  "blah",
					FieldElem: &Ptr{
						Varname: "z.Thing",
						Value: &BaseElem{
							Varname: "*z.Thing",
							Value:   Float64,
						},
					},
				},
				{
					FieldName: "Things",
					FieldTag:  "things",
					FieldElem: &Slice{
						Varname: "z.Things",
						Els: &Ptr{
							Varname: "z.Things[i]",
							Value: &BaseElem{
								Varname: "*z.Things[i]",
								Value:   String,
							},
						},
					},
				},
				{
					FieldName: "Imap",
					FieldTag:  "imap",
					FieldElem: &BaseElem{
						Varname: "z.Imap",
						Value:   MapStrStr,
					},
				},
				{
					FieldName: "Simple",
					FieldTag:  "simple",
					FieldElem: &BaseElem{
						Varname: "z.Simple",
						Value:   Int,
					},
				},
			},
		},
	}

	var tt *TestType
	e, err := GenerateTree(reflect.TypeOf(tt))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(e, want) {
		t.Error("Decoded elements not equal.")
		t.Errorf("Input: %s", want)
		t.Errorf("Output: %s", e)
	}

}*/
