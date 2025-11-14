//msgp:limit arrays:30 maps:20 marshal:true

package _generated

//go:generate msgp

// Test structures for marshal-time limit enforcement
type MarshalLimitTestData struct {
	SmallArray [5]int         `msg:"small_array"`
	TestSlice  []string       `msg:"test_slice"`
	TestMap    map[string]int `msg:"test_map"`
}
