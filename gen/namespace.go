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

	"go.uber.org/thriftrw/internal/goast"
)

// Namespace allows finding names that don't conflict with other names in
// certain scopes.
type Namespace interface {
	// Generates a new name based on the given name. Prefers the provided name
	// but the name may be altered if it conflicts with another name in the
	// current scope.
	//
	// The returned name is also reserved so that future names cannot conflict
	// with it.
	NewName(name string) string

	// Reserve the given name in this namespace. Fails with an error if the name
	// is already taken.
	Reserve(name string) error

	// Forget the given name. Nothing happens if this name wasn't already
	// taken in this namespace. The name will NOT be removed from parent
	// namespaces.
	Forget(name string)

	// Create a new Child namespace. The child namespace cannot use any names
	// defined in this namespace or any of its parent namespaces.
	Child() Namespace

	// Rotate returns the oldest known name that was used for the given
	// requested name, and moves it to the end of the list of names used for
	// that requested name.
	//
	// This may be used to access the names used in a template dynamically in a
	// deterministic order if recording them in local variables is not possible.
	Rotate(name string) (string, error)
}

// namespace helps reserve names within a scope with support for child
// namespaces that do not attempt to shadow names from the parent namespace.
type namespace struct {
	parent *namespace
	taken  map[string]struct{}

	// gave is a mapping from the name requested by the user to the names we
	// returned
	gave map[string][]string
}

// NewNamespace creates a new namespace.
func NewNamespace() Namespace {
	return &namespace{
		taken: make(map[string]struct{}),
		gave:  make(map[string][]string),
	}
}

func (n *namespace) isTaken(name string) bool {
	if _, ok := n.taken[name]; ok {
		return true
	}
	if n.parent != nil {
		return n.parent.isTaken(name)
	}
	return goast.IsReservedKeyword(name)
}

func (n *namespace) NewName(base string) string {
	// TODO(abg): Avoid clashing with Go keywords.
	name := base
	for i := 2; n.isTaken(name); i++ {
		name = fmt.Sprintf("%s%d", base, i)
	}
	n.taken[name] = struct{}{}
	n.gave[base] = append(n.gave[base], name)
	return name
}

func (n *namespace) Rotate(base string) (string, error) {
	names, ok := n.gave[base]
	if !ok || len(names) == 0 {
		return "", fmt.Errorf("%q is not a known name", base)
	}

	name := names[0]
	names = names[1:]
	n.gave[base] = append(names, name)
	return name, nil
}

func (n *namespace) Reserve(name string) error {
	// A single Go file can have multiple init() functions so we don't need to
	// reserve it.
	if name == "init" {
		return nil
	}

	if n.isTaken(name) {
		return namespaceError{name}
	}

	n.taken[name] = struct{}{}
	return nil
}

func (n *namespace) Forget(name string) {
	delete(n.taken, name)
}

func (n *namespace) Child() Namespace {
	return &namespace{
		parent: n,
		taken:  make(map[string]struct{}),
		gave:   make(map[string][]string),
	}
}

type namespaceError struct {
	Name string
}

func (e namespaceError) Error() string {
	return fmt.Sprintf("name %q is already taken", e.Name)
}
