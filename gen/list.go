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

// listGenerator generates logic to convert lists of arbitrary Thrift types to
// and from ValueLists.
type listGenerator struct{}

// ValueList generates a new ValueList type alias for the given list.
//
// The following is generated:
//
// 	type $valueListName []$valueType
//
// 	func (v $valueListName) ForEach(f func(wire.Value) error) error { ... }
//
// 	func (v $valueListName) Close() { ... }
//
// And $valueListName is returned. This may be used where a ValueList of the
// given type is expected.
func (l *listGenerator) ValueList(g Generator, spec *compile.ListSpec) (string, error) {
	name := valueListName(g, spec)
	err := g.EnsureDeclared(
		`
			<$wire := import "go.uber.org/thriftrw/wire">
			type <.Name> <typeReference .Spec>

			<$i := newVar "i">
			<$v := newVar "v">
			<$x := newVar "x">
			<$f := newVar "f">
			<$w := newVar "w">
			func (<$v> <.Name>) ForEach(<$f> func(<$wire>.Value) error) error {
				<if isPrimitiveType .Spec.ValueSpec ->
				for _, <$x> := range <$v> {
				<- else ->
				for <$i>, <$x> := range <$v> {
					if <$x> == nil {
						return <import "fmt">.Errorf("invalid [%v]: value is nil", <$i>)
					}
				<- end>
					<$w>, err := <toWire .Spec.ValueSpec $x>
					if err != nil {
						return err
					}
					err = <$f>(<$w>)
					if err != nil {
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
			Spec *compile.ListSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

// Reader generates a function to read a list of the given type from a
// wire.List.
//
// 	func $name(l wire.List) ($listType, error) {
// 		...
// 	}
//
// And returns its name.
func (l *listGenerator) Reader(g Generator, spec *compile.ListSpec) (string, error) {
	name := readerFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$wire := import "go.uber.org/thriftrw/wire">
			<$listType := typeReference .Spec>

			<$l := newVar "l">
			<$i := newVar "i">
			<$o := newVar "o">
			<$x := newVar "x">
			func <.Name>(<$l> <$wire>.ValueList) (<$listType>, error) {
				if <$l>.ValueType() != <typeCode .Spec.ValueSpec> {
					return nil, nil
				}

				<$o> := make(<$listType>, 0, <$l>.Size())
				err := <$l>.ForEach(func(<$x> <$wire>.Value) error {
					<$i>, err := <fromWire .Spec.ValueSpec $x>
					if err != nil {
						return err
					}
					<$o> = append(<$o>, <$i>)
					return nil
				})
				<$l>.Close()
				return <$o>, err
			}
		`,
		struct {
			Name string
			Spec *compile.ListSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

// Equals generates a function to compare lists of the given type
//
// 	func $name(lhs, rhs $listType) bool {
// 		...
// 	}
//
// And returns its name.
func (l *listGenerator) Equals(g Generator, spec *compile.ListSpec) (string, error) {
	name := equalsFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$listType := typeReference .Spec>

			<$lhs := newVar "lhs">
			<$rhs := newVar "rhs">
			func <.Name>(<$lhs>, <$rhs> <$listType>) bool {
				if len(<$lhs>) != len(<$rhs>) {
					return false
				}

				<$i := newVar "i">
				<$lv := newVar "lv">
				<$rv := newVar "rv">
				for <$i>, <$lv> := range <$lhs> {
					<$rv> := <$rhs>[<$i>]
					if !<equals .Spec.ValueSpec $lv $rv> {
						return false
					}
				}

				return true
			}
		`,
		struct {
			Name string
			Spec *compile.ListSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

// Slices are logged as JSON arrays.
func (l *listGenerator) zapMarshaler(
	g Generator,
	spec *compile.ListSpec,
	fieldValue string,
) (string, error) {
	name := zapperName(g, spec)
	if err := g.EnsureDeclared(
		`
			<$zapcore := import "go.uber.org/zap/zapcore">

			type <.Name> <typeReference .Type>
			<$l := newVar "l">
			<$v := newVar "v">
			<$enc := newVar "enc">
			// MarshalLogArray implements zapcore.ArrayMarshaler, enabling
			// fast logging of <.Name>.
			func (<$l> <.Name>) MarshalLogArray(<$enc> <$zapcore>.ArrayEncoder) (err error) {
				for _, <$v> := range <$l> {
					<zapEncodeBegin .Type.ValueSpec ->
						<$enc>.Append<zapEncoder .Type.ValueSpec>(<zapMarshaler .Type.ValueSpec $v>)
					<- zapEncodeEnd .Type.ValueSpec>
				}
				return err
			}
			`, struct {
			Name string
			Type *compile.ListSpec
		}{
			Name: name,
			Type: spec,
		},
	); err != nil {
		return "", err
	}
	return fmt.Sprintf("(%v)(%v)", name, fieldValue), nil
}
