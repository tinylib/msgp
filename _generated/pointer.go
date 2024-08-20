package _generated

import (
	"fmt"
	"time"

	"github.com/tinylib/msgp/msgp"
)

//go:generate msgp $GOFILE$

// Generate only pointer receivers:

//msgp:pointer

var mustNoInterf = []interface{}{
	Pointer0{},
	NamedBoolPointer(true),
	NamedIntPointer(0),
	NamedFloat64Pointer(0),
	NamedStringPointer(""),
	NamedMapStructPointer(nil),
	NamedMapStructPointer2(nil),
	NamedMapStringPointer(nil),
	NamedMapStringPointer2(nil),
	EmbeddableStructPointer{},
	EmbeddableStruct2Pointer{},
	PointerHalfFull{},
	PointerNoName{},
}

var mustHaveInterf = []interface{}{
	&Pointer0{},
	mustPtr(NamedBoolPointer(true)),
	mustPtr(NamedIntPointer(0)),
	mustPtr(NamedFloat64Pointer(0)),
	mustPtr(NamedStringPointer("")),
	mustPtr(NamedMapStructPointer(nil)),
	mustPtr(NamedMapStructPointer2(nil)),
	mustPtr(NamedMapStringPointer(nil)),
	mustPtr(NamedMapStringPointer2(nil)),
	&EmbeddableStructPointer{},
	&EmbeddableStruct2Pointer{},
	&PointerHalfFull{},
	&PointerNoName{},
}

func mustPtr[T any](v T) *T {
	return &v
}

func init() {
	for _, v := range mustNoInterf {
		if _, ok := v.(msgp.Marshaler); ok {
			panic(fmt.Sprintf("type %T supports interface", v))
		}
		if _, ok := v.(msgp.Encodable); ok {
			panic(fmt.Sprintf("type %T supports interface", v))
		}
	}
	for _, v := range mustHaveInterf {
		if _, ok := v.(msgp.Marshaler); !ok {
			panic(fmt.Sprintf("type %T does not support interface", v))
		}
		if _, ok := v.(msgp.Encodable); !ok {
			panic(fmt.Sprintf("type %T does not support interface", v))
		}
	}
}

type Pointer0 struct {
	ABool       bool       `msg:"abool"`
	AInt        int        `msg:"aint"`
	AInt8       int8       `msg:"aint8"`
	AInt16      int16      `msg:"aint16"`
	AInt32      int32      `msg:"aint32"`
	AInt64      int64      `msg:"aint64"`
	AUint       uint       `msg:"auint"`
	AUint8      uint8      `msg:"auint8"`
	AUint16     uint16     `msg:"auint16"`
	AUint32     uint32     `msg:"auint32"`
	AUint64     uint64     `msg:"auint64"`
	AFloat32    float32    `msg:"afloat32"`
	AFloat64    float64    `msg:"afloat64"`
	AComplex64  complex64  `msg:"acomplex64"`
	AComplex128 complex128 `msg:"acomplex128"`

	ANamedBool    bool    `msg:"anamedbool"`
	ANamedInt     int     `msg:"anamedint"`
	ANamedFloat64 float64 `msg:"anamedfloat64"`

	AMapStrStr map[string]string `msg:"amapstrstr"`

	APtrNamedStr *NamedString `msg:"aptrnamedstr"`

	AString      string `msg:"astring"`
	ANamedString string `msg:"anamedstring"`
	AByteSlice   []byte `msg:"abyteslice"`

	ASliceString      []string      `msg:"aslicestring"`
	ASliceNamedString []NamedString `msg:"aslicenamedstring"`

	ANamedStruct    NamedStruct  `msg:"anamedstruct"`
	APtrNamedStruct *NamedStruct `msg:"aptrnamedstruct"`

	AUnnamedStruct struct {
		A string `msg:"a"`
	} `msg:"aunnamedstruct"` // omitempty not supported on unnamed struct

	EmbeddableStruct `msg:",flatten"` // embed flat

	EmbeddableStruct2 `msg:"embeddablestruct2"` // embed non-flat

	AArrayInt [5]int `msg:"aarrayint"` // not supported

	ATime time.Time `msg:"atime"`
}

type (
	NamedBoolPointer       bool
	NamedIntPointer        int
	NamedFloat64Pointer    float64
	NamedStringPointer     string
	NamedMapStructPointer  map[string]Pointer0
	NamedMapStructPointer2 map[string]*Pointer0
	NamedMapStringPointer  map[string]NamedStringPointer
	NamedMapStringPointer2 map[string]*NamedStringPointer
)

type EmbeddableStructPointer struct {
	SomeEmbed string `msg:"someembed"`
}

type EmbeddableStruct2Pointer struct {
	SomeEmbed2 string `msg:"someembed2"`
}

type NamedStructPointer struct {
	A string `msg:"a"`
	B string `msg:"b"`
}

type PointerHalfFull struct {
	Field00 string `msg:"field00"`
	Field01 string `msg:"field01"`
	Field02 string `msg:"field02"`
	Field03 string `msg:"field03"`
}

type PointerNoName struct {
	ABool       bool       `msg:""`
	AInt        int        `msg:""`
	AInt8       int8       `msg:""`
	AInt16      int16      `msg:""`
	AInt32      int32      `msg:""`
	AInt64      int64      `msg:""`
	AUint       uint       `msg:""`
	AUint8      uint8      `msg:""`
	AUint16     uint16     `msg:""`
	AUint32     uint32     `msg:""`
	AUint64     uint64     `msg:""`
	AFloat32    float32    `msg:""`
	AFloat64    float64    `msg:""`
	AComplex64  complex64  `msg:""`
	AComplex128 complex128 `msg:""`

	ANamedBool    bool    `msg:""`
	ANamedInt     int     `msg:""`
	ANamedFloat64 float64 `msg:""`

	AMapStrF       map[string]NamedFloat64Pointer `msg:""`
	AMapStrStruct  map[string]PointerHalfFull     `msg:""`
	AMapStrStruct2 map[string]*PointerHalfFull    `msg:""`

	APtrNamedStr *NamedStringPointer `msg:""`

	AString    string `msg:""`
	AByteSlice []byte `msg:""`

	ASliceString      []string             `msg:""`
	ASliceNamedString []NamedStringPointer `msg:""`

	ANamedStruct    NamedStructPointer  `msg:""`
	APtrNamedStruct *NamedStructPointer `msg:""`

	AUnnamedStruct struct {
		A string `msg:""`
	} `msg:""` // omitempty not supported on unnamed struct

	EmbeddableStructPointer `msg:",flatten"` // embed flat

	EmbeddableStruct2Pointer `msg:""` // embed non-flat

	AArrayInt [5]int `msg:""` // not supported

	ATime time.Time     `msg:""`
	ADur  time.Duration `msg:""`
}
