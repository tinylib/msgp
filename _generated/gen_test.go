package _generated

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/tinylib/msgp/msgp"
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
	msgp.Encode(&buf, v)
	en := msgp.NewWriter(msgp.Nowhere)
	b.SetBytes(int64(buf.Len()))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.EncodeMsg(en)
	}
	en.Flush()
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
	msgp.Encode(&buf, v)
	dc := msgp.NewReader(msgp.NewEndlessReader(buf.Bytes(), b))
	b.SetBytes(int64(buf.Len()))
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.DecodeMsg(dc)
	}
}

func (a *TestType) Equal(b *TestType) bool {
	// compare times, appended, then zero out those
	// fields, perform a DeepEqual, and restore them
	ta, tb := a.Time, b.Time
	if !ta.Equal(tb) {
		return false
	}
	aa, ab := a.Appended, b.Appended
	if !bytes.Equal(aa, ab) {
		return false
	}
	if len(a.MapStringEmpty) == 0 && len(b.MapStringEmpty) == 0 {
		a.MapStringEmpty = nil
		b.MapStringEmpty = nil
	}
	if len(a.MapStringEmpty2) == 0 && len(b.MapStringEmpty2) == 0 {
		a.MapStringEmpty2 = nil
		b.MapStringEmpty2 = nil
	}

	a.Time, b.Time = time.Time{}, time.Time{}
	aa, ab = nil, nil
	ok := reflect.DeepEqual(a, b)
	a.Time, b.Time = ta, tb
	a.Appended, b.Appended = aa, ab
	return ok
}

// This covers the following cases:
//   - Recursive types
//   - Non-builtin identifiers (and recursive types)
//   - time.Time
//   - map[string]string
//   - anonymous structs
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
		Child:           nil,
		Time:            time.Now(),
		Appended:        msgp.Raw([]byte{}), // 'nil'
		MapStringEmpty:  map[string]struct{}{"Key": {}, "Key2": {}},
		MapStringEmpty2: map[string]EmptyStruct{"Key3": {}, "Key4": {}},
	}

	var buf bytes.Buffer

	err := msgp.Encode(&buf, tt)
	if err != nil {
		t.Fatal(err)
	}

	tnew := new(TestType)

	err = msgp.Decode(&buf, tnew)
	if err != nil {
		t.Error(err)
	}

	if !tt.Equal(tnew) {
		t.Logf("in:  %#v", tt)
		t.Logf("out: %#v", tnew)
		t.Fatal("objects not equal")
	}
	if tnew.Time.Location() != time.UTC {
		t.Errorf("time location not UTC: %v", tnew.Time.Location())
	}

	tanother := new(TestType)

	buf.Reset()
	msgp.Encode(&buf, tt)

	var left []byte
	left, err = tanother.UnmarshalMsg(buf.Bytes())
	if err != nil {
		t.Error(err)
	}
	if len(left) > 0 {
		t.Errorf("%d bytes left", len(left))
	}

	if !tt.Equal(tanother) {
		t.Logf("in: %v", tt)
		t.Logf("out: %v", tanother)
		t.Fatal("objects not equal")
	}
	if tanother.Time.Location() != time.UTC {
		t.Errorf("time location not UTC: %v", tanother.Time.Location())
	}

}

func TestIssue168(t *testing.T) {
	buf := bytes.Buffer{}
	test := TestObj{}

	msgp.Encode(&buf, &TestObj{ID1: "1", ID2: "2"})
	msgp.Decode(&buf, &test)

	if test.ID1 != "1" || test.ID2 != "2" {
		t.Fatalf("got back %+v", test)
	}
}

func TestIssue362(t *testing.T) {
	in := StructByteSlice{
		ABytes:      make([]byte, 0),
		AString:     make([]string, 0),
		ABool:       make([]bool, 0),
		AInt:        make([]int, 0),
		AInt8:       make([]int8, 0),
		AInt16:      make([]int16, 0),
		AInt32:      make([]int32, 0),
		AInt64:      make([]int64, 0),
		AUint:       make([]uint, 0),
		AUint8:      make([]uint8, 0),
		AUint16:     make([]uint16, 0),
		AUint32:     make([]uint32, 0),
		AUint64:     make([]uint64, 0),
		AFloat32:    make([]float32, 0),
		AFloat64:    make([]float64, 0),
		AComplex64:  make([]complex64, 0),
		AComplex128: make([]complex128, 0),
		AStruct:     make([]Fixed, 0),
	}

	b, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	var dst StructByteSlice
	_, err = dst.UnmarshalMsg(b)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, dst) {
		t.Fatalf("mismatch %#v != %#v", in, dst)
	}
	dst2 := StructByteSlice{}
	dec := msgp.NewReader(bytes.NewReader(b))
	err = dst2.DecodeMsg(dec)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, dst2) {
		t.Fatalf("mismatch %#v != %#v", in, dst2)
	}

	// Encode with nil
	zero := StructByteSlice{}
	b, err = zero.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	// Decode into dst that now has values...
	_, err = dst.UnmarshalMsg(b)
	if err != nil {
		t.Fatal(err)
	}
	// All should be nil now.
	if !reflect.DeepEqual(zero, dst) {
		t.Fatalf("mismatch %#v != %#v", zero, dst)
	}
	dec = msgp.NewReader(bytes.NewReader(b))
	err = dst2.DecodeMsg(dec)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(zero, dst2) {
		t.Fatalf("mismatch %#v != %#v", zero, dst2)
	}
}
