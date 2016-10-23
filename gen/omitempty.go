package gen

import (
	"fmt"
)

func emptyOmitter(p *printer, varname string) *omitEmpty {
	return &omitEmpty{
		p:       p,
		varname: varname,
	}
}

type omitEmpty struct {
	p       *printer
	varname string
}

func (s *omitEmpty) gStruct(st *Struct) {
	s.p.printf("false // struct values are never empty\n")
}

func (s *omitEmpty) gPtr(p *Ptr) {
	s.p.printf("(%s == nil) // pointer, omitempty\n", p.vname)
}

func (s *omitEmpty) gSlice(sl *Slice) {
	s.p.printf("%s", IsLenZero(sl.vname))
}

func (s *omitEmpty) gArray(a *Array) {
	s.p.printf("%s", IsLenZero(a.vname))
}

func (s *omitEmpty) gMap(m *Map) {
	s.p.printf("%s", IsLenZero(m.vname))
}

func (s *omitEmpty) gBase(b *BaseElem) {
	switch b.Value {
	case Bytes:
		s.p.printf("%s", IsLenZero(b.Varname()))
	case String:
		s.p.printf("%s", IsLenZero(b.Varname()))
	case Float32, Float64, Complex64, Complex128, Uint, Uint8, Uint16, Uint32, Uint64, Byte, Int, Int8, Int16, Int32, Int64:
		s.p.printf("%s", IsEmptyNumber(b.Varname()))
	case Bool:
		s.p.printf("%s", IsEmptyBool(b.Varname()))
	case Time: // time.Time
		s.p.printf("%s", IsEmptyTime(b.Varname()))
	case Intf: // interface{}
		// assume, for now, never empty. rarely makes sense to serialize these.
		fallthrough
	default:
		s.p.print("false\n")
	}
}

func IsEmptyNumber(f string) string {
	return fmt.Sprintf("(%s == 0) // number, omitempty\n",
		f)
}

func IsLenZero(f string) string {
	return fmt.Sprintf("(len(%s) == 0) // string, omitempty\n",
		f)
}

func IsEmptyBool(f string) string {
	return fmt.Sprintf("(!%s) // bool, omitempty\n",
		f)
}

func IsEmptyTime(f string) string {
	return fmt.Sprintf("(%s.IsZero()) // time.Time, omitempty\n",
		f)
}
