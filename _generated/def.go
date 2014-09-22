package _generated

//go:generate msgp

type TestType struct {
	F   *float64          `msg:"float"`
	Els map[string]string `msg:"elements"`
	Obj struct {
		ValueA string `msg:"value_a"`
		ValueB []byte `msg:"value_b"`
	} `msg:"object"`
	Child *TestType `msg:"child"`
}

type TestFast struct {
	Lat, Long, Alt float64
	Data           []byte
}

type TestHidden struct {
	A   string
	B   []float64
	Bad func(string) bool
}
