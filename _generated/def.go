package _generated

type TestType struct {
	F   *float64          `msg:"float"`
	Els map[string]string `msg:"elements"`
	Obj struct {
		ValueA string `msg:"value_a"`
		ValueB []byte `msg:"value_b"`
	} `msg:"object"`
	Child *TestType `msg:"child"`
}
