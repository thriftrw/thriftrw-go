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
	"sync"

	"go.uber.org/thriftrw/wire"
)

var (
	lazyFieldListPool = sync.Pool{New: func() interface{} {
		return new(lazyFieldList)
	}}
)

type lazyFieldList struct {
	reader *Reader
	offset int64
}

func borrowLazyFieldList(r *Reader) *lazyFieldList {
	return lazyFieldListPool.Get().(*lazyFieldList)
}

func (ll *lazyFieldList) Close() {
	ll.reader = nil
	lazyFieldListPool.Put(ll)
}

func (ll *lazyFieldList) ForEach(f func(wire.Field) error) error {
	off := ll.offset

	br := ll.reader
	typ, off, err := br.readByte(off)
	if err != nil {
		return err
	}

	for typ != 0 {
		var fid int16
		var val wire.Value

		fid, off, err = br.readInt16(off)
		if err != nil {
			return err
		}

		val, off, err = br.ReadValue(wire.Type(typ), off)
		if err != nil {
			return err
		}

		if err := f(wire.Field{ID: fid, Value: val}); err != nil {
			return err
		}

		typ, off, err = br.readByte(off)
		if err != nil {
			return err
		}
	}

	return nil
}
