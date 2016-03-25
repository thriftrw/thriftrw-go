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

	"github.com/thriftrw/thriftrw-go/gen/testdata/test"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

func TestPrimitiveRequiredStructWire(t *testing.T) {
	tests := []struct {
		s test.PrimitiveRequiredStruct
		v wire.Value
	}{
		{
			test.PrimitiveRequiredStruct{
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

		var s test.PrimitiveRequiredStruct
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
		var s test.PrimitiveRequiredStruct
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
		s test.PrimitiveOptionalStruct
		v wire.Value
	}{
		{
			test.PrimitiveOptionalStruct{
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
			test.PrimitiveOptionalStruct{BoolField: boolp(true)},
			singleFieldStruct(1, wire.NewValueBool(true)),
		},
		{
			test.PrimitiveOptionalStruct{ByteField: bytep(1)},
			singleFieldStruct(2, wire.NewValueI8(1)),
		},
		{
			test.PrimitiveOptionalStruct{Int16Field: int16p(2)},
			singleFieldStruct(3, wire.NewValueI16(2)),
		},
		{
			test.PrimitiveOptionalStruct{Int32Field: int32p(3)},
			singleFieldStruct(4, wire.NewValueI32(3)),
		},
		{
			test.PrimitiveOptionalStruct{Int64Field: int64p(4)},
			singleFieldStruct(5, wire.NewValueI64(4)),
		},
		{
			test.PrimitiveOptionalStruct{DoubleField: doublep(5.0)},
			singleFieldStruct(6, wire.NewValueDouble(5.0)),
		},
		{
			test.PrimitiveOptionalStruct{StringField: stringp("foo")},
			singleFieldStruct(7, wire.NewValueString("foo")),
		},
		{
			test.PrimitiveOptionalStruct{BinaryField: []byte("bar")},
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

		var s test.PrimitiveOptionalStruct
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestPrimitiveContainersRequired(t *testing.T) {
	tests := []struct {
		s test.PrimitiveContainersRequired
		v wire.Value
	}{
		{
			test.PrimitiveContainersRequired{
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

		var s test.PrimitiveContainersRequired
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestNestedStructsRequired(t *testing.T) {
	tests := []struct {
		s test.Frame
		v wire.Value
		o string
	}{
		{
			test.Frame{
				TopLeft: &test.Point{X: 1, Y: 2},
				Size:    &test.Size{Width: 100, Height: 200},
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

		var s test.Frame
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		assert.Equal(t, tt.o, tt.s.String())
	}
}

func TestNestedStructsOptional(t *testing.T) {
	tests := []struct {
		s test.User
		v wire.Value
		o string
	}{
		{
			test.User{Name: "Foo Bar"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
			}}),
			"User{Name: Foo Bar}",
		},
		{
			test.User{
				Name:    "Foo Bar",
				Contact: &test.ContactInfo{EmailAddress: "foo@example.com"},
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

		var s test.User
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
			&test.Frame{TopLeft: &test.Point{}},
			"Frame{Size: <nil>, TopLeft: Point{X: 0, Y: 0}}",
		},
		{
			&test.Frame{Size: &test.Size{}},
			"Frame{Size: Size{Height: 0, Width: 0}, TopLeft: <nil>}",
		},
		{
			&test.Edge{Start: &test.Point{X: 1, Y: 2}},
			"Edge{End: <nil>, Start: Point{X: 1, Y: 2}}",
		},
		{
			&test.Edge{End: &test.Point{X: 3, Y: 4}},
			"Edge{End: Point{X: 3, Y: 4}, Start: <nil>}",
		},
		{
			&test.Graph{},
			"Graph{Edges: []}",
		},
		{
			&test.Event{},
			"Event{UUID: <nil>}",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.o, tt.i.String())
	}
}

func TestBasicException(t *testing.T) {
	tests := []struct {
		s test.DoesNotExistException
		v wire.Value
	}{
		{
			test.DoesNotExistException{Key: "foo"},
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

		var s test.DoesNotExistException
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		err := error(&s) // should implement the error interface
		assert.Equal(t, "DoesNotExistException{Key: foo}", err.Error())
	}
}

func TestUnionSimple(t *testing.T) {
	tests := []struct {
		s test.Document
		v wire.Value
		o string
	}{
		{
			test.Document{Pdf: []byte{1, 2, 3}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
			}}),
			"Document{Pdf: [1 2 3]}",
		},
		{
			test.Document{PlainText: stringp("hello")},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueString("hello")},
			}}),
			"Document{PlainText: hello}",
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s test.Document
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		assert.Equal(t, tt.o, tt.s.String())
	}
}

func TestUnionComplex(t *testing.T) {
	tests := []struct {
		s test.ArbitraryValue
		v wire.Value
		o string
	}{
		{
			test.ArbitraryValue{BoolValue: boolp(true)},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
			}}),
			"ArbitraryValue{BoolValue: true}",
		},
		{
			test.ArbitraryValue{Int64Value: int64p(42)},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI64(42)},
			}}),
			"ArbitraryValue{Int64Value: 42}",
		},
		{
			test.ArbitraryValue{StringValue: stringp("hello")},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			"ArbitraryValue{StringValue: hello}",
		},
		{
			test.ArbitraryValue{ListValue: []*test.ArbitraryValue{
				{BoolValue: boolp(true)},
				{Int64Value: int64p(42)},
				{StringValue: stringp("hello")},
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 4, Value: wire.NewValueList(wire.List{
					ValueType: wire.TStruct,
					Size:      3,
					Items: wire.ValueListFromSlice([]wire.Value{
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
				})},
			}}),
			"ArbitraryValue{ListValue: [ArbitraryValue{BoolValue: true} ArbitraryValue{Int64Value: 42} ArbitraryValue{StringValue: hello}]}",
		},
		{
			test.ArbitraryValue{MapValue: map[string]*test.ArbitraryValue{
				"bool":   {BoolValue: boolp(true)},
				"int64":  {Int64Value: int64p(42)},
				"string": {StringValue: stringp("hello")},
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 5, Value: wire.NewValueMap(wire.Map{
					KeyType:   wire.TBinary,
					ValueType: wire.TStruct,
					Size:      3,
					Items: wire.MapItemListFromSlice([]wire.MapItem{
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
				})},
			}}),
			"",
		},
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s test.ArbitraryValue
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		if tt.o != "" {
			assert.Equal(t, tt.o, tt.s.String())
		}
	}
}

func TestStructFromWireUnrecognizedField(t *testing.T) {
	tests := []struct {
		desc string
		i    wire.Value
		o    test.ContactInfo
	}{
		{
			"unknown field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
				{ID: 2, Value: wire.NewValueI32(42)},
			}}),
			test.ContactInfo{EmailAddress: "foo"},
		},
		{
			"only unknown field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueString("bar")},
			}}),
			test.ContactInfo{},
		},
		{
			"wrong type for recognized field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(42)},
			}}),
			test.ContactInfo{},
		},
	}

	for _, tt := range tests {
		var o test.ContactInfo
		if assert.NoError(t, o.FromWire(tt.i)) {
			assert.Equal(t, tt.o, o)
		}
	}
}

func TestUnionFromWireInconsistencies(t *testing.T) {
	tests := []struct {
		desc    string
		input   wire.Value
		success *test.Document
		failure string
	}{
		{
			desc: "multiple recognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
				{ID: 2, Value: wire.NewValueString("hello")},
			}}),
			success: &test.Document{
				Pdf:       []byte{1, 2, 3},
				PlainText: stringp("hello"),
			},
			// TODO(abg): Once we're validating, this will become a failure
			// case. We want to fail if mutliple fields are set on a union.
		},
		{
			desc: "recognized and unrecognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			success: &test.Document{Pdf: []byte{1, 2, 3}},
		},
		{
			desc: "only unrecognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI32(42)}, // also a type mismatch
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			success: &test.Document{},
			// TODO(abg): If the union is empty, we need to fail the request
		},
		{
			desc:    "no fields",
			input:   wire.NewValueStruct(wire.Struct{}),
			success: &test.Document{},
			// TODO(abg): If the union is empty, we need to fail the request
		},
	}

	for _, tt := range tests {
		var o test.Document
		err := o.FromWire(tt.input)
		if tt.success != nil {
			if assert.NoError(t, err) {
				assert.Equal(t, tt.success, &o, tt.desc)
			}
		} else {
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), tt.failure, tt.desc)
			}
		}
	}
}
