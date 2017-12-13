package parse

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/tinylib/msgp/gen"
)

const linePrefix = "//msgp:"

// func(args, fileset)
type directive func([]string, *FileSet) error

// func(passName, args, printer)
type passDirective func(gen.Method, []string, *gen.Printer) error

// map of all recognized directives
//
// to add a directive, define a func([]string, *FileSet) error
// and then add it to this list.
var directives = map[string]directive{
	"shim":   applyShim,
	"ignore": ignore,
	"tuple":  astuple,
}

// passDirectives lists the directives that can be used together
// with a pass name. See func applyDirs for context.
var passDirectives = map[string]passDirective{
	"ignore": passignore,
}

func passignore(m gen.Method, typeNamePatterns []string, p *gen.Printer) error {
	pushstate(m.String())
	for _, a := range typeNamePatterns {
		p.ApplyDirective(m, gen.IgnoreTypename(a))
		infof("ignoring %s\n", a)
	}
	popstate()
	return nil
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

//msgp:shim {Type} as:{Newtype} using:{toFunc/fromFunc} mode:{Mode}
func applyShim(text []string, f *FileSet) error {
	if len(text) < 4 || len(text) > 5 {
		return fmt.Errorf("shim directive should have 3 or 4 arguments; found %d", len(text)-1)
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

	if len(text) == 5 {
		modestr := strings.TrimPrefix(strings.TrimSpace(text[4]), "mode:") // parse mode::{mode}
		switch modestr {
		case "cast":
			be.ShimMode = gen.Cast
		case "convert":
			be.ShimMode = gen.Convert
		default:
			return fmt.Errorf("invalid shim mode; found %s, expected 'cast' or 'convert", modestr)
		}
	}

	infof("%s -> %s\n", name, be.Value.String())
	f.findShim(name, be)

	return nil
}

//msgp:ignore {TypeA} {TypeB}...
// Parameter "text" includes the string "ignore" and possibly more
// strings that represent either exact type names or regexp patterns.
func ignore(text []string, f *FileSet) error {
	if len(text) < 2 {
		return nil
	}
	for _, typeNamePattern := range text[1:] {
		typeNamePattern = strings.TrimSpace(typeNamePattern)
		for k := range f.Identities {
			if gen.TypeNameMatches(typeNamePattern, f.Identities[k].TypeName()) {
				infof("ignoring %s\n", f.Identities[k].TypeName())
				delete(f.Identities, k)
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
		if el, ok := f.Identities[name]; ok {
			if st, ok := el.(*gen.Struct); ok {
				st.AsTuple = true
				infoln(name)
			} else {
				warnf("%s: only structs can be tuples\n", name)
			}
		}
	}
	return nil
}
