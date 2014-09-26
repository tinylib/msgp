package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/philhofer/msgp/gen"
	"github.com/philhofer/msgp/parse"
	"github.com/ttacon/chalk"
	"io"
	"os"
	"strings"
)

var (
	// command line flags
	out  string
	file string
	pkg  string
	test bool

	// these are the required imports
	injectImports []string = []string{
		"github.com/philhofer/msgp/enc",
		"io",
		"bytes",
	}
)

func init() {
	flag.StringVar(&out, "o", "", "output file")
	flag.StringVar(&file, "file", "", "input file")
	flag.StringVar(&pkg, "pkg", "", "output package")
	flag.BoolVar(&test, "test", false, "generate test and bench files")
}

func main() {
	flag.Parse()

	if file == "" {
		file = os.Getenv("GOFILE")
	}
	if pkg == "" {
		pkg = os.Getenv("GOPACKAGE")
	}

	if file == "" {
		fmt.Println(chalk.Red.Color("No file to parse."))
		os.Exit(1)
	}
	if pkg == "" {
		fmt.Println(chalk.Red.Color("No package specified."))
		os.Exit(1)
	}

	err := DoAll(pkg, file)
	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
}

// DoAll writes all methods using the associated file and package.
// (The package is only relevant for writing the new file's package declaration.)
func DoAll(gopkg string, gofile string) error {
	// location of the file to pase

	fmt.Printf(chalk.Magenta.Color("=========%12s      =========\n"), gofile)
	elems, err := parse.GetElems(gofile)
	if err != nil {
		return err
	}

	if len(elems) == 0 {
		fmt.Println(chalk.Magenta.Color("No structs requiring code generation were found..."))
		return nil
	}

	// new file name is old file name + _gen.go
	var newfile string
	if out != "" {
		newfile = out
	} else {
		newfile = strings.TrimSuffix(gofile, ".go") + "_gen.go"
	}
	// GENERATED FILES
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

	// TESTING FILES
	var twr *bufio.Writer
	var testfile string
	if test {
		testfile = strings.TrimSuffix(newfile, ".go") + "_test.go"
		tfl, err := os.Create(testfile)
		if err != nil {
			return err
		}
		defer tfl.Close()
		twr = bufio.NewWriter(tfl)
		err = writePkgHeader(twr, gopkg)
		if err != nil {
			return err
		}
		err = writeImportHeader(twr, "testing", "bytes", "github.com/philhofer/msgp/enc")
		if err != nil {
			return err
		}
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

		if test {
			err = gen.WriteTestNBench(twr, p.Value.Struct())
			if err != nil {
				wr.Flush()
				return err
			}
		}
	}

	fmt.Printf(chalk.Magenta.Color("OUTPUT ======> %s/%s"), gopkg, newfile)
	err = wr.Flush()
	if err != nil {
		return err
	}
	fmt.Print(chalk.Green.Color(" \u2713\n"))
	if test {
		fmt.Printf(chalk.Magenta.Color("TESTS =====> %s/%s"), gopkg, testfile)
		err = twr.Flush()
		if err != nil {
			return err
		}
	}
	fmt.Print(chalk.Green.Color(" \u2713\n"))
	return nil
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
