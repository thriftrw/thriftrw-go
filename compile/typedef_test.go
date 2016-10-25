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

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
	"go.uber.org/thriftrw/wire"
)

func parseTypedef(s string) *ast.Typedef {
	prog, err := idl.Parse([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("failure to parse: %v: %s", err, s))
	}

	if len(prog.Definitions) != 1 {
		panic("parseTypedef may be used to parse a single typedef only")
	}

	return prog.Definitions[0].(*ast.Typedef)
}

func TestCompileTypedef(t *testing.T) {
	tests := []struct {
		src   string
		scope Scope
		code  wire.Type
		spec  *TypedefSpec
	}{
		{
			`typedef i64 (js.type = "Long") timestamp (foo = "bar")`,
			nil,
			wire.TI64,
			&TypedefSpec{
				Name:        "timestamp",
				File:        "test.thrift",
				Target:      &I64Spec{Annotations: Annotations{"js.type": "Long"}},
				Annotations: Annotations{"foo": "bar"},
			},
		},
		{
			"typedef Bar Foo",
			scope("Bar", &TypedefSpec{
				Name:   "Bar",
				File:   "test.thrift",
				Target: &I32Spec{},
			}),
			wire.TI32,
			&TypedefSpec{
				Name: "Foo",
				File: "test.thrift",
				Target: &TypedefSpec{
					Name:   "Bar",
					File:   "test.thrift",
					Target: &I32Spec{},
				},
			},
		},
	}

	for _, tt := range tests {
		expected := mustLink(t, tt.spec, defaultScope)

		src := parseTypedef(tt.src)
		typedefSpec, err := compileTypedef("test.thrift", src)
		if !assert.NoError(t, err, tt.src) {
			continue
		}

		scope := scopeOrDefault(tt.scope)
		spec, err := typedefSpec.Link(scope)
		if assert.NoError(t, err) {
			assert.Equal(t, tt.code, spec.TypeCode())
			assert.Equal(t, expected, spec)
		}
	}
}

func TestCompileTypedefFailure(t *testing.T) {
	tests := []struct {
		desc     string
		src      string
		scope    Scope
		messages []string
	}{
		{
			"unknown type",
			"typedef foo bar",
			nil,
			[]string{
				`could not resolve reference "foo"`,
			},
		},
		{
			"conflicting annotations",
			`typedef i32 bar (a = "b", a, b)`,
			nil,
			[]string{
				`cannot compile "bar" on line 1:`,
				`annotation conflict: the name "a" has already been used on line 1`,
			},
		},
	}

	for _, tt := range tests {
		src := parseTypedef(tt.src)
		scope := scopeOrDefault(tt.scope)

		spec, err := compileTypedef("test.thrift", src)
		if err == nil {
			_, err = spec.Link(scope)
		}

		assert.Error(t, err, tt.desc)
		for _, msg := range tt.messages {
			assert.Contains(t, err.Error(), msg, tt.desc)
		}
	}
}
