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

// isHashable returns true if the given type is considered hashable by
// thriftrw.
//
// Only primitive types, enums, and typedefs of other hashable types are
// considered hashable.
func isHashable(t compile.TypeSpec) bool {
	return isPrimitiveType(t)
}

// setUsesMap returns true if the given set type is not annotated with
// (go.type = "slice") and the value of the set is considered hashable
// by thriftrw.
func setUsesMap(spec *compile.SetSpec) bool {
	return (spec.Annotations[goTypeKey] != sliceType) && isHashable(spec.ValueSpec)
}

// isPrimitiveType returns true if the given type is a primitive type.
// Primitive types, enums, and typedefs of primitive types are considered
// primitive.
//
// Note that binary is not considered a primitive type because it is
// represented as []byte in Go.
func isPrimitiveType(spec compile.TypeSpec) bool {
	spec = compile.RootTypeSpec(spec)
	switch spec.(type) {
	case *compile.BoolSpec, *compile.I8Spec, *compile.I16Spec, *compile.I32Spec,
		*compile.I64Spec, *compile.DoubleSpec, *compile.StringSpec:
		return true
	}

	_, isEnum := spec.(*compile.EnumSpec)
	return isEnum
}

// isReferenceType checks if the given TypeSpec represents a reference type.
//
// Sets, maps, lists, and slices are reference types.
func isReferenceType(spec compile.TypeSpec) bool {
	spec = compile.RootTypeSpec(spec)
	if _, ok := spec.(*compile.BinarySpec); ok {
		return true
	}

	switch spec.(type) {
	case *compile.MapSpec, *compile.ListSpec, *compile.SetSpec:
		return true
	default:
		return false
	}
}

func isStructType(spec compile.TypeSpec) bool {
	spec = compile.RootTypeSpec(spec)
	_, isStruct := spec.(*compile.StructSpec)
	return isStruct
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
	switch s := spec.(type) {
	case *compile.BoolSpec:
		return "bool", nil
	case *compile.I8Spec:
		return "int8", nil
	case *compile.I16Spec:
		return "int16", nil
	case *compile.I32Spec:
		return "int32", nil
	case *compile.I64Spec:
		return "int64", nil
	case *compile.DoubleSpec:
		return "float64", nil
	case *compile.StringSpec:
		return "string", nil
	case *compile.BinarySpec:
		return "[]byte", nil
	case *compile.MapSpec:
		k, err := typeReference(g, s.KeySpec)
		if err != nil {
			return "", err
		}
		v, err := typeReference(g, s.ValueSpec)
		if err != nil {
			return "", err
		}
		if !isHashable(s.KeySpec) {
			// unhashable type
			return fmt.Sprintf("[]struct{Key %s; Value %s}", k, v), nil
		}
		return fmt.Sprintf("map[%s]%s", k, v), nil
	case *compile.ListSpec:
		v, err := typeReference(g, s.ValueSpec)
		if err != nil {
			return "", err
		}
		return "[]" + v, nil
	case *compile.SetSpec:
		v, err := typeReference(g, s.ValueSpec)
		if err != nil {
			return "", err
		}
		// not annotated to be slice and hashable value type
		if setUsesMap(s) {
			return fmt.Sprintf("map[%s]struct{}", v), nil
		}
		return fmt.Sprintf("[]%s", v), nil
	case *compile.EnumSpec, *compile.StructSpec, *compile.TypedefSpec:
		return g.LookupTypeName(spec)
	default:
		panic(fmt.Sprintf("Unknown type (%T) %v", spec, spec))
	}
}

func equalsFuncName(g Generator, spec compile.TypeSpec) string {
	return fmt.Sprintf("_%s_Equals", g.MangleType(spec))
}

func equalsPtrFuncName(g Generator, spec compile.TypeSpec) string {
	return fmt.Sprintf("_%s_EqualsPtr", g.MangleType(spec))
}

func readerFuncName(g Generator, spec compile.TypeSpec) string {
	return fmt.Sprintf("_%s_Read", g.MangleType(spec))
}

func valueListName(g Generator, spec compile.TypeSpec) string {
	return fmt.Sprintf("_%s_ValueList", g.MangleType(spec))
}

// zapperName returns the name that should be used for wrapper types that
// implement zap.ObjectMarshaler or zap.ArrayMarshaler for the provided
// Thrift type.
func zapperName(g Generator, spec compile.TypeSpec) string {
	return fmt.Sprintf("_%s_Zapper", g.MangleType(spec))
}

// canBeConstant returns true if the given type can be a constant.
func canBeConstant(t compile.TypeSpec) bool {
	// Only primitives can use const declarations. Everything else has to be a
	// `var` declaration.
	return isPrimitiveType(t)
}
