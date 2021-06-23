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
	"testing"

	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	ts "go.uber.org/thriftrw/gen/internal/tests/structs"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectionsOfPrimitives(t *testing.T) {
	tests := []struct {
		desc string
		p    tc.PrimitiveContainers
		v    wire.Value
	}{
		// Lists /////////////////////////////////////////////////////////////
		{
			"empty list",
			tc.PrimitiveContainers{ListOfInts: []int64{}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID:    2,
				Value: wire.NewValueList(wire.ValueListFromSlice(wire.TI64, []wire.Value{})),
			}}}),
		},
		{
			"list of ints",
			tc.PrimitiveContainers{ListOfInts: []int64{1, 2, 3}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 2,
				Value: wire.NewValueList(
					wire.ValueListFromSlice(wire.TI64, []wire.Value{
						wire.NewValueI64(1),
						wire.NewValueI64(2),
						wire.NewValueI64(3),
					}),
				),
			}}}),
		},
		{
			"list of binary",
			tc.PrimitiveContainers{
				ListOfBinary: [][]byte{
					[]byte("foo"), {}, []byte("bar"), []byte("baz"),
				},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 1,
				Value: wire.NewValueList(
					wire.ValueListFromSlice(wire.TBinary, []wire.Value{
						wire.NewValueBinary([]byte("foo")),
						wire.NewValueBinary([]byte{}),
						wire.NewValueBinary([]byte("bar")),
						wire.NewValueBinary([]byte("baz")),
					}),
				),
			}}}),
		},
		// Sets //////////////////////////////////////////////////////////////
		{
			"empty set",
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 3,
				Value: wire.NewValueSet(
					wire.ValueListFromSlice(wire.TBinary, []wire.Value{}),
				),
			}}}),
		},
		{
			"set of strings",
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{
				"foo": {},
				"bar": {},
				"baz": {},
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 3,
				Value: wire.NewValueSet(
					wire.ValueListFromSlice(wire.TBinary, []wire.Value{
						wire.NewValueString("foo"),
						wire.NewValueString("bar"),
						wire.NewValueString("baz"),
					}),
				),
			}}}),
		},
		{
			"set of bytes",
			tc.PrimitiveContainers{SetOfBytes: map[int8]struct{}{
				-1:  {},
				1:   {},
				125: {},
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 4,
				Value: wire.NewValueSet(
					wire.ValueListFromSlice(wire.TI8, []wire.Value{
						wire.NewValueI8(-1),
						wire.NewValueI8(1),
						wire.NewValueI8(125),
					}),
				),
			}}}),
		},
		// Maps //////////////////////////////////////////////////////////////
		{
			"empty map",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 6,
				Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TBinary, wire.TBool, []wire.MapItem{}),
				),
			}}}),
		},
		{
			"map of int to string",
			tc.PrimitiveContainers{MapOfIntToString: map[int32]string{
				-1:    "foo",
				1234:  "bar",
				-9876: "baz",
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 5,
				Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TI32, wire.TBinary, []wire.MapItem{
						{Key: wire.NewValueI32(-1), Value: wire.NewValueString("foo")},
						{Key: wire.NewValueI32(1234), Value: wire.NewValueString("bar")},
						{Key: wire.NewValueI32(-9876), Value: wire.NewValueString("baz")},
					}),
				),
			}}}),
		},
		{
			"map of string to bool",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{
				"foo": true,
				"bar": false,
				"baz": true,
			}},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
				ID: 6,
				Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TBinary, wire.TBool, []wire.MapItem{
						{Key: wire.NewValueString("foo"), Value: wire.NewValueBool(true)},
						{Key: wire.NewValueString("bar"), Value: wire.NewValueBool(false)},
						{Key: wire.NewValueString("baz"), Value: wire.NewValueBool(true)},
					}),
				),
			}}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.p, tt.v, tt.desc)
		assert.True(t, tt.p.Equals(&tt.p), tt.desc)

		testRoundTripCombos(t, &tt.p, tt.v, tt.desc)
		assert.True(t, tt.p.Equals(&tt.p), tt.desc)
	}
}

func TestCollectionsOfPrimitivesEquals(t *testing.T) {
	tests := []struct {
		desc string
		p, q tc.PrimitiveContainers
	}{
		// Lists /////////////////////////////////////////////////////////////
		{
			"empty list",
			tc.PrimitiveContainers{ListOfInts: []int64{}},
			tc.PrimitiveContainers{ListOfInts: []int64{1, 2}},
		},
		{
			"list of ints",
			tc.PrimitiveContainers{ListOfInts: []int64{1, 2, 3}},
			tc.PrimitiveContainers{ListOfInts: []int64{1, 4, 3}},
		},
		{
			"list of binary",
			tc.PrimitiveContainers{
				ListOfBinary: [][]byte{
					[]byte("foo"), {}, []byte("bar"), []byte("baz"),
				},
			},
			tc.PrimitiveContainers{
				ListOfBinary: [][]byte{
					[]byte("foo"), {}, []byte("bar"), []byte("bazzinga"),
				},
			},
		},
		// Sets //////////////////////////////////////////////////////////////
		{
			"empty set",
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{}},
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{"foo": {}}},
		},
		{
			"set of strings",
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{
				"foo": {},
				"bar": {},
				"baz": {},
			}},
			tc.PrimitiveContainers{SetOfStrings: map[string]struct{}{
				"foobar": {},
			}},
		},
		{
			"set of bytes",
			tc.PrimitiveContainers{SetOfBytes: map[int8]struct{}{
				-1:  {},
				1:   {},
				125: {},
			}},
			tc.PrimitiveContainers{SetOfBytes: map[int8]struct{}{
				-1:  {},
				2:   {},
				125: {},
			}},
		},
		// Maps //////////////////////////////////////////////////////////////
		{
			"empty map",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{}},
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{"foo": false}},
		},
		{
			"map of int to string",
			tc.PrimitiveContainers{MapOfIntToString: map[int32]string{
				-1:    "foo",
				1234:  "bar",
				-9876: "baz",
			}},
			tc.PrimitiveContainers{MapOfIntToString: map[int32]string{
				-1:    "foo",
				1234:  "bar",
				-9876: "bazzinga",
			}},
		},
		{
			"map of string to bool",
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{
				"foo": true,
				"bar": false,
				"baz": true,
			}},
			tc.PrimitiveContainers{MapOfStringToBool: map[string]bool{
				"foo":      true,
				"bazzinga": true,
				"bar":      false,
			}},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.p.Equals(&tt.p), tt.desc)
		assert.False(t, tt.p.Equals(&tt.q), tt.desc)
	}
}

func TestEnumContainers(t *testing.T) {
	tests := []struct {
		r tc.EnumContainers
		v wire.Value
	}{
		{
			tc.EnumContainers{
				ListOfEnums: []te.EnumDefault{
					te.EnumDefaultFoo,
					te.EnumDefaultBar,
				},
			},
			singleFieldStruct(1, wire.NewValueList(
				wire.ValueListFromSlice(wire.TI32, []wire.Value{
					wire.NewValueI32(0),
					wire.NewValueI32(1),
				}),
			)),
		},
		{
			tc.EnumContainers{
				SetOfEnums: map[te.EnumWithValues]struct{}{
					te.EnumWithValuesX: {},
					te.EnumWithValuesZ: {},
				},
			},
			singleFieldStruct(2, wire.NewValueSet(
				wire.ValueListFromSlice(wire.TI32, []wire.Value{
					wire.NewValueI32(123),
					wire.NewValueI32(789),
				}),
			)),
		},
		{
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesP: 123,
					te.EnumWithDuplicateValuesQ: 456,
				},
			},
			singleFieldStruct(3, wire.NewValueMap(
				wire.MapItemListFromSlice(wire.TI32, wire.TI32, []wire.MapItem{
					{Key: wire.NewValueI32(0), Value: wire.NewValueI32(123)},
					{Key: wire.NewValueI32(-1), Value: wire.NewValueI32(456)},
				}),
			)),
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
			singleFieldStruct(3, wire.NewValueMap(

				wire.MapItemListFromSlice(wire.TI32, wire.TI32, []wire.MapItem{
					{Key: wire.NewValueI32(0), Value: wire.NewValueI32(123)},
					{Key: wire.NewValueI32(-1), Value: wire.NewValueI32(456)},
				}),
			)),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.r, tt.v, "EnumContainers")
		assert.True(t, tt.r.Equals(&tt.r), "EnumContainers equal")

		testRoundTripCombos(t, &tt.r, tt.v, "EnumContainers")
		assert.True(t, tt.r.Equals(&tt.r), "EnumContainers equal")
	}
}

func TestEnumContainersEquals(t *testing.T) {
	tests := []struct {
		r, s tc.EnumContainers
	}{
		{
			tc.EnumContainers{
				ListOfEnums: []te.EnumDefault{
					te.EnumDefaultFoo,
					te.EnumDefaultBar,
				},
			},
			tc.EnumContainers{
				ListOfEnums: []te.EnumDefault{
					te.EnumDefaultFoo,
				},
			},
		},
		{
			tc.EnumContainers{
				SetOfEnums: map[te.EnumWithValues]struct{}{
					te.EnumWithValuesX: {},
					te.EnumWithValuesZ: {},
				},
			},
			tc.EnumContainers{
				SetOfEnums: map[te.EnumWithValues]struct{}{
					te.EnumWithValuesX: {},
					te.EnumWithValuesY: {},
				},
			},
		},
		{
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesP: 123,
					te.EnumWithDuplicateValuesQ: 456,
				},
			},
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesP: 123,
					te.EnumWithDuplicateValuesQ: 789,
				},
			},
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
			tc.EnumContainers{
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesQ: 456,
				},
			},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.r.Equals(&tt.r), "EnumContainers equal")
		assert.False(t, tt.r.Equals(&tt.s), "EnumContainers unequal")
	}
}

func TestListOfStructs(t *testing.T) {
	tests := []struct {
		r ts.Graph
		v wire.Value
	}{
		{
			ts.Graph{Edges: []*ts.Edge{}},
			singleFieldStruct(1, wire.NewValueList(
				wire.ValueListFromSlice(wire.TStruct, nil),
			)),
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
			singleFieldStruct(1, wire.NewValueList(
				wire.ValueListFromSlice(wire.TStruct, []wire.Value{
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
								{ID: 1, Value: wire.NewValueDouble(3.0)},
								{ID: 2, Value: wire.NewValueDouble(4.0)},
							}}),
						},
					}}),
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{
							ID: 1,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueDouble(5.0)},
								{ID: 2, Value: wire.NewValueDouble(6.0)},
							}}),
						},
						{
							ID: 2,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueDouble(7.0)},
								{ID: 2, Value: wire.NewValueDouble(8.0)},
							}}),
						},
					}}),
					wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{
							ID: 1,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueDouble(9.0)},
								{ID: 2, Value: wire.NewValueDouble(10.0)},
							}}),
						},
						{
							ID: 2,
							Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
								{ID: 1, Value: wire.NewValueDouble(11.0)},
								{ID: 2, Value: wire.NewValueDouble(12.0)},
							}}),
						},
					}}),
				}),
			)),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.r, tt.v, "Graph")
		assert.True(t, tt.r.Equals(&tt.r), "Graph equal")

		testRoundTripCombos(t, &tt.r, tt.v, "Graph")
		assert.True(t, tt.r.Equals(&tt.r), "Graph equal")
	}
}

func TestListOfStructsEquals(t *testing.T) {
	tests := []struct {
		r, s ts.Graph
	}{
		{
			ts.Graph{Edges: []*ts.Edge{}},
			ts.Graph{Edges: []*ts.Edge{
				{
					StartPoint: &ts.Point{X: 1.0, Y: 2.0},
					EndPoint:   &ts.Point{X: 3.0, Y: 4.0},
				},
			}},
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
					EndPoint:   &ts.Point{X: 999.0, Y: 1000.0},
				},
			}},
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.r.Equals(&tt.r), "Graph equal")
		assert.False(t, tt.r.Equals(&tt.s), "Graph unequal")
	}
}

func TestCrazyTown(t *testing.T) {
	tests := []struct {
		desc string
		x    tc.ContainersOfContainers
		v    wire.Value
	}{
		{
			"ListOfLists",
			tc.ContainersOfContainers{
				ListOfLists: [][]int32{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueList(
					wire.ValueListFromSlice(wire.TList, []wire.Value{
						wire.NewValueList(
							wire.ValueListFromSlice(wire.TI32, []wire.Value{
								wire.NewValueI32(1),
								wire.NewValueI32(2),
								wire.NewValueI32(3),
							}),
						),
						wire.NewValueList(
							wire.ValueListFromSlice(wire.TI32, []wire.Value{
								wire.NewValueI32(4),
								wire.NewValueI32(5),
								wire.NewValueI32(6),
							}),
						),
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 2, Value: wire.NewValueList(
					wire.ValueListFromSlice(wire.TSet, []wire.Value{
						wire.NewValueSet(
							wire.ValueListFromSlice(wire.TI32, []wire.Value{
								wire.NewValueI32(1),
								wire.NewValueI32(2),
								wire.NewValueI32(3),
							}),
						),
						wire.NewValueSet(
							wire.ValueListFromSlice(wire.TI32, []wire.Value{
								wire.NewValueI32(4),
								wire.NewValueI32(5),
								wire.NewValueI32(6),
							}),
						),
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 3, Value: wire.NewValueList(
					wire.ValueListFromSlice(wire.TMap, []wire.Value{
						wire.NewValueMap(
							wire.MapItemListFromSlice(wire.TI32, wire.TI32, []wire.MapItem{
								{Key: wire.NewValueI32(1), Value: wire.NewValueI32(100)},
								{Key: wire.NewValueI32(2), Value: wire.NewValueI32(200)},
								{Key: wire.NewValueI32(3), Value: wire.NewValueI32(300)},
							}),
						),
						wire.NewValueMap(
							wire.MapItemListFromSlice(wire.TI32, wire.TI32, []wire.MapItem{
								{Key: wire.NewValueI32(4), Value: wire.NewValueI32(400)},
								{Key: wire.NewValueI32(5), Value: wire.NewValueI32(500)},
								{Key: wire.NewValueI32(6), Value: wire.NewValueI32(600)},
							}),
						),
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 4, Value: wire.NewValueSet(
					wire.ValueListFromSlice(wire.TSet, []wire.Value{
						wire.NewValueSet(
							wire.ValueListFromSlice(wire.TBinary, []wire.Value{
								wire.NewValueString("1"),
								wire.NewValueString("2"),
								wire.NewValueString("3"),
							}),
						),
						wire.NewValueSet(
							wire.ValueListFromSlice(wire.TBinary, []wire.Value{
								wire.NewValueString("4"),
								wire.NewValueString("5"),
								wire.NewValueString("6"),
							}),
						),
					}),
				)},
			}}),
		},
		{
			"SetOfLists",
			tc.ContainersOfContainers{
				SetOfLists: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
				},
			},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 5, Value: wire.NewValueSet(
					wire.ValueListFromSlice(wire.TList, []wire.Value{
						wire.NewValueList(
							wire.ValueListFromSlice(wire.TBinary, []wire.Value{
								wire.NewValueString("1"),
								wire.NewValueString("2"),
								wire.NewValueString("3"),
							}),
						),
						wire.NewValueList(
							wire.ValueListFromSlice(wire.TBinary, []wire.Value{
								wire.NewValueString("4"),
								wire.NewValueString("5"),
								wire.NewValueString("6"),
							}),
						),
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 6, Value: wire.NewValueSet(
					wire.ValueListFromSlice(wire.TMap, []wire.Value{
						wire.NewValueMap(
							wire.MapItemListFromSlice(wire.TBinary, wire.TBinary, []wire.MapItem{
								{Key: wire.NewValueString("1"), Value: wire.NewValueString("one")},
								{Key: wire.NewValueString("2"), Value: wire.NewValueString("two")},
								{Key: wire.NewValueString("3"), Value: wire.NewValueString("three")},
							}),
						),
						wire.NewValueMap(
							wire.MapItemListFromSlice(wire.TBinary, wire.TBinary, []wire.MapItem{
								{Key: wire.NewValueString("4"), Value: wire.NewValueString("four")},
								{Key: wire.NewValueString("5"), Value: wire.NewValueString("five")},
								{Key: wire.NewValueString("6"), Value: wire.NewValueString("six")},
							}),
						),
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 7, Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TMap, wire.TI64, []wire.MapItem{
						{
							Key: wire.NewValueMap(
								wire.MapItemListFromSlice(wire.TBinary, wire.TI32, []wire.MapItem{
									{Key: wire.NewValueString("1"), Value: wire.NewValueI32(1)},
									{Key: wire.NewValueString("2"), Value: wire.NewValueI32(2)},
									{Key: wire.NewValueString("3"), Value: wire.NewValueI32(3)},
								}),
							),
							Value: wire.NewValueI64(123),
						},
						{
							Key: wire.NewValueMap(
								wire.MapItemListFromSlice(wire.TBinary, wire.TI32, []wire.MapItem{
									{Key: wire.NewValueString("4"), Value: wire.NewValueI32(4)},
									{Key: wire.NewValueString("5"), Value: wire.NewValueI32(5)},
									{Key: wire.NewValueString("6"), Value: wire.NewValueI32(6)},
								}),
							),
							Value: wire.NewValueI64(456),
						},
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 8, Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TList, wire.TSet, []wire.MapItem{
						{
							Key: wire.NewValueList(
								wire.ValueListFromSlice(wire.TI32, []wire.Value{
									wire.NewValueI32(1),
									wire.NewValueI32(2),
									wire.NewValueI32(3),
								}),
							),
							Value: wire.NewValueSet(
								wire.ValueListFromSlice(wire.TI64, []wire.Value{
									wire.NewValueI64(1),
									wire.NewValueI64(2),
									wire.NewValueI64(3),
								}),
							),
						},
						{
							Key: wire.NewValueList(
								wire.ValueListFromSlice(wire.TI32, []wire.Value{
									wire.NewValueI32(4),
									wire.NewValueI32(5),
									wire.NewValueI32(6),
								}),
							),
							Value: wire.NewValueSet(
								wire.ValueListFromSlice(wire.TI64, []wire.Value{
									wire.NewValueI64(4),
									wire.NewValueI64(5),
									wire.NewValueI64(6),
								}),
							),
						},
					}),
				)},
			}}),
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
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 9, Value: wire.NewValueMap(
					wire.MapItemListFromSlice(wire.TSet, wire.TList, []wire.MapItem{
						{
							Key: wire.NewValueSet(
								wire.ValueListFromSlice(wire.TI32, []wire.Value{
									wire.NewValueI32(1),
									wire.NewValueI32(2),
									wire.NewValueI32(3),
								}),
							),
							Value: wire.NewValueList(
								wire.ValueListFromSlice(wire.TDouble, []wire.Value{
									wire.NewValueDouble(1.0),
									wire.NewValueDouble(2.0),
									wire.NewValueDouble(3.0),
								}),
							),
						},
						{
							Key: wire.NewValueSet(
								wire.ValueListFromSlice(wire.TI32, []wire.Value{
									wire.NewValueI32(4),
									wire.NewValueI32(5),
									wire.NewValueI32(6),
								}),
							),
							Value: wire.NewValueList(
								wire.ValueListFromSlice(wire.TDouble, []wire.Value{
									wire.NewValueDouble(4.0),
									wire.NewValueDouble(5.0),
									wire.NewValueDouble(6.0),
								}),
							),
						},
					}),
				)},
			}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.x, tt.v, tt.desc)
		assert.True(t, tt.x.Equals(&tt.x), tt.desc)

		testRoundTripCombos(t, &tt.x, tt.v, tt.desc)
		assert.True(t, tt.x.Equals(&tt.x), tt.desc)
	}
}

func TestCrazyTownEquals(t *testing.T) {
	tests := []struct {
		desc string
		x, y tc.ContainersOfContainers
	}{
		{
			"ListOfLists",
			tc.ContainersOfContainers{
				ListOfLists: [][]int32{
					{1, 2, 3},
					{4, 5, 6},
				},
			},
			tc.ContainersOfContainers{
				ListOfLists: [][]int32{
					{1, 2, 3},
					{4, 5, 6},
					{7, 8, 9},
				},
			},
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
			tc.ContainersOfContainers{
				ListOfSets: []map[int32]struct{}{
					{
						1:  struct{}{},
						2:  struct{}{},
						30: struct{}{},
					},
					{
						40: struct{}{},
						50: struct{}{},
						60: struct{}{},
					},
				},
			},
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
			tc.ContainersOfContainers{
				ListOfMaps: []map[int32]int32{
					{
						1: 100,
						2: 200,
						3: 300,
					},
					{
						4: 400,
						5: 500000,
						6: 600,
					},
				},
			},
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
					},
				},
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
			tc.ContainersOfContainers{
				SetOfLists: [][]string{
					{"1", "2", "3"},
					{"4", "500", "6"},
				},
			},
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
			tc.ContainersOfContainers{
				SetOfMaps: []map[string]string{
					{
						"1": "one",
						"2": "two",
						"3": "three",
					},
					{
						"4": "four",
						"5": "fiftyfive",
						"6": "six",
					},
				},
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
						Key:   map[string]int32{"1": 1, "2": 2, "3": 3},
						Value: 123,
					},
					{
						Key:   map[string]int32{"4": 4, "5": 5, "6": 6},
						Value: 456,
					},
				},
			},
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
						Key:   map[string]int32{"4": 4, "55": 5, "6": 6},
						Value: 456,
					},
				},
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
							404: {},
							5:   {},
							6:   {},
						},
					},
				},
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
						Value: []float64{1.0, 3.0},
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
		},
	}

	for _, tt := range tests {
		assert.True(t, tt.x.Equals(&tt.x), tt.desc)
		assert.False(t, tt.x.Equals(&tt.y), tt.desc)
	}
}

func TestContainerValidate(t *testing.T) {
	tests := []struct {
		desc      string
		value     thriftType
		wantError string
	}{
		{
			desc: "nil byte sub-array",
			value: &tc.PrimitiveContainers{
				ListOfBinary: [][]byte{
					{1, 2, 3},
					{},
					nil,
					{4, 5, 6},
				},
			},
			wantError: "invalid [2]: value is nil",
		},
		{
			desc: "nil second element of array",
			value: &ts.Graph{
				Edges: []*ts.Edge{
					{StartPoint: &ts.Point{X: 1, Y: 2}, EndPoint: &ts.Point{X: 3, Y: 4}},
					nil,
					{StartPoint: &ts.Point{X: 5, Y: 6}, EndPoint: &ts.Point{X: 7, Y: 8}},
				},
			},
			wantError: "invalid [1]: value is nil",
		},
		{
			desc: "nil set item",
			value: &tc.ContainersOfContainers{
				SetOfLists: [][]string{{}, nil},
			},
			wantError: "invalid set item: value is nil",
		},
		{
			desc: "nil map key",
			value: &tc.MapOfBinaryAndString{
				BinaryToString: []struct {
					Key   []byte
					Value string
				}{
					{Key: []byte("hello"), Value: "world"},
					{Key: nil, Value: "foo"},
				},
			},
			wantError: "invalid map key: value is nil",
		},
		{
			desc: "nil map value",
			value: &tc.MapOfBinaryAndString{
				StringToBinary: map[string][]byte{
					"hello": []byte("world"),
					"foo":   nil,
				},
			},
			wantError: "invalid [foo]: value is nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			value, err := tt.value.ToWire()
			if err == nil {
				err = wire.EvaluateValue(value) // lazy error
			}

			if assert.Error(t, err) {
				assert.Equal(t, tt.wantError, err.Error())
			}
		})
	}
}

func TestListOfBinaryReadNil(t *testing.T) {
	value := wire.NewValueStruct(wire.Struct{Fields: []wire.Field{{
		ID: 1,
		Value: wire.NewValueList(
			wire.ValueListFromSlice(wire.TBinary, []wire.Value{
				wire.NewValueBinary([]byte("foo")),
				wire.NewValueBinary(nil),
				wire.NewValueBinary([]byte("bar")),
				wire.NewValueBinary([]byte("baz")),
			}),
		),
	}}})

	var c tc.PrimitiveContainers
	require.NoError(t, c.FromWire(value))

	got, err := c.ToWire()
	require.NoError(t, err)
	require.NoError(t, wire.EvaluateValue(got))
	assert.True(t, wire.ValuesAreEqual(value, got))
}

func TestEmptyContainersRoundTrip(t *testing.T) {
	t.Run("required", func(t *testing.T) {
		give := tc.PrimitiveContainersRequired{
			ListOfStrings:      []string{},
			SetOfInts:          make(map[int32]struct{}),
			MapOfIntsToDoubles: make(map[int64]float64),
		}

		b, err := json.Marshal(give)
		require.NoError(t, err, "failed to encode to JSON")

		var decoded tc.PrimitiveContainersRequired
		require.NoError(t, json.Unmarshal(b, &decoded), "failed to decode JSON")

		assert.Equal(t, give, decoded)

		v, err := decoded.ToWire()
		require.NoError(t, err, "failed to convert to wire.Value")

		var got tc.PrimitiveContainersRequired
		require.NoError(t, got.FromWire(v), "failed to convert from wire.Value")

		assert.Equal(t, give, got)
	})

	t.Run("optional", func(t *testing.T) {
		give := tc.PrimitiveContainers{
			ListOfInts:       []int64{},
			SetOfStrings:     make(map[string]struct{}),
			MapOfIntToString: make(map[int32]string),
		}

		b, err := json.Marshal(give)
		require.NoError(t, err, "failed to encode to JSON")

		var decoded tc.PrimitiveContainers
		require.NoError(t, json.Unmarshal(b, &decoded), "failed to decode JSON")

		// We check individual fields because a full assert.Equal could mismatch
		// on nil vs empty slice.
		assert.Empty(t, decoded.ListOfInts)
		assert.Empty(t, decoded.SetOfStrings)
		assert.Empty(t, decoded.MapOfIntToString)

		v, err := decoded.ToWire()
		require.NoError(t, err, "failed to convert to wire.Value")

		var got tc.PrimitiveContainers
		require.NoError(t, got.FromWire(v), "failed to convert from wire.Value")

		assert.Empty(t, got.ListOfInts)
		assert.Empty(t, got.SetOfStrings)
		assert.Empty(t, got.MapOfIntToString)
	})
}
