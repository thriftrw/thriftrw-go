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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPascalCase(t *testing.T) {
	tests := []struct {
		allCaps bool
		words   []string
		output  string
	}{
		{
			words:  []string{"snake", "case"},
			output: "SnakeCase",
		},
		{
			words:  []string{"get", "ZIP", "code"},
			output: "GetZipCode",
		},
		{
			allCaps: true,
			words:   []string{"get", "ZIP", "code"},
			output:  "GetZIPCode",
		},
		{
			words:  []string{"IP"},
			output: "IP",
		},
		{
			allCaps: true,
			words:   []string{"VIP"},
			output:  "VIP",
		},
		{
			allCaps: false,
			words:   []string{"VIP"},
			output:  "Vip",
		},
		{
			allCaps: false,
			words:   []string{"MyEnum", "FOO"},
			output:  "MyEnumFoo",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.output, pascalCase(tt.allCaps, tt.words...))
	}
}

func TestGoCase(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{"snake_case", "SnakeCase"},
		{"foo__bar", "FooBar"},
		{"get_FooBar", "GetFooBar"},
		{"alreadyCamelCase", "AlreadyCamelCase"},
		{"get500Error", "Get500Error"},
		{"http_request", "HTTPRequest"},
		{"HTTPRequest", "HTTPRequest"},
		{"ALL_CAPS_WITH_UNDERSCORE", "AllCapsWithUnderscore"},
		{"get_user_id", "GetUserID"},
		{"GET_USER_ID", "GetUserID"},
		{"IP", "IP"},
		{"ZIPCode", "ZIPCode"},
		{"VIP", "VIP"}, // not a known abbreviation
	}

	for _, tt := range tests {
		assert.Equal(t, tt.output, goCase(tt.input))
	}
}

func TestConstantName(t *testing.T) {
	tests := []struct {
		give string
		want string
	}{
		{
			give: "VERSION",
			want: "Version",
		},
		{
			give: "API_VERSION",
			want: "APIVersion",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, constantName(tt.give))
	}
}
