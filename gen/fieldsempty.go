package gen

import (
	"fmt"
	"io"
)

func fieldsempty(w io.Writer) *fieldsEmpty {
	return &fieldsEmpty{
		p: printer{w: w},
	}
}

type fieldsEmpty struct {
	passes
	p     printer
	recvr string
}

func (e *fieldsEmpty) Method() Method { return FieldsEmpty }

func (e *fieldsEmpty) Execute(p Elem) error {
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

	e.recvr = fmt.Sprintf("%s %s", p.Varname(), imutMethodReceiver(p))

	next(e, p)
	return e.p.err
}

func (e *fieldsEmpty) gStruct(s *Struct) {
	e.p.printf("// fieldsNotEmpty supports omitempty tags\n")
	e.p.printf("func (%s) fieldsNotEmpty(isempty []bool) uint32 {", e.recvr)

	nfields := len(s.Fields)
	numOE := 0
	for i := range s.Fields {
		if s.Fields[i].OmitEmpty {
			numOE++
		}
	}
	if numOE == 0 {
		// no fields tagged with omitempty, just return the full field count.
		e.p.printf("\nreturn %d }\n", nfields)
		return
	}

	// remember this to avoid recomputing it in other passes.
	s.hasOmitEmptyTags = true

	om := emptyOmitter(&e.p, s.vname)

	e.p.printf("if len(isempty) == 0 { return %d }\n", nfields)
	e.p.printf("var fieldsInUse uint32 = %d\n", nfields)
	for i := range s.Fields {
		if s.Fields[i].OmitEmpty {
			e.p.printf("isempty[%d] = ", i)
			next(om, s.Fields[i].FieldElem)

			e.p.printf("if isempty[%d] { fieldsInUse-- ; }\n", i)
			//or: e.p.printf("if isempty[%d] { fieldsInUse-- ; fmt.Printf(\"\\n %s is not in use!\\n \")}\n", i, s.Fields[i].FieldTag)
		}
	}
	//e.p.printf("\n fmt.Printf(\"\\n\\n fieldsInUse=%%v\", fieldsInUse) \n\n")
	e.p.printf("\n return fieldsInUse \n}\n")
}

func (e *fieldsEmpty) gPtr(p *Ptr) {}

func (e *fieldsEmpty) gSlice(sl *Slice) {}

func (e *fieldsEmpty) gArray(a *Array) {}

func (e *fieldsEmpty) gMap(m *Map) {}

func (e *fieldsEmpty) gBase(b *BaseElem) {}
