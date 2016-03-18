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
	"go/token"
	"io"
	"sort"

	"github.com/thriftrw/thriftrw-go/compile"
)

// Options controls how code gets generated.
type Options struct {
	PackageName string
	Output      io.Writer
}

// Generate generates code based on the given options.
func Generate(m *compile.Module, o *Options) error {
	g := NewGenerator(o.PackageName)

	for _, constantName := range constantNames(m.Constants) {
		c := m.Constants[constantName]
		if err := Constant(g, c); err != nil {
			return err
		}
	}

	for _, typeName := range typeNames(m.Types) {
		t := m.Types[typeName]
		if err := TypeDefinition(g, t); err != nil {
			return err
		}
	}

	if err := g.Write(o.Output, token.NewFileSet()); err != nil {
		return err
	}

	return nil
}

// constantNames sorts the keys of a map of constants in a deterministic order.
func constantNames(constants map[string]*compile.Constant) []string {
	names := make([]string, 0, len(constants))
	for name := range constants {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// typenames sorts the keys of a map of types in a deterministic order.
func typeNames(types map[string]compile.TypeSpec) []string {
	names := make([]string, 0, len(types))
	for name := range types {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
