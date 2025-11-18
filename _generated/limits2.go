package _generated

//go:generate msgp

//msgp:limit arrays:200 maps:100

type LimitTestData2 struct {
	BigArray [20]int        `msg:"big_array"`
	BigMap   map[string]int `msg:"big_map"`
}

// Test field limits vs file limits precedence with different file limits
// File limits: arrays:200 maps:100
type FieldOverrideTestData2 struct {
	TightSlice   []int            `msg:"tight_slice,limit=30"`   // Field limit (30) < file limit (200)
	MediumSlice  []string         `msg:"medium_slice,limit=150"` // Field limit (150) < file limit (200)
	LooseSlice   []byte           `msg:"loose_slice,limit=300"`  // Field limit (300) > file limit (200)
	TightMap     map[string]int   `msg:"tight_map,limit=20"`     // Field limit (20) < file limit (100)
	MediumMap    map[int]string   `msg:"medium_map,limit=75"`    // Field limit (75) < file limit (100)
	LooseMap     map[string][]int `msg:"loose_map,limit=150"`    // Field limit (150) > file limit (100)
	DefaultSlice []int            `msg:"default_slice"`          // No field limit, uses file limit (200)
	DefaultMap   map[int]int      `msg:"default_map"`            // No field limit, uses file limit (100)
}
