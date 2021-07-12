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

package binary

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/internal/iface"
	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

// Default is the default implementation of the Thrift Binary Protocol.
var Default = new(Protocol)

// Protocol implements the Thrift Binary Protocol.
type Protocol struct {
	iface.Private
}

// Encode the given Value and write the result to the given Writer.
func (*Protocol) Encode(v wire.Value, w io.Writer) error {
	writer := BorrowWriter(w)
	err := writer.WriteValue(v)
	ReturnWriter(writer)
	return err
}

// Decode reads a Value of the given type from the given Reader.
func (*Protocol) Decode(r io.ReaderAt, t wire.Type) (wire.Value, error) {
	reader := NewReader(r)
	value, _, err := reader.ReadValue(t, 0)
	return value, err
}

// Writer builds a stream writer that writes to the provided stream using the
// Thrift Binary Protocol.
func (*Protocol) Writer(w io.Writer) stream.Writer {
	return NewStreamWriter(w)
}

// Reader builds a stream reader that reads from the provided stream using the
// Thrift Binary Protocol.
func (*Protocol) Reader(r io.Reader) stream.Reader {
	return NewStreamReader(r)
}

// EncodeEnveloped encodes the enveloped value and writes the result
// to the given Writer.
func (*Protocol) EncodeEnveloped(e wire.Envelope, w io.Writer) error {
	writer := BorrowWriter(w)
	err := writer.WriteEnveloped(e)
	ReturnWriter(writer)
	return err
}

// DecodeEnveloped reads an enveloped value from the given Reader.
// Enveloped values are assumed to be TStructs.
func (*Protocol) DecodeEnveloped(r io.ReaderAt) (wire.Envelope, error) {
	reader := NewReader(r)
	e, err := reader.ReadEnveloped()
	return e, err
}

// DecodeRequest specializes Decode and replaces DecodeEnveloped for the
// specific purpose of decoding request structs that may or may not have an
// envelope.
// This allows a Thrift request handler to transparently accept requests
// regardless of whether the caller submits an envelope.
// The caller specifies the expected envelope type, one of OneWay or Unary, on
// which the decoder asserts if the envelope is present.
//
// This is possible because we can distinguish an envelope from a bare request
// struct by looking at the first byte and the length of the message.
//
// 1. A message of length 1 containing only 0x00 can only be an empty struct.
// 0x00 is the type ID for STOP, indicating the end of the struct.
//
// 2. A message of length >1 starting with 0x00 can only be a non-strict
// envelope (not versioned), assuming the message name is less than 16MB long.
// In this case, the first four bytes indicate the length of the method name,
// which is unlikely to overflow into the high byte.
//
// 3. A message of length >1, where the first byte is <0 can only be a strict envelope.
// The MSB indicates that the message is versioned. Reading the first two bytes
// and masking out the MSB indicates the version number.
// At this time, there is only one version.
//
// 4. A message of length >1, where the first byte is >=0 can only be a bare
// struct starting with that field identifier. Valid field identifiers today
// are in the range 0x00-0x0f. There is some chance that a future version of
// the protocol will add more field types, but it is very unlikely that the
// field type will flow into the MSB (128 type identifiers, starting with the
// 15 valid types today).
func (p *Protocol) DecodeRequest(et wire.EnvelopeType, r io.ReaderAt) (wire.Value, Responder, error) {
	var buf [2]byte

	// If we fail to read two bytes, the only possible valid value is the empty struct.
	if count, _ := r.ReadAt(buf[0:2], 0); count < 2 {
		val, err := p.Decode(r, wire.TStruct)
		return val, NoEnvelopeResponder, err
	}

	// If length > 1, 0x00 is only a valid preamble for a non-strict enveloped request.
	if buf[0] == 0x00 {
		e, err := p.DecodeEnveloped(r)
		if err != nil {
			return wire.Value{}, NoEnvelopeResponder, err
		}
		if e.Type != et {
			return wire.Value{}, NoEnvelopeResponder, errUnexpectedEnvelopeType(e.Type)
		}
		return e.Value, &EnvelopeV0Responder{
			Name:  e.Name,
			SeqID: e.SeqID,
		}, nil
	}

	// Only strict (versioned) envelopes begin with the most significant bit set.
	// This could only be confused for a type identifier greater than 127
	// (beyond the 15 Thrift has at time of writing), or a message name longer
	// than 16MB.
	if buf[0]&0x80 > 0 {
		e, err := p.DecodeEnveloped(r)
		if err != nil {
			return wire.Value{}, NoEnvelopeResponder, err
		}
		if e.Type != et {
			return wire.Value{}, NoEnvelopeResponder, errUnexpectedEnvelopeType(e.Type)
		}
		return e.Value, &EnvelopeV1Responder{
			Name:  e.Name,
			SeqID: e.SeqID,
		}, nil
	}

	// All other patterns are either bare structs or invalid.
	// We delegate to the struct decoder to distinguish invalid type
	// identifiers, outside the 0-15 range.
	val, err := p.Decode(r, wire.TStruct)
	return val, NoEnvelopeResponder, err
}

// Responder captures how to respond to a request, concerning whether and what
// kind of envelope to use, how to match the sequence identifier of the
// corresponding request.
type Responder interface {
	EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error
}

// noEnvelopeResponder responds to a request without an envelope.
type noEnvelopeResponder struct{}

func (noEnvelopeResponder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	return Default.Encode(v, w)
}

// NoEnvelopeResponder responds to a request without an envelope.
var NoEnvelopeResponder = &noEnvelopeResponder{}

// EnvelopeV0Responder responds to requests with a non-strict (unversioned) envelope.
type EnvelopeV0Responder struct {
	Name  string
	SeqID int32
}

// EncodeResponse writes the response to the writer using a non-strict
// envelope.
func (r EnvelopeV0Responder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	writer := BorrowWriter(w)
	err := writer.WriteLegacyEnveloped(wire.Envelope{
		Name:  r.Name,
		Type:  t,
		SeqID: r.SeqID,
		Value: v,
	})
	ReturnWriter(writer)
	return err
}

// EnvelopeV1Responder responds to requests with a strict, version 1 envelope.
type EnvelopeV1Responder struct {
	Name  string
	SeqID int32
}

// EncodeResponse writes the response to the writer using a strict, version 1
// envelope.
func (r EnvelopeV1Responder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	writer := BorrowWriter(w)
	err := writer.WriteEnveloped(wire.Envelope{
		Name:  r.Name,
		Type:  t,
		SeqID: r.SeqID,
		Value: v,
	})
	ReturnWriter(writer)
	return err
}

type errUnexpectedEnvelopeType wire.EnvelopeType

func (e errUnexpectedEnvelopeType) Error() string {
	return fmt.Sprintf("unexpected envelope type: %v", wire.EnvelopeType(e))
}