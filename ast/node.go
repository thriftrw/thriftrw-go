// Copyright (c) 2017 Uber Technologies, Inc.
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

// Node is a single element in the Thrift AST.
//
// In addition to all Header, ConstantValue, Type, and Definition types, the
// following types are also AST nodes: *Annotation, ConstantMapItem,
// *EnumItem, *Field, *Function, *Program.
type Node interface {
	node()

	// Must call the given function on each child node.
	forEachChild(func(Node))
}

var _ Node = (*Annotation)(nil)
var _ Node = BaseType{}
var _ Node = (*Constant)(nil)
var _ Node = ConstantBoolean(true)
var _ Node = ConstantDouble(1.0)
var _ Node = ConstantInteger(1)
var _ Node = ConstantList{}
var _ Node = ConstantMap{}
var _ Node = ConstantMapItem{}
var _ Node = ConstantReference{}
var _ Node = ConstantString("hi")
var _ Node = (*Enum)(nil)
var _ Node = (*EnumItem)(nil)
var _ Node = (*Field)(nil)
var _ Node = (*Function)(nil)
var _ Node = (*Include)(nil)
var _ Node = ListType{}
var _ Node = MapType{}
var _ Node = (*Namespace)(nil)
var _ Node = (*Program)(nil)
var _ Node = (*Service)(nil)
var _ Node = SetType{}
var _ Node = (*Struct)(nil)
var _ Node = TypeReference{}
var _ Node = (*Typedef)(nil)