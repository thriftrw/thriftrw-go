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
	"fmt"
	"go.uber.org/thriftrw/wire"
	"go.uber.org/zap/zapcore"
	"strings"
)

// Plugin_Goodbye_Args represents the arguments for the Plugin.goodbye function.
//
// The arguments for goodbye are sent and received over the wire as this struct.
type Plugin_Goodbye_Args struct {
}

// ToWire translates a Plugin_Goodbye_Args struct into a Thrift-level intermediate
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
func (v *Plugin_Goodbye_Args) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a Plugin_Goodbye_Args struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a Plugin_Goodbye_Args struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v Plugin_Goodbye_Args
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *Plugin_Goodbye_Args) FromWire(w wire.Value) error {

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}

	return nil
}

// String returns a readable string representation of a Plugin_Goodbye_Args
// struct.
func (v *Plugin_Goodbye_Args) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [0]string
	i := 0

	return fmt.Sprintf("Plugin_Goodbye_Args{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this Plugin_Goodbye_Args match the
// provided Plugin_Goodbye_Args.
//
// This function performs a deep comparison.
func (v *Plugin_Goodbye_Args) Equals(rhs *Plugin_Goodbye_Args) bool {

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler. (TODO)
func (v *Plugin_Goodbye_Args) MarshalLogObject(enc zapcore.ObjectEncoder) error {

	return nil
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the arguments.
//
// This will always be "goodbye" for this struct.
func (v *Plugin_Goodbye_Args) MethodName() string {
	return "goodbye"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Call for this struct.
func (v *Plugin_Goodbye_Args) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

// Plugin_Goodbye_Helper provides functions that aid in handling the
// parameters and return values of the Plugin.goodbye
// function.
var Plugin_Goodbye_Helper = struct {
	// Args accepts the parameters of goodbye in-order and returns
	// the arguments struct for the function.
	Args func() *Plugin_Goodbye_Args

	// IsException returns true if the given error can be thrown
	// by goodbye.
	//
	// An error can be thrown by goodbye only if the
	// corresponding exception type was mentioned in the 'throws'
	// section for it in the Thrift file.
	IsException func(error) bool

	// WrapResponse returns the result struct for goodbye
	// given the error returned by it. The provided error may
	// be nil if goodbye did not fail.
	//
	// This allows mapping errors returned by goodbye into a
	// serializable result struct. WrapResponse returns a
	// non-nil error if the provided error cannot be thrown by
	// goodbye
	//
	//   err := goodbye(args)
	//   result, err := Plugin_Goodbye_Helper.WrapResponse(err)
	//   if err != nil {
	//     return fmt.Errorf("unexpected error from goodbye: %v", err)
	//   }
	//   serialize(result)
	WrapResponse func(error) (*Plugin_Goodbye_Result, error)

	// UnwrapResponse takes the result struct for goodbye
	// and returns the erorr returned by it (if any).
	//
	// The error is non-nil only if goodbye threw an
	// exception.
	//
	//   result := deserialize(bytes)
	//   err := Plugin_Goodbye_Helper.UnwrapResponse(result)
	UnwrapResponse func(*Plugin_Goodbye_Result) error
}{}

func init() {
	Plugin_Goodbye_Helper.Args = func() *Plugin_Goodbye_Args {
		return &Plugin_Goodbye_Args{}
	}

	Plugin_Goodbye_Helper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}

	Plugin_Goodbye_Helper.WrapResponse = func(err error) (*Plugin_Goodbye_Result, error) {
		if err == nil {
			return &Plugin_Goodbye_Result{}, nil
		}

		return nil, err
	}
	Plugin_Goodbye_Helper.UnwrapResponse = func(result *Plugin_Goodbye_Result) (err error) {
		return
	}

}

// Plugin_Goodbye_Result represents the result of a Plugin.goodbye function call.
//
// The result of a goodbye execution is sent and received over the wire as this struct.
type Plugin_Goodbye_Result struct {
}

// ToWire translates a Plugin_Goodbye_Result struct into a Thrift-level intermediate
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
func (v *Plugin_Goodbye_Result) ToWire() (wire.Value, error) {
	var (
		fields [0]wire.Field
		i      int = 0
	)

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a Plugin_Goodbye_Result struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a Plugin_Goodbye_Result struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v Plugin_Goodbye_Result
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *Plugin_Goodbye_Result) FromWire(w wire.Value) error {

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		}
	}

	return nil
}

// String returns a readable string representation of a Plugin_Goodbye_Result
// struct.
func (v *Plugin_Goodbye_Result) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [0]string
	i := 0

	return fmt.Sprintf("Plugin_Goodbye_Result{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this Plugin_Goodbye_Result match the
// provided Plugin_Goodbye_Result.
//
// This function performs a deep comparison.
func (v *Plugin_Goodbye_Result) Equals(rhs *Plugin_Goodbye_Result) bool {

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler. (TODO)
func (v *Plugin_Goodbye_Result) MarshalLogObject(enc zapcore.ObjectEncoder) error {

	return nil
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the result.
//
// This will always be "goodbye" for this struct.
func (v *Plugin_Goodbye_Result) MethodName() string {
	return "goodbye"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Reply for this struct.
func (v *Plugin_Goodbye_Result) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}
