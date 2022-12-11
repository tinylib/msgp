package main

//go:generate msgp

type Example struct {
	I int    `msg:"i"`
	S string `msg:"s"`
}

func main() {
	b := []byte{130, 161, 105, 1, 161, 115, 161, 50}

	var e2 Example
	extra, err := e2.UnmarshalMsg(b)
	if err != nil {
		panic(err)
	}
	if len(extra) != 0 {
		panic("unexpected extra")
	}

	if e2.I != 1 || e2.S != "2" {
		panic("not equal")
	}

	println("done, len(b):", len(b))
}
