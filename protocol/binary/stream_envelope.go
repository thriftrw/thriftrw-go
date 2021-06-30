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
	"errors"

	"go.uber.org/thriftrw/protocol/stream"
)

// WriteEnvelopeBegin writes the start of a strict envelope (contains an envelope version).
func (sw *StreamWriter) WriteEnvelopeBegin(eh stream.EnvelopeHeader) error {
	return errors.New("not implemented")
}

// WriteEnvelopeEnd writes the "end" of an envelope. Since there is no ending
// to an envelope, this is a no-op.
func (sw *StreamWriter) WriteEnvelopeEnd() error {
	return errors.New("not implemented")
}

// ReadEnvelopeBegin reads the start of an Apache Thrift envelope. Thrift supports
// two kinds of envelopes: strict, and non-strict. See ReadEnveloped method
// for more information on enveloping.
func (sw *StreamReader) ReadEnvelopeBegin() (stream.EnvelopeHeader, error) {
	return stream.EnvelopeHeader{}, errors.New("not implemented")
}

// ReadEnvelopeEnd reads the "end" of an envelope.  Since there is no real
// envelope end, this is a no-op.
func (sw *StreamReader) ReadEnvelopeEnd() error {
	return errors.New("not implemented")
}
