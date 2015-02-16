package _generated

import (
	"github.com/tinylib/msgp/msgp"
	"time"
)

//go:generate msgp -o generated.go

// All of the struct
// definitions in this
// file are fed to the code
// generator when `make test` is
// called, followed by an
// invocation of `go test -v` in this
// directory. A simple way of testing
// a struct definition is
// by adding it to this file.

// tests edge-cases with
// compiling size compilation.
type X struct {
	Values [32]byte    // should compile to 32*msgp.ByteSize; encoded as Bin
	Others [][32]int32 // should compile to len(x.Others)*32*msgp.Int32Size
	Matrix [][]int32   // should not optimize
}

type TestType struct {
	F   *float64          `msg:"float"`
	Els map[string]string `msg:"elements"`
	Obj struct {          // test anonymous struct
		ValueA string `msg:"value_a"`
		ValueB []byte `msg:"value_b"`
	} `msg:"object"`
	Child    *TestType   `msg:"child"`
	Time     time.Time   `msg:"time"`
	Any      interface{} `msg:"any"`
	Appended msgp.Raw    `msg:"appended"`
}

//msgp:tuple TestBench

type TestBench struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

//msgp:tuple TestFast

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
	*Embedded   // test embedded field
	Children    []Embedded
	PtrChildren []*Embedded
	Other       string
}

const eight = 8

type Things struct {
	Cmplx complex64                         `msg:"complex"` // test slices
	Vals  []int32                           `msg:"values"`
	Arr   [msgp.ExtensionPrefixSize]float64 `msg:"arr"`            // test const array and *ast.SelectorExpr as array size
	Arr2  [4]float64                        `msg:"arr2"`           // test basic lit array
	Ext   *msgp.RawExtension                `msg:"ext,extension"`  // test extension
	Oext  msgp.RawExtension                 `msg:"oext,extension"` // test extension reference
}

// test ignore directive (below)

//msgp:ignore Empty

type Empty struct{}

type MyEnum byte

const (
	A MyEnum = iota
	B
	C
	D
	invalid
)

// test shim directive (below)

//msgp:shim MyEnum as:string using:(MyEnum).String/myenumStr

func (m MyEnum) String() string {
	switch m {
	case A:
		return "A"
	case B:
		return "B"
	case C:
		return "C"
	case D:
		return "D"
	default:
		return "<invalid>"
	}
}

func myenumStr(s string) MyEnum {
	switch s {
	case "A":
		return A
	case "B":
		return B
	case "C":
		return C
	case "D":
		return D
	default:
		return invalid
	}
}

type Custom struct {
	Int   map[string]CustomInt `msg:"mapstrint"`
	Bts   CustomBytes          `msg:"bts"`
	Mp    map[string]*Embedded `msg:"mp"`
	Enums []MyEnum             `msg:"enums"` // test explicit enum shim
}
type CustomInt int
type CustomBytes []byte
