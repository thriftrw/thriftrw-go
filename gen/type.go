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

// fieldRequired indicates whether a field is required or not.
type fieldRequired int

// Whether the field is required or not.
const (
	Optional fieldRequired = iota // default
	Required
)

// TypeDefinition generates code for the given TypeSpec.
func (g *Generator) TypeDefinition(spec compile.TypeSpec) error {
	switch s := spec.(type) {
	case *compile.EnumSpec:
		return g.enum(s)
	case *compile.StructSpec:
		return g.structure(s)
	case *compile.TypedefSpec:
		return g.typedef(s)
	default:
		panic(fmt.Sprintf("%q is not a defined type", spec.ThriftName()))
	}
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
//
// req specifies whether the reference expects the type to be always present. If
// not required, the returned type will be a reference type so that it can be
// nil.
func typeReference(spec compile.TypeSpec, req fieldRequired) string {
	name := typeName(spec)
	if (req != Required && !isReferenceType(spec)) || isStructType(spec) {
		// Prepend "*" to the result if the field is not required and the type
		// isn't a reference type.
		name = "*" + name
	}
	return name
}

// typeName returns the name of the given type, whether it's a custom type or
// native.
func typeName(spec compile.TypeSpec) string {
	switch spec {
	case compile.BoolSpec:
		return "bool"
	case compile.I8Spec:
		return "int8"
	case compile.I16Spec:
		return "int16"
	case compile.I32Spec:
		return "int32"
	case compile.I64Spec:
		return "int64"
	case compile.DoubleSpec:
		return "float64"
	case compile.StringSpec:
		return "string"
	case compile.BinarySpec:
		return "[]byte"
	default:
		// Not a primitive type. Try checking if it's a container.
	}

	switch s := spec.(type) {
	case *compile.MapSpec:
		// TODO unhashable types
		return fmt.Sprintf(
			"map[%s]%s",
			typeReference(s.KeySpec, Required),
			typeReference(s.ValueSpec, Required),
		)
	case *compile.ListSpec:
		return "[]" + typeReference(s.ValueSpec, Required)
	case *compile.SetSpec:
		// TODO unhashable types
		return fmt.Sprintf("map[%s]struct{}", typeReference(s.ValueSpec, Required))
	default:
		// Custom defined type. The reference is just the name of the type then.
		return typeDeclName(spec)
	}
}

// typeDeclName returns the name that should be used to define the given type.
//
// This panics if the given TypeSpec is not a custom user-defined type.
func typeDeclName(spec compile.TypeSpec) string {
	switch spec.(type) {
	case *compile.EnumSpec, *compile.StructSpec, *compile.TypedefSpec:
		return goCase(spec.ThriftName())
	default:
		panic(fmt.Sprintf(
			"Type %q can't have a declaration name", spec.ThriftName(),
		))
	}
}
