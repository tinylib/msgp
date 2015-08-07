// +build gofuzz

package skipfuzz

import "github.com/tinylib/msgp/msgp"

func Fuzz(raw []byte) int {
	var err error
	for err == nil && len(raw) > 0 {
		raw, err = msgp.Skip(raw)
	}
	if err != nil {
		return 0
	}
	return 1
}
