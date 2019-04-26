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

// setGenerator generates logic to convert lists of arbitrary Thrift types to
// and from ValueLists.
type setGenerator struct{}

// ValueList generates a new ValueList type alias for the given set.
//
// The following is generated:
//
// 	type $valueListName map[$valueType]struct{}
//
// 	func (v $valueListName) ForEach(f func(wire.Value) error) error { ... }
//
// 	func (v $valueListName) Close() { ... }
//
// And $valueListName is returned. This may be used where a ValueList of the
// given type is expected.
func (s *setGenerator) ValueList(g Generator, spec *compile.SetSpec) (string, error) {
	name := valueListName(g, spec)
	err := g.EnsureDeclared(
		`
			<$wire := import "go.uber.org/thriftrw/wire">
			type <.Name> <typeReference .Spec>

			<$v := newVar "v">
			<$x := newVar "x">
			<$f := newVar "f">
			<$w := newVar "w">
			func (<$v> <.Name>) ForEach(<$f> func(<$wire>.Value) error) error {
				<- if setUsesMap .Spec ->
					for <$x> := range <$v> {
				<- else ->
					for _, <$x> := range <$v> {
				<- end ->
						<if not (isPrimitiveType .Spec.ValueSpec)>
							if <$x> == nil {
								return <import "fmt">.Errorf("invalid set item: value is nil")
							}
						<end ->

						<$w>, err := <toWire .Spec.ValueSpec $x>
						if err != nil {
							return err
						}

						if err := <$f>(<$w>); err != nil {
							return err
						}
					}
				return nil
			}

			func (<$v> <.Name>) Size() int {
				return len(<$v>)
			}

			func (<.Name>) ValueType() <$wire>.Type {
				return <typeCode .Spec.ValueSpec>
			}

			func (<.Name>) Close() {}
		`,
		struct {
			Name string
			Spec *compile.SetSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

func (s *setGenerator) Reader(g Generator, spec *compile.SetSpec) (string, error) {
	name := readerFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$wire := import "go.uber.org/thriftrw/wire">
			<$setType := typeReference .Spec>

			<$s := newVar "s">
			<$i := newVar "i">
			<$o := newVar "o">
			<$x := newVar "x">
			func <.Name>(<$s> <$wire>.ValueList) (<$setType>, error) {
				if <$s>.ValueType() != <typeCode .Spec.ValueSpec> {
					return nil, nil
				}

				<if setUsesMap .Spec>
					<$o> := make(<$setType>, <$s>.Size())
				<else>
					<$o> := make(<$setType>, 0, <$s>.Size())
				<end ->
				err := <$s>.ForEach(func(<$x> <$wire>.Value) error {
					<$i>, err := <fromWire .Spec.ValueSpec $x>
					if err != nil {
						return err
					}
					<if setUsesMap .Spec>
						<$o>[<$i>] = struct{}{}
					<else>
						<$o> = append(<$o>, <$i>)
					<end ->
					return nil
				})
				<$s>.Close()
				return <$o>, err
			}
		`,
		struct {
			Name string
			Spec *compile.SetSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

// Equals generates a function to compare sets of the given type
//
// func $name(lhs, rhs $setType) bool {
//      ...
// }
//
// And returns its name.
func (s *setGenerator) Equals(g Generator, spec *compile.SetSpec) (string, error) {
	name := equalsFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$setType := typeReference .Spec>

			<$lhs := newVar "lhs">
			<$rhs := newVar "rhs">
			func <.Name>(<$lhs>, <$rhs> <$setType>) bool {
				if len(<$lhs>) != len(<$rhs>) {
					return false
				}

				// if the values in the set are hashable they can be used
				// as keys in a map.
				<$o := newVar "o">
				<$x := newVar "x">
				<$y := newVar "y">
				<$ok := newVar "ok">
				<if setUsesMap .Spec>
					for <$x> := range <$rhs> {
						if _, <$ok> := <$lhs>[<$x>]; !<$ok> {
							return false
						}
					}
				<else>
					// Note if values are not hashable then this is O(n^2) in time complexity.
					for _, <$x> := range <$lhs> {
						<$ok> := false
						for _, <$y> := range <$rhs> {
							if <equals .Spec.ValueSpec $x $y> {
								<$ok> = true
								break
							}
						}
						if !<$ok> {
							return false
						}
					}
				<end>

				return true
			}
		`,
		struct {
			Name string
			Spec *compile.SetSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

func (s *setGenerator) zapMarshaler(
	g Generator,
	root *compile.SetSpec,
	fieldValue string,
) (string, error) {
	name := zapperName(g, root)
	if err := g.EnsureDeclared(
		`
			<$zapcore := import "go.uber.org/zap/zapcore">

			type <.Name> <typeReference .Type>
			<$s := newVar "s">
			<$v := newVar "v">
			<$enc := newVar "enc">
			// MarshalLogArray implements zapcore.ArrayMarshaler, enabling
			// fast logging of <.Name>.
			func (<$s> <.Name>) MarshalLogArray(<$enc> <$zapcore>.ArrayEncoder) (err error) {
				<- if setUsesMap .Type ->
					for <$v> := range <$s> {
				<else ->
					for _, <$v> := range <$s> {
				<end ->
					<zapEncodeBegin .Type.ValueSpec ->
						<$enc>.Append<zapEncoder .Type.ValueSpec>(<zapMarshaler .Type.ValueSpec $v>)
					<- zapEncodeEnd .Type.ValueSpec>
				}
				return err
			}
			`, struct {
			Name string
			Type *compile.SetSpec
		}{
			Name: name,
			Type: root,
		},
	); err != nil {
		return "", err
	}
	return fmt.Sprintf("(%v)(%v)", name, fieldValue), nil
}
