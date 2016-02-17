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

import "github.com/thriftrw/thriftrw-go/compile"

func (g *Generator) enum(spec *compile.EnumSpec) error {
	// TODO(abg) define an error type in the library for unrecognized enums.
	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$enumName := defName .Spec>
		type <$enumName> int32

		const (
		<range .Spec.Items>
			<$enumName><goCase .Name> <$enumName> = <.Value>
		<end>
		)

		<$v := newVar "v">
		func (<$v> <$enumName>) ToWire() <$wire>.Value {
			return <$wire>.NewValueI32(int32(<$v>))
		}

		<$w := newVar "w">
		func (<$v> *<$enumName>) FromWire(<$w> <$wire>.Value) error {
			*<$v> = (<$enumName>)(<$w>.GetI32());
			return nil
		}

		func <.Reader>(<$w> <$wire>.Value) (<$enumName>, error) {
			var <$v> <$enumName>
			err := <$v>.FromWire(<$w>)
			return <$v>, err
		}
		`,
		struct {
			Spec   *compile.EnumSpec
			Reader string
		}{Spec: spec, Reader: typeReader(spec)},
	)

	return wrapGenerateError(spec.Name, err)
}
