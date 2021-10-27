package main

import "encoding/json"

//go:generate msgp

type ID int64

type Example struct {
	ID
	I int    `msg:"i"`
	S string `msg:"s"`
}

func (id ID) GetID() int64 {
	return int64(id)
}

func main() {

	e := Example{
		ID: 1,
		I:  1,
		S:  "2",
	}

	println(e.GetID())

	b, err := e.MarshalMsg(nil)
	if err != nil {
		panic(err)
	}
	j, _ := json.Marshal(e)

	println("done, len(b):", len(b), string(b), string(j))
	e1 := Example{}

	json.Unmarshal(j, &e1)

	for i := 0; i < len(b); i++ {
		print(b[i], ",")
	}
	println()

}
