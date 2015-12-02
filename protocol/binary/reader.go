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

import (
	"fmt"
	"io"
	"math"

	"github.com/uber/thriftrw-go/wire"
)

type Reader struct {
	reader io.ReaderAt

	// This buffer is re-used every time we need a slice of up to 8 bytes.
	buffer [8]byte
}

func NewReader(r io.ReaderAt) Reader {
	return Reader{reader: r}
}

// For the reader, we keep track of the read offset manually everywhere so
// that we can implement lazy collections without extra allocations

func (br *Reader) skipStruct(off int64) (int64, error) {
	typ, off, err := br.readByte(off)
	if err != nil {
		return off, err
	}

	for typ != 0 {
		off += 2
		off, err = br.skipValue(wire.Type(typ), off)
		if err != nil {
			return off, err
		}

		typ, off, err = br.readByte(off)
		if err != nil {
			return off, err
		}
	}
	return off, err
}

func (br *Reader) skipMap(off int64) (int64, error) {
	kt_, off, err := br.readByte(off)
	if err != nil {
		return off, err
	}

	vt_, off, err := br.readByte(off)
	if err != nil {
		return off, err
	}

	kt := wire.Type(kt_)
	vt := wire.Type(vt_)

	count, off, err := br.readInt32(off)
	if err != nil {
		return off, err
	}

	for i := int32(0); i < count; i++ {
		off, err = br.skipValue(kt, off)
		if err != nil {
			return off, err
		}

		off, err = br.skipValue(vt, off)
		if err != nil {
			return off, err
		}
	}
	return off, err
}

func (br *Reader) skipList(off int64) (int64, error) {
	vt_, off, err := br.readByte(off)
	if err != nil {
		return off, err
	}
	vt := wire.Type(vt_)

	count, off, err := br.readInt32(off)
	if err != nil {
		return off, err
	}

	for i := int32(0); i < count; i++ {
		off, err = br.skipValue(vt, off)
		if err != nil {
			return off, err
		}
	}
	return off, err
}

func (br *Reader) skipValue(t wire.Type, off int64) (int64, error) {
	switch t {
	case wire.TBool:
		return off + 1, nil
	case wire.TByte:
		return off + 1, nil
	case wire.TDouble:
		return off + 8, nil
	case wire.TI16:
		return off + 2, nil
	case wire.TI32:
		return off + 4, nil
	case wire.TI64:
		return off + 8, nil
	case wire.TBinary:
		if length, off, err := br.readInt32(off); err != nil {
			return off, err
		} else {
			off += int64(length)
			return off, err
		}
	case wire.TStruct:
		return br.skipStruct(off)
	case wire.TMap:
		return br.skipMap(off)
	case wire.TSet:
		return br.skipList(off)
	case wire.TList:
		return br.skipList(off)
	default:
		return off, fmt.Errorf("unknown ttype %v", t)
	}
}

func (br *Reader) read(bs []byte, off int64) (int64, error) {
	n, err := br.reader.ReadAt(bs, off)
	off += int64(n)
	return off, err
}

func (br *Reader) readByte(off int64) (byte, int64, error) {
	bs := br.buffer[0:1]
	off, err := br.read(bs, off)
	return bs[0], off, err
}

func (br *Reader) readInt16(off int64) (int16, int64, error) {
	bs := br.buffer[0:2]
	off, err := br.read(bs, off)
	return int16(bigEndian.Uint16(bs)), off, err
}

func (br *Reader) readInt32(off int64) (int32, int64, error) {
	bs := br.buffer[0:4]
	off, err := br.read(bs, off)
	return int32(bigEndian.Uint32(bs)), off, err
}

func (br *Reader) readInt64(off int64) (int64, int64, error) {
	bs := br.buffer[0:8]
	off, err := br.read(bs, off)
	return int64(bigEndian.Uint64(bs)), off, err
}

func (br *Reader) readStruct(off int64) (wire.Struct, int64, error) {
	var fields []wire.Field
	// TODO lazy FieldList

	typ, off, err := br.readByte(off)
	if err != nil {
		return wire.Struct{}, off, err
	}

	for typ != 0 {
		var fid int16
		var val wire.Value

		fid, off, err = br.readInt16(off)
		if err != nil {
			return wire.Struct{}, off, err
		}

		val, off, err = br.ReadValue(wire.Type(typ), off)
		if err != nil {
			return wire.Struct{}, off, err
		}

		fields = append(fields, wire.Field{ID: fid, Value: val})

		typ, off, err = br.readByte(off)
		if err != nil {
			return wire.Struct{}, off, err
		}
	}
	return wire.Struct{Fields: fields}, off, err
}

func (br *Reader) readMap(off int64) (wire.Map, int64, error) {
	kt_, off, err := br.readByte(off)
	if err != nil {
		return wire.Map{}, off, err
	}

	vt_, off, err := br.readByte(off)
	if err != nil {
		return wire.Map{}, off, err
	}

	count, off, err := br.readInt32(off)
	if err != nil {
		return wire.Map{}, off, err
	}

	kt := wire.Type(kt_)
	vt := wire.Type(vt_)

	start := off
	for i := int32(0); i < count; i++ {
		off, err = br.skipValue(kt, off)
		if err != nil {
			return wire.Map{}, off, err
		}

		off, err = br.skipValue(vt, off)
		if err != nil {
			return wire.Map{}, off, err
		}
	}

	return wire.Map{
		KeyType:   kt,
		ValueType: vt,
		Size:      int(count),
		Items: lazyMapItemList{
			ktype:       kt,
			vtype:       vt,
			count:       count,
			reader:      br,
			startOffset: start,
		},
	}, off, err
}

func (br *Reader) readSet(off int64) (wire.Set, int64, error) {
	typ, off, err := br.readByte(off)
	if err != nil {
		return wire.Set{}, off, err
	}

	count, off, err := br.readInt32(off)
	if err != nil {
		return wire.Set{}, off, err
	}

	start := off
	for i := int32(0); i < count; i++ {
		off, err = br.skipValue(wire.Type(typ), off)
		if err != nil {
			return wire.Set{}, off, err
		}
	}

	return wire.Set{
		ValueType: wire.Type(typ),
		Size:      int(count),
		Items: lazyValueList{
			count:       count,
			typ:         wire.Type(typ),
			reader:      br,
			startOffset: start,
		},
	}, off, err
}

func (br *Reader) readList(off int64) (wire.List, int64, error) {
	typ, off, err := br.readByte(off)
	if err != nil {
		return wire.List{}, off, err
	}

	count, off, err := br.readInt32(off)
	if err != nil {
		return wire.List{}, off, err
	}

	start := off
	for i := int32(0); i < count; i++ {
		off, err = br.skipValue(wire.Type(typ), off)
		if err != nil {
			return wire.List{}, off, err
		}
	}

	return wire.List{
		ValueType: wire.Type(typ),
		Size:      int(count),
		Items: lazyValueList{
			count:       count,
			typ:         wire.Type(typ),
			reader:      br,
			startOffset: start,
		},
	}, off, err
}

func (br *Reader) ReadValue(t wire.Type, off int64) (wire.Value, int64, error) {
	switch t {
	case wire.TBool:
		b, off, err := br.readByte(off)
		if err != nil {
			return wire.Value{}, off, err
		}

		return wire.Value{Type: t, Bool: b == 1}, off, nil

	case wire.TByte:
		b, off, err := br.readByte(off)
		return wire.Value{Type: t, Byte: int8(b)}, off, err

	case wire.TDouble:
		value, off, err := br.readInt64(off)
		d := math.Float64frombits(uint64(value))
		return wire.Value{Type: t, Double: d}, off, err

	case wire.TI16:
		n, off, err := br.readInt16(off)
		return wire.Value{Type: t, I16: n}, off, err

	case wire.TI32:
		n, off, err := br.readInt32(off)
		return wire.Value{Type: t, I32: n}, off, err

	case wire.TI64:
		n, off, err := br.readInt64(off)
		return wire.Value{Type: t, I64: n}, off, err

	case wire.TBinary:
		length, off, err := br.readInt32(off)
		if err != nil {
			return wire.Value{}, off, err
		}

		bs := make([]byte, length)
		if length != 0 {
			off, err = br.read(bs, off)
		}
		return wire.Value{Type: t, Binary: bs}, off, err

	case wire.TStruct:
		s, off, err := br.readStruct(off)
		return wire.Value{Type: t, Struct: s}, off, err

	case wire.TMap:
		m, off, err := br.readMap(off)
		return wire.Value{Type: t, Map: m}, off, err

	case wire.TSet:
		s, off, err := br.readSet(off)
		return wire.Value{Type: t, Set: s}, off, err

	case wire.TList:
		l, off, err := br.readList(off)
		return wire.Value{Type: t, List: l}, off, err

	default:
		return wire.Value{}, off, fmt.Errorf("unknown ttype %v", t)
	}
}

type lazyValueList struct {
	count       int32
	typ         wire.Type
	reader      *Reader
	startOffset int64
}

func (ll lazyValueList) ForEach(f func(wire.Value) error) error {
	off := ll.startOffset

	var val wire.Value
	var err error
	for i := int32(0); i < ll.count; i++ {
		val, off, err = ll.reader.ReadValue(ll.typ, off)

		if err != nil {
			return err
		}

		if err := f(val); err != nil {
			return err
		}
	}
	return nil
}

type lazyMapItemList struct {
	ktype, vtype wire.Type
	count        int32
	reader       *Reader
	startOffset  int64
}

func (lm lazyMapItemList) ForEach(f func(wire.MapItem) error) error {
	off := lm.startOffset

	var k, v wire.Value
	var err error

	for i := int32(0); i < lm.count; i++ {
		k, off, err = lm.reader.ReadValue(lm.ktype, off)
		if err != nil {
			return err
		}

		v, off, err = lm.reader.ReadValue(lm.vtype, off)
		if err != nil {
			return err
		}

		item := wire.MapItem{Key: k, Value: v}
		if err := f(item); err != nil {
			return err
		}
	}
	return nil
}
