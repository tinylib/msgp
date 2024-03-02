package _generated

import (
	"bytes"
	"testing"
)

func TestOmitZero(t *testing.T) {

	t.Run("OmitZeroExt_not_empty", func(t *testing.T) {

		z := OmitZero0{AExt: OmitZeroExt{a: 1}}
		b, err := z.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Contains(b, []byte("aext")) {
			t.Errorf("expected to find aext in bytes %X", b)
		}
		z = OmitZero0{}
		_, err = z.UnmarshalMsg(b)
		if err != nil {
			t.Fatal(err)
		}
		if z.AExt.a != 1 {
			t.Errorf("z.AExt.a expected 1 but got %d", z.AExt.a)
		}

	})

	t.Run("OmitZeroExt_negative", func(t *testing.T) {

		z := OmitZero0{AExt: OmitZeroExt{a: -1}} // negative value should act as empty, via IsEmpty() call
		b, err := z.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Contains(b, []byte("aext")) {
			t.Errorf("expected to not find aext in bytes %X", b)
		}
		z = OmitZero0{}
		_, err = z.UnmarshalMsg(b)
		if err != nil {
			t.Fatal(err)
		}
		if z.AExt.a != 0 {
			t.Errorf("z.AExt.a expected 0 but got %d", z.AExt.a)
		}

	})

	t.Run("OmitZeroTuple", func(t *testing.T) {

		// make sure tuple encoding isn't affected by omitempty or omitzero

		z := OmitZero0{OmitZeroTuple: OmitZeroTuple{FieldA: "", FieldB: "", FieldC: "fcval"}}
		b, err := z.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		// verify the exact binary encoding, that the values follow each other without field names
		if !bytes.Contains(b, []byte{0xA0, 0xA0, 0xA5, 'f', 'c', 'v', 'a', 'l'}) {
			t.Errorf("failed to find expected bytes in %X", b)
		}
		z = OmitZero0{}
		_, err = z.UnmarshalMsg(b)
		if err != nil {
			t.Fatal(err)
		}
		if z.OmitZeroTuple.FieldA != "" ||
			z.OmitZeroTuple.FieldB != "" ||
			z.OmitZeroTuple.FieldC != "fcval" {
			t.Errorf("z.OmitZeroTuple unexpected value: %#v", z.OmitZeroTuple)
		}

	})
}
