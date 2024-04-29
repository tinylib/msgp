package _generated

import "time"

// external types to test replace directive

type (
	MapString   map[string]string
	MapAny      map[string]any
	SliceString []String
	SliceInt    []Int
	Array8      [8]byte
	Array16     [16]byte
	Int         int
	Uint        uint
	String      string
	Float32     float32
	Float64     float64
	Time        time.Time
	Duration    time.Duration
	Any         any

	StructA struct {
		StructB StructB
		Int     Int
	}

	StructB struct {
		StructC StructC
		Any     Any
		Array8  Array8
	}

	StructC struct {
		StructD StructD
		Float64 Float32
		Float32 Float64
	}

	StructD struct {
		Time      Time
		Duration  Duration
		MapString MapString
	}

	StructI struct {
		Int  *Int
		Uint *Uint
	}

	StructS struct {
		Slice SliceInt
	}
)
