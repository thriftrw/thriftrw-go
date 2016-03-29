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

package compile

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
