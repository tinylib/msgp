package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestIssue465(t *testing.T) {
	s := Issue465Str("b")
	st := Issue465Struct{A: "d"}
	issue459RoundTrip(t, Issue465{
		Str:     "a",
		PStr:    &s,
		Struct:  Issue465Struct{A: "c"},
		PStruct: &st,
		Bytes:   Issue465Bytes{1},
		Ext:     Issue465Ext{2},
		Unexp:   "e",
	}, Issue465{
		Str:     "a",
		PStr:    &s,
		Struct:  Issue465Struct{A: "c"},
		PStruct: &st,
		Bytes:   Issue465Bytes{1},
		Ext:     Issue465Ext{2},
		Unexp:   "e",
	})
	// allownil keeps nil and empty apart for the nilable types.
	issue459RoundTrip(t, Issue465{}, Issue465{})
	issue459RoundTrip(t, Issue465{Bytes: Issue465Bytes{}, Ext: Issue465Ext{}}, Issue465{Bytes: Issue465Bytes{}, Ext: Issue465Ext{}})
}

func TestIssue465Omit(t *testing.T) {
	filled := func() Issue465Omit {
		return Issue465Omit{Str: "a", Struct: Issue465Struct{A: "b"}, Unexp: "c", Bytes: Issue465Bytes{1}}
	}
	// The non-nilable fields are always written; the nil Bytes is omitted.
	var zero Issue465Omit
	b, err := zero.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	if sz, _, err := msgp.ReadMapHeaderBytes(b); err != nil || sz != 3 {
		t.Fatalf("map size %d (%v), want 3", sz, err)
	}
	// clearomitted zeroes fields missing from the input.
	for _, in := range [][]byte{b, msgp.AppendMapHeader(nil, 0)} {
		got := filled()
		if _, err := got.UnmarshalMsg(in); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, Issue465Omit{}) {
			t.Errorf("unmarshal %x: got %+v", in, got)
		}
		got = filled()
		if err := got.DecodeMsg(msgp.NewReader(bytes.NewReader(in))); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, Issue465Omit{}) {
			t.Errorf("decode %x: got %+v", in, got)
		}
	}
	issue459RoundTrip(t, filled(), filled())

	g := Issue465Gen[Int64, *Int64]{A: 1, Struct: Issue465GenStruct{A: "x"}}
	issue459RoundTrip(t, g, g)
	if _, err := g.UnmarshalMsg(msgp.AppendMapHeader(nil, 0)); err != nil {
		t.Fatal(err)
	}
	if g.Struct != (Issue465GenStruct{}) {
		t.Errorf("clearomitted: got %+v", g.Struct)
	}
}
