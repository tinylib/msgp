package main

import "testing"

func TestIssue275Tuples(t *testing.T) {
	mainFilename, err := generate(t, issue275Tuples)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}
	goExec(t, mainFilename, false) // exec go run
	goExec(t, mainFilename, true)  // exec go test
}

var issue275Tuples = `package main

import (
	"fmt"
	"os"
)

//go:generate msgp

//msgp:tuple Test1
type Test1 struct {
	Foo string
}

//msgp:tuple Test2
type Test2 struct {
	Foo string
	Bar string
}

//msgp:tuple Test3
type Test3 struct {
	Foo string
	Bar string
	Baz string
}

//msgp:vartuple Test
type Test struct {
	Foo string
	Bar string
}

func main() {
	t1 := Test1{Foo: "Foo1"}
	d1, err := t1.MarshalMsg(nil)
	if err != nil {
		fmt.Println("Test1 MarshalMsg failed:", err)
		os.Exit(1)
	}

	t2 := Test2{Foo: "Foo2", Bar: "Bar2"}
	d2, err := t2.MarshalMsg(nil)
	if err != nil {
		fmt.Println("Test2 MarshalMsg failed:", err)
		os.Exit(1)
	}

	t3 := Test3{Foo: "Foo3", Bar: "Bar3", Baz: "Baz3"}
	d3, err := t3.MarshalMsg(nil)
	if err != nil {
		fmt.Println("Test3 MarshalMsg failed:", err)
		os.Exit(1)
	}

	var msg Test

	msg = Test{}
	if _, err := msg.UnmarshalMsg(d1); err != nil {
		fmt.Println("Test UnmarshalMsg from Test1 failed:", err)
		os.Exit(1)
	}
	if msg.Foo != "Foo1" {
		fmt.Println("Test UnmarshalMsg from Test1 bad msg:", msg)
		os.Exit(1)
	}

	msg = Test{}
	if _, err := msg.UnmarshalMsg(d2); err != nil {
		fmt.Println("Test UnmarshalMsg from Test2 failed:", err)
		os.Exit(1)
	}
	if msg.Foo != "Foo2" || msg.Bar != "Bar2" {
		fmt.Println("Test UnmarshalMsg from Test2 bad msg:", msg)
		os.Exit(1)
	}

	msg = Test{}
	if _, err := msg.UnmarshalMsg(d3); err != nil {
		fmt.Println("Test UnmarshalMsg from Test3 failed:", err)
		os.Exit(1)
	}
	if msg.Foo != "Foo3" || msg.Bar != "Bar3" {
		fmt.Println("Test UnmarshalMsg from Test3 bad msg:", msg)
		os.Exit(1)
	}

	var msg1 Test1
	if _, err := msg1.UnmarshalMsg(d2); err != nil {
		if err.Error() != "msgp: wanted array of size 1; got 2" {
			fmt.Println("Test1 UnmarshalMsg from Test2 failed:", err)
			os.Exit(1)
		}
	}

	var msg2 Test2
	if _, err := msg2.UnmarshalMsg(d2); err != nil {
		fmt.Println("Test2 UnmarshalMsg from Test2 failed:", err)
		os.Exit(1)
	}
	if msg2.Foo != "Foo2" || msg2.Bar != "Bar2" {
		fmt.Println("Test3 UnmarshalMsg from Test2 bad msg:", msg)
		os.Exit(1)
	}

	var msg3 Test3
	if _, err := msg3.UnmarshalMsg(d2); err != nil {
		if err.Error() != "msgp: wanted array of size 3; got 2" {
			fmt.Println("Test3 UnmarshalMsg from Test2 failed:", err)
			os.Exit(1)
		}
	}
}
`
