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

package idl

import (
	"strings"
	"testing"

	. "go.uber.org/thriftrw/ast"

	"github.com/kr/pretty"
	"github.com/stretchr/testify/assert"
)

type parseCase struct {
	document string
	program  *Program
}

func assertParseCases(t *testing.T, tests []parseCase) {
	for _, tt := range tests {
		program, err := Parse([]byte(tt.document))
		if assert.NoError(t, err, "Parsing failed:\n%s", tt.document) {
			succ := assert.Equal(
				t, tt.program, program,
				"Got unexpected program when parsing:\n%s", tt.document,
			)
			if !succ {
				lines := pretty.Diff(tt.program, program)
				t.Log("\n\t" + strings.Join(lines, "\n\t"))
			}
		}
	}
}

func TestParseEmpty(t *testing.T) {
	program, err := Parse([]byte{})
	if assert.NoError(t, err, "%v", err) {
		assert.Equal(t, &Program{}, program)
	}
}

func TestParseComments(t *testing.T) {
	s := "// foo\n//\n//bar"
	_, err := Parse([]byte(s))
	assert.NoError(t, err, "Failed to parse:\n%s", s)
}

func TestParseErrors(t *testing.T) {
	tests := []struct {
		give       string
		wantErrors []string
	}{
		{
			give:       "namespace foo \x00",
			wantErrors: []string{"line 1: unknown token at index 14"},
		},
		{
			give:       `const string 42 = "foo"`,
			wantErrors: []string{"line 1:", "unexpected INTCONSTANT, expecting IDENTIFIER"},
		},
		{
			give:       `typedef foo bar baz`,
			wantErrors: []string{"line 1:", "unexpected IDENTIFIER"},
		},
		{
			give:       `typedef foo`,
			wantErrors: []string{"line 1:", "unexpected $end"},
		},
		{
			give:       `enum Foo {`,
			wantErrors: []string{"line 1:", "unexpected $end"},
		},
		{
			give:       `enum { }`,
			wantErrors: []string{"line 1:", "unexpected '{'"},
		},
		{
			give: `
				enum Foo {}
				include "bar.thrift"
			`,
			wantErrors: []string{"line 3:", "unexpected INCLUDE"},
		},
		{
			give:       `service Foo extends {}`,
			wantErrors: []string{"line 1:", "unexpected '{'"},
		},
		{
			give:       `service Foo Bar {}`,
			wantErrors: []string{"line 1:", "unexpected IDENTIFIER"},
		},
		{
			give:       `service Foo { void foo() () (foo = "bar") }`,
			wantErrors: []string{"line 1:", "unexpected '('"},
		},
		{
			give:       `service Foo { void foo() throws }`,
			wantErrors: []string{"line 1:", "unexpected '}'"},
		},
		{
			give:       `typedef string (foo =) UUID`,
			wantErrors: []string{"line 1:", "unexpected ')'"},
		},
		{
			give:       `union Operation { 1: Insert insert; 2: Delete delete }`,
			wantErrors: []string{"line 1:", `"delete" is a reserved keyword`},
		},
	}

	for _, tt := range tests {
		_, err := Parse([]byte(tt.give))
		if assert.Error(t, err, "expected error while parsing:\n%s", tt.give) {
			for _, msg := range tt.wantErrors {
				assert.Contains(t, err.Error(), msg, "error for %q must contain %q", tt.give, err.Error(), msg)
			}
		}
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []parseCase{
		{
			`
				include "foo.thrift"
				include t "bar.thrift"
			`,
			&Program{Headers: []Header{
				&Include{Path: "foo.thrift", Line: 2},
				&Include{Path: "bar.thrift", Name: "t", Line: 3},
			}},
		},
		{
			`
				namespace py bar
				namespace * foo
			`,
			&Program{Headers: []Header{
				&Namespace{Scope: "py", Name: "bar", Line: 2},
				&Namespace{Scope: "*", Name: "foo", Line: 3},
			}},
		},
		{
			`
				// defines shared types
				include "shared.thrift"
				namespace go foo_service  # go namespace should be foo_service

				/*
					common error types
				*/
				include "errors.thrift"

				/* python code goes to service.foo */
				namespace py services.foo
			`,
			&Program{
				Headers: []Header{
					&Include{Path: "shared.thrift", Line: 3},
					&Namespace{Scope: "go", Name: "foo_service", Line: 4},
					&Include{Path: "errors.thrift", Line: 9},
					&Namespace{Scope: "py", Name: "services.foo", Line: 12},
				},
			},
		},
	}
	assertParseCases(t, tests)
}

func TestParseConstants(t *testing.T) {
	tests := []parseCase{
		{
			`
				const i32 foo = 42
				const i64 bar = shared.baz;

				const string baz = "hello world";

				const double qux = 3.141592

				// def is reserved but def_ is not
				const double def_ = 1.23
			`,
			&Program{Definitions: []Definition{
				&Constant{
					Name:  "foo",
					Type:  BaseType{ID: I32TypeID, Line: 2},
					Value: ConstantInteger(42),
					Line:  2,
				},
				&Constant{
					Name: "bar",
					Type: BaseType{ID: I64TypeID, Line: 3},
					Value: ConstantReference{
						Name: "shared.baz",
						Line: 3,
					},
					Line: 3,
				},
				&Constant{
					Name:  "baz",
					Type:  BaseType{ID: StringTypeID, Line: 5},
					Value: ConstantString("hello world"),
					Line:  5,
				},
				&Constant{
					Name:  "qux",
					Type:  BaseType{ID: DoubleTypeID, Line: 7},
					Value: ConstantDouble(3.141592),
					Line:  7,
				},
				&Constant{
					Name:  "def_",
					Type:  BaseType{ID: DoubleTypeID, Line: 10},
					Value: ConstantDouble(1.23),
					Line:  10,
				},
			}},
		},
		{
			`const bool (foo = "a\nb") baz = true
			 const bool include_something = false`,
			&Program{Definitions: []Definition{
				&Constant{
					Name: "baz",
					Type: BaseType{
						ID:   BoolTypeID,
						Line: 1,
						Annotations: []*Annotation{
							{Name: "foo", Value: "a\nb", Line: 1},
						},
					},
					Value: ConstantBoolean(true),
					Line:  1,
				},
				&Constant{
					Name:  "include_something",
					Type:  BaseType{ID: BoolTypeID, Line: 2},
					Value: ConstantBoolean(false),
					Line:  2,
				},
			}},
		},
		{
			`
				const map<string (foo), i32> (baz = "qux") stuff = {
					"a": 1,
					"b": 2,
				}
				const list<list<i32>> list_of_lists = [
					[1, 2, 3]  # optional separator
					[4, 5, 6]
				];
				const Item const_struct = {
					"key": "foo",
					"value": 42,
				};
			`,
			&Program{Definitions: []Definition{
				&Constant{
					Name: "stuff",
					Type: MapType{
						KeyType: BaseType{
							ID:   StringTypeID,
							Line: 2,
							Annotations: []*Annotation{
								{Name: "foo", Value: "", Line: 2},
							},
						},
						ValueType: BaseType{ID: I32TypeID, Line: 2},
						Line:      2,
						Annotations: []*Annotation{
							{Name: "baz", Value: "qux", Line: 2},
						},
					},
					Value: ConstantMap{
						Items: []ConstantMapItem{
							{
								Key:   ConstantString("a"),
								Value: ConstantInteger(1),
								Line:  3,
							},
							{
								Key:   ConstantString("b"),
								Value: ConstantInteger(2),
								Line:  4,
							},
						},
						Line: 2,
					},
					Line: 2,
				},
				&Constant{
					Name: "list_of_lists",
					Type: ListType{
						ValueType: ListType{
							ValueType: BaseType{ID: I32TypeID, Line: 6},
							Line:      6,
						},
						Line: 6,
					},
					Value: ConstantList{
						Items: []ConstantValue{
							ConstantList{
								Items: []ConstantValue{
									ConstantInteger(1),
									ConstantInteger(2),
									ConstantInteger(3),
								},
								Line: 7,
							},
							ConstantList{
								Items: []ConstantValue{
									ConstantInteger(4),
									ConstantInteger(5),
									ConstantInteger(6),
								},
								Line: 8,
							},
						},
						Line: 6,
					},
					Line: 6,
				},
				&Constant{
					Name: "const_struct",
					Type: TypeReference{Name: "Item", Line: 10},
					Value: ConstantMap{
						Items: []ConstantMapItem{
							{
								Key:   ConstantString("key"),
								Value: ConstantString("foo"),
								Line:  11,
							},
							{
								Key:   ConstantString("value"),
								Value: ConstantInteger(42),
								Line:  12,
							},
						},
						Line: 10,
					},
					Line: 10,
				},
			}},
		},
		{
			`
				const string foo = 'a "b" c'
				const string bar = "a 'b' c"
			`,
			&Program{Definitions: []Definition{
				&Constant{
					Name:  "foo",
					Type:  BaseType{ID: StringTypeID, Line: 2},
					Value: ConstantString(`a "b" c`),
					Line:  2,
				},
				&Constant{
					Name:  "bar",
					Type:  BaseType{ID: StringTypeID, Line: 3},
					Value: ConstantString(`a 'b' c`),
					Line:  3,
				},
			}},
		},
	}
	assertParseCases(t, tests)
}

func TestParseTypedef(t *testing.T) {
	tests := []parseCase{
		{
			`
				typedef string UUID (length = "32");

				typedef i64 (js.type = "Date") Date

				typedef i8 foo
				typedef byte bar
			`,
			&Program{Definitions: []Definition{
				&Typedef{
					Name: "UUID",
					Type: BaseType{ID: StringTypeID, Line: 2},
					Annotations: []*Annotation{
						{
							Name:  "length",
							Value: "32",
							Line:  2,
						},
					},
					Line: 2,
				},
				&Typedef{
					Name: "Date",
					Type: BaseType{
						ID:   I64TypeID,
						Line: 4,
						Annotations: []*Annotation{
							{
								Name:  "js.type",
								Value: "Date",
								Line:  4,
							},
						},
					},
					Line: 4,
				},
				&Typedef{
					Name: "foo",
					Type: BaseType{ID: I8TypeID, Line: 6},
					Line: 6,
				},
				&Typedef{
					Name: "bar",
					Type: BaseType{ID: I8TypeID, Line: 7},
					Line: 7,
				},
			}},
		},
	}

	assertParseCases(t, tests)
}

func TestParseEnum(t *testing.T) {
	aValue := 42

	tests := []parseCase{
		{
			`
				enum EmptyEnum
				{
				}
			`,
			&Program{Definitions: []Definition{&Enum{Name: "EmptyEnum", Line: 2}}},
		},
		{
			`
				enum SillyEnum {
					foo (x), bar /*
					*/ baz = 42
					qux;
					quux
				} (_ = "__", foo = "bar")
			`,
			&Program{Definitions: []Definition{
				&Enum{
					Name: "SillyEnum",
					Items: []*EnumItem{
						{
							Name: "foo",
							Annotations: []*Annotation{
								{
									Name:  "x",
									Value: "",
									Line:  3,
								},
							},
							Line: 3,
						},
						{Name: "bar", Line: 3},
						{Name: "baz", Value: &aValue, Line: 4},
						{Name: "qux", Line: 5},
						{Name: "quux", Line: 6},
					},
					Annotations: []*Annotation{
						{Name: "_", Value: "__", Line: 7},
						{Name: "foo", Value: "bar", Line: 7},
					},
					Line: 2,
				},
			}},
		},
	}

	assertParseCases(t, tests)
}

func TestParseStruct(t *testing.T) {
	tests := []parseCase{
		{
			`
				struct EmptyStruct {}
				union EmptyUnion {}
				exception EmptyExc {}
			`,
			&Program{Definitions: []Definition{
				&Struct{Name: "EmptyStruct", Type: StructType, Line: 2},
				&Struct{Name: "EmptyUnion", Type: UnionType, Line: 3},
				&Struct{Name: "EmptyExc", Type: ExceptionType, Line: 4},
			}},
		},
		{
			`
				struct i128 {
					1: required i64 high
					2: required i64 low
				} (serializer = "Int128Serializer")

				union Contents {
					1: string (format = "markdown") plainText
					2: binary pdf (name = "pdfFile")
				}

				exception GreatSadness {
					1: optional string message
				}
			`,
			&Program{Definitions: []Definition{
				&Struct{
					Name: "i128",
					Type: StructType,
					Fields: []*Field{
						{
							ID:           1,
							Name:         "high",
							Type:         BaseType{ID: I64TypeID, Line: 3},
							Requiredness: Required,
							Line:         3,
						},
						{
							ID:           2,
							Name:         "low",
							Type:         BaseType{ID: I64TypeID, Line: 4},
							Requiredness: Required,
							Line:         4,
						},
					},
					Annotations: []*Annotation{
						{
							Name:  "serializer",
							Value: "Int128Serializer",
							Line:  5,
						},
					},
					Line: 2,
				},
				&Struct{
					Name: "Contents",
					Type: UnionType,
					Fields: []*Field{
						{
							ID:           1,
							Name:         "plainText",
							Requiredness: Unspecified,
							Type: BaseType{
								ID:   StringTypeID,
								Line: 8,
								Annotations: []*Annotation{
									{
										Name:  "format",
										Value: "markdown",
										Line:  8,
									},
								},
							},
							Line: 8,
						},
						{
							ID:   2,
							Name: "pdf",
							Type: BaseType{ID: BinaryTypeID, Line: 9},
							// Requiredness intentionally skipped because
							// zero-value for it is Unspecified.
							Annotations: []*Annotation{
								{
									Name:  "name",
									Value: "pdfFile",
									Line:  9,
								},
							},
							Line: 9,
						},
					},
					Line: 7,
				},
				&Struct{
					Name: "GreatSadness",
					Type: ExceptionType,
					Fields: []*Field{
						{
							ID:           1,
							Name:         "message",
							Type:         BaseType{ID: StringTypeID, Line: 13},
							Requiredness: Optional,
							Line:         13,
						},
					},
					Line: 12,
				},
			}},
		},
	}

	assertParseCases(t, tests)
}

func TestParseServices(t *testing.T) {
	tests := []parseCase{
		{
			`
				service EmptyService {}
				service AnotherEmptyService extends EmptyService {}
			`,
			&Program{Definitions: []Definition{
				&Service{Name: "EmptyService", Line: 2},
				&Service{
					Name: "AnotherEmptyService",
					Parent: &ServiceReference{
						Name: "EmptyService",
						Line: 3,
					},
					Line: 3,
				},
			}},
		},
		{
			`
				service KeyValue {
					oneway void
						empty()
							throws ()

					i32 something(
					) throws (1: GreatSadness sadness);

					void somethingElse(
						1: A a;
						2: B b;
					) (py.name = "something_else"),
				} (ttl.milliseconds = "200")
			`,
			&Program{Definitions: []Definition{
				&Service{
					Name: "KeyValue",
					Functions: []*Function{
						{
							Name:   "empty",
							OneWay: true,
							Line:   4,
						},
						{
							Name:       "something",
							ReturnType: BaseType{ID: I32TypeID, Line: 7},
							Exceptions: []*Field{
								{
									ID:   1,
									Name: "sadness",
									Type: TypeReference{
										Name: "GreatSadness",
										Line: 8,
									},
									Line: 8,
								},
							},
							OneWay: false,
							Line:   7,
						},
						{
							Name: "somethingElse",
							Parameters: []*Field{
								{
									ID:   1,
									Name: "a",
									Type: TypeReference{Name: "A", Line: 11},
									Line: 11,
								},
								{
									ID:   2,
									Name: "b",
									Type: TypeReference{Name: "B", Line: 12},
									Line: 12,
								},
							},
							Annotations: []*Annotation{
								{
									Name:  "py.name",
									Value: "something_else",
									Line:  13,
								},
							},
							Line: 10,
						},
					},
					Annotations: []*Annotation{
						{
							Name:  "ttl.milliseconds",
							Value: "200",
							Line:  14,
						},
					},
					Line: 2,
				},
			}},
		},
	}

	assertParseCases(t, tests)
}
