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
		// Assuming Reader has a method to read array length
		var x V
		length, err := m.ReadArrayHeader()
		if err != nil {
			yield(x, err)
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
		// Assuming Reader has a method to read array length
		var length uint32
		length, err = m.ReadArrayHeader()
		if err != nil {
			return
		}
		for range length {
			var key K
			switch v := any(key).(type) {
			case string:
				if v, err = m.ReadString(); err != nil {
					return
				}
				key = (any)(v).(K)
			case []byte:
				if v, err = m.ReadBytes(nil); err != nil {
					return
				}
				key = (any)(v).(K)
			case bool:
				if v, err = m.ReadBool(); err != nil {
					return
				}
				key = (any)(v).(K)
			case time.Time:
				if v, err = m.ReadTime(); err != nil {
					return
				}
				key = (any)(v).(K)
			case time.Duration:
				if v, err = m.ReadDuration(); err != nil {
					return
				}
				key = (any)(v).(K)
			case complex64:
				if v, err = m.ReadComplex64(); err != nil {
					return
				}
				key = (any)(v).(K)
			case complex128:
				if v, err = m.ReadComplex128(); err != nil {
					return
				}
				key = (any)(v).(K)
			case uint8:
				if v, err = m.ReadUint8(); err != nil {
					return
				}
				key = (any)(v).(K)
			case uint16:
				if v, err = m.ReadUint16(); err != nil {
					return
				}
				key = (any)(v).(K)
			case uint32:
				if v, err = m.ReadUint32(); err != nil {
					return
				}
				key = (any)(v).(K)
			case uint64:
				if v, err = m.ReadUint64(); err != nil {
					return
				}
				key = (any)(v).(K)
			case uint:
				if v, err = m.ReadUint(); err != nil {
					return
				}
				key = (any)(v).(K)
			case int8:
				if v, err = m.ReadInt8(); err != nil {
					return
				}
				key = (any)(v).(K)
			case int16:
				if v, err = m.ReadInt16(); err != nil {
					return
				}
				key = (any)(v).(K)
			case int32:
				if v, err = m.ReadInt32(); err != nil {
					return
				}
				key = (any)(v).(K)
			case int64:
				if v, err = m.ReadInt64(); err != nil {
					return
				}
				key = (any)(v).(K)
			case int:
				if v, err = m.ReadInt(); err != nil {
					return
				}
				key = (any)(v).(K)
			case float32:
				if v, err = m.ReadFloat32(); err != nil {
					return
				}
				key = (any)(v).(K)
			case float64:
				if v, err = m.ReadFloat64(); err != nil {
					return
				}
				key = (any)(v).(K)
			default:
				_ = v
				ptr := &key
				if dc, ok := any(ptr).(Decodable); ok {
					if err = dc.DecodeMsg(m); err != nil {
						return
					}
				} else {
					err = fmt.Errorf("cannot decode key into type %T", ptr)
					return
				}
			}

			var val V
			switch v := any(key).(type) {
			case string:
				if v, err = m.ReadString(); err != nil {
					return
				}
				val = (any)(v).(V)
			case []byte:
				if v, err = m.ReadBytes(nil); err != nil {
					return
				}
				val = (any)(v).(V)
			case bool:
				if v, err = m.ReadBool(); err != nil {
					return
				}
				val = (any)(v).(V)
			case time.Time:
				if v, err = m.ReadTime(); err != nil {
					return
				}
				val = (any)(v).(V)
			case time.Duration:
				if v, err = m.ReadDuration(); err != nil {
					return
				}
				val = (any)(v).(V)
			case complex64:
				if v, err = m.ReadComplex64(); err != nil {
					return
				}
				val = (any)(v).(V)
			case complex128:
				if v, err = m.ReadComplex128(); err != nil {
					return
				}
				val = (any)(v).(V)
			case int8:
				if v, err = m.ReadInt8(); err != nil {
					return
				}
				val = (any)(v).(V)
			case int16:
				if v, err = m.ReadInt16(); err != nil {
					return
				}
				val = (any)(v).(V)
			case int32:
				if v, err = m.ReadInt32(); err != nil {
					return
				}
				val = (any)(v).(V)
			case int64:
				if v, err = m.ReadInt64(); err != nil {
					return
				}
				val = (any)(v).(V)
			case int:
				if v, err = m.ReadInt(); err != nil {
					return
				}
				val = (any)(v).(V)
			case float32:
				if v, err = m.ReadFloat32(); err != nil {
					return
				}
				val = (any)(v).(V)
			case float64:
				if v, err = m.ReadFloat64(); err != nil {
					return
				}
				val = (any)(v).(V)
			case uint8:
				if v, err = m.ReadUint8(); err != nil {
					return
				}
				val = (any)(v).(V)
			case uint16:
				if v, err = m.ReadUint16(); err != nil {
					return
				}
				val = (any)(v).(V)
			case uint32:
				if v, err = m.ReadUint32(); err != nil {
					return
				}
				val = (any)(v).(V)
			case uint64:
				if v, err = m.ReadUint64(); err != nil {
					return
				}
				val = (any)(v).(V)
			case uint:
				if v, err = m.ReadUint(); err != nil {
					return
				}
				val = (any)(v).(V)

			default:
				_ = v
				ptr := &val
				if dc, ok := any(ptr).(Decodable); ok {
					if err = dc.DecodeMsg(m); err != nil {
						return
					}
				} else {
					err = fmt.Errorf("cannot decode value into type %T", ptr)
					return
				}
			}
			if !yield(key, val) {
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
		// Assuming Reader has a method to read array length
		var x V
		length, err := m.ReadArrayHeader()
		if err != nil {
			yield(x, err)
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
	sz, b, err := ReadArrayHeaderBytes(b)
	if err != nil {
		return nil, func() ([]byte, error) { return b, err }
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

// ReadMapBytes returns an iterator over key-value pairs of a map encoded in MessagePack bytes.
// The type parameters K and V specify the types of the keys and values in the map.
// The iterator yields pairs in wire order. After iteration completes (or stops early),
// call the returned tail function to get any remaining bytes and the first error encountered (if any).
// K and V must be one of the supported built-in types, or pointers to types implementing Unmarshaler.
// Byte slices will reference the same underlying data.
func ReadMapBytes[K MapKeyTypes, V MapValueTypes](b []byte) (iter.Seq2[K, V], func() (remain []byte, err error)) {
	sz, b, err := ReadMapHeaderBytes(b)
	if err != nil {
		return nil, func() ([]byte, error) { return b, err }
	}

	readKey := func() (K, error) {
		var key K
		switch any(key).(type) {
		case string:
			var v string
			v, b, err = ReadStringBytes(b)
			key = any(v).(K)
		case []byte:
			var v []byte
			v, b, err = ReadBytesZC(b)
			key = any(v).(K)
		case bool:
			var v bool
			v, b, err = ReadBoolBytes(b)
			key = any(v).(K)
		case time.Time:
			var v time.Time
			v, b, err = ReadTimeBytes(b)
			key = any(v).(K)
		case time.Duration:
			var v time.Duration
			v, b, err = ReadDurationBytes(b)
			key = any(v).(K)
		case complex64:
			var v complex64
			v, b, err = ReadComplex64Bytes(b)
			key = any(v).(K)
		case complex128:
			var v complex128
			v, b, err = ReadComplex128Bytes(b)
			key = any(v).(K)
		case uint8:
			var v uint8
			v, b, err = ReadUint8Bytes(b)
			key = any(v).(K)
		case uint16:
			var v uint16
			v, b, err = ReadUint16Bytes(b)
			key = any(v).(K)
		case uint32:
			var v uint32
			v, b, err = ReadUint32Bytes(b)
			key = any(v).(K)
		case uint64:
			var v uint64
			v, b, err = ReadUint64Bytes(b)
			key = any(v).(K)
		case uint:
			var v uint
			v, b, err = ReadUintBytes(b)
			key = any(v).(K)
		case int8:
			var v int8
			v, b, err = ReadInt8Bytes(b)
			key = any(v).(K)
		case int16:
			var v int16
			v, b, err = ReadInt16Bytes(b)
			key = any(v).(K)
		case int32:
			var v int32
			v, b, err = ReadInt32Bytes(b)
			key = any(v).(K)
		case int64:
			var v int64
			v, b, err = ReadInt64Bytes(b)
			key = any(v).(K)
		case int:
			var v int
			v, b, err = ReadIntBytes(b)
			key = any(v).(K)
		case float32:
			var v float32
			v, b, err = ReadFloat32Bytes(b)
			key = any(v).(K)
		case float64:
			var v float64
			v, b, err = ReadFloat64Bytes(b)
			key = any(v).(K)
		default:
			// Fallback for custom types implementing Unmarshaler
			ptr := &key
			if um, ok := any(ptr).(Unmarshaler); ok {
				b, err = um.UnmarshalMsg(b)
			} else {
				err = fmt.Errorf("cannot unmarshal key into type %T", ptr)
			}
		}
		return key, err
	}

	readVal := func() (V, error) {
		var val V
		switch any(val).(type) {
		case string:
			var v string
			v, b, err = ReadStringBytes(b)
			val = any(v).(V)
		case []byte:
			var v []byte
			v, b, err = ReadBytesZC(b)
			val = any(v).(V)
		case bool:
			var v bool
			v, b, err = ReadBoolBytes(b)
			val = any(v).(V)
		case time.Time:
			var v time.Time
			v, b, err = ReadTimeBytes(b)
			val = any(v).(V)
		case time.Duration:
			var v time.Duration
			v, b, err = ReadDurationBytes(b)
			val = any(v).(V)
		case complex64:
			var v complex64
			v, b, err = ReadComplex64Bytes(b)
			val = any(v).(V)
		case complex128:
			var v complex128
			v, b, err = ReadComplex128Bytes(b)
			val = any(v).(V)
		case uint8:
			var v uint8
			v, b, err = ReadUint8Bytes(b)
			val = any(v).(V)
		case uint16:
			var v uint16
			v, b, err = ReadUint16Bytes(b)
			val = any(v).(V)
		case uint32:
			var v uint32
			v, b, err = ReadUint32Bytes(b)
			val = any(v).(V)
		case uint64:
			var v uint64
			v, b, err = ReadUint64Bytes(b)
			val = any(v).(V)
		case uint:
			var v uint
			v, b, err = ReadUintBytes(b)
			val = any(v).(V)
		case int8:
			var v int8
			v, b, err = ReadInt8Bytes(b)
			val = any(v).(V)
		case int16:
			var v int16
			v, b, err = ReadInt16Bytes(b)
			val = any(v).(V)
		case int32:
			var v int32
			v, b, err = ReadInt32Bytes(b)
			val = any(v).(V)
		case int64:
			var v int64
			v, b, err = ReadInt64Bytes(b)
			val = any(v).(V)
		case int:
			var v int
			v, b, err = ReadIntBytes(b)
			val = any(v).(V)
		case float32:
			var v float32
			v, b, err = ReadFloat32Bytes(b)
			val = any(v).(V)
		case float64:
			var v float64
			v, b, err = ReadFloat64Bytes(b)
			val = any(v).(V)
		default:
			// Fallback for custom types implementing Unmarshaler
			ptr := &val
			if um, ok := any(ptr).(Unmarshaler); ok {
				b, err = um.UnmarshalMsg(b)
			} else {
				err = fmt.Errorf("cannot unmarshal value into type %T", ptr)
			}
		}
		return val, err
	}

	return func(yield func(K, V) bool) {
			for sz > 0 {
				k, e := readKey()
				if e != nil {
					err = e
					return
				}
				v, e := readVal()
				if e != nil {
					err = e
					return
				}
				if !yield(k, v) {
					return
				}
				sz--
			}
		}, func() ([]byte, error) {
			return b, err
		}
}
