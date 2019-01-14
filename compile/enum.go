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
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

// EnumSpec represents an enum defined in the Thrift file.
type EnumSpec struct {
	Name        string
	File        string
	Items       []EnumItem
	Annotations Annotations
	Doc         string
}

// EnumItem is a single item inside an enum.
type EnumItem struct {
	Name        string
	Value       int32
	Annotations Annotations
	Doc         string
}

// compileEnum compiles the given Enum AST into an EnumSpec.
func compileEnum(file string, src *ast.Enum) (*EnumSpec, error) {
	enumNS := newNamespace(caseInsensitive)
	prev := -1

	var items []EnumItem
	for _, astItem := range src.Items {
		if err := enumNS.claim(astItem.Name, astItem.Line); err != nil {
			return nil, compileError{
				Target: src.Name + "." + astItem.Name,
				Line:   astItem.Line,
				Reason: err,
			}
		}
		value := prev + 1
		if astItem.Value != nil {
			value = *astItem.Value
		}
		prev = value

		itemAnnotations, err := compileAnnotations(astItem.Annotations)
		if err != nil {
			return nil, compileError{
				Target: src.Name + "." + astItem.Name,
				Line:   astItem.Line,
				Reason: err,
			}
		}
		// TODO bounds check for value
		item := EnumItem{
			Name:        astItem.Name,
			Value:       int32(value),
			Doc:         astItem.Doc,
			Annotations: itemAnnotations,
		}
		items = append(items, item)
	}

	annotations, err := compileAnnotations(src.Annotations)
	if err != nil {
		return nil, compileError{
			Target: src.Name,
			Line:   src.Line,
			Reason: err,
		}
	}
	return &EnumSpec{
		Name:        src.Name,
		File:        file,
		Doc:         src.Doc,
		Items:       items,
		Annotations: annotations,
	}, nil
}

// LookupItem retrieves the item with the given name from the enum.
//
// Returns true or false indicating whether the result is valid or not.
func (e *EnumSpec) LookupItem(name string) (*EnumItem, bool) {
	for _, item := range e.Items {
		if item.Name == name {
			return &item, true
		}
	}
	return nil, false
}

// Link resolves any references made by the Enum.
func (e *EnumSpec) Link(scope Scope) (TypeSpec, error) {
	return e, nil // nothing to do
}

// ThriftName for EnumSpec
func (e *EnumSpec) ThriftName() string {
	return e.Name
}

// ThriftFile for EnumSpec
func (e *EnumSpec) ThriftFile() string {
	return e.File
}

// ForEachTypeReference for EnumSpec
func (e *EnumSpec) ForEachTypeReference(func(TypeSpec) error) error {
	return nil
}

// TypeCode for EnumSpec.
//
// Enums are represented as i32 over the wire.
func (e *EnumSpec) TypeCode() wire.Type {
	return wire.TI32
}

// ThriftAnnotations returns all associated annotations.
func (e *EnumSpec) ThriftAnnotations() Annotations {
	return e.Annotations
}

// ThriftName for EnumItem
func (e *EnumItem) ThriftName() string {
	return e.Name
}

// ThriftAnnotations returns all associated annotations.
func (e *EnumItem) ThriftAnnotations() Annotations {
	return e.Annotations
}
