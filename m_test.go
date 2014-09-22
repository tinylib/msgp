package main

import (
	"os"
	"os/exec"
	"testing"
)

// Here we generate code from ./_generated/def.go
// and then run the unit tests there on the generated code.
func TestGenerated(t *testing.T) {
	file := "def.go"
	pkg := "github.com/philhofer/msgp/_generated"
	err := DoAll(GOPATH, pkg, file)
	if err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "test", "-v", "github.com/philhofer/msgp/_generated")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}
