package _generated

import (
	"encoding/hex"
	"time"
)

//go:generate msgp -unexported -v

//msgp:maps binkeys

type ArrayMapKey [4]byte

//msgp:replace ExternalArr with:[4]byte

//msgp:replace ExternalString with:string

type mapKeyBytes2 [8]byte

//msgp:shim mapKeyBytes2 as:string using:hexEncode2/hexDecode2 witherr:false

type mapKeyShimmed time.Duration

//msgp:shim mapKeyShimmed as:[]byte using:durEncode/durDecode witherr:true

type MyStringType string

type MyMapKeyStruct2 struct {
	MapString     map[string]int         `msg:",allownil"`
	MapString2    map[MyStringType]int   `msg:",allownil"`
	MapString3    map[ExternalString]int `msg:",allownil"`
	MapString4    map[mapKeyBytes2]int   `msg:",allownil"`
	MapFloat32    map[float32]int        `msg:",allownil"`
	MapFloat64    map[float64]int        `msg:",allownil"`
	MapComplex64  map[complex64]int      `msg:",allownil"`
	MapComplex128 map[complex128]int     `msg:",allownil"`
	MapUint       map[uint]int           `msg:",allownil"`
	MapUint8      map[uint8]int          `msg:",allownil"`
	MapUint16     map[uint16]int         `msg:",allownil"`
	MapUint32     map[uint32]int         `msg:",allownil"`
	MapUint64     map[uint64]int         `msg:",allownil"`
	MapByte       map[byte]int           `msg:",allownil"`
	MapInt        map[int]int            `msg:",allownil"`
	MapInt8       map[int8]int           `msg:",allownil"`
	MapInt16      map[int16]int          `msg:",allownil"`
	MapInt32      map[int32]int          `msg:",allownil"`
	MapInt64      map[int64]int          `msg:",allownil"`
	MapBool       map[bool]int           `msg:",allownil"`
	MapMapInt     map[int]map[int]int    `msg:",allownil"`

	// Maps with array keys
	MapArray  map[[4]byte]int     `msg:",allownil"`
	MapArray2 map[ArrayMapKey]int `msg:",allownil"`
	MapArray3 map[ExternalArr]int `msg:",allownil"`
	MapArray4 map[[4]uint32]int   `msg:",allownil"`

	// Maps with shimmed types
	MapDuration map[mapKeyShimmed]int `msg:",allownil"`
}

func hexEncode2(b mapKeyBytes2) string {
	return hex.EncodeToString(b[:])
}

func hexDecode2(s string) mapKeyBytes2 {
	var b [8]byte
	_, err := hex.Decode(b[:], []byte(s))
	if err != nil {
		panic(err)
	}
	return b
}

func durEncode(v mapKeyShimmed) []byte {
	return []byte(time.Duration(v).String())
}

func durDecode(b []byte) (mapKeyShimmed, error) {
	v, err := time.ParseDuration(string(b))
	return mapKeyShimmed(v), err
}
