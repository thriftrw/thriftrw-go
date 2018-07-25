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
	commonCase := func(method string) (string, error) {
		return fmt.Sprintf("%v.Add%v(%q, %v)", encoder, method, fieldName, fieldValue), nil
	}

	root := compile.RootTypeSpec(spec)
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
	case *compile.SetSpec, *compile.ListSpec:
		// TODO: generate wrapper types for sets and use those here
		return commonCase("Reflected")

	// User-defined types
	case *compile.EnumSpec:
		return fmt.Sprintf("%v.zapObjectEncode(%v, %q)", fieldValue, encoder, fieldName), nil
	case *compile.StructSpec:
		return commonCase("Reflected")
	case *compile.TypedefSpec:
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
	// TODO: If spec is a typedef, cast to root type.
	//  fieldValue = ($rootName)($fieldValue)
	return zapObjectEncode(g, encoder, spec, fieldName, fieldValue)
}
