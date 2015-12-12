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

	// Close indicates that the caller is finished reading from the lazy list.
	Close()
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

	// Close indicates that the caller is finished reading from the lazy list.
	Close()
}

//////////////////////////////////////////////////////////////////////////////

// ValueListFromSlice builds a ValueList from the given slice of Values.
type ValueListFromSlice []Value

func (vs ValueListFromSlice) ForEach(f func(Value) error) error {
	for _, v := range vs {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

func (ValueListFromSlice) Close() {}

//////////////////////////////////////////////////////////////////////////////

// MapItemListFromSlice builds a MapItemList from the given slice of Values.
type MapItemListFromSlice []MapItem

func (vs MapItemListFromSlice) ForEach(f func(MapItem) error) error {
	for _, v := range vs {
		if err := f(v); err != nil {
			return err
		}
	}
	return nil
}

func (MapItemListFromSlice) Close() {}

//////////////////////////////////////////////////////////////////////////////

// ValueListToSlice builds a slice of values from the given ValueList.
//
// Capacity may be provided to set an initial capacity for the slice.
func ValueListToSlice(l ValueList, capacity int) []Value {
	if capacity < 0 {
		capacity = 0
	}

	items := make([]Value, 0, capacity)
	l.ForEach(func(v Value) error {
		items = append(items, v)
		return nil
	})
	return items
}

// MapItemListToSlice builds a slice of values from the given MapItemList.
//
// Capacity may be provided to set an initial capacity for the slice.
func MapItemListToSlice(l MapItemList, capacity int) []MapItem {
	if capacity < 0 {
		capacity = 0
	}

	items := make([]MapItem, 0, capacity)
	l.ForEach(func(v MapItem) error {
		items = append(items, v)
		return nil
	})
	return items
}
