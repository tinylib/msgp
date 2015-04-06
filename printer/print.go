package printer

import (
	"bufio"
	"fmt"
	"github.com/tinylib/msgp/gen"
	"github.com/tinylib/msgp/parse"
	"github.com/ttacon/chalk"
	"golang.org/x/tools/imports"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func infof(s string, v ...interface{}) {
	fmt.Printf(chalk.Magenta.Color(s), v...)
}

// PrintFile prints the methods for the provided list
// of elements to the given file name and canonical
// package path.
func PrintFile(file string, f *parse.FileSet, mode gen.Method) error {
	err := generate(file, f, mode)
	if err != nil {
		return err
	}
	infof(">>> Generated \"%s\"\n", file)
	err = format(file)
	if err != nil {
		return err
	}
	infof(">>> Formatted \"%s\"\n", file)
	infof(">>> Done.\n")
	return nil
}

func format(file string) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	data, err = imports.Process(file, data, nil)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0600)
}

func generate(file string, f *parse.FileSet, mode gen.Method) error {
	if mode&^gen.Test == 0 {
		return nil
	}

	outfile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outfile.Close()
	outwr := bufio.NewWriter(outfile)
	defer outwr.Flush()

	err = writePkgHeader(outwr, f.Package)
	if err != nil {
		return err
	}
	err = writeImportHeader(outwr, "github.com/tinylib/msgp/msgp")
	if err != nil {
		return err
	}

	var testwr *bufio.Writer
	if mode&gen.Test == gen.Test {
		tfname := strings.TrimSuffix(file, ".go") + "_test.go"
		testfile, err := os.Create(strings.TrimSuffix(file, ".go") + "_test.go")
		if err != nil {
			return err
		}
		infof(">>> Tests in \"%s\"\n", tfname)
		defer testfile.Close()
		testwr = bufio.NewWriter(testfile)
		defer testwr.Flush()
		err = writePkgHeader(testwr, f.Package)
		if err != nil {
			return err
		}
		err = writeImportHeader(testwr, "bytes", "github.com/tinylib/msgp/msgp", "testing")
	}
	return f.PrintTo(gen.NewPrinter(mode, outwr, testwr))
}

func writePkgHeader(w io.Writer, name string) error {
	_, err := fmt.Fprintln(w, "package", name, "\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(w, "// NOTE: THIS FILE WAS PRODUCED BY THE\n// MSGP CODE GENERATION TOOL (github.com/tinylib/msgp)\n// DO NOT EDIT\n\n")
	return err
}

func writeImportHeader(w io.Writer, imports ...string) error {
	_, err := io.WriteString(w, "import (\n")
	if err != nil {
		return err
	}
	for _, im := range imports {
		_, err = io.WriteString(w, fmt.Sprintf("\t%q\n", im))
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w, ")\n\n")
	return err
}
