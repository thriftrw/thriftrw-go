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

package compile

// Module represents a compiled Thrift module. It contains all information
// about all known types, constants, services, and includes from the Thrift
// file.
//
// ThriftPath is the absolute path to the Thrift file from which this module
// was compiled. All includes made by the Thrift file are relative to that
// path.
//
// The module name is usually just the basename of the ThriftPath.
type Module struct {
	Name       string
	ThriftPath string

	// Mapping from the /Thrift name/ to the compiled representation of
	// different definitions.

	Includes  map[string]*IncludedModule
	Constants map[string]*Constant
	Types     map[string]TypeSpec
	Services  map[string]*ServiceSpec

	Raw []byte // The raw IDL input.
}

// GetName for Module
func (m *Module) GetName() string {
	return m.Name
}

// LookupType for Module.
func (m *Module) LookupType(name string) (TypeSpec, error) {
	if t, ok := m.Types[name]; ok {
		return t, nil
	}

	return nil, lookupError{Name: name}
}

// LookupConstant for Module.
func (m *Module) LookupConstant(name string) (*Constant, error) {
	if c, ok := m.Constants[name]; ok {
		return c, nil
	}

	return nil, lookupError{Name: name}
}

// LookupService for Module.
func (m *Module) LookupService(name string) (*ServiceSpec, error) {
	if s, ok := m.Services[name]; ok {
		return s, nil
	}

	return nil, lookupError{Name: name}
}

// LookupInclude for Module.
func (m *Module) LookupInclude(name string) (Scope, error) {
	if s, ok := m.Includes[name]; ok {
		return s.Module, nil
	}

	return nil, lookupError{Name: name}
}

// Walk the module tree starting at the given module. This module and all its
// direct and transitive dependencies will be visited exactly once in an
// unspecified order. The walk will stop on the first error returned by `f`.
func (m *Module) Walk(f func(*Module) error) error {
	visited := make(map[string]struct{})

	toVisit := make([]*Module, 0, 100)
	toVisit = append(toVisit, m)

	for len(toVisit) > 0 {
		m := toVisit[0]
		toVisit = toVisit[1:]

		if _, ok := visited[m.ThriftPath]; ok {
			continue
		}

		visited[m.ThriftPath] = struct{}{}
		for _, inc := range m.Includes {
			toVisit = append(toVisit, inc.Module)
		}

		if err := f(m); err != nil {
			return err
		}
	}

	return nil
}

// IncludedModule represents an included module in the Thrift file.
//
// The name of the IncludedModule is the name under which the module is
// exposed in the Thrift file which included it. This is usually the same as
// the Module name except when our custom include-as syntax is used.
type IncludedModule struct {
	Name   string
	Module *Module
}
