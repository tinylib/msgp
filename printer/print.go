package printer

import (
	"bytes"
	"fmt"
	"github.com/tinylib/msgp/gen"
	"github.com/tinylib/msgp/parse"
	"github.com/ttacon/chalk"
	"golang.org/x/tools/imports"
	"io"
	"io/ioutil"
	"strings"
)

func infof(s string, v ...interface{}) {
	fmt.Printf(chalk.Magenta.Color(s), v...)
}

// PrintFile prints the methods for the provided list
// of elements to the given file name and canonical
// package path.
func PrintFile(file string, f *parse.FileSet, mode gen.Method) error {
	out, tests, err := generate(f, mode)
	if err != nil {
		return err
	}
	err = format(file, out.Bytes())
	if err != nil {
		return err
	}
	infof(">>> Wrote and formatted \"%s\"\n", file)
	if tests != nil {
		testfile := strings.TrimSuffix(file, ".go") + "_test.go"
		err = format(testfile, tests.Bytes())
		if err != nil {
			return err
		}
		infof(">>> Wrote and formatted \"%s\"\n", testfile)
	}
	infof(">>> Done.\n")
	return nil
}

func format(file string, data []byte) error {
	out, err := imports.Process(file, data, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, out, 0600)
}

func generate(f *parse.FileSet, mode gen.Method) (*bytes.Buffer, *bytes.Buffer, error) {
	outbuf := bytes.NewBuffer(make([]byte, 0, 4096))
	writePkgHeader(outbuf, f.Package)
	writeImportHeader(outbuf, "github.com/tinylib/msgp/msgp")

	var testbuf *bytes.Buffer
	var testwr io.Writer
	if mode&gen.Test == gen.Test {
		testbuf = bytes.NewBuffer(make([]byte, 0, 4096))
		writePkgHeader(testbuf, f.Package)
		if mode&(gen.Encode|gen.Decode) != 0 {
			writeImportHeader(testbuf, "bytes", "github.com/tinylib/msgp/msgp", "testing")
		} else {
			writeImportHeader(testbuf, "github.com/tinylib/msgp/msgp", "testing")
		}
		testwr = testbuf
	}
	return outbuf, testbuf, f.PrintTo(gen.NewPrinter(mode, outbuf, testwr))
}

func writePkgHeader(b *bytes.Buffer, name string) {
	b.WriteString("package ")
	b.WriteString(name)
	b.WriteByte('\n')
	b.WriteString("// NOTE: THIS FILE WAS PRODUCED BY THE\n// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)\n// DO NOT EDIT\n\n")
}

func writeImportHeader(b *bytes.Buffer, imports ...string) {
	b.WriteString("import (\n")
	for _, im := range imports {
		fmt.Fprintf(b, "\t%q\n", im)
	}
	b.WriteString(")\n\n")
}
