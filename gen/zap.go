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

func zapEncoder(g Generator, spec compile.TypeSpec) string {
	root := compile.RootTypeSpec(spec)

	switch root.(type) {
	// Primitives
	case *compile.BoolSpec:
		return ("Bool")
	case *compile.I8Spec:
		return ("Int8")
	case *compile.I16Spec:
		return ("Int16")
	case *compile.I32Spec:
		return ("Int32")
	case *compile.I64Spec:
		return ("Int64")
	case *compile.DoubleSpec:
		return ("Float64")
	case *compile.StringSpec:
		return ("String")
	case *compile.BinarySpec:
		return ("Binary")

		// Containers
	case *compile.MapSpec:
		// TODO: use objects if the key is a string or array if not.
		return ("Reflected")
	case *compile.SetSpec:
		// TODO: generate wrapper types for sets and use those here
		return ("Reflected")
	case *compile.ListSpec:

		// User-defined
	case *compile.EnumSpec:
		return ("Object")
	case *compile.StructSpec:
		return ("Reflected")
	default:
	}
	panic("Wat")
}

func zapMarshaler(g Generator, spec compile.TypeSpec, fieldValue string) (string, error) {
	root := compile.RootTypeSpec(spec)

	if _, ok := spec.(*compile.TypedefSpec); ok {
		// For typedefs, cast to the root type and rely on that functionality.
		rootName, err := typeReference(g, root)
		if err != nil {
			return "", err
		}
		fieldValue = fmt.Sprintf("(%v)(%v)", rootName, fieldValue)
	}

	if isPrimitiveType(spec) {
		return fieldValue, nil
	}

	switch root.(type) {
	case *compile.MapSpec:
		// TODO: use objects if the key is a string or array if not.
		return fieldValue, nil
	case *compile.SetSpec:
		// TODO: generate wrapper types for sets and use those here
		return fieldValue, nil
	case *compile.ListSpec:
		name := "_" + g.MangleType(spec) + "_Zapper"
		if err := g.EnsureDeclared(
			`
				type <.Name> <typeReference .Type>
				<$zapcore := import "go.uber.org/zap/zapcore">
				<$v := newVar "v">
				func (<$v> <.Name>) MarshalLogArray(enc <$zapcore>.ArrayEncoder) error {
					for _, x := range <$v> {
						enc.Append<zapEncoder .Type>(<zapMarshaler .Type.ValueSpec "x">)
					}
					return nil
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
		return fmt.Sprintf("(%v)(%v)", name, fieldValue), nil
	case *compile.StructSpec:
		return fieldValue, nil
	default:
	}
	panic("Wat")
}

func zapMarshalerPtr(g Generator, spec compile.TypeSpec, fieldValue string) (string, error) {
	if isPrimitiveType(spec) {
		fieldValue = "*" + fieldValue
	}
	return zapMarshaler(g, spec, fieldValue)
}
