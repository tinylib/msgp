MessagePack Code Generator
=======

[![forthebadge](http://forthebadge.com/badges/uses-badges.svg)](http://forthebadge.com)
[![forthebadge](http://forthebadge.com/badges/ages-12.svg)](http://forthebadge.com)

This is a code generation tool and serialization library for [MesssagePack](http://msgpack.org). It is targeted at the `go generate` [tool](http://tip.golang.org/cmd/go/#hdr-Generate_Go_files_by_processing_source). You can read more about MessagePack [in the wiki](http://github.com/philhofer/msgp/wiki), or at [msgpack.org](http://msgp.org).

### Why?

- Use Go as your schema language
- Speeeeeed (400MB/s on modern hardware)
- [JSON interop](http://godoc.org/github.com/philhofer/msgp/msgp#CopyToJSON)
- [User-defined types](http://github.com/philhofer/msgp/wiki/Using-Extensions)

### Quickstart

*** Note: `go generate` is a go 1.4+ feature. ***

In a source file, include the following directive:

```go
//go:generate msgp
```

The `msgp` command will generate serialization methods for all exported struct
definitions in the file. You will need to include that directive in every file that contains structs that 
need code generation.

You can [read more about the code generation options here](http://github.com/philhofer/msgp/wiki/Using-the-Code-Generator).

### Use

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

By default, the code generator will satisfy `msgp.Sizer`, `msgp.Encodable`, `msgp.Decodable`, 
`msgp.Marshaler`, and `msgp.Unmarshaler`. Carefully-designed applications can use these methods to do
marshalling/unmarshalling with zero allocations.

While `msgp.Marshaler` and `msgp.Unmarshaler` are quite similar to the standard library's
`json.Marshaler` and `json.Unmarshaler`, `msgp.Encodable` and `msgp.Decodable` are useful for 
stream serialization. (`*msgp.Writer` and `*msgp.Reader` are essentially protocol-aware versions
of `*bufio.Writer` and `*bufio.Reader`, respectively.)

### Features

 - Extremely fast generated code
 - JSON interoperability (see `msgp.CopyToJSON() and msgp.UnmarshalAsJSON()`)
 - Support for embedded fields, anonymous structs, and multi-field inline declarations
 - Identifier resolution (see below)
 - Native support for Go's `time.Time`, `complex64`, and `complex128` types 
 - Generation of both `[]byte`-oriented and `io.Reader/io.Writer`-oriented methods
 - Support for arbitrary type system extensions
 - [Preprocessor directives](http://github.com/philhofer/msgp/wiki/Preprocessor-Directives)

Because of (limited) identifier resolution, the code generator will still yield the
correct code for the following struct declaration:
```go
const Eight = 8
type MyInt int
type Data []byte

type Struct struct {
	Which  map[string]*MyInt `msg:"which"`
	Other  Data              `msg:"other"`
	Nums   [Eight]float64    `msg:"nums"`
}
```
As long as the declarations of `MyInt` and `Data` are in the same file as `Struct`, the parser will figure out that 
`MyInt` is really an `int`, and `Data` is really just a `[]byte`. The constant `Eight` does not have to be in the 
same file as the struct definition, but it does have to be in the same package (as the generated code will simply
use `Eight` as a literal under the assumption that the compiler will figure out what it is.) Note that this only works for "base" types (no composite types, although `[]byte` is supported as a special case.) Unresolved identifiers are (optimistically) 
assumed to be struct definitions in other files. (The parser will spit out warnings about unresolved identifiers.)

#### Extensions

MessagePack supports defining your own types through "extensions," which are just a tuple of
the data "type" (`int8`) and the raw binary. You [can see a worked example in the wiki.](http://github.com/philhofer/msgp/wiki/Using-Extensions)

### Status

Alpha. I _will_ break stuff.

You can read more about how `msgp` maps MessagePack types onto Go types [in the wiki](http://github.com/philhofer/msgp/wiki).

Here some of the known limitations/restrictions:

 - All fields of a struct that are not Go built-ins are assumed (optimistically) to have been seen by the code generator in another file. The generator will output a warning if it can't resolve an identifier in the file, or if it ignores an exported field. The generated code will fail to compile if you encounter this issue, so it shouldn't catch you by surprise.
 - Like most serializers, `chan` and `func` fields are ignored, as well as non-exported fields.
 - Methods are only generated for `struct` definitions. Chances are that we will keep things this way.
 - Encoding of `interface{}` is limited to built-ins or types that have explicit encoding methods.
 - _Maps must have `string` keys._ This is intentional (as it preserves JSON interop.) Although non-string map keys are not forbidden by the MessagePack standard, many serializers impose this restriction. (It also means *any* well-formed `struct` can be de-serialized into a `map[string]interface{}`.) The only exception to this rule is that the deserializers will allow you to read map keys encoded as `bin` types, due to the fact that some legacy encodings permitted this. (However, those values will still be cast to Go `string`s, and they will be converted to `str` types when re-encoded. It is the responsibility of the user to ensure that map keys are UTF-8 safe in this case.) The same rules hold true for JSON translation.
 - All variable-length objects (maps, strings, arrays, extensions, etc.) cannot have more than `(1<<32)-1` elements.

If the output compiles, then there's a pretty good chance things are fine. (Plus, we generate tests for you.) *Please, please, please* file an issue if you think the generator is writing broken code.

### Performance

If you like benchmarks, we're the [fastest and lowest-memory-footprint round-trip serializer for Go in this test.](https://github.com/alecthomas/go_serialization_benchmarks)

Speed is the core impetus for this project; there are already a number of reflection-based libraries that exist 
for MessagePack serialization. (Type safety and generated tests are an added bonus.)