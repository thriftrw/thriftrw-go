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

package frame

import (
	"encoding/binary"
	"io"
	"sync"

	"go.uber.org/atomic"
)

// Writer is a writer for framed messages.
type Writer struct {
	sync.Mutex

	closed atomic.Bool
	w      io.Writer
	buff   [4]byte
}

// NewWriter builds a new Writer which writes frames to the given io.Writer.
//
// If the io.Writer is a WriteCloser, its Close method will be called when the
// frame.Writer is closed.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Write writes the given frame to the Writer.
func (w *Writer) Write(b []byte) error {
	w.Lock()
	defer w.Unlock()

	// TODO(abg): Bounds check?
	binary.BigEndian.PutUint32(w.buff[:], uint32(len(b)))
	if _, err := w.w.Write(w.buff[:]); err != nil {
		return err
	}

	if len(b) == 0 {
		return nil
	}

	_, err := w.w.Write(b)
	return err
}

// Close closes the given Writeer.
func (w *Writer) Close() error {
	if w.closed.Swap(true) {
		return nil // already closed
	}
	if c, ok := w.w.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
