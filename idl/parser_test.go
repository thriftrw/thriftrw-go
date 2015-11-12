package idl

import (
	"testing"

	. "github.com/uber/thriftrw-go/ast"

	"github.com/stretchr/testify/assert"
)

type parseCase struct {
	program  *Program
	document string
}

func assertParseCases(t *testing.T, tests []parseCase) {
	for _, tt := range tests {
		program, err := Parse([]byte(tt.document))
		if assert.NoError(t, err, "Parsing failed:\n%s", tt.document) {
			assert.Equal(
				t, tt.program, program,
				"Got unexpected program when parsing:\n%s", tt.document,
			)
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
	}

	for _, tt := range tests {
		_, err := Parse([]byte(tt))
		assert.Error(t, err, "Expected error while parsing:\n%s", tt)
	}
}

func TestParseHeaders(t *testing.T) {
	tests := []parseCase{
		{
			&Program{Includes: []*Include{
				&Include{"foo.thrift", 2},
				&Include{"bar.thrift", 3},
			}},
			`
				include "foo.thrift"
				include "bar.thrift"
			`,
		},
		{
			&Program{Namespaces: []*Namespace{
				&Namespace{"py", "bar", 2},
				&Namespace{"*", "foo", 3},
			}},
			`
				namespace py bar
				namespace * foo
			`,
		},
		{
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
		},
	}
	assertParseCases(t, tests)
}

func TestParseConstants(t *testing.T) {
	tests := []parseCase{
		{
			&Program{Constants: []*Constant{
				&Constant{
					Name:  "foo",
					Type:  BaseType{ID: I32BaseTypeID},
					Value: ConstantValue(ConstantInteger(42)),
					Line:  2,
				},
				&Constant{
					Name: "bar",
					Type: BaseType{ID: I64BaseTypeID},
					Value: ConstantValue(
						ConstantReference{
							Name: "shared.baz",
							Line: 3,
						},
					),
					Line: 3,
				},
				&Constant{
					Name:  "baz",
					Type:  BaseType{ID: StringBaseTypeID},
					Value: ConstantValue(ConstantString("hello world")),
					Line:  5,
				},
				&Constant{
					Name:  "qux",
					Type:  BaseType{ID: DoubleBaseTypeID},
					Value: ConstantValue(ConstantDouble(3.141592)),
					Line:  7,
				},
			}},
			`
				const i32 foo = 42
				const i64 bar = shared.baz;

				const string baz = "hello world";

				const double qux = 3.141592
			`,
		},
		{
			&Program{Constants: []*Constant{
				&Constant{
					Name: "baz",
					Type: BaseType{
						ID: BoolBaseTypeID,
						Annotations: []*Annotation{
							&Annotation{Name: "foo", Value: "a\nb", Line: 1},
						},
					},
					Value: ConstantValue(ConstantBoolean(true)),
					Line:  1,
				},
				&Constant{
					Name:  "include_something",
					Type:  BaseType{ID: BoolBaseTypeID},
					Value: ConstantValue(ConstantBoolean(false)),
					Line:  2,
				},
			}},
			`const bool (foo = "a\nb") baz = true
			 const bool include_something = false`,
		},
		{
			&Program{Constants: []*Constant{
				&Constant{
					Name: "stuff",
					Type: MapType{
						KeyType: BaseType{
							ID: StringBaseTypeID,
							Annotations: []*Annotation{
								&Annotation{Name: "foo", Value: "bar", Line: 2},
							},
						},
						ValueType: BaseType{ID: I32BaseTypeID},
						Annotations: []*Annotation{
							&Annotation{Name: "baz", Value: "qux", Line: 2},
						},
					},
					Value: ConstantValue(ConstantMap{
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
					}),
					Line: 2,
				},
				&Constant{
					Name: "list_of_lists",
					Type: ListType{ValueType: ListType{
						ValueType: BaseType{ID: I32BaseTypeID},
					}},
					Value: ConstantValue(ConstantList{
						Items: []ConstantValue{
							ConstantValue(ConstantList{
								Items: []ConstantValue{
									ConstantInteger(1),
									ConstantInteger(2),
									ConstantInteger(3),
								},
							}),
							ConstantValue(ConstantList{
								Items: []ConstantValue{
									ConstantInteger(4),
									ConstantInteger(5),
									ConstantInteger(6),
								},
							}),
						},
					}),
					Line: 6,
				},
			}},
			`
				const map<string (foo = "bar"), i32> (baz = "qux") stuff = {
					"a": 1,
					"b": 2,
				}
				const list<list<i32>> list_of_lists = [
					[1, 2, 3]  # optional separator
					[4, 5, 6]
				];
			`,
		},
		{
			&Program{Constants: []*Constant{
				&Constant{
					Name:  "foo",
					Type:  BaseType{ID: StringBaseTypeID},
					Value: ConstantValue(ConstantString(`a "b" c`)),
					Line:  2,
				},
				&Constant{
					Name:  "bar",
					Type:  BaseType{ID: StringBaseTypeID},
					Value: ConstantValue(ConstantString(`a 'b' c`)),
					Line:  3,
				},
			}},
			`
				const string foo = 'a "b" c'
				const string bar = "a 'b' c"
			`,
		},
	}
	assertParseCases(t, tests)
}

func TestParseTypedef(t *testing.T) {
	tests := []parseCase{
		{
			&Program{Typedefs: []*Typedef{
				&Typedef{
					Name: "UUID",
					Type: BaseType{ID: StringBaseTypeID},
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
						ID: I64BaseTypeID,
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
			`
				typedef string UUID (length = "32");

				typedef i64 (js.type = "Date") Date
			`,
		},
	}

	assertParseCases(t, tests)
}

func TestParseEnum(t *testing.T) {
	aValue := 42

	tests := []parseCase{
		{
			&Program{Enums: []*Enum{&Enum{Name: "EmptyEnum", Line: 1}}},
			`enum EmptyEnum {}`,
		},
		{
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
			`
				enum SillyEnum {
					foo (x = "y"), bar /*
					*/ baz = 42
					qux;
					quux
				} (_ = "__", foo = "bar")
			`,
		},
	}

	assertParseCases(t, tests)
}
