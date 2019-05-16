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

package gen

import (
	"fmt"

	"go.uber.org/thriftrw/compile"
)

type mangler struct {
	// Thrift file -> Type name -> value name
	names map[string]map[string]string

	// all names for custom types that have been taken so far
	taken map[string]struct{}
}

func newMangler() *mangler {
	return &mangler{
		names: make(map[string]map[string]string),
		taken: make(map[string]struct{}),
	}
}

func (m *mangler) MangleType(spec compile.TypeSpec) string {
	switch s := spec.(type) {
	case *compile.MapSpec:
		return fmt.Sprintf(
			"Map_%s_%s", m.MangleType(s.KeySpec), m.MangleType(s.ValueSpec),
		)
	case *compile.ListSpec:
		return fmt.Sprintf("List_%s", m.MangleType(s.ValueSpec))
	case *compile.SetSpec:
		setType := "slice"
		if setUsesMap(s) {
			setType = "map"
		}

		return fmt.Sprintf("Set_%s_%vType", m.MangleType(s.ValueSpec), setType)
	}

	// Native primitive types have unique names
	thriftFile := spec.ThriftFile()
	if thriftFile == "" {
		return goCase(spec.ThriftName())
	}

	namesForFile, ok := m.names[thriftFile]
	if !ok {
		namesForFile = make(map[string]string)
		m.names[thriftFile] = namesForFile
	}

	if name, ok := namesForFile[spec.ThriftName()]; ok {
		return name
	}

	i := 0
	baseName := goCase(spec.ThriftName())
	name := baseName
	for _, taken := m.taken[name]; taken; {
		i++
		name = fmt.Sprintf("%s_%d", baseName, i)
		_, taken = m.taken[name]
	}

	m.taken[name] = struct{}{}
	namesForFile[spec.ThriftName()] = name
	return name
}
