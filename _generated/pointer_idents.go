package _generated

//go:generate msgp

//msgp:pointer
//msgp:binappend PointerIdentsVA

type PointerIdentsVA struct{ V string }

func (t PointerIdentsVA) AppendBinary(b []byte) ([]byte, error) { return append(b, t.V...), nil }
func (t *PointerIdentsVA) UnmarshalBinary(b []byte) error       { t.V = string(b); return nil }

type PointerIdents struct{ C []PointerIdentsVA }
