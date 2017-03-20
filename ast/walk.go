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

// Walk walks the AST depth-first with the given visitor, starting at the
// given node. The visitor's Visit function should return a non-nil visitor if
// it wants to visit the children of the node it was called with.
func Walk(v Visitor, n Node) {
	w := visitorWrapper{v: v}
	w.visit(nil, n)
}

// nodeStack of nodes visited in the order they were visited
type nodeStack []Node

func (ss nodeStack) Parent() Node {
	if len(ss) == 0 {
		return nil
	}
	return ss[len(ss)-1]
}

func (ss nodeStack) Ancestors() []Node {
	if len(ss) == 0 {
		return nil
	}

	ancestors := make([]Node, len(ss))
	for i, n := range ss {
		ancestors[len(ss)-1-i] = n
	}
	return ancestors
}

// Adapts a user-provided Visitor into an internal visitor so that we can call
// the visitChildren method on nodes.
type visitorWrapper struct{ v Visitor }

func (w visitorWrapper) visit(ss nodeStack, n Node) {
	if n == nil {
		return
	}

	// Note that visitorWrapper is passed by value so we're operating on a
	// copy of visitorWrapper.
	w.v = w.v.Visit(ss, n)
	if w.v == nil {
		return
	}

	ss = append(ss, n)
	n.visitChildren(ss, w)
}
