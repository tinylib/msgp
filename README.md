msgp
=======

This is a tool for serializing Go `struct`s using the [MesssagePack](http://msgpack.org) standard. It is targeted 
at the `go generate` tool.

### Status

Very alpha.

### Performance

As you might imagine, the generated code is quite a lot more performant than reflection-based serialization.

|  Method | Time | Heap Use | Heap Allocs |
|:-------:|:----:|:--------:|:-----------:|
| msgp codegen | 2553ns | 129 B | 5 allocs |
| [ugorji/go](http://github.com/ugorji/go) | 8467ns | 2015 B | 43 allocs |
| encoding/json | 11976ns | 2004 B | 33 allocs |