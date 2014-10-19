package enc

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.WriteMapHeader(5)

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

	var js bytes.Buffer
	_, err := UnmarshalAsJSON(&js, buf.Bytes())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", js.Bytes())
	mp := make(map[string]interface{})
	err = json.Unmarshal(js.Bytes(), &mp)
	if err != nil {
		t.Log(js.String())
		t.Fatalf("Error unmarshaling: %s", err)
	}

	if len(mp) != 5 {
		t.Errorf("map length should be %d, not %d", 4, len(mp))
	}

	so, ok := mp["thing_1"]
	if !ok || so != "a string object" {
		t.Errorf("expected %q; got %q", "a string object", so)
	}

	in, ok := mp["a map"]
	if !ok {
		t.Error("no key 'a map'")
	}
	if inm, ok := in.(map[string]interface{}); !ok {
		t.Error("inner map not type-assertable to map[string]interface{}")
	} else {
		inm1, ok := inm["internal_one"]
		if !ok || !reflect.DeepEqual(inm1, "blah") {
			t.Errorf("inner map field %q should be %q, not %q", "internal_one", "blah", inm1)
		}
	}
}

func BenchmarkUnmarshalAsJSON(b *testing.B) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
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
	enc.WriteString("part_A")
	enc.WriteString("part_B")

	enc.WriteString("a_second_map")
	enc.WriteMapStrStr(map[string]string{
		"internal_one": "blah",
		"internal_two": "blahhh...",
	})
	var js bytes.Buffer
	bts := buf.Bytes()
	_, err := UnmarshalAsJSON(&js, bts)
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(js.Bytes())))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		js.Reset()
		UnmarshalAsJSON(&js, bts)
	}
}
