package enc

import (
	"encoding/binary"
	"io"
	"math"
	"unsafe"
)

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

const (
	// Complex64 numbers are encoded as extension #3
	Complex64Extension = 3

	// Complex128 numbers are encoded as extension #4
	Complex128Extension = 4
)

type MsgWriter struct {
	w       io.Writer
	scratch [24]byte
}

func (mw *MsgWriter) WriteMapHeader(sz uint32) (n int, err error) {
	switch {
	case sz < 16:
		//fixmap is 1000XXXX: 10000000 | 0000XXXX
		mw.scratch[0] = (mfixmap | (uint8(sz) & 0x0f))
		return mw.w.Write(mw.scratch[:1])

	case sz < 1<<16-1:
		mw.scratch[0] = mmap16
		mw.scratch[1] = byte(sz >> 8)
		mw.scratch[2] = byte(sz)
		return mw.w.Write(mw.scratch[:3])

	default:
		mw.scratch[0] = mmap32
		binary.BigEndian.PutUint32(mw.scratch[1:], sz)
		return mw.w.Write(mw.scratch[:5])
	}
}

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

func (mw *MsgWriter) WriteArrayHeader(sz uint32) (n int, err error) {
	switch {
	case sz < 16:
		mw.scratch[0] = mfixarray | byte(sz)
		return mw.w.Write(mw.scratch[:1])
	case sz < math.MaxUint16:
		mw.scratch[0] = marray16
		binary.BigEndian.PutUint16(mw.scratch[1:], uint16(sz))
		return mw.w.Write(mw.scratch[:3])
	default:
		mw.scratch[0] = marray32
		binary.BigEndian.PutUint32(mw.scratch[1:], sz)
		return mw.w.Write(mw.scratch[:5])
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

func (mw *MsgWriter) WriteNil() (n int, err error) {
	mw.scratch[0] = mnil
	return mw.w.Write(mw.scratch[:1])
}

// WriteNil writes the 'nil' byte
func WriteNil(w io.Writer) (n int, err error) { return w.Write([]byte{mnil}) }

func (mw *MsgWriter) WriteFloat64(f float64) (n int, err error) {
	mw.scratch[0] = mfloat64
	bits := *(*uint64)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint64(mw.scratch[1:], bits)
	return mw.w.Write(mw.scratch[:9])
}

// WriteFloat64 writes a float64
func WriteFloat64(w io.Writer, f float64) (n int, err error) {
	var scratch [9]byte
	var bits uint64
	scratch[0] = mfloat64
	bits = *(*uint64)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint64(scratch[1:], bits)
	return w.Write(scratch[:])
}

func (mw *MsgWriter) WriteFloat32(f float32) (n int, err error) {
	mw.scratch[0] = mfloat32
	bits := *(*uint32)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint32(mw.scratch[1:], bits)
	return mw.w.Write(mw.scratch[:5])
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

func (mw *MsgWriter) WriteInt64(i int64) (n int, err error) {
	switch {
	case i < 0 && i > -32:
		mw.scratch[0] = byte(int8(i))
		return mw.w.Write(mw.scratch[:1])

	case i >= 0 && i < 128:
		mw.scratch[0] = (byte(i) & 0x7f)
		return mw.w.Write(mw.scratch[:1])

	case abs(i) < math.MaxInt8:
		mw.scratch[0] = mint8
		mw.scratch[1] = byte(int8(i))
		return mw.w.Write(mw.scratch[:2])

	case abs(i) < math.MaxInt16:
		mw.scratch[0] = mint16
		mw.scratch[1] = byte(int16(i >> 8))
		mw.scratch[2] = byte(int16(i))
		return mw.w.Write(mw.scratch[:3])

	case abs(i) < math.MaxInt32:
		mw.scratch[0] = mint32
		mw.scratch[1] = byte(int32(i >> 24))
		mw.scratch[2] = byte(int32(i >> 16))
		mw.scratch[3] = byte(int32(i >> 8))
		mw.scratch[4] = byte(int32(i))
		return mw.w.Write(mw.scratch[:5])

	default:
		mw.scratch[0] = mint64
		mw.scratch[1] = byte(i >> 56)
		mw.scratch[2] = byte(i >> 48)
		mw.scratch[3] = byte(i >> 40)
		mw.scratch[4] = byte(i >> 32)
		mw.scratch[5] = byte(i >> 24)
		mw.scratch[6] = byte(i >> 16)
		mw.scratch[7] = byte(i >> 8)
		mw.scratch[8] = byte(i)
		return mw.w.Write(mw.scratch[:9])
	}

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

func (mw *MsgWriter) WriteUint64(u uint64) (n int, err error) {
	switch {
	case u < (1 << 7):
		mw.scratch[0] = (byte(u) & 0x7f)
		return mw.w.Write(mw.scratch[:1])
	case u < math.MaxUint16:
		mw.scratch[0] = muint16
		binary.BigEndian.PutUint16(mw.scratch[1:], uint16(u))
		return mw.w.Write(mw.scratch[:3])
	case u < math.MaxUint32:
		mw.scratch[0] = muint32
		binary.BigEndian.PutUint32(mw.scratch[1:], uint32(u))
		return mw.w.Write(mw.scratch[:5])
	default:
		mw.scratch[0] = muint64
		binary.BigEndian.PutUint64(mw.scratch[1:], u)
		return mw.w.Write(mw.scratch[:9])
	}
}

// WriteUint writes a uint64
func WriteUint64(w io.Writer, u uint64) (n int, err error) {
	switch {
	case u < (1 << 7):
		return w.Write([]byte{byte(u) & 0x7f}) // "fixint"
	case u < math.MaxUint16:
		return w.Write([]byte{muint16, byte(u >> 8), byte(u)})
	case u < math.MaxUint32:
		return w.Write([]byte{muint32, byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)})
	default:
		return w.Write([]byte{muint64, byte(u >> 56), byte(u >> 48), byte(u >> 40), byte(u >> 32), byte(u >> 24), byte(u >> 16), byte(u >> 8), byte(u)})
	}
}

// Unsigned Integer Utilities
func WriteByte(w io.Writer, u byte) (int, error)     { return WriteUint8(w, uint8(u)) }
func WriteUint8(w io.Writer, u uint8) (int, error)   { return WriteUint64(w, uint64(u)) }
func WriteUint16(w io.Writer, u uint16) (int, error) { return WriteUint64(w, uint64(u)) }
func WriteUint32(w io.Writer, u uint32) (int, error) { return WriteUint64(w, uint64(u)) }
func WriteUint(w io.Writer, u uint) (int, error)     { return WriteUint64(w, uint64(u)) }

func (mw *MsgWriter) WriteBytes(b []byte) (n int, err error) {
	l := len(b)
	if l > math.MaxUint32 {
		panic("Cannot write bytes with length > MaxUint32")
	}
	sz := uint32(l)
	var nn int
	switch {
	case sz < math.MaxUint8:
		mw.scratch[0] = mbin8
		mw.scratch[1] = byte(sz)
		nn, err = mw.w.Write(mw.scratch[:2])
	case sz < math.MaxUint16:
		mw.scratch[0] = mbin16
		mw.scratch[1] = byte(sz >> 8)
		mw.scratch[2] = byte(sz)
		nn, err = mw.w.Write(mw.scratch[:3])
	default:
		mw.scratch[0] = mbin32
		mw.scratch[1] = byte(sz >> 24)
		mw.scratch[2] = byte(sz >> 16)
		mw.scratch[3] = byte(sz >> 8)
		mw.scratch[4] = byte(sz)
		nn, err = mw.w.Write(mw.scratch[:5])
	}
	n += nn
	if err != nil {
		return
	}
	nn, err = mw.w.Write(b)
	n += nn
	return
}

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

func (mw *MsgWriter) WriteBool(b bool) (n int, err error) {
	if b {
		mw.scratch[0] = mtrue
	} else {
		mw.scratch[1] = mfalse
	}
	return mw.w.Write(mw.scratch[:1])
}

// WriteBool writes the boolean 'b'
func WriteBool(w io.Writer, b bool) (n int, err error) {
	if b {
		return w.Write([]byte{mtrue})
	}
	return w.Write([]byte{mfalse})
}

func (mw *MsgWriter) WriteString(s string) (n int, err error) {
	var nn int
	l := len(s)
	if l > math.MaxUint32 {
		panic("Cannot write string with length greater than MaxUint32")
	}
	sz := uint32(l)
	switch {
	case sz < 32:
		mw.scratch[0] = byte((uint8(sz) & 0x1f) | mfixstr)
		nn, err = mw.w.Write(mw.scratch[:1])
	case sz < 256:
		mw.scratch[0] = mstr8
		mw.scratch[1] = byte(sz)
		nn, err = mw.w.Write(mw.scratch[:2])
	case sz < (1<<16)-1:
		mw.scratch[0] = mstr16
		binary.BigEndian.PutUint16(mw.scratch[1:], uint16(sz))
		nn, err = mw.w.Write(mw.scratch[:3])
	default:
		mw.scratch[0] = mstr32
		binary.BigEndian.PutUint32(mw.scratch[1:], sz)
		nn, err = mw.w.Write(mw.scratch[:5])
	}
	n += nn
	if err != nil {
		return
	}
	nn, err = io.WriteString(mw.w, s)
	n += nn
	return
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

func (mw *MsgWriter) WriteComplex64(f complex64) (n int, err error) {
	mw.scratch[0] = mfixext8
	mw.scratch[1] = Complex64Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint32(mw.scratch[2:], *(*uint32)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint32(mw.scratch[6:], *(*uint32)(unsafe.Pointer(&im)))
	return mw.w.Write(mw.scratch[:10])
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

func (mw *MsgWriter) WriteComplex128(f complex128) (n int, err error) {
	mw.scratch[0] = mfixext16
	mw.scratch[1] = Complex128Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint64(mw.scratch[2:], *(*uint64)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint64(mw.scratch[10:], *(*uint64)(unsafe.Pointer(&im)))
	return mw.w.Write(mw.scratch[:18])
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
