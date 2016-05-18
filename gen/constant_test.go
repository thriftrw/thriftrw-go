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

	tk "github.com/thriftrw/thriftrw-go/gen/testdata/constants"
	tc "github.com/thriftrw/thriftrw-go/gen/testdata/containers"
	te "github.com/thriftrw/thriftrw-go/gen/testdata/enums"
	tx "github.com/thriftrw/thriftrw-go/gen/testdata/exceptions"
	ts "github.com/thriftrw/thriftrw-go/gen/testdata/structs"
	tu "github.com/thriftrw/thriftrw-go/gen/testdata/unions"

	"github.com/stretchr/testify/assert"
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
					"foo": struct{}{},
					"bar": struct{}{},
				},
				SetOfBytes: map[int8]struct{}{
					1: struct{}{},
					2: struct{}{},
					3: struct{}{},
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
					te.EnumWithValuesX: struct{}{},
					te.EnumWithValuesY: struct{}{},
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
							1: struct{}{},
							2: struct{}{},
							3: struct{}{},
						},
					},
					{
						Key: []int32{4, 5, 6},
						Value: map[int64]struct{}{
							4: struct{}{},
							5: struct{}{},
							6: struct{}{},
						},
					},
				},
				MapOfSetToListOfDouble: []struct {
					Key   map[int32]struct{}
					Value []float64
				}{
					{
						Key: map[int32]struct{}{
							1: struct{}{},
							2: struct{}{},
							3: struct{}{},
						},
						Value: []float64{1.2, 3.4},
					},
					{
						Key: map[int32]struct{}{
							4: struct{}{},
							5: struct{}{},
							6: struct{}{},
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
					&ts.Edge{
						Start: &ts.Point{X: 1, Y: 2},
						End:   &ts.Point{X: 3, Y: 4},
					},
					&ts.Edge{
						Start: &ts.Point{X: 5, Y: 6},
						End:   &ts.Point{X: 7, Y: 8},
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
							"foo": &tu.ArbitraryValue{StringValue: stringp("bar")},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.o, tt.i, tt.name)
	}
}
