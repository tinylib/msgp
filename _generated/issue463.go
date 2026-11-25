package _generated

//go:generate msgp

type Issue463 struct {
	B []byte `msg:",omitempty,allownil"`
}
