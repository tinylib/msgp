//go:build go1.23

package msgp

import (
	"fmt"
	"iter"
	"time"
)

// ArrayExtraTypes is a type constraint that includes all types that can be
// decoded from an array.
// Even though 'any' type can be used, its pointer must implement the
// Decodable when using Reader or Unmarhaler interface when reading from bytes.
type ArrayExtraTypes interface {
	bool | string | []byte | complex64 | complex128 | time.Time | time.Duration | any
}

// ReadArray returns an iterator that can be used to iterate over the elements
// of an array in the MessagePack data while being read by the provided Reader.
// The type parameter V specifies the type of the elements in the array.
// The type parameter V must be ArrayExtraTypes or a type whose
// pointer implements the Decodable interface.
// Use ReadNumberArray for numbers.
// The returned iterator implements the iter.Seq[V] interface,
// allowing for sequential access to the array elements.
func ReadArray[V ArrayExtraTypes](m *Reader) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		// Check if nil
		if m.IsNil() {
			m.ReadNil()
			return
		}
		// Regular array.
		var x V
		length, err := m.ReadArrayHeader()
		if err != nil {
			yield(x, fmt.Errorf("cannot read array header: %w", err))
			return
		}
		switch any(x).(type) {
		case string:
			for range length {
				v, err := m.ReadString()
				if !yield(any(v).(V), err) {
					return
				}
			}
		case []byte:
			for range length {
				v, err := m.ReadBytes(nil)
				if !yield(any(v).(V), err) {
					return
				}
			}
		case bool:
			for range length {
				v, err := m.ReadBool()
				if !yield(any(v).(V), err) {
					return
				}
			}
		case time.Time:
			for range length {
				v, err := m.ReadTime()
				if !yield(any(v).(V), err) {
					return
				}
			}
		case time.Duration:
			for range length {
				v, err := m.ReadDuration()
				if !yield(any(v).(V), err) {
					return
				}
			}
		case complex64:
			for range length {
				v, err := m.ReadComplex64()
				if !yield(any(v).(V), err) {
					return
				}
			}
		case complex128:
			for range length {
				v, err := m.ReadComplex128()
				if !yield(any(v).(V), err) {
					return
				}
			}
		default:
			for range length {
				var v V
				ptr := &v
				if dc, ok := any(ptr).(Decodable); ok {
					err = dc.DecodeMsg(m)
					if !yield(v, err) {
						return
					}
				} else {
					yield(v, fmt.Errorf("cannot decode into type %T", ptr))
					return
				}
			}
		}
	}
}

// MapKeyTypes are possible key types. Usually this is a string.
// Even though 'any' type can be used, its pointer must implement the Decodable interface.
type MapKeyTypes interface {
	NumberTypes | bool | string | []byte | complex64 | complex128 | time.Time | time.Duration | any
}

// MapValueTypes are possible value types.
// Even though 'any' type can be used, its pointer must implement the Decodable interface.
type MapValueTypes interface {
	NumberTypes | bool | string | []byte | complex64 | complex128 | time.Time | time.Duration | any
}

// ReadMap returns an iterator that can be used to iterate over the elements
// of a map in the MessagePack data while being read by the provided Reader.
// The type parameters K and V specify the types of the keys and values in the map.
// The returned iterator implements the iter.Seq2[K, V] interface,
// allowing for sequential access to the map elements.
// The returned function can be used to read any error that occurred during iteration when iteration is done.
func ReadMap[K MapKeyTypes, V MapValueTypes](m *Reader) (iter.Seq2[K, V], func() error) {
	var err error
	return func(yield func(K, V) bool) {
		var sz uint32
		if m.IsNil() {
			err = m.ReadNil()
			return
		}
		sz, err = m.ReadMapHeader()
		if err != nil {
			err = fmt.Errorf("cannot read map header: %w", err)
			return
		}

		// Prepare per-type readers once, avoid switching for each element.
		// Key reader.
		var readKey func() (K, error)
		{
			var keyZero K
			switch any(keyZero).(type) {
			case string:
				readKey = func() (K, error) {
					v, e := m.ReadString()
					return any(v).(K), e
				}
			case []byte:
				readKey = func() (K, error) {
					v, e := m.ReadBytes(nil)
					return any(v).(K), e
				}
			case bool:
				readKey = func() (K, error) {
					v, e := m.ReadBool()
					return any(v).(K), e
				}
			case time.Time:
				readKey = func() (K, error) {
					v, e := m.ReadTime()
					return any(v).(K), e
				}
			case time.Duration:
				readKey = func() (K, error) {
					v, e := m.ReadDuration()
					return any(v).(K), e
				}
			case complex64:
				readKey = func() (K, error) {
					v, e := m.ReadComplex64()
					return any(v).(K), e
				}
			case complex128:
				readKey = func() (K, error) {
					v, e := m.ReadComplex128()
					return any(v).(K), e
				}
			case uint8:
				readKey = func() (K, error) {
					v, e := m.ReadUint8()
					return any(v).(K), e
				}
			case uint16:
				readKey = func() (K, error) {
					v, e := m.ReadUint16()
					return any(v).(K), e
				}
			case uint32:
				readKey = func() (K, error) {
					v, e := m.ReadUint32()
					return any(v).(K), e
				}
			case uint64:
				readKey = func() (K, error) {
					v, e := m.ReadUint64()
					return any(v).(K), e
				}
			case uint:
				readKey = func() (K, error) {
					v, e := m.ReadUint()
					return any(v).(K), e
				}
			case int8:
				readKey = func() (K, error) {
					v, e := m.ReadInt8()
					return any(v).(K), e
				}
			case int16:
				readKey = func() (K, error) {
					v, e := m.ReadInt16()
					return any(v).(K), e
				}
			case int32:
				readKey = func() (K, error) {
					v, e := m.ReadInt32()
					return any(v).(K), e
				}
			case int64:
				readKey = func() (K, error) {
					v, e := m.ReadInt64()
					return any(v).(K), e
				}
			case int:
				readKey = func() (K, error) {
					v, e := m.ReadInt()
					return any(v).(K), e
				}
			case float32:
				readKey = func() (K, error) {
					v, e := m.ReadFloat32()
					return any(v).(K), e
				}
			case float64:
				readKey = func() (K, error) {
					v, e := m.ReadFloat64()
					return any(v).(K), e
				}
			default:
				readKey = func() (K, error) {
					var k K
					ptr := &k
					if dc, ok := any(ptr).(Decodable); ok {
						return k, dc.DecodeMsg(m)
					}
					return k, fmt.Errorf("cannot decode key into type %T", ptr)
				}
			}
		}

		// Value reader.
		var readVal func() (V, error)
		{
			var valZero V
			switch any(valZero).(type) {
			case string:
				readVal = func() (V, error) {
					v, e := m.ReadString()
					return any(v).(V), e
				}
			case []byte:
				readVal = func() (V, error) {
					v, e := m.ReadBytes(nil)
					return any(v).(V), e
				}
			case bool:
				readVal = func() (V, error) {
					v, e := m.ReadBool()
					return any(v).(V), e
				}
			case time.Time:
				readVal = func() (V, error) {
					v, e := m.ReadTime()
					return any(v).(V), e
				}
			case time.Duration:
				readVal = func() (V, error) {
					v, e := m.ReadDuration()
					return any(v).(V), e
				}
			case complex64:
				readVal = func() (V, error) {
					v, e := m.ReadComplex64()
					return any(v).(V), e
				}
			case complex128:
				readVal = func() (V, error) {
					v, e := m.ReadComplex128()
					return any(v).(V), e
				}
			case int8:
				readVal = func() (V, error) {
					v, e := m.ReadInt8()
					return any(v).(V), e
				}
			case int16:
				readVal = func() (V, error) {
					v, e := m.ReadInt16()
					return any(v).(V), e
				}
			case int32:
				readVal = func() (V, error) {
					v, e := m.ReadInt32()
					return any(v).(V), e
				}
			case int64:
				readVal = func() (V, error) {
					v, e := m.ReadInt64()
					return any(v).(V), e
				}
			case int:
				readVal = func() (V, error) {
					v, e := m.ReadInt()
					return any(v).(V), e
				}
			case float32:
				readVal = func() (V, error) {
					v, e := m.ReadFloat32()
					return any(v).(V), e
				}
			case float64:
				readVal = func() (V, error) {
					v, e := m.ReadFloat64()
					return any(v).(V), e
				}
			case uint8:
				readVal = func() (V, error) {
					v, e := m.ReadUint8()
					return any(v).(V), e
				}
			case uint16:
				readVal = func() (V, error) {
					v, e := m.ReadUint16()
					return any(v).(V), e
				}
			case uint32:
				readVal = func() (V, error) {
					v, e := m.ReadUint32()
					return any(v).(V), e
				}
			case uint64:
				readVal = func() (V, error) {
					v, e := m.ReadUint64()
					return any(v).(V), e
				}
			case uint:
				readVal = func() (V, error) {
					v, e := m.ReadUint()
					return any(v).(V), e
				}
			default:
				readVal = func() (V, error) {
					var v V
					ptr := &v
					if dc, ok := any(ptr).(Decodable); ok {
						return v, dc.DecodeMsg(m)
					}
					return v, fmt.Errorf("cannot decode value into type %T", ptr)
				}
			}
		}

		for range sz {
			var k K
			k, err = readKey()
			if err != nil {
				err = fmt.Errorf("cannot read key: %w", err)
				return
			}
			var v V
			v, err = readVal()
			if err != nil {
				err = fmt.Errorf("cannot read value: %w", err)
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}, func() error { return err }
}

// NumberTypes is a type constraint that includes all numeric types.
type NumberTypes interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64 | float32 | float64
}

// ReadNumberArray returns an iterator that can be used to iterate over the elements
// of an array in the MessagePack data while being read by the provided Reader.
// The type parameter V specifies the type of the elements in the array.
// The returned iterator implements the iter.Seq[V] interface,
// allowing for sequential access to the array elements.
func ReadNumberArray[V NumberTypes](m *Reader) iter.Seq2[V, error] {
	return func(yield func(V, error) bool) {
		if m.IsNil() {
			m.ReadNil()
			return
		}
		var x V
		length, err := m.ReadArrayHeader()
		if err != nil {
			yield(x, fmt.Errorf("cannot read array header: %w", err))
			return
		}

		switch any(x).(type) {
		case uint8:
			for range length {
				v, err := m.ReadUint8()
				if !yield(V(v), err) {
					return
				}
			}
		case uint16:
			for range length {
				v, err := m.ReadUint16()
				if !yield(V(v), err) {
					return
				}
			}
		case uint32:
			for range length {
				v, err := m.ReadUint32()
				if !yield(V(v), err) {
					return
				}
			}
		case uint64:
			for range length {
				v, err := m.ReadUint64()
				if !yield(V(v), err) {
					return
				}
			}
		case uint:
			for range length {
				v, err := m.ReadUint()
				if !yield(V(v), err) {
					return
				}
			}
		case int8:
			for range length {
				v, err := m.ReadInt8()
				if !yield(V(v), err) {
					return
				}
			}
		case int16:
			for range length {
				v, err := m.ReadInt16()
				if !yield(V(v), err) {
					return
				}
			}
		case int32:
			for range length {
				v, err := m.ReadInt32()
				if !yield(V(v), err) {
					return
				}
			}
		case int64:
			for range length {
				v, err := m.ReadInt64()
				if !yield(V(v), err) {
					return
				}
			}
		case int:
			for range length {
				v, err := m.ReadInt()
				if !yield(V(v), err) {
					return
				}
			}
		case float32:
			for range length {
				v, err := m.ReadFloat32()
				if !yield(V(v), err) {
					return
				}
			}
		case float64:
			for range length {
				v, err := m.ReadFloat64()
				if !yield(V(v), err) {
					return
				}
			}
		default:
			panic("unreachable")
		}
	}
}

// ReadNumberArrayBytes returns an iterator that can be used to iterate over the elements
// of an array in the MessagePack data.
// The type parameter V specifies the type of the elements in the array.
// The returned iterator implements the iter.Seq[V] interface,
// allowing for sequential access to the array elements.
// After the iterator is exhausted, the remaining bytes in the buffer
// and any error can be read by calling the returned function.
func ReadNumberArrayBytes[V NumberTypes](b []byte) (iter.Seq[V], func() (remain []byte, err error)) {
	if IsNil(b) {
		b, err := ReadNilBytes(b)
		return func(yield func(V) bool) {}, func() ([]byte, error) { return b, err }
	}
	// Regular array.
	sz, b, err := ReadArrayHeaderBytes(b)
	if err != nil {
		return func(yield func(V) bool) {}, func() ([]byte, error) { return b, fmt.Errorf("cannot read array header: %w", err) }
	}

	var readValue func() (V, error)
	var v V
	switch any(v).(type) {
	case uint8:
		readValue = func() (V, error) {
			var val uint8
			val, b, err = ReadUint8Bytes(b)
			return V(val), err
		}
	case uint16:
		readValue = func() (V, error) {
			var val uint16
			val, b, err = ReadUint16Bytes(b)
			return V(val), err
		}
	case uint32:
		readValue = func() (V, error) {
			var val uint32
			val, b, err = ReadUint32Bytes(b)
			return V(val), err
		}
	case uint64:
		readValue = func() (V, error) {
			var val uint64
			val, b, err = ReadUint64Bytes(b)
			return V(val), err
		}
	case uint:
		readValue = func() (V, error) {
			var val uint
			val, b, err = ReadUintBytes(b)
			return V(val), err
		}
	case int8:
		readValue = func() (V, error) {
			var val int8
			val, b, err = ReadInt8Bytes(b)
			return V(val), err
		}
	case int16:
		readValue = func() (V, error) {
			var val int16
			val, b, err = ReadInt16Bytes(b)
			return V(val), err
		}
	case int32:
		readValue = func() (V, error) {
			var val int32
			val, b, err = ReadInt32Bytes(b)
			return V(val), err
		}
	case int64:
		readValue = func() (V, error) {
			var val int64
			val, b, err = ReadInt64Bytes(b)
			return V(val), err
		}
	case int:
		readValue = func() (V, error) {
			var val int
			val, b, err = ReadIntBytes(b)
			return V(val), err
		}
	case float32:
		readValue = func() (V, error) {
			var val float32
			val, b, err = ReadFloat32Bytes(b)
			return V(val), err
		}
	case float64:
		readValue = func() (V, error) {
			var val float64
			val, b, err = ReadFloat64Bytes(b)
			return V(val), err
		}
	default:
		panic("unreachable")
	}
	return func(yield func(V) bool) {
		for sz > 0 {
			v, err = readValue()
			if err != nil || !yield(v) {
				return
			}
			sz--
		}
	}, func() ([]byte, error) { return b, err }
}

// ReadArrayBytes returns an iterator that can be used to iterate over the elements
// of an array in the MessagePack data while being read by the provided Reader.
// The type parameter V specifies the type of the elements in the array.
// The type parameter V must be bool, string, []byte or a type whose
// pointer implements the Unmarshaler interface.
// Use ReadNumberArrayBytes for numbers.
// The returned iterator implements the iter.Seq[V] interface,
// allowing for sequential access to the array elements.
// Byte slices will reference the same underlying data.
// After the iterator is exhausted, the remaining bytes in the buffer
// and any error can be read by calling the returned function.
func ReadArrayBytes[V ArrayExtraTypes](b []byte) (iter.Seq[V], func() (remain []byte, err error)) {
	if IsNil(b) {
		b, err := ReadNilBytes(b)
		return func(yield func(V) bool) {}, func() ([]byte, error) { return b, err }
	}
	sz, b, err := ReadArrayHeaderBytes(b)
	if err != nil {
		return nil, func() ([]byte, error) { return b, err }
	}
	return func(yield func(V) bool) {
			var x V
			switch any(x).(type) {
			case string:
				for range sz {
					var v string
					v, b, err = ReadStringBytes(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			case []byte:
				for range sz {
					var v []byte
					v, b, err = ReadBytesZC(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			case bool:
				for range sz {
					var v bool
					v, b, err = ReadBoolBytes(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			case time.Time:
				for range sz {
					var v time.Time
					v, b, err = ReadTimeBytes(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			case time.Duration:
				for range sz {
					var v time.Duration
					v, b, err = ReadDurationBytes(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			case complex64:
				for range sz {
					var v complex64
					v, b, err = ReadComplex64Bytes(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			case complex128:
				for range sz {
					var v complex128
					v, b, err = ReadComplex128Bytes(b)
					if err != nil || !yield(any(v).(V)) {
						return
					}
				}
			default:
				for range sz {
					var v V
					ptr := &v
					if um, ok := any(ptr).(Unmarshaler); ok {
						b, err = um.UnmarshalMsg(b)
						if err != nil || !yield(v) {
							return
						}
					} else {
						err = fmt.Errorf("cannot unmarshal into type %T", ptr)
						return
					}
				}
			}
		}, func() (remain []byte, err error) {
			return b, err
		}
}

// ReadMapBytes returns an iterator over key/value pairs from a MessagePack map encoded in b.
// The iterator yields K,V pairs and this function also returns a closure to obtain the remaining bytes and any error.
// It avoids per-element type switches by precomputing readKey/readVal funcs based on K and V.
func ReadMapBytes[K MapKeyTypes, V MapValueTypes](b []byte) (iter.Seq2[K, V], func() (remain []byte, err error)) {
	var err error
	var sz uint32
	if IsNil(b) {
		b, err = ReadNilBytes(b)
		return func(yield func(K, V) bool) {}, func() ([]byte, error) { return b, err }
	}
	sz, b, err = ReadMapHeaderBytes(b)
	if err != nil {
		return func(yield func(K, V) bool) {}, func() ([]byte, error) { return b, err }
	}

	// Precompute key reader
	var readKey func() (K, error)
	{
		var keyZero K
		switch any(keyZero).(type) {
		case string:
			readKey = func() (K, error) {
				v, e, er := ReadStringBytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case []byte:
			// Map keys can be str or bin; use specialized helper that accepts both.
			readKey = func() (K, error) {
				v, e, er := ReadMapKeyZC(b)
				b, err = e, er
				return any(v).(K), err
			}
		case bool:
			readKey = func() (K, error) {
				v, e, er := ReadBoolBytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case time.Time:
			readKey = func() (K, error) {
				v, e, er := ReadTimeBytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case time.Duration:
			readKey = func() (K, error) {
				v, e, er := ReadDurationBytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case complex64:
			readKey = func() (K, error) {
				v, e, er := ReadComplex64Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case complex128:
			readKey = func() (K, error) {
				v, e, er := ReadComplex128Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case uint8:
			readKey = func() (K, error) {
				v, e, er := ReadUint8Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case uint16:
			readKey = func() (K, error) {
				v, e, er := ReadUint16Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case uint32:
			readKey = func() (K, error) {
				v, e, er := ReadUint32Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case uint64:
			readKey = func() (K, error) {
				v, e, er := ReadUint64Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case uint:
			readKey = func() (K, error) {
				v, e, er := ReadUintBytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case int8:
			readKey = func() (K, error) {
				v, e, er := ReadInt8Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case int16:
			readKey = func() (K, error) {
				v, e, er := ReadInt16Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case int32:
			readKey = func() (K, error) {
				v, e, er := ReadInt32Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case int64:
			readKey = func() (K, error) {
				v, e, er := ReadInt64Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case int:
			readKey = func() (K, error) {
				v, e, er := ReadIntBytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case float32:
			readKey = func() (K, error) {
				v, e, er := ReadFloat32Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		case float64:
			readKey = func() (K, error) {
				v, e, er := ReadFloat64Bytes(b)
				b, err = e, er
				return any(v).(K), err
			}
		default:
			readKey = func() (K, error) {
				var k K
				ptr := &k
				if um, ok := any(ptr).(Unmarshaler); ok {
					var e error
					b, e = um.UnmarshalMsg(b)
					return k, e
				}
				return k, fmt.Errorf("cannot unmarshal key into type %T", ptr)
			}
		}
	}

	// Precompute value reader
	var readVal func() (V, error)
	{
		var valZero V
		switch any(valZero).(type) {
		case string:
			readVal = func() (V, error) {
				v, e, er := ReadStringBytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case []byte:
			// For values, zero-copy read of bin/str payload.
			readVal = func() (V, error) {
				v, e, er := ReadBytesZC(b)
				b, err = e, er
				return any(v).(V), err
			}
		case bool:
			readVal = func() (V, error) {
				v, e, er := ReadBoolBytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case time.Time:
			readVal = func() (V, error) {
				v, e, er := ReadTimeBytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case time.Duration:
			readVal = func() (V, error) {
				v, e, er := ReadDurationBytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case complex64:
			readVal = func() (V, error) {
				v, e, er := ReadComplex64Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case complex128:
			readVal = func() (V, error) {
				v, e, er := ReadComplex128Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case int8:
			readVal = func() (V, error) {
				v, e, er := ReadInt8Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case int16:
			readVal = func() (V, error) {
				v, e, er := ReadInt16Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case int32:
			readVal = func() (V, error) {
				v, e, er := ReadInt32Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case int64:
			readVal = func() (V, error) {
				v, e, er := ReadInt64Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case int:
			readVal = func() (V, error) {
				v, e, er := ReadIntBytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case float32:
			readVal = func() (V, error) {
				v, e, er := ReadFloat32Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case float64:
			readVal = func() (V, error) {
				v, e, er := ReadFloat64Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case uint8:
			readVal = func() (V, error) {
				v, e, er := ReadUint8Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case uint16:
			readVal = func() (V, error) {
				v, e, er := ReadUint16Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case uint32:
			readVal = func() (V, error) {
				v, e, er := ReadUint32Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case uint64:
			readVal = func() (V, error) {
				v, e, er := ReadUint64Bytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		case uint:
			readVal = func() (V, error) {
				v, e, er := ReadUintBytes(b)
				b, err = e, er
				return any(v).(V), err
			}
		default:
			readVal = func() (V, error) {
				var v V
				ptr := &v
				if um, ok := any(ptr).(Unmarshaler); ok {
					var e error
					b, e = um.UnmarshalMsg(b)
					return v, e
				}
				return v, fmt.Errorf("cannot unmarshal value into type %T", ptr)
			}
		}
	}

	return func(yield func(K, V) bool) {
		for range sz {
			k, er := readKey()
			if er != nil {
				err = fmt.Errorf("cannot read map key: %w", er)
				return
			}
			v, er := readVal()
			if er != nil {
				err = fmt.Errorf("cannot read map value: %w", er)
				return
			}
			if !yield(k, v) {
				return
			}
		}
	}, func() ([]byte, error) { return b, err }
}
