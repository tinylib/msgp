package _generated

//go:generate msgp

//msgp:compactfloats

//msgp:ignore F64
type F64 float64

//msgp:replace F64 with:float64

type Floats struct {
	A     float64
	B     float32
	Slice []float64
	Map   map[string]float64
	F     F64
	OE    float64 `msg:",omitempty"`
}
