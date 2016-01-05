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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/uber/thriftrw-go/ast"
)

func TestLookupType(t *testing.T) {
	foo := &StructSpec{
		Name: "Foo",
		Type: ast.StructType,
		Fields: map[string]*FieldSpec{
			"value": {ID: 1, Name: "value", Type: I32Spec},
		},
	}

	tests := []struct {
		desc     string
		module   *Module
		name     string
		expected TypeSpec
	}{
		{
			"simple type lookup",
			&Module{Types: map[string]TypeSpec{"Foo": foo}},
			"Foo",
			foo,
		},
		{
			"included type lookup",
			&Module{
				Includes: map[string]*IncludedModule{
					"bar": {Name: "bar", Module: &Module{
						Types: map[string]TypeSpec{"Foo": foo},
					}},
				},
			},
			"bar.Foo",
			foo,
		},
	}

	for _, tt := range tests {
		expected := mustLink(t, tt.expected, scope())
		spec, err := tt.module.LookupType(tt.name)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, expected, spec, tt.desc)
		}
	}
}

func TestLookupTypeFailure(t *testing.T) {
	tests := []struct {
		desc     string
		module   *Module
		name     string
		messages []string
	}{
		{
			"unknown identifier",
			&Module{Name: "some_module"},
			"foo",
			[]string{`unknown identifier "foo" in "some_module"`},
		},
		{
			"unknown module",
			&Module{Name: "some_module"},
			"shared.UUID",
			[]string{
				`unknown identifier "shared.UUID" in "some_module"`,
				`unknown module "shared"`,
			},
		},
		{
			"unknown identifier in included module",
			&Module{
				Name: "foo",
				Includes: map[string]*IncludedModule{
					"shared": {
						Name:   "shared",
						Module: &Module{Name: "shared"},
					},
				},
			},
			"shared.UUID",
			[]string{
				`unknown identifier "shared.UUID" in "foo"`,
				`unknown identifier "UUID" in "shared"`,
			},
		},
		{
			"unknown identifier in included module of included module",
			&Module{
				Name: "foo",
				Includes: map[string]*IncludedModule{
					"bar": {
						Name: "bar",
						Module: &Module{
							Name: "bar",
							Includes: map[string]*IncludedModule{
								"baz": {
									Name:   "baz",
									Module: &Module{Name: "baz"},
								},
							},
						},
					},
				},
			},
			"bar.baz.Qux",
			[]string{
				`unknown identifier "bar.baz.Qux" in "foo"`,
				`unknown identifier "baz.Qux" in "bar"`,
				`unknown identifier "Qux" in "baz"`,
			},
		},
	}

	for _, tt := range tests {
		_, err := tt.module.LookupType(tt.name)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}

func TestLookupConstant(t *testing.T) {
	role := &EnumSpec{
		Name: "Role",
		Items: []EnumItem{
			{Name: "Disabled", Value: -1},
			{Name: "Enabled", Value: 1},
			{Name: "Moderator", Value: 2},
		},
	}

	tests := []struct {
		desc     string
		module   *Module
		name     string
		expected ast.ConstantValue
	}{
		{
			"simple constant lookup",
			&Module{
				Name: "foo",
				Constants: map[string]*Constant{
					"Version": {
						Name:  "Version",
						Type:  I32Spec,
						Value: ast.ConstantInteger(42),
					},
				},
			},
			"Version",
			ast.ConstantInteger(42),
		},
		{
			"included constant lookup",
			&Module{
				Name: "foo",
				Includes: map[string]*IncludedModule{
					"shared": {Name: "shared", Module: &Module{
						Name: "shared",
						Constants: map[string]*Constant{
							"DefaultUser": {
								Name:  "DefaultUser",
								Type:  StringSpec,
								Value: ast.ConstantString("anonymous"),
							},
						},
					}},
				},
			},
			"shared.DefaultUser",
			ast.ConstantString("anonymous"),
		},
		{
			"enum constant lookup",
			&Module{Name: "foo", Types: map[string]TypeSpec{"Role": role}},
			"Role.Moderator",
			ast.ConstantInteger(2),
		},
		{
			"included enum constant lookup",
			&Module{
				Name: "foo",
				Includes: map[string]*IncludedModule{
					"shared": {Name: "shared", Module: &Module{
						Name:  "shared",
						Types: map[string]TypeSpec{"Role": role},
					}},
				},
			},
			"shared.Role.Disabled",
			ast.ConstantInteger(-1),
		},
	}

	for _, tt := range tests {
		spec, err := tt.module.LookupConstant(tt.name)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, tt.expected, spec, tt.desc)
		}
	}
}

func TestLookupConstantFailure(t *testing.T) {
	tests := []struct {
		desc     string
		module   *Module
		name     string
		messages []string
	}{
		{
			"unknown identifier",
			&Module{Name: "bar"},
			"foo",
			[]string{`unknown identifier "foo" in "bar"`},
		},
		{
			"unknown module",
			&Module{Name: "bar"},
			"shared.DEFAULT_UUID",
			[]string{
				`unknown identifier "shared.DEFAULT_UUID" in "bar"`,
				`unknown module "shared"`,
			},
		},
		{
			"unknown identifier in included module",
			&Module{
				Name: "foo",
				Includes: map[string]*IncludedModule{
					"shared": {
						Name:   "shared",
						Module: &Module{Name: "shared"},
					},
				},
			},
			"shared.DEFAULT_UUID",
			[]string{
				`unknown identifier "shared.DEFAULT_UUID" in "foo"`,
				`unknown identifier "DEFAULT_UUID" in "shared"`,
			},
		},
		{
			"unknown item in enum",
			&Module{
				Name: "foo",
				Types: map[string]TypeSpec{
					"Foo": &EnumSpec{
						Name: "Foo",
						Items: []EnumItem{
							{Name: "A", Value: 1},
							{Name: "B", Value: 2},
						},
					},
				},
			},
			"Foo.C",
			[]string{
				`unknown identifier "Foo.C" in "foo"`,
				`enum "Foo" does not have an item named "C"`,
			},
		},
	}

	for _, tt := range tests {
		_, err := tt.module.LookupConstant(tt.name)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}

func TestLookupService(t *testing.T) {
	someService := &ServiceSpec{
		Name: "SomeService",
		Functions: map[string]*FunctionSpec{
			"someMethod": {
				Name: "someMethod",
				ArgsSpec: map[string]*FieldSpec{
					"someArg": {
						ID:       1,
						Name:     "someArg",
						Type:     I32Spec,
						Required: true,
					},
				},
			},
		},
	}

	tests := []struct {
		desc     string
		module   *Module
		name     string
		expected *ServiceSpec
	}{
		{
			"simple service lookup",
			&Module{
				Name:     "foo",
				Services: map[string]*ServiceSpec{"SomeService": someService},
			},
			"SomeService",
			someService,
		},
		{
			"included service lookup",
			&Module{
				Name: "foo",
				Includes: map[string]*IncludedModule{
					"bar": {
						Name: "bar",
						Module: &Module{
							Name: "bar",
							Services: map[string]*ServiceSpec{
								"SomeService": someService,
							},
						},
					},
				},
			},
			"bar.SomeService",
			someService,
		},
	}

	for _, tt := range tests {
		// invalid test if the sample doesn't link with an empty scope
		require.NoError(
			t,
			tt.expected.Link(scope()),
			"test service must link with an empty scope",
		)

		spec, err := tt.module.LookupService(tt.name)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, tt.expected, spec, tt.desc)
		}
	}
}

func TestLookupServiceFailure(t *testing.T) {
	tests := []struct {
		desc     string
		module   *Module
		name     string
		messages []string
	}{
		{
			"unknown service",
			&Module{Name: "foo"},
			"Foo",
			[]string{`unknown identifier "Foo" in "foo"`},
		},
		{
			"unknown module",
			&Module{Name: "bar"},
			"foo.Bar",
			[]string{
				`unknown identifier "foo.Bar" in "bar"`,
				`unknown module "foo"`,
			},
		},
		{
			"unknown service in included module",
			&Module{
				Name: "bar",
				Includes: map[string]*IncludedModule{
					"foo": {Name: "foo", Module: &Module{Name: "foo"}},
				},
			},
			"foo.Bar",
			[]string{
				`unknown identifier "foo.Bar" in "bar"`,
				`unknown identifier "Bar" in "foo"`,
			},
		},
	}

	for _, tt := range tests {
		_, err := tt.module.LookupService(tt.name)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}
