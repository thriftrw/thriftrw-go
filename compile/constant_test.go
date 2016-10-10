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
	"go.uber.org/thriftrw/idl"
)

func parseConstant(s string) *ast.Constant {
	prog, err := idl.Parse([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("failure to parse: %v: %s", err, s))
	}

	if len(prog.Definitions) != 1 {
		panic("parseConstant may be used to parse single constants only")
	}

	return prog.Definitions[0].(*ast.Constant)
}

func TestCompileConstant(t *testing.T) {
	y := &Constant{
		Name:  "y",
		Type:  &StringSpec{},
		Value: ConstantString("bar"),
	}

	tests := []struct {
		src      string
		scope    Scope
		constant *Constant
	}{
		{
			"const i32 version = 1",
			nil,
			&Constant{
				Name:  "version",
				File:  "test.thrift",
				Type:  &I32Spec{},
				Value: ConstantInt(1),
			},
		},
		{
			`const string foo = "hello world"`,
			nil,
			&Constant{
				Name:  "foo",
				File:  "test.thrift",
				Type:  &StringSpec{},
				Value: ConstantString("hello world"),
			},
		},
		{
			`const list<string> foo = ["hello", "world"]`,
			nil,
			&Constant{
				Name: "foo",
				File: "test.thrift",
				Type: &ListSpec{ValueSpec: &StringSpec{}},
				Value: ConstantList{
					ConstantString("hello"),
					ConstantString("world"),
				},
			},
		},
		{
			`const list<string> foo = ["x", y]`,
			scope("y", y),
			&Constant{
				Name: "foo",
				File: "test.thrift",
				Type: &ListSpec{ValueSpec: &StringSpec{}},
				Value: ConstantList{
					ConstantString("x"),
					ConstReference{Target: y},
				},
			},
		},
	}

	for _, tt := range tests {
		scope := scopeOrDefault(tt.scope)
		src := parseConstant(tt.src)
		require.NoError(
			t,
			tt.constant.Link(scope),
			"invalid test: expected spec failed to link",
		)

		constant, err := compileConstant("test.thrift", src)
		if assert.NoError(t, err) && assert.NoError(t, constant.Link(scope)) {
			assert.Equal(t, tt.constant, constant)
		}
	}
}

func TestCompileConstantFailure(t *testing.T) {
	y := &Constant{
		Name:  "y",
		Type:  &StringSpec{},
		Value: ConstantString("bar"),
	}

	tests := []struct {
		src      string
		scope    Scope
		messages []string
	}{
		{
			`const list<string> foo = [x]`,
			scope("y", y),
			[]string{
				`cannot compile "foo"`,
				`could not resolve reference "x"`,
			},
		},
		{
			`const Foo foo = {"bar": "baz"}`,
			nil,
			[]string{
				`cannot compile "foo"`,
				`could not resolve reference "Foo"`,
			},
		},
		{
			`const map<string, string> something = {foo: "bar"}`,
			nil,
			[]string{
				`cannot compile "something"`,
				`could not resolve reference "foo"`,
			},
		},
		{
			`const map<string, string> something = {"foo": bar}`,
			nil,
			[]string{
				`cannot compile "something"`,
				`could not resolve reference "bar"`,
			},
		},
	}

	for _, tt := range tests {
		src := parseConstant(tt.src)
		scope := scopeOrDefault(tt.scope)
		constant, err := compileConstant("test.thrift", src)
		if assert.NoError(t, err) {
			err := constant.Link(scope)
			if assert.Error(t, err) {
				for _, msg := range tt.messages {
					assert.Contains(t, err.Error(), msg)
				}
			}
		}
	}
}
