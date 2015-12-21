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

	. "github.com/uber/thriftrw-go/ast"

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

func TestParseErrors(t *testing.T) {
	tests := []string{
		"namespace foo \x00",
		`const string 42 = "foo"`,
		`typedef foo bar baz`,
		`typedef foo`,
		`enum Foo {`,
		`enum { }`,
		`
			enum Foo {}
			include "bar.thrift"
		`,
		`service Foo extends {}`,
		`service Foo Bar {}`,
		`service Foo { void foo() () (foo = "bar") }`,
		`service Foo { void foo() throws }`,
		`typedef string (foo =) UUID`,
	}

	for _, tt := range tests {
		_, err := Parse([]byte(tt))
		assert.Error(t, err, "Expected error while parsing:\n%s", tt)
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
				&Include{Path: "foo.thrift", ILine: 2},
				&Include{Path: "bar.thrift", Name: "t", ILine: 3},
			}},
		},
		{
			`
				namespace py bar
				namespace * foo
			`,
			&Program{Headers: []Header{
				&Namespace{"py", "bar", 2},
				&Namespace{"*", "foo", 3},
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
					&Include{Path: "shared.thrift", ILine: 3},
					&Namespace{"go", "foo_service", 4},
					&Include{Path: "errors.thrift", ILine: 9},
					&Namespace{"py", "services.foo", 12},
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
			`,
			&Program{Definitions: []Definition{
				&Constant{
					CName: "foo",
					Type:  BaseType{ID: I32TypeID},
					Value: ConstantInteger(42),
					CLine: 2,
				},
				&Constant{
					CName: "bar",
					Type:  BaseType{ID: I64TypeID},
					Value: ConstantReference{
						Name: "shared.baz",
						Line: 3,
					},
					CLine: 3,
				},
				&Constant{
					CName: "baz",
					Type:  BaseType{ID: StringTypeID},
					Value: ConstantString("hello world"),
					CLine: 5,
				},
				&Constant{
					CName: "qux",
					Type:  BaseType{ID: DoubleTypeID},
					Value: ConstantDouble(3.141592),
					CLine: 7,
				},
			}},
		},
		{
			`const bool (foo = "a\nb") baz = true
			 const bool include_something = false`,
			&Program{Definitions: []Definition{
				&Constant{
					CName: "baz",
					Type: BaseType{
						ID: BoolTypeID,
						Annotations: []*Annotation{
							&Annotation{Name: "foo", Value: "a\nb", Line: 1},
						},
					},
					Value: ConstantBoolean(true),
					CLine: 1,
				},
				&Constant{
					CName: "include_something",
					Type:  BaseType{ID: BoolTypeID},
					Value: ConstantBoolean(false),
					CLine: 2,
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
					CName: "stuff",
					Type: MapType{
						KeyType: BaseType{
							ID: StringTypeID,
							Annotations: []*Annotation{
								&Annotation{Name: "foo", Value: "", Line: 2},
							},
						},
						ValueType: BaseType{ID: I32TypeID},
						Annotations: []*Annotation{
							&Annotation{Name: "baz", Value: "qux", Line: 2},
						},
					},
					Value: ConstantMap{
						Items: []ConstantMapItem{
							ConstantMapItem{
								Key:   ConstantString("a"),
								Value: ConstantInteger(1),
							},
							ConstantMapItem{
								Key:   ConstantString("b"),
								Value: ConstantInteger(2),
							},
						},
					},
					CLine: 2,
				},
				&Constant{
					CName: "list_of_lists",
					Type: ListType{ValueType: ListType{
						ValueType: BaseType{ID: I32TypeID},
					}},
					Value: ConstantList{
						Items: []ConstantValue{
							ConstantList{
								Items: []ConstantValue{
									ConstantInteger(1),
									ConstantInteger(2),
									ConstantInteger(3),
								},
							},
							ConstantList{
								Items: []ConstantValue{
									ConstantInteger(4),
									ConstantInteger(5),
									ConstantInteger(6),
								},
							},
						},
					},
					CLine: 6,
				},
				&Constant{
					CName: "const_struct",
					Type:  TypeReference{Name: "Item", Line: 10},
					Value: ConstantMap{Items: []ConstantMapItem{
						ConstantMapItem{
							Key:   ConstantString("key"),
							Value: ConstantString("foo"),
						},
						ConstantMapItem{
							Key:   ConstantString("value"),
							Value: ConstantInteger(42),
						},
					}},
					CLine: 10,
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
					CName: "foo",
					Type:  BaseType{ID: StringTypeID},
					Value: ConstantString(`a "b" c`),
					CLine: 2,
				},
				&Constant{
					CName: "bar",
					Type:  BaseType{ID: StringTypeID},
					Value: ConstantString(`a 'b' c`),
					CLine: 3,
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
			`,
			&Program{Definitions: []Definition{
				&Typedef{
					TName: "UUID",
					Type:  BaseType{ID: StringTypeID},
					Annotations: []*Annotation{
						&Annotation{
							Name:  "length",
							Value: "32",
							Line:  2,
						},
					},
					TLine: 2,
				},
				&Typedef{
					TName: "Date",
					Type: BaseType{
						ID: I64TypeID,
						Annotations: []*Annotation{
							&Annotation{
								Name:  "js.type",
								Value: "Date",
								Line:  4,
							},
						},
					},
					TLine: 4,
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
			&Program{Definitions: []Definition{&Enum{EName: "EmptyEnum", ELine: 2}}},
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
					EName: "SillyEnum",
					Items: []*EnumItem{
						&EnumItem{
							Name: "foo",
							Annotations: []*Annotation{
								&Annotation{
									Name:  "x",
									Value: "",
									Line:  3,
								},
							},
							Line: 3,
						},
						&EnumItem{Name: "bar", Line: 3},
						&EnumItem{Name: "baz", Value: &aValue, Line: 4},
						&EnumItem{Name: "qux", Line: 5},
						&EnumItem{Name: "quux", Line: 6},
					},
					Annotations: []*Annotation{
						&Annotation{Name: "_", Value: "__", Line: 7},
						&Annotation{Name: "foo", Value: "bar", Line: 7},
					},
					ELine: 2,
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
				&Struct{SName: "EmptyStruct", Type: StructType, SLine: 2},
				&Struct{SName: "EmptyUnion", Type: UnionType, SLine: 3},
				&Struct{SName: "EmptyExc", Type: ExceptionType, SLine: 4},
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
					SName: "i128",
					Type:  StructType,
					Fields: []*Field{
						&Field{
							ID:           1,
							Name:         "high",
							Type:         BaseType{ID: I64TypeID},
							Requiredness: Required,
							Line:         3,
						},
						&Field{
							ID:           2,
							Name:         "low",
							Type:         BaseType{ID: I64TypeID},
							Requiredness: Required,
							Line:         4,
						},
					},
					Annotations: []*Annotation{
						&Annotation{
							Name:  "serializer",
							Value: "Int128Serializer",
							Line:  5,
						},
					},
					SLine: 2,
				},
				&Struct{
					SName: "Contents",
					Type:  UnionType,
					Fields: []*Field{
						&Field{
							ID:           1,
							Name:         "plainText",
							Requiredness: Unspecified,
							Type: BaseType{
								ID: StringTypeID,
								Annotations: []*Annotation{
									&Annotation{
										Name:  "format",
										Value: "markdown",
										Line:  8,
									},
								},
							},
							Line: 8,
						},
						&Field{
							ID:   2,
							Name: "pdf",
							Type: BaseType{ID: BinaryTypeID},
							// Requiredness intentionally skipped because
							// zero-value for it is Unspecified.
							Annotations: []*Annotation{
								&Annotation{
									Name:  "name",
									Value: "pdfFile",
									Line:  9,
								},
							},
							Line: 9,
						},
					},
					SLine: 7,
				},
				&Struct{
					SName: "GreatSadness",
					Type:  ExceptionType,
					Fields: []*Field{
						&Field{
							ID:           1,
							Name:         "message",
							Type:         BaseType{ID: StringTypeID},
							Requiredness: Optional,
							Line:         13,
						},
					},
					SLine: 12,
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
				&Service{SName: "EmptyService", SLine: 2},
				&Service{
					SName: "AnotherEmptyService",
					Parent: &ServiceReference{
						Name: "EmptyService",
						Line: 3,
					},
					SLine: 3,
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
					SName: "KeyValue",
					Functions: []*Function{
						&Function{
							Name:   "empty",
							OneWay: true,
							Line:   4,
						},
						&Function{
							Name:       "something",
							ReturnType: BaseType{ID: I32TypeID},
							Exceptions: []*Field{
								&Field{
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
						&Function{
							Name: "somethingElse",
							Parameters: []*Field{
								&Field{
									ID:   1,
									Name: "a",
									Type: TypeReference{Name: "A", Line: 11},
									Line: 11,
								},
								&Field{
									ID:   2,
									Name: "b",
									Type: TypeReference{Name: "B", Line: 12},
									Line: 12,
								},
							},
							Annotations: []*Annotation{
								&Annotation{
									Name:  "py.name",
									Value: "something_else",
									Line:  13,
								},
							},
							Line: 10,
						},
					},
					Annotations: []*Annotation{
						&Annotation{
							Name:  "ttl.milliseconds",
							Value: "200",
							Line:  14,
						},
					},
					SLine: 2,
				},
			}},
		},
	}

	assertParseCases(t, tests)
}
