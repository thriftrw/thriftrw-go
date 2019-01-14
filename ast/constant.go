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

// ConstantValue unifies the different types representing constant values in
// Thrift files.
type ConstantValue interface {
	Node

	constantValue()
}

func (ConstantBoolean) node()   {}
func (ConstantInteger) node()   {}
func (ConstantString) node()    {}
func (ConstantDouble) node()    {}
func (ConstantReference) node() {}
func (ConstantMap) node()       {}
func (ConstantList) node()      {}

func (ConstantBoolean) visitChildren(nodeStack, visitor)   {}
func (ConstantInteger) visitChildren(nodeStack, visitor)   {}
func (ConstantString) visitChildren(nodeStack, visitor)    {}
func (ConstantDouble) visitChildren(nodeStack, visitor)    {}
func (ConstantReference) visitChildren(nodeStack, visitor) {}

func (ConstantBoolean) constantValue()   {}
func (ConstantInteger) constantValue()   {}
func (ConstantString) constantValue()    {}
func (ConstantDouble) constantValue()    {}
func (ConstantReference) constantValue() {}
func (ConstantMap) constantValue()       {}
func (ConstantList) constantValue()      {}

func (l ConstantList) visitChildren(ss nodeStack, v visitor) {
	for _, item := range l.Items {
		v.visit(ss, item)
	}
}

func (m ConstantMap) visitChildren(ss nodeStack, v visitor) {
	for _, item := range m.Items {
		v.visit(ss, item)
	}
}

func (i ConstantMapItem) visitChildren(ss nodeStack, v visitor) {
	v.visit(ss, i.Key)
	v.visit(ss, i.Value)
}

func (m ConstantMap) lineNumber() int       { return m.Line }
func (i ConstantMapItem) lineNumber() int   { return i.Line }
func (l ConstantList) lineNumber() int      { return l.Line }
func (r ConstantReference) lineNumber() int { return r.Line }

// ConstantBoolean is a boolean value specified in the Thrift file.
//
//   true
//   false
type ConstantBoolean bool

// ConstantInteger is an integer value specified in the Thrift file.
//
//   42
type ConstantInteger int64

// ConstantString is a string literal specified in the Thrift file.
//
//   "hello world"
type ConstantString string

// ConstantDouble is a floating point value specified in the Thrift file.
//
//   1.234
type ConstantDouble float64

// ConstantMap is a map literal from the Thrift file.
//
// 	{"a": 1, "b": 2}
//
// Note that map literals can also be used to build structs.
type ConstantMap struct {
	Items []ConstantMapItem
	Line  int
}

// ConstantMapItem is a single item in a ConstantMap.
type ConstantMapItem struct {
	Key, Value ConstantValue
	Line       int
}

func (ConstantMapItem) node() {}

// ConstantList is a list literal from the Thrift file.
//
// 	[1, 2, 3]
type ConstantList struct {
	Items []ConstantValue
	Line  int
}

// ConstantReference is a reference to another constant value defined in the
// Thrift file.
//
// 	foo.bar
type ConstantReference struct {
	// Name of the referenced value.
	Name string

	// Line number on which this reference was made.
	Line int
}
