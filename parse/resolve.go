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

	case *gen.Slice:
		out = append(out, findUnresolved(g.(*gen.Slice).Els)...)

	case *gen.BaseElem:
		base := g.(*gen.BaseElem).Value
		if base == gen.IDENT {
			tp, ok := globalIdents[g.(*gen.BaseElem).Ident]
			if !ok {
				out = append(out, g.(*gen.BaseElem).Ident)
			} else {

				// determine if the identity satisfies
				// is IDENT as a consequence of processing
				_, ok = globalProcessed[g.(*gen.BaseElem).Ident]

				// TYPE LOWERING
				// this is where a type like `type Custom int`
				// gets lowered to an `int` and explicitly
				// converted in the generated code

				if !ok && tp != gen.IDENT {
					// Lower type one level
					i := g.(*gen.BaseElem).Ident
					*(g.(*gen.BaseElem)) = gen.BaseElem{
						Value:   tp,
						Ident:   i, // save identifier
						Convert: true,
					}
				} else if !ok && tp == gen.IDENT {
					out = append(out, g.(*gen.BaseElem).Ident)
				}
			}
		}

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

	}
	return out
}
