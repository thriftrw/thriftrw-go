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

// Services generates code for all services into a single file and stores the code
// in the generator to be written.
func Services(g Generator, services map[string]*compile.ServiceSpec) error {
	for _, serviceName := range sortStringKeys(services) {
		s := services[serviceName]
		for _, functionName := range sortStringKeys(s.Functions) {
			function := s.Functions[functionName]
			if err := ServiceFunction(g, s, function); err != nil {
				return fmt.Errorf(
					"could not generate types for %s.%s: %v",
					s.Name, functionName, err)
			}
		}
	}

	return nil
}

// ServiceFunction generates code for the given function of the given service.
func ServiceFunction(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	argsName := functionNamePrefix(s, f) + "Args"
	argsGen := fieldGroupGenerator{
		Namespace: NewNamespace(),
		Name:      argsName,
		Fields:    compile.FieldGroup(f.ArgsSpec),
		Doc: fmt.Sprintf(
			"%v represents the arguments for the %v.%v function.\n\n"+
				"The arguments for %v are sent and received over the wire as this struct.",
			argsName, s.Name, f.Name, f.Name,
		),
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
			Doc:  fmt.Sprintf("Value returned by %v after a successful execution.", f.Name),
		})
	}
	resultFields = append(resultFields, f.ResultSpec.Exceptions...)

	resultName := functionNamePrefix(s, f) + "Result"
	resultDoc := fmt.Sprintf(
		"%v represents the result of a %v.%v function call.\n\n"+
			"The result of a %v execution is sent and received over the wire as this struct.",
		resultName, s.Name, f.Name, f.Name,
	)
	if f.ResultSpec.ReturnType != nil {
		resultDoc += fmt.Sprintf("\n\nSuccess is set only if the function did not throw an exception.")
	}

	resultGen := fieldGroupGenerator{
		Namespace:       NewNamespace(),
		Name:            resultName,
		Fields:          resultFields,
		IsUnion:         true,
		AllowEmptyUnion: f.ResultSpec.ReturnType == nil,
		Doc:             resultDoc,
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
		<- $params := newNamespace ->
		<- range .ArgsSpec>
			<- if .Required>
				<$params.NewName .Name> <typeReference .Type>,
			<- else>
				<$params.NewName .Name> <typeReferencePtr .Type>,
			<- end ->
		<end>
		`, f)
}

func functionHelper(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) error {
	return g.DeclareFromTemplate(
		`
		<$f := .Function>
		<$prefix := namePrefix .Service $f>

		// <$prefix>Helper provides functions that aid in handling the
		// parameters and return values of the <.Service.Name>.<$f.Name>
		// function.
		var <$prefix>Helper = struct{
			// Args accepts the parameters of <$f.Name> in-order and returns
			// the arguments struct for the function.
			Args func(<params $f>) *<$prefix>Args
			<if not $f.OneWay>
				// IsException returns true if the given error can be thrown
				// by <$f.Name>.
				//
				// An error can be thrown by <$f.Name> only if the
				// corresponding exception type was mentioned in the 'throws'
				// section for it in the Thrift file.
				IsException func(error) bool
				<if $f.ResultSpec.ReturnType>
					// WrapResponse returns the result struct for <$f.Name>
					// given its return value and error.
					//
					// This allows mapping values and errors returned by
					// <$f.Name> into a serializable result struct.
					// WrapResponse returns a non-nil error if the provided
					// error cannot be thrown by <$f.Name>
					//
					//   value, err := <$f.Name>(args)
					//   result, err := <$prefix>Helper.WrapResponse(value, err)
					//   if err != nil {
					//     return fmt.Errorf("unexpected error from <$f.Name>: %v", err)
					//   }
					//   serialize(result)
					WrapResponse func(<typeReference $f.ResultSpec.ReturnType>, error) (*<$prefix>Result, error)

					// UnwrapResponse takes the result struct for <$f.Name>
					// and returns the value or error returned by it.
					//
					// The error is non-nil only if <$f.Name> threw an
					// exception.
					//
					//   result := deserialize(bytes)
					//   value, err := <$prefix>Helper.UnwrapResponse(result)
					UnwrapResponse func(*<$prefix>Result) (<typeReference $f.ResultSpec.ReturnType>, error)
				<else>
					// WrapResponse returns the result struct for <$f.Name>
					// given the error returned by it. The provided error may
					// be nil if <$f.Name> did not fail.
					//
					// This allows mapping errors returned by <$f.Name> into a
					// serializable result struct. WrapResponse returns a
					// non-nil error if the provided error cannot be thrown by
					// <$f.Name>
					//
					//   err := <$f.Name>(args)
					//   result, err := <$prefix>Helper.WrapResponse(err)
					//   if err != nil {
					//     return fmt.Errorf("unexpected error from <$f.Name>: %v", err)
					//   }
					//   serialize(result)
					WrapResponse func(error) (*<$prefix>Result, error)

					// UnwrapResponse takes the result struct for <$f.Name>
					// and returns the erorr returned by it (if any).
					//
					// The error is non-nil only if <$f.Name> threw an
					// exception.
					//
					//   result := deserialize(bytes)
					//   err := <$prefix>Helper.UnwrapResponse(result)
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
		`func(err error) bool {
			switch err.(type) {
			<range .ResultSpec.Exceptions ->
				case <typeReferencePtr .Type>:
					return true
			<end ->
			default:
				return false
			}
		}`, f)
}

// functionNewArgs generates an expression which provides the NewArgs function
// for the given Thrift function.
func functionNewArgs(g Generator, s *compile.ServiceSpec, f *compile.FunctionSpec) (string, error) {
	return g.TextTemplate(
		`
		<- $f := .Function ->
		<- $prefix := namePrefix .Service $f ->
		<- $params := newNamespace ->
		func(
			<- range $f.ArgsSpec>
				<- if .Required>
					<$params.NewName .Name> <typeReference .Type>,
				<- else>
					<$params.NewName .Name> <typeReferencePtr .Type>,
				<- end ->
			<end>
		) *<$prefix>Args {
			return &<$prefix>Args{
			<range $f.ArgsSpec>
				<- if .Required ->
					<goCase .Name>: <$params.Rotate .Name>,
				<- else ->
					<goCase .Name>: <$params.Rotate .Name>,
				<- end>
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
		<- $f := .Function ->
		<- $prefix := namePrefix .Service $f ->

		<- if $f.ResultSpec.ReturnType ->
			func(success <typeReference $f.ResultSpec.ReturnType>, err error) (*<$prefix>Result, error) {
				if err == nil {
					<if isPrimitiveType $f.ResultSpec.ReturnType ->
						return &<$prefix>Result{Success: &success}, nil
					<- else ->
						return &<$prefix>Result{Success: success}, nil
					<- end>
				}
		<- else ->
			func(err error) (*<$prefix>Result, error) {
				if err == nil {
					return &<$prefix>Result{}, nil
				}
		<- end>
				<if $f.ResultSpec.Exceptions>
					switch e := err.(type) {
						<range $f.ResultSpec.Exceptions ->
						case <typeReferencePtr .Type>:
							if e == nil {
								return nil, <import "errors">.New("WrapResponse received non-nil error type with nil value for <$prefix>Result.<goCase .Name>")
							}
							return &<$prefix>Result{<goCase .Name>: e}, nil
						<end ->
					}
				<end>
				return nil, err
			}`,
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
		<- $f := .Function ->
		<- $prefix := namePrefix .Service $f ->

		<- if $f.ResultSpec.ReturnType ->
			func(result *<$prefix>Result) (success <typeReference $f.ResultSpec.ReturnType>, err error) {
		<- else ->
			func(result *<$prefix>Result) (err error) {
		<- end>
				<range $f.ResultSpec.Exceptions ->
					if result.<goCase .Name> != nil {
						err = result.<goCase .Name>
						return
					}
				<end ->

				<if $f.ResultSpec.ReturnType>
					if result.Success != nil {
						<- if isPrimitiveType $f.ResultSpec.ReturnType>
							success = *result.Success
						<- else>
							success = result.Success
						<- end>
						return
					}

					err = <import "errors">.New("expected a non-void result")
				<end ->
				return
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

		// MethodName returns the name of the Thrift function as specified in
		// the IDL, for which this struct represent the arguments.
		//
		// This will always be "<$f.MethodName>" for this struct.
		func (<$v> *<$prefix>Args) MethodName() string {
			return "<$f.MethodName>"
		}

		// EnvelopeType returns the kind of value inside this struct.
		//
		// This will always be <$f.CallType.String> for this struct.
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

		// MethodName returns the name of the Thrift function as specified in
		// the IDL, for which this struct represent the result.
		//
		// This will always be "<$f.MethodName>" for this struct.
		func (<$v> *<$prefix>Result) MethodName() string {
			return "<$f.MethodName>"
		}

		// EnvelopeType returns the kind of value inside this struct.
		//
		// This will always be Reply for this struct.
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
