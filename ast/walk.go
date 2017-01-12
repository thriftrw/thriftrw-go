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
	walk(nil, v, n)
}

// nodeStack of nodes visited in the order they were visited
type nodeStack []Node

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

func walk(stack nodeStack, v Visitor, node Node) {
	if v = v.Visit(stack, node); v == nil {
		return
	}

	// After visiting the node, if we still have a non-nil visitor, we need to
	// visit its children recursively.

	stack = append(stack, node)
	switch n := node.(type) {
	case BaseType:
		walkAnnotations(stack, v, n.Annotations)
	case *Constant:
		walk(stack, v, n.Type)
		walk(stack, v, n.Value)
	case ConstantList:
		for _, item := range n.Items {
			walk(stack, v, item)
		}
	case ConstantMap:
		for _, item := range n.Items {
			walk(stack, v, item)
		}
	case ConstantMapItem:
		walk(stack, v, n.Key)
		walk(stack, v, n.Value)
	case *Enum:
		for _, item := range n.Items {
			walk(stack, v, item)
		}
		walkAnnotations(stack, v, n.Annotations)
	case *EnumItem:
		walkAnnotations(stack, v, n.Annotations)
	case *Field:
		walk(stack, v, n.Type)
		walk(stack, v, n.Default)
		walkAnnotations(stack, v, n.Annotations)
	case *Function:
		walk(stack, v, n.ReturnType)
		walkFields(stack, v, n.Parameters)
		walkFields(stack, v, n.Exceptions)
		walkAnnotations(stack, v, n.Annotations)
	case ListType:
		walk(stack, v, n.ValueType)
		walkAnnotations(stack, v, n.Annotations)
	case MapType:
		walk(stack, v, n.KeyType)
		walk(stack, v, n.ValueType)
		walkAnnotations(stack, v, n.Annotations)
	case *Program:
		for _, h := range n.Headers {
			walk(stack, v, h)
		}
		for _, d := range n.Definitions {
			walk(stack, v, d)
		}
	case *Service:
		for _, function := range n.Functions {
			walk(stack, v, function)
		}
		walkAnnotations(stack, v, n.Annotations)
	case SetType:
		walk(stack, v, n.ValueType)
		walkAnnotations(stack, v, n.Annotations)
	case *Struct:
		walkFields(stack, v, n.Fields)
		walkAnnotations(stack, v, n.Annotations)
	case *Typedef:
		walk(stack, v, n.Type)
		walkAnnotations(stack, v, n.Annotations)
	}
}

func walkFields(stack nodeStack, v Visitor, fs []*Field) {
	for _, f := range fs {
		walk(stack, v, f)
	}
}

func walkAnnotations(stack nodeStack, v Visitor, anns []*Annotation) {
	for _, ann := range anns {
		walk(stack, v, ann)
	}
}
