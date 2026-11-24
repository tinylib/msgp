package _generated

import "testing"

func TestPtrPtrRoundTrip(t *testing.T) {
	in1, in2 := &PtrPtrInner{A: 1}, &PtrPtrInner{A: 2}
	pin2 := &in2
	ext := &PtrPtrExt{B: "x"}
	in := PtrPtr{
		F:   &in1,
		FFF: &pin2,
		M:   map[string]**PtrPtrInner{"k": &in1},
		S:   []**PtrPtrInner{&in2},
		E:   &ext,
		ME:  map[string]**PtrPtrExt{"k": &ext},
	}
	issue459RoundTrip(t, in, in)
}
