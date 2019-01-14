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
	"fmt"

	"go.uber.org/thriftrw/ast"
)

// Constant represents a single named constant value from the Thrift file.
type Constant struct {
	linkOnce

	Name  string
	File  string
	Doc   string
	Type  TypeSpec
	Value ConstantValue
}

// compileConstant builds a Constant from the given AST constant.
func compileConstant(file string, src *ast.Constant) (*Constant, error) {
	typ, err := compileTypeReference(src.Type)
	if err != nil {
		return nil, err
	}

	return &Constant{
		Name:  src.Name,
		File:  file,
		Type:  typ,
		Doc:   src.Doc,
		Value: compileConstantValue(src.Value),
	}, nil
}

// Link resolves any references made by the constant.
func (c *Constant) Link(scope Scope) (err error) {
	if c.linked() {
		return nil
	}

	if c.Type, err = c.Type.Link(scope); err != nil {
		return compileError{Target: c.Name, Reason: err}
	}

	if c.Value, err = c.Value.Link(scope, c.Type); err != nil {
		return compileError{Target: c.Name, Reason: err}
	}

	return nil
}

func (c *Constant) String() string {
	return fmt.Sprintf("Constant(%s %s)", c.Type.ThriftName(), c.Name)
}
