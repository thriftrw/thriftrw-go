// Copyright (c) 2016 Uber Technologies, Inc.
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

package gen

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"

	"github.com/thriftrw/thriftrw-go/compile"
)

// Service generates code for the given service.
//
// Returns a map from file name to contents for that file. The file names are
// relative to the package directory for the service.
func Service(g Generator, s *compile.ServiceSpec) (map[string]*bytes.Buffer, error) {
	files := make(map[string]*bytes.Buffer)

	// TODO inherited service functions

	for _, functionName := range sortStringKeys(s.Functions) {
		function := s.Functions[functionName]
		if err := ServiceFunction(g, s, function); err != nil {
			return nil, fmt.Errorf(
				"could not generate types for %s.%s: %v",
				s.Name, functionName, err)
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, token.NewFileSet()); err != nil {
			return nil, fmt.Errorf("could not write %s.%s: %v", s.Name, functionName, err)
		}

		// TODO check for conflicts
		files[strings.ToLower(functionName)+".go"] = buff
	}

	return files, nil
}

// ServiceFunction generates code for the given function of the given service.
func ServiceFunction(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	if f.OneWay {
		return fmt.Errorf(
			"oneway functions are not yet supported: %s.%s is oneway",
			s.Name, f.Name)
	}

	argsGen := fieldGroupGenerator{
		Name:   goCase(f.Name) + "Args",
		Fields: compile.FieldGroup(f.ArgsSpec),
	}
	if err := argsGen.Generate(g); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}

	resultFields := make(compile.FieldGroup, 0, len(f.ResultSpec.Exceptions)+1)
	if f.ResultSpec.ReturnType != nil {
		resultFields = append(resultFields, &compile.FieldSpec{
			ID:   0,
			Name: "success",
			Type: f.ResultSpec.ReturnType,
		})
	}
	resultFields = append(resultFields, f.ResultSpec.Exceptions...)

	resultGen := fieldGroupGenerator{
		Name:   goCase(f.Name) + "Result",
		Fields: resultFields,
	}
	if err := resultGen.Generate(g); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}

	if err := functionHelper(g, f); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}

	// TODO(abg): If we receive unknown exceptions over the wire, we need to
	// throw a generic error.
	return nil
}

// functionParams returns a named parameter list for the given function.
func functionParams(g Generator, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<$params := newNamespace>
		<range .ArgsSpec>
			<if .Required>
				<$params.NewName .Name> <typeReference .Type>,
			<else>
				<$params.NewName .Name> <typeReferencePtr .Type>,
			<end>
		<end>
        `, f)
}

func functionHelper(g Generator, f *compile.FunctionSpec) error {
	return g.DeclareFromTemplate(
		`
		<$name := goCase .Name>

		var <$name>Helper = struct{
			IsException func(error) bool

			Args func(<params .>) *<$name>Args

			<if .ResultSpec.ReturnType>
				WrapResponse func(
					<typeReference .ResultSpec.ReturnType>,
					error) (*<$name>Result, error)
				UnwrapResponse func(*<$name>Result) (
					<typeReference .ResultSpec.ReturnType>, error)
			<else>
				WrapResponse func(error) (*<$name>Result, error)
				UnwrapResponse func(*<$name>Result) error
			<end>
		}{}

		func init() {
			<$name>Helper.IsException = <isException .>
			<$name>Helper.Args = <newArgs .>
			<$name>Helper.WrapResponse = <wrapResponse .>
			<$name>Helper.UnwrapResponse = <unwrapResponse .>
		}
		`,
		f,
		TemplateFunc("params", functionParams),
		TemplateFunc("isException", functionIsException),
		TemplateFunc("newArgs", functionNewArgs),
		TemplateFunc("wrapResponse", functionWrapResponse),
		TemplateFunc("unwrapResponse", functionUnwrapResponse),
	)
}

// functionIsException generates an expression that provides the IsException
// function for the given Thrift function.
func functionIsException(g Generator, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		func(err error) bool {
			switch err.(type) {
			<range .ResultSpec.Exceptions>
				case <typeReferencePtr .Type>:
					return true
			<end>
			default:
				return false
			}
		}
		`, f)
}

// functionNewArgs generates an expression which provides the NewArgs function
// for the given Thrift function.
func functionNewArgs(g Generator, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<$params := newNamespace>
		func(
			<range .ArgsSpec>
				<if .Required>
					<$params.NewName .Name> <typeReference .Type>,
				<else>
					<$params.NewName .Name> <typeReferencePtr .Type>,
				<end>
			<end>
		) *<goCase .Name>Args {
			return &<goCase .Name>Args{
			<range .ArgsSpec>
				<if .Required>
					<goCase .Name>: <$params.Rotate .Name>,
				<else>
					<goCase .Name>: <$params.Rotate .Name>,
				<end>
			<end>
			}
		}
		`, f)
}

// functionWrapResponse generates an expression that provides the WrapResponse
// function for the given Thrift function.
func functionWrapResponse(g Generator, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<$name := goCase .Name>
		<if .ResultSpec.ReturnType>
			func(success <typeReference .ResultSpec.ReturnType>,
				err error) (*<$name>Result, error) {
				if err == nil {
					<if isPrimitiveType .ResultSpec.ReturnType>
						return &<$name>Result{Success: &success}, nil
					<else>
						return &<$name>Result{Success: success}, nil
					<end>
				}
		<else>
			func(err error) (*<$name>Result, error) {
				if err == nil {
					return &<$name>Result{}, nil
				}
		<end>
				<if .ResultSpec.Exceptions>
					switch e := err.(type) {
						<range .ResultSpec.Exceptions>
						case <typeReferencePtr .Type>:
							if e == nil {
								return nil, <import "errors">.New(
									"WrapResponse received non-nil error type with nil value for <$name>Result.<goCase .Name>")
							}
							return &<$name>Result{<goCase .Name>: e}, nil
						<end>
					}
				<end>
				return nil, err
			}
		`, f)
}

// functionUnwrapResponse generates an expression that provides the
// UnwrapResponse function for the given Thrift function.
func functionUnwrapResponse(g Generator, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<if .ResultSpec.ReturnType>
			func(result *<goCase .Name>Result) (
				success <typeReference .ResultSpec.ReturnType>,
				err error) {
		<else>
			func(result *<goCase .Name>Result) (err error) {
		<end>
				<range .ResultSpec.Exceptions>
					if result.<goCase .Name> != nil {
						err = result.<goCase .Name>
						return
					}
				<end>

				// TODO unrecognized exceptions

				<if .ResultSpec.ReturnType>
					if result.Success != nil {
						<if isPrimitiveType .ResultSpec.ReturnType>
							success = *result.Success
						<else>
							success = result.Success
						<end>
						return
					}

					// TODO library-level error type
					err = <import "errors">.New("expected a non-void result")
					return
				<else>
					return
				<end>

			}
		`, f)
}
