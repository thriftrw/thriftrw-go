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

func TestTypedefI64(t *testing.T) {
	tests := []struct {
		x test.Timestamp
		v wire.Value
	}{
		{
			test.Timestamp(1),
			wire.NewValueI64(1),
		},
		{
			test.Timestamp(-1),
			wire.NewValueI64(-1),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.x.ToWire()),
			"%v.ToWire() != %v", tt.x, tt.v)

		var x test.Timestamp
		if assert.NoError(t, x.FromWire(tt.v)) {
			assert.Equal(t, tt.x, x)
		}
	}
}

func TestTypedefString(t *testing.T) {
	tests := []struct {
		x test.State
		v wire.Value
	}{
		{
			test.State("hello"),
			wire.NewValueString("hello"),
		},
		{
			test.State("world"),
			wire.NewValueString("world"),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.x.ToWire()),
			"%v.ToWire() != %v", tt.x, tt.v)

		var x test.State
		if assert.NoError(t, x.FromWire(tt.v)) {
			assert.Equal(t, tt.x, x)
		}
	}
}

func TestTypedefBinary(t *testing.T) {
	tests := []struct {
		x test.Pdf
		v wire.Value
	}{
		{
			test.Pdf{1, 2, 3},
			wire.NewValueBinary([]byte{1, 2, 3}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.x.ToWire()),
			"%v.ToWire() != %v", tt.x, tt.v)

		var x test.Pdf
		if assert.NoError(t, x.FromWire(tt.v)) {
			assert.Equal(t, tt.x, x)
		}
	}
}

func TestTypedefStruct(t *testing.T) {
	tests := []struct {
		x *test.UUID
		v wire.Value
	}{
		{
			(*test.UUID)(&test.I128{1, 2}),
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI64(1)},
				{ID: 2, Value: wire.NewValueI64(2)},
			}}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.x.ToWire()),
			"%v.ToWire() != %v", tt.x, tt.v)

		var x test.UUID
		if assert.NoError(t, x.FromWire(tt.v)) {
			assert.Equal(t, tt.x, &x)
		}
	}
}

func TestTypedefContainer(t *testing.T) {
	tests := []struct {
		x test.EventGroup
		v wire.Value
	}{
		{
			test.EventGroup{
				&test.Event{
					UUID: &test.UUID{High: 100, Low: 200},
					Time: (*test.Timestamp)(int64p(42)),
				},
				&test.Event{
					UUID: &test.UUID{High: 0, Low: 42},
					Time: (*test.Timestamp)(int64p(100)),
				},
			},
			wire.NewValueList(wire.List{
				ValueType: wire.TStruct,
				Size:      2,
				Items: wire.ValueListFromSlice([]wire.Value{
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueI64(100)},
							{ID: 2, Value: wire.NewValueI64(200)},
						}})},
						{ID: 2, Value: wire.NewValueI64(42)},
					}}),
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueI64(0)},
							{ID: 2, Value: wire.NewValueI64(42)},
						}})},
						{ID: 2, Value: wire.NewValueI64(100)},
					}}),
				}),
			}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.x.ToWire()),
			"%v.ToWire() != %v", tt.x, tt.v)

		var x test.EventGroup
		if assert.NoError(t, x.FromWire(tt.v)) {
			assert.Equal(t, tt.x, x)
		}
	}
}
