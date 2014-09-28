msgp
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

The `msgp` command will generate `Unmarshal`, `Marshal`, `WriteTo`, and `ReadFrom` methods for all exported struct
definitions in the file. You will need to include that directive in every file that contains structs that 
need code generation. The generated files will be named {filename}_gen.go by default (but can 
be overridden with the `-o` flag.)

Field names can be overridden in much the same way as the `encoding/json` package. For example:

```go
type Person struct {
	Name       string `msg:"name"`
	Address    string `msg:"address"`
	Age        int    `msg:"age"`
	Hidden     string `msg:"-"` // this field is ignored
	unexported bool             // this field is also ignored
}
```

Running `msgp` on a file containing the above declaration would generate the following methods:

```go
func (z *Person) Marshal() ([]byte, error)

func (z *Person) Unmarshal(b []byte) error

func (z *Person) EncodeMsg(w io.Writer) (n int, err error)

func (z *Person) DecodeMsg(r io.Reader) (n int, err error)
```

The `msgp/enc` package has a function called `CopyToJSON` which can take MessagePack-encoded binary
and translate it directly into JSON. It has reasonably high performance, and works much the same way that `io.Copy` does.


### Features

 - Proper handling of pointers, nullables, and recursive types
 - Backwards-compatible decoding logic
 - Efficient and readable generated code
 - JSON interoperability
 - Support for embedded fields, anonymous structs, and multi-field inline declarations
 - Identifier resolution (see below)
 - Support for Go's `time.Time`, `complex64`, and `complex128` types through extensions
 - Generation of both `[]byte`-oriented and `io.Reader/io.Writer`-oriented methods.

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


### Status

Very alpha. Here are the known limitations:

 - All fields of a struct that are not Go built-ins are assumed to satisfy the `enc.MsgEncoder` and `enc.MsgDecoder`
   interfaces. This will only *actually* be the case if the declaration for that type is in a file touched by the code generator.
   The generator will output a warning if it can't resolve an identifier in the file, or if it ignores an exported field.
 - Like most serializers, `chan` and `func` fields are ignored, as well as non-exported fields.
 - Methods are only generated for `struct` definitions.
 - *Encoding/decoding of `interface{}` and doesn't work yet*.
 - Static arrays are not yet supported.

There may be other problems. I want this code to be in beta by the go 1.4 release.

### Performance

As you might imagine, the generated code is quite a lot more performant than reflection-based serialization; generally 
you can expect a ~4x speed improvement and an order of magnitude fewer memory allocations when compared to other reflection-based messagepack serialization libraries. YMMV.

If you like benchmarks, we're the [second-fastest and lowest-memory-footprint round-trip serializer for Go in this test.](https://github.com/alecthomas/go_serialization_benchmarks)
