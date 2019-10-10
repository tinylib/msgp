package _generated

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func mustEncodeToJSON(o msgp.Encodable) string {
	var buf bytes.Buffer
	var err error

	en := msgp.NewWriter(&buf)
	err = o.EncodeMsg(en)
	if err != nil {
		panic(err)
	}
	en.Flush()

	var outbuf bytes.Buffer
	_, err = msgp.CopyToJSON(&outbuf, &buf)
	if err != nil {
		panic(err)
	}

	return outbuf.String()
}

func TestOmitEmpty0(t *testing.T) {

	var s string

	var oe0a OmitEmpty0

	s = mustEncodeToJSON(&oe0a)
	if s != `{"aunnamedstruct":{},"aarrayint":[0,0,0,0,0]}` {
		t.Errorf("wrong result: %s", s)
	}

	var oe0b OmitEmpty0
	oe0b.AString = "teststr"
	s = mustEncodeToJSON(&oe0b)
	if s != `{"astring":"teststr","aunnamedstruct":{},"aarrayint":[0,0,0,0,0]}` {
		t.Errorf("wrong result: %s", s)
	}

	// more than 15 fields filled in
	var oe0c OmitEmpty0
	oe0c.ABool = true
	oe0c.AInt = 1
	oe0c.AInt8 = 1
	oe0c.AInt16 = 1
	oe0c.AInt32 = 1
	oe0c.AInt64 = 1
	oe0c.AUint = 1
	oe0c.AUint8 = 1
	oe0c.AUint16 = 1
	oe0c.AUint32 = 1
	oe0c.AUint64 = 1
	oe0c.AFloat32 = 1
	oe0c.AFloat64 = 1
	oe0c.AComplex64 = complex(1, 1)
	oe0c.AComplex128 = complex(1, 1)
	oe0c.AString = "test"
	oe0c.ANamedBool = true
	oe0c.ANamedInt = 1
	oe0c.ANamedFloat64 = 1

	var buf bytes.Buffer
	en := msgp.NewWriter(&buf)
	err := oe0c.EncodeMsg(en)
	if err != nil {
		t.Fatal(err)
	}
	en.Flush()
	de := msgp.NewReader(&buf)
	var oe0d OmitEmpty0
	err = oe0d.DecodeMsg(de)
	if err != nil {
		t.Fatal(err)
	}

	// spot check some fields
	if oe0c.AFloat32 != oe0d.AFloat32 {
		t.Fail()
	}
	if oe0c.ANamedBool != oe0d.ANamedBool {
		t.Fail()
	}
	if oe0c.AInt64 != oe0d.AInt64 {
		t.Fail()
	}

}

// TestOmitEmptyHalfFull tests mixed omitempty and not
func TestOmitEmptyHalfFull(t *testing.T) {

	var s string

	var oeA OmitEmptyHalfFull

	s = mustEncodeToJSON(&oeA)
	if s != `{"field01":"","field03":""}` {
		t.Errorf("wrong result: %s", s)
	}

	var oeB OmitEmptyHalfFull
	oeB.Field02 = "val2"
	s = mustEncodeToJSON(&oeB)
	if s != `{"field01":"","field02":"val2","field03":""}` {
		t.Errorf("wrong result: %s", s)
	}

	var oeC OmitEmptyHalfFull
	oeC.Field03 = "val3"
	s = mustEncodeToJSON(&oeC)
	if s != `{"field01":"","field03":"val3"}` {
		t.Errorf("wrong result: %s", s)
	}
}

// TestOmitEmptyLotsOFields tests the case of > 64 fields (triggers the bitmask needing to be an array instead of a single value)
func TestOmitEmptyLotsOFields(t *testing.T) {

	var s string

	var oeLotsA OmitEmptyLotsOFields

	s = mustEncodeToJSON(&oeLotsA)
	if s != `{}` {
		t.Errorf("wrong result: %s", s)
	}

	var oeLotsB OmitEmptyLotsOFields
	oeLotsB.Field04 = "val4"
	s = mustEncodeToJSON(&oeLotsB)
	if s != `{"field04":"val4"}` {
		t.Errorf("wrong result: %s", s)
	}

	var oeLotsC OmitEmptyLotsOFields
	oeLotsC.Field64 = "val64"
	s = mustEncodeToJSON(&oeLotsC)
	if s != `{"field64":"val64"}` {
		t.Errorf("wrong result: %s", s)
	}

}

func BenchmarkOmitEmpty10AllEmpty(b *testing.B) {

	en := msgp.NewWriter(ioutil.Discard)
	var s OmitEmpty10

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := s.EncodeMsg(en)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkOmitEmpty10AllFull(b *testing.B) {

	en := msgp.NewWriter(ioutil.Discard)
	var s OmitEmpty10
	s.Field00 = "this is the value of field00"
	s.Field01 = "this is the value of field01"
	s.Field02 = "this is the value of field02"
	s.Field03 = "this is the value of field03"
	s.Field04 = "this is the value of field04"
	s.Field05 = "this is the value of field05"
	s.Field06 = "this is the value of field06"
	s.Field07 = "this is the value of field07"
	s.Field08 = "this is the value of field08"
	s.Field09 = "this is the value of field09"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := s.EncodeMsg(en)
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkNotOmitEmpty10AllEmpty(b *testing.B) {

	en := msgp.NewWriter(ioutil.Discard)
	var s NotOmitEmpty10

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := s.EncodeMsg(en)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNotOmitEmpty10AllFull(b *testing.B) {

	en := msgp.NewWriter(ioutil.Discard)
	var s NotOmitEmpty10
	s.Field00 = "this is the value of field00"
	s.Field01 = "this is the value of field01"
	s.Field02 = "this is the value of field02"
	s.Field03 = "this is the value of field03"
	s.Field04 = "this is the value of field04"
	s.Field05 = "this is the value of field05"
	s.Field06 = "this is the value of field06"
	s.Field07 = "this is the value of field07"
	s.Field08 = "this is the value of field08"
	s.Field09 = "this is the value of field09"

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := s.EncodeMsg(en)
		if err != nil {
			b.Fatal(err)
		}
	}
}
