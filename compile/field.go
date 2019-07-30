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

package compile

import (
	"fmt"
	"math"

	"go.uber.org/thriftrw/ast"
)

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

// isRequired checks if a field should be required based on the
// fieldRequiredness setting. An error is returned if the specified requiredness
// is disallowed by this configuration.
func (r fieldRequiredness) isRequired(src *ast.Field) (bool, error) {
	switch r {
	case explicitRequiredness:
		if src.Requiredness == ast.Unspecified {
			return false, requirednessRequiredError{
				FieldName: src.Name,
				Line:      src.Line,
			}
		}
	case noRequiredFields:
		if src.Requiredness == ast.Required {
			return false, cannotBeRequiredError{
				FieldName: src.Name,
				Line:      src.Line,
			}
		}
	default:
		// do nothing
	}

	// A field is considered required only if it was marked required AND it
	// does not have a default value.
	return (src.Requiredness == ast.Required && src.Default == nil), nil
}

// fieldOptions controls the behavior of field compilation.
//
// requiredness may be used to control how fields treat required/optional
// specifiers and how it behaves when it is absent.
//
// disallowDefaultValue specifies whether the field is allowed to have a default
// value.
type fieldOptions struct {
	requiredness         fieldRequiredness
	disallowDefaultValue bool
}

// FieldSpec represents a single field of a struct or parameter list.
type FieldSpec struct {
	ID          int16
	Name        string
	Type        TypeSpec
	Required    bool
	Doc         string
	Default     ConstantValue
	Annotations Annotations
}

// compileField compiles the given Field source into a FieldSpec.
func compileField(src *ast.Field, options fieldOptions) (*FieldSpec, error) {
	if src.ID < 1 || src.ID > math.MaxInt16 {
		return nil, fieldIDOutOfBoundsError{ID: src.ID, Name: src.Name}
	}

	required, err := options.requiredness.isRequired(src)
	if err != nil {
		return nil, err
	}

	if options.disallowDefaultValue && src.Default != nil {
		return nil, defaultValueNotAllowedError{
			FieldName: src.Name,
			Line:      src.Line,
		}
	}

	annotations, err := compileAnnotations(src.Annotations)
	if err != nil {
		return nil, compileError{
			Target: src.Name,
			Line:   src.Line,
			Reason: err,
		}
	}

	typ, err := compileTypeReference(src.Type)
	if err != nil {
		return nil, compileError{
			Target: src.Name,
			Line:   src.Line,
			Reason: err,
		}
	}

	return &FieldSpec{
		// TODO(abg): perform bounds check on field ID
		ID:          int16(src.ID),
		Name:        src.Name,
		Type:        typ,
		Doc:         src.Doc,
		Required:    required,
		Default:     compileConstantValue(src.Default),
		Annotations: annotations,
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
		f.Default, err = f.Default.Link(scope, f.Type)
	}
	return err
}

// ThriftAnnotations returns all associated annotations.
func (f *FieldSpec) ThriftAnnotations() Annotations {
	return f.Annotations
}

// FieldGroup represents a collection of fields for struct-like types.
type FieldGroup []*FieldSpec

// compileFields compiles a collection of AST fields into a FieldGroup.
func compileFields(src []*ast.Field, options fieldOptions) (FieldGroup, error) {
	fieldsNS := newNamespace(caseSensitive)
	usedIDs := make(map[int16]string)

	fields := make([]*FieldSpec, 0, len(src))
	for _, astField := range src {
		if err := fieldsNS.claim(astField.Name, astField.Line); err != nil {
			return nil, compileError{
				Target: astField.Name,
				Line:   astField.Line,
				Reason: err,
			}
		}

		field, err := compileField(astField, options)
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

		fields = append(fields, field)
		usedIDs[field.ID] = field.Name
	}

	return FieldGroup(fields), nil
}

// FindByName retrieves the FieldSpec for the field with the given name.
func (fg FieldGroup) FindByName(name string) (*FieldSpec, error) {
	for _, field := range fg {
		if field.Name == name {
			return field, nil
		}
	}
	return nil, fmt.Errorf("unknown field %v", name)
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

// ForEachTypeReference applies the given function on each TypeSpec in the
// FieldGroup.
func (fg FieldGroup) ForEachTypeReference(f func(TypeSpec) error) error {
	for _, field := range fg {
		if err := f(field.Type); err != nil {
			return err
		}
	}
	return nil
}
