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

package compile

import (
	"fmt"

	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/wire"
)

// TypeSpec contains information about Thrift types.
type TypeSpec interface {
	Unit

	// TypeCode is the wire-level Thrift Type associated with this Type.
	TypeCode() wire.Type
}

// DefinedTypeSpec contains information about types that map directly to types
// defined in the Thrift file.
type DefinedTypeSpec interface {
	TypeSpec

	// ThriftName is the name of the given object as it appears in the Thrift
	// file.
	ThriftName() string
}

// ResolveType resolves the given AST type to a specific compiled TypeSpec.
func ResolveType(typ ast.Type, scope Scope) (TypeSpec, error) {
	var spec TypeSpec

	switch t := typ.(type) {
	case ast.BaseType:
		spec = resolveBaseType(t)
	case ast.MapType:
		spec = NewMapSpec(t)
	case ast.ListType:
		spec = NewListSpec(t)
	case ast.SetType:
		spec = NewSetSpec(t)
	case ast.TypeReference:
		var err error
		spec, err = scope.LookupType(t.Name)
		if err != nil {
			return nil, referenceError{
				Target: t.Name,
				Line:   t.Line,
				Reason: err,
			}
		}
	default:
		panic(fmt.Sprintf("unknown type %v", typ))
	}

	if err := spec.Compile(scope); err != nil {
		return nil, err
	}
	return spec, nil
}

func resolveBaseType(t ast.BaseType) TypeSpec {
	switch t.ID {
	case ast.BoolTypeID:
		return BoolSpec
	case ast.I8TypeID:
		return I8Spec
	case ast.I16TypeID:
		return I16Spec
	case ast.I32TypeID:
		return I32Spec
	case ast.I64TypeID:
		return I64Spec
	case ast.DoubleTypeID:
		return DoubleSpec
	case ast.StringTypeID:
		return StringSpec
	case ast.BinaryTypeID:
		return BinarySpec
	default:
		panic(fmt.Sprintf("unknown base type %v", t))
	}
}
