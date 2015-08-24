package gen

import (
	"io"
	"strconv"
)

const (
	structArraySizeVar = "ssz"
	structMapSizeVar   = "isz"
	mapSizeVar         = "msz"
	sliceSizeVar       = "xsz"
	arraySizeVar       = "asz"
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
}

func (d *decodeGen) Method() Method { return Decode }

func (d *decodeGen) needsField() {
	if d.hasfield {
		return
	}
	d.p.print("\nvar field []byte; _ = field")
	d.hasfield = true
}

func (d *decodeGen) Execute(p Elem) error {
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
	return
}

func (d *decodeGen) assignAndCheck(name string, typ string) {
	if !d.p.ok() {
		return
	}
	d.p.printf("\n%s, err = dc.Read%s()", name, typ)
	d.p.print(errcheck)
}

func (d *decodeGen) structAsTuple(s *Struct) {
	nfields := len(s.Fields)

	d.p.print("\n{")
	d.p.declare(structArraySizeVar, u32)
	d.assignAndCheck(structArraySizeVar, arrayHeader)
	d.p.arrayCheck(strconv.Itoa(nfields), structArraySizeVar)
	d.p.print("\n}")
	for i := range s.Fields {
		if !d.p.ok() {
			return
		}
		next(d, s.Fields[i].FieldElem)
	}
}

func (d *decodeGen) structAsMap(s *Struct) {
	d.needsField()
	d.p.declare(structMapSizeVar, u32)
	d.assignAndCheck(structMapSizeVar, mapHeader)

	d.p.print("\nfor isz > 0 {\nisz--")
	d.assignAndCheck("field", mapKey)
	d.p.print("\nswitch msgp.UnsafeString(field) {")
	for i := range s.Fields {
		d.p.printf("\ncase \"%s\":", s.Fields[i].FieldTag)
		next(d, s.Fields[i].FieldElem)
		if !d.p.ok() {
			return
		}
	}
	d.p.print("\ndefault:\nerr = dc.Skip()")
	d.p.print(errcheck)
	d.p.closeblock() // close switch
	d.p.closeblock() // close for loop
}

func (d *decodeGen) gBase(b *BaseElem) {
	if !d.p.ok() {
		return
	}

	// open block for 'tmp'
	if b.Convert {
		d.p.printf("\n{ var tmp %s", b.BaseType())
	}

	vname := b.Varname()  // e.g. "z.FieldOne"
	bname := b.BaseName() // e.g. "Float64"

	// handle special cases
	// for object type.
	switch b.Value {
	case Bytes:
		if b.Convert {
			d.p.printf("\ntmp, err = dc.ReadBytes([]byte(%s))", vname)
		} else {
			d.p.printf("\n%s, err = dc.ReadBytes(%s)", vname, vname)
		}
	case IDENT:
		d.p.printf("\nerr = %s.DecodeMsg(dc)", vname)
	case Ext:
		d.p.printf("\nerr = dc.ReadExtension(%s)", vname)
	default:
		if b.Convert {
			d.p.printf("\ntmp, err = dc.Read%s()", bname)
		} else {
			d.p.printf("\n%s, err = dc.Read%s()", vname, bname)
		}
	}

	// close block for 'tmp'
	if b.Convert {
		d.p.printf("\n%s = %s(tmp)\n}", vname, b.FromBase())
	}

	d.p.print(errcheck)
}

func (d *decodeGen) gMap(m *Map) {
	if !d.p.ok() {
		return
	}

	// resize or allocate map
	d.p.declare(mapSizeVar, u32)
	d.assignAndCheck(mapSizeVar, mapHeader)
	d.p.resizeMap(mapSizeVar, m)

	// for element in map, read string/value
	// pair and assign
	d.p.print("\nfor msz > 0 {\nmsz--")
	d.p.declare(m.Keyidx, "string")
	d.p.declare(m.Validx, m.Value.TypeName())
	d.assignAndCheck(m.Keyidx, stringTyp)
	next(d, m.Value)
	d.p.mapAssign(m)
	d.p.closeblock()
}

func (d *decodeGen) gSlice(s *Slice) {
	if !d.p.ok() {
		return
	}
	d.p.declare(sliceSizeVar, u32)
	d.assignAndCheck(sliceSizeVar, arrayHeader)
	d.p.resizeSlice(sliceSizeVar, s)
	d.p.rangeBlock(s.Index, s.Varname(), d, s.Els)
}

func (d *decodeGen) gArray(a *Array) {
	if !d.p.ok() {
		return
	}

	// special case if we have [const]byte
	if be, ok := a.Els.(*BaseElem); ok && (be.Value == Byte || be.Value == Uint8) {
		d.p.printf("\nerr = dc.ReadExactBytes(%s[:])", a.Varname())
		d.p.print(errcheck)
		return
	}

	d.p.declare(arraySizeVar, u32)
	d.assignAndCheck(arraySizeVar, arrayHeader)
	d.p.arrayCheck(a.Size, arraySizeVar)

	d.p.rangeBlock(a.Index, a.Varname(), d, a.Els)
}

func (d *decodeGen) gPtr(p *Ptr) {
	if !d.p.ok() {
		return
	}
	d.p.print("\nif dc.IsNil() {")
	d.p.print("\nerr = dc.ReadNil()")
	d.p.print(errcheck)
	d.p.printf("\n%s = nil\n} else {", p.Varname())
	d.p.initPtr(p)
	next(d, p.Value)
	d.p.closeblock()
}
