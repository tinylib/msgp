package main

import "testing"

func TestIssue430Nested(t *testing.T) {
	mainFilename, err := generate(t, issue430Nested)
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	goExec(t, mainFilename, false) // exec go run
	goExec(t, mainFilename, true)  // exec go test
}

var issue430Nested = `package main

import (
	"fmt"
	"os"
)

//go:generate msgp

//msgp:vartuple Bag
type Bag struct {
	ItemID  uint32
	ItemNum uint32
}

//msgp:vartuple User
type User struct {
	Name string
	Bag  []*Bag
}

func main() {
	u := User{Name: "user", Bag: []*Bag{{1, 2}, {3, 4}}}
	data, err := u.MarshalMsg(nil)
	if err != nil {
		fmt.Println("User MarshalMsg failed:", err)
		os.Exit(1)
	}

	var user User
	if _, err := user.UnmarshalMsg(data); err != nil {
		fmt.Println("User UnmarshalMsg failed:", err)
		os.Exit(1)
	}
	if user.Name != "user" ||
		user.Bag[0] == nil ||
		user.Bag[0].ItemID != 1 ||
		user.Bag[0].ItemNum != 2 ||
		user.Bag[1] == nil ||
		user.Bag[1].ItemID != 3 ||
		user.Bag[1].ItemNum != 4 {
		fmt.Println("User UnmarshalMsg bad user:", user)
		os.Exit(1)
	}
}
`
