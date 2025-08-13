package _generated

import (
	"bytes"
	"math"
	"math/rand/v2"
	"reflect"
	"strconv"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestBinKeys(t *testing.T) {
	rng := rand.New(rand.NewPCG(0, 0))

	// Generate a bunch of random maps
	for i := range 500 {
		var test MyMapKeyStruct2
		if i != 0 {
			// Don't add anything to the first object
			test.MapString = make(map[string]int)
			test.MapString2 = make(map[MyStringType]int)
			test.MapString3 = make(map[ExternalString]int)
			test.MapFloat32 = make(map[float32]int)
			test.MapFloat64 = make(map[float64]int)
			test.MapComplex64 = make(map[complex64]int)
			test.MapComplex128 = make(map[complex128]int)
			test.MapUint = make(map[uint]int)
			test.MapUint8 = make(map[uint8]int)
			test.MapUint16 = make(map[uint16]int)
			test.MapUint32 = make(map[uint32]int)
			test.MapUint64 = make(map[uint64]int)
			test.MapByte = make(map[byte]int)
			test.MapInt = make(map[int]int)
			test.MapInt8 = make(map[int8]int)
			test.MapInt16 = make(map[int16]int)
			test.MapInt32 = make(map[int32]int)
			test.MapInt64 = make(map[int64]int)
			test.MapBool = make(map[bool]int)
			test.MapArray = make(map[[4]byte]int)
			test.MapArray2 = make(map[ArrayMapKey]int)
			test.MapArray3 = make(map[ExternalArr]int)
			test.MapArray4 = make(map[[4]uint32]int)
			test.MapDuration = make(map[mapKeyShimmed]int)
			test.MapMapInt = make(map[int]map[int]int)

			for range rng.IntN(50) {
				test.MapString[string(strconv.Itoa(rng.IntN(math.MaxInt32)))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapString2[MyStringType(strconv.Itoa(rng.IntN(math.MaxInt32)))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapString3[ExternalString(strconv.Itoa(rng.IntN(math.MaxInt32)))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapFloat32[float32(rng.Float32())] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapFloat64[rng.Float64()] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				c := complex(float32(rng.Float32()), float32(rng.Float32()))
				test.MapComplex64[c] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				c := complex(rng.Float64(), rng.Float64())
				test.MapComplex128[c] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapUint[uint(rng.Uint64())] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapUint8[uint8(rng.Uint64())] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapUint16[uint16(rng.Uint64())] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapUint32[rng.Uint32()] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapUint64[rng.Uint64()] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapByte[byte(uint8(rng.Uint64()))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapInt[rng.IntN(math.MaxInt32)] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapInt8[int8(rng.IntN(int(math.MaxInt8)+1))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapInt16[int16(rng.IntN(int(math.MaxInt16)+1))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapInt32[int32(rng.IntN(int(math.MaxInt32)))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				// Use only non-negative values to stay within IntN capabilities
				test.MapInt64[int64(rng.IntN(int(math.MaxInt32)))] = rng.IntN(100)
			}
			if rng.IntN(2) == 0 {
				test.MapBool[true] = rng.IntN(100)
			}
			if rng.IntN(2) == 0 {
				test.MapBool[false] = rng.IntN(100)
			}

			for range rng.IntN(50) {
				var k [4]byte
				for i := 0; i < 4; i++ {
					k[i] = uint8(rng.Uint64())
				}
				test.MapArray[k] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				var k ArrayMapKey
				for i := 0; i < 4; i++ {
					k[i] = uint8(rng.Uint64())
				}
				test.MapArray2[k] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				var k ExternalArr
				for i := 0; i < 4; i++ {
					k[i] = uint8(rng.Uint64())
				}
				test.MapArray3[k] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				var k [4]uint32
				for i := 0; i < 4; i++ {
					k[i] = rng.Uint32()
				}
				test.MapArray4[k] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapDuration[mapKeyShimmed(rng.IntN(math.MaxInt32))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				dst := make(map[int]int, 50)
				test.MapMapInt[rng.Int()] = dst
				for range rng.IntN(50) {
					dst[rng.Int()] = rng.IntN(100)
				}
			}
		}
		var encoded [][]byte
		b, err := test.MarshalMsg(nil)
		if err != nil {
			t.Fatal(err)
		}
		encoded = append(encoded, b)
		var buf bytes.Buffer
		en := msgp.NewWriter(&buf)
		err = test.EncodeMsg(en)
		if err != nil {
			t.Fatal(err)
		}
		err = en.Flush()
		if err != nil {
			t.Fatal(err)
		}
		encoded = append(encoded, buf.Bytes())
		for _, enc := range encoded {
			var decoded, decoded2 MyMapKeyStruct2
			_, err = decoded.UnmarshalMsg(enc)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(&decoded, &test) {
				t.Errorf("decoded != test")
			}
			dec := msgp.NewReader(bytes.NewReader(enc))
			err = decoded2.DecodeMsg(dec)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(&decoded2, &test) {
				t.Errorf("decoded2 != test")
			}
		}
	}
}
