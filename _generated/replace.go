package _generated

import "encoding/json"

//go:generate msgp
//msgp:replace Any with:any
//msgp:replace MapString with:CompatibleMapString
//msgp:replace MapAny with:map[string]any
//msgp:replace SliceString with:[]string
//msgp:replace SliceInt with:CompatibleSliceInt
//msgp:replace Array8 with:CompatibleArray8
//msgp:replace Array16 with:[16]byte
//msgp:replace String with:string
//msgp:replace Int with:CompatibleInt
//msgp:replace Uint with:uint
//msgp:replace Float32 with:CompatibleFloat32
//msgp:replace Float64 with:CompatibleFloat64
//msgp:replace Time with:time.Time
//msgp:replace Duration with:time.Duration
//msgp:replace StructA with:CompatibleStructA
//msgp:replace StructB with:CompatibleStructB
//msgp:replace StructC with:CompatibleStructC
//msgp:replace StructD with:CompatibleStructD
//msgp:replace StructI with:CompatibleStructI
//msgp:replace StructS with:CompatibleStructS

type (
	CompatibleMapString map[string]string
	CompatibleArray8    [8]byte
	CompatibleInt       int
	CompatibleFloat32   float32
	CompatibleFloat64   float64
	CompatibleSliceInt  []Int

	// Doesn't work
	// CompatibleTime time.Time

	CompatibleStructA struct {
		StructB StructB
		Int     Int
	}

	CompatibleStructB struct {
		StructC StructC
		Any     Any
		Array8  Array8
	}

	CompatibleStructC struct {
		StructD StructD
		Float64 Float32
		Float32 Float64
	}

	CompatibleStructD struct {
		Time      Time
		Duration  Duration
		MapString MapString
	}

	CompatibleStructI struct {
		Int  *Int
		Uint *Uint
	}

	CompatibleStructS struct {
		Slice SliceInt
	}

	Dummy struct {
		StructA StructA
		StructI StructI
		StructS StructS
		Array16 Array16
		Uint    Uint
		String  String
	}
)

//msgp:replace json.Number with:string

type NumberJSONSampleReplace struct {
	Single json.Number
	Array  []json.Number
	Map    map[string]json.Number
}
