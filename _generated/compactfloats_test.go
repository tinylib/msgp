package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestCompactFloats(t *testing.T) {
	// Constant that can be represented in f32 without loss
	const f32ok = -1e2
	allF32 := Floats{
		A:     f32ok,
		B:     f32ok,
		Slice: []float64{f32ok, f32ok},
		Map:   map[string]float64{"a": f32ok},
		F:     f32ok,
		OE:    f32ok,
	}
	asF32 := float32(f32ok)
	wantF32 := map[string]any{"A": asF32, "B": asF32, "F": asF32, "Map": map[string]any{"a": asF32}, "OE": asF32, "Slice": []any{asF32, asF32}}

	enc, err := allF32.MarshalMsg(nil)
	if err != nil {
		t.Error(err)
	}
	i, _, _ := msgp.ReadIntfBytes(enc)
	got := i.(map[string]any)
	if !reflect.DeepEqual(got, wantF32) {
		t.Errorf("want: %v, got: %v (diff may be types)", wantF32, got)
	}

	var buf bytes.Buffer
	en := msgp.NewWriter(&buf)
	allF32.EncodeMsg(en)
	en.Flush()
	enc = buf.Bytes()
	i, _, _ = msgp.ReadIntfBytes(enc)
	got = i.(map[string]any)
	if !reflect.DeepEqual(got, wantF32) {
		t.Errorf("want: %v, got: %v (diff may be types)", wantF32, got)
	}

	const f64ok = -10e64
	allF64 := Floats{
		A:     f64ok,
		B:     f32ok,
		Slice: []float64{f64ok, f64ok},
		Map:   map[string]float64{"a": f64ok},
		F:     f64ok,
		OE:    f64ok,
	}
	asF64 := float64(f64ok)
	wantF64 := map[string]any{"A": asF64, "B": asF32, "F": asF64, "Map": map[string]any{"a": asF64}, "OE": asF64, "Slice": []any{asF64, asF64}}

	enc, err = allF64.MarshalMsg(nil)
	if err != nil {
		t.Error(err)
	}
	i, _, _ = msgp.ReadIntfBytes(enc)
	got = i.(map[string]any)
	if !reflect.DeepEqual(got, wantF64) {
		t.Errorf("want: %v, got: %v (diff may be types)", wantF64, got)
	}

	buf.Reset()
	en = msgp.NewWriter(&buf)
	allF64.EncodeMsg(en)
	en.Flush()
	enc = buf.Bytes()
	i, _, _ = msgp.ReadIntfBytes(enc)
	got = i.(map[string]any)
	if !reflect.DeepEqual(got, wantF64) {
		t.Errorf("want: %v, got: %v (diff may be types)", wantF64, got)
	}
}
