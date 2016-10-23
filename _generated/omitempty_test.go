package _generated

import (
	"bytes"
	"fmt"
	"github.com/shurcooL/go-goon"
	"github.com/tinylib/msgp/msgp"
	"testing"
)

var VerboseOmitTest = false

func TestMissingNilledOutWhenUnmarshallingOrDecodingNilIntoNestedStructs(t *testing.T) {

	underTest := "UnmarshalMsg"
	for passNo := 0; passNo < 2; passNo++ {
		if passNo == 1 {
			underTest = "UnmarshalMsg"
		}

		// passNo 0 => UnmarshalMsg
		// passNo 1 => DecodeMsg
		//
		// Given a tree of structs with element points
		// that is three levels deep, when omitempty fields
		// are omitted from the wire (msgpack data), we
		// should re-use the existing structs, maps, and slices,
		// and shrink them without re-allocation.

		s := TopNester{
			TopId:     43,
			Greetings: "greetings",
			Bullwinkle: &Rocky{
				Road: "Long and winding",
				Bugs: &Bunny{
					Carrots: []int{41, 8},
					Sayings: map[string]string{"whatsup": "doc"},
					BunnyId: 88,
				},
				Moose: &Moose{
					Trees:   []int{0, 1, 2, 3},
					Sayings: map[string]string{"one": "singular sensation"},
					Id:      2,
				},
			},
			MyIntArray:  [3]int{4, 5, 6},
			MyByteArray: [3]byte{1, 2, 3},
			MyMap:       map[string]string{"key": "to my heart"},
			MyArrayMap:  [3]map[string]string{{"key1": "to my heart1"}, {"key2": "to my heart2"}, {}},
		}

		// so pointers should not change upon decoding from nil
		pGreet := &s.Greetings
		pBull := &s.Bullwinkle
		pBugs := &s.Bullwinkle.Bugs
		pCarrots := &s.Bullwinkle.Bugs.Carrots
		pSay1 := &s.Bullwinkle.Bugs.Sayings
		pMoose := &s.Bullwinkle.Moose
		pTree := &s.Bullwinkle.Moose.Trees
		pSay2 := &s.Bullwinkle.Moose.Sayings

		if VerboseOmitTest {
			fmt.Printf("\n ======== BEGIN goon.Dump of TopNester *BEFORE* %s:\n", underTest)
			goon.Dump(s)
			fmt.Printf("\n ======== END goon.Dump of TopNester *BEFORE* %s\n", underTest)
		}
		var o []byte
		var err error
		nilMsg := []byte{0xc0}
		if passNo == 0 {
			if VerboseOmitTest {
				fmt.Printf("\n testing UnmarshalMsg ******************\n")
			}
			o, err = s.UnmarshalMsg(nilMsg)
			if err != nil {
				panic(err)
			}
		} else {
			if VerboseOmitTest {
				fmt.Printf("\n testing DecodeMsg ******************\n")
			}
			nilMsg = []byte{0xc0}
			dc := msgp.NewReader(bytes.NewBuffer(nilMsg))
			err := s.DecodeMsg(dc) // msgp: attempted to decode type "nil" with method for "map"
			if err != nil {
				panic(err)
			}
		}

		if VerboseOmitTest {
			fmt.Printf("\n ======== BEGIN goon.Dump of TopNester AFTER %s:\n", underTest)
			goon.Dump(s)
			fmt.Printf("\n ======== END goon.Dump of TopNester AFTER %s\n", underTest)
		}

		if len(o) != 0 {
			t.Fatal("nilMsg should have been consumed")
		}

		for i := range s.MyIntArray {
			if s.MyIntArray[i] != 0 {
				panic("shoud have been set to 0")
			}
		}
		for i := range s.MyByteArray {
			if s.MyIntArray[i] != 0 {
				panic("shoud have been set to 0")
			}
		}

		if pGreet != &s.Greetings {
			t.Fatal("pGreet differed from original")
		}
		if pBull != &s.Bullwinkle {
			t.Fatal("pBull differed from original")
		}
		if pBugs != &s.Bullwinkle.Bugs {
			t.Fatal("pBugs differed from original")
		}
		if pCarrots != &s.Bullwinkle.Bugs.Carrots {
			t.Fatal("pCarrots differed from original")
		}
		if pSay1 != &s.Bullwinkle.Bugs.Sayings {
			t.Fatal("pSay1 differed from original")
		}
		if pMoose != &s.Bullwinkle.Moose {
			t.Fatal("pMoose differed from original")
		}
		if pTree != &s.Bullwinkle.Moose.Trees {
			t.Fatal("pTree differed from original")
		}
		if pSay2 != &s.Bullwinkle.Moose.Sayings {
			t.Fatal("pSay2 differed from original")
		}

		// and yet, the maps and slices should be size 0,
		// the strings empty, the integers zeroed out.

		// TopNester
		if s.TopId != 0 {
			t.Fatal("s.TopId should be 0")
		}
		if len(s.Greetings) != 0 {
			t.Fatal("s.Grettings should be len 0")
		}
		if s.Bullwinkle == nil {
			t.Fatal("s.Bullwinkle should not be nil")
		}

		// TopNester.Bullwinkle
		if s.Bullwinkle.Bugs == nil {
			t.Fatal("s.Bullwinkle.Bugs should not be nil")
		}
		if len(s.Bullwinkle.Road) != 0 {
			t.Fatal("s.Bullwinkle.Road should be len 0")
		}
		if s.Bullwinkle.Moose == nil {
			t.Fatal("s.Bullwinkle.Moose should not be nil")
		}

		// TopNester.Bullwinkle.Bugs
		if len(s.Bullwinkle.Bugs.Carrots) != 0 {
			panic("this is wrong")
		}
		if len(s.Bullwinkle.Bugs.Sayings) != 0 {
			panic("this is wrong")
		}
		if s.Bullwinkle.Bugs.BunnyId != 0 {
			panic("this is wrong")
		}

		// TopNester.Bullwinkle.Moose
		if len(s.Bullwinkle.Moose.Trees) != 0 {
			panic("this is wrong")
		}
		if len(s.Bullwinkle.Moose.Sayings) != 0 {
			panic("this is wrong")
		}
		if s.Bullwinkle.Moose.Id != 0 {
			panic("this is wrong")
		}

		// MyMap should have been emptied out
		if len(s.MyMap) != 0 {
			panic("MyMap should have been len 0 after decode from nil")
		}

		// each map in MyArrayMap should have been emptied
		for k := range s.MyArrayMap {
			if len(s.MyArrayMap[k]) != 0 {
				panic(fmt.Sprintf("MyArrayMap[%v] should have been len 0 after decode from nil", k))
			}
		}

	} // end passNo loop
}
