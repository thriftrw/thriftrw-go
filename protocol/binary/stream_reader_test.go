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

package binary

import "testing"

// byteGenerator generates an infinite stream of val
type byteGenerator struct {
	val byte
}

// Read is the io.Reader implementation
func (bg *byteGenerator) Read(p []byte) (n int, err error) {
	p[0] = bg.val
	return 1, nil
}

// ReadByte is the io.ByteReader implementation
func (bg *byteGenerator) ReadByte() (byte, error) {
	return bg.val, nil
}

func BenchmarkReadInt8(b *testing.B) {
	input := &byteGenerator{val: 0xff}
	reader := NewStreamReader(input)
	defer reader.Close()

	for i:= 0; i<b.N; i++ {
		reader.ReadInt8()
	}
}

func BenchmarkReadInt16(b *testing.B) {
	input := &byteGenerator{val: 0xff}
	reader := NewStreamReader(input)
	defer reader.Close()

	for i:= 0; i<b.N; i++ {
		reader.ReadInt16()
	}
}

func BenchmarkReadInt32(b *testing.B) {
	input := &byteGenerator{val: 0xff}
	reader := NewStreamReader(input)
	defer reader.Close()

	for i:= 0; i<b.N; i++ {
		reader.ReadInt32()
	}
}

func BenchmarkReadInt64(b *testing.B) {
	input := &byteGenerator{val: 0xff}
	reader := NewStreamReader(input)
	defer reader.Close()

	for i:= 0; i<b.N; i++ {
		reader.ReadInt64()
	}
}
