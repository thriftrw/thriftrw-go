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

package wire

// FieldList is a list of fields of a struct.
type FieldList interface {
	// ForEach calls the provided function on all serializable fields in a
	// struct.
	//
	// If the provided function returns a non-nil error, ForEach halts and
	// returns the error.
	ForEach(func(Field) error) error

	// Close releases any resources associated with the FieldList. The
	// FieldList MUST NOT be used again after this invocation.
	Close()
}

// FieldListFromSlice builds a FieldList from the provided slice.
func FieldListFromSlice(fields []Field) FieldList {
	return fieldList(fields)
}

type fieldList []Field

func (fl fieldList) ForEach(f func(Field) error) error {
	for _, field := range fl {
		if err := f(field); err != nil {
			return err
		}
	}
	return nil
}

func (fl fieldList) Close() {}

// FieldListToSlice converts a FieldList to a list of fields.
func FieldListToSlice(fl FieldList) []Field {
	if fields, ok := fl.(fieldList); ok {
		// Optimize FieldListToSlice(FieldListFromSlice(...))
		return fields
	}

	var fields []Field
	// ignore the error; this cannot fail
	_ = fl.ForEach(func(f Field) error {
		fields = append(fields, f)
		return nil
	})
	return fields
}
