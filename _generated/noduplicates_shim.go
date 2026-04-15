package _generated

import "strings"

//go:generate msgp -v

//msgp:maps shim
//msgp:noduplicates

//msgp:shim NoDupShimKey as:string using:noDupShimEnc/noDupShimDec

// NoDupShimKey is a key type that normalizes to lowercase.
type NoDupShimKey string

func noDupShimEnc(k NoDupShimKey) string { return strings.ToLower(string(k)) }
func noDupShimDec(s string) NoDupShimKey { return NoDupShimKey(strings.ToLower(s)) }

// NoDupShimMap is a map with shimmed keys that normalize to lowercase.
// Two different wire strings (e.g. "Foo" and "foo") shim to the same key.
type NoDupShimMap map[NoDupShimKey]int

// NoDupShimStruct has a shimmed-key map as a field.
type NoDupShimStruct struct {
	Name string                  `msg:"name"`
	Data map[NoDupShimKey]string `msg:"data"`
}
