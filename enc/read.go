package enc

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

var (
	// ErrNil is returned when reading
	// a value encoded as 'nil'
	ErrNil     = errors.New("value encoded as nil")
	readerPool sync.Pool
)

func init() {
	readerPool.New = func() interface{} {
		return &MsgReader{}
	}
}

func popReader(r io.Reader) *MsgReader {
	p := readerPool.Get().(*MsgReader)
	if p.r == nil {
		p.r = bufio.NewReaderSize(r, 32)
		return p
	}
	p.Reset(r)
	return p
}

func pushReader(m *MsgReader) {
	readerPool.Put(m)
}

func Done(m *MsgReader) {
	if m == nil {
		return
	}
	pushReader(m)
}

type MsgDecoder interface {
	DecodeMsg(r io.Reader) (int, error)
}

func NewDecoder(r io.Reader) *MsgReader {
	return popReader(r)
}

func NewDecoderSize(r io.Reader, sz int) *MsgReader {
	return &MsgReader{
		r:     bufio.NewReaderSize(r, sz),
		under: r,
	}
}

type MsgReader struct {
	r       *bufio.Reader
	under   io.Reader
	leader  [18]byte
	scratch []byte // recycled []byte for temporary storage
}

// Read implements the standard io.Reader method
func (m *MsgReader) Read(p []byte) (int, error) {
	return m.r.Read(p)
}

// ReadFull implements io.ReadFull
func (m *MsgReader) ReadFull(p []byte) (int, error) {
	return io.ReadFull(m.r, p)
}

func (m *MsgReader) Reset(r io.Reader) {
	m.r.Reset(r)
	m.under = r
}

// is the next byte 'nil'?
func (m *MsgReader) IsNil() bool {
	k, _ := m.nextKind()
	return k == knull
}

// Skip skips over the next object
func (m *MsgReader) Skip() (n int, err error) {
	var k kind
	var nn int
	var sz uint32
	k, err = m.nextKind()
	if err != nil {
		return
	}
	switch k {
	case kmap:
		sz, nn, err = m.ReadMapHeader()
		n += nn
		if err != nil {
			return
		}
		// maps have 2n fields
		sz *= 2
		for i := uint32(0); i < sz; i++ {
			nn, err = m.Skip()
			n += nn
			if err != nil {
				return
			}
		}
		return

	case karray:
		_, nn, err = m.ReadArrayHeader()
		n += nn
		if err != nil {
			return
		}
		for i := uint32(0); i < sz; i++ {
			nn, err = m.Skip()
			n += nn
			if err != nil {
				return
			}
		}
		return

	case kint:
		_, n, err = m.ReadInt64()
		return

	case kuint:
		_, n, err = m.ReadUint64()
		return

	case kfloat32:
		_, n, err = m.ReadFloat32()
		return

	case kfloat64:
		_, n, err = m.ReadFloat64()
		return

	case kextension:
		return m.skipExtension()

	case kbytes:
		return m.skipBytes()

	case kstring:
		return m.skipString()

	case knull:
		return m.ReadNil()

	case kbool:
		_, n, err = m.ReadBool()
		return

	case invalid:
		err = fmt.Errorf("unknown leading byte: 0x%x", byte(k))
		return

	default:
		err = errors.New("unknown leading type...?")
		return
	}
}

func (m *MsgReader) skipExtension() (n int, err error) {
	var lead byte
	var nn int
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
	case mfixext1:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		return

	case mfixext2:
		nn, err = io.ReadFull(m.r, m.leader[:3])
		n += nn
		return

	case mfixext4:
		nn, err = io.ReadFull(m.r, m.leader[:5])
		n += nn
		return

	case mfixext8:
		nn, err = io.ReadFull(m.r, m.leader[:9])
		n += nn
		return

	case mfixext16:
		nn, err = io.ReadFull(m.r, m.leader[:17])
		n += nn
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

	default:
		err = fmt.Errorf("unexpected byte 0x%x for extension", lead)
		return

	}

	var cn int64
	cn, err = io.CopyN(ioutil.Discard, m.r, int64(read))
	n += int(cn)
	return
}

func (m *MsgReader) skipBytes() (n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
	case mnil:
		return
	case mbin8:
		nn, err = io.ReadFull(m.r, m.leader[:1])
		n += nn
		if err != nil {
			return
		}
		read = int(m.leader[0])
	case mbin16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint16(m.leader[:]))
	case mbin32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[:]))
	default:
		err = fmt.Errorf("bad byte %x for []byte", m.leader[0])
		return
	}
	var cn int64
	cn, err = io.CopyN(ioutil.Discard, m.r, int64(read))
	n += int(cn)
	return
}

func (m *MsgReader) skipString() (n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
	case mnil:
		return
	case mstr8:
		nn, err = io.ReadFull(m.r, m.leader[:1])
		n += nn
		if err != nil {
			return
		}
		read = int(m.leader[0])
	case mstr16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint16(m.leader[:]))
	case mstr32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[:]))
	default:
		// try fixstr - first bits should be 101
		if isfixstr(lead) {
			read = int(rfixstr(lead))
		} else {
			err = fmt.Errorf("unexpected byte %x for string", lead)
			return
		}
	}
	if read == 0 {
		return
	}
	var cn int64
	cn, err = io.CopyN(ioutil.Discard, m.r, int64(read))
	n += int(cn)
	return
}

func (m *MsgReader) ReadMapHeader() (sz uint32, n int, err error) {
	var lead byte
	var nn int
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	if isfixmap(lead) {
		sz = uint32(rfixmap(lead))
		return
	}
	switch lead {
	case mnil:
		err = ErrNil
		return
	case mmap16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		usz := binary.BigEndian.Uint16(m.leader[:])
		sz = uint32(usz)
		return
	case mmap32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		sz = binary.BigEndian.Uint32(m.leader[:])
		return
	default:
		err = fmt.Errorf("unexpected byte 0x%x for map", lead)
		return
	}
}

func (m *MsgReader) ReadArrayHeader() (sz uint32, n int, err error) {
	var lead byte
	var nn int
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	if isfixarray(lead) {
		sz = uint32(rfixarray(lead))
		return
	}
	switch lead {
	case mnil:
		return // sz = 0

	case marray16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		usz := binary.BigEndian.Uint16(m.leader[:])
		sz = uint32(usz)
		return

	case marray32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		sz = binary.BigEndian.Uint32(m.leader[:])
		return

	default:
		err = fmt.Errorf("unexpected byte 0x%x for array", lead)
		return
	}
}

func (m *MsgReader) ReadNil() (int, error) {
	lead, err := m.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if lead != mnil {
		return 1, fmt.Errorf("unexpected byte %x for Nil", lead)
	}
	return 1, nil
}

func (m *MsgReader) ReadFloat64() (f float64, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:9])
	if err != nil {
		return
	}
	if m.leader[0] != mfloat64 {
		if m.leader[0] == mnil {
			err = ErrNil
			return
		}
		err = fmt.Errorf("msgp/enc: unexpected byte %x for Float64; expected %x", m.leader[0], mfloat64)
		return
	}
	bits := binary.BigEndian.Uint64(m.leader[1:])
	f = *(*float64)(unsafe.Pointer(&bits))
	return
}

func (m *MsgReader) ReadFloat32() (f float32, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:5])
	if err != nil {
		return
	}
	if m.leader[0] != mfloat32 {
		if m.leader[0] == mnil {
			err = ErrNil
			return
		}
		err = fmt.Errorf("msgp/enc: unexpected byte %x for Float32; expected %x", m.leader[0], mfloat64)
		return
	}
	bits := binary.BigEndian.Uint32(m.leader[1:])
	f = *(*float32)(unsafe.Pointer(&bits))
	return
}

func (m *MsgReader) ReadBool() (bool, int, error) {
	lead, err := m.r.ReadByte()
	if err != nil {
		return false, 0, err
	}
	switch lead {
	case mnil:
		return false, 1, ErrNil
	case mtrue:
		return true, 1, nil
	case mfalse:
		return false, 1, nil
	default:
		return false, 1, fmt.Errorf("unexpected byte %x for bool", lead)
	}
}

func (m *MsgReader) ReadInt64() (i int64, n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	if isfixint(lead) {
		i = int64(rfixint(lead))
		return
	}
	// try to decode negative fixnum
	if isnfixint(lead) {
		i = int64(rnfixint(lead))
		return
	}
	switch lead {
	case mnil:
		err = ErrNil
		return

	case mint8:
		var next byte
		next, err = m.r.ReadByte()
		if err != nil {
			return
		}
		n += 1
		i = int64(int8(next))
		return

	case mint16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		i = int64((int16(m.leader[0]) << 8) | int16(m.leader[1]))
		return

	case mint32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		i = int64((int32(m.leader[0]) << 24) | (int32(m.leader[1]) << 16) | (int32(m.leader[2]) << 8) | (int32(m.leader[3])))
		return

	case mint64:
		nn, err = io.ReadFull(m.r, m.leader[:8])
		n += nn
		i |= int64(m.leader[0]) << 56
		i |= int64(m.leader[1]) << 48
		i |= int64(m.leader[2]) << 40
		i |= int64(m.leader[3]) << 32
		i |= int64(m.leader[4]) << 24
		i |= int64(m.leader[5]) << 16
		i |= int64(m.leader[6]) << 8
		i |= int64(m.leader[7])
		return

	default:
		err = fmt.Errorf("unknown leading byte 0x%x for int", lead)
		return
	}
}

func (m *MsgReader) ReadInt32() (i int32, n int, err error) {
	var in int64
	in, n, err = m.ReadInt64()
	if in > math.MaxInt32 || in < math.MinInt32 {
		err = fmt.Errorf("%d overflows int32", in)
		return
	}
	i = int32(in)
	return
}

func (m *MsgReader) ReadInt16() (i int16, n int, err error) {
	var in int64
	in, n, err = m.ReadInt64()
	if in > math.MaxInt16 || in < math.MinInt16 {
		err = fmt.Errorf("%d overflows int16", in)
		return
	}
	i = int16(in)
	return
}

func (m *MsgReader) ReadInt8() (i int8, n int, err error) {
	var in int64
	in, n, err = m.ReadInt64()
	if in > math.MaxInt8 || in < math.MinInt8 {
		err = fmt.Errorf("%d overflows in8", in)
		return
	}
	i = int8(in)
	return
}

func (m *MsgReader) ReadInt() (i int, n int, err error) {
	var in int64
	var sin int32
	switch unsafe.Sizeof(i) {
	case unsafe.Sizeof(in):
		in, n, err = m.ReadInt64()
		i = int(in)
		return
	case unsafe.Sizeof(sin):
		sin, n, err = m.ReadInt32()
		i = int(sin)
		return
	}
	panic("impossible")
}

func (m *MsgReader) ReadUint64() (u uint64, n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	if isfixint(lead) {
		u = uint64(rfixint(lead))
		return
	}
	switch lead {
	case mnil:
		err = ErrNil
		return

	case muint8:
		var next byte
		next, err = m.r.ReadByte()
		if err != nil {
			return
		}
		n += 1
		u = uint64(next)
		return

	case muint16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		usz := binary.BigEndian.Uint16(m.leader[:])
		u = uint64(usz)
		return

	case muint32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		usz := binary.BigEndian.Uint32(m.leader[:])
		u = uint64(usz)
		return

	case muint64:
		nn, err = io.ReadFull(m.r, m.leader[:8])
		n += nn
		u = binary.BigEndian.Uint64(m.leader[:])
		return

	default:
		err = fmt.Errorf("unexpected byte %x for Uint", lead)
		return

	}
}

func (m *MsgReader) ReadUint32() (u uint32, n int, err error) {
	var in uint64
	in, n, err = m.ReadUint64()
	if in > math.MaxUint32 {
		err = fmt.Errorf("%d overflows uint32", in)
		return
	}
	u = uint32(in)
	return
}

func (m *MsgReader) ReadUint16() (u uint16, n int, err error) {
	var in uint64
	in, n, err = m.ReadUint64()
	if in > math.MaxUint16 {
		err = fmt.Errorf("%d overflows uint16", in)
	}
	u = uint16(in)
	return
}

func (m *MsgReader) ReadUint8() (u uint8, n int, err error) {
	var in uint64
	in, n, err = m.ReadUint64()
	if in > math.MaxUint8 {
		err = fmt.Errorf("%d overflows uint8", in)
	}
	u = uint8(in)
	return
}

func (m *MsgReader) ReadUint() (u uint, n int, err error) {
	var l uint64
	var s uint32
	switch unsafe.Sizeof(u) {
	case unsafe.Sizeof(s):
		s, n, err = m.ReadUint32()
		u = uint(s)
		return
	case unsafe.Sizeof(l):
		l, n, err = m.ReadUint64()
		u = uint(l)
		return
	}
	panic("impossible")
}

func (m *MsgReader) ReadBytes(scratch []byte) (b []byte, n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
	case mnil:
		return
	case mbin8:
		nn, err = io.ReadFull(m.r, m.leader[:1])
		n += nn
		if err != nil {
			return
		}
		read = int(m.leader[0])
	case mbin16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint16(m.leader[:]))
	case mbin32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[:]))
	default:
		err = fmt.Errorf("bad byte %x for []byte", m.leader[0])
		return
	}
	b, nn, err = readN(m.r, scratch, read)
	n += nn
	return
}

// opportunistic readfull
func readN(r io.Reader, scratch []byte, c int) (b []byte, n int, err error) {
	if scratch == nil || cap(scratch) == 0 || c > cap(scratch) {
		b = make([]byte, c)
		n, err = io.ReadFull(r, b)
		return
	} else {
		n, err = io.ReadFull(r, scratch[0:c])
		b = scratch[0:c]
		return
	}
}

func (m *MsgReader) ReadStringAsBytes(scratch []byte) (b []byte, n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
	case mnil:
		return
	case mstr8:
		nn, err = io.ReadFull(m.r, m.leader[:1])
		n += nn
		if err != nil {
			return
		}
		read = int(m.leader[0])
	case mstr16:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint16(m.leader[:]))
	case mstr32:
		nn, err = io.ReadFull(m.r, m.leader[:4])
		n += nn
		if err != nil {
			return
		}
		read = int(binary.BigEndian.Uint32(m.leader[:]))
	default:
		// try fixstr - first bits should be 101
		if isfixstr(lead) {
			read = int(rfixstr(lead))
		} else {
			err = fmt.Errorf("unexpected byte %x for string", lead)
			return
		}
	}
	if read == 0 {
		return
	}

	b, nn, err = readN(m.r, scratch, read)
	n += nn
	return
}

func (m *MsgReader) ReadString() (string, int, error) {
	var n int
	var err error
	m.scratch, n, err = m.ReadStringAsBytes(m.scratch)
	return string(m.scratch), n, err
}

func (m *MsgReader) ReadComplex64() (f complex64, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:10])
	if err != nil {
		return
	}
	if m.leader[0] != mfixext8 {
		if m.leader[0] == mnil {
			err = ErrNil
			return
		}
		err = fmt.Errorf("unexpected byte %x for complex64", m.leader[0])
	}
	if m.leader[1] != Complex64Extension {
		err = fmt.Errorf("unexpected byte %x for complex64 extension", m.leader[1])
	}
	rlb := binary.BigEndian.Uint32(m.leader[2:])
	imb := binary.BigEndian.Uint32(m.leader[6:])
	f = complex(*(*float32)(unsafe.Pointer(&rlb)), *(*float32)(unsafe.Pointer(&imb)))
	return
}

func (m *MsgReader) ReadComplex128() (f complex128, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:18])
	if err != nil {
		return
	}
	if m.leader[0] != mfixext16 {
		if m.leader[0] == mnil {
			err = ErrNil
			return
		}
		err = fmt.Errorf("unexpected byte %x for complex128", m.leader[0])
		return
	}
	if m.leader[1] != Complex128Extension {
		err = fmt.Errorf("unexpected byte %x for complex128 extension", m.leader[1])
		return
	}
	rlb := binary.BigEndian.Uint64(m.leader[2:])
	imb := binary.BigEndian.Uint64(m.leader[10:])
	f = complex(*(*float64)(unsafe.Pointer(&rlb)), *(*float64)(unsafe.Pointer(&imb)))
	return
}

// ReadMapStrStr CANNOT be passed a nil map
func (m *MsgReader) ReadMapStrStr(mp map[string]string) (n int, err error) {
	var nn int
	var sz uint32
	sz, nn, err = m.ReadMapHeader()
	n += nn
	if err == ErrNil || sz == 0 {
		err = nil
		mp = nil
		return
	}
	if err != nil {
		return
	}
	for key, _ := range mp {
		delete(mp, key)
	}
	for i := uint32(0); i < sz; i++ {
		var key string
		var val string
		key, nn, err = m.ReadString()
		n += nn
		if err != nil {
			return
		}
		val, nn, err = m.ReadString()
		n += nn
		if err != nil {
			return
		}
		mp[key] = val
	}
	return
}

type Extension struct {
	Type byte
	Data []byte
}

func (m *MsgReader) ReadExtension(e *Extension) (n int, err error) {
	var lead byte
	var nn int
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
	case mfixext1:
		nn, err = io.ReadFull(m.r, m.leader[:2])
		n += nn
		if err != nil {
			return
		}
		e.Type = m.leader[0]
		e.Data = append(e.Data[0:0], m.leader[1])
		return

	case mfixext2:
		nn, err = io.ReadFull(m.r, m.leader[:3])
		n += nn
		if err != nil {
			return
		}
		e.Type = m.leader[0]
		e.Data = append(e.Data[0:0], m.leader[1:3]...)
		return

	case mfixext4:
		nn, err = io.ReadFull(m.r, m.leader[:5])
		n += nn
		if err != nil {
			return
		}
		e.Type = m.leader[0]
		e.Data = append(e.Data[0:0], m.leader[1:5]...)
		return

	case mfixext8:
		nn, err = io.ReadFull(m.r, m.leader[:9])
		n += nn
		if err != nil {
			return
		}
		e.Type = m.leader[0]
		e.Data = append(e.Data[0:0], m.leader[1:9]...)
		return

	case mfixext16:
		nn, err = io.ReadFull(m.r, m.leader[:17])
		n += nn
		if err != nil {
			return
		}
		e.Type = m.leader[0]
		e.Data = append(e.Data[0:0], m.leader[1:17]...)
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

	e.Data, nn, err = readN(m.r, e.Data, read)
	n += nn
	return
}

func (m *MsgReader) ReadMapStrIntf(mp map[string]interface{}) (n int, err error) {
	var nn int
	var sz uint32
	sz, nn, err = m.ReadMapHeader()
	n += nn
	if err == ErrNil || sz == 0 {
		err = nil
		mp = nil
		return
	}
	if err != nil {
		return
	}
	if mp != nil {
		for key, _ := range mp {
			delete(mp, key)
		}
	} else {
		mp = make(map[string]interface{}, int(sz))
	}
	for i := uint32(0); i < sz; i++ {
		var key string
		var val interface{}
		key, nn, err = m.ReadString()
		n += nn
		if err != nil {
			return
		}
		val, nn, err = m.readInterface()
		n += nn
		if err != nil {
			return
		}
		mp[key] = val
	}
	return
}

func (m *MsgReader) ReadTime() (t time.Time, n int, err error) {
	// read all 18 bytes
	n, err = io.ReadFull(m.r, m.leader[:])
	if err != nil {
		return
	}
	if m.leader[0] != mfixext16 {
		err = fmt.Errorf("unexpected byte 0x%x for time.Time", m.leader[0])
		return
	}
	if m.leader[1] != TimeExtension {
		err = fmt.Errorf("unexpected extension type %d for time.Time", m.leader[1])
		return
	}
	err = t.UnmarshalBinary(m.leader[2:17]) // wants 15 bytes; last byte is 0
	return
}

func (m *MsgReader) ReadIdent(d MsgDecoder) (n int, err error) {
	return d.DecodeMsg(m.r)
}

func (m *MsgReader) ReadIntf() (i interface{}, n int, err error) {
	return m.readInterface()
}

func (m *MsgReader) readInterface() (i interface{}, n int, err error) {
	var k kind
	k, err = m.nextKind()
	if err != nil {
		return
	}
	switch k {
	case kbool:
		i, n, err = m.ReadBool()
		return

	case kint:
		i, n, err = m.ReadInt64()
		return

	case kuint:
		i, n, err = m.ReadUint64()
		return

	case kbytes:
		i, n, err = m.ReadBytes(nil)
		return

	case kstring:
		i, n, err = m.ReadString()
		return

	case kextension:
		e := new(Extension)
		n, err = m.ReadExtension(e)
		if err != nil {
			return
		}
		if e.Type == Complex128Extension && len(e.Data) == 16 {
			rlbits := binary.BigEndian.Uint64(e.Data[0:])
			imbits := binary.BigEndian.Uint64(e.Data[8:])
			rl := *(*float64)(unsafe.Pointer(&rlbits))
			im := *(*float64)(unsafe.Pointer(&imbits))
			i = complex(rl, im)
			return
		}
		if e.Type == Complex64Extension && len(e.Data) == 8 {
			rlbits := binary.BigEndian.Uint64(e.Data[0:])
			imbits := binary.BigEndian.Uint64(e.Data[8:])
			rl := *(*float32)(unsafe.Pointer(&rlbits))
			im := *(*float32)(unsafe.Pointer(&imbits))
			i = complex(rl, im)
			return
		}
		if e.Type == TimeExtension && len(e.Data) == 16 {
			var t time.Time
			err = t.UnmarshalBinary(e.Data[:15])
			i = t
			return
		}
		i = e
		return

	case kmap:
		mp := make(map[string]interface{})
		n, err = m.ReadMapStrIntf(mp)
		i = mp
		return

	case knull:
		n, err = m.ReadNil()
		i = nil
		return

	case kfloat32:
		i, n, err = m.ReadFloat32()
		return

	case kfloat64:
		i, n, err = m.ReadFloat64()
		return

	case karray:
		var sz uint32
		var nn int
		sz, nn, err = m.ReadArrayHeader()
		n += nn
		if err != nil {
			return
		}
		out := make([]interface{}, int(sz))
		for j := range out {
			out[j], nn, err = m.readInterface()
			n += nn
			if err != nil {
				return
			}
		}
		i = out
		return

	default:
		return nil, 0, errors.New("unrecognized type")

	}
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
