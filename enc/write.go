package enc

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"unsafe"
)

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

var (
	mapStrStr  reflect.Type
	mapStrIntf reflect.Type
	btsType    reflect.Type
)

func init() {
	var mss map[string]string
	mapStrStr = reflect.TypeOf(mss)
	var msi map[string]interface{}
	mapStrIntf = reflect.TypeOf(msi)
	var bts []byte
	btsType = reflect.TypeOf(bts)
}

const (
	// Complex64 numbers are encoded as extension #3
	Complex64Extension = 3

	// Complex128 numbers are encoded as extension #4
	Complex128Extension = 4
)

type MsgEncoder interface {
	EncodeMsg(io.Writer) (int, error)
}

type MsgWriter struct {
	w       io.Writer
	scratch [24]byte
}

func NewEncoder(w io.Writer) *MsgWriter {
	return &MsgWriter{
		w: w,
	}
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

func (mw *MsgWriter) WriteNil() (n int, err error) {
	mw.scratch[0] = mnil
	return mw.w.Write(mw.scratch[:1])
}

func (mw *MsgWriter) WriteFloat64(f float64) (n int, err error) {
	mw.scratch[0] = mfloat64
	bits := *(*uint64)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint64(mw.scratch[1:], bits)
	return mw.w.Write(mw.scratch[:9])
}

func (mw *MsgWriter) WriteFloat32(f float32) (n int, err error) {
	mw.scratch[0] = mfloat32
	bits := *(*uint32)(unsafe.Pointer(&f))
	binary.BigEndian.PutUint32(mw.scratch[1:], bits)
	return mw.w.Write(mw.scratch[:5])
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

// Integer Utilities
func (m *MsgWriter) WriteInt8(i int8) (int, error)   { return m.WriteInt64(int64(i)) }
func (m *MsgWriter) WriteInt16(i int16) (int, error) { return m.WriteInt64(int64(i)) }
func (m *MsgWriter) WriteInt32(i int32) (int, error) { return m.WriteInt64(int64(i)) }
func (m *MsgWriter) WriteInt(i int) (int, error)     { return m.WriteInt64(int64(i)) }

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

// Unsigned Integer Utilities
func (m *MsgWriter) WriteByte(u byte) (int, error)     { return m.WriteUint8(uint8(u)) }
func (m *MsgWriter) WriteUint8(u uint8) (int, error)   { return m.WriteUint64(uint64(u)) }
func (m *MsgWriter) WriteUint16(u uint16) (int, error) { return m.WriteUint64(uint64(u)) }
func (m *MsgWriter) WriteUint32(u uint32) (int, error) { return m.WriteUint64(uint64(u)) }
func (m *MsgWriter) WriteUint(u uint) (int, error)     { return m.WriteUint64(uint64(u)) }

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

func (mw *MsgWriter) WriteBool(b bool) (n int, err error) {
	if b {
		mw.scratch[0] = mtrue
	} else {
		mw.scratch[1] = mfalse
	}
	return mw.w.Write(mw.scratch[:1])
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

func (mw *MsgWriter) WriteComplex64(f complex64) (n int, err error) {
	mw.scratch[0] = mfixext8
	mw.scratch[1] = Complex64Extension
	rl := real(f)
	im := imag(f)
	binary.BigEndian.PutUint32(mw.scratch[2:], *(*uint32)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint32(mw.scratch[6:], *(*uint32)(unsafe.Pointer(&im)))
	return mw.w.Write(mw.scratch[:10])
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

// WriteExtension writes an extension type (tuple of (type, binary))
//func WriteExtension(w io.Writer, typ byte, bts []byte) (n int, err error) {/
//
//}

func (mw *MsgWriter) WriteMapStrStr(mp map[string]string) (n int, err error) {
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

func (mw *MsgWriter) WriteMapStrIntf(mp map[string]interface{}) (n int, err error) {
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
		nn, err = mw.Encode(val)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (mw *MsgWriter) Encode(v interface{}) (n int, err error) {
	if enc, ok := v.(MsgEncoder); ok {
		return enc.EncodeMsg(mw.w)
	}
	return mw.writeVal(reflect.ValueOf(v))
}

func (mw *MsgWriter) writeMap(v reflect.Value) (n int, err error) {
	tp := v.Type()
	if tp.AssignableTo(mapStrStr) || tp.ConvertibleTo(mapStrStr) {
		return mw.WriteMapStrStr(v.Interface().(map[string]string))
	} else if tp.AssignableTo(mapStrIntf) || tp.ConvertibleTo(mapStrIntf) {
		var nn int
		sz := uint32(v.Len())
		nn, err = mw.WriteMapHeader(sz)
		n += nn
		if err != nil {
			return
		}
		keys := v.MapKeys()
		for _, key := range keys {
			nn, err = mw.WriteString(key.String())
			n += nn
			if err != nil {
				return
			}
			nn, err = mw.writeVal(v.MapIndex(key))
			n += nn
			if err != nil {
				return
			}
		}
		return
	} else {
		return 0, fmt.Errorf("type %s not supported", tp)
	}
}

func (mw *MsgWriter) writeSlice(v reflect.Value) (n int, err error) {
	var nn int
	if v.Type().AssignableTo(btsType) {
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
		nn, err = mw.writeVal(v.Index(int(i)))
		n += nn
		if err != nil {
			return
		}
	}
	return
}

func (mw *MsgWriter) writeStruct(v reflect.Value) (n int, err error) {
	if enc, ok := v.Interface().(MsgEncoder); ok {
		return enc.EncodeMsg(mw.w)
	}

	sz := uint32(v.NumField())
	var nn int
	nn, err = mw.WriteArrayHeader(sz)
	n += nn
	if err != nil {
		return
	}
	for i := uint32(0); i < sz; i++ {
		k := v.Field(int(i)).Interface().(*reflect.StructField)
		var name string
		msg := k.Tag.Get("msg")
		if msg != "" {
			name = msg
		} else {
			name = k.Name
		}
		nn, err = mw.WriteString(name)
		n += nn
		if err != nil {
			return
		}
		// TODO
	}
	return
}

func (mw *MsgWriter) writeVal(v reflect.Value) (n int, err error) {
	// shortcut for nil values
	if v.IsNil() {
		return mw.WriteNil()
	}
	var nn int
	switch v.Kind() {
	case reflect.Bool:
		nn, err = mw.WriteBool(v.Bool())

	case reflect.Float32, reflect.Float64:
		nn, err = mw.WriteFloat64(v.Float())

	case reflect.Complex64, reflect.Complex128:
		nn, err = mw.WriteComplex128(v.Complex())

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8:
		nn, err = mw.WriteInt64(v.Int())

	case reflect.Interface, reflect.Ptr:
		nn, err = mw.writeVal(v.Elem())

	case reflect.Map:
		nn, err = mw.writeMap(v)

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8:
		nn, err = mw.WriteUint64(v.Uint())

	case reflect.String:
		nn, err = mw.WriteString(v.String())

	case reflect.Slice, reflect.Array:
		nn, err = mw.writeSlice(v)

	case reflect.Struct:
		nn, err = mw.writeStruct(v)

	}
	n += nn
	return
}
