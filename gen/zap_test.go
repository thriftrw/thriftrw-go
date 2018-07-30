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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	"go.uber.org/zap/zapcore"
)

func jsonToComparableMap(jsonMap string) (map[string]interface{}, error) {
	var retMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonMap), &retMap); err != nil {
		return nil, err
	}
	return retMap, nil
}

func mapTojson(mapThing map[string]interface{}) (string, error) {
	retBytes, err := json.Marshal(mapThing)
	if err != nil {
		return "", err
	}
	return string(retBytes), nil
}

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
				"bar": {},
				"baz": {},
			}},
			o{"setOfStrings": a{"foo", "bar", "baz"}},
		},
		{
			"set of bytes",
			tc.PrimitiveContainers{SetOfBytes: map[int8]struct{}{
				-1:  {},
				1:   {},
				125: {},
			}},
			o{"setOfBytes": a{int8(-1), int8(1), int8(125)}},
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
				-1:    "foo",
				1234:  "bar",
				-9876: "baz",
			}},
			o{"mapOfIntToString": a{
				o{"key": int32(-1), "value": "foo"},
				o{"key": int32(1234), "value": "bar"},
				o{"key": int32(-9876), "value": "baz"},
			}},
		},
		{
			"map of string to bool",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{
				"foo": true,
				"bar": false,
				"baz": true,
			}},
			o{"mapOfStringToBool": o{"bar": false, "baz": true, "foo": true}},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.p.MarshalLogObject(mapEncoder)
		// mapEncoder.AddObject("actual", &tt.p)
		// t.Log(mapTojson(mapEncoder.Fields))
		assert.Equal(t, tt.v, mapEncoder.Fields)
		// DeepEquals mapEncoder.Fields and tt.v
	}
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
					te.EnumDefaultFoo,
					te.EnumDefaultBar,
				},
			},
			o{"listOfEnums": a{
				o{"value": int32(0), "name": "Foo"},
				o{"value": int32(1), "name": "Bar"}},
			},
		},
		{
			tc.EnumContainers{
				SetOfEnums: map[te.EnumWithValues]struct{}{
					te.EnumWithValuesX: {},
					te.EnumWithValuesZ: {},
				},
			},
			o{"setOfEnums": a{
				o{"name": "X", "value": 123},
				o{"name": "Z", "value": 789}}},
		},
		{
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesP: 123,
					te.EnumWithDuplicateValuesQ: 456,
				},
			},
			o{"mapOfEnums": a{
				o{"key": o{"value": 0, "name": "P"}, "value": 123},
				o{"key": o{"name": "Q", "value": -1}, "value": 456}}},
		},
		{
			// this is the same as the one above except we're using "R" intsead
			// of "P" (they both have the same value)
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesR: 123,
					te.EnumWithDuplicateValuesQ: 456,
				},
			},
			o{"mapOfEnums": a{
				o{"key": o{"value": -1, "name": "Q"}, "value": 456},
				o{"key": o{"value": 0, "name": "P"}, "value": 123}}},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.r.MarshalLogObject(mapEncoder)
		// mapEncoder.AddObject("actual", &tt.p)
		// t.Log(mapTojson(mapEncoder.Fields))
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
					"endPoint": o{"x": float64(11), "y": float64(12)}}},
			},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.r.MarshalLogObject(mapEncoder)
		// mapEncoder.AddObject("actual", &tt.p)
		// t.Log(mapTojson(mapEncoder.Fields))
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
			o{},
		},
		{
			"ListOfSets",
			tc.ContainersOfContainers{
				ListOfSets: []map[int32]struct{}{
					{
						1: struct{}{},
						2: struct{}{},
						3: struct{}{},
					},
					{
						4: struct{}{},
						5: struct{}{},
						6: struct{}{},
					},
				},
			},
			o{},
		},
		{
			"ListOfMaps",
			tc.ContainersOfContainers{
				ListOfMaps: []map[int32]int32{
					{
						1: 100,
						2: 200,
						3: 300,
					},
					{
						4: 400,
						5: 500,
						6: 600,
					},
				},
			},
			o{},
		},
		{
			"SetOfSets",
			tc.ContainersOfContainers{
				SetOfSets: []map[string]struct{}{
					{
						"1": struct{}{},
						"2": struct{}{},
						"3": struct{}{},
					},
					{
						"4": struct{}{},
						"5": struct{}{},
						"6": struct{}{},
					},
				},
			},
			o{},
		},
		{
			"SetOfLists",
			tc.ContainersOfContainers{
				SetOfLists: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
				},
			},
			o{},
		},
		{
			"SetOfMaps",
			tc.ContainersOfContainers{
				SetOfMaps: []map[string]string{
					{
						"1": "one",
						"2": "two",
						"3": "three",
					},
					{
						"4": "four",
						"5": "five",
						"6": "six",
					},
				},
			},
			o{},
		},
		{
			"MapOfMapToInt",
			tc.ContainersOfContainers{
				MapOfMapToInt: []struct {
					Key   map[string]int32
					Value int64
				}{
					{
						Key:   map[string]int32{"1": 1, "2": 2, "3": 3},
						Value: 123,
					},
					{
						Key:   map[string]int32{"4": 4, "5": 5, "6": 6},
						Value: 456,
					},
				},
			},
			o{},
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
							2: {},
							3: {},
						},
					},
					{
						Key: []int32{4, 5, 6},
						Value: map[int64]struct{}{
							4: {},
							5: {},
							6: {},
						},
					},
				},
			},
			o{},
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
							1: {},
							2: {},
							3: {},
						},
						Value: []float64{1.0, 2.0, 3.0},
					},
					{
						Key: map[int32]struct{}{
							4: {},
							5: {},
							6: {},
						},
						Value: []float64{4.0, 5.0, 6.0},
					},
				},
			},
			o{},
		},
	}

	for _, tt := range tests {
		mapEncoder := zapcore.NewMapObjectEncoder()
		tt.r.MarshalLogObject(mapEncoder)
		// mapEncoder.AddObject("actual", &tt.p)
		// t.Log(mapTojson(mapEncoder.Fields))
		assert.Equal(t, tt.v, mapEncoder.Fields)
	}
}
