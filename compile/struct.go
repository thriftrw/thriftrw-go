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

package compile

import "github.com/uber/thriftrw-go/ast"

// StructSpec represents a structure defined in the Thrift file.
type StructSpec struct {
	Name   string
	Fields map[string]FieldSpec
}

// FieldSpec represents a single field of a struct or parameter list.
type FieldSpec struct {
	ID       int16
	Name     string
	Type     TypeSpec
	Required bool
	Default  ast.ConstantValue
}

// ExceptionSpec represents an exception defined in the Thrift file.
type ExceptionSpec StructSpec

// UnionSpec represents a union defined in the Thrift file.
type UnionSpec StructSpec
