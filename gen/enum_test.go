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

	"github.com/thriftrw/thriftrw-go/gen/testdata"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

func TestValueOfEnumDefault(t *testing.T) {
	tests := []struct {
		e testdata.EnumDefault
		i int32
	}{
		{testdata.EnumDefaultFoo, 0},
		{testdata.EnumDefaultBar, 1},
		{testdata.EnumDefaultBaz, 2},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestValueOfEnumWithValues(t *testing.T) {
	tests := []struct {
		e testdata.EnumWithValues
		i int32
	}{
		{testdata.EnumWithValuesX, 123},
		{testdata.EnumWithValuesY, 456},
		{testdata.EnumWithValuesZ, 789},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEnumDefaultWire(t *testing.T) {
	tests := []struct {
		e testdata.EnumDefault
		v wire.Value
	}{
		{testdata.EnumDefaultFoo, wire.NewValueI32(0)},
		{testdata.EnumDefaultBar, wire.NewValueI32(1)},
		{testdata.EnumDefaultBaz, wire.NewValueI32(2)},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.v, tt.e.ToWire())

		var e testdata.EnumDefault
		if assert.NoError(t, e.FromWire(tt.v)) {
			assert.Equal(t, tt.e, e)
		}
	}
}

func TestValueOfEnumWithDuplicateValues(t *testing.T) {
	tests := []struct {
		e testdata.EnumWithDuplicateValues
		i int32
	}{
		{testdata.EnumWithDuplicateValuesP, 0},
		{testdata.EnumWithDuplicateValuesQ, -1},
		{testdata.EnumWithDuplicateValuesR, 0},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEnumWithDuplicateValuesWire(t *testing.T) {
	tests := []struct {
		e testdata.EnumWithDuplicateValues
		v wire.Value
	}{
		{testdata.EnumWithDuplicateValuesP, wire.NewValueI32(0)},
		{testdata.EnumWithDuplicateValuesQ, wire.NewValueI32(-1)},
		{testdata.EnumWithDuplicateValuesR, wire.NewValueI32(0)},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.v, tt.e.ToWire())

		var e testdata.EnumWithDuplicateValues
		if assert.NoError(t, e.FromWire(tt.v)) {
			assert.Equal(t, tt.e, e)
		}
	}
}

func TestUnknownEnumValue(t *testing.T) {
	var e testdata.EnumDefault
	if assert.NoError(t, e.FromWire(wire.NewValueI32(42))) {
		assert.Equal(t, testdata.EnumDefault(42), e)
	}
}
