// Copyright (c) 2019 Uber Technologies, Inc.
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
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	"go.uber.org/thriftrw/wire"
)

// TestListRequired tests that empty slice and nil allocate an empty list after deserialization via
// FromWire.
func TestListRequired(t *testing.T) {
	tests := []struct {
		desc   string
		p      thriftType // object we run through ToWire and Encode/Decode to assure no changes in round trip.
		v      wire.Value // value we expect after ToWire and Encode/Decode.
		expect []string   // expected after allocating a new p object.
	}{
		{
			"required nil list",
			&tc.ListOfRequiredPrimitives{
				ListOfStrings: nil,
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			[]string{}, // we allocate a new slice if it doesn't exist.
		},
		{
			"required empty list",
			&tc.ListOfRequiredPrimitives{
				ListOfStrings: []string{},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			[]string{},
		},
	}
	for _, tt := range tests {
		// newV should be wire compatible to tt.v after round trip.
		b, newV := assertSerialization(t, tt.p, tt.v, tt.desc)
		assert.True(t, b)
		newP := tt.p.(*tc.ListOfRequiredPrimitives)
		assert.True(t, tt.p.(*tc.ListOfRequiredPrimitives).Equals(newP))
		// Allocate a new instance to serialize from Thrift representation.
		gotX := new(tc.ListOfRequiredPrimitives)
		if assert.NoError(t, gotX.FromWire(newV), tt.desc) { // assure no errors and convert newP.
			assert.Equal(t, tt.expect, gotX.ListOfStrings)
		}
	}
}

// TestListOptional tests that optional lists whether on part of Thrift struct that has a single
// list field or not returns an empty struct for nil and an empty typed slice if an empty slice is
// passed in.
func TestListOptional(t *testing.T) {
	tests := []struct {
		desc string
		p    thriftType
		v    wire.Value
	}{
		{
			"optional nil list",
			&tc.PrimitiveContainers{
				ListOfInts: nil,
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
		},
		{
			"optional nil list only",
			&tc.ListOfOptionalPrimitives{
				ListOfStrings: nil,
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
		},
		{
			"optional empty list",
			&tc.PrimitiveContainers{
				ListOfInts: []int64{},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    2,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TI64, []wire.Value{})),
				},
			}}),
		},
		{
			"optional empty list only",
			&tc.ListOfOptionalPrimitives{
				ListOfStrings: []string{},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
		},
	}
	for _, tt := range tests {
		assertRoundTrip(t, tt.p, tt.v, tt.desc)
	}
}
