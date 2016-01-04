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

import "github.com/uber/thriftrw-go/ast"

// Constant represents a single named constant value from the Thrift file.
type Constant struct {
	linkOnce

	Name  string
	Type  TypeSpec
	Value ast.ConstantValue
}

// compileConstant builds a Constant from the given AST constant.
func compileConstant(src *ast.Constant) *Constant {
	return &Constant{
		Name:  src.Name,
		Type:  compileType(src.Type),
		Value: src.Value,
	}
}

// Link resolves any references made by the constant.
func (c *Constant) Link(scope Scope) (err error) {
	if c.linked() {
		return nil
	}

	if c.Type, err = c.Type.Link(scope); err != nil {
		return err
	}
	err = verifyConstantValue(c.Value, scope)
	// TODO(abg): validate that the constant matches the TypeSpec
	return err
}

// LinkConstantValue ensures that all references made by the given constant
// value are valid.
func verifyConstantValue(v ast.ConstantValue, scope Scope) error {
	// TODO(abg): We'll need a separate ConstantValue type that tracks whether
	// that constant has already been linked/verified to break cycles.

	switch c := v.(type) {
	case ast.ConstantReference:
		// Note that ConstantReferences are not resolved to their target values.
		// We only verify that the references are valid. We do this because we
		// may want constant references to be actual references in the generated
		// code.
		if _, err := scope.LookupConstant(c.Name); err != nil {
			return referenceError{
				Target: c.Name,
				Line:   c.Line,
				Reason: err,
			}
		}
	case ast.ConstantMap:
		for _, item := range c.Items {
			if err := verifyConstantValue(item.Key, scope); err != nil {
				return err
			}
			if err := verifyConstantValue(item.Value, scope); err != nil {
				return err
			}
		}
	case ast.ConstantList:
		for _, item := range c.Items {
			if err := verifyConstantValue(item, scope); err != nil {
				return err
			}
		}
	default:
		// primitive constant. do nothing.
	}
	return nil
}
