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

import (
	"github.com/thriftrw/thriftrw-go/ast"
	"github.com/thriftrw/thriftrw-go/wire"
)

// StructSpec represents a structure defined in the Thrift file.
type StructSpec struct {
	linkOnce

	Name   string
	File   string
	Type   ast.StructureType
	Fields FieldGroup
}

// compileStruct compiles a struct AST into a StructSpec.
func compileStruct(file string, src *ast.Struct) (*StructSpec, error) {
	opts := fieldOptions{requiredness: explicitRequiredness}

	if src.Type == ast.UnionType {
		opts.requiredness = noRequiredFields
		opts.disallowDefaultValue = true
	}

	fields, err := compileFields(src.Fields, opts)
	if err != nil {
		return nil, compileError{
			Target: src.Name,
			Line:   src.Line,
			Reason: err,
		}
	}

	return &StructSpec{
		Name:   src.Name,
		File:   file,
		Type:   src.Type,
		Fields: fields,
	}, nil
}

// Link links together all references in the StructSpec.
func (s *StructSpec) Link(scope Scope) (TypeSpec, error) {
	if s.linked() {
		return s, nil
	}

	err := s.Fields.Link(scope)
	return s, err
}

// TypeCode for structs.
func (s *StructSpec) TypeCode() wire.Type {
	return wire.TStruct
}

// ThriftName of the StructSpec.
func (s *StructSpec) ThriftName() string {
	return s.Name
}

// ThriftFile of the StructSpec.
func (s *StructSpec) ThriftFile() string {
	return s.File
}

// IsExceptionType returns true if the StructSpec represents an exception
// declaration.
func (s *StructSpec) IsExceptionType() bool {
	return s.Type == ast.ExceptionType
}
