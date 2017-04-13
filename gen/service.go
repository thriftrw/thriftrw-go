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

package gen

import (
	"bytes"
	"fmt"
	"go/token"
	"strings"

	"go.uber.org/thriftrw/compile"
)

// Service generates code for the given service.
//
// Returns a map from file name to contents for that file. The file names are
// relative to the package directory for the service.
func Service(g Generator, s *compile.ServiceSpec) (map[string]*bytes.Buffer, error) {
	files := make(map[string]*bytes.Buffer)

	for _, functionName := range sortStringKeys(s.Functions) {
		fileName := fmt.Sprintf("%s_%s.go", strings.ToLower(s.Name), strings.ToLower(functionName))

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
		files[fileName] = buff
	}

	return files, nil
}

// ServiceFunction generates code for the given function of the given service.
func ServiceFunction(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	argsGen := fieldGroupGenerator{
		Namespace: NewNamespace(),
		Name:      functionNamePrefix(s, f) + "Args",
		Fields:    compile.FieldGroup(f.ArgsSpec),
	}
	if err := argsGen.Generate(g); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}
	if err := functionArgsEnveloper(g, s, f); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}

	if err := functionHelper(g, s, f); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}

	if f.ResultSpec == nil {
		return nil
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
		Namespace:       NewNamespace(),
		Name:            functionNamePrefix(s, f) + "Result",
		Fields:          resultFields,
		IsUnion:         true,
		AllowEmptyUnion: f.ResultSpec.ReturnType == nil,
	}
	if err := resultGen.Generate(g); err != nil {
		return wrapGenerateError(fmt.Sprintf("%s.%s", s.Name, f.Name), err)
	}
	if err := functionResponseEnveloper(g, s, f); err != nil {
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

func functionHelper(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	return g.DeclareFromTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>

		var <$prefix>Helper = struct{
			Args func(<params $f>) *<$prefix>Args
			<if not $f.OneWay>
				IsException func(error) bool
				<if $f.ResultSpec.ReturnType>
					WrapResponse func(
						<typeReference $f.ResultSpec.ReturnType>,
						error) (*<$prefix>Result, error)
					UnwrapResponse func(*<$prefix>Result) (
						<typeReference $f.ResultSpec.ReturnType>, error)
				<else>
					WrapResponse func(error) (*<$prefix>Result, error)
					UnwrapResponse func(*<$prefix>Result) error
				<end>
			<end>
		}{}

		func init() {
			<$prefix>Helper.Args = <newArgs .Service $f>
			<if not $f.OneWay>
				<$prefix>Helper.IsException = <isException $f>
				<$prefix>Helper.WrapResponse = <wrapResponse .Service $f>
				<$prefix>Helper.UnwrapResponse = <unwrapResponse .Service $f>
			<end>
		}
		`,
		struct {
			Service  *compile.ServiceSpec
			Function *compile.FunctionSpec
		}{
			Service:  s,
			Function: f,
		},
		TemplateFunc("params", functionParams),
		TemplateFunc("isException", functionIsException),
		TemplateFunc("newArgs", functionNewArgs),
		TemplateFunc("wrapResponse", functionWrapResponse),
		TemplateFunc("unwrapResponse", functionUnwrapResponse),
		TemplateFunc("namePrefix", functionNamePrefix),
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
func functionNewArgs(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>
		<$params := newNamespace>
		func(
			<range $f.ArgsSpec>
				<if .Required>
					<$params.NewName .Name> <typeReference .Type>,
				<else>
					<$params.NewName .Name> <typeReferencePtr .Type>,
				<end>
			<end>
		) *<$prefix>Args {
			return &<$prefix>Args{
			<range $f.ArgsSpec>
				<if .Required>
					<goCase .Name>: <$params.Rotate .Name>,
				<else>
					<goCase .Name>: <$params.Rotate .Name>,
				<end>
			<end>
			}
		}
		`,
		struct {
			Service  *compile.ServiceSpec
			Function *compile.FunctionSpec
		}{
			Service:  s,
			Function: f,
		},
		TemplateFunc("namePrefix", functionNamePrefix))
}

// functionWrapResponse generates an expression that provides the WrapResponse
// function for the given Thrift function.
func functionWrapResponse(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>

		<if $f.ResultSpec.ReturnType>
			func(success <typeReference $f.ResultSpec.ReturnType>,
				err error) (*<$prefix>Result, error) {
				if err == nil {
					<if isPrimitiveType $f.ResultSpec.ReturnType>
						return &<$prefix>Result{Success: &success}, nil
					<else>
						return &<$prefix>Result{Success: success}, nil
					<end>
				}
		<else>
			func(err error) (*<$prefix>Result, error) {
				if err == nil {
					return &<$prefix>Result{}, nil
				}
		<end>
				<if $f.ResultSpec.Exceptions>
					switch e := err.(type) {
						<range $f.ResultSpec.Exceptions>
						case <typeReferencePtr .Type>:
							if e == nil {
								return nil, <import "errors">.New(
									"WrapResponse received non-nil error type with nil value for <$prefix>Result.<goCase .Name>")
							}
							return &<$prefix>Result{<goCase .Name>: e}, nil
						<end>
					}
				<end>
				return nil, err
			}
		`,
		struct {
			Service  *compile.ServiceSpec
			Function *compile.FunctionSpec
		}{
			Service:  s,
			Function: f,
		},
		TemplateFunc("namePrefix", functionNamePrefix))
}

// functionUnwrapResponse generates an expression that provides the
// UnwrapResponse function for the given Thrift function.
func functionUnwrapResponse(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>

		<if $f.ResultSpec.ReturnType>
			func(result *<$prefix>Result) (
				success <typeReference $f.ResultSpec.ReturnType>,
				err error) {
		<else>
			func(result *<$prefix>Result) (err error) {
		<end>
				<range $f.ResultSpec.Exceptions>
					if result.<goCase .Name> != nil {
						err = result.<goCase .Name>
						return
					}
				<end>

				// TODO unrecognized exceptions

				<if $f.ResultSpec.ReturnType>
					if result.Success != nil {
						<if isPrimitiveType $f.ResultSpec.ReturnType>
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
		`, struct {
			Service  *compile.ServiceSpec
			Function *compile.FunctionSpec
		}{
			Service:  s,
			Function: f,
		},
		TemplateFunc("namePrefix", functionNamePrefix))
}

func functionArgsEnveloper(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	// TODO: Figure out naming conflicts with user fields.
	return g.DeclareFromTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>

		<$wire := import "go.uber.org/thriftrw/wire">
		<$v := newVar "v">

		func (<$v> *<$prefix>Args) MethodName() string {
			return "<$f.MethodName>"
		}

		func (<$v> *<$prefix>Args) EnvelopeType() <$wire>.EnvelopeType {
			return <$wire>.<$f.CallType.String>
		}
		`, struct {
			Service  *compile.ServiceSpec
			Function *compile.FunctionSpec
		}{
			Service:  s,
			Function: f,
		},
		TemplateFunc("namePrefix", functionNamePrefix))

}

func functionResponseEnveloper(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	return g.DeclareFromTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>

		<$wire := import "go.uber.org/thriftrw/wire">
		<$v := newVar "v">

		func (<$v> *<$prefix>Result) MethodName() string {
			return "<$f.MethodName>"
		}

		func (<$v> *<$prefix>Result) EnvelopeType() <$wire>.EnvelopeType {
			return <$wire>.Reply
		}
		`, struct {
			Service  *compile.ServiceSpec
			Function *compile.FunctionSpec
		}{
			Service:  s,
			Function: f,
		},
		TemplateFunc("namePrefix", functionNamePrefix))

}

func functionNamePrefix(s *compile.ServiceSpec, f *compile.FunctionSpec) string {
	return fmt.Sprintf("%s_%s_", goCase(s.Name), goCase(f.Name))
}
