package plugin

import (
	"sort"
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
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/bar",
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
			desc: "type reference with annotations",
			template: `
				package foo

				var foo <formatType .> = nil
				var fooAnnotations = map[string]string{
					<range $pair := typeAnnotations .>"<$pair.Key>": "<$pair.Value>",
					<end>
				}
			`,
			data: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/bar",
					Annotations: map[string]string{
						"foo": "bar",
						"baz": "bat",
					},
				},
			},
			options: []TemplateOption{
				TemplateFunc("typeAnnotations", typeAnnotations),
			},
			wantBody: unlines(
				`package foo`,
				``,
				`import "go.uber.org/thriftrw/bar"`,
				``,
				`var foo bar.Foo = nil`,
				`var fooAnnotations = map[string]string{`,
				"\t"+`"baz": "bat",`,
				"\t"+`"foo": "bar",`,
				`}`,
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
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/bar",
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
		{
			desc: "import with dash",
			template: `
				package hello

				<$yarpc := import "github.com/yarpc/yarpc-go">

				func Hello() <$yarpc>.ReqMeta {
					return nil
				}
			`,
			wantBody: unlines(
				`package hello`,
				``,
				`import yarpc "github.com/yarpc/yarpc-go"`,
				``,
				`func Hello() yarpc.ReqMeta {`,
				`	return nil`,
				`}`,
			),
		},
		{
			desc: "import with dash and conflict",
			template: `
				package hello

				<$foo1 := import "go.uber.org/thriftrw/foo-bar">
				<$foo2 := import "go.uber.org/thriftrw/foo_bar">

				var x <$foo1>.Foo1 = <$foo2>.Foo2
			`,
			wantBody: unlines(
				`package hello`,
				``,
				`import (`,
				`	foo_bar "go.uber.org/thriftrw/foo-bar"`,
				`	foo_bar2 "go.uber.org/thriftrw/foo_bar"`,
				`)`,
				``,
				`var x foo_bar.Foo1 = foo_bar2.Foo2`,
			),
		},
	}

	for _, tt := range tests {
		got, err := GoFileFromTemplate("test.go", tt.template, tt.data, tt.options...)
		if tt.wantError != "" {
			assert.Contains(t, err.Error(), tt.wantError, tt.desc)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, string(got), tt.desc)
		}
	}
}

// unlines joins the given lines with newlines in between followied by a
// trailing newline.
func unlines(lines ...string) string {
	return strings.Join(lines, "\n") + "\n"
}

// annotationPair is a key/value pair of an annotation.
//
// This is needed as Golang maps are not sorted, and we want deterministic
// output for our generated files, so we return annotations as a slice of
// annotationPairs sorted on key.
type annotationPair struct {
	Key   string
	Value string
}

// typeAnnotations returns the annotations for the api.Type.
//
// Only api.TypeReferences have annotations, so this returns nil if the given
// api.Type is not an api.TypeReference.
func typeAnnotations(t *api.Type) []annotationPair {
	if t.ReferenceType != nil && len(t.ReferenceType.Annotations) > 0 {
		var keys []string
		for key := range t.ReferenceType.Annotations {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		annotationPairs := make([]annotationPair, 0, len(keys))
		for _, key := range keys {
			annotationPairs = append(annotationPairs, annotationPair{
				Key:   key,
				Value: t.ReferenceType.Annotations[key],
			})
		}
		return annotationPairs
	}
	return nil
}
