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

// EnumSpec represents an enum defined in the Thrift file.
type EnumSpec struct {
	compileOnce

	Items []EnumItem
	src   *ast.Enum
}

// EnumItem is a single item inside an enum.
type EnumItem struct {
	Name  string
	Value ast.ConstantValue
}

// ThriftName for EnumSpec
func (e *EnumSpec) ThriftName() string {
	return e.src.Name
}

// TypeCode for EnumSpec.
//
// Enums are represented as i32 over the wire.
func (e *EnumSpec) TypeCode() wire.Type {
	return wire.TI32
}

// NewEnumSpec creates a new uncompiled EnumSpec from the given AST
// definition.
func NewEnumSpec(src *ast.Enum) *EnumSpec {
	return &EnumSpec{src: src}
}

// Compile compiles the EnumSpec.
func (e *EnumSpec) Compile(scope Scope) error {
	if e.compiled() {
		return nil
	}

	enumNS := newNamespace(caseInsensitive)
	prev := -1
	var items []EnumItem
	for _, astItem := range e.src.Items {
		if err := enumNS.claim(astItem.Name, astItem.Line); err != nil {
			return compileError{
				Target: e.ThriftName() + "." + astItem.Name,
				Line:   astItem.Line,
				Reason: err,
			}
		}

		value := prev + 1
		if astItem.Value != nil {
			value = *astItem.Value
		}
		prev = value

		item := EnumItem{
			Name:  astItem.Name,
			Value: ast.ConstantInteger(value),
		}
		items = append(items, item)
	}

	e.Items = items
	return nil
}
