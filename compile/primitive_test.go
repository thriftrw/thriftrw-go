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

	"go.uber.org/thriftrw/ast"

	"github.com/stretchr/testify/assert"
)

func TestCompileBaseType(t *testing.T) {
	tests := []struct {
		desc      string
		give      ast.BaseType
		want      TypeSpec
		wantError []string
	}{
		{
			desc: "bool",
			give: ast.BaseType{ID: ast.BoolTypeID},
			want: &BoolSpec{},
		},
		{
			desc: "i8",
			give: ast.BaseType{ID: ast.I8TypeID},
			want: &I8Spec{},
		},
		{
			desc: "i16",
			give: ast.BaseType{ID: ast.I16TypeID},
			want: &I16Spec{},
		},
		{
			desc: "i32",
			give: ast.BaseType{ID: ast.I32TypeID},
			want: &I32Spec{},
		},
		{
			desc: "i64",
			give: ast.BaseType{ID: ast.I64TypeID},
			want: &I64Spec{},
		},
		{
			desc: "double",
			give: ast.BaseType{ID: ast.DoubleTypeID},
			want: &DoubleSpec{},
		},
		{
			desc: "string",
			give: ast.BaseType{ID: ast.StringTypeID},
			want: &StringSpec{},
		},
		{
			desc: "binary",
			give: ast.BaseType{ID: ast.BinaryTypeID},
			want: &BinarySpec{},
		},

		// With annotations (success)
		{
			desc: `bool (x = "y")`,
			give: ast.BaseType{
				ID: ast.BoolTypeID,
				Annotations: []*ast.Annotation{
					{
						Name:  "x",
						Value: "y",
					},
				},
			},
			want: &BoolSpec{Annotations: Annotations{"x": "y"}},
		},
		{
			desc: `i8 (x = "\"")`,
			give: ast.BaseType{
				ID: ast.I8TypeID,
				Annotations: []*ast.Annotation{
					{
						Name:  "x",
						Value: `"`,
					},
				},
			},
			want: &I8Spec{Annotations: Annotations{"x": `"`}},
		},
		{
			desc: `i16 (a, b = "c")`,
			give: ast.BaseType{
				ID: ast.I16TypeID,
				Annotations: []*ast.Annotation{
					{Name: "a"},
					{Name: "b", Value: "c"},
				},
			},
			want: &I16Spec{Annotations: Annotations{
				"a": "",
				"b": "c",
			}},
		},
		{
			desc: `i32 ()`,
			give: ast.BaseType{
				ID:          ast.I32TypeID,
				Annotations: []*ast.Annotation{},
			},
			want: &I32Spec{},
		},
		{
			desc: `i64 (js.type = "Long")`,
			give: ast.BaseType{
				ID: ast.I64TypeID,
				Annotations: []*ast.Annotation{
					{Name: "js.type", Value: "Long"},
				},
			},
			want: &I64Spec{Annotations: Annotations{"js.type": "Long"}},
		},
		{
			desc: `double (java.type = "BigDecimal")`,
			give: ast.BaseType{
				ID: ast.DoubleTypeID,
				Annotations: []*ast.Annotation{
					{Name: "java.type", Value: "BigDecimal"},
				},
			},
			want: &DoubleSpec{Annotations: Annotations{"java.type": "BigDecimal"}},
		},
		{
			desc: `string (encoding = "utf8")`,
			give: ast.BaseType{
				ID: ast.StringTypeID,
				Annotations: []*ast.Annotation{
					{
						Name:  "encoding",
						Value: "utf8",
					},
				},
			},
			want: &StringSpec{Annotations: Annotations{"encoding": "utf8"}},
		},
		{
			desc: `binary (max_length = 42)`,
			give: ast.BaseType{
				ID: ast.BinaryTypeID,
				Annotations: []*ast.Annotation{
					{
						Name:  "max_length",
						Value: "42",
					},
				},
			},
			want: &BinarySpec{Annotations: Annotations{"max_length": "42"}},
		},

		// With annotations (failure)
		{
			desc: `bool (a, b, a)`,
			give: ast.BaseType{
				ID: ast.BoolTypeID,
				Annotations: []*ast.Annotation{
					{Name: "a"},
					{Name: "b"},
					{Name: "a"},
				},
			},
			wantError: []string{
				`annotation conflict: the name "a" has already been used on line 0`,
			},
		},
	}

	for _, tt := range tests {
		got, err := compileBaseType(tt.give)
		if len(tt.wantError) > 0 {
			if assert.Error(t, err, tt.desc) {
				for _, msg := range tt.wantError {
					assert.Contains(t, err.Error(), msg, tt.desc)
				}
			}
		} else {
			if assert.NoError(t, err, tt.desc) {
				assert.Equal(t, tt.want, got, tt.desc)
			}
		}
	}
}
