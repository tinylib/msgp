package main

import "testing"

func TestEmptyStructCollections(t *testing.T) {
	filename, err := generate(t, emptyStructCollections)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	goExec(t, filename, false)
	goExec(t, filename, true)
}

var emptyStructCollections = `package main

import (
	"bytes"
	"fmt"
	"reflect"

	"github.com/tinylib/msgp/msgp"
)

type Empty struct{}
type ZeroArray [0]Empty
type NestedArray [2][1]Empty

type Collection struct {
	Named []Empty
	Anonymous []struct{}
	Array [2]Empty
}

func main() {
	input := Collection{
		Named: []Empty{{}, {}},
		Anonymous: []struct{}{{}, {}, {}},
	}
	data, err := input.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	if size := input.Msgsize(); size < len(data) {
		panic(fmt.Sprintf("Msgsize %d is smaller than encoded length %d", size, len(data)))
	}
	var decoded Collection
	if rest, err := decoded.UnmarshalMsg(data); err != nil || len(rest) != 0 {
		panic(fmt.Sprintf("UnmarshalMsg: rest %x, error %v", rest, err))
	}
	if !reflect.DeepEqual(input, decoded) {
		panic(fmt.Sprintf("UnmarshalMsg got %#v, want %#v", decoded, input))
	}

	var stream bytes.Buffer
	if err := msgp.Encode(&stream, &input); err != nil {
		panic(err)
	}
	if !bytes.Equal(stream.Bytes(), data) {
		panic("stream encoding differs from MarshalMsg")
	}
	var streamed Collection
	if err := msgp.Decode(&stream, &streamed); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(input, streamed) {
		panic(fmt.Sprintf("DecodeMsg got %#v, want %#v", streamed, input))
	}
}
`
