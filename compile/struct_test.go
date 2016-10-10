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

func parseStruct(s string) *ast.Struct {
	prog, err := idl.Parse([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("failure to parse: %v: %s", err, s))
	}

	if len(prog.Definitions) != 1 {
		panic("parseStruct may be used to parse single structs only")
	}

	return prog.Definitions[0].(*ast.Struct)
}

func TestCompileStructSuccess(t *testing.T) {
	tests := []struct {
		src          string
		scope        Scope
		requiredness fieldRequiredness
		spec         *StructSpec
	}{
		{
			"struct Health { 1: optional bool healthy = true }",
			nil,
			explicitRequiredness,
			&StructSpec{
				Name: "Health",
				File: "test.thrift",
				Type: ast.StructType,
				Fields: FieldGroup{
					{
						ID:       1,
						Name:     "healthy",
						Type:     &BoolSpec{},
						Required: false,
						Default:  ConstantBool(true),
					},
				},
			},
		},
		{
			"struct Health { 1: bool healthy = true }",
			nil,
			defaultToOptional,
			&StructSpec{
				Name: "Health",
				File: "test.thrift",
				Type: ast.StructType,
				Fields: FieldGroup{
					{
						ID:       1,
						Name:     "healthy",
						Type:     &BoolSpec{},
						Required: false,
						Default:  ConstantBool(true),
					},
				},
			},
		},
		{
			`exception KeyNotFoundError {
				1: required string message
				2: optional Key key
			}`,
			scope("Key", &TypedefSpec{Name: "Key", Target: &StringSpec{}}),
			explicitRequiredness,
			&StructSpec{
				Name: "KeyNotFoundError",
				File: "test.thrift",
				Type: ast.ExceptionType,
				Fields: FieldGroup{
					{
						ID:       1,
						Name:     "message",
						Type:     &StringSpec{},
						Required: true,
					},
					{
						ID:       2,
						Name:     "key",
						Type:     &TypedefSpec{Name: "Key", Target: &StringSpec{}},
						Required: false,
					},
				},
			},
		},
		{
			`union Body {
				1234: optional string plainText
				5678: binary richText
			}`,
			nil,
			explicitRequiredness,
			&StructSpec{
				Name: "Body",
				File: "test.thrift",
				Type: ast.UnionType,
				Fields: FieldGroup{
					{
						ID:       1234,
						Name:     "plainText",
						Type:     &StringSpec{},
						Required: false,
					},
					{
						ID:       5678,
						Name:     "richText",
						Type:     &BinarySpec{},
						Required: false,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		expected := mustLink(t, tt.spec, defaultScope)

		src := parseStruct(tt.src)
		structSpec, err := compileStruct("test.thrift", src, tt.requiredness)
		scope := scopeOrDefault(tt.scope)
		if assert.NoError(t, err) {
			spec, err := structSpec.Link(scope)
			assert.NoError(t, err)
			assert.Equal(t, expected, spec)
			assert.Equal(t, wire.TStruct, spec.TypeCode())
		}
	}
}

func TestCompileStructFailure(t *testing.T) {
	tests := []struct {
		desc     string
		src      string
		messages []string
	}{
		{
			"optional/required is required for structs and exceptions",
			"struct Foo { 1: string bar }",
			[]string{
				`field "bar"`,
				"not marked required or optional",
			},
		},
		{
			"optional/required is required for structs and exceptions",
			"exception Foo { 1: string bar }",
			[]string{
				`field "bar"`,
				"not marked required or optional",
			},
		},
		{
			"unions cannot have required fields",
			"union Foo { 1: required string bar }",
			[]string{
				`field "bar"`,
				"marked as required but it cannot be required",
			},
		},
		{
			"unions cannot have default values",
			`
				union Foo {
					1: string a
					2: binary b
					3: i32 c = 42
				}
			`,
			[]string{`field "c"`, "cannot have a default value"},
		},
		{
			"field name conflict",
			`struct Foo {
				1: required string bar
				2: optional string baz
				3: optional i32 bar
			}`,
			[]string{`the name "bar" has already been used`},
		},
		{
			"field ID conflict",
			`struct Foo {
				1: optional string foo
				1: optional string bar
			}`,
			[]string{`field "foo" has already used ID 1`},
		},
	}

	for _, tt := range tests {
		src := parseStruct(tt.src)
		_, err := compileStruct("test.thrift", src, explicitRequiredness)

		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}

func TestLinkStructFailure(t *testing.T) {
	tests := []struct {
		desc     string
		src      string
		scope    Scope
		messages []string
	}{
		{
			"unknown field type",
			"struct Foo { 1: optional Bar bar }",
			nil,
			[]string{`could not resolve reference "Bar"`},
		},
		{
			"unknown constant as default value",
			`
				struct Foo {
					1: optional string foo = DEFAULT_FOO
				}
			`,
			nil,
			[]string{
				`could not resolve reference "DEFAULT_FOO"`,
			},
		},
	}

	for _, tt := range tests {
		src := parseStruct(tt.src)
		scope := scopeOrDefault(tt.scope)

		spec, err := compileStruct("test.thrift", src, explicitRequiredness)
		if assert.NoError(t, err, tt.desc) {
			_, err := spec.Link(scope)
			if assert.Error(t, err, tt.desc) {
				for _, msg := range tt.messages {
					assert.Contains(t, err.Error(), msg, tt.desc)
				}
			}
		}
	}
}
