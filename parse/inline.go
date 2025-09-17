package parse

import (
	"sort"
	"strings"

	"github.com/tinylib/msgp/gen"
)

// This file defines when and how we
// propagate type information from
// one type declaration to another.
// After the processing pass, every
// non-primitive type is marshalled/unmarshalled/etc.
// through a function call. Here, we propagate
// the type information into the caller's type
// tree *if* the child type is simple enough.
//
// For example, types like
//
//    type A [4]int
//
// will get pushed into parent methods,
// whereas types like
//
//    type B [3]map[string]struct{A, B [4]string}
//
// will not.

// this is an approximate measure
// of the number of children in a node
const maxComplex = 5

// begin recursive search for identities with the
// given name and replace them with e
func (fs *FileSet) findShim(id string, e gen.Elem, addID bool) {
	for name, el := range fs.Identities {
		pushstate(name)
		switch el := el.(type) {
		case *gen.Struct:
			for i := range el.Fields {
				fs.nextShim(&el.Fields[i].FieldElem, id, e)
			}
		case *gen.Array:
			fs.nextShim(&el.Els, id, e)
		case *gen.Slice:
			fs.nextShim(&el.Els, id, e)
		case *gen.Map:
			fs.nextShim(&el.Value, id, e)
			if el.Key != nil {
				fs.nextShim(&el.Key, id, e)
			}
		case *gen.Ptr:
			fs.nextShim(&el.Value, id, e)
		}
		popstate()
	}
	if addID {
		fs.Identities[id] = e
	}
}

func (fs *FileSet) nextShim(ref *gen.Elem, id string, e gen.Elem) {
	if (*ref).TypeName() == id {
		vn := (*ref).Varname()
		*ref = e.Copy()
		(*ref).SetVarname(vn)
	} else {
		switch el := (*ref).(type) {
		case *gen.Struct:
			for i := range el.Fields {
				fs.nextShim(&el.Fields[i].FieldElem, id, e)
			}
		case *gen.Array:
			fs.nextShim(&el.Els, id, e)
		case *gen.Slice:
			fs.nextShim(&el.Els, id, e)
		case *gen.Map:
			fs.nextShim(&el.Value, id, e)
			if el.Key != nil {
				fs.nextShim(&el.Key, id, e)
			}
		case *gen.Ptr:
			fs.nextShim(&el.Value, id, e)
		}
	}
}

// propInline identifies and inlines candidates
func (fs *FileSet) propInline() {
	type gelem struct {
		name string
		el   gen.Elem
	}

	all := make([]gelem, 0, len(fs.Identities))

	for name, el := range fs.Identities {
		all = append(all, gelem{name: name, el: el})
	}

	// make sure we process inlining determinstically:
	// start with the least-complex elems;
	// use identifier names as a tie-breaker
	sort.Slice(all, func(i, j int) bool {
		ig, jg := &all[i], &all[j]
		ic, jc := ig.el.Complexity(), jg.el.Complexity()
		return ic < jc || (ic == jc && ig.name < jg.name)
	})

	for i := range all {
		name := all[i].name
		pushstate(name)
		switch el := all[i].el.(type) {
		case *gen.Struct:
			for i := range el.Fields {
				fs.nextInline(&el.Fields[i].FieldElem, name, el.TypeParams())
			}
		case *gen.Array:
			fs.nextInline(&el.Els, name, el.TypeParams())
		case *gen.Slice:
			fs.nextInline(&el.Els, name, el.TypeParams())
		case *gen.Map:
			fs.nextInline(&el.Value, name, el.TypeParams())
		case *gen.Ptr:
			fs.nextInline(&el.Value, name, el.TypeParams())
		}
		popstate()
	}
}

const fatalloop = `detected infinite recursion in inlining loop!
Please file a bug at github.com/tinylib/msgp/issues!
Thanks!
`

func (fs *FileSet) nextInline(ref *gen.Elem, root string, params gen.GenericTypeParams) {
	switch el := (*ref).(type) {
	case *gen.BaseElem:
		// ensure that we're not inlining
		// a type into itself
		typ := el.TypeName()
		if el.Value == gen.IDENT && typ != root {
			if node, ok := fs.Identities[typ]; ok && node.Complexity() < maxComplex {
				infof("inlining %s\n", typ)

				// This should never happen; it will cause
				// infinite recursion.
				if node == *ref {
					panic(fatalloop)
				}

				*ref = node.Copy()
				fs.nextInline(ref, node.TypeName(), params)
			} else if !ok && !el.Resolved() {
				if params.ToPointerMap[typ] == "" && (!strings.Contains(typ, "[") || !strings.Contains(typ, "]")) {
					// this is the point at which we're sure that
					// we've got a type that isn't a primitive,
					// a library builtin, or a processed type
					warnf("unresolved identifier: %s\n", typ)
				}
			}
		}
	case *gen.Struct:
		for i := range el.Fields {
			fs.nextInline(&el.Fields[i].FieldElem, root, el.TypeParams())
		}
	case *gen.Array:
		fs.nextInline(&el.Els, root, params)
	case *gen.Slice:
		fs.nextInline(&el.Els, root, params)
	case *gen.Map:
		fs.nextInline(&el.Value, root, params)
	case *gen.Ptr:
		fs.nextInline(&el.Value, root, params)
	default:
		panic("bad elem type")
	}
}
