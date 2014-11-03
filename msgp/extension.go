package msgp

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	// Complex64 numbers are encoded as extension #3
	Complex64Extension = 3

	// Complex128 numbers are encoded as extension #4
	Complex128Extension = 4

	// time.Time values are encoded as extension #5
	TimeExtension = 5
)

var (
	extensionReg map[int8]func() Extension
)

func init() {
	extensionReg = make(map[int8]func() Extension)
}

// RegisterExtension registers extensions so that they
// can be initialized and returned by methods that
// decode `interface{}` values. This should only
// be called during initialization. f() should return
// a newly-initialized zero value of the extension. Keep in
// mind that extensions 3, 4, and 5 are reserved for
// complex64, complex128, and time.Time, respectively,
// and that MessagePack reserves extension types from -127 to -1.
//
// For example, if you wanted to register a user-defined struct:
//
//  msgp.RegisterExtension(10, func() msgp.Extension { &MyExtension{} })
//
func RegisterExtension(typ int8, f func() Extension) {
	extensionReg[typ] = f
}

func errExt(got int8, wanted int8) error {
	return fmt.Errorf("msgp: error decoding extension: wanted type %d; got type %d", wanted, got)
}

// Extension is the interface fulfilled
// by types that want to define their
// own binary encoding.
type Extension interface {
	// ExtensionType should return
	// a int8 that identifies the concrete
	// type of the extension. (Types <0 are
	// officially reserved by the MessagePack
	// specifications.)
	ExtensionType() int8

	// Len should return the length
	// of the data to be encoded
	Len() int

	// MarshalBinaryTo should copy
	// the data into the supplied slice,
	// assuming that the slice has length Len()
	MarshalBinaryTo([]byte) error

	UnmarshalBinary([]byte) error
}

// RawExtension implements the Extension interface
type RawExtension struct {
	Type int8
	Data []byte
}

func (r *RawExtension) ExtensionType() int8 { return r.Type }
func (r *RawExtension) Len() int            { return len(r.Data) }
func (r *RawExtension) MarshalBinaryTo(d []byte) error {
	copy(d, r.Data)
	return nil
}

func (r *RawExtension) UnmarshalBinary(b []byte) error {
	if cap(r.Data) >= len(b) {
		r.Data = r.Data[0:len(b)]
	} else {
		r.Data = make([]byte, len(b))
	}
	copy(r.Data, b)
	return nil
}

func (mw *Writer) WriteExtension(e Extension) (int, error) {
	l := e.Len()
	var (
		write bool
		err   error
	)
	switch l {
	case 0:
		mw.scratch[0] = mext8
		mw.scratch[1] = 0
		mw.scratch[2] = byte(e.ExtensionType())
		return mw.w.Write(mw.scratch[:3])
	case 1:
		mw.scratch[0] = mfixext1
		mw.scratch[1] = byte(e.ExtensionType())
		write = true
	case 2:
		mw.scratch[0] = mfixext2
		mw.scratch[1] = byte(e.ExtensionType())
		write = true
	case 4:
		mw.scratch[0] = mfixext4
		mw.scratch[1] = byte(e.ExtensionType())
		write = true
	case 8:
		mw.scratch[0] = mfixext8
		mw.scratch[1] = byte(e.ExtensionType())
		write = true
	case 16:
		mw.scratch[0] = mfixext16
		mw.scratch[1] = byte(e.ExtensionType())
		write = true
	}
	var nn int
	if write {
		nn, err = mw.w.Write(mw.scratch[:2])
	} else {
		switch {
		case l < math.MaxUint8:
			mw.scratch[0] = mext8
			mw.scratch[1] = byte(uint8(l))
			mw.scratch[2] = byte(e.ExtensionType())
			nn, err = mw.w.Write(mw.scratch[:3])
		case l < math.MaxUint16:
			mw.scratch[0] = mext16
			binary.BigEndian.PutUint16(mw.scratch[1:], uint16(l))
			mw.scratch[3] = byte(e.ExtensionType())
			nn, err = mw.w.Write(mw.scratch[:4])
		default:
			mw.scratch[0] = mext32
			binary.BigEndian.PutUint32(mw.scratch[1:], uint32(l))
			mw.scratch[5] = byte(e.ExtensionType())
			nn, err = mw.w.Write(mw.scratch[:6])
		}
	}
	if err != nil {
		return nn, err
	}
	var b []byte
	if e.Len() <= cap(mw.scratch) {
		b = mw.scratch[0:e.Len()]
	} else {
		b = make([]byte, e.Len())
	}
	err = e.MarshalBinaryTo(b)
	if err != nil {
		return nn, err
	}
	n := nn
	nn, err = mw.w.Write(b)
	return n + nn, err
}

// peek at the extension type, assuming the next
// kind to be read is Extension
func (m *Reader) peekExtensionType() (int8, error) {
	l, err := m.r.Peek(6)
	if err != nil {
		return 0, err
	}
	switch l[0] {
	case mfixext1, mfixext2, mfixext4, mfixext8, mfixext16:
		return int8(l[1]), nil
	case mext8:
		return int8(l[2]), nil
	case mext16:
		return int8(l[3]), nil
	case mext32:
		return int8(l[5]), nil
	default:
		return 0, fmt.Errorf("unexpected leading byte %x for extension", l[0])
	}
}

// peekExtension peeks at the extension encoding type
// (must guarantee at least 3 bytes in 'b')
func peekExtension(b []byte) (int8, error) {
	switch b[0] {
	case mfixext1, mfixext2, mfixext4, mfixext8, mfixext16:
		return int8(b[1]), nil
	case mext8:
		return int8(b[2]), nil
	case mext16:
		if len(b) < 4 {
			return 0, ErrShortBytes
		}
		return int8(b[3]), nil
	case mext32:
		if len(b) < 5 {
			return 0, ErrShortBytes
		}
		return int8(b[5]), nil
	default:
		return 0, fmt.Errorf("unexpected leading byte %x for extension", b[0])
	}
}

func (m *Reader) ReadExtension(e Extension) (n int, err error) {
	var lead byte
	var nn int
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n++
	var read int
	switch lead {
	case mfixext1:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[0]) != e.ExtensionType() {
			err = errExt(int8(m.leader[0]), e.ExtensionType())
			return
		}
		err = e.UnmarshalBinary(m.leader[1:2])
		return

	case mfixext2:
		nn, err = io.ReadFull(m.r, m.leader[:3])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[0]) != e.ExtensionType() {
			err = errExt(int8(m.leader[0]), e.ExtensionType())
			return
		}
		err = e.UnmarshalBinary(m.leader[1:3])
		return

	case mfixext4:
		nn, err = io.ReadFull(m.r, m.leader[:5])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[0]) != e.ExtensionType() {
			err = errExt(int8(m.leader[0]), e.ExtensionType())
			return
		}
		err = e.UnmarshalBinary(m.leader[1:5])
		return

	case mfixext8:
		nn, err = io.ReadFull(m.r, m.leader[:9])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[0]) != e.ExtensionType() {
			err = errExt(int8(m.leader[0]), e.ExtensionType())
			return
		}
		err = e.UnmarshalBinary(m.leader[1:9])
		return

	case mfixext16:
		nn, err = io.ReadFull(m.r, m.leader[:17])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[0]) != e.ExtensionType() {
			err = errExt(int8(m.leader[0]), e.ExtensionType())
			return
		}
		err = e.UnmarshalBinary(m.leader[1:17])
		return

	case mext8:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[1]) != e.ExtensionType() {
			err = errExt(int8(m.leader[0]), e.ExtensionType())
			return
		}
		read = int(uint8(m.leader[0]))

	case mext16:
		nn, err = io.ReadFull(m.r, m.leader[:3])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[2]) != e.ExtensionType() {
			err = errExt(int8(m.leader[2]), e.ExtensionType())
			return
		}
		read = int(binary.BigEndian.Uint16(m.leader[0:2]))

	case mext32:
		nn, err = io.ReadFull(m.r, m.leader[:5])
		n += nn
		if err != nil {
			return
		}
		if int8(m.leader[4]) != e.ExtensionType() {
			err = errExt(int8(m.leader[4]), e.ExtensionType())
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[0:4]))

	}

	var bts []byte
	bts, nn, err = readN(m.r, m.scratch, read)
	n += nn
	if err != nil {
		return
	}
	err = e.UnmarshalBinary(bts)
	return
}

func AppendExtension(b []byte, e Extension) ([]byte, error) {
	l := e.Len()
	o, n := ensure(b, ExtensionPrefixSize+l)
	switch l {
	case 0:
		o[n] = mext8
		o[n+1] = 0
		o[n+2] = byte(e.ExtensionType())
		return o[:n+3], nil
	case 1:
		o[n] = mfixext1
		o[n+1] = byte(e.ExtensionType())
		n += 2
	case 2:
		o[n] = mfixext2
		o[n+1] = byte(e.ExtensionType())
		n += 2
	case 4:
		o[n] = mfixext4
		o[n+1] = byte(e.ExtensionType())
		n += 2
	case 8:
		o[n] = mfixext8
		o[n+1] = byte(e.ExtensionType())
		n += 2
	case 16:
		o[n] = mfixext16
		o[n+1] = byte(e.ExtensionType())
		n += 2
	}
	switch {
	case l < math.MaxUint8:
		o[n] = mext8
		o[n+1] = byte(uint8(l))
		o[n+2] = byte(e.ExtensionType())
		n += 3
	case l < math.MaxUint16:
		o[n] = mext16
		binary.BigEndian.PutUint16(o[n+1:], uint16(l))
		o[n+3] = byte(e.ExtensionType())
		n += 4
	default:
		o[n] = mext32
		binary.BigEndian.PutUint32(o[n+1:], uint32(l))
		o[n+5] = byte(e.ExtensionType())
		n += 6
	}
	return o[:n+l], e.MarshalBinaryTo(o[n : n+l])
}

func ReadExtensionBytes(b []byte, e Extension) ([]byte, error) {
	l := len(b)
	if l < 3 {
		return b, ErrShortBytes
	}
	lead := b[0]
	var (
		sz  int // size of 'data'
		off int // offset of 'data'
		typ int8
	)
	switch lead {
	case mfixext1:
		typ = int8(b[1])
		sz = 1
		off = 2
	case mfixext2:
		typ = int8(b[1])
		sz = 2
		off = 2
	case mfixext4:
		typ = int8(b[1])
		sz = 4
		off = 2
	case mfixext8:
		typ = int8(b[1])
		sz = 8
		off = 2
	case mfixext16:
		typ = int8(b[1])
		sz = 16
		off = 2
	case mext8:
		sz = int(uint8(b[1]))
		typ = int8(b[2])
		off = 3
		if sz == 0 {
			return b[3:], e.UnmarshalBinary(b[3:3])
		}
	case mext16:
		if l < 4 {
			return b, ErrShortBytes
		}
		sz = int(binary.BigEndian.Uint16(b[1:]))
		typ = int8(b[3])
		off = 4
	case mext32:
		if l < 6 {
			return b, ErrShortBytes
		}
		sz = int(binary.BigEndian.Uint32(b[1:]))
		typ = int8(b[5])
		off = 6
	}

	if typ != e.ExtensionType() {
		return b, errExt(typ, e.ExtensionType())
	}

	// the data of the extension starts
	// at 'off' and is 'sz' bytes long
	if len(b[off:]) < sz {
		return b, ErrShortBytes
	}
	return b[off+sz:], e.UnmarshalBinary(b[off : off+sz])
}
