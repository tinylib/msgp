package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/philhofer/msgp/gen"
	"github.com/philhofer/msgp/parse"
	"github.com/ttacon/chalk"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	// command line flags
	out     string // output file
	file    string // input file (or directory)
	pkg     string // output package name
	encode  bool   // write io.Writer/io.Reader-based methods
	marshal bool   // write []byte-based methods
	tests   bool   // write test file

	// marshal/unmarshal imports
	injectImports = []string{
		"github.com/philhofer/msgp/msgp",
	}

	// testing imports
	testImport = []string{
		"testing",
		"bytes",
		"github.com/philhofer/msgp/msgp",
	}
)

func init() {
	flag.StringVar(&out, "o", "", "output file")
	flag.StringVar(&file, "file", "", "input file")
	flag.StringVar(&pkg, "pkg", "", "output package")
	flag.BoolVar(&encode, "io", true, "create Encode and Decode methods")
	flag.BoolVar(&marshal, "marshal", true, "create Marshal and Unmarshal methods")
	flag.BoolVar(&tests, "tests", true, "create tests and benchmarks")
}

func main() {
	flag.Parse()

	// GOFILE and GOPACKAGE are
	// set by `go generate`
	if file == "" {
		file = os.Getenv("GOFILE")
		if file == "" {
			fmt.Println(chalk.Red.Color("No file to parse."))
			os.Exit(1)
		}
	}
	if pkg == "" {
		pkg = os.Getenv("GOPACKAGE")
	}

	if !encode && !marshal {
		fmt.Println(chalk.Red.Color("No methods to generate; -io=false AND -marshal=false"))
		os.Exit(1)
	}

	if err := DoAll(pkg, file, marshal, encode, tests); err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
}

func isDirectory(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func newFilename(old string, pkg string, isDir bool) string {
	if out != "" {
		if pre := strings.TrimPrefix(out, old); len(pre) > 0 &&
			!strings.HasSuffix(out, ".go") {
			return filepath.Join(old, out)
		}
		return out
	}
	if isDir {
		old = filepath.Join(old, pkg)
	}
	// new file name is old file name + _gen.go
	return strings.TrimSuffix(old, ".go") + "_gen.go"
}

func writeMarshal(els []gen.Elem, dst *bufio.Writer, tst *bufio.Writer, buf *bytes.Buffer) error {
	for _, el := range els {
		p, ok := el.(*gen.Ptr)
		if !ok || p.Value.Type() != gen.StructType {
			continue
		}
		err := gen.WriteMarshalUnmarshal(dst, p, buf)
		if err != nil {
			return err
		}
		if tst != nil {
			err = gen.WriteMarshalUnmarshalTests(tst, p.Value.Struct(), buf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func writeEncode(els []gen.Elem, dst *bufio.Writer, tst *bufio.Writer, buf *bytes.Buffer) error {
	for _, el := range els {
		p, ok := el.(*gen.Ptr)
		if !ok || p.Value.Type() != gen.StructType {
			continue
		}
		err := gen.WriteEncodeDecode(dst, p, buf)
		if err != nil {
			return err
		}
		if tst != nil {
			err = gen.WriteEncodeDecodeTests(tst, p.Value.Struct(), buf)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// DoAll writes all methods using the associated file and package.
// (The package is only relevant for writing the new file's package declaration.)
func DoAll(gopkg string, gofile string, marshal bool, encode bool, tests bool) error {
	var (
		testwr *bufio.Writer // location to write tests, if applicable
		outwr  *bufio.Writer // location to write methods
	)

	// ...nothing to do!
	if !marshal && !encode {
		return nil
	}

	isDir := isDirectory(gofile)

	if isDir {
		fmt.Printf(chalk.Magenta.Color("========= %s =========\n"), filepath.Clean(gofile))
	} else {
		fmt.Printf(chalk.Magenta.Color("========= %s =========\n"), gofile)
	}

	elems, pkgName, err := parse.GetElems(gofile)
	if err != nil {
		return err
	}

	// use the parsed
	// package name if it
	// isn't set from $GOPACKAGE
	// or -pkg
	if gopkg == "" {
		gopkg = pkgName
	}

	// no need to continue if
	// we don't need to generate anything
	if len(elems) == 0 {
		fmt.Println(chalk.Magenta.Color("No structs requiring code generation were found..."))
		return nil
	}

	newfile := newFilename(gofile, pkgName, isDir)

	// GENERATED FILES

	//////////////////
	/// MAIN FILE ////
	file, err := os.Create(newfile)
	if err != nil {
		return err
	}
	defer file.Close()

	outwr = bufio.NewWriter(file)

	err = writePkgHeader(outwr, gopkg)
	if err != nil {
		return err
	}

	err = writeImportHeader(outwr, injectImports...)
	if err != nil {
		return err
	}

	///////////////////

	///////////////////
	// TESTING FILE  //
	var testfile string
	if tests {
		testfile = strings.TrimSuffix(newfile, ".go") + "_test.go"
		tfl, err := os.Create(testfile)
		if err != nil {
			return err
		}
		defer tfl.Close()
		testwr = bufio.NewWriter(tfl)
		err = writePkgHeader(testwr, gopkg)
		if err != nil {
			return err
		}
		err = writeImportHeader(testwr, testImport...)
		if err != nil {
			return err
		}
	}

	var buf bytes.Buffer
	if encode {
		if err := writeEncode(elems, outwr, testwr, &buf); err != nil {
			return err
		}
	}
	if marshal {
		if err := writeMarshal(elems, outwr, testwr, &buf); err != nil {
			return err
		}
	}

	fmt.Printf(chalk.Magenta.Color("OUTPUT ======> %s "), newfile)
	err = outwr.Flush()
	if err != nil {
		return err
	}
	fmt.Print(chalk.Green.Color("\u2713\n"))
	if tests {
		fmt.Printf(chalk.Magenta.Color("TESTS =====> %s "), testfile)
		err = testwr.Flush()
		if err != nil {
			return err
		}
		fmt.Print(chalk.Green.Color("\u2713\n"))
	}
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
