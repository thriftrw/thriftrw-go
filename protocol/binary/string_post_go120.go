// Copyright (c) 2024 Uber Technologies, Inc.
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

//go:build go1.20
// +build go1.20

package binary

import "unsafe"

// ReadString reads a Thrift encoded string.
func (sr *StreamReader) ReadString() (string, error) {
	bs, err := sr.ReadBinary()
	// It is safe to use "unsafe" here because there are no
	// mutable references to bs.
	return unsafe.String(unsafe.SliceData(bs), len(bs)), err
}

// WriteString encodes a string
func (sw *StreamWriter) WriteString(s string) error {
	if err := sw.WriteInt32(int32(len(s))); err != nil {
		return err
	}
	// It is safe to use "unsafe" here because there are no
	// mutable references to the byte slice b.
	// sw.write() delegates to the underlying io.Writer,
	// and according to its documentation, "Write must
	// not modify the slice data, even temporarily."
	b := unsafe.Slice(unsafe.StringData(s), len(s))
	return sw.write(b)
}
