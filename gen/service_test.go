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
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"go.uber.org/thriftrw/envelope"
	tx "go.uber.org/thriftrw/gen/internal/tests/exceptions"
	tv "go.uber.org/thriftrw/gen/internal/tests/services"
	tu "go.uber.org/thriftrw/gen/internal/tests/unions"
	"go.uber.org/thriftrw/internal/envelope/envelopetest"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// serviceType is the args or result struct for a thrift function.
type serviceType interface {
	fmt.Stringer
	envelope.Enveloper

	FromWire(wire.Value) error
}

func TestServiceArgsAndResult(t *testing.T) {
	tests := []struct {
		desc         string
		x            serviceType
		value        wire.Value
		methodName   string
		envelopeType wire.EnvelopeType
	}{
		{
			desc: "setValue args",
			x: &tv.KeyValue_SetValue_Args{
				Key:   (*tv.Key)(stringp("foo")),
				Value: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueBool(true)},
					}}),
				},
			}}),
			methodName:   "setValue",
			envelopeType: wire.Call,
		},
		{
			desc:         "setValue result",
			x:            &tv.KeyValue_SetValue_Result{},
			value:        wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			methodName:   "setValue",
			envelopeType: wire.Reply,
		},
		{
			desc: "getValue args",
			x:    &tv.KeyValue_GetValue_Args{Key: (*tv.Key)(stringp("foo"))},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
			methodName:   "getValue",
			envelopeType: wire.Call,
		},
		{
			desc: "getValue result success",
			x: &tv.KeyValue_GetValue_Result{
				Success: &tu.ArbitraryValue{Int64Value: int64p(42)},
			},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 0,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 2, Value: wire.NewValueI64(42)},
					}}),
				},
			}}),
			methodName:   "getValue",
			envelopeType: wire.Reply,
		},
		{
			desc: "getValue result failure",
			x: &tv.KeyValue_GetValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
			methodName:   "getValue",
			envelopeType: wire.Reply,
		},
		{
			desc: "deleteValue args",
			x:    &tv.KeyValue_DeleteValue_Args{Key: (*tv.Key)(stringp("foo"))},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
			methodName:   "deleteValue",
			envelopeType: wire.Call,
		},
		{
			desc:         "deleteValue result success",
			x:            &tv.KeyValue_DeleteValue_Result{},
			value:        wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			methodName:   "deleteValue",
			envelopeType: wire.Reply,
		},
		{
			desc: "deleteValue result failure",
			x: &tv.KeyValue_DeleteValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
			methodName:   "deleteValue",
			envelopeType: wire.Reply,
		},
		{
			desc: "deleteValue result failure 2",
			x: &tv.KeyValue_DeleteValue_Result{
				InternalError: &tv.InternalError{Message: stringp("foo")},
			},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
			methodName:   "deleteValue",
			envelopeType: wire.Reply,
		},
		{
			desc: "size result",
			x:    &tv.KeyValue_Size_Result{Success: int64p(42)},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 0, Value: wire.NewValueI64(42)},
			}}),
			methodName:   "size",
			envelopeType: wire.Reply,
		},
		{
			desc:         "oneway empty args",
			x:            &tv.Cache_Clear_Args{},
			value:        wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			methodName:   "clear",
			envelopeType: wire.OneWay,
		},
		{
			desc: "oneway with args",
			x:    &tv.Cache_ClearAfter_Args{DurationMS: ptr.Int64(42)},
			value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueI64(42)},
			}}),
			methodName:   "clearAfter",
			envelopeType: wire.OneWay,
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, tt.x, tt.value, tt.desc)
		assert.Equal(t, tt.methodName, tt.x.MethodName(), tt.desc)
		assert.Equal(t, tt.envelopeType, tt.x.EnvelopeType(), tt.desc)
	}
}

func TestServiceArgs(t *testing.T) {
	tests := []struct {
		input  interface{}
		output interface{}
	}{
		{
			input: tv.KeyValue_SetValue_Helper.Args(
				(*tv.Key)(stringp("foo")),
				&tu.ArbitraryValue{BoolValue: boolp(true)},
			),
			output: &tv.KeyValue_SetValue_Args{
				Key:   (*tv.Key)(stringp("foo")),
				Value: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
		},
		{
			input:  tv.KeyValue_GetValue_Helper.Args((*tv.Key)(stringp("foo"))),
			output: &tv.KeyValue_GetValue_Args{Key: (*tv.Key)(stringp("foo"))},
		},
		{
			input:  tv.KeyValue_DeleteValue_Helper.Args((*tv.Key)(stringp("foo"))),
			output: &tv.KeyValue_DeleteValue_Args{Key: (*tv.Key)(stringp("foo"))},
		},
		{
			input:  tv.Cache_Clear_Helper.Args(),
			output: &tv.Cache_Clear_Args{},
		},
		{
			input:  tv.Cache_ClearAfter_Helper.Args(ptr.Int64(42)),
			output: &tv.Cache_ClearAfter_Args{DurationMS: ptr.Int64(42)},
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.output, tt.input)
	}
}

func TestServiceIsException(t *testing.T) {
	tests := []struct {
		isException func(error) bool
		err         error
		expected    bool
	}{
		{
			isException: tv.KeyValue_SetValue_Helper.IsException,
			err:         &tx.DoesNotExistException{Key: "foo"},
			expected:    false,
		},
		{
			isException: tv.KeyValue_SetValue_Helper.IsException,
			err:         errors.New("some error"),
			expected:    false,
		},
		{
			isException: tv.KeyValue_GetValue_Helper.IsException,
			err:         &tx.DoesNotExistException{Key: "foo"},
			expected:    true,
		},
		{
			isException: tv.KeyValue_GetValue_Helper.IsException,
			err:         errors.New("some error"),
			expected:    false,
		},
		{
			isException: tv.KeyValue_DeleteValue_Helper.IsException,
			err:         &tv.InternalError{},
			expected:    true,
		},
		{
			isException: tv.KeyValue_DeleteValue_Helper.IsException,
			err:         &tv.InternalError{Message: stringp("foo")},
			expected:    true,
		},
		{
			isException: tv.KeyValue_DeleteValue_Helper.IsException,
			err:         &tx.DoesNotExistException{Key: "foo"},
			expected:    true,
		},
		{
			isException: tv.KeyValue_DeleteValue_Helper.IsException,
			err:         errors.New("some error"),
			expected:    false,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.isException(tt.err))
	}
}

func TestWrapResponse(t *testing.T) {
	tests := []struct {
		desc           string
		run            func() (interface{}, error)
		expectedResult interface{}
		expectedError  error
	}{
		{
			desc: "setValue success",
			run: func() (interface{}, error) {
				return tv.KeyValue_SetValue_Helper.WrapResponse(nil)
			},
			expectedResult: &tv.KeyValue_SetValue_Result{},
		},
		{
			desc: "setValue failure",
			run: func() (interface{}, error) {
				return tv.KeyValue_SetValue_Helper.WrapResponse(errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
		{
			desc: "getValue success",
			run: func() (interface{}, error) {
				return tv.KeyValue_GetValue_Helper.WrapResponse(&tu.ArbitraryValue{BoolValue: boolp(true)}, nil)
			},
			expectedResult: &tv.KeyValue_GetValue_Result{
				Success: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
		},
		{
			desc: "getValue application error",
			run: func() (interface{}, error) {
				return tv.KeyValue_GetValue_Helper.WrapResponse(nil, &tx.DoesNotExistException{Key: "foo"})
			},
			expectedResult: &tv.KeyValue_GetValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
		},
		{
			desc: "getValue failure",
			run: func() (interface{}, error) {
				return tv.KeyValue_GetValue_Helper.WrapResponse(nil, errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
		{
			desc: "deleteValue success",
			run: func() (interface{}, error) {
				return tv.KeyValue_DeleteValue_Helper.WrapResponse(nil)
			},
			expectedResult: &tv.KeyValue_DeleteValue_Result{},
		},
		{
			desc: "deleteValue application error (1)",
			run: func() (interface{}, error) {
				return tv.KeyValue_DeleteValue_Helper.WrapResponse(&tx.DoesNotExistException{Key: "foo"})
			},
			expectedResult: &tv.KeyValue_DeleteValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
		},
		{
			desc: "deleteValue application error (2)",
			run: func() (interface{}, error) {
				return tv.KeyValue_DeleteValue_Helper.WrapResponse(&tv.InternalError{})
			},
			expectedResult: &tv.KeyValue_DeleteValue_Result{
				InternalError: &tv.InternalError{},
			},
		},
		{
			desc: "deleteValue failure",
			run: func() (interface{}, error) {
				return tv.KeyValue_DeleteValue_Helper.WrapResponse(errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
		{
			desc: "size success",
			run: func() (interface{}, error) {
				return tv.KeyValue_Size_Helper.WrapResponse(42, nil)
			},
			expectedResult: &tv.KeyValue_Size_Result{Success: int64p(42)},
		},
		{
			desc: "size failure",
			run: func() (interface{}, error) {
				return tv.KeyValue_Size_Helper.WrapResponse(42, errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
	}

	for _, tt := range tests {
		result, err := tt.run()
		if tt.expectedError != nil {
			assert.Equal(t, tt.expectedError, err, tt.desc)
		} else {
			assert.Equal(t, tt.expectedResult, result, tt.desc)
		}
	}
}

func TestWrapResponseTypedNilError(t *testing.T) {
	_, err := tv.KeyValue_GetValue_Helper.WrapResponse(nil, (*tx.DoesNotExistException)(nil))
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(),
			"received non-nil error type with nil value for KeyValue_GetValue_Result.DoesNotExist")
	}
}

func TestUnwrapResponse(t *testing.T) {
	tests := []struct {
		desc           string
		unwrapResponse interface{}
		resultArg      interface{}

		expectedReturn interface{}
		expectedError  error
	}{
		{
			desc:           "setValue success",
			unwrapResponse: tv.KeyValue_SetValue_Helper.UnwrapResponse,
			resultArg:      &tv.KeyValue_SetValue_Result{},
		},
		{
			desc:           "getValue success",
			unwrapResponse: tv.KeyValue_GetValue_Helper.UnwrapResponse,
			resultArg: &tv.KeyValue_GetValue_Result{
				Success: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
			expectedReturn: &tu.ArbitraryValue{BoolValue: boolp(true)},
		},
		{
			desc:           "getValue failure",
			unwrapResponse: tv.KeyValue_GetValue_Helper.UnwrapResponse,
			resultArg: &tv.KeyValue_GetValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "getValue failure with success set",
			unwrapResponse: tv.KeyValue_GetValue_Helper.UnwrapResponse,
			resultArg: &tv.KeyValue_GetValue_Result{
				Success:      &tu.ArbitraryValue{BoolValue: boolp(true)},
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			// exception takes precedence over success
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "deleteValue success",
			unwrapResponse: tv.KeyValue_DeleteValue_Helper.UnwrapResponse,
			resultArg:      &tv.KeyValue_DeleteValue_Result{},
		},
		{
			desc:           "deleteValue failure (1)",
			unwrapResponse: tv.KeyValue_DeleteValue_Helper.UnwrapResponse,
			resultArg: &tv.KeyValue_DeleteValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "deleteValue failure (2)",
			unwrapResponse: tv.KeyValue_DeleteValue_Helper.UnwrapResponse,
			resultArg: &tv.KeyValue_DeleteValue_Result{
				InternalError: &tv.InternalError{},
			},
			expectedError: &tv.InternalError{},
		},
		{
			desc:           "deleteValue failure multiple values set",
			unwrapResponse: tv.KeyValue_DeleteValue_Helper.UnwrapResponse,
			resultArg: &tv.KeyValue_DeleteValue_Result{
				DoesNotExist:  &tx.DoesNotExistException{Key: "foo"},
				InternalError: &tv.InternalError{},
			},
			// lower field ID takes precedence
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "size success",
			unwrapResponse: tv.KeyValue_Size_Helper.UnwrapResponse,
			resultArg:      &tv.KeyValue_Size_Result{Success: int64p(42)},
			expectedReturn: int64(42),
		},
	}

	for _, tt := range tests {
		unwrapResponse := reflect.ValueOf(tt.unwrapResponse)
		out := unwrapResponse.Call([]reflect.Value{reflect.ValueOf(tt.resultArg)})

		var returnValue, err interface{}
		if len(out) == 1 {
			err = out[0].Interface()
		} else {
			returnValue = out[0].Interface()
			err = out[1].Interface()
		}

		if tt.expectedError != nil {
			assert.Equal(t, tt.expectedError, err, tt.desc)
		} else {
			assert.Equal(t, tt.expectedReturn, returnValue, tt.desc)
		}
	}
}

func TestServiceTypesEnveloper(t *testing.T) {
	getResponse, err := tv.KeyValue_GetValue_Helper.WrapResponse(&tu.ArbitraryValue{BoolValue: boolp(true)}, nil)
	require.NoError(t, err, "Failed to get successful GetValue response")

	tests := []struct {
		s            envelope.Enveloper
		wantEnvelope wire.Envelope
	}{
		{
			s: tv.KeyValue_GetValue_Helper.Args((*tv.Key)(stringp("foo"))),
			wantEnvelope: wire.Envelope{
				Name: "getValue",
				Type: wire.Call,
			},
		},
		{
			s: getResponse,
			wantEnvelope: wire.Envelope{
				Name: "getValue",
				Type: wire.Reply,
			},
		},
		{
			s: tv.KeyValue_DeleteValue_Helper.Args((*tv.Key)(stringp("foo"))),
			wantEnvelope: wire.Envelope{
				Name: "deleteValue",
				Type: wire.Call,
			},
		},
	}

	for _, tt := range tests {
		buf := &bytes.Buffer{}
		err := envelope.Write(protocol.Binary, buf, 1234, tt.s)
		require.NoError(t, err, "envelope.Write for %v failed", tt)

		// Decode the payload and validate the payload.
		reader := bytes.NewReader(buf.Bytes())
		envelope, err := protocol.Binary.DecodeEnveloped(reader)
		require.NoError(t, err, "Failed to read enveloped data for %v", tt)

		expected := tt.wantEnvelope
		expected.SeqID = 1234
		expected.Value, err = tt.s.ToWire()
		if assert.NoError(t, err, "Error serializing %v", tt.s) {
			envelopetest.AssertEqual(t, expected, envelope, "envelope mismatch for %v", tt)
		}
	}
}

func TestArgsAndResultValidation(t *testing.T) {
	tests := []struct {
		desc        string
		serialize   thriftType
		deserialize wire.Value
		typ         reflect.Type // must be set if serialize is not
		wantError   string
	}{
		{
			desc: "SetValue: args: value: empty",
			serialize: tv.KeyValue_SetValue_Helper.Args(
				(*tv.Key)(stringp("foo")),
				&tu.ArbitraryValue{},
			),
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
				{
					ID:    2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
				},
			}}),
			wantError: "ArbitraryValue should have exactly one field: got 0 fields",
		},
		{
			desc: "SetValueV2: args: missing value",
			typ:  reflect.TypeOf(tv.KeyValue_SetValueV2_Args{}),
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueBool(true)},
					}}),
				},
			}}),
			wantError: "field Key of KeyValue_SetValueV2_Args is required",
		},
		{
			desc: "SetValueV2: args: missing key",
			typ:  reflect.TypeOf(tv.KeyValue_SetValueV2_Args{}),
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
			wantError: "field Value of KeyValue_SetValueV2_Args is required",
		},
		{
			desc:        "getValue: result: empty",
			serialize:   &tv.KeyValue_GetValue_Result{},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			wantError:   "KeyValue_GetValue_Result should have exactly one field: got 0 fields",
		},
		{
			desc: "getValue: result: multiple",
			serialize: &tv.KeyValue_GetValue_Result{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
				Success:      &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 0,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueBool(true)},
					}}),
				},
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
			wantError: "KeyValue_GetValue_Result should have exactly one field: got 2 fields",
		},
		{
			desc: "deleteValue: result: multiple",
			serialize: &tv.KeyValue_DeleteValue_Result{
				DoesNotExist:  &tx.DoesNotExistException{Key: "foo"},
				InternalError: &tv.InternalError{},
			},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
				{
					ID:    2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
				},
			}}),
			wantError: "KeyValue_DeleteValue_Result should have at most one field: got 2 fields",
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

		x := reflect.New(typ).Interface().(serviceType)
		if err := x.FromWire(tt.deserialize); assert.Errorf(t, err, "%v: expected failure but got %v", tt.desc, x) {
			assert.Contains(t, err.Error(), tt.wantError, tt.desc)
		}
	}
}
