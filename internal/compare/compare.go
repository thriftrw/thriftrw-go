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

// Modules looks for removed methods and added required fields.
func (p *Pass) Modules(from, to *compile.Module) {
	for name, fromService := range from.Services {
		p.service(fromService, to.Services[name])
	}
	// p.services(from, to)
	p.checkRequiredFields(from, to)
}

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

func (p *Pass) checkRequiredFields(fromModule, toModule *compile.Module) {
	for n, spec := range toModule.Types {
		fromSpec, ok := fromModule.Types[n]
		if !ok {
			// This is a new Type, which is backwards compatible.
			continue
		}
		if s, ok := spec.(*compile.StructSpec); ok {
			// Match on Type names. Here we hit a limitation, that if someone
			// renames the struct and then adds a new field, we don't really have
			// a good way of tracking it.
			if fromStructSpec, ok := fromSpec.(*compile.StructSpec); ok {
				p.structSpecs(fromStructSpec, s, filepath.Base(fromModule.ThriftPath))
			}
		}
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
	for f := range from.Functions {
		if _, ok := to.Functions[f]; !ok {
			p.Report(Diagnostic{
				File:    filepath.Base(from.File),
				Message: fmt.Sprintf("removing method %s in service %s is not backwards compatible", f, from.Name),
			})
		}
	}
}

// Services compares two service definitions.
// func (p *Pass) services(fromModule, toModule *compile.Module) {
// 	for n, fromService := range fromModule.Services {
// 		toServ, ok := toModule.Services[n]
// 		if !ok {
// 			// Service was deleted, which is not backwards compatible.
// 			p.Report(Diagnostic{
// 				File:    filepath.Base(fromModule.ThriftPath), // toModule could have been deleted.
// 				Message: fmt.Sprintf("deleting service %s is not backwards compatible", n),
// 			})
// 			// Do not need to check its functions since it was deleted.
// 			continue
// 		}
// 		for f := range fromService.Functions {
// 			if _, ok := toServ.Functions[f]; !ok {
// 				p.Report(Diagnostic{
// 					File:    filepath.Base(fromModule.ThriftPath),
// 					Message: fmt.Sprintf("removing method %s in service %s is not backwards compatible", f, n),
// 				})
// 			}
// 		}
// 	}
// }
