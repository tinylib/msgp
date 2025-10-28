package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"strconv"
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
	"shim":          applyShim,
	"replace":       replace,
	"ignore":        ignore,
	"tuple":         astuple,
	"vartuple":      asvartuple,
	"compactfloats": compactfloats,
	"clearomitted":  clearomitted,
	"newtime":       newtime,
	"timezone":      newtimezone,
}

// map of all recognized directives which will be applied
// before process() is called
//
// to add an early directive, define a func([]string, *FileSet) error
// and then add it to this list.
var earlyDirectives = map[string]directive{
	"tag":     tag,
	"pointer": pointer,
	"maps":    maps,
}

var passDirectives = map[string]passDirective{
	"ignore": passignore,
}

func passignore(m gen.Method, text []string, p *gen.Printer) error {
	pushstate(m.String())
	for _, a := range text {
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
			if after, ok := strings.CutPrefix(line.Text, linePrefix); ok {
				out = append(out, after)
			}
		}
	}
	return out
}

//msgp:shim {Type} as:{NewType} using:{toFunc/fromFunc} mode:{Mode}
func applyShim(text []string, f *FileSet) error {
	if len(text) < 4 {
		return fmt.Errorf("shim directive should have at least 3 arguments; found %d", len(text)-1)
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

	text = text[4:]
	for len(text) > 0 {
		arg := strings.TrimSpace(text[0])
		switch {
		case strings.HasPrefix(arg, "mode:"):
			modestr := strings.TrimPrefix(arg, "mode:") // parse mode::{mode}
			switch modestr {
			case "cast":
				be.ShimMode = gen.Cast
			case "convert":
				be.ShimMode = gen.Convert
			default:
				return fmt.Errorf("invalid shim mode; found %s, expected 'cast' or 'convert", modestr)
			}
		case strings.HasPrefix(arg, "witherr:"):
			haserr, err := strconv.ParseBool(strings.TrimPrefix(arg, "witherr:"))
			if err != nil {
				return fmt.Errorf("invalid witherr directive; found %s, expected 'true' or 'false'", arg)
			}
			be.ShimErrs = haserr
		default:
			return fmt.Errorf("invalid shim directive; found %s", arg)
		}
		text = text[1:]
	}

	infof("%s -> %s\n", name, be.Value.String())
	f.findShim(name, be, true)

	return nil
}

//msgp:replace {Type} with:{NewType}
func replace(text []string, f *FileSet) error {
	if len(text) != 3 {
		return fmt.Errorf("replace directive should have only 2 arguments; found %d", len(text)-1)
	}

	name := text[1]
	replacement := strings.TrimPrefix(strings.TrimSpace(text[2]), "with:")

	expr, err := parser.ParseExpr(replacement)
	if err != nil {
		return err
	}
	e := f.parseExpr(expr)
	e.AlwaysPtr(&f.pointerRcv)

	if be, ok := e.(*gen.BaseElem); ok {
		be.Convert = true
		be.Alias(name)
		if be.Value == gen.IDENT {
			be.ShimToBase = "(*" + replacement + ")"
			be.Needsref(true)
		}
	}

	infof("%s -> %s\n", name, replacement)
	f.findShim(name, e, false)

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
			infof("ignoring %s\n", name)
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
				infof(name)
			} else {
				warnf("%s: only structs can be tuples\n", name)
			}
		}
	}
	return nil
}

//msgp:vartuple {TypeA} {TypeB}...
func asvartuple(text []string, f *FileSet) error {
	if len(text) < 2 {
		return nil
	}
	for _, item := range text[1:] {
		name := strings.TrimSpace(item)
		if el, ok := f.Identities[name]; ok {
			if st, ok := el.(*gen.Struct); ok {
				st.AsTuple = true
				st.AsVarTuple = true
				infof(name)
			} else {
				warnf("%s: only structs can be tuples\n", name)
			}
		}
	}
	return nil
}

//msgp:tag {tagname}
func tag(text []string, f *FileSet) error {
	if len(text) != 2 {
		return nil
	}
	f.tagName = strings.TrimSpace(text[1])
	infof("using field tag %q\n", f.tagName)
	return nil
}

//msgp:pointer
func pointer(text []string, f *FileSet) error {
	infof("using pointer receiver\n")
	f.pointerRcv = true
	return nil
}

//msgp:compactfloats
func compactfloats(text []string, f *FileSet) error {
	infof("using compact floats\n")
	f.CompactFloats = true
	return nil
}

//msgp:clearomitted
func clearomitted(text []string, f *FileSet) error {
	infof("clearing omitted fields\n")
	f.ClearOmitted = true
	return nil
}

//msgp:newtime
func newtime(text []string, f *FileSet) error {
	infof("using new time encoding\n")
	f.NewTime = true
	return nil
}

func maps(text []string, f *FileSet) (err error) {
	for _, arg := range text[1:] {
		arg = strings.ToLower(strings.TrimSpace(arg))
		switch {
		case strings.HasPrefix(arg, "shim"):
			arg = strings.TrimPrefix(arg, "shim")
			if len(arg) == 0 {
				f.AllowMapShims = true
				continue
			}
			f.AllowMapShims, err = strconv.ParseBool(strings.TrimPrefix(arg, ":"))
			if err != nil {
				return fmt.Errorf("invalid shim directive; found %s, expected 'true' or 'false'", arg)
			}
		case strings.HasPrefix(arg, "binkeys"):
			arg = strings.TrimPrefix(arg, "binkeys")
			if len(arg) == 0 {
				f.AllowBinMaps = true
				continue
			}
			f.AllowBinMaps, err = strconv.ParseBool(strings.TrimPrefix(arg, ":"))
			if err != nil {
				return fmt.Errorf("invalid binkeys directive; found %s, expected 'true' or 'false'", arg)
			}
		case strings.HasPrefix(arg, "autoshim"):
			arg = strings.TrimPrefix(arg, "autoshim")
			if len(arg) == 0 {
				f.AutoMapShims = true
				continue
			}
			f.AutoMapShims, err = strconv.ParseBool(strings.TrimPrefix(arg, ":"))
			if err != nil {
				return fmt.Errorf("invalid autoshim directive; found %s, expected 'true' or 'false'", arg)
			}
		default:
			if err != nil {
				return fmt.Errorf("invalid autoshim directive; found %s, expected 'true' or 'false'", arg)
			}
		}
	}
	if f.AllowBinMaps && f.AutoMapShims {
		warnf("both binkeys and autoshim are enabled; ignoring autoshim\n")
		f.AutoMapShims = false
	}
	infof("shim:%t binkeys:%t autoshim:%t\n", f.AllowMapShims, f.AllowBinMaps, f.AutoMapShims)
	return nil
}

//msgp:timezone
func newtimezone(text []string, f *FileSet) error {
	if len(text) != 2 {
		return fmt.Errorf("timezone directive should have only 1 argument; found %d", len(text)-1)
	}
	switch strings.ToLower(strings.TrimSpace(text[1])) {
	case "local":
		f.AsUTC = false
	case "utc":
		f.AsUTC = true
	default:
		return fmt.Errorf("timezone directive should be either 'local' or 'utc'; found %q", text[1])
	}
	infof("using timezone %q\n", text[1])
	return nil
}
