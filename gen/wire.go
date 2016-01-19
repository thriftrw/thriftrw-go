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

package gen

import (
	"fmt"

	"github.com/uber/thriftrw-go/compile"
)

// toWire generates a call to the given variable of the given type.
func (g *Generator) toWire(spec compile.TypeSpec, varName string) (string, error) {
	wire := g.Import("github.com/uber/thriftrw-go/wire")
	switch spec {
	case compile.BoolSpec:
		return fmt.Sprintf("%s.NewValueBool(%s)", wire, varName), nil
	case compile.I8Spec:
		return fmt.Sprintf("%s.NewValueI8(%s)", wire, varName), nil
	case compile.I16Spec:
		return fmt.Sprintf("%s.NewValueI16(%s)", wire, varName), nil
	case compile.I32Spec:
		return fmt.Sprintf("%s.NewValueI32(%s)", wire, varName), nil
	case compile.I64Spec:
		return fmt.Sprintf("%s.NewValueI64(%s)", wire, varName), nil
	case compile.DoubleSpec:
		return fmt.Sprintf("%s.NewValueDouble(%s)", wire, varName), nil
	case compile.StringSpec:
		return fmt.Sprintf("%s.NewValueString(%s)", wire, varName), nil
	case compile.BinarySpec:
		return fmt.Sprintf("%s.NewValueBinary(%s)", wire, varName), nil
	default:
		// Not a primitive type. It's probably a container or a custom type.
	}

	switch spec.(type) {
	case *compile.MapSpec:
		// TODO unhashable types
		// TODO generate MapItemList alias if necessary
		return g.TextTemplate(
			`<.Wire>.NewValueMap(<.Wire>.Map{
				KeyType: TODO,
				ValueType: TODO,
				Size: len(<.Name>),
				Items: TODO(<.Name>),
			})`,
			struct {
				Wire string
				Name string
			}{Wire: wire, Name: varName},
		)
	case *compile.ListSpec:
		return g.TextTemplate(
			`<.Wire>.NewValueList(<.Wire>.List{
				ValueType: TODO,
				Size: len(<.Name>),
				Items: TODO(<.Name>),
			})`,
			struct {
				Wire string
				Name string
			}{Wire: wire, Name: varName},
		)
	case *compile.SetSpec:
		// TODO unhashable types
		return g.TextTemplate(
			`<.Wire>.NewValueSet(<.Wire>.Set{
				ValueType: TODO,
				Size: len(<.Name>),
				Items: TODO(<.Name>),
			})`,
			struct {
				Wire string
				Name string
			}{Wire: wire, Name: varName},
		)
	default:
		// Custom defined type
		return fmt.Sprintf("%s.ToWire()", varName), nil
	}
}

func (g *Generator) fromWire(spec compile.TypeSpec, target string, value string) (string, error) {
	switch spec {
	case compile.BoolSpec:
		return fmt.Sprintf("%s = %s.GetBool()", target, value), nil
	case compile.I8Spec:
		return fmt.Sprintf("%s = %s.GetI8()", target, value), nil
	case compile.I16Spec:
		return fmt.Sprintf("%s = %s.GetI16()", target, value), nil
	case compile.I32Spec:
		return fmt.Sprintf("%s = %s.GetI32()", target, value), nil
	case compile.I64Spec:
		return fmt.Sprintf("%s = %s.GetI64()", target, value), nil
	case compile.DoubleSpec:
		return fmt.Sprintf("%s = %s.GetDouble()", target, value), nil
	case compile.StringSpec:
		return fmt.Sprintf("%s = %s.GetString()", target, value), nil
	case compile.BinarySpec:
		return fmt.Sprintf("%s = %s.GetBinary()", target, value), nil
	default:
		// Not a primitive type. It's probably a container or a custom type.
	}

	switch spec.(type) {
	case *compile.MapSpec:
		return fmt.Sprintf("%s = %s.GetList().TODO()", target, value), nil
	case *compile.ListSpec:
		return fmt.Sprintf("%s = %s.GetMap().TODO()", target, value), nil
	case *compile.SetSpec:
		return fmt.Sprintf("%s = %s.GetSet().TODO()", target, value), nil
	default:
		// TODO read errors
		return fmt.Sprintf("%s.FromWire(%s)", target, value), nil
	}
}
