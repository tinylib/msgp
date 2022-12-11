//go:build (amd64 && go1.18) || (darwin && go1.18)
// +build amd64,go1.18 darwin,go1.18

package tinygotest

// NOTE: this is intended to be run with `go test` not `tinygo test`,
// as it then in turn performs various `tinygo build` commands and
// verifies the output.

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// empty gives us a baseline to compare sizes (and to prove that tinygo itself isn't broken)

func TestEmptyRun(t *testing.T) {
	// build and run native executable
	tinygoBuild(t, "empty")
	run(t, "empty", "nativebin")
}

func TestEmptyBuild(t *testing.T) {
	// build for various environments
	tinygoBuild(t, "empty", buildOnlyTargets...)
	reportSizes(t, "empty")
}

// roundtrip provides reasonable coverage for various cases, but is
// fairly bloated in terms of size compared to simpler examples

func TestRoundtripRun(t *testing.T) {
	// build and run native executable
	goGenerate(t, "roundtrip")
	tinygoBuild(t, "roundtrip")
	run(t, "roundtrip", "nativebin")
}

func TestRoundtripBuild(t *testing.T) {
	// build for various environments
	goGenerate(t, "roundtrip")
	tinygoBuild(t, "roundtrip", buildOnlyTargets...)
	reportSizes(t, "roundtrip")
}

// simple_roundtrip is a minimal marshal+unmarshal example and should show
// the baseline size when only using minimal functionality and marshal+unmarshal

func TestSimpleRoundtripRun(t *testing.T) {
	// build and run native executable
	goGenerate(t, "simple_roundtrip")
	tinygoBuild(t, "simple_roundtrip")
	run(t, "simple_roundtrip", "nativebin")
}

func TestSimpleRoundtripBuild(t *testing.T) {
	// build for various environments
	goGenerate(t, "simple_roundtrip")
	tinygoBuild(t, "simple_roundtrip", buildOnlyTargets...)
	sizes := reportSizes(t, "simple_roundtrip")
	// this was ~13k at the time this test was written,
	// if it bloats up to over 20k, we assume something is wrong
	// (e.g. "fmt" or other large package being using)
	sz := sizes["arduino-nano33.bin"]
	if sz > 20000 {
		t.Errorf("arduino-nano33.bin is larger than expected: %d", sz)
	}
}

// simple_marshal is just the MarshalMsg part

func TestSimpleMarshalRun(t *testing.T) {
	// build and run native executable
	goGenerate(t, "simple_marshal")
	tinygoBuild(t, "simple_marshal")
	run(t, "simple_marshal", "nativebin")
}

func TestSimpleMarshalBuild(t *testing.T) {
	// build for various environments
	goGenerate(t, "simple_marshal")
	tinygoBuild(t, "simple_marshal", buildOnlyTargets...)
	reportSizes(t, "simple_marshal")
}

// simple_unmarshal is just the UnmarshalMsg part

func TestSimpleUnmarshalRun(t *testing.T) {
	// build and run native executable
	goGenerate(t, "simple_unmarshal")
	tinygoBuild(t, "simple_unmarshal")
	run(t, "simple_unmarshal", "nativebin")
}

func TestSimpleUnmarshalBuild(t *testing.T) {
	// build for various environments
	goGenerate(t, "simple_unmarshal")
	tinygoBuild(t, "simple_unmarshal", buildOnlyTargets...)
	reportSizes(t, "simple_unmarshal")
}

// simple_bytes_append is AppendX() methods without code generation

func TestSimpleBytesAppendRun(t *testing.T) {
	// build and run native executable
	goGenerate(t, "simple_bytes_append")
	tinygoBuild(t, "simple_bytes_append")
	run(t, "simple_bytes_append", "nativebin")
}

func TestSimpleBytesAppendBuild(t *testing.T) {
	// build for various environments
	goGenerate(t, "simple_bytes_append")
	tinygoBuild(t, "simple_bytes_append", buildOnlyTargets...)
	reportSizes(t, "simple_bytes_append")
}

// simple_bytes_append is ReadX() methods without code generation

func TestSimpleBytesReadRun(t *testing.T) {
	// build and run native executable
	goGenerate(t, "simple_bytes_read")
	tinygoBuild(t, "simple_bytes_read")
	run(t, "simple_bytes_read", "nativebin")
}

func TestSimpleBytesReadBuild(t *testing.T) {
	// build for various environments
	goGenerate(t, "simple_bytes_read")
	tinygoBuild(t, "simple_bytes_read", buildOnlyTargets...)
	reportSizes(t, "simple_bytes_read")
}

// do builds for these tinygo boards/environments to make sure
// it at least compiles
var buildOnlyTargets = []string{
	"arduino-nano33", // ARM Cortex-M0, SAMD21
	"feather-m4",     // ARM Cortex-M4, SAMD51
	"wasm",           // WebAssembly

	// "arduino-nano",       // AVR - roundtrip build currently fails with: could not store type code number inside interface type code
	// "esp32-coreboard-v2", // ESP - xtensa seems to require additional setup, currently errors with: No available targets are compatible with triple "xtensa"
}

func run(t *testing.T, dir, exe string, args ...string) {
	t.Helper()
	cmd := exec.Command("./" + exe)
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Dir = filepath.Join(wd, "testdata", dir)
	cmd.Args = args
	b, err := cmd.CombinedOutput()
	if len(bytes.TrimSpace(b)) > 0 {
		t.Logf("%s: %s %v; output: %s", dir, exe, args, b)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func goGenerate(t *testing.T, dir string) {
	t.Helper()
	t.Logf("%s: go generate", dir)
	cmd := exec.Command("go", "generate")
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Dir = filepath.Join(wd, "testdata", dir)
	b, err := cmd.CombinedOutput()
	if len(bytes.TrimSpace(b)) > 0 {
		t.Logf("%s: go generate output: %s", dir, b)
	}
	if err != nil {
		t.Fatal(err)
	}
}

func tinygoBuild(t *testing.T, dir string, targets ...string) {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	dirabs := filepath.Join(wd, "testdata", dir)

	// no targets specified implies just build the native executable
	if len(targets) == 0 {
		targets = []string{""}
	}

	for _, tgt := range targets {

		ext := ".bin"
		if tgt == "wasm" {
			ext = ".wasm"
		}

		var args []string
		if tgt == "" { // empty target means the native platform
			args = []string{"build", "-o=nativebin", "."}
		} else {
			args = []string{"build", "-target=" + tgt, "-o=" + tgt + ext, "."}
		}

		t.Logf("%s: tinygo %v", dir, args)
		cmd := exec.Command("tinygo", args...)
		cmd.Dir = dirabs
		b, err := cmd.CombinedOutput()
		if len(bytes.TrimSpace(b)) > 0 {
			t.Logf("%s: tinygo build %v; output: %s", dir, args, b)
		}
		if err != nil {
			t.Fatal(err)
		}
	}
}

var spacePad = strings.Repeat(" ", 64)

func reportSizes(t *testing.T, dir string) (ret map[string]int64) {
	ret = make(map[string]int64)
	var fnl []string
	for _, gl := range []string{"*.bin", "*.wasm"} {
		fnl2, err := filepath.Glob(filepath.Join("testdata", dir, gl))
		if err != nil {
			t.Fatal(err)
		}
		fnl = append(fnl, fnl2...)
	}
	for _, fn := range fnl {
		st, err := os.Stat(fn)
		if err != nil {
			t.Fatal(err)
		}
		_, jfn := filepath.Split(fn)
		t.Logf("size report - %s %6d bytes", (jfn + spacePad)[:24], st.Size())
		ret[jfn] = st.Size()
	}
	return ret
}
