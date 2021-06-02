// Copyright (c) 2021 Uber Technologies, Inc.
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

// Package protocol provides implementations of different Thrift protocols.
package protocol

import (
	"io"

	"go.uber.org/thriftrw/wire"
)

// Protocol defines a specific way for a Thrift value to be encoded or
// decoded.
type Protocol interface {
	// Encode the given Value and write the result to the given Writer.
	Encode(v wire.Value, w io.Writer) error

	// EncodeEnveloped encodes the enveloped value and writes the result
	// to the given Writer.
	EncodeEnveloped(e wire.Envelope, w io.Writer) error

	// Decode reads a Value of the given type from the given Reader.
	Decode(r io.ReaderAt, t wire.Type) (wire.Value, error)

	// DecodeEnveloped reads an enveloped value from the given Reader.
	// Enveloped values are assumed to be TStructs.
	DecodeEnveloped(r io.ReaderAt) (wire.Envelope, error)
}

// EnvelopeAgnosticProtocol defines a specific way for a Thrift value to be
// encoded or decoded, additionally being able to decode requests without prior
// knowledge of whether the request is enveloped.
//
// The Binary protocol in particular can be upcast to EnvelopeAgnosticProtocol.
type EnvelopeAgnosticProtocol interface {
	Protocol

	// DecodeRequest reads an enveloped or un-enveloped struct from the given
	// ReaderAt.
	// This allows a Thrift request handler to transparently accept requests
	// regardless of whether the caller is configured to submit envelopes.
	// The caller specifies the expected envelope type, one of OneWay or Unary,
	// on which the decoder asserts if the envelope is present.
	DecodeRequest(et wire.EnvelopeType, r io.ReaderAt) (wire.Value, Responder, error)
}

// Responder captures how to respond to a request, concerning whether and what
// kind of envelope to use, how to match the sequence identifier of the
// corresponding request.
type Responder interface {
	// EncodeResponse writes a response value to the writer, with the envelope
	// style of the corresponding request.
	// The EnvelopeType should be either wire.Reply or wire.Exception.
	EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error
}

// StreamProtocol defines a specific way for a Thrift value to be encoded or
// decoded, implemented in a streaming fashion.
type StreamProtocol interface {
	// StreamWriter returns a streaming implementation of an encoder for a
	// Thrift value.
	StreamWriter(w io.Writer) Writer
	// StreamReader returns a streaming implementation of a decoder for a
	// Thrift value.
	StreamReader(r io.Reader) Reader
}

// FieldHeader defines the metadata needed to define the beginning of a field
// in a Thrift value.
type FieldHeader struct {
	ID   int16
	Type wire.Type
}

// MapHeader defines the metadata needed to define the beginning of a map in a
// Thrift value.
type MapHeader struct {
	KeyType   wire.Type
	ValueType wire.Type
	Length    int
}

// SetHeader defines the metadata needed to define the beginning of a set in a
// Thrift value.
type SetHeader struct {
	Length int
	Type   wire.Type
}

// ListHeader defines the metadata needed to define the beginning of a list in a
// Thrift value.
type ListHeader struct {
	Length int
	Type   wire.Type
}

// Writer defines an encoder for a Thrift value, implemented in a streaming
// fashion.
type Writer interface {
	WriteBool(b bool) error
	WriteInt8(i int8) error
	WriteInt16(i int16) error
	WriteInt32(i int32) error
	WriteInt64(i int64) error
	WriteString(s string) error
	WriteDouble(f float64) error
	WriteBinary(b []byte) error
	WriteStructBegin() error
	WriteStructEnd() error
	WriteFieldBegin(f FieldHeader) error
	WriteFieldEnd() error
	WriteMapBegin(m MapHeader) error
	WriteMapEnd() error
	WriteSetBegin(s SetHeader) error
	WriteSetEnd() error
	WriteListBegin(l ListHeader) error
	WriteListEnd() error
}

// Reader defines an decoder for a Thrift value, implemented in a streaming
// fashion.
type Reader interface {
	ReadBool() (bool, error)
	ReadInt8() (int8, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt64() (int64, error)
	ReadString() (string, error)
	ReadDouble() (float64, error)
	ReadBinary() ([]byte, error)
	ReadStructBegin() error
	ReadStructEnd() error
	ReadFieldBegin() (FieldHeader, bool, error)
	ReadFieldEnd() error
	ReadListBegin() (ListHeader, error)
	ReadListEnd() error
	ReadSetBegin() (SetHeader, error)
	ReadSetEnd() error
	ReadMapBegin() (MapHeader, error)
	ReadMapEnd() error
	// Skip skips over the bytes of the wire type and any applicable headers.
	Skip(w wire.Type) error
}
