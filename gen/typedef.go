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

// typedef generates code for the given typedef.
func typedef(g Generator, spec *compile.TypedefSpec) error {
	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">
		<$typedefType := typeReference .Spec Required>

		type <defName .Spec> <typeName .Spec.Target>

		<$v := newVar "v">
		<$x := newVar "x">
		func (<$v> <$typedefType>) ToWire() <$wire>.Value {
			<$x> := (<typeReference .Spec.Target Required>)(<$v>)
			return <toWire .Spec.Target $x>
		}

		<$w := newVar "w">
		<if isStructType .Spec>
			func (<$v> <$typedefType>) FromWire(<$w> <$wire>.Value) error {
				return (<typeReference .Spec.Target Required>)(<$v>).FromWire(<$w>)
			}
		<else>
			func (<$v> *<$typedefType>) FromWire(<$w> <$wire>.Value) error {
				<$x>, err := <fromWire .Spec.Target $w>
				*<$v> = (<$typedefType>)(<$x>)
				return err
			}
		<end>

		func <.Reader>(<$w> <$wire>.Value) (<$typedefType>, error) {
			<if isStructType .Spec>
				<$x>, err := <fromWire .Spec.Target $w>
				return (<$typedefType>)(<$x>), err
			<else>
				var <$x> <$typedefType>
				err := <$x>.FromWire(<$w>)
				return <$x>, err
			<end>
		}
		`,
		struct {
			Spec   *compile.TypedefSpec
			Reader string
		}{Spec: spec, Reader: typeReader(spec)},
	)
	// TODO(abg): To/FromWire.
	return wrapGenerateError(spec.Name, err)
}
