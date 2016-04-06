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

	tc "github.com/thriftrw/thriftrw-go/gen/testdata/containers"
	tx "github.com/thriftrw/thriftrw-go/gen/testdata/exceptions"
	ts "github.com/thriftrw/thriftrw-go/gen/testdata/structs"
	td "github.com/thriftrw/thriftrw-go/gen/testdata/typedefs"
	tu "github.com/thriftrw/thriftrw-go/gen/testdata/unions"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
)

func TestPrimitiveRequiredStructWire(t *testing.T) {
	tests := []struct {
		s ts.PrimitiveRequiredStruct
		v wire.Value
	}{
		{
			ts.PrimitiveRequiredStruct{
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

		var s ts.PrimitiveRequiredStruct
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
		var s ts.PrimitiveRequiredStruct
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
		s ts.PrimitiveOptionalStruct
		v wire.Value
	}{
		{
			ts.PrimitiveOptionalStruct{
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
			ts.PrimitiveOptionalStruct{BoolField: boolp(true)},
			singleFieldStruct(1, wire.NewValueBool(true)),
		},
		{
			ts.PrimitiveOptionalStruct{ByteField: bytep(1)},
			singleFieldStruct(2, wire.NewValueI8(1)),
		},
		{
			ts.PrimitiveOptionalStruct{Int16Field: int16p(2)},
			singleFieldStruct(3, wire.NewValueI16(2)),
		},
		{
			ts.PrimitiveOptionalStruct{Int32Field: int32p(3)},
			singleFieldStruct(4, wire.NewValueI32(3)),
		},
		{
			ts.PrimitiveOptionalStruct{Int64Field: int64p(4)},
			singleFieldStruct(5, wire.NewValueI64(4)),
		},
		{
			ts.PrimitiveOptionalStruct{DoubleField: doublep(5.0)},
			singleFieldStruct(6, wire.NewValueDouble(5.0)),
		},
		{
			ts.PrimitiveOptionalStruct{StringField: stringp("foo")},
			singleFieldStruct(7, wire.NewValueString("foo")),
		},
		{
			ts.PrimitiveOptionalStruct{BinaryField: []byte("bar")},
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

		var s ts.PrimitiveOptionalStruct
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestPrimitiveContainersRequired(t *testing.T) {
	tests := []struct {
		s tc.PrimitiveContainersRequired
		v wire.Value
	}{
		{
			tc.PrimitiveContainersRequired{
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

		var s tc.PrimitiveContainersRequired
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}
	}
}

func TestNestedStructsRequired(t *testing.T) {
	tests := []struct {
		s ts.Frame
		v wire.Value
		o string
	}{
		{
			ts.Frame{
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
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s ts.Frame
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		assert.Equal(t, tt.o, tt.s.String())
	}
}

func TestNestedStructsOptional(t *testing.T) {
	tests := []struct {
		s ts.User
		v wire.Value
		o string
	}{
		{
			ts.User{Name: "Foo Bar"},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("Foo Bar")},
			}}),
			"User{Name: Foo Bar}",
		},
		{
			ts.User{
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
	}

	for _, tt := range tests {
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s ts.User
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
			&ts.Frame{TopLeft: &ts.Point{}},
			"Frame{TopLeft: Point{X: 0, Y: 0}, Size: <nil>}",
		},
		{
			&ts.Frame{Size: &ts.Size{}},
			"Frame{TopLeft: <nil>, Size: Size{Width: 0, Height: 0}}",
		},
		{
			&ts.Edge{Start: &ts.Point{X: 1, Y: 2}},
			"Edge{Start: Point{X: 1, Y: 2}, End: <nil>}",
		},
		{
			&ts.Edge{End: &ts.Point{X: 3, Y: 4}},
			"Edge{Start: <nil>, End: Point{X: 3, Y: 4}}",
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
		assert.True(
			t,
			wire.ValuesAreEqual(tt.v, tt.s.ToWire()),
			"%v.ToWire() != %v", tt.s, tt.v,
		)

		var s tx.DoesNotExistException
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		err := error(&s) // should implement the error interface
		assert.Equal(t, "DoesNotExistException{Key: foo}", err.Error())
	}
}

func TestUnionSimple(t *testing.T) {
	tests := []struct {
		s tu.Document
		v wire.Value
		o string
	}{
		{
			tu.Document{Pdf: []byte{1, 2, 3}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBinary([]byte{1, 2, 3})},
			}}),
			"Document{Pdf: [1 2 3]}",
		},
		{
			tu.Document{PlainText: stringp("hello")},
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

		var s tu.Document
		if assert.NoError(t, s.FromWire(tt.v)) {
			assert.Equal(t, tt.s, s)
		}

		assert.Equal(t, tt.o, tt.s.String())
	}
}

func TestUnionComplex(t *testing.T) {
	tests := []struct {
		s tu.ArbitraryValue
		v wire.Value
		o string
	}{
		{
			tu.ArbitraryValue{BoolValue: boolp(true)},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueBool(true)},
			}}),
			"ArbitraryValue{BoolValue: true}",
		},
		{
			tu.ArbitraryValue{Int64Value: int64p(42)},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI64(42)},
			}}),
			"ArbitraryValue{Int64Value: 42}",
		},
		{
			tu.ArbitraryValue{StringValue: stringp("hello")},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			"ArbitraryValue{StringValue: hello}",
		},
		{
			tu.ArbitraryValue{ListValue: []*tu.ArbitraryValue{
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
			tu.ArbitraryValue{MapValue: map[string]*tu.ArbitraryValue{
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

		var s tu.ArbitraryValue
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
		o    ts.ContactInfo
	}{
		{
			"unknown field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
				{ID: 2, Value: wire.NewValueI32(42)},
			}}),
			ts.ContactInfo{EmailAddress: "foo"},
		},
		{
			"only unknown field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueString("bar")},
			}}),
			ts.ContactInfo{},
		},
		{
			"wrong type for recognized field",
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(42)},
			}}),
			ts.ContactInfo{},
		},
	}

	for _, tt := range tests {
		var o ts.ContactInfo
		if assert.NoError(t, o.FromWire(tt.i)) {
			assert.Equal(t, tt.o, o)
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
			success: &tu.Document{
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
			success: &tu.Document{Pdf: []byte{1, 2, 3}},
		},
		{
			desc: "only unrecognized fields",
			input: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueI32(42)}, // also a type mismatch
				{ID: 3, Value: wire.NewValueString("hello")},
			}}),
			success: &tu.Document{},
			// TODO(abg): If the union is empty, we need to fail the request
		},
		{
			desc:    "no fields",
			input:   wire.NewValueStruct(wire.Struct{}),
			success: &tu.Document{},
			// TODO(abg): If the union is empty, we need to fail the request
		},
	}

	for _, tt := range tests {
		var o tu.Document
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

func TestEmptyStruct(t *testing.T) {
	var x, y ts.EmptyStruct
	v := wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}})

	assert.Equal(t, x.ToWire(), v)
	if assert.NoError(t, y.FromWire(v)) {
		assert.Equal(t, x, y)
	}
}

func TestEmptyUnion(t *testing.T) {
	var x, y tu.EmptyUnion
	v := wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}})

	assert.Equal(t, x.ToWire(), v)
	if assert.NoError(t, y.FromWire(v)) {
		assert.Equal(t, x, y)
	}
}

func TestEmptyException(t *testing.T) {
	var x, y tx.EmptyException
	v := wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}})

	assert.Equal(t, x.ToWire(), v)
	if assert.NoError(t, y.FromWire(v)) {
		assert.Equal(t, x, y)
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
				Start: &ts.Point{X: 1, Y: 2},
				End:   &ts.Point{X: 3, Y: 4},
			},
			`{"start":{"x":1,"y":2},"end":{"x":3,"y":4}}`,
		},
		{
			&ts.Edge{Start: &ts.Point{X: 1, Y: 1}},
			`{"start":{"x":1,"y":1},"end":null}`,
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
		{&tu.Document{Pdf: td.Pdf("hello")}, `{"pdf":"aGVsbG8="}`},
		{&tu.Document{Pdf: td.Pdf{}}, `{"pdf":""}`},
		{
			&tu.Document{PlainText: stringp("hello")},
			`{"pdf":null,"plainText":"hello"}`,
		},
		{&tu.Document{PlainText: stringp("")}, `{"pdf":null,"plainText":""}`},
		{
			&tu.ArbitraryValue{BoolValue: boolp(true)},
			`{"boolValue":true,"listValue":null,"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{BoolValue: boolp(false)},
			`{"boolValue":false,"listValue":null,"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{Int64Value: int64p(42)},
			`{"int64Value":42,"listValue":null,"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{Int64Value: int64p(0)},
			`{"int64Value":0,"listValue":null,"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{StringValue: stringp("foo")},
			`{"stringValue":"foo","listValue":null,"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{StringValue: stringp("")},
			`{"stringValue":"","listValue":null,"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{ListValue: []*tu.ArbitraryValue{
				{BoolValue: boolp(true)},
				{Int64Value: int64p(42)},
				{StringValue: stringp("foo")},
			}},
			`{"listValue":[` +
				`{"boolValue":true,"listValue":null,"mapValue":null},` +
				`{"int64Value":42,"listValue":null,"mapValue":null},` +
				`{"stringValue":"foo","listValue":null,"mapValue":null}` +
				`],"mapValue":null}`,
		},
		{
			&tu.ArbitraryValue{MapValue: map[string]*tu.ArbitraryValue{
				"bool":   {BoolValue: boolp(true)},
				"int64":  {Int64Value: int64p(42)},
				"string": {StringValue: stringp("foo")},
			}},
			`{"listValue":null,"mapValue":{` +
				`"bool":{"boolValue":true,"listValue":null,"mapValue":null},` +
				`"int64":{"int64Value":42,"listValue":null,"mapValue":null},` +
				`"string":{"stringValue":"foo","listValue":null,"mapValue":null}` +
				`}}`,
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
