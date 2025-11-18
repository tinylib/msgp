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

// Test field limits vs file limits precedence
// File limits: arrays:100 maps:50
type FieldOverrideTestData struct {
	TightSlice  []int          `msg:"tight_slice,limit=10"`   // Field limit (10) < file limit (100)
	LooseSlice  []string       `msg:"loose_slice,limit=200"`  // Field limit (200) > file limit (100)
	TightMap    map[string]int `msg:"tight_map,limit=5"`      // Field limit (5) < file limit (50)
	LooseMap    map[int]string `msg:"loose_map,limit=80"`     // Field limit (80) > file limit (50)
	DefaultSlice []byte        `msg:"default_slice"`          // No field limit, uses file limit (100)
	DefaultMap   map[string]string `msg:"default_map"`        // No field limit, uses file limit (50)
}
