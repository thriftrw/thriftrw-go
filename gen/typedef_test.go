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

	td "github.com/thriftrw/thriftrw-go/gen/testdata/typedefs"
	"github.com/thriftrw/thriftrw-go/wire"
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
		x td.Pdf
		v wire.Value
	}{
		{
			td.Pdf{1, 2, 3},
			wire.NewValueBinary([]byte{1, 2, 3}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, "Pdf")
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
		assertRoundTrip(t, &tt.x, tt.v, "EventGroup")
	}
}
