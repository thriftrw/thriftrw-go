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

package wire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func vi32(i int32) Value {
	return NewValueI32(i)
}

func vbinary(s string) Value {
	return NewValueBinary([]byte(s))
}

func vlist(typ Type, vs ...Value) Value {
	return NewValueList(ValueListFromSlice(typ, vs))
}

func vset(typ Type, vs ...Value) Value {
	return NewValueSet(ValueListFromSlice(typ, vs))
}

func vmap(kt, vt Type, items ...MapItem) Value {
	return NewValueMap(MapItemListFromSlice(kt, vt, items))
}

func vitem(k, v Value) MapItem {
	return MapItem{Key: k, Value: v}
}

func TestEquals(t *testing.T) {
	tests := []struct{ l, r Value }{
		{
			// set with hashable items
			vset(TBinary, vbinary("1"), vbinary("2"), vbinary("3")),
			vset(TBinary, vbinary("3"), vbinary("1"), vbinary("2")),
		},
		{
			// map with hashable keys
			vmap(
				TI32, TSet,
				vitem(vi32(1), vset(TI32, vi32(2), vi32(3))),
				vitem(vi32(4), vset(TI32, vi32(5), vi32(6))),
				vitem(vi32(7), vset(TI32, vi32(8), vi32(9))),
			),
			vmap(
				TI32, TSet,
				vitem(vi32(7), vset(TI32, vi32(9), vi32(8))),
				vitem(vi32(1), vset(TI32, vi32(3), vi32(2))),
				vitem(vi32(4), vset(TI32, vi32(6), vi32(5))),
			),
		},
		{
			// set with unhashable items
			vset(
				TList,
				vlist(TI32, vi32(1), vi32(2), vi32(3), vi32(4)),
				vlist(TI32, vi32(5), vi32(6), vi32(7), vi32(8)),
				vlist(TI32, vi32(9), vi32(0), vi32(1), vi32(2)),
			),
			vset(
				TList,
				vlist(TI32, vi32(9), vi32(0), vi32(1), vi32(2)),
				vlist(TI32, vi32(1), vi32(2), vi32(3), vi32(4)),
				vlist(TI32, vi32(5), vi32(6), vi32(7), vi32(8)),
			),
		},
		{
			// map with unhashable keys
			vmap(
				TSet, TI32,
				vitem(vset(TI32, vi32(1), vi32(2), vi32(3)), vi32(1)),
				vitem(vset(TI32, vi32(4), vi32(5), vi32(6)), vi32(2)),
				vitem(vset(TI32, vi32(7), vi32(8), vi32(9)), vi32(3)),
			),
			vmap(
				TSet, TI32,
				vitem(vset(TI32, vi32(6), vi32(4), vi32(5)), vi32(2)),
				vitem(vset(TI32, vi32(8), vi32(9), vi32(7)), vi32(3)),
				vitem(vset(TI32, vi32(3), vi32(2), vi32(1)), vi32(1)),
			),
		},
	}

	for _, tt := range tests {
		assert.True(
			t, ValuesAreEqual(tt.l, tt.r),
			"Values should be equal:\n\t   %v\n\t!= %v", tt.l, tt.r,
		)
	}
}
