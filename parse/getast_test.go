package parse

import (
	"fmt"
	"testing"
)

//go:generate msgp

type AliasType = string

type Func1 func() error

type Func2 = func() int

type TestStruct struct {
	Str   string
	Alias AliasType
}

func TestParseType(t *testing.T) {
	f, err := File("getast_test.go", false)
	if err != nil {
		panic(err)
	}
	for k, v := range f.Identities {
		fmt.Printf("Identity: %s: varname: %s, typename: %s, complexity: %d\n", k, v.Varname(), v.TypeName(), v.Complexity())
	}
	for k, v := range f.Aliases {
		fmt.Printf("Alias: %s: type: %v\n", k, v)
	}
}
