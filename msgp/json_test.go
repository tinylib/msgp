package msgp

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
	"unicode/utf8"
)

func TestCopyJSON(t *testing.T) {
	var buf bytes.Buffer
	enc := NewWriter(&buf)
	const mapLength = 6
	enc.WriteMapHeader(mapLength)

	enc.WriteString("thing_1")
	enc.WriteString("a string object")

	enc.WriteString("a_map")
	enc.WriteMapHeader(2)
	enc.WriteString("float_a")
	enc.WriteFloat32(1.0)
	enc.WriteString("int_b")
	enc.WriteInt64(-100)

	enc.WriteString("some bytes")
	enc.WriteBytes([]byte("here are some bytes"))
	enc.WriteString("a bool")
	enc.WriteBool(true)

	enc.WriteString("a map")
	enc.WriteMapStrStr(map[string]string{
		"internal_one": "blah",
		"internal_two": "blahhh...",
	})
	enc.WriteString("float64")
	const encodedFloat64 = 1672209023
	enc.WriteFloat64(encodedFloat64)
	enc.Flush()

	var js bytes.Buffer
	_, err := CopyToJSON(&js, &buf)
	if err != nil {
		t.Fatal(err)
	}
	mp := make(map[string]any)
	err = json.Unmarshal(js.Bytes(), &mp)
	if err != nil {
		t.Log(js.String())
		t.Fatalf("Error unmarshaling: %s", err)
	}

	if len(mp) != mapLength {
		t.Errorf("map length should be %d, not %d", mapLength, len(mp))
	}

	so, ok := mp["thing_1"]
	if !ok || so != "a string object" {
		t.Errorf("expected %q; got %q", "a string object", so)
	}

	in, ok := mp["a map"]
	if !ok {
		t.Error("no key 'a map'")
	}
	if inm, ok := in.(map[string]any); !ok {
		t.Error("inner map not type-assertable to map[string]interface{}")
	} else {
		inm1, ok := inm["internal_one"]
		if !ok || !reflect.DeepEqual(inm1, "blah") {
			t.Errorf("inner map field %q should be %q, not %q", "internal_one", "blah", inm1)
		}
	}
	if actual := mp["float64"]; float64(encodedFloat64) != actual.(float64) {
		t.Errorf("expected %G, got %G", float64(encodedFloat64), actual)
	}
}

// Encoder should generate valid utf-8 even if passed bad input
func TestCopyJSONNegativeUTF8(t *testing.T) {
	// Single string with non-compliant utf-8 byte
	stringWithBadUTF8 := []byte{
		0xa1, 0xe0,
	}

	src := bytes.NewBuffer(stringWithBadUTF8)

	var js bytes.Buffer
	_, err := CopyToJSON(&js, src)
	if err != nil {
		t.Fatal(err)
	}

	// Even though we provided bad input, should have escaped the naughty character
	if !utf8.Valid(js.Bytes()) {
		t.Errorf("Expected JSON to be valid utf-8 even when provided bad input")
	}

	// Expect a bad character string
	expected := `"\ufffd"`
	if js.String() != expected {
		t.Errorf("Expected: '%s', got: '%s'", expected, js.String())
	}
}

func BenchmarkCopyToJSON(b *testing.B) {
	var buf bytes.Buffer
	enc := NewWriter(&buf)
	enc.WriteMapHeader(4)

	enc.WriteString("thing_1")
	enc.WriteString("a string object")

	enc.WriteString("a_first_map")
	enc.WriteMapHeader(2)
	enc.WriteString("float_a")
	enc.WriteFloat32(1.0)
	enc.WriteString("int_b")
	enc.WriteInt64(-100)

	enc.WriteString("an array")
	enc.WriteArrayHeader(2)
	enc.WriteBool(true)
	enc.WriteUint(2089)

	enc.WriteString("a_second_map")
	enc.WriteMapStrStr(map[string]string{
		"internal_one": "blah",
		"internal_two": "blahhh...",
	})
	enc.Flush()

	var js bytes.Buffer
	bts := buf.Bytes()
	_, err := CopyToJSON(&js, &buf)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(js.Bytes())))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		js.Reset()
		CopyToJSON(&js, bytes.NewReader(bts))
	}
}

func BenchmarkStdlibJSON(b *testing.B) {
	obj := map[string]any{
		"thing_1": "a string object",
		"a_first_map": map[string]any{
			"float_a": float32(1.0),
			"float_b": -100,
		},
		"an array": []any{
			"part_A",
			"part_B",
		},
		"a_second_map": map[string]any{
			"internal_one": "blah",
			"internal_two": "blahhh...",
		},
	}
	var js bytes.Buffer
	err := json.NewEncoder(&js).Encode(&obj)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(js.Bytes())))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		js.Reset()
		json.NewEncoder(&js).Encode(&obj)
	}
}
