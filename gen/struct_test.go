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
	"fmt"
	"testing"

	"github.com/thriftrw/thriftrw-go/gen/testdata"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

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

		assert.NotPanics(t, func() { tt.s.String() })

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

		assert.NotPanics(t, func() { tt.s.String() })

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

		assert.NotPanics(t, func() { tt.s.String() })

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
		o string
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
			"Frame{Size: Size{Height: 200, Width: 100}, TopLeft: Point{X: 1, Y: 2}}",
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

		assert.Equal(t, tt.o, tt.s.String())
	}
}

func TestNestedStructsOptional(t *testing.T) {
	tests := []struct {
		s testdata.User
		v wire.Value
		o string
	}{
		{
			testdata.User{Name: "Foo Bar"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
			}}),
			"User{Contact: <nil>, Name: Foo Bar}",
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
			"User{Contact: ContactInfo{EmailAddress: foo@example.com}, Name: Foo Bar}",
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

		assert.Equal(t, tt.o, tt.s.String())
	}
}

func TestStructStringWithMissingRequiredFields(t *testing.T) {
	tests := []struct {
		i fmt.Stringer
		o string
	}{
		{
			&testdata.Frame{TopLeft: &testdata.Point{}},
			"Frame{Size: <nil>, TopLeft: Point{X: 0, Y: 0}}",
		},
		{
			&testdata.Frame{Size: &testdata.Size{}},
			"Frame{Size: Size{Height: 0, Width: 0}, TopLeft: <nil>}",
		},
		{
			&testdata.Edge{Start: &testdata.Point{X: 1, Y: 2}},
			"Edge{End: <nil>, Start: Point{X: 1, Y: 2}}",
		},
		{
			&testdata.Edge{End: &testdata.Point{X: 3, Y: 4}},
			"Edge{End: Point{X: 3, Y: 4}, Start: <nil>}",
		},
		{
			&testdata.Graph{},
			"Graph{Edges: []}",
		},
		{
			&testdata.Event{},
			"Event{UUID: <nil>}",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.o, tt.i.String())
	}
}

func TestBasicException(t *testing.T) {
	tests := []struct {
		s testdata.DoesNotExistException
		v wire.Value
	}{
		{
			testdata.DoesNotExistException{Key: "foo"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s testdata.DoesNotExistException
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		err := error(&s) // should implement the error interface
		assert.Equal(t, "DoesNotExistException{Key: foo}", err.Error())
	}
}
