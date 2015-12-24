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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/wire"
)

func TestResolveBaseType(t *testing.T) {
	tests := []struct {
		input    ast.BaseType
		wireType wire.Type
	}{
		{ast.BaseType{ID: ast.BoolTypeID}, wire.TBool},
		{ast.BaseType{ID: ast.I8TypeID}, wire.TI8},
		{ast.BaseType{ID: ast.I16TypeID}, wire.TI16},
		{ast.BaseType{ID: ast.I32TypeID}, wire.TI32},
		{ast.BaseType{ID: ast.I64TypeID}, wire.TI64},
		{ast.BaseType{ID: ast.DoubleTypeID}, wire.TDouble},
		{ast.BaseType{ID: ast.StringTypeID}, wire.TBinary},
		{ast.BaseType{ID: ast.BinaryTypeID}, wire.TBinary},
	}

	for _, tt := range tests {
		spec := compileType(tt.input)
		linked, err := spec.Link(scope())

		assert.NoError(t, err)
		assert.Equal(t, tt.wireType, spec.TypeCode())
		assert.Equal(t, tt.wireType, linked.TypeCode())
	}
}

func TestResolveInvalidBaseType(t *testing.T) {
	assert.Panics(t, func() {
		compileType(ast.BaseType{ID: ast.BaseTypeID(42)})
	})
}
