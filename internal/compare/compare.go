// Copyright (c) 2021 Uber Technologies, Inc.
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

package compare

import (
	"fmt"
	"path/filepath"
	"strings"

	"go.uber.org/thriftrw/compile"
)

// Diagnostic is a message associated with an error and a file name.
type Diagnostic struct {
	FilePath string // FilePath where error was discovered.
	Message  string // Message contains error message.
}

func (d *Diagnostic) String() string {
	return fmt.Sprintf("%s:%s", d.FilePath, d.Message)
}

// Pass provides all reported errors.
type Pass struct {
	lints  []Diagnostic
	GitDir string
}

// Report reports an error.
func (p *Pass) Report(d Diagnostic) {
	p.lints = append(p.lints, d)
}

// Lints returns all errors.
func (p *Pass) Lints() []Diagnostic {
	return p.lints
}

func (p *Pass) String() string {
	var b strings.Builder
	for _, l := range p.lints {
		_, _ = fmt.Fprintf(&b, "%s\n", l.String())
	}

	return b.String()
}

// CompareModules looks for removed methods and added required fields.
func (p *Pass) CompareModules(from, to *compile.Module) {
	for name, fromService := range from.Services {
		p.service(fromService, to.Services[name])
	}

	file := p.getRelativePath(from.ThriftPath)
	for n, fromType := range from.Types {
		p.typ(fromType, to.Types[n], file)
	}
}

func (p *Pass) typ(from, to compile.TypeSpec, file string) {
	if f, ok := from.(*compile.StructSpec); ok {
		t, ok := to.(*compile.StructSpec)
		if !ok {
			// A struct was deleted which is ok if it's unused or
			// it's usage was also removed.
			return
		}
		p.structSpecs(f, t, file)
	}
}

func (p *Pass) requiredField(fromField, toField *compile.FieldSpec, to *compile.StructSpec, file string) {
	fromRequired := fromField.Required
	if !fromRequired && toField.Required {
		p.Report(Diagnostic{
			FilePath: file,
			Message: fmt.Sprintf(
				"changing an optional field %q in %q to required",
				toField.ThriftName(), to.ThriftName()),
		})
	}
}

// StructSpecs compares two structs defined in a Thrift file.
func (p *Pass) structSpecs(from, to *compile.StructSpec, file string) {
	fields := make(map[int16]*compile.FieldSpec, len(from.Fields))
	// Assume that these two should be compared.
	for _, f := range from.Fields {
		fields[f.ID] = f
	}
	for _, toField := range to.Fields {
		if fromField, ok := fields[toField.ID]; ok {
			p.requiredField(fromField, toField, to, file)
		} else if toField.Required {
			p.Report(Diagnostic{
				FilePath: file,
				Message: fmt.Sprintf("adding a required field %q to %q",
					toField.ThriftName(), to.ThriftName()),
			})
		}
	}
}

func (p *Pass) service(from, to *compile.ServiceSpec) {
	if to == nil {
		// Service was deleted, which is not backwards compatible.
		p.Report(Diagnostic{
			FilePath: filepath.Base(from.File), // toModule could have been deleted.
			Message:  fmt.Sprintf("deleting service %q", from.Name),
		})

		return
	}
	file := p.getRelativePath(from.File)
	for n := range from.Functions {
		p.function(to.Functions[n], n, file, from.Name)
	}
}

// getRelativePath returns a relative path to a file or
// fallbacks to file name for cases when it was deleted.
func (p *Pass) getRelativePath(filePath string) string {
	if file, err := filepath.Rel(p.GitDir, filePath); err == nil {
		return file
	}
	// If a file was deleted, then we will not be able to
	// find a relative path to it.
	return filepath.Base(filePath)
}

func (p *Pass) function(to *compile.FunctionSpec, fn string, path string, service string) {
	file := p.getRelativePath(path)
	if to == nil {
		p.Report(Diagnostic{
			FilePath: file,
			Message:  fmt.Sprintf("removing method %q in service %q", fn, service),
		})
	}
}
