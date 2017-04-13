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
var Binary Protocol

func init() {
	Binary = binaryProtocol{}
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
