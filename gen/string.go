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

package gen

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	// boundary matches word boundaries for camel case expressions.
	//
	// 	FooBar => Fo(o)(Bar)
	// 	HTTPRequest => HTT(P)(Request)
	//
	// We use this to detect if an object is already camel case.
	boundary = regexp.MustCompile(`([^_])([A-Z][a-z]+)`)
)

// goCase converts strings into standard camel case with the first letter being
// upper case.
func goCase(s string) string {
	if len(s) == 0 {
		panic(fmt.Sprintf("%q is not a valid identifier", s))
	}

	if boundary.MatchString(s) {
		// Already camelCase. Just ensure that the first letter is upper case.
		return strings.Title(s)
	}

	chunks := strings.Split(s, "_")
	for i, c := range chunks {
		chunks[i] = strings.Title(strings.ToLower(c))
	}

	// TODO(abg): If the string ends with "Id", this should replace that with
	// "ID".

	return strings.Join(chunks, "")
}
