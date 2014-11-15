package msgp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
	"unsafe"
)

var (
	// ErrShortBytes is returned when the
	// slice being decoded is too short to
	// contain the contents of the message
	ErrShortBytes = errors.New("msgp: too few bytes left to read object")
)

// IsNil returns true if len(b)>0 and
// the leading byte is a 'nil' MessagePack
// byte; false otherwise
func IsNil(b []byte) bool {
	if len(b) > 0 && b[0] == mnil {
		return true
	}
	return false
}

// IsError returns whether or not
// an error belongs to the msgp package.
// (This is useful for determining if an
// error was an i/o error as opposed to
// an encoding error.)
func IsError(err error) bool {
	if err != nil {
		if err == ErrShortBytes {
			return true
		}
		switch err.(type) {
		case IntOverflow, UintOverflow, TypeError,
			ArrayError, InvalidPrefixError, ExtensionTypeError:
			return true
		default:
			return strings.HasPrefix(err.Error(), "msgp")
		}
	}
	return false
}

// IntOverflow is returned when a call
// would downcast an integer to a type
// with too few bits to hold its value
type IntOverflow struct {
	Value         int64 // the value of the integer
	FailedBitsize int   // the bit size that the int64 could not fit into
}

// Error implements the error interface
func (i IntOverflow) Error() string {
	return fmt.Sprintf("msgp: %d overflows int%d", i.Value, i.FailedBitsize)
}

// UintOverflow is returned when a call
// would downcast an unsigned integer to a type
// with too few bits to hold its value
type UintOverflow struct {
	Value         uint64 // value of the uint
	FailedBitsize int    // the bit size that couldn't fit the value
}

// Error implements the error interface
func (u UintOverflow) Error() string {
	return fmt.Sprintf("msgp: %d overflows uint%d", u.Value, u.FailedBitsize)
}

// ReadMapHeaderBytes reads a map header size
// from 'b' and returns the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a map)
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
		err = TypeError{Method: MapType, Encoded: getType(lead)}
		return
	}
}

// ReadMapKeyZC attempts to read a map key
// from 'b' and returns the key bytes and the remaining bytes
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a str or bin)
func ReadMapKeyZC(b []byte) ([]byte, []byte, error) {
	o, b, err := ReadStringZC(b)
	if err != nil {
		if tperr, ok := err.(TypeError); ok && tperr.Encoded == BinType {
			return ReadBytesZC(b)
		}
		return nil, b, err
	}
	return o, b, nil
}

// ReadArrayHeaderBytes attempts to read
// the array header size off of 'b' and return
// the size and remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not an array)
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
		err = TypeError{Method: ArrayType, Encoded: getType(lead)}
		return
	}
}

// ReadNilBytes tries to read a "nil" byte
// off of 'b' and return the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a 'nil')
func ReadNilBytes(b []byte) ([]byte, error) {
	if len(b) < 1 {
		return nil, ErrShortBytes
	}
	if b[0] != mnil {
		return b, TypeError{Method: NilType, Encoded: getType(b[0])}
	}
	return b[1:], nil
}

// ReadFloat64Bytes tries to read a float64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a float64)
func ReadFloat64Bytes(b []byte) (f float64, o []byte, err error) {
	if len(b) < 9 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfloat64 {
		err = TypeError{Method: Float64Type, Encoded: getType(b[0])}
		return
	}

	copy((*(*[8]byte)(unsafe.Pointer(&f)))[:], b[1:])
	o = b[9:]
	return
}

// ReadFloat32Bytes tries to read a float64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a float32)
func ReadFloat32Bytes(b []byte) (f float32, o []byte, err error) {
	if len(b) < 5 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfloat32 {
		err = TypeError{Method: Float32Type, Encoded: getType(b[0])}
		return
	}

	copy((*(*[4]byte)(unsafe.Pointer(&f)))[:], b[1:])
	o = b[5:]
	return
}

// ReadBoolBytes tries to read a float64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a bool)
func ReadBoolBytes(b []byte) (bool, []byte, error) {
	if len(b) < 1 {
		return false, b, ErrShortBytes
	}
	switch b[0] {
	case mtrue:
		return true, b[1:], nil
	case mfalse:
		return false, b[1:], nil
	default:
		return false, b, TypeError{Method: BoolType, Encoded: getType(b[0])}
	}
}

// ReadInt64Bytes tries to read an int64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError (not a int)
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
		err = TypeError{Method: IntType, Encoded: getType(lead)}
		return
	}
}

// ReadInt32Bytes tries to read an int32
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a int)
// - IntOverflow{} (value doesn't fit in int32)
func ReadInt32Bytes(b []byte) (int32, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt32 || i < math.MinInt32 {
		return 0, o, IntOverflow{Value: i, FailedBitsize: 32}
	}
	return int32(i), o, err
}

// ReadInt16Bytes tries to read an int16
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a int)
// - IntOverflow{} (value doesn't fit in int16)
func ReadInt16Bytes(b []byte) (int16, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt16 || i < math.MinInt16 {
		return 0, o, IntOverflow{Value: i, FailedBitsize: 16}
	}
	return int16(i), o, err
}

// ReadInt8Bytes tries to read an int16
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a int)
// - IntOverflow{} (value doesn't fit in int8)
func ReadInt8Bytes(b []byte) (int8, []byte, error) {
	i, o, err := ReadInt64Bytes(b)
	if i > math.MaxInt8 || i < math.MinInt8 {
		return 0, o, IntOverflow{Value: i, FailedBitsize: 8}
	}
	return int8(i), o, err
}

// ReadIntBytes tries to read an int
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a int)
// - IntOverflow{} (value doesn't fit in int; 32-bit platforms only)
func ReadIntBytes(b []byte) (int, []byte, error) {
	var v int
	i, o, err := ReadInt64Bytes(b)
	if unsafe.Sizeof(v) == 4 {
		if i > math.MaxInt32 || i < math.MinInt32 {
			return 0, o, IntOverflow{Value: i, FailedBitsize: 32}
		}
	}
	v = int(i)
	return v, o, err
}

// ReadUint64Bytes tries to read a uint64
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a uint)
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
		err = TypeError{Method: UintType, Encoded: getType(lead)}
		return
	}
}

// ReadUint32Bytes tries to read a uint32
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a uint)
// - UintOverflow{} (value too large for uint32)
func ReadUint32Bytes(b []byte) (uint32, []byte, error) {
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint32 {
		return 0, nil, UintOverflow{Value: v, FailedBitsize: 32}
	}
	return uint32(v), o, err
}

// ReadUint16Bytes tries to read a uint16
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a uint)
// - UintOverflow{} (value too large for uint16)
func ReadUint16Bytes(b []byte) (uint16, []byte, error) {
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint16 {
		return 0, nil, UintOverflow{Value: v, FailedBitsize: 16}
	}
	return uint16(v), o, err
}

// ReadUint8Bytes tries to read a uint8
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a uint)
// - UintOverflow{} (value too large for uint8)
func ReadUint8Bytes(b []byte) (uint8, []byte, error) {
	v, o, err := ReadUint64Bytes(b)
	if v > math.MaxUint8 {
		return 0, nil, UintOverflow{Value: v, FailedBitsize: 8}
	}
	return uint8(v), o, err
}

// ReadUintBytes tries to read a uint
// from 'b' and return the value and the remaining bytes.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a uint)
// - UintOverflow{} (value too large for uint; 32-bit platforms only)
func ReadUintBytes(b []byte) (uint, []byte, error) {
	var l uint
	v, o, err := ReadUint64Bytes(b)
	if unsafe.Sizeof(l) == 4 && v > math.MaxUint32 {
		return 0, nil, UintOverflow{Value: v, FailedBitsize: 32}
	}
	l = uint(v)
	return l, o, err
}

// ReadByteBytes is analagous to ReadUint8Bytes
func ReadByteBytes(b []byte) (byte, []byte, error) {
	return ReadUint8Bytes(b)
}

// ReadBytesBytes reads a 'bin' object
// from 'b' and returns its vaue and
// the remaining bytes in 'b'.
// Possible errors:
// - ErrShortBytes (too few bytes)
// - TypeError{} (not a 'bin' object)
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
		err = TypeError{Method: BinType, Encoded: getType(lead)}
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
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError{} (object not 'bin')
func ReadBytesZC(b []byte) (v []byte, o []byte, err error) {
	return readBytesBytes(b, nil, true)
}

// ReadStringZC reads a messagepack string field
// without copying. The returned []byte points
// to the same memory as the input slice.
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError{} (object not 'str')
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
			err = TypeError{Method: StrType, Encoded: getType(lead)}
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

// ReadStringBytes reads a 'str' object
// from 'b' and returns its value and the
// remaining bytes in 'b'.
// Possible errors:
// - ErrShortBytes (b not long enough)
// - TypeError{} (not 'str' type)
func ReadStringBytes(b []byte) (string, []byte, error) {
	v, o, err := ReadStringZC(b)
	return string(v), o, err
}

// ReadComplex128Bytes reads a complex128
// extension object from 'b' and returns the
// remaining bytes.
// Possible errors:
// - ErrShortBytes (not enough bytes in 'b')
// - TypeError{} (object not a complex128)
// - ExtensionTypeError{} (object an extension of the correct size, but not a complex128)
func ReadComplex128Bytes(b []byte) (c complex128, o []byte, err error) {
	if len(b) < 18 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext16 {
		v := getType(b[0])
		err = TypeError{Method: Complex128Type, Encoded: v}
		return
	}

	if int8(b[1]) != Complex128Extension {
		err = errExt(int8(b[1]), Complex128Extension)
		return
	}
	copy((*(*[16]byte)(unsafe.Pointer(&c)))[:], b[2:])
	o = b[18:]
	return
}

// ReadComplex64Bytes reads a complex64
// extension object from 'b' and returns the
// remaining bytes.
// Possible errors:
// - ErrShortBytes (not enough bytes in 'b')
// - TypeError{} (object not a complex64)
// - ExtensionTypeError{} (object an extension of the correct size, but not a complex64)
func ReadComplex64Bytes(b []byte) (c complex64, o []byte, err error) {
	if len(b) < 10 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext8 {
		v := getType(b[0])
		err = TypeError{Method: Complex64Type, Encoded: v}
		return
	}

	if b[1] != Complex64Extension {
		err = errExt(int8(b[1]), Complex64Extension)
		return
	}
	copy((*(*[8]byte)(unsafe.Pointer(&c)))[:], b[2:])
	o = b[10:]
	return
}

// ReadTimeBytes reads a time.Time
// extension object from 'b' and returns the
// remaining bytes.
// Possible errors:
// - ErrShortBytes (not enough bytes in 'b')
// - TypeError{} (object not a complex64)
// - ExtensionTypeError{} (object an extension of the correct size, but not a time.Time)
func ReadTimeBytes(b []byte) (t time.Time, o []byte, err error) {
	if len(b) < 18 {
		err = ErrShortBytes
		return
	}

	if b[0] != mfixext16 {
		v := getType(b[0])
		err = TypeError{Method: TimeType, Encoded: v}
		return
	}

	if int8(b[1]) != TimeExtension {
		err = errExt(int8(b[1]), TimeExtension)
		return
	}

	err = t.UnmarshalBinary(b[2:17])
	o = b[18:]
	return
}

// ReadMapStrIntfBytes reads a map[string]interface{}
// out of 'b' and returns the map and remaining bytes.
// If 'old' is non-nil, the values will be read into that map.
func ReadMapStrIntfBytes(b []byte, old map[string]interface{}) (v map[string]interface{}, o []byte, err error) {
	var sz uint32
	o = b
	sz, o, err = ReadMapHeaderBytes(o)

	if err != nil {
		return
	}

	if old != nil {
		for key := range old {
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

// ReadIntfBytes attempts to read
// the next object out of 'b' as a raw interface{} and
// return the remaining bytes.
func ReadIntfBytes(b []byte) (i interface{}, o []byte, err error) {
	if len(b) < 1 {
		err = ErrShortBytes
		return
	}

	k := getType(b[0])

	switch k {
	case MapType:
		i, o, err = ReadMapStrIntfBytes(b, nil)
		return

	case ArrayType:
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

	case Float32Type:
		i, o, err = ReadFloat32Bytes(b)
		return

	case Float64Type:
		i, o, err = ReadFloat64Bytes(b)
		return

	case IntType:
		i, o, err = ReadInt64Bytes(b)
		return

	case UintType:
		i, o, err = ReadUint64Bytes(b)
		return

	case BoolType:
		i, o, err = ReadBoolBytes(b)
		return

	case ExtensionType:
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

	case NilType:
		o, err = ReadNilBytes(b)
		return

	case BinType:
		i, o, err = ReadBytesBytes(b, nil)
		return

	case StrType:
		i, o, err = ReadStringBytes(b)
		return

	default:
		err = InvalidPrefixError(b[0])
		return
	}
}

func getType(v byte) Type {

	// fixed encoding
	switch {
	case isfixmap(v):
		return MapType
	case isfixarray(v):
		return ArrayType
	case isfixint(v), isnfixint(v):
		return IntType
	case isfixstr(v):
		return StrType
	}

	// var encoding
	switch v {
	case mmap16, mmap32:
		return MapType
	case marray16, marray32:
		return ArrayType
	case mfloat32:
		return Float32Type
	case mfloat64:
		return Float64Type
	case mint8, mint16, mint32, mint64:
		return IntType
	case muint8, muint16, muint32, muint64:
		return UintType
	case mfixext1, mfixext2, mfixext4, mfixext8, mfixext16, mext8, mext16, mext32:
		return ExtensionType
	case mstr8, mstr16, mstr32:
		return StrType
	case mbin8, mbin16, mbin32:
		return BinType
	case mnil:
		return NilType
	case mfalse, mtrue:
		return BoolType
	default:
		return InvalidType
	}
}

// Skip skips the next object in 'b' and
// returns the remaining bytes. If the object
// is a map or array, all of its elements
// will be skipped.
// Possible Errors:
// - ErrShortBytes (not enough bytes in b)
// - InvalidPrefixError (bad encoding)
func Skip(b []byte) ([]byte, error) {
	sz, asz, err := getSize(b)
	if err != nil {
		return b, err
	}
	b, err = skipN(b, sz)
	if err != nil {
		return b, err
	}
	for i := 0; i < asz; i++ {
		b, err = Skip(b)
		if err != nil {
			return b, err
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
