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

	// marshal/unmarshal imports
	injectImports []string = []string{
		"github.com/philhofer/msgp/enc",
		"io",
		"bytes",
	}

	// testing imports
	testImport []string = []string{
		"testing",
		"bytes",
		"github.com/philhofer/msgp/enc",
	}
)

func init() {
	flag.StringVar(&out, "o", "", "output file")
	flag.StringVar(&file, "file", "", "input file")
	flag.StringVar(&pkg, "pkg", "", "output package")
}

func main() {
	flag.Parse()

	// GOFILE and GOPACKAGE are
	// set by `go generate`
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
	elems, pkgName, err := parse.GetElems(gofile)
	if err != nil {
		return err
	}
	if len(gopkg) == 0 {
		gopkg = pkgName
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
		// small sanity check if gofile == . or dir
		// let's just stat it again, not too costly
		if fInfo, err := os.Stat(gofile); err == nil && fInfo.IsDir() {
			gofile = pkgName
		}
		newfile = strings.TrimSuffix(gofile, ".go") + "_gen.go"
	}

	// GENERATED FILES

	//////////////////
	/// MAIN FILE ////
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

	err = writeImportHeader(wr, injectImports...)
	if err != nil {
		return err
	}
	///////////////////

	///////////////////
	// TESTING FILE  //
	testfile := strings.TrimSuffix(newfile, ".go") + "_test.go"
	tfl, err := os.Create(testfile)
	if err != nil {
		return err
	}
	defer tfl.Close()
	twr := bufio.NewWriter(tfl)
	err = writePkgHeader(twr, gopkg)
	if err != nil {
		return err
	}
	err = writeImportHeader(twr, testImport...)
	if err != nil {
		return err
	}
	//////////////////

	for _, el := range elems {
		p, ok := el.(*gen.Ptr)
		if !ok {
			continue
		}
		// propogate names to
		// child elements of struct
		gen.Propogate(p, "z")

		// write Marshal()
		err = gen.WriteMarshalMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		// write Unmarshal()
		err = gen.WriteUnmarshalMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		// write WriteTo()
		err = gen.WriteEncoderMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		// write ReadFrom()
		err = gen.WriteDecoderMethod(wr, p)
		if err != nil {
			wr.Flush()
			return err
		}

		// write Test{{Type}}, BenchEncode{{Type}}, BenchDecode{{Type}}
		err = gen.WriteTestNBench(twr, p.Value.Struct())
		if err != nil {
			wr.Flush()
			return err
		}
	}

	fmt.Printf(chalk.Magenta.Color("OUTPUT ======> %s/%s "), gopkg, newfile)
	err = wr.Flush()
	if err != nil {
		return err
	}
	fmt.Print(chalk.Green.Color("\u2713\n"))
	fmt.Printf(chalk.Magenta.Color("TESTS =====> %s/%s "), gopkg, testfile)
	err = twr.Flush()
	if err != nil {
		return err
	}
	fmt.Print(chalk.Green.Color("\u2713\n"))
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
