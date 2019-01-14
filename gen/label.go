// Copyright (c) 2019 Uber Technologies, Inc.
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

package gen

import "go.uber.org/thriftrw/compile"

// goLabelKey is a Thrift annotation that allows overriding the text
// formatting of an enum item or a struct field.
//
// Specifically, enum items will use the value specified in the goLabelKey
// annotation in place of the enum item name when they are serialized to/from
// user-readable strings.
//
// Similarly, struct fields will use the value specified in goLabelKey instead
// of the field name if the struct is serialized to JSON, logged to Zap, etc.
const goLabelKey = "go.label"

// entityLabel returns the string label that should be used for the given
// entity.  This is the human-readable name of the entity that *does not*
// affect the names of the generated Go-level entities, only their string
// representations.
func entityLabel(e compile.NamedEntity) string {
	anns := e.ThriftAnnotations()
	if val := anns[goLabelKey]; len(val) > 0 {
		return val
	}
	return e.ThriftName()
}
