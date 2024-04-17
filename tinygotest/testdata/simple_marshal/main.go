package main

//go:generate msgp

type Example struct {
	I int    `msg:"i"`
	S string `msg:"s"`
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

	println("done, len(b):", len(b))
	for i := 0; i < len(b); i++ {
		print(b[i], " ")
	}
	println()
}
