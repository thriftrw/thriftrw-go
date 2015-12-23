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

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/idl"
)

// Compile parses and compiles the Thrift file at the given path and any other
// Thrift file it includes.
func Compile(path string) (*Module, error) {
	return newCompiler().compile(path)
}

// compiler is responsible for compiling Thrift files.
type compiler struct {
	// Map from file path to Module representing that file.
	Modules map[string]*Module
}

func newCompiler() compiler {
	return compiler{
		Modules: make(map[string]*Module),
	}
}

func (c compiler) compile(p string) (*Module, error) {
	m, err := c.load(p)
	if err != nil {
		return m, err
	}

	// TODO(abg): compile includes
	// TODO(abg): compile constants
	// TODO(abg): compile services
	// TODO(abg): might be worth accumulating compile errors with a max count

	for _, typ := range m.Types {
		if err := typ.Compile(m); err != nil {
			return m, err
		}
	}

	return m, nil
}

// load populates the compiler with information from the given Thrift file.
//
// The types aren't actually compiled in this step.
func (c compiler) load(p string) (*Module, error) {
	p, err := filepath.Abs(p)
	if err != nil {
		return nil, err
	}

	if m, ok := c.Modules[p]; ok {
		// Already loaded.
		return m, nil
	}

	s, err := ioutil.ReadFile(p)
	if err != nil {
		// TODO(abg): real error type instead of strings
		return nil, fmt.Errorf("error reading %s: %s", p, err)
	}

	prog, err := idl.Parse(s)
	if err != nil {
		// TODO(abg): real error type instead of strings
		return nil, fmt.Errorf("error parsing %s: %s", p, err)
	}

	m := &Module{
		Name:       fileBaseName(p),
		ThriftPath: p,
		Includes:   make(map[string]*IncludedModule),
		Constants:  make(map[string]Constant),
		Types:      make(map[string]TypeSpec),
		Services:   make(map[string]*Service),
	}
	c.Modules[p] = m
	// the module is added to the map before processing includes to break
	// cyclic includes.

	if err := c.gather(m, prog); err != nil {
		// TODO(abg): Real error types intsead of string
		return nil, fmt.Errorf("failed to compile %s: %s", p, err)
	}
	return m, nil
}

// gather populates the Module for the given program with knowledge about all
// definitions from the Thrift file.
//
// It recursively processes includes, relying on load() to break cycles.
//
// prog is the parsed representation of it, and m is the Module representing
// this file.
func (c compiler) gather(m *Module, prog *ast.Program) error {
	// Namespace of items defined in the Thrift file.
	//
	// This is not shared with the Go namespace because we will capitalize
	// names and possibly allow overriding them with annotations.
	thriftNS := newNamespace(caseSensitive)

	// Process all included modules first.
	for _, h := range prog.Headers {
		header, ok := h.(*ast.Include)
		if !ok {
			continue
		}

		include, err := c.include(m, header)
		if err != nil {
			return err
		}

		if err := thriftNS.claim(include.Name, header.Line); err != nil {
			return includeError{
				Include: header,
				Reason:  err,
			}
		}

		m.Includes[include.Name] = include
	}

	for _, d := range prog.Definitions {
		if err := thriftNS.claim(d.Info().Name, d.Info().Line); err != nil {
			return definitionError{Definition: d, Reason: err}
		}

		switch definition := d.(type) {
		case *ast.Constant:
			// TODO
		case *ast.Typedef:
			// TODO
		case *ast.Enum:
			enum := NewEnumSpec(definition)
			m.Types[enum.ThriftName()] = enum
		case *ast.Struct:
			// TODO
		case *ast.Service:
			// TODO
		}
	}

	return nil
}

// include loads the file specified by the given include in the given Module.
//
// The path to the file is relative to the ThriftPath of the given module.
func (c compiler) include(m *Module, include *ast.Include) (*IncludedModule, error) {
	if len(include.Name) > 0 {
		// TODO(abg): Add support for include-as flag somewhere.
		return nil, includeError{
			Include: include,
			Reason:  errors.New("include-as syntax is currently disabled"),
		}
	}

	ipath := filepath.Join(filepath.Dir(m.ThriftPath), include.Path)
	incM, err := c.load(ipath)
	if err != nil {
		return nil, includeError{Include: include, Reason: err}
	}

	return &IncludedModule{Name: fileBaseName(include.Path), Module: incM}, nil
}
