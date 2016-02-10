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
	"github.com/thriftrw/thriftrw-go/ast"
)

func TestLinkConstantReference(t *testing.T) {
	role := &EnumSpec{
		Name: "Role",
		Items: []EnumItem{
			{Name: "Disabled", Value: -1},
			{Name: "Enabled", Value: 1},
			{Name: "Moderator", Value: 2},
		},
	}

	version := &Constant{
		Name:  "Version",
		Type:  I32Spec,
		Value: ConstantInt(42),
	}

	defaultUser := &Constant{
		Name:  "DefaultUser",
		Type:  StringSpec,
		Value: ConstantString("anonymous"),
	}

	tests := []struct {
		desc     string
		scope    Scope
		name     string
		expected ConstantValue
	}{
		{
			"simple constant lookup",
			scope("Version", version),
			"Version",
			ConstReference{Target: version},
		},
		{
			"included constant lookup",
			scope("shared", scope("DefaultUser", defaultUser)),
			"shared.DefaultUser",
			ConstReference{Target: defaultUser},
		},
		{
			"enum constant lookup",
			scope("Role", role),
			"Role.Moderator",
			EnumItemReference{
				Enum: role,
				Item: EnumItem{Name: "Moderator", Value: 2},
			},
		},
		{
			"included enum constant lookup",
			scope("shared", scope("Role", role)),
			"shared.Role.Disabled",
			EnumItemReference{
				Enum: role,
				Item: EnumItem{Name: "Disabled", Value: -1},
			},
		},
	}

	for _, tt := range tests {
		expected, err := tt.expected.Link(scope())
		require.NoError(t, err, "Test constant value must link without errors")

		scope := scopeOrDefault(tt.scope)
		got, err := constantReference(ast.ConstantReference{Name: tt.name}).Link(scope)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, expected, got)
		}
	}
}

func TestLinkConstantReferenceFailure(t *testing.T) {
	tests := []struct {
		desc     string
		scope    Scope
		name     string
		messages []string
	}{
		{
			"unknown identifier",
			scope("bar"),
			"foo",
			[]string{`could not resolve reference "foo" in "bar"`},
		},
		{
			"unknown module",
			scope("bar"),
			"shared.DEFAULT_UUID",
			[]string{
				`could not resolve reference "shared.DEFAULT_UUID" in "bar"`,
				`unknown module "shared"`,
			},
		},
		{
			"unknown identifier in included module",
			scope("foo", "shared", scope("shared")),
			"shared.DEFAULT_UUID",
			[]string{
				`could not resolve reference "shared.DEFAULT_UUID" in "foo"`,
				`could not resolve reference "DEFAULT_UUID" in "shared"`,
			},
		},
		{
			"unknown item in enum",
			scope("foo",
				"Foo", &EnumSpec{
					Name: "Foo",
					Items: []EnumItem{
						{Name: "A", Value: 1},
						{Name: "B", Value: 2},
					},
				},
			),
			"Foo.C",
			[]string{
				`could not resolve reference "Foo.C" in "foo"`,
				`enum "Foo" does not have an item named "C"`,
			},
		},
	}

	for _, tt := range tests {
		scope := scopeOrDefault(tt.scope)
		_, err := constantReference(ast.ConstantReference{Name: tt.name}).Link(scope)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}
