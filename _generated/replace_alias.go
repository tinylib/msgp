package _generated

import "github.com/tinylib/msgp/msgp"

//go:generate msgp

//msgp:shim ReplAliasShim as:string using:replAliasToString/stringToReplAlias
//msgp:replace ReplChainOuter with:[]ReplChainInner
//msgp:replace ReplChainInner with:[]int
//msgp:replace ReplTarget with:[]int
//msgp:replace ReplInList with:string
//msgp:replace ReplStruct with:struct{A[]int}
//msgp:replace ReplGenStr with:string

type ReplAliasSt struct{ A string }

// A shim on an alias name.
type ReplAliasShim = ReplAliasSt

func replAliasToString(v ReplAliasShim) string { return v.A }
func stringToReplAlias(s string) ReplAliasShim { return ReplAliasShim{A: s} }

// An alias chain.
type ReplStrs = []ReplStr
type ReplStr = string

// A replace on the target of an alias.
type ReplOnTarget = ReplTarget

// A replaced name inside another alias.
type ReplList = []ReplInList

type ReplAlias struct {
	Chain  ReplChainOuter // replaces in either order
	Strs   ReplStrs
	Shim   ReplAliasShim
	Target ReplOnTarget
	List   ReplList
	PS     *ReplStruct
}

type ReplGenInner struct{ B int }

// Names that are not generic must not get the type parameters.
type ReplGeneric[Key any, Val msgp.RTFor[Key]] struct {
	K  Key
	M  map[string]int
	PI *ReplGenInner
	R  ReplGenStr
}
