msgp
=======

This is a code generation tool for serializing Go `struct`s using the [MesssagePack](http://msgpack.org) standard. It is targeted 
at the `go generate` [tool](http://tip.golang.org/cmd/go/#hdr-Generate_Go_files_by_processing_source).

### Use

In a source file, include the following directive:

```go
//go:generate msgp
```

The `msgp` command will generate `Unmarshal`, `Marshal`, `EncodeMsg`, and `DecodeMsg` methods for all exported struct
definitions in the file. You will need to include that directive in every file that contains structs that 
need code generation. The generated files will be named {filename}_gen.go.

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
 - Efficient and readable generated code
 - JSON interoperability
 - Support for embedded fields (see Benchmark 1) and inline comma-separated declarations (see Benchmark 2)
 - Support for Go's `complex64` and `complex128` types
 - Generation of both `[]byte`-oriented and `io.Reader/io.Writer`-oriented methods.


### Status

Very alpha. Here are the known limitations:

 - The only currently supported map types are `map[string]string` and `map[string]interface{}`.
 - All fields of a struct that are not Go built-ins are assumed to satisfy the `enc.MsgEncoder` and `enc.MsgDecoder`
   interface. This will only *actually* be the case if the declaration for that type is in a file touched by the code generator.
 - Like most serializers, `chan` and `func` fields are ignored, as well as non-exported fields.

We will (soon) offer support for `map[string]T`, where `T` is any other supported type. We do *not* plan on 
supporting maps keyed on non-string fields, as it would severly limit JSON interoperability.

### Performance

As you might imagine, the generated code is quite a lot more performant than reflection-based serialization. Here 
are the two built-in benchmarking cases provided. Each benchmark writes the struct to a buffer and then extracts 
it again.


##### Benchmark 1

Type specification:
```go
type TestType struct {
	F   *float64          `msg:"float"`
	Els map[string]string `msg:"elements"`
	Obj struct {
		ValueA string `msg:"value_a"`
		ValueB []byte `msg:"value_b"`
	} `msg:"object"`
	Child *TestType `msg:"child"`
}
```

|  Method | Time | Heap Use | Heap Allocs |
|:-------:|:----:|:--------:|:-----------:|
| msgp codegen | 2553ns | 129 B | 5 allocs |
| [ugorji/go](http://github.com/ugorji/go) | 8467ns | 2015 B | 43 allocs |
| encoding/json | 11976ns | 2004 B | 33 allocs |



##### Benchmark 2

Type specification:
```go
type TestFast struct {
	Lat, Long, Alt float64
	Data []byte
}
```
|  Method | Time | Heap Use | Heap Allocs |
|:-------:|:----:|:--------:|:-----------:|
| msgp codegen | 1151ns | 57 B | 1 allocs |
| [ugorji/go](http://github.com/ugorji/go) | 5199ns | 1227 B | 27 allocs |
| encoding/json | 8120ns | 1457 B | 11 allocs |
