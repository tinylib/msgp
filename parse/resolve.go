package parse

import (
	"github.com/philhofer/msgp/gen"
)

// findUnresolved finds identifiers and attempts
// to match them with other known (custom) identifiers
func findUnresolved(g gen.Elem) []string {
	var out []string

	switch g.(type) {
	case *gen.Ptr:
		out = append(out, findUnresolved(g.(*gen.Ptr).Value)...)
		return out

	case *gen.Slice:
		out = append(out, findUnresolved(g.(*gen.Slice).Els)...)
		return out

	case *gen.BaseElem:
		base := g.(*gen.BaseElem).Value
		if base == gen.IDENT {
			id := g.(*gen.BaseElem).Ident
			tp, ok := globalIdents[id]
			if !ok {
				out = append(out, g.(*gen.BaseElem).Ident)
				return out
			}

			// determine if the identity satisfies
			// is IDENT as a consequence of processing
			_, ok = globalProcessed[g.(*gen.BaseElem).Ident]
			if ok {
				return out
			}

			// TYPE LOWERING
			// this is where a type like `type Custom int`
			// gets lowered to an `int` and explicitly
			// converted in the generated code

			if tp != gen.IDENT {
				// Lower type one level
				i := g.(*gen.BaseElem).Ident
				*(g.(*gen.BaseElem)) = gen.BaseElem{
					Value:   tp,
					Ident:   i, // save identifier
					Convert: true,
				}
			} else {
				out = append(out, g.(*gen.BaseElem).Ident)
			}
		}
		return out

	case *gen.Struct:
		nm := g.(*gen.Struct).Name
		_, ok := globalIdents[nm]

		// we have to check that the name is
		// not empty; otherwise we flag anonymous structs
		if !ok && nm != "" {
			out = append(out, nm)
		}

		for _, field := range g.(*gen.Struct).Fields {
			out = append(out, findUnresolved(field.FieldElem)...)
		}
		return out

	case *gen.Map:
		out = append(out, findUnresolved(g.(*gen.Map).Value)...)
		return out

	default:
		return out
	}
}
