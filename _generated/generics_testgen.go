package _generated

import "github.com/tinylib/msgp/msgp"

//go:generate msgp

// The tests generated for GenericTestgenB follow the comment for the
// skipped generic type and must not be commented out by it.
type GenericTestgenA[T any, P msgp.RTFor[T]] struct{ A T }

type GenericTestgenB struct{ B int }
