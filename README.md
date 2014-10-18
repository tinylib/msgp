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
need code generation. The generated files will be named {filename}_gen.go by default (but can 
be overridden with the `-o` flag.) Additionally, it will write tests and benchmarks in {{filename}}_gen_test.go.

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
func (z *Person) MarshalMsg() ([]byte, error)

func (z *Person) UnmarshalMsg(b []byte) ([]byte, error)

func (z *Person) EncodeTo(en *enc.MsgWriter) (n int, err error)

func (z *Person) DecodeFrom(dc *enc.MsgReader) (n int, err error)

func (z *Person) EncodeMsg(w io.Writer) (n int, err error)

func (z *Person) DecodeMsg(r io.Reader) (n int, err error)
```

Each method is optimized for a certain use-case, depending on whether or not the user
can afford to recycle `*enc.MsgWriter` and `*enc.MsgReader` object, and whether or not
the user is dealing with `io.Reader/io.Writer` or `[]byte`. (Pro tip: `EncodeTo` is the
fastest marshaller, and `UnmarshalMsg` is the fastest unmarshaller. Also, `*enc.MsgReader`s
are buffered, and you can create one with an arbitrary buffer size with `enc.NewDecoderSize`.)

The `msgp/enc` package has utility functions for transforming MessagePack to JSON directly,
both from `io.Reader`s and `[]byte`.


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

Alpha. Here are the known limitations/restrictions:

 - All fields of a struct that are not Go built-ins are assumed to satisfy the `enc.MsgEncoder` and `enc.MsgDecoder`
   interfaces. This will only *actually* be the case if the declaration for that type is in a file touched by the code generator. The generator will output a warning if it can't resolve an identifier in the file, or if it ignores an exported field. The generated code will fail to compile if you encounter this issue, so it shouldn't catch you by surprise.
 - Like most serializers, `chan` and `func` fields are ignored, as well as non-exported fields.
 - Methods are only generated for `struct` definitions. Chances are that we will keep things this way.
 - Encoding/decoding of `interface{}` can be flaky. It works fine for go builtins, but don't count on it working 
   well for anything beyond that.
 - Maps must have `string` keys. This is intentional (as it preserves JSON interop.) The generator will yell
   at you if you try to use something else. Although non-string map keys are not explicitly forbidden by the messagepack
   standard, many serializers impose this restriction.

I have no idea whether or not this generator will work with your code. Luckily, if the output compiles, then 
there's a pretty good chance things are fine. (Plus, we generate tests for you.) *Please, please, please* file 
an issue if you think the generator is writing broken code.

### Performance

As you might imagine, the generated code is quite a lot more performant than reflection-based serialization; generally 
you can expect a ~4x speed improvement and an order of magnitude fewer memory allocations when compared to other reflection-based messagepack serialization libraries. (Typically, the `EncodeTo` method does zero allocations. `DecodeFrom` can do as little as 1 allocation if the struct fields are already initialized.) YMMV.

If you like benchmarks, we're the [second-fastest and lowest-memory-footprint round-trip serializer for Go in this test.](https://github.com/alecthomas/go_serialization_benchmarks)
