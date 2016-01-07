// Copyright (c) 2015 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package gen

import (
	"go/ast"
	"go/token"

	"github.com/uber/thriftrw-go/compile"
)

func (g *Generator) enum(spec *compile.EnumSpec) {
	enumName := typeDeclName(spec)

	// type $enumName int32
	g.defineType(enumName, ast.NewIdent("int32"))

	specs := make([]ast.Spec, len(spec.Items))
	for i, item := range spec.Items {
		specs[i] = g.enumItem(spec, item)
	}

	g.appendDecl(
		// const (
		// 	$itemName $enumName = $value
		// )
		&ast.GenDecl{
			Tok: token.CONST,
			// just need a non-zero lparen position for the formatter to put
			// parens in the correct places
			Lparen: token.Pos(1),
			Specs:  specs,
		},
	)

	g.enumToWire(spec)
}

func (g *Generator) enumItem(enum *compile.EnumSpec, item compile.EnumItem) ast.Spec {
	// TODO convert all caps
	enumName := typeDeclName(enum)
	name := enumName + capitalize(item.Name)

	// $itemName $enumName = $enumName($itemValue)
	return &ast.ValueSpec{
		Type:   typeReference(enum, false),
		Names:  []*ast.Ident{ast.NewIdent(name)},
		Values: []ast.Expr{intLiteral(int64(item.Value))},
	}
}

func (g *Generator) enumToWire(enum *compile.EnumSpec) {
	wire := g.Import("github.com/uber/thriftrw-go/wire")
	newValueI32 := wire.Select("NewValueI32")

	decl := &ast.FuncDecl{
		Recv: fieldList(fields(typeReference(enum, false), "x")),
		Name: ast.NewIdent("ToWire"),
		Type: &ast.FuncType{
			Params:  fieldList(),
			Results: fieldList(fields(wire.Select("Value").E)),
		},
		Body: block(
			&ast.ReturnStmt{
				Results: []ast.Expr{
					newValueI32.Call(ident("int32").Call(ident("x"))).E,
				},
			},
		),
	}

	g.appendDecl(decl)
}

func block(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: stmts}
}

func fieldList(fs ...*ast.Field) *ast.FieldList {
	return &ast.FieldList{List: fs}
}

func fields(typ ast.Expr, fields ...string) *ast.Field {
	names := make([]*ast.Ident, len(fields))
	for i, f := range fields {
		names[i] = ast.NewIdent(f)
	}

	return &ast.Field{Names: names, Type: typ}
}
