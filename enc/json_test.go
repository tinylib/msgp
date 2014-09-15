package enc

import (
	"bytes"
	"testing"
)

func TestCopyJSON(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.WriteMapHeader(4)

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

	enc.WriteString("a map")
	enc.WriteMapStrStr(map[string]string{
		"internal_one": "blah",
		"internal_two": "blahhh...",
	})

	var js bytes.Buffer
	_, err := CopyToJSON(&js, &buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(js.String())
}
