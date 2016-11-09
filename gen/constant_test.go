// Copyright (c) 2016 Uber Technologies, Inc.
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

	tk "go.uber.org/thriftrw/gen/testdata/constants"
	tc "go.uber.org/thriftrw/gen/testdata/containers"
	te "go.uber.org/thriftrw/gen/testdata/enums"
	tx "go.uber.org/thriftrw/gen/testdata/exceptions"
	tok "go.uber.org/thriftrw/gen/testdata/other_constants"
	ts "go.uber.org/thriftrw/gen/testdata/structs"
	td "go.uber.org/thriftrw/gen/testdata/typedefs"
	tu "go.uber.org/thriftrw/gen/testdata/unions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConstants(t *testing.T) {
	enumDefaultBaz := te.EnumDefaultBaz

	tests := []struct {
		name string
		i    interface{}
		o    interface{}
	}{
		{
			"primitiveContainers",
			tk.PrimitiveContainers,
			&tc.PrimitiveContainers{
				ListOfInts: []int64{1, 2, 3},
				SetOfStrings: map[string]struct{}{
					"foo": {},
					"bar": {},
				},
				SetOfBytes: map[int8]struct{}{
					1: {},
					2: {},
					3: {},
				},
				MapOfIntToString: map[int32]string{
					1: "1",
					2: "2",
					3: "3",
				},
				MapOfStringToBool: map[string]bool{
					"1": false,
					"2": true,
					"3": true,
				},
			},
		},
		{
			"enumContainers",
			tk.EnumContainers,
			&tc.EnumContainers{
				ListOfEnums: []te.EnumDefault{te.EnumDefaultBar, te.EnumDefaultFoo},
				SetOfEnums: map[te.EnumWithValues]struct{}{
					te.EnumWithValuesX: {},
					te.EnumWithValuesY: {},
				},
				MapOfEnums: map[te.EnumWithDuplicateValues]int32{
					te.EnumWithDuplicateValuesP: 1,
					te.EnumWithDuplicateValuesQ: 2,
				},
			},
		},
		{
			"containersOfContainers",
			tk.ContainersOfContainers,
			&tc.ContainersOfContainers{
				ListOfLists: [][]int32{
					{1, 2, 3},
					{4, 5, 6},
				},
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
				ListOfMaps: []map[int32]int32{
					{
						1: 2,
						3: 4,
						5: 6,
					},
					{
						7:  8,
						9:  10,
						11: 12,
					},
				},
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
				SetOfLists: [][]string{
					{"1", "2", "3"},
					{"4", "5", "6"},
				},
				SetOfMaps: []map[string]string{
					{
						"1": "2",
						"3": "4",
						"5": "6",
					},
					{
						"7":  "8",
						"9":  "10",
						"11": "12",
					},
				},
				MapOfMapToInt: []struct {
					Key   map[string]int32
					Value int64
				}{
					{
						Key:   map[string]int32{"1": 1, "2": 2, "3": 3},
						Value: 100,
					},
					{
						Key:   map[string]int32{"4": 4, "5": 5, "6": 6},
						Value: 200,
					},
				},
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
						Value: []float64{1.2, 3.4},
					},
					{
						Key: map[int32]struct{}{
							4: {},
							5: {},
							6: {},
						},
						Value: []float64{5.6, 7.8},
					},
				},
			},
		},
		{
			"structWithOptionalEnum",
			tk.StructWithOptionalEnum,
			&te.StructWithOptionalEnum{E: &enumDefaultBaz},
		},
		{
			"emptyException",
			tk.EmptyException,
			&tx.EmptyException{},
		},
		{
			"graph",
			tk.Graph,
			&ts.Graph{
				Edges: []*ts.Edge{
					{
						StartPoint: &ts.Point{X: 1, Y: 2},
						EndPoint:   &ts.Point{X: 3, Y: 4},
					},
					{
						StartPoint: &ts.Point{X: 5, Y: 6},
						EndPoint:   &ts.Point{X: 7, Y: 8},
					},
				},
			},
		},
		{
			"arbitraryValue",
			tk.ArbitraryValue,
			&tu.ArbitraryValue{
				ListValue: []*tu.ArbitraryValue{
					{BoolValue: boolp(true)},
					{Int64Value: int64p(2)},
					{StringValue: stringp("hello")},
					{
						MapValue: map[string]*tu.ArbitraryValue{
							"foo": {StringValue: stringp("bar")},
						},
					},
				},
			},
		},
		{
			"lastNode",
			tk.LastNode,
			&ts.Node{Value: 3},
		},
		{
			"node",
			tk.Node,
			&ts.Node{Value: 1, Tail: &ts.List{Value: 2, Tail: &ts.List{Value: 3}}},
		},
		{
			"i128",
			tk.I128,
			&td.I128{High: 1234, Low: 5678},
		},
		{
			"uuid",
			tk.UUID,
			&td.UUID{High: 1234, Low: 5678},
		},
		{
			"beginningOfTime",
			tk.BeginningOfTime,
			td.Timestamp(0),
		},
		{
			"frameGroup",
			tk.FrameGroup,
			td.FrameGroup{
				&ts.Frame{
					TopLeft: &ts.Point{X: 1, Y: 2},
					Size:    &ts.Size{Width: 100, Height: 200},
				},
				&ts.Frame{
					TopLeft: &ts.Point{X: 3, Y: 4},
					Size:    &ts.Size{Width: 300, Height: 400},
				},
			},
		},
		{
			"myEnum",
			tk.MyEnum,
			td.MyEnum(te.EnumWithValuesY),
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.o, tt.i, tt.name)
	}
}

func TestConstantsMutation(t *testing.T) {
	originalX := tok.SomePoint.X

	wireVal, err := tk.Graph.ToWire()
	require.NoError(t, err)
	var g ts.Graph
	err = g.FromWire(wireVal)
	require.NoError(t, err)
	assert.Equal(t, g.Edges[0].StartPoint.X, originalX)

	tok.SomePoint.X += 42.0
	assert.Equal(t, tok.SomePoint.X, originalX+42.0)

	wireVal, err = tk.Graph.ToWire()
	require.NoError(t, err)
	g = ts.Graph{}
	err = g.FromWire(wireVal)
	require.NoError(t, err)
	assert.Equal(t, g.Edges[0].StartPoint.X, originalX)
}
