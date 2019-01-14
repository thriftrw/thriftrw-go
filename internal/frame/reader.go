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
	"bytes"
	"encoding/binary"
	"io"
	"sync"

	"go.uber.org/atomic"
)

// Maximum frame size for which we pre-allocate buffers.
var _fastPathFrameSize int64 = 10 * 1024 * 1024 // 10 MB

// Reader is a reader for framed messages.
type Reader struct {
	sync.Mutex

	closed atomic.Bool
	r      io.Reader
	buff   [4]byte
}

// NewReader builds a new Reader which reads frames from the given io.Reader.
//
// If the io.Reader is a ReadCloser, its Close method will be called when the
// frame.Reader is closed.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// Read reads the next frame from the Reader.
func (r *Reader) Read() ([]byte, error) {
	r.Lock()
	defer r.Unlock()

	if _, err := io.ReadFull(r.r, r.buff[:]); err != nil {
		return nil, err
	}

	length := int64(binary.BigEndian.Uint32(r.buff[:]))
	if length < _fastPathFrameSize {
		return r.readFastPath(length)
	}

	var buff bytes.Buffer
	_, err := io.CopyN(&buff, r.r, length)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}

func (r *Reader) readFastPath(l int64) ([]byte, error) {
	buff := make([]byte, l)
	if l == 0 {
		return buff, nil
	}
	_, err := io.ReadFull(r.r, buff)
	return buff, err
}

// Close closes the given Reader.
func (r *Reader) Close() error {
	if r.closed.Swap(true) {
		return nil // already closed
	}

	if c, ok := r.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
