package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tinylib/msgp/gen"
	"github.com/tinylib/msgp/parse"
	"github.com/tinylib/msgp/printer"
	"github.com/ttacon/chalk"
)

var (
	out     = flag.String("o", "", "output file")
	file    = flag.String("file", "", "input file")
	encode  = flag.Bool("io", true, "create Encode and Decode methods")
	marshal = flag.Bool("marshal", true, "create Marshal and Unmarshal methods")
	tests   = flag.Bool("tests", true, "create tests and benchmarks")
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

	var mode gen.Method
	if *encode {
		mode |= (gen.Encode | gen.Decode | gen.Size)
	}
	if *marshal {
		mode |= (gen.Marshal | gen.Unmarshal)
	}
	if *tests {
		mode |= gen.Test
	}

	if mode&^gen.Test == 0 {
		fmt.Println(chalk.Red.Color("No methods to generate; -io=false && -marshal=false"))
		os.Exit(1)
	}

	if err := Run(*file, mode); err != nil {
		fmt.Println(chalk.Red.Color(err.Error()))
		os.Exit(1)
	}
}

// Run writes all methods using the associated file/path and package.
// (The package is only relevant for writing the new file's package declaration.)
func Run(gofile string, mode gen.Method) error {
	if mode&^gen.Test == 0 {
		return nil
	}
	fmt.Println(chalk.Magenta.Color("======== MessagePack Code Generator ======="))
	fmt.Printf(chalk.Magenta.Color(">>> Input: \"%s\"...\n"), gofile)
	fs, err := parse.File(gofile)
	if err != nil {
		return err
	}

	if len(fs.Identities) == 0 {
		fmt.Println(chalk.Magenta.Color("No types requiring code generation were found!"))
		return nil
	}

	return printer.PrintFile(newFilename(gofile, fs.Package), fs, mode)
}

func newFilename(old string, pkg string) string {
	if *out != "" {
		if pre := strings.TrimPrefix(*out, old); len(pre) > 0 &&
			!strings.HasSuffix(*out, ".go") {
			return filepath.Join(old, *out)
		}
		return *out
	}

	if fi, err := os.Stat(old); err == nil && fi.IsDir() {
		old = filepath.Join(old, pkg)
	}
	// new file name is old file name + _gen.go
	return strings.TrimSuffix(old, ".go") + "_gen.go"
}
