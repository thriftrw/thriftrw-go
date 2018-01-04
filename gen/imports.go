// Copyright (c) 2018 Uber Technologies, Inc.
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
	"path/filepath"
	"strings"
	"unicode"
)

// importer is responsible for managing imports for the code generator and
// ensuring that we don't end up with naming conflicts in imports.
type importer struct {
	ns      Namespace
	imports map[string]*ast.ImportSpec
}

// newImporter builds a new importer.
func newImporter(ns Namespace) importer {
	return importer{
		ns:      ns,
		imports: make(map[string]*ast.ImportSpec),
	}
}

// AddImportSpec allows adding existing import specs to the importer.
//
// An error is returned if there's a naming conflict.
func (i importer) AddImportSpec(spec *ast.ImportSpec) error {
	path := spec.Path.Value
	name := filepath.Base(path)
	if spec.Name != nil {
		name = spec.Name.Name
	}

	if err := i.ns.Reserve(name); err != nil {
		return err
	}

	i.imports[path] = spec
	return nil
}

// Import ensures that the generated module has the given module imported and
// returns the name that should be used by the generated code to reference items
// defined in the module.
func (i importer) Import(path string) string {
	if imp, ok := i.imports[path]; ok {
		if imp.Name != nil {
			return imp.Name.Name
		}
		return filepath.Base(path)
	}

	// Find a name, preferring the base name
	// TODO what if the package name is not the base name?
	baseName := filepath.Base(path)
	name := i.ns.NewName(sanitizeImportName(baseName))
	astImport := &ast.ImportSpec{Path: stringLiteral(path)}
	if name != baseName {
		astImport.Name = ast.NewIdent(name)
	}

	i.imports[path] = astImport
	return name
}

func sanitizeImportName(s string) string {
	// special handling for common "foo-go" pattern
	if strings.HasSuffix(s, "-go") {
		s = s[:len(s)-3]
	}

	return strings.Map(func(c rune) rune {
		switch {
		case unicode.IsLetter(c), unicode.IsDigit(c):
			return c
		default:
			return '_'
		}
	}, s)
}

// importDecl builds an import declation from the given list of imports.
func (i importer) importDecl() ast.Decl {
	imports := i.imports
	if imports == nil || len(imports) == 0 {
		return nil
	}

	specs := make([]ast.Spec, 0, len(imports))
	for _, iname := range sortStringKeys(imports) {
		imp := imports[iname]
		specs = append(specs, imp)
	}

	decl := &ast.GenDecl{Tok: token.IMPORT, Specs: specs}
	if len(specs) > 1 {
		// Just need a non-zero value for Lparen to get the parens added.
		decl.Lparen = token.Pos(1)
	}

	return decl
}
