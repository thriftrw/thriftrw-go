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

	// ThriftName is the name of the type as it appears in the Thrift file.
	ThriftName() string
}

// typeSpecReference is a dummy TypeSpec that represents a reference to another
// TypeSpec. These will be replaced with actual TypeSpecs during the Link()
// step.
type typeSpecReference ast.TypeReference

// Link replaces the typeSpecReference with an actual linked TypeSpec.
func (r typeSpecReference) Link(scope Scope) (TypeSpec, error) {
	src := ast.TypeReference(r)
	t, err := scope.LookupType(src.Name)
	if err == nil {
		return t.Link(scope)
	}

	mname, iname := splitInclude(src.Name)
	if len(mname) == 0 {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	includedScope, err := scope.LookupInclude(mname)
	if err != nil {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    unrecognizedModuleError{Name: mname, Reason: err},
		}
	}

	t, err = typeSpecReference{Name: iname}.Link(includedScope)
	if err != nil {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	return t, nil
}

// TypeCode on an unresolved typeSpecReference will cause a system panic.
func (r typeSpecReference) TypeCode() wire.Type {
	panic(fmt.Sprintf(
		"TypeCode() requested for unresolved TypeSpec reference %v."+
			"Make sure you called Link().", r,
	))
}

// ThriftName is the name of the typeSpecReference as it appears in the Thrift
// file.
func (r typeSpecReference) ThriftName() string {
	return r.Name
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
		return compileBaseType(t)
	case ast.MapType:
		return compileMapType(t)
	case ast.ListType:
		return compileListType(t)
	case ast.SetType:
		return compileSetType(t)
	case ast.TypeReference:
		return typeSpecReference(t)
	default:
		panic(fmt.Sprintf("unknown type %v", typ))
	}
}
