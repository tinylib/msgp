package _generated

import (
	"bytes"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

// TestIssue463 checks that UnmarshalMsg handles omitempty,allownil like DecodeMsg.
func TestIssue463(t *testing.T) {
	for _, want := range [][]byte{nil, {}, {1}} {
		b := msgp.AppendMapHeader(nil, 1)
		b = msgp.AppendString(b, "B")
		if want == nil {
			b = msgp.AppendNil(b)
		} else {
			b = msgp.AppendBytes(b, want)
		}
		var u, d Issue463
		if _, err := u.UnmarshalMsg(b); err != nil {
			t.Fatalf("%#v: unmarshal: %v", want, err)
		}
		if err := d.DecodeMsg(msgp.NewReader(bytes.NewReader(b))); err != nil {
			t.Fatalf("%#v: decode: %v", want, err)
		}
		for name, got := range map[string][]byte{"unmarshal": u.B, "decode": d.B} {
			if (got == nil) != (want == nil) || !bytes.Equal(got, want) {
				t.Errorf("%s: got %#v, want %#v", name, got, want)
			}
		}
	}
}
