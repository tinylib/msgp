//msgp:limit arrays:100 maps:50

package _generated

//go:generate msgp

// Test structures for limit directive
type LimitedData struct {
	SmallArray [10]int        `msg:"small_array"`
	LargeSlice []byte         `msg:"large_slice"`
	SmallMap   map[string]int `msg:"small_map"`
}

type UnlimitedData struct {
	BigArray [1000]int        `msg:"big_array"`
	BigSlice []string         `msg:"big_slice"`
	BigMap   map[string][]int `msg:"big_map"`
}

type LimitTestData struct {
	SmallArray [10]int        `msg:"small_array"`
	LargeSlice []byte         `msg:"large_slice"`
	SmallMap   map[string]int `msg:"small_map"`
}
