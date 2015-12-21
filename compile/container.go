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
	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/wire"
)

// MapSpec represents a key-value mapping between two types.
type MapSpec struct {
	KeySpec, ValueSpec TypeSpec

	compiled bool
	src      ast.MapType
}

// NewMapSpec constructs a new MapSpec from the given AST.
func NewMapSpec(src ast.MapType) *MapSpec {
	return &MapSpec{src: src, compiled: false}
}

// Compile resolves the type references in the MapSpec.
func (m *MapSpec) Compile(scope Scope) error {
	if m.compiled {
		return nil
	}
	m.compiled = true

	var err error

	m.KeySpec, err = ResolveType(m.src.KeyType, scope)
	if err != nil {
		return err
	}

	m.ValueSpec, err = ResolveType(m.src.ValueType, scope)
	return err
}

// TypeCode for MapSpec
func (m *MapSpec) TypeCode() wire.Type {
	return wire.TMap
}

//////////////////////////////////////////////////////////////////////////////

// ListSpec represents lists of values of the same type.
type ListSpec struct {
	ValueSpec TypeSpec

	compiled bool
	src      ast.ListType
}

// NewListSpec constructs a new ListSpec from the given AST.
func NewListSpec(src ast.ListType) *ListSpec {
	return &ListSpec{src: src, compiled: false}
}

// Compile resolves the type references in the ListSpec.
func (m *ListSpec) Compile(scope Scope) error {
	if m.compiled {
		return nil
	}
	m.compiled = true

	spec, err := ResolveType(m.src.ValueType, scope)
	m.ValueSpec = spec
	return err
}

// TypeCode for ListSpec
func (m *ListSpec) TypeCode() wire.Type {
	return wire.TList
}

//////////////////////////////////////////////////////////////////////////////

// SetSpec represents sets of values of the same type.
type SetSpec struct {
	ValueSpec TypeSpec

	compiled bool
	src      ast.SetType
}

// NewSetSpec constructs a new SetSpec from the given AST.
func NewSetSpec(src ast.SetType) *SetSpec {
	return &SetSpec{src: src, compiled: false}
}

// Compile resolves the type references in the SetSpec.
func (m *SetSpec) Compile(scope Scope) error {
	if m.compiled {
		return nil
	}
	m.compiled = true

	spec, err := ResolveType(m.src.ValueType, scope)
	m.ValueSpec = spec
	return err
}

// TypeCode for SetSpec
func (m *SetSpec) TypeCode() wire.Type {
	return wire.TSet
}
