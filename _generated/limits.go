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
	TightSlice   []int             `msg:"tight_slice,limit=10"`  // Field limit (10) < file limit (100)
	LooseSlice   []string          `msg:"loose_slice,limit=200"` // Field limit (200) > file limit (100)
	TightMap     map[string]int    `msg:"tight_map,limit=5"`     // Field limit (5) < file limit (50)
	LooseMap     map[int]string    `msg:"loose_map,limit=80"`    // Field limit (80) > file limit (50)
	DefaultSlice []byte            `msg:"default_slice"`         // No field limit, uses file limit (100)
	DefaultMap   map[string]string `msg:"default_map"`           // No field limit, uses file limit (50)
}

// Test allownil functionality with limits
// File limits: arrays:100 maps:50
type AllowNilTestData struct {
	NilSlice        []byte                 `msg:"nil_slice,allownil"`                           // Can be nil, uses file limit (100)
	NilTightSlice   []int                  `msg:"nil_tight_slice,allownil,limit=5"`             // Can be nil, field limit (5) < file limit (100)
	NilLooseSlice   []string               `msg:"nil_loose_slice,allownil,limit=200"`           // Can be nil, field limit (200) > file limit (100)
	NilZCSlice      []byte                 `msg:"nil_zc_slice,allownil,zerocopy"`               // Can be nil, zerocopy, uses file limit (100)
	NilZCTightSlice []byte                 `msg:"nil_zc_tight_slice,allownil,zerocopy,limit=5"` // Can be nil, zerocopy, field limit (5)
	NilMap          map[string]int         `msg:"nil_map,allownil"`                             // Can be nil, uses file limit (50)
	NilTightMap     map[string]int         `msg:"nil_tight_map,allownil,limit=3"`               // Can be nil, field limit (3) < file limit (50)
	NilLooseMap     map[int]string         `msg:"nil_loose_map,allownil,limit=75"`              // Can be nil, field limit (75) > file limit (50)
	RegularSlice    []byte                 `msg:"regular_slice"`                                // Cannot be nil, uses file limit (100)
	RegularMap      map[string]int         `msg:"regular_map"`                                  // Cannot be nil, uses file limit (50)
	ObjSlice        []LimitedData          `msg:"obj_slice"`                                    // Cannot be nil, uses file limit (100)
	ObjMap          map[string]LimitedData `msg:"obj_map"`                                      // Cannot be nil, uses file limit (50)
}

type ObjSlices []LimitedData
type ObjMaps map[string]LimitedData

// Test allownil with zero-sized allocations
type AllowNilZeroTestData struct {
	NilBytes     []byte `msg:"nil_bytes,allownil"` // Can be nil, should allocate zero-sized slice
	RegularBytes []byte `msg:"regular_bytes"`      // Cannot be nil
}

type StandaloneBytesLimited []byte
