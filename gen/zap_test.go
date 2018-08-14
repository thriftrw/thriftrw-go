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
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	td "go.uber.org/thriftrw/gen/internal/tests/typedefs"
	"go.uber.org/zap/zapcore"
)

func TestCollectionsOfPrimitivesZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}

	b64 := func(byteString string) string {
		return base64.StdEncoding.EncodeToString([]byte(byteString))
	}

	tests := []struct {
		desc string
		p    tc.PrimitiveContainers
		v    interface{}
	}{
		// Lists /////////////////////////////////////////////////////////////
		{
			"empty list",
			tc.PrimitiveContainers{ListOfInts: []int64{}},
			o{"listOfInts": a{}},
		},
		{
			"list of ints",
			tc.PrimitiveContainers{ListOfInts: []int64{1, 2, 3}},
			o{"listOfInts": a{int64(1), int64(2), int64(3)}},
		},
		{
			"list of binary",
			tc.PrimitiveContainers{
				ListOfBinary: [][]byte{
					[]byte("foo"), {}, []byte("bar"), []byte("baz"),
				},
			},
			// base64
			o{"listOfBinary": a{b64("foo"), b64(""), b64("bar"), b64("baz")}},
		},
		// Sets //////////////////////////////////////////////////////////////
		{
			"empty set",
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{}},
			o{"setOfStrings": a{}},
		},
		{
			"set of strings",
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{
				"foo": {},
			}},
			o{"setOfStrings": a{"foo"}},
		},
		{
			"set of bytes",
			tc.PrimitiveContainers{SetOfBytes: map[int8]struct{}{
				125: {},
			}},
			o{"setOfBytes": a{int8(125)}},
		},
		// Maps //////////////////////////////////////////////////////////////
		{
			"empty map",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{}},
			o{"mapOfStringToBool": o{}},
		},
		{
			"map of int to string",
			tc.PrimitiveContainers{MapOfIntToString: map[int32]string{
				1234: "bar",
			}},
			o{"mapOfIntToString": a{
				o{"key": int32(1234),
					"value": "bar"},
			}},
		},
		{
			"map of string to bool",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{
				"foo": false,
				"bar": true,
				"baz": true,
			}},
			o{"mapOfStringToBool": o{"foo": false, "bar": true, "baz": true}},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.p.MarshalLogObject(mapEncoder)
		assert.Equal(t, tt.v, mapEncoder.Fields)
	}
}

func TestNondeterministicZapLogging(t *testing.T) {
	mapEncoder := zapcore.NewMapObjectEncoder()
	// case SetOfStrings
	test1 := tc.PrimitiveContainers{
		SetOfStrings: map[string]struct{}{
			"foo":   {},
			"bar":   {},
			"baz":   {},
			"hello": {},
			"world": {},
		},
	}
	test1.MarshalLogObject(mapEncoder)
	expected1 := []string{"foo", "bar", "baz", "hello", "world"}
	assert.ElementsMatch(t, mapEncoder.Fields["setOfStrings"], expected1)

	// case SetOfBytes
	test2 := tc.PrimitiveContainers{
		SetOfBytes: map[int8]struct{}{
			125: {},
			1:   {},
			5:   {},
			73:  {},
			42:  {},
		},
	}
	test2.MarshalLogObject(mapEncoder)
	expected2 := []int8{1, 5, 42, 73, 125}
	assert.ElementsMatch(t, mapEncoder.Fields["setOfBytes"], expected2)

	// case MapOfIntToString
	test3 := tc.PrimitiveContainers{
		MapOfIntToString: map[int32]string{
			125: "foo",
			1:   "bar",
			5:   "baz",
			73:  "hello",
			42:  "world",
		},
	}
	test3.MarshalLogObject(mapEncoder)
	expected3 := []map[string]interface{}{
		{"key": int32(125), "value": "foo"},
		{"key": int32(1), "value": "bar"},
		{"key": int32(5), "value": "baz"},
		{"key": int32(73), "value": "hello"},
		{"key": int32(42), "value": "world"},
	}
	assert.ElementsMatch(t, mapEncoder.Fields["mapOfIntToString"], expected3)
}

func TestEnumContainersZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}

	tests := []struct {
		r tc.EnumContainers
		v interface{}
	}{
		{
			tc.EnumContainers{
				ListOfEnums: []te.EnumDefault{
					te.EnumDefaultBar,
				},
			},
			o{"listOfEnums": a{
				o{"name": "Bar", "value": int32(1)}},
			},
		},
		{
			tc.EnumContainers{
				SetOfEnums: map[te.EnumWithValues]struct{}{
					te.EnumWithValuesZ: {},
				},
			},
			o{"setOfEnums": a{
				o{"name": "Z", "value": int32(789)}},
			},
		},
		{
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesP: 123,
				},
			},
			o{"mapOfEnums": a{
				o{"key": o{"value": int32(0), "name": "P"}, "value": int32(123)}}},
		},
		{
			// unknown enum name
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValues(1523): 123,
				},
			},
			o{"mapOfEnums": a{
				o{"key": o{"value": int32(1523)}, "value": int32(123)}}},
		},
		{
			// this is the same as the one above except we're using "R" intsead
			// of "P" (they both have the same value)
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesR: 123,
				},
			},
			o{"mapOfEnums": a{
				o{"key": o{"value": int32(0), "name": "P"}, "value": int32(123)}}},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.r.MarshalLogObject(mapEncoder)
		assert.Equal(t, tt.v, mapEncoder.Fields)
	}
}

func TestListOfStructsZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}

	tests := []struct {
		r ts.Graph
		v interface{}
	}{
		{
			ts.Graph{Edges: []*ts.Edge{}},
			o{"edges": a{}},
		},
		{
			ts.Graph{Edges: []*ts.Edge{
				{
					StartPoint: &ts.Point{X: 1.0, Y: 2.0},
					EndPoint:   &ts.Point{X: 3.0, Y: 4.0},
				},
				{
					StartPoint: &ts.Point{X: 5.0, Y: 6.0},
					EndPoint:   &ts.Point{X: 7.0, Y: 8.0},
				},
				{
					StartPoint: &ts.Point{X: 9.0, Y: 10.0},
					EndPoint:   &ts.Point{X: 11.0, Y: 12.0},
				},
			}},
			o{"edges": a{
				o{"startPoint": o{"x": float64(1), "y": float64(2)},
					"endPoint": o{"x": float64(3), "y": float64(4)}},
				o{"startPoint": o{"x": float64(5), "y": float64(6)},
					"endPoint": o{"y": float64(8), "x": float64(7)}},
				o{"startPoint": o{"x": float64(9), "y": float64(10)},
					"endPoint": o{"x": float64(11), "y": float64(12)}},
			},
			},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.r.MarshalLogObject(mapEncoder)
		assert.Equal(t, tt.v, mapEncoder.Fields)
	}
}

func TestCrazyTownZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}

	tests := []struct {
		desc string
		r    tc.ContainersOfContainers
		v    interface{}
	}{
		{
			"ListOfLists",
			tc.ContainersOfContainers{
				ListOfLists: [][]int32{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			o{"listOfLists": a{
				a{int32(1), int32(2), int32(3)},
				a{int32(4), int32(5), int32(6)}},
			},
		},
		{
			"ListOfSets",
			tc.ContainersOfContainers{
				ListOfSets: []map[int32]struct{}{
					{
						2: struct{}{},
					},
					{
						5: struct{}{},
					},
				},
			},
			o{"listOfSets": a{
				a{int32(2)},
				a{int32(5)}},
			},
		},
		{
			"ListOfMaps",
			tc.ContainersOfContainers{
				ListOfMaps: []map[int32]int32{
					{
						2: 200,
					},
					{
						5: 500,
					},
				},
			},
			o{"listOfMaps": a{
				a{
					o{"key": int32(2),
						"value": int32(200)}},
				a{
					o{"key": int32(5),
						"value": int32(500)}}},
			},
		},
		{
			"SetOfSets",
			tc.ContainersOfContainers{
				SetOfSets: []map[string]struct{}{
					{
						"2": struct{}{},
					},
					{
						"5": struct{}{},
					},
				},
			},
			o{"setOfSets": a{
				a{"2"},
				a{"5"}},
			},
		},
		{
			"SetOfLists",
			tc.ContainersOfContainers{
				SetOfLists: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
				},
			},
			o{"setOfLists": a{
				a{"1", "2", "3"},
				a{"4", "5", "6"}},
			},
		},
		{
			"SetOfMaps",
			tc.ContainersOfContainers{
				SetOfMaps: []map[string]string{
					{
						"2": "two",
					},
					{
						"6": "six",
					},
				},
			},
			o{"setOfMaps": a{
				o{"2": "two"},
				o{"6": "six"}},
			},
		},
		{
			"MapOfMapToInt",
			tc.ContainersOfContainers{
				MapOfMapToInt: []struct {
					Key   map[string]int32
					Value int64
				}{
					{
						Key:   map[string]int32{"2": 2},
						Value: 123,
					},
					{
						Key:   map[string]int32{"4": 4},
						Value: 135,
					},
				},
			},
			o{"mapOfMapToInt": a{
				o{"key": o{"2": int32(2)},
					"value": int64(123)},
				o{"key": o{"4": int32(4)},
					"value": int64(135)}},
			},
		},
		{
			"MapOfListToSet",
			tc.ContainersOfContainers{
				MapOfListToSet: []struct {
					Key   []int32
					Value map[int64]struct{}
				}{
					{
						Key: []int32{1, 2, 3},
						Value: map[int64]struct{}{
							1: {},
						},
					},
				},
			},
			o{"mapOfListToSet": a{
				o{"value": a{int64(1)},
					"key": a{int32(1), int32(2), int32(3)}}},
			},
		},
		{
			"MapOfSetToListOfDouble",
			tc.ContainersOfContainers{
				MapOfSetToListOfDouble: []struct {
					Key   map[int32]struct{}
					Value []float64
				}{
					{
						Key: map[int32]struct{}{
							2: {},
						},
						Value: []float64{1.1, 2.2, 3.3},
					},
					{
						Key: map[int32]struct{}{
							4: {},
						},
						Value: []float64{2.2, 3.3},
					},
				},
			},
			o{"mapOfSetToListOfDouble": a{
				o{"key": a{int32(2)},
					"value": a{float64(1.1), float64(2.2), float64(3.3)}},
				o{"key": a{int32(4)},
					"value": a{float64(2.2), float64(3.3)}}},
			},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.r.MarshalLogObject(mapEncoder)
		assert.Equal(t, tt.v, mapEncoder.Fields)
	}
}

func TestOptOutOfZap(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}

	mapEncoder := zapcore.NewMapObjectEncoder()
	test := ts.ZapOptOutStruct{
		Name:   "foo",
		Optout: "bar",
	}
	test.MarshalLogObject(mapEncoder)
	expected := o{"name": "foo"}
	assert.Equal(t, expected, mapEncoder.Fields)
}

func TestTypedefsZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}

	// test alias of primitive
	mapEncoder := zapcore.NewMapObjectEncoder()
	testState := td.State("hello")
	test1 := td.DefaultPrimitiveTypedef{State: &testState}
	test1.MarshalLogObject(mapEncoder)
	expected1 := o{"state": "hello"}
	assert.Equal(t, expected1, mapEncoder.Fields)

	// test alias of struct
	mapEncoder = zapcore.NewMapObjectEncoder()
	test2 := td.UUID{High: 123, Low: 456}
	test2.MarshalLogObject(mapEncoder)
	expected2 := o{"high": int64(123), "low": int64(456)}
	assert.Equal(t, expected2, mapEncoder.Fields)

	// test alias of list of structs
	mapEncoder = zapcore.NewMapObjectEncoder()
	testUUID := td.UUID(td.I128{High: 123, Low: 456})
	testTimestamp := td.Timestamp(123)
	test3 := td.EventGroup([]*td.Event{
		&td.Event{
			UUID: &testUUID,
			Time: &testTimestamp,
		},
	})
	mapEncoder.AddArray("addTypedefArrayTest", test3)
	expected3 := o{"addTypedefArrayTest": a{o{"uuid": o{"high": int64(123), "low": int64(456)}, "time": int64(123)}}}
	assert.Equal(t, expected3, mapEncoder.Fields)

	// test alias of set
	b64 := func(byteString string) string {
		return base64.StdEncoding.EncodeToString([]byte(byteString))
	}
	mapEncoder = zapcore.NewMapObjectEncoder()
	test4 := td.BinarySet([][]byte{
		[]byte("foo"), {}, []byte("bar"), []byte("baz"),
	})
	mapEncoder.AddArray("addTypedefSetTest", test4)
	// base64
	expected4 := o{"addTypedefSetTest": a{b64("foo"), b64(""), b64("bar"), b64("baz")}}
	assert.Equal(t, expected4, mapEncoder.Fields)

	// test alias of enums
	mapEncoder = zapcore.NewMapObjectEncoder()
	test5 := td.MyEnum(te.EnumWithValuesX)
	test5.MarshalLogObject(mapEncoder)
	expected5 := o{"value": int32(123), "name": "X"}
	assert.Equal(t, expected5, mapEncoder.Fields)

	// test map with typedef of string key
	mapEncoder = zapcore.NewMapObjectEncoder()
	test6 := td.StateMap(map[td.State]int64{
		td.State("foo"): 1,
		td.State("bar"): 2,
	})
	test6.MarshalLogObject(mapEncoder)
	expected6 := o{"foo": int64(1), "bar": int64(2)}
	assert.Equal(t, expected6, mapEncoder.Fields)
}

func TestEnumWithLabelZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}
	type i = int32

	tests := []struct {
		desc string
		p    te.EnumWithLabel
		v    interface{}
	}{
		{
			desc: "with label",
			p:    te.EnumWithLabelUsername,
			v:    o{"name": "surname", "value": i(0)},
		},
		{
			desc: "empty label",
			p:    te.EnumWithLabelSalt,
			v:    o{"name": "SALT", "value": i(2)},
		},
		{
			desc: "unspecified label",
			p:    te.EnumWithLabelSugar,
			v:    o{"name": "SUGAR", "value": i(3)},
		},
		{
			desc: "keyword label",
			p:    te.EnumWithLabelNaive4N1,
			v:    o{"name": "function", "value": i(5)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			enc := zapcore.NewMapObjectEncoder()
			tt.p.MarshalLogObject(enc)
			assert.Equal(t, tt.v, enc.Fields)
		})
	}
}
