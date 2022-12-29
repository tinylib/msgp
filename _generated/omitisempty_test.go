package _generated

import (
	"bytes"
	"testing"
)

func TestOmitIsEmpty(t *testing.T) {

	t.Run("OmitIsEmptyExt_not_empty", func(t *testing.T) {

		z := OmitIsEmpty0{AExt: OmitIsEmptyExt{a: 1}}
		b, err := z.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(b, []byte("aext")) {
			t.Errorf("expected to find aext in bytes %X", b)
		}
		z = OmitIsEmpty0{}
		_, err = z.UnmarshalMsg(b)
		if err != nil {
			t.Fatal(err)
		}
		if z.AExt.a != 1 {
			t.Errorf("z.AExt.a expected 1 but got %d", z.AExt.a)
		}

	})

	t.Run("OmitIsEmptyExt_negative", func(t *testing.T) {

		z := OmitIsEmpty0{AExt: OmitIsEmptyExt{a: -1}} // negative value should act as empty, via IsEmpty() call
		b, err := z.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(b, []byte("aext")) {
			t.Errorf("expected to not find aext in bytes %X", b)
		}
		z = OmitIsEmpty0{}
		_, err = z.UnmarshalMsg(b)
		if err != nil {
			t.Fatal(err)
		}
		if z.AExt.a != 0 {
			t.Errorf("z.AExt.a expected 0 but got %d", z.AExt.a)
		}

	})
}
