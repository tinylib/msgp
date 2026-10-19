package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/tinylib/msgp/gen"
)

func TestStreamTestsWithoutSizer(t *testing.T) {
	for _, tc := range []struct {
		name        string
		directive   string
		customSizer string
	}{
		{name: "ignored", directive: "//msgp:size ignore Stream"},
		{name: "regex", directive: "//msgp:size ignore regex:^Stream$"},
		{name: "generated"},
		{name: "custom", directive: "//msgp:size ignore Stream", customSizer: `
var sizeCalls int
func (s *Stream) Msgsize() int { sizeCalls++; return 1 }
`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			filename := filepath.Join(dir, "main.go")
			custom := tc.customSizer
			if tc.name == "generated" || custom != "" {
				custom += "\nvar _ msgp.Sizer = (*Stream)(nil)\n"
			}
			content := fmt.Sprintf(streamWithoutSizer, tc.directive, custom)
			if err := os.WriteFile(filename, []byte(content), 0600); err != nil {
				t.Fatal(err)
			}
			// The CLI's -marshal=false mode still supports stream methods and tests.
			mode := gen.Encode | gen.Decode | gen.Size | gen.Test
			if err := Run(filename, mode, false); err != nil {
				t.Fatalf("generate failed: %v", err)
			}
			goExec(t, filename, false)
			args := []string{"test", "-v", filename, filepath.Join(dir, "main_gen.go"), filepath.Join(dir, "main_gen_test.go")}
			if tc.customSizer != "" {
				probe := filepath.Join(dir, "sizer_test.go")
				if err := os.WriteFile(probe, []byte(customSizerTest), 0600); err != nil {
					t.Fatal(err)
				}
				args = append(args, probe)
			}
			output, err := exec.Command("go", args...).CombinedOutput()
			if err != nil {
				t.Fatalf("generated stream tests failed: %v\n%s", err, output)
			}
			if tc.customSizer != "" && !bytes.Contains(output, []byte("WARNING:")) {
				t.Fatalf("generated stream test did not warn about undersized custom Msgsize: %s", output)
			}
		})
	}
}

const streamWithoutSizer = `package main

import (
    "bytes"
    "github.com/tinylib/msgp/msgp"
)

%s

type Stream struct { A int }
%s

func main() {
    input := Stream{A: 42}
    var stream bytes.Buffer
    if err := msgp.Encode(&stream, &input); err != nil { panic(err) }
    if !bytes.Equal(stream.Bytes(), []byte{0x81, 0xa1, 'A', 42}) {
        panic("unexpected wire bytes")
    }
    var output Stream
    if err := msgp.Decode(&stream, &output); err != nil { panic(err) }
    if output != input { panic("stream round trip changed value") }
}
`

const customSizerTest = `package main

import "testing"

func TestGeneratedStreamTestChecksCustomSizer(t *testing.T) {
    sizeCalls = 0
    TestEncodeDecodeStream(t)
    if sizeCalls != 1 {
        t.Fatalf("generated stream test called Msgsize %d times, want 1", sizeCalls)
    }
}
`
