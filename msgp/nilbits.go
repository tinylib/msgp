package msgp

import (
//"fmt"
//"runtime/debug"
)

// NilBitsStack is a helper for Unmarshal
// methods to track where we are when
// deserializing from empty/nil/missing
// fields.
type NilBitsStack struct {
	// simulate getting nils on the wire
	AlwaysNil     bool
	LifoAlwaysNil []bool
	LifoBts       [][]byte
}

func (r *NilBitsStack) IsNil(bts []byte) bool {
	if r.AlwaysNil {
		return true
	}
	if len(bts) != 0 && bts[0] == mnil {
		return true
	}
	return false
}

// OnlyNilSlice is a slice that contains
// only the msgpack nil (0xc0) bytes.
var OnlyNilSlice = []byte{mnil}

// AlwaysNilString returns a string representation
// of the internal state of the NilBitsStack for
// debugging purposes.
func (r *NilBitsStack) AlwaysNilString() string {
	s := "bottom: "
	for _, v := range r.LifoAlwaysNil {
		if v {
			s += "T"
		} else {
			s += "f"
		}
	}
	return s
}

// PushAlwaysNil will set r.AlwaysNil to true
// and store bts on the internal stack.
func (r *NilBitsStack) PushAlwaysNil(bts []byte) []byte {
	//fmt.Printf("PushAlwaysNil(), pre we are '%v'. Called from: stack is '%v'\n", r.AlwaysNilString(), string(debug.Stack()))
	// save current state
	r.LifoBts = append(r.LifoBts, bts)
	r.LifoAlwaysNil = append(r.LifoAlwaysNil, r.AlwaysNil)
	// set reader r to always return nils
	r.AlwaysNil = true

	return OnlyNilSlice
}

// PopAlwaysNil pops the last []byte off the internal
// stack and returns it. If the stack is empty
// we panic.
func (r *NilBitsStack) PopAlwaysNil() (bts []byte) {
	//fmt.Printf("my NilBitsTack.PopAlwaysNil() called!! ...  stack is '%v'\n", string(debug.Stack()))
	//	defer func() {
	//		fmt.Printf("len of bts returned by PopAlwaysNil: %v, and debug string: '%v'\n",
	//			len(bts), r.AlwaysNilString())
	//	}()
	n := len(r.LifoAlwaysNil)
	//fmt.Printf("\n Reader.PopAlwaysNil() called! qlen = %d, '%s'\n",
	//	n, r.AlwaysNilString())
	if n == 0 {
		panic("PopAlwaysNil called on empty lifo")
	}
	a := r.LifoAlwaysNil[n-1]
	r.LifoAlwaysNil = r.LifoAlwaysNil[:n-1]
	r.AlwaysNil = a

	bts = r.LifoBts[n-1]
	r.LifoBts = r.LifoBts[:n-1]
	return bts
}

func ShowFound(found []bool) string {
	s := "["
	for i := range found {
		if found[i] {
			s += "1,"
		} else {
			s += "0,"
		}
	}
	s += "]"
	return s
}
