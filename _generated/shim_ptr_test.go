package _generated

import "testing"

func TestShimPtrRoundTrip(t *testing.T) {
	in := ShimPtr{
		S:  ShimPtrStr{A: "a"},
		B:  ShimPtrBytes{B: "b"},
		E:  ShimPtrErr{C: "c"},
		MS: map[string]ShimPtrStr{"k": {A: "d"}},
		PS: &ShimPtrStr{A: "e"},
		PE: &ShimPtrErr{C: "f"},
	}
	issue459RoundTrip(t, in, in)

	issue459RoundTrip(t, ShimPtr{S: ShimPtrStr{}}, ShimPtr{S: ShimPtrStr{}, MS: map[string]ShimPtrStr{}})

	in.E.C = "bad"
	b, err := in.MarshalMsg(nil)
	if err != nil {
		t.Fatal(err)
	}
	var out ShimPtr
	if _, err := out.UnmarshalMsg(b); err == nil {
		t.Fatal("want shim error")
	}
}
