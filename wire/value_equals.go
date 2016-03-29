// Copyright (c) 2016 Uber Technologies, Inc.
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
	"bytes"
	"errors"
	"fmt"
)

// notEqualError is a sentinel error used while iterating through ValueLists
// to indicate that two values did not match.
var notEqualError = errors.New("Values are not equal")

// ValuesAreEqual checks if two values are equal.
func ValuesAreEqual(left, right Value) bool {
	if left.typ != right.typ {
		return false
	}

	switch left.typ {
	case TBool:
		return left.tbool == right.tbool
	case TI8:
		return left.ti8 == right.ti8
	case TDouble:
		return left.tdouble == right.tdouble
	case TI16:
		return left.ti16 == right.ti16
	case TI32:
		return left.ti32 == right.ti32
	case TI64:
		return left.ti64 == right.ti64
	case TBinary:
		return bytes.Equal(left.tbinary, right.tbinary)
	case TStruct:
		return StructsAreEqual(left.tstruct, right.tstruct)
	case TMap:
		return MapsAreEqual(left.tmap, right.tmap)
	case TSet:
		return SetsAreEqual(left.tset, right.tset)
	case TList:
		return ListsAreEqual(left.tlist, right.tlist)
	default:
		return false
	}
}

// StructsAreEqual checks if two structs are equal.
func StructsAreEqual(left, right Struct) bool {
	if len(left.Fields) != len(right.Fields) {
		return false
	}

	// Fields are unordered so we need to build a map to actually compare
	// them.

	leftFields := left.fieldMap()
	rightFields := right.fieldMap()

	for i, lvalue := range leftFields {
		if rvalue, ok := rightFields[i]; !ok {
			return false
		} else if !ValuesAreEqual(lvalue, rvalue) {
			return false
		}
	}

	return true
}

// SetsAreEqual checks if two sets are equal.
func SetsAreEqual(left, right Set) bool {
	if left.ValueType != right.ValueType {
		return false
	}
	if left.Size != right.Size {
		return false
	}

	if isHashable(left.ValueType) {
		return setsArEqualHashable(left.Size, left.Items, right.Items)
	} else {
		return setsAreEqualUnhashable(left.Size, left.Items, right.Items)
	}
}

// setsArEqualHashable checks if two unordered ValueLists are equal, provided
// that they contain items that are hashable -- that is, the items can be used
// as keys in a map.
func setsArEqualHashable(size int, l, r ValueList) bool {
	m := make(map[interface{}]bool, size)
	l.ForEach(func(v Value) error {
		m[toHashable(v)] = true
		return nil
	})

	return notEqualError != r.ForEach(func(v Value) error {
		if _, ok := m[toHashable(v)]; !ok {
			return notEqualError
		}
		return nil
	})
}

// setsAreEqualUnhashable checks if two unordered ValueLists are equal for
// types that are not hashable. Note that this is O(n^2) in time complexity.
func setsAreEqualUnhashable(size int, l, r ValueList) bool {
	lItems := ValueListToSlice(l, size)

	return notEqualError != r.ForEach(func(rItem Value) error {
		matched := false
		for _, lItem := range lItems {
			if ValuesAreEqual(lItem, rItem) {
				matched = true
				break
			}
		}
		if !matched {
			return notEqualError
		}
		return nil
	})
}

// MapsAreEqual checks if two maps are equal.
func MapsAreEqual(left, right Map) bool {
	if left.KeyType != right.KeyType {
		return false
	}
	if left.ValueType != right.ValueType {
		return false
	}
	if left.Size != right.Size {
		return false
	}

	if isHashable(left.KeyType) {
		return mapsAreEqualHashable(left.Size, left.Items, right.Items)
	} else {
		return mapsAreEqualUnhashable(left.Size, left.Items, right.Items)
	}
}

func mapsAreEqualHashable(size int, l, r MapItemList) bool {
	m := make(map[interface{}]Value, size)

	l.ForEach(func(item MapItem) error {
		m[toHashable(item.Key)] = item.Value
		return nil
	})

	return notEqualError != r.ForEach(func(item MapItem) error {
		if lValue, ok := m[toHashable(item.Key)]; !ok {
			return notEqualError
		} else {
			if !ValuesAreEqual(lValue, item.Value) {
				return notEqualError
			}
			return nil
		}
	})
}

func mapsAreEqualUnhashable(size int, l, r MapItemList) bool {
	lItems := MapItemListToSlice(l, size)

	return notEqualError != r.ForEach(func(rItem MapItem) error {
		matched := false
		for _, lItem := range lItems {
			if !ValuesAreEqual(lItem.Key, rItem.Key) {
				continue
			}
			if !ValuesAreEqual(lItem.Value, rItem.Value) {
				continue
			}
			matched = true
		}

		if !matched {
			return notEqualError
		}
		return nil
	})
}

func isHashable(t Type) bool {
	switch t {
	case TBool, TI8, TDouble, TI16, TI32, TI64, TBinary:
		return true
	default:
		return false
	}
}

func toHashable(v Value) interface{} {
	switch v.Type() {
	case TBool, TI8, TDouble, TI16, TI32, TI64:
		return v.Get()
	case TBinary:
		return string(v.GetBinary())
	default:
		panic(fmt.Sprintf("value is not hashable: %v", v))
	}
}

// ListsAreEqual checks if two lists are equal.
func ListsAreEqual(left, right List) bool {
	if left.ValueType != right.ValueType {
		return false
	}
	if left.Size != right.Size {
		return false
	}

	leftItems := ValueListToSlice(left.Items, left.Size)
	rightItems := ValueListToSlice(right.Items, right.Size)

	for i, lv := range leftItems {
		rv := rightItems[i]
		if !ValuesAreEqual(lv, rv) {
			return false
		}
	}

	return true
}
