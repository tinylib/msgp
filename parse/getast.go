package parse

import (
	"fmt"
	"github.com/philhofer/msgp/gen"
	"github.com/ttacon/chalk"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"strings"
)

// just syntactic sugar
type flag struct{}

var set = flag{}

// A FileSet is the in-memory representation of a
// parsed file. The same FileSet will always
// generate the same element (gen.Elem) set.
type FileSet struct {
	Package    string              // package name
	Specs      []*ast.TypeSpec     // type specs in file
	Directives []string            // preprocessor directives
	Identities map[string]gen.Base // alias types (e.g. type Flag uint32)

	processed map[string]flag  // processed type decls
	shims     map[string]*shim // shims
	tuples    map[string]flag  // tuples
}

// File parses a file at the relative path
// provided and produces a new *FileSet.
// (No exported structs is considered an error.)
func File(name string) (*FileSet, error) {
	var files []*ast.File
	var finfo os.FileInfo
	var err error
	var pkg string
	fset := token.NewFileSet()
	finfo, err = os.Stat(name)
	if err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		pkgs, err := parser.ParseDir(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		var one *ast.Package
		if len(pkgs) > 1 {
			return nil, fmt.Errorf("multiple packages in directory: %s", name)
		}
		for _, nm := range pkgs {
			one = nm
			break
		}
		pkg = one.Name
		files = make([]*ast.File, 0, len(one.Files))
		for _, file := range one.Files {
			files = append(files, file)
		}
	} else {
		var f *ast.File
		f, err = parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		files = []*ast.File{f}
		pkg = f.Name.Name
	}

	var comments []string
	for _, fl := range files {
		comments = append(comments, yieldComments(fl.Comments)...)
	}

	// drop non-exported fields
	for _, fl := range files {
		ast.FileExports(fl)
	}

	fs := &FileSet{
		Package:    pkg,
		Specs:      make([]*ast.TypeSpec, 0, 8), // pre-allocate some space
		Directives: comments,
		Identities: make(map[string]gen.Base),
		processed:  make(map[string]flag),
		shims:      make(map[string]*shim),
		tuples:     make(map[string]flag),
	}

	// get specs from each *ast.File
	for _, fl := range files {
		fs.getTypeSpecs(fl)
	}

	if len(fs.Specs) == 0 {
		return nil, fmt.Errorf("no exported definitions in %s", name)
	}

	return fs, nil
}

// ApplyDirectives applies all of the preprocessor
// directives to the file set in the order that they
// appear in the source file.
func (f *FileSet) ApplyDirectives() {
	for _, d := range f.Directives {
		chunks := strings.Split(d, " ")
		if len(chunks) > 0 {
			if fn, ok := directives[chunks[0]]; ok {
				err := fn(chunks, f)
				if err != nil {
					warnf("error applying directive: %s\n", err)
				}
			}
		}
	}
}

// Process processes the file set into generator "trees" that
// can be used for code generation.
func (f *FileSet) Process() []gen.Elem {
	g := make([]gen.Elem, 0, len(f.Specs))

	// process each element, then
	// resolve identifiers. we have
	// to do this in two passes, b/c
	// types are added to the "processed"
	// list as we generate elements.

	// generate elements
	for _, spec := range f.Specs {
		e := f.genElem(spec)
		if e != nil {
			g = append(g, e)
		}
	}
	// resolve typedefs
	var unresolved []string
	for _, el := range g {
		un := f.findUnresolved(el)
		if len(un) > 0 {
			unresolved = append(unresolved, un...)
		}
	}
	// warn about unresolved identifiers
	if len(unresolved) > 0 {
		warnln("unresolved identifiers:")
		for _, u := range unresolved {
			warnf("\t-> %q\n", u)
		}
	}

	// propogate variable names
	for _, e := range g {
		e.SetVarname("z")
	}

	return g
}

// GetElems creates a FileSet from 'filename' and
// returns the processed elements.
func GetElems(filename string) ([]gen.Elem, string, error) {
	fs, err := File(filename)
	if err != nil {
		return nil, "", err
	}
	fs.ApplyDirectives()
	g := fs.Process()
	return g, fs.Package, nil
}

// getTypeSpecs extracts all of the *ast.TypeSpecs in the file.
func (fs *FileSet) getTypeSpecs(f *ast.File) {

	// check all declarations...
	for i := range f.Decls {

		// for GenDecls...
		if g, ok := f.Decls[i].(*ast.GenDecl); ok {

			// and check the specs...
			for _, s := range g.Specs {

				// for ast.TypeSpecs....
				if ts, ok := s.(*ast.TypeSpec); ok {
					//out = append(out, ts)
					fs.Specs = append(fs.Specs, ts)

					// record identifier
					switch a := ts.Type.(type) {
					case *ast.StructType:
						fs.Identities[ts.Name.Name] = gen.IDENT

					case *ast.Ident:
						// we will resolve this later
						fs.Identities[ts.Name.Name] = pullIdent(a.Name)

					case *ast.ArrayType:
						switch k := a.Elt.(type) {
						case *ast.Ident:
							if k.Name == "byte" && a.Len == nil {
								fs.Identities[ts.Name.Name] = gen.Bytes
							} else {
								fs.Identities[ts.Name.Name] = gen.IDENT
							}
						default:
							fs.Identities[ts.Name.Name] = gen.IDENT
						}

					case *ast.StarExpr:
						fs.Identities[ts.Name.Name] = gen.IDENT

					case *ast.MapType:
						fs.Identities[ts.Name.Name] = gen.IDENT

					}
				}
			}
		}
	}
}

// genElem creates the gen.Elem out of an
// ast.TypeSpec. Right now the only supported
// TypeSpec.Type is *ast.StructType. Unsupported
// types will yield a 'nil' return value.
func (fs *FileSet) genElem(in *ast.TypeSpec) gen.Elem {
	if v, ok := in.Type.(*ast.StructType); ok {
		fmt.Printf(chalk.Green.Color("parsing %s..."), in.Name.Name)
		p := &gen.Ptr{
			Value: &gen.Struct{
				Name:   in.Name.Name, // ast.Ident
				Fields: fs.parseFieldList(v.Fields),
			},
		}

		// mark type as processed
		fs.processed[in.Name.Name] = set

		// use as tuple if marked
		if _, ok := fs.tuples[in.Name.Name]; ok {
			p.Value.(*gen.Struct).AsTuple = true
		}

		if len(p.Value.(*gen.Struct).Fields) == 0 {
			fmt.Printf(chalk.Red.Color(" has no exported fields \u2717\n")) // X
			return nil
		}
		fmt.Print(chalk.Green.Color("  \u2713\n")) // check
		return p
	}
	return nil // all non-*ast.StructType elements are unsupported
}

// this is where most of the magic happens
func (fs *FileSet) parseFieldList(fl *ast.FieldList) []gen.StructField {
	if fl == nil || fl.NumFields() == 0 {
		return nil
	}
	out := make([]gen.StructField, 0, fl.NumFields())
	for _, field := range fl.List {
		fds := fs.getField(field)
		if len(fds) > 0 {
			out = append(out, fds...)
		}
	}
	return out
}

// translate *ast.Field into []gen.StructField
func (fs *FileSet) getField(f *ast.Field) []gen.StructField {
	sf := make([]gen.StructField, 1)
	var extension bool
	// parse tag; otherwise field name is field tag
	if f.Tag != nil {
		body := reflect.StructTag(strings.Trim(f.Tag.Value, "`")).Get("msg")
		tags := strings.Split(body, ",")
		if len(tags) == 2 && tags[1] == "extension" {
			extension = true
		}
		// ignore "-" fields
		if tags[0] == "-" {
			return nil
		}
		sf[0].FieldTag = tags[0]
	}

	ex := fs.parseExpr(f.Type)
	if ex == nil {
		return nil
	}

	// parse field name
	switch len(f.Names) {
	case 0:
		sf[0].FieldName = embedded(f.Type)
	case 1:
		sf[0].FieldName = f.Names[0].Name
	default:
		// this is for a multiple in-line declaration,
		// e.g. type A struct { One, Two int }
		sf = sf[0:0]
		for _, nm := range f.Names {
			sf = append(sf, gen.StructField{
				FieldTag:  nm.Name,
				FieldName: nm.Name,
				FieldElem: fs.parseExpr(f.Type),
			})
		}
		return sf
	}
	sf[0].FieldElem = ex
	if sf[0].FieldTag == "" {
		sf[0].FieldTag = sf[0].FieldName
	}

	// validate extension
	if extension {
		switch ex.Type() {
		case gen.PtrType:
			if ex.Ptr().Value.Type() == gen.BaseType {
				ex.Ptr().Value.Base().Value = gen.Ext
			} else {
				warnf(" (\u26a0 field %q couldn't be cast as an extension)", sf[0].FieldName)
				return nil
			}
		case gen.BaseType:
			ex.Base().Value = gen.Ext
		default:
			warnf(" (\u26a0 field %q couldn't be cast as an extension)", sf[0].FieldName)
			return nil
		}
	}
	return sf
}

// extract embedded field name
func embedded(f ast.Expr) string {
	switch f := f.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.StarExpr:
		return embedded(f.X)
	default:
		// other possibilities (like selector expressions)
		// are disallowed; we can't reasonably know
		// their type
		return ""
	}
}

// stringify a field type name
func stringify(e ast.Expr) string {
	switch e := e.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.StarExpr:
		return "*" + stringify(e.X)
	case *ast.SelectorExpr:
		return stringify(e.X) + "." + e.Sel.Name
	case *ast.ArrayType:
		if e.Len == nil {
			return "[]" + stringify(e.Elt)
		}
		return fmt.Sprintf("[%s]%s", stringify(e.Len), stringify(e.Elt))
	case *ast.InterfaceType:
		if e.Methods == nil || e.Methods.NumFields() == 0 {
			return "interface{}"
		}
	}
	return ""
}

// recursively translate ast.Expr to gen.Elem; nil means type not supported
// expected input types:
// - *ast.MapType (map[T]J)
// - *ast.Ident (name)
// - *ast.ArrayType ([(sz)]T)
// - *ast.StarExpr (*T)
// - *ast.StructType (struct {})
// - *ast.SelectorExpr (a.B)
// - *ast.InterfaceType (interface {})
func (fs *FileSet) parseExpr(e ast.Expr) gen.Elem {

	// check for shim; apply and return
	if len(fs.shims) > 0 {
		s := stringify(e)
		if shm, ok := fs.shims[s]; ok {
			return &gen.BaseElem{
				Value:        shm.tp,
				Convert:      true,
				Ident:        s,
				ShimToBase:   shm.to,
				ShimFromBase: shm.from,
			}
		}
	}

	switch e := e.(type) {

	case *ast.MapType:
		if k, ok := e.Key.(*ast.Ident); ok && k.Name == "string" {
			if in := fs.parseExpr(e.Value); in != nil {
				return &gen.Map{Value: in}
			}
		}
		return nil

	case *ast.Ident:
		b := &gen.BaseElem{
			Value: pullIdent(e.Name),
		}
		if b.Value == gen.IDENT {
			b.Ident = (e.Name)
		}
		return b

	case *ast.ArrayType:

		// special case for []byte
		if e.Len == nil {
			if i, ok := e.Elt.(*ast.Ident); ok && i.Name == "byte" {
				return &gen.BaseElem{Value: gen.Bytes}
			}
		}

		// return early if we don't know
		// what the slice element type is
		els := fs.parseExpr(e.Elt)
		if els == nil {
			return nil
		}

		// array and not a slice
		if e.Len != nil {
			switch s := e.Len.(type) {
			case *ast.BasicLit:
				return &gen.Array{
					Size: s.Value,
					Els:  els,
				}

			case *ast.Ident:
				return &gen.Array{
					Size: s.String(),
					Els:  els,
				}

			default: // TODO: support *ast.SelectorExpr; requires custom import(s)
				return nil
			}
		}
		return &gen.Slice{Els: els}

	case *ast.StarExpr:
		if v := fs.parseExpr(e.X); v != nil {
			return &gen.Ptr{Value: v}
		}
		return nil

	case *ast.StructType:
		if fields := fs.parseFieldList(e.Fields); len(fields) > 0 {
			return &gen.Struct{Fields: fields}
		}
		return nil

	case *ast.SelectorExpr:
		// special case for time.Time; others go to Ident
		if im, ok := e.X.(*ast.Ident); ok {
			if e.Sel.Name == "Time" && im.Name == "time" {
				return &gen.BaseElem{Value: gen.Time}
			} else {
				return &gen.BaseElem{
					Value: gen.IDENT,
					Ident: im.Name + "." + e.Sel.Name,
				}
			}
		}
		return nil

	case *ast.InterfaceType:
		// support `interface{}`
		if len(e.Methods.List) == 0 {
			return &gen.BaseElem{Value: gen.Intf}
		}
		return nil

	default: // other types not supported
		return nil
	}
}

// convert an identity to a base type
func pullIdent(name string) gen.Base {
	switch name {
	case "string":
		return gen.String
	case "[]byte":
		return gen.Bytes
	case "byte":
		return gen.Byte
	case "int":
		return gen.Int
	case "int8":
		return gen.Int8
	case "int16":
		return gen.Int16
	case "int32":
		return gen.Int32
	case "int64":
		return gen.Int64
	case "uint":
		return gen.Uint
	case "uint8":
		return gen.Uint8
	case "uint16":
		return gen.Uint16
	case "uint32":
		return gen.Uint32
	case "uint64":
		return gen.Uint64
	case "bool":
		return gen.Bool
	case "float64":
		return gen.Float64
	case "float32":
		return gen.Float32
	case "complex64":
		return gen.Complex64
	case "complex128":
		return gen.Complex128
	case "time.Time":
		return gen.Time
	case "interface{}":
		return gen.Intf
	case "msgp.Extension", "Extension":
		return gen.Ext
	default:
		// unrecognized identity
		return gen.IDENT
	}
}

func infof(s string, v ...interface{})  { fmt.Printf(chalk.Green.Color(s), v...) }
func warnf(s string, v ...interface{})  { fmt.Printf(chalk.Yellow.Color(s), v...) }
func warnln(s string)                   { fmt.Println(chalk.Yellow.Color(s)) }
func fatalf(s string, v ...interface{}) { fmt.Printf(chalk.Red.Color(s), v...) }
