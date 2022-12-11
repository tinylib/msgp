package main

import "github.com/tinylib/msgp/msgp"

type Example struct {
	I int
	S string
}

var buf [64]byte

func main() {
	e := Example{
		I: 1,
		S: "2",
	}

	b := buf[:0]
	b = msgp.AppendMapHeader(b, 2)
	b = msgp.AppendString(b, "i")
	b = msgp.AppendInt(b, e.I)
	b = msgp.AppendString(b, "s")
	b = msgp.AppendString(b, e.S)

	println("done, len(b):", len(b))
	for i := 0; i < len(b); i++ {
		print(b[i], " ")
	}
	println()
}
