package _generated

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestRuneEncodeDecode(t *testing.T) {
	tt := &TestType{}
	r := 'r'
	rp := &r
	tt.Rune = r
	tt.RunePtr = &r
	tt.RunePtrPtr = &rp
	tt.RuneSlice = []rune{'a', 'b', '😳'}

	var buf bytes.Buffer
	wrt := msgp.NewWriter(&buf)
	if err := tt.EncodeMsg(wrt); err != nil {
		t.Errorf("%v", err)
	}
	wrt.Flush()

	var out TestType
	rdr := msgp.NewReader(&buf)
	if err := (&out).DecodeMsg(rdr); err != nil {
		t.Errorf("%v", err)
	}
	if r != out.Rune {
		t.Errorf("rune mismatch: expected %c found %c", r, out.Rune)
	}
	if r != *out.RunePtr {
		t.Errorf("rune ptr mismatch: expected %c found %c", r, *out.RunePtr)
	}
	if r != **out.RunePtrPtr {
		t.Errorf("rune ptr ptr mismatch: expected %c found %c", r, **out.RunePtrPtr)
	}
	if !reflect.DeepEqual(tt.RuneSlice, out.RuneSlice) {
		t.Errorf("rune slice mismatch")
	}
}

func TestRuneMarshalUnmarshal(t *testing.T) {
	tt := &TestType{}
	r := 'r'
	rp := &r
	tt.Rune = r
	tt.RunePtr = &r
	tt.RunePtrPtr = &rp
	tt.RuneSlice = []rune{'a', 'b', '😳'}

	bts, err := tt.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}

	var out TestType
	if _, err := (&out).UnmarshalMsg(bts); err != nil {
		t.Errorf("%v", err)
	}
	if r != out.Rune {
		t.Errorf("rune mismatch: expected %c found %c", r, out.Rune)
	}
	if r != *out.RunePtr {
		t.Errorf("rune ptr mismatch: expected %c found %c", r, *out.RunePtr)
	}
	if r != **out.RunePtrPtr {
		t.Errorf("rune ptr ptr mismatch: expected %c found %c", r, **out.RunePtrPtr)
	}
	if !reflect.DeepEqual(tt.RuneSlice, out.RuneSlice) {
		t.Errorf("rune slice mismatch")
	}
}

func TestJSONNumber(t *testing.T) {
	test := NumberJSONSample{
		Single: "-42",
		Array:  []json.Number{"0", "-0", "1", "-1", "0.1", "-0.1", "1234", "-1234", "12.34", "-12.34", "12E0", "12E1", "12e34", "12E-0", "12e+1", "12e-34", "-12E0", "-12E1", "-12e34", "-12E-0", "-12e+1", "-12e-34", "1.2E0", "1.2E1", "1.2e34", "1.2E-0", "1.2e+1", "1.2e-34", "-1.2E0", "-1.2E1", "-1.2e34", "-1.2E-0", "-1.2e+1", "-1.2e-34", "0E0", "0E1", "0e34", "0E-0", "0e+1", "0e-34", "-0E0", "-0E1", "-0e34", "-0E-0", "-0e+1", "-0e-34"},
		Map: map[string]json.Number{
			"a": json.Number("50.2"),
		},
	}

	// This is not guaranteed to be symmetric
	encoded, err := test.MarshalMsg(nil)
	if err != nil {
		t.Errorf("%v", err)
	}
	var v NumberJSONSample
	_, err = v.UnmarshalMsg(encoded)
	if err != nil {
		t.Errorf("%v", err)
	}
	// Test two values
	if v.Single != "-42" {
		t.Errorf("want %v, got %v", "-42", v.Single)
	}
	if v.Map["a"] != "50.2" {
		t.Errorf("want %v, got %v", "50.2", v.Map["a"])
	}

	var jsBuf bytes.Buffer
	remain, err := msgp.UnmarshalAsJSON(&jsBuf, encoded)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(remain) != 0 {
		t.Errorf("remain should be empty")
	}
	wantjs := `{"Single":-42,"Array":[0,0,1,-1,0.1,-0.1,1234,-1234,12.34,-12.34,12,120,120000000000000000000000000000000000,12,120,0.0000000000000000000000000000000012,-12,-120,-120000000000000000000000000000000000,-12,-120,-0.0000000000000000000000000000000012,1.2,12,12000000000000000000000000000000000,1.2,12,0.00000000000000000000000000000000012,-1.2,-12,-12000000000000000000000000000000000,-1.2,-12,-0.00000000000000000000000000000000012,0,0,0,0,0,0,-0,-0,-0,-0,-0,-0],"Map":{"a":50.2}}`
	if jsBuf.String() != wantjs {
		t.Errorf("jsBuf.String() = \n%s, want \n%s", jsBuf.String(), wantjs)
	}
	// Test encoding
	var buf bytes.Buffer
	en := msgp.NewWriter(&buf)
	err = test.EncodeMsg(en)
	if err != nil {
		t.Errorf("%v", err)
	}
	en.Flush()
	encoded = buf.Bytes()

	dc := msgp.NewReader(&buf)
	err = v.DecodeMsg(dc)
	if err != nil {
		t.Errorf("%v", err)
	}
	// Test two values
	if v.Single != "-42" {
		t.Errorf("want %s, got %s", "-42", v.Single)
	}
	if v.Map["a"] != "50.2" {
		t.Errorf("want %s, got %s", "50.2", v.Map["a"])
	}

	jsBuf.Reset()
	remain, err = msgp.UnmarshalAsJSON(&jsBuf, encoded)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(remain) != 0 {
		t.Errorf("remain should be empty")
	}
	if jsBuf.String() != wantjs {
		t.Errorf("jsBuf.String() = \n%s, want \n%s", jsBuf.String(), wantjs)
	}

	// Try interface encoder
	jd := json.NewDecoder(&jsBuf)
	jd.UseNumber()
	var jsIntf map[string]any
	err = jd.Decode(&jsIntf)
	if err != nil {
		t.Errorf("%v", err)
	}
	// Ensure we encode correctly
	_ = (jsIntf["Single"]).(json.Number)

	fromInt, err := msgp.AppendIntf(nil, jsIntf)
	if err != nil {
		t.Errorf("%v", err)
	}

	// Take the value from the JSON interface encoder and unmarshal back into our struct.
	v = NumberJSONSample{}
	_, err = v.UnmarshalMsg(fromInt)
	if err != nil {
		t.Errorf("%v", err)
	}
	if v.Single != "-42" {
		t.Errorf("want %s, got %s", "-42", v.Single)
	}
	if v.Map["a"] != "50.2" {
		t.Errorf("want %s, got %s", "50.2", v.Map["a"])
	}
}
