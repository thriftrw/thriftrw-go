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
	"io"
	"sync"
)

// Client provides bidirectional outgoing framed communication.
//
// It allows sending framed requests where each request has a corresponding
// response. Only one active request is allowed at a time. Other requests are
// blocked while a request is ongoing.
type Client struct {
	sync.Mutex

	r *Reader
	w *Writer
}

// NewClient builds a new Client which uses the given writer to send requests
// and the given reader to read their responses.
func NewClient(w io.Writer, r io.Reader) *Client {
	return &Client{
		r: NewReader(r),
		w: NewWriter(w),
	}
}

// Send sends the given frame and returns its response.
func (c *Client) Send(b []byte) ([]byte, error) {
	c.Lock()
	defer c.Unlock()

	if err := c.w.Write(b); err != nil {
		return nil, err
	}
	return c.r.Read()
}
