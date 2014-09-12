package enc

import (
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func ReadMapHeader(r io.Reader) (sz uint32, n int, err error) {
	var lead [1]byte
	var nn int
	nn, err = io.ReadFull(r, lead[:])
	n += nn
	if err != nil {
		return
	}
	switch uint8(lead[0]) {
	case mmap16:
		var scratch [2]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		if err != nil {
			return
		}
		usz := binary.BigEndian.Uint16(scratch[:])
		sz = uint32(usz)
		return

	case mmap32:
		var scratch [4]byte
		nn, err = r.Read(scratch[:])
		n += nn
		if err != nil {
			return
		}
		sz = binary.BigEndian.Uint32(scratch[:])
		return

	default:
		// fixmap starts with nibble 1000
		if uint8(lead[0])&0xf0 != mfixmap {
			err = fmt.Errorf("unexpected byte %x for fixmap", lead[0])
			return
		}
		// length in last 4 bits
		var inner uint8 = uint8(lead[0]) & 0x0f
		sz = uint32(inner)
		return

	}
}

func ReadArrayHeader(r io.Reader) (sz uint32, n int, err error) {
	var lead [1]byte
	var nn int
	nn, err = io.ReadFull(r, lead[:])
	n += nn
	if err != nil {
		return
	}
	switch lead[0] {
	case marray32:
		var scratch [4]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		if err != nil {
			return
		}
		sz = binary.BigEndian.Uint32(scratch[:])
		return

	case marray16:
		var scratch [2]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		if err != nil {
			return
		}
		usz := binary.BigEndian.Uint16(scratch[:])
		sz = uint32(usz)
		return

	default:
		// decode fixarray
		if (lead[0] & 0xf0) != mfixarray {
			err = fmt.Errorf("unexpected byte %x for fixarray", lead[0])
			return
		}
		var inner byte = lead[0] & 0x0f
		sz = uint32(inner)
		return
	}
}

func ReadNil(r io.Reader) (n int, err error) {
	var lead [1]byte
	n, err = io.ReadFull(r, lead[:])
	if err != nil {
		return
	}
	if lead[0] != mnil {
		err = fmt.Errorf("unexpected byte %x for Nil", lead[0])
	}
	return
}

// ReadFloat64 reads a float64 out of 'r'
func ReadFloat64(r io.Reader) (f float64, n int, err error) {
	var bts [9]byte
	n, err = io.ReadFull(r, bts[:])
	if err != nil {
		return
	}
	if bts[0] != mfloat64 {
		err = fmt.Errorf("msgp/enc: unexpected byte %x for Float64; expected %x", bts[0], mfloat64)
		return
	}
	bits := binary.BigEndian.Uint64(bts[1:])
	f = *(*float64)(unsafe.Pointer(&bits))
	return
}

// ReadFloat32 read a float32 out of 'r'
func ReadFloat32(r io.Reader) (f float32, n int, err error) {
	var bts [5]byte
	n, err = io.ReadFull(r, bts[:])
	if err != nil {
		return
	}
	if bts[0] != mfloat32 {
		err = fmt.Errorf("msgp/enc: unexpected byte %x for Float64; expected %x", bts[0], mfloat64)
		return
	}
	bits := binary.BigEndian.Uint32(bts[1:])
	f = *(*float32)(unsafe.Pointer(&bits))
	return
}

func ReadBool(r io.Reader) (bool, int, error) {
	var out [1]byte
	n, err := io.ReadFull(r, out[:])
	if err != nil {
		return false, n, err
	}
	switch out[0] {
	case mtrue:
		return true, n, err
	case mfalse:
		return false, n, err
	default:
		return false, n, fmt.Errorf("Unexpected byte %x for bool", out[0])
	}
}

func ReadInt64(r io.Reader) (i int64, n int, err error) {
	var lead [1]byte
	var nn int
	nn, err = io.ReadFull(r, lead[:])
	n += nn
	if err != nil {
		return
	}
	switch lead[0] {
	case mint8:
		var scratch [1]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		i = int64(scratch[0])
		return

	case mint16:
		var scratch [2]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		i = int64((int16(scratch[0]) << 8) | (int16(scratch[1])))
		return

	case mint32:
		var scratch [4]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		i = int64(int32(scratch[0]<<24) | (int32(scratch[1]) << 16) | (int32(scratch[2]) << 8) | (int32(scratch[3])))
		return

	case mint64:
		var scratch [8]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		i |= int64(scratch[0]) << 56
		i |= int64(scratch[1]) << 48
		i |= int64(scratch[2]) << 40
		i |= int64(scratch[3]) << 32
		i |= int64(scratch[4]) << 24
		i |= int64(scratch[5]) << 16
		i |= int64(scratch[6]) << 8
		i |= int64(scratch[7])
		return

	default:
		// try to decode positive fixnum
		if lead[0]&0x7f == 0 {
			i = int64(lead[0] & 0x80)
			return
		}
		// try to decode negative fixnum
		if lead[0]&0xe0 == 0 {
			i = int64(int8(lead[0]))
			return
		}
		err = fmt.Errorf("unknown leading byte %x for int", lead[0])
		return

	}
}

func ReadInt(r io.Reader) (int, int, error) {
	i, n, err := ReadInt64(r)
	return int(i), n, err
}

func ReadInt8(r io.Reader) (int8, int, error) {
	i, n, err := ReadInt64(r)
	return int8(i), n, err
}

func ReadInt16(r io.Reader) (int16, int, error) {
	i, n, err := ReadInt64(r)
	return int16(i), n, err
}

func ReadInt32(r io.Reader) (int32, int, error) {
	i, n, err := ReadInt64(r)
	return int32(i), n, err
}

func ReadUint64(r io.Reader) (u uint64, n int, err error) {
	var lead [1]byte
	var nn int
	nn, err = io.ReadFull(r, lead[:])
	n += nn
	if err != nil {
		return
	}

	switch lead[0] {
	case muint8:
		nn, err = io.ReadFull(r, lead[:])
		n += nn
		u = uint64(lead[0])
		return

	case muint16:
		var scratch [2]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		u = uint64(binary.BigEndian.Uint16(scratch[:]))
		return

	case muint32:
		var scratch [4]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		u = uint64(binary.BigEndian.Uint32(scratch[:]))
		return

	case muint64:
		var scratch [8]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		u = binary.BigEndian.Uint64(scratch[:])
		return

	default:
		// try positive fixnum:
		if lead[0]&0x7f == 0 {
			u = uint64(lead[0] & 0x80)
			return
		}
		err = fmt.Errorf("unexpected byte %x for Uint", lead[0])
		return
	}
}

func ReadUint(r io.Reader) (uint, int, error) {
	u, n, err := ReadUint64(r)
	return uint(u), n, err
}

func ReadByte(r io.Reader) (byte, int, error) {
	u, n, err := ReadUint64(r)
	return byte(u), n, err
}

func ReadUint8(r io.Reader) (uint8, int, error) {
	u, n, err := ReadUint64(r)
	return uint8(u), n, err
}

func ReadUint16(r io.Reader) (uint16, int, error) {
	u, n, err := ReadUint64(r)
	return uint16(u), n, err
}

func ReadUint32(r io.Reader) (uint32, int, error) {
	u, n, err := ReadUint64(r)
	return uint32(u), n, err
}

// ReadBytes reads bytes from 'r', optionally using 'scratch'
// as a buffer to overwrite.
func ReadBytes(r io.Reader, scratch []byte) (b []byte, n int, err error) {
	b = scratch[0:0]

	var lead [2]byte
	var nn int
	nn, err = io.ReadFull(r, lead[:])
	n += nn
	if err != nil {
		return
	}

	var read int // bytes to read
	switch lead[0] {
	case mbin8:
		read = int(lead[1])
	case mbin16:
		var rest [1]byte
		nn, err = io.ReadFull(r, rest[:])
		n += nn
		if err != nil {
			return
		}
		read = int(uint16(lead[1]<<8) | (uint16(rest[0])))
	case mbin32:
		var rest [3]byte
		nn, err = io.ReadFull(r, rest[:])
		n += nn
		if err != nil {
			return
		}
		read = int(uint32(lead[1]<<24) | (uint32(rest[0]) << 16) | (uint32(rest[1]) << 8) | (uint32(rest[2])))
	}

	b, nn, err = readN(r, scratch, read)
	n += nn
	return
}

func readN(r io.Reader, scratch []byte, c int) (b []byte, n int, err error) {
	if scratch == nil || cap(scratch) == 0 {
		b = make([]byte, c)
		n, err = io.ReadFull(r, b)
		return
	}

	if c == 0 {
		return scratch[0:0], 0, nil
	}

	b = scratch[0:0]
	var stack [512]byte
	if c < len(stack) {
		n, err = io.ReadAtLeast(r, stack[:], c)
		b = append(b, stack[:n]...)
		return
	}

	var nn int
	for n < c {
		if (c - n) < len(stack) {
			nn, err = r.Read(stack[:(c - n)])
		} else {
			nn, err = r.Read(stack[:])
		}
		n += nn
		b = append(b, stack[:n]...)
		if err != nil {
			return
		}
	}
	return
}

func ReadStringAsBytes(r io.Reader, scratch []byte) (b []byte, n int, err error) {
	b = scratch[0:0]
	var lead [1]byte
	var nn int
	nn, err = io.ReadFull(r, lead[:])
	n += nn
	if err != nil {
		return
	}
	var read int
	switch lead[0] {
	case mstr8:
		nn, err = io.ReadFull(r, lead[:])
		n += nn
		if err != nil {
			return
		}
		read = int(uint8(lead[0]))

	case mstr16:
		var scratch [2]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		if err != nil {
			return
		}
		read = int((uint16(scratch[0]) << 8) | (uint16(scratch[1])))

	case mstr32:
		var scratch [4]byte
		nn, err = io.ReadFull(r, scratch[:])
		n += nn
		if err != nil {
			return
		}
		read = int((uint32(scratch[0]) << 24) | (uint32(scratch[1]) << 16) | (uint32(scratch[2]) << 8) | (uint32(scratch[3])))

	default:
		// try fixstr - first bits should be 101
		if lead[0]&0xe0 == mfixstr {
			read = int(uint8(lead[0]) & 0x1f)
		} else {
			err = fmt.Errorf("unexpected byte %x for string", lead[0])
			return
		}
	}
	if read == 0 {
		return
	}

	b, nn, err = readN(r, scratch, read)
	n += nn
	return
}

func ReadString(r io.Reader) (string, int, error) {
	bts, n, err := ReadStringAsBytes(r, nil)
	return string(bts), n, err
}

func ReadComplex64(r io.Reader) (f complex64, n int, err error) {
	var scratch [10]byte
	n, err = io.ReadFull(r, scratch[:])
	if err != nil {
		return
	}
	if scratch[0] != mfixext8 {
		err = fmt.Errorf("unexpected byte %x for complex64", scratch[0])
		return
	}
	if scratch[1] != Complex64Extension {
		err = fmt.Errorf("unexpected byte %x for complex64 extension", scratch[1])
		return
	}

	var rlbits uint32
	var imbits uint32
	rlbits = binary.BigEndian.Uint32(scratch[2:])
	imbits = binary.BigEndian.Uint32(scratch[6:])
	var rl float32 = *(*float32)(unsafe.Pointer(&rlbits))
	var im float32 = *(*float32)(unsafe.Pointer(&imbits))
	f = complex(rl, im)
	return
}

func ReadComplex128(r io.Reader) (f complex128, n int, err error) {
	var scratch [18]byte
	n, err = io.ReadFull(r, scratch[:])
	if err != nil {
		return
	}
	if scratch[0] != mfixext16 {
		err = fmt.Errorf("unexpected byte %x for complex128", scratch[0])
	}
	if scratch[1] != Complex128Extension {
		err = fmt.Errorf("unexpected byte %x for complex128 extension", scratch[1])
	}
	var rlbits uint64
	var imbits uint64
	rlbits = binary.BigEndian.Uint64(scratch[2:])
	imbits = binary.BigEndian.Uint64(scratch[10:])
	var rl float64 = *(*float64)(unsafe.Pointer(&rlbits))
	var im float64 = *(*float64)(unsafe.Pointer(&imbits))
	f = complex(rl, im)
	return
}

// UnsafeString returns the byte slice as a string
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(&b[0])), Len: len(b)}))
}

// UnsafeBytes returns the string as a byte slice
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  len(s),
		Cap:  len(s),
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
	}))
}
