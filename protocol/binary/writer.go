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
	"fmt"
	"io"
	"math"
	"sync"

	"go.uber.org/thriftrw/wire"
)

var writerPool = sync.Pool{New: func() interface{} {
	writer := &Writer{}
	writer.writeValue = writer.WriteValue
	writer.writeMapItem = writer.realWriteMapItem
	writer.writeField = writer.realWriteField
	return writer
}}

// Writer implements basic logic for writing the Thrift Binary Protocol to an
// io.Writer.
type Writer struct {
	writer io.Writer

	// This buffer is re-used every time we need a slice of up to 8 bytes.
	buffer [8]byte

	// NOTE:
	// This is a hack to avoid memory allocation in closures. Passing the
	// bound WriteValue or realWriteMapItem methods into a function results in
	// a memory allocation because the system doesn't know we're going to
	// reuse the closure. So we create that bound reference in advance when
	// the writer is created.
	writeValue   func(wire.Value) error
	writeField   func(wire.Field) error
	writeMapItem func(wire.MapItem) error
}

// BorrowWriter fetches a Writer from the system that will write its output to
// the given io.Writer.
//
// This Writer must be returned back using ReturnWriter.
func BorrowWriter(w io.Writer) *Writer {
	writer := writerPool.Get().(*Writer)
	writer.writer = w
	return writer
}

// ReturnWriter returns a previously borrowed Writer back to the system.
func ReturnWriter(w *Writer) {
	w.writer = nil
	writerPool.Put(w)
}

func (bw *Writer) write(bs []byte) error {
	_, err := bw.writer.Write(bs)
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

func (bw *Writer) writeString(s string) error {
	if err := bw.writeInt32(int32(len(s))); err != nil {
		return err
	}

	_, err := io.WriteString(bw.writer, s)
	return err
}

func (bw *Writer) realWriteField(f wire.Field) error {
	// type:1
	if err := bw.writeByte(byte(f.Value.Type())); err != nil {
		return err
	}

	// id:2
	if err := bw.writeInt16(f.ID); err != nil {
		return err
	}

	// value
	if err := bw.WriteValue(f.Value); err != nil {
		// TODO(abg): Figure out better error handling story. We need access
		// to the underlying error object if it's a network error.
		return fmt.Errorf(
			"failed to write field %d (%v): %s",
			f.ID, f.Value.Type(), err,
		)
	}

	return nil
}

func (bw *Writer) writeFieldList(fs wire.FieldList) error {
	if err := fs.ForEach(bw.writeField); err != nil {
		return err
	}
	// TODO(abg): It looks like write* functions for lazy collections do
	// not call Close. This has not been an issue during serialization
	// because the struct-based wrapper implementations don't have any
	// resources to free up, but we should still call close.
	return bw.writeByte(0) // end struct
}

func (bw *Writer) realWriteMapItem(item wire.MapItem) error {
	if err := bw.WriteValue(item.Key); err != nil {
		return err
	}
	return bw.WriteValue(item.Value)
}

func (bw *Writer) writeMap(m wire.MapItemList) error {
	// ktype:1
	if err := bw.writeByte(byte(m.KeyType())); err != nil {
		return err
	}

	// vtype:1
	if err := bw.writeByte(byte(m.ValueType())); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(m.Size())); err != nil {
		return err
	}

	return m.ForEach(bw.writeMapItem)
}

func (bw *Writer) writeSet(s wire.ValueList) error {
	// vtype:1
	if err := bw.writeByte(byte(s.ValueType())); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(s.Size())); err != nil {
		return err
	}

	return s.ForEach(bw.writeValue)
}

func (bw *Writer) writeList(l wire.ValueList) error {
	// vtype:1
	if err := bw.writeByte(byte(l.ValueType())); err != nil {
		return err
	}

	// length:4
	if err := bw.writeInt32(int32(l.Size())); err != nil {
		return err
	}

	return l.ForEach(bw.writeValue)
}

// WriteValue writes the given Thrift value to the underlying stream using the
// Thrift Binary Protocol.
func (bw *Writer) WriteValue(v wire.Value) error {
	switch v.Type() {
	case wire.TBool:
		if v.GetBool() {
			return bw.writeByte(1)
		}
		return bw.writeByte(0)

	case wire.TI8:
		return bw.writeByte(byte(v.GetI8()))

	case wire.TDouble:
		value := math.Float64bits(v.GetDouble())
		return bw.writeInt64(int64(value))

	case wire.TI16:
		return bw.writeInt16(v.GetI16())

	case wire.TI32:
		return bw.writeInt32(v.GetI32())

	case wire.TI64:
		return bw.writeInt64(v.GetI64())

	case wire.TBinary:
		b := v.GetBinary()
		if err := bw.writeInt32(int32(len(b))); err != nil {
			return err
		}
		return bw.write(b)

	case wire.TStruct:
		return bw.writeFieldList(v.GetFieldList())

	case wire.TMap:
		return bw.writeMap(v.GetMap())

	case wire.TSet:
		return bw.writeSet(v.GetSet())

	case wire.TList:
		return bw.writeList(v.GetList())

	default:
		return fmt.Errorf("unknown ttype %v", v.Type())
	}
}
