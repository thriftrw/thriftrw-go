// Copyright (c) 2016 Uber Technologies, Inc.
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

// listGenerator generates logic to convert lists of arbitrary Thrift types to
// and from ValueLists.
type listGenerator struct {
	hasReaders
	hasLazyLists
}

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
	name := "_" + valueName(spec) + "_ValueList"
	if l.HasLazyList(name) {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			type <.Name> <typeReference .Spec>

			<$v := newVar "v">
			<$x := newVar "x">
			<$f := newVar "f">
			func (<$v> <.Name>) ForEach(<$f> func(<$wire>.Value) error) error {
				for _, <$x> := range <$v> {
					err := <$f>(<toWire .Spec.ValueSpec $x>)
					if err != nil {
						return err
					}
				}
				return nil
			}

			func (<$v> <.Name>) Close() {}
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
	name := "_" + valueName(spec) + "_Read"
	if l.HasReader(name) {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			<$listType := typeReference .Spec>

			<$l := newVar "l">
			<$i := newVar "i">
			<$o := newVar "o">
			<$x := newVar "x">
			func <.Name>(<$l> <$wire>.List) (<$listType>, error) {
				if <$l>.ValueType != <typeCode .Spec.ValueSpec> {
					return nil, nil
				}

				<$o> := make(<$listType>, 0, <$l>.Size)
				err := <$l>.Items.ForEach(func(<$x> <$wire>.Value) error {
					<$i>, err := <fromWire .Spec.ValueSpec $x>
					if err != nil {
						return err
					}
					<$o> = append(<$o>, <$i>)
					return nil
				})
				<$l>.Items.Close()
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
