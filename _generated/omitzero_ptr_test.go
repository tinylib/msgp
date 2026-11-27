package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestOmitZeroNilPtr(t *testing.T) {
	for _, tc := range []struct {
		in     OmitZeroPtr
		fields uint32
	}{
		{OmitZeroPtr{X: 1}, 1},
		{OmitZeroPtr{Val: &OmitZeroVal{}, Ptr: &OmitZeroPtrRecv{}, X: 1}, 1},
		{OmitZeroPtr{Val: &OmitZeroVal{A: 2}, Ptr: &OmitZeroPtrRecv{A: 3}, X: 1}, 3},
	} {
		b, err := tc.in.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		var buf bytes.Buffer
		if err := msgp.Encode(&buf, &tc.in); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf.Bytes(), b) {
			t.Fatalf("EncodeMsg and MarshalMsg differ for %+v", tc.in)
		}
		if n, _, err := msgp.ReadMapHeaderBytes(b); err != nil || n != tc.fields {
			t.Fatalf("%+v: %d fields (%v), want %d", tc.in, n, err, tc.fields)
		}
		var out OmitZeroPtr
		if _, err := out.UnmarshalMsg(b); err != nil {
			t.Fatal(err)
		}
		if tc.fields == 3 && !reflect.DeepEqual(out, tc.in) {
			t.Fatalf("got %+v, want %+v", out, tc.in)
		}
	}
}
