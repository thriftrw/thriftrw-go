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

	ts "go.uber.org/thriftrw/gen/testdata/structs"
	td "go.uber.org/thriftrw/gen/testdata/typedefs"
	"go.uber.org/thriftrw/wire"
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
	}
}

func TestTypedefStruct(t *testing.T) {
	tests := []struct {
		x *td.UUID
		v wire.Value
	}{
		{
			(*td.UUID)(&td.I128{1, 2}),
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI64(1)},
				{ID: 2, Value: wire.NewValueI64(2)},
			}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, tt.x, tt.v, "UUID")
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
	}
}
