package msgp

import (
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"sync"
	"time"
	"unsafe"
)

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

// Sizer is an interface implemented
// by types that can estimate their
// size when MessagePack encoded.
// This interface is optional, but
// encoding/marshaling implementations
// may use this as a way to pre-allocate
// memory for serialization.
type Sizer interface {
	Msgsize() int
}

var (
	// Nowhere is an io.Writer to nowhere
	Nowhere io.Writer = nwhere{}

	btsType    = reflect.TypeOf(([]byte)(nil))
	writerPool = sync.Pool{
		New: func() interface{} {
			return &Writer{buf: make([]byte, 0, 1024)}
		},
	}
)

func popWriter(w io.Writer) *Writer {
	wr := writerPool.Get().(*Writer)
	wr.Reset(w)
	return wr
}

func pushWriter(wr *Writer) {
	if wr != nil && cap(wr.buf) > 256 {
		wr.w = nil
		writerPool.Put(wr)
	}
}

// FreeW frees a writer for use
// by other processes. It is not necessary
// to call FreeW on a writer. However, maintaining
// a reference to a *Writer after calling FreeW on
// it will cause undefined behavior.
func FreeW(w *Writer) { pushWriter(w) }

// Require ensures that cap(old)-len(old) >= extra
func Require(old []byte, extra int) []byte {
	if cap(old)-len(old) >= extra {
		return old
	}
	if len(old) == 0 {
		return make([]byte, 0, extra)
	}
	n := make([]byte, len(old), cap(old)-len(old)+extra)
	copy(n, old)
	return n
}

// nowhere writer
type nwhere struct{}

func (n nwhere) Write(p []byte) (int, error) { return len(p), nil }

// Marshaler is the interface implemented
// by types that know how to marshal themselves
// as MessagePack. MarshalMsg appends the marshalled
// form of the object to the provided
// byte slice, returning the extended
// slice and any errors encountered.
type Marshaler interface {
	MarshalMsg([]byte) ([]byte, error)
}

// Encodable is the interface implemented
// by types that know how to write themselves
// as MessagePack
type Encodable interface {
	EncodeMsg(*Writer) error
}

// Writer is a buffered writer
// that can be used to write
// MessagePack objects to an io.Writer.
// You must call *Writer.Flush() in order
// to flush all of the buffered data
// to the underlying writer.
type Writer struct {
	w   io.Writer
	buf []byte // buffered data; [0:len(buf)] is valid
}

// NewWriter returns a new *Writer.
func NewWriter(w io.Writer) *Writer {
	if wr, ok := w.(*Writer); ok {
		return wr
	}
	return popWriter(w)
}

// NewWriterSize returns a writer with a custom buffer size.
func NewWriterSize(w io.Writer, sz int) *Writer {
	if sz < 16 {
		sz = 16
	}
	return &Writer{
		w:   w,
		buf: make([]byte, 0, sz),
	}
}

// Encode encodes an Encodable to an io.Writer.
func Encode(w io.Writer, e Encodable) error {
	wr := NewWriter(w)
	err := e.EncodeMsg(wr)
	if err == nil {
		err = wr.Flush()
	}
	FreeW(wr)
	return err
}

// Write writes a marshaler to an io.Writer.
func Write(w io.Writer, m Marshaler) error {
	wr := NewWriter(w)
	err := wr.Encode(m)
	if err == nil {
		err = wr.Flush()
	}
	FreeW(wr)
	return err
}

func (mw *Writer) flush() error {
	if len(mw.buf) > 0 {
		n, err := mw.w.Write(mw.buf)
		if err != nil {
			// copy unwritten data
			// back to index 0
			if n > 0 {
				mw.buf = mw.buf[:copy(mw.buf[0:], mw.buf[n:])]
			}
			return err
		}
		mw.buf = mw.buf[0:0]
		return nil
	}
	return nil
}

// Flush flushes all of the buffered
// data to the underlying writer.
func (mw *Writer) Flush() error { return mw.flush() }

// Buffered returns the number bytes in the write buffer
func (mw *Writer) Buffered() int { return len(mw.buf) }

func (mw *Writer) avail() int { return cap(mw.buf) - len(mw.buf) }

func (mw *Writer) require(n int) (int, error) {
	l := len(mw.buf)
	c := cap(mw.buf)
	if c-l >= n {
		mw.buf = mw.buf[:l+n] // grow by 'n'; return old offset
		return l, nil
	}
	err := mw.flush()
	if err != nil {
		return 0, err
	}
	// after flush,
	// len(mw.buf) = 0
	if n > c {
		mw.buf = make([]byte, n)
		return 0, nil
	}
	mw.buf = mw.buf[:n]
	return 0, nil
}

// push one byte onto the buffer
func (mw *Writer) push(b byte) error {
	l := len(mw.buf)
	if l == cap(mw.buf) {
		err := mw.flush()
		if err != nil {
			return err
		}
		l = 0
	}
	mw.buf = mw.buf[:l+1]
	mw.buf[l] = b
	return nil
}

// Write implements io.Writer, and writes
// data directly to the buffer.
func (mw *Writer) Write(p []byte) (int, error) {
	l := len(p)
	if mw.avail() >= l {
		o := len(mw.buf)
		mw.buf = mw.buf[:o+l]
		copy(mw.buf[o:], p)
		return l, nil
	}
	err := mw.flush()
	if err != nil {
		return 0, err
	}
	if l > cap(mw.buf) {
		return mw.w.Write(p)
	}
	mw.buf = mw.buf[:l]
	copy(mw.buf, p)
	return l, nil
}

// implements io.WriteString
func (mw *Writer) writeString(s string) error {
	l := len(s)

	// we have space; copy
	if mw.avail() >= l {
		o := len(mw.buf)
		mw.buf = mw.buf[:o+l]
		copy(mw.buf[o:], s)
		return nil
	}

	// we need to flush one way
	// or another.
	if err := mw.flush(); err != nil {
		return err
	}

	// shortcut: big strings go
	// straight to the underlying writer
	if l > cap(mw.buf) {
		_, err := io.WriteString(mw.w, s)
		return err
	}

	mw.buf = mw.buf[:l]
	copy(mw.buf, s)
	return nil
}

// Encode writes a Marshaler to the writer.
// Users should attempt to ensure that the
// encoded size of the marshaler is less than
// the total capacity of the buffer in order
// to avoid the writer having to re-allocate
// the entirety of the buffer.
func (mw *Writer) Encode(m Marshaler) error {
	if s, ok := m.(Sizer); ok {
		// check for available space
		sz := s.Msgsize()
		if sz > mw.avail() {
			err := mw.flush()
			if err != nil {
				return err
			}
			if sz > cap(mw.buf) {
				mw.buf = make([]byte, 0, sz)
			}
		}
	} else if mw.avail() < 3*cap(mw.buf)/4 {
		// flush if more than 3/4 full
		err := mw.flush()
		if err != nil {
			return err
		}
	}
	var err error
	old := mw.buf
	mw.buf, err = m.MarshalMsg(mw.buf)
	if err != nil {
		mw.buf = old
		return err
	}
	return nil
}

// Reset changes the underlying writer used by the MsgWriter
func (mw *Writer) Reset(w io.Writer) {
	mw.w = w
	mw.buf = mw.buf[0:0]
}

// WriteMapHeader writes a map header of the given
// size to the writer
func (mw *Writer) WriteMapHeader(sz uint32) error {
	switch {
	case sz < 16:
		return mw.push(wfixmap(uint8(sz)))

	case sz < 1<<16-1:
		o, err := mw.require(3)
		if err != nil {
			return err
		}
		prefixu16(mw.buf[o:], mmap16, uint16(sz))
		return nil

	default:
		o, err := mw.require(5)
		if err != nil {
			return err
		}
		prefixu32(mw.buf[o:], mmap32, sz)
		return nil
	}
}

// WriteArrayHeader writes an array header of the
// given size to the writer
func (mw *Writer) WriteArrayHeader(sz uint32) error {
	switch {
	case sz < 16:
		return mw.push(wfixarray(uint8(sz)))

	case sz < math.MaxUint16:
		o, err := mw.require(3)
		if err != nil {
			return err
		}
		prefixu16(mw.buf[o:], marray16, uint16(sz))
		return nil

	default:
		o, err := mw.require(5)
		if err != nil {
			return err
		}
		prefixu32(mw.buf[o:], marray32, sz)
		return nil
	}
}

// WriteNil writes a nil byte to the buffer
func (mw *Writer) WriteNil() error {
	return mw.push(mnil)
}

// WriteFloat64 writes a float64 to the writer
func (mw *Writer) WriteFloat64(f float64) error {
	o, err := mw.require(9)
	if err != nil {
		return err
	}
	mw.buf[o] = mfloat64
	memcpy8(unsafe.Pointer(&mw.buf[o+1]), unsafe.Pointer(&f))
	return nil
}

// WriteFloat32 writes a float32 to the writer
func (mw *Writer) WriteFloat32(f float32) error {
	o, err := mw.require(5)
	if err != nil {
		return err
	}
	mw.buf[o] = mfloat32
	memcpy4(unsafe.Pointer(&mw.buf[o+1]), unsafe.Pointer(&f))
	return nil
}

// WriteInt64 writes an int64 to the writer
func (mw *Writer) WriteInt64(i int64) error {
	a := abs(i)
	switch {
	case i < 0 && i > -32:
		return mw.push(wnfixint(int8(i)))

	case i >= 0 && i < 128:
		return mw.push(wfixint(uint8(i)))

	case a < math.MaxInt8:
		o, err := mw.require(2)
		if err != nil {
			return err
		}
		mw.buf[o] = mint8
		mw.buf[o+1] = byte(int8(i))
		return nil

	case a < math.MaxInt16:
		o, err := mw.require(3)
		if err != nil {
			return err
		}
		putMint16(mw.buf[o:], int16(i))
		return nil

	case a < math.MaxInt32:
		o, err := mw.require(5)
		if err != nil {
			return err
		}
		putMint32(mw.buf[o:], int32(i))
		return nil

	default:
		o, err := mw.require(9)
		if err != nil {
			return err
		}
		putMint64(mw.buf[o:], i)
		return nil
	}

}

// WriteInt8 writes an int8 to the writer
func (mw *Writer) WriteInt8(i int8) error { return mw.WriteInt64(int64(i)) }

// WriteInt16 writes an int16 to the writer
func (mw *Writer) WriteInt16(i int16) error { return mw.WriteInt64(int64(i)) }

// WriteInt32 writes an int32 to the writer
func (mw *Writer) WriteInt32(i int32) error { return mw.WriteInt64(int64(i)) }

// WriteInt writes an int to the writer
func (mw *Writer) WriteInt(i int) error { return mw.WriteInt64(int64(i)) }

// WriteUint64 writes a uint64 to the writer
func (mw *Writer) WriteUint64(u uint64) error {
	switch {
	case u < (1 << 7):
		return mw.push(wfixint(uint8(u)))

	case u < math.MaxUint8:
		o, err := mw.require(2)
		if err != nil {
			return err
		}
		mw.buf[o] = muint8
		mw.buf[o+1] = byte(uint8(u))
		return nil
	case u < math.MaxUint16:
		o, err := mw.require(3)
		if err != nil {
			return err
		}
		putMuint16(mw.buf[o:], uint16(u))
		return nil
	case u < math.MaxUint32:
		o, err := mw.require(5)
		if err != nil {
			return err
		}
		putMuint32(mw.buf[o:], uint32(u))
		return nil
	default:
		o, err := mw.require(9)
		if err != nil {
			return err
		}
		putMuint64(mw.buf[o:], u)
		return nil
	}
}

// WriteByte is analagous to WriteUint8
func (mw *Writer) WriteByte(u byte) error { return mw.WriteUint8(uint8(u)) }

// WriteUint8 writes a uint8 to the writer
func (mw *Writer) WriteUint8(u uint8) error { return mw.WriteUint64(uint64(u)) }

// WriteUint16 writes a uint16 to the writer
func (mw *Writer) WriteUint16(u uint16) error { return mw.WriteUint64(uint64(u)) }

// WriteUint32 writes a uint32 to the writer
func (mw *Writer) WriteUint32(u uint32) error { return mw.WriteUint64(uint64(u)) }

// WriteUint writes a uint to the writer
func (mw *Writer) WriteUint(u uint) error { return mw.WriteUint64(uint64(u)) }

// WriteBytes writes binary as 'bin' to the writer
func (mw *Writer) WriteBytes(b []byte) error {
	sz := uint32(len(b))

	// write size
	switch {
	case sz < math.MaxUint8:
		o, err := mw.require(2)
		if err != nil {
			return err
		}
		mw.buf[o] = mbin8
		mw.buf[o+1] = byte(uint8(sz))

	case sz < math.MaxUint16:
		o, err := mw.require(3)
		if err != nil {
			return err
		}
		prefixu16(mw.buf[o:], mbin16, uint16(sz))

	default:
		o, err := mw.require(5)
		if err != nil {
			return err
		}
		prefixu32(mw.buf[o:], mbin32, sz)
	}

	// write body
	_, err := mw.Write(b)
	return err
}

// WriteBool writes a bool to the writer
func (mw *Writer) WriteBool(b bool) error {
	if b {
		return mw.push(mtrue)
	}
	return mw.push(mfalse)
}

// WriteString writes a messagepack string to the writer.
// (This is NOT an implementation of io.StringWriter)
func (mw *Writer) WriteString(s string) error {
	sz := uint32(len(s))

	// write size
	switch {
	case sz < 32:
		err := mw.push(wfixstr(uint8(sz)))
		if err != nil {
			return err
		}
	case sz < math.MaxUint8:
		o, err := mw.require(2)
		if err != nil {
			return err
		}
		mw.buf[o] = mstr8
		mw.buf[o+1] = byte(uint8(sz))
	case sz < math.MaxUint16:
		o, err := mw.require(3)
		if err != nil {
			return err
		}
		prefixu16(mw.buf[o:], mstr16, uint16(sz))
	default:
		o, err := mw.require(5)
		if err != nil {
			return err
		}
		prefixu32(mw.buf[o:], mstr32, sz)
	}

	// write body
	return mw.writeString(s)
}

// WriteComplex64 writes a complex64 to the writer
func (mw *Writer) WriteComplex64(f complex64) error {
	o, err := mw.require(10)
	if err != nil {
		return err
	}
	mw.buf[o] = mfixext8
	mw.buf[o+1] = Complex64Extension
	memcpy8(unsafe.Pointer(&mw.buf[o+2]), unsafe.Pointer(&f))
	return nil
}

// WriteComplex128 writes a complex128 to the writer
func (mw *Writer) WriteComplex128(f complex128) error {
	o, err := mw.require(18)
	if err != nil {
		return err
	}
	mw.buf[o] = mfixext16
	mw.buf[o+1] = Complex128Extension
	memcpy16(unsafe.Pointer(&mw.buf[o+2]), unsafe.Pointer(&f))
	return nil
}

// WriteMapStrStr writes a map[string]string to the writer
func (mw *Writer) WriteMapStrStr(mp map[string]string) (err error) {
	err = mw.WriteMapHeader(uint32(len(mp)))
	if err != nil {
		return
	}
	for key, val := range mp {
		err = mw.WriteString(key)
		if err != nil {
			return
		}
		err = mw.WriteString(val)
		if err != nil {
			return
		}
	}
	return nil
}

// WriteMapStrIntf writes a map[string]interface to the writer
func (mw *Writer) WriteMapStrIntf(mp map[string]interface{}) (err error) {
	err = mw.WriteMapHeader(uint32(len(mp)))
	if err != nil {
		return
	}
	for key, val := range mp {
		err = mw.WriteString(key)
		if err != nil {
			return
		}
		err = mw.WriteIntf(val)
		if err != nil {
			return
		}
	}
	return nil
}

// WriteIdent is a shim for e.EncodeMsg
func (mw *Writer) WriteIdent(e Encodable) error {
	return e.EncodeMsg(mw)
}

// WriteTime writes a time.Time object to the wire
func (mw *Writer) WriteTime(t time.Time) error {
	var bts []byte
	var err error
	bts, err = t.MarshalBinary()
	if err != nil {
		return err
	}
	o, err := mw.require(18)
	if err != nil {
		return err
	}
	mw.buf[o] = mfixext16                   // byte 0 is mfixext16
	mw.buf[o+1] = byte(int8(TimeExtension)) // TimeExtension byte
	copy(mw.buf[o+2:], bts)                 // time is 15 bytes
	mw.buf[o+17] = 0                        // last byte is 0
	return nil
}

// WriteIntf writes the concrete type of 'v'.
// WriteIntf will error if 'v' is not one of the following:
//  - A bool, float, string, []byte, int, uint, or complex
//  - A map of supported types (with string keys)
//  - An array or slice of supported types
//  - A pointer to a supported type
//  - A type that satisfies the msgp.Encodable interface
//  - A type that satisfies the msgp.Extension interface
func (mw *Writer) WriteIntf(v interface{}) error {
	if enc, ok := v.(Encodable); ok {
		return enc.EncodeMsg(mw)
	}
	if mar, ok := v.(Marshaler); ok {
		return mw.Encode(mar)
	}
	if ext, ok := v.(Extension); ok {
		return mw.WriteExtension(ext)
	}
	if v == nil {
		return mw.WriteNil()
	}
	switch v.(type) {
	case bool:
		return mw.WriteBool(v.(bool))
	case float32:
		return mw.WriteFloat32(v.(float32))
	case float64:
		return mw.WriteFloat64(v.(float64))
	case complex64:
		return mw.WriteComplex64(v.(complex64))
	case complex128:
		return mw.WriteComplex128(v.(complex128))
	case uint8:
		return mw.WriteUint8(v.(uint8))
	case uint16:
		return mw.WriteUint16(v.(uint16))
	case uint32:
		return mw.WriteUint32(v.(uint32))
	case uint64:
		return mw.WriteUint64(v.(uint64))
	case uint:
		return mw.WriteUint(v.(uint))
	case int8:
		return mw.WriteInt8(v.(int8))
	case int16:
		return mw.WriteInt16(v.(int16))
	case int32:
		return mw.WriteInt32(v.(int32))
	case int64:
		return mw.WriteInt64(v.(int64))
	case int:
		return mw.WriteInt(v.(int))
	case string:
		return mw.WriteString(v.(string))
	case []byte:
		return mw.WriteBytes(v.([]byte))
	case map[string]string:
		return mw.WriteMapStrStr(v.(map[string]string))
	case map[string]interface{}:
		return mw.WriteMapStrIntf(v.(map[string]interface{}))
	case time.Time:
		return mw.WriteTime(v.(time.Time))
	}

	val := reflect.ValueOf(v)
	if !isSupported(val.Kind()) || !val.IsValid() {
		return fmt.Errorf("msgp: type %s not supported", val)
	}

	switch val.Kind() {
	case reflect.Ptr:
		if val.IsNil() {
			return mw.WriteNil()
		}
		return mw.WriteIntf(val.Elem().Interface())
	case reflect.Slice:
		return mw.writeSlice(val)
	case reflect.Map:
		return mw.writeMap(val)
	}
	return fmt.Errorf("msgp: type %s not supported", val.Type())
}

func (mw *Writer) writeMap(v reflect.Value) (err error) {
	if v.Elem().Kind() != reflect.String {
		return errors.New("msgp: map keys must be strings")
	}
	ks := v.MapKeys()
	err = mw.WriteMapHeader(uint32(len(ks)))
	if err != nil {
		return
	}
	for _, key := range ks {
		val := v.MapIndex(key)
		err = mw.WriteString(key.String())
		if err != nil {
			return
		}
		err = mw.WriteIntf(val.Interface())
		if err != nil {
			return
		}
	}
	return
}

func (mw *Writer) writeSlice(v reflect.Value) (err error) {
	// is []byte
	if v.Type().ConvertibleTo(btsType) {
		return mw.WriteBytes(v.Bytes())
	}

	sz := uint32(v.Len())
	err = mw.WriteArrayHeader(sz)
	if err != nil {
		return
	}
	for i := uint32(0); i < sz; i++ {
		err = mw.WriteIntf(v.Index(int(i)).Interface())
		if err != nil {
			return
		}
	}
	return
}

func (mw *Writer) writeStruct(v reflect.Value) error {
	if enc, ok := v.Interface().(Encodable); ok {
		return enc.EncodeMsg(mw)
	}
	if mar, ok := v.Interface().(Marshaler); ok {
		return mw.Encode(mar)
	}
	return fmt.Errorf("msgp: unsupported type: %s", v.Type())
}

func (mw *Writer) writeVal(v reflect.Value) error {
	if !isSupported(v.Kind()) {
		return fmt.Errorf("msgp: msgp/enc: type %q not supported", v.Type())
	}

	// shortcut for nil values
	if v.IsNil() {
		return mw.WriteNil()
	}
	switch v.Kind() {
	case reflect.Bool:
		return mw.WriteBool(v.Bool())

	case reflect.Float32, reflect.Float64:
		return mw.WriteFloat64(v.Float())

	case reflect.Complex64, reflect.Complex128:
		return mw.WriteComplex128(v.Complex())

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		return mw.WriteInt64(v.Int())

	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			mw.WriteNil()
		}
		return mw.writeVal(v.Elem())

	case reflect.Map:
		return mw.writeMap(v)

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		return mw.WriteUint64(v.Uint())

	case reflect.String:
		return mw.WriteString(v.String())

	case reflect.Slice, reflect.Array:
		return mw.writeSlice(v)

	case reflect.Struct:
		return mw.writeStruct(v)

	}
	return fmt.Errorf("msgp: msgp/enc: type %q not supported", v.Type())
}

// is the reflect.Kind encodable?
func isSupported(k reflect.Kind) bool {
	switch k {
	case reflect.Func, reflect.Chan, reflect.Invalid, reflect.UnsafePointer:
		return false
	default:
		return true
	}
}

// GuessSize guesses the size of the underlying
// value of 'i'. If the underlying value is not
// a simple builtin (or []byte), GuessSize defaults
// to 512.
func GuessSize(i interface{}) int {
	if s, ok := i.(Sizer); ok {
		return s.Msgsize()
	} else if e, ok := i.(Extension); ok {
		return ExtensionPrefixSize + e.Len()
	} else if i == nil {
		return NilSize
	}

	switch i.(type) {
	case float64:
		return Float64Size
	case float32:
		return Float32Size
	case uint8, uint16, uint32, uint64, uint:
		return UintSize
	case int8, int16, int32, int64, int:
		return IntSize
	case []byte:
		return BytesPrefixSize + len(i.([]byte))
	case string:
		return StringPrefixSize + len(i.(string))
	case complex64:
		return Complex64Size
	case complex128:
		return Complex128Size
	case bool:
		return BoolSize
	case map[string]interface{}:
		s := MapHeaderSize
		for key, val := range i.(map[string]interface{}) {
			s += StringPrefixSize + len(key)
			s += GuessSize(val)
		}
		return s
	case map[string]string:
		s := MapHeaderSize
		for key, val := range i.(map[string]string) {
			s += 2*StringPrefixSize + len(key) + len(val)
		}
		return s
	default:
		return 512
	}
}
