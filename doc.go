// msgp is a code generation tool for
// creating methods to serialize and de-serialize
// Go data structures to and from MessagePack.
//
// This package is targeted at the `go generate` tool.
// To use it, include the following directive in a
// go source file:
//
//     //go:generate msgp
//
// The go generate tool should set the proper environment variables for
// the generator to execute without any command-line flags. However, the
// following options are supported, if you need them:
//
//  -o = output file name (default is {filename}_gen.go)
//  -file = input file name (default is $GOPATH/src/$GOPACKAGE/$GOFILE, which are set by the `go generate` command)
//  -pkg = output package name (default is $GOPACKAGE)
//  -io = satisfy the `msgp.Decoder` and `msgp.Encoder` interfaces (default is true)
//  -marshal = satisfy the `msgp.Marshaler` and `msgp.Unmarshaler` interfaces (default is true)
//  -tests = generate tests and benchmarks (default is true)
//
// For more information, please read README.md
//
package main
