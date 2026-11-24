package _generated

//go:generate msgp

//msgp:replace Issue465Str with:[]byte
//msgp:replace Issue465Ext with:[]byte
//msgp:replace issue465Unexported with:[]byte
//msgp:shim Issue465Struct as:[]byte using:issue465StructToBytes/issue465StructFromBytes
//msgp:shim Issue465Bytes as:[]byte using:issue465BytesToBytes/issue465BytesFromBytes

// Not nilable, so allownil does not apply.
type Issue465Str string

// Unexported declarations are dropped before parsing without -unexported,
// but must still be found to be not nilable.
type issue465Unexported string

type Issue465Struct struct{ A string }

func issue465StructToBytes(v Issue465Struct) []byte   { return []byte(v.A) }
func issue465StructFromBytes(b []byte) Issue465Struct { return Issue465Struct{A: string(b)} }

// Nilable, so allownil still applies.
type Issue465Bytes []byte

func issue465BytesToBytes(v Issue465Bytes) []byte   { return v }
func issue465BytesFromBytes(b []byte) Issue465Bytes { return b }

type Issue465 struct {
	Str     Issue465Str        `msg:",allownil"`
	PStr    *Issue465Str       `msg:",allownil"`
	Struct  Issue465Struct     `msg:",allownil"`
	PStruct *Issue465Struct    `msg:",allownil"`
	Bytes   Issue465Bytes      `msg:",allownil"`
	Ext     Issue465Ext        `msg:",allownil"`
	Unexp   issue465Unexported `msg:",allownil"`
}
