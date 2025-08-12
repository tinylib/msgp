package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/tinylib/msgp/gen"
)

const showGeneratedFile = false

func generate(t *testing.T, content string) (string, error) {
	tempDir := t.TempDir()

	mainFilename := filepath.Join(tempDir, "main.go")

	fd, err := os.OpenFile(mainFilename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0o600)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	if _, err := fd.WriteString(content); err != nil {
		return "", err
	}

	mode := gen.Encode | gen.Decode | gen.Size | gen.Marshal | gen.Unmarshal | gen.Test
	if err := Run(mainFilename, mode, false); err != nil {
		return "", err
	}

	if showGeneratedFile {
		mainGenFilename := strings.TrimSuffix(mainFilename, ".go") + "_gen.go"
		content, err := os.ReadFile(mainGenFilename)
		if err != nil {
			return "", err
		}
		t.Logf("generated %s content:\n%s", mainGenFilename, content)
	}

	return mainFilename, nil
}

func goExec(t *testing.T, mainFilename string, test bool) {
	mainGenFilename := strings.TrimSuffix(mainFilename, ".go") + "_gen.go"

	args := []string{"run", mainFilename, mainGenFilename}
	if test {
		mainGenTestFilename := strings.TrimSuffix(mainFilename, ".go") + "_gen_test.go"
		args = []string{"test", mainFilename, mainGenFilename, mainGenTestFilename}
	}

	output, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("go run failed: %v, output:\n%s", err, output)
	}
}
