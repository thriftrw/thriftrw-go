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
	"testing"

	tec "go.uber.org/thriftrw/gen/testdata/enum_conflict"
	te "go.uber.org/thriftrw/gen/testdata/enums"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
)

func TestValueOfEnumDefault(t *testing.T) {
	tests := []struct {
		e te.EnumDefault
		i int32
	}{
		{te.EnumDefaultFoo, 0},
		{te.EnumDefaultBar, 1},
		{te.EnumDefaultBaz, 2},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEqualsEnumDefault(t *testing.T) {
	tests := []struct {
		lhs, rhs te.EnumDefault
		want     bool
	}{
		{te.EnumDefaultFoo, te.EnumDefaultFoo, true},
		{te.EnumDefaultBar, te.EnumDefaultBar, true},
		{te.EnumDefaultBaz, te.EnumDefaultBaz, true},
		{te.EnumDefaultFoo, te.EnumDefaultBar, false},
		{te.EnumDefaultBaz, te.EnumDefaultBar, false},
		{te.EnumDefaultBaz, te.EnumDefaultFoo, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(tt.rhs), tt.want)
	}
}

func TestValueOfEnumWithValues(t *testing.T) {
	tests := []struct {
		e te.EnumWithValues
		i int32
	}{
		{te.EnumWithValuesX, 123},
		{te.EnumWithValuesY, 456},
		{te.EnumWithValuesZ, 789},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEqualsEnumWithValues(t *testing.T) {
	tests := []struct {
		lhs, rhs te.EnumWithValues
		want     bool
	}{
		{te.EnumWithValuesX, te.EnumWithValuesX, true},
		{te.EnumWithValuesY, te.EnumWithValuesY, true},
		{te.EnumWithValuesZ, te.EnumWithValuesZ, true},
		{te.EnumWithValuesX, te.EnumWithValuesY, false},
		{te.EnumWithValuesY, te.EnumWithValuesZ, false},
		{te.EnumWithValuesZ, te.EnumWithValuesX, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(tt.rhs), tt.want)
	}
}

func TestEnumDefaultWire(t *testing.T) {
	tests := []struct {
		e te.EnumDefault
		v wire.Value
	}{
		{te.EnumDefaultFoo, wire.NewValueI32(0)},
		{te.EnumDefaultBar, wire.NewValueI32(1)},
		{te.EnumDefaultBaz, wire.NewValueI32(2)},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.e, tt.v, "EnumDefault")
	}
}

func TestValueOfEnumWithDuplicateValues(t *testing.T) {
	tests := []struct {
		e te.EnumWithDuplicateValues
		i int32
	}{
		{te.EnumWithDuplicateValuesP, 0},
		{te.EnumWithDuplicateValuesQ, -1},
		{te.EnumWithDuplicateValuesR, 0},
	}
	for _, tt := range tests {
		assert.Equal(t, int32(tt.e), tt.i, "Value for %v does not match", tt.e)
	}
}

func TestEqualsEnumWithDuplicateValues(t *testing.T) {
	tests := []struct {
		lhs, rhs te.EnumWithDuplicateValues
		want     bool
	}{
		{te.EnumWithDuplicateValuesP, te.EnumWithDuplicateValuesP, true},
		{te.EnumWithDuplicateValuesQ, te.EnumWithDuplicateValuesQ, true},
		{te.EnumWithDuplicateValuesR, te.EnumWithDuplicateValuesR, true},
		{te.EnumWithDuplicateValuesP, te.EnumWithDuplicateValuesQ, false},
		{te.EnumWithDuplicateValuesQ, te.EnumWithDuplicateValuesR, false},
		{te.EnumWithDuplicateValuesR, te.EnumWithDuplicateValuesP, true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(tt.rhs), tt.want)
	}
}

func TestEnumWithDuplicateValuesWire(t *testing.T) {
	tests := []struct {
		e te.EnumWithDuplicateValues
		v wire.Value
	}{
		{te.EnumWithDuplicateValuesP, wire.NewValueI32(0)},
		{te.EnumWithDuplicateValuesQ, wire.NewValueI32(-1)},
		{te.EnumWithDuplicateValuesR, wire.NewValueI32(0)},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.e, tt.v, "EnumWithDuplicateValues")
	}
}

func TestUnknownEnumValue(t *testing.T) {
	var e te.EnumDefault
	if assert.NoError(t, e.FromWire(wire.NewValueI32(42))) {
		assert.Equal(t, te.EnumDefault(42), e)
	}
}

func TestOptionalEnum(t *testing.T) {
	foo := te.EnumDefaultFoo

	tests := []struct {
		s te.StructWithOptionalEnum
		v wire.Value
	}{
		{
			te.StructWithOptionalEnum{E: &foo},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI32(0)},
			}}),
		},
		{
			te.StructWithOptionalEnum{},
			wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, &tt.s, tt.v, "StructWithOptionalEnum")
	}
}

func TestEqualsOptionalEnum(t *testing.T) {
	foo := te.EnumDefaultFoo
	bar := te.EnumDefaultBar

	tests := []struct {
		lhs, rhs te.StructWithOptionalEnum
		want     bool
	}{
		{te.StructWithOptionalEnum{}, te.StructWithOptionalEnum{}, true},
		{te.StructWithOptionalEnum{E: &foo}, te.StructWithOptionalEnum{E: &foo}, true},
		{te.StructWithOptionalEnum{E: &foo}, te.StructWithOptionalEnum{E: &bar}, false},
		{te.StructWithOptionalEnum{E: &foo}, te.StructWithOptionalEnum{}, false},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.lhs.Equals(&tt.rhs), tt.want)
	}
}

func TestEnumString(t *testing.T) {
	tests := []struct {
		give fmt.Stringer
		want string
	}{
		{
			te.EmptyEnum(42),
			"EmptyEnum(42)",
		},
		{
			te.EnumDefaultFoo,
			"Foo",
		},
		{
			te.EnumDefault(-1),
			"EnumDefault(-1)",
		},
		{
			te.EnumWithDuplicateValuesP,
			"P",
		},
		{
			te.EnumWithDuplicateValuesR,
			"P", // same value as P
		},
		{
			te.EnumWithDuplicateValuesQ,
			"Q",
		},
		{
			te.LowerCaseEnumContaining,
			"containing",
		},
		{
			te.LowerCaseEnumLowerCase,
			"lower_case",
		},
		{
			te.LowerCaseEnumItems,
			"items",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, tt.give.String())
	}
}

func TestJson(t *testing.T) {
	tests := []struct {
		title string
		e     te.EnumWithDuplicateValues
		json  string
	}{
		{`Q <-> "Q"`, te.EnumWithDuplicateValuesQ, `"Q"`},
		{`P <-> "P"`, te.EnumWithDuplicateValuesP, `"P"`},
		{`Q <-> "P"`, te.EnumWithDuplicateValuesR, `"P"`}, // not a typo.
		{`42 <-> 42`, 42, `42`},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			b, err := json.Marshal(tt.e)
			if assert.NoError(t, err, "unabled to marshal enum") {
				assert.Equal(t, b, []byte(tt.json), "json output doesn't match")

				var e te.EnumWithDuplicateValues = 99
				err = json.Unmarshal([]byte(tt.json), &e)
				if assert.NoError(t, err, "unabled to unmarshal enum") {
					assert.Equal(t, tt.e, e, "parsed json doesn't match")
				}
			}
		})
	}
}

func TestInvalidJSON(t *testing.T) {
	tests := []struct {
		title       string
		json        string
		errContains string
	}{
		{`empty`, ``, "unexpected end of JSON input"},
		{`ident`, `P`, "invalid character 'P'"},
		{`not terminated`, `"P`, "unexpected end of JSON input"},
		{`not started`, `P"`, "invalid character 'P'"},
		{`notint`, `42.`, "invalid character"},
		{`toobig`, ` 2147483648`, "overflow"},
		{`toosmall`, `-2147483649`, "underflow"},
		{`justbigenough`, ` 2147483647`, ""},
		{`justsmallenough`, `-2147483648`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			var e te.EnumWithDuplicateValues
			err := json.Unmarshal([]byte(tt.json), &e)
			if tt.errContains != "" {
				if assert.Error(t, err, "must NOT be able to unmarshal this") {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err, "must be able to unmarshal this")
			}
		})
	}
}

func TestEnumValuesCanBeListed(t *testing.T) {
	values := te.EnumDefault_Values()
	assert.Equal(t, values, []te.EnumDefault{te.EnumDefaultFoo, te.EnumDefaultBar, te.EnumDefaultBaz})
}

func TestUnmarshalTextReturnsValue(t *testing.T) {
	var v te.EnumDefault
	err := v.UnmarshalText([]byte("Foo"))
	assert.NoError(t, err)
	assert.Equal(t, te.EnumDefaultFoo, v)
}

func TestEnumValueCanBecomePointer(t *testing.T) {
	v := te.EnumDefaultFoo
	assert.Equal(t, &v, te.EnumDefaultFoo.Ptr())
}

func TestUnmarshalTextReturnsErrorOnInvalidValue(t *testing.T) {
	var v te.EnumDefault
	err := v.UnmarshalText([]byte("blah"))
	assert.Error(t, err)
}

func TestTextMarshaler(t *testing.T) {
	tests := []struct {
		title    string
		giveInt  te.EnumDefault
		wantText string
	}{
		{`0 <-> Foo`, te.EnumDefaultFoo, "Foo"},
		{`1 <-> Bar`, te.EnumDefaultBar, "Bar"},
		{`2 <-> Baz`, te.EnumDefaultBaz, "Baz"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			b, err := tt.giveInt.MarshalText()
			if assert.NoError(t, err, "unable to marshal enum") {
				assert.Equal(t, []byte(tt.wantText), b, "text output doesn't match")

				var v te.EnumDefault
				err = v.UnmarshalText(b)
				if assert.NoError(t, err, "unable to unmarshall enum value") {
					assert.Equal(t, tt.giveInt, v, "enum output doesn't match")
				}
			}
		})
	}
}

func TestTextMarshalerUnknownValue(t *testing.T) {
	v := te.EnumDefault(42)
	b, err := v.MarshalText()
	if assert.NoError(t, err, "unable to marshall unknown value") {
		assert.Equal(t, `42`, string(b), "text output doesn't match")
		err = v.UnmarshalText(b)
		assert.Error(t, err, "unmarshaling unknown enum value should return error")
	}
}

func TestEnumAccessors(t *testing.T) {
	t.Run("Records", func(t *testing.T) {
		t.Run("set", func(t *testing.T) {
			rt := tec.RecordTypeName
			r := tec.Records{RecordType: &rt}
			assert.Equal(t, tec.RecordTypeName, r.GetRecordType())
		})

		t.Run("unset", func(t *testing.T) {
			var r tec.Records
			assert.Equal(t, tec.DefaultRecordType, r.GetRecordType())
		})
	})

	t.Run("StructWithOptionalEnum", func(t *testing.T) {
		t.Run("set", func(t *testing.T) {
			e := te.EnumDefaultBar
			s := te.StructWithOptionalEnum{E: &e}
			assert.Equal(t, te.EnumDefaultBar, s.GetE())
		})

		t.Run("unset", func(t *testing.T) {
			var s te.StructWithOptionalEnum
			assert.Equal(t, te.EnumDefaultFoo, s.GetE())
		})
	})
}
