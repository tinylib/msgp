package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestNamedShim(t *testing.T) {
	for _, pointer := range []*NamedShimOriginal{nil, {NamedShimString: "pointer"}} {
		name := "pointer"
		if pointer == nil {
			name = "nil pointer"
		}
		t.Run(name, func(t *testing.T) {
			in := NamedShimContainer{
				Value:   NamedShimOriginal{NamedShimString: "value"},
				Pointer: pointer,
				Plain:   "plain",
				Slice:   []NamedShimOriginal{{NamedShimString: "slice"}},
				Array:   [1]NamedShimOriginal{{NamedShimString: "array"}},
				Map:     map[string]NamedShimOriginal{"key": {NamedShimString: "map"}},
			}
			// Assemble the expected map encoding independently of generated methods.
			want := msgp.AppendMapHeader(nil, 6)
			want = msgp.AppendString(want, "Value")
			want = msgp.AppendMapHeader(want, 1)
			want = msgp.AppendString(want, "Value")
			want = msgp.AppendString(want, "value")
			want = msgp.AppendString(want, "Pointer")
			if pointer == nil {
				want = msgp.AppendNil(want)
			} else {
				want = msgp.AppendMapHeader(want, 1)
				want = msgp.AppendString(want, "Value")
				want = msgp.AppendString(want, "pointer")
			}
			want = msgp.AppendString(want, "Plain")
			want = msgp.AppendString(want, "plain")
			for _, field := range []struct{ name, value string }{{"Slice", "slice"}, {"Array", "array"}} {
				want = msgp.AppendString(want, field.name)
				want = msgp.AppendArrayHeader(want, 1)
				want = msgp.AppendMapHeader(want, 1)
				want = msgp.AppendString(want, "Value")
				want = msgp.AppendString(want, field.value)
			}
			want = msgp.AppendString(want, "Map")
			want = msgp.AppendMapHeader(want, 1)
			want = msgp.AppendString(want, "key")
			want = msgp.AppendMapHeader(want, 1)
			want = msgp.AppendString(want, "Value")
			want = msgp.AppendString(want, "map")

			t.Run("MarshalMsg", func(t *testing.T) {
				got, err := in.MarshalMsg(nil)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want) {
					t.Fatalf("encoding = %x, want %x", got, want)
				}
			})
			t.Run("EncodeMsg", func(t *testing.T) {
				var buf bytes.Buffer
				writer := msgp.NewWriter(&buf)
				if err := in.EncodeMsg(writer); err != nil {
					t.Fatal(err)
				}
				if err := writer.Flush(); err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(buf.Bytes(), want) {
					t.Fatalf("encoding = %x, want %x", buf.Bytes(), want)
				}
			})
			t.Run("UnmarshalMsg", func(t *testing.T) {
				var out NamedShimContainer
				left, err := out.UnmarshalMsg(want)
				if err != nil {
					t.Fatal(err)
				}
				if len(left) != 0 || !reflect.DeepEqual(out, in) {
					t.Fatalf("decoded %#v with %x left, want %#v", out, left, in)
				}
			})
			t.Run("DecodeMsg", func(t *testing.T) {
				var out NamedShimContainer
				if err := out.DecodeMsg(msgp.NewReader(bytes.NewReader(want))); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(out, in) {
					t.Fatalf("decoded %#v, want %#v", out, in)
				}
			})
			t.Run("Msgsize", func(t *testing.T) {
				if got := in.Msgsize(); got < len(want) {
					t.Fatalf("Msgsize %d under-reports wire size %d", got, len(want))
				}
			})
		})
	}
}

func TestNamedShimDecodeError(t *testing.T) {
	// Valid MessagePack with a boolean where the named serializer expects a string.
	data := msgp.AppendMapHeader(nil, 1)
	data = msgp.AppendString(data, "Value")
	data = msgp.AppendMapHeader(data, 1)
	data = msgp.AppendString(data, "Value")
	data = msgp.AppendBool(data, true)

	t.Run("UnmarshalMsg", func(t *testing.T) {
		var out NamedShimContainer
		_, err := out.UnmarshalMsg(data)
		if cause, ok := msgp.Cause(err).(msgp.TypeError); !ok || cause.Method != msgp.StrType || cause.Encoded != msgp.BoolType {
			t.Fatalf("expected string/boolean type error, got %v", err)
		}
	})
	t.Run("DecodeMsg", func(t *testing.T) {
		var out NamedShimContainer
		err := out.DecodeMsg(msgp.NewReader(bytes.NewReader(data)))
		if cause, ok := msgp.Cause(err).(msgp.TypeError); !ok || cause.Method != msgp.StrType || cause.Encoded != msgp.BoolType {
			t.Fatalf("expected string/boolean type error, got %v", err)
		}
	})
}
