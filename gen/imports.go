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

// importer is responsible for managing imports for the code generator and
// ensuring that we don't end up with naming conflicts in imports.
type importer struct {
	usedNames map[string]struct{}
	imports   map[string]*ast.ImportSpec
}

// newImporter builds a new importer.
func newImporter() importer {
	return importer{
		usedNames: make(map[string]struct{}),
		imports:   make(map[string]*ast.ImportSpec),
	}
}

// addImportSpec adds an ImportSpec to the importer.
//
// No conflict resolution is performed.
func (i importer) addImportSpec(spec *ast.ImportSpec) {
	path := spec.Path.Value
	name := filepath.Base(path)
	if spec.Name != nil {
		name = spec.Name.Name
	}

	i.usedNames[name] = struct{}{}
	i.imports[path] = spec
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

	baseName := filepath.Base(path)
	name := baseName

	// If there's a name collision, use the format _$baseName$counter to find
	// a non-conflicting name.
	if _, conflict := i.usedNames[name]; conflict {
		for counter := 0; conflict; counter++ {
			name = fmt.Sprintf("_%s%d", baseName, counter)
			_, conflict = i.usedNames[name]
		}
	}

	astImport := &ast.ImportSpec{Path: stringLiteral(path)}
	if name != baseName {
		astImport.Name = ast.NewIdent(name)
	}
	i.addImportSpec(astImport)

	return name
}

// importDecl builds an import declation from the given list of imports.
func (i importer) importDecl() ast.Decl {
	imports := i.imports
	if imports == nil || len(imports) == 0 {
		return nil
	}

	specs := make([]ast.Spec, 0, len(imports))
	for _, imp := range imports {
		specs = append(specs, imp)
	}

	decl := &ast.GenDecl{Tok: token.IMPORT, Specs: specs}
	if len(specs) > 1 {
		// Just need a non-zero value for Lparen to get the parens added.
		decl.Lparen = token.Pos(1)
	}

	return decl
}
