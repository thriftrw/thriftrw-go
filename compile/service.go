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
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/wire"
)

// ServiceSpec is a collection of named functions.
type ServiceSpec struct {
	linkOnce

	Name        string
	File        string
	Parent      *ServiceSpec
	Functions   map[string]*FunctionSpec
	Annotations Annotations

	parentSrc *ast.ServiceReference
}

func compileService(file string, src *ast.Service) (*ServiceSpec, error) {
	serviceNS := newNamespace(caseInsensitive)

	functions := make(map[string]*FunctionSpec)
	for _, astFunction := range src.Functions {
		if err := serviceNS.claim(astFunction.Name, astFunction.Line); err != nil {
			return nil, compileError{
				Target: src.Name + "." + astFunction.Name,
				Line:   astFunction.Line,
				Reason: err,
			}
		}

		function, err := compileFunction(astFunction)
		if err != nil {
			return nil, compileError{
				Target: src.Name + "." + astFunction.Name,
				Line:   astFunction.Line,
				Reason: err,
			}
		}

		functions[function.Name] = function
	}

	annotations, err := compileAnnotations(src.Annotations)
	if err != nil {
		return nil, compileError{
			Target: src.Name,
			Line:   src.Line,
			Reason: err,
		}
	}

	return &ServiceSpec{
		Name:        src.Name,
		File:        file,
		Functions:   functions,
		Annotations: annotations,
		parentSrc:   src.Parent,
	}, nil
}

// resolveService resolves a ServiceReference in the given scope.
func resolveService(src ast.ServiceReference, scope Scope) (*ServiceSpec, error) {
	s, err := scope.LookupService(src.Name)
	if err == nil {
		err = s.Link(scope)
		return s, err
	}

	mname, iname := splitInclude(src.Name)
	if len(mname) == 0 {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	includedScope, err := getIncludedScope(scope, mname)
	if err != nil {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	return resolveService(ast.ServiceReference{Name: iname}, includedScope)
}

// Link resolves any references made by the given service.
func (s *ServiceSpec) Link(scope Scope) error {
	if s.linked() {
		return nil
	}

	if s.parentSrc != nil {
		parent, err := resolveService(*s.parentSrc, scope)
		if err != nil {
			return compileError{
				Target: s.Name,
				Reason: referenceError{
					Target:    s.parentSrc.Name,
					Line:      s.parentSrc.Line,
					ScopeName: scope.GetName(),
					Reason:    err,
				},
			}
		}

		if err := parent.Link(scope); err != nil {
			return compileError{Target: s.Name, Reason: err}
		}

		s.Parent = parent
		s.parentSrc = nil
	}

	for _, function := range s.Functions {
		if err := function.Link(scope); err != nil {
			return compileError{
				Target: s.Name + "." + function.Name,
				Reason: err,
			}
		}
	}

	return nil
}

// ThriftFile is the Thrift file in which this service was defined.
func (s *ServiceSpec) ThriftFile() string {
	return s.File
}

// FunctionSpec is a single function inside a Service.
type FunctionSpec struct {
	linkOnce

	Name        string
	ArgsSpec    ArgsSpec
	ResultSpec  *ResultSpec // nil if OneWay is true
	OneWay      bool
	Annotations Annotations
}

func compileFunction(src *ast.Function) (*FunctionSpec, error) {
	args, err := compileArgSpec(src.Parameters)
	if err != nil {
		return nil, compileError{
			Target: src.Name,
			Line:   src.Line,
			Reason: err,
		}
	}

	var result *ResultSpec
	if src.OneWay {
		// oneway can't have a return type or exceptions
		if src.ReturnType != nil || len(src.Exceptions) > 0 {
			return nil, oneWayCannotReturnError{Name: src.Name}
		}
	} else {
		result, err = compileResultSpec(src.ReturnType, src.Exceptions)
		if err != nil {
			return nil, compileError{
				Target: src.Name,
				Line:   src.Line,
				Reason: err,
			}
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

	return &FunctionSpec{
		Name:        src.Name,
		ArgsSpec:    args,
		ResultSpec:  result,
		Annotations: annotations,
		OneWay:      src.OneWay,
	}, nil
}

// Link resolves any references made by the given function.
func (f *FunctionSpec) Link(scope Scope) error {
	if f.linked() {
		return nil
	}

	if err := f.ArgsSpec.Link(scope); err != nil {
		return compileError{Target: f.Name, Reason: err}
	}

	if f.ResultSpec != nil {
		if err := f.ResultSpec.Link(scope); err != nil {
			return compileError{Target: f.Name, Reason: err}
		}
	}

	return nil
}

// MethodName returns the method name for this function.
func (f *FunctionSpec) MethodName() string {
	return f.Name
}

// CallType returns the envelope type that is used when making enveloped
// requests for this function.
func (f *FunctionSpec) CallType() wire.EnvelopeType {
	if f.OneWay {
		return wire.OneWay
	}
	return wire.Call
}

// ArgsSpec contains information about a Function's arguments.
type ArgsSpec FieldGroup

func compileArgSpec(args []*ast.Field) (ArgsSpec, error) {
	fields, err := compileFields(
		args,
		fieldOptions{requiredness: defaultToOptional},
	)
	return ArgsSpec(fields), err
}

// Link resolves references made by the ArgsSpec.
func (as ArgsSpec) Link(scope Scope) error {
	return FieldGroup(as).Link(scope)
}

// ResultSpec contains information about a Function's result type.
type ResultSpec struct {
	ReturnType TypeSpec
	Exceptions FieldGroup
}

func compileResultSpec(returnType ast.Type, exceptions []*ast.Field) (*ResultSpec, error) {
	var excFields FieldGroup

	if len(exceptions) > 0 {
		var err error
		excFields, err = compileFields(
			exceptions,
			fieldOptions{
				requiredness:         noRequiredFields,
				disallowDefaultValue: true,
			},
		)
		if err != nil {
			return nil, err
		}
	}

	typ, err := compileTypeReference(returnType)
	if err != nil {
		return nil, err
	}

	return &ResultSpec{
		ReturnType: typ,
		Exceptions: excFields,
	}, nil
}

// Link resolves any references made by the return type or exceptions in the
// ResultSpec.
func (rs *ResultSpec) Link(scope Scope) (err error) {
	if rs.ReturnType != nil {
		rs.ReturnType, err = rs.ReturnType.Link(scope)
		if err != nil {
			return err
		}
	}

	if err := rs.Exceptions.Link(scope); err != nil {
		return err
	}

	// verify that everything listed under throws is an exception.
	for _, exception := range rs.Exceptions {
		spec, ok := exception.Type.(*StructSpec)
		if !ok || spec.Type != ast.ExceptionType {
			return notAnExceptionError{
				FieldName: exception.ThriftName(),
				TypeName:  exception.Type.ThriftName(),
			}
		}
	}

	return nil
}
