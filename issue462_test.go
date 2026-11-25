package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tinylib/msgp/gen"
)

func TestIssue462PointerToReplacedComposite(t *testing.T) {
	tempDir := t.TempDir()
	mainFilename := filepath.Join(tempDir, "main.go")
	extFilename := filepath.Join(tempDir, "ext.go")
	if err := os.WriteFile(mainFilename, []byte(issue462Main), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(extFilename, []byte(issue462Ext), 0o600); err != nil {
		t.Fatal(err)
	}

	mode := gen.Encode | gen.Decode | gen.Size | gen.Marshal | gen.Unmarshal
	if err := Run(mainFilename, mode, false); err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	mainGenFilename := strings.TrimSuffix(mainFilename, ".go") + "_gen.go"
	output, err := exec.Command("go", "run", extFilename, mainFilename, mainGenFilename).CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v, output:\n%s", err, output)
	}
}

const issue462Ext = `package main

type MapAny map[string]any
type SliceString []string
`

const issue462Main = `package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

//msgp:replace MapAny with:map[string]any
//msgp:replace SliceString with:[]string

type S struct {
	M *MapAny
	S *SliceString
}

func main() {
	in := S{
		M: &MapAny{"a": "x"},
		S: &SliceString{"y", "z"},
	}

	data, err := in.MarshalMsg(nil)
	if err != nil {
		fmt.Println("marshal failed:", err)
		os.Exit(1)
	}
	var out S
	if _, err := out.UnmarshalMsg(data); err != nil {
		fmt.Println("unmarshal failed:", err)
		os.Exit(1)
	}
	if !reflect.DeepEqual(in, out) {
		fmt.Printf("unmarshal mismatch: got %#v want %#v\n", out, in)
		os.Exit(1)
	}

	var buf bytes.Buffer
	if err := msgp.Encode(&buf, &in); err != nil {
		fmt.Println("encode failed:", err)
		os.Exit(1)
	}
	out = S{}
	if err := msgp.Decode(&buf, &out); err != nil {
		fmt.Println("decode failed:", err)
		os.Exit(1)
	}
	if !reflect.DeepEqual(in, out) {
		fmt.Printf("decode mismatch: got %#v want %#v\n", out, in)
		os.Exit(1)
	}

}
`
