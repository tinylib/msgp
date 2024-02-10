package _generated

import (
	"time"
)

//go:generate msgp

//msgp:shim time.Time as:string using:timetostr/strtotime
type T struct {
	T time.Time
}

func timetostr(t time.Time) string {
	return t.Format(time.RFC3339)
}

func strtotime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
