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

// setGenerator generates logic to convert lists of arbitrary Thrift types to
// and from ValueLists.
type setGenerator struct {
	// Set of ValueLists that have already been generated.
	valueLists map[string]struct{}

	// Set of readers that have already been generated.
	readers map[string]struct{}
}

func newSetGenerator() setGenerator {
	return setGenerator{
		valueLists: make(map[string]struct{}),
		readers:    make(map[string]struct{}),
	}
}

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
func (s setGenerator) ValueList(g Generator, spec *compile.SetSpec) (string, error) {
	// TODO(abg): Unhashable types
	name := "_" + valueName(spec) + "_ValueList"
	if _, ok := s.valueLists[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			type <.Name> <typeReference .Spec Required>

			<$v := newVar "v">
			<$x := newVar "x">
			<$f := newVar "f">
			func (<$v> <.Name>) ForEach(<$f> func(<$wire>.Value) error) error {
				for <$x> := range <$v> {
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
			Spec *compile.SetSpec
		}{Name: name, Spec: spec},
	)
	if err != nil {
		return "", generateError{
			Name:   typeReference(spec, Required),
			Reason: err,
		}
	}

	s.valueLists[name] = struct{}{}
	return name, nil
}

func (s setGenerator) Reader(g Generator, spec *compile.SetSpec) (string, error) {
	name := "_" + valueName(spec) + "_Read"
	if _, ok := s.readers[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			<$setType := typeReference .Spec Required>

			<$s := newVar "s">
			<$i := newVar "i">
			<$o := newVar "o">
			<$x := newVar "x">
			func <.Name>(<$s> <$wire>.Set) (<$setType>, error) {
				if <$s>.ValueType != <typeCode .Spec.ValueSpec> {
					return nil, nil
				}

				<$o> := make(<$setType>, <$s>.Size)
				err := <$s>.Items.ForEach(func(<$x> <$wire>.Value) error {
					<$i>, err := <fromWire .Spec.ValueSpec $x>
					if err != nil {
						return err
					}
					<$o>[<$i>] = struct{}{}
					return nil
				})
				<$s>.Items.Close()
				return <$o>, err
			}
		`,
		struct {
			Name string
			Spec *compile.SetSpec
		}{Name: name, Spec: spec},
	)

	if err != nil {
		return "", generateError{
			Name:   typeReference(spec, Required),
			Reason: err,
		}
	}

	s.readers[name] = struct{}{}
	return name, nil
}
