package msgp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestLocate(t *testing.T) {
	var buf bytes.Buffer
	en := NewWriter(&buf)
	en.WriteMapHeader(2)
	en.WriteString("thing_one")
	en.WriteString("value_one")
	en.WriteString("thing_two")
	en.WriteFloat64(2.0)

	field := Locate("thing_one", buf.Bytes())
	if len(field) == 0 {
		t.Fatal("field not found")
	}

	var zbuf bytes.Buffer
	NewWriter(&zbuf).WriteString("value_one")

	if !bytes.Equal(zbuf.Bytes(), field) {
		t.Errorf("got %q; wanted %q", field, zbuf.Bytes())
	}

	zbuf.Reset()
	NewWriter(&zbuf).WriteFloat64(2.0)
	field = Locate("thing_two", buf.Bytes())
	if len(field) == 0 {
		t.Fatal("field not found")
	}
	if !bytes.Equal(zbuf.Bytes(), field) {
		t.Errorf("got %q; wanted %q", field, zbuf.Bytes())
	}

	field = Locate("nope", buf.Bytes())
	if len(field) != 0 {
		t.Fatalf("wanted %q; got %q", ErrFieldNotFound, nil)
	}

}

func TestReplace(t *testing.T) {
	// there are 4 cases that need coverage:
	//  - new value is smaller than old value
	//  - new value is the same size as the old value
	//  - new value is larger than old, but fits within cap(b)
	//  - new value is larger than old, and doesn't fit within cap(b)

	var buf bytes.Buffer
	en := NewWriter(&buf)
	en.WriteMapHeader(3)
	en.WriteString("thing_one")
	en.WriteString("value_one")
	en.WriteString("thing_two")
	en.WriteFloat64(2.0)
	en.WriteString("some_bytes")
	en.WriteBytes([]byte("here are some bytes"))

	// same-size replacement
	var fbuf bytes.Buffer
	NewWriter(&fbuf).WriteFloat64(4.0)

	// replace 2.0 with 4.0 in field two
	raw, err := Replace("thing_two", buf.Bytes(), fbuf.Bytes())
	if err != nil {
		t.Fatal(err)
	}

	m := make(map[string]interface{})
	m, _, err = ReadMapStrIntfBytes(raw, m)
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m["thing_two"], 4.0) {
		t.Errorf("wanted %v; got %v", 4.0, m["thing_two"])
	}

	// smaller-size replacement
	// replace 2.0 with []byte("hi!")
	fbuf.Reset()
	NewWriter(&fbuf).WriteBytes([]byte("hi!"))
	raw, err = Replace("thing_two", raw, fbuf.Bytes())
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	m, _, err = ReadMapStrIntfBytes(raw, m)
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m["thing_two"], []byte("hi!")) {
		t.Errorf("wanted %v; got %v", []byte("hi!"), m["thing_two"])
	}

	// larger-size replacement
	fbuf.Reset()
	NewWriter(&fbuf).WriteBytes([]byte("some even larger bytes than before"))
	raw, err = Replace("some_bytes", raw, fbuf.Bytes())
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	m, _, err = ReadMapStrIntfBytes(raw, m)
	if err != nil {
		t.Logf("%q", raw)
		t.Fatal(err)
	}

	if !reflect.DeepEqual(m["some_bytes"], []byte("some even larger bytes than before")) {
		t.Errorf("wanted %v; got %v", []byte("hello there!"), m["some_bytes"])
	}
}

func BenchmarkLocate(b *testing.B) {
	var buf bytes.Buffer
	en := NewWriter(&buf)
	en.WriteMapHeader(3)
	en.WriteString("thing_one")
	en.WriteString("value_one")
	en.WriteString("thing_two")
	en.WriteFloat64(2.0)
	en.WriteString("thing_three")
	en.WriteBytes([]byte("hello!"))

	raw := buf.Bytes()
	// bytes/s will be the number of bytes traversed per unit of time
	field := Locate("thing_three", raw)
	b.SetBytes(int64(len(raw) - len(field)))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Locate("thing_three", raw)
	}
}
