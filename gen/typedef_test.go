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

	tss "go.uber.org/thriftrw/gen/internal/tests/set_to_slice"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	td "go.uber.org/thriftrw/gen/internal/tests/typedefs"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
)

func TestTypedefI64(t *testing.T) {
	tests := []struct {
		x td.Timestamp
		v wire.Value
	}{
		{
			td.Timestamp(1),
			wire.NewValueI64(1),
		},
		{
			td.Timestamp(-1),
			wire.NewValueI64(-1),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "Timestamp")
		assert.True(t, tt.x.Equals(tt.x), "Timestamp equal")

		testRoundTripCombos(t, &tt.x, tt.v, "Timestamp")
		assert.True(t, tt.x.Equals(tt.x), "Timestamp equal")
	}
}

func TestTypedefI64Equals(t *testing.T) {
	tests := []struct {
		x, y td.Timestamp
	}{
		{
			td.Timestamp(1),
			td.Timestamp(-1),
		},
		{
			td.Timestamp(-1),
			td.Timestamp(1),
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "Timestamp equal")
		assert.False(t, tt.x.Equals(tt.y), "Timestamp unequal")
	}
}

func TestTypedefString(t *testing.T) {
	tests := []struct {
		x td.State
		v wire.Value
	}{
		{
			td.State("hello"),
			wire.NewValueString("hello"),
		},
		{
			td.State("world"),
			wire.NewValueString("world"),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "State")
		assert.True(t, tt.x.Equals(tt.x), "State equal")

		testRoundTripCombos(t, &tt.x, tt.v, "State")
		assert.True(t, tt.x.Equals(tt.x), "State equal")
	}
}

func TestTypedefStringEquals(t *testing.T) {
	tests := []struct {
		x, y td.State
	}{
		{
			td.State("hello"),
			td.State("world"),
		},
		{
			td.State("world"),
			td.State("hello"),
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "State equal")
		assert.False(t, tt.x.Equals(tt.y), "State unequal")
	}
}

func TestTypedefBinary(t *testing.T) {
	tests := []struct {
		x td.PDF
		v wire.Value
	}{
		{
			td.PDF{1, 2, 3},
			wire.NewValueBinary([]byte{1, 2, 3}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "PDF")
		assert.True(t, tt.x.Equals(tt.x))

		testRoundTripCombos(t, &tt.x, tt.v, "PDF")
		assert.True(t, tt.x.Equals(tt.x))
	}
}

func TestTypedefBinaryEquals(t *testing.T) {
	tests := []struct {
		x, y td.PDF
	}{
		{
			td.PDF{1, 2, 3},
			td.PDF{1, 3, 5},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x))
		assert.False(t, tt.x.Equals(tt.y))
	}
}

func TestTypedefStruct(t *testing.T) {
	tests := []struct {
		x *td.UUID
		v wire.Value
	}{
		{
			(*td.UUID)(&td.I128{High: 1, Low: 2}),
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI64(1)},
				{ID: 2, Value: wire.NewValueI64(2)},
			}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, tt.x, tt.v, "UUID")
		assert.True(t, tt.x.Equals(tt.x), "UUID equal")

		testRoundTripCombos(t, tt.x, tt.v, "UUID")
		assert.True(t, tt.x.Equals(tt.x), "UUID equal")
	}
}

func TestTypedefStructEquals(t *testing.T) {
	tests := []struct {
		x, y *td.UUID
	}{
		{
			(*td.UUID)(&td.I128{High: 1, Low: 2}),
			(*td.UUID)(&td.I128{High: 3, Low: 1}),
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "UUID equal")
		assert.False(t, tt.x.Equals(tt.y), "UUID unequal")
	}
}

func TestTypedefContainer(t *testing.T) {
	tests := []struct {
		x td.EventGroup
		v wire.Value
	}{
		{
			td.EventGroup{
				&td.Event{
					UUID: &td.UUID{High: 100, Low: 200},
					Time: (*td.Timestamp)(int64p(42)),
				},
				&td.Event{
					UUID: &td.UUID{High: 0, Low: 42},
					Time: (*td.Timestamp)(int64p(100)),
				},
			},
			wire.NewValueList(
				wire.ValueListFromSlice(wire.TStruct, []wire.Value{
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
			),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "EventGroup")
		assert.True(t, tt.x.Equals(tt.x), "EventGroup equal")

		testRoundTripCombos(t, &tt.x, tt.v, "EventGroup")
		assert.True(t, tt.x.Equals(tt.x), "EventGroup equal")
	}
}

func TestTypedefContainerEquals(t *testing.T) {
	tests := []struct {
		x, y td.EventGroup
	}{
		{
			td.EventGroup{
				&td.Event{
					UUID: &td.UUID{High: 100, Low: 200},
					Time: (*td.Timestamp)(int64p(42)),
				},
				&td.Event{
					UUID: &td.UUID{High: 0, Low: 42},
					Time: (*td.Timestamp)(int64p(100)),
				},
			},
			td.EventGroup{
				&td.Event{
					UUID: &td.UUID{High: 100, Low: 200},
					Time: (*td.Timestamp)(int64p(42)),
				},
				&td.Event{
					UUID: &td.UUID{High: 0, Low: 42},
					Time: (*td.Timestamp)(int64p(99)),
				},
			},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "EventGroup equal")
		assert.False(t, tt.x.Equals(tt.y), "EventGroup unequal")
	}
}

func TestUnhashableSetAlias(t *testing.T) {
	tests := []struct {
		x td.FrameGroup
		v wire.Value
	}{
		{
			td.FrameGroup{},
			wire.NewValueSet(
				wire.ValueListFromSlice(wire.TStruct, []wire.Value{}),
			),
		},
		{
			td.FrameGroup{
				&ts.Frame{TopLeft: &ts.Point{X: 1, Y: 2}, Size: &ts.Size{Width: 3, Height: 4}},
				&ts.Frame{TopLeft: &ts.Point{X: 5, Y: 6}, Size: &ts.Size{Width: 7, Height: 8}},
			},
			wire.NewValueSet(
				wire.ValueListFromSlice(wire.TStruct, []wire.Value{
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(1)},
							{ID: 2, Value: wire.NewValueDouble(2)},
						}})},
						{ID: 2, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(3)},
							{ID: 2, Value: wire.NewValueDouble(4)},
						}})},
					}}),
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(5)},
							{ID: 2, Value: wire.NewValueDouble(6)},
						}})},
						{ID: 2, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(7)},
							{ID: 2, Value: wire.NewValueDouble(8)},
						}})},
					}}),
				}),
			),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "FrameGroup")
		assert.True(t, tt.x.Equals(tt.x), "FrameGroup equal")

		testRoundTripCombos(t, &tt.x, tt.v, "FrameGroup")
		assert.True(t, tt.x.Equals(tt.x), "FrameGroup equal")
	}
}

func TestUnhashableSetAliasEquals(t *testing.T) {
	tests := []struct {
		x, y td.FrameGroup
	}{
		{
			td.FrameGroup{},
			td.FrameGroup{
				&ts.Frame{TopLeft: &ts.Point{X: 1, Y: 2}, Size: &ts.Size{Width: 3, Height: 4}},
			},
		},
		{
			td.FrameGroup{
				&ts.Frame{TopLeft: &ts.Point{X: 1, Y: 2}, Size: &ts.Size{Width: 3, Height: 4}},
				&ts.Frame{TopLeft: &ts.Point{X: 5, Y: 6}, Size: &ts.Size{Width: 7, Height: 8}},
			},
			td.FrameGroup{
				&ts.Frame{TopLeft: &ts.Point{X: 1, Y: 2}, Size: &ts.Size{Width: 30, Height: 40}},
				&ts.Frame{TopLeft: &ts.Point{X: 5, Y: 6}, Size: &ts.Size{Width: 7, Height: 8}},
			},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "FrameGroup equal")
		assert.False(t, tt.x.Equals(tt.y), "FrameGroup unequal")
	}
}

func TestUnhashableMapKeyAlias(t *testing.T) {
	tests := []struct {
		x td.PointMap
		v wire.Value
	}{
		{
			td.PointMap{},
			wire.NewValueMap(
				wire.MapItemListFromSlice(wire.TStruct, wire.TStruct, []wire.MapItem{}),
			),
		},
		{
			td.PointMap{
				{
					Key:   &ts.Point{X: 1, Y: 2},
					Value: &ts.Point{X: 3, Y: 4},
				},
				{
					Key:   &ts.Point{X: 5, Y: 6},
					Value: &ts.Point{X: 7, Y: 8},
				},
				{
					Key:   &ts.Point{X: 9, Y: 10},
					Value: &ts.Point{X: 11, Y: 12},
				},
			},
			wire.NewValueMap(
				wire.MapItemListFromSlice(wire.TStruct, wire.TStruct, []wire.MapItem{
					{
						Key: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(1)},
							{ID: 2, Value: wire.NewValueDouble(2)},
						}}),
						Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(3)},
							{ID: 2, Value: wire.NewValueDouble(4)},
						}}),
					},
					{
						Key: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(5)},
							{ID: 2, Value: wire.NewValueDouble(6)},
						}}),
						Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(7)},
							{ID: 2, Value: wire.NewValueDouble(8)},
						}}),
					},
					{
						Key: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(9)},
							{ID: 2, Value: wire.NewValueDouble(10)},
						}}),
						Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueDouble(11)},
							{ID: 2, Value: wire.NewValueDouble(12)},
						}}),
					},
				}),
			),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "PointMap")
		assert.True(t, tt.x.Equals(tt.x), "PointMap equal")

		testRoundTripCombos(t, &tt.x, tt.v, "PointMap")
		assert.True(t, tt.x.Equals(tt.x), "PointMap equal")
	}
}

func TestUnhashableMapKeyAliasEquals(t *testing.T) {
	tests := []struct {
		x, y td.PointMap
	}{
		{
			td.PointMap{},
			td.PointMap{
				{
					Key:   &ts.Point{X: 1, Y: 2},
					Value: &ts.Point{X: 3, Y: 4},
				},
			},
		},
		{
			td.PointMap{
				{
					Key:   &ts.Point{X: 1, Y: 2},
					Value: &ts.Point{X: 3, Y: 4},
				},
				{
					Key:   &ts.Point{X: 5, Y: 6},
					Value: &ts.Point{X: 7, Y: 8},
				},
				{
					Key:   &ts.Point{X: 9, Y: 10},
					Value: &ts.Point{X: 11, Y: 12},
				},
			},
			td.PointMap{
				{
					Key:   &ts.Point{X: 1, Y: 2},
					Value: &ts.Point{X: 3, Y: 4},
				},
				{
					Key:   &ts.Point{X: 5, Y: 6},
					Value: &ts.Point{X: 70, Y: 80},
				},
				{
					Key:   &ts.Point{X: 9, Y: 10},
					Value: &ts.Point{X: 11, Y: 12},
				},
			},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "PointMap equal")
		assert.False(t, tt.x.Equals(tt.y), "PointMap unequal")
	}
}

func TestBinarySet(t *testing.T) {
	tests := []struct {
		x td.BinarySet
		v wire.Value
	}{
		{
			td.BinarySet{},
			wire.NewValueSet(
				wire.ValueListFromSlice(wire.TBinary, []wire.Value{}),
			),
		},
		{
			td.BinarySet{
				{1, 2, 3},
				{4, 5, 6},
			},
			wire.NewValueSet(
				wire.ValueListFromSlice(wire.TBinary, []wire.Value{
					wire.NewValueBinary([]byte{1, 2, 3}),
					wire.NewValueBinary([]byte{4, 5, 6}),
				}),
			),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "BinarySet")
		assert.True(t, tt.x.Equals(tt.x), "BinarySet equal")

		testRoundTripCombos(t, &tt.x, tt.v, "BinarySet")
		assert.True(t, tt.x.Equals(tt.x), "BinarySet equal")
	}
}

func TestBinarySetEquals(t *testing.T) {
	tests := []struct {
		x, y td.BinarySet
	}{
		{
			td.BinarySet{},
			td.BinarySet{
				{1, 2, 3},
			},
		},
		{
			td.BinarySet{
				{1, 2, 3},
				{4, 5, 6},
			},
			td.BinarySet{
				{10, 20, 30},
			},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(tt.x), "BinarySet equal")
		assert.False(t, tt.x.Equals(tt.y), "BinarySet unequal")
	}
}

func TestTypedefAccessors(t *testing.T) {
	t.Run("Timestamp", func(t *testing.T) {
		t.Run("set", func(t *testing.T) {
			timestamp := td.Timestamp(42)
			s := td.Event{Time: &timestamp}
			assert.Equal(t, td.Timestamp(42), s.GetTime())
		})
		t.Run("unset", func(t *testing.T) {
			var s td.Event
			assert.Equal(t, td.Timestamp(0), s.GetTime())
		})
	})

	t.Run("DefaultPrimitiveTypedef", func(t *testing.T) {
		t.Run("set", func(t *testing.T) {
			value := td.State("hi")
			s := td.DefaultPrimitiveTypedef{State: &value}
			assert.Equal(t, td.State("hi"), s.GetState())
		})
		t.Run("unset", func(t *testing.T) {
			var s td.DefaultPrimitiveTypedef
			assert.Equal(t, td.State("hello"), s.GetState())
		})
	})
}

func TestTypedefPtr(t *testing.T) {
	assert.Equal(t, td.State("foo"), *td.State("foo").Ptr())
}

func TestTypedefAnnotatedSetToSlice(t *testing.T) {
	a := tss.StringList{"foo"}
	b := tss.StringList{"foo"}
	c := tss.MyStringList{"foo"}
	d := tss.MyStringList{"foo"}
	e := tss.AnotherStringList{"foo"}
	f := tss.AnotherStringList{"foo"}
	g := tss.StringListList{{"foo"}}
	l := wire.NewValueSet(
		wire.ValueListFromSlice(wire.TBinary, []wire.Value{
			wire.NewValueString("foo"),
		}),
	)
	ll := wire.NewValueSet(
		wire.ValueListFromSlice(wire.TSet, []wire.Value{l}),
	)
	s := "[foo]"

	assertRoundTrip(t, &a, l, "StringList")
	testRoundTripCombos(t, &a, l, "StringList")
	assert.True(t, a.Equals(b))
	assert.Equal(t, s, a.String())

	assertRoundTrip(t, &c, l, "MyStringList")
	testRoundTripCombos(t, &c, l, "MyStringList")
	assert.True(t, c.Equals(d))
	assert.Equal(t, s, c.String())

	assertRoundTrip(t, &e, l, "AnotherStringList")
	testRoundTripCombos(t, &e, l, "AnotherStringList")
	assert.True(t, e.Equals(f))
	assert.Equal(t, s, e.String())

	assertRoundTrip(t, &g, ll, "StringListList")
	testRoundTripCombos(t, &g, ll, "StringListList")
	assert.Equal(t, "[[foo]]", g.String())

	testRoundTripCombos(t, &a, l, "StringList")
	assert.True(t, a.Equals(b))
	assert.Equal(t, s, a.String())

	testRoundTripCombos(t, &c, l, "MyStringList")
	assert.True(t, c.Equals(d))
	assert.Equal(t, s, c.String())

	testRoundTripCombos(t, &e, l, "AnotherStringList")
	assert.True(t, e.Equals(f))
	assert.Equal(t, s, e.String())

	testRoundTripCombos(t, &g, ll, "StringListList")
	assert.Equal(t, "[[foo]]", g.String())
}
