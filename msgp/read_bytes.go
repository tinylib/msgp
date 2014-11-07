package msgp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"time"
	"unsafe"
)

var (
	// ErrShortBytes is returned when the
	// slice being decoded is too short to
	// contain the contents of the message
	ErrShortBytes = errors.New("msgp: too few bytes left to read object")
)

func IsNil(b []byte) bool {
	if len(b) > 0 && b[0] == mnil {
		return true
	}
	return false
}

func ReadMapHeaderBytes(b []byte) (sz uint32, o []byte, err error) {
	l := len(b)
	if l < 1 {
		err = ErrShortBytes
		return
	}

	lead := b[0]
	if isfixmap(lead) {
		sz = uint32(rfixmap(lead))
		o = b[1:]
		return
	}

	switch lead {
	case mnil:
		return 0, b, ErrNil

	case mmap16:
		if l < 3 {
			err = ErrShortBytes
			return
		}
		sz = uint32(binary.BigEndian.Uint16(b[1:]))
		o = b[3:]
		return

	case mmap32:
		if l < 5 {
			err = ErrShortBytes
			return
		}
		sz = binary.BigEndian.Uint32(b[1:])
		o = b[5:]
		return

	default:
		err = TypeError{Type: MapType, Prefix: lead}
		return
	}
}

func ReadMapKeyZC(b []byte) ([]byte, []byte, error) {
	if len(b) < 2 {
		return nil, b, ErrShortBytes
	}
	k := getKind(b[0])
	if k == kstring {
		return ReadStringZC(b)
	} else if k == kbytes {
		return ReadBytesZC(b)
	}
	return nil, b, fmt.Errorf("msgp: %q not convertible to map key (string)", k)
}

func ReadArrayHeaderBytes(b []byte) (sz uint32, o []byte, err error) {
	if len(b) < 1 {
		return 0, nil, ErrShortBytes
	}
	lead := b[0]
	if isfixarray(lead) {
		sz = uint32(rfixarray(lead))
		o = b[1:]
		return
	}

	switch lead {
	case mnil:
		return 0, b, ErrNil

	case marray16:
		if len(b) < 3 {
			err = ErrShortBytes
			return
		}
		sz = uint32(binary.BigEndian.Uint16(b[1:]))
		o = b[3:]
		return

	case marray32:
		if len(b) < 5 {
			err = ErrShortBytes
			return
		}
		sz = binary.BigEndian.Uint32(b[1:])
		o = b[5:]
		return

	default:
		err = TypeError{Type: ArrayType, Prefix: lead}
		return
	}
}

func ReadNilBytes(b []byte) ([]byte, error) {
	if len(b) < 1 {
		return nil, ErrShortBytes
	}
	if b[0] != mnil {
		return b, TypeError{Type: NilType, Prefix: b[0]}
	}
	return b[1:], nil
}

func ReadFloat64Bytes(b []byte) (f float64, o []byte, err error) {
	if len(b) < 9 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfloat64 {
		err = TypeError{Type: Float64Type, Prefix: b[0]}
		return
	}

	f = *(*float64)(unsafe.Pointer(&b[1]))
	o = b[9:]
	return
}

func ReadFloat32Bytes(b []byte) (f float32, o []byte, err error) {
	if len(b) < 5 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfloat32 {
		err = TypeError{Type: Float32Type, Prefix: b[0]}
		return
	}

	f = *(*float32)(unsafe.Pointer(&b[1]))
	o = b[5:]
	return
}

func ReadBoolBytes(b []byte) (v bool, o []byte, err error) {
	if len(b) < 1 {
		err = ErrShortBytes
		return
	}

	switch b[0] {
	case mtrue:
		v = true
		o = b[1:]
		return

	case mfalse:
		v = false
		o = b[1:]
		return

	default:
		err = TypeError{Type: BoolType, Prefix: b[0]}
		return

	}
}

func ReadInt64Bytes(b []byte) (i int64, o []byte, err error) {
	l := len(b)
	if l < 1 {
		return 0, nil, ErrShortBytes
	}

	lead := b[0]

	if isfixint(lead) {
		i = int64(rfixint(lead))
		o = b[1:]
		return
	}

	if isnfixint(lead) {
		i = int64(rnfixint(lead))
		o = b[1:]
		return
	}

	switch lead {
	case mnil:
		err = ErrNil
		return

	case mint8:
		if l < 2 {
			err = ErrShortBytes
			return
		}
		i = int64(int8(b[1]))
		o = b[2:]
		return

	case mint16:
		if l < 3 {
			err = ErrShortBytes
			return
		}
		i = int64((int16(b[1]) << 8) | int16(b[2]))
		o = b[3:]
		return

	case mint32:
		if l < 5 {
			err = ErrShortBytes
			return
		}
		i = int64((int32(b[1]) << 24) | (int32(b[2]) << 16) | (int32(b[3]) << 8) | int32(b[4]))
		o = b[5:]
		return

	case mint64:
		if l < 9 {
			err = ErrShortBytes
			return
		}
		i |= (int64(b[1]) << 56)
		i |= (int64(b[2]) << 48)
		i |= (int64(b[3]) << 40)
		i |= (int64(b[4]) << 36)
		i |= (int64(b[5]) << 24)
		i |= (int64(b[6]) << 16)
		i |= (int64(b[7]) << 8)
		i |= (int64(b[8]))
		o = b[9:]
		return

	default:
		err = TypeError{Type: IntType, Prefix: lead}
		return
	}
}

func ReadInt32Bytes(b []byte) (int32, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, o, fmt.Errorf("msgp: %d overflows int32", i)
	}
	return int32(i), o, err
}

func ReadInt16Bytes(b []byte) (int16, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt16 || i < math.MinInt16 {
		return 0, o, fmt.Errorf("msgp: %d overflows int16", i)
	}
	return int16(i), o, err
}

func ReadInt8Bytes(b []byte) (int8, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt8 || i < math.MinInt8 {
		return 0, o, fmt.Errorf("msgp: %d overflows int8", i)
	}
	return int8(i), o, err
}

func ReadIntBytes(b []byte) (int, []byte, error) {
	var v int
	i, o, err := ReadInt64Bytes(b)
	if unsafe.Sizeof(v) == 32 {
		if i > math.MaxInt32 || i < math.MinInt32 {
			return 0, o, fmt.Errorf("msgp: %d overflows int(32)", i)
		}
	}
	v = int(i)
	return v, o, err
}

func ReadUint64Bytes(b []byte) (u uint64, o []byte, err error) {
	l := len(b)
	if l < 1 {
		return 0, nil, ErrShortBytes
	}

	lead := b[0]
	if isfixint(lead) {
		u = uint64(rfixint(lead))
		o = b[1:]
		return
	}

	switch lead {
	case muint8:
		if l < 2 {
			err = ErrShortBytes
			return
		}
		u = uint64(b[1])
		o = b[2:]
		return

	case muint16:
		if l < 3 {
			err = ErrShortBytes
			return
		}
		u = uint64(binary.BigEndian.Uint16(b[1:]))
		o = b[3:]
		return

	case muint32:
		if l < 5 {
			err = ErrShortBytes
			return
		}
		u = uint64(binary.BigEndian.Uint32(b[1:]))
		o = b[5:]
		return

	case muint64:
		if l < 9 {
			err = ErrShortBytes
			return
		}
		u = binary.BigEndian.Uint64(b[1:])
		o = b[9:]
		return

	default:
		err = TypeError{Type: UintType, Prefix: lead}
		return
	}
}

func ReadUint32Bytes(b []byte) (uint32, []byte, error) {
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint32 {
		return 0, nil, fmt.Errorf("msgp: %d overflows uint32", v)
	}
	return uint32(v), o, err
}

func ReadUint16Bytes(b []byte) (uint16, []byte, error) {
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint16 {
		return 0, nil, fmt.Errorf("msgp: %d overflows uint16", v)
	}
	return uint16(v), o, err
}

func ReadUint8Bytes(b []byte) (uint8, []byte, error) {
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint8 {
		return 0, nil, fmt.Errorf("msgp: %d overflows uint8", v)
	}
	return uint8(v), o, err
}

func ReadUintBytes(b []byte) (uint, []byte, error) {
	var l uint
	v, o, err := ReadUint64Bytes(b)
	if unsafe.Sizeof(l) == 32 && v > math.MaxUint32 {
		return 0, nil, fmt.Errorf("msgp: %d overflows uint(32)", v)
	}
	l = uint(v)
	return l, o, err
}

func ReadByteBytes(b []byte) (byte, []byte, error) {
	return ReadUint8Bytes(b)
}

func ReadBytesBytes(b []byte, scratch []byte) (v []byte, o []byte, err error) {
	return readBytesBytes(b, scratch, false)
}

func readBytesBytes(b []byte, scratch []byte, zc bool) (v []byte, o []byte, err error) {
	l := len(b)
	if l < 1 {
		return nil, nil, ErrShortBytes
	}

	lead := b[0]
	var read int
	switch lead {
	case mbin8:
		if l < 2 {
			err = ErrShortBytes
			return
		}

		read = int(b[1])
		b = b[2:]

	case mbin16:
		if l < 3 {
			err = ErrShortBytes
			return
		}
		read = int(binary.BigEndian.Uint16(b[1:]))
		b = b[3:]

	case mbin32:
		if l < 5 {
			err = ErrShortBytes
			return
		}
		read = int(binary.BigEndian.Uint32(b[1:]))
		b = b[5:]

	default:
		err = TypeError{Type: BinType, Prefix: lead}
		return
	}

	if len(b) < read {
		err = ErrShortBytes
		return
	}

	// zero-copy
	if zc {
		v = b[0:read]
		o = b[read:]
		return
	}

	if cap(scratch) >= read {
		v = scratch[0:read]
	} else {
		v = make([]byte, read)
	}

	n := copy(v, b)
	o = b[n:]
	return
}

// ReadBytesZC extracts the messagepack-encoded
// binary field without copying. The returned []byte
// points to the same memory as the input slice.
func ReadBytesZC(b []byte) (v []byte, o []byte, err error) {
	return readBytesBytes(b, nil, true)
}

// ReadStringZC reads a messagepack string field
// without copying. The returned []byte points
// to the same memory as the input slice.
func ReadStringZC(b []byte) (v []byte, o []byte, err error) {
	l := len(b)
	if l < 1 {
		return nil, nil, ErrShortBytes
	}

	lead := b[0]
	var read int

	if isfixstr(lead) {
		read = int(rfixstr(lead))
		b = b[1:]
	} else {
		switch lead {
		case mstr8:
			if l < 2 {
				err = ErrShortBytes
				return
			}
			read = int(b[1])
			b = b[2:]

		case mstr16:
			if l < 3 {
				err = ErrShortBytes
				return
			}
			read = int(binary.BigEndian.Uint16(b[1:]))
			b = b[3:]

		case mstr32:
			if l < 5 {
				err = ErrShortBytes
				return
			}
			read = int(binary.BigEndian.Uint32(b[1:]))
			b = b[5:]

		default:
			err = TypeError{Type: StrType, Prefix: lead}
			return
		}
	}

	if len(b) < read {
		err = ErrShortBytes
		return
	}

	v = b[0:read]
	o = b[read:]
	return
}

func ReadStringBytes(b []byte) (string, []byte, error) {
	v, o, err := ReadStringZC(b)
	return string(v), o, err
}

func ReadComplex128Bytes(b []byte) (c complex128, o []byte, err error) {
	if len(b) < 18 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext16 {
		err = fmt.Errorf("msgp: unexpected extension length (%x) for complex128", b[0])
		return
	}

	if b[1] != Complex128Extension {
		err = fmt.Errorf("msgp: unexpected byte %x for complex128", b[1])
		return
	}
	c = *(*complex128)(unsafe.Pointer(&b[2]))
	o = b[18:]
	return
}

func ReadComplex64Bytes(b []byte) (c complex64, o []byte, err error) {
	if len(b) < 10 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext8 {
		err = fmt.Errorf("msgp: unexpected extension length (%x) for complex64", b[0])
		return
	}

	if b[1] != Complex64Extension {
		err = fmt.Errorf("msgp: unexpected byte %x for complex64", b[1])
		return
	}
	c = *(*complex64)(unsafe.Pointer(&b[2]))
	o = b[10:]
	return
}

func ReadTimeBytes(b []byte) (t time.Time, o []byte, err error) {
	if len(b) < 18 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext16 {
		err = fmt.Errorf("msgp: unexpected extension length (%x) for time.Time", b[0])
		return
	}

	if int8(b[1]) != TimeExtension {
		err = fmt.Errorf("msgp: unexpected byte %x for time extension", b[1])
		return
	}

	err = t.UnmarshalBinary(b[2:17])
	o = b[18:]
	return
}

func ReadMapStrIntfBytes(b []byte, old map[string]interface{}) (v map[string]interface{}, o []byte, err error) {
	var sz uint32
	o = b
	sz, o, err = ReadMapHeaderBytes(o)

	if err != nil {
		return
	}

	if old != nil {
		for key, _ := range old {
			delete(old, key)
		}
		v = old
	} else {
		v = make(map[string]interface{}, int(sz))
	}

	for z := uint32(0); z < sz; z++ {
		if len(o) < 1 {
			err = ErrShortBytes
			return
		}
		var key string
		switch getKind(o[0]) {
		case kstring:
			key, o, err = ReadStringBytes(o)
			if err != nil {
				return
			}
		case kbytes:
			var bts []byte
			bts, o, err = ReadBytesZC(o)
			if err != nil {
				return
			}
			key = string(bts)
		default:
			err = TypeError{Type: StrType, Prefix: o[0]}
			return
		}
		var val interface{}
		val, o, err = ReadIntfBytes(o)
		if err != nil {
			return
		}
		v[key] = val
	}
	return
}

func ReadIntfBytes(b []byte) (i interface{}, o []byte, err error) {
	if len(b) < 1 {
		err = ErrShortBytes
		return
	}

	k := getKind(b[0])

	switch k {
	case kmap:
		i, o, err = ReadMapStrIntfBytes(b, nil)
		return

	case karray:
		var sz uint32
		sz, b, err = ReadArrayHeaderBytes(b)
		if err != nil {
			return
		}
		j := make([]interface{}, int(sz))
		i = j
		for d := range j {
			j[d], b, err = ReadIntfBytes(b)
			if err != nil {
				return
			}
		}
		return

	case kfloat32:
		i, o, err = ReadFloat32Bytes(b)
		return

	case kfloat64:
		i, o, err = ReadFloat64Bytes(b)
		return

	case kint:
		i, o, err = ReadInt64Bytes(b)
		return

	case kuint:
		i, o, err = ReadUint64Bytes(b)
		return

	case kbool:
		i, o, err = ReadBoolBytes(b)
		return

	case kextension:
		if len(b) < 3 {
			err = ErrShortBytes
			return
		}
		var t int8
		t, err = peekExtension(b)
		if err != nil {
			return
		}
		switch t {
		case TimeExtension:
			i, o, err = ReadTimeBytes(b)
			return
		case Complex128Extension:
			i, o, err = ReadComplex128Bytes(b)
			return
		case Complex64Extension:
			i, o, err = ReadComplex64Bytes(b)
			return
		}
		// use a user-defined extension,
		// if it's been registered
		f, ok := extensionReg[t]
		if ok {
			e := f()
			o, err = ReadExtensionBytes(b, e)
			i = e
			return
		}
		// last resort is a raw extension
		e := RawExtension{}
		e.Type = int8(t)
		o, err = ReadExtensionBytes(b, &e)
		i = &e
		return

	case knull:
		o, err = ReadNilBytes(b)
		return

	case kbytes:
		i, o, err = ReadBytesBytes(b, nil)
		return

	case kstring:
		i, o, err = ReadStringBytes(b)
		return

	default:
		err = InvalidPrefixError(b[0])
		return
	}
}

func getKind(v byte) kind {

	// fixed encoding
	switch {
	case isfixmap(v):
		return kmap
	case isfixarray(v):
		return karray
	case isfixint(v), isnfixint(v):
		return kint
	case isfixstr(v):
		return kstring
	}

	// var encoding
	switch v {
	case mmap16, mmap32:
		return kmap
	case marray16, marray32:
		return karray
	case mfloat32:
		return kfloat32
	case mfloat64:
		return kfloat64
	case mint8, mint16, mint32, mint64:
		return kint
	case muint8, muint16, muint32, muint64:
		return kuint
	case mfixext1, mfixext2, mfixext4, mfixext8, mfixext16, mext8, mext16, mext32:
		return kextension
	case mstr8, mstr16, mstr32:
		return kstring
	case mbin8, mbin16, mbin32:
		return kbytes
	case mnil:
		return knull
	case mfalse, mtrue:
		return kbool
	default:
		return invalid
	}
}

func Skip(b []byte) ([]byte, error) {
	sz, asz, err := getSize(b)
	if err != nil {
		return b, err
	}
	if sz > 0 {
		b, err = skipN(b, sz)
		if err != nil {
			return b, err
		}
	}
	if asz > 0 {
		for i := 0; i < asz; i++ {
			b, err = Skip(b)
			if err != nil {
				return b, err
			}
		}
	}
	return b, nil
}

func skipN(b []byte, n int) ([]byte, error) {
	if len(b) < n {
		return nil, ErrShortBytes
	}
	return b[n:], nil
}

// returns (skip N bytes, skip M objects, error)
func getSize(b []byte) (int, int, error) {
	if len(b) < 1 {
		return 0, 0, ErrShortBytes
	}

	lead := b[0]

	switch {
	case isfixarray(lead):
		return 1, int(rfixarray(lead)), nil

	case isfixmap(lead):
		return 1, 2 * int(rfixmap(lead)), nil

	case isfixstr(lead):
		return int(rfixstr(lead)) + 1, 0, nil

	case isfixint(lead):
		return 1, 0, nil

	case isnfixint(lead):
		return 1, 0, nil
	}

	switch lead {

	// the following objects
	// are always the same
	// number of bytes (including
	// the leading byte)
	case mnil, mfalse, mtrue:
		return 1, 0, nil
	case mint64, muint64, mfloat64:
		return 9, 0, nil
	case mint32, muint32, mfloat32:
		return 5, 0, nil
	case mint8, muint8:
		return 2, 0, nil
	case mfixext1:
		return 3, 0, nil
	case mfixext2:
		return 4, 0, nil
	case mfixext4:
		return 6, 0, nil
	case mfixext8:
		return 10, 0, nil
	case mfixext16:
		return 18, 0, nil

	// the following objects
	// need to be skipped N bytes
	// plus the lead byte and size bytes
	case mbin8, mstr8:
		if len(b) < 2 {
			return 0, 0, ErrShortBytes
		}
		return int(uint8(b[1])) + 2, 0, nil

	case mbin16, mstr16:
		if len(b) < 3 {
			return 0, 0, ErrShortBytes
		}
		return int(binary.BigEndian.Uint16(b[1:])) + 3, 0, nil

	case mbin32, mstr32:
		if len(b) < 5 {
			return 0, 0, ErrShortBytes
		}
		return int(binary.BigEndian.Uint32(b[1:])) + 5, 0, nil

	// variable extensions
	// require 1 extra byte
	// to skip
	case mext8:
		if len(b) < 3 {
			return 0, 0, ErrShortBytes
		}
		return int(uint8(b[1])) + 3, 0, nil

	case mext16:
		if len(b) < 4 {
			return 0, 0, ErrShortBytes
		}
		return int(binary.BigEndian.Uint16(b[1:])) + 4, 0, nil

	case mext32:
		if len(b) < 6 {
			return 0, 0, ErrShortBytes
		}
		return int(binary.BigEndian.Uint32(b[1:])) + 6, 0, nil

	// arrays skip lead byte,
	// size byte, N objects
	case marray16:
		if len(b) < 3 {
			return 0, 0, ErrShortBytes
		}
		return 3, int(binary.BigEndian.Uint16(b[1:])), nil

	case marray32:
		if len(b) < 5 {
			return 0, 0, ErrShortBytes
		}
		return 5, int(binary.BigEndian.Uint32(b[1:])), nil

	// maps skip lead byte,
	// size byte, 2N objects
	case mmap16:
		if len(b) < 3 {
			return 0, 0, ErrShortBytes
		}
		return 3, 2 * (int(binary.BigEndian.Uint16(b[1:]))), nil

	case mmap32:
		if len(b) < 5 {
			return 0, 0, ErrShortBytes
		}
		return 5, 2 * (int(binary.BigEndian.Uint32(b[1:]))), nil

	default:
		return 0, 0, InvalidPrefixError(lead)

	}
}
