package parse

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestNilable(t *testing.T) {
	const src = `package p
type Str string
type Named Str
type St struct{ A int }
type Alias = St
type Bytes []byte
type Arr [4]byte
type Ptr *int
type unexported string
type Cyc1 Cyc2
type Cyc2 Cyc1
`
	f, err := parser.ParseFile(token.NewFileSet(), "p.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	fs := &FileSet{declared: make(map[string]ast.Expr)}
	fs.addDeclared(f)
	ast.FileExports(f)
	for name, want := range map[string]bool{
		"Str": false, "Named": false, "St": false, "Alias": false, "Arr": false,
		"unexported": false, "int": false, "string": false,
		"Bytes": true, "Ptr": true, "any": true, "error": true,
		"Unknown": true, "ext.Bytes": true, "Cyc1": true,
	} {
		if got := fs.nilable(ast.NewIdent(name)); got != want {
			t.Errorf("nilable(%s) = %v, want %v", name, got, want)
		}
	}
}
