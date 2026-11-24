package _generated

import "github.com/tinylib/msgp/msgp"

//go:generate msgp

//msgp:replace ReplChainOuter with:[]ReplChainInner
//msgp:replace ReplChainInner with:[]int
//msgp:replace ReplTarget with:[]int
//msgp:replace ReplInList with:string
//msgp:replace ReplGenStr with:string

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
	Target ReplOnTarget
	List   ReplList
}

type ReplGenInner struct{ B int }

// Names that are not generic must not get the type parameters.
type ReplGeneric[Key any, Val msgp.RTFor[Key]] struct {
	K  Key
	M  map[string]int
	PI *ReplGenInner
	R  ReplGenStr
}
