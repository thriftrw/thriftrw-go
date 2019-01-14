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

// EqualsGenerator is responsible for generating code that knows how
// to compare the equality of two Thrift types and their Value representations.
type equalsGenerator struct {
	mapG  mapGenerator
	setG  setGenerator
	listG listGenerator
}

// Equals generates a string comparing rhs to the given lhs.
// Equals generates an expression of type bool.
func (e *equalsGenerator) Equals(g Generator, spec compile.TypeSpec, lhs, rhs string) (string, error) {
	if isPrimitiveType(spec) {
		if _, isEnum := spec.(*compile.EnumSpec); !isEnum {
			return fmt.Sprintf("(%s == %s)", lhs, rhs), nil
		}
	}

	switch s := spec.(type) {
	case *compile.BinarySpec:
		bytes := g.Import("bytes")
		return fmt.Sprintf("%s.Equal(%s, %s)", bytes, lhs, rhs), nil
	case *compile.MapSpec:
		equals, err := e.mapG.Equals(g, s)
		return fmt.Sprintf("%s(%s, %s)", equals, lhs, rhs), err
	case *compile.ListSpec:
		equals, err := e.listG.Equals(g, s)
		return fmt.Sprintf("%s(%s, %s)", equals, lhs, rhs), err
	case *compile.SetSpec:
		equals, err := e.setG.Equals(g, s)
		return fmt.Sprintf("%s(%s, %s)", equals, lhs, rhs), err
	default:
		// Custom defined type
		return fmt.Sprintf("%s.Equals(%s)", lhs, rhs), nil
	}
}

// EqualsPtr is the same as Equals except `lhs` and `rhs` are expected to be a
// reference to a value of the given type.
func (e *equalsGenerator) EqualsPtr(g Generator, spec compile.TypeSpec, lhs, rhs string) (string, error) {
	if !isPrimitiveType(spec) {
		// Everything else is a reference type that has a Equals method on it.
		return g.TextTemplate(
			`((<.LHS> == nil && <.RHS> == nil) || (<.LHS> != nil && <.RHS> != nil && <equals .Spec .LHS .RHS>))`,
			struct {
				Spec compile.TypeSpec
				LHS  string
				RHS  string
			}{Spec: spec, LHS: lhs, RHS: rhs},
		)
	}

	name := equalsPtrFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$type := typeReference .Spec>
			<$lhs := newVar "lhs">
			<$rhs := newVar "rhs">
			func <.Name>(<$lhs>, <$rhs> *<$type>) bool {
				<- $x := newVar "x" ->
				<- $y := newVar "y">
				if <$lhs> != nil && <$rhs> != nil {
					// Call Equals method after dereferencing the pointers
					<$x> := *<$lhs>
					<$y> := *<$rhs>
					return <equals .Spec $x $y>
				}
				return <$lhs> == nil && <$rhs> == nil
			}
		`,
		struct {
			Name string
			Spec compile.TypeSpec
		}{Name: name, Spec: spec},
	)
	return fmt.Sprintf("%s(%s, %s)", name, lhs, rhs), err
}
