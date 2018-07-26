// Copyright (c) 2018 Uber Technologies, Inc.
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

func zapObjectEncode(
	g Generator,
	encoder string,
	spec compile.TypeSpec,
	fieldName string,
	fieldValue string,
) (string, error) {
	root := compile.RootTypeSpec(spec)

	if _, ok := spec.(*compile.TypedefSpec); ok {
		// For typedefs, cast to the root type and rely on that functionality.
		rootName, err := typeReference(g, root)
		if err != nil {
			return "", err
		}
		fieldValue = fmt.Sprintf("(%v)(%v)", rootName, fieldValue)
	}

	commonCase := func(method string) (string, error) {
		return fmt.Sprintf("%v.Add%v(%q, %v)", encoder, method, fieldName, fieldValue), nil
	}

	switch root.(type) {
	// Primitives
	case *compile.BoolSpec:
		return commonCase("Bool")
	case *compile.I8Spec:
		return commonCase("Int8")
	case *compile.I16Spec:
		return commonCase("Int16")
	case *compile.I32Spec:
		return commonCase("Int32")
	case *compile.I64Spec:
		return commonCase("Int64")
	case *compile.DoubleSpec:
		return commonCase("Float64")
	case *compile.StringSpec:
		return commonCase("String")
	case *compile.BinarySpec:
		return commonCase("Binary")

	// Containers
	case *compile.MapSpec:
		// TODO: use objects if the key is a string or array if not.
		return commonCase("Reflected")
	case *compile.SetSpec:
		// TODO: generate wrapper types for sets and use those here
		return commonCase("Reflected")
	case *compile.ListSpec:
		name := "_" + g.MangleType(spec) + "_Zapper"
		if err := g.EnsureDeclared(
			`
				type <.Name> <typeReference .Type>
				<$zapcore := import "go.uber.org/zap/zapcore">
				<$v := newVar "v">
				func (<$v> <.Name>) MarshalLogArray(enc <$zapcore>.ArrayEncoder) {
					for _, x := range <$v> {
						<zapObjectEncode "enc" .Type.ValueSpec "someName" "x">
					}
				}
				`, struct {
				Name string
				Type compile.TypeSpec
			}{
				Name: name,
				Type: root,
			},
		); err != nil {
			return "", err
		}
		// TODO: generate wrapper types for sets and use those here
		return fmt.Sprintf("%v.AddArray(%q, (%v)(%v))", encoder, fieldName, name, fieldValue), nil
	// User-defined types
	case *compile.EnumSpec:
		return fmt.Sprintf("%v.AddObject(%q, %v)", encoder, fieldName, fieldValue), nil
	case *compile.StructSpec:
		return commonCase("Reflected")
	default:
		panic("Wat")
	}
}

func zapObjectEncodePtr(
	g Generator,
	encoder string,
	spec compile.TypeSpec,
	fieldName string,
	fieldValue string,
) (string, error) {
	if isPrimitiveType(spec) {
		fieldValue = "*" + fieldValue
	}
	return zapObjectEncode(g, encoder, spec, fieldName, fieldValue)
}
