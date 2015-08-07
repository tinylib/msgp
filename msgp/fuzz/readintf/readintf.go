package readintf

import "github.com/tinylib/msgp/msgp"

func Fuzz(data []byte) int {
	_, _, err := msgp.ReadIntfBytes(data)
	if err != nil {
		return 0
	}
	return 1
}
