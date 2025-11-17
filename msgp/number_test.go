package msgp

import (
	"bytes"
	"math"
	"testing"
)

func TestNumber(t *testing.T) {
	n := Number{}

	if n.Type() != IntType {
		t.Errorf("expected zero-value type to be %s; got %s", IntType, n.Type())
	}

	if n.String() != "0" {
		t.Errorf("expected Number{}.String() to be \"0\" but got %q", n.String())
	}

	n.AsInt(248)
	i, ok := n.Int()
	if !ok || i != 248 || n.Type() != IntType || n.String() != "248" {
		t.Errorf("%d in; %d out!", 248, i)
	}

	n.AsFloat64(3.141)
	f, ok := n.Float()
	if !ok || f != 3.141 || n.Type() != Float64Type || n.String() != "3.141" {
		t.Errorf("%f in; %f out!", 3.141, f)
	}

	n.AsUint(40000)
	u, ok := n.Uint()
	if !ok || u != 40000 || n.Type() != UintType || n.String() != "40000" {
		t.Errorf("%d in; %d out!", 40000, u)
	}

	nums := []any{
		float64(3.14159),
		int64(-29081),
		uint64(90821983),
		float32(3.141),
	}

	var dat []byte
	var buf bytes.Buffer
	wr := NewWriter(&buf)
	for _, n := range nums {
		dat, _ = AppendIntf(dat, n)
		wr.WriteIntf(n)
	}
	wr.Flush()

	mout := make([]Number, len(nums))
	dout := make([]Number, len(nums))

	rd := NewReader(&buf)
	unm := dat
	for i := range nums {
		var err error
		unm, err = mout[i].UnmarshalMsg(unm)
		if err != nil {
			t.Fatal("unmarshal error:", err)
		}
		err = dout[i].DecodeMsg(rd)
		if err != nil {
			t.Fatal("decode error:", err)
		}
		if mout[i] != dout[i] {
			t.Errorf("for %#v, got %#v from unmarshal and %#v from decode", nums[i], mout[i], dout[i])
		}
	}

	buf.Reset()
	var odat []byte
	for i := range nums {
		var err error
		odat, err = mout[i].MarshalMsg(odat)
		if err != nil {
			t.Fatal("marshal error:", err)
		}
		err = dout[i].EncodeMsg(wr)
		if err != nil {
			t.Fatal("encode error", err)
		}
	}
	wr.Flush()

	if !bytes.Equal(dat, odat) {
		t.Errorf("marshal: expected output %#v; got %#v", dat, odat)
	}

	if !bytes.Equal(dat, buf.Bytes()) {
		t.Errorf("encode: expected output %#v; got %#v", dat, buf.Bytes())
	}
}

func TestConv32(t *testing.T) {
	for i := -1 << 24; i < 1<<25; i++ {
		x := math.Float32bits(float32(i))
		exp := int(x>>23)&255 - 127
		mant := x & ((1 << 23) - 1)
		isExact := false
		if exp >= 23 || (exp == -127 && mant == 0) {
			isExact = true
		}
		// Only exp >= 0 can be exact integer
		if exp >= 0 && exp <= 23 {
			mantissaMask := uint32((1 << (23 - exp)) - 1)
			isExact = (mant & mantissaMask) == 0
		}

		n := Number{bits: uint64(x), typ: Float32Type}
		got := n.isExactInt()
		if got != isExact {
			t.Errorf("n.IsExactInt(): got %t, want %t", got, isExact)
		}
		n = Number{bits: uint64(math.Float32bits(float32(i) + 0.1)), typ: Float32Type}
		got = n.isExactInt()
		if got != false && i > -2097152 && i < 2097152 {
			val, ok := n.Float()
			t.Fatalf("n.IsExactInt(%f): got %t, want %t, ok: %v", val, got, false, ok)
		}
	}
}

func TestConv64(t *testing.T) {
	for i := -1 << 30; i < 1<<30; i++ {
		if testing.Short() {
			i += 8
		}
		x := math.Float64bits(float64(i))
		exp := int(x>>52)&2047 - 1023
		mant := x & ((1 << 52) - 1)

		isExact := false
		if exp >= 52 || (exp == -1023 && mant == 0) {
			isExact = true
		}
		// Only exp >= 0 can be exact integer
		if exp >= 0 && exp <= 52 {
			mantissaMask := uint64((1 << (52 - exp)) - 1)
			isExact = (mant & mantissaMask) == 0
		}
		n := Number{bits: x, typ: Float64Type}
		got := n.isExactInt()
		if got != isExact {
			t.Errorf("n.IsExactInt(): got %t, want %t", got, isExact)
		}
		if !got {
			t.Fatal(i)
		}
		n = Number{bits: math.Float64bits(float64(i) + 0.1), typ: Float64Type}
		got = n.isExactInt()
		if got != false {
			val, ok := n.Float()
			t.Fatalf("n.IsExactInt(%f): got %t, want %t, ok: %v", val, got, false, ok)
		}
	}
}
