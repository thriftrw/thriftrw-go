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

import "fmt"

// Scope represents a queryable compilation scope.
//
// All Lookup methods must only return types defined in this scope. References
// to items defined in included modules will be resolved by the caller.
type Scope interface {
	// GetName retrieves gets this scope.
	GetName() string

	// Retrieve a type defined in this scope.
	LookupType(name string) (TypeSpec, error)

	// Retrieve a service defined in this scope.
	LookupService(name string) (*ServiceSpec, error)

	// Retrieve a constant defined in this scope.
	LookupConstant(name string) (*Constant, error)

	// Retrieve an included scope.
	LookupInclude(name string) (Scope, error)
}

// getIncludedScope retrieves an included scope from the given scope.
func getIncludedScope(scope Scope, name string) (Scope, error) {
	included, err := scope.LookupInclude(name)
	if err != nil {
		return nil, unrecognizedModuleError{Name: name, Reason: err}
	}
	return included, nil
}

// EmptyScope returns a Scope that fails all lookups.
func EmptyScope(name string) Scope {
	return emptyScope{name}
}

type emptyScope struct{ name string }

func (e emptyScope) GetName() string { return e.name }

func (emptyScope) LookupType(name string) (TypeSpec, error) {
	return nil, fmt.Errorf("unknown type: %v", name)
}

func (emptyScope) LookupService(name string) (*ServiceSpec, error) {
	return nil, fmt.Errorf("unknown service: %v", name)
}

func (emptyScope) LookupConstant(name string) (*Constant, error) {
	return nil, fmt.Errorf("unknown constant: %v", name)
}

func (emptyScope) LookupInclude(name string) (Scope, error) {
	return nil, fmt.Errorf("unknown include: %v", name)
}
