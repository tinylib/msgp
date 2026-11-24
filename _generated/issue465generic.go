package _generated

import "github.com/tinylib/msgp/msgp"

//go:generate msgp

//msgp:shim Issue465GenStruct as:[]byte using:issue465GenStructToBytes/issue465GenStructFromBytes
//msgp:clearomitted

type Issue465GenStruct struct{ A string }

func issue465GenStructToBytes(v Issue465GenStruct) []byte   { return []byte(v.A) }
func issue465GenStructFromBytes(b []byte) Issue465GenStruct { return Issue465GenStruct{A: string(b)} }

// clearomitted must not add the type parameters to the zero value.
type Issue465Gen[T any, P msgp.RTFor[T]] struct {
	A      T
	Struct Issue465GenStruct `msg:",omitempty"`
}
