package main

import (
	"io"
	"os"

	"github.com/tinylib/msgp/msgp"
)

func main() {
	_, err := msgp.NewReader(os.Stdin).WriteToJSON(os.Stdout)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
}
