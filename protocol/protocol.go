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
