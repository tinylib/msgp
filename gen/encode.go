package gen

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/tinylib/msgp/msgp"
)

func encode(w io.Writer) *encodeGen {
	return &encodeGen{
		p: printer{w: w},
	}
}

type encodeGen struct {
	passes
	p    printer
	fuse []byte
	ctx  *Context
}

func (e *encodeGen) Method() Method { return Encode }

func (e *encodeGen) Apply(dirs []string) error {
	return nil
}

func (e *encodeGen) writeAndCheck(typ string, argfmt string, arg any) {
	if e.ctx.compFloats && typ == "Float64" {
		typ = "Float"
	}
	if e.ctx.newTime && typ == "Time" {
		typ = "TimeExt"
	}

	e.p.printf("\nerr = en.Write%s(%s)", typ, fmt.Sprintf(argfmt, arg))
	e.p.wrapErrCheck(e.ctx.ArgsStr())
}

func (e *encodeGen) writeAndCheckWithArrayLimit(typ string, argfmt string, arg any) {
	e.writeAndCheck(typ, argfmt, arg)
	if e.ctx.marshalLimits && e.ctx.arrayLimit != math.MaxUint32 {
		e.p.printf("\nif %s > %slimitArrays {", fmt.Sprintf(argfmt, arg), e.ctx.limitPrefix)
		e.p.printf("\nerr = msgp.ErrLimitExceeded")
		e.p.printf("\nreturn")
		e.p.printf("\n}")
	}
}

func (e *encodeGen) writeAndCheckWithMapLimit(typ string, argfmt string, arg any) {
	e.writeAndCheck(typ, argfmt, arg)
	if e.ctx.marshalLimits && e.ctx.mapLimit != math.MaxUint32 {
		e.p.printf("\nif %s > %slimitMaps {", fmt.Sprintf(argfmt, arg), e.ctx.limitPrefix)
		e.p.printf("\nerr = msgp.ErrLimitExceeded")
		e.p.printf("\nreturn")
		e.p.printf("\n}")
	}
}

func (e *encodeGen) fuseHook() {
	if len(e.fuse) > 0 {
		e.appendraw(e.fuse)
		e.fuse = e.fuse[:0]
	}
}

func (e *encodeGen) Fuse(b []byte) {
	if len(e.fuse) > 0 {
		e.fuse = append(e.fuse, b...)
	} else {
		e.fuse = b
	}
}

// binaryEncodeCall generates code for marshaler interfaces
func (e *encodeGen) binaryEncodeCall(vname, method, writeType, arg string) {
	bts := randIdent()
	e.p.printf("\nvar %s []byte", bts)
	if arg == "" {
		e.p.printf("\n%s, err = %s.%s()", bts, vname, method)
	} else {
		e.p.printf("\n%s, err = %s.%s(%s)", bts, vname, method, arg)
	}
	e.p.wrapErrCheck(e.ctx.ArgsStr())
	if writeType == "String" {
		e.writeAndCheck(writeType, literalFmt, "string("+bts+")")
	} else {
		e.writeAndCheck(writeType, literalFmt, bts)
	}
}

func (e *encodeGen) Execute(p Elem, ctx Context) error {
	e.ctx = &ctx
	if !e.p.ok() {
		return e.p.err
	}
	p = e.applyall(p)
	if p == nil {
		return nil
	}
	if !IsPrintable(p) {
		return nil
	}

	e.p.comment("EncodeMsg implements msgp.Encodable")
	rcv := imutMethodReceiver(p)
	ogVar := p.Varname()
	if p.AlwaysPtr(nil) {
		rcv = methodReceiver(p)
	}
	e.p.printf("\nfunc (%s %s) EncodeMsg(en *msgp.Writer) (err error) {", ogVar, rcv)
	next(e, p)
	if p.AlwaysPtr(nil) {
		p.SetVarname(ogVar)
	}
	e.p.nakedReturn()
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
}

func (e *encodeGen) tuple(s *Struct) {
	nfields := len(s.Fields)
	data := msgp.AppendArrayHeader(nil, uint32(nfields))
	e.p.printf("\n// array header, size %d", nfields)
	e.Fuse(data)
	e.fuseHook()
	for i := range s.Fields {
		if !e.p.ok() {
			return
		}
		fieldElem := s.Fields[i].FieldElem
		anField := s.Fields[i].HasTagPart("allownil") && fieldElem.AllowNil()
		if anField {
			e.p.printf("\nif %s { // allownil: if nil", fieldElem.IfZeroExpr())
			e.p.printf("\nerr = en.WriteNil(); if err != nil { return; }")
			e.p.printf("\n} else {")
		}
		SetIsAllowNil(fieldElem, anField)
		e.ctx.PushString(s.Fields[i].FieldName)
		setTypeParams(s.Fields[i].FieldElem, s.typeParams)
		next(e, s.Fields[i].FieldElem)
		e.ctx.Pop()
		if anField {
			e.p.print("\n}") // close if statement
		}
	}
}

func (e *encodeGen) appendraw(bts []byte) {
	e.p.print("\nerr = en.Append(")
	for i, b := range bts {
		if i != 0 {
			e.p.print(", ")
		}
		e.p.printf("0x%x", b)
	}
	e.p.print(")\nif err != nil { return }")
}

func (e *encodeGen) structmap(s *Struct) {
	oeIdentPrefix := randIdent()

	var data []byte
	nfields := len(s.Fields)
	bm := bmask{
		bitlen:  nfields,
		varname: oeIdentPrefix + "Mask",
	}

	omitempty := s.AnyHasTagPart("omitempty")
	omitzero := s.AnyHasTagPart("omitzero")
	var closeZero bool
	var fieldNVar string
	if omitempty || omitzero {

		fieldNVar = oeIdentPrefix + "Len"

		e.p.printf("\n// check for omitted fields")
		e.p.printf("\n%s := uint32(%d)", fieldNVar, nfields)
		e.p.printf("\n%s", bm.typeDecl())
		e.p.printf("\n_ = %s", bm.varname)
		for i, sf := range s.Fields {
			if !e.p.ok() {
				return
			}
			if ize := sf.FieldElem.IfZeroExpr(); ize != "" && sf.HasTagPart("omitempty") {
				e.p.printf("\nif %s {", ize)
				e.p.printf("\n%s--", fieldNVar)
				e.p.printf("\n%s", bm.setStmt(i))
				e.p.printf("\n}")
			} else if sf.HasTagPart("omitzero") {
				e.p.printf("\nif %s.IsZero() {", sf.FieldElem.Varname())
				e.p.printf("\n%s--", fieldNVar)
				e.p.printf("\n%s", bm.setStmt(i))
				e.p.printf("\n}")
			}
		}

		e.p.printf("\n// variable map header, size %s", fieldNVar)
		e.p.varWriteMapHeader("en", fieldNVar, nfields)
		e.p.print("\nif err != nil { return }")
		if !e.p.ok() {
			return
		}

		// Skip block, if no fields are set.
		if nfields > 1 {
			e.p.printf("\n\n// skip if no fields are to be emitted")
			e.p.printf("\nif %s != 0 {", fieldNVar)
			closeZero = true
		}

	} else {

		// non-omit version
		data = msgp.AppendMapHeader(nil, uint32(nfields))
		e.p.printf("\n// map header, size %d", nfields)
		e.Fuse(data)
		if len(s.Fields) == 0 {
			e.p.printf("\n_ = %s", s.vname)
			e.fuseHook()
		}

	}

	for i := range s.Fields {
		if !e.p.ok() {
			return
		}

		// if field is omitempty or omitzero, wrap with if statement based on the emptymask
		oeField := (omitempty || omitzero) &&
			((s.Fields[i].HasTagPart("omitempty") && s.Fields[i].FieldElem.IfZeroExpr() != "") ||
				s.Fields[i].HasTagPart("omitzero"))
		if oeField {
			e.p.printf("\nif %s == 0 { // if not omitted", bm.readExpr(i))
		}

		data = msgp.AppendString(nil, s.Fields[i].FieldTag)
		e.p.printf("\n// write %q", s.Fields[i].FieldTag)
		e.Fuse(data)
		e.fuseHook()
		fieldElem := s.Fields[i].FieldElem
		anField := !oeField && s.Fields[i].HasTagPart("allownil") && fieldElem.AllowNil()
		if anField {
			e.p.printf("\nif %s { // allownil: if nil", s.Fields[i].FieldElem.IfZeroExpr())
			e.p.printf("\nerr = en.WriteNil(); if err != nil { return; }")
			e.p.printf("\n} else {")
		}
		SetIsAllowNil(fieldElem, anField)

		e.ctx.PushString(s.Fields[i].FieldName)
		setTypeParams(s.Fields[i].FieldElem, s.typeParams)
		next(e, s.Fields[i].FieldElem)
		e.ctx.Pop()

		if oeField || anField {
			e.p.print("\n}") // close if statement
		}
	}
	if closeZero {
		e.p.printf("\n}") // close if statement
	}
}

func (e *encodeGen) gMap(m *Map) {
	if !e.p.ok() {
		return
	}
	e.fuseHook()
	vname := m.Varname()
	e.writeAndCheckWithMapLimit(mapHeader, lenAsUint32, vname)

	e.p.printf("\nfor %s, %s := range %s {", m.Keyidx, m.Validx, vname)
	if m.Key != nil {
		if m.AllowBinMaps {
			e.ctx.PushVar(m.Keyidx)
			m.Key.SetVarname(m.Keyidx)
			next(e, m.Key)
			e.ctx.Pop()
		} else {
			keyIdx := m.Keyidx
			if key, ok := m.Key.(*BaseElem); ok {
				if m.AutoMapShims && CanAutoShim[key.Value] {
					keyIdx = fmt.Sprintf("msgp.AutoShim{}.%sString(%s(%s))", key.Value.String(), strings.ToLower(key.Value.String()), keyIdx)
				} else if key.Value == String {
					keyIdx = fmt.Sprintf("%s(%s)", key.ToBase(), keyIdx)
				} else if key.alias != "" {
					keyIdx = fmt.Sprintf("string(%s)", keyIdx)
				}
			}
			e.writeAndCheck(stringTyp, literalFmt, keyIdx)
		}
	} else {
		e.writeAndCheck(stringTyp, literalFmt, m.Keyidx)
	}
	e.ctx.PushVar(m.Keyidx)
	m.Value.SetIsAllowNil(false)
	setTypeParams(m.Value, m.typeParams)
	next(e, m.Value)
	e.ctx.Pop()
	e.p.closeblock()
}

func (e *encodeGen) gPtr(s *Ptr) {
	if !e.p.ok() {
		return
	}
	e.fuseHook()
	e.p.printf("\nif %s == nil { err = en.WriteNil(); if err != nil { return; } } else {", s.Varname())
	if s.typeParams.TypeParams != "" {
		tp := s.typeParams
		tp.isPtr = true
		s.Value.SetTypeParams(tp)
	}
	next(e, s.Value)
	e.p.closeblock()
}

func (e *encodeGen) gSlice(s *Slice) {
	if !e.p.ok() {
		return
	}
	e.fuseHook()
	e.writeAndCheckWithArrayLimit(arrayHeader, lenAsUint32, s.Varname())
	setTypeParams(s.Els, s.typeParams)
	e.p.rangeBlock(e.ctx, s.Index, s.Varname(), e, s.Els)
}

func (e *encodeGen) gArray(a *Array) {
	if !e.p.ok() {
		return
	}
	e.fuseHook()
	// shortcut for [const]byte
	if be, ok := a.Els.(*BaseElem); ok && (be.Value == Byte || be.Value == Uint8) {
		e.p.printf("\nerr = en.WriteBytes((%s)[:])", a.Varname())
		e.p.wrapErrCheck(e.ctx.ArgsStr())
		return
	}

	e.writeAndCheck(arrayHeader, literalFmt, coerceArraySize(a.Size))
	setTypeParams(a.Els, a.typeParams)
	e.p.rangeBlock(e.ctx, a.Index, a.Varname(), e, a.Els)
}

func (e *encodeGen) gBase(b *BaseElem) {
	if !e.p.ok() {
		return
	}
	e.fuseHook()
	vname := b.Varname()
	if b.Convert {
		if b.ShimMode == Cast {
			vname = tobaseConvert(b)
		} else {
			vname = randIdent()
			e.p.printf("\nvar %s %s", vname, b.BaseType())
			e.p.printf("\n%s, err = %s", vname, tobaseConvert(b))
			e.p.wrapErrCheck(e.ctx.ArgsStr())
		}
	}
	switch b.Value {
	case AInt64, AInt32, AUint64, AUint32, ABool:
		t := strings.TrimPrefix(b.BaseName(), "atomic.")
		e.writeAndCheck(t, literalFmt, strings.TrimPrefix(vname, "*")+".Load()")
	case BinaryMarshaler:
		e.binaryEncodeCall(vname, "MarshalBinary", "Bytes", "")
	case TextMarshalerBin:
		e.binaryEncodeCall(vname, "MarshalText", "Bytes", "")
	case TextMarshalerString:
		e.binaryEncodeCall(vname, "MarshalText", "String", "")
	case BinaryAppender:
		// We do not know if the interface is implemented on pointer or value.
		vname = strings.Trim(vname, "*()")
		e.writeAndCheck("BinaryAppender", literalFmt, vname)
	case TextAppenderBin:
		vname = strings.Trim(vname, "*()")
		e.writeAndCheck("TextAppender", literalFmt, vname)
	case TextAppenderString:
		vname = strings.Trim(vname, "*()")
		e.writeAndCheck("TextAppenderString", literalFmt, vname)
	case IDENT: // unknown identity
		dst := b.BaseType()
		if b.typeParams.isPtr {
			dst = "*" + dst
		}

		// Strip type parameters from dst for lookup in ToPointerMap
		lookupKey := stripTypeParams(dst)
		if idx := strings.Index(dst, "["); idx != -1 {
			lookupKey = dst[:idx]
		}

		if remap := b.typeParams.ToPointerMap[lookupKey]; remap != "" {
			vname = fmt.Sprintf(remap, vname)
		}
		e.p.printf("\nerr = %s.EncodeMsg(en)", vname)
		e.p.wrapErrCheck(e.ctx.ArgsStr())
	default:
		e.writeAndCheck(b.BaseName(), literalFmt, vname)
	}
}
