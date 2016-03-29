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

	te "github.com/thriftrw/thriftrw-go/gen/testdata/enums"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

func TestValueOfEnumDefault(t *testing.T) {
	tests := []struct {
		e te.EnumDefault
		i int32
	}{
		{te.EnumDefaultFoo, 0},
		{te.EnumDefaultBar, 1},
		{te.EnumDefaultBaz, 2},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestValueOfEnumWithValues(t *testing.T) {
	tests := []struct {
		e te.EnumWithValues
		i int32
	}{
		{te.EnumWithValuesX, 123},
		{te.EnumWithValuesY, 456},
		{te.EnumWithValuesZ, 789},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEnumDefaultWire(t *testing.T) {
	tests := []struct {
		e te.EnumDefault
		v wire.Value
	}{
		{te.EnumDefaultFoo, wire.NewValueI32(0)},
		{te.EnumDefaultBar, wire.NewValueI32(1)},
		{te.EnumDefaultBaz, wire.NewValueI32(2)},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.v, tt.e.ToWire())

		var e te.EnumDefault
		if assert.NoError(t, e.FromWire(tt.v)) {
			assert.Equal(t, tt.e, e)
		}
	}
}

func TestValueOfEnumWithDuplicateValues(t *testing.T) {
	tests := []struct {
		e te.EnumWithDuplicateValues
		i int32
	}{
		{te.EnumWithDuplicateValuesP, 0},
		{te.EnumWithDuplicateValuesQ, -1},
		{te.EnumWithDuplicateValuesR, 0},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEnumWithDuplicateValuesWire(t *testing.T) {
	tests := []struct {
		e te.EnumWithDuplicateValues
		v wire.Value
	}{
		{te.EnumWithDuplicateValuesP, wire.NewValueI32(0)},
		{te.EnumWithDuplicateValuesQ, wire.NewValueI32(-1)},
		{te.EnumWithDuplicateValuesR, wire.NewValueI32(0)},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.v, tt.e.ToWire())

		var e te.EnumWithDuplicateValues
		if assert.NoError(t, e.FromWire(tt.v)) {
			assert.Equal(t, tt.e, e)
		}
	}
}

func TestUnknownEnumValue(t *testing.T) {
	var e te.EnumDefault
	if assert.NoError(t, e.FromWire(wire.NewValueI32(42))) {
		assert.Equal(t, te.EnumDefault(42), e)
	}
}
