package _generated

import "testing"

func TestIssue462RoundTrip(t *testing.T) {
	m := Issue462Map{"b": 2}
	pm := &m
	n := 13
	p := Issue462Ptr(&n)
	alias := Issue462Alias{5, 4, 3, 2}
	in := Issue462{
		Map:       Issue462Map{"a": 1},
		PMap:      &Issue462Map{"c": 3},
		PPMap:     &pm,
		Slice:     Issue462Slice{"d"},
		PSlice:    &Issue462Slice{"e", "f"},
		Bytes:     Issue462Bytes{1, 2, 3, 4},
		PBytes:    &Issue462Bytes{5, 6, 7, 8},
		PInts:     &Issue462Ints{9, 10, 11},
		MapPSlice: map[string]*Issue462Slice{"g": {"h"}},
		SlicePMap: []*Issue462Map{{"i": 12}},
		Ptr:       p,
		PPtr:      &p,
		Alias:     Issue462Alias{1, 2, 3, 4},
		PAlias:    &alias,
		Names:     Issue462Names{"j"},
		PNames:    &Issue462Names{"k"},
	}
	issue459RoundTrip(t, in, in)

	ps := &Issue462Slice{"p"}
	g := Issue462Gen[Int64, *Int64]{
		A:      7,
		PBytes: &Issue462Bytes{9, 8, 7, 6},
		Slice:  Issue462Slice{"m"},
		Names:  Issue462Names{"n"},
		PNames: &Issue462Names{"o"},
		PPS:    &ps,
		MPS:    map[string]*Issue462Slice{"q": {"r"}},
	}
	issue459RoundTrip(t, g, g)
}
