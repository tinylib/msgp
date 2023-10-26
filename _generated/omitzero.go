package _generated

import "time"

//go:generate msgp

// check some specific cases for omitzero

type OmitZero0 struct {
	AStruct    OmitZeroA       `msg:"astruct,omitempty"`    // leave this one omitempty
	BStruct    OmitZeroA       `msg:"bstruct,omitzero"`     // and compare to this
	AStructPtr *OmitZeroA      `msg:"astructptr,omitempty"` // a pointer case omitempty
	BStructPtr *OmitZeroA      `msg:"bstructptr,omitzero"`  // a pointer case omitzero
	AExt       OmitZeroExt     `msg:"aext,omitzero"`        // external type case
	AExtPtr    *OmitZeroExtPtr `msg:"aextptr,omitzero"`     // external type pointer case

	// more
	APtrNamedStr       *NamedStringOZ                     `msg:"aptrnamedstr,omitzero"`
	ANamedStruct       NamedStructOZ                      `msg:"anamedstruct,omitzero"`
	APtrNamedStruct    *NamedStructOZ                     `msg:"aptrnamedstruct,omitzero"`
	EmbeddableStruct   `msg:",flatten,omitzero"`          // embed flat
	EmbeddableStructOZ `msg:"embeddablestruct2,omitzero"` // embed non-flat
	ATime              time.Time                          `msg:"atime,omitzero"`

	OmitZeroTuple OmitZeroTuple `msg:"ozt"` // the inside of a tuple should ignore both omitempty and omitzero
}

type OmitZeroA struct {
	A string        `msg:"a,omitempty"`
	B NamedStringOZ `msg:"b,omitzero"`
	C NamedStringOZ `msg:"c,omitzero"`
}

func (o *OmitZeroA) IsZero() bool {
	if o == nil {
		return true
	}
	return *o == (OmitZeroA{})
}

type NamedStructOZ struct {
	A string `msg:"a,omitempty"`
	B string `msg:"b,omitempty"`
}

func (ns *NamedStructOZ) IsZero() bool {
	if ns == nil {
		return true
	}
	return *ns == (NamedStructOZ{})
}

type NamedStringOZ string

func (ns *NamedStringOZ) IsZero() bool {
	if ns == nil {
		return true
	}
	return *ns == ""
}

type EmbeddableStructOZ struct {
	SomeEmbed string `msg:"someembed2,omitempty"`
}

func (es EmbeddableStructOZ) IsZero() bool { return es == (EmbeddableStructOZ{}) }

type EmbeddableStructOZ2 struct {
	SomeEmbed2 string `msg:"someembed2,omitempty"`
}

func (es EmbeddableStructOZ2) IsZero() bool { return es == (EmbeddableStructOZ2{}) }

//msgp:tuple OmitZeroTuple

// OmitZeroTuple is flagged for tuple output, it should ignore all omitempty and omitzero functionality
// since it's fundamentally incompatible.
type OmitZeroTuple struct {
	FieldA string        `msg:"fielda,omitempty"`
	FieldB NamedStringOZ `msg:"fieldb,omitzero"`
	FieldC NamedStringOZ `msg:"fieldc,omitzero"`
}

type OmitZero1 struct {
	T1 OmitZeroTuple `msg:"t1"`
}
