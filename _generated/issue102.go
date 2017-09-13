package _generated

//go:generate msgp

type Issue102 struct{}

type Issue102deep struct {
	A int
	X struct{}
	Y struct{}
	Z int
}
