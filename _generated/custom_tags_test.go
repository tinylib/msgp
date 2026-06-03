package _generated

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestCustomTags(t *testing.T) {
	ts := CustomTags{
		Field1: "v1",
		Field2: "v2",
		Field3: "v3",
		Field4: "v4",
		Field5: "v5",
		Field6: "v6",
		Field7: "v7",
		Field8: "v8",
	}
	wantKeys := map[string]any{
		"f1_primary":  "v1",
		"f2_fallback": "v2",
		"f3_primary":  "v3",
		"f4_msg":      "v4",
		"f5_msgpack":  "v5",
		"f6_fallback": "v6",
		"f7_fallback": "v7",
		"f8_primary":  "v8",
	}

	t.Run("EncodeDecode", func(t *testing.T) {
		var b bytes.Buffer
		if err := msgp.Encode(&b, &ts); err != nil {
			t.Fatal(err)
		}
		var got CustomTags
		if err := msgp.Decode(&b, &got); err != nil {
			t.Fatal(err)
		}
		if got != ts {
			t.Errorf("got %+v, want %+v", got, ts)
		}
	})

	t.Run("MarshalUnmarshal", func(t *testing.T) {
		buf, err := ts.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		var got CustomTags
		left, err := got.UnmarshalMsg(buf)
		if err != nil {
			t.Fatal(err)
		}
		if len(left) != 0 {
			t.Errorf("%d bytes left after unmarshal", len(left))
		}
		if got != ts {
			t.Errorf("got %+v, want %+v", got, ts)
		}
	})

	t.Run("WireKeys", func(t *testing.T) {
		buf, err := ts.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		var jsonBuf bytes.Buffer
		if _, err := msgp.UnmarshalAsJSON(&jsonBuf, buf); err != nil {
			t.Fatal(err)
		}
		got := map[string]any{}
		if err := json.Unmarshal(jsonBuf.Bytes(), &got); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(got, wantKeys) {
			t.Errorf("got %v, want %v", got, wantKeys)
		}
	})

	t.Run("Msgsize", func(t *testing.T) {
		buf, err := ts.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		if est := ts.Msgsize(); est < len(buf) {
			t.Errorf("Msgsize %d underestimates actual %d", est, len(buf))
		}
	})
}
