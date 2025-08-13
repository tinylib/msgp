package _generated

//go:generate msgp -v

//msgp:maps autoshim

//msgp:replace ExternalString with:string
//msgp:replace ExternalInt with:int

type MyMapKeyStruct3 struct {
	MapString  map[string]int         `msg:",allownil"`
	MapString3 map[ExternalString]int `msg:",allownil"`
	MapString2 map[ExternalInt]int    `msg:",allownil"`
	MapFloat32 map[float32]int        `msg:",allownil"`
	MapFloat64 map[float64]int        `msg:",allownil"`
	MapUint    map[uint]int           `msg:",allownil"`
	MapUint8   map[uint8]int          `msg:",allownil"`
	MapUint16  map[uint16]int         `msg:",allownil"`
	MapUint32  map[uint32]int         `msg:",allownil"`
	MapUint64  map[uint64]int         `msg:",allownil"`
	MapByte    map[byte]int           `msg:",allownil"`
	MapInt     map[int]int            `msg:",allownil"`
	MapInt8    map[int8]int           `msg:",allownil"`
	MapInt16   map[int16]int          `msg:",allownil"`
	MapInt32   map[int32]int          `msg:",allownil"`
	MapInt64   map[int64]int          `msg:",allownil"`
	MapBool    map[bool]int           `msg:",allownil"`
	MapMapInt  map[int]map[int]int    `msg:",allownil"`

	// Maps with unsupported base types.
	MapArray      map[[4]byte]int    `msg:",allownil"` // Ignored
	MapArray4     map[[4]uint32]int  `msg:",allownil"` // Ignored
	MapComplex64  map[complex64]int  `msg:",allownil"` // Ignored
	MapComplex128 map[complex128]int `msg:",allownil"` // Ignored
}
