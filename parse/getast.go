package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/tinylib/msgp/gen"
)

// TypeInfo holds both the type expression and its generic type parameters
type TypeInfo struct {
	Type       ast.Expr       // The actual type expression
	TypeParams *ast.FieldList // Generic type parameters
}

// A FileSet is the in-memory representation of a
// parsed file.
type FileSet struct {
	Package       string               // package name
	Specs         map[string]ast.Expr  // type specs in file
	TypeInfos     map[string]*TypeInfo // type specs with generic info
	Identities    map[string]gen.Elem  // processed from specs
	Aliased       map[string]string    // Aliased types.
	Directives    []string             // raw preprocessor directives
	Imports       []*ast.ImportSpec    // imports
	CompactFloats bool                 // Use smaller floats when feasible
	ClearOmitted  bool                 // Set omitted fields to zero value
	NewTime       bool                 // Set to use -1 extension for time.Time
	AsUTC         bool                 // Set timezone to UTC instead of local
	AllowMapShims bool                 // Allow map keys to be shimmed (default true)
	AllowBinMaps  bool                 // Allow maps with binary keys to be used (default false)
	AutoMapShims  bool                 // Automatically shim map keys of builtin types(default false)
	ArrayLimit    uint32               // Maximum array/slice size allowed during deserialization
	MapLimit      uint32               // Maximum map size allowed during deserialization
	MarshalLimits bool                 // Whether to enforce limits during marshaling
	LimitPrefix   string               // Unique prefix for limit constants to avoid collisions

	tagName    string // tag to read field names from
	pointerRcv bool   // generate with pointer receivers.
}

// File parses a file at the relative path
// provided and produces a new *FileSet.
// If you pass in a path to a directory, the entire
// directory will be parsed.
// If unexport is false, only exported identifiers are included in the FileSet.
// If the resulting FileSet would be empty, an error is returned.
func File(name string, unexported bool, directives []string) (*FileSet, error) {
	pushstate(name)
	defer popstate()
	fs := &FileSet{
		Specs:      make(map[string]ast.Expr),
		TypeInfos:  make(map[string]*TypeInfo),
		Identities: make(map[string]gen.Elem),
		Directives: append([]string{}, directives...),
		ArrayLimit: math.MaxUint32,
		MapLimit:   math.MaxUint32,
	}

	fset := token.NewFileSet()
	finfo, err := os.Stat(name)
	if err != nil {
		return nil, err
	}
	if finfo.IsDir() {
		pkgs, err := parser.ParseDir(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		if len(pkgs) != 1 {
			return nil, fmt.Errorf("multiple packages in directory: %s", name)
		}
		var one *ast.Package
		for _, nm := range pkgs {
			one = nm
			break
		}
		fs.Package = one.Name
		for _, fl := range one.Files {
			pushstate(fl.Name.Name)
			fs.Directives = append(fs.Directives, yieldComments(fl.Comments)...)
			if !unexported {
				ast.FileExports(fl)
			}
			fs.getTypeSpecs(fl)
			popstate()
		}
	} else {
		f, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
		if err != nil {
			return nil, err
		}
		fs.Package = f.Name.Name
		fs.Directives = append(fs.Directives, yieldComments(f.Comments)...)
		if !unexported {
			ast.FileExports(f)
		}
		fs.getTypeSpecs(f)
	}

	if len(fs.Specs) == 0 {
		return nil, fmt.Errorf("no definitions in %s", name)
	}

	fs.applyEarlyDirectives()
	fs.process()
	fs.applyDirectives()
	fs.propInline()

	return fs, nil
}

// applyDirectives applies all of the directives that
// are known to the parser. additional method-specific
// directives remain in f.Directives
func (fs *FileSet) applyDirectives() {
	newdirs := make([]string, 0, len(fs.Directives))
	for _, d := range fs.Directives {
		chunks := strings.Split(d, " ")
		if len(chunks) > 0 {
			if fn, ok := directives[chunks[0]]; ok {
				pushstate(chunks[0])
				err := fn(chunks, fs)
				if err != nil {
					warnf("directive error: %s", err)
				}
				popstate()
			} else {
				newdirs = append(newdirs, d)
			}
		}
	}
	// Apply aliases last, so we don't overrule any manually specified replace directives.
	for _, d := range fs.Aliased {
		chunks := strings.Split(d, " ")
		if len(chunks) > 0 {
			if fn, ok := directives[chunks[0]]; ok {
				pushstate(chunks[0])
				err := fn(chunks, fs)
				if err != nil {
					warnf("directive error: %s", err)
				}
				popstate()
			} else {
				newdirs = append(newdirs, d)
			}
		}
	}
	fs.Directives = newdirs
}

// applyEarlyDirectives applies all early directives needed before process() is called.
// additional directives remain in f.Directives for future processing
func (fs *FileSet) applyEarlyDirectives() {
	newdirs := make([]string, 0, len(fs.Directives))
	for _, d := range fs.Directives {
		parts := strings.Split(d, " ")
		if len(parts) == 0 {
			continue
		}
		if fn, ok := earlyDirectives[parts[0]]; ok {
			pushstate(parts[0])
			err := fn(parts, fs)
			if err != nil {
				warnf("early directive error: %s", err)
			}
			popstate()
		} else {
			newdirs = append(newdirs, d)
		}
	}
	fs.Directives = newdirs
}

// A linkset is a graph of unresolved
// identities.
//
// Since gen.Ident can only represent
// one level of type indirection (e.g. Foo -> uint8),
// type declarations like `type Foo Bar`
// aren't resolve-able until we've processed
// everything else.
//
// The goal of this dependency resolution
// is to distill the type declaration
// into just one level of indirection.
// In other words, if we have:
//
//	type A uint64
//	type B A
//	type C B
//	type D C
//
// ... then we want to end up
// figuring out that D is just a uint64.
type linkset map[string]*gen.BaseElem

func (fs *FileSet) resolve(ls linkset) {
	progress := true
	for progress && len(ls) > 0 {
		progress = false
		for name, elem := range ls {
			real, ok := fs.Identities[elem.TypeName()]
			if ok {
				// copy the old type descriptor,
				// alias it to the new value,
				// and insert it into the resolved
				// identities list
				progress = true
				nt := real.Copy()
				nt.Alias(name)
				fs.Identities[name] = nt
				delete(ls, name)
			}
		}
	}

	// what's left can't be resolved
	for name, elem := range ls {
		warnf("couldn't resolve type %s (%s)\n", name, elem.TypeName())
	}
}

// formatTypeParams converts an AST FieldList to a string representation
func formatTypeParams(params *ast.FieldList) string {
	if params == nil || params.NumFields() == 0 {
		return ""
	}

	var paramStrs []string
	for _, field := range params.List {
		// Convert underscores to _RTn where n is the number of the parameter
		convert := isrtfor(field.Type)

		// Each field can have multiple names (e.g., T, U constraint)
		for _, name := range field.Names {
			if convert && name.Name == "_" {
				name.Name = fmt.Sprintf("_RT%d", len(paramStrs)+1)
			}
			// For method receivers, we only include the type parameter name
			// The constraints are defined in the type declaration, not the method receiver
			paramStrs = append(paramStrs, name.Name)
		}
	}

	return "[" + strings.Join(paramStrs, ", ") + "]"
}

// isrtfor returns whether the provided expression is a msgp.RTFor[T] pattern.
func isrtfor(t ast.Expr) bool { return strings.HasPrefix(stringify(t), "msgp.RTFor[") }

// findRTForInInterface recursively searches for msgp.RTFor[T] patterns within interface types
func findRTForInInterface(iface *ast.InterfaceType) []string {
	var rtfors []string
	if iface.Methods == nil {
		return rtfors
	}

	for _, method := range iface.Methods.List {
		// Check if this is an embedded interface/type
		if len(method.Names) == 0 {
			if isrtfor(method.Type) {
				rtfors = append(rtfors, stringify(method.Type))
			}
			// Recursively check nested interfaces
			if nestedIface, ok := method.Type.(*ast.InterfaceType); ok {
				rtfors = append(rtfors, findRTForInInterface(nestedIface)...)
			}
		}
	}
	return rtfors
}

// formatTypeParams converts an AST FieldList to a string representation.
// For 'Foo[T any, P msgp.RTFor[T]]' will return {"T": "P"}.
func getMspTypeParams(params *ast.FieldList) map[string]string {
	if params == nil || params.NumFields() == 0 {
		return nil
	}

	paramStrs := make(map[string]string)
	for _, field := range params.List {

		// Handle simple msgp.RTFor[T] constraints
		if isrtfor(field.Type) {
			t := strings.TrimSuffix(strings.TrimPrefix(stringify(field.Type), "msgp.RTFor["), "]")
			for _, name := range field.Names {
				paramStrs[t] = name.Name + "(&%s)"
				paramStrs["*"+t] = name.Name + "(%s)"
				paramStrs[name.Name] = "%s"
				infof("found generic type %s, with roundtrippper %s\n", t, name.Name)
			}
			continue
		}

		// Handle complex interface constraints that embed msgp.RTFor[T]
		if iface, ok := field.Type.(*ast.InterfaceType); ok {
			rtfors := findRTForInInterface(iface)
			for _, rtfor := range rtfors {
				t := strings.TrimSuffix(strings.TrimPrefix(rtfor, "msgp.RTFor["), "]")
				for _, name := range field.Names {
					paramStrs[t] = name.Name + "(&%s)"
					paramStrs["*"+t] = name.Name + "(%s)"
					paramStrs[name.Name] = "%s"
					infof("found generic type %s, with roundtrippper %s (in complex interface)\n", t, name.Name)
				}
			}
		}
	}

	return paramStrs
}

// process takes the contents of f.Specs and
// uses them to populate f.Identities
func (fs *FileSet) process() {
	deferred := make(linkset)
parse:
	for name, def := range fs.Specs {
		pushstate(name)
		el := fs.parseExpr(def)
		if el == nil {
			warnf("failed to parse")
			popstate()
			continue parse
		}
		el.AlwaysPtr(&fs.pointerRcv)

		// Apply type parameters if available
		if typeInfo, ok := fs.TypeInfos[name]; ok && typeInfo.TypeParams != nil {
			typeParamsStr := formatTypeParams(typeInfo.TypeParams)
			ptrMap := getMspTypeParams(typeInfo.TypeParams)
			if typeParamsStr != "" && ptrMap != nil {
				el.SetTypeParams(gen.GenericTypeParams{
					TypeParams:   typeParamsStr,
					ToPointerMap: ptrMap,
				})
			}
		}

		// push unresolved identities into
		// the graph of links and resolve after
		// we've handled every possible named type.
		if be, ok := el.(*gen.BaseElem); ok && be.Value == gen.IDENT {
			deferred[name] = be
			popstate()
			continue parse
		}
		el.Alias(name)
		fs.Identities[name] = el
		popstate()
	}

	if len(deferred) > 0 {
		fs.resolve(deferred)
	}
}

func strToMethod(s string) gen.Method {
	switch s {
	case "encode":
		return gen.Encode
	case "decode":
		return gen.Decode
	case "test":
		return gen.Test
	case "size":
		return gen.Size
	case "marshal":
		return gen.Marshal
	case "unmarshal":
		return gen.Unmarshal
	default:
		return 0
	}
}

func (fs *FileSet) applyDirs(p *gen.Printer) {
	// apply directives of the form
	//
	// 	//msgp:encode ignore {{TypeName}}
	//
loop:
	for _, d := range fs.Directives {
		chunks := strings.Split(d, " ")
		if len(chunks) > 1 {
			for i := range chunks {
				chunks[i] = strings.TrimSpace(chunks[i])
			}
			m := strToMethod(chunks[0])
			if m == 0 {
				warnf("unknown pass name: %q\n", chunks[0])
				continue loop
			}
			if fn, ok := passDirectives[chunks[1]]; ok {
				pushstate(chunks[1])
				err := fn(m, chunks[2:], p)
				if err != nil {
					warnf("error applying directive: %s\n", err)
				}
				popstate()
			} else {
				warnf("unrecognized directive %q\n", chunks[1])
			}
		} else {
			warnf("empty directive: %q\n", d)
		}
	}
	p.CompactFloats = fs.CompactFloats
	p.ClearOmitted = fs.ClearOmitted
	p.NewTime = fs.NewTime
	p.AsUTC = fs.AsUTC
	p.ArrayLimit = fs.ArrayLimit
	p.MapLimit = fs.MapLimit
	p.MarshalLimits = fs.MarshalLimits
	p.LimitPrefix = fs.LimitPrefix
}

func (fs *FileSet) PrintTo(p *gen.Printer) error {
	fs.applyDirs(p)
	names := make([]string, 0, len(fs.Identities))
	for name := range fs.Identities {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		el := fs.Identities[name]
		el.SetVarname("z")
		pushstate(el.TypeName())
		err := p.Print(el)
		popstate()
		if err != nil {
			return err
		}
	}
	return nil
}

// getTypeSpecs extracts all of the *ast.TypeSpecs in the file
// into fs.Identities, but does not set the actual element
func (fs *FileSet) getTypeSpecs(f *ast.File) {
	// collect all imports...
	fs.Imports = append(fs.Imports, f.Imports...)

	// check all declarations...
	for i := range f.Decls {
		// for GenDecls...
		if g, ok := f.Decls[i].(*ast.GenDecl); ok {
			// and check the specs...
			for _, s := range g.Specs {
				// for ast.TypeSpecs....
				if ts, ok := s.(*ast.TypeSpec); ok {
					// Handle type aliases, by adding a "replace" directive.
					if ts.Assign != 0 {
						if fs.Aliased == nil {
							fs.Aliased = make(map[string]string)
						}
						fs.Aliased[ts.Name.Name] = fmt.Sprintf("replace %s with:%s", ts.Name.Name, stringify(ts.Type))
						continue
					}

					switch ts.Type.(type) {
					// this is the list of parse-able
					// type specs
					case *ast.ArrayType,
						*ast.StarExpr,
						*ast.Ident,
						*ast.StructType,
						*ast.MapType:
						fs.Specs[ts.Name.Name] = ts.Type
						// Store type info (no type params for non-struct types yet)
						fs.TypeInfos[ts.Name.Name] = &TypeInfo{
							Type:       ts.Type,
							TypeParams: ts.TypeParams,
						}
					}
				}
			}
		}
	}
}

func fieldName(f *ast.Field) string {
	switch len(f.Names) {
	case 0:
		return stringify(f.Type)
	case 1:
		return f.Names[0].Name
	default:
		return f.Names[0].Name + " (and others)"
	}
}

func (fs *FileSet) parseFieldList(fl *ast.FieldList) []gen.StructField {
	if fl == nil || fl.NumFields() == 0 {
		return nil
	}
	out := make([]gen.StructField, 0, fl.NumFields())
	for _, field := range fl.List {
		pushstate(fieldName(field))
		fds := fs.getField(field)
		if len(fds) > 0 {
			out = append(out, fds...)
		} else {
			warnf("ignored")
		}
		popstate()
	}
	return out
}

// translate *ast.Field into []gen.StructField
func (fs *FileSet) getField(f *ast.Field) []gen.StructField {
	sf := make([]gen.StructField, 1)
	var extension, flatten bool
	// parse tag; otherwise field name is field tag
	if f.Tag != nil {
		var body string
		if fs.tagName != "" {
			body = reflect.StructTag(strings.Trim(f.Tag.Value, "`")).Get(fs.tagName)
		}
		if body == "" {
			body = reflect.StructTag(strings.Trim(f.Tag.Value, "`")).Get("msg")
		}
		if body == "" {
			body = reflect.StructTag(strings.Trim(f.Tag.Value, "`")).Get("msgpack")
		}
		tags := strings.Split(body, ",")
		if len(tags) >= 2 {
			for _, tag := range tags[1:] {
				switch tag {
				case "extension":
					extension = true
				case "flatten":
					flatten = true
				default:
					// Check for limit=N format
					if strings.HasPrefix(tag, "limit=") {
						limitStr := strings.TrimPrefix(tag, "limit=")
						if limit, err := strconv.ParseUint(limitStr, 10, 32); err == nil {
							sf[0].FieldLimit = uint32(limit)
						} else {
							warnf("invalid limit value in field tag: %s", limitStr)
						}
					}
				}
			}
		}
		// ignore "-" fields
		if tags[0] == "-" {
			return nil
		}
		sf[0].FieldTag = tags[0]
		sf[0].FieldTagParts = tags
		sf[0].RawTag = f.Tag.Value
	}

	ex := fs.parseExpr(f.Type)
	if ex == nil {
		return nil
	}

	// parse field name
	switch len(f.Names) {
	case 0:
		if flatten {
			return fs.getFieldsFromEmbeddedStruct(f.Type)
		} else {
			sf[0].FieldName = embedded(f.Type)
		}
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
				FieldElem: ex.Copy(),
			})
		}
		return sf
	}
	sf[0].FieldElem = ex
	if sf[0].FieldTag == "" {
		sf[0].FieldTag = sf[0].FieldName
		if len(sf[0].FieldTagParts) <= 1 {
			sf[0].FieldTagParts = []string{sf[0].FieldTag}
		} else {
			sf[0].FieldTagParts = append([]string{sf[0].FieldName}, sf[0].FieldTagParts[1:]...)
		}
	}

	// validate extension
	if extension {
		switch ex := ex.(type) {
		case *gen.Ptr:
			if b, ok := ex.Value.(*gen.BaseElem); ok {
				b.Value = gen.Ext
			} else {
				warnf("couldn't cast to extension.")
				return nil
			}
		case *gen.BaseElem:
			ex.Value = gen.Ext
		default:
			warnf("couldn't cast to extension.")
			return nil
		}
	}
	return sf
}

func (fs *FileSet) getFieldsFromEmbeddedStruct(f ast.Expr) []gen.StructField {
	switch f := f.(type) {
	case *ast.Ident:
		s := fs.Specs[f.Name]
		switch s := s.(type) {
		case *ast.StructType:
			return fs.parseFieldList(s.Fields)
		default:
			return nil
		}
	default:
		// other possibilities are disallowed
		return nil
	}
}

// extract embedded field name
//
// so, for a struct like
//
//		type A struct {
//			io.Writer
//	 }
//
// we want "Writer"
func embedded(f ast.Expr) string {
	switch f := f.(type) {
	case *ast.Ident:
		return f.Name
	case *ast.StarExpr:
		return embedded(f.X)
	case *ast.SelectorExpr:
		return f.Sel.Name
	default:
		// other possibilities are disallowed
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
	case *ast.BasicLit:
		return e.Value
	case *ast.IndexExpr:
		// Single type argument: Generic[T]
		return fmt.Sprintf("%s[%s]", stringify(e.X), stringify(e.Index))
	case *ast.IndexListExpr:
		// Multiple type arguments: Generic[A,B,...]
		args := make([]string, 0, len(e.Indices))
		for _, ix := range e.Indices {
			args = append(args, stringify(ix))
		}
		return fmt.Sprintf("%s[%s]", stringify(e.X), strings.Join(args, ","))
	}
	return "<BAD>"
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
	switch e := e.(type) {

	case *ast.MapType:
		switch k := e.Key.(type) {
		case *ast.Ident:
			switch k.Name {
			case "string":
				if in := fs.parseExpr(e.Value); in != nil {
					return &gen.Map{Value: in, AllowBinMaps: fs.AllowBinMaps, AllowMapShims: fs.AllowMapShims, AutoMapShims: fs.AutoMapShims}
				}
				warnf("%s: map keys of type  are not supported\n", stringify(e.Key))
			default:
				if !fs.AllowMapShims && !fs.AllowBinMaps && !fs.AutoMapShims {
					warnf("map keys of type %s are not supported without binary keys or shimming\n", stringify(e.Key))
					return nil
				}
				// Allow for other types, assuming they will be shimmed later.
				key := fs.parseExpr(k)
				// Types that aren't idents are native types and cannot currently be used as map keys.
				switch k := key.(type) {
				case *gen.BaseElem:
					switch k.Value {
					case gen.IDENT:
						if in := fs.parseExpr(e.Value); in != nil {
							return &gen.Map{Value: in, Key: key, AllowBinMaps: fs.AllowBinMaps, AllowMapShims: fs.AllowMapShims, AutoMapShims: fs.AutoMapShims}
						}
						warnf("map keys of type %s are not supported\n", k.TypeName())
						// Exclude types that cannot be used as native map keys.
					case gen.Bytes:
						warnf("map keys of type %s are not supported\n", k.TypeName())
					default:
						if in := fs.parseExpr(e.Value); (fs.AllowBinMaps || (fs.AutoMapShims && gen.CanAutoShim[k.Value])) && in != nil {
							return &gen.Map{Value: in, Key: key, AllowBinMaps: fs.AllowBinMaps, AllowMapShims: fs.AllowMapShims, AutoMapShims: fs.AutoMapShims}
						}
						warnf("map keys of type %s are not supported without binary keys or shimming\n", k.TypeName())
					}
				default:
					warnf("map keys of type %s are not supported\n", k.TypeName())
				}
				return nil
			}
		case *ast.ArrayType:
			if fs.AllowBinMaps {
				key := fs.parseExpr(k)
				if in := fs.parseExpr(e.Value); fs.AllowBinMaps && in != nil {
					return &gen.Map{Value: in, Key: key, AllowBinMaps: fs.AllowBinMaps, AllowMapShims: fs.AllowMapShims, AutoMapShims: fs.AutoMapShims}
				}
			}
			warnf("array map keys (type %s) are not supported without binary keys or shimming\n", stringify(e.Key))
		default:
			warnf("array map key type not supported\n")
		}
		return nil

	case *ast.Ident:
		b := gen.Ident(e.Name)

		// work to resolve this expression
		// can be done later, once we've resolved
		// everything else.
		if b.Value == gen.IDENT {
			if _, ok := fs.Specs[e.Name]; !ok && fs.Aliased[e.Name] == "" {
				// This can be a generic type.
				warnf("possible non-local identifier: %s\n", e.Name)
			}
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

			case *ast.SelectorExpr:
				return &gen.Array{
					Size: stringify(s),
					Els:  els,
				}

			default:
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
		return &gen.Struct{Fields: fs.parseFieldList(e.Fields)}

	case *ast.SelectorExpr:
		return gen.Ident(stringify(e))

	case *ast.InterfaceType:
		// support `interface{}`
		if len(e.Methods.List) == 0 {
			return &gen.BaseElem{Value: gen.Intf}
		}
		return nil

	case *ast.IndexExpr:
		// Treat a generic instantiation like an identifier of the instantiated name.
		// Example: GenericTest2[T] -> "GenericTest2[T]"
		return gen.Ident(stringify(e))

	case *ast.IndexListExpr:
		// Treat a generic instantiation with multiple args similarly.
		return gen.Ident(stringify(e))

	default: // other types not supported
		return nil
	}
}

var Logf func(s string, v ...any)

func infof(s string, v ...any) {
	if Logf != nil {
		pushstate(s)
		Logf("info: "+strings.Join(logctx, ": "), v...)
		popstate()
	}
}

func warnf(s string, v ...any) {
	if Logf != nil {
		pushstate(s)
		Logf("warn: "+strings.Join(logctx, ": "), v...)
		popstate()
	}
}

var logctx []string

// push logging state
func pushstate(s string) {
	logctx = append(logctx, s)
}

// pop logging state
func popstate() {
	logctx = logctx[:len(logctx)-1]
}
