package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestIssue459RoundTrip(t *testing.T) {
	u := uint64(7)
	pu := &u
	m := map[string]uint64{"x": 1}
	in := Issue459{
		Count:       3,
		Buckets:     &Issue459Buckets{"a": 1, "b": 2},
		Map:         &map[string]uint64{"c": 3},
		MapStruct:   &map[string]Issue459Inner{"d": {A: 4}},
		MapPtr:      &map[string]*uint64{"e": &u},
		MapIntKey:   &map[int]string{5: "f"},
		Slice:       &[]string{"g", "h"},
		NamedSlice:  &Issue459Slice{"i"},
		Array:       &[3]int{1, 2, 3},
		DoublePtr:   &pu,
		DoubleMap:   func() **map[string]uint64 { p := &m; return &p }(),
		PtrSliceMap: &[]map[string]uint64{{"j": 6}},
	}

	b, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var out Issue459
	if _, err := out.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("unmarshal mismatch:\n got %+v\nwant %+v", out, in)
	}

	var buf bytes.Buffer
	if err := msgp.Encode(&buf, &in); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(buf.Bytes(), b) {
		t.Fatal("EncodeMsg and MarshalMsg output differ")
	}
	out = Issue459{}
	if err := msgp.Decode(&buf, &out); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("decode mismatch:\n got %+v\nwant %+v", out, in)
	}
}

func TestIssue459Duplicates(t *testing.T) {
	for _, field := range []string{"Buckets", "MapIntKey"} {
		t.Run(field, func(t *testing.T) {
			b := msgp.AppendMapHeader(nil, 1)
			b = msgp.AppendString(b, field)
			b = msgp.AppendMapHeader(b, 2)
			for range 2 {
				b = msgp.AppendString(b, "1")
				if field == "Buckets" {
					b = msgp.AppendUint64(b, 1)
				} else {
					b = msgp.AppendString(b, "v")
				}
			}
			var v Issue459
			_, err := v.UnmarshalMsg(b)
			expectDup(t, err)
			expectDup(t, v.DecodeMsg(msgp.NewReader(bytes.NewReader(b))))
		})
	}
}
