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

func TestCompileAnnotation(t *testing.T) {
	tests := []struct {
		desc string
		src  []*ast.Annotation
		spec Annotations
	}{
		{
			"simple annotation",
			[]*ast.Annotation{
				{
					Name:  "test",
					Value: "ok",
					Line:  1,
				},
			},
			Annotations{
				"test": "ok",
			},
		},
	}

	for _, tt := range tests {
		spec, err := compileAnnotations(tt.src)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, tt.spec, spec, tt.desc)
		}
	}
}

func TestCompileAnnotationFailure(t *testing.T) {
	tests := []struct {
		desc     string
		src      []*ast.Annotation
		messages []string
	}{
		{
			"duplicate annotation",
			[]*ast.Annotation{
				{
					Name:  "test",
					Value: "ok",
					Line:  1,
				},
				{
					Name:  "test",
					Value: "abcd",
					Line:  1,
				},
			},
			[]string{
				`the name "test" has already been used`,
			},
		},
	}

	for _, tt := range tests {
		_, err := compileAnnotations(tt.src)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}
