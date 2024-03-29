// Copyright (c) 2021 Uber Technologies, Inc.
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

package idl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl/internal"
)

func TestPos(t *testing.T) {
	tests := []struct {
		node ast.Node
		pos  *ast.Position
		want ast.Position
	}{
		{
			node: &ast.Struct{Line: 10},
			want: ast.Position{Line: 10},
		},
		{
			node: ast.ConstantString("s"),
			want: ast.Position{Line: 0},
		},
		{
			node: ast.ConstantString("s"),
			pos:  &ast.Position{Line: 1},
			want: ast.Position{Line: 1},
		},
	}

	for _, tt := range tests {
		i := &Info{}
		if tt.pos != nil {
			i.nodePositions = internal.NodePositions{tt.node: *tt.pos}
		}
		assert.Equal(t, tt.want, i.Pos(tt.node))
	}
}
