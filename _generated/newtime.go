package _generated

import "time"

//go:generate msgp -v

//msgp:newtime
//msgp:timezone local

type NewTime struct {
	T     time.Time
	Array []time.Time
	Map   map[string]time.Time
}

func (t1 NewTime) Equal(t2 NewTime) bool {
	if !t1.T.Equal(t2.T) {
		return false
	}
	if len(t1.Array) != len(t2.Array) {
		return false
	}
	for i := range t1.Array {
		if !t1.Array[i].Equal(t2.Array[i]) {
			return false
		}
	}
	if len(t1.Map) != len(t2.Map) {
		return false
	}
	for k, v := range t1.Map {
		if !t2.Map[k].Equal(v) {
			return false
		}
	}
	return true
}
