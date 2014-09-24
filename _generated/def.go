package _generated

//go:generate msgp

// All of the struct
// definitions in this
// file are fed to the code
// generator when `make test` is
// called, followed by an
// invocation of `go test -v` in this
// directory. A simple way of testing
// a struct definition is
// by adding it to this file.

type TestType struct {
	F   *float64          `msg:"float"`
	Els map[string]string `msg:"elements"`
	Obj struct {          // test anonymous struct
		ValueA string `msg:"value_a"`
		ValueB []byte `msg:"value_b"`
	} `msg:"object"`
	Child *TestType `msg:"child"`
}

type TestFast struct {
	Lat, Long, Alt float64 // test inline decl
	Data           []byte
}

type TestHidden struct {
	A   string
	B   []float64
	Bad func(string) bool // This results in a warning: field "Bad" unsupported
}

type Embedded struct {
	*Embedded // test embedded field
	Other     string
}

type Things struct {
	Cmplx []complex64 `msg:"complexes"` // test slices
	Vals  []int32     `msg:"values"`
}

type Empty struct{}

type CustomInt int
type CustomBytes []byte

type Custom struct {
	Int []CustomInt
	Bts CustomBytes
	Mp  map[string]*Embedded
}
