package parse

type TestType struct {
	M     map[string]string `msg:"map"`
	V     *float64          `msg:"float"`
	Thing struct {
		A []byte `msg:"a"`
		B []byte `msg:"b"`
	} `msg:"thing"`
}
