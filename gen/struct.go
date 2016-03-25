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

package gen

import (
	"github.com/thriftrw/thriftrw-go/ast"
	"github.com/thriftrw/thriftrw-go/compile"
)

// structGenerator generates code to serialize and deserialize structs.
type structGenerator struct {
	hasReaders
}

func (s *structGenerator) Reader(g Generator, spec *compile.StructSpec) (string, error) {
	name := "_" + goCase(spec.ThriftName()) + "_Read"
	if s.HasReader(name) {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$v := newVar "v">
		<$w := newVar "w">
		func <.Name>(<$w> <$wire>.Value) (<typeReference .Spec>, error) {
			var <$v> <typeName .Spec>
			err := <$v>.FromWire(<$w>)
			return &<$v>, err
		}
		`,
		struct {
			Name string
			Spec *compile.StructSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

func structure(g Generator, spec *compile.StructSpec) error {
	fg := fieldGroupGenerator{
		Name:   goCase(spec.Name),
		Fields: spec.Fields,
	}

	if err := fg.DefineStruct(g); err != nil {
		return wrapGenerateError(spec.ThriftName(), err)
	}

	if err := fg.ToWire(g); err != nil {
		return wrapGenerateError(spec.ThriftName(), err)
	}

	if err := fg.FromWire(g); err != nil {
		return wrapGenerateError(spec.ThriftName(), err)
	}

	if err := fg.String(g); err != nil {
		return wrapGenerateError(spec.ThriftName(), err)
	}

	if spec.Type == ast.ExceptionType {
		err := g.DeclareFromTemplate(
			`
			<$v := newVar "v">
			func (<$v> *<typeName .>) Error() string {
				return <$v>.String()
			}
			`, spec)
		if err != nil {
			return wrapGenerateError(spec.ThriftName(), err)
		}
	}

	return nil
	// TODO(abg): JSON tags for generated structs
	// TODO(abg): For all struct types, handle the case where fields are named
	// ToWire or FromWire.
	// TODO(abg): For exceptions, handle the case where a field is named
	// Error.
}
