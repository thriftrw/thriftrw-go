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

import "github.com/uber/thriftrw-go/compile"

func (g *Generator) structure(spec *compile.StructSpec) error {
	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/uber/thriftrw-go/wire">
		<$structName := defName .>

		type <$structName> struct {
		<range .Fields>
			<goCase .Name> <typeReference .Type (required .Required)>
		<end>
		}


		<$v := newVar "v">
		func (<$v> *<$structName>) ToWire() <$wire>.Value {
			return <$wire>.NewValueStruct(
				<$wire>.Struct{
					[]<$wire>.Field{
					<range .Fields>
						// TODO handle optional fields and nil values
						<$f := printf "%s.%s" $v (goCase .Name)>
						{ID: <.ID>, Value: <toWire .Type $f>},
					<end>
					},
				},
			)
		}

		<$w := newVar "w">
		func (<$v> *<$structName>) FromWire(<$w> <$wire>.Value) error {
			<$f := newVar "f">
			for _, <$f> := range <$w>.GetStruct().Fields {
				switch <$f>.ID {
				<range .Fields>
				case <.ID>:
					if <$f>.Value.Type() == nil { // TODO
						<$t := printf "%s.%s" $v (goCase .Name)>
						<fromWire .Type $t (printf "%s.Value" $f)>
						// TODO read errors
					}
				<end>
				}
			}
		}
		`,
		spec,
	)
	// TODO(abg): JSON tags for generated structs
	// TODO(abg): ToWire/FromWire for all fields

	return wrapGenerateError(spec.Name, err)
	// TODO methods
}
