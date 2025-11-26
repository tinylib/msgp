package _generated

import (
	"encoding/json"
	"os"
	"time"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp -v -o generated.go

// All of the struct
// definitions in this
// file are fed to the code
// generator when `make test` is
// called, followed by an
// invocation of `go test -v` in this
// directory. A simple way of testing
// a struct definition is
// by adding it to this file.

//msgp:timezone utc

type Block [32]byte

// test fixed-size struct
// size compilation
type Fixed struct {
	A float64
	B bool
}

type StandaloneBytes []byte

type AliasedType = Fixed
type AliasedType2 = *Fixed
type AliasedType3 = uint64
type AliasedType4 = *uint64

// tests edge-cases with
// compiling size compilation.
type X struct {
	Values         [32]byte    // should compile to 32*msgp.ByteSize; encoded as Bin
	ValuesPtr      *[32]byte   // check (*)[:] deref
	More           Block       // should be identical to the above
	Others         [][32]int32 // should compile to len(x.Others)*32*msgp.Int32Size
	Matrix         [][]int32   // should not optimize
	ManyFixed      []Fixed
	WeirdTag       string                       `msg:"\x0b"`
	ZCBytes        []byte                       `msg:",zerocopy"`
	ZCBytesAN      []byte                       `msg:",zerocopy,allownil"`
	ZCBytesOE      []byte                       `msg:",zerocopy,omitempty"`
	ZCBytesSlice   [][]byte                     `msg:",zerocopy"`
	ZCBytesArr     [2][]byte                    `msg:",zerocopy"`
	ZCBytesMap     map[string][]byte            `msg:",zerocopy"`
	ZCBytesMapDeep map[string]map[string][]byte `msg:",zerocopy"`
	CustomBytes    CustomBytes                  `msg:",zerocopy"`
	Renamed1       AliasedType
	Renamed2       AliasedType2
	Renamed3       AliasedType3
	Renamed4       AliasedType4
}

type TestType struct {
	F   *float64          `msg:"float"`
	Els map[string]string `msg:"elements"`
	Obj struct { // test anonymous struct
		ValueA string `msg:"value_a"`
		ValueB []byte `msg:"value_b"`
	} `msg:"object"`
	Child           *TestType   `msg:"child"`
	Time            time.Time   `msg:"time"`
	Any             interface{} `msg:"any"`
	Appended        msgp.Raw    `msg:"appended"`
	Num             msgp.Number `msg:"num"`
	Byte            byte
	Rune            rune
	RunePtr         *rune
	RunePtrPtr      **rune
	RuneSlice       []rune
	Slice1          []string
	Slice2          []string
	SlicePtr        *[]string
	MapStringEmpty  map[string]struct{}
	MapStringEmpty2 map[string]EmptyStruct
}

//msgp:tuple Object
type Object struct {
	ObjectNo        string   `msg:"objno"`
	Slice1          []string `msg:"slice1"`
	Slice2          []string `msg:"slice2"`
	MapMap          map[string]map[string]string
	MapStringEmpty  map[string]struct{}
	MapStringEmpty2 map[string]EmptyStruct
	ZCBytes         []byte                       `msg:",zerocopy"`
	ZCBytesAN       []byte                       `msg:",zerocopy,allownil"`
	ZCBytesOE       []byte                       `msg:",zerocopy,omitempty"`
	ZCBytesSlice    [][]byte                     `msg:",zerocopy"`
	ZCBytesArr      [2][]byte                    `msg:",zerocopy"`
	ZCBytesMap      map[string][]byte            `msg:",zerocopy"`
	ZCBytesMapDeep  map[string]map[string][]byte `msg:",zerocopy"`
	CustomBytes     CustomBytes                  `msg:",zerocopy"`
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

// Test nested aliases
type (
	FastAlias      TestFast
	AliasContainer struct {
		Fast FastAlias
	}
)

// Test dependency resolution
type (
	IntA int
	IntB IntA
	IntC IntB
)

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

//msgp:shim SpecialID as:[]byte using:toBytes/fromBytes

type (
	SpecialID string
	TestObj   struct{ ID1, ID2 SpecialID }
)

func toBytes(id SpecialID) []byte   { return []byte(string(id)) }
func fromBytes(id []byte) SpecialID { return SpecialID(string(id)) }

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

//msgp:shim *os.File as:string using:filetostr/filefromstr

func filetostr(f *os.File) string {
	return f.Name()
}

func filefromstr(s string) *os.File {
	f, _ := os.Open(s)
	return f
}

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

// test pass-specific directive
//msgp:decode ignore Insane

type Insane [3]map[string]struct{ A, B CustomInt }

type Custom struct {
	Bts   CustomBytes          `msg:"bts"`
	Mp    map[string]*Embedded `msg:"mp"`
	Enums []MyEnum             `msg:"enums"` // test explicit enum shim
	Some  FileHandle           `msg:"file_handle"`
}

type Files []*os.File

type FileHandle struct {
	Relevant Files  `msg:"files"`
	Name     string `msg:"name"`
}

type (
	CustomInt   int
	CustomBytes []byte
)

type Wrapper struct {
	Tree *Tree
}

type Tree struct {
	Children []Tree
	Element  int
	Parent   *Wrapper
}

// Ensure all different widths of integer can be used as constant keys.
const (
	ConstantInt    int    = 8
	ConstantInt8   int8   = 8
	ConstantInt16  int16  = 8
	ConstantInt32  int32  = 8
	ConstantInt64  int64  = 8
	ConstantUint   uint   = 8
	ConstantUint8  uint8  = 8
	ConstantUint16 uint16 = 8
	ConstantUint32 uint32 = 8
	ConstantUint64 uint64 = 8
)

type ArrayConstants struct {
	ConstantInt    [ConstantInt]string
	ConstantInt8   [ConstantInt8]string
	ConstantInt16  [ConstantInt16]string
	ConstantInt32  [ConstantInt32]string
	ConstantInt64  [ConstantInt64]string
	ConstantUint   [ConstantUint]string
	ConstantUint8  [ConstantUint8]string
	ConstantUint16 [ConstantUint16]string
	ConstantUint32 [ConstantUint32]string
	ConstantUint64 [ConstantUint64]string
	ConstantHex    [0x16]string
	ConstantOctal  [0o7]string
}

// Ensure non-msg struct tags work:
// https://github.com/tinylib/msgp/issues/201

type NonMsgStructTags struct {
	A      []string `json:"fooJSON" msg:"fooMsgp"`
	B      string   `json:"barJSON"`
	C      []string `json:"bazJSON" msg:"-"`
	Nested []struct {
		A          []string `json:"a"`
		B          string   `json:"b"`
		C          []string `json:"c"`
		VeryNested []struct {
			A []string            `json:"a"`
			B []string            `msg:"bbbb" xml:"-"`
			C map[string]struct{} `msg:"cccc"`
		}
	}
}

type EmptyStruct struct{}

type StructByteSlice struct {
	ABytes      []byte       `msg:",allownil"`
	AString     []string     `msg:",allownil"`
	ABool       []bool       `msg:",allownil"`
	AInt        []int        `msg:",allownil"`
	AInt8       []int8       `msg:",allownil"`
	AInt16      []int16      `msg:",allownil"`
	AInt32      []int32      `msg:",allownil"`
	AInt64      []int64      `msg:",allownil"`
	AUint       []uint       `msg:",allownil"`
	AUint8      []uint8      `msg:",allownil"`
	AUint16     []uint16     `msg:",allownil"`
	AUint32     []uint32     `msg:",allownil"`
	AUint64     []uint64     `msg:",allownil"`
	AFloat32    []float32    `msg:",allownil"`
	AFloat64    []float64    `msg:",allownil"`
	AComplex64  []complex64  `msg:",allownil"`
	AComplex128 []complex128 `msg:",allownil"`
	AStruct     []Fixed      `msg:",allownil"`
}

type NumberJSONSample struct {
	Single json.Number
	Array  []json.Number
	Map    map[string]json.Number
	OE     json.Number `msg:",omitempty"`
}

type Flobbity struct {
	A Flobs `msg:"a,omitempty"`
	B Flobs `msg:"b,omitempty"`
}

type Flobs []Flob

type Flob struct {
	X Numberwang `msg:"x"`
	Y int8       `msg:"y"`
	Z int8       `msg:"z"`
	W int32      `msg:"w"`
}

type Numberwang int8

//msgp:ignore ExternalString
type ExternalString string
type ExternalArr [4]byte
type ExternalInt int

//msgp:ignore regex:IGNORE

type RegexIGNORETest struct{}

// Will fail to compile if also generated
func (z *RegexIGNORETest) Msgsize() int {
	return 0
}

//msgp:size ignore regex:IGNSIZE

type RegexIGNSIZETest struct{}

// Will fail to compile if also generated
func (z *RegexIGNSIZETest) Msgsize() int {
	return 0
}
