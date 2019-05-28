package test

//go:generate msgp

type AliasType = string

type DefType string

type Func1 func() error

type Func2 = func() int

type TestStruct struct {
	Str   string
	Alias AliasType
}
