package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/tinylib/msgp/gen"
	"github.com/tinylib/msgp/parse"
	"github.com/ttacon/chalk"
	"golang.org/x/tools/imports"
)

var (
	// command line flags
	out     = flag.String("o", "", "output file")
	file    = flag.String("file", "", "input file")
	pkg     = flag.String("pkg", "", "output package")
	encode  = flag.Bool("io", true, "create Encode and Decode methods")
	marshal = flag.Bool("marshal", true, "create Marshal and Unmarshal methods")
	tests   = flag.Bool("tests", true, "create tests and benchmarks")

	// marshal/unmarshal imports
	injectImports = []string{
		"github.com/tinylib/msgp/msgp",
	}

	// testing imports
	testImport = []string{
		"testing",
		"bytes",
		"github.com/tinylib/msgp/msgp",
	}
)

func main() {
	flag.Parse()

	// GOFILE and GOPACKAGE are
	// set by `go generate`
	if *file == "" {
		*file = os.Getenv("GOFILE")
		if *file == "" {
			fmt.Println(chalk.Red.Color("No file to parse."))
			os.Exit(1)
		}
	}

	if *pkg == "" {
		*pkg = os.Getenv("GOPACKAGE")
	}

	if !(*encode) && !(*marshal) {
		fmt.Println(chalk.Red.Color("No methods to generate; -io=false && -marshal=false"))
		os.Exit(1)
	}

	// test file generation can run at
	// the same time as method generation
	if *tests && runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(2)
	}

	if err := DoAll(*pkg, *file, *marshal, *encode, *tests); err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
}

func isDirectory(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func newFilename(old string, pkg string, isDir bool) string {
	if *out != "" {
		if pre := strings.TrimPrefix(*out, old); len(pre) > 0 &&
			!strings.HasSuffix(*out, ".go") {
			return filepath.Join(old, *out)
		}
		return *out
	}
	if isDir {
		old = filepath.Join(old, pkg)
	}
	// new file name is old file name + _gen.go
	return strings.TrimSuffix(old, ".go") + "_gen.go"
}

func generate(els []gen.Elem, dst *bufio.Writer, ms ...gen.Method) error {
	gs := make([]gen.Generator, len(ms))
	for i := range gs {
		gs[i] = gen.New(ms[i], dst)
	}
	for _, el := range els {
		p, ok := el.(*gen.Ptr)
		if !ok {
			continue
		}
		for _, g := range gs {
			if err := g.Execute(p); err != nil {
				return err
			}
		}
	}
	return nil
}

func makeTests(newfile string, pkg string, elems []gen.Elem, marshal bool, encode bool) error {
	testfile := strings.TrimSuffix(newfile, ".go") + "_test.go"
	tfl, err := os.Create(testfile)
	if err != nil {
		return err
	}
	defer tfl.Close()
	testwr := bufio.NewWriter(tfl)
	err = writePkgHeader(testwr, pkg)
	if err != nil {
		return err
	}
	err = writeImportHeader(testwr, testImport...)
	if err != nil {
		return err
	}
	if encode {
		err = generate(elems, testwr, gen.EncodeTest)
		if err != nil {
			return err
		}
	}
	if marshal {
		err = generate(elems, testwr, gen.MarshalTest)
		if err != nil {
			return err
		}
	}
	err = testwr.Flush()
	if err != nil {
		return err
	}
	err = processFileImports(tfl, testfile)
	if err != nil {
		return err
	}
	fmt.Printf(chalk.Magenta.Color("TESTS  ======> %s \u2713\n"), testfile)
	return nil
}

// DoAll writes all methods using the associated file and package.
// (The package is only relevant for writing the new file's package declaration.)
func DoAll(gopkg string, gofile string, marshal bool, encode bool, tests bool) error {
	// ...nothing to do!
	if !marshal && !encode {
		return nil
	}

	fmt.Printf(chalk.Magenta.Color("========= %s =========\n"), filepath.Clean(gofile))

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

	newfile := newFilename(gofile, pkgName, isDirectory(gofile))

	// test file is generated in a separate goroutine
	var testc chan error
	if tests {
		testc = make(chan error, 1)
		go func(c chan error) {
			c <- makeTests(newfile, gopkg, elems, marshal, encode)
		}(testc)
	}

	file, err := os.Create(newfile)
	if err != nil {
		return err
	}
	defer file.Close()

	outwr := bufio.NewWriter(file)

	err = writePkgHeader(outwr, gopkg)
	if err != nil {
		return err
	}

	err = writeImportHeader(outwr, injectImports...)
	if err != nil {
		return err
	}

	if encode {
		err = generate(elems, outwr, gen.Encode, gen.Decode)
		if err != nil {
			return err
		}
	}
	if marshal {
		err = generate(elems, outwr, gen.Marshal, gen.Unmarshal, gen.Size)
		if err != nil {
			return err
		}
	}
	err = outwr.Flush()
	if err != nil {
		return err
	}
	err = processFileImports(file, newfile)
	if err != nil {
		return err
	}

	fmt.Printf(chalk.Magenta.Color("OUTPUT ======> %s \u2713\n"), newfile)

	if tests {
		err = <-testc
	}
	return err
}

func processFileImports(f *os.File, fName string) error {
	fStat, err := f.Stat()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(fName)
	if err != nil {
		return err
	}

	data, err = imports.Process(fName, data, nil)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fName, data, fStat.Mode())
	return err
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
