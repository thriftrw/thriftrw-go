// Copyright (c) 2015 Uber Technologies, Inc.
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

package gen

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/protocol/binary"
	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

// thriftType is implemented by all generated types that know how to encode
// and decode themselves.
//
// Having this interface allows tests that use reflection to be more readable
// by relying on the interface rather than reflection to call methods.
type thriftType interface {
	fmt.Stringer

	ToWire() (wire.Value, error)
	FromWire(wire.Value) error
	Encode(stream.Writer) error
	Decode(stream.Reader) error
}

func streamDecodeWireType(t *testing.T, wv wire.Value, tt thriftType) error {
	t.Helper()

	var buf bytes.Buffer
	require.NoError(t, binary.Default.Encode(wv, &buf))

	r := bytes.NewReader(buf.Bytes())
	sr := binary.Default.Reader(r)
	defer func() {
		assert.NoError(t, sr.Close())
	}()

	// Only IO errors would cause Decode to error early - since we shouldn't
	// expect any of those, the full contents of the raw bytes should be read out.
	defer func() {
		assert.Zero(t, r.Len(), "expected to be end of read")
	}()

	return tt.Decode(sr)
}
