package _generated

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

// benchmark encoding a small, "fast" type.
// the point here is to see how much garbage
// is generated intrinsically by the encoding/
// decoding process as opposed to the nature
// of the struct.
func BenchmarkFastEncode(b *testing.B) {
	v := &TestFast{
		Lat:  40.12398,
		Long: -41.9082,
		Alt:  201.08290,
		Data: []byte("whaaaaargharbl"),
	}
	var buf bytes.Buffer
	n, _ := v.WriteTo(&buf)
	b.SetBytes(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		v.WriteTo(&buf)
	}

	// should return ~500ns and 1 alloc
	// anything else means a regression
}

// benchmark decoding a small, "fast" type.
// the point here is to see how much garbage
// is generated intrinsically by the encoding/
// decoding process as opposed to the nature
// of the struct.
func BenchmarkFastDecode(b *testing.B) {
	v := &TestFast{
		Lat:  40.12398,
		Long: -41.9082,
		Alt:  201.08290,
		Data: []byte("whaaaaargharbl"),
	}

	var buf bytes.Buffer
	n, _ := v.WriteTo(&buf)
	rd := bytes.NewReader(buf.Bytes())
	b.SetBytes(n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rd.Seek(0, 0) // reset
		v.ReadFrom(rd)
	}

	// should return ~700ns and 2 allocs
	// anything else means a regression
}

// This covers the following cases:
//  - Recursive types
//  - Non-builtin identifiers (and recursive types)
//  - time.Time
//  - map[string]string
//  - anonymous structs
//
func Test1EncodeDecode(t *testing.T) {
	f := 32.00
	tt := &TestType{
		F: &f,
		Els: map[string]string{
			"thing_one": "one",
			"thing_two": "two",
		},
		Obj: struct {
			ValueA string `msg:"value_a"`
			ValueB []byte `msg:"value_b"`
		}{
			ValueA: "here's the first inner value",
			ValueB: []byte("here's the second inner value"),
		},
		Child: nil,
		Time:  time.Now(),
	}

	var buf bytes.Buffer

	_, err := tt.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}

	tnew := new(TestType)

	_, err = tnew.ReadFrom(&buf)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(tt, tnew) {
		t.Logf("in: %v", tt)
		t.Logf("out: %v", tnew)
		t.Fatal("objects not equal")
	}
}
