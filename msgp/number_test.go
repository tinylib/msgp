package msgp

import (
	"testing"
)

func TestNumber(t *testing.T) {

	n := Number{}

	if n.Type() != IntType {
		t.Errorf("expect zero-value type to be %s; got %s", IntType, n.Type())
	}

	n.AsInt(248)
	i, ok := n.Int()
	if !ok || i != 248 || n.Type() != IntType {
		t.Errorf("%d in; %d out!", 248, i)
	}

	n.AsFloat64(3.141)
	f, ok := n.Float()
	if !ok || f != 3.141 || n.Type() != Float64Type {
		t.Errorf("%f in; %f out!", 3.141, f)
	}

	n.AsUint(40000)
	u, ok := n.Uint()
	if !ok || u != 40000 || n.Type() != UintType {
		t.Errorf("%d in; %d out!", 40000, u)
	}

	nums := []interface{}{
		float64(3.14159),
		int64(-29081),
		uint64(90821983),
		float32(3.141),
	}

	var dat []byte
	for _, n := range nums {
		dat, _ = AppendIntf(dat, n)
	}

	for range nums {
		var err error
		dat, err = n.UnmarshalMsg(dat)
		if err != nil {
			t.Fatal("decode error:", err)
		}
	}

}
