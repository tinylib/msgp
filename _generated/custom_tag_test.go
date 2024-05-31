package _generated

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"bytes"

	"github.com/tinylib/msgp/msgp"
)

func TestCustomTag(t *testing.T) {
	t.Run("File Scope", func(t *testing.T) {
		ts := CustomTag{
			Foo: "foostring13579",
			Bar: 999_999}
		encDecCustomTag(t, ts, "mytag")
	})
}

func encDecCustomTag(t *testing.T, testStruct msgp.Encodable, tag string) {
	var b bytes.Buffer
	msgp.Encode(&b, testStruct)

	// Check tag names using JSON as an intermediary layer
	// TODO: is there a way to avoid the JSON layer? We'd need to directly decode raw msgpack -> map[string]any
	refJSON, err := json.Marshal(testStruct)
	if err != nil {
		t.Error(fmt.Sprintf("error encoding struct as JSON: %v", err))
	}
	ref := make(map[string]any)
	// Encoding and decoding the original struct via JSON is necessary
	// for field comparisons to work, since JSON -> map[string]any
	// relies on type inferences such as all numbers being float64s
	json.Unmarshal(refJSON, &ref)

	var encJSON bytes.Buffer
	msgp.UnmarshalAsJSON(&encJSON, b.Bytes())
	encoded := make(map[string]any)
	json.Unmarshal(encJSON.Bytes(), &encoded)

	tsType := reflect.TypeOf(testStruct)
	for i := 0; i < tsType.NumField(); i++ {
		// Check encoded field name
		field := tsType.Field(i)
		encodedValue, ok := encoded[field.Tag.Get(tag)]
		if !ok {
			t.Error("missing encoded value for field", field.Name)
			continue
		}
		// Check encoded field value (against original value post-JSON enc + dec)
		jsonName, ok := field.Tag.Lookup("json")
		if !ok {
			jsonName = field.Name
		}
		refValue := ref[jsonName]
		if !reflect.DeepEqual(refValue, encodedValue) {
			t.Error(fmt.Sprintf("incorrect encoded value for field %s. reference: %v, encoded: %v",
				field.Name, refValue, encodedValue))
		}
	}
}
