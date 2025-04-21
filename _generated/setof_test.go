package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestSetOf_MarshalMsg(t *testing.T) {
	b, err := setofFilled.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var got SetOf
	if _, err := got.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, setofFilled) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, setofFilled)
	}
	// nil values...
	var empty SetOf
	b, err = empty.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := got.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, empty) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, empty)
	}
	// empty values...
	b, err = setofZeroLen.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := got.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, setofZeroLen) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, setofZeroLen)
	}
}

func TestSetOf_EncodeMsg(t *testing.T) {
	var buf bytes.Buffer
	enc := msgp.NewWriter(&buf)
	dec := msgp.NewReader(&buf)
	err := setofFilled.EncodeMsg(enc)
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	var got SetOf
	if err := got.DecodeMsg(dec); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, setofFilled) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, setofFilled)
	}
	// nil values...
	var empty SetOf
	err = empty.EncodeMsg(enc)
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	if err := got.DecodeMsg(dec); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, empty) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, empty)
	}
	// empty values...
	err = setofZeroLen.EncodeMsg(enc)
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	if err := got.DecodeMsg(dec); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, setofZeroLen) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, setofZeroLen)
	}
}

func TestSetOfTuple_MarshalMsg(t *testing.T) {
	want := SetOfTuple(setofFilled)
	b, err := want.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var got SetOfTuple
	if _, err := got.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, want)
	}
	// nil values...
	var empty SetOfTuple
	b, err = empty.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := got.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, empty) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, empty)
	}
	// empty values...
	want = SetOfTuple(setofZeroLen)
	b, err = want.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := got.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, want)
	}
}

func TestSetOfTuple_EncodeMsg(t *testing.T) {
	var buf bytes.Buffer
	enc := msgp.NewWriter(&buf)
	dec := msgp.NewReader(&buf)
	want := SetOfTuple(setofFilled)
	err := want.EncodeMsg(enc)
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	var got SetOfTuple
	if err := got.DecodeMsg(dec); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, want)
	}
	// nil values...
	var empty SetOfTuple
	err = empty.EncodeMsg(enc)
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	if err := got.DecodeMsg(dec); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, empty) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, empty)
	}
	// empty values...
	want = SetOfTuple(setofZeroLen)
	err = want.EncodeMsg(enc)
	if err != nil {
		t.Fatal(err)
	}
	enc.Flush()
	if err := got.DecodeMsg(dec); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", got, want)
	}
}

func TestSetOfSorted_MarshalMsg(t *testing.T) {
	a, err := setofSortedFilled.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	b, err := setofSortedFilled.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", a, b)
	}
}

func TestSetOfSorted_EncodeMsg(t *testing.T) {
	var a, b bytes.Buffer
	encA := msgp.NewWriter(&a)
	encB := msgp.NewWriter(&b)

	err := setofSortedFilled.EncodeMsg(encA)
	if err != nil {
		t.Fatal(err)
	}
	encA.Flush()
	err = setofSortedFilled.EncodeMsg(encB)
	if err != nil {
		t.Fatal(err)
	}
	encB.Flush()

	if !bytes.Equal(a.Bytes(), b.Bytes()) {
		t.Fatalf("unmarshaled value does not match original: %+v != %+v", a, b)
	}
}
