// Copyright (c) 2021 Uber Technologies, Inc.
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

// StreamGenerator generates code that knows how to encode and decode Thrift
// objects in a streaming fashion.
type StreamGenerator struct {
	mapG  mapGenerator
	setG  setGenerator
	listG listGenerator

	enumG    enumGenerator
	structG  structGenerator
	typedefG typedefGenerator
}

// Encode generates code that knows how to serialize Thrift types into bytes.
func (sg *StreamGenerator) Encode(g Generator, spec compile.TypeSpec, varName string, sw string) (string, error) {
	switch s := spec.(type) {
	case *compile.BoolSpec:
		return fmt.Sprintf("%s.WriteBool(%s)", sw, varName), nil
	case *compile.I8Spec:
		return fmt.Sprintf("%s.WriteInt8(%s)", sw, varName), nil
	case *compile.I16Spec:
		return fmt.Sprintf("%s.WriteInt16(%s)", sw, varName), nil
	case *compile.I32Spec:
		return fmt.Sprintf("%s.WriteInt32(%s)", sw, varName), nil
	case *compile.I64Spec:
		return fmt.Sprintf("%s.WriteInt64(%s)", sw, varName), nil
	case *compile.DoubleSpec:
		return fmt.Sprintf("%s.WriteDouble(%s)", sw, varName), nil
	case *compile.StringSpec:
		return fmt.Sprintf("%s.WriteString(%s)", sw, varName), nil
	case *compile.BinarySpec:
		return fmt.Sprintf("%s.WriteBinary(%s)", sw, varName), nil
	case *compile.MapSpec:
		encoder, err := sg.mapG.Encoder(g, s)
		return fmt.Sprintf("%s(%s, %s)", encoder, varName, sw), err
	case *compile.ListSpec:
		encoder, err := sg.listG.Encoder(g, s)
		return fmt.Sprintf("%s(%s, %s)", encoder, varName, sw), err
	case *compile.SetSpec:
		encoder, err := sg.setG.Encoder(g, s)
		return fmt.Sprintf("%s(%s, %s)", encoder, varName, sw), err
	default:
		return fmt.Sprintf("%s.Encode(%s)", varName, sw), nil
	}
}

// EncodePtr is the same as Encode except varName is expected to be a reference
// to a value of the given type.
func (sg *StreamGenerator) EncodePtr(g Generator, spec compile.TypeSpec, varName string, sw string) (string, error) {
	switch spec.(type) {
	case *compile.BoolSpec, *compile.I8Spec, *compile.I16Spec, *compile.I32Spec,
		*compile.I64Spec, *compile.DoubleSpec, *compile.StringSpec:
		return sg.Encode(g, spec, fmt.Sprintf("*(%s)", varName), sw)
	default:
		// Everything else is either a reference type or has an Encode method
		// on it that does automatic dereferencing.
		return sg.Encode(g, spec, varName, sw)
	}
}

// Decode generates an expression that can deserialize Thrift data into their
// raw types.
func (sg *StreamGenerator) Decode(g Generator, spec compile.TypeSpec, reader string) (string, error) {
	switch s := spec.(type) {
	case *compile.BoolSpec:
		return fmt.Sprintf("%s.ReadBool()", reader), nil
	case *compile.I8Spec:
		return fmt.Sprintf("%s.ReadInt8()", reader), nil
	case *compile.I16Spec:
		return fmt.Sprintf("%s.ReadInt16()", reader), nil
	case *compile.I32Spec:
		return fmt.Sprintf("%s.ReadInt32()", reader), nil
	case *compile.I64Spec:
		return fmt.Sprintf("%s.ReadInt64()", reader), nil
	case *compile.DoubleSpec:
		return fmt.Sprintf("%s.ReadDouble()", reader), nil
	case *compile.StringSpec:
		return fmt.Sprintf("%s.ReadString()", reader), nil
	case *compile.BinarySpec:
		return fmt.Sprintf("%s.ReadBinary()", reader), nil
	case *compile.MapSpec:
		decoder, err := sg.mapG.Decoder(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", decoder, reader), nil
	case *compile.ListSpec:
		decoder, err := sg.listG.Decoder(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", decoder, reader), nil
	case *compile.SetSpec:
		decoder, err := sg.setG.Decoder(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", decoder, reader), nil
	case *compile.TypedefSpec:
		decoder, err := sg.typedefG.Decoder(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", decoder, reader), nil
	case *compile.EnumSpec:
		decoder, err := sg.enumG.Decoder(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", decoder, reader), nil
	case *compile.StructSpec:
		decoder, err := sg.structG.Decoder(g, s)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s(%s)", decoder, reader), nil
	default:
		panic(fmt.Sprintf("Unknown TypeSpec (%T) %v", spec, spec))
	}
}

// DecodePtr generates an expression that assigns the "lhs" to a pointer of the
// decoded value.
func (sg *StreamGenerator) DecodePtr(g Generator, spec compile.TypeSpec, lhs string, reader string) (string, error) {
	if !isPrimitiveType(spec) {
		out, err := sg.Decode(g, spec, reader)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s, err = %s", lhs, out), err
	}
	return g.TextTemplate(
		`
		<- $x := newVar "x" ->
		var <$x> <typeReference .Spec>
		<$x>, err = <decode .Spec .Reader>
		<.LHS> = &<$x ->
		`,
		struct {
			Spec   compile.TypeSpec
			LHS    string
			Reader string
		}{Spec: spec, LHS: lhs, Reader: reader},
	)
}
