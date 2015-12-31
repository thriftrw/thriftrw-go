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
	// Link resolves references to other types in this TypeSpecs to actual
	// TypeSpecs from the given Scope.
	Link(scope Scope) (TypeSpec, error)

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

// typeSpecReference is a dummy TypeSpec that represents a reference to another
// TypeSpec. These will be replaced with actual TypeSpecs during the Link()
// step.
type typeSpecReference struct {
	Name string
	Line int
}

// Link replaces the typeSpecReference with an actual linked TypeSpec.
func (t *typeSpecReference) Link(scope Scope) (TypeSpec, error) {
	spec, err := scope.LookupType(t.Name)
	if err != nil {
		return nil, referenceError{Target: t.Name, Line: t.Line, Reason: err}
	}
	return spec.Link(scope)
}

// TypeCode on an unresolved typeSpecReference will cause a system panic.
func (t *typeSpecReference) TypeCode() wire.Type {
	panic(fmt.Sprintf(
		"TypeCode() requested for unresolved TypeSpec reference %v."+
			"Make sure you called Link().", t,
	))
}

// ThriftName is the name of the typeSpecReference as it appears in the Thrift
// file.
func (t *typeSpecReference) ThriftName() string {
	return t.Name
}

// compileType compiles the given AST type reference into a TypeSpec.
//
// The returned TypeSpec may need to be linked eventually.
func compileType(typ ast.Type) TypeSpec {
	if typ == nil {
		return nil
	}
	switch t := typ.(type) {
	case ast.BaseType:
		return resolveBaseType(t)
	case ast.MapType:
		return compileMapType(t)
	case ast.ListType:
		return compileListType(t)
	case ast.SetType:
		return compileSetType(t)
	case ast.TypeReference:
		return &typeSpecReference{Name: t.Name, Line: t.Line}
	default:
		panic(fmt.Sprintf("unknown type %v", typ))
	}
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
