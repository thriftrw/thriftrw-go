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

package compile

import (
	"path/filepath"
	"strings"
	"unicode"
)

// capitalize changes the first letter of a string to upper case.
func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	x := []rune(s)
	return string(unicode.ToUpper(x[0])) + string(x[1:])
	// TODO(abg): surely there has to be a better way
}

// fileBaseName returns the base name of the given file without the extension.
func fileBaseName(p string) string {
	p = filepath.Base(p)
	if i := strings.LastIndexByte(p, '.'); i > 0 {
		p = p[:i]
	}
	return p
}

// splitInclude splits the given string at the first ".".
//
// If the given string doesn't have a '.', the second returned string contains
// the entirety of it.
func splitInclude(s string) (string, string) {
	if i := strings.IndexRune(s, '.'); i > 0 {
		return s[:i], s[i+1:]
	}
	return "", s
}
