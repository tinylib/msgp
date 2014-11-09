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
	injectImports []string = []string{
		"github.com/philhofer/msgp/msgp",
	}

	// testing imports
	testImport []string = []string{
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
	}
	if pkg == "" {
		pkg = os.Getenv("GOPACKAGE")
	}

	if file == "" {
		fmt.Println(chalk.Red.Color("No file to parse."))
		os.Exit(1)
	}

	if !encode && !marshal {
		fmt.Println(chalk.Red.Color("No methods to generate; -io=false AND -marshal=false"))
		os.Exit(1)
	}

	err := DoAll(pkg, file, marshal, encode, tests)
	if err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
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

	// new file name is old file name + _gen.go
	var isDir bool
	if fInfo, err := os.Stat(gofile); err == nil && fInfo.IsDir() {
		isDir = true
	}

	if isDir {
		fmt.Printf(chalk.Magenta.Color("========= %s =========\n"), filepath.Clean(gofile))
	} else {
		fmt.Printf(chalk.Magenta.Color("========= %s =========\n"), gofile)
	}

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

	var newfile string
	if out != "" {
		newfile = out
		if pre := strings.TrimPrefix(out, gofile); len(pre) > 0 &&
			!strings.HasSuffix(out, ".go") {
			newfile = filepath.Join(gofile, out)
		}
	} else {
		// small sanity check if gofile == . or dir
		// let's just stat it again, not too costly
		if isDir {
			gofile = filepath.Join(gofile, pkgName)
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

	//////////////////
	var buf bytes.Buffer
	for _, el := range elems {
		p, ok := el.(*gen.Ptr)
		if !ok || p.Value.Type() != gen.StructType {
			continue
		}
		// propogate names to
		// child elements of struct
		p.SetVarname("z")

		if marshal {
			// write MarshalMsg()
			err = gen.WriteMarshalUnmarshal(outwr, p, &buf)
			if err != nil {
				outwr.Flush()
				return err
			}

			if tests {
				err = gen.WriteMarshalUnmarshalTests(testwr, p.Value.Struct(), &buf)
				if err != nil {
					testwr.Flush()
					return err
				}
			}
		}

		if encode {
			err = gen.WriteEncodeDecode(outwr, p, &buf)
			if err != nil {
				outwr.Flush()
				return err
			}

			if tests {
				err = gen.WriteEncodeDecodeTests(testwr, p.Value.Struct(), &buf)
				if err != nil {
					testwr.Flush()
					return err
				}
			}
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
