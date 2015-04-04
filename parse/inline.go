package parse

import (
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

// propInline identifies and inlines candidates
func (f *FileSet) propInline() {
	for name, el := range f.Identities {
		switch el := el.(type) {
		case *gen.Struct:
			for i := range el.Fields {
				f.nextInline(&el.Fields[i].FieldElem, name)
			}
		case *gen.Array:
			f.nextInline(&el.Els, name)
		case *gen.Slice:
			f.nextInline(&el.Els, name)
		case *gen.Map:
			f.nextInline(&el.Value, name)
		case *gen.Ptr:
			f.nextInline(&el.Value, name)
		}
	}
}

func (f *FileSet) nextInline(ref *gen.Elem, root string) {
	switch el := (*ref).(type) {
	case *gen.BaseElem:
		// ensure that we're not inlining
		// a type into itself
		typ := el.TypeName()
		if el.Value == gen.IDENT && typ != root {
			if node, ok := f.Identities[typ]; ok && node.Complexity() < maxComplex {

				// in order to ensure bottom-up inlining, we need
				// to make sure this node has already had all its children
				// inlined.
				f.nextInline(&node, root)

				infof("inlining methods for %s into %s...\n", typ, root)
				*ref = node.Copy()
			} else if !ok && !el.Resolved() {
				// this is the point at which we're sure that
				// we've got a type that isn't a primitive,
				// a library builtin, or a processed type
				warnf(" \u26a0 WARNING: unresolved identifier: %s\n", typ)
			}
		}
	case *gen.Struct:
		for i := range el.Fields {
			f.nextInline(&el.Fields[i].FieldElem, root)
		}
	case *gen.Array:
		f.nextInline(&el.Els, root)
	case *gen.Slice:
		f.nextInline(&el.Els, root)
	case *gen.Map:
		f.nextInline(&el.Value, root)
	case *gen.Ptr:
		f.nextInline(&el.Value, root)
	default:
		panic("bad elem type")
	}
}
