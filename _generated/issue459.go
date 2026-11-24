package _generated

//go:generate msgp

//msgp:maps autoshim
//msgp:noduplicates

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
}

type Issue459Inner struct {
	A int
}

type Issue459AllowNil struct {
	Map       *map[string]uint64  `msg:",allownil"`
	Buckets   *Issue459Buckets    `msg:",allownil"`
	Slice     *[]string           `msg:",allownil"`
	Bytes     *[]byte             `msg:",allownil"`
	DoubleMap **map[string]uint64 `msg:",allownil"`
	Uint      *uint64             `msg:",allownil"`
	Array     *[3]int             `msg:",allownil"`
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
}
