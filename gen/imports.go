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
	"fmt"
	"go/ast"
	"go/token"
	"path/filepath"
)

// Expr helps build complex expressions.
type Expr struct {
	E ast.Expr
}

func ident(name string) Expr {
	return Expr{ast.NewIdent(name)}
}

// Select can be used to select an item on an imported module.
func (e Expr) Select(name string) Expr {
	return Expr{&ast.SelectorExpr{X: e.E, Sel: ast.NewIdent(name)}}
}

func (e Expr) call(args ...ast.Expr) Expr {
	return Expr{&ast.CallExpr{Fun: e.E, Args: args}}
}

// Call TODO
func (e Expr) Call(args ...Expr) Expr {
	exprs := make([]ast.Expr, len(args))
	for i, expr := range args {
		exprs[i] = expr.E
	}
	return e.call(exprs...)
}

type importer struct {
	usedImpNames map[string]struct{}
	imps         map[string]*ast.ImportSpec
}

func newImporter() importer {
	return importer{
		usedImpNames: make(map[string]struct{}),
		imps:         make(map[string]*ast.ImportSpec),
	}
}

// Import ensures that the generated module has the given module imported and
// returns the name that should be used by the generated code to reference items
// defined in the module.
func (i importer) Import(path string) Expr {
	if imp, ok := i.imps[path]; ok {
		if imp.Name != nil {
			return Expr{imp.Name}
		}
		return ident(filepath.Base(path))
	}

	baseName := filepath.Base(path)
	name := baseName

	// If there's a name collision, use the format _$baseName$counter to find
	// a non-conflicting name.
	if _, conflict := i.usedImpNames[name]; conflict {
		counter := 0
		for conflict {
			counter++
			name = fmt.Sprintf("_%s%d", baseName, counter)
			_, conflict = i.usedImpNames[name]
		}
	}

	astName := ast.NewIdent(name)
	astImport := &ast.ImportSpec{Path: stringLiteral(path)}
	if name != baseName {
		astImport.Name = astName
	}

	i.usedImpNames[name] = struct{}{}
	i.imps[path] = astImport
	return Expr{astName}
}

// importDecl builds an import declation from the given list of imports.
func (i importer) importDecl() ast.Decl {
	imps := i.imps
	if imps == nil || len(imps) == 0 {
		return nil
	}

	specs := make([]ast.Spec, 0, len(imps))
	for _, imp := range imps {
		specs = append(specs, imp)
	}

	decl := &ast.GenDecl{Tok: token.IMPORT, Specs: specs}
	if len(specs) > 1 {
		// Just need a non-zero value for Lparen to get the parens added.
		decl.Lparen = token.Pos(1)
	}

	return decl
}
