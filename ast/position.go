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

package ast

import "strconv"

// Position represents a position in the parsed document.
// Line and column numbers are 1-based.
type Position struct {
	Line   int
	Column int
}

func (p Position) String() string {
	s := strconv.Itoa(p.Line)
	if c := p.Column; c > 0 {
		s += ":" + strconv.Itoa(c)
	}
	return s
}

// Pos attempts to return the position of a Node in the parsed document.
// For most use cases, prefer to use idl.Info to access positional information.
func Pos(n Node) (Position, bool) {
	if np, ok := n.(nodeWithPosition); ok {
		return np.pos(), true
	}
	return Position{}, false
}

// LineNumber returns the line in the file at which the given node was defined
// or 0 if the Node does not record its line number.
func LineNumber(n Node) int {
	if np, ok := n.(nodeWithPosition); ok {
		return np.pos().Line
	}
	return 0
}

// Nodes which know their document position can implement this interface.
type nodeWithPosition interface {
	Node

	pos() Position
}

var _ nodeWithPosition = (*Annotation)(nil)
var _ nodeWithPosition = BaseType{}
var _ nodeWithPosition = (*Constant)(nil)
var _ nodeWithPosition = ConstantList{}
var _ nodeWithPosition = ConstantMap{}
var _ nodeWithPosition = ConstantMapItem{}
var _ nodeWithPosition = ConstantReference{}
var _ nodeWithPosition = (*Enum)(nil)
var _ nodeWithPosition = (*EnumItem)(nil)
var _ nodeWithPosition = (*Field)(nil)
var _ nodeWithPosition = (*Function)(nil)
var _ nodeWithPosition = (*Include)(nil)
var _ nodeWithPosition = (*CppInclude)(nil)
var _ nodeWithPosition = ListType{}
var _ nodeWithPosition = MapType{}
var _ nodeWithPosition = (*Namespace)(nil)
var _ nodeWithPosition = (*Service)(nil)
var _ nodeWithPosition = SetType{}
var _ nodeWithPosition = (*Struct)(nil)
var _ nodeWithPosition = (*Typedef)(nil)
var _ nodeWithPosition = TypeReference{}
