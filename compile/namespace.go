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

package compile

import (
	"fmt"
	"strings"
)

type namespaceType func(string) string

var (
	caseSensitive   = namespaceType(func(s string) string { return s })
	caseInsensitive = namespaceType(strings.ToLower)
)

// namespace helps dole out names and avoid conflicts.
//
// Claim with the same name twice on a namespace will result in a nameConflict
// error.
type namespace struct {
	transform func(string) string
	names     map[string]int
}

// newNamespace instantiates a new namespace.
//
// 	ns := newNamespace(caseSensitive)
// 	ns.claim("foo")
func newNamespace(t namespaceType) namespace {
	return namespace{
		transform: t,
		names:     make(map[string]int),
	}
}

// claim requests the given name in the namespace. If the name is already
// claimed, an error will be returned.
func (n namespace) claim(name string, line int) error {
	s := n.transform(name)
	if line, ok := n.names[s]; ok {
		return nameConflict{name: name, line: line}
	}
	n.names[s] = line
	return nil
}

// nameConflict is raised when the name for an identifier conflicts with a
// name that has already been used.
type nameConflict struct {
	name string
	line int
}

func (e nameConflict) Error() string {
	return fmt.Sprintf(
		"the name %q has already been used on line %d", e.name, e.line,
	)
}
