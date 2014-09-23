package parse

import (
	"errors"
	"fmt"
	"github.com/philhofer/msgp/gen"
	"github.com/ttacon/chalk"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

var (
	// this records a set of all the
	// identifiers in the file that are
	// not go builtins. identities not
	// in this set after the first pass
	// of processing are "unknown" identifiers.
	globalIdents map[string]struct{}
)

func init() {
	globalIdents = make(map[string]struct{})
}

// GetAST simply creates the ast out of a filename and filters
// out non-exported elements.
func GetAST(filename string) (f *ast.File, err error) {
	fset := token.NewFileSet()
	f, err = parser.ParseFile(fset, filename, nil, parser.AllErrors)
	if err != nil {
		return
	}
	if !ast.FileExports(f) {
		f, err = nil, errors.New("no exports in file")
	}
	return
}

// GetElems gets the generator elements out of a file (may be nil)
func GetElems(filename string) ([]gen.Elem, error) {
	f, err := GetAST(filename)
	if err != nil {
		return nil, err
	}
	var out []gen.Elem
	specs := GetTypeSpecs(f)
	if specs == nil {
		return nil, nil
	}

	for i := range specs {
		el := GenElem(specs[i])
		if el != nil {
			out = append(out, el)
		}
	}

	var ptd bool
	for _, o := range out {
		unr := findUnresolved(o)
		if unr != nil {
			if !ptd {
				fmt.Println(chalk.Yellow.Color("Non-local or unresolved identifiers:"))
				ptd = true
			}
			for _, u := range unr {
				fmt.Printf(chalk.Yellow.Color(" -> %q\n"), u)
			}
		}
	}

	return out, nil
}

// should return a list of *ast.TypeSpec we are interested in
func GetTypeSpecs(f *ast.File) []*ast.TypeSpec {
	var out []*ast.TypeSpec

	// check all declarations...
	for i := range f.Decls {

		// for GenDecls...
		if g, ok := f.Decls[i].(*ast.GenDecl); ok {

			// and check the specs...
			for _, s := range g.Specs {

				// for ast.TypeSpecs....
				if ts, ok := s.(*ast.TypeSpec); ok {
					out = append(out, ts)
				}
			}
		}
	}
	return out
}

// GenElem creates the gen.Elem out of an
// ast.TypeSpec. Right now the only supported
// TypeSpec.Type is *ast.StructType
func GenElem(in *ast.TypeSpec) gen.Elem {
	// handle supported types
	switch in.Type.(type) {

	case *ast.StructType:
		v := in.Type.(*ast.StructType)
		fmt.Printf(chalk.Green.Color("parsing %s..."), in.Name.Name)
		p := &gen.Ptr{
			Value: &gen.Struct{
				Name:   in.Name.Name, // ast.Ident
				Fields: parseFieldList(v.Fields),
			},
		}

		// add the struct to our recognized identifiers
		globalIdents[in.Name.Name] = struct{}{}

		if len(p.Value.(*gen.Struct).Fields) == 0 {
			fmt.Printf(chalk.Red.Color(" has no exported fields \u2717\n")) // X
			return nil
		}
		fmt.Print(chalk.Green.Color("  \u2713\n")) // check
		return p

	default:
		return nil

	}
}

func parseFieldList(fl *ast.FieldList) []gen.StructField {
	if fl == nil || fl.NumFields() == 0 {
		return nil
	}
	out := make([]gen.StructField, 0, fl.NumFields())

for_fields:
	for _, field := range fl.List {
		var sf gen.StructField
		// field name

		switch len(field.Names) {
		case 1:
			sf.FieldName = field.Names[0].Name
		case 0:
			sf.FieldName = embedded(field.Type)
		default:
			// inline multiple field declaration
			for _, nm := range field.Names {
				out = append(out, gen.StructField{
					FieldTag:  nm.Name,
					FieldName: nm.Name,
					// we have to duplicate the field
					// type here, or otherwise the FieldElems
					// will be the same pointer
					FieldElem: parseExpr(field.Type),
				})
			}
			continue for_fields
		}

		// field tag
		if field.Tag != nil {
			// we need to trim the leading and trailing ` characters for
			// to convert to reflect.StructTag
			sf.FieldTag = reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get("msg")
		}
		if sf.FieldTag == "" {
			sf.FieldTag = sf.FieldName
		} else if sf.FieldTag == "-" {
			// deliberately ignore field
			continue for_fields
		}

		e := parseExpr(field.Type)
		if e == nil {
			// unsupported type
			fmt.Printf(chalk.Yellow.Color(" (\u26a0 field %q unsupported)"), sf.FieldName)
			continue
		}
		sf.FieldElem = e

		out = append(out, sf)
	}
	return out
}

// extract embedded field name
func embedded(f ast.Expr) string {
	switch f.(type) {
	case *ast.Ident:
		return f.(*ast.Ident).Name
	case *ast.StarExpr:
		return embedded(f.(*ast.StarExpr).X)
	default:
		// not sure what to do here,
		// i don't think any other types
		// are allowed to be embedded
		return ""
	}
}

// go from ast.Expr to gen.Elem; nil means type not supported
func parseExpr(e ast.Expr) gen.Elem {
	switch e.(type) {

	case *ast.MapType:
		// we only support map[string]string and map[string]interface{} right now

		switch e.(*ast.MapType).Key.(type) {
		case *ast.Ident:
			switch e.(*ast.MapType).Key.(*ast.Ident).Name {
			case "string":
				switch e.(*ast.MapType).Value.(*ast.Ident).Name {
				case "string":
					return &gen.BaseElem{
						Value: gen.MapStrStr,
					}
				case "interface{}":
					return &gen.BaseElem{
						Value: gen.MapStrIntf,
					}
				default:
					return nil
				}
			default:
				return nil
			}
		default:
			// we don't support non-string map keys
			return nil
		}

	case *ast.Ident:
		b := &gen.BaseElem{
			Value: pullIdent(e.(*ast.Ident).Name),
		}
		if b.Value == gen.IDENT {
			b.Ident = (e.(*ast.Ident).Name)
		}
		return b

	case *ast.ArrayType:
		// special case for []byte
		switch e.(*ast.ArrayType).Elt.(type) {
		case *ast.Ident:
			i := e.(*ast.ArrayType).Elt.(*ast.Ident)
			if i.Name == "byte" {
				return &gen.BaseElem{
					Value: gen.Bytes,
				}
			} else {
				return &gen.Slice{
					Els: parseExpr(e.(*ast.ArrayType).Elt),
				}
			}
		default:
			return &gen.Slice{
				Els: parseExpr(e.(*ast.ArrayType).Elt),
			}

		}

	case *ast.StarExpr:
		return &gen.Ptr{
			Value: parseExpr(e.(*ast.StarExpr).X),
		}

	case *ast.StructType:
		return &gen.Struct{
			Fields: parseFieldList(e.(*ast.StructType).Fields),
		}

	default: // other types not supported
		return nil
	}
}

func pullIdent(name string) gen.Base {
	switch name {
	case "string":
		return gen.String
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
	default:
		// unrecognized identity
		return gen.IDENT
	}
}
