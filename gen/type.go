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

	"github.com/thriftrw/thriftrw-go/compile"
)

// TypeDefinition generates code for the given TypeSpec.
func TypeDefinition(g Generator, spec compile.TypeSpec) error {
	switch s := spec.(type) {
	case *compile.EnumSpec:
		return enum(g, s)
	case *compile.StructSpec:
		return structure(g, s)
	case *compile.TypedefSpec:
		return typedef(g, s)
	default:
		panic(fmt.Sprintf("%q is not a defined type", spec.ThriftName()))
	}
}

// isPrimitiveType returns true if the given type is a primitive type.
//
// Note that binary is not considered a primitive type because it is
// represented as []byte in Go.
func isPrimitiveType(spec compile.TypeSpec) bool {
	switch spec {
	case compile.BoolSpec, compile.I8Spec, compile.I16Spec, compile.I32Spec,
		compile.I64Spec, compile.DoubleSpec, compile.StringSpec:
		return true
	}

	switch s := spec.(type) {
	case *compile.TypedefSpec:
		return isPrimitiveType(s.Target)
	}

	return false
}

// isReferenceType checks if the given TypeSpec represents a reference type.
//
// Sets, maps, lists, and slices are reference types.
func isReferenceType(spec compile.TypeSpec) bool {
	if spec == compile.BinarySpec {
		return true
	}

	switch s := spec.(type) {
	case *compile.MapSpec,
		*compile.ListSpec,
		*compile.SetSpec:
		return true
	case *compile.TypedefSpec:
		return isReferenceType(s.Target)
	default:
		return false
	}
}

func isStructType(spec compile.TypeSpec) bool {
	switch s := spec.(type) {
	case *compile.StructSpec:
		return true
	case *compile.TypedefSpec:
		return isStructType(s.Target)
	default:
		return false
	}
}

// typeReference returns a string representation of a reference to the given
// type.
func typeReference(g Generator, spec compile.TypeSpec) (string, error) {
	name, err := typeName(g, spec)
	if err != nil {
		return "", err
	}
	if isStructType(spec) {
		// Prepend "*" to the result if the field is not required and the type
		// isn't a reference type.
		name = "*" + name
	}
	return name, nil
}

// typeReferencePtr returns a strung representing a reference to a pointer of
// the given type. The pointer prefix is not added for types that are already
// reference types.
func typeReferencePtr(g Generator, spec compile.TypeSpec) (string, error) {
	ref, err := typeName(g, spec)
	if err != nil {
		return "", err
	}
	if !isReferenceType(spec) {
		// need * prefix for everything but map, string, and list.
		return "*" + ref, nil
	}
	return ref, nil
}

// typeName returns the name of the given type, whether it's a custom type or
// native.
func typeName(g Generator, spec compile.TypeSpec) (string, error) {
	switch spec {
	case compile.BoolSpec:
		return "bool", nil
	case compile.I8Spec:
		return "int8", nil
	case compile.I16Spec:
		return "int16", nil
	case compile.I32Spec:
		return "int32", nil
	case compile.I64Spec:
		return "int64", nil
	case compile.DoubleSpec:
		return "float64", nil
	case compile.StringSpec:
		return "string", nil
	case compile.BinarySpec:
		return "[]byte", nil
	default:
		// Not a primitive type. Try checking if it's a container.
	}

	switch s := spec.(type) {
	case *compile.MapSpec:
		// TODO unhashable types
		k, err := typeReference(g, s.KeySpec)
		if err != nil {
			return "", err
		}
		v, err := typeReference(g, s.ValueSpec)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]%s", k, v), nil
	case *compile.ListSpec:
		v, err := typeReference(g, s.ValueSpec)
		if err != nil {
			return "", err
		}
		return "[]" + v, nil
	case *compile.SetSpec:
		// TODO unhashable types
		v, err := typeReference(g, s.ValueSpec)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("map[%s]struct{}", v), nil
	case *compile.EnumSpec, *compile.StructSpec, *compile.TypedefSpec:
		return g.LookupTypeName(spec)
	default:
		panic(fmt.Sprintf("Unknown type (%T) %v", spec, spec))
	}
}
