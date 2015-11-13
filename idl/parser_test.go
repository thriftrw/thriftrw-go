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
				include "bar.thrift"
			`,
			&Program{Includes: []*Include{
				&Include{"foo.thrift", 2},
				&Include{"bar.thrift", 3},
			}},
		},
		{
			`
				namespace py bar
				namespace * foo
			`,
			&Program{Namespaces: []*Namespace{
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
				Includes: []*Include{
					&Include{"shared.thrift", 3},
					&Include{"errors.thrift", 9},
				},
				Namespaces: []*Namespace{
					&Namespace{"go", "foo_service", 4},
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
			&Program{Constants: []*Constant{
				&Constant{
					Name:  "foo",
					Type:  BaseType{ID: I32TypeID},
					Value: ConstantInteger(42),
					Line:  2,
				},
				&Constant{
					Name: "bar",
					Type: BaseType{ID: I64TypeID},
					Value: ConstantReference{
						Name: "shared.baz",
						Line: 3,
					},
					Line: 3,
				},
				&Constant{
					Name:  "baz",
					Type:  BaseType{ID: StringTypeID},
					Value: ConstantString("hello world"),
					Line:  5,
				},
				&Constant{
					Name:  "qux",
					Type:  BaseType{ID: DoubleTypeID},
					Value: ConstantDouble(3.141592),
					Line:  7,
				},
			}},
		},
		{
			`const bool (foo = "a\nb") baz = true
			 const bool include_something = false`,
			&Program{Constants: []*Constant{
				&Constant{
					Name: "baz",
					Type: BaseType{
						ID: BoolTypeID,
						Annotations: []*Annotation{
							&Annotation{Name: "foo", Value: "a\nb", Line: 1},
						},
					},
					Value: ConstantBoolean(true),
					Line:  1,
				},
				&Constant{
					Name:  "include_something",
					Type:  BaseType{ID: BoolTypeID},
					Value: ConstantBoolean(false),
					Line:  2,
				},
			}},
		},
		{
			`
				const map<string (foo = "bar"), i32> (baz = "qux") stuff = {
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
			&Program{Constants: []*Constant{
				&Constant{
					Name: "stuff",
					Type: MapType{
						KeyType: BaseType{
							ID: StringTypeID,
							Annotations: []*Annotation{
								&Annotation{Name: "foo", Value: "bar", Line: 2},
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
					Line: 2,
				},
				&Constant{
					Name: "list_of_lists",
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
					Line: 6,
				},
				&Constant{
					Name: "const_struct",
					Type: TypeReference{Name: "Item", Line: 10},
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
					Line: 10,
				},
			}},
		},
		{
			`
				const string foo = 'a "b" c'
				const string bar = "a 'b' c"
			`,
			&Program{Constants: []*Constant{
				&Constant{
					Name:  "foo",
					Type:  BaseType{ID: StringTypeID},
					Value: ConstantString(`a "b" c`),
					Line:  2,
				},
				&Constant{
					Name:  "bar",
					Type:  BaseType{ID: StringTypeID},
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
			`,
			&Program{Typedefs: []*Typedef{
				&Typedef{
					Name: "UUID",
					Type: BaseType{ID: StringTypeID},
					Annotations: []*Annotation{
						&Annotation{
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
						ID: I64TypeID,
						Annotations: []*Annotation{
							&Annotation{
								Name:  "js.type",
								Value: "Date",
								Line:  4,
							},
						},
					},
					Line: 4,
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
			&Program{Enums: []*Enum{&Enum{Name: "EmptyEnum", Line: 2}}},
		},
		{
			`
				enum SillyEnum {
					foo (x = "y"), bar /*
					*/ baz = 42
					qux;
					quux
				} (_ = "__", foo = "bar")
			`,
			&Program{Enums: []*Enum{
				&Enum{
					Name: "SillyEnum",
					Items: []*EnumItem{
						&EnumItem{
							Name: "foo",
							Annotations: []*Annotation{
								&Annotation{
									Name:  "x",
									Value: "y",
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
			&Program{Structs: []*Struct{
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
			&Program{Structs: []*Struct{
				&Struct{
					Name: "i128",
					Type: StructType,
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
					Line: 2,
				},
				&Struct{
					Name: "Contents",
					Type: UnionType,
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
					Line: 7,
				},
				&Struct{
					Name: "GreatSadness",
					Type: ExceptionType,
					Fields: []*Field{
						&Field{
							ID:           1,
							Name:         "message",
							Type:         BaseType{ID: StringTypeID},
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
			&Program{Services: []*Service{
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
			&Program{Services: []*Service{
				&Service{
					Name: "KeyValue",
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
					Line: 2,
				},
			}},
		},
	}

	assertParseCases(t, tests)
}
