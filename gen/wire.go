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

// WireGenerator is responsible for generating code that knows how to convert
// between Thrift types and their Value representations.
type WireGenerator struct {
	mapG  mapGenerator
	setG  setGenerator
	listG listGenerator

	enumG    enumGenerator
	structG  structGenerator
	typedefG typedefGenerator
}

// ToWire generates an expression of type (Value, error) object containing the
// wire representation of the variable $varName of type $spec or an error.
func (w *WireGenerator) ToWire(g Generator, spec compile.TypeSpec, varName string) (string, error) {
	wire := g.Import("go.uber.org/thriftrw/wire")
	switch s := spec.(type) {
	case *compile.BoolSpec:
		return fmt.Sprintf("%s.NewValueBool(%s), error(nil)", wire, varName), nil
	case *compile.I8Spec:
		return fmt.Sprintf("%s.NewValueI8(%s), error(nil)", wire, varName), nil
	case *compile.I16Spec:
		return fmt.Sprintf("%s.NewValueI16(%s), error(nil)", wire, varName), nil
	case *compile.I32Spec:
		return fmt.Sprintf("%s.NewValueI32(%s), error(nil)", wire, varName), nil
	case *compile.I64Spec:
		return fmt.Sprintf("%s.NewValueI64(%s), error(nil)", wire, varName), nil
	case *compile.DoubleSpec:
		return fmt.Sprintf("%s.NewValueDouble(%s), error(nil)", wire, varName), nil
	case *compile.StringSpec:
		return fmt.Sprintf("%s.NewValueString(%s), error(nil)", wire, varName), nil
	case *compile.BinarySpec:
		return fmt.Sprintf("%s.NewValueBinary(%s), error(nil)", wire, varName), nil
	case *compile.MapSpec:
		mapItemList, err := w.mapG.ItemList(g, s)
		if err != nil {
			return "", err
		}

		return g.TextTemplate(
			`<.Wire>.NewValueMap(<.MapItemList>(<.Name>)), error(nil)`,
			struct {
				Wire        string
				Name        string
				Spec        *compile.MapSpec
				MapItemList string
			}{Wire: wire, Name: varName, Spec: s, MapItemList: mapItemList},
		)
	case *compile.ListSpec:
		valueList, err := w.listG.ValueList(g, s)
		if err != nil {
			return "", err
		}

		return g.TextTemplate(
			`<.Wire>.NewValueList(<.ValueList>(<.Name>)), error(nil)`,
			struct {
				Wire      string
				Name      string
				Spec      *compile.ListSpec
				ValueList string
			}{Wire: wire, Name: varName, Spec: s, ValueList: valueList},
		)
	case *compile.SetSpec:
		valueList, err := w.setG.ValueList(g, s)
		if err != nil {
			return "", err
		}

		return g.TextTemplate(
			`<.Wire>.NewValueSet(<.ValueList>(<.Name>)), error(nil)`,
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

// ToWirePtr is the same as ToWire expect `varName` is expected to be a
// reference to a value of the given type.
func (w *WireGenerator) ToWirePtr(g Generator, spec compile.TypeSpec, varName string) (string, error) {
	switch spec.(type) {
	case *compile.BoolSpec, *compile.I8Spec, *compile.I16Spec, *compile.I32Spec,
		*compile.I64Spec, *compile.DoubleSpec, *compile.StringSpec:
		return w.ToWire(g, spec, fmt.Sprintf("*(%s)", varName))
	default:
		// Everything else is either a reference type or has a ToWire method
		// on it that does automatic dereferencing.
		return w.ToWire(g, spec, varName)
	}
}

// FromWire generates an expression of type ($spec, error) which reads the Value
// at $value into a $spec.
func (w *WireGenerator) FromWire(g Generator, spec compile.TypeSpec, value string) (string, error) {
	switch s := spec.(type) {
	case *compile.BoolSpec:
		return fmt.Sprintf("%s.GetBool(), error(nil)", value), nil
	case *compile.I8Spec:
		return fmt.Sprintf("%s.GetI8(), error(nil)", value), nil
	case *compile.I16Spec:
		return fmt.Sprintf("%s.GetI16(), error(nil)", value), nil
	case *compile.I32Spec:
		return fmt.Sprintf("%s.GetI32(), error(nil)", value), nil
	case *compile.I64Spec:
		return fmt.Sprintf("%s.GetI64(), error(nil)", value), nil
	case *compile.DoubleSpec:
		return fmt.Sprintf("%s.GetDouble(), error(nil)", value), nil
	case *compile.StringSpec:
		return fmt.Sprintf("%s.GetString(), error(nil)", value), nil
	case *compile.BinarySpec:
		return fmt.Sprintf("%s.GetBinary(), error(nil)", value), nil
	case *compile.MapSpec:
		reader, err := w.mapG.Reader(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s.GetMap())", reader, value), nil
	case *compile.ListSpec:
		reader, err := w.listG.Reader(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s.GetList())", reader, value), nil
	case *compile.SetSpec:
		reader, err := w.setG.Reader(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s.GetSet())", reader, value), nil
	case *compile.TypedefSpec:
		reader, err := w.typedefG.Reader(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", reader, value), nil
	case *compile.EnumSpec:
		reader, err := w.enumG.Reader(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", reader, value), nil
	case *compile.StructSpec:
		reader, err := w.structG.Reader(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", reader, value), nil
	default:
		panic(fmt.Sprintf("Unknown TypeSpec (%T) %v", spec, spec))
	}
}

// FromWirePtr generates a string assigning the given Value to the given lhs,
// which is a pointer to a value of the given type.
//
// A variable err of type error MUST be in scope and will be assigned the
// parse error, if any.
func (w *WireGenerator) FromWirePtr(g Generator, spec compile.TypeSpec, lhs string, value string) (string, error) {
	if !isPrimitiveType(spec) {
		// Everything else can be assigned to directly.
		out, err := w.FromWire(g, spec, value)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s, err = %s", lhs, out), err
	}
	return g.TextTemplate(
		`
			<- $x := newVar "x" ->
			var <$x> <typeReference .Spec>
			<$x>, err = <fromWire .Spec .Value>
			<.LHS> = &<$x ->
			`,
		struct {
			Spec  compile.TypeSpec
			LHS   string
			Value string
		}{Spec: spec, LHS: lhs, Value: value},
	)
}

// TypeCode gets an expression of type 'wire.Type' that represents the
// over-the-wire type code for the given TypeSpec.
func TypeCode(g Generator, spec compile.TypeSpec) string {
	wire := g.Import("go.uber.org/thriftrw/wire")
	spec = compile.RootTypeSpec(spec)

	switch spec.(type) {
	case *compile.BoolSpec:
		return fmt.Sprintf("%s.TBool", wire)
	case *compile.I8Spec:
		return fmt.Sprintf("%s.TI8", wire)
	case *compile.I16Spec:
		return fmt.Sprintf("%s.TI16", wire)
	case *compile.I32Spec:
		return fmt.Sprintf("%s.TI32", wire)
	case *compile.I64Spec:
		return fmt.Sprintf("%s.TI64", wire)
	case *compile.DoubleSpec:
		return fmt.Sprintf("%s.TDouble", wire)
	case *compile.StringSpec, *compile.BinarySpec:
		return fmt.Sprintf("%s.TBinary", wire)
	case *compile.MapSpec:
		return fmt.Sprintf("%s.TMap", wire)
	case *compile.ListSpec:
		return fmt.Sprintf("%s.TList", wire)
	case *compile.SetSpec:
		return fmt.Sprintf("%s.TSet", wire)
	case *compile.EnumSpec:
		return fmt.Sprintf("%s.TI32", wire)
	case *compile.StructSpec:
		return fmt.Sprintf("%s.TStruct", wire)
	default:
		panic(fmt.Sprintf("unknown type spec %v (type %T)", spec, spec))
	}
}
