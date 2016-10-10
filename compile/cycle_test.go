// Copyright (c) 2016 Uber Technologies, Inc.
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

	"go.uber.org/thriftrw/ast"

	"github.com/stretchr/testify/assert"
)

// Helper to create typedefs with cyclic references inside the test table.
func withTypedef(f func(t *TypedefSpec)) TypeSpec {
	t := &TypedefSpec{File: "test.thrift"}
	f(t)
	return t
}

func TestFindTypeCycles(t *testing.T) {
	tests := []struct {
		desc string
		typ  TypeSpec
		msgs []string
	}{
		{
			desc: "self-referential typedef",
			typ: withTypedef(func(t *TypedefSpec) {
				t.Name = "foo"
				t.Target = t
			}),
			msgs: []string{
				"found a type reference cycle",
				"   foo",
				"-> foo",
			},
		},
		{
			desc: "mutually recursive references",
			typ: withTypedef(func(x *TypedefSpec) {
				x.Name = "foo"
				x.Target = &TypedefSpec{
					Name:   "bar",
					File:   "test.thrift",
					Target: x,
				}
			}),
			msgs: []string{
				"found a type reference cycle",
				"   foo",
				"-> bar",
				"-> foo",
			},
		},
		{
			desc: "recurse from map",
			typ: withTypedef(func(t *TypedefSpec) {
				t.Name = "foo"
				t.Target = &MapSpec{
					KeySpec:   &StringSpec{},
					ValueSpec: t,
				}
			}),
			msgs: []string{
				"found a type reference cycle",
				"   foo",
				"-> map<string, foo>",
				"-> foo",
			},
		},
		{
			desc: "recurse from list",
			typ: withTypedef(func(t *TypedefSpec) {
				t.Name = "foo"
				t.Target = &ListSpec{ValueSpec: t}
			}),
			msgs: []string{
				"found a type reference cycle",
				"   foo",
				"-> list<foo>",
				"-> foo",
			},
		},
		{
			desc: "recurse from set",
			typ: withTypedef(func(t *TypedefSpec) {
				t.Name = "foo"
				t.Target = &SetSpec{ValueSpec: t}
			}),
			msgs: []string{
				"found a type reference cycle",
				"   foo",
				"-> set<foo>",
				"-> foo",
			},
		},
		{
			desc: "recurse from struct",
			typ: withTypedef(func(t *TypedefSpec) {
				t.Name = "foo"
				t.Target = &StructSpec{
					Name: "bar",
					File: "test.thrift",
					Type: ast.StructType,
					Fields: FieldGroup{
						{
							ID:       1,
							Name:     "foo",
							Type:     t,
							Required: true,
						},
					},
				}
			}),
		},
	}

	for _, tt := range tests {
		typ := mustLink(t, tt.typ, defaultScope)

		err := findTypeCycles(typ)
		if len(tt.msgs) > 0 {
			if assert.Error(t, err, tt.desc) {
				for _, msg := range tt.msgs {
					assert.Contains(t, err.Error(), msg)
				}
			}
		} else {
			assert.NoError(t, err, tt.desc)
		}
	}
}
