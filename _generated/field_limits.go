package _generated

//go:generate msgp

// Aliased types for testing
type AliasedSlice []string
type AliasedMap map[string]bool
type AliasedIntSlice []int

// Test structures for field-level limit tags
type FieldLimitTestData struct {
	SmallSlice []int          `msg:"small_slice,limit=5"`
	LargeSlice []string       `msg:"large_slice,limit=100"`
	SmallMap   map[string]int `msg:"small_map,limit=3"`
	LargeMap   map[int]string `msg:"large_map,limit=20"`
	NoLimit    []byte         `msg:"no_limit"`  // Uses file-level limits if any
	FixedArray [10]int        `msg:"fixed_array,limit=2"` // Should be ignored
}

// Test structure with aliased types and field limits
type AliasedFieldLimitTestData struct {
	SmallAliasedSlice AliasedSlice    `msg:"small_aliased_slice,limit=3"`
	LargeAliasedSlice AliasedSlice    `msg:"large_aliased_slice,limit=50"`
	SmallAliasedMap   AliasedMap      `msg:"small_aliased_map,limit=2"`
	LargeAliasedMap   AliasedMap      `msg:"large_aliased_map,limit=25"`
	IntSliceAlias     AliasedIntSlice `msg:"int_slice_alias,limit=10"`
	NoLimitAlias      AliasedSlice    `msg:"no_limit_alias"` // Uses file-level limits
}