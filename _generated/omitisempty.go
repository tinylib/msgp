package _generated

import "time"

//go:generate msgp

// check some specific cases for omitisempty

type OmitIsEmpty0 struct {
	ABool      bool               `msg:"abool,omitisempty"`      // unaliased primitives should act exactly as omitempty
	ABoolPtr   *bool              `msg:"aboolptr,omitisempty"`   // pointers to primitives will do unnecessary interface check for now but function correctly
	AStruct    OmitIsEmptyA       `msg:"astruct,omitempty"`      // leave this one omitempty
	BStruct    OmitIsEmptyA       `msg:"bstruct,omitisempty"`    // and compare to this
	AStructPtr *OmitIsEmptyA      `msg:"astructptr,omitempty"`   // a pointer case omitempty
	BStructPtr *OmitIsEmptyA      `msg:"bstructptr,omitisempty"` // a pointer case omitisempty
	AExt       OmitIsEmptyExt     `msg:"aext,omitisempty"`       // external type case
	AExtPtr    *OmitIsEmptyExtPtr `msg:"aextptr,omitisempty"`    // external type pointer case

	// check a bunch of other types just to make sure none of the output breaks on compilation

	AInt              int               `msg:"aint,omitisempty"`
	AInt8             int8              `msg:"aint8,omitisempty"`
	AInt16            int16             `msg:"aint16,omitisempty"`
	AInt32            int32             `msg:"aint32,omitisempty"`
	AInt64            int64             `msg:"aint64,omitisempty"`
	AUint             uint              `msg:"auint,omitisempty"`
	AUint8            uint8             `msg:"auint8,omitisempty"`
	AUint16           uint16            `msg:"auint16,omitisempty"`
	AUint32           uint32            `msg:"auint32,omitisempty"`
	AUint64           uint64            `msg:"auint64,omitisempty"`
	AFloat32          float32           `msg:"afloat32,omitisempty"`
	AFloat64          float64           `msg:"afloat64,omitisempty"`
	AComplex64        complex64         `msg:"acomplex64,omitisempty"`
	AComplex128       complex128        `msg:"acomplex128,omitisempty"`
	ANamedBool        bool              `msg:"anamedbool,omitisempty"`
	ANamedInt         int               `msg:"anamedint,omitisempty"`
	ANamedFloat64     float64           `msg:"anamedfloat64,omitisempty"`
	AMapStrStr        map[string]string `msg:"amapstrstr,omitisempty"`
	APtrNamedStr      *NamedString      `msg:"aptrnamedstr,omitisempty"`
	AString           string            `msg:"astring,omitisempty"`
	ANamedString      string            `msg:"anamedstring,omitisempty"`
	AByteSlice        []byte            `msg:"abyteslice,omitisempty"`
	ASliceString      []string          `msg:"aslicestring,omitisempty"`
	ASliceNamedString []NamedString     `msg:"aslicenamedstring,omitisempty"`
	ANamedStruct      NamedStruct       `msg:"anamedstruct,omitisempty"`
	APtrNamedStruct   *NamedStruct      `msg:"aptrnamedstruct,omitisempty"`
	AUnnamedStruct    struct {
		A string `msg:"a,omitisempty"`
	} `msg:"aunnamedstruct,omitisempty"` // omitisempty not supported on unnamed struct
	EmbeddableStruct  `msg:",flatten,omitisempty"`          // embed flat
	EmbeddableStruct2 `msg:"embeddablestruct2,omitisempty"` // embed non-flat
	AArrayInt         [5]int                                `msg:"aarrayint,omitisempty"` // not supported
	ATime             time.Time                             `msg:"atime,omitisempty"`
}

type OmitIsEmptyA struct {
	A int    `msg:"a,omitempty"`
	B int    `msg:"b,omitisempty"`
	C string `msg:"c,omitisempty"`
}
