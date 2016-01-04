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

// fieldRequiredness controls how fields treat required/optional specifiers.
type fieldRequiredness int

const (
	// defaultToOptional states that the field MAY have a required/optional
	// specifier, and it will default to optional if not specified.
	//
	// This is used for function parameters.
	defaultToOptional fieldRequiredness = iota // default

	// explicitRequiredness states that the field must explicitly specify
	// whether it is required or optional.
	//
	// This is used for all structs and exceptions.
	explicitRequiredness

	// noRequiredFields states that the field MUST NOT have any fields marked
	// as required but it may have fields marked as optional.
	//
	// This is used primarily for unions.
	noRequiredFields
)

// FieldSpec represents a single field of a struct or parameter list.
type FieldSpec struct {
	ID       int16
	Name     string
	Type     TypeSpec
	Required bool
	Default  ast.ConstantValue
}

// compileField compiles the given Field source into a FieldSpec.
//
// requireRequiredness specifies whether the field must explicitly specify
// whether it's required or optional. If this is false and the field is not
// explicitly marked required, it will be treated as optional.
func compileField(src *ast.Field, req fieldRequiredness) (*FieldSpec, error) {
	switch req {
	case explicitRequiredness:
		if src.Requiredness == ast.Unspecified {
			return nil, requirednessRequiredError{
				FieldName: src.Name,
				Line:      src.Line,
			}
		}
	case noRequiredFields:
		if src.Requiredness == ast.Required {
			return nil, cannotBeRequiredError{
				FieldName: src.Name,
				Line:      src.Line,
			}
		}
	default:
		// do nothing
	}

	return &FieldSpec{
		// TODO(abg): perform bounds check on field ID
		ID:       int16(src.ID),
		Name:     src.Name,
		Type:     compileType(src.Type),
		Required: src.Requiredness == ast.Required,
		Default:  src.Default,
	}, nil
}

// ThriftName is the name of the FieldSpec as it appears in the Thrift file.
func (f *FieldSpec) ThriftName() string {
	return f.Name
}

// Link links together any references made by the FieldSpec.
func (f *FieldSpec) Link(scope Scope) (err error) {
	if f.Type, err = f.Type.Link(scope); err != nil {
		return err
	}
	if f.Default != nil {
		err = verifyConstantValue(f.Default, scope)
	}
	return err
}

// FieldGroup represents a collection of fields for struct-like types.
type FieldGroup map[string]*FieldSpec

// compileFields compiles a collection of AST fields into a FieldGroup.
func compileFields(src []*ast.Field, req fieldRequiredness) (FieldGroup, error) {
	fieldsNS := newNamespace(caseInsensitive)
	usedIDs := make(map[int16]string)

	fields := make(map[string]*FieldSpec)
	for _, astField := range src {
		if err := fieldsNS.claim(astField.Name, astField.Line); err != nil {
			return nil, compileError{
				Target: astField.Name,
				Line:   astField.Line,
				Reason: err,
			}
		}

		field, err := compileField(astField, req)
		if err != nil {
			return nil, compileError{
				Target: astField.Name,
				Line:   astField.Line,
				Reason: err,
			}
		}

		if conflictingField, ok := usedIDs[field.ID]; ok {
			return nil, compileError{
				Target: astField.Name,
				Line:   astField.Line,
				Reason: fieldIDConflictError{
					ID:   field.ID,
					Name: conflictingField,
				},
			}
		}

		fields[field.Name] = field
		usedIDs[field.ID] = field.Name
	}

	return FieldGroup(fields), nil
}

// Link resolves references made by fields inside the FieldGroup.
func (fg FieldGroup) Link(scope Scope) error {
	for _, field := range fg {
		if err := field.Link(scope); err != nil {
			return err
		}
	}

	return nil
}
