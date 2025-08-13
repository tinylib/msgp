package _generated

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand/v2"
	"reflect"
	"strconv"
	"testing"

	"github.com/tinylib/msgp/msgp"
)

func TestShimmedKeys(t *testing.T) {
	rng := rand.New(rand.NewPCG(0, 0))
	// Generate a bunch of random maps
	for i := range 500 {
		var test MyMapKeyStruct
		if i != 0 {
			// Don't add anything to the first object
			test.Map = make(map[mapKey]int)
			test.MapS = make(map[mapKeyString]int)
			test.MapX = make(map[ExternalString]int)
			test.MapB = make(map[mapKeyBytes]int)
			for range rng.IntN(50) {
				test.Map[mapKey(rng.IntN(math.MaxInt32))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapS[mapKeyString(strconv.Itoa(rng.IntN(math.MaxInt32)))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				test.MapX[ExternalString(strconv.Itoa(rng.IntN(math.MaxInt32)))] = rng.IntN(100)
			}
			for range rng.IntN(50) {
				var tmp [8]byte
				binary.LittleEndian.PutUint64(tmp[:], rng.Uint64())
				test.MapB[tmp] = rng.IntN(100)
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
			var decoded, decoded2 MyMapKeyStruct
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
