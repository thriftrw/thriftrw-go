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

	switch s := spec.(type) {
	case *compile.MapSpec:
		// TODO unhashable types
		mapItemList, err := g.mapItemList(s)
		if err != nil {
			return "", err
		}

		return g.TextTemplate(
			`<.Wire>.NewValueMap(<.Wire>.Map{
				KeyType: <typeCode .Spec.KeySpec>,
				ValueType: <typeCode .Spec.ValueSpec>,
				Size: len(<.Name>),
				Items: <.MapItemList>(<.Name>),
			})`,
			struct {
				Wire        string
				Name        string
				Spec        *compile.MapSpec
				MapItemList string
			}{Wire: wire, Name: varName, Spec: s, MapItemList: mapItemList},
		)
	case *compile.ListSpec:
		valueList, err := g.listValueList(s)
		if err != nil {
			return "", err
		}

		return g.TextTemplate(
			`<.Wire>.NewValueList(<.Wire>.List{
				ValueType: <typeCode .Spec.ValueSpec>,
				Size: len(<.Name>),
				Items: <.ValueList>(<.Name>),
			})`,
			struct {
				Wire      string
				Name      string
				Spec      *compile.ListSpec
				ValueList string
			}{Wire: wire, Name: varName, Spec: s, ValueList: valueList},
		)
	case *compile.SetSpec:
		valueList, err := g.setValueList(s)
		if err != nil {
			return "", err
		}

		// TODO unhashable types
		return g.TextTemplate(
			`<.Wire>.NewValueSet(<.Wire>.Set{
				ValueType: <typeCode .Spec.ValueSpec>,
				Size: len(<.Name>),
				Items: <.ValueList>(<.Name>),
			})`,
			struct {
				Wire      string
				Name      string
				Spec      *compile.SetSpec
				ValueList string
			}{Wire: wire, Name: varName, Spec: s, ValueList: valueList},
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

// typeCode gets a value of type 'wire.Type' that represents the over-the-wire
// type code for the given TypeSpec.
func (g *Generator) typeCode(spec compile.TypeSpec) string {
	wire := g.Import("github.com/uber/thriftrw-go/wire")

	switch spec {
	case compile.BoolSpec:
		return fmt.Sprintf("%s.TBool", wire)
	case compile.I8Spec:
		return fmt.Sprintf("%s.TI8", wire)
	case compile.I16Spec:
		return fmt.Sprintf("%s.TI16", wire)
	case compile.I32Spec:
		return fmt.Sprintf("%s.TI32", wire)
	case compile.I64Spec:
		return fmt.Sprintf("%s.TI64", wire)
	case compile.DoubleSpec:
		return fmt.Sprintf("%s.TDouble", wire)
	case compile.StringSpec, compile.BinarySpec:
		return fmt.Sprintf("%s.TBinary", wire)
	default:
		// Not a primitive type
	}

	switch s := spec.(type) {
	case *compile.MapSpec:
		return fmt.Sprintf("%s.TMap", wire)
	case *compile.ListSpec:
		return fmt.Sprintf("%s.TList", wire)
	case *compile.SetSpec:
		return fmt.Sprintf("%s.TSet", wire)
	case *compile.TypedefSpec:
		return g.typeCode(s.Target)
	case *compile.EnumSpec:
		return fmt.Sprintf("%s.TI32", wire)
	case *compile.StructSpec:
		return fmt.Sprintf("%s.TStruct", wire)
	default:
		panic(fmt.Sprintf("unknown type spec %v (type %T)", spec, spec))
	}
}
