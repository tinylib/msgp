package main

import "github.com/tinylib/msgp/msgp"

type Example struct {
	I int
	S string
}

func main() {
	b := []byte{130, 161, 105, 1, 161, 115, 161, 50}

	e := Example{}

	sz, b, err := msgp.ReadMapHeaderBytes(b)
	if err != nil {
		panic(err)
	}
	if sz != 2 {
		panic("bad sz")
	}

	for i := 0; i < int(sz); i++ {
		var sb []byte
		var i32 int32
		var err error
		sb, b, err = msgp.ReadStringZC(b)
		switch string(sb) {
		case "i":
			i32, b, err = msgp.ReadInt32Bytes(b)
			if err != nil {
				panic(err)
			}
			e.I = int(i32)
		case "s":
			sb, b, err = msgp.ReadStringZC(b)
			if err != nil {
				panic(err)
			}
			e.S = string(sb)
		default:
			panic("unexpected field:" + string(sb))
		}
	}

	if len(b) != 0 {
		panic("unexpected extra")
	}

	if e.I != 1 || e.S != "2" {
		panic("not equal")
	}

	println("done")
}
