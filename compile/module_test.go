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
			&Module{},
			"foo",
			[]string{`unknown identifier "foo"`},
		},
		{
			"unknown module",
			&Module{},
			"shared.UUID",
			[]string{
				`unknown identifier "shared.UUID"`,
				`unknown module "shared"`,
			},
		},
		{
			"unknown identifier in included module",
			&Module{Includes: map[string]*IncludedModule{
				"shared": {
					Name:   "shared",
					Module: &Module{},
				},
			}},
			"shared.UUID",
			[]string{
				`unknown identifier "shared.UUID"`,
				`unknown identifier "UUID"`,
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
				Includes: map[string]*IncludedModule{
					"shared": {Name: "shared", Module: &Module{
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
			&Module{Types: map[string]TypeSpec{"Role": role}},
			"Role.Moderator",
			ast.ConstantInteger(2),
		},
		{
			"included enum constant lookup",
			&Module{
				Includes: map[string]*IncludedModule{
					"shared": {Name: "shared", Module: &Module{
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
			&Module{},
			"foo",
			[]string{`unknown identifier "foo"`},
		},
		{
			"unknown module",
			&Module{},
			"shared.DEFAULT_UUID",
			[]string{
				`unknown identifier "shared.DEFAULT_UUID"`,
				`unknown module "shared"`,
			},
		},
		{
			"unknown identifier in included module",
			&Module{Includes: map[string]*IncludedModule{
				"shared": {
					Name:   "shared",
					Module: &Module{},
				},
			}},
			"shared.DEFAULT_UUID",
			[]string{
				`unknown identifier "shared.DEFAULT_UUID"`,
				`unknown identifier "DEFAULT_UUID"`,
			},
		},
		{
			"unknown item in enum",
			&Module{Types: map[string]TypeSpec{
				"Foo": &EnumSpec{
					Name: "Foo",
					Items: []EnumItem{
						{Name: "A", Value: 1},
						{Name: "B", Value: 2},
					},
				},
			}},
			"Foo.C",
			[]string{
				`unknown identifier "Foo.C"`,
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
