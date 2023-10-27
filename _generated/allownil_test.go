package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestAllownil(t *testing.T) {
	tt := &NamedStructAN{
		A: []string{},
		B: nil,
	}
	var buf bytes.Buffer

	err := msgp.Encode(&buf, tt)
	if err != nil {
		t.Fatal(err)
	}
	in := buf.Bytes()

	for _, tnew := range []*NamedStructAN{{}, {A: []string{}}, {B: []string{}}} {
		err = msgp.Decode(bytes.NewBuffer(in), tnew)
		if err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(tt, tnew) {
			t.Logf("in:  %#v", tt)
			t.Logf("out: %#v", tnew)
			t.Fatal("objects not equal")
		}
	}

	in, err = tt.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, tanother := range []*NamedStructAN{{}, {A: []string{}}, {B: []string{}}} {
		var left []byte
		left, err = tanother.UnmarshalMsg(in)
		if err != nil {
			t.Error(err)
		}
		if len(left) > 0 {
			t.Errorf("%d bytes left", len(left))
		}

		if !reflect.DeepEqual(tt, tanother) {
			t.Logf("in:  %#v", tt)
			t.Logf("out: %#v", tanother)
			t.Fatal("objects not equal")
		}
	}
}
