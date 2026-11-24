package _generated

import "testing"

func TestIssue465(t *testing.T) {
	s := Issue465Str("b")
	st := Issue465Struct{A: "d"}
	issue459RoundTrip(t, Issue465{
		Str:     "a",
		PStr:    &s,
		Struct:  Issue465Struct{A: "c"},
		PStruct: &st,
		Bytes:   Issue465Bytes{1},
		Ext:     Issue465Ext{2},
		Unexp:   "e",
	}, Issue465{
		Str:     "a",
		PStr:    &s,
		Struct:  Issue465Struct{A: "c"},
		PStruct: &st,
		Bytes:   Issue465Bytes{1},
		Ext:     Issue465Ext{2},
		Unexp:   "e",
	})
	// allownil keeps nil and empty apart for the nilable types.
	issue459RoundTrip(t, Issue465{}, Issue465{})
	issue459RoundTrip(t, Issue465{Bytes: Issue465Bytes{}, Ext: Issue465Ext{}}, Issue465{Bytes: Issue465Bytes{}, Ext: Issue465Ext{}})
}
