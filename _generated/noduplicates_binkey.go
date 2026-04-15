package _generated

//go:generate msgp -unexported -v

//msgp:maps binkeys
//msgp:noduplicates

// NoDupBinKeyInt is a map with int keys (binary-encoded).
type NoDupBinKeyInt map[int]string

// NoDupBinKeyUint32 is a map with uint32 keys (binary-encoded).
type NoDupBinKeyUint32 map[uint32]int

// NoDupBinKeyStruct has both string-keyed and binary-keyed maps.
type NoDupBinKeyStruct struct {
	Name   string         `msg:"name"`
	IntMap map[int]string `msg:"int_map"`
	StrMap map[string]int `msg:"str_map"`
}

// NoDupBinKeyArray is a map with array keys.
type NoDupBinKeyArray map[[2]byte]int
