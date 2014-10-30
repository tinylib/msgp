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
	// type of the extension
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
	t := byte(e.ExtensionType())
	l := len(bts)
	mw.scratch[1] = t
	var write bool
	switch l {
	case 0:
		mw.scratch[0] = mext8
		mw.scratch[2] = 0
		return mw.w.Write(mw.scratch[:3])
	case 1:
		mw.scratch[0] = mfixext1
		write = true
	case 2:
		mw.scratch[0] = mfixext2
		write = true
	case 4:
		mw.scratch[0] = mfixext4
		write = true
	case 8:
		mw.scratch[0] = mfixext8
		write = true
	case 16:
		mw.scratch[0] = mfixext16
		write = true
	}
	if write {
		var n int
		n, err = mw.w.Write(mw.scratch[:3])
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
		mw.scratch[2] = byte(uint8(l))
		nn, err = mw.w.Write(mw.scratch[:3])
	case l < math.MaxUint16:
		mw.scratch[0] = mext16
		binary.BigEndian.PutUint16(mw.scratch[2:], uint16(l))
		nn, err = mw.w.Write(mw.scratch[:4])
	default:
		mw.scratch[0] = mext32
		binary.BigEndian.PutUint32(mw.scratch[2:], uint32(l))
		nn, err = mw.w.Write(mw.scratch[:6])
	}
	if err != nil {
		return nn, err
	}
	n := nn
	nn, err = mw.w.Write(bts)
	return n + nn, err
}

func (m *MsgReader) peekExtensionType() (int8, error) {
	l, err := m.r.Peek(2)
	if err == nil {
		return int8(l[1]), nil
	}
	return 0, err
}

func (m *MsgReader) ReadExtension(e Extension) (n int, err error) {
	var lead byte
	var nn int
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n++
	var typ byte
	typ, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n++
	if int8(typ) != e.ExtensionType() {
		err = errExt(int8(typ), e.ExtensionType())
		return
	}
	var read int
	switch lead {
	case mfixext1:
		nn, err = io.ReadFull(m.r, m.leader[:1])
		n += nn
		if err != nil {
			return
		}
		err = e.UnmarshalBinary(m.leader[:1])
		return

	case mfixext2:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		err = e.UnmarshalBinary(m.leader[:2])
		return

	case mfixext4:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		err = e.UnmarshalBinary(m.leader[:4])
		return

	case mfixext8:
		nn, err = io.ReadFull(m.r, m.leader[:8])
		n += nn
		if err != nil {
			return
		}
		err = e.UnmarshalBinary(m.leader[:8])
		return

	case mfixext16:
		nn, err = io.ReadFull(m.r, m.leader[:16])
		n += nn
		if err != nil {
			return
		}
		err = e.UnmarshalBinary(m.leader[:16])
		return

	case mext8:
		lead, err = m.r.ReadByte()
		if err != nil {
			return
		}
		n += 1
		read = int(uint8(lead))

	case mext16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[:]))

	case mext32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[:]))

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
	o, n := ensure(b, 5+l)
	o[n+1] = byte(e.ExtensionType())
	switch l {
	case 0:
		o[n] = mext8
		o[n+2] = 0
		return o[:n+3], nil
	case 1:
		o[n] = mfixext1
		o[n+2] = bts[0]
		return o[:n+3], nil
	case 2:
		o[n] = mfixext2
		copy(o[n+2:], bts)
		return o[:n+4], nil
	case 4:
		o[n] = mfixext4
		copy(o[n+2:], bts)
		return o[:n+6], nil
	case 8:
		o[n] = mfixext8
		copy(o[n+2:], bts)
		return o[:n+10], nil
	case 16:
		o[n] = mfixext16
		copy(o[n+2:], bts)
		return o[:n+18], nil
	}
	switch {
	case l < math.MaxUint8:
		o[n] = mext8
		o[n+2] = byte(uint8(l))
		n += 3
	case l < math.MaxUint16:
		o[n] = mext16
		binary.BigEndian.PutUint16(o[n+2:], uint16(l))
		n += 4
	default:
		o[n] = mext32
		binary.BigEndian.PutUint32(o[n+2:], uint32(l))
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
	if int8(b[1]) != e.ExtensionType() {
		return b, errExt(int8(b[1]), e.ExtensionType())
	}
	var (
		sz  int // size of 'data'
		off int // offset of 'data'
	)
	switch lead {
	case mfixext1:
		sz = 1
		off = 2
	case mfixext2:
		sz = 2
		off = 2
	case mfixext4:
		sz = 4
		off = 2
	case mfixext8:
		sz = 8
		off = 2
	case mfixext16:
		sz = 16
		off = 2
	case mext8:
		sz = int(uint8(b[2]))
		off = 3
	case mext16:
		if l < 4 {
			err = ErrShortBytes
			return
		}
		sz = int(binary.BigEndian.Uint16(b[2:]))
		off = 4
	case mext32:
		if l < 6 {
			err = ErrShortBytes
			return
		}
		sz = int(binary.BigEndian.Uint32(b[2:]))
		off = 6
	}

	// the data of the extension starts
	// at 'off' and is 'sz' bytes long

	err = e.UnmarshalBinary(b[off : off+sz])
	return b[:off+sz], err
}
