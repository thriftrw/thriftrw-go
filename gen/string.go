// Copyright (c) 2026 Uber Technologies, Inc.
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

package gen

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/gocase"
)

func isAllCaps(s string) bool {
	return gocase.IsAllCaps(s)
}

func pascalCase(allowAllCaps bool, words ...string) string {
	return gocase.PascalCase(allowAllCaps, words...)
}

func constantName(s string) string {
	return gocase.PascalCase(false /* all caps */, strings.Split(s, "_")...)
}

func goCase(s string) string {
	return gocase.GoCase(s)
}

// goNameAnnotation returns ("", nil) if there is no "go.name" annotation.
func goNameAnnotation(e compile.NamedEntity) (string, error) {
	name, ok := e.ThriftAnnotations()["go.name"]
	if !ok {
		return "", nil
	}

	c, _ := utf8.DecodeRuneInString(name)
	capitalized := unicode.IsLetter(c) && unicode.IsUpper(c)
	underscore := strings.Contains(name, "_")

	if !capitalized || underscore {
		var emsg []string
		if underscore {
			emsg = append(emsg, "contains underscores")
		}
		if !capitalized {
			emsg = append(emsg, "is not capitalized")
		}

		return "", fmt.Errorf("%q (from go.name annotation) is not a Go style public identifier (%s), suggestion: %q)", name, strings.Join(emsg, ", "), goCase(name))
	}

	return name, nil
}

func goNameForNamedEntity(e compile.NamedEntity) (name string, fromAnnotation bool, err error) {
	fromAnnotation = true
	name, err = goNameAnnotation(e)
	if err == nil && name == "" {
		name = goCase(e.ThriftName())
		fromAnnotation = false
	}
	return name, fromAnnotation, err
}

func goName(e compile.NamedEntity) (string, error) {
	name, _, err := goNameForNamedEntity(e)
	return name, err
}

var commonInitialisms = gocase.CommonInitialisms
