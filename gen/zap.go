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
	"fmt"

	"go.uber.org/thriftrw/compile"
)

// NoZapLabel allows opt out of Zap logging for struct fields.
// Fields of Thrift structs will use this annotation to opt-out of being logged
// when that struct is logged. i.e.
//
// 	struct ZapOptOutStruct {
// 		1: required string name
// 		2: required string optout (go.nolog)
// 	}
//
// The above struct will be logged without the optout string.
const NoZapLabel = "go.nolog"

type zapGenerator struct {
	mapG  mapGenerator
	setG  setGenerator
	listG listGenerator
}

// zapEncoder returns the Zap type name of the root spec, determining what type
// the Zap marshaler needs to log it as (i.e. AddString, AppendObject, etc.)
func (z *zapGenerator) zapEncoder(g Generator, spec compile.TypeSpec) string {
	root := compile.RootTypeSpec(spec)

	switch t := root.(type) {
	// Primitives
	case *compile.BoolSpec:
		return "Bool"
	case *compile.I8Spec:
		return "Int8"
	case *compile.I16Spec:
		return "Int16"
	case *compile.I32Spec:
		return "Int32"
	case *compile.I64Spec:
		return "Int64"
	case *compile.DoubleSpec:
		return "Float64"
	case *compile.StringSpec:
		return "String"
	case *compile.BinarySpec:
		return "String" // encode binary as a string and log as string

	// Containers
	case *compile.MapSpec:
		switch compile.RootTypeSpec(t.KeySpec).(type) {
		case *compile.StringSpec:
			return "Object"
		default:
			return "Array"
		}
	case *compile.SetSpec, *compile.ListSpec:
		return "Array"

	// User-defined
	case *compile.EnumSpec, *compile.StructSpec:
		return "Object"
	}
	panic(root)
}

// zapMarshaler takes a TypeSpec, evaluates whether there are underlying elements
// that require more Zap implementation to log everything, and returns a string
// that properly casts the fieldValue, if needed, for logging.
//
// This should be used in conjunction with zapEncoder:
//
//   v := ...
//   enc.Add<zapEncoder .Type>("foo", <zapMarshaler .Type "v">)
//
func (z *zapGenerator) zapMarshaler(g Generator, spec compile.TypeSpec, fieldValue string) (string, error) {
	root := compile.RootTypeSpec(spec)

	if _, ok := spec.(*compile.TypedefSpec); ok {
		// For typedefs, cast to the root type and rely on that functionality.
		rootName, err := typeReference(g, root)
		if err != nil {
			return "", err
		}
		fieldValue = fmt.Sprintf("(%v)(%v)", rootName, fieldValue)
	}

	if isPrimitiveType(spec) {
		return fieldValue, nil
	}

	switch t := root.(type) {
	case *compile.BinarySpec:
		// There is no AppendBinary for ArrayEncoder, so we opt for encoding it ourselves and
		// appending it as a string. We also use AddString instead of AddBinary for ObjectEncoder
		// for consistency.
		base64 := g.Import("encoding/base64")
		return fmt.Sprintf("%v.StdEncoding.EncodeToString(%v)", base64, fieldValue), nil
	case *compile.MapSpec:
		return z.mapG.zapMarshaler(g, t, fieldValue)
	case *compile.SetSpec:
		return z.setG.zapMarshaler(g, t, fieldValue)
	case *compile.ListSpec:
		return z.listG.zapMarshaler(g, t, fieldValue)
	case *compile.StructSpec:
		return fieldValue, nil
	}
	panic(root)
}

// zapMarshalerPtr will dereference the pointer and call zapMarshal on it.
func (z *zapGenerator) zapMarshalerPtr(g Generator, spec compile.TypeSpec, fieldValue string) (string, error) {
	if isPrimitiveType(spec) {
		fieldValue = "*" + fieldValue
	}
	return z.zapMarshaler(g, spec, fieldValue)
}

// zapEncodeBegin/End handle any logging that can error and add error handling logic.
//
// Make sure that an `err` variable is declared when this is called.
func (z *zapGenerator) zapEncodeBegin(g Generator, spec compile.TypeSpec) string {
	root := compile.RootTypeSpec(spec)

	switch root.(type) {
	// Non-primitives
	case *compile.MapSpec, *compile.SetSpec, *compile.ListSpec, *compile.EnumSpec,
		*compile.StructSpec:
		return fmt.Sprintf("err = %v.Append(err, ", g.Import("go.uber.org/multierr"))
	}
	return ""
}

func (z *zapGenerator) zapEncodeEnd(spec compile.TypeSpec) string {
	root := compile.RootTypeSpec(spec)

	switch root.(type) {
	// Non-primitives
	case *compile.MapSpec, *compile.SetSpec, *compile.ListSpec, *compile.EnumSpec,
		*compile.StructSpec:
		return ")"
	}
	return ""
}

func zapOptOut(spec *compile.FieldSpec) bool {
	_, ok := spec.Annotations[NoZapLabel]
	return ok
}
