package enc

import (
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"
)

// ensures 'sz' extra bytes; returns offset and 'ok'
func grow(b []byte, sz int) ([]byte, int, bool) {
	if cap(b)-len(b) >= sz {
		return b[:len(b)+sz], len(b), true
	}
	return nil, 0, false
}

func realloc(b []byte, sz int) ([]byte, int) {
	nsz := (2 * cap(b)) + sz
	g := make([]byte, nsz)
	n := copy(g, b)
	return g[:n+sz], n
}

func ensure(b []byte, sz int) ([]byte, int) {
	o, n, ok := grow(b, sz)
	if !ok {
		return realloc(b, sz)
	}
	return o, n
}

func AppendMapHeader(b []byte, sz uint32) []byte {
	o, n := ensure(b, MapHeaderSize)
	switch {
	case sz < 16:
		o[n] = wfixmap(uint8(sz))
		return o[:n+1]

	case sz < math.MaxUint16:
		u := uint16(sz)
		o[n] = mmap16
		binary.BigEndian.PutUint16(o[n+1:], u)
		return o[:n+3]

	default:
		o[n] = mmap32
		binary.BigEndian.PutUint32(o[n+1:], sz)
		return o[:n+5]
	}
}

func AppendArrayHeader(b []byte, sz uint32) []byte {
	o, n := ensure(b, ArrayHeaderSize)
	switch {
	case sz < 16:
		o[n] = wfixarray(uint8(sz))
		return o[:n+1]

	case sz < math.MaxUint16:
		u := uint16(sz)
		o[n] = marray16
		binary.BigEndian.PutUint16(o[n+1:], u)
		return o[:n+3]

	default:
		o[n] = marray32
		binary.BigEndian.PutUint32(o[n+1:], sz)
		return o[:n+5]
	}
}

func AppendNil(b []byte) []byte { return append(b, mnil) }

func AppendFloat64(b []byte, f float64) []byte {
	o, n := ensure(b, Float64Size)
	o[n] = mfloat64
	binary.BigEndian.PutUint64(o[n+1:], *(*uint64)(unsafe.Pointer(&f)))
	return o[:n+9]
}

func AppendFloat32(b []byte, f float32) []byte {
	o, n := ensure(b, Float32Size)
	o[n] = mfloat32
	binary.BigEndian.PutUint32(o[n+1:], *(*uint32)(unsafe.Pointer(&f)))
	return o[:n+5]
}

func AppendInt64(b []byte, i int64) []byte {
	o, n := ensure(b, Int64Size)
	switch {
	case i < 0 && i > -32:
		o[n] = byte(int8(i))
		return o[:n+1]

	case i >= 0 && i < 128:
		o[n] = (byte(i) & 0x7f)
		return o[:n+1]

	case abs(i) < math.MaxInt8:
		o[n] = mint8
		o[n+1] = byte(int8(i))
		return o[:n+2]

	case abs(i) < math.MaxInt16:
		o[n] = mint16
		o[n+1] = byte(int16(i >> 8))
		o[n+2] = byte(int16(i))
		return o[:n+3]

	case abs(i) < math.MaxInt32:
		o[n] = mint32
		o[n+1] = byte(int32(i >> 24))
		o[n+2] = byte(int32(i >> 16))
		o[n+3] = byte(int32(i >> 8))
		o[n+4] = byte(int32(i))
		return o[:n+5]

	default:
		o[n] = mint64
		o[n+1] = byte(i >> 56)
		o[n+2] = byte(i >> 48)
		o[n+3] = byte(i >> 40)
		o[n+4] = byte(i >> 32)
		o[n+5] = byte(i >> 24)
		o[n+6] = byte(i >> 16)
		o[n+7] = byte(i >> 8)
		o[n+8] = byte(i)
		return o[:n+9]
	}
}

// TODO: expand AppendIntXx and AppendUintXx into the proper code; we can remove some unnecessary branches

func AppendInt(b []byte, i int) []byte     { return AppendInt64(b, int64(i)) }
func AppendInt8(b []byte, i int8) []byte   { return AppendInt64(b, int64(i)) }
func AppendInt16(b []byte, i int16) []byte { return AppendInt64(b, int64(i)) }
func AppendInt32(b []byte, i int32) []byte { return AppendInt64(b, int64(i)) }

func AppendUint64(b []byte, u uint64) []byte {
	o, n := ensure(b, Uint64Size)
	switch {
	case u < (1 << 7):
		o[n] = byte(u) & 0x7f
		return o[:n+1]

	case u < math.MaxUint8:
		o[n] = muint8
		o[n+1] = byte(uint8(u))
		return o[:n+2]

	case u < math.MaxUint16:
		o[n] = muint16
		binary.BigEndian.PutUint16(o[n+1:], uint16(u))
		return o[:n+3]

	case u < math.MaxUint32:
		o[n] = muint32
		binary.BigEndian.PutUint32(o[n+1:], uint32(u))
		return o[:n+5]

	default:
		o[n] = muint64
		binary.BigEndian.PutUint64(o[n+1:], u)
		return o[:n+9]

	}
}

func AppendUint(b []byte, u uint) []byte     { return AppendUint64(b, uint64(u)) }
func AppendUint8(b []byte, u uint8) []byte   { return AppendUint64(b, uint64(u)) }
func AppendUint16(b []byte, u uint16) []byte { return AppendUint64(b, uint64(u)) }
func AppendUint32(b []byte, u uint32) []byte { return AppendUint64(b, uint64(u)) }

func AppendBytes(b []byte, bts []byte) []byte {
	sz := uint32(len(bts))
	o, n := ensure(b, BytesPrefixSize+len(bts))
	switch {
	case sz < math.MaxUint8:
		o[n] = mbin8
		o[n+1] = byte(uint8(sz))
		n += 2

	case sz < math.MaxUint16:
		o[n] = mbin16
		binary.BigEndian.PutUint16(o[n+1:], uint16(sz))
		n += 3

	default:
		o[n] = mbin32
		binary.BigEndian.PutUint32(o[n+1:], sz)
		n += 5
	}
	j := copy(o[n:], bts)
	return o[:n+j]
}

func AppendBool(b []byte, t bool) []byte {
	if t {
		return append(b, mtrue)
	}
	return append(b, mfalse)
}

func AppendString(b []byte, s string) []byte {
	sz := uint32(len(s))
	o, n := ensure(b, BytesPrefixSize+len(s))
	switch {
	case sz < 32:
		o[n] = wfixstr(uint8(sz))
		n++

	case sz < math.MaxUint8:
		o[n] = mstr8
		o[n+1] = byte(uint8(sz))
		n += 2

	case sz < math.MaxUint16:
		o[n] = mstr16
		binary.BigEndian.PutUint16(o[n+1:], uint16(sz))
		n += 3

	default:
		o[n] = mstr32
		binary.BigEndian.PutUint32(o[n+1:], sz)
		n += 5
	}
	j := copy(o[n:], s)
	return o[:n+j]
}

func AppendComplex64(b []byte, c complex64) []byte {
	o, n := ensure(b, Complex64Size)
	o[n] = mfixext8
	o[n+1] = Complex64Extension
	rl := real(c)
	im := imag(c)
	binary.BigEndian.PutUint32(o[n+2:], *(*uint32)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint32(o[n+6:], *(*uint32)(unsafe.Pointer(&im)))
	return o[:n+10]
}

func AppendComplex128(b []byte, c complex128) []byte {
	o, n := ensure(b, Complex64Size)
	o[n] = mfixext16
	o[n+1] = Complex128Extension
	rl := real(c)
	im := imag(c)
	binary.BigEndian.PutUint64(o[n+2:], *(*uint64)(unsafe.Pointer(&rl)))
	binary.BigEndian.PutUint64(o[n+10:], *(*uint64)(unsafe.Pointer(&im)))
	return o[:n+18]
}

func AppendTime(b []byte, t time.Time) []byte {
	o, n := ensure(b, TimeSize)
	bts, _ := t.MarshalBinary()
	o[n] = mfixext16
	o[n+1] = TimeExtension
	copy(o[n+2:], bts)
	return o[:n+18]
}

func AppendMapStrStr(b []byte, m map[string]string) []byte {
	sz := uint32(len(m))
	b = AppendMapHeader(b, sz)
	for key, val := range m {
		b = AppendString(b, key)
		b = AppendString(b, val)
	}
	return b
}

func AppendMapStrIntf(b []byte, m map[string]interface{}) ([]byte, error) {
	sz := uint32(len(m))
	b = AppendMapHeader(b, sz)
	var err error
	for key, val := range m {
		b = AppendString(b, key)
		b, err = AppendIntf(b, val)
		if err != nil {
			return b, err
		}
	}
	return b, nil
}

func AppendExtension(b []byte, e Extension) []byte {
	o, n := ensure(b, 5+len(e.Data))
	l := len(e.Data)
	switch l {
	case 0:
		o[n] = mext8
		o[n+1] = e.Type
		o[n+2] = 0
		return o[:n+3]
	case 1:
		o[n] = mfixext1
		o[n+1] = e.Type
		o[n+2] = e.Data[0]
		return o[:n+3]
	case 2:
		o[n] = mfixext2
		o[n+1] = e.Type
		copy(o[n+2:], e.Data)
		return o[:n+4]
	case 4:
		o[n] = mfixext4
		o[n+1] = e.Type
		copy(o[n+2:], e.Data)
		return o[:n+6]
	case 8:
		o[n] = mfixext8
		o[n+1] = e.Type
		copy(o[n+2:], e.Data)
		return o[:n+10]
	case 16:
		o[n] = mfixext16
		o[n+1] = e.Type
		copy(o[n+2:], e.Data)
		return o[:n+18]
	}
	switch {
	case l < math.MaxUint8:
		o[n] = mext8
		o[n+1] = e.Type
		o[n+2] = byte(uint8(l))
		n += 3
	case l < math.MaxUint16:
		o[n] = mext16
		o[n+1] = e.Type
		binary.BigEndian.PutUint16(o[n+2:], uint16(l))
		n += 4
	default:
		o[n] = mext32
		o[n+1] = e.Type
		binary.BigEndian.PutUint32(o[n+2:], uint32(l))
		n += 6
	}
	x := copy(o[n:], e.Data)
	return o[:n+x]
}

func AppendIntf(b []byte, i interface{}) ([]byte, error) {
	if m, ok := i.(MsgMarshaler); ok {
		return m.AppendMsg(b)
	}

	if i == nil {
		return AppendNil(b), nil
	}

	// all the concrete types
	// for which we have methods
	switch i.(type) {
	case bool:
		return AppendBool(b, i.(bool)), nil
	case float32:
		return AppendFloat32(b, i.(float32)), nil
	case float64:
		return AppendFloat64(b, i.(float64)), nil
	case complex64:
		return AppendComplex64(b, i.(complex64)), nil
	case complex128:
		return AppendComplex128(b, i.(complex128)), nil
	case string:
		return AppendString(b, i.(string)), nil
	case []byte:
		return AppendBytes(b, i.([]byte)), nil
	case int8:
		return AppendInt8(b, i.(int8)), nil
	case int16:
		return AppendInt16(b, i.(int16)), nil
	case int32:
		return AppendInt32(b, i.(int32)), nil
	case int64:
		return AppendInt64(b, i.(int64)), nil
	case int:
		return AppendInt64(b, int64(i.(int))), nil
	case uint:
		return AppendUint64(b, uint64(i.(uint))), nil
	case uint8:
		return AppendUint8(b, i.(uint8)), nil
	case uint16:
		return AppendUint16(b, i.(uint16)), nil
	case uint32:
		return AppendUint32(b, i.(uint32)), nil
	case uint64:
		return AppendUint64(b, i.(uint64)), nil
	case map[string]interface{}:
		return AppendMapStrIntf(b, i.(map[string]interface{}))
	case map[string]string:
		return AppendMapStrStr(b, i.(map[string]string)), nil
	case []interface{}:
		j := i.([]interface{})
		b = AppendArrayHeader(b, uint32(len(j)))
		var err error
		for _, k := range j {
			b, err = AppendIntf(b, k)
			if err != nil {
				return b, err
			}
		}
		return b, nil
	}

	var err error
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		l := v.Len()
		b = AppendArrayHeader(b, uint32(l))
		for i := 0; i < l; i++ {
			b, err = AppendIntf(b, v.Index(i).Interface())
			if err != nil {
				return b, err
			}
		}
		return b, nil

	// TODO: maybe some struct fiddling?

	default:
		return b, fmt.Errorf("msgp/enc: type %q not supported", v.Type())
	}
}
