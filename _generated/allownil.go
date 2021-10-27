package _generated

import "time"

//go:generate msgp

type AllowNil0 struct {
	ABool       bool       `msg:"abool,allownil"`
	AInt        int        `msg:"aint,allownil"`
	AInt8       int8       `msg:"aint8,allownil"`
	AInt16      int16      `msg:"aint16,allownil"`
	AInt32      int32      `msg:"aint32,allownil"`
	AInt64      int64      `msg:"aint64,allownil"`
	AUint       uint       `msg:"auint,allownil"`
	AUint8      uint8      `msg:"auint8,allownil"`
	AUint16     uint16     `msg:"auint16,allownil"`
	AUint32     uint32     `msg:"auint32,allownil"`
	AUint64     uint64     `msg:"auint64,allownil"`
	AFloat32    float32    `msg:"afloat32,allownil"`
	AFloat64    float64    `msg:"afloat64,allownil"`
	AComplex64  complex64  `msg:"acomplex64,allownil"`
	AComplex128 complex128 `msg:"acomplex128,allownil"`

	ANamedBool    bool    `msg:"anamedbool,allownil"`
	ANamedInt     int     `msg:"anamedint,allownil"`
	ANamedFloat64 float64 `msg:"anamedfloat64,allownil"`

	AMapStrStr map[string]string `msg:"amapstrstr,allownil"`

	APtrNamedStr *NamedString `msg:"aptrnamedstr,allownil"`

	AString      string `msg:"astring,allownil"`
	ANamedString string `msg:"anamedstring,allownil"`
	AByteSlice   []byte `msg:"abyteslice,allownil"`

	ASliceString      []string      `msg:"aslicestring,allownil"`
	ASliceNamedString []NamedString `msg:"aslicenamedstring,allownil"`

	ANamedStruct    NamedStructAN  `msg:"anamedstruct,allownil"`
	APtrNamedStruct *NamedStructAN `msg:"aptrnamedstruct,allownil"`

	AUnnamedStruct struct {
		A string `msg:"a,allownil"`
	} `msg:"aunnamedstruct,allownil"` // allownil not supported on unnamed struct

	EmbeddableStructAN `msg:",flatten,allownil"` // embed flat

	EmbeddableStruct2AN `msg:"embeddablestruct2,allownil"` // embed non-flat

	AArrayInt [5]int `msg:"aarrayint,allownil"` // not supported

	ATime time.Time `msg:"atime,allownil"`
}

type AllowNil1 struct {
	ABool       []bool       `msg:"abool,allownil"`
	AInt        []int        `msg:"aint,allownil"`
	AInt8       []int8       `msg:"aint8,allownil"`
	AInt16      []int16      `msg:"aint16,allownil"`
	AInt32      []int32      `msg:"aint32,allownil"`
	AInt64      []int64      `msg:"aint64,allownil"`
	AUint       []uint       `msg:"auint,allownil"`
	AUint8      []uint8      `msg:"auint8,allownil"`
	AUint16     []uint16     `msg:"auint16,allownil"`
	AUint32     []uint32     `msg:"auint32,allownil"`
	AUint64     []uint64     `msg:"auint64,allownil"`
	AFloat32    []float32    `msg:"afloat32,allownil"`
	AFloat64    []float64    `msg:"afloat64,allownil"`
	AComplex64  []complex64  `msg:"acomplex64,allownil"`
	AComplex128 []complex128 `msg:"acomplex128,allownil"`

	ANamedBool    []bool    `msg:"anamedbool,allownil"`
	ANamedInt     []int     `msg:"anamedint,allownil"`
	ANamedFloat64 []float64 `msg:"anamedfloat64,allownil"`

	AMapStrStr map[string]string `msg:"amapstrstr,allownil"`

	APtrNamedStr *NamedString `msg:"aptrnamedstr,allownil"`

	AString      []string `msg:"astring,allownil"`
	ANamedString []string `msg:"anamedstring,allownil"`
	AByteSlice   []byte   `msg:"abyteslice,allownil"`

	ASliceString      []string      `msg:"aslicestring,allownil"`
	ASliceNamedString []NamedString `msg:"aslicenamedstring,allownil"`

	ANamedStruct    NamedStructAN  `msg:"anamedstruct,allownil"`
	APtrNamedStruct *NamedStructAN `msg:"aptrnamedstruct,allownil"`

	AUnnamedStruct struct {
		A []string `msg:"a,allownil"`
	} `msg:"aunnamedstruct,allownil"`

	*EmbeddableStructAN `msg:",flatten,allownil"` // embed flat

	*EmbeddableStruct2AN `msg:"embeddablestruct2,allownil"` // embed non-flat

	AArrayInt [5]int `msg:"aarrayint,allownil"` // not supported

	ATime *time.Time `msg:"atime,allownil"`
}

type EmbeddableStructAN struct {
	SomeEmbed []string `msg:"someembed,allownil"`
}

type EmbeddableStruct2AN struct {
	SomeEmbed2 []string `msg:"someembed2,allownil"`
}

type NamedStructAN struct {
	A []string `msg:"a,allownil"`
	B []string `msg:"b,allownil"`
}

type AllowNilHalfFull struct {
	Field00 []string `msg:"field00,allownil"`
	Field01 []string `msg:"field01"`
	Field02 []string `msg:"field02,allownil"`
	Field03 []string `msg:"field03"`
}

type AllowNilLotsOFields struct {
	Field00 []string `msg:"field00,allownil"`
	Field01 []string `msg:"field01,allownil"`
	Field02 []string `msg:"field02,allownil"`
	Field03 []string `msg:"field03,allownil"`
	Field04 []string `msg:"field04,allownil"`
	Field05 []string `msg:"field05,allownil"`
	Field06 []string `msg:"field06,allownil"`
	Field07 []string `msg:"field07,allownil"`
	Field08 []string `msg:"field08,allownil"`
	Field09 []string `msg:"field09,allownil"`
	Field10 []string `msg:"field10,allownil"`
	Field11 []string `msg:"field11,allownil"`
	Field12 []string `msg:"field12,allownil"`
	Field13 []string `msg:"field13,allownil"`
	Field14 []string `msg:"field14,allownil"`
	Field15 []string `msg:"field15,allownil"`
	Field16 []string `msg:"field16,allownil"`
	Field17 []string `msg:"field17,allownil"`
	Field18 []string `msg:"field18,allownil"`
	Field19 []string `msg:"field19,allownil"`
	Field20 []string `msg:"field20,allownil"`
	Field21 []string `msg:"field21,allownil"`
	Field22 []string `msg:"field22,allownil"`
	Field23 []string `msg:"field23,allownil"`
	Field24 []string `msg:"field24,allownil"`
	Field25 []string `msg:"field25,allownil"`
	Field26 []string `msg:"field26,allownil"`
	Field27 []string `msg:"field27,allownil"`
	Field28 []string `msg:"field28,allownil"`
	Field29 []string `msg:"field29,allownil"`
	Field30 []string `msg:"field30,allownil"`
	Field31 []string `msg:"field31,allownil"`
	Field32 []string `msg:"field32,allownil"`
	Field33 []string `msg:"field33,allownil"`
	Field34 []string `msg:"field34,allownil"`
	Field35 []string `msg:"field35,allownil"`
	Field36 []string `msg:"field36,allownil"`
	Field37 []string `msg:"field37,allownil"`
	Field38 []string `msg:"field38,allownil"`
	Field39 []string `msg:"field39,allownil"`
	Field40 []string `msg:"field40,allownil"`
	Field41 []string `msg:"field41,allownil"`
	Field42 []string `msg:"field42,allownil"`
	Field43 []string `msg:"field43,allownil"`
	Field44 []string `msg:"field44,allownil"`
	Field45 []string `msg:"field45,allownil"`
	Field46 []string `msg:"field46,allownil"`
	Field47 []string `msg:"field47,allownil"`
	Field48 []string `msg:"field48,allownil"`
	Field49 []string `msg:"field49,allownil"`
	Field50 []string `msg:"field50,allownil"`
	Field51 []string `msg:"field51,allownil"`
	Field52 []string `msg:"field52,allownil"`
	Field53 []string `msg:"field53,allownil"`
	Field54 []string `msg:"field54,allownil"`
	Field55 []string `msg:"field55,allownil"`
	Field56 []string `msg:"field56,allownil"`
	Field57 []string `msg:"field57,allownil"`
	Field58 []string `msg:"field58,allownil"`
	Field59 []string `msg:"field59,allownil"`
	Field60 []string `msg:"field60,allownil"`
	Field61 []string `msg:"field61,allownil"`
	Field62 []string `msg:"field62,allownil"`
	Field63 []string `msg:"field63,allownil"`
	Field64 []string `msg:"field64,allownil"`
	Field65 []string `msg:"field65,allownil"`
	Field66 []string `msg:"field66,allownil"`
	Field67 []string `msg:"field67,allownil"`
	Field68 []string `msg:"field68,allownil"`
	Field69 []string `msg:"field69,allownil"`
}

type AllowNil10 struct {
	Field00 []string `msg:"field00,allownil"`
	Field01 []string `msg:"field01,allownil"`
	Field02 []string `msg:"field02,allownil"`
	Field03 []string `msg:"field03,allownil"`
	Field04 []string `msg:"field04,allownil"`
	Field05 []string `msg:"field05,allownil"`
	Field06 []string `msg:"field06,allownil"`
	Field07 []string `msg:"field07,allownil"`
	Field08 []string `msg:"field08,allownil"`
	Field09 []string `msg:"field09,allownil"`
}

type NotAllowNil10 struct {
	Field00 []string `msg:"field00"`
	Field01 []string `msg:"field01"`
	Field02 []string `msg:"field02"`
	Field03 []string `msg:"field03"`
	Field04 []string `msg:"field04"`
	Field05 []string `msg:"field05"`
	Field06 []string `msg:"field06"`
	Field07 []string `msg:"field07"`
	Field08 []string `msg:"field08"`
	Field09 []string `msg:"field09"`
}

type AllowNilOmitEmpty struct {
	Field00 []string `msg:"field00,allownil,omitempty"`
	Field01 []string `msg:"field01,allownil"`
}

type AllowNilOmitEmpty2 struct {
	Field00 []string `msg:"field00,allownil,omitempty"`
	Field01 []string `msg:"field01,allownil,omitempty"`
}
