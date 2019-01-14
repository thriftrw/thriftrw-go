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

func mapItemListName(g Generator, spec *compile.MapSpec) string {
	return fmt.Sprintf("_%s_MapItemList", g.MangleType(spec))
}

// mapGenerator generates logic to convert lists of arbitrary Thrift types to
// and from MapItemLists.
type mapGenerator struct{}

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
func (m *mapGenerator) ItemList(g Generator, spec *compile.MapSpec) (string, error) {
	name := mapItemListName(g, spec)
	err := g.EnsureDeclared(
		`
			<$wire := import "go.uber.org/thriftrw/wire">
			type <.Name> <typeReference .Spec>

			<$m := newVar "m">
			<$f := newVar "f">
			<$k := newVar "k">
			<$v := newVar "v">
			<$i := newVar "i">
			<$kw := newVar "kw">
			<$vw := newVar "vw">
			func (<$m> <.Name>) ForEach(<$f> func(<$wire>.MapItem) error) error {
				<- if isHashable .Spec.KeySpec ->
					for <$k>, <$v> := range <$m> {
				<else ->
					for _, <$i> := range <$m> {
						<$k> := <$i>.Key
						<$v> := <$i>.Value
				<end>
						<- if not (isPrimitiveType .Spec.KeySpec) ->
							if <$k> == nil {
								return <import "fmt">.Errorf("invalid map key: value is nil")
							}
						<end ->
						<- if not (isPrimitiveType .Spec.ValueSpec) ->
							if <$v> == nil {
								return <import "fmt">.Errorf("invalid [%v]: value is nil", <$k>)
							}
						<end ->

						<$kw>, err := <toWire .Spec.KeySpec $k>
						if err != nil {
							return err
						}

						<$vw>, err := <toWire .Spec.ValueSpec $v>
						if err != nil {
							return err
						}
						err = <$f>(<$wire>.MapItem{Key: <$kw>, Value: <$vw>})
						if err != nil {
							return err
						}
					}
				return nil
			}

			func (<$m> <.Name>) Size() int {
				return len(<$m>)
			}

			func (<.Name>) KeyType() <$wire>.Type {
				return <typeCode .Spec.KeySpec>
			}

			func (<.Name>) ValueType() <$wire>.Type {
				return <typeCode .Spec.ValueSpec>
			}

			func (<.Name>) Close() {}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

func (m *mapGenerator) Reader(g Generator, spec *compile.MapSpec) (string, error) {
	name := readerFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$wire := import "go.uber.org/thriftrw/wire">
			<$mapType := typeReference .Spec>

			<$m := newVar "m">
			<$o := newVar "o">
			<$x := newVar "x">
			<$k := newVar "k">
			<$v := newVar "v">
			func <.Name>(<$m> <$wire>.MapItemList) (<$mapType>, error) {
				if <$m>.KeyType() != <typeCode .Spec.KeySpec> {
					return nil, nil
				}

				if <$m>.ValueType() != <typeCode .Spec.ValueSpec> {
					return nil, nil
				}

				<if isHashable .Spec.KeySpec>
					<$o> := make(<$mapType>, <$m>.Size())
				<else>
					<$o> := make(<$mapType>, 0, <$m>.Size())
				<end ->
				err := <$m>.ForEach(func(<$x> <$wire>.MapItem) error {
					<$k>, err := <fromWire .Spec.KeySpec (printf "%s.Key" $x)>
					if err != nil {
						return err
					}

					<$v>, err := <fromWire .Spec.ValueSpec (printf "%s.Value" $x)>
					if err != nil {
						return err
					}

					<if isHashable .Spec.KeySpec>
						<$o>[<$k>] = <$v>
					<else>
						<$o> = append(<$o>, struct {
							Key <typeReference .Spec.KeySpec>
							Value <typeReference .Spec.ValueSpec>
						}{<$k>, <$v>})
					<end ->
					return nil
				})
				<$m>.Close()
				return <$o>, err
			}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

// Equals generates a function to compare maps of the given type
//
// 	func $name(lhs, rhs $mapType) bool {
// 		...
// 	}
//
// And returns its name.
func (m *mapGenerator) Equals(g Generator, spec *compile.MapSpec) (string, error) {
	if !isHashable(spec.KeySpec) {
		return m.equalsUnhashable(g, spec)
	}

	name := equalsFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$mapType := typeReference .Spec>

			<$lhs := newVar "lhs">
			<$rhs := newVar "rhs">
			func <.Name>(<$lhs>, <$rhs> <$mapType>) bool {
				if len(<$lhs>) != len(<$rhs>) {
					return false
				}

				<$lk := newVar "lk">
				<$lv := newVar "lv">
				<$rv := newVar "rv">
				<$ok := newVar "ok">
				for <$lk>, <$lv> := range <$lhs> {
					<$rv>, <$ok> := <$rhs>[<$lk>]
					if !<$ok> {
						return false
					}
					if !<equals .Spec.ValueSpec $lv $rv> {
						return false
					}
				}
				return true
			}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

func (m *mapGenerator) equalsUnhashable(g Generator, spec *compile.MapSpec) (string, error) {
	name := equalsFuncName(g, spec)
	err := g.EnsureDeclared(
		`
			<$mapType := typeReference .Spec>

			<$lhs := newVar "lhs">
			<$rhs := newVar "rhs">
			func <.Name>(<$lhs>, <$rhs> <$mapType>) bool {
				if len(<$lhs>) != len(<$rhs>) {
					return false
				}

				<$i := newVar "i">
				<$j := newVar "j">
				<$lk := newVar "lk">
				<$lv := newVar "lv">
				<$rk := newVar "rk">
				<$rv := newVar "rv">
				<$ok := newVar "ok">
				for _, <$i> := range <$lhs> {
					<$lk> := <$i>.Key
					<$lv> := <$i>.Value
					<$ok> := false
					for _, <$j> := range <$rhs> {
						<$rk> := <$j>.Key
						<$rv> := <$j>.Value
						if !<equals .Spec.KeySpec $lk $rk> {
							continue
						}

						if !<equals .Spec.ValueSpec $lv $rv> {
							return false
						}
						<$ok> = true
						break
					}

					if !<$ok> {
						return false
					}
				}
				return true
			}
		`,
		struct {
			Name string
			Spec *compile.MapSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

// Maps are logged as objects if the key is a string or a typedef of a
// string. If the key is not a string, maps are logged as arrays of
// objects with a key and value.
//
//   map[string]int32{"foo": 1, "bar": 2}
//   => {"foo": 1, "bar": 2}
//
//   map[int32]string{1: "foo", 2: "bar"}
//   => [{"key": 1, "value": "foo"}, {"key": 2, "value": "bar"}]
//
func (m *mapGenerator) zapMarshaler(
	g Generator,
	root *compile.MapSpec,
	fieldValue string,
) (string, error) {
	name := zapperName(g, root)
	switch compile.RootTypeSpec(root.KeySpec).(type) {
	case *compile.StringSpec:
		return m.zapStringKeyMarshaler(g, name, root, fieldValue)
	default:
		return m.zapNonstringKeyMarshaler(g, name, root, fieldValue)
	}
}

func (m *mapGenerator) zapStringKeyMarshaler(
	g Generator,
	name string,
	root *compile.MapSpec,
	fieldValue string,
) (string, error) {
	err := g.EnsureDeclared(
		`
			<$zapcore := import "go.uber.org/zap/zapcore">

			type <.Name> <typeReference .Type>
			<$m := newVar "m">
			<$k := newVar "k">
			<$v := newVar "v">
			<$enc := newVar "enc">
			// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
			// fast logging of <.Name>.
			func (<$m> <.Name>) MarshalLogObject(<$enc> <$zapcore>.ObjectEncoder) (err error) {
				for <$k>, <$v> := range <$m> {
					<zapEncodeBegin .Type.ValueSpec ->
						<$enc>.Add<zapEncoder .Type.ValueSpec>((string)(<$k>), <zapMarshaler .Type.ValueSpec $v>)
					<- zapEncodeEnd .Type.ValueSpec>
				}
				return err
			}
			`, struct {
			Name string
			Type *compile.MapSpec
		}{
			Name: name,
			Type: root,
		},
	)
	return fmt.Sprintf("(%v)(%v)", name, fieldValue), err
}

func (m *mapGenerator) zapNonstringKeyMarshaler(
	g Generator,
	name string,
	root *compile.MapSpec,
	fieldValue string,
) (string, error) {
	if err := g.EnsureDeclared(
		`
			<$zapcore := import "go.uber.org/zap/zapcore">
			<$multierr := import "go.uber.org/multierr">

			type <.Name> <typeReference .Type>
			<$m := newVar "m">
			<$k := newVar "k">
			<$v := newVar "v">
			<$i := newVar "i">
			<$enc := newVar "enc">
			// MarshalLogArray implements zapcore.ArrayMarshaler, enabling
			// fast logging of <.Name>.
			func (<$m> <.Name>) MarshalLogArray(<$enc> <$zapcore>.ArrayEncoder) (err error) {
				<- if isHashable .Type.KeySpec ->
					for <$k>, <$v> := range <$m> {
				<else ->
					for _, <$i> := range <$m> {
						<$k> := <$i>.Key
						<$v> := <$i>.Value
				<end ->
					err = <$multierr>.Append(err, <$enc>.AppendObject(<zapMapItemMarshaler .Type $k $v>))
				}
				return err
			}
			`, struct {
			Name string
			Type *compile.MapSpec
		}{
			Name: name,
			Type: root,
		},
		TemplateFunc("zapMapItemMarshaler", m.zapMapItemMarshaler),
	); err != nil {
		return "", err
	}
	return fmt.Sprintf("(%v)(%v)", name, fieldValue), nil
}

func (m *mapGenerator) zapMapItemMarshaler(
	g Generator,
	mapSpec *compile.MapSpec,
	keyVar string,
	valueVar string,
) (string, error) {
	name := fmt.Sprintf("_%s_Item_Zapper", g.MangleType(mapSpec))
	if err := g.EnsureDeclared(
		`
			<$zapcore := import "go.uber.org/zap/zapcore">

			type <.Name> struct {
				Key   <typeReference .KeyType>
				Value <typeReference .ValueType>
			}
			<$v := newVar "v">
			<$key := printf "%s.%s" $v "Key">
			<$val := printf "%s.%s" $v "Value">
			<$enc := newVar "enc">
			// MarshalLogArray implements zapcore.ArrayMarshaler, enabling
			// fast logging of <.Name>.
			func (<$v> <.Name>) MarshalLogObject(<$enc> <$zapcore>.ObjectEncoder) (err error) {
				<zapEncodeBegin .KeyType ->
					<$enc>.Add<zapEncoder .KeyType>("key", <zapMarshaler .KeyType $key>)
				<- zapEncodeEnd .KeyType>
				<zapEncodeBegin .ValueType ->
					<$enc>.Add<zapEncoder .ValueType>("value", <zapMarshaler .ValueType $val>)
				<- zapEncodeEnd .ValueType>
				return err
			}
			`, struct {
			Name      string
			KeyType   compile.TypeSpec
			ValueType compile.TypeSpec
		}{
			Name:      name,
			KeyType:   mapSpec.KeySpec,
			ValueType: mapSpec.ValueSpec,
		},
	); err != nil {
		return "", err
	}
	return fmt.Sprintf("%v{Key: %v, Value: %v}", name, keyVar, valueVar), nil
}
