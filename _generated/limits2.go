package _generated

//go:generate msgp

//msgp:limit arrays:200 maps:100

type LimitTestData2 struct {
	BigArray [20]int        `msg:"big_array"`
	BigMap   map[string]int `msg:"big_map"`
}
