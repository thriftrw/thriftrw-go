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

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

func TestCompileContainers(t *testing.T) {
	tests := []struct {
		desc  string
		give  ast.Type
		scope Scope

		want             TypeSpec
		wantTypeCode     wire.Type
		wantCompileError []string
		wantLinkError    []string
	}{
		// list
		{
			desc:         "list<i32>",
			give:         ast.ListType{ValueType: ast.BaseType{ID: ast.I32TypeID}},
			want:         &ListSpec{ValueSpec: &I32Spec{}},
			wantTypeCode: wire.TList,
		},
		{
			desc: "typedef list<i32> Foo; list<Foo>",
			give: ast.ListType{ValueType: ast.TypeReference{Name: "Foo"}},
			scope: scope("Foo", &TypedefSpec{
				Name:   "Foo",
				Target: &ListSpec{ValueSpec: &I32Spec{}},
			}),
			want: &ListSpec{
				ValueSpec: &TypedefSpec{
					Name:   "Foo",
					Target: &ListSpec{ValueSpec: &I32Spec{}},
				},
			},
			wantTypeCode: wire.TList,
		},
		{
			desc: "list<i64> (obfuscate)",
			give: ast.ListType{
				ValueType: ast.BaseType{ID: ast.I64TypeID},
				Annotations: []*ast.Annotation{
					{
						Name:  "obfuscate",
						Value: "",
					},
				},
			},
			want: &ListSpec{
				ValueSpec:   &I64Spec{},
				Annotations: Annotations{"obfuscate": ""},
			},
			wantTypeCode: wire.TList,
		},
		{
			desc: `list<string (a = "b")> (c = "d")`,
			give: ast.ListType{
				ValueType: ast.BaseType{
					ID: ast.StringTypeID,
					Annotations: []*ast.Annotation{
						{Name: "a", Value: "b"},
					},
				},
				Annotations: []*ast.Annotation{
					{Name: "c", Value: "d"},
				},
			},
			want: &ListSpec{
				ValueSpec:   &StringSpec{Annotations: Annotations{"a": "b"}},
				Annotations: Annotations{"c": "d"},
			},
			wantTypeCode: wire.TList,
		},
		{
			desc: `list<string> (c, c = "d")`,
			give: ast.ListType{
				ValueType: ast.BaseType{ID: ast.StringTypeID},
				Annotations: []*ast.Annotation{
					{Name: "c"},
					{Name: "c", Value: "d"},
				},
			},
			wantCompileError: []string{`annotation conflict: the name "c" has already been used`},
		},
		{
			desc:          "list<Foo>",
			give:          ast.ListType{ValueType: ast.TypeReference{Name: "Foo"}},
			wantLinkError: []string{`could not resolve reference "Foo"`},
		},

		// map
		{
			desc: "typedef string Key; typedef i64 Value; map<Key, Value>",
			give: ast.MapType{
				KeyType:   ast.TypeReference{Name: "Key"},
				ValueType: ast.TypeReference{Name: "Value"},
			},
			scope: scope(
				"Key", &TypedefSpec{Name: "Key", Target: &StringSpec{}},
				"Value", &TypedefSpec{Name: "Value", Target: &I64Spec{}},
			),
			want: &MapSpec{
				KeySpec:   &TypedefSpec{Name: "Key", Target: &StringSpec{}},
				ValueSpec: &TypedefSpec{Name: "Value", Target: &I64Spec{}},
			},
			wantTypeCode: wire.TMap,
		},
		{
			desc: "map<i8, Foo>",
			give: ast.MapType{
				KeyType:   ast.BaseType{ID: ast.I8TypeID},
				ValueType: ast.TypeReference{Name: "Foo"},
			},
			wantLinkError: []string{`could not resolve reference "Foo"`},
		},
		{
			desc: "map<Foo, i8>",
			give: ast.MapType{
				KeyType:   ast.TypeReference{Name: "Foo"},
				ValueType: ast.BaseType{ID: ast.I8TypeID},
			},
			wantLinkError: []string{`could not resolve reference "Foo"`},
		},

		// set
		{
			desc:         "typedef string Foo; set<Foo>",
			give:         ast.SetType{ValueType: ast.TypeReference{Name: "Foo"}},
			scope:        scope("Foo", &TypedefSpec{Name: "Foo", Target: &StringSpec{}}),
			want:         &SetSpec{ValueSpec: &TypedefSpec{Name: "Foo", Target: &StringSpec{}}},
			wantTypeCode: wire.TSet,
		},
		{
			desc:          "set<Foo>",
			give:          ast.SetType{ValueType: ast.TypeReference{Name: "Foo"}},
			wantLinkError: []string{`could not resolve reference "Foo"`},
		},
	}

	for _, tt := range tests {
		spec, err := compileTypeReference(tt.give)
		if len(tt.wantCompileError) > 0 {
			if assert.Error(t, err, tt.desc) {
				for _, msg := range tt.wantCompileError {
					assert.Contains(t, err.Error(), msg, tt.desc)
				}
			}
			continue
		} else if !assert.NoError(t, err, tt.desc) {
			continue
		}

		scope := scopeOrDefault(tt.scope)
		spec, err = spec.Link(scope)
		if len(tt.wantLinkError) > 0 {
			if assert.Error(t, err, tt.desc) {
				for _, msg := range tt.wantLinkError {
					assert.Contains(t, err.Error(), msg, tt.desc)
				}
			}
			continue
		} else if !assert.NoError(t, err, tt.desc) {
			continue
		}

		want := mustLink(t, tt.want, defaultScope)
		assert.Equal(t, tt.wantTypeCode, spec.TypeCode(), tt.desc)
		assert.Equal(t, want, spec, tt.desc)
	}
}
