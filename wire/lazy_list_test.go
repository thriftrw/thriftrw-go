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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValueListFromSliceAll(t *testing.T) {
	slice := []Value{
		Value{Type: TI32, I32: 1},
		Value{Type: TI32, I32: 2},
		Value{Type: TI32, I32: 3},
	}

	i := 0
	ValueListFromSlice(slice).ForEach(func(v Value) error {
		assert.Equal(t, slice[i], v)
		i++

		return nil
	})

	assert.Equal(t, 3, i)
}

func TestValueListFromSliceBreak(t *testing.T) {
	slice := []Value{
		Value{Type: TI32, I32: 1},
		Value{Type: TI32, I32: 2},
		Value{Type: TI32, I32: 3},
	}

	i := 0
	ValueListFromSlice(slice).ForEach(func(v Value) error {
		assert.Equal(t, slice[i], v)
		i++
		if i == 2 {
			// Break after processing I32: 2
			return fmt.Errorf("fail")
		} else {
			return nil
		}
	})

	assert.Equal(t, 2, i)
}

//////////////////////////////////////////////////////////////////////////////

func TestMapItemListFromSliceAll(t *testing.T) {
	slice := []MapItem{
		MapItem{
			Key:   Value{Type: TI32, I32: 1},
			Value: Value{Type: TI64, I64: 101},
		},
		MapItem{
			Key:   Value{Type: TI32, I32: 2},
			Value: Value{Type: TI64, I64: 102},
		},
		MapItem{
			Key:   Value{Type: TI32, I32: 3},
			Value: Value{Type: TI64, I64: 103},
		},
	}

	i := 0
	MapItemListFromSlice(slice).ForEach(func(v MapItem) error {
		assert.Equal(t, slice[i], v)
		i++

		return nil
	})

	assert.Equal(t, 3, i)
}

func TestMapItemListFromSliceBreak(t *testing.T) {
	slice := []MapItem{
		MapItem{
			Key:   Value{Type: TI32, I32: 1},
			Value: Value{Type: TI64, I64: 101},
		},
		MapItem{
			Key:   Value{Type: TI32, I32: 2},
			Value: Value{Type: TI64, I64: 102},
		},
		MapItem{
			Key:   Value{Type: TI32, I32: 3},
			Value: Value{Type: TI64, I64: 103},
		},
	}

	i := 0
	MapItemListFromSlice(slice).ForEach(func(v MapItem) error {
		assert.Equal(t, slice[i], v)
		i++
		if i == 2 {
			// Break after processing I32: 2
			return fmt.Errorf("fail")
		} else {
			return nil
		}
	})

	assert.Equal(t, 2, i)
}
