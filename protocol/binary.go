// Copyright (c) 2017 Uber Technologies, Inc.
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

package protocol

import (
	"fmt"
	"io"

	"go.uber.org/thriftrw/protocol/binary"
	"go.uber.org/thriftrw/wire"
)

// Binary implements the Thrift Binary Protocol.
// Binary can be cast up to EnvelopeAgnosticProtocol to support DecodeRequest.
var Binary Protocol

// EnvelopeAgnosticBinary implements the Thrift Binary Protocol, using
// DecodeRequest for request bodies that may or may not have an envelope.
// This in turn produces a responder with an EncodeResponse method so a handler
// can reply in-kind.
//
// EnvelopeAgnosticBinary makes some practical assumptions about messages
// to be able to distinguish enveloped, versioned enveloped, and not-enveloped
// messages reliably:
//
//  1. No message will use an envelope version greater than 0x7fff.  This would
//  flip the bit that makes versioned envelopes recognizable.  The only
//  envelope version we recognize today is version 1.
//
//  2. No message with an unversioned envelope will have a message name
//  (procedure name) longer than 0x00ffffffff (three bytes of length prefix).
//  This would roll into the byte that distinguishes the type of the first
//  field of an un-enveloped struct.  This would require a 16MB procedure name.
var EnvelopeAgnosticBinary EnvelopeAgnosticProtocol

type errUnexpectedEnvelopeType wire.EnvelopeType

func (e errUnexpectedEnvelopeType) Error() string {
	return fmt.Sprintf("unexpected envelope type: %v", wire.EnvelopeType(e))
}

func init() {
	Binary = binaryProtocol{}
	EnvelopeAgnosticBinary = binaryProtocol{}
}

type binaryProtocol struct{}

func (binaryProtocol) Encode(v wire.Value, w io.Writer) error {
	writer := binary.BorrowWriter(w)
	err := writer.WriteValue(v)
	binary.ReturnWriter(writer)
	return err
}

func (binaryProtocol) Decode(r io.ReaderAt, t wire.Type) (wire.Value, error) {
	reader := binary.NewReader(r)
	value, _, err := reader.ReadValue(t, 0)
	return value, err
}

func (binaryProtocol) EncodeEnveloped(e wire.Envelope, w io.Writer) error {
	writer := binary.BorrowWriter(w)
	err := writer.WriteEnveloped(e)
	binary.ReturnWriter(writer)
	return err
}

func (binaryProtocol) DecodeEnveloped(r io.ReaderAt) (wire.Envelope, error) {
	reader := binary.NewReader(r)
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
func (b binaryProtocol) DecodeRequest(et wire.EnvelopeType, r io.ReaderAt) (wire.Value, EnvelopeSpecificResponder, error) {
	var buf [2]byte

	// If we fail to read two bytes, the only possible valid value is the empty struct.
	if _, err := r.ReadAt(buf[0:2], 0); err != nil {
		val, err := b.Decode(r, wire.TStruct)
		if err != nil {
			return wire.Value{}, _noEnvelopeResponder, err
		}
		return val, _noEnvelopeResponder, nil
	}

	// If length > 1, 0x00 is only a valid preamble for a non-strict enveloped request.
	if buf[0] == 0x00 {
		e, err := b.DecodeEnveloped(r)
		if err != nil {
			return wire.Value{}, _noEnvelopeResponder, err
		}
		if e.Type != et {
			return wire.Value{}, _noEnvelopeResponder, errUnexpectedEnvelopeType(e.Type)
		}
		return e.Value, &envelopeV0Responder{
			Name:  e.Name,
			SeqID: e.SeqID,
		}, nil
	}

	// Only strict (versioned) envelopes begin with the most significant bit set.
	// This could only be confused for a type identifier greater than 127
	// (beyond the 15 Thrift has at time of writing), or a message name longer
	// than 16MB.
	if buf[0]&0x80 > 0 {
		e, err := b.DecodeEnveloped(r)
		if err != nil {
			return wire.Value{}, _noEnvelopeResponder, err
		}
		if e.Type != et {
			return wire.Value{}, _noEnvelopeResponder, errUnexpectedEnvelopeType(e.Type)
		}
		return e.Value, &envelopeV1Responder{
			Name:  e.Name,
			SeqID: e.SeqID,
		}, nil
	}

	// All other patterns are either bare structs or invalid.
	// We delegate to the struct decoder to distinguish invalid type
	// identifiers, outside the 0-15 range.
	val, err := b.Decode(r, wire.TStruct)
	if err != nil {
		return wire.Value{}, _noEnvelopeResponder, err
	}
	return val, _noEnvelopeResponder, nil
}

// noEnvelopeResponder responds to a request without an envelope.
type noEnvelopeResponder struct{}

func (noEnvelopeResponder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	return Binary.Encode(v, w)
}

var _noEnvelopeResponder EnvelopeSpecificResponder = &noEnvelopeResponder{}

// envelopeV0Responder responds to requests with a non-strict (unversioned) envelope.
type envelopeV0Responder struct {
	Name  string
	SeqID int32
}

func (r envelopeV0Responder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	writer := binary.BorrowWriter(w)
	err := writer.WriteLegacyEnveloped(wire.Envelope{
		Name:  r.Name,
		Type:  t,
		SeqID: r.SeqID,
		Value: v,
	})
	binary.ReturnWriter(writer)
	return err
}

// envelopeV1Responder responds to requests with a strict, version 1 envelope.
type envelopeV1Responder struct {
	Name  string
	SeqID int32
}

func (r envelopeV1Responder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	writer := binary.BorrowWriter(w)
	err := writer.WriteEnveloped(wire.Envelope{
		Name:  r.Name,
		Type:  t,
		SeqID: r.SeqID,
		Value: v,
	})
	binary.ReturnWriter(writer)
	return err
}
