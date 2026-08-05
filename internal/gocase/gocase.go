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

// Package gocase converts Thrift identifiers into PascalCase Go identifiers,
// handling common initialisms (ID, UUID, HTTP, ...) the same way the rest of
// thriftrw's code generation does. It is shared by the gen and plugin packages
// so that plugins can produce identifiers that match thriftrw-generated code.
package gocase

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// GoCase converts a Thrift identifier into a PascalCase Go identifier.
//
// A single all-caps word is left unchanged (e.g. "FOO" stays "FOO"), but
// SCREAMING_SNAKE_CASE with underscores is title-cased (e.g. "FOO_BAR"
// becomes "FooBar"). Known initialisms like "id" are uppercased ("ID").
func GoCase(s string) string {
	if len(s) == 0 {
		panic(fmt.Sprintf("%q is not a valid identifier", s))
	}

	words := strings.Split(s, "_")
	return PascalCase(len(words) == 1 /* allowAllCaps */, words...)
}

// PascalCase combines the given words using PascalCase.
//
// If allowAllCaps is true, when an all-caps word that is not a known
// abbreviation is encountered, it is left unchanged. Otherwise, it is
// Titlecased.
func PascalCase(allowAllCaps bool, words ...string) string {
	for i, chunk := range words {
		if len(chunk) == 0 {
			continue
		}

		init := strings.ToUpper(chunk)
		if _, ok := CommonInitialisms[init]; ok {
			words[i] = init
			continue
		}

		if IsAllCaps(chunk) && !allowAllCaps {
			words[i] = strings.Title(strings.ToLower(chunk))
			continue
		}

		head, headIndex := utf8.DecodeRuneInString(chunk)
		words[i] = string(unicode.ToUpper(head)) + string(chunk[headIndex:])
	}

	return strings.Join(words, "")
}

// IsAllCaps checks if a string contains all capital letters only. Non-letters
// are not considered.
func IsAllCaps(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

// CommonInitialisms is the set of known abbreviations that are kept fully
// uppercased when converting to Go identifiers.
//
// This set is taken from https://github.com/golang/lint/blob/master/lint.go#L692
var CommonInitialisms = map[string]bool{
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
