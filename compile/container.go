// Copyright (c) 2016 Uber Technologies, Inc.
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
	"fmt"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

// MapSpec represents a key-value mapping between two types.
type MapSpec struct {
	linkOnce
	nativeThriftType

	KeySpec, ValueSpec TypeSpec
}

// compileMapType compiles the given MapType AST into a MapSpec.
func compileMapType(src ast.MapType) *MapSpec {
	return &MapSpec{
		KeySpec:   compileTypeReference(src.KeyType),
		ValueSpec: compileTypeReference(src.ValueType),
	}
}

// Link resolves the type references in the MapSpec.
func (m *MapSpec) Link(scope Scope) (TypeSpec, error) {
	if m.linked() {
		return m, nil
	}

	var err error
	m.KeySpec, err = m.KeySpec.Link(scope)
	if err != nil {
		return m, err
	}

	m.ValueSpec, err = m.ValueSpec.Link(scope)
	if err != nil {
		return m, err
	}

	return m, nil
}

// ThriftName for MapSpec
func (m *MapSpec) ThriftName() string {
	return fmt.Sprintf(
		"map<%s, %s>", m.KeySpec.ThriftName(), m.ValueSpec.ThriftName(),
	)
}

// TypeCode for MapSpec
func (m *MapSpec) TypeCode() wire.Type {
	return wire.TMap
}

// ForEachTypeReference for MapSpec
func (m *MapSpec) ForEachTypeReference(f func(TypeSpec) error) error {
	if err := f(m.KeySpec); err != nil {
		return err
	}
	if err := f(m.ValueSpec); err != nil {
		return err
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////

// ListSpec represents lists of values of the same type.
type ListSpec struct {
	linkOnce
	nativeThriftType

	ValueSpec TypeSpec
}

// compileListSpec compiles the given ListType AST into a ListSpec.
func compileListType(src ast.ListType) *ListSpec {
	return &ListSpec{ValueSpec: compileTypeReference(src.ValueType)}
}

// Link resolves the type references in the ListSpec.
func (l *ListSpec) Link(scope Scope) (TypeSpec, error) {
	if l.linked() {
		return l, nil
	}

	var err error
	l.ValueSpec, err = l.ValueSpec.Link(scope)
	return l, err
}

// TypeCode for ListSpec
func (l *ListSpec) TypeCode() wire.Type {
	return wire.TList
}

// ThriftName for ListSpec
func (l *ListSpec) ThriftName() string {
	return fmt.Sprintf("list<%s>", l.ValueSpec.ThriftName())
}

// ForEachTypeReference for ListSpec
func (l *ListSpec) ForEachTypeReference(f func(TypeSpec) error) error {
	return f(l.ValueSpec)
}

//////////////////////////////////////////////////////////////////////////////

// SetSpec represents sets of values of the same type.
type SetSpec struct {
	linkOnce
	nativeThriftType

	ValueSpec TypeSpec
}

// compileSetSpec compiles the given SetType AST into a SetSpec.
func compileSetType(src ast.SetType) *SetSpec {
	return &SetSpec{ValueSpec: compileTypeReference(src.ValueType)}
}

// Link resolves the type references in the SetSpec.
func (s *SetSpec) Link(scope Scope) (TypeSpec, error) {
	if s.linked() {
		return s, nil
	}

	var err error
	s.ValueSpec, err = s.ValueSpec.Link(scope)
	return s, err
}

// TypeCode for SetSpec
func (s *SetSpec) TypeCode() wire.Type {
	return wire.TSet
}

// ThriftName for SetSpec.
func (s *SetSpec) ThriftName() string {
	return fmt.Sprintf("set<%s>", s.ValueSpec.ThriftName())
}

// ForEachTypeReference for SetSpec
func (s *SetSpec) ForEachTypeReference(f func(TypeSpec) error) error {
	return f(s.ValueSpec)
}
