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

package protocol

import "go.uber.org/thriftrw/wire"

// This file doesn't actually contain any tests. It just contains helpers for
// constructing complex Value objects during protocol.

func vbool(b bool) wire.Value {
	return wire.NewValueBool(b)
}

func vi8(b int8) wire.Value {
	return wire.NewValueI8(b)
}

func vi16(i int16) wire.Value {
	return wire.NewValueI16(i)
}

func vi32(i int32) wire.Value {
	return wire.NewValueI32(i)
}

func vi64(i int64) wire.Value {
	return wire.NewValueI64(i)
}

func vdouble(f float64) wire.Value {
	return wire.NewValueDouble(f)
}

func vbinary(s string) wire.Value {
	return wire.NewValueBinary([]byte(s))
}

func vstruct(fs ...wire.Field) wire.Value {
	return wire.NewValueStruct(wire.Struct{Fields: fs})
}

func vfield(id int16, v wire.Value) wire.Field {
	return wire.Field{ID: id, Value: v}
}

func vlist(typ wire.Type, vs ...wire.Value) wire.Value {
	return wire.NewValueList(wire.ValueListFromSlice(typ, vs))
}

func vset(typ wire.Type, vs ...wire.Value) wire.Value {
	return wire.NewValueSet(wire.ValueListFromSlice(typ, vs))
}

func vmap(kt, vt wire.Type, items ...wire.MapItem) wire.Value {
	return wire.NewValueMap(wire.MapItemListFromSlice(kt, vt, items))
}

func vitem(k, v wire.Value) wire.MapItem {
	return wire.MapItem{Key: k, Value: v}
}
