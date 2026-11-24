package _generated

//go:generate msgp

type OmitZeroVal struct{ A int }

// IsZero has a value receiver, so it must not be called on a nil pointer.
func (z OmitZeroVal) IsZero() bool { return z.A == 0 }

type OmitZeroPtrRecv struct{ A int }

func (z *OmitZeroPtrRecv) IsZero() bool { return z.A == 0 }

type OmitZeroPtr struct {
	Val *OmitZeroVal     `msg:",omitzero"`
	Ptr *OmitZeroPtrRecv `msg:",omitzero"`
	X   int
}
