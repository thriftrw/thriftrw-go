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

// +build gofuzz

package binary

import (
	"bytes"
	"fmt"

	"github.com/uber/thriftrw-go/wire"
)

func Fuzz(data []byte) int {
	reader := NewReader(bytes.NewReader(data))
	value, pos, err := reader.ReadValue(wire.TStruct, 0)
	if err != nil || pos != int64(len(data)) {
		return 0
	}
	if err := evaluate(value); err != nil {
		return 0
	}

	buffer := bytes.Buffer{}
	writer := BorrowWriter(&buffer)
	if err := writer.WriteValue(value); err != nil {
		panic(fmt.Sprintf("error encoding %s: %s", value, err))
	}
	ReturnWriter(writer)

	encoded := buffer.Bytes()
	if !bytesSame(data, encoded) {
		panic(fmt.Sprintf(
			"encoding mismatch for %s:\n\t   %#v (got)\n\t!= %#v (expected)\n",
			value, encoded, data,
		))
	}

	return 1
}

func bytesSame(ls, rs []byte) bool {
	if len(ls) != len(rs) {
		return false
	}
	for i := 0; i < len(ls); i++ {
		if ls[i] != rs[i] {
			return false
		}
	}
	return true
}

// fully evaluate a value that contains lazy lists, etc.
func evaluate(v wire.Value) error {
	switch v.Type {
	case wire.TBool:
		return nil
	case wire.TByte:
		return nil
	case wire.TDouble:
		return nil
	case wire.TI16:
		return nil
	case wire.TI32:
		return nil
	case wire.TI64:
		return nil
	case wire.TBinary:
		return nil
	case wire.TStruct:
		for _, f := range v.Struct.Fields {
			if err := evaluate(f.Value); err != nil {
				return err
			}
		}
		return nil
	case wire.TMap:
		return v.Map.Items.ForEach(func(item wire.MapItem) error {
			if err := evaluate(item.Key); err != nil {
				return err
			}
			if err := evaluate(item.Value); err != nil {
				return err
			}
			return nil
		})
	case wire.TSet:
		return v.Set.Items.ForEach(evaluate)
	case wire.TList:
		return v.List.Items.ForEach(evaluate)
	default:
		return fmt.Errorf("unknown type %s", v.Type)
	}
}
