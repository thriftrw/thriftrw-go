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
	"reflect"
	"testing"

	"github.com/thriftrw/thriftrw-go/envelope"
	tx "github.com/thriftrw/thriftrw-go/gen/testdata/exceptions"
	tv "github.com/thriftrw/thriftrw-go/gen/testdata/services"
	"github.com/thriftrw/thriftrw-go/gen/testdata/services/service/keyvalue"
	tu "github.com/thriftrw/thriftrw-go/gen/testdata/unions"
	"github.com/thriftrw/thriftrw-go/protocol"
	"github.com/thriftrw/thriftrw-go/wire"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceArgsAndResult(t *testing.T) {
	tests := []struct {
		desc string
		x    thriftType
		v    wire.Value
	}{
		{
			desc: "setValue args",
			x: &keyvalue.SetValueArgs{
				Key:   (*tv.Key)(stringp("foo")),
				Value: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueBool(true)},
					}}),
				},
			}}),
		},
		{
			desc: "setValue result",
			x:    &keyvalue.SetValueResult{},
			v:    wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
		},
		{
			desc: "getValue args",
			x:    &keyvalue.GetValueArgs{Key: (*tv.Key)(stringp("foo"))},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
		},
		{
			desc: "getValue result success",
			x: &keyvalue.GetValueResult{
				Success: &tu.ArbitraryValue{Int64Value: int64p(42)},
			},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 0,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 2, Value: wire.NewValueI64(42)},
					}}),
				},
			}}),
		},
		{
			desc: "getValue result failure",
			x: &keyvalue.GetValueResult{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
		},
		{
			desc: "deleteValue args",
			x:    &keyvalue.DeleteValueArgs{Key: (*tv.Key)(stringp("foo"))},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
		},
		{
			desc: "deleteValue result success",
			x:    &keyvalue.DeleteValueResult{},
			v:    wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
		},
		{
			desc: "deleteValue result failure",
			x: &keyvalue.DeleteValueResult{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 1,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
		},
		{
			desc: "deleteValue result failure 2",
			x: &keyvalue.DeleteValueResult{
				InternalError: &tv.InternalError{Message: stringp("foo")},
			},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueString("foo")},
					}}),
				},
			}}),
		},
		{
			desc: "size result",
			x:    &keyvalue.SizeResult{Success: int64p(42)},
			v: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 0, Value: wire.NewValueI64(42)},
			}}),
		},
	}

	for _, tt := range tests {
		assertRoundTrip(t, tt.x, tt.v, tt.desc)
	}
}

func TestServiceArgs(t *testing.T) {
	tests := []struct {
		input  interface{}
		output interface{}
	}{
		{
			input: keyvalue.SetValueHelper.Args(
				(*tv.Key)(stringp("foo")),
				&tu.ArbitraryValue{BoolValue: boolp(true)},
			),
			output: &keyvalue.SetValueArgs{
				Key:   (*tv.Key)(stringp("foo")),
				Value: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
		},
		{
			input:  keyvalue.GetValueHelper.Args((*tv.Key)(stringp("foo"))),
			output: &keyvalue.GetValueArgs{Key: (*tv.Key)(stringp("foo"))},
		},
		{
			input:  keyvalue.DeleteValueHelper.Args((*tv.Key)(stringp("foo"))),
			output: &keyvalue.DeleteValueArgs{Key: (*tv.Key)(stringp("foo"))},
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
			isException: keyvalue.SetValueHelper.IsException,
			err:         &tx.DoesNotExistException{Key: "foo"},
			expected:    false,
		},
		{
			isException: keyvalue.SetValueHelper.IsException,
			err:         errors.New("some error"),
			expected:    false,
		},
		{
			isException: keyvalue.GetValueHelper.IsException,
			err:         &tx.DoesNotExistException{Key: "foo"},
			expected:    true,
		},
		{
			isException: keyvalue.GetValueHelper.IsException,
			err:         errors.New("some error"),
			expected:    false,
		},
		{
			isException: keyvalue.DeleteValueHelper.IsException,
			err:         &tv.InternalError{},
			expected:    true,
		},
		{
			isException: keyvalue.DeleteValueHelper.IsException,
			err:         &tv.InternalError{Message: stringp("foo")},
			expected:    true,
		},
		{
			isException: keyvalue.DeleteValueHelper.IsException,
			err:         &tx.DoesNotExistException{Key: "foo"},
			expected:    true,
		},
		{
			isException: keyvalue.DeleteValueHelper.IsException,
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
				return keyvalue.SetValueHelper.WrapResponse(nil)
			},
			expectedResult: &keyvalue.SetValueResult{},
		},
		{
			desc: "setValue failure",
			run: func() (interface{}, error) {
				return keyvalue.SetValueHelper.WrapResponse(errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
		{
			desc: "getValue success",
			run: func() (interface{}, error) {
				return keyvalue.GetValueHelper.WrapResponse(&tu.ArbitraryValue{BoolValue: boolp(true)}, nil)
			},
			expectedResult: &keyvalue.GetValueResult{
				Success: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
		},
		{
			desc: "getValue application error",
			run: func() (interface{}, error) {
				return keyvalue.GetValueHelper.WrapResponse(nil, &tx.DoesNotExistException{Key: "foo"})
			},
			expectedResult: &keyvalue.GetValueResult{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
		},
		{
			desc: "getValue failure",
			run: func() (interface{}, error) {
				return keyvalue.GetValueHelper.WrapResponse(nil, errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
		{
			desc: "deleteValue success",
			run: func() (interface{}, error) {
				return keyvalue.DeleteValueHelper.WrapResponse(nil)
			},
			expectedResult: &keyvalue.DeleteValueResult{},
		},
		{
			desc: "deleteValue application error (1)",
			run: func() (interface{}, error) {
				return keyvalue.DeleteValueHelper.WrapResponse(&tx.DoesNotExistException{Key: "foo"})
			},
			expectedResult: &keyvalue.DeleteValueResult{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
		},
		{
			desc: "deleteValue application error (2)",
			run: func() (interface{}, error) {
				return keyvalue.DeleteValueHelper.WrapResponse(&tv.InternalError{})
			},
			expectedResult: &keyvalue.DeleteValueResult{
				InternalError: &tv.InternalError{},
			},
		},
		{
			desc: "deleteValue failure",
			run: func() (interface{}, error) {
				return keyvalue.DeleteValueHelper.WrapResponse(errors.New("foo"))
			},
			expectedError: errors.New("foo"),
		},
		{
			desc: "size success",
			run: func() (interface{}, error) {
				return keyvalue.SizeHelper.WrapResponse(42, nil)
			},
			expectedResult: &keyvalue.SizeResult{Success: int64p(42)},
		},
		{
			desc: "size failure",
			run: func() (interface{}, error) {
				return keyvalue.SizeHelper.WrapResponse(42, errors.New("foo"))
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
	_, err := keyvalue.GetValueHelper.WrapResponse(nil, (*tx.DoesNotExistException)(nil))
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(),
			"received non-nil error type with nil value for GetValueResult.DoesNotExist")
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
			unwrapResponse: keyvalue.SetValueHelper.UnwrapResponse,
			resultArg:      &keyvalue.SetValueResult{},
		},
		{
			desc:           "getValue success",
			unwrapResponse: keyvalue.GetValueHelper.UnwrapResponse,
			resultArg: &keyvalue.GetValueResult{
				Success: &tu.ArbitraryValue{BoolValue: boolp(true)},
			},
			expectedReturn: &tu.ArbitraryValue{BoolValue: boolp(true)},
		},
		{
			desc:           "getValue failure",
			unwrapResponse: keyvalue.GetValueHelper.UnwrapResponse,
			resultArg: &keyvalue.GetValueResult{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "getValue failure with success set",
			unwrapResponse: keyvalue.GetValueHelper.UnwrapResponse,
			resultArg: &keyvalue.GetValueResult{
				Success:      &tu.ArbitraryValue{BoolValue: boolp(true)},
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			// exception takes precedence over success
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "deleteValue success",
			unwrapResponse: keyvalue.DeleteValueHelper.UnwrapResponse,
			resultArg:      &keyvalue.DeleteValueResult{},
		},
		{
			desc:           "deleteValue failure (1)",
			unwrapResponse: keyvalue.DeleteValueHelper.UnwrapResponse,
			resultArg: &keyvalue.DeleteValueResult{
				DoesNotExist: &tx.DoesNotExistException{Key: "foo"},
			},
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "deleteValue failure (2)",
			unwrapResponse: keyvalue.DeleteValueHelper.UnwrapResponse,
			resultArg: &keyvalue.DeleteValueResult{
				InternalError: &tv.InternalError{},
			},
			expectedError: &tv.InternalError{},
		},
		{
			desc:           "deleteValue failure multiple values set",
			unwrapResponse: keyvalue.DeleteValueHelper.UnwrapResponse,
			resultArg: &keyvalue.DeleteValueResult{
				DoesNotExist:  &tx.DoesNotExistException{Key: "foo"},
				InternalError: &tv.InternalError{},
			},
			// lower field ID takes precedence
			expectedError: &tx.DoesNotExistException{Key: "foo"},
		},
		{
			desc:           "size success",
			unwrapResponse: keyvalue.SizeHelper.UnwrapResponse,
			resultArg:      &keyvalue.SizeResult{Success: int64p(42)},
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
	getResponse, err := keyvalue.GetValueHelper.WrapResponse(&tu.ArbitraryValue{BoolValue: boolp(true)}, nil)
	require.NoError(t, err, "Failed to get successful GetValue response")

	tests := []struct {
		s            envelope.Enveloper
		wantEnvelope wire.Envelope
	}{
		{
			s: keyvalue.GetValueHelper.Args((*tv.Key)(stringp("foo"))),
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
			s: keyvalue.DeleteValueHelper.Args((*tv.Key)(stringp("foo"))),
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
			assert.Equal(t, expected, envelope, "Envelope mismatch for %v", tt)
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
			serialize: keyvalue.SetValueHelper.Args(
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
			wantError: "ArbitraryValue should receive exactly one field value: received 0 values",
		},
		{
			desc: "SetValueV2: args: missing value",
			typ:  reflect.TypeOf(keyvalue.SetValueV2Args{}),
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{
					ID: 2,
					Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
						{ID: 1, Value: wire.NewValueBool(true)},
					}}),
				},
			}}),
			wantError: "field Key of SetValueV2Args is required",
		},
		{
			desc: "SetValueV2: args: missing key",
			typ:  reflect.TypeOf(keyvalue.SetValueV2Args{}),
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
				{ID: 1, Value: wire.NewValueString("foo")},
			}}),
			wantError: "field Value of SetValueV2Args is required",
		},
		{
			desc:        "getValue: result: empty",
			serialize:   &keyvalue.GetValueResult{},
			deserialize: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{}}),
			wantError:   "GetValueResult should receive exactly one field value: received 0 values",
		},
		{
			desc: "getValue: result: multiple",
			serialize: &keyvalue.GetValueResult{
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
			wantError: "GetValueResult should receive exactly one field value: received 2 values",
		},
		{
			desc: "deleteValue: result: multiple",
			serialize: &keyvalue.DeleteValueResult{
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
			wantError: "DeleteValueResult should receive at most one field value: received 2 values",
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

		x := reflect.New(typ)
		args := []reflect.Value{reflect.ValueOf(tt.deserialize)}
		e := x.MethodByName("FromWire").Call(args)[0].Interface()
		if assert.NotNil(t, e, "%v: expected failure but got %v", tt.desc, x) {
			assert.Contains(t, e.(error).Error(), tt.wantError, tt.desc)
		}
	}
}
