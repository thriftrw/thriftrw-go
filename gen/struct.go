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

func structure(g Generator, spec *compile.StructSpec) error {
	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">
		<$structName := defName .Spec>
		<$structRef := typeReference .Spec Required>

		type <$structName> struct {
		<range .Spec.Fields>
			<goCase .Name> <typeReference .Type (required .Required)>
		<end>
		}


		<$v := newVar "v">
		<$i := newVar "i">
		<$fields := newVar "fs">
		func (<$v> <$structRef>) ToWire() <$wire>.Value {
			var <$fields> [<len .Spec.Fields>]<$wire>.Field
			<$i> := 0

			<range .Spec.Fields>
				<$f := printf "%s.%s" $v (goCase .Name)>

				<if .Required>
					<$wVal := toWire .Type $f>
					<$fields>[<$i>] = <$wire>.Field{ID: <.ID>, Value: <$wVal>}
					<$i>++
				<else>
					if <$f> != nil {
						<if or (isReferenceType .Type) (isStructType .Type)>
							<$fields>[<$i>] = <$wire>.Field{
								ID: <.ID>,
								Value: <toWire .Type $f>,
							}
						<else>
							<$fields>[<$i>] = <$wire>.Field{
								ID: <.ID>,
								Value: <toWire .Type (printf "*%s" $f)>,
							}
						<end>
						<$i>++
					}
				<end>
			<end>

			return <$wire>.NewValueStruct(
				<$wire>.Struct{Fields: <$fields>[:<$i>]},
			)
		}

		<$w := newVar "w">
		func (<$v> <$structRef>) FromWire(<$w> <$wire>.Value) error {
			var err error
			<$f := newVar "f">
			for _, <$f> := range <$w>.GetStruct().Fields {
				switch <$f>.ID {
				<range .Spec.Fields>
				case <.ID>:
					if <$f>.Value.Type() == <typeCode .Type> {
						<$valueErr := fromWire .Type (printf "%s.Value" $f)>
						<if or .Required (or (isReferenceType .Type) (isStructType .Type))>
							<$v>.<goCase .Name>, err = <$valueErr>
						<else>
							*<$v>.<goCase .Name>, err = <$valueErr>
						<end>

						if err != nil {
							return err
						}
					}
				<end>
				}
			}

			// TODO(abg): Check that all required fields were set.
			return nil
		}

		func <.Reader>(<$w> <$wire>.Value) (<$structRef>, error) {
			var <$v> <$structName>
			err := <$v>.FromWire(<$w>)
			return &<$v>, err
		}

		// TODO(abg): Generate Error() if exception.
		`,
		struct {
			Spec   *compile.StructSpec
			Reader string
		}{Spec: spec, Reader: typeReader(spec)},
	)
	// TODO(abg): JSON tags for generated structs

	return wrapGenerateError(spec.Name, err)
}
