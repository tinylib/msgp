package main

import "testing"

func TestIssue463OmitEmptyAllowNil(t *testing.T) {
	mainFilename, err := generate(t, issue463OmitEmptyAllowNil)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	goExec(t, mainFilename, false) // exec go run
	goExec(t, mainFilename, true)  // exec go test
}

var issue463OmitEmptyAllowNil = `package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/tinylib/msgp/msgp"
)

type S struct {
	B []byte ` + "`" + `msg:",omitempty,allownil"` + "`" + `
}

func payload(isNil bool, value []byte) []byte {
	b := msgp.AppendMapHeader(nil, 1)
	b = msgp.AppendString(b, "B")
	if isNil {
		return msgp.AppendNil(b)
	}
	return msgp.AppendBytes(b, value)
}

func check(name string, data []byte, wantNil bool, want []byte) bool {
	var u, d S
	if _, err := u.UnmarshalMsg(data); err != nil {
		fmt.Printf("%s: unmarshal error: %v\n", name, err)
		return false
	}
	if err := d.DecodeMsg(msgp.NewReader(bytes.NewReader(data))); err != nil {
		fmt.Printf("%s: decode error: %v\n", name, err)
		return false
	}
	for method, got := range map[string][]byte{"unmarshal": u.B, "decode": d.B} {
		if (got == nil) != wantNil || !bytes.Equal(got, want) {
			fmt.Printf("%s: %s got %#v, want %#v (nil %v)\n", name, method, got, want, wantNil)
			return false
		}
	}
	return true
}

func main() {
	ok := true
	ok = check("nil", payload(true, nil), true, nil) && ok
	ok = check("empty", payload(false, []byte{}), false, []byte{}) && ok
	ok = check("full", payload(false, []byte{1}), false, []byte{1}) && ok
	if !ok {
		os.Exit(1)
	}
}
`
