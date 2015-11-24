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
	"strings"
)

// Value holds the over-the-wire representation of a Thrift value.
//
// The Type of the value determines which field in the Value is valid.
type Value struct {
	Type Type

	Bool   bool
	Byte   int8
	Double float64
	I16    int16
	I32    int32
	I64    int64
	Binary []byte
	Struct Struct
	Map    Map
	Set    Set
	List   List
}

func (v Value) String() string {
	switch v.Type {
	case TBool:
		return fmt.Sprint(v.Bool)
	case TByte:
		return fmt.Sprint(v.Byte)
	case TDouble:
		return fmt.Sprint(v.Double)
	case TI16:
		return fmt.Sprint(v.I16)
	case TI32:
		return fmt.Sprint(v.I32)
	case TI64:
		return fmt.Sprint(v.I64)
	case TBinary:
		return fmt.Sprint(v.Binary)
	case TStruct:
		return v.Struct.String()
	case TMap:
		return v.Map.String()
	case TSet:
		return v.Set.String()
	case TList:
		return v.List.String()
	default:
		return fmt.Sprintf("%#v", v)
	}
}

// Struct provides a wire-level representation of a struct.
//
// At this level, structs don't have names or named fields.
type Struct struct {
	Fields []Field
}

func (s Struct) String() string {
	fields := make([]string, len(s.Fields))
	for i, field := range s.Fields {
		fields[i] = field.String()
	}
	return fmt.Sprintf("{%s}", strings.Join(fields, ", "))
}

// Field is a single field inside a Struct.
type Field struct {
	ID    int16
	Value Value
}

func (f Field) String() string {
	return fmt.Sprintf("%v: %v", f.ID, f.Value)
}

// Set is a set of values.
type Set struct {
	ValueType Type
	Size      int
	Items     ValueList
}

func (s Set) String() string {
	items := make([]string, 0, s.Size)
	s.Items.ForEach(func(item Value) error {
		items = append(items, item.String())
		return nil
	})

	return fmt.Sprintf("[set]%v{%s}", s.ValueType, strings.Join(items, ", "))
}

// List is a list of values.
type List struct {
	ValueType Type
	Size      int
	Items     ValueList
}

func (l List) String() string {
	items := make([]string, 0, l.Size)
	l.Items.ForEach(func(item Value) error {
		items = append(items, item.String())
		return nil
	})

	return fmt.Sprintf("[]%v{%s}", l.ValueType, strings.Join(items, ", "))
}

// Map is a collection of key-value pairs.
type Map struct {
	KeyType   Type
	ValueType Type
	Size      int
	Items     MapItemList
}

func (m Map) String() string {
	items := make([]string, 0, m.Size)
	m.Items.ForEach(func(item MapItem) error {
		items = append(items, item.String())
		return nil
	})

	return fmt.Sprintf(
		"map[%v]%v{%s}", m.KeyType, m.ValueType, strings.Join(items, ", "),
	)
}

// MapItem is a single item in a Map.
type MapItem struct {
	Key   Value
	Value Value
}

func (mi MapItem) String() string {
	return fmt.Sprintf("%v: %v", mi.Key, mi.Value)
}
