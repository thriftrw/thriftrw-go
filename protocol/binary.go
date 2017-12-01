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
func (b binaryProtocol) DecodeRequest(r io.ReaderAt) (wire.Value, EnvelopeSpecificResponder, error) {
	// Strip request envelopes if present.
	// 1. A message of length 1 containing only 0x00 can only be a bare,
	// un-enveloped empty struct.
	// 2. A message of length >1 starting with 0x00 can only be enveloped
	// without a version (non-strict)
	// 3. A message of length >1 with MSB 0x80 set can only be enveloped with a
	// version (strict)
	// 4. A message of length >1 starting with 0x00-0x0f can only be the
	// beginning of an un-enveloped field.
	var buf [2]byte
	// If we fail to read two bytes, the only possible valid value is the empty struct.
	if _, err := r.ReadAt(buf[0:2], 0); err != nil {
		val, err := b.Decode(r, wire.TStruct)
		if err != nil {
			return wire.Value{}, noEnvelopeResponder, err
		}
		return val, noEnvelopeResponder, nil
	}

	// 0x00 is a valid preamble for a non-strict enveloped request.
	if buf[0] == 0x00 {
		e, err := b.DecodeEnveloped(r)
		if err != nil {
			return wire.Value{}, noEnvelopeResponder, err
		}
		return e.Value, &envelopeV0Responder{
			Name:  e.Name,
			SeqID: e.SeqID,
		}, nil
	}

	// 0x80 is a valid preamble for a strict enveloped request.
	if buf[0]&0x80 > 0 {
		e, err := b.DecodeEnveloped(r)
		if err != nil {
			return wire.Value{}, noEnvelopeResponder, err
		}
		return e.Value, &envelopeV1Responder{
			Name:  e.Name,
			SeqID: e.SeqID,
		}, nil
	}

	// All other patterns are either bare structs or invalid.
	val, err := b.Decode(r, wire.TStruct)
	if err != nil {
		return wire.Value{}, noEnvelopeResponder, err
	}
	return val, noEnvelopeResponder, nil
}

// noEnvelopeResponder responds to a request without an envelope.

type _noEnvelopeResponder struct{}

func (_noEnvelopeResponder) EncodeResponse(v wire.Value, t wire.EnvelopeType, w io.Writer) error {
	return Binary.Encode(v, w)
}

var noEnvelopeResponder EnvelopeSpecificResponder = &_noEnvelopeResponder{}

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
