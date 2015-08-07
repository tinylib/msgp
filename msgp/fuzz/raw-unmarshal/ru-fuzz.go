package ru

import "github.com/tinylib/msgp/msgp"

func Fuzz(data []byte) int {
	var r msgp.Raw
	_, err := r.UnmarshalMsg(data)
	if err != nil {
		return 0
	}
	return 1
}
