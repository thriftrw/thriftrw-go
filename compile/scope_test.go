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

import "fmt"

var defaultScope = EmptyScope("fake")

// fakeScope is an implementation of Scope for testing. Instances may be
// constructed easily with the scope() function.
type fakeScope struct {
	name      string
	types     map[string]TypeSpec
	services  map[string]*ServiceSpec
	constants map[string]*Constant
	includes  map[string]Scope
}

func (s fakeScope) GetName() string {
	return s.name
}

func (s fakeScope) LookupType(name string) (TypeSpec, error) {
	if t, ok := s.types[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("unknown type: %s", name)
}

func (s fakeScope) LookupService(name string) (*ServiceSpec, error) {
	if svc, ok := s.services[name]; ok {
		return svc, nil
	}
	return nil, fmt.Errorf("unknown service: %s", name)
}

func (s fakeScope) LookupConstant(name string) (*Constant, error) {
	if c, ok := s.constants[name]; ok {
		return c, nil
	}
	return nil, fmt.Errorf("unknown constant: %s", name)
}

func (s fakeScope) LookupInclude(name string) (Scope, error) {
	if i, ok := s.includes[name]; ok {
		return i, nil
	}
	return nil, fmt.Errorf("unknown include: %s", name)
}

// scopeOrDefault accepts an optional scope as an argument and returns a default
// empty scope if the given scope was nil.
func scopeOrDefault(s Scope) Scope {
	if s == nil {
		s = defaultScope
	}
	return s
}

// Helper to construct Scopes from the given pairs of items.
//
// An even number of items must be given.
//
// An odd number of arguments may be given if the first one is the scope name.
// Otherwise a default scope name of "fake" will be used.
func scope(args ...interface{}) Scope {
	scopeName := "fake"
	if len(args)%2 == 1 {
		scopeName = args[0].(string)
		args = args[1:]
	}

	if len(args)%2 != 0 {
		panic("scope() expects an even number of arguments after the name")
	}

	scope := fakeScope{
		name:      scopeName,
		types:     make(map[string]TypeSpec),
		services:  make(map[string]*ServiceSpec),
		constants: make(map[string]*Constant),
		includes:  make(map[string]Scope),
	}

	var name string

	flag := false
	for _, arg := range args {
		flag = !flag
		if flag {
			name = arg.(string)
			continue
		}

		switch v := arg.(type) {
		case TypeSpec:
			scope.types[name] = v
		case *Constant:
			scope.constants[name] = v
		case *ServiceSpec:
			scope.services[name] = v
		case Scope:
			scope.includes[name] = v
		default:
			panic(fmt.Sprintf(
				"value %v of unknown type %T with name %s", arg, arg, name,
			))
		}
	}
	return scope
}
