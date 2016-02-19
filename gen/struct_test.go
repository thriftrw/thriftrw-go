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

func boolp(x bool) *bool         { return &x }
func bytep(x int8) *int8         { return &x }
func int16p(x int16) *int16      { return &x }
func int32p(x int32) *int32      { return &x }
func int64p(x int64) *int64      { return &x }
func doublep(x float64) *float64 { return &x }
func stringp(x string) *string   { return &x }

func TestPrimitiveRequiredStructWire(t *testing.T) {
	tests := []struct {
		s testdata.PrimitiveRequiredStruct
		v wire.Value
	}{
		{
			testdata.PrimitiveRequiredStruct{
				BoolField:   true,
				ByteField:   1,
				Int16Field:  2,
				Int32Field:  3,
				Int64Field:  4,
				DoubleField: 5.0,
				StringField: "foo",
				BinaryField: []byte("bar"),
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s testdata.PrimitiveRequiredStruct
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestPrimitiveRequiredMissingFields(t *testing.T) {
	tests := []struct {
		v    wire.Value
		msgs []string
	}{}
	// TODO add cases here once we're validating that required fields are
	// present.

	for _, tt := range tests {
		var s testdata.PrimitiveRequiredStruct
		err := s.FromWire(tt.v)
		if assert.Error(t, err) {
			for _, m := range tt.msgs {
				assert.Contains(t, err.Error(), m)
			}
		}
	}
}

func TestPrimitiveOptionalStructWire(t *testing.T) {
	singleFieldStruct := func(id int16, value wire.Value) wire.Value {
		return wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
			{ID: id, Value: value},
		}})
	}

	tests := []struct {
		s testdata.PrimitiveOptionalStruct
		v wire.Value
	}{
		{
			testdata.PrimitiveOptionalStruct{
				BoolField:   boolp(true),
				ByteField:   bytep(1),
				Int16Field:  int16p(2),
				Int32Field:  int32p(3),
				Int64Field:  int64p(4),
				DoubleField: doublep(5.0),
				StringField: stringp("foo"),
				BinaryField: []byte("bar"),
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
		},
		{
			testdata.PrimitiveOptionalStruct{BoolField: boolp(true)},
			singleFieldStruct(1, wire.NewValueBool(true)),
		},
		{
			testdata.PrimitiveOptionalStruct{ByteField: bytep(1)},
			singleFieldStruct(2, wire.NewValueI8(1)),
		},
		{
			testdata.PrimitiveOptionalStruct{Int16Field: int16p(2)},
			singleFieldStruct(3, wire.NewValueI16(2)),
		},
		{
			testdata.PrimitiveOptionalStruct{Int32Field: int32p(3)},
			singleFieldStruct(4, wire.NewValueI32(3)),
		},
		{
			testdata.PrimitiveOptionalStruct{Int64Field: int64p(4)},
			singleFieldStruct(5, wire.NewValueI64(4)),
		},
		{
			testdata.PrimitiveOptionalStruct{DoubleField: doublep(5.0)},
			singleFieldStruct(6, wire.NewValueDouble(5.0)),
		},
		{
			testdata.PrimitiveOptionalStruct{StringField: stringp("foo")},
			singleFieldStruct(7, wire.NewValueString("foo")),
		},
		{
			testdata.PrimitiveOptionalStruct{BinaryField: []byte("bar")},
			singleFieldStruct(8, wire.NewValueBinary([]byte("bar"))),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s testdata.PrimitiveOptionalStruct
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestPrimitiveContainersRequired(t *testing.T) {
	tests := []struct {
		s testdata.PrimitiveContainersRequired
		v wire.Value
	}{
		{
			testdata.PrimitiveContainersRequired{
				ListOfStrings:      []string{"foo", "bar", "baz"},
				SetOfInts:          map[int32]struct{}{1: struct{}{}, 2: struct{}{}},
				MapOfIntsToDoubles: map[int64]float64{1: 2.0, 3: 4.0},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueList(wire.List{
						ValueType: wire.TBinary,
						Size:      3,
						Items: wire.ValueListFromSlice([]wire.Value{
							wire.NewValueString("foo"),
							wire.NewValueString("bar"),
							wire.NewValueString("baz"),
						}),
					}),
				},
				{
					ID: 2,
					Value: wire.NewValueSet(wire.Set{
						ValueType: wire.TI32,
						Size:      2,
						Items: wire.ValueListFromSlice([]wire.Value{
							wire.NewValueI32(1),
							wire.NewValueI32(2),
						}),
					}),
				},
				{
					ID: 3,
					Value: wire.NewValueMap(wire.Map{
						KeyType:   wire.TI64,
						ValueType: wire.TDouble,
						Size:      2,
						Items: wire.MapItemListFromSlice([]wire.MapItem{
							wire.MapItem{
								Key:   wire.NewValueI64(1),
								Value: wire.NewValueDouble(2.0),
							},
							wire.MapItem{
								Key:   wire.NewValueI64(3),
								Value: wire.NewValueDouble(4.0),
							},
						}),
					}),
				},
			}}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s testdata.PrimitiveContainersRequired
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestNestedStructsRequired(t *testing.T) {
	tests := []struct {
		s testdata.Frame
		v wire.Value
	}{
		{
			testdata.Frame{
				TopLeft: &testdata.Point{X: 1, Y: 2},
				Size:    &testdata.Size{Width: 100, Height: 200},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueDouble(1.0)},
						{ID: 2, Value: wire.NewValueDouble(2.0)},
					}}),
				},
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueDouble(100.0)},
						{ID: 2, Value: wire.NewValueDouble(200.0)},
					}}),
				},
			}}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s testdata.Frame
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestNestedStructsOptional(t *testing.T) {
	tests := []struct {
		s testdata.User
		v wire.Value
	}{
		{
			testdata.User{Name: "Foo Bar"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
			}}),
		},
		{
			testdata.User{
				Name:    "Foo Bar",
				Contact: &testdata.ContactInfo{EmailAddress: "foo@example.com"},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
				{ID: 2, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("foo@example.com")},
				}})},
			}}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s testdata.User
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}
