package goon_test

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/shurcooL/go/gists/gist5259939"

	"github.com/shurcooL/go-goon"
)

func foo(bar int) int { return bar * 2 }

func ExampleLarge() {
	fset := token.NewFileSet()
	if file, err := parser.ParseFile(fset, gist5259939.GetThisGoSourceFilepath(), nil, 0); nil == err {
		for _, d := range file.Decls {
			if f, ok := d.(*ast.FuncDecl); ok {
				goon.Dump(f)
				break
			}
		}
	}

	// Output:
	// (*ast.FuncDecl)(&ast.FuncDecl{
	// 	Doc:  (*ast.CommentGroup)(nil),
	// 	Recv: (*ast.FieldList)(nil),
	// 	Name: (*ast.Ident)(&ast.Ident{
	// 		NamePos: (token.Pos)(149),
	// 		Name:    (string)("foo"),
	// 		Obj: (*ast.Object)(&ast.Object{
	// 			Kind: (ast.ObjKind)(5),
	// 			Name: (string)("foo"),
	// 			Decl: (*ast.FuncDecl)(already_shown),
	// 			Data: (interface{})(nil),
	// 			Type: (interface{})(nil),
	// 		}),
	// 	}),
	// 	Type: (*ast.FuncType)(&ast.FuncType{
	// 		Func: (token.Pos)(144),
	// 		Params: (*ast.FieldList)(&ast.FieldList{
	// 			Opening: (token.Pos)(152),
	// 			List: ([]*ast.Field)([]*ast.Field{
	// 				(*ast.Field)(&ast.Field{
	// 					Doc: (*ast.CommentGroup)(nil),
	// 					Names: ([]*ast.Ident)([]*ast.Ident{
	// 						(*ast.Ident)(&ast.Ident{
	// 							NamePos: (token.Pos)(153),
	// 							Name:    (string)("bar"),
	// 							Obj: (*ast.Object)(&ast.Object{
	// 								Kind: (ast.ObjKind)(4),
	// 								Name: (string)("bar"),
	// 								Decl: (*ast.Field)(already_shown),
	// 								Data: (interface{})(nil),
	// 								Type: (interface{})(nil),
	// 							}),
	// 						}),
	// 					}),
	// 					Type: (*ast.Ident)(&ast.Ident{
	// 						NamePos: (token.Pos)(157),
	// 						Name:    (string)("int"),
	// 						Obj:     (*ast.Object)(nil),
	// 					}),
	// 					Tag:     (*ast.BasicLit)(nil),
	// 					Comment: (*ast.CommentGroup)(nil),
	// 				}),
	// 			}),
	// 			Closing: (token.Pos)(160),
	// 		}),
	// 		Results: (*ast.FieldList)(&ast.FieldList{
	// 			Opening: (token.Pos)(0),
	// 			List: ([]*ast.Field)([]*ast.Field{
	// 				(*ast.Field)(&ast.Field{
	// 					Doc:   (*ast.CommentGroup)(nil),
	// 					Names: ([]*ast.Ident)(nil),
	// 					Type: (*ast.Ident)(&ast.Ident{
	// 						NamePos: (token.Pos)(162),
	// 						Name:    (string)("int"),
	// 						Obj:     (*ast.Object)(nil),
	// 					}),
	// 					Tag:     (*ast.BasicLit)(nil),
	// 					Comment: (*ast.CommentGroup)(nil),
	// 				}),
	// 			}),
	// 			Closing: (token.Pos)(0),
	// 		}),
	// 	}),
	// 	Body: (*ast.BlockStmt)(&ast.BlockStmt{
	// 		Lbrace: (token.Pos)(166),
	// 		List: ([]ast.Stmt)([]ast.Stmt{
	// 			(*ast.ReturnStmt)(&ast.ReturnStmt{
	// 				Return: (token.Pos)(168),
	// 				Results: ([]ast.Expr)([]ast.Expr{
	// 					(*ast.BinaryExpr)(&ast.BinaryExpr{
	// 						X: (*ast.Ident)(&ast.Ident{
	// 							NamePos: (token.Pos)(175),
	// 							Name:    (string)("bar"),
	// 							Obj: (*ast.Object)(&ast.Object{
	// 								Kind: (ast.ObjKind)(4),
	// 								Name: (string)("bar"),
	// 								Decl: (*ast.Field)(&ast.Field{
	// 									Doc: (*ast.CommentGroup)(nil),
	// 									Names: ([]*ast.Ident)([]*ast.Ident{
	// 										(*ast.Ident)(&ast.Ident{
	// 											NamePos: (token.Pos)(153),
	// 											Name:    (string)("bar"),
	// 											Obj:     (*ast.Object)(already_shown),
	// 										}),
	// 									}),
	// 									Type: (*ast.Ident)(&ast.Ident{
	// 										NamePos: (token.Pos)(157),
	// 										Name:    (string)("int"),
	// 										Obj:     (*ast.Object)(nil),
	// 									}),
	// 									Tag:     (*ast.BasicLit)(nil),
	// 									Comment: (*ast.CommentGroup)(nil),
	// 								}),
	// 								Data: (interface{})(nil),
	// 								Type: (interface{})(nil),
	// 							}),
	// 						}),
	// 						OpPos: (token.Pos)(179),
	// 						Op:    (token.Token)(14),
	// 						Y: (*ast.BasicLit)(&ast.BasicLit{
	// 							ValuePos: (token.Pos)(181),
	// 							Kind:     (token.Token)(5),
	// 							Value:    (string)("2"),
	// 						}),
	// 					}),
	// 				}),
	// 			}),
	// 		}),
	// 		Rbrace: (token.Pos)(183),
	// 	}),
	// })
	//
}
