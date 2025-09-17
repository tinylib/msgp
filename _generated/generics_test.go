package _generated

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestGenericsMarshal(t *testing.T) {
	x := GenericTest[Fixed, *Fixed]{}
	x.B = &Fixed{A: 1.5}
	x.C = append(x.C, Fixed{A: 2.5})
	x.D = map[string]Fixed{"hello": Fixed{A: 2.5}}
	x.E.A = Fixed{A: 2.5}
	x.F = append(x.F, GenericTest2[Fixed, *Fixed, string]{A: Fixed{A: 3.5}})
	x.G = map[string]GenericTest2[Fixed, *Fixed, string]{"hello": {A: Fixed{A: 3.5}}}

	bts, err := x.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	got := GenericTest[Fixed, *Fixed]{}
	got.B = x.B // We must initialize this.
	*got.B = Fixed{}
	bts, err = got.UnmarshalMsg(bts)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(x, got) {
		t.Errorf("\n got=%#v\nwant=%#v", got, x)
	}
}

func TestGenericsEncode(t *testing.T) {
	x := GenericTest[Fixed, *Fixed]{}
	x.B = &Fixed{A: 1.5}
	x.C = append(x.C, Fixed{A: 2.5})
	x.D = map[string]Fixed{"hello": Fixed{A: 2.5}}
	x.E.A = Fixed{A: 2.5}
	x.F = append(x.F, GenericTest2[Fixed, *Fixed, string]{A: Fixed{A: 3.5}})
	x.G = map[string]GenericTest2[Fixed, *Fixed, string]{"hello": {A: Fixed{A: 3.5}}}

	var buf bytes.Buffer
	w := msgp.NewWriter(&buf)
	err := x.EncodeMsg(w)
	if err != nil {
		t.Fatal(err)
	}
	w.Flush()
	got := GenericTest[Fixed, *Fixed]{}
	got.B = x.B // We must initialize this.
	*got.B = Fixed{}
	r := msgp.NewReader(&buf)
	err = got.DecodeMsg(r)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(x, got) {
		t.Errorf("\n got=%#v\nwant=%#v", got, x)
	}
}
