// Code generated by thriftrw v1.13.0. DO NOT EDIT.
// @generated

// Copyright (c) 2018 Uber Technologies, Inc.
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

package api

import (
	"errors"
	"fmt"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
	"strings"
)

// ServiceGenerator_Generate_Args represents the arguments for the ServiceGenerator.generate function.
//
// The arguments for generate are sent and received over the wire as this struct.
type ServiceGenerator_Generate_Args struct {
	Request *GenerateServiceRequest `json:"request,omitempty"`
}

// ToWire translates a ServiceGenerator_Generate_Args struct into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// An error is returned if the struct or any of its fields failed to
// validate.
//
//   x, err := v.ToWire()
//   if err != nil {
//     return err
//   }
//
//   if err := binaryProtocol.Encode(x, writer); err != nil {
//     return err
//   }
func (v *ServiceGenerator_Generate_Args) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.Request != nil {
		w, err = v.Request.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _GenerateServiceRequest_Read(w wire.Value) (*GenerateServiceRequest, error) {
	var v GenerateServiceRequest
	err := v.FromWire(w)
	return &v, err
}

// FromWire deserializes a ServiceGenerator_Generate_Args struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a ServiceGenerator_Generate_Args struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v ServiceGenerator_Generate_Args
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *ServiceGenerator_Generate_Args) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.Request, err = _GenerateServiceRequest_Read(field.Value)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

// String returns a readable string representation of a ServiceGenerator_Generate_Args
// struct.
func (v *ServiceGenerator_Generate_Args) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [1]string
	i := 0
	if v.Request != nil {
		fields[i] = fmt.Sprintf("Request: %v", v.Request)
		i++
	}

	return fmt.Sprintf("ServiceGenerator_Generate_Args{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this ServiceGenerator_Generate_Args match the
// provided ServiceGenerator_Generate_Args.
//
// This function performs a deep comparison.
func (v *ServiceGenerator_Generate_Args) Equals(rhs *ServiceGenerator_Generate_Args) bool {
	if !((v.Request == nil && rhs.Request == nil) || (v.Request != nil && rhs.Request != nil && v.Request.Equals(rhs.Request))) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, allowing
// fast logging of ServiceGenerator_Generate_Args.
func (v *ServiceGenerator_Generate_Args) MarshalLogObject(enc zapcore.ObjectEncoder) error {

	if v.Request != nil {
		if err := enc.AddObject("request", v.Request); err != nil {
			return err
		}
	}

	return nil
}

// GetRequest returns the value of Request if it is set or its
// zero value if it is unset.
func (v *ServiceGenerator_Generate_Args) GetRequest() (o *GenerateServiceRequest) {
	if v.Request != nil {
		return v.Request
	}

	return
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the arguments.
//
// This will always be "generate" for this struct.
func (v *ServiceGenerator_Generate_Args) MethodName() string {
	return "generate"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Call for this struct.
func (v *ServiceGenerator_Generate_Args) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

// ServiceGenerator_Generate_Helper provides functions that aid in handling the
// parameters and return values of the ServiceGenerator.generate
// function.
var ServiceGenerator_Generate_Helper = struct {
	// Args accepts the parameters of generate in-order and returns
	// the arguments struct for the function.
	Args func(
		request *GenerateServiceRequest,
	) *ServiceGenerator_Generate_Args

	// IsException returns true if the given error can be thrown
	// by generate.
	//
	// An error can be thrown by generate only if the
	// corresponding exception type was mentioned in the 'throws'
	// section for it in the Thrift file.
	IsException func(error) bool

	// WrapResponse returns the result struct for generate
	// given its return value and error.
	//
	// This allows mapping values and errors returned by
	// generate into a serializable result struct.
	// WrapResponse returns a non-nil error if the provided
	// error cannot be thrown by generate
	//
	//   value, err := generate(args)
	//   result, err := ServiceGenerator_Generate_Helper.WrapResponse(value, err)
	//   if err != nil {
	//     return fmt.Errorf("unexpected error from generate: %v", err)
	//   }
	//   serialize(result)
	WrapResponse func(*GenerateServiceResponse, error) (*ServiceGenerator_Generate_Result, error)

	// UnwrapResponse takes the result struct for generate
	// and returns the value or error returned by it.
	//
	// The error is non-nil only if generate threw an
	// exception.
	//
	//   result := deserialize(bytes)
	//   value, err := ServiceGenerator_Generate_Helper.UnwrapResponse(result)
	UnwrapResponse func(*ServiceGenerator_Generate_Result) (*GenerateServiceResponse, error)
}{}

func init() {
	ServiceGenerator_Generate_Helper.Args = func(
		request *GenerateServiceRequest,
	) *ServiceGenerator_Generate_Args {
		return &ServiceGenerator_Generate_Args{
			Request: request,
		}
	}

	ServiceGenerator_Generate_Helper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}

	ServiceGenerator_Generate_Helper.WrapResponse = func(success *GenerateServiceResponse, err error) (*ServiceGenerator_Generate_Result, error) {
		if err == nil {
			return &ServiceGenerator_Generate_Result{Success: success}, nil
		}

		return nil, err
	}
	ServiceGenerator_Generate_Helper.UnwrapResponse = func(result *ServiceGenerator_Generate_Result) (success *GenerateServiceResponse, err error) {

		if result.Success != nil {
			success = result.Success
			return
		}

		err = errors.New("expected a non-void result")
		return
	}

}

// ServiceGenerator_Generate_Result represents the result of a ServiceGenerator.generate function call.
//
// The result of a generate execution is sent and received over the wire as this struct.
//
// Success is set only if the function did not throw an exception.
type ServiceGenerator_Generate_Result struct {
	// Value returned by generate after a successful execution.
	Success *GenerateServiceResponse `json:"success,omitempty"`
}

// ToWire translates a ServiceGenerator_Generate_Result struct into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// An error is returned if the struct or any of its fields failed to
// validate.
//
//   x, err := v.ToWire()
//   if err != nil {
//     return err
//   }
//
//   if err := binaryProtocol.Encode(x, writer); err != nil {
//     return err
//   }
func (v *ServiceGenerator_Generate_Result) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.Success != nil {
		w, err = v.Success.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 0, Value: w}
		i++
	}

	if i != 1 {
		return wire.Value{}, fmt.Errorf("ServiceGenerator_Generate_Result should have exactly one field: got %v fields", i)
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _GenerateServiceResponse_Read(w wire.Value) (*GenerateServiceResponse, error) {
	var v GenerateServiceResponse
	err := v.FromWire(w)
	return &v, err
}

// FromWire deserializes a ServiceGenerator_Generate_Result struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a ServiceGenerator_Generate_Result struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v ServiceGenerator_Generate_Result
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *ServiceGenerator_Generate_Result) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 0:
			if field.Value.Type() == wire.TStruct {
				v.Success, err = _GenerateServiceResponse_Read(field.Value)
				if err != nil {
					return err
				}

			}
		}
	}

	count := 0
	if v.Success != nil {
		count++
	}
	if count != 1 {
		return fmt.Errorf("ServiceGenerator_Generate_Result should have exactly one field: got %v fields", count)
	}

	return nil
}

// String returns a readable string representation of a ServiceGenerator_Generate_Result
// struct.
func (v *ServiceGenerator_Generate_Result) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [1]string
	i := 0
	if v.Success != nil {
		fields[i] = fmt.Sprintf("Success: %v", v.Success)
		i++
	}

	return fmt.Sprintf("ServiceGenerator_Generate_Result{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this ServiceGenerator_Generate_Result match the
// provided ServiceGenerator_Generate_Result.
//
// This function performs a deep comparison.
func (v *ServiceGenerator_Generate_Result) Equals(rhs *ServiceGenerator_Generate_Result) bool {
	if !((v.Success == nil && rhs.Success == nil) || (v.Success != nil && rhs.Success != nil && v.Success.Equals(rhs.Success))) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, allowing
// fast logging of ServiceGenerator_Generate_Result.
func (v *ServiceGenerator_Generate_Result) MarshalLogObject(enc zapcore.ObjectEncoder) error {

	if v.Success != nil {
		if err := enc.AddObject("success", v.Success); err != nil {
			return err
		}
	}

	return nil
}

// GetSuccess returns the value of Success if it is set or its
// zero value if it is unset.
func (v *ServiceGenerator_Generate_Result) GetSuccess() (o *GenerateServiceResponse) {
	if v.Success != nil {
		return v.Success
	}

	return
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the result.
//
// This will always be "generate" for this struct.
func (v *ServiceGenerator_Generate_Result) MethodName() string {
	return "generate"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Reply for this struct.
func (v *ServiceGenerator_Generate_Result) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}
