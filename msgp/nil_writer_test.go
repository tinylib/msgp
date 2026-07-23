package msgp

import "testing"

func TestNewWriterNil(t *testing.T) {
	w := NewWriter(nil)
	if err := w.WriteString("hi"); err != nil {
		t.Fatal(err)
	}
	if err := w.Flush(); err != nil {
		t.Fatal(err)
	}
}
