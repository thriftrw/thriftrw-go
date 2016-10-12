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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

func mustCompileTypeReference(t *testing.T, typ ast.Type) TypeSpec {
	result, err := compileTypeReference(typ)
	require.NoError(t, err, "failed to compile type reference %v: %v", typ, err)
	return result
}

// mustLink links the given TypeSpec inside the given scope, failing the test
// if the spec does not link successfully.
func mustLink(t *testing.T, spec TypeSpec, scope Scope) TypeSpec {
	result, err := spec.Link(scope)
	require.NoError(t, err, "failed to link %v inside %v", spec, scope)
	return result
}

func TestCompileTypeWithNil(t *testing.T) {
	// make sure compileType(nil) doesn't explode.
	assert.Nil(t, mustCompileTypeReference(t, nil))
}

func TestResolveBaseType(t *testing.T) {
	tests := []struct {
		input    ast.BaseType
		wireType wire.Type
	}{
		{ast.BaseType{ID: ast.BoolTypeID}, wire.TBool},
		{ast.BaseType{ID: ast.I8TypeID}, wire.TI8},
		{ast.BaseType{ID: ast.I16TypeID}, wire.TI16},
		{ast.BaseType{ID: ast.I32TypeID}, wire.TI32},
		{ast.BaseType{ID: ast.I64TypeID}, wire.TI64},
		{ast.BaseType{ID: ast.DoubleTypeID}, wire.TDouble},
		{ast.BaseType{ID: ast.StringTypeID}, wire.TBinary},
		{ast.BaseType{ID: ast.BinaryTypeID}, wire.TBinary},
	}

	for _, tt := range tests {
		spec, err := compileTypeReference(tt.input)
		if !assert.NoError(t, err) {
			continue
		}

		linked, err := spec.Link(defaultScope)
		assert.NoError(t, err)
		assert.Equal(t, tt.wireType, spec.TypeCode())
		assert.Equal(t, tt.wireType, linked.TypeCode())
	}
}

func TestResolveInvalidBaseType(t *testing.T) {
	assert.Panics(t, func() {
		compileTypeReference(ast.BaseType{ID: ast.BaseTypeID(42)})
	})
}

func TestLinkTypeReference(t *testing.T) {
	foo := &StructSpec{
		Name: "Foo",
		Type: ast.StructType,
		Fields: FieldGroup{
			{ID: 1, Name: "value", Type: &I32Spec{}},
		},
	}

	tests := []struct {
		desc     string
		scope    Scope
		name     string
		expected TypeSpec
	}{
		{"simple type lookup", scope("Foo", foo), "Foo", foo},
		{
			"included type lookup",
			scope("bar", scope("Foo", foo)),
			"bar.Foo",
			foo,
		},
	}

	for _, tt := range tests {
		expected := mustLink(t, tt.expected, defaultScope)
		scope := scopeOrDefault(tt.scope)

		spec, err := typeSpecReference(ast.TypeReference{Name: tt.name}).Link(scope)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, expected, spec, tt.desc)
		}
	}
}

func TestLinkTypeReferenceFailure(t *testing.T) {
	tests := []struct {
		desc     string
		scope    Scope
		name     string
		messages []string
	}{
		{
			"unknown identifier",
			scope("some_module"),
			"foo",
			[]string{`could not resolve reference "foo" in "some_module"`},
		},
		{
			"unknown module",
			scope("some_module"),
			"shared.UUID",
			[]string{
				`could not resolve reference "shared.UUID" in "some_module"`,
				`unknown module "shared"`,
			},
		},
		{
			"unknown identifier in included module",
			scope("foo", "shared", scope("shared")),
			"shared.UUID",
			[]string{
				`could not resolve reference "shared.UUID" in "foo"`,
				`could not resolve reference "UUID" in "shared"`,
			},
		},
		{
			"unknown identifier in included module of included module",
			scope("foo", "bar", scope("bar", "baz", scope("baz"))),
			"bar.baz.Qux",
			[]string{
				`could not resolve reference "bar.baz.Qux" in "foo"`,
				`could not resolve reference "baz.Qux" in "bar"`,
				`could not resolve reference "Qux" in "baz"`,
			},
		},
	}

	for _, tt := range tests {
		scope := scopeOrDefault(tt.scope)
		_, err := typeSpecReference(ast.TypeReference{Name: tt.name}).Link(scope)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}

func TestRootTypeSpec(t *testing.T) {
	tests := []struct {
		desc string
		give TypeSpec
	}{
		// Primitives
		{desc: "BoolSpec", give: &BoolSpec{}},
		{desc: "I8Spec", give: &I8Spec{}},
		{desc: "I16Spec", give: &I16Spec{}},
		{desc: "I32Spec", give: &I32Spec{}},
		{desc: "I64Spec", give: &I64Spec{}},
		{desc: "DoubleSpec", give: &DoubleSpec{}},
		{desc: "StringSpec", give: &StringSpec{}},
		{desc: "BinarySpec", give: &BinarySpec{}},

		// Containers
		{
			desc: "ListSpec",
			give: &ListSpec{
				ValueSpec: &TypedefSpec{Name: "UUID", Target: &StringSpec{}},
			},
		},
		{
			desc: "SetSpec",
			give: &SetSpec{
				ValueSpec: &TypedefSpec{Name: "UUID", Target: &StringSpec{}},
			},
		},
		{
			desc: "MapSpec",
			give: &MapSpec{
				KeySpec:   &StringSpec{},
				ValueSpec: &TypedefSpec{Name: "UUID", Target: &StringSpec{}},
			},
		},

		// Structs and enums
		{
			desc: "StructSpec",
			give: &StructSpec{
				Name: "Foo",
				Type: ast.StructType,
				Fields: FieldGroup{
					{
						ID:   1,
						Name: "a",
						Type: &StringSpec{},
					},
					{
						ID:   2,
						Name: "b",
						Type: &TypedefSpec{Name: "UUID", Target: &StringSpec{}},
					},
				},
			},
		},
		{
			desc: "EnumSpec",
			give: &EnumSpec{
				Name: "Numbers",
				Items: []EnumItem{
					{Name: "One", Value: 1},
					{Name: "Two", Value: 2},
				},
			},
		},
	}

	for _, tt := range tests {
		give, err := tt.give.Link(defaultScope)
		if !assert.NoError(t, err, "%v: failed to link", tt.desc) {
			continue
		}

		spec := give
		for i := 0; i < 10; i++ {
			assert.Equal(t, give, RootTypeSpec(spec), "%v: level %d", tt.desc, i)

			spec = &TypedefSpec{Name: fmt.Sprintf("foo%d", i), Target: spec}
			spec, err = spec.Link(defaultScope)
			if !assert.NoError(t, err, "%v: failed to link level %d", tt.desc, i) {
				break
			}
		}
	}
}
