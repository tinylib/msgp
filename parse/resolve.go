package parse

import (
	"github.com/tinylib/msgp/gen"
)

// findUnresolved finds identifiers and attempts
// to match them with other known (custom) identifiers.
// any unrecognized identifiers left over after a second
// pass are returned.
func (fs *FileSet) findUnresolved(g gen.Elem) []string {

	switch g := g.(type) {
	case *gen.Ptr:
		return fs.findUnresolved(g.Value)

	case *gen.Slice:
		return fs.findUnresolved(g.Els)

	case *gen.Array:
		return fs.findUnresolved(g.Els)

	case *gen.BaseElem:
		if g.Value == gen.IDENT { // type is unrecognized
			id := g.Ident
			if tp, ok := fs.Identities[id]; ok {

				// skip types that the code generator has seen
				_, ok = fs.processed[id]
				if ok {
					return nil
				}

				// if we have found another identity
				if tp != gen.IDENT {
					// Lower type one level
					i := g.Ident
					*g = gen.BaseElem{
						Value:   tp,   // "true" type
						Ident:   i,    // identifier name
						Convert: true, // requires explicit conversion
					}
					return nil
				}
			}
			return []string{g.Ident}
		}
		return nil

	case *gen.Struct:
		var out []string
		nm := g.Name
		_, ok := fs.Identities[nm]

		// we have to check that the name is
		// not empty (b/c of anonymous embedded structs)
		if !ok && nm != "" {
			out = append(out, nm)
		}

		for _, field := range g.Fields {
			if u := fs.findUnresolved(field.FieldElem); len(u) > 0 {
				out = append(out, u...)
			}
		}
		return out

	case *gen.Map:
		return fs.findUnresolved(g.Value)

	default:
		return nil
	}
}
