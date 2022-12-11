package main

import (
	"bytes"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

type SomeStruct struct {
	A string `msg:"a"`
}

type EmbeddedStruct struct {
	A string `msg:"a"`
}

// Example provides a decent variety of types and features and
// lets us test for basic functionality in TinyGo.
type Example struct {
	Int64       int64      `msg:"int64"`
	Uint64      uint64     `msg:"uint64"`
	Int32       int32      `msg:"int32"`
	Uint32      uint32     `msg:"uint32"`
	Int16       int32      `msg:"int16"`
	Uint16      uint32     `msg:"uint16"`
	Int8        int32      `msg:"int8"`
	Byte        byte       `msg:"byte"`
	Float64     float64    `msg:"float64"`
	Float32     float32    `msg:"float32"`
	String      string     `msg:"string"`
	ByteSlice   []byte     `msg:"byte_slice"`
	StringSlice []string   `msg:"string_slice"`
	IntArray    [2]int     `msg:"int_array"`
	SomeStruct  SomeStruct `msg:"some_struct"`

	EmbeddedStruct

	Omitted string `msg:"-"`

	OmitEmptyString string `msg:"omit_empty_string,omitempty"`
}

// Setup populuates the struct with test data
func (e *Example) Setup() {
	e.Int64 = 10
	e.Uint64 = 11
	e.Int32 = 12
	e.Uint32 = 13
	e.Int16 = 14
	e.Uint16 = 15
	e.Int8 = 16
	e.Byte = 17
	e.Float64 = 18.1
	e.Float32 = 19.2
	e.String = "astr"
	e.ByteSlice = []byte("bstr")
	e.StringSlice = []string{"a", "b"}
	e.IntArray = [...]int{20, 21}
	e.SomeStruct = SomeStruct{A: "x"}

	e.EmbeddedStruct = EmbeddedStruct{A: "y"}

	e.Omitted = "nope"

	e.OmitEmptyString = "here"
}

func (e *Example) Eq(e2 *Example) bool {
	if e.Int64 != e2.Int64 {
		return false
	}
	if e.Uint64 != e2.Uint64 {
		return false
	}
	if e.Int32 != e2.Int32 {
		return false
	}
	if e.Uint32 != e2.Uint32 {
		return false
	}
	if e.Int16 != e2.Int16 {
		return false
	}
	if e.Uint16 != e2.Uint16 {
		return false
	}
	if e.Int8 != e2.Int8 {
		return false
	}
	if e.Byte != e2.Byte {
		return false
	}
	if e.Float64 != e2.Float64 {
		return false
	}
	if e.Float32 != e2.Float32 {
		return false
	}
	if e.String != e2.String {
		return false
	}
	if bytes.Compare(e.ByteSlice, e2.ByteSlice) != 0 {
		return false
	}
	if len(e.StringSlice) != len(e2.StringSlice) {
		return false
	}
	for i := 0; i < len(e.StringSlice); i++ {
		if e.StringSlice[i] != e2.StringSlice[i] {
			return false
		}
	}
	if len(e.IntArray) != len(e2.IntArray) {
		return false
	}
	for i := 0; i < len(e.IntArray); i++ {
		if e.IntArray[i] != e2.IntArray[i] {
			return false
		}
	}
	if e.SomeStruct.A != e2.SomeStruct.A {
		return false
	}

	if e.EmbeddedStruct.A != e2.EmbeddedStruct.A {
		return false
	}

	if e.Omitted != e2.Omitted {
		return false
	}

	if e.OmitEmptyString != e2.OmitEmptyString {
		return false
	}

	return true
}

var buf [256]byte

func main() {
	var e Example
	e.Setup()

	b, err := e.MarshalMsg(buf[:0])
	if err != nil {
		panic(err)
	}
	b1 := b

	var e2 Example
	_, err = e2.UnmarshalMsg(b)
	if err != nil {
		panic(err)
	}

	println("marshal/unmarshal done: ", &e2)

	e.Omitted = ""
	if !e.Eq(&e2) {
		panic("comparison after marshal/unmarhsal failed")
	}

	var wbuf bytes.Buffer
	mw := msgp.NewWriterSize(&wbuf, 64)

	e.Omitted = "other"
	err = e.EncodeMsg(mw)
	if err != nil {
		panic(err)
	}
	mw.Flush()

	e2 = Example{}
	mr := msgp.NewReaderSize(bytes.NewReader(wbuf.Bytes()), 64)
	err = e2.DecodeMsg(mr)
	if err != nil {
		panic(err)
	}

	println("writer/reader done: ", &e2)
	e.Omitted = ""
	if !e.Eq(&e2) {
		panic("comparison after writer/reader failed")
	}

	if bytes.Compare(wbuf.Bytes(), b1) != 0 {
		panic("writer and marshal produced different results")
	}
}
