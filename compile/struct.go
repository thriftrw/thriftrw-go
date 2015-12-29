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
	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/wire"
)

// StructSpec represents a structure defined in the Thrift file.
type StructSpec struct {
	linkOnce

	Name   string
	Type   ast.StructureType
	Fields map[string]*FieldSpec
}

// compileStruct compiles a struct AST into a StructSpec.
func compileStruct(src *ast.Struct) (*StructSpec, error) {
	structNS := newNamespace(caseInsensitive)

	requiredness := explicitRequiredness
	if src.Type == ast.UnionType {
		requiredness = noRequiredFields
	}

	fields := make(map[string]*FieldSpec)
	usedFieldIDs := make(map[int16]string)
	for _, astField := range src.Fields {
		if err := structNS.claim(astField.Name, astField.Line); err != nil {
			return nil, compileError{
				Target: src.Name + "." + astField.Name,
				Line:   astField.Line,
				Reason: err,
			}
		}

		field, err := compileField(astField, requiredness)
		if err != nil {
			return nil, compileError{
				Target: src.Name + "." + astField.Name,
				Line:   astField.Line,
				Reason: err,
			}
		}

		if conflictName, ok := usedFieldIDs[field.ID]; ok {
			return nil, compileError{
				Target: src.Name + "." + astField.Name,
				Line:   astField.Line,
				Reason: fieldIDConflictError{
					ID:   field.ID,
					Name: conflictName,
				},
			}
		}

		usedFieldIDs[field.ID] = field.Name
		fields[field.Name] = field
	}

	return &StructSpec{
		Name:   src.Name,
		Type:   src.Type,
		Fields: fields,
	}, nil
}

// Link links together all references in the StructSpec.
func (s *StructSpec) Link(scope Scope) (TypeSpec, error) {
	if s.linked() {
		return s, nil
	}

	for _, field := range s.Fields {
		if err := field.Link(scope); err != nil {
			return s, err
		}
	}

	return s, nil
}

// TypeCode for structs.
func (s *StructSpec) TypeCode() wire.Type {
	return wire.TStruct
}

// ThriftName of the StructSpec.
func (s *StructSpec) ThriftName() string {
	return s.Name
}
