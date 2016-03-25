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

	"github.com/thriftrw/thriftrw-go/gen/testdata/test"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

func TestValueOfEnumDefault(t *testing.T) {
	tests := []struct {
		e test.EnumDefault
		i int32
	}{
		{test.EnumDefaultFoo, 0},
		{test.EnumDefaultBar, 1},
		{test.EnumDefaultBaz, 2},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestValueOfEnumWithValues(t *testing.T) {
	tests := []struct {
		e test.EnumWithValues
		i int32
	}{
		{test.EnumWithValuesX, 123},
		{test.EnumWithValuesY, 456},
		{test.EnumWithValuesZ, 789},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEnumDefaultWire(t *testing.T) {
	tests := []struct {
		e test.EnumDefault
		v wire.Value
	}{
		{test.EnumDefaultFoo, wire.NewValueI32(0)},
		{test.EnumDefaultBar, wire.NewValueI32(1)},
		{test.EnumDefaultBaz, wire.NewValueI32(2)},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.v, tt.e.ToWire())

		var e test.EnumDefault
		if assert.NoError(t, e.FromWire(tt.v)) {
			assert.Equal(t, tt.e, e)
		}
	}
}

func TestValueOfEnumWithDuplicateValues(t *testing.T) {
	tests := []struct {
		e test.EnumWithDuplicateValues
		i int32
	}{
		{test.EnumWithDuplicateValuesP, 0},
		{test.EnumWithDuplicateValuesQ, -1},
		{test.EnumWithDuplicateValuesR, 0},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEnumWithDuplicateValuesWire(t *testing.T) {
	tests := []struct {
		e test.EnumWithDuplicateValues
		v wire.Value
	}{
		{test.EnumWithDuplicateValuesP, wire.NewValueI32(0)},
		{test.EnumWithDuplicateValuesQ, wire.NewValueI32(-1)},
		{test.EnumWithDuplicateValuesR, wire.NewValueI32(0)},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.v, tt.e.ToWire())

		var e test.EnumWithDuplicateValues
		if assert.NoError(t, e.FromWire(tt.v)) {
			assert.Equal(t, tt.e, e)
		}
	}
}

func TestUnknownEnumValue(t *testing.T) {
	var e test.EnumDefault
	if assert.NoError(t, e.FromWire(wire.NewValueI32(42))) {
		assert.Equal(t, test.EnumDefault(42), e)
	}
}
