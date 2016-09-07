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

import (
	"strings"

	"github.com/thriftrw/thriftrw-go/compile"
)

// enumGenerator generates code to serialize and deserialize enums.
type enumGenerator struct {
	hasReaders
}

func (e *enumGenerator) Reader(g Generator, spec *compile.EnumSpec) (string, error) {
	name := "_" + goCase(spec.ThriftName()) + "_Read"
	if e.HasReader(name) {
		return name, nil
	}

	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$v := newVar "v">
		<$w := newVar "w">
		func <.Name>(<$w> <$wire>.Value) (<typeName .Spec>, error) {
			var <$v> <typeName .Spec>
			err := <$v>.FromWire(<$w>)
			return <$v>, err
		}
		`,
		struct {
			Name string
			Spec *compile.EnumSpec
		}{Name: name, Spec: spec},
	)

	return name, wrapGenerateError(spec.ThriftName(), err)
}

func enum(g Generator, spec *compile.EnumSpec) error {
	items := enumUniqueItems(spec.Items)

	// TODO(abg) define an error type in the library for unrecognized enums.
	err := g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$enumName := typeName .Spec>
		type <$enumName> int32

		<if .Spec.Items>
			const (
			<$enum := .Spec>
			<range .Spec.Items>
				<enumItemName $enum .Name> <$enumName> = <.Value>
			<end>
			)
		<end>

		<$v := newVar "v">
		func (<$v> <$enumName>) ToWire() (<$wire>.Value, error) {
			return <$wire>.NewValueI32(int32(<$v>)), nil
		}

		<$w := newVar "w">
		func (<$v> *<$enumName>) FromWire(<$w> <$wire>.Value) error {
			*<$v> = (<$enumName>)(<$w>.GetI32());
			return nil
		}

		func (<$v> <$enumName>) String() string {
			<$w> := int32(<$v>)
			<if len .Spec.Items>
				switch <$w> {
				<range .UniqueItems>
					case <.Value>:
						return "<.Name>"
				<end>
				}
			<end>
			return fmt.Sprintf("<$enumName>(%d)", <$w>)
		}
		`,
		struct {
			Spec        *compile.EnumSpec
			UniqueItems []compile.EnumItem
		}{
			Spec:        spec,
			UniqueItems: items,
		},
		TemplateFunc("enumItemName", enumItemName),
	)

	return wrapGenerateError(spec.Name, err)
}

// enumItemName returns the Go name that should be used for an enum item with
// the given Thrift name.
func enumItemName(g Generator, spec compile.TypeSpec, itemName string) (string, error) {
	enumName, err := typeName(g, spec)
	if err != nil {
		return "", err
	}
	return enumName + pascalCase(false /* all caps */, strings.Split(itemName, "_")...), nil
}

// enumUniqueItems returns a subset of the given list of enum items where
// there are no value collisions between items.
func enumUniqueItems(items []compile.EnumItem) []compile.EnumItem {
	used := make(map[int32]struct{}, len(items))
	filtered := items[:0] // zero-alloc filtering
	for _, i := range items {
		if _, isUsed := used[i.Value]; isUsed {
			continue
		}
		filtered = append(filtered, i)
		used[i.Value] = struct{}{}
	}
	return filtered
}
