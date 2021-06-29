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

	"go.uber.org/thriftrw/ast"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	c := &Config{}
	prog, err := c.Parse([]byte{})
	if assert.NoError(t, err) {
		assert.Equal(t, &ast.Program{}, prog)
	}
}

func TestInfoPos(t *testing.T) {
	c := &Config{Info: &Info{}}
	prog, err := c.Parse([]byte(`const string a = 'a';`))
	if assert.NoError(t, err) {
		assert.Equal(t, ast.Position{Line: 0}, c.Info.Pos(prog))
		assert.Equal(t, ast.Position{Line: 1}, c.Info.Pos(prog.Definitions[0]))
		if assert.IsType(t, &ast.Constant{}, prog.Definitions[0]) {
			cv := prog.Definitions[0].(*ast.Constant).Value
			assert.Equal(t, ast.Position{Line: 1, Column: 18}, c.Info.Pos(cv))
		}
	}
}
