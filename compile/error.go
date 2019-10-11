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
	"strings"

	"go.uber.org/thriftrw/ast"
)

// fileReadError is raised when there's an error reading a file.
type fileReadError struct {
	Path   string
	Reason error
}

func (e fileReadError) Error() string {
	return fmt.Sprintf("could not read file %q: %v", e.Path, e.Reason)
}

// parseError is raised when there's an error parsing a Thrift file.
type parseError struct {
	Path   string
	Reason error
}

func (e parseError) Error() string {
	return fmt.Sprintf("could not parse file %q: %v", e.Path, e.Reason)
}

type fileCompileError struct {
	Path   string
	Reason error
}

func (e fileCompileError) Error() string {
	return fmt.Sprintf("could not compile file %q: %v", e.Path, e.Reason)
}

// includeAsDisabledError is raised when the user attempts to use the include-as
// syntax without explicitly enabling it.
type includeAsDisabledError struct{}

func (e includeAsDisabledError) Error() string {
	return "include-as syntax is currently disabled"
}

// includeHyphenatedFileNameError is raised when the user attempts to
// include hyphenated file names.
type includeHyphenatedFileNameError struct{}

func (e includeHyphenatedFileNameError) Error() string {
	return "cannot include hyphenated Thrift files"
}

// includeError is raised when there is an error including another Thrift
// file.
type includeError struct {
	Include *ast.Include
	Reason  error
}

func (e includeError) Error() string {
	return fmt.Sprintf(
		"cannot include %q as %q on line %d: %v",
		e.Include.Path, e.Include.Name, e.Include.Line, e.Reason,
	)
}

// definitionError is raised when there was an error compiling a definition
// from the Thrift file.
type definitionError struct {
	Definition ast.Definition
	Reason     error
}

func (e definitionError) Error() string {
	return fmt.Sprintf(
		"cannot define %q on line %d: %v",
		e.Definition.Info().Name, e.Definition.Info().Line, e.Reason,
	)
}

// compileError is a general error raised while trying to compile components
// of the Thrift file.
type compileError struct {
	Target string
	Line   int
	Reason error
}

func (e compileError) Error() string {
	msg := fmt.Sprintf("cannot compile %q", e.Target)
	if e.Line > 0 {
		msg += fmt.Sprintf(" on line %d", e.Line)
	}
	if e.Reason != nil {
		msg += fmt.Sprintf(": %v", e.Reason)
	}
	return msg
}

// referenceError is raised when there's an error resolving a reference.
type referenceError struct {
	Target    string
	Line      int
	ScopeName string
	Reason    error
}

func (e referenceError) Error() string {
	msg := fmt.Sprintf("could not resolve reference %q", e.Target)
	if e.Line > 0 {
		msg += fmt.Sprintf(" on line %d", e.Line)
	}
	if len(e.ScopeName) > 0 {
		msg += fmt.Sprintf(" in %q", e.ScopeName)
	}
	if e.Reason != nil {
		msg += fmt.Sprintf(": %v", e.Reason)
	}
	return msg
}

type unrecognizedModuleError struct {
	Name   string
	Reason error
}

func (e unrecognizedModuleError) Error() string {
	msg := fmt.Sprintf("unknown module %q", e.Name)
	if e.Reason != nil {
		msg += fmt.Sprintf(": %v", e.Reason)
	}
	return msg
}

type unrecognizedEnumItemError struct {
	EnumName string
	ItemName string
}

func (e unrecognizedEnumItemError) Error() string {
	return fmt.Sprintf(
		"enum %q does not have an item named %q", e.EnumName, e.ItemName,
	)
}

// lookupError is raised by Module if the Lookup* functions are called with
// unknown values.
type lookupError struct {
	Name string
}

func (e lookupError) Error() string {
	return fmt.Sprintf("unknown identifier %q", e.Name)
}

type requirednessRequiredError struct {
	FieldName string
	Line      int
}

func (e requirednessRequiredError) Error() string {
	return fmt.Sprintf(
		"field %q on line %d is not marked required or optional",
		e.FieldName, e.Line,
	)
}

type cannotBeRequiredError struct {
	FieldName string
	Line      int
}

func (e cannotBeRequiredError) Error() string {
	return fmt.Sprintf(
		"field %q on line %d is marked as required but it cannot be required",
		e.FieldName, e.Line,
	)
}

type defaultValueNotAllowedError struct {
	FieldName string
	Line      int
}

func (e defaultValueNotAllowedError) Error() string {
	return fmt.Sprintf(
		"field %q on line %d cannot have a default value", e.FieldName, e.Line,
	)
}

type fieldIDConflictError struct {
	ID   int16
	Name string
}

func (e fieldIDConflictError) Error() string {
	return fmt.Sprintf("field %q has already used ID %d", e.Name, e.ID)
}

type fieldIDOutOfBoundsError struct {
	ID   int
	Name string
}

func (e fieldIDOutOfBoundsError) Error() string {
	return fmt.Sprintf(
		"field ID %v of %q is out of bounds: "+
			"field IDs must be in the range [1, 32767]", e.ID, e.Name)
}

type oneWayCannotReturnError struct {
	Name string
}

func (e oneWayCannotReturnError) Error() string {
	return fmt.Sprintf(
		"function %q cannot return values or raise exceptions: %q is oneway",
		e.Name, e.Name,
	)
}

type notAnExceptionError struct {
	TypeName  string
	FieldName string
}

func (e notAnExceptionError) Error() string {
	return fmt.Sprintf(
		"field %q with type %q is not an exception", e.FieldName, e.TypeName,
	)
}

type typeReferenceCycleError struct {
	Nodes []TypeSpec
}

func (e typeReferenceCycleError) Error() string {
	// Outputs:
	//
	// 	found a type reference cycle:
	// 	    foo (a.thrift)
	// 	 -> bar (b.thrift)
	// 	 -> foo (a.thrift)
	//
	// File names are omitted if all types are from the same file.

	files := make(map[string]struct{})
	for _, t := range e.Nodes {
		file := t.ThriftFile()
		if file != "" {
			files[t.ThriftFile()] = struct{}{}
		}
	}
	includeFileName := len(files) > 1

	lines := make([]string, 0, len(e.Nodes)+1)
	lines = append(lines, "found a type reference cycle:")
	for i, t := range e.Nodes {
		line := " "
		if i == 0 {
			line += "   "
		} else {
			line += "-> "
		}

		file := t.ThriftFile()
		if file != "" && includeFileName {
			line += fmt.Sprintf("%v (%v)", t.ThriftName(), file)
		} else {
			line += t.ThriftName()
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// Failure to cast a Constantvalue to a specific type.
type constantValueCastError struct {
	Value  ConstantValue
	Type   TypeSpec
	Reason error // optional
}

func (e constantValueCastError) Error() string {
	s := fmt.Sprintf("cannot cast %v to %q", e.Value, e.Type.ThriftName())
	if e.Reason != nil {
		s += fmt.Sprintf(": %v", e.Reason)
	}
	return s
}

// Failure to cast a specific field of a struct literal.
type constantStructFieldCastError struct {
	FieldName string
	Reason    error
}

func (e constantStructFieldCastError) Error() string {
	return fmt.Sprintf("failed to cast field %q: %v", e.FieldName, e.Reason)
}

// Failure to cast a value referenced by a named constant.
type constantCastError struct {
	Name   string
	Reason error
}

func (e constantCastError) Error() string {
	return fmt.Sprintf("failed to cast constant %q: %v", e.Name, e.Reason)
}

type annotationConflictError struct {
	Reason error
}

func (e annotationConflictError) Error() string {
	return fmt.Sprintf("annotation conflict: %v", e.Reason)
}
