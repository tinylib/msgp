package enc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"
	"unsafe"
)

var (
	ErrShortBytes = errors.New("too few bytes")
)

// MsgUnmarshaler is the interface
// fulfilled by types that know how to
// unmarshal themselves from raw MessagePack
// bytes
type MsgUnmarshaler interface {
	// UnmarshalMsg works similarly to encoding.BinaryUnmarshaler,
	// except that it returns the bytes left over in the slice
	UnmarshalMsg([]byte) ([]byte, error)
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
		err = ErrNil
		return

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
		err = fmt.Errorf("unexpected byte %x for map", lead)
		return
	}
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
		err = ErrNil
		return

	case marray16:
		if len(b) < 2 {
			err = ErrShortBytes
			return
		}
		sz = uint32(binary.BigEndian.Uint16(b))
		o = b[2:]
		return

	case marray32:
		if len(b) < 4 {
			err = ErrShortBytes
			return
		}
		sz = binary.BigEndian.Uint32(b)
		o = b[4:]
		return

	default:
		err = fmt.Errorf("unexpected byte 0x%x for array", lead)
		return
	}
}

func ReadNilBytes(b []byte) ([]byte, error) {
	if len(b) < 1 {
		return nil, ErrShortBytes
	}
	if b[0] != mnil {
		return b, fmt.Errorf("unexpected byte %x for Nil", b[0])
	}
	return b[1:], nil
}

func ReadFloat64Bytes(b []byte) (f float64, o []byte, err error) {
	if len(b) < 9 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfloat64 {
		err = fmt.Errorf("unexpected byte %x for float64", b[0])
		return
	}

	bits := binary.BigEndian.Uint64(b[1:])
	f = *(*float64)(unsafe.Pointer(&bits))
	o = b[9:]
	return
}

func ReadFloat32Bytes(b []byte) (f float32, o []byte, err error) {
	if len(b) < 5 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfloat32 {
		err = fmt.Errorf("unexpected byte %x for float32", b[0])
		return
	}

	bits := binary.BigEndian.Uint32(b[1:])
	f = *(*float32)(unsafe.Pointer(&bits))
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
		err = fmt.Errorf("unexpected byte %x for bool", b[0])
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
		err = fmt.Errorf("unexpected byte %x for int64", lead)
		return
	}
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
		err = fmt.Errorf("unexpected byte %x for uint64", lead)
		return
	}
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
		err = fmt.Errorf("unexpected byte %x for bytes", lead)
		return
	}

	if len(b) < read {
		err = ErrShortBytes
		return
	}

	// zero-copy
	if zc {
		v = b[0:read]
		b = b[read:]
		return
	}

	if cap(scratch) >= read {
		v = scratch[0:read]
	} else {
		v = make([]byte, read)
	}

	n := copy(v, b)
	b = b[n:]
	return
}

func ReadZCBytes(b []byte) (v []byte, o []byte, err error) {
	return readBytesBytes(b, nil, true)
}

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
			err = fmt.Errorf("unexpected byte %x for string", lead)
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
		err = fmt.Errorf("unexpected extension type (%x) for fixext16", b[0])
		return
	}

	if b[1] != Complex128Extension {
		err = fmt.Errorf("unexpected byte %x for complex128", b[1])
		return
	}

	rbits := binary.BigEndian.Uint64(b[2:])
	ibits := binary.BigEndian.Uint64(b[10:])

	c = complex(*(*float64)(unsafe.Pointer(&rbits)), *(*float64)(unsafe.Pointer(&ibits)))
	o = b[18:]
	return
}

func ReadComplex64Bytes(b []byte) (c complex64, o []byte, err error) {
	if len(b) < 10 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext8 {
		err = fmt.Errorf("unexpected extension type (%x) for fixext8", b[0])
		return
	}

	if b[1] != Complex64Extension {
		err = fmt.Errorf("unexpected byte %x for complex64", b[1])
		return
	}

	rbits := binary.BigEndian.Uint32(b[2:])
	ibits := binary.BigEndian.Uint32(b[6:])
	c = complex(*(*float32)(unsafe.Pointer(&rbits)), *(*float32)(unsafe.Pointer(&ibits)))
	o = b[10:]
	return
}

func ReadTimeBytes(b []byte) (t time.Time, o []byte, err error) {
	if len(b) < 18 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext16 {
		err = fmt.Errorf("unexpected extension type (%x) for fixext16", b[0])
		return
	}

	if b[1] != TimeExtension {
		err = fmt.Errorf("unexpected byte %x for time extension", b[1])
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
		var key string
		key, o, err = ReadStringBytes(o)
		if err != nil {
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

	k := getKind(b[1])

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
		// TODO
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
		err = fmt.Errorf("unrecognized leading byte: %x", b[1])
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
