package main

//go:generate msgp

type Example struct {
	I int    `msg:"i"`
	S string `msg:"s"`
	// NOTE: floats omitted intentionally, many MCUs don't have
	// native support for it so the binary size bloats just due to
	// the software float implementation
}

func main() {
	e := Example{
		I: 1,
		S: "2",
	}
	b, err := e.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}

	var e2 Example
	extra, err := e2.UnmarshalMsg(b)
	if err != nil {
		panic(err)
	}
	if len(extra) != 0 {
		panic("unexpected extra")
	}

	if e.I != e2.I || e.S != e2.S {
		panic("not equal")
	}

	println("done, len(b):", len(b))
}
