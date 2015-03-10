package printer

import (
	"bufio"
	"fmt"
	"github.com/tinylib/msgp/gen"
	"github.com/ttacon/chalk"
	"golang.org/x/tools/imports"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// Mode determines what code
// is printed for each element.
type Mode uint32

const (
	Zero    Mode = 0
	Marshal      = 1 << iota // Generate MarshalMsg and UnmarshalMsg
	Encode                   // Generate EncodeMsg and DecodeMsg
	Test                     // Generate tests and benchmarks

	All = Marshal | Encode | Test // Generate everything
)

func infof(s string, v ...interface{}) {
	fmt.Printf(chalk.Magenta.Color(s), v...)
}

// PrintFile prints the methods for the provided list
// of elements to the given file name and canonical
// package path.
func PrintFile(file string, pkg string, els []gen.Elem, mode Mode) error {
	err := generate(file, pkg, els, mode)
	if err != nil {
		return err
	}
	err = format(file)
	if err != nil {
		return err
	}
	infof(">>> Done.\n")
	return nil
}

func format(file string) error {
	infof(">>> Formatting \"%s\"...\n", file)
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

func generate(file string, pkg string, els []gen.Elem, mode Mode) error {
	infof(">>> Generating \"%s\"...\n", file)
	if mode == Zero {
		return nil
	}

	outfile, err := os.Create(file)
	if err != nil {
		return err
	}
	defer outfile.Close()
	outwr := bufio.NewWriter(outfile)
	defer outwr.Flush()

	runner := make([]gen.Generator, 0, 7)
	if mode&Marshal == Marshal {
		runner = append(runner, gen.New(gen.Marshal, outwr), gen.New(gen.Size, outwr), gen.New(gen.Unmarshal, outwr))
	}
	if mode&Encode == Encode {
		runner = append(runner, gen.New(gen.Encode, outwr), gen.New(gen.Decode, outwr))
	}
	err = writePkgHeader(outwr, pkg)
	if err != nil {
		return err
	}
	err = writeImportHeader(outwr, "github.com/tinylib/msgp/msgp")
	if err != nil {
		return err
	}

	if mode&Test == Test {
		tfname := strings.TrimSuffix(file, ".go") + "_test.go"
		testfile, err := os.Create(strings.TrimSuffix(file, ".go") + "_test.go")
		if err != nil {
			return err
		}
		infof(">>> Generating \"%s\"...\n", tfname)
		defer testfile.Close()
		testwr := bufio.NewWriter(testfile)
		defer testwr.Flush()
		if mode&Marshal == Marshal {
			runner = append(runner, gen.New(gen.MarshalTest, testwr))
		}
		if mode&Encode == Encode {
			runner = append(runner, gen.New(gen.EncodeTest, testwr))
		}
		err = writePkgHeader(testwr, pkg)
		if err != nil {
			return err
		}
		err = writeImportHeader(testwr, "bytes", "github.com/tinylib/msgp/msgp", "testing")
	}

	for _, g := range runner {
		for _, e := range els {
			if err := g.Execute(e); err != nil {
				return err
			}
		}
	}
	return nil
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
