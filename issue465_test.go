package main

import "testing"

func TestIssue465AllowNilNonNilableBytes(t *testing.T) {
	mainFilename, err := generate(t, issue465AllowNilNonNilableBytes)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	goExec(t, mainFilename, false) // exec go run
	goExec(t, mainFilename, true)  // exec go test
}

var issue465AllowNilNonNilableBytes = `package main

import (
	"bytes"
	"fmt"
	"os"
	"reflect"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp

//msgp:replace issue465Str with:[]byte
//msgp:replace Issue465Bytes with:[]byte
//msgp:shim Issue465VT as:[]byte using:vtToBytes/vtFromBytes

type issue465Str string

type Issue465Bytes []byte

type Issue465VT struct{ A string }

func vtToBytes(v Issue465VT) []byte   { return []byte(v.A) }
func vtFromBytes(b []byte) Issue465VT { return Issue465VT{A: string(b)} }

type S struct {
	V  issue465Str   ` + "`" + `msg:",allownil"` + "`" + `
	PV *issue465Str  ` + "`" + `msg:",allownil"` + "`" + `
	W  Issue465VT    ` + "`" + `msg:",allownil"` + "`" + `
	B  Issue465Bytes ` + "`" + `msg:",allownil"` + "`" + `
}

func checkRoundTrip(name string, in S) bool {
	data, err := in.MarshalMsg(nil)
	if err != nil {
		fmt.Printf("%s: marshal error: %v\n", name, err)
		return false
	}
	var out S
	if _, err := out.UnmarshalMsg(data); err != nil {
		fmt.Printf("%s: unmarshal error: %v\n", name, err)
		return false
	}
	if !reflect.DeepEqual(in, out) {
		fmt.Printf("%s: unmarshal got %#v, want %#v\n", name, out, in)
		return false
	}

	var buf bytes.Buffer
	if err := msgp.Encode(&buf, &in); err != nil {
		fmt.Printf("%s: encode error: %v\n", name, err)
		return false
	}
	out = S{}
	if err := msgp.Decode(&buf, &out); err != nil {
		fmt.Printf("%s: decode error: %v\n", name, err)
		return false
	}
	if !reflect.DeepEqual(in, out) {
		fmt.Printf("%s: decode got %#v, want %#v\n", name, out, in)
		return false
	}
	return true
}

func main() {
	v := issue465Str("p")
	ok := true
	ok = checkRoundTrip("filled", S{V: "a", PV: &v, W: Issue465VT{A: "b"}, B: Issue465Bytes{1}}) && ok
	ok = checkRoundTrip("nil bytes", S{}) && ok
	ok = checkRoundTrip("empty bytes", S{B: Issue465Bytes{}}) && ok
	if !ok {
		os.Exit(1)
	}
}
`
