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
	"go.uber.org/thriftrw/ast"
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
		Type:  &I32Spec{},
		Value: ConstantInt(42),
	}

	defaultUser := &Constant{
		Name:  "DefaultUser",
		Type:  &StringSpec{},
		Value: ConstantString("anonymous"),
	}

	tests := []struct {
		desc  string
		scope Scope
		name  string

		expected ConstantValue
		typ      TypeSpec
	}{
		{
			"simple constant lookup",
			scope("Version", version),
			"Version",
			ConstReference{Target: version},
			&I32Spec{},
		},
		{
			"included constant lookup",
			scope("shared", scope("DefaultUser", defaultUser)),
			"shared.DefaultUser",
			ConstReference{Target: defaultUser},
			&StringSpec{},
		},
		{
			"enum constant lookup",
			scope("Role", role),
			"Role.Moderator",
			EnumItemReference{
				Enum: role,
				Item: &EnumItem{Name: "Moderator", Value: 2},
			},
			role,
		},
		{
			"included enum constant lookup",
			scope("shared", scope("Role", role)),
			"shared.Role.Disabled",
			EnumItemReference{
				Enum: role,
				Item: &EnumItem{Name: "Disabled", Value: -1},
			},
			role,
		},
	}

	for _, tt := range tests {
		expected, err := tt.expected.Link(defaultScope, tt.typ)
		require.NoError(t, err, "Test constant value must link without errors")

		scope := scopeOrDefault(tt.scope)
		got, err := constantReference(ast.ConstantReference{Name: tt.name}).Link(scope, tt.typ)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, expected, got)
		}
	}
}

func TestCastConstants(t *testing.T) {
	role := &EnumSpec{
		Name: "Role",
		Items: []EnumItem{
			{Name: "Disabled", Value: -1},
			{Name: "Enabled", Value: 1},
			{Name: "Moderator", Value: 2},
		},
	}

	someStruct := &StructSpec{
		Name: "SomeStruct",
		Type: ast.StructType,
		Fields: FieldGroup{
			{
				ID:       1,
				Name:     "someRequiredField",
				Required: true,
				Type:     &I32Spec{},
			},
			{
				ID:   2,
				Name: "someOptionalField",
				Type: &StringSpec{},
			},
			{
				ID:       3,
				Name:     "someFieldWithADefault",
				Required: true,
				Type:     &I64Spec{},
				Default:  ConstantInt(42),
			},
		},
	}

	tests := []struct {
		desc  string
		scope Scope
		typ   TypeSpec
		give  ConstantValue

		want      ConstantValue
		wantError string
	}{
		{
			desc: "ConstantBool",
			typ:  &BoolSpec{},
			give: ConstantBool(true),
			want: ConstantBool(true),
		},
		{
			desc: "ConstantInt: bool (false)",
			typ:  &BoolSpec{},
			give: ConstantInt(0),
			want: ConstantBool(false),
		},
		{
			desc: "ConstantInt: bool (true)",
			typ:  &BoolSpec{},
			give: ConstantInt(1),
			want: ConstantBool(true),
		},
		{
			desc:      "ConstantInt: bool (failure)",
			typ:       &BoolSpec{},
			give:      ConstantInt(2),
			wantError: "the value must be 0 or 1",
		},
		{
			desc: "ConstantInt: i8",
			typ:  &I8Spec{},
			give: ConstantInt(42),
			want: ConstantInt(42),
		},
		{
			desc: "ConstantInt: i16",
			typ:  &I16Spec{},
			give: ConstantInt(42),
			want: ConstantInt(42),
		},
		{
			desc: "ConstantInt: i32",
			typ:  &I32Spec{},
			give: ConstantInt(42),
			want: ConstantInt(42),
		},
		{
			desc: "ConstantInt: i64",
			typ:  &I64Spec{},
			give: ConstantInt(42),
			want: ConstantInt(42),
		},
		{
			desc: "ConstantInt: double",
			typ:  &DoubleSpec{},
			give: ConstantInt(42),
			want: ConstantDouble(42.0),
		},
		{
			desc: "ConstantInt: enum (negative)",
			typ:  role,
			give: ConstantInt(-1),
			want: EnumItemReference{
				Enum: role,
				Item: &role.Items[0], // Disabled
			},
		},
		{
			desc: "ConstantInt: enum (positive)",
			typ:  role,
			give: ConstantInt(2),
			want: EnumItemReference{
				Enum: role,
				Item: &role.Items[2], // Moderator
			},
		},
		{
			desc:      "ConstantInt: enum (failure)",
			typ:       role,
			give:      ConstantInt(3),
			wantError: `3 is not a valid value for enum "Role"`,
		},
		{
			desc:      "ConstantInt: failure",
			typ:       &StringSpec{},
			give:      ConstantInt(1),
			wantError: `cannot cast 1 to "string"`,
		},
		{
			desc: "ConstantString",
			typ:  &StringSpec{},
			give: ConstantString("foo"),
			want: ConstantString("foo"),
		},
		{
			desc: "ConstantDouble",
			typ:  &DoubleSpec{},
			give: ConstantDouble(42.0),
			want: ConstantDouble(42.0),
		},
		{
			desc: "ConstantStruct: all fields",
			typ:  someStruct,
			give: &ConstantStruct{
				Fields: map[string]ConstantValue{
					"someRequiredField":     ConstantInt(100),
					"someOptionalField":     ConstantString("hello"),
					"someFieldWithADefault": ConstantInt(1),
				},
			},
			want: &ConstantStruct{
				Fields: map[string]ConstantValue{
					"someRequiredField":     ConstantInt(100),
					"someOptionalField":     ConstantString("hello"),
					"someFieldWithADefault": ConstantInt(1),
				},
			},
		},
		{
			desc: "ConstantStruct: with default",
			typ:  someStruct,
			give: &ConstantStruct{
				Fields: map[string]ConstantValue{
					"someRequiredField": ConstantInt(100),
				},
			},
			want: &ConstantStruct{
				Fields: map[string]ConstantValue{
					"someRequiredField":     ConstantInt(100),
					"someFieldWithADefault": ConstantInt(42),
				},
			},
		},
		{
			desc:      "ConstantStruct: failure",
			typ:       someStruct,
			give:      &ConstantStruct{Fields: map[string]ConstantValue{}},
			wantError: `"someRequiredField" is a required field`,
		},
		{
			desc: "ConstantStruct: field casting failure",
			typ:  someStruct,
			give: &ConstantStruct{
				Fields: map[string]ConstantValue{
					"someRequiredField": ConstantString("foo"),
				},
			},
			wantError: `failed to cast field "someRequiredField": cannot cast foo to "i32"`,
		},
		{
			desc: "ConstantMap",
			typ:  &MapSpec{KeySpec: &StringSpec{}, ValueSpec: &I32Spec{}},
			give: ConstantMap{
				{Key: ConstantString("hello"), Value: ConstantInt(100)},
				{Key: ConstantString("world"), Value: ConstantInt(200)},
			},
			want: ConstantMap{
				{Key: ConstantString("hello"), Value: ConstantInt(100)},
				{Key: ConstantString("world"), Value: ConstantInt(200)},
			},
		},
		{
			desc: "ConstantMap: struct",
			typ:  someStruct,
			give: ConstantMap{
				{
					Key:   ConstantString("someRequiredField"),
					Value: ConstantInt(100),
				},
			},
			want: &ConstantStruct{
				Fields: map[string]ConstantValue{
					"someRequiredField":     ConstantInt(100),
					"someFieldWithADefault": ConstantInt(42),
				},
			},
		},
		{
			desc: "ConstantMap: struct (invalid key)",
			typ:  someStruct,
			give: ConstantMap{
				{
					Key:   ConstantInt(100),
					Value: ConstantInt(200),
				},
			},
			wantError: "100 is not a string: all keys must be strings",
		},
		{
			desc: "ConstantSet",
			typ:  &SetSpec{ValueSpec: &I32Spec{}},
			give: ConstantSet{ConstantInt(1), ConstantInt(2), ConstantInt(3)},
			want: ConstantSet{ConstantInt(1), ConstantInt(2), ConstantInt(3)},
		},
		{
			desc: "ConstantList",
			typ:  &ListSpec{ValueSpec: &I32Spec{}},
			give: ConstantList{ConstantInt(1), ConstantInt(2), ConstantInt(3)},
			want: ConstantList{ConstantInt(1), ConstantInt(2), ConstantInt(3)},
		},
		{
			desc: "ConstantReference",
			typ:  &I32Spec{},
			give: ConstReference{Target: &Constant{
				Name:  "Version",
				Type:  &I32Spec{},
				Value: ConstantInt(42),
			}},
			want: ConstReference{Target: &Constant{
				Name:  "Version",
				Type:  &I32Spec{},
				Value: ConstantInt(42),
			}},
		},
		{
			desc: "ConstantReference: mismatch",
			typ:  &DoubleSpec{},
			give: ConstReference{Target: &Constant{
				Name:  "Version",
				Type:  &I32Spec{},
				Value: ConstantInt(42),
			}},
			want: ConstantDouble(42.0),
		},
	}

	for _, tt := range tests {
		var err error
		tt.typ, err = tt.typ.Link(defaultScope)
		require.NoError(t, err, "'typ' must link without errors")

		if tt.want != nil {
			tt.want, err = tt.want.Link(defaultScope, tt.typ)
			require.NoError(t, err, "'want' must link without errors")
		}

		scope := scopeOrDefault(tt.scope)
		got, err := tt.give.Link(scope, tt.typ)
		if tt.wantError != "" {
			if assert.Error(t, err, tt.desc) {
				assert.Contains(t, err.Error(), tt.wantError, tt.desc)
			}
		} else {
			if assert.NoError(t, err, tt.desc) {
				assert.Equal(t, tt.want, got, tt.desc)
			}
		}
	}
}

func TestLinkConstantReferenceFailure(t *testing.T) {
	foo := &EnumSpec{
		Name: "Foo",
		Items: []EnumItem{
			{Name: "A", Value: 1},
			{Name: "B", Value: 2},
		},
	}
	tests := []struct {
		desc     string
		scope    Scope
		name     string
		messages []string
		typ      TypeSpec
	}{
		{
			"unknown identifier",
			scope("bar"),
			"foo",
			[]string{`could not resolve reference "foo" in "bar"`},
			&StringSpec{},
		},
		{
			"unknown module",
			scope("bar"),
			"shared.DEFAULT_UUID",
			[]string{
				`could not resolve reference "shared.DEFAULT_UUID" in "bar"`,
				`unknown module "shared"`,
			},
			&StringSpec{},
		},
		{
			"unknown identifier in included module",
			scope("foo", "shared", scope("shared")),
			"shared.DEFAULT_UUID",
			[]string{
				`could not resolve reference "shared.DEFAULT_UUID" in "foo"`,
				`could not resolve reference "DEFAULT_UUID" in "shared"`,
			},
			&StringSpec{},
		},
		{
			"unknown item in enum",
			scope("foo",
				"Foo", foo,
			),
			"Foo.C",
			[]string{
				`could not resolve reference "Foo.C" in "foo"`,
				`enum "Foo" does not have an item named "C"`,
			},
			foo,
		},
	}

	for _, tt := range tests {
		scope := scopeOrDefault(tt.scope)
		_, err := constantReference(ast.ConstantReference{Name: tt.name}).Link(scope, tt.typ)
		if assert.Error(t, err, tt.desc) {
			for _, msg := range tt.messages {
				assert.Contains(t, err.Error(), msg, tt.desc)
			}
		}
	}
}
