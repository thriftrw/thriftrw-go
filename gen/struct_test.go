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
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	tx "go.uber.org/thriftrw/gen/internal/tests/exceptions"
	"go.uber.org/thriftrw/gen/internal/tests/noerror"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	td "go.uber.org/thriftrw/gen/internal/tests/typedefs"
	tu "go.uber.org/thriftrw/gen/internal/tests/unions"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
)

func TestStructRoundTripAndString(t *testing.T) {
	tests := []struct {
		desc string
		x    thriftType
		v    wire.Value
		s    string
	}{
		{
			"PrimitiveRequiredStruct",
			&ts.PrimitiveRequiredStruct{
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
			"",
		},
		{
			"PrimitiveOptionalStruct: all fields",
			&ts.PrimitiveOptionalStruct{
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
			"",
		},
		{
			"PrimitiveOptionalStruct: bool",
			&ts.PrimitiveOptionalStruct{BoolField: boolp(true)},
			singleFieldStruct(1, wire.NewValueBool(true)),
			"",
		},
		{
			"PrimitiveOptionalStruct: byte",
			&ts.PrimitiveOptionalStruct{ByteField: bytep(1)},
			singleFieldStruct(2, wire.NewValueI8(1)),
			"",
		},
		{
			"PrimitiveOptionalStruct: int16",
			&ts.PrimitiveOptionalStruct{Int16Field: int16p(2)},
			singleFieldStruct(3, wire.NewValueI16(2)),
			"",
		},
		{
			"PrimitiveOptionalStruct: int32",
			&ts.PrimitiveOptionalStruct{Int32Field: int32p(3)},
			singleFieldStruct(4, wire.NewValueI32(3)),
			"",
		},
		{
			"PrimitiveOptionalStruct: int64",
			&ts.PrimitiveOptionalStruct{Int64Field: int64p(4)},
			singleFieldStruct(5, wire.NewValueI64(4)),
			"",
		},
		{
			"PrimitiveOptionalStruct: double",
			&ts.PrimitiveOptionalStruct{DoubleField: doublep(5.0)},
			singleFieldStruct(6, wire.NewValueDouble(5.0)),
			"",
		},
		{
			"PrimitiveOptionalStruct: string",
			&ts.PrimitiveOptionalStruct{StringField: stringp("foo")},
			singleFieldStruct(7, wire.NewValueString("foo")),
			"",
		},
		{
			"PrimitiveOptionalStruct: binary",
			&ts.PrimitiveOptionalStruct{BinaryField: []byte("bar")},
			singleFieldStruct(8, wire.NewValueBinary([]byte("bar"))),
			"",
		},
		{
			"PrimitiveContainersRequired",
			&tc.PrimitiveContainersRequired{
				ListOfStrings:      []string{"foo", "bar", "baz"},
				SetOfInts:          map[int32]struct{}{1: {}, 2: {}},
				MapOfIntsToDoubles: map[int64]float64{1: 2.0, 3: 4.0},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TBinary, []wire.Value([]wire.Value{
							wire.NewValueString("foo"),
							wire.NewValueString("bar"),
							wire.NewValueString("baz"),
						})),
					),
				},
				{
					ID: 2,
					Value: wire.NewValueSet(
						wire.ValueListFromSlice(wire.TI32, []wire.Value{
							wire.NewValueI32(1),
							wire.NewValueI32(2),
						}),
					),
				},
				{
					ID: 3,
					Value: wire.NewValueMap(
						wire.MapItemListFromSlice(wire.TI64, wire.TDouble, []wire.MapItem{
							{
								Key:   wire.NewValueI64(1),
								Value: wire.NewValueDouble(2.0),
							},
							{
								Key:   wire.NewValueI64(3),
								Value: wire.NewValueDouble(4.0),
							},
						}),
					),
				},
			}}),
			"",
		},
		{
			"Frame",
			&ts.Frame{
				TopLeft: &ts.Point{X: 1, Y: 2},
				Size:    &ts.Size{Width: 100, Height: 200},
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
			"Frame{TopLeft: Point{X: 1, Y: 2}, Size: Size{Width: 100, Height: 200}}",
		},
		{
			"User: optional field missing",
			&ts.User{Name: "Foo Bar"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
			}}),
			"User{Name: Foo Bar}",
		},
		{
			"User: optional field present",
			&ts.User{
				Name:    "Foo Bar",
				Contact: &ts.ContactInfo{EmailAddress: "foo@example.com"},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
				{ID: 2, Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("foo@example.com")},
				}})},
			}}),
			"User{Name: Foo Bar, Contact: ContactInfo{EmailAddress: foo@example.com}}",
		},
		{
			"List: self-referential struct",
			&ts.List{Value: 1, Tail: &ts.List{Value: 2}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(1)},
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueI32(2)},
					}}),
				},
			}}),
			"Node{Value: 1, Tail: Node{Value: 2}}",
		},
		{
			"Document: PDF",
			&tu.Document{Pdf: []byte{1, 2, 3}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
			}}),
			"Document{Pdf: [1 2 3]}",
		},
		{
			"Document: PlainText",
			&tu.Document{PlainText: stringp("hello")},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueString("hello")},
			}}),
			"Document{PlainText: hello}",
		},
		{
			"ArbitraryValue: bool",
			&tu.ArbitraryValue{BoolValue: boolp(true)},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
			}}),
			"ArbitraryValue{BoolValue: true}",
		},
		{
			"ArbitraryValue: i64",
			&tu.ArbitraryValue{Int64Value: int64p(42)},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI64(42)},
			}}),
			"ArbitraryValue{Int64Value: 42}",
		},
		{
			"ArbitraryValue: string",
			&tu.ArbitraryValue{StringValue: stringp("hello")},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			"ArbitraryValue{StringValue: hello}",
		},
		{
			"ArbitraryValue: list",
			&tu.ArbitraryValue{ListValue: []*tu.ArbitraryValue{
				{BoolValue: boolp(true)},
				{Int64Value: int64p(42)},
				{StringValue: stringp("hello")},
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 4, Value: wire.NewValueList(
					wire.ValueListFromSlice(wire.TStruct, []wire.Value{
						wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 1, Value: wire.NewValueBool(true)},
						}}),
						wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 2, Value: wire.NewValueI64(42)},
						}}),
						wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{ID: 3, Value: wire.NewValueString("hello")},
						}}),
					}),
				)},
			}}),
			"ArbitraryValue{ListValue: [ArbitraryValue{BoolValue: true} ArbitraryValue{Int64Value: 42} ArbitraryValue{StringValue: hello}]}",
		},
		{
			"ArbitraryValue: map",
			&tu.ArbitraryValue{MapValue: map[string]*tu.ArbitraryValue{
				"bool":   {BoolValue: boolp(true)},
				"int64":  {Int64Value: int64p(42)},
				"string": {StringValue: stringp("hello")},
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 5, Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TBinary, wire.TStruct, []wire.MapItem{
						{
							Key: wire.NewValueString("bool"),
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueBool(true)},
							}}),
						},
						{
							Key: wire.NewValueString("int64"),
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 2, Value: wire.NewValueI64(42)},
							}}),
						},
						{
							Key: wire.NewValueString("string"),
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 3, Value: wire.NewValueString("hello")},
							}}),
						},
					}),
				)},
			}}),
			"",
		},
		{
			"EmptyStruct",
			&ts.EmptyStruct{},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			"",
		},
		{
			"EmptyUnion",
			&tu.EmptyUnion{},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			"",
		},
		{
			"EmptyException",
			&tx.EmptyException{},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			"",
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, tt.x, tt.v, tt.desc)
		if tt.s != "" {
			assert.Equal(t, tt.s, tt.x.String(), "ToString: %v", tt.desc)
		} else {
			assert.NotPanics(t, func() { _ = tt.x.String() }, "ToString: %v", tt.desc)
		}
	}
}

func TestPrimitiveRequiredMissingFields(t *testing.T) {
	tests := []struct {
		desc      string
		v         wire.Value
		wantError string
	}{
		{
			"bool",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field BoolField of PrimitiveRequiredStruct is required",
		},
		{
			"byte",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field ByteField of PrimitiveRequiredStruct is required",
		},
		{
			"int16",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field Int16Field of PrimitiveRequiredStruct is required",
		},
		{
			"int32",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field Int32Field of PrimitiveRequiredStruct is required",
		},
		{
			"int64",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field Int64Field of PrimitiveRequiredStruct is required",
		},
		{
			"double",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 7, Value: wire.NewValueString("foo")},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field DoubleField of PrimitiveRequiredStruct is required",
		},
		{
			"string",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 8, Value: wire.NewValueBinary([]byte("bar"))},
			}}),
			"field StringField of PrimitiveRequiredStruct is required",
		},
		{
			"binary",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI8(1)},
				{ID: 3, Value: wire.NewValueI16(2)},
				{ID: 4, Value: wire.NewValueI32(3)},
				{ID: 5, Value: wire.NewValueI64(4)},
				{ID: 6, Value: wire.NewValueDouble(5.0)},
				{ID: 7, Value: wire.NewValueString("foo")},
			}}),
			"field BinaryField of PrimitiveRequiredStruct is required",
		},
	}

	for _, tt := range tests {
		var s ts.PrimitiveRequiredStruct
		err := s.FromWire(tt.v)
		if assert.Error(t, err, tt.desc) {
			assert.Contains(t, err.Error(), tt.wantError, tt.desc)
		}
	}
}

func TestStructStringWithNil(t *testing.T) {
	var f *ts.Frame
	assert.Equal(t, "<nil>", f.String())
}

func TestStructStringWithMissingRequiredFields(t *testing.T) {
	tests := []struct {
		i fmt.Stringer
		o string
	}{
		{
			&ts.Frame{TopLeft: &ts.Point{}},
			"Frame{TopLeft: Point{X: 0, Y: 0}, Size: <nil>}",
		},
		{
			&ts.Frame{Size: &ts.Size{}},
			"Frame{TopLeft: <nil>, Size: Size{Width: 0, Height: 0}}",
		},
		{
			&ts.Edge{StartPoint: &ts.Point{X: 1, Y: 2}},
			"Edge{StartPoint: Point{X: 1, Y: 2}, EndPoint: <nil>}",
		},
		{
			&ts.Edge{EndPoint: &ts.Point{X: 3, Y: 4}},
			"Edge{StartPoint: <nil>, EndPoint: Point{X: 3, Y: 4}}",
		},
		{
			&ts.Graph{},
			"Graph{Edges: []}",
		},
		{
			&td.Event{},
			"Event{UUID: <nil>}",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.o, tt.i.String())
	}
}

func TestBasicException(t *testing.T) {
	tests := []struct {
		s tx.DoesNotExistException
		v wire.Value
	}{
		{
			tx.DoesNotExistException{Key: "foo"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.s, tt.v, "DoesNotExistException")

		err := error(&tt.s) // should implement the error interface
		assert.Equal(t, "DoesNotExistException{Key: foo}", err.Error())
	}
}

func TestStructFromWireUnrecognizedField(t *testing.T) {
	tests := []struct {
		desc string
		give wire.Value

		want      ts.ContactInfo
		wantError string
	}{
		{
			desc: "unknown field",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
				{ID: 2, Value: wire.NewValueI32(42)},
			}}),
			want: ts.ContactInfo{EmailAddress: "foo"},
		},
		{
			desc: "only unknown field",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueString("bar")},
			}}),
			wantError: "field EmailAddress of ContactInfo is required",
		},
		{
			desc: "wrong type for recognized field",
			give: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(42)},
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
			want: ts.ContactInfo{EmailAddress: "foo"},
		},
	}

	for _, tt := range tests {
		var o ts.ContactInfo
		err := o.FromWire(tt.give)
		if tt.wantError != "" {
			if assert.Error(t, err, tt.desc) {
				assert.Contains(t, err.Error(), tt.wantError)
			}
		} else {
			if assert.NoError(t, err, tt.desc) {
				assert.Equal(t, tt.want, o)
			}
		}
	}
}

func TestUnionFromWireInconsistencies(t *testing.T) {
	tests := []struct {
		desc    string
		input   wire.Value
		success *tu.Document
		failure string
	}{
		{
			desc: "multiple recognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
				{ID: 2, Value: wire.NewValueString("hello")},
			}}),
			failure: "should have exactly one field: got 2 fields",
		},
		{
			desc: "recognized and unrecognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			success: &tu.Document{Pdf: []byte{1, 2, 3}},
		},
		{
			desc: "recognized field duplicates",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
				{ID: 1, Value: wire.NewValueBinary([]byte{4, 5, 6})},
			}}),
			success: &tu.Document{Pdf: []byte{4, 5, 6}},
		},
		{
			desc: "only unrecognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI32(42)}, // also a type mismatch
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			failure: "should have exactly one field: got 0 fields",
		},
		{
			desc:    "no fields",
			input:   wire.NewValueStruct(wire.Struct{}),
			failure: "should have exactly one field: got 0 fields",
		},
	}

	for _, tt := range tests {
		var o tu.Document
		err := o.FromWire(tt.input)
		if tt.success != nil {
			if assert.NoError(t, err, tt.desc) {
				assert.Equal(t, tt.success, &o, tt.desc)
			}
		} else {
			if assert.Error(t, err, tt.desc) {
				assert.Contains(t, err.Error(), tt.failure, tt.desc)
			}
		}
	}
}

func TestStructWithDefaults(t *testing.T) {
	enumDefaultFoo := te.EnumDefaultFoo
	enumDefaultBar := te.EnumDefaultBar
	enumDefaultBaz := te.EnumDefaultBaz

	tests := []struct {
		give     *ts.DefaultsStruct
		giveWire wire.Value

		wantToWire   wire.Value
		wantFromWire *ts.DefaultsStruct
	}{
		{
			give:     &ts.DefaultsStruct{},
			giveWire: wire.NewValueStruct(wire.Struct{}),

			wantToWire: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(100)},
				{ID: 2, Value: wire.NewValueI32(200)},
				{ID: 3, Value: wire.NewValueI32(1)},
				{ID: 4, Value: wire.NewValueI32(2)},
				{
					ID: 5,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TBinary, []wire.Value{
							wire.NewValueString("hello"),
							wire.NewValueString("world"),
						}),
					),
				},
				{
					ID: 6,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TDouble, []wire.Value{
							wire.NewValueDouble(1.0),
							wire.NewValueDouble(2.0),
							wire.NewValueDouble(3.0),
						}),
					),
				},
				{
					ID: 7,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
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
				{
					ID: 8,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
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
								{ID: 1, Value: wire.NewValueDouble(3.0)},
								{ID: 2, Value: wire.NewValueDouble(4.0)},
							}}),
						},
					}}),
				},
			}}),
			wantFromWire: &ts.DefaultsStruct{
				RequiredPrimitive: int32p(100),
				OptionalPrimitive: int32p(200),
				RequiredEnum:      &enumDefaultBar,
				OptionalEnum:      &enumDefaultBaz,
				RequiredList:      []string{"hello", "world"},
				OptionalList:      []float64{1.0, 2.0, 3.0},
				RequiredStruct: &ts.Frame{
					TopLeft: &ts.Point{X: 1.0, Y: 2.0},
					Size:    &ts.Size{Width: 100.0, Height: 200.0},
				},
				OptionalStruct: &ts.Edge{
					StartPoint: &ts.Point{X: 1.0, Y: 2.0},
					EndPoint:   &ts.Point{X: 3.0, Y: 4.0},
				},
			},
		},
		{
			give: &ts.DefaultsStruct{
				RequiredPrimitive: int32p(0),
				OptionalEnum:      &enumDefaultFoo,
				RequiredList:      []string{},
			},
			giveWire: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(0)},
				{ID: 4, Value: wire.NewValueI32(0)},
				{
					ID: 5,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TBinary, []wire.Value{}),
					),
				},
			}}),

			wantToWire: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(0)},
				{ID: 2, Value: wire.NewValueI32(200)},
				{ID: 3, Value: wire.NewValueI32(1)},
				{ID: 4, Value: wire.NewValueI32(0)},
				{
					ID: 5,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TBinary, []wire.Value{}),
					),
				},
				{
					ID: 6,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TDouble, []wire.Value{
							wire.NewValueDouble(1.0),
							wire.NewValueDouble(2.0),
							wire.NewValueDouble(3.0),
						}),
					),
				},
				{
					ID: 7,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
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
				{
					ID: 8,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
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
								{ID: 1, Value: wire.NewValueDouble(3.0)},
								{ID: 2, Value: wire.NewValueDouble(4.0)},
							}}),
						},
					}}),
				},
			}}),
			wantFromWire: &ts.DefaultsStruct{
				RequiredPrimitive: int32p(0),
				OptionalPrimitive: int32p(200),
				RequiredEnum:      &enumDefaultBar,
				OptionalEnum:      &enumDefaultFoo,
				RequiredList:      []string{},
				OptionalList:      []float64{1.0, 2.0, 3.0},
				RequiredStruct: &ts.Frame{
					TopLeft: &ts.Point{X: 1.0, Y: 2.0},
					Size:    &ts.Size{Width: 100.0, Height: 200.0},
				},
				OptionalStruct: &ts.Edge{
					StartPoint: &ts.Point{X: 1.0, Y: 2.0},
					EndPoint:   &ts.Point{X: 3.0, Y: 4.0},
				},
			},
		},
	}

	for _, tt := range tests {
		if gotWire, err := tt.give.ToWire(); assert.NoError(
			t, err, "%v.ToWire() failed", tt.give) {
			assert.True(
				t, wire.ValuesAreEqual(tt.wantToWire, gotWire),
				"%v.ToWire() != %v", tt.give, tt.wantToWire)
		}

		var gotFromWire ts.DefaultsStruct
		if err := gotFromWire.FromWire(tt.giveWire); assert.NoError(t, err) {
			assert.Equal(t, tt.wantFromWire, &gotFromWire)
		}
	}
}

func TestStructJSON(t *testing.T) {
	tests := []struct {
		v interface{}
		j string
	}{
		{&ts.Point{X: 0, Y: 0}, `{"x":0,"y":0}`},
		{&ts.Point{X: 1, Y: 2}, `{"x":1,"y":2}`},
		{
			&ts.Edge{
				StartPoint: &ts.Point{X: 1, Y: 2},
				EndPoint:   &ts.Point{X: 3, Y: 4},
			},
			`{"startPoint":{"x":1,"y":2},"endPoint":{"x":3,"y":4}}`,
		},
		{
			&ts.Edge{StartPoint: &ts.Point{X: 1, Y: 1}},
			`{"startPoint":{"x":1,"y":1},"endPoint":null}`,
		},
		{
			&ts.Edge{
				StartPoint: &ts.Point{X: 1, Y: 1},
				EndPoint:   &ts.Point{X: 1, Y: 1},
			},
			`{"startPoint":{"x":1,"y":1},"endPoint":{"x":1,"y":1}}`,
		},
		{&ts.User{Name: ""}, `{"name":""}`},
		{&ts.User{Name: "foo"}, `{"name":"foo"}`},
		{
			&ts.User{
				Name:    "foo",
				Contact: &ts.ContactInfo{EmailAddress: "bar@example.com"},
			},
			`{"name":"foo","contact":{"emailAddress":"bar@example.com"}}`,
		},
		{
			&ts.User{
				Name:    "foo",
				Contact: &ts.ContactInfo{EmailAddress: ""},
			},
			`{"name":"foo","contact":{"emailAddress":""}}`,
		},
		{&tu.EmptyUnion{}, "{}"},
		{&tu.Document{Pdf: td.PDF("hello")}, `{"pdf":"aGVsbG8="}`},
		{
			&tu.Document{PlainText: stringp("hello")},
			`{"plainText":"hello"}`,
		},
		{&tu.Document{PlainText: stringp("")}, `{"plainText":""}`},
		{
			&tu.ArbitraryValue{BoolValue: boolp(true)},
			`{"boolValue":true}`,
		},
		{
			&tu.ArbitraryValue{BoolValue: boolp(false)},
			`{"boolValue":false}`,
		},
		{
			&tu.ArbitraryValue{Int64Value: int64p(42)},
			`{"int64Value":42}`,
		},
		{
			&tu.ArbitraryValue{Int64Value: int64p(0)},
			`{"int64Value":0}`,
		},
		{
			&tu.ArbitraryValue{StringValue: stringp("foo")},
			`{"stringValue":"foo"}`,
		},
		{
			&tu.ArbitraryValue{StringValue: stringp("")},
			`{"stringValue":""}`,
		},
		{
			&tu.ArbitraryValue{ListValue: []*tu.ArbitraryValue{
				{BoolValue: boolp(true)},
				{Int64Value: int64p(42)},
				{StringValue: stringp("foo")},
			}},
			`{"listValue":[` +
				`{"boolValue":true},` +
				`{"int64Value":42},` +
				`{"stringValue":"foo"}` +
				`]}`,
		},
		{
			&tu.ArbitraryValue{MapValue: map[string]*tu.ArbitraryValue{
				"bool":   {BoolValue: boolp(true)},
				"int64":  {Int64Value: int64p(42)},
				"string": {StringValue: stringp("foo")},
			}},
			`{"mapValue":{` +
				`"bool":{"boolValue":true},` +
				`"int64":{"int64Value":42},` +
				`"string":{"stringValue":"foo"}` +
				`}}`,
		},
		{
			&ts.List{Value: 0, Tail: &ts.List{Value: 1}},
			`{"value":0,"tail":{"value":1}}`,
		},
		{
			&ts.Rename{Default: "foo", CamelCase: "bar"},
			`{"default":"foo","snake_case":"bar"}`,
		},
		{
			&ts.Omit{Serialized: "foo", Hidden: ""},
			`{"serialized":"foo"}`,
		},
	}

	for _, tt := range tests {
		encoded, err := json.Marshal(tt.v)
		if assert.NoError(t, err, "failed to JSON encode %v", tt.v) {
			assert.Equal(t, tt.j, string(encoded))
		}

		v := reflect.New(reflect.TypeOf(tt.v).Elem()).Interface()
		if assert.NoError(t, json.Unmarshal([]byte(tt.j), v), "failed to decode %q", tt.j) {
			assert.Equal(t, tt.v, v)
		}
	}
}

func TestOptionalEmptyListJSON(t *testing.T) {
	give := tu.ArbitraryValue{ListValue: []*tu.ArbitraryValue{}}
	want := tu.ArbitraryValue{ListValue: nil}
	j := `{}`

	encoded, err := json.Marshal(give)
	require.NoError(t, err, "failed to encode to JSON")
	assert.Equal(t, j, string(encoded))

	var get tu.ArbitraryValue
	require.NoError(t, json.Unmarshal(encoded, &get), "failed to decode JSON")
	assert.Equal(t, want, get)
}

func TestOptionalEmptyMapJSON(t *testing.T) {
	give := tu.ArbitraryValue{MapValue: map[string]*tu.ArbitraryValue{}}
	want := tu.ArbitraryValue{MapValue: nil}
	j := `{}`

	encoded, err := json.Marshal(give)
	require.NoError(t, err, "failed to encode to JSON")
	assert.Equal(t, j, string(encoded))

	var get tu.ArbitraryValue
	require.NoError(t, json.Unmarshal(encoded, &get), "failed to decode JSON")
	assert.Equal(t, want, get)
}

func TestJSONOmitBehaviour(t *testing.T) {
	omit := ts.Omit{Serialized: "foo", Hidden: "bar"}
	b, err := json.Marshal(&omit)

	assert.NoError(t, err, "should marshal")
	assert.Equal(t, b, []byte(`{"serialized":"foo"}`))

	omit2 := ts.Omit{}
	err = json.Unmarshal([]byte(`{"serialized":"foo","hidden":"bar"}`), &omit2)

	assert.NoError(t, err, "should unmarshal")
	assert.Equal(t, omit2.Serialized, "foo")
	assert.Equal(t, omit2.Hidden, "")
}

func TestStructGoTags(t *testing.T) {
	gt := &ts.GoTags{
		Foo:                "shouldomit",
		FooBar:             "shouldnotomit",
		FooBarWithSpace:    "shouldalsonotomit",
		FooBarWithRequired: "shouldrequired",
	}
	b, err := json.Marshal(gt)
	assert.NoError(t, err, "should marshal")
	assert.JSONEq(t, string(b), `{"foobar":"shouldnotomit", "foobarWithSpace":"shouldalsonotomit","foobarWithRequired":"shouldrequired"}`)

	foo, _ := reflect.TypeOf(gt).Elem().FieldByName("Foo")
	assert.Equal(t, `json:"-" foo:"bar"`, string(foo.Tag))

	foobar, _ := reflect.TypeOf(gt).Elem().FieldByName("FooBar")
	assert.Equal(t, `json:"foobar,option1,option2,required" bar:"foo,option1" foo:"foobar"`, string(foobar.Tag))

	foobarWithSpace, _ := reflect.TypeOf(gt).Elem().FieldByName("FooBarWithSpace")
	assert.Equal(t, `json:"foobarWithSpace,required" foo:"foo bar foobar barfoo"`, string(foobarWithSpace.Tag))

	bar, _ := reflect.TypeOf(gt).Elem().FieldByName("Bar")
	assert.Equal(t, `json:"Bar,omitempty" bar:"foo"`, string(bar.Tag))

	foobarWithOmitEmpty, _ := reflect.TypeOf(gt).Elem().FieldByName("FooBarWithOmitEmpty")
	assert.Equal(t, `json:"foobarWithOmitEmpty,omitempty"`, string(foobarWithOmitEmpty.Tag))

	foobarWithRequired, _ := reflect.TypeOf(gt).Elem().FieldByName("FooBarWithRequired")
	assert.Equal(t, `json:"foobarWithRequired,required"`, string(foobarWithRequired.Tag))
}

func TestStructValidation(t *testing.T) {
	tests := []struct {
		desc        string
		serialize   thriftType
		deserialize wire.Value
		typ         reflect.Type // must be set if serialize is not
		wantError   string
	}{
		{
			desc: "Point: missing X",
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueDouble(42)},
			}}),
			typ:       reflect.TypeOf(ts.Point{}),
			wantError: "field Y of Point is required",
		},
		{
			desc: "Point: missing Y",
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueDouble(42)},
			}}),
			typ:       reflect.TypeOf(ts.Point{}),
			wantError: "field X of Point is required",
		},
		{
			desc: "Size: missing width",
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueDouble(42)},
			}}),
			typ:       reflect.TypeOf(ts.Size{}),
			wantError: "field Height of Size is required",
		},
		{
			desc: "Size: missing height",
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueDouble(42)},
			}}),
			typ:       reflect.TypeOf(ts.Size{}),
			wantError: "field Width of Size is required",
		},
		{
			desc:      "Frame: missing topLeft",
			serialize: &ts.Frame{Size: &ts.Size{Width: 1, Height: 2}},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueDouble(1)},
						{ID: 2, Value: wire.NewValueDouble(2)},
					}}),
				},
			}}),
			wantError: "field TopLeft of Frame is required",
		},
		{
			desc:      "Frame: missing Size",
			serialize: &ts.Frame{TopLeft: &ts.Point{X: 1, Y: 2}},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueDouble(1)},
						{ID: 2, Value: wire.NewValueDouble(2)},
					}}),
				},
			}}),
			wantError: "field Size of Frame is required",
		},
		{
			desc: "Frame: topLeft: missing y",
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueDouble(1)},
					}}),
				},
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueDouble(1)},
						{ID: 2, Value: wire.NewValueDouble(2)},
					}}),
				},
			}}),
			typ:       reflect.TypeOf(ts.Frame{}),
			wantError: "field Y of Point is required",
		},
		{
			desc:        "Graph: missing edges",
			serialize:   &ts.Graph{Edges: nil},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			wantError:   "field Edges of Graph is required",
		},
		{
			desc: "Graph: edges: misssing end",
			serialize: &ts.Graph{
				Edges: []*ts.Edge{
					{StartPoint: &ts.Point{X: 1, Y: 2}, EndPoint: &ts.Point{X: 3, Y: 4}},
					{StartPoint: &ts.Point{X: 5, Y: 6}},
				},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TStruct, []wire.Value{
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{
									ID: 1,
									Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
										{ID: 1, Value: wire.NewValueDouble(1)},
										{ID: 2, Value: wire.NewValueDouble(2)},
									}}),
								},
								{
									ID: 2,
									Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
										{ID: 1, Value: wire.NewValueDouble(3)},
										{ID: 2, Value: wire.NewValueDouble(4)},
									}}),
								},
							}}),
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{
									ID: 1,
									Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
										{ID: 1, Value: wire.NewValueDouble(1)},
										{ID: 2, Value: wire.NewValueDouble(2)},
									}}),
								},
							}}),
						}),
					),
				},
			}}),
			wantError: "field EndPoint of Edge is required",
		},
		{
			desc: "User: contact: missing emailAddress",
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("hello")},
				{
					ID:    2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
				},
			}}),
			typ:       reflect.TypeOf(ts.User{}),
			wantError: "field EmailAddress of ContactInfo is required",
		},
		{
			desc: "PrimitiveContainersRequired: missing list",
			serialize: &tc.PrimitiveContainersRequired{
				SetOfInts: map[int32]struct{}{
					1: {},
					2: {},
					3: {},
				},
				MapOfIntsToDoubles: map[int64]float64{1: 2.3, 4: 5.6},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 2,
					Value: wire.NewValueSet(
						wire.ValueListFromSlice(wire.TI32, []wire.Value{
							wire.NewValueI32(1),
							wire.NewValueI32(2),
							wire.NewValueI32(3),
						}),
					),
				},
				{
					ID: 3,
					Value: wire.NewValueMap(
						wire.MapItemListFromSlice(wire.TI64, wire.TDouble, []wire.MapItem{
							{
								Key:   wire.NewValueI64(1),
								Value: wire.NewValueDouble(2.3),
							},
							{
								Key:   wire.NewValueI64(4),
								Value: wire.NewValueDouble(5.6),
							},
						}),
					),
				},
			}}),
			wantError: "field ListOfStrings of PrimitiveContainersRequired is required",
		},
		{
			desc: "PrimitiveContainersRequired: missing set",
			serialize: &tc.PrimitiveContainersRequired{
				ListOfStrings:      []string{"hello", "world"},
				MapOfIntsToDoubles: map[int64]float64{1: 2.3, 4: 5.6},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TBinary, []wire.Value{
							wire.NewValueString("hello"),
							wire.NewValueString("world"),
						}),
					),
				},
				{
					ID: 3,
					Value: wire.NewValueMap(
						wire.MapItemListFromSlice(wire.TI64, wire.TDouble, []wire.MapItem{
							{
								Key:   wire.NewValueI64(1),
								Value: wire.NewValueDouble(2.3),
							},
							{
								Key:   wire.NewValueI64(4),
								Value: wire.NewValueDouble(5.6),
							},
						}),
					),
				},
			}}),
			wantError: "field SetOfInts of PrimitiveContainersRequired is required",
		},
		{
			desc: "PrimitiveContainersRequired: missing map",
			serialize: &tc.PrimitiveContainersRequired{
				ListOfStrings: []string{"hello", "world"},
				SetOfInts: map[int32]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TBinary, []wire.Value{
							wire.NewValueString("hello"),
							wire.NewValueString("world"),
						}),
					),
				},
				{
					ID: 2,
					Value: wire.NewValueSet(
						wire.ValueListFromSlice(wire.TI32, []wire.Value{
							wire.NewValueI32(1),
							wire.NewValueI32(2),
							wire.NewValueI32(3),
						}),
					),
				},
			}}),
			wantError: "field MapOfIntsToDoubles of PrimitiveContainersRequired is required",
		},
		{
			desc:        "Document: empty",
			serialize:   &tu.Document{},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			wantError:   "Document should have exactly one field: got 0 fields",
		},
		{
			desc: "Document: multiple",
			serialize: &tu.Document{
				Pdf:       td.PDF{},
				PlainText: stringp("hello"),
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueBinary([]byte{}),
				},
				{
					ID:    2,
					Value: wire.NewValueString("hello"),
				},
			}}),
			wantError: "Document should have exactly one field: got 2 fields",
		},
		{
			desc:        "ArbitraryValue: empty",
			serialize:   &tu.ArbitraryValue{},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			wantError:   "ArbitraryValue should have exactly one field: got 0 fields",
		},
		{
			desc: "ArbitraryValue: primitives",
			serialize: &tu.ArbitraryValue{
				BoolValue:   boolp(true),
				Int64Value:  int64p(42),
				StringValue: stringp(""),
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI64(42)},
				{ID: 3, Value: wire.NewValueString("")},
			}}),
			wantError: "ArbitraryValue should have exactly one field: got 3 fields",
		},
		{
			desc: "ArbitraryValue: full",
			serialize: &tu.ArbitraryValue{
				BoolValue:   boolp(true),
				Int64Value:  int64p(42),
				StringValue: stringp(""),
				ListValue: []*tu.ArbitraryValue{
					{BoolValue: boolp(true)},
					{Int64Value: int64p(42)},
					{StringValue: stringp("")},
				},
				MapValue: map[string]*tu.ArbitraryValue{
					"bool":   {BoolValue: boolp(true)},
					"int":    {Int64Value: int64p(42)},
					"string": {StringValue: stringp("")},
				},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
				{ID: 2, Value: wire.NewValueI64(42)},
				{ID: 3, Value: wire.NewValueString("")},
				{
					ID: 4,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TStruct, []wire.Value{
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueBool(true)},
							}}),
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 2, Value: wire.NewValueI64(42)},
							}}),
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 3, Value: wire.NewValueString("")},
							}}),
						}),
					),
				},
				{
					ID: 5,
					Value: wire.NewValueMap(
						wire.MapItemListFromSlice(wire.TBinary, wire.TStruct, []wire.MapItem{
							{
								Key: wire.NewValueString("bool"),
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 1, Value: wire.NewValueBool(true)},
								}}),
							},
							{
								Key: wire.NewValueString("int"),
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 2, Value: wire.NewValueI64(42)},
								}}),
							},
							{
								Key: wire.NewValueString("string"),
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 3, Value: wire.NewValueString("")},
								}}),
							},
						}),
					),
				},
			}}),
			wantError: "ArbitraryValue should have exactly one field: got 5 fields",
		},
		{
			desc: "ArbitraryValue: error inside a list",
			serialize: &tu.ArbitraryValue{
				ListValue: []*tu.ArbitraryValue{
					{BoolValue: boolp(true)},
					{},
				},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 4,
					Value: wire.NewValueList(
						wire.ValueListFromSlice(wire.TStruct, []wire.Value{
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueBool(true)},
							}}),
							wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
						})),
				},
			}}),
			wantError: "ArbitraryValue should have exactly one field: got 0 fields",
		},
		{
			desc: "ArbitraryValue: error inside a map value",
			serialize: &tu.ArbitraryValue{
				MapValue: map[string]*tu.ArbitraryValue{
					"bool":  {BoolValue: boolp(true)},
					"empty": {},
				},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 5,
					Value: wire.NewValueMap(
						wire.MapItemListFromSlice(wire.TBinary, wire.TStruct, []wire.MapItem{
							{
								Key: wire.NewValueString("bool"),
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 1, Value: wire.NewValueBool(true)},
								}}),
							},
							{
								Key:   wire.NewValueString("empty"),
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
							},
						}),
					),
				},
			}}),
			wantError: "ArbitraryValue should have exactly one field: got 0 fields",
		},
		{
			desc: "FrameGroup: error inside a set",
			serialize: &td.FrameGroup{
				&ts.Frame{
					TopLeft: &ts.Point{X: 1, Y: 2},
					Size:    &ts.Size{Width: 3, Height: 4},
				},
				&ts.Frame{TopLeft: &ts.Point{X: 5, Y: 6}},
			},
			deserialize: wire.NewValueSet(
				wire.ValueListFromSlice(wire.TStruct, []wire.Value{
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{
							ID: 1,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{
									ID:    1,
									Value: wire.NewValueDouble(1),
								},
								{
									ID:    2,
									Value: wire.NewValueDouble(2),
								},
							}}),
						},
						{
							ID: 2,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{
									ID:    1,
									Value: wire.NewValueDouble(3),
								},
								{
									ID:    2,
									Value: wire.NewValueDouble(4),
								},
							}}),
						},
					}}),
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{
							ID: 1,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{
									ID:    1,
									Value: wire.NewValueDouble(5),
								},
								{
									ID:    2,
									Value: wire.NewValueDouble(6),
								},
							}}),
						},
					}}),
				}),
			),
			wantError: "field Size of Frame is required",
		},
		{
			desc: "EdgeMap: error inside a map key",
			serialize: &td.EdgeMap{
				{
					Key:   &ts.Edge{StartPoint: &ts.Point{X: 1, Y: 2}},
					Value: &ts.Edge{StartPoint: &ts.Point{X: 3, Y: 4}, EndPoint: &ts.Point{X: 5, Y: 6}},
				},
			},
			deserialize: wire.NewValueMap(
				wire.MapItemListFromSlice(wire.TStruct, wire.TStruct, []wire.MapItem{
					{
						Key: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{
								ID: 1,
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 1, Value: wire.NewValueDouble(1)},
									{ID: 2, Value: wire.NewValueDouble(2)},
								}}),
							},
						}}),
						Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
							{
								ID: 1,
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 1, Value: wire.NewValueDouble(3)},
									{ID: 2, Value: wire.NewValueDouble(4)},
								}}),
							},
							{
								ID: 2,
								Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
									{ID: 1, Value: wire.NewValueDouble(5)},
									{ID: 2, Value: wire.NewValueDouble(6)},
								}}),
							},
						}}),
					},
				}),
			),
			wantError: "field EndPoint of Edge is required",
		},
	}

	for _, tt := range tests {
		var typ reflect.Type
		if tt.serialize != nil {
			typ = reflect.TypeOf(tt.serialize).Elem()
			v, err := tt.serialize.ToWire()
			if err == nil {
				err = wire.EvaluateValue(v)
			}
			if assert.Error(t, err, "%v: expected failure but got %v", tt.desc, v) {
				assert.Contains(t, err.Error(), tt.wantError, tt.desc)
			}
		} else {
			typ = tt.typ
		}

		if typ == nil {
			t.Fatalf("invalid test %q: either typ or serialize must be set", tt.desc)
		}

		x := reflect.New(typ).Interface().(thriftType)
		if err := x.FromWire(tt.deserialize); assert.Error(t, err, "%v: expected failure but got %v", tt.desc, x) {
			assert.Contains(t, err.Error(), tt.wantError, tt.desc)
		}
	}
}

func TestStructAccessors(t *testing.T) {
	t.Run("User", func(t *testing.T) {
		t.Run("Personal", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				u := ts.User{
					Personal: &ts.PersonalInfo{
						Age: ptr.Int32(30),
					},
				}
				assert.True(t, u.GetPersonal().IsSetAge())
				assert.Equal(t, int32(30), u.GetPersonal().GetAge())
			})
			t.Run("unset", func(t *testing.T) {
				var u *ts.User
				assert.False(t, u.GetPersonal().IsSetAge())
				assert.Equal(t, int32(0), u.GetPersonal().GetAge())
			})
		})
	})
	t.Run("DoesNotExistException", func(t *testing.T) {
		t.Run("set", func(t *testing.T) {
			err := tx.DoesNotExistException{Error2: ptr.String("foo")}
			assert.True(t, err.IsSetError2())
			assert.Equal(t, "foo", err.GetError2())
		})
		t.Run("unset", func(t *testing.T) {
			var err tx.DoesNotExistException
			assert.False(t, err.IsSetError2())
			assert.Equal(t, "", err.GetError2())
		})
	})
	t.Run("DefaultsStruct", func(t *testing.T) {
		t.Run("RequiredPrimitive", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				s := ts.DefaultsStruct{RequiredPrimitive: ptr.Int32(42)}
				assert.True(t, s.IsSetRequiredPrimitive())
				assert.Equal(t, int32(42), s.GetRequiredPrimitive())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetRequiredPrimitive())
				assert.Equal(t, int32(100), s.GetRequiredPrimitive())
			})
		})

		t.Run("OptionalPrimitive", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				s := ts.DefaultsStruct{OptionalPrimitive: ptr.Int32(100)}
				assert.True(t, s.IsSetOptionalPrimitive())
				assert.Equal(t, int32(100), s.GetOptionalPrimitive())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetOptionalPrimitive())
				assert.Equal(t, int32(200), s.GetOptionalPrimitive())
			})
		})

		t.Run("RequiredEnum", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				e := te.EnumDefaultFoo
				s := ts.DefaultsStruct{RequiredEnum: &e}
				assert.True(t, s.IsSetRequiredEnum())
				assert.Equal(t, te.EnumDefaultFoo, s.GetRequiredEnum())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetRequiredEnum())
				assert.Equal(t, te.EnumDefaultBar, s.GetRequiredEnum())
			})
		})

		t.Run("OptionalEnum", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				e := te.EnumDefaultFoo
				s := ts.DefaultsStruct{OptionalEnum: &e}
				assert.True(t, s.IsSetOptionalEnum())
				assert.Equal(t, te.EnumDefaultFoo, s.GetOptionalEnum())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetOptionalEnum())
				assert.Equal(t, te.EnumDefaultBaz, s.GetOptionalEnum())
			})
		})

		t.Run("RequiredList", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				lst := []string{"foo", "bar", "baz"}
				s := ts.DefaultsStruct{RequiredList: lst}
				assert.True(t, s.IsSetRequiredList())
				assert.Equal(t, lst, s.GetRequiredList())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetRequiredList())
				assert.Equal(t, []string{"hello", "world"}, s.GetRequiredList())
			})
		})
		t.Run("OptionalList", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				lst := []float64{0, 1, 2}
				s := ts.DefaultsStruct{OptionalList: lst}
				assert.True(t, s.IsSetOptionalList())
				assert.Equal(t, lst, s.GetOptionalList())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetOptionalList())
				assert.Equal(t, []float64{1, 2, 3}, s.GetOptionalList())
			})
		})
		t.Run("RequiredStruct", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				f := &ts.Frame{
					TopLeft: &ts.Point{X: 1, Y: 1},
					Size:    &ts.Size{Width: 1, Height: 1},
				}
				s := ts.DefaultsStruct{RequiredStruct: f}
				assert.True(t, s.IsSetRequiredStruct())
				assert.Equal(t, f, s.GetRequiredStruct())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetRequiredStruct())
				assert.Equal(t,
					&ts.Frame{
						TopLeft: &ts.Point{X: 1, Y: 2},
						Size:    &ts.Size{Width: 100, Height: 200},
					},
					s.GetRequiredStruct(),
				)

			})
		})
		t.Run("OptionalStruct", func(t *testing.T) {
			t.Run("set", func(t *testing.T) {
				e := &ts.Edge{
					StartPoint: &ts.Point{X: 1.0, Y: 1.0},
					EndPoint:   &ts.Point{X: 1.0, Y: 1.0},
				}
				s := ts.DefaultsStruct{OptionalStruct: e}
				assert.True(t, s.IsSetOptionalStruct())
				assert.Equal(t, e, s.GetOptionalStruct())
			})
			t.Run("unset", func(t *testing.T) {
				var s ts.DefaultsStruct
				assert.False(t, s.IsSetOptionalStruct())
				assert.Equal(t,
					&ts.Edge{
						StartPoint: &ts.Point{X: 1.0, Y: 2.0},
						EndPoint:   &ts.Point{X: 3.0, Y: 4.0},
					},
					s.GetOptionalStruct(),
				)
			})
		})
	})
}

func TestEmptyPrimitivesRoundTrip(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		give := ts.PrimitiveRequiredStruct{
			BoolField:   false,
			ByteField:   0,
			Int16Field:  0,
			Int32Field:  0,
			Int64Field:  0,
			DoubleField: 0.0,
			StringField: "",
			BinaryField: []byte{},
		}

		b, err := json.Marshal(give)
		require.NoError(t, err, "failed to encode to JSON")

		var decoded ts.PrimitiveRequiredStruct
		require.NoError(t, json.Unmarshal(b, &decoded), "failed to decode JSON")
		assert.Equal(t, give, decoded)

		v, err := decoded.ToWire()
		require.NoError(t, err, "failed to convert to wire.Value")

		var got ts.PrimitiveRequiredStruct
		require.NoError(t, got.FromWire(v), "failed to convert from wire.Value")
		assert.Equal(t, give, got)
	})

	t.Run("optional", func(t *testing.T) {
		give := ts.PrimitiveOptionalStruct{
			BoolField:   ptr.Bool(false),
			ByteField:   ptr.Int8(0),
			Int16Field:  ptr.Int16(0),
			Int32Field:  ptr.Int32(0),
			Int64Field:  ptr.Int64(0),
			DoubleField: ptr.Float64(0.0),
			StringField: ptr.String(""),
			// BinaryField is a slice so we won't check it here to avoid nil
			// vs empty slice mismatch.
		}

		b, err := json.Marshal(give)
		require.NoError(t, err, "failed to encode to JSON")

		var decoded ts.PrimitiveOptionalStruct
		require.NoError(t, json.Unmarshal(b, &decoded), "failed to decode JSON")
		assert.Equal(t, give, decoded)

		v, err := decoded.ToWire()
		require.NoError(t, err, "failed to convert to wire.Value")

		var got ts.PrimitiveOptionalStruct
		require.NoError(t, got.FromWire(v), "failed to convert from wire.Value")
		assert.Equal(t, give, got)
	})
}

func TestStructLabel(t *testing.T) {
	// Convenience type to build map[string]interface{}.
	type attrs = map[string]interface{}

	tests := []struct {
		desc   string
		give   ts.StructLabels
		json   string
		logged attrs
	}{
		{
			desc:   "keyword as label",
			give:   ts.StructLabels{IsRequired: ptr.Bool(true)},
			json:   `{"required": true}`,
			logged: attrs{"required": true},
		},
		{
			desc:   "JSON tag overrides label",
			give:   ts.StructLabels{Foo: ptr.String("foo")},
			json:   `{"not_bar": "foo"}`,
			logged: attrs{"bar": "foo"},
		},
		{
			desc:   "empty label",
			give:   ts.StructLabels{Qux: ptr.String("foo")},
			json:   `{"qux": "foo"}`,
			logged: attrs{"qux": "foo"},
		},
		{
			desc:   "all caps",
			give:   ts.StructLabels{Quux: ptr.String("foo")},
			json:   `{"QUUX": "foo"}`,
			logged: attrs{"QUUX": "foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Run("JSON", func(t *testing.T) {
				b, err := json.Marshal(tt.give)
				require.NoError(t, err)

				require.JSONEq(t, tt.json, string(b))

				var got ts.StructLabels
				require.NoError(t, json.Unmarshal(b, &got))
				assert.Equal(t, tt.give, got)
			})

			t.Run("logging", func(t *testing.T) {
				enc := zapcore.NewMapObjectEncoder()
				tt.give.MarshalLogObject(enc)
				assert.Equal(t, tt.logged, enc.Fields)
			})
		})
	}
}

func TestNoErrorFlagExceptions(t *testing.T) {
	var noErrorException interface{}
	noErrorException = &noerror.NoErrorException{}
	_, ok := noErrorException.(error)
	assert.False(t, ok)
}
