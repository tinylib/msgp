package _generated

import "time"

//go:generate msgp

type OmitEmpty0 struct {
	ABool       bool       `msg:"abool,omitempty"`
	AInt        int        `msg:"aint,omitempty"`
	AInt8       int8       `msg:"aint8,omitempty"`
	AInt16      int16      `msg:"aint16,omitempty"`
	AInt32      int32      `msg:"aint32,omitempty"`
	AInt64      int64      `msg:"aint64,omitempty"`
	AUint       uint       `msg:"auint,omitempty"`
	AUint8      uint8      `msg:"auint8,omitempty"`
	AUint16     uint16     `msg:"auint16,omitempty"`
	AUint32     uint32     `msg:"auint32,omitempty"`
	AUint64     uint64     `msg:"auint64,omitempty"`
	AFloat32    float32    `msg:"afloat32,omitempty"`
	AFloat64    float64    `msg:"afloat64,omitempty"`
	AComplex64  complex64  `msg:"acomplex64,omitempty"`
	AComplex128 complex128 `msg:"acomplex128,omitempty"`

	ANamedBool    bool    `msg:"anamedbool,omitempty"`
	ANamedInt     int     `msg:"anamedint,omitempty"`
	ANamedFloat64 float64 `msg:"anamedfloat64,omitempty"`

	AMapStrStr map[string]string `msg:"amapstrstr,omitempty"`

	APtrNamedStr *NamedString `msg:"aptrnamedstr,omitempty"`

	AString      string `msg:"astring,omitempty"`
	ANamedString string `msg:"anamedstring,omitempty"`
	AByteSlice   []byte `msg:"abyteslice,omitempty"`

	ASliceString      []string      `msg:"aslicestring,omitempty"`
	ASliceNamedString []NamedString `msg:"aslicenamedstring,omitempty"`

	ANamedStruct    NamedStruct  `msg:"anamedstruct,omitempty"`
	APtrNamedStruct *NamedStruct `msg:"aptrnamedstruct,omitempty"`

	AUnnamedStruct struct {
		A string `msg:"a,omitempty"`
	} `msg:"aunnamedstruct,omitempty"` // omitempty not supported on unnamed struct

	EmbeddableStruct `msg:",flatten,omitempty"` // embed flat

	EmbeddableStruct2 `msg:"embeddablestruct2,omitempty"` // embed non-flat

	AArrayInt [5]int `msg:"aarrayint,omitempty"` // not supported

	ATime time.Time `msg:"atime,omitempty"`
}

type (
	NamedBool    bool
	NamedInt     int
	NamedFloat64 float64
	NamedString  string
)

type EmbeddableStruct struct {
	SomeEmbed string `msg:"someembed,omitempty"`
}

type EmbeddableStruct2 struct {
	SomeEmbed2 string `msg:"someembed2,omitempty"`
}

type NamedStruct struct {
	A string `msg:"a,omitempty"`
	B string `msg:"b,omitempty"`
}

type OmitEmptyHalfFull struct {
	Field00 string `msg:"field00,omitempty"`
	Field01 string `msg:"field01"`
	Field02 string `msg:"field02,omitempty"`
	Field03 string `msg:"field03"`
}

type OmitEmptyLotsOFields struct {
	Field00 string `msg:"field00,omitempty"`
	Field01 string `msg:"field01,omitempty"`
	Field02 string `msg:"field02,omitempty"`
	Field03 string `msg:"field03,omitempty"`
	Field04 string `msg:"field04,omitempty"`
	Field05 string `msg:"field05,omitempty"`
	Field06 string `msg:"field06,omitempty"`
	Field07 string `msg:"field07,omitempty"`
	Field08 string `msg:"field08,omitempty"`
	Field09 string `msg:"field09,omitempty"`
	Field10 string `msg:"field10,omitempty"`
	Field11 string `msg:"field11,omitempty"`
	Field12 string `msg:"field12,omitempty"`
	Field13 string `msg:"field13,omitempty"`
	Field14 string `msg:"field14,omitempty"`
	Field15 string `msg:"field15,omitempty"`
	Field16 string `msg:"field16,omitempty"`
	Field17 string `msg:"field17,omitempty"`
	Field18 string `msg:"field18,omitempty"`
	Field19 string `msg:"field19,omitempty"`
	Field20 string `msg:"field20,omitempty"`
	Field21 string `msg:"field21,omitempty"`
	Field22 string `msg:"field22,omitempty"`
	Field23 string `msg:"field23,omitempty"`
	Field24 string `msg:"field24,omitempty"`
	Field25 string `msg:"field25,omitempty"`
	Field26 string `msg:"field26,omitempty"`
	Field27 string `msg:"field27,omitempty"`
	Field28 string `msg:"field28,omitempty"`
	Field29 string `msg:"field29,omitempty"`
	Field30 string `msg:"field30,omitempty"`
	Field31 string `msg:"field31,omitempty"`
	Field32 string `msg:"field32,omitempty"`
	Field33 string `msg:"field33,omitempty"`
	Field34 string `msg:"field34,omitempty"`
	Field35 string `msg:"field35,omitempty"`
	Field36 string `msg:"field36,omitempty"`
	Field37 string `msg:"field37,omitempty"`
	Field38 string `msg:"field38,omitempty"`
	Field39 string `msg:"field39,omitempty"`
	Field40 string `msg:"field40,omitempty"`
	Field41 string `msg:"field41,omitempty"`
	Field42 string `msg:"field42,omitempty"`
	Field43 string `msg:"field43,omitempty"`
	Field44 string `msg:"field44,omitempty"`
	Field45 string `msg:"field45,omitempty"`
	Field46 string `msg:"field46,omitempty"`
	Field47 string `msg:"field47,omitempty"`
	Field48 string `msg:"field48,omitempty"`
	Field49 string `msg:"field49,omitempty"`
	Field50 string `msg:"field50,omitempty"`
	Field51 string `msg:"field51,omitempty"`
	Field52 string `msg:"field52,omitempty"`
	Field53 string `msg:"field53,omitempty"`
	Field54 string `msg:"field54,omitempty"`
	Field55 string `msg:"field55,omitempty"`
	Field56 string `msg:"field56,omitempty"`
	Field57 string `msg:"field57,omitempty"`
	Field58 string `msg:"field58,omitempty"`
	Field59 string `msg:"field59,omitempty"`
	Field60 string `msg:"field60,omitempty"`
	Field61 string `msg:"field61,omitempty"`
	Field62 string `msg:"field62,omitempty"`
	Field63 string `msg:"field63,omitempty"`
	Field64 string `msg:"field64,omitempty"`
	Field65 string `msg:"field65,omitempty"`
	Field66 string `msg:"field66,omitempty"`
	Field67 string `msg:"field67,omitempty"`
	Field68 string `msg:"field68,omitempty"`
	Field69 string `msg:"field69,omitempty"`
}

type OmitEmpty10 struct {
	Field00 string `msg:"field00,omitempty"`
	Field01 string `msg:"field01,omitempty"`
	Field02 string `msg:"field02,omitempty"`
	Field03 string `msg:"field03,omitempty"`
	Field04 string `msg:"field04,omitempty"`
	Field05 string `msg:"field05,omitempty"`
	Field06 string `msg:"field06,omitempty"`
	Field07 string `msg:"field07,omitempty"`
	Field08 string `msg:"field08,omitempty"`
	Field09 string `msg:"field09,omitempty"`
}

type NotOmitEmpty10 struct {
	Field00 string `msg:"field00"`
	Field01 string `msg:"field01"`
	Field02 string `msg:"field02"`
	Field03 string `msg:"field03"`
	Field04 string `msg:"field04"`
	Field05 string `msg:"field05"`
	Field06 string `msg:"field06"`
	Field07 string `msg:"field07"`
	Field08 string `msg:"field08"`
	Field09 string `msg:"field09"`
}
