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

package ast

// Nodes which know the line number they were defined on can implement this
// interface.
type nodeWithLine interface {
	Node

	lineNumber() int
}

// LineNumber returns the line in the file at which the given node was defined
// or 0 if the Node does not record its line number.
func LineNumber(n Node) int {
	if nl, ok := n.(nodeWithLine); ok {
		return nl.lineNumber()
	}
	return 0
}

var _ nodeWithLine = (*Annotation)(nil)
var _ nodeWithLine = BaseType{}
var _ nodeWithLine = (*Constant)(nil)
var _ nodeWithLine = ConstantList{}
var _ nodeWithLine = ConstantMap{}
var _ nodeWithLine = ConstantMapItem{}
var _ nodeWithLine = ConstantReference{}
var _ nodeWithLine = (*Enum)(nil)
var _ nodeWithLine = (*EnumItem)(nil)
var _ nodeWithLine = (*Field)(nil)
var _ nodeWithLine = (*Function)(nil)
var _ nodeWithLine = (*Include)(nil)
var _ nodeWithLine = ListType{}
var _ nodeWithLine = MapType{}
var _ nodeWithLine = (*Namespace)(nil)
var _ nodeWithLine = (*Service)(nil)
var _ nodeWithLine = SetType{}
var _ nodeWithLine = (*Struct)(nil)
var _ nodeWithLine = TypeReference{}
var _ nodeWithLine = (*Typedef)(nil)
