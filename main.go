package main

import (
	"bufio"
	"fmt"
	"github.com/philhofer/msgp/gen"
	"github.com/philhofer/msgp/parse"
	"io"
	"os"
	"strings"
)

var (
	// exports from go:generate
	GOFILE string
	GOPKG  string
	GOPATH string
)

func init() {
	GOFILE = os.Getenv("GOFILE")
	GOPKG = os.Getenv("GOPKG")
	GOPATH = os.Getenv("GOPATH")
}

func main() {
	// location of the file to pase
	loc := GOPATH + "/src/" + GOPKG + "/" + GOFILE

	elems, err := parse.GetElems(loc)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// new file name is old file name + _gen.go
	newfile := strings.TrimSuffix(loc, ".go") + "_gen.go"

	file, err := os.Create(newfile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	wr := bufio.NewWriter(file)

	err = writePkgHeader(wr, GOPKG)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = writeImportHeader(wr, "github.com/philhofer/msgp/enc", "io")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, el := range elems {
		p, ok := el.(*gen.Ptr)
		if !ok {
			continue
		}
		err = gen.WriteEncoderMethod(wr, p)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		err = gen.WriteDecoderMethod(wr, p)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	err = wr.Flush()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return
}

func writePkgHeader(w io.Writer, name string) error {
	_, err := io.WriteString(w, fmt.Sprintf("package %s\n\n", name))
	return err
}

func writeImportHeader(w io.Writer, imports ...string) error {
	start := "import(\n"
	for _, im := range imports {
		start += fmt.Sprintf("\t%q\n", im)
	}
	start += "\n)\n\n"
	_, err := io.WriteString(w, start)
	return err
}
