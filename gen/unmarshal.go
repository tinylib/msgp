package gen

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func unmarshal(w io.Writer) *unmarshalGen {
	return &unmarshalGen{
		p: printer{w: w},
	}
}

type unmarshalGen struct {
	passes
	p        printer
	hasfield bool
	ctx      *Context
}

func (u *unmarshalGen) Method() Method { return Unmarshal }

func (u *unmarshalGen) needsField() {
	if u.hasfield {
		return
	}
	u.p.print("\nvar field []byte; _ = field")
	u.hasfield = true
}

func (u *unmarshalGen) Execute(p Elem, ctx Context) error {
	u.hasfield = false
	u.ctx = &ctx
	if !u.p.ok() {
		return u.p.err
	}
	p = u.applyall(p)
	if p == nil {
		return nil
	}
	if !IsPrintable(p) {
		return nil
	}

	u.p.comment("UnmarshalMsg implements msgp.Unmarshaler")

	u.p.printf("\nfunc (%s %s) UnmarshalMsg(bts []byte) (o []byte, err error) {", p.Varname(), methodReceiver(p))
	next(u, p)
	u.p.print("\no = bts")
	u.p.nakedReturn()
	unsetReceiver(p)
	return u.p.err
}

// does assignment to the variable "name" with the type "base"
func (u *unmarshalGen) assignAndCheck(name string, base string) {
	if !u.p.ok() {
		return
	}
	u.p.printf("\n%s, bts, err = msgp.Read%sBytes(bts)", name, base)
	u.p.wrapErrCheck(u.ctx.ArgsStr())
}

func (u *unmarshalGen) assignArray(name string, base string, fieldLimit uint32) {
	if !u.p.ok() {
		return
	}
	u.p.printf("\n%s, bts, err = msgp.Read%sBytes(bts)", name, base)
	u.p.wrapErrCheck(u.ctx.ArgsStr())

	// Determine effective limit: field limit > context field limit > file limit
	var limit uint32
	var limitName string

	if fieldLimit > 0 {
		// Explicit field limit passed as parameter
		limit = fieldLimit
		limitName = fmt.Sprintf("%d", fieldLimit)
	} else if u.ctx.currentFieldArrayLimit != math.MaxUint32 {
		// Field limit from context (set during field processing)
		limit = u.ctx.currentFieldArrayLimit
		limitName = fmt.Sprintf("%d", u.ctx.currentFieldArrayLimit)
	} else if u.ctx.arrayLimit != math.MaxUint32 {
		// File-level limit
		limit = u.ctx.arrayLimit
		limitName = fmt.Sprintf("%slimitArrays", u.ctx.limitPrefix)
	}

	if limit > 0 && limit != math.MaxUint32 {
		u.p.printf("\nif %s > %s {", name, limitName)
		u.p.printf("\nerr = msgp.ErrLimitExceeded")
		u.p.printf("\nreturn")
		u.p.printf("\n}")
	}
}

func (u *unmarshalGen) assignMap(name string, base string, fieldLimit uint32) {
	if !u.p.ok() {
		return
	}
	u.p.printf("\n%s, bts, err = msgp.Read%sBytes(bts)", name, base)
	u.p.wrapErrCheck(u.ctx.ArgsStr())

	// Determine effective limit: field limit > context field limit > file limit
	var limit uint32
	var limitName string

	if fieldLimit > 0 {
		// Explicit field limit passed as parameter
		limit = fieldLimit
		limitName = fmt.Sprintf("%d", fieldLimit)
	} else if u.ctx.currentFieldMapLimit != math.MaxUint32 {
		// Field limit from context (set during field processing)
		limit = u.ctx.currentFieldMapLimit
		limitName = fmt.Sprintf("%d", u.ctx.currentFieldMapLimit)
	} else if u.ctx.mapLimit != math.MaxUint32 {
		// File-level limit
		limit = u.ctx.mapLimit
		limitName = fmt.Sprintf("%slimitMaps", u.ctx.limitPrefix)
	}

	if limit > 0 && limit != math.MaxUint32 {
		u.p.printf("\nif %s > %s {", name, limitName)
		u.p.printf("\nerr = msgp.ErrLimitExceeded")
		u.p.printf("\nreturn")
		u.p.printf("\n}")
	}
}

// Returns whether a nil check should be done
func (u *unmarshalGen) readBytesWithLimit(refname, lowered string, zerocopy bool, fieldLimit uint32) bool {
	if !u.p.ok() {
		return false
	}

	// Determine effective limit: field limit > context field limit > file limit
	var limit uint32
	var limitName string

	if fieldLimit > 0 {
		// Explicit field limit passed as parameter
		limit = fieldLimit
		limitName = fmt.Sprintf("%d", fieldLimit)
	} else if u.ctx.currentFieldArrayLimit != math.MaxUint32 {
		// Field limit from context (set during field processing)
		limit = u.ctx.currentFieldArrayLimit
		limitName = fmt.Sprintf("%d", u.ctx.currentFieldArrayLimit)
	} else if u.ctx.arrayLimit != math.MaxUint32 {
		// File-level limit
		limit = u.ctx.arrayLimit
		limitName = fmt.Sprintf("%slimitArrays", u.ctx.limitPrefix)
	}

	// Choose reading strategy based on whether limits exist
	if limit > 0 && limit != math.MaxUint32 {
		// Limits exist - use header-first security approach
		sz := randIdent()
		u.p.printf("\nvar %s uint32", sz)
		u.p.printf("\n%s, bts, err = msgp.ReadBytesHeader(bts)", sz)
		u.p.wrapErrCheck(u.ctx.ArgsStr())

		// Check size against limit before allocating
		u.p.printf("\nif %s > %s {", sz, limitName)
		u.p.printf("\nerr = msgp.ErrLimitExceeded")
		u.p.printf("\nreturn")
		u.p.printf("\n}")

		// Now safely read the data
		if zerocopy {
			u.p.printf("\nif uint32(len(bts)) < %s {", sz)
			u.p.printf("\nerr = msgp.ErrShortBytes")
			u.p.printf("\nreturn")
			u.p.printf("\n}")
			u.p.printf("\n%s = bts[:%s]", refname, sz)
			u.p.printf("\nbts = bts[%s:]", sz)
		} else {
			if refname != lowered {
				u.p.printf("\n%s = %s", refname, lowered)
			}
			u.p.printf("\nif %s == nil || uint32(cap(%s)) < %s {", refname, refname, sz)
			u.p.printf("\n%s = make([]byte, %s)", refname, sz)
			u.p.printf("\n} else {")
			u.p.printf("\n%s = %s[:%s]", refname, refname, sz)
			u.p.printf("\n}")

			u.p.printf("\nif uint32(len(bts)) < %s {", sz)
			u.p.printf("\nerr = msgp.ErrShortBytes")
			u.p.printf("\nreturn")
			u.p.printf("\n}")
			u.p.printf("\ncopy(%s, bts[:%s])", refname, sz)
			u.p.printf("\nbts = bts[%s:]", sz)
		}
		return false
	} else {
		// No limits - use original direct reading approach for efficiency
		if zerocopy {
			u.p.printf("\n%s, bts, err = msgp.ReadBytesZC(bts)", refname)
		} else {
			u.p.printf("\n%s, bts, err = msgp.ReadBytesBytes(bts, %s)", refname, lowered)
		}
		u.p.wrapErrCheck(u.ctx.ArgsStr())
		return !zerocopy
	}
}

func (u *unmarshalGen) gStruct(s *Struct) {
	if !u.p.ok() {
		return
	}
	if s.AsTuple {
		u.tuple(s)
	} else {
		u.mapstruct(s)
	}
}

func (u *unmarshalGen) tuple(s *Struct) {
	// open block
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignAndCheck(sz, arrayHeader)
	if s.AsVarTuple {
		u.p.printf("\nif %[1]s == 0 {\no = bts\nreturn\n}", sz)
	} else {
		u.p.arrayCheck(strconv.Itoa(len(s.Fields)), sz)
	}
	for i := range s.Fields {
		if !u.p.ok() {
			return
		}
		u.ctx.PushString(s.Fields[i].FieldName)
		fieldElem := s.Fields[i].FieldElem
		anField := s.Fields[i].HasTagPart("allownil") && fieldElem.AllowNil()

		// Set field-specific limits in context based on struct field's FieldLimit
		if s.Fields[i].FieldLimit > 0 {
			// Apply same limit to both arrays and maps for this field
			u.ctx.SetFieldLimits(s.Fields[i].FieldLimit, s.Fields[i].FieldLimit)
		} else {
			u.ctx.ClearFieldLimits()
		}

		if anField {
			u.p.printf("\nif msgp.IsNil(bts) {\nbts = bts[1:]\n%s = nil\n} else {", fieldElem.Varname())
		}
		SetIsAllowNil(fieldElem, anField)
		if s.Fields[i].HasTagPart("zerocopy") {
			setRecursiveZC(fieldElem, true)
		}
		setTypeParams(fieldElem, s.typeParams)
		next(u, fieldElem)

		// Clear field limits after processing
		u.ctx.ClearFieldLimits()

		if s.Fields[i].HasTagPart("zerocopy") {
			setRecursiveZC(fieldElem, false)
		}

		u.ctx.Pop()
		if anField {
			u.p.printf("\n}")
		}
		if s.AsVarTuple {
			u.p.printf("\nif %[1]s--; %[1]s == 0 {\no = bts\nreturn\n}", sz)
		}
	}
	if s.AsVarTuple {
		u.p.printf("\nfor ; %[1]s > 0; %[1]s-- {\nbts, err = msgp.Skip(bts)\nif err != nil {\nerr = msgp.WrapError(err)\nreturn\n}\n}", sz)
	}
}

// setRecursiveZC will alloc zerocopy for byte fields that are present.
func setRecursiveZC(e Elem, enable bool) {
	if base, ok := e.(*BaseElem); ok {
		base.zerocopy = enable
	}
	if el, ok := e.(*Slice); ok {
		setRecursiveZC(el.Els, enable)
	}
	if el, ok := e.(*Array); ok {
		setRecursiveZC(el.Els, enable)
	}
	if el, ok := e.(*Map); ok {
		setRecursiveZC(el.Value, enable)
	}
}

func (u *unmarshalGen) mapstruct(s *Struct) {
	u.needsField()
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignMap(sz, mapHeader, 0)

	oeCount := s.CountFieldTagPart("omitempty") + s.CountFieldTagPart("omitzero")
	if !u.ctx.clearOmitted {
		oeCount = 0
	}
	bm := bmask{
		bitlen:  oeCount,
		varname: sz + "Mask",
	}
	if oeCount > 0 {
		// Declare mask
		u.p.printf("\n%s", bm.typeDecl())
		u.p.printf("\n_ = %s", bm.varname)
	}
	// Index to field idx of each emitted
	oeEmittedIdx := []int{}

	u.p.printf("\nfor %s > 0 {", sz)
	u.p.printf("\n%s--; field, bts, err = msgp.ReadMapKeyZC(bts)", sz)
	u.p.wrapErrCheck(u.ctx.ArgsStr())
	u.p.print("\nswitch msgp.UnsafeString(field) {")
	for i := range s.Fields {
		if !u.p.ok() {
			return
		}
		u.p.printf("\ncase %q:", s.Fields[i].FieldTag)
		u.ctx.PushString(s.Fields[i].FieldName)

		fieldElem := s.Fields[i].FieldElem
		anField := s.Fields[i].HasTagPart("allownil") && fieldElem.AllowNil()

		// Set field-specific limits in context based on struct field's FieldLimit
		if s.Fields[i].FieldLimit > 0 {
			// Apply same limit to both arrays and maps for this field
			u.ctx.SetFieldLimits(s.Fields[i].FieldLimit, s.Fields[i].FieldLimit)
		} else {
			u.ctx.ClearFieldLimits()
		}

		if anField {
			u.p.printf("\nif msgp.IsNil(bts) {\nbts = bts[1:]\n%s = nil\n} else {", fieldElem.Varname())
		}
		SetIsAllowNil(fieldElem, anField)
		if s.Fields[i].HasTagPart("zerocopy") {
			setRecursiveZC(fieldElem, true)
		}
		setTypeParams(fieldElem, s.typeParams)

		next(u, fieldElem)

		// Clear field limits after processing
		u.ctx.ClearFieldLimits()

		if s.Fields[i].HasTagPart("zerocopy") {
			setRecursiveZC(fieldElem, false)
		}
		u.ctx.Pop()
		if oeCount > 0 && (s.Fields[i].HasTagPart("omitempty") || s.Fields[i].HasTagPart("omitzero")) {
			u.p.printf("\n%s", bm.setStmt(len(oeEmittedIdx)))
			oeEmittedIdx = append(oeEmittedIdx, i)
		}
		if anField {
			u.p.printf("\n}")
		}
	}
	u.p.print("\ndefault:\nbts, err = msgp.Skip(bts)")
	u.p.wrapErrCheck(u.ctx.ArgsStr())
	u.p.print("\n}\n}") // close switch and for loop
	if oeCount > 0 {
		u.p.printf("\n// Clear omitted fields.\n")
		if bm.bitlen > 1 {
			u.p.printf("if %s {\n", bm.notAllSet())
		}
		for bitIdx, fieldIdx := range oeEmittedIdx {
			fieldElem := s.Fields[fieldIdx].FieldElem

			u.p.printf("if %s == 0 {\n", bm.readExpr(bitIdx))
			fze := fieldElem.ZeroExpr()
			if fze != "" {
				u.p.printf("%s = %s\n", fieldElem.Varname(), fze)
			} else {
				u.p.printf("%s = %s{}\n", fieldElem.Varname(), fieldElem.TypeName())
			}
			u.p.printf("}\n")
		}
		if bm.bitlen > 1 {
			u.p.printf("}")
		}
	}
}

// binaryUnmarshalCall generates code for unmarshaling marshaler/appender interfaces
func (u *unmarshalGen) binaryUnmarshalCall(refname, unmarshalMethod, readType string) {
	tmpBytes := randIdent()
	refname = strings.Trim(refname, "(*)")

	u.p.printf("\nvar %s []byte", tmpBytes)
	if readType == "String" {
		u.p.printf("\n%s, bts, err = msgp.ReadStringZC(bts)", tmpBytes)
	} else {
		u.p.printf("\n%s, bts, err = msgp.ReadBytesZC(bts)", tmpBytes)
	}
	u.p.wrapErrCheck(u.ctx.ArgsStr())
	u.p.printf("\nerr = %s.%s(%s)", refname, unmarshalMethod, tmpBytes)
}

func (u *unmarshalGen) gBase(b *BaseElem) {
	if !u.p.ok() {
		return
	}

	refname := b.Varname() // assigned to
	lowered := b.Varname() // passed as argument
	// begin 'tmp' block
	if b.Convert && b.Value != IDENT { // we don't need block for 'tmp' in case of IDENT
		refname = randIdent()
		lowered = b.ToBase() + "(" + lowered + ")"
		u.p.printf("\n{\nvar %s %s", refname, b.BaseType())
	}
	nilCheck := false
	switch b.Value {
	case Bytes:
		nilCheck = u.readBytesWithLimit(refname, lowered, b.zerocopy, 0)
	case Ext:
		u.p.printf("\nbts, err = msgp.ReadExtensionBytes(bts, %s)", lowered)
	case BinaryMarshaler, BinaryAppender:
		u.binaryUnmarshalCall(refname, "UnmarshalBinary", "Bytes")
	case TextMarshalerBin, TextAppenderBin:
		u.binaryUnmarshalCall(refname, "UnmarshalText", "Bytes")
	case TextMarshalerString, TextAppenderString:
		u.binaryUnmarshalCall(refname, "UnmarshalText", "String")
	case IDENT:
		if b.Convert {
			lowered = b.ToBase() + "(" + lowered + ")"
		}
		dst := b.BaseType()
		if b.typeParams.isPtr {
			dst = "*" + dst
		}
		if remap := b.typeParams.ToPointerMap[stripTypeParams(dst)]; remap != "" {
			lowered = fmt.Sprintf(remap, lowered)
		}
		u.p.printf("\nbts, err = %s.UnmarshalMsg(bts)", lowered)
	case Time:
		if u.ctx.asUTC {
			u.p.printf("\n%s, bts, err = msgp.Read%sUTCBytes(bts)", refname, b.BaseName())
		} else {
			u.p.printf("\n%s, bts, err = msgp.Read%sBytes(bts)", refname, b.BaseName())
		}
	case AInt64, AInt32, AUint64, AUint32, ABool:
		tmp := randIdent()
		t := strings.TrimPrefix(b.BaseName(), "atomic.")
		u.p.printf("\n var %s %s", tmp, strings.ToLower(t))
		u.p.printf("\n%s, bts, err = msgp.Read%sBytes(bts)", tmp, t)
		u.p.printf("\n%s.Store(%s)", strings.TrimPrefix(refname, "*"), tmp)

	default:
		u.p.printf("\n%s, bts, err = msgp.Read%sBytes(bts)", refname, b.BaseName())
	}
	if b.Value != Bytes {
		u.p.wrapErrCheck(u.ctx.ArgsStr())
	}

	if nilCheck && b.AllowNil() {
		// Ensure that 0 sized slices are allocated.
		// We are inside the path where the value wasn't nil.
		u.p.printf("\nif %s == nil {\n%s = make([]byte, 0)\n}", refname, refname)
	}

	// close 'tmp' block
	if b.Convert && b.Value != IDENT {
		if b.ShimMode == Cast && !b.ShimErrs {
			u.p.printf("\n%s = %s(%s)\n", b.Varname(), b.FromBase(), refname)
		} else {
			u.p.printf("\n%s, err = %s(%s)\n", b.Varname(), b.FromBase(), refname)
			u.p.wrapErrCheck(u.ctx.ArgsStr())
		}
		u.p.printf("}")
	}
}

func (u *unmarshalGen) gArray(a *Array) {
	if !u.p.ok() {
		return
	}

	// special case for [const]byte objects
	// see decode.go for symmetry
	if be, ok := a.Els.(*BaseElem); ok && be.Value == Byte {
		u.p.printf("\nbts, err = msgp.ReadExactBytes(bts, (%s)[:])", a.Varname())
		u.p.wrapErrCheck(u.ctx.ArgsStr())
		return
	}

	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignAndCheck(sz, arrayHeader)
	u.p.arrayCheck(coerceArraySize(a.Size), sz)
	setTypeParams(a.Els, a.typeParams)
	u.p.rangeBlock(u.ctx, a.Index, a.Varname(), u, a.Els)
}

func (u *unmarshalGen) gSlice(s *Slice) {
	if !u.p.ok() {
		return
	}
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignArray(sz, arrayHeader, 0)
	if s.isAllowNil {
		u.p.resizeSliceNoNil(sz, s)
	} else {
		u.p.resizeSlice(sz, s)
	}
	setTypeParams(s.Els, s.typeParams)
	u.p.rangeBlock(u.ctx, s.Index, s.Varname(), u, s.Els)
}

func (u *unmarshalGen) gMap(m *Map) {
	if !u.p.ok() {
		return
	}
	sz := randIdent()
	u.p.declare(sz, u32)
	u.assignMap(sz, mapHeader, 0)

	// allocate or clear map
	u.p.resizeMap(sz, m)

	// We likely need a field.
	// Add now to not be inside for scope.
	u.needsField()

	// loop and get key,value
	u.p.printf("\nfor %s > 0 {", sz)
	u.p.printf("\nvar %s %s; %s--", m.Validx, m.Value.TypeName(), sz)
	m.readKey(u.ctx, u.p, u, u.assignAndCheck)
	u.ctx.PushVar(m.Keyidx)
	m.Value.SetIsAllowNil(false)
	setTypeParams(m.Value, m.typeParams)
	next(u, m.Value)
	u.ctx.Pop()
	u.p.mapAssign(m)
	u.p.closeblock()
}

func (u *unmarshalGen) gPtr(p *Ptr) {
	u.p.printf("\nif msgp.IsNil(bts) { bts, err = msgp.ReadNilBytes(bts); if err != nil { return }; %s = nil; } else { ", p.Varname())
	u.p.initPtr(p)
	if p.typeParams.TypeParams != "" {
		tp := p.typeParams
		tp.isPtr = true
		p.Value.SetTypeParams(tp)
	}
	next(u, p.Value)
	u.p.closeblock()
}
