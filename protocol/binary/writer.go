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

// Package binary implements the Thrift Binary protocol.
package binary

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/uber/thriftrw-go/wire"
)

var bigEndian = binary.BigEndian

// Writer implements basic logic for writing the the Thrift Binary Protocol to
// an io.Writer.
type Writer struct {
	Writer io.Writer

	// This buffer is re-used every time we need a slice of up to 8 bytes.
	buffer [8]byte
}

// Write writes the given slice of bytes.
func (bw *Writer) Write(bs []byte) error {
	_, err := bw.Writer.Write(bs)
	return err
}

// WriteByte writes out a single byte.
func (bw *Writer) WriteByte(b byte) error {
	bs := bw.buffer[0:1]
	bs[0] = b
	return bw.Write(bs)
}

// WriteInt16 writes the given 16-bit integer using big endian byte ordering.
func (bw *Writer) WriteInt16(n int16) error {
	bs := bw.buffer[0:2]
	bigEndian.PutUint16(bs, uint16(n))
	return bw.Write(bs)
}

// WriteInt32 writes the given 32-bit integer using big endian byte ordering.
func (bw *Writer) WriteInt32(n int32) error {
	bs := bw.buffer[0:4]
	bigEndian.PutUint32(bs, uint32(n))
	return bw.Write(bs)
}

// WriteInt64 writes the given 64-bit integer using big endian byte ordering.
func (bw *Writer) WriteInt64(n int64) error {
	bs := bw.buffer[0:8]
	bigEndian.PutUint64(bs, uint64(n))
	return bw.Write(bs)
}

// WriteValue writes out the given Thrift value.
func (w *Writer) WriteValue(v wire.Value) error {
	switch v.Type {
	case wire.TBool:
		if v.Bool {
			return w.WriteByte(1)
		} else {
			return w.WriteByte(0)
		}

	case wire.TByte:
		return w.WriteByte(byte(v.Byte))

	case wire.TDouble:
		value := math.Float64bits(v.Double)
		return w.WriteInt64(int64(value))

	case wire.TI16:
		return w.WriteInt16(v.I16)

	case wire.TI32:
		return w.WriteInt32(v.I32)

	case wire.TI64:
		return w.WriteInt64(v.I64)

	case wire.TBinary:
		if err := w.WriteInt32(int32(len(v.Binary))); err != nil {
			return err
		}
		return w.Write(v.Binary)

	case wire.TStruct:
		for _, f := range v.Struct.Fields {
			// type:1
			if err := w.WriteByte(byte(f.Value.Type)); err != nil {
				return err
			}

			// id:2
			if err := w.WriteInt16(f.ID); err != nil {
				return err
			}

			// value
			if err := w.WriteValue(f.Value); err != nil {
				return fmt.Errorf(
					"failed to write field %d (%v): %s",
					f.ID, f.Value.Type, err,
				)
			}

		}
		return w.WriteByte(0) // end struct

	case wire.TMap:
		// ktype:1
		if err := w.WriteByte(byte(v.Map.KeyType)); err != nil {
			return err
		}

		// vtype:1
		if err := w.WriteByte(byte(v.Map.ValueType)); err != nil {
			return err
		}

		// length:4
		if err := w.WriteInt32(int32(len(v.Map.Items))); err != nil {
			return err
		}

		for _, item := range v.Map.Items {
			if err := w.WriteValue(item.Key); err != nil {
				return err
			}
			if err := w.WriteValue(item.Value); err != nil {
				return err
			}
		}

		return nil

	case wire.TSet:
		// vtype:1
		if err := w.WriteByte(byte(v.Set.ValueType)); err != nil {
			return err
		}

		// length:4
		if err := w.WriteInt32(int32(len(v.Set.Items))); err != nil {
			return err
		}

		for _, item := range v.Set.Items {
			if err := w.WriteValue(item); err != nil {
				return err
			}
		}

		return nil

	case wire.TList:
		// vtype:1
		if err := w.WriteByte(byte(v.List.ValueType)); err != nil {
			return err
		}

		// length:4
		if err := w.WriteInt32(int32(len(v.List.Items))); err != nil {
			return err
		}

		for _, item := range v.List.Items {
			if err := w.WriteValue(item); err != nil {
				return err
			}
		}

		return nil

	default:
		return fmt.Errorf("unknown ttype %v", v.Type)
	}
}
