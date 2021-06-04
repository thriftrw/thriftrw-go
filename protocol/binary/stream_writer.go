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
	"io"
	"math"
	"sync"

	"go.uber.org/thriftrw/internal/iface"
	"go.uber.org/thriftrw/protocol/stream"
)

var streamWriterPool = sync.Pool{
	New: func() interface{} {
		writer := &StreamWriter{}
		return writer
	}}

// StreamWriter implements basic logic for writing the Thrift Binary Protocol
// to an io.Writer.
type StreamWriter struct {
	// Private implementation to disallow custom implementations of
	// the Writer interface
	iface.Impl

	writer io.Writer

	// This buffer is re-used every time we need a slice of up to 8 bytes.
	buffer [8]byte
}

// BorrowStreamWriter fetches a StreamWriter from the system that will write
// its output to the given io.Writer.
//
// This StreamWriter must be returned back using ReturnStreamWriter.
func BorrowStreamWriter(w io.Writer) *StreamWriter {
	streamWriter := streamWriterPool.Get().(*StreamWriter)
	streamWriter.writer = w
	return streamWriter
}

// ReturnStreamWriter returns a previously borrowed StreamWriter back to the
// system.
func ReturnStreamWriter(sw *StreamWriter) {
	sw.writer = nil
	streamWriterPool.Put(sw)
}

func (bw *StreamWriter) write(bs []byte) error {
	_, err := bw.writer.Write(bs)
	return err
}

func (bw *StreamWriter) writeByte(b byte) error {
	bs := bw.buffer[0:1]
	bs[0] = b
	return bw.write(bs)
}

func (bw *StreamWriter) writeInt16(n int16) error {
	bs := bw.buffer[0:2]
	bigEndian.PutUint16(bs, uint16(n))
	return bw.write(bs)
}

func (bw *StreamWriter) writeInt32(n int32) error {
	bs := bw.buffer[0:4]
	bigEndian.PutUint32(bs, uint32(n))
	return bw.write(bs)
}

func (bw *StreamWriter) writeInt64(n int64) error {
	bs := bw.buffer[0:8]
	bigEndian.PutUint64(bs, uint64(n))
	return bw.write(bs)
}

func (bw *StreamWriter) writeString(s string) error {
	if err := bw.writeInt32(int32(len(s))); err != nil {
		return err
	}

	_, err := io.WriteString(bw.writer, s)
	return err
}

// WriteBool encodes a boolean
func (bw *StreamWriter) WriteBool(b bool) error {
	if b {
		return bw.writeByte(1)
	}
	return bw.writeByte(0)
}

// WriteInt8 encodes an int8
func (bw *StreamWriter) WriteInt8(i int8) error {
	return bw.writeByte(byte(i))
}

// WriteInt16 encodes an int16
func (bw *StreamWriter) WriteInt16(i int16) error {
	return bw.writeInt16(i)
}

// WriteInt32 encodes an int32
func (bw *StreamWriter) WriteInt32(i int32) error {
	return bw.writeInt32(i)
}

// WriteInt64 encodes an int64
func (bw *StreamWriter) WriteInt64(i int64) error {
	return bw.writeInt64(i)
}

// WriteString encodes a string
func (bw *StreamWriter) WriteString(s string) error {
	return bw.writeString(s)
}

// WriteDouble encodes a double
func (bw *StreamWriter) WriteDouble(d float64) error {
	value := math.Float64bits(d)
	return bw.writeInt64(int64(value))
}

// WriteBinary encodes binary
func (bw *StreamWriter) WriteBinary(b []byte) error {
	if err := bw.writeInt32(int32(len(b))); err != nil {
		return err
	}
	return bw.write(b)
}

// WriteFieldBegin marks the beginning of a new field in a struct. The first
// byte denotes the type and the next two bytes denote the field id.
func (bw *StreamWriter) WriteFieldBegin(f stream.FieldHeader) error {
	// type:1
	if err := bw.writeByte(byte(f.Type)); err != nil {
		return err
	}

	// id:2
	if err := bw.writeInt16(f.ID); err != nil {
		return err
	}

	return nil
}

// WriteFieldEnd denotes the end of a field. No-op.
func (bw *StreamWriter) WriteFieldEnd() error {
	return nil
}

// WriteStructBegin denotes the beginning of a struct. No-op.
func (bw *StreamWriter) WriteStructBegin() error {
	return nil
}

// WriteStructEnd uses the zero byte to mark the end of a struct.
func (bw *StreamWriter) WriteStructEnd() error {
	return bw.writeByte(0) // end struct
}

// WriteListBegin marks the beginning of a new list. The first byte denotes
// the type of the items and the next four bytes denote the length of the list.
func (bw *StreamWriter) WriteListBegin(l stream.ListHeader) error {
	// vtype:1
	if err := bw.writeByte(byte(l.Type)); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(l.Length)); err != nil {
		return err
	}

	return nil
}

// WriteListEnd marks the end of a list. No-op.
func (bw *StreamWriter) WriteListEnd() error {
	return nil
}

// WriteSetBegin marks the beginning of a new set. The first byte denotes
// the type of the items and the next four bytes denote the length of the set.
func (bw *StreamWriter) WriteSetBegin(s stream.SetHeader) error {
	// vtype:1
	if err := bw.writeByte(byte(s.Type)); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(s.Length)); err != nil {
		return err
	}

	return nil
}

// WriteSetEnd marks the end of a set. No-op.
func (bw *StreamWriter) WriteSetEnd() error {
	return nil
}

// WriteMapBegin marks the beginning of a new map. The first byte denotes
// the type of the keys, the second byte denotes the type of the values,
// and the next four bytes denote the length of the map.
func (bw *StreamWriter) WriteMapBegin(m stream.MapHeader) error {
	// ktype:1
	if err := bw.writeByte(byte(m.KeyType)); err != nil {
		return err
	}

	// vtype:1
	if err := bw.writeByte(byte(m.ValueType)); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(m.Length)); err != nil {
		return err
	}

	return nil
}

// WriteMapEnd marks the end of a map. No-op.
func (bw *StreamWriter) WriteMapEnd() error {
	return nil
}
