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
	"fmt"

	"github.com/uber/thriftrw-go/ast"
)

// includeAsDisabledError is raised when the user attempts to use the include-as
// syntax without explicitly enabling it.
type includeAsDisabledError struct{}

func (e includeAsDisabledError) Error() string {
	return "include-as syntax is currently disabled"
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
	return fmt.Sprintf(
		"cannot compile %q on line %d: %v", e.Target, e.Line, e.Reason,
	)
}

// referenceError is raised when there's an error resolving a reference.
type referenceError struct {
	Target string
	Line   int
	Reason error
}

func (e referenceError) Error() string {
	return fmt.Sprintf(
		"could not resolve reference %q on line %d: %v",
		e.Target, e.Line, e.Reason,
	)
}

type unrecognizedModuleError struct {
	Name string
}

func (e unrecognizedModuleError) Error() string {
	return fmt.Sprintf("unknown module %q", e.Name)
}

// lookupError is raised when an unknown identifier is requested via the
// Lookup* methods.
type lookupError struct {
	Name   string
	Reason error
}

func (e lookupError) Error() string {
	msg := fmt.Sprintf("unknown identifier %q", e.Name)
	if e.Reason != nil {
		msg = fmt.Sprintf("%s: %v", msg, e.Reason)
	}
	return msg
}

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
