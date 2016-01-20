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

	"github.com/uber/thriftrw-go/compile"
)

// listValueList generates a new ValueList type alias for the given list.
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
func (g *Generator) listValueList(spec *compile.ListSpec) (string, error) {
	name := "_" + valueName(spec) + "_ValueList"
	if _, ok := g.listValueLists[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/uber/thriftrw-go/wire">
			type <.Name> <typeReference .Spec Required>

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
	if err != nil {
		return "", generateError{
			Name:   typeReference(spec, Required),
			Reason: err,
		}
	}

	g.listValueLists[name] = struct{}{}
	return name, nil
}

func (g *Generator) listReader(spec *compile.ListSpec) (string, error) {
	name := "_" + valueName(spec) + "_Read"
	if _, ok := g.listReaders[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/uber/thriftrw-go/wire">
			<$listType := typeReference .Spec Required>

			<$l := newVar "l">
			<$i := newVar "i">
			<$o := newVar "o">
			<$x := newVar "x">
			func <.Name>(<$l> <$wire>.List) <$listType> {
				if <$l>.ValueType != <typeCode .Spec.ValueSpec> {
					return nil
				}

				<$o> := make(<$listType>, 0, <$l>.Size)
				<$l>.Items.ForEach(func(<$x> <$wire>.Value) error {
					var <$i> <typeReference .Spec.ValueSpec Required>
					<fromWire .Spec.ValueSpec $i $x>
					// TODO error handling
					<$o> = append(<$o>, <$i>)
					return nil
				})
				<$l>.Items.Close()
				return <$o>
			}
		`,
		struct {
			Name string
			Spec *compile.ListSpec
		}{Name: name, Spec: spec},
	)

	if err != nil {
		return "", generateError{
			Name:   typeReference(spec, Required),
			Reason: err,
		}
	}

	g.listReaders[name] = struct{}{}
	return name, nil
}

// setValueList generates a new ValueList type alias for the given set.
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
func (g *Generator) setValueList(spec *compile.SetSpec) (string, error) {
	// TODO(abg): Unhashable types
	name := "_" + valueName(spec) + "_ValueList"
	if _, ok := g.setValueLists[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/uber/thriftrw-go/wire">
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

	g.setValueLists[name] = struct{}{}
	return name, nil
}

func (g *Generator) setReader(spec *compile.SetSpec) (string, error) {
	name := "_" + valueName(spec) + "_Read"
	if _, ok := g.setReaders[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/uber/thriftrw-go/wire">
			<$setType := typeReference .Spec Required>

			<$s := newVar "s">
			<$i := newVar "i">
			<$o := newVar "o">
			<$x := newVar "x">
			func <.Name>(<$s> <$wire>.Set) <$setType> {
				if <$s>.ValueType != <typeCode .Spec.ValueSpec> {
					return nil
				}

				<$o> := make(<$setType>, <$s>.Size)
				<$s>.Items.ForEach(func(<$x> <$wire>.Value) error {
					var <$i> <typeReference .Spec.ValueSpec Required>
					<fromWire .Spec.ValueSpec $i $x>
					// TODO error handling
					<$o>[<$i>] = struct{}{}
					return nil
				})
				<$s>.Items.Close()
				return <$o>
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

	g.setReaders[name] = struct{}{}
	return name, nil
}

// mapItemList generates a new MapItemList type alias for the given map.
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
func (g *Generator) mapItemList(spec *compile.MapSpec) (string, error) {
	// TODO(abg): Unhashable types
	name := "_" + valueName(spec) + "_MapItemList"
	if _, ok := g.mapItemLists[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/uber/thriftrw-go/wire">
			type <.Name> <typeReference .Spec Required>

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
			Name:   typeReference(spec, Required),
			Reason: err,
		}
	}

	g.mapItemLists[name] = struct{}{}
	return name, nil
}

func (g *Generator) mapReader(spec *compile.MapSpec) (string, error) {
	name := "_" + valueName(spec) + "_Read"
	if _, ok := g.mapReaders[name]; ok {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
			<$wire := import "github.com/uber/thriftrw-go/wire">
			<$mapType := typeReference .Spec Required>

			<$m := newVar "m">
			<$o := newVar "o">
			<$x := newVar "x">
			<$k := newVar "k">
			<$v := newVar "v">
			func <.Name>(<$m> <$wire>.Map) <$mapType> {
				if <$m>.KeyType != <typeCode .Spec.KeySpec> {
					return nil
				}

				if <$m>.ValueType != <typeCode .Spec.ValueSpec> {
					return nil
				}

				<$o> := make(<$mapType>, <$m>.Size)
				<$m>.Items.ForEach(func(<$x> <$wire>.MapItem) error {
					var <$k> <typeReference .Spec.KeySpec Required>
					var <$v> <typeReference .Spec.ValueSpec Required>
					<fromWire .Spec.KeySpec $k (printf "%s.Key" $x)>
					<fromWire .Spec.ValueSpec $v (printf "%s.Value" $x)>
					// TODO error handling
					<$o>[<$k>] = <$v>
					return nil
				})
				<$m>.Items.Close()
				return <$o>
			}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)

	if err != nil {
		return "", generateError{
			Name:   typeReference(spec, Required),
			Reason: err,
		}
	}

	g.mapReaders[name] = struct{}{}
	return name, nil
}

func valueName(spec compile.TypeSpec) string {
	switch spec {
	case compile.BoolSpec:
		return "Bool"
	case compile.I8Spec:
		return "I8"
	case compile.I16Spec:
		return "I16"
	case compile.I32Spec:
		return "I32"
	case compile.I64Spec:
		return "I64"
	case compile.DoubleSpec:
		return "Double"
	case compile.StringSpec:
		return "String"
	case compile.BinarySpec:
		return "Binary"
	default:
		// Not a primitive type
	}

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
		return typeDeclName(spec)
	}
}
