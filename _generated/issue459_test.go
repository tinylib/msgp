package _generated

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

// issue459RoundTrip checks that in comes back as want through
// MarshalMsg/UnmarshalMsg and EncodeMsg/DecodeMsg.
func issue459RoundTrip[T any, P interface {
	*T
	msgp.Marshaler
	msgp.Unmarshaler
	msgp.Encodable
	msgp.Decodable
}](t *testing.T, in, want T) {
	t.Helper()
	b, err := P(&in).MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := msgp.Encode(&buf, P(&in)); err != nil {
		t.Fatal(err)
	}
	var out T
	if _, err := P(&out).UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("unmarshal mismatch:\n got %+v\nwant %+v", out, want)
	}
	out = *new(T)
	if err := msgp.Decode(&buf, P(&out)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, want) {
		t.Fatalf("decode mismatch:\n got %+v\nwant %+v", out, want)
	}
}

func TestIssue459RoundTrip(t *testing.T) {
	u := uint64(7)
	pu := &u
	m := map[string]uint64{"x": 1}
	pm := &m
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
		DoubleMap:   &pm,
		PtrSliceMap: &[]map[string]uint64{{"j": 6}},
		Limited:     &[]byte{1, 2, 3},
	}
	issue459RoundTrip(t, in, in)
}

func TestIssue459Limit(t *testing.T) {
	b := msgp.AppendMapHeader(nil, 1)
	b = msgp.AppendString(b, "Limited")
	b = msgp.AppendBytes(b, make([]byte, 11))
	var v Issue459
	if _, err := v.UnmarshalMsg(b); !errors.Is(err, msgp.ErrLimitExceeded) {
		t.Fatalf("unmarshal: want ErrLimitExceeded, got %v", err)
	}
	if err := v.DecodeMsg(msgp.NewReader(bytes.NewReader(b))); !errors.Is(err, msgp.ErrLimitExceeded) {
		t.Fatalf("decode: want ErrLimitExceeded, got %v", err)
	}
}

func TestIssue459AllowNil(t *testing.T) {
	var (
		nilMap   map[string]uint64
		nilPtr   *map[string]uint64
		pNilMap  = &nilMap
		empty    = map[string]uint64{}
		pEmpty   = &empty
		filled   = map[string]uint64{"a": 1}
		pFilled  = &filled
		u        uint64
		nb       Issue459Buckets
		ns       []string
		nilBytes []byte
	)
	cases := []struct {
		name string
		in   Issue459AllowNil
		keep bool // decodes back to in, otherwise to all nil
	}{
		{name: "nil"},
		{
			name: "pointer to nil",
			in:   Issue459AllowNil{Map: &nilMap, Buckets: &nb, Slice: &ns, Bytes: &nilBytes, DoubleMap: &nilPtr},
		},
		{
			name: "pointer to pointer to nil",
			in:   Issue459AllowNil{DoubleMap: &pNilMap},
		},
		{
			name: "pointer to empty",
			in: Issue459AllowNil{
				Map: &map[string]uint64{}, Buckets: &Issue459Buckets{}, Slice: &[]string{},
				Bytes: &[]byte{}, DoubleMap: &pEmpty, Uint: &u, Array: &[3]int{}, Shim: &Issue459Shim{},
			},
			keep: true,
		},
		{
			name: "filled",
			in: Issue459AllowNil{
				Map: &map[string]uint64{"b": 2}, Buckets: &Issue459Buckets{"c": 3}, Slice: &[]string{"d"},
				Bytes: &[]byte{4}, DoubleMap: &pFilled, Uint: &u, Array: &[3]int{5, 6, 7}, Shim: &Issue459Shim{A: "e"},
			},
			keep: true,
		},
	}
	for _, c := range cases {
		var want Issue459AllowNil
		if c.keep {
			want = c.in
		}
		t.Run(c.name, func(t *testing.T) {
			issue459RoundTrip(t, c.in, want)
			issue459RoundTrip(t, Issue459AllowNilTuple(c.in), Issue459AllowNilTuple(want))
		})
	}
}

func TestIssue459OmitAllowNil(t *testing.T) {
	var (
		nilBytes *[]byte
		nilSlice *[]string
		nilMap   *map[string]uint64
		bs       []byte
		ss       []string
		m        map[string]uint64
		pbs, pss = &bs, &ss
		pm       = &m
		eb, es   = &[]byte{}, &[]string{}
		em       = &map[string]uint64{}
		fb, fs   = &[]byte{1}, &[]string{"a"}
		fm       = &map[string]uint64{"b": 2}
	)
	cases := []struct {
		name     string
		in, want Issue459OmitAllowNil
	}{
		{name: "nil"},
		{
			name: "pointer to nil pointer",
			in:   Issue459OmitAllowNil{Bytes: &nilBytes, Slice: &nilSlice, Map: &nilMap},
		},
		{
			// omitempty turns allownil off on encode, so nil goes out as empty.
			name: "pointer to pointer to nil",
			in:   Issue459OmitAllowNil{Bytes: &pbs, Slice: &pss, Map: &pm},
			want: Issue459OmitAllowNil{Bytes: &eb, Slice: &es, Map: &em},
		},
		{
			name: "filled",
			in:   Issue459OmitAllowNil{Bytes: &fb, Slice: &fs, Map: &fm},
			want: Issue459OmitAllowNil{Bytes: &fb, Slice: &fs, Map: &fm},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			issue459RoundTrip(t, c.in, c.want)
		})
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
