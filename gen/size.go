package gen

import (
	"fmt"
	"strconv"
)

type sizeState uint8

const (
	// need to write "s = ..."
	assign sizeState = iota

	// need to write "s += ..."
	add

	// can just append "+ ..."
	expr
)

type sizeGen struct {
	p     printer
	state sizeState
}

func builtinSize(typ string) string {
	return "msgp." + typ + "Size"
}

// this lets us chain together addition
// operations where possible
func (s *sizeGen) addConstant(sz string) {
	if !s.p.ok() {
		return
	}

	switch s.state {
	case assign:
		s.p.print("\ns = " + sz)
		s.state = expr
		return
	case add:
		s.p.print("\ns += " + sz)
		s.state = expr
		return
	case expr:
		s.p.print(" + " + sz)
		return
	}

	panic("unknown size state")
}

func (s *sizeGen) Execute(p *Ptr) error {
	if !s.p.ok() {
		return s.p.err
	}
	st := p.Value.(*Struct)

	s.p.printf("\nfunc (%s *%s) Msgsize() (s int) {", p.Varname(), st.Name)
	s.state = assign
	s.gStruct(st)
	s.p.nakedReturn()
	return s.p.err
}

func (s *sizeGen) gStruct(st *Struct) {
	if !s.p.ok() {
		return
	}
	if st.AsTuple {
		s.addConstant(builtinSize(arrayHeader))
		for i := range st.Fields {
			if !s.p.ok() {
				return
			}
			next(s, st.Fields[i].FieldElem)
		}
	} else {
		s.addConstant(builtinSize(mapHeader))
		for i := range st.Fields {
			s.addConstant(builtinSize("StringPrefix"))
			s.addConstant(strconv.Itoa(len(st.Fields[i].FieldTag)))
			next(s, st.Fields[i].FieldElem)
		}
	}
}

func (s *sizeGen) gPtr(p *Ptr) {
	s.state = add // inner must use add
	s.p.printf("\nif %s == nil {\ns += msgp.NilSize\n} else {", p.Varname())
	next(s, p.Value)
	s.state = add // closing block; reset to add
	s.p.closeblock()
}

func (s *sizeGen) gSlice(sl *Slice) {
	if !s.p.ok() {
		return
	}
	s.addConstant(builtinSize(arrayHeader))

	if b, ok := sl.Els.(*BaseElem); ok && fixedSize(b.Value) {
		// for something like float64, we can write (len(a)*msgp.Float64Size)
		s.addConstant(fmt.Sprintf("(%s * %s)", sizeExpr(b), lenExpr(sl)))
		return

	}
	// add inside the range block, and immediately after
	s.state = add
	s.p.rangeBlock(sl.Index, sl.Varname(), s, sl.Els)
	s.state = add
}

func (s *sizeGen) gArray(a *Array) {
	if !s.p.ok() {
		return
	}
	s.addConstant(builtinSize(arrayHeader))
	if b, ok := a.Els.(*BaseElem); ok && fixedSize(b.Value) {
		s.addConstant(fmt.Sprintf("(%s * (%s))", a.Size, sizeExpr(b)))
		return
	}
	s.state = add
	s.p.rangeBlock(a.Index, a.Varname(), s, a.Els)
	s.state = add
}

func (s *sizeGen) gMap(m *Map) {
	s.addConstant(builtinSize(mapHeader))
	vn := m.Varname()
	s.p.printf("\nif %s != nil {", vn)
	s.p.printf("\nfor %s, %s := range %s {", m.Keyidx, m.Validx, vn)
	s.p.printf("\n_ = %s", m.Validx) // we may not use the value
	s.p.printf("\ns += msgp.StringPrefixSize + len(%s)", m.Keyidx)
	s.state = expr
	next(s, m.Value)
	s.p.closeblock()
	s.p.closeblock()
	s.state = add
}

func (s *sizeGen) gBase(b *BaseElem) {
	if !s.p.ok() {
		return
	}
	s.addConstant(sizeExpr(b))
}

func lenExpr(sl *Slice) string {
	return "len(" + sl.Varname() + ")"
}

// is a given primitive always the same (max)
// size on the wire?
func fixedSize(p Primitive) bool {
	switch p {
	case Intf, Ext, IDENT, Bytes, String:
		return false
	default:
		return true
	}
}

// strip reference from extension
func stripRef(s string) string {
	if s[0] == '&' {
		return s[1:]
	}
	return s
}

// print size expression of a variable name
func sizeExpr(b *BaseElem) string {
	vname := b.Varname()
	if b.Convert {
		vname = tobaseConvert(b)
	}
	switch b.Value {
	case Ext:
		return "msgp.ExtensionPrefixSize + " + stripRef(vname) + ".Len()"
	case Intf:
		return "msgp.GuessSize(" + vname + ")"
	case IDENT:
		return vname + ".Msgsize()"
	case Bytes:
		return "msgp.BytesPrefixSize + len(" + vname + ")"
	case String:
		return "msgp.StringPrefixSize + len(" + vname + ")"
	default:
		return builtinSize(b.BaseName())
	}
}
