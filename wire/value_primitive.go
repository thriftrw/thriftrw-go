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

// ToPrimitive converts a Value into a primitive representation of the data it
// contains.
//
// This is meant for debugging purposes only.
//
// Panics if keys in a map or items in a set aren't hashable.
func ToPrimitive(v Value) interface{} {
	return toPrimitive(v, false)
}

// Helper to determine if the value we want needs to be hashable.
//
// For most types this has no difference, but for binary, we'll return a
// string if the result needs to be hashable.
func toPrimitive(v Value, hashable bool) interface{} {
	switch v.Type {
	case TBool:
		return v.Bool
	case TByte:
		return v.Byte
	case TDouble:
		return v.Double
	case TI16:
		return v.I16
	case TI32:
		return v.I32
	case TI64:
		return v.I64
	case TBinary:
		if !hashable {
			return v.Binary
		}
		return string(v.Binary)
	case TStruct:
		s := make(map[int16]interface{})
		for _, f := range v.Struct.Fields {
			s[f.ID] = ToPrimitive(f.Value)
		}
		return s
	case TMap:
		m := make(map[interface{}]interface{})
		v.Map.Items.ForEach(func(item MapItem) error {
			m[toPrimitive(item.Key, true)] = ToPrimitive(item.Value)
			return nil
		})
		return m
	case TSet:
		s := make(map[interface{}]bool)
		v.Set.Items.ForEach(func(v Value) error {
			s[toPrimitive(v, true)] = true
			return nil
		})
		return s
	case TList:
		l := make([]interface{}, 0, v.List.Size)
		v.List.Items.ForEach(func(v Value) error {
			l = append(l, ToPrimitive(v))
			return nil
		})
		return l
	default:
		return v // unrecognized
	}
}
