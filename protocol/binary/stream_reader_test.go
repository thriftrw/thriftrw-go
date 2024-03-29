// Copyright (c) 2023 Uber Technologies, Inc.
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
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkReadString(b *testing.B) {
	var streamBuff bytes.Buffer

	// Encode with Streaming protocol
	w := NewStreamWriter(&streamBuff)
	require.NoError(b, w.WriteString("the quick brown fox jumps over the lazy dog"))
	require.NoError(b, w.Close())

	encoded := streamBuff.Bytes()

	enc := bytes.NewReader(encoded)
	sr := NewStreamReader(enc)
	defer sr.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sr.ReadString()
		enc.Seek(0, io.SeekStart)
	}
}
