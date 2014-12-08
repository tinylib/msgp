package msgp

import (
	"errors"
	"github.com/philhofer/fwd"
	"io"
	"math"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

// where we keep old *Readers
var readerPool = sync.Pool{New: func() interface{} { return &Reader{} }}

// Type is a MessagePack wire type,
// including this package's built-in
// extension types.
type Type byte

// MessagePack Types
const (
	InvalidType Type = iota

	// MessagePack built-in types

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

	// pseudo-types provided
	// by extensions

	Complex64Type
	Complex128Type
	TimeType
)

// String implements fmt.Stringer
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

// FreeR frees a reader for use
// by other processes. It is not necessary
// to call FreeR on a reader. However, maintaining
// a reference to a *Reader after calling FreeR on
// it will cause undefined behavior. Typically, this
// function should only be used by the code generator.
func FreeR(m *Reader) {
	if m != nil {
		readerPool.Put(m)
	}
}

// Unmarshaler is the interface fulfilled
// by objects that know how to unmarshal
// themselves from MessagePack.
// UnmarshalMsg unmarshals the object
// from binary, returing any leftover
// bytes and any errors encountered.
type Unmarshaler interface {
	UnmarshalMsg([]byte) ([]byte, error)
}

// Decodable is the interface fulfilled
// by objects that know how to read
// themselves from a *Reader.
type Decodable interface {
	DecodeMsg(*Reader) error
}

// Decode decodes 'd' from 'r'.
func Decode(r io.Reader, d Decodable) error {
	rd := NewReader(r)
	err := d.DecodeMsg(rd)
	FreeR(rd)
	return err
}

// NewReader returns a *Reader that
// reads from the provided reader. The
// reader will be buffered.
func NewReader(r io.Reader) *Reader {
	p := readerPool.Get().(*Reader)
	if p.r == nil {
		p.r = fwd.NewReaderSize(r, 1024)
		return p
	}
	p.Reset(r)
	return p
}

// NewReaderSize returns a *Reader with a buffer of the given size.
// (This is vastly preferable to passing the decoder a reader that is already buffered.)
func NewReaderSize(r io.Reader, sz int) *Reader {
	return &Reader{r: fwd.NewReaderSize(r, sz)}
}

// Reader wraps an io.Reader and provides
// methods to read MessagePack-encoded values
// from it. Readers are buffered.
type Reader struct {
	r       *fwd.Reader
	scratch []byte // recycled []byte for temporary storage
}

// Read implements `io.Reader`
func (m *Reader) Read(p []byte) (int, error) {
	return m.r.Read(p)
}

// ReadFull implements `io.ReadFull`
func (m *Reader) ReadFull(p []byte) (int, error) {
	return m.r.ReadFull(p)
}

// Reset resets the underlying reader.
func (m *Reader) Reset(r io.Reader) { m.r.Reset(r) }

// Buffered returns the number of bytes currently in the read buffer.
func (m *Reader) Buffered() int { return m.r.Buffered() }

// NextType returns the next object type to be decoded.
func (m *Reader) NextType() (Type, error) {
	p, err := m.r.Peek(1)
	if err != nil {
		return InvalidType, err
	}
	t := getType(p[0])
	if t == InvalidType {
		return InvalidType, InvalidPrefixError(p[0])
	}
	if t == ExtensionType {
		v, err := m.peekExtensionType()
		if err != nil {
			return InvalidType, err
		}
		switch v {
		case Complex64Extension:
			return Complex64Type, nil
		case Complex128Extension:
			return Complex128Type, nil
		case TimeExtension:
			return TimeType, nil
		}
	}
	return t, nil
}

// IsNil returns whether or not
// the next byte is a null messagepack byte
func (m *Reader) IsNil() bool {
	p, err := m.r.Peek(1)
	return err == nil && p[0] == mnil
}

// returns (obj size, obj elements, error)
// only maps and arrays have non-zero obj elements
func getNextSize(r *fwd.Reader) (uint32, uint32, error) {
	var (
		p   []byte
		err error
	)
	p, err = r.Peek(1)
	if err != nil {
		return 0, 0, err
	}
	lead := p[0]

	switch {
	case isfixarray(lead):
		return 1, uint32(rfixarray(lead)), nil

	case isfixmap(lead):
		return 1, 2 * uint32(rfixmap(lead)), nil

	case isfixstr(lead):
		return uint32(rfixstr(lead)) + 1, 0, nil

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
		p, err = r.Peek(2)
		if err != nil {
			return 0, 0, err
		}
		return uint32(uint8(p[1])) + 2, 0, nil

	case mbin16, mstr16:
		p, err = r.Peek(3)
		if err != nil {
			return 0, 0, err
		}
		return uint32(big.Uint16(p[1:])) + 3, 0, nil

	case mbin32, mstr32:
		p, err = r.Peek(5)
		if err != nil {
			return 0, 0, err
		}
		return big.Uint32(p[1:]) + 5, 0, nil

	// variable extensions
	// require 1 extra byte
	// to skip
	case mext8:
		p, err = r.Peek(3)
		if err != nil {
			return 0, 0, err
		}
		return uint32(uint8(p[1])) + 3, 0, nil

	case mext16:
		p, err = r.Peek(4)
		if err != nil {
			return 0, 0, err
		}
		return uint32(big.Uint16(p[1:])) + 4, 0, nil

	case mext32:
		p, err = r.Peek(6)
		if err != nil {
			return 0, 0, err
		}
		return big.Uint32(p[1:]) + 6, 0, nil

	// arrays skip lead byte,
	// size byte, N objects
	case marray16:
		p, err = r.Peek(3)
		if err != nil {
			return 0, 0, err
		}
		return 3, uint32(big.Uint16(p[1:])), nil

	case marray32:
		p, err = r.Peek(5)
		if err != nil {
			return 0, 0, err
		}
		return 5, big.Uint32(p[1:]), nil

	// maps skip lead byte,
	// size byte, 2N objects
	case mmap16:
		p, err = r.Peek(3)
		if err != nil {
			return 0, 0, err
		}
		return 3, 2 * uint32(big.Uint16(p[1:])), nil

	case mmap32:
		p, err = r.Peek(5)
		if err != nil {
			return 0, 0, err
		}
		return 5, 2 * big.Uint32(p[1:]), nil

	default:
		return 0, 0, InvalidPrefixError(lead)

	}
}

// Skip skips over the next object, regardless of
// its type. If it is an array or map, the whole array
// or map will be skipped.
func (m *Reader) Skip() error {
	var (
		v   uint32
		o   uint32
		err error
		p   []byte
	)

	// we can use the faster
	// method if we have enough
	// buffered data
	if m.r.Buffered() >= 5 {
		p, err = m.r.Peek(5)
		if err != nil {
			return err
		}
		v, o, err = getSize(p)
		if err != nil {
			return err
		}
	} else {
		v, o, err = getNextSize(m.r)
		if err != nil {
			return err
		}
	}

	// 'v' is always non-zero
	// if err == nil
	_, err = m.r.Skip(int(v))
	if err != nil {
		return err
	}

	// for maps and slices, skip elements
	for x := uint32(0); x < o; x++ {
		err = m.Skip()
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadMapHeader reads the next object
// as a map header and returns the size
// of the map and the number of bytes written.
// It will return a TypeError{} if the next
// object is not a map.
func (m *Reader) ReadMapHeader() (sz uint32, err error) {
	var p []byte
	var lead byte
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	lead = p[0]
	if isfixmap(lead) {
		sz = uint32(rfixmap(lead))
		_, err = m.r.Skip(1)
		return
	}
	switch lead {
	case mmap16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		usz := big.Uint16(p[1:])
		sz = uint32(usz)
		return
	case mmap32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		sz = big.Uint32(p[1:])
		return
	default:
		err = badPrefix(MapType, lead)
		return
	}
}

// ReadMapKey reads either a 'str' or 'bin' field from
// the reader and returns the value as a []byte. It uses
// scratch for storage if it is large enough.
func (m *Reader) ReadMapKey(scratch []byte) ([]byte, error) {
	out, err := m.ReadStringAsBytes(scratch)
	if err != nil {
		if tperr, ok := err.(TypeError); ok && tperr.Encoded == BinType {
			return m.ReadBytes(scratch)
		}
		return nil, err
	}
	return out, nil
}

// ReadArrayHeader reads the next object as an
// array header and returns the size of the array
// and the number of bytes read.
func (m *Reader) ReadArrayHeader() (sz uint32, err error) {
	var lead byte
	var p []byte
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	lead = p[0]
	if isfixarray(lead) {
		sz = uint32(rfixarray(lead))
		_, err = m.r.Skip(1)
		return
	}
	switch lead {
	case marray16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		sz = uint32(big.Uint16(p[1:]))
		return

	case marray32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		sz = big.Uint32(p[1:])
		return

	default:
		err = badPrefix(ArrayType, lead)
		return
	}
}

// ReadNil reads a 'nil' MessagePack byte from the reader
func (m *Reader) ReadNil() error {
	p, err := m.r.Peek(1)
	if err != nil {
		return err
	}
	if p[0] != mnil {
		return badPrefix(NilType, p[0])
	}
	_, err = m.r.Skip(1)
	return err
}

// ReadFloat64 reads a float64 from the reader.
// (If the value on the wire is encoded as a float32,
// it will be up-cast to a float64.)
func (m *Reader) ReadFloat64() (f float64, err error) {
	var p []byte
	p, err = m.r.Peek(9)
	if err != nil {
		// we'll allow a coversion from float32 to float64,
		// since we don't loose any precision
		if err == io.EOF && len(p) > 0 && p[0] == mfloat32 {
			ef, err := m.ReadFloat32()
			return float64(ef), err
		}
		return
	}
	if p[0] != mfloat64 {
		// see above
		if p[0] == mfloat32 {
			ef, err := m.ReadFloat32()
			return float64(ef), err
		}
		err = badPrefix(Float64Type, p[0])
		return
	}
	memcpy8(unsafe.Pointer(&f), unsafe.Pointer(&p[1]))
	_, err = m.r.Skip(9)
	return
}

// ReadFloat32 reads a float32 from the reader
func (m *Reader) ReadFloat32() (f float32, err error) {
	var p []byte
	p, err = m.r.Peek(5)
	if err != nil {
		return
	}
	if p[0] != mfloat32 {
		err = badPrefix(Float32Type, p[0])
		return
	}
	memcpy4(unsafe.Pointer(&f), unsafe.Pointer(&p[1]))
	_, err = m.r.Skip(5)
	return
}

// ReadBool reads a bool from the reader
func (m *Reader) ReadBool() (b bool, err error) {
	var p []byte
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	switch p[0] {
	case mtrue:
		b = true
	case mfalse:
	default:
		err = badPrefix(BoolType, p[0])
		return
	}
	_, err = m.r.Skip(1)
	return
}

// ReadInt64 reads an int64 from the reader
func (m *Reader) ReadInt64() (i int64, err error) {
	var p []byte
	var lead byte
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	lead = p[0]

	if isfixint(lead) {
		i = int64(rfixint(lead))
		_, err = m.r.Skip(1)
		return
	} else if isnfixint(lead) {
		i = int64(rnfixint(lead))
		_, err = m.r.Skip(1)
		return
	}

	switch lead {
	case mint8:
		p, err = m.r.Next(2)
		if err != nil {
			return
		}
		i = int64(getMint8(p))
		return

	case mint16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		i = int64(getMint16(p))
		return

	case mint32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		i = int64(getMint32(p))
		return

	case mint64:
		p, err = m.r.Next(9)
		if err != nil {
			return
		}
		i = getMint64(p)
		return

	default:
		err = badPrefix(IntType, lead)
		return
	}
}

// ReadInt32 reads an int32 from the reader
func (m *Reader) ReadInt32() (i int32, err error) {
	var in int64
	in, err = m.ReadInt64()
	if in > math.MaxInt32 || in < math.MinInt32 {
		err = IntOverflow{Value: in, FailedBitsize: 32}
		return
	}
	i = int32(in)
	return
}

// ReadInt16 reads an int16 from the reader
func (m *Reader) ReadInt16() (i int16, err error) {
	var in int64
	in, err = m.ReadInt64()
	if in > math.MaxInt16 || in < math.MinInt16 {
		err = IntOverflow{Value: in, FailedBitsize: 16}
		return
	}
	i = int16(in)
	return
}

// ReadInt8 reads an int8 from the reader
func (m *Reader) ReadInt8() (i int8, err error) {
	var in int64
	in, err = m.ReadInt64()
	if in > math.MaxInt8 || in < math.MinInt8 {
		err = IntOverflow{Value: in, FailedBitsize: 8}
		return
	}
	i = int8(in)
	return
}

// ReadInt reads an int from the reader
func (m *Reader) ReadInt() (i int, err error) {
	if unsafe.Sizeof(i) == 4 {
		var in int32
		in, err = m.ReadInt32()
		i = int(in)
		return
	}
	var in int64
	in, err = m.ReadInt64()
	i = int(in)
	return
}

// ReadUint64 reads a uint64 from the reader
func (m *Reader) ReadUint64() (u uint64, err error) {
	var p []byte
	var lead byte
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	lead = p[0]
	if isfixint(lead) {
		u = uint64(rfixint(lead))
		_, err = m.r.Skip(1)
		return
	}
	switch lead {
	case muint8:
		p, err = m.r.Next(2)
		if err != nil {
			return
		}
		u = uint64(getMuint8(p))
		return

	case muint16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		u = uint64(getMuint16(p))
		return

	case muint32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		u = uint64(getMuint32(p))
		return

	case muint64:
		p, err = m.r.Next(9)
		if err != nil {
			return
		}
		u = getMuint64(p)
		return

	default:
		err = badPrefix(UintType, lead)
		return

	}
}

// ReadUint32 reads a uint32 from the reader
func (m *Reader) ReadUint32() (u uint32, err error) {
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint32 {
		err = UintOverflow{Value: in, FailedBitsize: 32}
		return
	}
	u = uint32(in)
	return
}

// ReadUint16 reads a uint16 from the reader
func (m *Reader) ReadUint16() (u uint16, err error) {
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint16 {
		err = UintOverflow{Value: in, FailedBitsize: 16}
		return
	}
	u = uint16(in)
	return
}

// ReadUint8 reads a uint8 from the reader
func (m *Reader) ReadUint8() (u uint8, err error) {
	var in uint64
	in, err = m.ReadUint64()
	if in > math.MaxUint8 {
		err = UintOverflow{Value: in, FailedBitsize: 8}
		return
	}
	u = uint8(in)
	return
}

// ReadUint reads a uint from the reader
func (m *Reader) ReadUint() (u uint, err error) {
	if unsafe.Sizeof(u) == 4 {
		var un uint32
		un, err = m.ReadUint32()
		u = uint(un)
		return
	}
	var un uint64
	un, err = m.ReadUint64()
	u = uint(un)
	return
}

// ReadBytes reads a MessagePack 'bin' object
// from the reader and returns its value. It may
// use 'scratch' for storage if it is non-nil.
func (m *Reader) ReadBytes(scratch []byte) (b []byte, err error) {
	var p []byte
	var lead byte
	p, err = m.r.Peek(2)
	if err != nil {
		return
	}
	lead = p[0]
	var read int64
	switch lead {
	case mbin8:
		read = int64(p[1])
		m.r.Skip(2)
	case mbin16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		read = int64(big.Uint16(p[1:]))
	case mbin32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		read = int64(big.Uint32(p[1:]))
	default:
		err = badPrefix(BinType, lead)
		return
	}
	if int64(cap(scratch)) < read {
		b = make([]byte, read)
	} else {
		b = scratch[0:read]
	}
	_, err = m.r.ReadFull(b)
	return
}

// ReadStringAsBytes reads a MessagePack 'str' (utf-8) string
// and returns its value as bytes. It may use 'scratch' for storage
// if it is non-nil.
func (m *Reader) ReadStringAsBytes(scratch []byte) (b []byte, err error) {
	var p []byte
	var lead byte
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	lead = p[0]
	var read int64

	if isfixstr(lead) {
		read = int64(rfixstr(lead))
		m.r.Skip(1)
		goto fill
	}

	switch lead {
	case mstr8:
		p, err = m.r.Next(2)
		if err != nil {
			return
		}
		read = int64(uint8(p[1]))
	case mstr16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		read = int64(big.Uint16(p[1:]))
	case mstr32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		read = int64(big.Uint32(p[1:]))
	default:
		err = badPrefix(StrType, lead)
		return
	}
fill:
	if int64(cap(scratch)) < read {
		b = make([]byte, read)
	} else {
		b = scratch[0:read]
	}
	_, err = m.r.ReadFull(b)
	return
}

// ReadString reads a utf-8 string from the reader
func (m *Reader) ReadString() (s string, err error) {
	var p []byte
	var lead byte
	var read int64
	p, err = m.r.Peek(1)
	if err != nil {
		return
	}
	lead = p[0]

	if isfixstr(lead) {
		read = int64(rfixstr(lead))
		m.r.Skip(1)
		goto fill
	}

	switch lead {
	case mstr8:
		p, err = m.r.Next(2)
		if err != nil {
			return
		}
		read = int64(uint8(p[1]))
	case mstr16:
		p, err = m.r.Next(3)
		if err != nil {
			return
		}
		read = int64(big.Uint16(p[1:]))
	case mstr32:
		p, err = m.r.Next(5)
		if err != nil {
			return
		}
		read = int64(big.Uint32(p[1:]))
	default:
		err = badPrefix(StrType, lead)
		return
	}
fill:
	if read == 0 {
		s, err = "", nil
		return
	}
	// reading into the memory
	// that will become the string
	// itself has vastly superior
	// worst-case performance, because
	// the reader buffer doesn't have
	// to be large enough to hold the string.
	out := make([]byte, read)
	_, err = m.r.ReadFull(out)
	if err != nil {
		return
	}
	s = UnsafeString(out)
	return
}

// ReadComplex64 reads a complex64 from the reader
func (m *Reader) ReadComplex64() (f complex64, err error) {
	var p []byte
	p, err = m.r.Peek(10)
	if err != nil {
		return
	}
	if p[0] != mfixext8 {
		err = badPrefix(Complex64Type, p[0])
		return
	}
	if int8(p[1]) != Complex64Extension {
		err = errExt(int8(p[1]), Complex64Extension)
		return
	}
	memcpy8(unsafe.Pointer(&f), unsafe.Pointer(&p[2]))
	_, err = m.r.Skip(10)
	return
}

// ReadComplex128 reads a complex128 from the reader
func (m *Reader) ReadComplex128() (f complex128, err error) {
	var p []byte
	p, err = m.r.Peek(18)
	if err != nil {
		return
	}
	if p[0] != mfixext16 {
		err = badPrefix(Complex128Type, p[0])
		return
	}
	if int8(p[1]) != Complex128Extension {
		err = errExt(int8(p[1]), Complex128Extension)
		return
	}
	memcpy16(unsafe.Pointer(&f), unsafe.Pointer(&p[2]))
	_, err = m.r.Skip(18)
	return
}

// ReadMapStrIntf reads a MessagePack map into a map[string]interface{}.
// (You must pass a non-nil map into the function.)
func (m *Reader) ReadMapStrIntf(mp map[string]interface{}) (err error) {
	var sz uint32
	sz, err = m.ReadMapHeader()
	if err != nil {
		return
	}
	for key := range mp {
		delete(mp, key)
	}
	var scratch []byte
	for i := uint32(0); i < sz; i++ {
		var val interface{}
		scratch, err = m.ReadMapKey(scratch)
		if err != nil {
			return
		}
		val, err = m.ReadIntf()
		if err != nil {
			return
		}
		mp[string(scratch)] = val
	}
	return
}

// ReadTime reads a time.Time object from the reader.
func (m *Reader) ReadTime() (t time.Time, err error) {
	var p []byte
	p, err = m.r.Peek(18)
	if err != nil {
		return
	}
	if p[0] != mfixext16 {
		err = badPrefix(TimeType, p[0])
		return
	}
	if int8(p[1]) != TimeExtension {
		err = errExt(int8(p[1]), TimeExtension)
		return
	}
	err = t.UnmarshalBinary(p[2:17]) // wants 15 bytes; last byte is 0
	if err == nil {
		_, err = m.r.Skip(18)
	}
	return
}

// ReadIdent reads data into an object that implements the msgp.Decoder interface
func (m *Reader) ReadIdent(d Decodable) error {
	return d.DecodeMsg(m)
}

// ReadIntf reads out the next object as a raw interface{}.
// Arrays are decoded as []interface{}, and maps are decoded
// as map[string]interface{}. Integers are decoded as int64
// and unsigned integers are decoded as uint64.
func (m *Reader) ReadIntf() (i interface{}, err error) {
	var t Type
	t, err = m.NextType()
	if err != nil {
		return
	}
	switch t {
	case BoolType:
		i, err = m.ReadBool()
		return

	case IntType:
		i, err = m.ReadInt64()
		return

	case UintType:
		i, err = m.ReadUint64()
		return

	case BinType:
		i, err = m.ReadBytes(nil)
		return

	case StrType:
		i, err = m.ReadString()
		return

	case Complex64Type:
		i, err = m.ReadComplex64()
		return

	case Complex128Type:
		i, err = m.ReadComplex128()
		return

	case TimeType:
		i, err = m.ReadTime()
		return

	case ExtensionType:
		var t int8
		t, err = m.peekExtensionType()
		if err != nil {
			return
		}
		f, ok := extensionReg[t]
		if ok {
			e := f()
			err = m.ReadExtension(e)
			i = e
			return
		}
		var e RawExtension
		e.Type = t
		err = m.ReadExtension(&e)
		i = &e
		return

	case MapType:
		mp := make(map[string]interface{})
		err = m.ReadMapStrIntf(mp)
		i = mp
		return

	case NilType:
		err = m.ReadNil()
		i = nil
		return

	case Float32Type:
		i, err = m.ReadFloat32()
		return

	case Float64Type:
		i, err = m.ReadFloat64()
		return

	case ArrayType:
		var sz uint32
		sz, err = m.ReadArrayHeader()

		if err != nil {
			return
		}
		out := make([]interface{}, int(sz))
		for j := range out {
			out[j], err = m.ReadIntf()
			if err != nil {
				return
			}
		}
		i = out
		return

	default:
		// we shouldn't get here; NextType() will error
		return nil, errors.New("msgp: unrecognized type error")

	}
}

// UnsafeString returns the byte slice as a volatile string
// THIS SHOULD ONLY BE USED BY THE CODE GENERATOR
// THIS IS EVIL CODE
// YOU HAVE BEEN WARNED
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(&b[0])), Len: len(b)}))
}

// UnsafeBytes returns the string as a byte slice
// THIS SHOULD ONLY BE USED BY THE CODE GENERATOR
// THIS IS EVIL CODE
// YOU HAVE BEEN WARNED
func UnsafeBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Len:  len(s),
		Cap:  len(s),
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&s))).Data,
	}))
}
