package msgp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
	"unsafe"
)

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

type sizer interface {
	Maxsize() int
}

var (
	btsType reflect.Type
)

func init() {
	var bts []byte
	btsType = reflect.TypeOf(bts)
}

// Marshaler is the interface implemented
// by types that know how to marshal themselves
// as MessagePack
type Marshaler interface {
	// AppendMsg appends the marshalled
	// form of the object to the provided
	// byte slice, returning the extended
	// slice and any errors encountered
	AppendMsg([]byte) ([]byte, error)

	// MarshalMsg returns a new []byte
	// containing the MessagePack-encoded
	// form of the object
	MarshalMsg() ([]byte, error)
}

// Encoder is the interface implemented
// by types that know how to write themselves
// as MessagePack
type Encoder interface {
	// EncodeMsg writes the object in MessagePack
	// format to the writer, returning the number
	// of bytes written and any errors encountered.
	EncodeMsg(io.Writer) (int, error)

	// EncodeTo writes the object to a *Writer,
	// returning the number of bytes written, and
	// any errors encountered.
	EncodeTo(*Writer) (int, error)
}

// Writer is the object used by
// Encoders to write themselves
// to an io.Writer
type Writer struct {
	w       io.Writer
	scratch [24]byte
}

// NewWriter returns a Writer, which can be
// used by Encoder to serialize themselves
// to the provided io.Writer. It does no buffering.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		w: w,
	}
}

// Write implements the standard io.Writer interface
func (mw *Writer) Write(p []byte) (int, error) {
	return mw.w.Write(p)
}

// Reset changes the underlying writer used by the MsgWriter
func (mw *Writer) Reset(w io.Writer) {
	mw.w = w
}

func (mw *Writer) WriteMapHeader(sz uint32) (n int, err error) {
	switch {
	case sz < 16:
		mw.scratch[0] = wfixmap(uint8(sz))
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

func (mw *Writer) WriteArrayHeader(sz uint32) (n int, err error) {
	switch {
	case sz < 16:
		mw.scratch[0] = wfixarray(uint8(sz))
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

func (mw *Writer) WriteNil() (n int, err error) {
	mw.scratch[0] = mnil
	return mw.w.Write(mw.scratch[:1])
}

func (mw *Writer) WriteFloat64(f float64) (n int, err error) {
	mw.scratch[0] = mfloat64
	copy(mw.scratch[1:], (*(*[8]byte)(unsafe.Pointer(&f)))[:])
	return mw.w.Write(mw.scratch[:9])
}

func (mw *Writer) WriteFloat32(f float32) (n int, err error) {
	mw.scratch[0] = mfloat32
	copy(mw.scratch[1:], (*(*[4]byte)(unsafe.Pointer(&f)))[:])
	return mw.w.Write(mw.scratch[:5])
}

func (mw *Writer) WriteInt64(i int64) (n int, err error) {
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

func (m *Writer) WriteInt8(i int8) (int, error)   { return m.WriteInt64(int64(i)) }
func (m *Writer) WriteInt16(i int16) (int, error) { return m.WriteInt64(int64(i)) }
func (m *Writer) WriteInt32(i int32) (int, error) { return m.WriteInt64(int64(i)) }
func (m *Writer) WriteInt(i int) (int, error)     { return m.WriteInt64(int64(i)) }

func (mw *Writer) WriteUint64(u uint64) (n int, err error) {
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

func (m *Writer) WriteByte(u byte) (int, error)     { return m.WriteUint8(uint8(u)) }
func (m *Writer) WriteUint8(u uint8) (int, error)   { return m.WriteUint64(uint64(u)) }
func (m *Writer) WriteUint16(u uint16) (int, error) { return m.WriteUint64(uint64(u)) }
func (m *Writer) WriteUint32(u uint32) (int, error) { return m.WriteUint64(uint64(u)) }
func (m *Writer) WriteUint(u uint) (int, error)     { return m.WriteUint64(uint64(u)) }

func (mw *Writer) WriteBytes(b []byte) (n int, err error) {
	l := len(b)
	if l > math.MaxUint32 {
		panic("msgp: cannot write bytes with length > MaxUint32")
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

func (mw *Writer) WriteBool(b bool) (n int, err error) {
	if b {
		mw.scratch[0] = mtrue
	} else {
		mw.scratch[0] = mfalse
	}
	return mw.w.Write(mw.scratch[:1])
}

func (mw *Writer) WriteString(s string) (n int, err error) {
	var nn int
	l := len(s)
	if l > math.MaxUint32 {
		panic("msgp: cannot write string with length greater than MaxUint32")
	}
	sz := uint32(l)
	switch {
	case sz < 32:
		mw.scratch[0] = wfixstr(uint8(sz))
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
	nn, err = mw.w.Write(UnsafeBytes(s))
	n += nn
	return
}

func (mw *Writer) WriteComplex64(f complex64) (n int, err error) {
	mw.scratch[0] = mfixext8
	mw.scratch[1] = Complex64Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint32(mw.scratch[2:], *(*uint32)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint32(mw.scratch[6:], *(*uint32)(unsafe.Pointer(&im)))
	return mw.w.Write(mw.scratch[:10])
}

func (mw *Writer) WriteComplex128(f complex128) (n int, err error) {
	mw.scratch[0] = mfixext16
	mw.scratch[1] = Complex128Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint64(mw.scratch[2:], *(*uint64)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint64(mw.scratch[10:], *(*uint64)(unsafe.Pointer(&im)))
	return mw.w.Write(mw.scratch[:18])
}

func (mw *Writer) WriteMapStrStr(mp map[string]string) (n int, err error) {
	var nn int
	nn, err = mw.WriteMapHeader(uint32(len(mp)))
	n += nn
	if err != nil {
		return
	}
	for key, val := range mp {
		nn, err = mw.WriteString(key)
		n += nn
		if err != nil {
			return
		}
		nn, err = mw.WriteString(val)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (mw *Writer) WriteMapStrIntf(mp map[string]interface{}) (n int, err error) {
	var nn int
	nn, err = mw.WriteMapHeader(uint32(len(mp)))
	n += nn
	if err != nil {
		return
	}
	for key, val := range mp {
		nn, err = mw.WriteString(key)
		n += nn
		if err != nil {
			return
		}
		nn, err = mw.WriteIntf(val)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (mw *Writer) WriteIdent(e Encoder) (n int, err error) {
	return e.EncodeTo(mw)
}

func (mw *Writer) WriteTime(t time.Time) (n int, err error) {
	var bts []byte
	bts, err = t.MarshalBinary()
	if err != nil {
		return
	}
	mw.scratch[0] = mfixext16                 // byte 0 is mfixext16
	mw.scratch[1] = byte(int8(TimeExtension)) // TimeExtension byte
	copy(mw.scratch[2:], bts)                 // time is 15 bytes
	mw.scratch[17] = 0                        // last byte is 0
	return mw.w.Write(mw.scratch[:18])
}

// WriteIntf writes the concrete type of 'v'.
// WriteIntf will error if 'v' is not one of the following:
//  - A bool, float, string, []byte, int, uint, or complex
//  - A map of supported types (with string keys)
//  - An array or slice of supported types
//  - A pointer to a supported type
//  - A type that satisfies the msgp.Encoder interface
//  - A type that satisfies the msgp.Extension interface
func (mw *Writer) WriteIntf(v interface{}) (n int, err error) {
	if enc, ok := v.(Encoder); ok {
		return enc.EncodeTo(mw)
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
		return 0, fmt.Errorf("msgp: type %s not supported", val)
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
	return 0, fmt.Errorf("msgp: type %s not supported", val.Type())
}

func (mw *Writer) writeMap(v reflect.Value) (n int, err error) {
	if v.Elem().Kind() != reflect.String {
		return 0, errors.New("msgp: map keys must be strings")
	}
	ks := v.MapKeys()
	var nn int
	nn, err = mw.WriteMapHeader(uint32(len(ks)))
	n += nn
	if err != nil {
		return
	}
	for _, key := range ks {
		val := v.MapIndex(key)

		nn, err = mw.WriteString(key.String())
		n += nn
		if err != nil {
			return
		}
		nn, err = mw.WriteIntf(val.Interface())
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (mw *Writer) writeSlice(v reflect.Value) (n int, err error) {
	var nn int
	if v.Type().ConvertibleTo(btsType) {
		nn, err = mw.WriteBytes(v.Bytes())
		n += nn
		if err != nil {
			return
		}
	}
	sz := uint32(v.Len())
	nn, err = mw.WriteArrayHeader(sz)
	n += nn
	if err != nil {
		return
	}
	for i := uint32(0); i < sz; i++ {
		nn, err = mw.WriteIntf(v.Index(int(i)).Interface())
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (mw *Writer) writeStruct(v reflect.Value) (n int, err error) {
	if enc, ok := v.Interface().(Encoder); ok {
		return enc.EncodeTo(mw)
	}
	return 0, fmt.Errorf("msgp: unsupported type: %s", v.Type())
}

func (mw *Writer) writeVal(v reflect.Value) (n int, err error) {
	if !isSupported(v.Kind()) {
		return 0, fmt.Errorf("msgp: msgp/enc: type %q not supported", v.Type())
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
	return 0, fmt.Errorf("msgp: msgp/enc: type %q not supported", v.Type())
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
	if s, ok := i.(sizer); ok {
		return s.Maxsize()
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
