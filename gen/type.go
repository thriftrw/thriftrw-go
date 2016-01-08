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

// TypeDefinition generates code for the given TypeSpec.
func (g *Generator) TypeDefinition(spec compile.TypeSpec) {
	switch s := spec.(type) {
	case *compile.EnumSpec:
		g.enum(s)
	case *compile.StructSpec:
		g.structure(s)
	case *compile.TypedefSpec:
		g.typedef(s)
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

	switch spec.(type) {
	case *compile.MapSpec:
		return true
	case *compile.ListSpec:
		return true
	case *compile.SetSpec:
		return true
	}

	return false
}

// typeReference returns a string representation of a reference to the given
// type.
//
// ptr specifies whether the reference should be a pointer. It will not be a
// pointer for types that are already reference types.
func typeReference(spec compile.TypeSpec, ptr bool) (result string) {
	if ptr && !isReferenceType(spec) {
		// If requested, prepend "*" to the result if the type isn't a reference
		// type.
		defer func() {
			result = "*" + result
		}()
	}

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
		return "double64"
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
			typeReference(s.KeySpec, false),
			typeReference(s.ValueSpec, false),
		)
	case *compile.ListSpec:
		return "[]" + typeReference(s.ValueSpec, false)
	case *compile.SetSpec:
		// TODO unhashable types
		return fmt.Sprintf("map[%s]struct{}", typeReference(s.ValueSpec, false))
	default:
		// Custom defined type. The reference is just the name of the type then.
		return typeDeclName(spec)
	}
}

// typeDeclName returns the name that should be used to define the given type.
//
// This panics if the given TypeSpec is not a custom user-defined type.
func typeDeclName(spec compile.TypeSpec) string {
	switch s := spec.(type) {
	case *compile.EnumSpec:
		return goCase(s.Name)
	case *compile.StructSpec:
		return goCase(s.Name)
	case *compile.TypedefSpec:
		return goCase(s.Name)
	default:
		panic(fmt.Sprintf(
			"Type %q doesn't can't have a declaration name", spec.ThriftName(),
		))
	}
}
