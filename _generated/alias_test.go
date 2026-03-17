package _generated

import (
	"bytes"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestAliasUnmarshal(t *testing.T) {
	v1 := AliasedFieldsV1{
		CurrentName: "v1",
		MultiAlias:  1,
		NoAlias:     true,
	}
	v1Marshalled, err := v1.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	v2 := AliasedFieldsV2{
		CurrentName: "v2",
		MultiAlias:  2,
		NoAlias:     true,
	}
	v2Marshalled, err := v2.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	v3 := AliasedFieldsV3{
		CurrentName: "v3",
		MultiAlias:  3,
		NoAlias:     true,
	}
	v3Marshalled, err := v3.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}

	type deserializer interface {
		msgp.Unmarshaler
		msgp.Decodable
	}
	unmarshal := func(t testing.TB, unmarshaler deserializer, bts []byte) {
		t.Helper()

		remainingBytes, err := unmarshaler.UnmarshalMsg(bts)
		if err != nil {
			t.Fatal(err)
		}
		if len(remainingBytes) > 0 {
			t.Fatal("unexpected remaining bytes")
		}
	}
	decode := func(t testing.TB, decoder deserializer, bts []byte) {
		t.Helper()

		reader := msgp.NewReader(bytes.NewReader(bts))
		err := decoder.DecodeMsg(reader)
		if err != nil {
			t.Fatal(err)
		}

		if reader.Buffered() != 0 {
			t.Fatal("unexpected remaining bytes")
		}
	}

	for _, deserializeFn := range []func(testing.TB, deserializer, []byte){unmarshal, decode} {
		t.Run("v3 with aliases supports all", func(t *testing.T) {
			var v3Res AliasedFieldsV3
			deserializeFn(t, &v3Res, v1Marshalled)
			if v3Res.CurrentName != v1.CurrentName || v3Res.MultiAlias != v1.MultiAlias || v3Res.NoAlias != v1.NoAlias {
				t.Fatalf("mismatch: %+v != %+v", v3Res, v1)
			}

			deserializeFn(t, &v3Res, v2Marshalled)
			if v3Res.CurrentName != v2.CurrentName || v3Res.MultiAlias != v2.MultiAlias || v3Res.NoAlias != v2.NoAlias {
				t.Fatalf("mismatch: %+v != %+v", v3Res, v1)
			}

			deserializeFn(t, &v3Res, v3Marshalled)
			if v3Res.CurrentName != v3.CurrentName || v3Res.MultiAlias != v3.MultiAlias || v3Res.NoAlias != v3.NoAlias {
				t.Fatalf("mismatch: %+v != %+v", v3Res, v1)
			}
		})

		t.Run("v2 supports v1", func(t *testing.T) {
			var v2Res AliasedFieldsV2
			deserializeFn(t, &v2Res, v1Marshalled)
			if v2Res.CurrentName != v1.CurrentName || v2Res.MultiAlias != v1.MultiAlias || v2Res.NoAlias != v1.NoAlias {
				t.Fatalf("mismatch: %+v != %+v", v2Res, v1)
			}
		})
	}
}
