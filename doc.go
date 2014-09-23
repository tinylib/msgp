// msgp is a code generation tool for
// creating methods to serialize and de-serialize
// Go data structures. Currently the only supported
// encoding is MessagePack (http://msgpack.org).
//
// This package is targeted at the `go generate` tool.
// To use it, include the following directive in a
// go source file:
//
//     //go:generate msgp
//
// The tool will read all of the struct type declarations
// in the file and create the following methods:
//
//     (z *T) Unmarshal(b []byte) error
//
//     (z *T) Marshal() ([]byte, error)
//
//     (z *T) WriteTo(w io.Writer) (int64, error)
//
//     (z *T) ReadFrom(r io.Reader) (int64, error)
//
// There are currently three supported command line flags:
//
//     -o        change the name of the output file (by default, {filename}_gen.go)
//     -file     change the name of the input file (by default, the file with the go:generate directive)
//     -pkg      change the name of the package header to write (by default, the same as the file with the go:generate directive)
//
//
//
package main
