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

func parseService(s string) *ast.Service {
	prog, err := idl.Parse([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("failure to parse: %v: %s", err, s))
	}

	if len(prog.Definitions) != 1 {
		panic("parseService may be used to parse single services only")
	}

	return prog.Definitions[0].(*ast.Service)
}

func TestCompileService(t *testing.T) {
	keyDoesNotExistSpec := &StructSpec{
		Name:   "KeyDoesNotExist",
		File:   "test.thrift",
		Type:   ast.ExceptionType,
		Fields: make(FieldGroup, 0),
	}

	internalErrorSpec := &StructSpec{
		Name:   "InternalServiceError",
		File:   "test.thrift",
		Type:   ast.ExceptionType,
		Fields: make(FieldGroup, 0),
	}

	keyValueSpec := &ServiceSpec{
		Name: "KeyValue",
		File: "test.thrift",
		Functions: map[string]*FunctionSpec{
			"setValue": {
				Name: "setValue",
				ArgsSpec: ArgsSpec{
					{
						ID:   1,
						Name: "key",
						Type: &StringSpec{},
					},
					{
						ID:   2,
						Name: "value",
						Type: &BinarySpec{},
					},
				},
				ResultSpec: &ResultSpec{},
			},
			"getValue": {
				Name: "getValue",
				ArgsSpec: ArgsSpec{
					{
						ID:   1,
						Name: "key",
						Type: &StringSpec{},
					},
				},
				ResultSpec: &ResultSpec{
					ReturnType: &BinarySpec{},
					Exceptions: FieldGroup{
						{
							ID:   1,
							Name: "doesNotExist",
							Type: keyDoesNotExistSpec,
						},
						{
							ID:   2,
							Name: "internalError",
							Type: internalErrorSpec,
						},
					},
				},
			},
		},
	}

	annotatedSpec := &ServiceSpec{
		Name: "AnnotatedService",
		File: "test.thrift",
		Functions: map[string]*FunctionSpec{
			"setValue": {
				Name: "setValue",
				ArgsSpec: ArgsSpec{
					{
						ID:   1,
						Name: "key",
						Type: &StringSpec{},
					},
					{
						ID:   2,
						Name: "value",
						Type: &BinarySpec{},
					},
				},
				Annotations: Annotations{
					"test": "ok",
				},
				ResultSpec: &ResultSpec{},
			},
		},
		Annotations: Annotations{
			"test": "abcde",
		},
	}

	tests := []struct {
		desc  string
		src   string
		scope Scope
		spec  *ServiceSpec
	}{
		{
			"empty service",
			"service Foo {}",
			nil,
			&ServiceSpec{
				Name:      "Foo",
				File:      "test.thrift",
				Functions: make(map[string]*FunctionSpec),
			},
		},
		{
			"simple service",
			`
				service KeyValue {
					void setValue(1: string key, 2: binary value)
					binary getValue(1: string key)
						throws (
							1: KeyDoesNotExist doesNotExist,
							2: InternalServiceError internalError
						)
				}
			`,
			scope(
				"KeyDoesNotExist", keyDoesNotExistSpec,
				"InternalServiceError", internalErrorSpec,
			),
			keyValueSpec,
		},
		{
			"service annotations",
			`
				service AnnotatedService {
					void setValue(1: string key, 2: binary value) (test = "ok")
				} (test = "abcde")
			`,
			scope(),
			annotatedSpec,
		},
		{
			"service inheritance",
			`
				service BulkKeyValue extends KeyValue {
					void setValues(1: map<string, binary> items)
				}
			`,
			scope("KeyValue", keyValueSpec),
			&ServiceSpec{
				Name:   "BulkKeyValue",
				File:   "test.thrift",
				Parent: keyValueSpec,
				Functions: map[string]*FunctionSpec{
					"setValues": {
						Name: "setValues",
						ArgsSpec: ArgsSpec{
							{
								ID:   1,
								Name: "items",
								Type: &MapSpec{
									KeySpec:   &StringSpec{},
									ValueSpec: &BinarySpec{},
								},
							},
						},
						ResultSpec: &ResultSpec{},
					},
				},
			},
		},
		{
			"included service inheritance",
			"service AnotherKeyValue extends shared.KeyValue {}",
			scope("shared", scope("KeyValue", keyValueSpec)),
			&ServiceSpec{
				Name:      "AnotherKeyValue",
				File:      "test.thrift",
				Parent:    keyValueSpec,
				Functions: make(map[string]*FunctionSpec),
			},
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

func TestCompileServiceFailure(t *testing.T) {
	tests := []struct {
		desc     string
		src      string
		messages []string
	}{
		{
			"duplicate function name",
			`
				service Foo {
					void foo()
					void bar()
					i32 foo()
				}
			`,
			[]string{
				`the name "foo" has already been used on line 3`,
			},
		},
		{
			"duplicate in arg list",
			`
				service Foo {
					void bar(
						1: string foo,
						2: binary bar,
						3: i32 foo
					)
				}
			`,
			[]string{
				`cannot compile "bar"`,
				`the name "foo" has already been used`,
			},
		},
		{
			"duplicate in exception list",
			`
				service Foo {
					i32 bar(1: string foo) throws (
						1: KeyDoesNotExist error,
						2: InternalServiceError error,
					)
				}
			`,
			[]string{
				`cannot compile "bar"`,
				`the name "error" has already been used`,
			},
		},
		{
			"exceptions cannot have default values",
			`
				service Foo {
					void bar() throws (
						1: KeyDoesNotExist doesNotExist,
						2: InternalServiceError internalError = {
							'message': 'An internal error occurred'
						}
					)
				}
			`,
			[]string{`field "internalError"`, "cannot have a default value"},
		},
		{
			"oneway cannot return",
			"service Foo { oneway i32 bar() }",
			[]string{`function "bar" cannot return values`},
		},
		{
			"oneway cannot raise",
			`
				service Foo {
					oneway void bar()
						throws (1: KeyDoesNotExistError err)
				}
			`,
			[]string{`function "bar" cannot`, "raise exceptions"},
		},
		{
			"duplicate annotation name",
			`
				service AnnotatedService {

				} (test = "mytest", test = "t")
			`,
			[]string{
				`the name "test" has already been used`,
			},
		},
		{
			"duplicate annotation name on function",
			`
				service AnnotatedService {
					i32 bar() (functest = "test", functest = "t")
				} (test = "mytest")
			`,
			[]string{
				`the name "functest" has already been used`,
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

func TestLinkServiceFailure(t *testing.T) {
	tests := []struct {
		desc     string
		src      string
		scope    Scope
		messages []string
	}{
		{
			"inherit unknown service",
			"service Foo extends Bar {}",
			nil,
			[]string{
				`cannot compile "Foo"`,
				`could not resolve reference "Bar"`,
			},
		},
		{
			"inherit service from unknown module",
			"service SomeService extends foo.Bar {}",
			scope("bar"),
			[]string{
				`could not resolve reference "foo.Bar"`, `in "bar"`,
				`unknown module "foo"`,
			},
		},
		{
			"unknown service inherited from included module",
			"service SomeService extends foo.Bar {}",
			scope("bar", "foo", scope("foo")),
			[]string{
				`could not resolve reference "foo.Bar"`, `in "bar"`,
				`could not resolve reference "Bar" in "foo"`,
			},
		},
		{
			"unknown reference from inherited service",
			"service Foo extends Bar {}",
			scope(
				"Bar", &ServiceSpec{
					Name:      "Bar",
					Functions: make(map[string]*FunctionSpec),
					parentSrc: &ast.ServiceReference{Name: "Baz"},
				},
			),
			[]string{
				`cannot compile "Foo"`,
				`cannot compile "Bar"`,
				`could not resolve reference "Baz"`,
			},
		},
		{
			"can throw exceptions only",
			`
				service Foo {
					void foo()
						throws (
							1: SomeException a
							2: NotAnException b
						)
				}
			`,
			scope(
				"SomeException", &StructSpec{
					Name:   "SomeException",
					Type:   ast.ExceptionType,
					Fields: make(FieldGroup, 0),
				},
				"NotAnException", &StructSpec{
					Name:   "NotAnException",
					Type:   ast.StructType,
					Fields: make(FieldGroup, 0),
				},
			),
			[]string{
				`field "b" with type "NotAnException" is not an exception`,
			},
		},
		{
			"unknown return type",
			"service Foo { Baz bar() }",
			nil,
			[]string{
				`cannot compile "Foo.bar"`,
				`cannot compile "bar"`,
				`could not resolve reference "Baz"`,
			},
		},
		{
			"unknown exception type",
			`
				service Foo {
					i32 bar()
						throws (1: UnknownExceptionError error)
				}
			`,
			nil,
			[]string{
				`cannot compile "Foo.bar"`,
				`cannot compile "bar"`,
				`could not resolve reference "UnknownExceptionError"`,
			},
		},
	}

	for _, tt := range tests {
		src := parseService(tt.src)
		scope := scopeOrDefault(tt.scope)

		spec, err := compileService("test.thrift", src)
		if assert.NoError(t, err, tt.desc) {
			if err := spec.Link(scope); assert.Error(t, err) {
				for _, msg := range tt.messages {
					assert.Contains(t, err.Error(), msg, tt.desc)
				}
			}
		}
	}
}
