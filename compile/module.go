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

package compile

import "fmt"

// Module represents a compiled Thrift module. It contains all information
// about all known types, constants, services, and includes from the Thrift
// file.
//
// ThriftPath is the path to the Thrift file from which this module was
// compiled. All includes made by the Thrift file are relative to that path.
//
// The module name is usually just the basename of the ThriftPath.
type Module struct {
	Name       string
	ThriftPath string

	// Mapping from the /Thrift name/ to the compiled representation of
	// different definitions.

	Includes  map[string]*IncludedModule
	Constants map[string]Constant
	Types     map[string]TypeSpec
	Services  map[string]*Service
}

// LookupType TODO
func (m *Module) LookupType(name string) (TypeSpec, error) {
	if t, ok := m.Types[name]; ok {
		return t, nil
	}

	if mname, name := splitInclude(name); len(mname) > 0 {
		if included, ok := m.Includes[mname]; ok {
			return included.Module.LookupType(name)
		}
	}

	return nil, lookupError{Name: name}
}

// LookupService TODO
func (m *Module) LookupService(name string) (*Service, error) {
	return nil, nil // TODO
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

// TODO(abg): Add support for include-as syntax.

// lookupError is raised when an unknown identifier is requested via the
// Lookup* methods.
type lookupError struct {
	Name string
}

func (e lookupError) Error() string {
	return fmt.Sprintf("unrecognized identifier '%s'", e.Name)
}
