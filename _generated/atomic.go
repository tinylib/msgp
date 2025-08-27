package _generated

import "sync/atomic"

//go:generate msgp

//msgp:clearomitted

type AtomicTest struct {
	A       atomic.Int64
	B       atomic.Uint64
	C       atomic.Int32
	D       atomic.Uint32
	E       atomic.Bool
	Omitted struct {
		A atomic.Int64  `msg:",omitempty"`
		B atomic.Uint64 `msg:",omitempty"`
		C atomic.Int32  `msg:",omitempty"`
		D atomic.Uint32 `msg:",omitempty"`
		E atomic.Bool   `msg:",omitempty"`
	}
	Ptr struct {
		A       *atomic.Int64
		B       *atomic.Uint64
		C       *atomic.Int32
		D       *atomic.Uint32
		E       *atomic.Bool
		Omitted struct {
			A *atomic.Int64  `msg:",omitempty"`
			B *atomic.Uint64 `msg:",omitempty"`
			C *atomic.Int32  `msg:",omitempty"`
			D *atomic.Uint32 `msg:",omitempty"`
			E *atomic.Bool   `msg:",omitempty"`
		}
	}
	F []atomic.Int32
	G []*atomic.Int32
	// We use pointers, as values don't make sense for maps.
	H map[string]*atomic.Int32
}
