MessagePack Code Generator
=======

[![forthebadge](http://forthebadge.com/badges/uses-badges.svg)](http://forthebadge.com)
[![forthebadge](http://forthebadge.com/badges/certified-snoop-lion.svg)](http://forthebadge.com)

This is a code generation tool for serializing Go `struct`s using the [MesssagePack](http://msgpack.org) standard. It is targeted 
at the `go generate` [tool](http://tip.golang.org/cmd/go/#hdr-Generate_Go_files_by_processing_source).

### Use

*** Note: `go generate` is a go 1.4+ feature. ***

In a source file, include the following directive:

```go
//go:generate msgp
```

The `msgp` command will generate serialization methods for all exported struct
definitions in the file. You will need to include that directive in every file that contains structs that 
need code generation.

#### Options

The following command-line flags are supported:

 - `-o` - output file name (default is `{filename}_gen.go`)
 - `-file` - input file name (default is `$GOPATH/src/$GOPACKAGE/$GOFILE`, which are set by the `go generate` command)
 - `-pkg` - output package name (default is `$GOPACKAGE`)
 - `-encode` - satisfy the `msgp.Decoder` and `msgp.Encoder` interfaces (default is `true`)
 - `-marshal` - satisfy the `msgp.Marshaler` and `msgp.Unmarshaler` interfaces (default is `true`)
 - `-tests` - generate tests and benchmarks (default is `true`)

#### Use

Field names can be set in much the same way as the `encoding/json` package. For example:

```go
type Person struct {
	Name       string `msg:"name"`
	Address    string `msg:"address"`
	Age        int    `msg:"age"`
	Hidden     string `msg:"-"` // this field is ignored
	unexported bool             // this field is also ignored
}
```

Running `msgp` on a file containing the above declaration would generate methods 
to fulfill the following interfaces (which are defined in the `/msgp` subdirectory):

```go
// Marshaler is the interface implemented
// by types that know how to marshal themselves
// as MessagePack
type Marshaler interface {
	// AppendMsg appends the marshalled
	// form of the object to the provided
	// byte slice, returning the extended
	// slice and any errors encountered
	AppendMsg([]byte) ([]byte, error)

	// MarshalMsg returns a new []byte
	// containing the MessagePack-encoded
	// form of the object
	MarshalMsg() ([]byte, error)
}

// Encoder is the interface implemented
// by types that know how to write themselves
// as MessagePack
type Encoder interface {
	// EncodeMsg writes the object in MessagePack
	// format to the writer, returning the number
	// of bytes written and any errors encountered.
	EncodeMsg(io.Writer) (int, error)

	// EncodeTo writes the object to a *Writer,
	// returning the number of bytes written, and
	// any errors encountered.
	EncodeTo(*Writer) (int, error)
}

// Unmarshaler is the interface fulfilled
// byte objects that know how to unmarshal
// themselves from MessagePack
type Unmarshaler interface {
	// UnmarshalMsg unmarshals the object
	// from binary, returing any leftover
	// bytes and any errors encountered
	UnmarshalMsg([]byte) ([]byte, error)
}

// Decoder is the interface fulfilled
// by objects that know how to decode themselves
// from MessagePack
type Decoder interface {
	// DecodeMsg decodes the object
	// from an io.Reader, returning
	// the number of bytes read and
	// any errors encountered
	DecodeMsg(io.Reader) (int, error)

	// DecodeFrom decodes the object
	// using an existing *Reader,
	// returning the number of bytes
	// read and any errors encountered
	DecodeFrom(*Reader) (int, error)
}
```

Each method is optimized for a certain use-case, depending on whether or not the user
can afford to recycle `*msgp.Writer` and `*msgp.Reader` object, and whether or not
the user is dealing with `io.Reader/io.Writer` or `[]byte`. Carefully-designed applications 
can use these methods to do marshalling/unmarshalling with zero allocations.

The `/msgp` package has utility functions for transforming MessagePack to JSON directly,
both from `io.Reader`s and `[]byte`. There are also utilities for in-place manipulation of
raw MessagePack.

### Features

 - Extremely fast generated code
 - JSON interoperability (see `msgp.CopyToJSON() and msgp.UnmarshalAsJSON()`)
 - Support for embedded fields, anonymous structs, and multi-field inline declarations
 - Identifier resolution (see below)
 - Native support for Go's `time.Time`, `complex64`, and `complex128` types 
 - Generation of both `[]byte`-oriented and `io.Reader/io.Writer`-oriented methods
 - Support for arbitrary type system extensions (see below)

For example, this will work fine:
```go
type MyInt int
type Data []byte
type Struct struct {
	Which  map[string]*MyInt `msg:"which"`
	Other  Data              `msg:"other"`
}
```
As long as the declarations of `MyInt` and `Data` are in the same file, the parser will figure out that 
`MyInt` is really an `int`, and `Data` is really just a `[]byte`. Note that this only works for "base" types 
(no slices or maps, although `[]byte` is supported as a special case.) Unresolved identifiers are (optimistically) 
assumed to be struct definitions in other files. (The parser will spit out warnings about unresolved identifiers.)

#### Extensions

MessagePack supports defining your own types through "extensions," which are just a tuple of
the data "type" (`int8`) and the raw binary. Extensions should satisfy the `msgp.Extension` interface:

```go
// Extension is the interface fulfilled
// by types that want to define their
// own binary encoding.
type Extension interface {
	// ExtensionType should return
	// a int8 that identifies the concrete
	// type of the extension. (Types <0 are
	// officially reserved by the MessagePack
	// specifications.)
	ExtensionType() int8

	// Len should return the length
	// of the data to be encoded
	Len() int

	// MarshalBinaryTo should copy
	// the data into the supplied slice,
	// assuming that the slice has length Len()
	MarshalBinaryTo([]byte) error

	UnmarshalBinary([]byte) error
}
```
A simple implementation of `Extension` is the `msgp.RawExtension` type.

In order to tell the generator to use the field as an extension, you must include the "extension"
annotation in the struct tag:

```go
type Thing struct {
	Blah msgp.RawExtension `msg:"blah,extension"`
}
```

Users should use `msgp.RegisterExtension` to register type system extensions. Registering extensions 
allows the `interface{}` decoders and the JSON translators to use the concrete type of your extension 
when decoding/encoding.

### Status

Alpha. I _will_ break stuff.

Here are the known limitations/restrictions:

 - All fields of a struct that are not Go built-ins are assumed (optimistically) to have been seen by the code generator in another file. The generator will output a warning if it can't resolve an identifier in the file, or if it ignores an exported field. The generated code will fail to compile if you encounter this issue, so it shouldn't catch you by surprise.
 - Like most serializers, `chan` and `func` fields are ignored, as well as non-exported fields.
 - Methods are only generated for `struct` definitions. Chances are that we will keep things this way.
 - Encoding of `interface{}` is limited to built-ins or types that have explicit encoding methods.
 - _Maps must have `string` keys._ This is intentional (as it preserves JSON interop.) Although non-string map keys are not explicitly forbidden by the MessagePack standard, many serializers impose this restriction. (This restriction also means *any* well-formed `struct` can be de-serialized into a `map[string]interface{}`.)
 - All variable-length objects (maps, strings, arrays, extensions, etc.) cannot have more than `(1<<32)-1` elements.

If the output compiles, then there's a pretty good chance things are fine. (Plus, we generate tests for you.) *Please, please, please* file an issue if you think the generator is writing broken code.

### Performance

If you like benchmarks, we're the [fastest and lowest-memory-footprint round-trip serializer for Go in this test.](https://github.com/alecthomas/go_serialization_benchmarks)

Speed is the core impetus for this project; there are already a number of reflection-based libraries that exist 
for MessagePack serialization. (Type safety and generated tests are an added bonus.)