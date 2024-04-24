// Copyright (c) 2024 Uber Technologies, Inc.
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

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

type (
	// BoolSpec is the TypeSpec for bool types in a Thrift file.
	BoolSpec struct {
		nativeThriftType

		Annotations Annotations
	}

	// I8Spec is the TypeSpec for i8/byte types in a Thrift file.
	I8Spec struct {
		nativeThriftType

		Annotations Annotations
	}

	// I16Spec is the TypeSpec for i16 types in a Thrift file.
	I16Spec struct {
		nativeThriftType

		Annotations Annotations
	}

	// I32Spec is the TypeSpec for i32 types in a Thrift file.
	I32Spec struct {
		nativeThriftType

		Annotations Annotations
	}

	// I64Spec is the TypeSpec for i64 types in a Thrift file.
	I64Spec struct {
		nativeThriftType

		Annotations Annotations
	}

	// DoubleSpec is the TypeSpec for double types in a Thrift file.
	DoubleSpec struct {
		nativeThriftType

		Annotations Annotations
	}

	// StringSpec is the TypeSpec for string types in a Thrift file.
	StringSpec struct {
		nativeThriftType

		Annotations Annotations
	}

	// BinarySpec is the TypeSpec for binary types in a Thrift file.
	BinarySpec struct {
		nativeThriftType

		Annotations Annotations
	}
)

// TypeCode returns TBool.
func (*BoolSpec) TypeCode() wire.Type { return wire.TBool }

// TypeCode returns TI8.
func (*I8Spec) TypeCode() wire.Type { return wire.TI8 }

// TypeCode returns TI16.
func (*I16Spec) TypeCode() wire.Type { return wire.TI16 }

// TypeCode returns TI32.
func (*I32Spec) TypeCode() wire.Type { return wire.TI32 }

// TypeCode returns TI64.
func (*I64Spec) TypeCode() wire.Type { return wire.TI64 }

// TypeCode returns TDouble.
func (*DoubleSpec) TypeCode() wire.Type { return wire.TDouble }

// TypeCode returns TBinary since Thrift strings are sent as binary over the
// wire.
func (*StringSpec) TypeCode() wire.Type { return wire.TBinary }

// TypeCode returns TBinary.
func (*BinarySpec) TypeCode() wire.Type { return wire.TBinary }

// ThriftName returns "bool".
func (*BoolSpec) ThriftName() string { return "bool" }

// ThriftName returns "byte".
func (*I8Spec) ThriftName() string { return "byte" }

// ThriftName returns "i16".
func (*I16Spec) ThriftName() string { return "i16" }

// ThriftName returns "i32".
func (*I32Spec) ThriftName() string { return "i32" }

// ThriftName returns "i64".
func (*I64Spec) ThriftName() string { return "i64" }

// ThriftName returns "double".
func (*DoubleSpec) ThriftName() string { return "double" }

// ThriftName returns "string".
func (*StringSpec) ThriftName() string { return "string" }

// ThriftName returns "binary".
func (*BinarySpec) ThriftName() string { return "binary" }

// Link is a no-op for primitives.
func (t *BoolSpec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *I8Spec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *I16Spec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *I32Spec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *I64Spec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *DoubleSpec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *StringSpec) Link(Scope) (TypeSpec, error) { return t, nil }

// Link is a no-op for primitives.
func (t *BinarySpec) Link(Scope) (TypeSpec, error) { return t, nil }

// ForEachTypeReference is a no-op for primitives.
func (*BoolSpec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*I8Spec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*I16Spec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*I32Spec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*I64Spec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*DoubleSpec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*StringSpec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ForEachTypeReference is a no-op for primitives.
func (*BinarySpec) ForEachTypeReference(func(TypeSpec) error) error { return nil }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *BoolSpec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *I8Spec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *I16Spec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *I32Spec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *I64Spec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *DoubleSpec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *StringSpec) ThriftAnnotations() Annotations { return t.Annotations }

// ThriftAnnotations returns the Thrift annotations specified with the
// reference to this type.
func (t *BinarySpec) ThriftAnnotations() Annotations { return t.Annotations }

// compileBaseType compiles a base type reference in the AST to a primitive
// TypeSpec.
func compileBaseType(t ast.BaseType) (TypeSpec, error) {
	annots, err := compileAnnotations(t.Annotations)
	if err != nil {
		return nil, err
	}

	switch t.ID {
	case ast.BoolTypeID:
		return &BoolSpec{Annotations: annots}, nil
	case ast.I8TypeID:
		return &I8Spec{Annotations: annots}, nil
	case ast.I16TypeID:
		return &I16Spec{Annotations: annots}, nil
	case ast.I32TypeID:
		return &I32Spec{Annotations: annots}, nil
	case ast.I64TypeID:
		return &I64Spec{Annotations: annots}, nil
	case ast.DoubleTypeID:
		return &DoubleSpec{Annotations: annots}, nil
	case ast.StringTypeID:
		return &StringSpec{Annotations: annots}, nil
	case ast.BinaryTypeID:
		return &BinarySpec{Annotations: annots}, nil
	default:
		panic(fmt.Sprintf("unknown base type %v", t))
	}
}
