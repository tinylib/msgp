package gen

import (
	"fmt"
)

type encodeGen struct {
	p printer
}

func (e *encodeGen) writeAndCheck(typ string, argfmt string, arg interface{}) {
	e.p.printf("\nerr = en.Write%s(%s)", typ, fmt.Sprintf(argfmt, arg))
	e.p.print(errcheck)
}

func (e *encodeGen) Execute(p Elem) error {
	if !e.p.ok() {
		return e.p.err
	}
	if !p.Printable() {
		return nil
	}

	e.p.comment("EncodeMsg implements msgp.Encodable")

	e.p.printf("\nfunc (%s %s) EncodeMsg(en *msgp.Writer) (err error) {", p.Varname(), methodReceiver(p, true))
	next(e, p)
	e.p.nakedReturn()
	unsetReceiver(p)
	return e.p.err
}

func (e *encodeGen) gStruct(s *Struct) {
	if !e.p.ok() {
		return
	}
	if s.AsTuple {
		e.tuple(s)
	} else {
		e.structmap(s)
	}
	return
}

func (e *encodeGen) tuple(s *Struct) {
	nfields := len(s.Fields)
	e.writeAndCheck(arrayHeader, intFmt, nfields)
	for i := range s.Fields {
		if !e.p.ok() {
			return
		}
		next(e, s.Fields[i].FieldElem)
	}
}

func (e *encodeGen) structmap(s *Struct) {
	nfields := len(s.Fields)
	e.writeAndCheck(mapHeader, intFmt, nfields)
	for i := range s.Fields {
		if !e.p.ok() {
			return
		}
		e.writeAndCheck(stringTyp, quotedFmt, s.Fields[i].FieldTag)
		next(e, s.Fields[i].FieldElem)
	}
}

func (e *encodeGen) gMap(m *Map) {
	if !e.p.ok() {
		return
	}
	vname := m.Varname()
	e.writeAndCheck(mapHeader, lenAsUint32, vname)

	e.p.printf("\nfor %s, %s := range %s {", m.Keyidx, m.Validx, vname)
	e.writeAndCheck(stringTyp, literalFmt, m.Keyidx)
	next(e, m.Value)
	e.p.closeblock()
}

func (e *encodeGen) gPtr(s *Ptr) {
	if !e.p.ok() {
		return
	}

	e.p.printf("\nif %s == nil { err = en.WriteNil(); if err != nil { return; } } else {", s.Varname())
	next(e, s.Value)
	e.p.closeblock()
}

func (e *encodeGen) gSlice(s *Slice) {
	if !e.p.ok() {
		return
	}
	e.writeAndCheck(arrayHeader, lenAsUint32, s.Varname())
	e.p.rangeBlock(s.Index, s.Varname(), e, s.Els)
}

func (e *encodeGen) gArray(a *Array) {
	if !e.p.ok() {
		return
	}
	// shortcut for [const]byte
	if be, ok := a.Els.(*BaseElem); ok && (be.Value == Byte || be.Value == Uint8) {
		e.p.printf("\nerr = en.WriteBytes(%s[:])", a.Varname())
		e.p.print(errcheck)
		return
	}

	e.writeAndCheck(arrayHeader, literalFmt, a.Size)
	e.p.rangeBlock(a.Index, a.Varname(), e, a.Els)
}

func (e *encodeGen) gBase(b *BaseElem) {
	if !e.p.ok() {
		return
	}
	vname := b.Varname()
	if b.Convert {
		vname = tobaseConvert(b)
	}

	if b.Value == IDENT { // unknown identity
		e.p.printf("\nerr = %s.EncodeMsg(en)", vname)
		e.p.print(errcheck)
	} else { // typical case
		e.writeAndCheck(b.BaseName(), literalFmt, vname)
	}
}
