package _generated

import "github.com/tinylib/msgp/msgp"

//go:generate msgp

//msgp:replace Issue462Map with:map[string]int
//msgp:replace Issue462Slice with:[]string
//msgp:replace Issue462Bytes with:[4]byte
//msgp:replace Issue462Ints with:[3]int
//msgp:replace Issue462Ptr with:*int

// A manual replace of a type alias wins over the alias.
//msgp:replace Issue462Alias with:[4]byte

type Issue462Alias = Issue462Bytes

type Issue462Names = []string

type Issue462 struct {
	Map       Issue462Map
	PMap      *Issue462Map
	PPMap     **Issue462Map
	Slice     Issue462Slice
	PSlice    *Issue462Slice
	Bytes     Issue462Bytes
	PBytes    *Issue462Bytes
	PInts     *Issue462Ints
	MapPSlice map[string]*Issue462Slice
	SlicePMap []*Issue462Map
	Ptr       Issue462Ptr
	PPtr      *Issue462Ptr
	Alias     Issue462Alias
	PAlias    *Issue462Alias
	Names     Issue462Names
	PNames    *Issue462Names
}

// Replaced and aliased names must not get the type parameters.
type Issue462Gen[T any, P msgp.RTFor[T]] struct {
	A      T
	PBytes *Issue462Bytes
	Slice  Issue462Slice
	Names  Issue462Names
	PNames *Issue462Names
	PPS    **Issue462Slice
	MPS    map[string]*Issue462Slice
}
