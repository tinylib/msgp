package _generated

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestIssue464ZeroCopy(t *testing.T) {
	v := func(s string) *[]byte { b := []byte(s); return &b }
	pv := v("b")
	in := Issue464{
		Bytes:      v("a"),
		BytesPP:    &pv,
		SlicePtr:   []*[]byte{v("c")},
		MapPtr:     map[string]*[]byte{"k": v("d")},
		NotZC:      v("e"),
		LimitBytes: v("f"),
	}
	b, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var out Issue464
	if _, err := out.UnmarshalMsg(b); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("got %+v, want %+v", out, in)
	}
	inBuf := func(p []byte) bool {
		start := uintptr(unsafe.Pointer(unsafe.SliceData(b)))
		x := uintptr(unsafe.Pointer(unsafe.SliceData(p)))
		return x >= start && x < start+uintptr(len(b))
	}
	for name, p := range map[string][]byte{
		"Bytes":      *out.Bytes,
		"BytesPP":    **out.BytesPP,
		"SlicePtr":   *out.SlicePtr[0],
		"MapPtr":     *out.MapPtr["k"],
		"LimitBytes": *out.LimitBytes,
	} {
		if !inBuf(p) {
			t.Errorf("%s: zerocopy field does not reference the input", name)
		}
	}
	if inBuf(*out.NotZC) {
		t.Error("NotZC references the input")
	}
}
