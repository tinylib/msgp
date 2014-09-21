package main

import (
	"bufio"
	"github.com/philhofer/msgp/gen"
	"github.com/philhofer/msgp/parse"
	"os"
	"os/exec"
	"testing"
)

func TestGenerated(t *testing.T) {
	els, err := parse.GetElems("./_generated/def.go")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Create("./_generated/gen.go")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	wr := bufio.NewWriter(file)

	err = writePkgHeader(wr, "_generated")
	if err != nil {
		t.Fatal(err)
	}
	err = writeImportHeader(wr, "github.com/philhofer/msgp/enc", "io", "bytes")
	if err != nil {
		t.Fatal(err)
	}

	for _, el := range els {
		p, ok := el.(*gen.Ptr)
		if !ok {
			t.Errorf("encountered non-Ptr: %v", el)
			continue
		}

		gen.Propogate(p, "z")

		err = gen.WriteMarshalMethod(wr, p)
		if err != nil {
			t.Error(err)
		}
		err = gen.WriteUnmarshalMethod(wr, p)
		if err != nil {
			t.Error(err)
		}
		err = gen.WriteEncoderMethod(wr, p)
		if err != nil {
			t.Error(err)
		}
		err = gen.WriteDecoderMethod(wr, p)
		if err != nil {
			t.Error(err)
		}
	}

	err = wr.Flush()
	if err != nil {
		t.Error(err)
	}

	cmd := exec.Command("go", "test", "-v", "-bench=.", "github.com/philhofer/msgp/_generated")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}
