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
	"fmt"

	"github.com/thriftrw/thriftrw-go/compile"
)

// mapGenerator generates logic to convert lists of arbitrary Thrift types to
// and from MapItemLists.
type mapGenerator struct {
	// Set of MapItemLIsts that have already been generated.
	itemLists map[string]struct{}

	// Set of readers that have already been generated.
	readers map[string]struct{}
}

func newMapGenerator() mapGenerator {
	return mapGenerator{
		itemLists: make(map[string]struct{}),
		readers:   make(map[string]struct{}),
	}
}

// MapItemList generates a new MapItemList type alias for the given map.
//
// The following is generated:
//
// 	type $mapItemListName map[$keyType]$valueType
//
// 	func (v $mapItemListName) ForEach(f func(wire.MapItem) error) error { ... }
//
// 	func (v $mapItemListName) Close() { ... }
//
// And $mapItemListName is returned. This may be used where a MapItemList of the
// given type is expected.
func (m mapGenerator) ItemList(g Generator, spec *compile.MapSpec) (string, error) {
	// TODO(abg): Unhashable types
	name := "_" + valueName(spec) + "_MapItemList"
	if _, ok := m.itemLists[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			type <.Name> <typeReference .Spec>

			<$m := newVar "m">
			<$f := newVar "f">
			<$k := newVar "k">
			<$v := newVar "v">
			func (<$m> <.Name>) ForEach(<$f> func(<$wire>.MapItem) error) error {
				for <$k>, <$v> := range <$m> {
					err := <$f>(<$wire>.MapItem{
						Key: <toWire .Spec.KeySpec $k>,
						Value: <toWire .Spec.ValueSpec $v>,
					})
					if err != nil {
						return err
					}
				}
				return nil
			}

			func (<$m> <.Name>) Close() {}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)
	if err != nil {
		return "", generateError{
			Name:   typeReference(spec),
			Reason: err,
		}
	}

	m.itemLists[name] = struct{}{}
	return name, nil
}

func (m mapGenerator) Reader(g Generator, spec *compile.MapSpec) (string, error) {
	name := "_" + valueName(spec) + "_Read"
	if _, ok := m.readers[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/thriftrw/thriftrw-go/wire">
			<$mapType := typeReference .Spec>

			<$m := newVar "m">
			<$o := newVar "o">
			<$x := newVar "x">
			<$k := newVar "k">
			<$v := newVar "v">
			func <.Name>(<$m> <$wire>.Map) (<$mapType>, error) {
				if <$m>.KeyType != <typeCode .Spec.KeySpec> {
					return nil, nil
				}

				if <$m>.ValueType != <typeCode .Spec.ValueSpec> {
					return nil, nil
				}

				<$o> := make(<$mapType>, <$m>.Size)
				err := <$m>.Items.ForEach(func(<$x> <$wire>.MapItem) error {
					<$k>, err := <fromWire .Spec.KeySpec (printf "%s.Key" $x)>
					if err != nil {
						return err
					}

					<$v>, err := <fromWire .Spec.ValueSpec (printf "%s.Value" $x)>
					if err != nil {
						return err
					}

					<$o>[<$k>] = <$v>
					return nil
				})
				<$m>.Items.Close()
				return <$o>, err
			}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)

	if err != nil {
		return "", generateError{
			Name:   typeReference(spec),
			Reason: err,
		}
	}

	m.readers[name] = struct{}{}
	return name, nil
}

func valueName(spec compile.TypeSpec) string {
	switch s := spec.(type) {
	case *compile.MapSpec:
		return fmt.Sprintf(
			"Map_%s_%s", valueName(s.KeySpec), valueName(s.ValueSpec),
		)
	case *compile.ListSpec:
		return fmt.Sprintf("List_%s", valueName(s.ValueSpec))
	case *compile.SetSpec:
		return fmt.Sprintf("Set_%s", valueName(s.ValueSpec))
	default:
		return goCase(spec.ThriftName())
	}
}
