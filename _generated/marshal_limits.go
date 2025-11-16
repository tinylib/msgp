//msgp:limit arrays:30 maps:20 marshal:true

package _generated

//go:generate msgp

// Test structures for marshal-time limit enforcement
type MarshalLimitTestData struct {
	SmallArray [5]int         `msg:"small_array"`
	TestSlice  []string       `msg:"test_slice"`
	TestMap    map[string]int `msg:"test_map"`
}

// Test field limits vs file limits precedence with marshal:true
// File limits: arrays:30 maps:20 marshal:true
type MarshalFieldOverrideTestData struct {
	TightSlice   []int            `msg:"tight_slice,limit=10"`  // Field limit (10) < file limit (30)
	LooseSlice   []string         `msg:"loose_slice,limit=50"`  // Field limit (50) > file limit (30)
	TightMap     map[string]int   `msg:"tight_map,limit=5"`     // Field limit (5) < file limit (20)
	LooseMap     map[int]string   `msg:"loose_map,limit=40"`    // Field limit (40) > file limit (20)
	DefaultSlice []byte           `msg:"default_slice"`         // No field limit, uses file limit (30)
	DefaultMap   map[string]byte  `msg:"default_map"`           // No field limit, uses file limit (20)
}
