package _generated

//go:generate msgp

//msgp:maps autoshim
//msgp:noduplicates
//msgp:shim Issue459Shim as:[]byte using:issue459ShimToBytes/issue459ShimFromBytes

type Issue459Buckets map[string]uint64

type Issue459Slice []string

type Issue459 struct {
	Count       uint64
	Buckets     *Issue459Buckets
	Map         *map[string]uint64
	MapStruct   *map[string]Issue459Inner
	MapPtr      *map[string]*uint64
	MapIntKey   *map[int]string
	Slice       *[]string
	NamedSlice  *Issue459Slice
	Array       *[3]int
	DoublePtr   **uint64
	DoubleMap   **map[string]uint64
	PtrSliceMap *[]map[string]uint64
	Limited     *[]byte `msg:",limit=10"`
}

type Issue459Inner struct {
	A int
}

// Issue459Shim is not nilable, so allownil must not apply through a pointer to it.
type Issue459Shim struct{ A string }

func issue459ShimToBytes(v Issue459Shim) []byte   { return []byte(v.A) }
func issue459ShimFromBytes(b []byte) Issue459Shim { return Issue459Shim{A: string(b)} }

type Issue459AllowNil struct {
	Map       *map[string]uint64  `msg:",allownil"`
	Buckets   *Issue459Buckets    `msg:",allownil"`
	Slice     *[]string           `msg:",allownil"`
	Bytes     *[]byte             `msg:",allownil"`
	DoubleMap **map[string]uint64 `msg:",allownil"`
	Uint      *uint64             `msg:",allownil"`
	Array     *[3]int             `msg:",allownil"`
	Shim      *Issue459Shim       `msg:",allownil"`
}

//msgp:tuple Issue459AllowNilTuple

type Issue459AllowNilTuple struct {
	Map       *map[string]uint64  `msg:",allownil"`
	Buckets   *Issue459Buckets    `msg:",allownil"`
	Slice     *[]string           `msg:",allownil"`
	Bytes     *[]byte             `msg:",allownil"`
	DoubleMap **map[string]uint64 `msg:",allownil"`
	Uint      *uint64             `msg:",allownil"`
	Array     *[3]int             `msg:",allownil"`
	Shim      *Issue459Shim       `msg:",allownil"`
}

// omitempty disables allownil on encode but not on decode.
type Issue459OmitAllowNil struct {
	Bytes **[]byte            `msg:",omitempty,allownil"`
	Slice **[]string          `msg:",omitempty,allownil"`
	Map   **map[string]uint64 `msg:",omitempty,allownil"`
}
