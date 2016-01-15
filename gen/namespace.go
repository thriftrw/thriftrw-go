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

import "fmt"

// namespace helps reserve names within a scope with support for child
// namespaces that do not attempt to shadow names from the parent namespace.
type namespace struct {
	parent *namespace
	taken  map[string]struct{}
}

func newNamespace() *namespace {
	return &namespace{taken: make(map[string]struct{})}
}

func (n *namespace) isTaken(name string) bool {
	_, ok := n.taken[name]
	if ok {
		return true
	}
	if n.parent != nil {
		return n.parent.isTaken(name)
	}
	return false
}

// New generates a new name based on the given name.
func (n *namespace) NewName(base string) string {
	// TODO(abg): Avoid clashing with Go keywords.
	name := base
	for i := 2; n.isTaken(name); i++ {
		name = fmt.Sprintf("%s%d", base, i)
	}
	n.taken[name] = struct{}{}
	return name
}

// Reserve reserves the given name with the namespace.
//
// An error is returned if the name was already taken.
func (n *namespace) Reserve(name string) error {
	if n.isTaken(name) {
		return namespaceError{name}
	}
	n.taken[name] = struct{}{}
	return nil
}

// Child creates a new child namespace from the given namespace. The child
// namespace will not be allowed to use names that area already taken by this
// namespace or any of its parents.
func (n *namespace) Child() *namespace {
	return &namespace{parent: n, taken: make(map[string]struct{})}
}

type namespaceError struct {
	Name string
}

func (e namespaceError) Error() string {
	return fmt.Sprintf("name %q is already taken", e.Name)
}
