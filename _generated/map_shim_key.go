package _generated

import (
	"encoding/hex"
	"strconv"
)

//go:generate msgp -unexported -v

//msgp:maps shim

//msgp:shim mapKey as:string using:mapKeyToString/stringToMapKey witherr:true

type mapKey uint64

type mapKeyString string

//msgp:shim mapKeyBytes as:string using:hexEncode/hexDecode witherr:false

type mapKeyBytes [8]byte

//msgp:replace ExternalString with:string

type MyMapKeyStruct struct {
	Map      map[mapKey]int         `msg:",allownil"` // Keys are converted to strings via shim
	MapS     map[mapKeyString]int   `msg:",allownil"` // Keys are strings, with conversion
	MapX     map[ExternalString]int `msg:",allownil"` // Keys is an external type to the file, replaced with a cast
	MapB     map[mapKeyBytes]int    `msg:",allownil"` // Keys are bytes, converted to hex strings via shim
	MapUint  map[uint64]int         `msg:",allownil"` // Will be ignored (as current behavior)
	Original map[string]int         `msg:",allownil"` // Original map, for comparison purposes

	// Should all be ignored:
	MapFloat32    map[float32]int    `msg:",allownil"`
	MapFloat64    map[float64]int    `msg:",allownil"`
	MapComplex64  map[complex64]int  `msg:",allownil"`
	MapComplex128 map[complex128]int `msg:",allownil"`
	MapUint2      map[uint]int       `msg:",allownil"`
	MapUint8      map[uint8]int      `msg:",allownil"`
	MapUint16     map[uint16]int     `msg:",allownil"`
	MapUint32     map[uint32]int     `msg:",allownil"`
	MapUint64     map[uint64]int     `msg:",allownil"`
	MapByte       map[byte]int       `msg:",allownil"`
	MapInt        map[int]int        `msg:",allownil"`
	MapInt8       map[int8]int       `msg:",allownil"`
	MapInt16      map[int16]int      `msg:",allownil"`
	MapInt32      map[int32]int      `msg:",allownil"`
	MapInt64      map[int64]int      `msg:",allownil"`
	MapBool       map[bool]int       `msg:",allownil"`
}

// stringToMapKey and mapKeyToString are shim functions that convert between
// mapKey and string.
func stringToMapKey(s string) (mapKey, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return mapKey(v), nil
}

func mapKeyToString(k mapKey) string {
	return strconv.FormatUint(uint64(k), 10)
}

func hexEncode(b mapKeyBytes) string {
	return hex.EncodeToString(b[:])
}

func hexDecode(s string) mapKeyBytes {
	var b [8]byte
	_, err := hex.Decode(b[:], []byte(s))
	if err != nil {
		panic(err)
	}
	return b
}
