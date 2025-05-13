package _generated

import "github.com/tinylib/msgp/msgp/setof"

//go:generate msgp
type SetOf struct {
	String        setof.String
	StringSorted  setof.StringSorted
	Int           setof.Int
	IntSorted     setof.IntSorted
	Int64         setof.Int64
	Int64Sorted   setof.Int64Sorted
	Int32         setof.Int32
	Int32Sorted   setof.Int32Sorted
	Int16         setof.Int16
	Int16Sorted   setof.Int16Sorted
	Int8          setof.Int8
	Int8Sorted    setof.Int8Sorted
	Byte          setof.Byte
	ByteSorted    setof.ByteSorted
	Uint          setof.Uint
	UintSorted    setof.UintSorted
	Uint64        setof.Uint64
	Uint64Sorted  setof.Uint64Sorted
	Uint32        setof.Uint32
	Uint32Sorted  setof.Uint32Sorted
	Uint16        setof.Uint16
	Uint16Sorted  setof.Uint16Sorted
	Uint8         setof.Uint8
	Uint8Sorted   setof.Uint8Sorted
	Float32       setof.Float32
	Float32Sorted setof.Float32Sorted
	Float64       setof.Float64
	Float64Sorted setof.Float64Sorted
}

//msgp:tuple SetOfTuple
type SetOfTuple struct {
	String        setof.String
	StringSorted  setof.StringSorted
	Int           setof.Int
	IntSorted     setof.IntSorted
	Int64         setof.Int64
	Int64Sorted   setof.Int64Sorted
	Int32         setof.Int32
	Int32Sorted   setof.Int32Sorted
	Int16         setof.Int16
	Int16Sorted   setof.Int16Sorted
	Int8          setof.Int8
	Int8Sorted    setof.Int8Sorted
	Byte          setof.Byte
	ByteSorted    setof.ByteSorted
	Uint          setof.Uint
	UintSorted    setof.UintSorted
	Uint64        setof.Uint64
	Uint64Sorted  setof.Uint64Sorted
	Uint32        setof.Uint32
	Uint32Sorted  setof.Uint32Sorted
	Uint16        setof.Uint16
	Uint16Sorted  setof.Uint16Sorted
	Uint8         setof.Uint8
	Uint8Sorted   setof.Uint8Sorted
	Float32       setof.Float32
	Float32Sorted setof.Float32Sorted
	Float64       setof.Float64
	Float64Sorted setof.Float64Sorted
}

var setofFilled = SetOf{
	String:        map[string]struct{}{"a": {}, "b": {}, "c": {}},
	StringSorted:  map[string]struct{}{"d": {}, "e": {}, "f": {}},
	Int:           map[int]struct{}{1: {}, 2: {}, 3: {}},
	IntSorted:     map[int]struct{}{4: {}, 5: {}, 6: {}},
	Int64:         map[int64]struct{}{1: {}, 2: {}, 3: {}},
	Int64Sorted:   map[int64]struct{}{4: {}, 5: {}, 6: {}},
	Int32:         map[int32]struct{}{1: {}, 2: {}, 3: {}},
	Int32Sorted:   map[int32]struct{}{4: {}, 5: {}, 6: {}},
	Int16:         map[int16]struct{}{1: {}, 2: {}, 3: {}},
	Int16Sorted:   map[int16]struct{}{4: {}, 5: {}, 6: {}},
	Int8:          map[int8]struct{}{1: {}, 2: {}, 3: {}},
	Int8Sorted:    map[int8]struct{}{4: {}, 5: {}, 6: {}},
	Byte:          map[byte]struct{}{1: {}, 2: {}, 3: {}},
	ByteSorted:    map[byte]struct{}{4: {}, 5: {}, 6: {}},
	Uint:          map[uint]struct{}{1: {}, 2: {}, 3: {}},
	UintSorted:    map[uint]struct{}{4: {}, 5: {}, 6: {}},
	Uint64:        map[uint64]struct{}{1: {}, 2: {}, 3: {}},
	Uint64Sorted:  map[uint64]struct{}{4: {}, 5: {}, 6: {}},
	Uint32:        map[uint32]struct{}{1: {}, 2: {}, 3: {}},
	Uint32Sorted:  map[uint32]struct{}{4: {}, 5: {}, 6: {}},
	Uint16:        map[uint16]struct{}{1: {}, 2: {}, 3: {}},
	Uint16Sorted:  map[uint16]struct{}{4: {}, 5: {}, 6: {}},
	Uint8:         map[uint8]struct{}{1: {}, 2: {}, 3: {}},
	Uint8Sorted:   map[uint8]struct{}{4: {}, 5: {}, 6: {}},
	Float32:       map[float32]struct{}{1: {}, 2: {}, 3: {}},
	Float32Sorted: map[float32]struct{}{4: {}, 5: {}, 6: {}},
	Float64:       map[float64]struct{}{1: {}, 2: {}, 3: {}},
	Float64Sorted: map[float64]struct{}{4: {}, 5: {}, 6: {}},
}

var setofZeroLen = SetOf{
	String:        map[string]struct{}{},
	StringSorted:  map[string]struct{}{},
	Int:           map[int]struct{}{},
	IntSorted:     map[int]struct{}{},
	Int64:         map[int64]struct{}{},
	Int64Sorted:   map[int64]struct{}{},
	Int32:         map[int32]struct{}{},
	Int32Sorted:   map[int32]struct{}{},
	Int16:         map[int16]struct{}{},
	Int16Sorted:   map[int16]struct{}{},
	Int8:          map[int8]struct{}{},
	Int8Sorted:    map[int8]struct{}{},
	Byte:          map[byte]struct{}{},
	ByteSorted:    map[byte]struct{}{},
	Uint:          map[uint]struct{}{},
	UintSorted:    map[uint]struct{}{},
	Uint64:        map[uint64]struct{}{},
	Uint64Sorted:  map[uint64]struct{}{},
	Uint32:        map[uint32]struct{}{},
	Uint32Sorted:  map[uint32]struct{}{},
	Uint16:        map[uint16]struct{}{},
	Uint16Sorted:  map[uint16]struct{}{},
	Uint8:         map[uint8]struct{}{},
	Uint8Sorted:   map[uint8]struct{}{},
	Float32:       map[float32]struct{}{},
	Float32Sorted: map[float32]struct{}{},
	Float64:       map[float64]struct{}{},
	Float64Sorted: map[float64]struct{}{},
}

var setofSortedFilled = SetOf{
	StringSorted:  map[string]struct{}{"d": {}, "e": {}, "f": {}},
	IntSorted:     map[int]struct{}{4: {}, 5: {}, 6: {}},
	Int64Sorted:   map[int64]struct{}{4: {}, 5: {}, 6: {}},
	Int32Sorted:   map[int32]struct{}{4: {}, 5: {}, 6: {}},
	Int16Sorted:   map[int16]struct{}{4: {}, 5: {}, 6: {}},
	Int8Sorted:    map[int8]struct{}{4: {}, 5: {}, 6: {}},
	ByteSorted:    map[byte]struct{}{4: {}, 5: {}, 6: {}},
	UintSorted:    map[uint]struct{}{4: {}, 5: {}, 6: {}},
	Uint64Sorted:  map[uint64]struct{}{4: {}, 5: {}, 6: {}},
	Uint32Sorted:  map[uint32]struct{}{4: {}, 5: {}, 6: {}},
	Uint16Sorted:  map[uint16]struct{}{4: {}, 5: {}, 6: {}},
	Uint8Sorted:   map[uint8]struct{}{4: {}, 5: {}, 6: {}},
	Float32Sorted: map[float32]struct{}{4: {}, 5: {}, 6: {}},
	Float64Sorted: map[float64]struct{}{4: {}, 5: {}, 6: {}},
}
