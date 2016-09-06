// Copyright (c) 2016 Uber Technologies, Inc.
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
func pascalCase(words ...string) string {
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
		if isAllCaps(chunk) && len(words) > 1 {
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

// enumItemName returns the Go name that should be used for an enum item with
// the given Thrift name.
func enumItemName(enumName string, itemName string) string {
	words := []string{enumName}
	words = append(words, strings.Split(itemName, "_")...)
	return pascalCase(words...)
}

// goCase converts strings into PascalCase.
func goCase(s string) string {
	if len(s) == 0 {
		panic(fmt.Sprintf("%q is not a valid identifier", s))
	}

	return pascalCase(strings.Split(s, "_")...)
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
