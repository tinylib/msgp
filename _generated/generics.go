package _generated

import "github.com/tinylib/msgp/msgp"

//go:generate msgp -v

type Int64 int

type GenericTest[T any, P msgp.RTFor[T]] struct {
	A  T
	B  P
	C  []T          `msg:",allownil"`
	D  map[string]T `msg:",allownil"`
	E  GenericTest2[T, P, string]
	F  []GenericTest2[T, P, string]          `msg:",allownil"`
	G  map[string]GenericTest2[T, P, string] `msg:",allownil"`
	AP *T
	CP []*T          `msg:",allownil"`
	DP map[string]*T `msg:",allownil"`
	EP *GenericTest2[T, P, string]
	FP []*GenericTest2[T, P, string]          `msg:",allownil"`
	GP map[string]*GenericTest2[T, P, string] `msg:",allownil"`
}

type GenericTest2[T any, P msgp.RTFor[T], B any] struct {
	A T
}

//msgp:tuple GenericTuple
type GenericTuple[T any, P msgp.RTFor[T], B any] struct {
	A  T
	B  P
	C  []T          `msg:",allownil"`
	D  map[string]T `msg:",allownil"`
	E  GenericTest2[T, P, string]
	F  []GenericTest2[T, P, string]          `msg:",allownil"`
	G  map[string]GenericTest2[T, P, string] `msg:",allownil"`
	AP *T
	CP []*T          `msg:",allownil"`
	DP map[string]*T `msg:",allownil"`
	EP *GenericTest2[T, P, string]
	FP []*GenericTest2[T, P, string]          `msg:",allownil"`
	GP map[string]*GenericTest2[T, P, string] `msg:",allownil"`
}

// Type that doesn't have any fields using the generic type should just output valid code.
type GenericTestUnused[T any] struct {
	A string
}

type GenericTestTwo[A, B any, AP msgp.RTFor[A], BP msgp.RTFor[B]] struct {
	A  A
	B  AP
	C  []A          `msg:",allownil"`
	D  map[string]T `msg:",allownil"`
	E  GenericTest2[A, AP, string]
	F  []GenericTest2[A, AP, string]          `msg:",allownil"`
	G  map[string]GenericTest2[A, AP, string] `msg:",allownil"`
	AP *A
	CP []*A          `msg:",allownil"`
	DP map[string]*A `msg:",allownil"`
	EP *GenericTest2[A, AP, string]
	FP []*GenericTest2[A, AP, string]          `msg:",allownil"`
	GP map[string]*GenericTest2[A, AP, string] `msg:",allownil"`

	A2  B
	B2  BP
	C2  []B          `msg:",allownil"`
	D2  map[string]B `msg:",allownil"`
	E2  GenericTest2[B, BP, string]
	F2  []GenericTest2[B, BP, string]          `msg:",allownil"`
	G2  map[string]GenericTest2[B, BP, string] `msg:",allownil"`
	AP2 *B
	CP2 []*B          `msg:",allownil"`
	DP2 map[string]*B `msg:",allownil"`
	EP2 *GenericTest2[B, BP, string]
	FP2 []*GenericTest2[B, BP, string]          `msg:",allownil"`
	GP2 map[string]*GenericTest2[B, BP, string] `msg:",allownil"`
}

type GenericTest3List[A, B any, AP msgp.RTFor[A], BP msgp.RTFor[B]] []GenericTest3[A, B, AP, BP]
type GenericTest3Map[A, B any, AP msgp.RTFor[A], BP msgp.RTFor[B]] map[string]GenericTest3[A, B, AP, BP]

type GenericTest3Array[A, B any, AP msgp.RTFor[A], BP msgp.RTFor[B]] [5]GenericTest3[A, B, AP, BP]

type GenericTest3[A, B any, _ msgp.RTFor[A], _ msgp.RTFor[B]] struct {
	A A
	B B
}

// GenericTest4 has the msgp.RTFor constraint as a sub-constraint.
type GenericTest4[T any, P interface {
	*T
	msgp.RTFor[T]
}] struct {
	Totals [60]T
}
