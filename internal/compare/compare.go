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
	File    string // File where error was discovered
	Message string // Message contains error message.
}

func (d *Diagnostic) String() string {
	return fmt.Sprintf("%s:%s", d.File, d.Message)
}

// Pass provides all reported errors.
type Pass struct {
	lints []Diagnostic
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

// Modules looks for removed methods and added required fields.
func (p *Pass) Modules(from, to *compile.Module) {
	for name, fromService := range from.Services {
		p.service(fromService, to.Services[name])
	}

	for n, fromType := range from.Types {
		p.typ(fromType, to.Types[n], filepath.Base(from.ThriftPath))
	}
}

func (p *Pass) typ(from, to compile.TypeSpec, file string) {
	switch f := from.(type) {
	case *compile.StructSpec:
		t, ok := to.(*compile.StructSpec)
		if !ok {
			// This is a new Type, which is backwards compatible.
			return
		}
		p.structSpecs(f, t, file)
	}
}

// StructSpecs compares two structs defined in a Thrift file.
func (p *Pass) structSpecs(from, to *compile.StructSpec, file string) {
	fields := make(map[int16]*compile.FieldSpec, len(from.Fields))
	// Assume that these two should be compared.
	for _, f := range from.Fields {
		// Capture state of all fields here.
		fields[f.ID] = f
	}

	for _, toField := range to.Fields {
		if fromField, ok := fields[toField.ID]; ok {
			fromRequired := fromField.Required
			toRequired := toField.Required
			if !fromRequired && toRequired {
				p.Report(Diagnostic{
					File: file,
					Message: fmt.Sprintf(
						"changing an optional field %s in %s to required is not backwards compatible",
						toField.ThriftName(), to.ThriftName()),
				})
			}
		} else if toField.Required {
			p.Report(Diagnostic{
				File: file,
				Message: fmt.Sprintf("adding a required field %s to %s is not backwards compatible",
					toField.ThriftName(), to.ThriftName()),
			})
		}
	}
}

func (p *Pass) service(from, to *compile.ServiceSpec) {
	if to == nil {
		// Service was deleted, which is not backwards compatible.
		p.Report(Diagnostic{
			File:    filepath.Base(from.File), // toModule could have been deleted.
			Message: fmt.Sprintf("deleting service %s is not backwards compatible", from.Name),
		})

		return
	}
	for n := range from.Functions {
		p.function(to.Functions[n], n, filepath.Base(from.File), from.Name)
	}
}

func (p *Pass) function(to *compile.FunctionSpec, fn string, file string, service string) {
	if to == nil {
		p.Report(Diagnostic{
			File:    file,
			Message: fmt.Sprintf("removing method %s in service %s is not backwards compatible", fn, service),
		})
	}
}
