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
//      list is empty: should encode into a wire.Value that holds an empty list.
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
		desc     string
		give     *tc.ListOfRequiredPrimitives
		want     wire.Value
		wantList []string
	}{
		{
			desc: "required nil list: wire.Value with empty list",
			give: &tc.ListOfRequiredPrimitives{
				ListOfStrings: nil,
			},
			want: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			wantList: []string{}, // we allocate a new slice if it doesn't exist.
		},
		{
			desc: "required empty list: wire.Value with empty list",
			give: &tc.ListOfRequiredPrimitives{
				ListOfStrings: []string{},
			},
			want: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			wantList: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			w, err := tt.give.ToWire()
			assert.NoError(t, err, "failed to serialize: %encodedValue", tt.give)
			assert.True(t, wire.ValuesAreEqual(tt.want, w))

			// Round trip them all.
			freshestV, b := assertBinaryRoundTrip(t, w, tt.desc)
			assert.True(t, b, "failed round trip")
			assert.True(t, tt.give.Equals(tt.give))
			// Allocate a new instance to serialize from Thrift representation.
			x := new(tc.ListOfRequiredPrimitives)
			if assert.NoError(t, x.FromWire(freshestV), tt.desc) {
				assert.Equal(t, tt.wantList, x.ListOfStrings)
			}
		})
	}
}

func TestListRequiredFromWire(t *testing.T) {
	tests := []struct {
		desc string
		give wire.Value
		want *tc.ListOfRequiredPrimitives
	}{
		{
			desc: "empty list field decodes into an empty slice",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			want: &tc.ListOfRequiredPrimitives{
				ListOfStrings: []string{},
			},
		},
	}
	for _, tt := range tests {
		x := new(tc.ListOfRequiredPrimitives)
		if assert.NoError(t, x.FromWire(tt.give), tt.desc) {
			assert.Equal(t, tt.want, x)
			_, b := assertBinaryRoundTrip(t, tt.give, tt.desc)
			assert.True(t, b, "failed round trip")
		}
	}
}

// Error if required list is missing in the wire representation.
func TestListRequiredFromWireError(t *testing.T) {
	tests := []struct {
		desc      string
		give      wire.Value
		wantError string
	}{
		{
			desc:      "empty list field decodes into empty",
			give:      wire.NewValueStruct(wire.Struct{}),
			wantError: "field ListOfStrings of ListOfRequiredPrimitives is required",
		},
	}
	for _, tt := range tests {
		x := new(tc.ListOfRequiredPrimitives)
		err := x.FromWire(tt.give)
		if assert.Error(t, err, tt.desc) {
			assert.Equal(t, tt.wantError, err.Error())
		}
	}
}

// TestListOptionalToWire tests optional serialization cases.
func TestListOptionalToWire(t *testing.T) {
	tests := []struct {
		desc string
		give wire.Value
		want *tc.ListOfOptionalPrimitives
	}{
		{
			desc: "optional nil list: no list field",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			want: &tc.ListOfOptionalPrimitives{
				ListOfStrings: nil,
			},
		},
		{
			desc: "optional empty list: with list field",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			want: &tc.ListOfOptionalPrimitives{
				ListOfStrings: []string{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			// assertRoundTrip does more than we need as we only need to test wire.Value FromWire response.
			assertRoundTrip(t, tt.want, tt.give, tt.desc)
		})
	}
}

func TestListOptionalFromWire(t *testing.T) {
	tests := []struct {
		desc string
		give wire.Value
		want []string
	}{
		{
			desc: "empty list field wire representation decodes into an empty slice",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueList(wire.ValueListFromSlice(wire.TBinary, []wire.Value{})),
				},
			}}),
			want: []string{},
		},
		{
			desc: "nil list field decodes to nil",
			give: wire.NewValueStruct(wire.Struct{}),
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			x := new(tc.ListOfOptionalPrimitives)
			if assert.NoError(t, x.FromWire(tt.give), tt.desc) {
				assert.Equal(t, tt.want, x.ListOfStrings)
				_, b := assertBinaryRoundTrip(t, tt.give, tt.desc)
				assert.True(t, b, "failed round trip")
			}
		})
	}
}
