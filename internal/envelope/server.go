// Copyright (c) 2019 Uber Technologies, Inc.
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

package envelope

import (
	"bytes"
	"fmt"

	"go.uber.org/thriftrw/internal/envelope/exception"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"
)

// ErrUnknownMethod is raised by Handlers to indicate that the given method is
// invalid.
type ErrUnknownMethod string

func (e ErrUnknownMethod) Error() string {
	return fmt.Sprintf("unknown method %q", string(e))
}

// Handler handles enveloped requests.
type Handler interface {
	// Receives a request to the given method name and returns the response
	// body.
	//
	// Implementations should return ErrUnknownMethod if the method is invalid.
	Handle(name string, body wire.Value) (wire.Value, error)
}

// Server allows a Handler to process bytes.
type Server struct {
	p protocol.Protocol
	h Handler
}

// NewServer builds a new server.
func NewServer(p protocol.Protocol, h Handler) Server {
	return Server{p: p, h: h}
}

// Handle handles the given binary payload.
func (s Server) Handle(data []byte) ([]byte, error) {
	request, err := s.p.DecodeEnveloped(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	response := wire.Envelope{
		Name:  request.Name,
		SeqID: request.SeqID,
		Type:  wire.Reply,
	}

	response.Value, err = s.h.Handle(request.Name, request.Value)
	if err != nil {
		response.Type = wire.Exception
		switch err.(type) {
		case ErrUnknownMethod:
			response.Value, err = tappExc(err, exception.ExceptionTypeUnknownMethod)
		default:
			response.Value, err = tappExc(err, exception.ExceptionTypeInternalError)
		}

		if err != nil {
			return nil, err
		}
	}

	var buff bytes.Buffer
	if err := s.p.EncodeEnveloped(response, &buff); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

// Helper to build TApplicationException wire.Values
func tappExc(err error, typ exception.ExceptionType) (wire.Value, error) {
	return (&exception.TApplicationException{
		Message: ptr.String(err.Error()),
		Type:    &typ,
	}).ToWire()
}
