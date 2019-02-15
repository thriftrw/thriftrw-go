// Copyright (c) 2019 Uber Technologies, Inc.
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

package plugin

import (
	"bytes"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"sort"
	"text/template"

	"go.uber.org/thriftrw/internal/goast"
	"go.uber.org/thriftrw/plugin/api"

	"golang.org/x/tools/go/ast/astutil"
)

// TemplateOption provides optional arguments to GoFileFromTemplate.
type TemplateOption struct {
	apply func(*goFileGenerator)
}

// TemplateFunc is a TemplateOption that makes a function available in the
// template.
//
// The function may be anything accepted by text/template.
//
// 	GoFileFromTemplate(
// 		filename,
// 		`package <lower "HELLO">`,
// 		TemplateFunc("lower", strings.ToLower),
// 	)
func TemplateFunc(name string, f interface{}) TemplateOption {
	return TemplateOption{apply: func(t *goFileGenerator) {
		t.templateFuncs[name] = f
	}}
}

// GoFileImportPath is a TemplateOption that specifies the intended absolute
// import path for the file being generated.
//
// 	GoFileFromTemplate(
// 		filename,
// 		mytemplate,
// 		GoFileImportPath("go.uber.org/thriftrw/myservice"),
// 	)
//
// If specified, this changes the behavior of the `formatType` template
// function to NOT import this package and instead use the types directly
// since they are available in the same package.
func GoFileImportPath(path string) TemplateOption {
	return TemplateOption{apply: func(t *goFileGenerator) {
		t.importPath = path
	}}
}

// goFileGenerator generates a single Go file.
type goFileGenerator struct {
	importPath    string
	templateFuncs template.FuncMap

	// Names of known globals. All global variables share this namespace.
	globals map[string]struct{}

	// import path -> import name
	imports map[string]string
}

func newGoFileGenerator(opts []TemplateOption) *goFileGenerator {
	t := goFileGenerator{
		templateFuncs: make(template.FuncMap),
		globals:       make(map[string]struct{}),
		imports:       make(map[string]string),
	}
	for _, opt := range opts {
		opt.apply(&t)
	}
	return &t
}

func (g *goFileGenerator) isGlobalTaken(name string) bool {
	_, taken := g.globals[name]
	return taken || goast.IsReservedKeyword(name)
}

// Import the given import path and return the imported name for this package.
func (g *goFileGenerator) Import(path string) string {
	if name, ok := g.imports[path]; ok {
		return name
	}

	name := goast.DeterminePackageName(path)

	// Find an import name that does not conflict with any known globals.
	importedName := name
	for i := 2; ; i++ {
		if !g.isGlobalTaken(importedName) {
			break
		}
		importedName = fmt.Sprintf("%s%d", name, i)
	}

	g.imports[path] = importedName
	g.globals[importedName] = struct{}{}
	return importedName
}

// FormatType formats the given api.Type into a Go type, importing packages
// necessary to reference this type.
func (g *goFileGenerator) FormatType(t *api.Type) (string, error) {
	switch {
	case t.SimpleType != nil:
		switch *t.SimpleType {
		case api.SimpleTypeBool:
			return "bool", nil
		case api.SimpleTypeByte:
			return "byte", nil
		case api.SimpleTypeInt8:
			return "int8", nil
		case api.SimpleTypeInt16:
			return "int16", nil
		case api.SimpleTypeInt32:
			return "int32", nil
		case api.SimpleTypeInt64:
			return "int64", nil
		case api.SimpleTypeFloat64:
			return "float64", nil
		case api.SimpleTypeString:
			return "string", nil
		case api.SimpleTypeStructEmpty:
			return "struct{}", nil
		default:
			return "", fmt.Errorf("unknown simple type: %v", *t.SimpleType)
		}
	case t.SliceType != nil:
		v, err := g.FormatType(t.SliceType)
		return "[]" + v, err
	case t.KeyValueSliceType != nil:
		k, err := g.FormatType(t.KeyValueSliceType.Left)
		if err != nil {
			return "", err
		}

		v, err := g.FormatType(t.KeyValueSliceType.Right)
		return fmt.Sprintf("[]struct{Key %v; Value %v}", k, v), err
	case t.MapType != nil:
		k, err := g.FormatType(t.MapType.Left)
		if err != nil {
			return "", err
		}

		v, err := g.FormatType(t.MapType.Right)
		return fmt.Sprintf("map[%v]%v", k, v), err
	case t.ReferenceType != nil:
		if g.importPath == t.ReferenceType.ImportPath {
			// Target is in the same package. No need to import.
			return t.ReferenceType.Name, nil
		}

		importName := g.Import(t.ReferenceType.ImportPath)
		return importName + "." + t.ReferenceType.Name, nil
	case t.PointerType != nil:
		v, err := g.FormatType(t.PointerType)
		return "*" + v, err
	default:
		return "", fmt.Errorf("unknown type: %v", t)
	}
}

// Generates a Go file with the given name using the provided template and
// template data.
func (g *goFileGenerator) Generate(filename, tmpl string, data interface{}) ([]byte, error) {
	funcs := template.FuncMap{
		"import":     g.Import,
		"formatType": g.FormatType,
	}
	for k, v := range g.templateFuncs {
		funcs[k] = v
	}

	t, err := template.New(filename).Delims("<", ">").Funcs(funcs).Parse(tmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %q: %v", filename, err)
	}

	var buff bytes.Buffer
	if err := t.Execute(&buff, data); err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, filename, buff.Bytes(), parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated code: %v:\n%s", err, buff.String())
	}

	if len(f.Imports) > 0 {
		return nil, fmt.Errorf(
			"plain imports are not allowed with GoFileFromTemplate: use the import function")
	}

	importPaths := make([]string, 0, len(g.imports))
	for path := range g.imports {
		importPaths = append(importPaths, path)
	}
	sort.Strings(importPaths)
	for _, path := range importPaths {
		astutil.AddNamedImport(fset, f, g.imports[path], path)
	}

	buff = bytes.Buffer{}
	if err := format.Node(&buff, fset, f); err != nil {
		return nil, err // TODO wrap error
	}

	return buff.Bytes(), nil
}

// GoFileFromTemplate generates a Go file from the given template and template
// data.
//
// The templating system follows the text/template templating format but with "<"
// and ">" as the delimiters.
//
// The following functions are provided inside the template:
//
// import: Use this if you need to import other packages. Import may be called
// anywhere in the template with an import path to ensure that that package is
// imported in the generated file. The import is automatically converted into a
// named import if there's a conflict with another import. This returns the
// imported name of the package. Use the return value of this function to
// reference the imported package.
//
// 	<$wire := import "go.uber.org/thriftrw/wire">
// 	var value <$wire>.Value
//
// formatType: Formats an api.Type into a Go type representation, automatically
// importing packages needed for type references. By default, this imports all
// packages referenced in the api.Type. If the GoFileImportPath option is
// specified, types from that package will not be imported and instead, will be
// assumed to be available in the same package.
//
// 	var value <formatType .Type>
//
// More functions may be added to the template using the TemplateFunc template
// option. If the name of a TemplateFunc conflicts with a pre-defined function,
// the TemplateFunc takes precedence.
//
// Code generated by this is automatically reformatted to comply with gofmt.
func GoFileFromTemplate(filename, tmpl string, data interface{}, opts ...TemplateOption) ([]byte, error) {
	return newGoFileGenerator(opts).Generate(filename, tmpl, data)
}
