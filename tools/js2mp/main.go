package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/tinylib/msgp/msgp"
)

func fatal(err error) {
	io.WriteString(os.Stderr, err.Error())
	os.Exit(1)
}

func main() {
	var out interface{}
	dec := json.NewDecoder(os.Stdin)
	wr := msgp.NewWriter(os.Stdout)
	var err error
	for err == nil {
		err = dec.Decode(&out)
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				fatal(err)
			}
		}
		err = wr.WriteIntf(out)
	}
	if err != nil {
		fatal(err)
	}
	err = wr.Flush()
	if err != nil {
		fatal(err)
	}
}
