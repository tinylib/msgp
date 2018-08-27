// msgp is a code generation tool for
// creating methods to serialize and de-serialize
// Go data structures to and from MessagePack.
//
// This package is targeted at the `go generate` tool.
// To use it, include the following directive in a
// go source file with types requiring source generation:
//
//     //go:generate msgp
//
// The go generate tool should set the proper environment variables for
// the generator to execute without any command-line flags. However, the
// following options are supported, if you need them:
//
//  -o = output file name (default is {input}_gen.go)
//  -file = input file name (or directory; default is $GOFILE, which is set by the `go generate` command)
//  -io = satisfy the `msgp.Decodable` and `msgp.Encodable` interfaces (default is true)
//  -marshal = satisfy the `msgp.Marshaler` and `msgp.Unmarshaler` interfaces (default is true)
//  -tests = generate tests and benchmarks (default is true)
//
// For more information, please read README.md, and the wiki at github.com/tinylib/msgp
//
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
)

var (
	out        = flag.String("o", "", "output file")
	verbose    = flag.Bool("v", false, "verbose")
	file       = flag.String("file", "", "input file")
	encode     = flag.Bool("io", true, "create Encode and Decode methods")
	marshal    = flag.Bool("marshal", true, "create Marshal and Unmarshal methods")
	tests      = flag.Bool("tests", true, "create tests and benchmarks")
	unexported = flag.Bool("unexported", false, "also process unexported types")
)

func fatalf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, f, args...)
	os.Exit(1)
}

func fatalln(line string) {
	fmt.Fprintln(os.Stderr, line)
	os.Exit(1)
}

func logf(f string, args ...interface{}) {
	if *verbose {
		fmt.Fprintf(os.Stderr, f, args...)
	}
}

func logln(line string) {
	if *verbose {
		fmt.Fprintln(os.Stderr, line)
	}
}

type logger struct{}

func (l *logger) Logf(f string, args ...interface{}) {
	logf(f, args...)
}

func main() {
	flag.Parse()

	// GOFILE is set by go generate
	if *file == "" {
		*file = os.Getenv("GOFILE")
		if *file == "" {
			fatalln("msgp: no input file given")
		}
	}

	var mode gen.Method
	if *encode {
		mode |= (gen.Encode | gen.Decode | gen.Size)
	}
	if *marshal {
		mode |= (gen.Marshal | gen.Unmarshal | gen.Size)
	}
	if *tests {
		mode |= gen.Test
	}

	if mode&^gen.Test == 0 {
		fatalln("msgp: no methods to generate; -io=false && -marshal=false")
		os.Exit(1)
	}

	if err := Run(*file, mode, *unexported); err != nil {
		fatalln("msgp: " + err.Error())
	}
}

// Run writes all methods using the associated file or path, e.g.
//
//	err := msgp.Run("path/to/myfile.go", gen.Size|gen.Marshal|gen.Unmarshal|gen.Test, false)
//
func Run(gofile string, mode gen.Method, unexported bool) error {
	if mode&^gen.Test == 0 {
		logln("msgp: no code to generate?")
		return nil
	}
	logf("msgp: input: \"%s\"\n", gofile)
	fs, err := parse.File(gofile, &logger{}, unexported)
	if err != nil {
		return err
	}

	if len(fs.Identities) == 0 {
		logln("msgp: no types requiring code generation were found")
		return nil
	}

	outfile := newFilename(gofile, fs.Package)
	logf("msgp: writing %s\n", outfile)
	return printer.PrintFile(outfile, fs, mode)
}

// picks a new file name based on input flags and input filename(s).
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
