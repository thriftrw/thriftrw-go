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

// Testing the following contract.
//
// ToWire
//  Required list field
//      list is nil: should encode into a wire.Value that holds an empty list.
//      list is empty: should also encode into a wire.Value that holds an empty list.
//
//  Optional list field
//      list is nil: should encode into a wire.Value that does not hold the list field.
//      list is empty: should also encode into a wire.Value that holds an empty list.
//
// FromWire
//  Required list field
//      wire.Value has an empty list field: decodes into empty.
//      wire.Value is missing the list field: fails.
//
//  Optional list field
//      wire.Value has an empty list field: decodes into empty.
//      wire.Value is missing the list field: decodes to nil.
//
// We also test that full round trip works as expected with protocol en/decoding.

// TestListRequired tests that empty slice and nil allocate an empty list after deserialization via
// FromWire.
func TestListRequiredToWire(t *testing.T) {
	tests := []struct {
		desc       string
		p          *tc.ListOfRequiredPrimitives
		offTheWire wire.Value
		expect     []string
	}{
		{
			"required nil list: wire.Value with empty list",
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
			"required empty list: wire.Value with empty list",
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
		// required
		w, err := tt.p.ToWire()
		assert.NoError(t, err, "failed to serialize: %encodedValue", tt.p)
		assert.True(t, wire.ValuesAreEqual(tt.offTheWire, w))

		// Round trip them all.
		freshestV, b := assertBinaryRoundTrip(t, w, tt.desc)
		assert.True(t, b, "failed round trip")
		assert.True(t, tt.p.Equals(tt.p))
		// Allocate a new instance to serialize from Thrift representation.
		x := new(tc.ListOfRequiredPrimitives)
		if assert.NoError(t, x.FromWire(freshestV), tt.desc) {
			assert.Equal(t, tt.expect, x.ListOfStrings)
		}
	}
}

func TestListRequiredFromWire(t *testing.T) {
	tests := []struct {
		desc       string
		offTheWire wire.Value
		expected   *tc.ListOfRequiredPrimitives
	}{
		{
			"empty list field decodes into an empty slice",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			&tc.ListOfRequiredPrimitives{
				ListOfStrings: []string{},
			},
		},
	}
	for _, tt := range tests {
		x := new(tc.ListOfRequiredPrimitives)
		if assert.NoError(t, x.FromWire(tt.offTheWire), tt.desc) {
			assert.Equal(t, tt.expected, x)
			_, b := assertBinaryRoundTrip(t, tt.offTheWire, tt.desc)
			assert.True(t, b, "failed round trip")
		}
	}
}

// Error if required list is missing in the wire representation.
func TestListRequiredFromWireError(t *testing.T) {
	tests := []struct {
		desc       string
		offTheWire wire.Value
		wantError  string
	}{
		{
			"empty list field decodes into empty",
			wire.NewValueStruct(wire.Struct{}),
			"field ListOfStrings of ListOfRequiredPrimitives is required",
		},
	}
	for _, tt := range tests {
		x := new(tc.ListOfRequiredPrimitives)
		err := x.FromWire(tt.offTheWire)
		if assert.Error(t, err, tt.desc) {
			assert.Equal(t, tt.wantError, err.Error())
		}
	}
}

// TestListOptionalToWire tests optional serialization cases.
func TestListOptionalToWire(t *testing.T) {
	tests := []struct {
		desc       string
		offTheWire wire.Value
		p          *tc.ListOfOptionalPrimitives
	}{
		{
			"optional nil list: no list field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			&tc.ListOfOptionalPrimitives{
				ListOfStrings: nil,
			},
		},
		{
			"optional empty list: with list field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			&tc.ListOfOptionalPrimitives{
				ListOfStrings: []string{},
			},
		},
	}
	for _, tt := range tests {
		assertRoundTrip(t, tt.p, tt.offTheWire, tt.desc)
	}
}

func TestListOptionalFromWire(t *testing.T) {
	tests := []struct {
		desc       string
		offTheWire wire.Value
		expected   []string
	}{
		{
			"empty list field wire representation decodes into an empty slice",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			[]string{},
		},
		{
			"nil list field decodes to nil",
			wire.NewValueStruct(wire.Struct{}),
			nil,
		},
	}
	for _, tt := range tests {
		x := new(tc.ListOfOptionalPrimitives)
		if assert.NoError(t, x.FromWire(tt.offTheWire), tt.desc) {
			assert.Equal(t, tt.expected, x.ListOfStrings)
			_, b := assertBinaryRoundTrip(t, tt.offTheWire, tt.desc)
			assert.True(t, b, "failed round trip")
		}
	}
}
