package parse

import (
	"fmt"
	"github.com/philhofer/msgp/gen"
	"go/ast"
	"strings"
)

const linePrefix = "//msgp:"

type directive func([]string, *FileSet) error

// map of all recognized directives
//
// to add a directive, define a func([]string, *FileSet) error
// and then add it to this list.
var directives = map[string]directive{
	"shim":   applyShim,
	"ignore": ignore,
	"tuple":  astuple,
}

type shim struct {
	tp   gen.Base
	to   string // toShim function name
	from string // fromShim function name
}

// find all comment lines that begin with //msgp:
func yieldComments(c []*ast.CommentGroup) []string {
	var out []string
	for _, cg := range c {
		for _, line := range cg.List {
			if strings.HasPrefix(line.Text, linePrefix) {
				out = append(out, strings.TrimPrefix(line.Text, linePrefix))
			}
		}
	}
	return out
}

//msgp:shim {Type} as:{Newtype} using:{toFunc/fromFunc}
func applyShim(text []string, f *FileSet) error {
	if len(text) != 4 {
		return fmt.Errorf("shim directive should have 3 arguments; found %d", len(text)-1)
	}
	if f.shims == nil {
		f.shims = make(map[string]*shim)
	}

	name := text[1]
	if _, ok := f.shims[name]; ok {
		return fmt.Errorf("shim already exists for %s", name)
	}
	tp := pullIdent(strings.TrimPrefix(strings.TrimSpace(text[2]), "as:")) // parse as::{base}

	usestr := strings.TrimPrefix(strings.TrimSpace(text[3]), "using:") // parse using::{method/method}

	methods := strings.Split(usestr, "/")
	if len(methods) != 2 {
		return fmt.Errorf("expected 2 using::{} methods; found %d (%q)", len(methods), text[3])
	}
	infof("applying shim for %s -> %s ...\n", name, tp.String())
	f.shims[name] = &shim{
		tp:   tp,
		to:   methods[0],
		from: methods[1],
	}
	return nil
}

//msgp:ignore {TypeA} {TypeB}...
func ignore(text []string, f *FileSet) error {
	if len(text) < 2 {
		return nil
	}
	for _, item := range text[1:] {
		name := strings.TrimSpace(item)
		for i, dec := range f.Specs {
			if dec != nil && dec.Name != nil && name == dec.Name.Name {
				// delete spec
				f.Specs, f.Specs[i], f.Specs[len(f.Specs)-1] = f.Specs[:len(f.Specs)-1], f.Specs[len(f.Specs)-1], nil
				infof("ignoring: %s...\n", name)
			}
		}
	}
	return nil
}

//msgp:tuple {TypeA} {TypeB}...
func astuple(text []string, f *FileSet) error {
	if len(text) < 2 {
		return nil
	}
	for _, item := range text[1:] {
		name := strings.TrimSpace(item)
		for _, dec := range f.Specs {
			if dec != nil && dec.Name != nil && name == dec.Name.Name {
				f.tuples[name] = set
				infof("using type %s as tuple...\n", name)
			}
		}
	}
	return nil
}
