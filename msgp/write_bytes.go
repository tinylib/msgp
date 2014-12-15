package msgp

import (
	"fmt"
	"math"
	"reflect"
	"time"
	"unsafe"
)

// ErrUnsupportedType is returned
// when a bad argument is supplied
// to a function that takes `interface{}`.
type ErrUnsupportedType struct {
	T reflect.Type
}

// Error implements error
func (e *ErrUnsupportedType) Error() string { return fmt.Sprintf("msgp: type %q not supported", e.T) }

// Resumable returns 'true' for ErrUnsupportedType
func (e *ErrUnsupportedType) Resumable() bool { return true }

// ensure 'sz' extra bytes in 'b' btw len(b) and cap(b)
func ensure(b []byte, sz int) ([]byte, int) {
	l := len(b)
	c := cap(b)
	if c-l < sz {
		o := make([]byte, (2*c)+sz) // exponential growth
		n := copy(o, b)
		return o[:n+sz], n
	}
	return b[:l+sz], l
}

// AppendMapHeader appends a map header with the
// given size to the slice
func AppendMapHeader(b []byte, sz uint32) []byte {
	switch {
	case sz < 16:
		return append(b, wfixmap(uint8(sz)))

	case sz < math.MaxUint16:
		o, n := ensure(b, 3)
		prefixu16(o[n:], mmap16, uint16(sz))
		return o

	default:
		o, n := ensure(b, 5)
		prefixu32(o[n:], mmap32, sz)
		return o
	}
}

// AppendArrayHeader appends an array header with
// the given size to the slice
func AppendArrayHeader(b []byte, sz uint32) []byte {
	switch {
	case sz < 16:
		return append(b, wfixarray(uint8(sz)))

	case sz < math.MaxUint16:
		o, n := ensure(b, 3)
		prefixu16(o[n:], marray16, uint16(sz))
		return o

	default:
		o, n := ensure(b, 5)
		prefixu32(o[n:], marray32, sz)
		return o
	}
}

// AppendNil appends a 'nil' byte to the slice
func AppendNil(b []byte) []byte { return append(b, mnil) }

// AppendFloat64 appends a float64 to the slice
func AppendFloat64(b []byte, f float64) []byte {
	o, n := ensure(b, Float64Size)
	o[n] = mfloat64
	memcpy8(unsafe.Pointer(&o[n+1]), unsafe.Pointer(&f))
	return o
}

// AppendFloat32 appends a float32 to the slice
func AppendFloat32(b []byte, f float32) []byte {
	o, n := ensure(b, Float32Size)
	o[n] = mfloat32
	memcpy4(unsafe.Pointer(&o[n+1]), unsafe.Pointer(&f))
	return o
}

// AppendInt64 appends an int64 to the slice
func AppendInt64(b []byte, i int64) []byte {
	a := abs(i)
	switch {
	case i < 0 && i > -32:
		return append(b, wnfixint(int8(i)))

	case i >= 0 && i < 128:
		return append(b, wfixint(uint8(i)))

	case a < math.MaxInt8:
		o, n := ensure(b, 2)
		putMint8(o[n:], int8(i))
		return o

	case a < math.MaxInt16:
		o, n := ensure(b, 3)
		putMint16(o[n:], int16(i))
		return o

	case a < math.MaxInt32:
		o, n := ensure(b, 5)
		putMint32(o[n:], int32(i))
		return o

	default:
		o, n := ensure(b, 9)
		putMint64(o[n:], i)
		return o
	}
}

// AppendInt appends an int to the slice
func AppendInt(b []byte, i int) []byte { return AppendInt64(b, int64(i)) }

// AppendInt8 appends an int8 to the slice
func AppendInt8(b []byte, i int8) []byte { return AppendInt64(b, int64(i)) }

// AppendInt16 appends an int16 to the slice
func AppendInt16(b []byte, i int16) []byte { return AppendInt64(b, int64(i)) }

// AppendInt32 appends an int32 to the slice
func AppendInt32(b []byte, i int32) []byte { return AppendInt64(b, int64(i)) }

// AppendUint64 appends a uint64 to the slice
func AppendUint64(b []byte, u uint64) []byte {
	switch {
	case u < (1 << 7):
		return append(b, wfixint(uint8(u)))

	case u < math.MaxUint8:
		o, n := ensure(b, 2)
		putMuint8(o[n:], uint8(u))
		return o

	case u < math.MaxUint16:
		o, n := ensure(b, 3)
		putMuint16(o[n:], uint16(u))
		return o

	case u < math.MaxUint32:
		o, n := ensure(b, 5)
		putMuint32(o[n:], uint32(u))
		return o

	default:
		o, n := ensure(b, 9)
		putMuint64(o[n:], u)
		return o

	}
}

// AppendUint appends a uint to the slice
func AppendUint(b []byte, u uint) []byte { return AppendUint64(b, uint64(u)) }

// AppendUint8 appends a uint8 to the slice
func AppendUint8(b []byte, u uint8) []byte { return AppendUint64(b, uint64(u)) }

// AppendByte is analagous to AppendUint8
func AppendByte(b []byte, u byte) []byte { return AppendUint8(b, uint8(u)) }

// AppendUint16 appends a uint16 to the slice
func AppendUint16(b []byte, u uint16) []byte { return AppendUint64(b, uint64(u)) }

// AppendUint32 appends a uint32 to the slice
func AppendUint32(b []byte, u uint32) []byte { return AppendUint64(b, uint64(u)) }

// AppendBytes appends bytes to the slice as MessagePack 'bin' data
func AppendBytes(b []byte, bts []byte) []byte {
	sz := uint32(len(bts))
	o, n := ensure(b, BytesPrefixSize+len(bts))
	switch {
	case sz < math.MaxUint8:
		prefixu8(o[n:], mbin8, uint8(sz))
		n += 2
	case sz < math.MaxUint16:
		prefixu16(o[n:], mbin16, uint16(sz))
		n += 3
	default:
		prefixu32(o[n:], mbin32, sz)
		n += 5
	}
	return o[:n+copy(o[n:], bts)]
}

// AppendBool appends a bool to the slice
func AppendBool(b []byte, t bool) []byte {
	if t {
		return append(b, mtrue)
	}
	return append(b, mfalse)
}

// AppendString appends a string as a MessagePack 'str' to the slice
func AppendString(b []byte, s string) []byte {
	sz := uint32(len(s))
	o, n := ensure(b, BytesPrefixSize+len(s))
	switch {
	case sz < 32:
		o[n] = wfixstr(uint8(sz))
		n++
	case sz < math.MaxUint8:
		prefixu8(o[n:], mstr8, uint8(sz))
		n += 2
	case sz < math.MaxUint16:
		prefixu16(o[n:], mstr16, uint16(sz))
		n += 3
	default:
		prefixu32(o[n:], mstr32, sz)
		n += 5
	}
	return o[:n+copy(o[n:], s)]
}

// AppendComplex64 appends a complex64 to the slice as a MessagePack extension
func AppendComplex64(b []byte, c complex64) []byte {
	o, n := ensure(b, Complex64Size)
	o[n] = mfixext8
	o[n+1] = Complex64Extension
	memcpy8(unsafe.Pointer(&o[n+2]), unsafe.Pointer(&c))
	return o
}

// AppendComplex128 appends a complex128 to the slice as a MessagePack extension
func AppendComplex128(b []byte, c complex128) []byte {
	o, n := ensure(b, Complex128Size)
	o[n] = mfixext16
	o[n+1] = Complex128Extension
	memcpy16(unsafe.Pointer(&o[n+2]), unsafe.Pointer(&c))
	return o
}

// AppendTime appends a time.Time to the slice as a MessagePack extension
func AppendTime(b []byte, t time.Time) []byte {
	o, n := ensure(b, TimeSize)
	bts, _ := t.MarshalBinary()
	o[n] = mfixext16
	o[n+1] = TimeExtension
	copy(o[n+2:], bts)
	return o
}

// AppendMapStrStr appends a map[string]string to the slice
// as a MessagePack map with 'str'-type keys and values
func AppendMapStrStr(b []byte, m map[string]string) []byte {
	sz := uint32(len(m))
	b = AppendMapHeader(b, sz)
	for key, val := range m {
		b = AppendString(b, key)
		b = AppendString(b, val)
	}
	return b
}

// AppendMapStrIntf appends a map[string]interface{} to the slice
// as a MessagePack map with 'str'-type keys.
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

// AppendIntf appends the concrete type of 'i' to the
// provided []byte. 'i' must be one of the following:
//  - 'nil'
//  - A bool, float, string, []byte, int, uint, or complex
//  - A map[string]interface{} or map[string]string
//  - A []T, where T is another supported type
//  - A *T, where T is another supported type
//  - A type that satisfieds the msgp.Marshaler interface
//  - A type that satisfies the msgp.Extension interface
func AppendIntf(b []byte, i interface{}) ([]byte, error) {
	if m, ok := i.(Marshaler); ok {
		return m.MarshalMsg(b)
	}
	if ext, ok := i.(Extension); ok {
		return AppendExtension(b, ext)
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

	default:
		return b, &ErrUnsupportedType{T: v.Type()}
	}
}
