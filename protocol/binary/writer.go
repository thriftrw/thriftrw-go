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

// Writer implements basic logic for writing the Thrift Binary Protocol to an
// io.Writer.
type Writer struct {
	Writer io.Writer

	// This buffer is re-used every time we need a slice of up to 8 bytes.
	buffer [8]byte
}

func (bw *Writer) write(bs []byte) error {
	_, err := bw.Writer.Write(bs)
	return err
}

func (bw *Writer) writeByte(b byte) error {
	bs := bw.buffer[0:1]
	bs[0] = b
	return bw.write(bs)
}

func (bw *Writer) writeInt16(n int16) error {
	bs := bw.buffer[0:2]
	bigEndian.PutUint16(bs, uint16(n))
	return bw.write(bs)
}

func (bw *Writer) writeInt32(n int32) error {
	bs := bw.buffer[0:4]
	bigEndian.PutUint32(bs, uint32(n))
	return bw.write(bs)
}

func (bw *Writer) writeInt64(n int64) error {
	bs := bw.buffer[0:8]
	bigEndian.PutUint64(bs, uint64(n))
	return bw.write(bs)
}

func (bw *Writer) writeField(f wire.Field) error {
	// type:1
	if err := bw.writeByte(byte(f.Value.Type)); err != nil {
		return err
	}

	// id:2
	if err := bw.writeInt16(f.ID); err != nil {
		return err
	}

	// value
	if err := bw.WriteValue(f.Value); err != nil {
		return fmt.Errorf(
			"failed to write field %d (%v): %s",
			f.ID, f.Value.Type, err,
		)
	}

	return nil
}

func (bw *Writer) writeStruct(s wire.Struct) error {
	for _, f := range s.Fields {
		if err := bw.writeField(f); err != nil {
			return err
		}
	}
	return bw.writeByte(0) // end struct
}

func (bw *Writer) writeMap(m wire.Map) error {
	// ktype:1
	if err := bw.writeByte(byte(m.KeyType)); err != nil {
		return err
	}

	// vtype:1
	if err := bw.writeByte(byte(m.ValueType)); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(len(m.Items))); err != nil {
		return err
	}

	for _, item := range m.Items {
		if err := bw.WriteValue(item.Key); err != nil {
			return err
		}
		if err := bw.WriteValue(item.Value); err != nil {
			return err
		}
	}
	return nil
}

func (bw *Writer) writeSet(s wire.Set) error {
	// vtype:1
	if err := bw.writeByte(byte(s.ValueType)); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(len(s.Items))); err != nil {
		return err
	}

	for _, item := range s.Items {
		if err := bw.WriteValue(item); err != nil {
			return err
		}
	}
	return nil
}

func (bw *Writer) writeList(l wire.List) error {
	// vtype:1
	if err := bw.writeByte(byte(l.ValueType)); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(len(l.Items))); err != nil {
		return err
	}

	for _, item := range l.Items {
		if err := bw.WriteValue(item); err != nil {
			return err
		}
	}
	return nil
}

// WriteValue writes out the given Thrift value.
func (bw *Writer) WriteValue(v wire.Value) error {
	switch v.Type {
	case wire.TBool:
		if v.Bool {
			return bw.writeByte(1)
		}
		return bw.writeByte(0)

	case wire.TByte:
		return bw.writeByte(byte(v.Byte))

	case wire.TDouble:
		value := math.Float64bits(v.Double)
		return bw.writeInt64(int64(value))

	case wire.TI16:
		return bw.writeInt16(v.I16)

	case wire.TI32:
		return bw.writeInt32(v.I32)

	case wire.TI64:
		return bw.writeInt64(v.I64)

	case wire.TBinary:
		if err := bw.writeInt32(int32(len(v.Binary))); err != nil {
			return err
		}
		return bw.write(v.Binary)

	case wire.TStruct:
		return bw.writeStruct(v.Struct)

	case wire.TMap:
		return bw.writeMap(v.Map)

	case wire.TSet:
		return bw.writeSet(v.Set)

	case wire.TList:
		return bw.writeList(v.List)

	default:
		return fmt.Errorf("unknown ttype %v", v.Type)
	}
}
