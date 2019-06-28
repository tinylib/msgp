package _generated

import "sync"

//go:generate msgp
//go:generate go vet

type Foo struct {
	I    struct{}
	lock sync.Mutex
}
