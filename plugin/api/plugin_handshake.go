// Code generated by thriftrw v0.6.0
// @generated

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

package api

import (
	"errors"
	"fmt"
	"go.uber.org/thriftrw/wire"
	"strings"
)

type Plugin_Handshake_Args struct {
	Request *HandshakeRequest `json:"request,omitempty"`
}

func (v *Plugin_Handshake_Args) ToWire() (wire.Value, error) {
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

func _HandshakeRequest_Read(w wire.Value) (*HandshakeRequest, error) {
	var v HandshakeRequest
	err := v.FromWire(w)
	return &v, err
}

func (v *Plugin_Handshake_Args) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.Request, err = _HandshakeRequest_Read(field.Value)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (v *Plugin_Handshake_Args) String() string {
	var fields [1]string
	i := 0
	if v.Request != nil {
		fields[i] = fmt.Sprintf("Request: %v", v.Request)
		i++
	}
	return fmt.Sprintf("Plugin_Handshake_Args{%v}", strings.Join(fields[:i], ", "))
}

func (v *Plugin_Handshake_Args) MethodName() string {
	return "handshake"
}

func (v *Plugin_Handshake_Args) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

var Plugin_Handshake_Helper = struct {
	Args           func(request *HandshakeRequest) *Plugin_Handshake_Args
	IsException    func(error) bool
	WrapResponse   func(*HandshakeResponse, error) (*Plugin_Handshake_Result, error)
	UnwrapResponse func(*Plugin_Handshake_Result) (*HandshakeResponse, error)
}{}

func init() {
	Plugin_Handshake_Helper.Args = func(request *HandshakeRequest) *Plugin_Handshake_Args {
		return &Plugin_Handshake_Args{Request: request}
	}
	Plugin_Handshake_Helper.IsException = func(err error) bool {
		switch err.(type) {
		default:
			return false
		}
	}
	Plugin_Handshake_Helper.WrapResponse = func(success *HandshakeResponse, err error) (*Plugin_Handshake_Result, error) {
		if err == nil {
			return &Plugin_Handshake_Result{Success: success}, nil
		}
		return nil, err
	}
	Plugin_Handshake_Helper.UnwrapResponse = func(result *Plugin_Handshake_Result) (success *HandshakeResponse, err error) {
		if result.Success != nil {
			success = result.Success
			return
		}
		err = errors.New("expected a non-void result")
		return
	}
}

type Plugin_Handshake_Result struct {
	Success *HandshakeResponse `json:"success,omitempty"`
}

func (v *Plugin_Handshake_Result) ToWire() (wire.Value, error) {
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
		return wire.Value{}, fmt.Errorf("Plugin_Handshake_Result should have exactly one field: got %v fields", i)
	}
	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _HandshakeResponse_Read(w wire.Value) (*HandshakeResponse, error) {
	var v HandshakeResponse
	err := v.FromWire(w)
	return &v, err
}

func (v *Plugin_Handshake_Result) FromWire(w wire.Value) error {
	var err error
	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 0:
			if field.Value.Type() == wire.TStruct {
				v.Success, err = _HandshakeResponse_Read(field.Value)
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
		return fmt.Errorf("Plugin_Handshake_Result should have exactly one field: got %v fields", count)
	}
	return nil
}

func (v *Plugin_Handshake_Result) String() string {
	var fields [1]string
	i := 0
	if v.Success != nil {
		fields[i] = fmt.Sprintf("Success: %v", v.Success)
		i++
	}
	return fmt.Sprintf("Plugin_Handshake_Result{%v}", strings.Join(fields[:i], ", "))
}

func (v *Plugin_Handshake_Result) MethodName() string {
	return "handshake"
}

func (v *Plugin_Handshake_Result) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}
