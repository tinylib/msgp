package _generated

import (
	"bytes"
	"os"
	"testing"
)

// Issue 94: shims were not propogated recursively,
// which caused shims that weren't at the top level
// to be silently ignored.
//
// The following line will generate an error after
// the code is generated if the generated code doesn't
// have the right identifier in it.
func TestIssue94(t *testing.T) {
	b, err := os.ReadFile("issue94_gen.go")
	if err != nil {
		t.Fatal(err)
	}
	const want = "timetostr"
	if !bytes.Contains(b, []byte(want)) {
		t.Errorf("generated code did not contain %q", want)
	}
}
