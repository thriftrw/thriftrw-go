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

// ValueList represents a collection of Value objects as an iteration through
// it. This helps us avoid the cost of allocating memory for all collections
// passing through the system.
type ValueList interface {
	// ForEach calls the given function on each element of the list.
	//
	// If any call fails with an error, that error is returned and the
	// iteration is stopped.
	ForEach(f func(Value) error) error
}

// MapItemList represents a collection of MapItem objects as an iteration
// through it. This helps us avoid the cost of allocating memory for all
// collections passing through the system.
type MapItemList interface {
	// ForEach calls the given function on each element of the list.
	//
	// If any call fails with an error, that error is returned and the
	// iteration is stopped.
	ForEach(f func(MapItem) error) error
}

//////////////////////////////////////////////////////////////////////////////

type valueListFromSlice struct {
	values []Value
}

// ValueListFromSlice builds a ValueList from the given slice of Values.
func ValueListFromSlice(ls []Value) ValueList {
	return valueListFromSlice{ls}
}

func (v valueListFromSlice) ForEach(f func(Value) error) error {
	for _, v := range v.values {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

//////////////////////////////////////////////////////////////////////////////

type mapItemListFromSlice struct {
	items []MapItem
}

// MapItemListFromSlice builds a MapItemList from the given slice of Values.
func MapItemListFromSlice(ls []MapItem) MapItemList {
	return mapItemListFromSlice{ls}
}

func (v mapItemListFromSlice) ForEach(f func(MapItem) error) error {
	for _, v := range v.items {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}
