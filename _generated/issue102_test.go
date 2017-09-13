package _generated

import (
	"testing"
)

// Issue #102: serialize empty struct

func TestIssue102(t *testing.T) {
	var x, y Issue102
	b, err := x.MarshalMsg(nil)
	if err != nil {
		t.Error(err)
	}
	_, err = y.UnmarshalMsg(b)
	if err != nil {
		t.Error(err)
	}
}

func TestIssue102deep(t *testing.T) {
	var x, y Issue102deep
	b, err := x.MarshalMsg(nil)
	if err != nil {
		t.Error(err)
	}
	_, err = y.UnmarshalMsg(b)
	if err != nil {
		t.Error(err)
	}
}
