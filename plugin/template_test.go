package plugin

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/plugin/api"
)

func TestGoFileFromTemplate(t *testing.T) {
	tests := []struct {
		desc     string
		template string
		data     interface{}
		options  []TemplateOption

		wantBody  string
		wantError string
	}{
		{
			desc:     "simple",
			template: "package foo",
			wantBody: unlines("package foo"),
		},
		{
			desc: "type reference",
			template: `
				package foo

				var foo <formatType .> = nil
			`,
			data: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:    "Foo",
					Package: "go.uber.org/thriftrw/bar",
				},
			},
			wantBody: unlines(
				`package foo`,
				``,
				`import "go.uber.org/thriftrw/bar"`,
				``,
				`var foo bar.Foo = nil`,
			),
		},
		{
			desc: "type reference in the same package",
			template: `
				package bar

				func hello() <formatType .> {
					return nil
				}
			`,
			data: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:    "Foo",
					Package: "go.uber.org/thriftrw/bar",
				},
			},
			options: []TemplateOption{
				GoFileImportPath("go.uber.org/thriftrw/bar"),
			},
			wantBody: unlines(
				`package bar`,
				``,
				`func hello() Foo {`,
				`	return nil`,
				`}`,
			),
		},
		{
			desc: "import",
			template: `
				package hello

				<$foo := import "go.uber.org/thriftrw/plugin">
				<$bar := import "go.uber.org/thriftrw/hello">

				func main() {
					<$foo>.Main(<$bar>.Baz)
				}
			`,
			wantBody: unlines(
				`package hello`,
				``,
				`import (`,
				`	"go.uber.org/thriftrw/hello"`,
				`	"go.uber.org/thriftrw/plugin"`,
				`)`,
				``,
				`func main() {`,
				`	plugin.Main(hello.Baz)`,
				`}`,
			),
		},
		{
			desc: "import conflicts",
			template: `
				package hello

				<$foo := import "context">
				<$bar := import "golang.org/x/net/context">

				// foo does stuff
				func foo() <$foo>.Context { return nil }
				func bar() <$bar>.Context { return nil }
			`,
			wantBody: unlines(
				`package hello`,
				``,
				`import (`,
				`	"context"`,
				`	context2 "golang.org/x/net/context"`,
				`)`,
				``,
				`// foo does stuff`,
				`func foo() context.Context  { return nil }`,
				`func bar() context2.Context { return nil }`,
			),
		},
		{
			desc: "import twice",
			template: `
				package foo

				<$fmt := import "fmt">

				func ErrFail(err error) error {
					return <import "fmt">.Errorf("great sadness: %v", err)
				}
			`,
			wantBody: unlines(
				`package foo`,
				``,
				`import "fmt"`,
				``,
				`func ErrFail(err error) error {`,
				`	return fmt.Errorf("great sadness: %v", err)`,
				`}`,
			),
		},
		{
			desc:      "invalid template",
			template:  `<import "`,
			wantError: `failed to parse template "test.go":`,
		},
		{
			desc:      "invalid Go code",
			template:  `func main() {}`,
			wantError: `failed to parse generated code: test.go:`,
		},
		{
			desc: "explicit import",
			template: `
				package main

				import "fmt"
			`,
			wantError: "plain imports are not allowed with GoFileFromTemplate: " +
				"use the import function",
		},
		{
			desc: "import keyword",
			template: `
				package hello

				<$foo := import "go.uber.org/thriftrw/range">

				type foo struct {
					<$foo>.Range
				}
			`,
			wantBody: unlines(
				`package hello`,
				``,
				`import range2 "go.uber.org/thriftrw/range"`,
				``,
				`type foo struct {`,
				`	range2.Range`,
				`}`,
			),
		},
	}

	for _, tt := range tests {
		got, err := GoFileFromTemplate("test.go", tt.template, tt.data, tt.options...)
		if tt.wantError != "" {
			assert.Contains(t, err.Error(), tt.wantError, tt.desc)
		} else {
			assert.Equal(t, tt.wantBody, string(got), tt.desc)
		}
	}
}

// unlines joins the given lines with newlines in between followied by a
// trailing newline.
func unlines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}
