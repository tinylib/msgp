package _generated

import (
	"encoding/json"
	"time"
)

//go:generate msgp

//msgp:clearomitted

// check some specific cases for omitzero

type ClearOmitted0 struct {
	AStruct    ClearOmittedA  `msg:"astruct,omitempty"`    // leave this one omitempty
	BStruct    ClearOmittedA  `msg:"bstruct,omitzero"`     // and compare to this
	AStructPtr *ClearOmittedA `msg:"astructptr,omitempty"` // a pointer case omitempty
	BStructPtr *ClearOmittedA `msg:"bstructptr,omitzero"`  // a pointer case omitzero
	AExt       OmitZeroExt    `msg:"aext,omitzero"`        // external type case

	// more
	APtrNamedStr        *NamedStringCO                     `msg:"aptrnamedstr,omitzero"`
	ANamedStruct        NamedStructCO                      `msg:"anamedstruct,omitzero"`
	APtrNamedStruct     *NamedStructCO                     `msg:"aptrnamedstruct,omitzero"`
	EmbeddableStructCO  `msg:",flatten,omitzero"`          // embed flat
	EmbeddableStructCO2 `msg:"embeddablestruct2,omitzero"` // embed non-flat
	ATime               time.Time                          `msg:"atime,omitzero"`
	ASlice              []int                              `msg:"aslice,omitempty"`
	AMap                map[string]int                     `msg:"amap,omitempty"`
	ABin                []byte                             `msg:"abin,omitempty"`
	AInt                int                                `msg:"aint,omitempty"`
	AString             string                             `msg:"atring,omitempty"`
	Adur                time.Duration                      `msg:"adur,omitempty"`
	AJSON               json.Number                        `msg:"ajson,omitempty"`
	AnAny               any                                `msg:"anany,omitempty"`

	ClearOmittedTuple ClearOmittedTuple `msg:"ozt"` // the inside of a tuple should ignore both omitempty and omitzero
}

type ClearOmittedA struct {
	A string        `msg:"a,omitempty"`
	B NamedStringCO `msg:"b,omitzero"`
	C NamedStringCO `msg:"c,omitzero"`
}

func (o *ClearOmittedA) IsZero() bool {
	if o == nil {
		return true
	}
	return *o == (ClearOmittedA{})
}

type NamedStructCO struct {
	A string `msg:"a,omitempty"`
	B string `msg:"b,omitempty"`
}

func (ns *NamedStructCO) IsZero() bool {
	if ns == nil {
		return true
	}
	return *ns == (NamedStructCO{})
}

type NamedStringCO string

func (ns *NamedStringCO) IsZero() bool {
	if ns == nil {
		return true
	}
	return *ns == ""
}

type EmbeddableStructCO struct {
	SomeEmbed string `msg:"someembed2,omitempty"`
}

func (es EmbeddableStructCO) IsZero() bool { return es == (EmbeddableStructCO{}) }

type EmbeddableStructCO2 struct {
	SomeEmbed2 string `msg:"someembed2,omitempty"`
}

func (es EmbeddableStructCO2) IsZero() bool { return es == (EmbeddableStructCO2{}) }

//msgp:tuple ClearOmittedTuple

// ClearOmittedTuple is flagged for tuple output, it should ignore all omitempty and omitzero functionality
// since it's fundamentally incompatible.
type ClearOmittedTuple struct {
	FieldA string        `msg:"fielda,omitempty"`
	FieldB NamedStringCO `msg:"fieldb,omitzero"`
	FieldC NamedStringCO `msg:"fieldc,omitzero"`
}

type ClearOmitted1 struct {
	T1 ClearOmittedTuple `msg:"t1"`
}
