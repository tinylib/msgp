package main

import (
	"bufio"
	"fmt"
	"github.com/philhofer/msgp/gen"
	"github.com/philhofer/msgp/parse"
	"github.com/ttacon/chalk"
	"io"
	"os"
	"strings"
)

var (
	// exports from go:generate
	GOFILE string
	GOPKG  string
	GOPATH string

	// these are the required imports
	injectImports []string = []string{
		"github.com/philhofer/msgp/enc",
		"io",
		"bytes",
	}
)

func init() {
	GOFILE = os.Getenv("GOFILE")
	GOPKG = os.Getenv("GOPACKAGE")
}

func main() {
	err := DoAll(GOPKG, GOFILE)
	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
	fmt.Println(chalk.Magenta.Color("MSGP: Done."))
}

// DoAll writes all methods using the associated GOPATH, GOPACKAGE,
// and GOFILE variables.
func DoAll(gopkg string, gofile string) error {
	// location of the file to pase

	fmt.Printf(chalk.Magenta.Color("MSGP: using %s/%s\n"), gopkg, gofile)
	elems, err := parse.GetElems(gofile)
	if err != nil {
		return err
	}

	// new file name is old file name + _gen.go
	newfile := strings.TrimSuffix(gofile, ".go") + "_gen.go"

	fmt.Printf(chalk.Magenta.Color("MSGP: Writing file %s/%s\n"), gopkg, newfile)
	file, err := os.Create(newfile)
	if err != nil {
		return err
	}
	defer file.Close()

	wr := bufio.NewWriter(file)

	err = writePkgHeader(wr, gopkg)
	if err != nil {
		return err
	}

	err = writeImports(wr)
	if err != nil {
		return err
	}

	for _, el := range elems {
		p, ok := el.(*gen.Ptr)
		if !ok {
			continue
		}
		// propogate names to
		// child elements of struct
		gen.Propogate(p, "z")

		err = gen.WriteMarshalMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		err = gen.WriteUnmarshalMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		err = gen.WriteEncoderMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		err = gen.WriteDecoderMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}
	}
	err = wr.Flush()
	return err
}

func writePkgHeader(w io.Writer, name string) error {
	_, err := io.WriteString(w, fmt.Sprintf("package %s\n\n", name))
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, "// NOTE: THIS FILE WAS PRODUCED BY THE\n// MSGP CODE GENERATION TOOL (github.com/philhofer/msgp)\n// DO NOT EDIT\n\n")

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

func writeImports(w io.Writer) error {
	return writeImportHeader(w, injectImports...)
}
