package parse

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/tinylib/msgp/gen"
)

const linePrefix = "//msgp:"

type directive func([]string, *FileSet) error

// map of all recognized directives
//
// to add a directive, define a func([]string, *FileSet) error
// and then add it to this list.
var directives = map[string]directive{
	"shim":                 applyShim,
	"ignore":               ignore,
	"tuple":                astuple,
	"encodeValueReceivers": encodeValueReceivers,
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

	name := text[1]
	be := gen.Ident(strings.TrimPrefix(strings.TrimSpace(text[2]), "as:")) // parse as::{base}
	if name[0] == '*' {
		name = name[1:]
		be.Needsref(true)
	}
	be.Alias(name)

	usestr := strings.TrimPrefix(strings.TrimSpace(text[3]), "using:") // parse using::{method/method}

	methods := strings.Split(usestr, "/")
	if len(methods) != 2 {
		return fmt.Errorf("expected 2 using::{} methods; found %d (%q)", len(methods), text[3])
	}

	be.ShimToBase = methods[0]
	be.ShimFromBase = methods[1]

	infof("applying shim for %s to %s...\n", name, be.Value.String())

	f.Identities[name] = be

	return nil
}

//msgp:ignore {TypeA} {TypeB}...
func ignore(text []string, f *FileSet) error {
	if len(text) < 2 {
		return nil
	}
	for _, item := range text[1:] {
		name := strings.TrimSpace(item)
		if _, ok := f.Identities[name]; ok {
			delete(f.Identities, name)
			infof("ignoring: %s...\n", name)
		}
	}
	return nil
}

//msgp:encodeValueReceivers
func encodeValueReceivers(text []string, f *FileSet) error {
	if len(text) != 1 {
		return fmt.Errorf("encodeValueReceivers directive should have 0 argument; found %+v", text)
	}
	f.EncodeValueReceivers = true
	return nil
}

//msgp:tuple {TypeA} {TypeB}...
func astuple(text []string, f *FileSet) error {
	if len(text) < 2 {
		return nil
	}
	for _, item := range text[1:] {
		name := strings.TrimSpace(item)
		if el, ok := f.Identities[name]; ok {
			if st, ok := el.(*gen.Struct); ok {
				st.AsTuple = true
				infof("using type %s as tuple...\n", name)
			}
		}
	}
	return nil
}
