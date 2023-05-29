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
	"fmt"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tc "go.uber.org/thriftrw/gen/internal/tests/containers"
	te "go.uber.org/thriftrw/gen/internal/tests/enums"
	tz "go.uber.org/thriftrw/gen/internal/tests/nozap"
	tss "go.uber.org/thriftrw/gen/internal/tests/set_to_slice"
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
		require.NoError(t, tt.p.MarshalLogObject(mapEncoder))
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
	require.NoError(t, test1.MarshalLogObject(mapEncoder))
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
	require.NoError(t, test2.MarshalLogObject(mapEncoder))
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
	require.NoError(t, test3.MarshalLogObject(mapEncoder))
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
		require.NoError(t, tt.r.MarshalLogObject(mapEncoder))
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
		require.NoError(t, tt.r.MarshalLogObject(mapEncoder))
		assert.Equal(t, tt.v, mapEncoder.Fields)
	}
}

func TestStructWithFieldTypeAnnotatedSetToSlice(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
	type a = []interface{}

	s := tss.Bar{
		RequiredInt32ListField:         []int32{1},
		OptionalStringListField:        []string{"foo"},
		RequiredTypedefStringListField: tss.StringList{"foo"},
		OptionalTypedefStringListField: tss.StringList{"foo"},
		RequiredFooListField: []*tss.Foo{
			{
				StringField: "foo",
			},
		},
		RequiredTypedefFooListField: tss.FooList{
			{
				StringField: "foo",
			},
		},
		RequiredStringListListField:        [][]string{{"foo"}},
		RequiredTypedefStringListListField: [][]string{{"foo"}},
	}
	v := o{
		"requiredInt32ListField":         a{int32(1)},
		"optionalStringListField":        a{"foo"},
		"requiredTypedefStringListField": a{"foo"},
		"optionalTypedefStringListField": a{"foo"},
		"requiredFooListField": a{
			o{
				"stringField": "foo",
			},
		},
		"requiredTypedefFooListField": a{
			o{
				"stringField": "foo",
			},
		},
		"requiredStringListListField":        a{a{"foo"}},
		"requiredTypedefStringListListField": a{a{"foo"}},
	}
	mapEncoder := zapcore.NewMapObjectEncoder()
	require.NoError(t, s.MarshalLogObject(mapEncoder))
	assert.Equal(t, v, mapEncoder.Fields)
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
		require.NoError(t, tt.r.MarshalLogObject(mapEncoder))
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
	require.NoError(t, test.MarshalLogObject(mapEncoder))
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
	require.NoError(t, test1.MarshalLogObject(mapEncoder))
	expected1 := o{"state": "hello"}
	assert.Equal(t, expected1, mapEncoder.Fields)

	// test alias of struct
	mapEncoder = zapcore.NewMapObjectEncoder()
	test2 := td.UUID{High: 123, Low: 456}
	require.NoError(t, test2.MarshalLogObject(mapEncoder))
	expected2 := o{"high": int64(123), "low": int64(456)}
	assert.Equal(t, expected2, mapEncoder.Fields)

	// test alias of list of structs
	mapEncoder = zapcore.NewMapObjectEncoder()
	testUUID := td.UUID(td.I128{High: 123, Low: 456})
	testTimestamp := td.Timestamp(123)
	test3 := td.EventGroup([]*td.Event{
		{
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
	require.NoError(t, test5.MarshalLogObject(mapEncoder))
	expected5 := o{"value": int32(123), "name": "X"}
	assert.Equal(t, expected5, mapEncoder.Fields)

	// test map with typedef of string key
	mapEncoder = zapcore.NewMapObjectEncoder()
	test6 := td.StateMap(map[td.State]int64{
		td.State("foo"): 1,
		td.State("bar"): 2,
	})
	require.NoError(t, test6.MarshalLogObject(mapEncoder))
	expected6 := o{"foo": int64(1), "bar": int64(2)}
	assert.Equal(t, expected6, mapEncoder.Fields)

	// test set annotated with (go.type = "slice")
	mapEncoder = zapcore.NewMapObjectEncoder()
	test7 := tss.StringList{"foo"}
	expected7 := o{"addTypedefSetToSliceTest": a{"foo"}}
	err := mapEncoder.AddArray("addTypedefSetToSliceTest", test7)
	require.NoError(t, err)
	assert.Equal(t, expected7, mapEncoder.Fields)

	mapEncoder = zapcore.NewMapObjectEncoder()
	test8 := tss.MyStringList{"foo"}
	expected8 := o{"addTypedefSetToSliceTest": a{"foo"}}
	err = mapEncoder.AddArray("addTypedefSetToSliceTest", test8)
	require.NoError(t, err)
	assert.Equal(t, expected8, mapEncoder.Fields)
}

func TestEnumWithLabelZapLogging(t *testing.T) {
	// These types are created to ease building map[string]interface{}
	type o = map[string]interface{}
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
			require.NoError(t, tt.p.MarshalLogObject(enc))
			assert.Equal(t, tt.v, enc.Fields)
		})
	}
}

func TestNoZapLogging(t *testing.T) {
	tests := []struct {
		desc string
		p    interface{}
	}{
		{
			desc: "enum no zap",
			p:    tz.EnumDefaultFoo,
		},
		{
			desc: "struct no zap",
			p:    tz.PrimitiveRequiredStruct{},
		},
		{
			desc: "typedef of map no zap",
			p:    tz.StringMap{},
		},
		{
			desc: "typedef of struct no zap",
			p:    tz.Primitives{},
		},
		{
			desc: "typedef of list no zap",
			p:    tz.StringList{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if _, ok := tt.p.(zapcore.ObjectMarshaler); ok {
				t.Error("should not generate zap functions")
			}
		})
	}
}

func TestMarshallingErrorCollation(t *testing.T) {
	t.Run("struct", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		x := &ts.Edge{
			StartPoint: &ts.Point{},
			EndPoint:   &ts.Point{},
		}

		enc := NewMockObjectEncoder(ctrl)
		enc.EXPECT().
			AddObject(gomock.Any(), gomock.Any()).
			DoAndReturn(func(name string, _ zapcore.ObjectMarshaler) error {
				return fmt.Errorf("failed to add %q", name)
			}).
			Times(2)

		err := x.MarshalLogObject(enc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `failed to add "startPoint"`)
		assert.Contains(t, err.Error(), `failed to add "endPoint"`)
	})

	t.Run("array", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		graph := &ts.Graph{Edges: []*ts.Edge{{}, {}, {}}}

		objEnc := NewMockObjectEncoder(ctrl)
		arrEnc := NewMockArrayEncoder(ctrl)

		objEnc.EXPECT().
			AddArray("edges", gomock.Any()).
			DoAndReturn(func(_ string, arr zapcore.ArrayMarshaler) error {
				return arr.MarshalLogArray(arrEnc)
			})

		pos := 0
		arrEnc.EXPECT().
			AppendObject(gomock.Any()).
			DoAndReturn(func(zapcore.ObjectMarshaler) error {
				pos++
				return fmt.Errorf("failed to append object %v", pos)
			}).
			Times(3)

		err := graph.MarshalLogObject(objEnc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "failed to append object 1")
		assert.Contains(t, err.Error(), "failed to append object 2")
	})

	t.Run("map/string key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		x := ts.UserMap{
			"foo": &ts.User{},
			"bar": &ts.User{},
			"baz": &ts.User{},
		}

		enc := NewMockObjectEncoder(ctrl)
		enc.EXPECT().
			AddObject(gomock.Any(), gomock.Any()).
			DoAndReturn(func(name string, _ zapcore.ObjectMarshaler) error {
				return fmt.Errorf("failed to add %q", name)
			}).
			Times(3)

		err := x.MarshalLogObject(enc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `failed to add "foo"`)
		assert.Contains(t, err.Error(), `failed to add "bar"`)
		assert.Contains(t, err.Error(), `failed to add "baz"`)
	})

	t.Run("map/non-string key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		x := td.PointMap{
			{
				Key:   &ts.Point{},
				Value: &ts.Point{},
			},
			{
				Key:   &ts.Point{},
				Value: &ts.Point{},
			},
		}

		idx := 0
		enc := NewMockArrayEncoder(ctrl)
		enc.EXPECT().
			AppendObject(gomock.Any()).
			DoAndReturn(func(obj zapcore.ObjectMarshaler) error {
				idx++

				objEnc := NewMockObjectEncoder(ctrl)
				keyCall := objEnc.EXPECT().AddObject("key", gomock.Any())
				valueCall := objEnc.EXPECT().AddObject("value", gomock.Any())

				// Fail the value for the first item and the key for the
				// second item.

				if idx == 1 {
					keyCall.Return(fmt.Errorf("failed to add key for item %v", idx))
				} else {
					valueCall.Return(fmt.Errorf("failed to add value for item %v", idx))
				}

				return obj.MarshalLogObject(objEnc)
			}).
			Times(2)

		err := x.MarshalLogArray(enc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `failed to add key for item 1`)
		assert.Contains(t, err.Error(), `failed to add value for item 2`)
	})

	t.Run("set", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fs := td.FrameGroup{{}, {}}

		idx := 0
		enc := NewMockArrayEncoder(ctrl)
		enc.EXPECT().
			AppendObject(gomock.Any()).
			DoAndReturn(func(obj zapcore.ObjectMarshaler) error {
				idx++
				return fmt.Errorf("could not add item %v", idx)
			}).
			Times(2)

		err := fs.MarshalLogArray(enc)
		require.Error(t, err)
		assert.Contains(t, err.Error(), `could not add item 1`)
		assert.Contains(t, err.Error(), `could not add item 2`)
	})
}

func TestLogNilStruct(t *testing.T) {
	enc := zapcore.NewMapObjectEncoder()

	var x *ts.Edge
	require.NoError(t, x.MarshalLogObject(enc))
	assert.Empty(t, enc.Fields)
}
