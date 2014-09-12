package enc

import (
	"encoding/binary"
	"io"
	"math"
	"unsafe"
)

const (
	// Complex64 numbers are encoded as extension #3
	Complex64Extension = 3

	// Complex128 numbers are encoded as extension #4
	Complex128Extension = 4
)

// WriteMapHeader writes the header for a map of size 'sz'
func WriteMapHeader(w io.Writer, sz uint32) (n int, err error) {
	switch {
	case sz < 16:
		//fixmap is 1000XXXX: 10000000 | 0000XXXX
		return w.Write([]byte{(mfixmap | (uint8(sz) & 0x0f))})

	case sz < 1<<16-1:
		return w.Write([]byte{mmap16, byte(sz >> 8), byte(sz)})

	default:
		return w.Write([]byte{mmap32, byte(sz >> 24), byte(sz >> 16), byte(sz >> 8), byte(sz)})
	}
}

// WriteArrayHeader writes the header for an array of size 'sz'
func WriteArrayHeader(w io.Writer, sz uint32) (n int, err error) {
	switch {
	case sz < 16:
		return w.Write([]byte{mfixarray | byte(sz)})
	case sz < math.MaxUint16:
		return w.Write([]byte{marray16, byte(sz >> 8), byte(sz)})
	default:
		return w.Write([]byte{marray32, byte(sz >> 24), byte(sz >> 16), byte(sz >> 8), byte(sz)})
	}
}

// WriteNil writes the 'nil' byte
func WriteNil(w io.Writer) (n int, err error) { return w.Write([]byte{mnil}) }

// WriteFloat64 writes a float64
func WriteFloat64(w io.Writer, f float64) (n int, err error) {
	var scratch [9]byte
	var bits uint64
	scratch[0] = mfloat64
	bits = *(*uint64)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint64(scratch[1:], bits)
	return w.Write(scratch[:])
}

// WriteFloat32 writes a float32
func WriteFloat32(w io.Writer, f float32) (n int, err error) {
	var scratch [5]byte
	var bits uint32
	scratch[0] = mfloat32
	bits = *(*uint32)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint32(scratch[1:], bits)
	return w.Write(scratch[:])
}

// WriteInt writes an int64
func WriteInt64(w io.Writer, i int64) (n int, err error) {
	var lead []byte

	if i >= 0 {
		switch {
		case i < 128:
			lead = []byte{byte(i) & 0x7f}
		case i < math.MaxInt8:
			lead = []byte{mint8, byte(i)}
		case i < math.MaxInt16:
			lead = []byte{mint16, byte(i >> 8), byte(i)}
		case i < math.MaxInt32:
			lead = []byte{mint32, byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
		default:
			lead = []byte{mint64, byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32), byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
		}
	} else {
		switch {
		case i > -32:
			lead = []byte{byte(int8(i))}
		case i > math.MinInt8:
			lead = []byte{mint8, byte(int8(i))}
		case i > math.MinInt16:
			lead = []byte{mint16, byte(int16(i) >> 8), byte(int16(i))}
		case i > math.MinInt32:
			lead = []byte{mint32, byte(int32(i) >> 24), byte(int32(i) >> 16), byte(int32(i) >> 8), byte(int32(i))}
		default:
			lead = []byte{mint64, byte(i >> 56), byte(i >> 48), byte(i >> 40), byte(i >> 32), byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
		}
	}

	return w.Write(lead)
}

// Integer Utilities
func WriteInt8(w io.Writer, i int8) (int, error)   { return WriteInt64(w, int64(i)) }
func WriteInt16(w io.Writer, i int16) (int, error) { return WriteInt64(w, int64(i)) }
func WriteInt32(w io.Writer, i int32) (int, error) { return WriteInt64(w, int64(i)) }
func WriteInt(w io.Writer, i int) (int, error)     { return WriteInt64(w, int64(i)) }

// WriteUint writes a uint64
func WriteUint64(w io.Writer, u uint64) (n int, err error) {
	var lead []byte

	switch {
	case u < 256:
		lead = []byte{byte(u) & 0x7f} // "fixint"
	case u < math.MaxUint16:
		lead = []byte{muint16, byte(u >> 8), byte(u)}
	case u < math.MaxUint32:
		lead = []byte{muint32, byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)}
	default:
		lead = []byte{muint64, byte(u >> 56), byte(u >> 48), byte(u >> 40), byte(u >> 32), byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)}
	}

	return w.Write(lead)
}

// Unsigned Integer Utilities
func WriteByte(w io.Writer, u byte) (int, error)     { return WriteUint8(w, uint8(u)) }
func WriteUint8(w io.Writer, u uint8) (int, error)   { return WriteUint64(w, uint64(u)) }
func WriteUint16(w io.Writer, u uint16) (int, error) { return WriteUint64(w, uint64(u)) }
func WriteUint32(w io.Writer, u uint32) (int, error) { return WriteUint64(w, uint64(u)) }
func WriteUint(w io.Writer, u uint) (int, error)     { return WriteUint64(w, uint64(u)) }

// WriteBytes writes 'b'
func WriteBytes(w io.Writer, b []byte) (n int, err error) {
	l := len(b)
	if l > math.MaxUint32 {
		panic("Cannot write bytes with length larger than MaxUint32")
	}
	bits := uint32(l)
	var lead []byte
	switch {
	case bits > math.MaxUint16:
		lead = []byte{mbin32, byte(bits >> 24), byte(bits >> 16), byte(bits >> 8), byte(bits)}
	case bits > math.MaxUint8:
		tb := uint16(bits)
		lead = []byte{mbin16, byte(tb >> 8), byte(tb)}
	default:
		lead = []byte{mbin8, byte(bits)}
	}
	var nn int
	nn, err = w.Write(lead)
	n += nn
	if err != nil {
		return
	}
	nn, err = w.Write(b)
	n += nn
	return
}

// WriteBool writes the boolean 'b'
func WriteBool(w io.Writer, b bool) (n int, err error) {
	if b {
		return w.Write([]byte{mtrue})
	}
	return w.Write([]byte{mfalse})
}

// WriteString writes the string 's'
func WriteString(w io.Writer, s string) (n int, err error) {
	var nn int
	sz := len(s)
	if sz > math.MaxUint32 {
		panic("Cannot write string with length greater than MaxUint32")
	}
	szbits := uint32(sz)

	switch {
	case szbits < 32:
		nn, err = w.Write([]byte{(byte(uint8(szbits)&0x1f) | mfixstr)}) // 5-bit size: 101XXXXX
	case szbits < 256:
		nn, err = w.Write([]byte{mstr8, byte(szbits)})
	case szbits < (1<<16)-1:
		nn, err = w.Write([]byte{mstr16, byte(uint16(szbits) >> 8), byte(uint16(szbits))})
	default:
		nn, err = w.Write([]byte{mstr32, byte(szbits >> 24), byte(szbits >> 16), byte(szbits >> 8), byte(szbits)})
	}

	n += nn
	if err != nil {
		return
	}
	if sz == 0 {
		return
	}
	nn, err = io.WriteString(w, s) // TODO: use unsafe.Pointer?
	n += nn
	return
}

// WriteComplex64 writes a complex64 as an Extension type with code 3
func WriteComplex64(w io.Writer, f complex64) (n int, err error) {
	var bts [10]byte
	bts[0] = mfixext8
	bts[1] = Complex64Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint32(bts[2:], *(*uint32)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint32(bts[6:], *(*uint32)(unsafe.Pointer(&im)))
	return w.Write(bts[:])
}

// WriteComplex128 writes a complex128 as an Extension type with code 4
func WriteComplex128(w io.Writer, f complex128) (n int, err error) {
	var bts [18]byte
	bts[0] = mfixext16
	bts[1] = Complex128Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint64(bts[2:], *(*uint64)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint64(bts[10:], *(*uint64)(unsafe.Pointer(&im)))
	return w.Write(bts[:])
}

// WriteExtension writes an extension type (tuple of (type, binary))
//func WriteExtension(w io.Writer, typ byte, bts []byte) (n int, err error) {/
//
//}
