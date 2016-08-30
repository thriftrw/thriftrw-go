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
	"github.com/stretchr/testify/require"
)

func TestCompileAnnotation(t *testing.T) {

	annotatedSpec := &ServiceSpec{
		Name:      "AnnotatedService",
		File:      "test.thrift",
		Functions: make(map[string]*FunctionSpec),
		Annotations: []*AnnotationSpec{
			{
				Name:  "test",
				Value: "test",
			},
		},
	}

	tests := []struct {
		desc  string
		src   string
		scope Scope
		spec  *ServiceSpec
	}{
		{
			"service annotations",
			`
				service AnnotatedService {

				} (test = "test")
			`,
			scope(),
			annotatedSpec,
		},
	}

	for _, tt := range tests {
		require.NoError(
			t, tt.spec.Link(defaultScope),
			"invalid test: service must with an empty scope",
		)
		scope := scopeOrDefault(tt.scope)

		src := parseService(tt.src)
		spec, err := compileService("test.thrift", src)
		if assert.NoError(t, err, tt.desc) {
			if assert.NoError(t, spec.Link(scope), tt.desc) {
				assert.Equal(t, tt.spec, spec, tt.desc)
			}
		}
	}
}

func TestCompileAnnotationFailure(t *testing.T) {

	tests := []struct {
		desc     string
		src      string
		messages []string
	}{
		{
			"service annotations",
			`
				service AnnotatedService {

				} (test = "test", test = "t")
			`,
			[]string{
				`the name "test" has already been used`,
			},
		},
	}

	for _, tt := range tests {
		src := parseService(tt.src)
		_, err := compileService("test.thrift", src)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}
