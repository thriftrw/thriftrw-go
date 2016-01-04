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

// TypeSpecs for primitive Thrift types.
var (
	BoolSpec   = primitiveTypeSpec{"bool", wire.TBool}
	I8Spec     = primitiveTypeSpec{"byte", wire.TI8}
	I16Spec    = primitiveTypeSpec{"i16", wire.TI16}
	I32Spec    = primitiveTypeSpec{"i32", wire.TI32}
	I64Spec    = primitiveTypeSpec{"i64", wire.TI64}
	DoubleSpec = primitiveTypeSpec{"double", wire.TDouble}
	StringSpec = primitiveTypeSpec{"string", wire.TBinary}
	BinarySpec = primitiveTypeSpec{"binary", wire.TBinary}
)

type primitiveTypeSpec struct {
	Name string
	Code wire.Type
	// TODO(abg): We'll want to expose type annotations here
}

// TypeCode of the primitive type.
func (t primitiveTypeSpec) TypeCode() wire.Type {
	return t.Code
}

// ThriftName of the primitive type.
func (t primitiveTypeSpec) ThriftName() string {
	return t.Name
}

// Link for primitive types is a no-op because primitive types don't make any
// references.
func (t primitiveTypeSpec) Link(Scope) (TypeSpec, error) {
	return t, nil
}

// compileBaseType compiles a base type reference in the AST to a primitive
// TypeSpec.
func compileBaseType(t ast.BaseType) primitiveTypeSpec {
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
