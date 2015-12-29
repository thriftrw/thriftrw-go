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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/wire"
)

func TestCompileList(t *testing.T) {
	tests := []struct {
		desc   string
		input  ast.ListType
		scope  Scope
		output *ListSpec
	}{
		{
			"list<i32>",
			ast.ListType{ValueType: ast.BaseType{ID: ast.I32TypeID}},
			nil,
			&ListSpec{ValueSpec: I32Spec},
		},
		{
			"typedef list<i32> Foo; list<Foo>",
			ast.ListType{ValueType: ast.TypeReference{Name: "Foo"}},
			scope("Foo", &TypedefSpec{
				Name:   "Foo",
				Target: &ListSpec{ValueSpec: I32Spec},
			}),
			&ListSpec{
				ValueSpec: &TypedefSpec{
					Name:   "Foo",
					Target: &ListSpec{ValueSpec: I32Spec},
				},
			},
		},
	}
	for _, tt := range tests {
		output := mustLink(t, tt.output, scope())

		scope := scopeOrDefault(tt.scope)
		spec, err := compileListType(tt.input).Link(scope)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, wire.TList, spec.TypeCode(), tt.desc)
			assert.Equal(t, output, spec, tt.desc)
		}
	}
}

func TestCompileMap(t *testing.T) {
	tests := []struct {
		desc   string
		input  ast.MapType
		scope  Scope
		output *MapSpec
	}{
		{
			"typedef string Key; typedef i64 Value; map<Key, Value>",
			ast.MapType{
				KeyType:   ast.TypeReference{Name: "Key"},
				ValueType: ast.TypeReference{Name: "Value"},
			},
			scope(
				"Key", &TypedefSpec{Name: "Key", Target: StringSpec},
				"Value", &TypedefSpec{Name: "Value", Target: I64Spec},
			),
			&MapSpec{
				KeySpec:   &TypedefSpec{Name: "Key", Target: StringSpec},
				ValueSpec: &TypedefSpec{Name: "Value", Target: I64Spec},
			},
		},
	}

	for _, tt := range tests {
		output := mustLink(t, tt.output, scope())

		scope := scopeOrDefault(tt.scope)
		spec, err := compileMapType(tt.input).Link(scope)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, wire.TMap, spec.TypeCode(), tt.desc)
			assert.Equal(t, output, spec, tt.desc)
		}
	}
}

func TestCompileSet(t *testing.T) {
	tests := []struct {
		desc   string
		input  ast.SetType
		scope  Scope
		output *SetSpec
	}{
		{
			"typedef string Foo; set<Foo>",
			ast.SetType{ValueType: ast.TypeReference{Name: "Foo"}},
			scope("Foo", &TypedefSpec{Name: "Foo", Target: StringSpec}),
			&SetSpec{ValueSpec: &TypedefSpec{Name: "Foo", Target: StringSpec}},
		},
	}

	for _, tt := range tests {
		output := mustLink(t, tt.output, scope())

		scope := scopeOrDefault(tt.scope)
		spec, err := compileSetType(tt.input).Link(scope)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, wire.TSet, spec.TypeCode(), tt.desc)
			assert.Equal(t, output, spec, tt.desc)
		}
	}
}
