package gen

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/tinylib/msgp/cfg"
)

func unmarshal(w io.Writer, cfg *cfg.MsgpConfig) *unmarshalGen {
	return &unmarshalGen{
		p:   printer{w: w},
		cfg: cfg,
	}
}

type unmarshalGen struct {
	passes
	p        printer
	hasfield bool
	depth    int
	cfg      *cfg.MsgpConfig
	post     postDefs
}

func (d *unmarshalGen) postLines() {
	lines := strings.Join(d.post.endlines, "\n")
	d.p.printf("\n%s\n", lines)
	d.post.reset()
}

func (u *unmarshalGen) Method() Method { return Unmarshal }

func (u *unmarshalGen) needsField() {
	if u.hasfield {
		return
	}
	u.p.print("\nvar field []byte; _ = field")
	u.hasfield = true
}

func (u *unmarshalGen) Execute(p Elem) error {
	u.hasfield = false
	if !u.p.ok() {
		return u.p.err
	}
	if !IsPrintable(p) {
		return nil
	}

	u.p.comment("UnmarshalMsg implements msgp.Unmarshaler")

	u.p.printf("\nfunc (%s %s) UnmarshalMsg(bts []byte) (o []byte, err error) {", p.Varname(), methodReceiver(p))
	// u.p.printf("\nvar nbs msgp.NilBitsStack;\nvar sawTopNil bool\n if msgp.IsNil(bts) {\n 	sawTopNil = true\n fmt.Printf(\"len of bts pre push: %%v\\n\", len(bts));	bts = nbs.PushAlwaysNil(bts[1:]);\n	fmt.Printf(\"len of bts post push: %%v\\n\", len(bts));\n   }\n")
	u.p.printf("\nvar nbs msgp.NilBitsStack;\nvar sawTopNil bool\n if msgp.IsNil(bts) {\n 	sawTopNil = true\n  bts = nbs.PushAlwaysNil(bts[1:]);\n	}\n")
	next(u, p)
	u.p.print("\n	if sawTopNil {bts = nbs.PopAlwaysNil()}\n o = bts")
	u.p.nakedReturn()
	unsetReceiver(p)
	u.postLines()
	return u.p.err
}

// does assignment to the variable "name" with the type "base"
func (u *unmarshalGen) assignAndCheck(name string, base string) {
	if !u.p.ok() {
		return
	}
	u.p.printf("\n%s, bts, err = nbs.Read%sBytes(bts)", name, base)
	u.p.print(errcheck)
}

func (u *unmarshalGen) gStruct(s *Struct) {
	u.depth++
	defer func() {
		u.depth--
	}()

	if !u.p.ok() {
		return
	}
	if s.AsTuple {
		u.tuple(s)
	} else {
		u.mapstruct(s)
	}
	return
}

func (u *unmarshalGen) tuple(s *Struct) {

	// open block
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignAndCheck(sz, arrayHeader)
	u.p.arrayCheck(strconv.Itoa(len(s.Fields)), sz, "")
	for i := range s.Fields {
		if !u.p.ok() {
			return
		}
		next(u, s.Fields[i].FieldElem)
	}
}

func (u *unmarshalGen) mapstruct(s *Struct) {
	n := len(s.Fields)
	if n == 0 {
		return
	}
	u.needsField()
	k := genSerial()
	tmpl, nStr := genUnmarshalMsgTemplate(k)

	fieldOrder := fmt.Sprintf("\n var unmarshalMsgFieldOrder%s = []string{", nStr)
	for i := range s.Fields {
		fieldOrder += fmt.Sprintf("%q,", s.Fields[i].FieldTag)
	}
	fieldOrder += "}\n"

	varname := strings.Replace(s.TypeName(), "\n", ";", -1)
	u.post.add(varname, "\n// fields of %s%s", varname, fieldOrder)

	u.p.printf("\n const maxFields%s = %d\n", nStr, n)

	found := "found" + nStr
	u.p.printf(tmpl)

	for i := range s.Fields {
		u.p.printf("\ncase \"%s\":", s.Fields[i].FieldTag)
		u.p.printf("\n%s[%d]=true;", found, i)
		u.depth++
		next(u, s.Fields[i].FieldElem)
		u.depth--
		if !u.p.ok() {
			return
		}
	}
	u.p.print("\ndefault:\nbts, err = msgp.Skip(bts)")
	u.p.print(errcheck)
	u.p.print("\n}\n}") // close switch and for loop

	u.p.printf("\n if nextMiss%s != -1 { bts = nbs.PopAlwaysNil(); }\n", nStr)
}

func (u *unmarshalGen) gBase(b *BaseElem) {
	if !u.p.ok() {
		return
	}

	refname := b.Varname() // assigned to
	lowered := b.Varname() // passed as argument
	if b.Convert {
		// begin 'tmp' block
		refname = randIdent()
		lowered = b.ToBase() + "(" + lowered + ")"
		u.p.printf("\n{\nvar %s %s", refname, b.BaseType())
	}

	switch b.Value {
	case Bytes:
		u.p.printf("\n if nbs.AlwaysNil || msgp.IsNil(bts) {\n if !nbs.AlwaysNil { bts = bts[1:]  }\n  %s = %s[:0]} else { %s, bts, err = nbs.ReadBytesBytes(bts, %s)\n", refname, refname, refname, lowered)
		u.p.print(errcheck)
		u.p.closeblock()
	case Ext:
		vn := b.Varname()[1:]
		u.p.printf("\n if nbs.AlwaysNil || msgp.IsNil(bts) { if !nbs.AlwaysNil { bts = bts[1:] }\n    %s = msgp.RawExtension{}  \n} else {\n bts, err = nbs.ReadExtensionBytes(bts, %s) \n", vn, lowered)
		u.p.print(errcheck)
		u.p.closeblock()
	case IDENT:
		u.p.printf("\n  bts, err = %s.UnmarshalMsg(bts);", lowered)
		u.p.print(errcheck)
	default:
		//		u.p.printf("\n if nbs.AlwaysNil || msgp.IsNil(bts) { if !nbs.AlwaysNil { bts=bts[1:]}\n   %s \n} else {  %s, bts, err = nbs.Read%sBytes(bts)\n", b.ZeroLiteral(refname), refname, b.BaseName())
		u.p.printf("\n %s, bts, err = nbs.Read%sBytes(bts)\n", refname, b.BaseName())
	}
	if b.Convert {
		// close 'tmp' block
		u.p.printf("\n%s = %s(%s)\n}", b.Varname(), b.FromBase(), refname)
	}
}

func (u *unmarshalGen) gArray(a *Array) {
	if !u.p.ok() {
		return
	}

	// special case for [const]byte objects
	// see decode.go for symmetry
	if be, ok := a.Els.(*BaseElem); ok && be.Value == Byte {
		u.p.printf("\nbts, err = nbs.ReadExactBytes(bts, %s[:])", a.Varname())
		u.p.print(errcheck)
		return
	}

	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignAndCheck(sz, arrayHeader)
	u.p.arrayCheck(a.Size, sz, "!nbs.IsNil(bts) && ")
	u.p.rangeBlock(a.Index, a.Varname(), u, a.Els)
}

func (u *unmarshalGen) gSlice(s *Slice) {
	if !u.p.ok() {
		return
	}
	u.p.printf("\n if nbs.AlwaysNil { %s \n} else {\n",
		s.ZeroLiteral(`(`+s.Varname()+`)`))
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignAndCheck(sz, arrayHeader)
	u.p.resizeSlice(sz, s)
	u.p.rangeBlock(s.Index, s.Varname(), u, s.Els)
	u.p.closeblock()
}

func (u *unmarshalGen) gMap(m *Map) {
	u.depth++
	defer func() {
		u.depth--
	}()

	if !u.p.ok() {
		return
	}
	u.p.printf("\n if nbs.AlwaysNil { %s \n} else {\n",
		m.ZeroLiteral(m.Varname()))
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignAndCheck(sz, mapHeader)

	// allocate or clear map
	u.p.resizeMap(sz, m)

	// loop and get key,value
	u.p.printf("\nfor %s > 0 {", sz)
	u.p.printf("\nvar %s string; var %s %s; %s--", m.Keyidx, m.Validx, m.Value.TypeName(), sz)
	u.assignAndCheck(m.Keyidx, stringTyp)
	next(u, m.Value)
	u.p.mapAssign(m)
	u.p.closeblock()
	u.p.closeblock()
}

func (u *unmarshalGen) gPtr(p *Ptr) {
	vname := p.Varname()

	base, isBase := p.Value.(*BaseElem)
	if isBase {
		//u.p.printf("\n // we have a BaseElem: %#v  \n", base)
		switch base.Value {
		case IDENT:
			//u.p.printf("\n // we have an IDENT: \n")

			u.p.printf("\n if nbs.AlwaysNil { ")
			u.p.printf("\n if %s != nil { \n", vname)

			niller := fmt.Sprintf("; %s.UnmarshalMsg(msgp.OnlyNilSlice);", vname)

			u.p.printf("%s\n}\n } else { \n // not nbs.AlwaysNil \n", niller)
			u.p.printf("if msgp.IsNil(bts) { bts = bts[1:]; if nil != %s { \n %s}", vname, niller)
			u.p.printf("} else { \n // not nbs.AlwaysNil and not IsNil(bts): have something to read \n")
			u.p.initPtr(p)
			next(u, p.Value)
			u.p.closeblock()
			u.p.closeblock()
			return
		case Ext:
			//u.p.printf("\n // we have an base.Value of Ext: \n")

			u.p.printf("\n if (nbs.AlwaysNil || msgp.IsNil(bts)) { \n // don't try to re-use extension pointers\n if !nbs.AlwaysNil { bts=bts[1:]  }\n %s = nil } else {\n // we have data \n", vname)
			u.p.initPtr(p)
			u.p.printf("\n  bts, err = nbs.ReadExtensionBytes(bts, %s)\n", vname)
			u.p.print(errcheck)
			u.p.closeblock()
			return
		}
	}

	u.p.printf("\nif msgp.IsNil(bts) { bts = bts[1:]; %s = nil; } else { ", vname)
	u.p.initPtr(p)
	next(u, p.Value)
	u.p.closeblock()
}
