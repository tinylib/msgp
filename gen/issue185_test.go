package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
	"text/template"
)

const issue185Tpl = `
package issue185

//go:generate msgp

type Test1 struct {
	Foo string
	Bar string
	{{ if .Extra }}Baz []string{{ end }}
	Qux string
}

type Test2 struct {
	Foo string
	Bar string
	Baz string
}
`

// Ensure that consistent identifiers are generated on a per-method basis by msgp.
//
// Also eusnre that no duplicate identifiers appear in a method.
//
// structs are currently processed alphabetically by msgp. this test relies on
// that property.
//
func TestIssue185Idents(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "msgp-")
	if err != nil {
		t.Fatalf("could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tfile := filepath.Join(tempDir, "issue185.go")
	genFile := filepath.Join(tempDir, "issue185_gen.go")

	tpl := template.Must(template.New("").Parse(issue185Tpl))
	tplData := struct{ Extra bool }{}

	tplData.Extra = false
	if err := goGenerateTpl(tempDir, tfile, tpl, tplData); err != nil {
		t.Fatalf("could not generate code: %v", err)
	}

	// generate the code, extract the generated variable names, mapped to function name
	varsBefore, err := extractVars(genFile)
	if err != nil {
		t.Fatalf("could not extract before vars: %v", err)
	}

	// regenerate the code with extra field(s), extract the generated variable
	// names, mapped to function name
	tplData.Extra = true
	if err := goGenerateTpl(tempDir, tfile, tpl, tplData); err != nil {
		t.Fatalf("could not generate code: %v", err)
	}
	varsAfter, err := extractVars(genFile)
	if err != nil {
		t.Fatalf("could not extract after vars: %v", err)
	}

	for _, fn := range []string{"Test2.DecodeMsg", "Test2.EncodeMsg", "Test2.Msgsize", "Test2.MarshalMsg", "Test2.UnmarshalMsg"} {
		if !reflect.DeepEqual(varsBefore.Value(fn), varsAfter.Value(fn)) {
			t.Fatalf("vars unexpectedly changed for %s", fn)
		}
	}

	for fn, vars := range varsAfter {
		sort.Strings(vars)
		for i := 0; i < len(vars)-1; i++ {
			if vars[i] == vars[i+1] {
				t.Fatalf("duplicate var %s in function %s", vars[i], fn)
			}
		}
	}
}

type varVisitor struct {
	vars []string
	fset *token.FileSet
}

func (v *varVisitor) Visit(node ast.Node) (w ast.Visitor) {
	gen, ok := node.(*ast.GenDecl)
	if !ok {
		return v
	}
	for _, spec := range gen.Specs {
		if vspec, ok := spec.(*ast.ValueSpec); ok {
			for _, n := range vspec.Names {
				v.vars = append(v.vars, n.Name)
			}
		}
	}
	return v
}

type extractedVars map[string][]string

func (e extractedVars) Value(key string) []string {
	if v, ok := e[key]; ok {
		return v
	}
	panic(fmt.Errorf("unknown key %s", key))
}

func extractVars(file string) (extractedVars, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, file, nil, 0)
	if err != nil {
		return nil, err
	}

	vars := make(map[string][]string)
	for _, d := range f.Decls {
		switch d := d.(type) {
		case *ast.FuncDecl:
			sn := ""
			switch rt := d.Recv.List[0].Type.(type) {
			case *ast.Ident:
				sn = rt.Name
			case *ast.StarExpr:
				sn = rt.X.(*ast.Ident).Name
			default:
				panic("unknown receiver type")
			}

			key := fmt.Sprintf("%s.%s", sn, d.Name.Name)
			vis := &varVisitor{fset: fset}
			ast.Walk(vis, d.Body)
			vars[key] = vis.vars
		}
	}
	return vars, nil
}

func goGenerateTpl(cwd, tfile string, tpl *template.Template, tplData interface{}) error {
	outf, err := os.OpenFile(tfile, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outf.Close()

	if err := tpl.Execute(outf, tplData); err != nil {
		return err
	}

	cmd := exec.Command("go", "generate")
	cmd.Dir = cwd
	return cmd.Run()
}
