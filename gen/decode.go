package gen

import (
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

func decode(w io.Writer) *decodeGen {
	return &decodeGen{
		p:        printer{w: w},
		hasfield: false,
	}
}

type decodeGen struct {
	passes
	p        printer
	hasfield bool
	ctx      *Context
}

func (d *decodeGen) Method() Method { return Decode }

func (d *decodeGen) needsField() {
	if d.hasfield {
		return
	}
	d.p.print("\nvar field []byte; _ = field")
	d.hasfield = true
}

func (d *decodeGen) Execute(p Elem, ctx Context) error {
	d.ctx = &ctx
	p = d.applyall(p)
	if p == nil {
		return nil
	}
	d.hasfield = false
	if !d.p.ok() {
		return d.p.err
	}

	if !IsPrintable(p) {
		return nil
	}

	d.p.comment("DecodeMsg implements msgp.Decodable")

	d.p.printf("\nfunc (%s %s) DecodeMsg(dc *msgp.Reader) (err error) {", p.Varname(), methodReceiver(p))
	next(d, p)
	d.p.nakedReturn()
	unsetReceiver(p)
	return d.p.err
}

func (d *decodeGen) gStruct(s *Struct) {
	if !d.p.ok() {
		return
	}
	if s.AsTuple {
		d.structAsTuple(s)
	} else {
		d.structAsMap(s)
	}
}

func (d *decodeGen) assignAndCheck(name string, typ string) {
	if !d.p.ok() {
		return
	}
	d.p.printf("\n%s, err = dc.Read%s()", name, typ)
	d.p.wrapErrCheck(d.ctx.ArgsStr())
}

func (d *decodeGen) assignArray(name string, typ string, fieldLimit uint32) {
	if !d.p.ok() {
		return
	}
	d.p.printf("\n%s, err = dc.Read%s()", name, typ)
	d.p.wrapErrCheck(d.ctx.ArgsStr())

	// Determine effective limit: field limit > context field limit > file limit
	var limit uint32
	var limitName string

	if fieldLimit > 0 {
		// Explicit field limit passed as parameter
		limit = fieldLimit
		limitName = fmt.Sprintf("%d", fieldLimit)
	} else if d.ctx.currentFieldArrayLimit != math.MaxUint32 {
		// Field limit from context (set during field processing)
		limit = d.ctx.currentFieldArrayLimit
		limitName = fmt.Sprintf("%d", d.ctx.currentFieldArrayLimit)
	} else if d.ctx.arrayLimit != math.MaxUint32 {
		// File-level limit
		limit = d.ctx.arrayLimit
		limitName = fmt.Sprintf("%slimitArrays", d.ctx.limitPrefix)
	}

	if limit > 0 && limit != math.MaxUint32 {
		d.p.printf("\nif %s > %s {", name, limitName)
		d.p.printf("\nerr = msgp.ErrLimitExceeded")
		d.p.printf("\nreturn")
		d.p.printf("\n}")
	}
}

func (d *decodeGen) assignMap(name string, typ string, fieldLimit uint32) {
	if !d.p.ok() {
		return
	}
	d.p.printf("\n%s, err = dc.Read%s()", name, typ)
	d.p.wrapErrCheck(d.ctx.ArgsStr())

	// Determine effective limit: field limit > context field limit > file limit
	var limit uint32
	var limitName string

	if fieldLimit > 0 {
		// Explicit field limit passed as parameter
		limit = fieldLimit
		limitName = fmt.Sprintf("%d", fieldLimit)
	} else if d.ctx.currentFieldMapLimit != math.MaxUint32 {
		// Field limit from context (set during field processing)
		limit = d.ctx.currentFieldMapLimit
		limitName = fmt.Sprintf("%d", d.ctx.currentFieldMapLimit)
	} else if d.ctx.mapLimit != math.MaxUint32 {
		// File-level limit
		limit = d.ctx.mapLimit
		limitName = fmt.Sprintf("%slimitMaps", d.ctx.limitPrefix)
	}

	if limit > 0 && limit != math.MaxUint32 {
		d.p.printf("\nif %s > %s {", name, limitName)
		d.p.printf("\nerr = msgp.ErrLimitExceeded")
		d.p.printf("\nreturn")
		d.p.printf("\n}")
	}
}

// readBytesWithLimit will read bytes into vname.
// Returns field to check for nil.
func (d *decodeGen) readBytesWithLimit(vname string, fieldLimit uint32) string {
	if !d.p.ok() {
		return ""
	}

	// Determine effective limit: field limit > context field limit > file limit
	var limit uint32
	var limitName string

	if fieldLimit > 0 {
		// Explicit field limit passed as parameter
		limit = fieldLimit
		limitName = fmt.Sprintf("%d", fieldLimit)
	} else if d.ctx.currentFieldArrayLimit != math.MaxUint32 {
		// Field limit from context (set during field processing)
		limit = d.ctx.currentFieldArrayLimit
		limitName = fmt.Sprintf("%d", d.ctx.currentFieldArrayLimit)
	} else if d.ctx.arrayLimit != math.MaxUint32 {
		// File-level limit
		limit = d.ctx.arrayLimit
		limitName = fmt.Sprintf("%slimitArrays", d.ctx.limitPrefix)
	}

	// Choose reading strategy based on whether limits exist
	if limit > 0 && limit != math.MaxUint32 {
		// Limits exist - use header-first security approach
		sz := randIdent()
		d.p.printf("\nvar %s uint32", sz)
		d.p.printf("\n%s, err = dc.ReadBytesHeader()", sz)
		d.p.wrapErrCheck(d.ctx.ArgsStr())

		// Check size against limit before allocating
		d.p.printf("\nif %s > %s {", sz, limitName)
		d.p.printf("\nerr = msgp.ErrLimitExceeded")
		d.p.printf("\nreturn")
		d.p.printf("\n}")

		// Allocate and read the data
		// regular field - ensure always allocated, even for size 0
		d.p.printf("\nif %s == nil || uint32(cap(%s)) < %s {", vname, vname, sz)
		d.p.printf("\n%s = make([]byte, %s)", vname, sz)
		d.p.printf("\n} else {")
		d.p.printf("\n%s = %s[:%s]", vname, vname, sz)
		d.p.printf("\n}")
		d.p.printf("\n_, err = dc.ReadFull(%s)", vname)
		return ""
	} else {
		// No limits - use original direct reading approach for efficiency
		d.p.printf("\n%s, err = dc.ReadBytes(%s)", vname, vname)
		return vname
	}
}

func (d *decodeGen) structAsTuple(s *Struct) {
	sz := randIdent()
	d.p.declare(sz, u32)
	d.assignArray(sz, arrayHeader, 0)
	if s.AsVarTuple {
		d.p.printf("\nif %[1]s == 0 { return }", sz)
	} else {
		d.p.arrayCheck(strconv.Itoa(len(s.Fields)), sz)
	}
	for i := range s.Fields {
		if !d.p.ok() {
			return
		}
		fieldElem := s.Fields[i].FieldElem
		anField := s.Fields[i].HasTagPart("allownil") && fieldElem.AllowNil()

		// Set field-specific limits in context based on struct field's FieldLimit
		if s.Fields[i].FieldLimit > 0 {
			// Apply same limit to both arrays and maps for this field
			d.ctx.SetFieldLimits(s.Fields[i].FieldLimit, s.Fields[i].FieldLimit)
		} else {
			d.ctx.ClearFieldLimits()
		}

		if anField {
			d.p.print("\nif dc.IsNil() {")
			d.p.print("\nerr = dc.ReadNil()")
			d.p.wrapErrCheck(d.ctx.ArgsStr())
			d.p.printf("\n%s = nil\n} else {", s.Fields[i].FieldElem.Varname())
		}
		SetIsAllowNil(fieldElem, anField)
		d.ctx.PushString(s.Fields[i].FieldName)
		setTypeParams(fieldElem, s.typeParams)
		next(d, fieldElem)

		// Clear field limits after processing
		d.ctx.ClearFieldLimits()

		d.ctx.Pop()
		if anField {
			d.p.printf("\n}") // close if statement
		}
		if s.AsVarTuple {
			d.p.printf("\nif %[1]s--; %[1]s == 0 { return }", sz)
		}
	}
	if s.AsVarTuple {
		d.p.printf("\nfor ; %[1]s > 0; %[1]s-- {\nif err = dc.Skip(); err != nil {\nerr = msgp.WrapError(err)\nreturn\n}\n}", sz)
	}
}

func (d *decodeGen) structAsMap(s *Struct) {
	d.needsField()
	sz := randIdent()
	d.p.declare(sz, u32)
	d.assignMap(sz, mapHeader, 0)

	oeCount := s.CountFieldTagPart("omitempty") + s.CountFieldTagPart("omitzero")
	if !d.ctx.clearOmitted {
		oeCount = 0
	}
	bm := bmask{
		bitlen:  oeCount,
		varname: sz + "Mask",
	}
	if oeCount > 0 {
		// Declare mask
		d.p.printf("\n%s", bm.typeDecl())
		d.p.printf("\n_ = %s", bm.varname)
	}
	// Index to field idx of each emitted
	oeEmittedIdx := []int{}

	d.p.printf("\nfor %s > 0 {\n%s--", sz, sz)
	d.assignAndCheck("field", mapKey)
	d.p.print("\nswitch msgp.UnsafeString(field) {")
	for i := range s.Fields {
		d.ctx.PushString(s.Fields[i].FieldName)
		d.p.printf("\ncase %q:", s.Fields[i].FieldTag)
		fieldElem := s.Fields[i].FieldElem
		anField := s.Fields[i].HasTagPart("allownil") && fieldElem.AllowNil()

		// Set field-specific limits in context based on struct field's FieldLimit
		if s.Fields[i].FieldLimit > 0 {
			// Apply same limit to both arrays and maps for this field
			d.ctx.SetFieldLimits(s.Fields[i].FieldLimit, s.Fields[i].FieldLimit)
		} else {
			d.ctx.ClearFieldLimits()
		}

		if anField {
			d.p.print("\nif dc.IsNil() {")
			d.p.print("\nerr = dc.ReadNil()")
			d.p.wrapErrCheck(d.ctx.ArgsStr())
			d.p.printf("\n%s = nil\n} else {", fieldElem.Varname())
		}
		SetIsAllowNil(fieldElem, anField)
		setTypeParams(fieldElem, s.typeParams)
		next(d, fieldElem)

		// Clear field limits after processing
		d.ctx.ClearFieldLimits()

		if oeCount > 0 && (s.Fields[i].HasTagPart("omitempty") || s.Fields[i].HasTagPart("omitzero")) {
			d.p.printf("\n%s", bm.setStmt(len(oeEmittedIdx)))
			oeEmittedIdx = append(oeEmittedIdx, i)
		}
		d.ctx.Pop()
		if !d.p.ok() {
			return
		}
		if anField {
			d.p.printf("\n}") // close if statement
		}
	}
	d.p.print("\ndefault:\nerr = dc.Skip()")
	d.p.wrapErrCheck(d.ctx.ArgsStr())

	d.p.closeblock() // close switch
	d.p.closeblock() // close for loop

	if oeCount > 0 {
		d.p.printf("\n// Clear omitted fields.\n")
		if bm.bitlen > 1 {
			d.p.printf("if %s {\n", bm.notAllSet())
		}
		for bitIdx, fieldIdx := range oeEmittedIdx {
			fieldElem := s.Fields[fieldIdx].FieldElem

			d.p.printf("if %s == 0 {\n", bm.readExpr(bitIdx))
			fze := fieldElem.ZeroExpr()
			if fze != "" {
				d.p.printf("%s = %s\n", fieldElem.Varname(), fze)
			} else {
				d.p.printf("%s = %s{}\n", fieldElem.Varname(), fieldElem.TypeName())
			}
			d.p.printf("}\n")
		}
		if bm.bitlen > 1 {
			d.p.printf("}")
		}
	}
}

func (d *decodeGen) readBytesConvertWithLimit(tmp string, allowNil bool, receiverVar string) {
	if !d.p.ok() {
		return
	}

	// Check if limits exist to decide on reading strategy
	if d.ctx.currentFieldArrayLimit != math.MaxUint32 || d.ctx.arrayLimit != math.MaxUint32 {
		// Limits exist - use header-first approach for security
		sz := randIdent()
		d.p.printf("\nvar %s uint32", sz)
		d.p.printf("\n%s, err = dc.ReadBytesHeader()", sz)
		d.p.wrapErrCheck(d.ctx.ArgsStr())

		// Check array limits for bytes (use currentFieldArrayLimit or arrayLimit)
		if d.ctx.currentFieldArrayLimit != math.MaxUint32 {
			d.p.printf("\nif %s > %d {", sz, d.ctx.currentFieldArrayLimit)
			d.p.printf("\nerr = msgp.ErrLimitExceeded")
			d.p.printf("\nreturn")
			d.p.printf("\n}")
		} else if d.ctx.arrayLimit != math.MaxUint32 {
			d.p.printf("\nif %s > %slimitArrays {", sz, d.ctx.limitPrefix)
			d.p.printf("\nerr = msgp.ErrLimitExceeded")
			d.p.printf("\nreturn")
			d.p.printf("\n}")
		}

		// Allocate and read with type conversion
		if tmp != receiverVar {
			d.p.printf("\n%s = %s", tmp, receiverVar)
		}
		d.p.printf("\nif %s == nil || uint32(cap(%s)) < %s {", tmp, tmp, sz)
		d.p.printf("\n%s = make([]byte, %s)", tmp, sz)
		d.p.printf("\n} else {")
		d.p.printf("\n%s = %s[:%s]", tmp, tmp, sz)
		d.p.printf("\n}")
		d.p.printf("\n_, err = dc.ReadFull(%s)", tmp)
	} else {
		// No limits - use original efficient approach with receiver cast as destination
		d.p.printf("\n%s, err = dc.ReadBytes(%s)", tmp, receiverVar)
	}
}

func (d *decodeGen) gBase(b *BaseElem) {
	if !d.p.ok() {
		return
	}

	// open block for 'tmp'
	var tmp string
	lowered := b.Varname()             // passed as argument
	if b.Convert && b.Value != IDENT { // we don't need block for 'tmp' in case of IDENT
		tmp = randIdent()
		lowered = b.ToBase() + "(" + lowered + ")"
		d.p.printf("\n{ var %s %s", tmp, b.BaseType())
	}

	vname := b.Varname()  // e.g. "z.FieldOne"
	bname := b.BaseName() // e.g. "Float64"
	checkNil := vname     // Name of var to check for nil
	alwaysRef := vname

	// make sure we always reference the pointer
	if strings.Contains(alwaysRef, "*") {
		alwaysRef = strings.Trim(alwaysRef, "*()")
	} else if !b.parentIsPtr {
		alwaysRef = "&" + vname
	}

	// handle special cases
	// for object type.
	switch b.Value {
	case Bytes:
		if b.Convert {
			d.readBytesConvertWithLimit(tmp, b.AllowNil(), lowered)
			checkNil = tmp
		} else {
			checkNil = d.readBytesWithLimit(vname, 0)
		}
	case BinaryMarshaler, BinaryAppender:
		d.p.printf("\nerr = dc.ReadBinaryUnmarshal(%s)", alwaysRef)
	case TextMarshalerBin, TextAppenderBin:
		d.p.printf("\nerr = dc.ReadTextUnmarshal(%s)", alwaysRef)
	case TextMarshalerString, TextAppenderString:
		d.p.printf("\nerr = dc.ReadTextUnmarshalString(%s)", alwaysRef)
	case IDENT:
		dst := b.BaseType()
		if b.typeParams.isPtr {
			dst = "*" + dst
		}

		if b.Convert {
			if remap := b.typeParams.ToPointerMap[stripTypeParams(dst)]; remap != "" {
				vname = fmt.Sprintf(remap, vname)
			}
			lowered := b.ToBase() + "(" + vname + ")"
			d.p.printf("\nerr = %s.DecodeMsg(dc)", lowered)
		} else {
			if remap := b.typeParams.ToPointerMap[stripTypeParams(dst)]; remap != "" {
				vname = fmt.Sprintf(remap, vname)
			}
			d.p.printf("\nerr = %s.DecodeMsg(dc)", vname)
		}
	case Ext:
		d.p.printf("\nerr = dc.ReadExtension(%s)", vname)
	case AInt64, AInt32, AUint64, AUint32, ABool:
		tmp := randIdent()
		t := strings.TrimPrefix(b.BaseName(), "atomic.")
		d.p.printf("\n var %s %s", tmp, strings.ToLower(t))
		d.p.printf("\n%s, err = dc.Read%s()", tmp, t)
		d.p.printf("\n%s.Store(%s)", strings.TrimPrefix(vname, "*"), tmp)
	default:
		if b.Value == Time && d.ctx.asUTC {
			bname += "UTC"
		}
		if b.Convert {
			d.p.printf("\n%s, err = dc.Read%s()", tmp, bname)
		} else {
			d.p.printf("\n%s, err = dc.Read%s()", vname, bname)
		}
	}
	d.p.wrapErrCheck(d.ctx.ArgsStr())

	if checkNil != "" && b.AllowNil() {
		// Ensure that 0 sized slices are allocated.
		d.p.printf("\nif %s == nil {\n%s = make([]byte, 0)\n}", checkNil, checkNil)
	}

	// close block for 'tmp'
	if b.Convert && b.Value != IDENT {
		if b.ShimMode == Cast && !b.ShimErrs {
			d.p.printf("\n%s = %s(%s)\n}", vname, b.FromBase(), tmp)
		} else {
			d.p.printf("\n%s, err = %s(%s)\n}", vname, b.FromBase(), tmp)
			d.p.wrapErrCheck(d.ctx.ArgsStr())
		}
	}
}

func (d *decodeGen) gMap(m *Map) {
	if !d.p.ok() {
		return
	}
	sz := randIdent()

	// resize or allocate map
	d.p.declare(sz, u32)
	d.assignMap(sz, mapHeader, 0)
	d.p.resizeMap(sz, m)

	// for element in map, read string/value
	// pair and assign
	d.needsField()
	d.p.printf("\nfor %s > 0 {\n%s--", sz, sz)
	m.readKey(d.ctx, d.p, d, d.assignAndCheck)
	d.p.declare(m.Validx, m.Value.TypeName())
	d.ctx.PushVar(m.Keyidx)
	m.Value.SetIsAllowNil(false)
	setTypeParams(m.Value, m.typeParams)
	next(d, m.Value)
	d.p.mapAssign(m)
	d.ctx.Pop()
	d.p.closeblock()
}

func (d *decodeGen) gSlice(s *Slice) {
	if !d.p.ok() {
		return
	}
	sz := randIdent()
	d.p.declare(sz, u32)
	d.assignArray(sz, arrayHeader, 0)
	if s.isAllowNil {
		d.p.resizeSliceNoNil(sz, s)
	} else {
		d.p.resizeSlice(sz, s)
	}
	setTypeParams(s.Els, s.typeParams)
	d.p.rangeBlock(d.ctx, s.Index, s.Varname(), d, s.Els)
}

func (d *decodeGen) gArray(a *Array) {
	if !d.p.ok() {
		return
	}

	// special case if we have [const]byte
	if be, ok := a.Els.(*BaseElem); ok && (be.Value == Byte || be.Value == Uint8) {
		d.p.printf("\nerr = dc.ReadExactBytes((%s)[:])", a.Varname())
		d.p.wrapErrCheck(d.ctx.ArgsStr())
		return
	}
	sz := randIdent()
	d.p.declare(sz, u32)
	d.assignAndCheck(sz, arrayHeader)
	d.p.arrayCheck(coerceArraySize(a.Size), sz)
	setTypeParams(a.Els, a.typeParams)
	d.p.rangeBlock(d.ctx, a.Index, a.Varname(), d, a.Els)
}

func (d *decodeGen) gPtr(p *Ptr) {
	if !d.p.ok() {
		return
	}
	d.p.print("\nif dc.IsNil() {")
	d.p.print("\nerr = dc.ReadNil()")
	d.p.wrapErrCheck(d.ctx.ArgsStr())
	d.p.printf("\n%s = nil\n} else {", p.Varname())
	d.p.initPtr(p)
	if p.typeParams.TypeParams != "" {
		tp := p.typeParams
		tp.isPtr = true
		p.Value.SetTypeParams(tp)
	}
	if be, ok := p.Value.(*BaseElem); ok {
		be.parentIsPtr = true
		defer func() { be.parentIsPtr = false }()
	}
	next(d, p.Value)
	d.p.closeblock()
}
