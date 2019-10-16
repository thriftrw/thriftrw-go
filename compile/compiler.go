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

import (
	"path/filepath"
	"strings"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
)

// Compile parses and compiles the Thrift file at the given path and any other
// Thrift file it includes.
func Compile(path string, opts ...Option) (*Module, error) {
	c := newCompiler()
	for _, opt := range opts {
		opt(&c)
	}

	m, err := c.load(path)
	if err != nil {
		return nil, err
	}

	err = m.Walk(func(m *Module) error {
		if err := c.link(m); err != nil {
			return compileError{
				Target: m.ThriftPath,
				Reason: err,
			}
		}
		return nil
	})
	return m, err
}

// compiler is responsible for compiling Thrift files.
type compiler struct {
	// fs is the interface used to interact with the filesystem.
	fs FS
	// nonStrict will compile Thrift files that do not pass strict validation.
	nonStrict bool
	// Map from file path to Module representing that file.
	Modules map[string]*Module
}

func newCompiler() compiler {
	return compiler{
		fs:      realFS{},
		Modules: make(map[string]*Module),
	}
}

func (c compiler) link(m *Module) error {
	// TODO(abg): might be worth accumulating compile errors with a max count

	// make a copy so that we can modify the list of types as we're iterating
	// through it.
	types := make(map[string]TypeSpec)
	for name, typ := range m.Types {
		types[name] = typ
	}

	for name, typ := range types {
		var err error
		m.Types[name], err = typ.Link(m)
		if err != nil {
			return compileError{Target: name, Reason: err}
		}
	}

	for name, constant := range m.Constants {
		if err := constant.Link(m); err != nil {
			return compileError{Target: name, Reason: err}
		}
	}

	for name, service := range m.Services {
		if err := service.Link(m); err != nil {
			return compileError{Target: name, Reason: err}
		}
	}

	// Find cycles in typedefs
	for name, t := range types {
		if _, ok := t.(*TypedefSpec); !ok {
			continue
		}

		if err := findTypeCycles(t); err != nil {
			return compileError{Target: name, Reason: err}
		}
	}

	return nil
}

// load populates the compiler with information from the given Thrift file.
//
// The types aren't actually compiled in this step.
func (c compiler) load(p string) (*Module, error) {
	p, err := c.fs.Abs(p)
	if err != nil {
		return nil, err
	}

	if m, ok := c.Modules[p]; ok {
		// Already loaded.
		return m, nil
	}

	s, err := c.fs.Read(p)
	if err != nil {
		return nil, fileReadError{Path: p, Reason: err}
	}

	prog, err := idl.Parse(s)
	if err != nil {
		return nil, parseError{Path: p, Reason: err}
	}

	m := &Module{
		Name:       fileBaseName(p),
		ThriftPath: p,
		Includes:   make(map[string]*IncludedModule),
		Constants:  make(map[string]*Constant),
		Types:      make(map[string]TypeSpec),
		Services:   make(map[string]*ServiceSpec),
	}

	m.Raw = s
	c.Modules[p] = m
	// the module is added to the map before processing includes to break
	// cyclic includes.

	if err := c.gather(m, prog); err != nil {
		return nil, fileCompileError{Path: p, Reason: err}
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
			constant, err := compileConstant(m.ThriftPath, definition)
			if err != nil {
				return definitionError{Definition: d, Reason: err}
			}
			m.Constants[constant.Name] = constant
		case *ast.Typedef:
			typedef, err := compileTypedef(m.ThriftPath, definition)
			if err != nil {
				return definitionError{Definition: d, Reason: err}
			}
			m.Types[typedef.ThriftName()] = typedef
		case *ast.Enum:
			enum, err := compileEnum(m.ThriftPath, definition)
			if err != nil {
				return definitionError{Definition: d, Reason: err}
			}
			m.Types[enum.ThriftName()] = enum
		case *ast.Struct:
			requiredness := explicitRequiredness
			if c.nonStrict {
				requiredness = defaultToOptional
			}
			s, err := compileStruct(m.ThriftPath, definition, requiredness)
			if err != nil {
				return definitionError{Definition: d, Reason: err}
			}
			m.Types[s.ThriftName()] = s
		case *ast.Service:
			service, err := compileService(m.ThriftPath, definition)
			if err != nil {
				return definitionError{Definition: d, Reason: err}
			}
			m.Services[service.Name] = service
		}
	}

	return nil
}

// include loads the file specified by the given include in the given Module.
//
// The path to the file is relative to the ThriftPath of the given module.
// Including hyphenated file names will error.
func (c compiler) include(m *Module, include *ast.Include) (*IncludedModule, error) {
	if len(include.Name) > 0 {
		// TODO(abg): Add support for include-as flag somewhere.
		return nil, includeError{
			Include: include,
			Reason:  includeAsDisabledError{},
		}
	}

	if strings.Contains(fileBaseName(include.Path), "-") {
		return nil, includeError{
			Include: include,
			Reason:  includeHyphenatedFileNameError{},
		}
	}

	ipath := filepath.Join(filepath.Dir(m.ThriftPath), include.Path)
	incM, err := c.load(ipath)
	if err != nil {
		return nil, includeError{Include: include, Reason: err}
	}

	return &IncludedModule{Name: fileBaseName(include.Path), Module: incM}, nil
}
