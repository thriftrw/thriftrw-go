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
	"go/format"
	"go/token"
	"io"
)

// Generator tracks code generation state as we generate the output.
type Generator struct {
	importer

	decls []ast.Decl

	// TODO use something to group related decls together

	// TODO(abg) We will keep track of needed map/list/set types and their
	// to/from value implementations here
}

// NewGenerator sets up a new generator for Go code.
func NewGenerator() *Generator {
	return &Generator{importer: newImporter()}
}

// TODO mutliple modules

func (g *Generator) Write(w io.Writer, fs *token.FileSet) error {
	// TODO newlines between decls
	// TODO constants first, types next, and functions after that
	// TODO sorting

	decls := make([]ast.Decl, 0, 1+len(g.decls))
	importDecl := g.importDecl()
	if importDecl != nil {
		decls = append(decls, importDecl)
	}
	decls = append(decls, g.decls...)

	file := &ast.File{
		Decls: decls,
		Name:  ast.NewIdent("todo"), // TODO
	}
	return format.Node(w, fs, file)
}

// appendDecl appends a new declaration to the generator.
func (g *Generator) appendDecl(decl ast.Decl) {
	g.decls = append(g.decls, decl)
}

// defineType defines a new type for the generator.
//
//  type $name $def
func (g *Generator) defineType(name string, def ast.Expr) {
	g.appendDecl(
		&ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{Name: ast.NewIdent(name), Type: def},
			},
		},
	)
}

// declareConstant declares a new constant with the given type and name.
//
//  const $name $typ = $value
func (g *Generator) declareConstant(name string, typ ast.Expr, value ast.Expr) {
	g.appendDecl(
		&ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{&ast.ValueSpec{
				Type:   typ,
				Names:  []*ast.Ident{ast.NewIdent(name)},
				Values: []ast.Expr{value},
			}},
		},
	)
}
