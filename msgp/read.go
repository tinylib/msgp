package msgp

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
	ErrNil = errors.New("value encoded as nil")

	readerPool sync.Pool
)

// ArrayError is an error returned
// when decoding a fix-sized array
// of the wrong size
type ArrayError struct {
	Wanted uint32
	Got    uint32
}

func (a ArrayError) Error() string {
	return fmt.Sprintf("wanted array of size %d; got %d", a.Wanted, a.Got)
}

type Type byte

const (
	InvalidType Type = iota
	StrType
	BinType
	MapType
	ArrayType
	Float64Type
	Float32Type
	BoolType
	IntType
	UintType
	NilType
	ExtensionType
)

func (t Type) String() string {
	switch t {
	case StrType:
		return "str"
	case BinType:
		return "bin"
	case MapType:
		return "map"
	case ArrayType:
		return "array"
	case Float64Type:
		return "float64"
	case Float32Type:
		return "float32"
	case BoolType:
		return "bool"
	case UintType:
		return "uint"
	case IntType:
		return "int"
	case ExtensionType:
		return "ext"
	case NilType:
		return "nil"
	default:
		return "<invalid>"
	}
}

// TypeError represents a decoding type error
type TypeError struct {
	Type   Type
	Prefix byte
}

func (t TypeError) Error() string {
	return fmt.Sprintf("msgp: prefix 0x%x invalid for type %q", t.Prefix, t.Type)
}

// InvalidPrefixError is returned when a bad encoding
// uses a prefix that is not recognized in the MessagePack standard
type InvalidPrefixError byte

func (i InvalidPrefixError) Error() string {
	return fmt.Sprintf("msgp: unrecognized type prefix 0x%x", byte(i))
}

func init() {
	readerPool.New = func() interface{} {
		return &Reader{}
	}
}

func popReader(r io.Reader) *Reader {
	p := readerPool.Get().(*Reader)
	if p.r == nil {
		p.r = bufio.NewReaderSize(r, 32)
		return p
	}
	p.Reset(r)
	return p
}

func pushReader(m *Reader) {
	readerPool.Put(m)
}

// Done recycles a reader
// to be used by other processes.
func Done(m *Reader) {
	if m == nil {
		return
	}
	pushReader(m)
}

// Unmarshaler is the interface fulfilled
// byte objects that know how to unmarshal
// themselves from MessagePack
type Unmarshaler interface {
	// UnmarshalMsg unmarshals the object
	// from binary, returing any leftover
	// bytes and any errors encountered
	UnmarshalMsg([]byte) ([]byte, error)
}

// Decoder is the interface fulfilled
// by objects that know how to decode themselves
// from MessagePack
type Decoder interface {
	// DecodeMsg decodes the object
	// from an io.Reader, returning
	// the number of bytes read and
	// any errors encountered
	DecodeMsg(io.Reader) (int, error)

	// DecodeFrom decodes the object
	// using an existing *Reader,
	// returning the number of bytes
	// read and any errors encountered
	DecodeFrom(*Reader) (int, error)
}

// NewDecoder returns a *Reader that
// reads from the provided reader.
func NewReader(r io.Reader) *Reader {
	return popReader(r)
}

// NewReaderSize returns a *Reader with
// a buffer of the given size. (This is preferable to
// passing the decoder a reader that is already buffered.)
func NewReaderSize(r io.Reader, sz int) *Reader {
	return &Reader{
		r: bufio.NewReaderSize(r, sz),
	}
}

// Reader wraps an io.Reader and provides
// methods to read MessagePack-encoded values
// from it. Readers are buffered.
type Reader struct {
	r       *bufio.Reader
	leader  [18]byte
	scratch []byte // recycled []byte for temporary storage
}

// Read implements the standard io.Reader method
func (m *Reader) Read(p []byte) (int, error) {
	return m.r.Read(p)
}

// ReadFull implements io.ReadFull
func (m *Reader) ReadFull(p []byte) (int, error) {
	return io.ReadFull(m.r, p)
}

// Reset resets the underlying reader
func (m *Reader) Reset(r io.Reader) {
	m.r.Reset(r)
}

// is the next byte 'nil'?
func (m *Reader) IsNil() bool {
	k, _ := m.nextKind()
	return k == knull
}

// Skip skips over the next object
func (m *Reader) Skip() (n int, err error) {
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

	default:
		// we shouldn't get here; nextKind() errors
		err = errors.New("msgp: Skip() decode error")
		return
	}
}

func (m *Reader) skipExtension() (n int, err error) {
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
		err = TypeError{Type: ExtensionType, Prefix: lead}
		return

	}

	var cn int64
	cn, err = io.CopyN(ioutil.Discard, m.r, int64(read))
	n += int(cn)
	return
}

func (m *Reader) skipBytes() (n int, err error) {
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
		err = TypeError{Type: BinType, Prefix: m.leader[0]}
		return
	}
	var cn int64
	cn, err = io.CopyN(ioutil.Discard, m.r, int64(read))
	n += int(cn)
	return
}

func (m *Reader) skipString() (n int, err error) {
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
			err = TypeError{Type: StrType, Prefix: lead}
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

func (m *Reader) ReadMapHeader() (sz uint32, n int, err error) {
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
		err = TypeError{Type: MapType, Prefix: lead}
		return
	}
}

func (m *Reader) ReadMapKey(scratch []byte) ([]byte, int, error) {
	k, err := m.nextKind()
	if err != nil {
		return nil, 0, err
	}
	if k == kstring {
		return m.ReadStringAsBytes(scratch)
	} else if k == kbytes {
		return m.ReadBytes(scratch)
	}
	return nil, 0, fmt.Errorf("msgp: %q not convertible to map key (string)", k)
}

func (m *Reader) ReadArrayHeader() (sz uint32, n int, err error) {
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
		err = TypeError{Type: ArrayType, Prefix: lead}
		return
	}
}

func (m *Reader) ReadNil() (int, error) {
	lead, err := m.r.ReadByte()
	if err != nil {
		return 0, err
	}
	if lead != mnil {
		return 1, TypeError{Type: NilType, Prefix: lead}
	}
	return 1, nil
}

func (m *Reader) ReadFloat64() (f float64, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:9])
	if err != nil {
		return
	}
	if m.leader[0] != mfloat64 {
		err = TypeError{Type: Float64Type, Prefix: m.leader[0]}
		return
	}
	f = *(*float64)(unsafe.Pointer(&m.leader[1]))
	return
}

func (m *Reader) ReadFloat32() (f float32, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:5])
	if err != nil {
		return
	}
	if m.leader[0] != mfloat32 {
		err = TypeError{Type: Float32Type, Prefix: m.leader[0]}
		return
	}
	f = *(*float32)(unsafe.Pointer(&m.leader[1]))
	return
}

func (m *Reader) ReadBool() (bool, int, error) {
	lead, err := m.r.ReadByte()
	if err != nil {
		return false, 0, err
	}
	switch lead {
	case mtrue:
		return true, 1, nil
	case mfalse:
		return false, 1, nil
	default:
		return false, 1, TypeError{Type: BoolType, Prefix: lead}
	}
}

func (m *Reader) ReadInt64() (i int64, n int, err error) {
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
		err = TypeError{Type: IntType, Prefix: lead}
		return
	}
}

func (m *Reader) ReadInt32() (i int32, n int, err error) {
	var in int64
	in, n, err = m.ReadInt64()
	if in > math.MaxInt32 || in < math.MinInt32 {
		err = fmt.Errorf("msgp: %d overflows int32", in)
		return
	}
	i = int32(in)
	return
}

func (m *Reader) ReadInt16() (i int16, n int, err error) {
	var in int64
	in, n, err = m.ReadInt64()
	if in > math.MaxInt16 || in < math.MinInt16 {
		err = fmt.Errorf("msgp: %d overflows int16", in)
		return
	}
	i = int16(in)
	return
}

func (m *Reader) ReadInt8() (i int8, n int, err error) {
	var in int64
	in, n, err = m.ReadInt64()
	if in > math.MaxInt8 || in < math.MinInt8 {
		err = fmt.Errorf("msgp: %d overflows int8", in)
		return
	}
	i = int8(in)
	return
}

func (m *Reader) ReadInt() (i int, n int, err error) {
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
	panic("impossible Sizeof(int)")
}

func (m *Reader) ReadUint64() (u uint64, n int, err error) {
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
		err = TypeError{Type: UintType, Prefix: lead}
		return

	}
}

func (m *Reader) ReadUint32() (u uint32, n int, err error) {
	var in uint64
	in, n, err = m.ReadUint64()
	if in > math.MaxUint32 {
		err = fmt.Errorf("msgp: %d overflows uint32", in)
		return
	}
	u = uint32(in)
	return
}

func (m *Reader) ReadUint16() (u uint16, n int, err error) {
	var in uint64
	in, n, err = m.ReadUint64()
	if in > math.MaxUint16 {
		err = fmt.Errorf("msgp: %d overflows uint16", in)
	}
	u = uint16(in)
	return
}

func (m *Reader) ReadUint8() (u uint8, n int, err error) {
	var in uint64
	in, n, err = m.ReadUint64()
	if in > math.MaxUint8 {
		err = fmt.Errorf("msgp: %d overflows uint8", in)
	}
	u = uint8(in)
	return
}

func (m *Reader) ReadUint() (u uint, n int, err error) {
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
	panic("impossible Sizeof(uint)")
}

func (m *Reader) ReadBytes(scratch []byte) (b []byte, n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
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
		err = TypeError{Type: BinType, Prefix: lead}
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

func (m *Reader) ReadStringAsBytes(scratch []byte) (b []byte, n int, err error) {
	var nn int
	var lead byte
	lead, err = m.r.ReadByte()
	if err != nil {
		return
	}
	n += 1
	var read int
	switch lead {
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
			err = TypeError{Type: StrType, Prefix: lead}
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

func (m *Reader) ReadString() (string, int, error) {
	var n int
	var err error
	m.scratch, n, err = m.ReadStringAsBytes(m.scratch)
	return string(m.scratch), n, err
}

func (m *Reader) ReadComplex64() (f complex64, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:10])
	if err != nil {
		return
	}
	if m.leader[0] != mfixext8 {
		err = fmt.Errorf("msgp: unexpected byte %x for complex64", m.leader[0])
		return
	}
	if m.leader[1] != Complex64Extension {
		err = fmt.Errorf("msgp: unexpected byte %x for complex64 extension", m.leader[1])
	}
	f = *(*complex64)(unsafe.Pointer(&m.leader[2]))
	return
}

func (m *Reader) ReadComplex128() (f complex128, n int, err error) {
	n, err = io.ReadFull(m.r, m.leader[:18])
	if err != nil {
		return
	}
	if m.leader[0] != mfixext16 {
		err = fmt.Errorf("msgp: unexpected byte %x for complex128", m.leader[0])
		return
	}
	if m.leader[1] != Complex128Extension {
		err = fmt.Errorf("msgp: unexpected byte %x for complex128 extension", m.leader[1])
		return
	}
	f = *(*complex128)(unsafe.Pointer(&m.leader[2]))
	return
}

// ReadMapStrStr CANNOT be passed a nil map
func (m *Reader) ReadMapStrStr(mp map[string]string) (n int, err error) {
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
	var key []byte
	var val []byte
	for i := uint32(0); i < sz; i++ {
		// we'll accept both 'str' and 'bin'
		// for strings here, since some legacy
		// encodings use 'bin' for strings
		key, nn, err = m.ReadMapKey(key)
		n += nn
		if err != nil {
			return
		}
		val, nn, err = m.ReadMapKey(val)
		n += nn
		if err != nil {
			return
		}
		mp[string(key)] = string(val)
	}
	return
}

func (m *Reader) ReadMapStrIntf(mp map[string]interface{}) (n int, err error) {
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
	var scratch []byte
	for i := uint32(0); i < sz; i++ {
		var val interface{}
		scratch, nn, err = m.ReadMapKey(scratch)
		n += nn
		if err != nil {
			return
		}
		val, nn, err = m.readInterface()
		n += nn
		if err != nil {
			return
		}
		mp[string(scratch)] = val
	}
	return
}

func (m *Reader) ReadTime() (t time.Time, n int, err error) {
	// read all 18 bytes
	n, err = io.ReadFull(m.r, m.leader[:])
	if err != nil {
		return
	}
	if m.leader[0] != mfixext16 {
		err = fmt.Errorf("msgp: unexpected byte 0x%x for time.Time", m.leader[0])
		return
	}
	if int8(m.leader[1]) != TimeExtension {
		err = fmt.Errorf("msgp: unexpected extension type %d for time.Time", int8(m.leader[1]))
		return
	}
	err = t.UnmarshalBinary(m.leader[2:17]) // wants 15 bytes; last byte is 0
	return
}

func (m *Reader) ReadIdent(d Decoder) (n int, err error) {
	return d.DecodeFrom(m)
}

func (m *Reader) ReadIntf() (i interface{}, n int, err error) {
	return m.readInterface()
}

func (m *Reader) readInterface() (i interface{}, n int, err error) {
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
		var t int8
		t, err = m.peekExtensionType()
		if err != nil {
			return
		}
		switch t {
		case Complex64Extension:
			i, n, err = m.ReadComplex64()
			return
		case Complex128Extension:
			i, n, err = m.ReadComplex128()
			return
		case TimeExtension:
			i, n, err = m.ReadTime()
			return
		default:
			f, ok := extensionReg[t]
			if ok {
				e := f()
				n, err = m.ReadExtension(e)
				i = e
				return
			}
			var e RawExtension
			e.Type = t
			n, err = m.ReadExtension(&e)
			i = &e
			return
		}

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
		// we shouldn't get here; nextKind() will error
		return nil, 0, errors.New("msgp: unrecognized type error")

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
