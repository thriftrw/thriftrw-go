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
)

func TestCapitalize(t *testing.T) {
	tests := []struct{ input, output string }{
		{"", ""},
		{"foo", "Foo"},
		{" foo", " foo"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.output, capitalize(tt.input))
	}
}

func TestFileBaseName(t *testing.T) {
	tests := []struct{ input, output string }{
		{"foo.bar", "foo"},
		{"foo/bar.thrift", "bar"},
		{"foo/bar-baz.thrift", "bar-baz"},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.output, fileBaseName(tt.input))
	}
}

func TestSplitInclude(t *testing.T) {
	tests := []struct {
		input       string
		outputLeft  string
		outputRight string
	}{
		{"UUID", "", "UUID"},
		{"common.UUID", "common", "UUID"},
		{"common.types.UUID", "common", "types.UUID"},
	}

	for _, tt := range tests {
		left, right := splitInclude(tt.input)
		assert.Equal(t, tt.outputLeft, left)
		assert.Equal(t, tt.outputRight, right)
	}
}
