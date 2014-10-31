package enc

import (
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func errExt(got int8, wanted int8) error {
	return fmt.Errorf("msgp/enc: error decoding extension: wanted type %d; got type %d", wanted, got)
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

	// Extensions must satisfy the
	// encoding.BinaryMarshaler and
	// encoding.BinaryUnmarshaler
	// interfaces
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

// RawExtension implements the Extension interface
type RawExtension struct {
	Type int8
	Data []byte
}

func (r *RawExtension) ExtensionType() int8 { return r.Type }
func (r *RawExtension) MarshalBinary() ([]byte, error) {
	d := make([]byte, len(r.Data))
	copy(d, r.Data)
	return d, nil
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

func (mw *MsgWriter) WriteExtension(e Extension) (int, error) {
	bts, err := e.MarshalBinary()
	if err != nil {
		return 0, err
	}
	l := len(bts)
	var write bool
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
	if write {
		var n int
		n, err = mw.w.Write(mw.scratch[:2])
		if err != nil {
			return n, err
		}
		n, err = mw.w.Write(bts)
		return n + 2, err
	}

	var nn int
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
	if err != nil {
		return nn, err
	}
	n := nn
	nn, err = mw.w.Write(bts)
	return n + nn, err
}

// peek at the extension type, assuming the next
// kind to be read is Extension
func (m *MsgReader) peekExtensionType() (int8, error) {
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

func (m *MsgReader) ReadExtension(e Extension) (n int, err error) {
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
	bts, nn, err = readN(m.r, nil, read)
	n += nn
	if err != nil {
		return
	}
	err = e.UnmarshalBinary(bts)
	return
}

func AppendExtension(b []byte, e Extension) ([]byte, error) {
	bts, err := e.MarshalBinary()
	if err != nil {
		return b, err
	}
	l := len(bts)
	o, n := ensure(b, 6+l)
	switch l {
	case 0:
		o[n] = mext8
		o[n+1] = 0
		o[n+2] = byte(e.ExtensionType())
		return o[:n+3], nil
	case 1:
		o[n] = mfixext1
		o[n+1] = byte(e.ExtensionType())
		o[n+2] = bts[0]
		return o[:n+3], nil
	case 2:
		o[n] = mfixext2
		o[n+1] = byte(e.ExtensionType())
		copy(o[n+2:], bts)
		return o[:n+4], nil
	case 4:
		o[n] = mfixext4
		o[n+1] = byte(e.ExtensionType())
		copy(o[n+2:], bts)
		return o[:n+6], nil
	case 8:
		o[n] = mfixext8
		o[n+1] = byte(e.ExtensionType())
		copy(o[n+2:], bts)
		return o[:n+10], nil
	case 16:
		o[n] = mfixext16
		o[n+1] = byte(e.ExtensionType())
		copy(o[n+2:], bts)
		return o[:n+18], nil
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
	x := copy(o[n:], bts)
	return o[:n+x], nil
}

func ReadExtensionBytes(b []byte, e Extension) (o []byte, err error) {
	l := len(b)
	if l < 3 {
		err = ErrShortBytes
		return
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
	case mext16:
		if l < 4 {
			err = ErrShortBytes
			return
		}
		sz = int(binary.BigEndian.Uint16(b[1:]))
		typ = int8(b[3])
		off = 4
	case mext32:
		if l < 6 {
			err = ErrShortBytes
			return
		}
		sz = int(binary.BigEndian.Uint32(b[1:]))
		typ = int8(b[5])
		off = 6
	}

	if typ != e.ExtensionType() {
		err = errExt(typ, e.ExtensionType())
		return
	}

	// the data of the extension starts
	// at 'off' and is 'sz' bytes long
	if len(b[off:]) < sz {
		err = ErrShortBytes
		return
	}
	err = e.UnmarshalBinary(b[off : off+sz])
	return b[:off+sz], err
}
