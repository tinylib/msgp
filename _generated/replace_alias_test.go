package _generated

import "testing"

func TestReplaceAlias(t *testing.T) {
	in := ReplAlias{
		Chain:  ReplChainOuter{{1}},
		Strs:   ReplStrs{"a"},
		Shim:   ReplAliasShim{A: "b"},
		Target: ReplOnTarget{2},
		List:   ReplList{"c"},
		PS:     &ReplStruct{A: []int{3}},
	}
	issue459RoundTrip(t, in, in)

	g := ReplGeneric[Int64, *Int64]{K: 4, M: map[string]int{"d": 5}, PI: &ReplGenInner{B: 6}, R: "e"}
	issue459RoundTrip(t, g, g)
}
