// Copyright (c) 2019 Uber Technologies, Inc.
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
)

// isAllCaps checks if a string contains all capital letters only. Non-letters
// are not considered.
func isAllCaps(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// pascalCase combines the given words using PascalCase.
//
// If allowAllCaps is true, when an all-caps word that is not a known
// abbreviation is encountered, it is left unchanged. Otherwise, it is
// Titlecased.
func pascalCase(allowAllCaps bool, words ...string) string {
	for i, chunk := range words {
		if len(chunk) == 0 {
			// foo__bar
			continue
		}

		// known initalism
		init := strings.ToUpper(chunk)
		if _, ok := commonInitialisms[init]; ok {
			words[i] = init
			continue
		}

		// Was SCREAMING_SNAKE_CASE and not a known initialism so Titlecase it.
		if isAllCaps(chunk) && !allowAllCaps {
			// A single ALLCAPS word does not count as SCREAMING_SNAKE_CASE.
			// There must be at least one underscore.
			words[i] = strings.Title(strings.ToLower(chunk))
			continue
		}

		// Just another word, but could already be camelCased somehow, so just
		// change the first letter.
		head, headIndex := utf8.DecodeRuneInString(chunk)
		words[i] = string(unicode.ToUpper(head)) + string(chunk[headIndex:])
	}

	return strings.Join(words, "")
}

func constantName(s string) string {
	return pascalCase(false /* all caps */, strings.Split(s, "_")...)
}

// goCase converts strings into PascalCase.
func goCase(s string) string {
	if len(s) == 0 {
		panic(fmt.Sprintf("%q is not a valid identifier", s))
	}

	words := strings.Split(s, "_")
	return pascalCase(len(words) == 1 /* all caps */, words...)
	// goCase allows all caps only if the string is a single all caps word.
	// That is, "FOO" is allowed but "FOO_BAR" is changed to "FooBar".
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

// This set is taken from https://github.com/golang/lint/blob/master/lint.go#L692
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}
