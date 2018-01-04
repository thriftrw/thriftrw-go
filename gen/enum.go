// Copyright (c) 2018 Uber Technologies, Inc.
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

	"go.uber.org/thriftrw/compile"
)

// enumGenerator generates code to serialize and deserialize enums.
type enumGenerator struct{}

func (e *enumGenerator) Reader(g Generator, spec *compile.EnumSpec) (string, error) {
	name := readerFuncName(g, spec)
	err := g.EnsureDeclared(
		`
		<$wire := import "go.uber.org/thriftrw/wire">

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
		<$bytes := import "bytes">
		<$fmt := import "fmt">
		<$json := import "encoding/json">
		<$math := import "math">
		<$strconv := import "strconv">

		<$wire := import "go.uber.org/thriftrw/wire">

		<$enumName := goName .Spec>
		<formatDoc .Spec.Doc>type <$enumName> int32

		<if .Spec.Items>
			const (
			<range .Spec.Items>
				<- formatDoc .Doc><enumItemName $enumName .> <$enumName> = <.Value>
			<end>
			)

		// <$enumName>_Values returns all recognized values of <$enumName>.
		func <$enumName>_Values() []<$enumName> {
			return []<$enumName>{
				<range .Spec.Items>
					<- enumItemName $enumName .>,
				<end>
			}
		}
		<end>

		<$v := newVar "v">
		// UnmarshalText tries to decode <$enumName> from a byte slice
		// containing its name.
		<- if .Spec.Items>
		//
		//   var <$v> <$enumName>
		//   err := <$v>.UnmarshalText([]byte("<(index .Spec.Items 0).Name>"))
		<- end>
		func (<$v> *<$enumName>) UnmarshalText(value []byte) error {
			switch string(value) {
			<- $enum := .Spec ->
			<range .Spec.Items ->
				case "<.Name>":
					*<$v> = <enumItemName $enumName .>
					return nil
			<end ->
				default:
					return <$fmt>.Errorf("unknown enum value %q for %q", value, "<$enumName>")
			}
		}

		// Ptr returns a pointer to this enum value.
		func (<$v> <$enumName>) Ptr() *<$enumName> {
			return &<$v>
		}

		// ToWire translates <$enumName> into a Thrift-level intermediate
		// representation. This intermediate representation may be serialized
		// into bytes using a ThriftRW protocol implementation.
		//
		// Enums are represented as 32-bit integers over the wire.
		func (<$v> <$enumName>) ToWire() (<$wire>.Value, error) {
			return <$wire>.NewValueI32(int32(<$v>)), nil
		}

		<$w := newVar "w">
		// FromWire deserializes <$enumName> from its Thrift-level
		// representation.
		//
		//   x, err := binaryProtocol.Decode(reader, wire.TI32)
		//   if err != nil {
		//     return <$enumName>(0), err
		//   }
		//
		//   var <$v> <$enumName>
		//   if err := <$v>.FromWire(x); err != nil {
		//     return <$enumName>(0), err
		//   }
		//   return <$v>, nil
		func (<$v> *<$enumName>) FromWire(<$w> <$wire>.Value) error {
			*<$v> = (<$enumName>)(<$w>.GetI32());
			return nil
		}

		// String returns a readable string representation of <$enumName>.
		func (<$v> <$enumName>) String() string {
			<$w> := int32(<$v>)
			<if len .Spec.Items ->
				switch <$w> {
				<range .UniqueItems ->
					case <.Value>:
						return "<.Name>"
				<end ->
				}
			<end ->
			return fmt.Sprintf("<$enumName>(%d)", <$w>)
		}

		<$rhs := newVar "rhs">
		// Equals returns true if this <$enumName> value matches the provided
		// value.
		func (<$v> <$enumName>) Equals(<$rhs> <$enumName>) bool {
			return <$v> == <$rhs>
		}

		// MarshalJSON serializes <$enumName> into JSON.
		//
		// If the enum value is recognized, its name is returned. Otherwise,
		// its integer value is returned.
		//
		// This implements json.Marshaler.
		func (<$v> <$enumName>) MarshalJSON() ([]byte, error) {
			<if len .Spec.Items ->
				switch int32(<$v>) {
				<range .UniqueItems ->
					case <.Value>:
						return ([]byte)("\"<.Name>\""), nil
				<end ->
				}
			<end ->
			return ([]byte)(<$strconv>.FormatInt(int64(<$v>), 10)), nil
		}

		<$text := newVar "text">
		// UnmarshalJSON attempts to decode <$enumName> from its JSON
		// representation.
		//
		// This implementation supports both, numeric and string inputs. If a
		// string is provided, it must be a known enum name.
		//
		// This implements json.Unmarshaler.
		func (<$v> *<$enumName>) UnmarshalJSON(<$text> []byte) error {
			<- $d := newVar "d" ->
			<- $t := newVar "t" ->

			<$d> := <$json>.NewDecoder(<$bytes>.NewReader(<$text>))
			<$d>.UseNumber()
			<$t>, err := <$d>.Token()
			if err != nil {
				return err
			}

			switch <$w> := <$t>.(type) {
			case <$json>.Number:
				<- $x := newVar "x" ->
				<$x>, err := <$w>.Int64()
				if err != nil {
					return err
				}
				if <$x> <">"> <$math>.MaxInt32 {
					return <$fmt>.Errorf("enum overflow from JSON %q for %q", <$text>, "<$enumName>")
				}
				if <$x> <"<"> <$math>.MinInt32 {
					return <$fmt>.Errorf("enum underflow from JSON %q for %q", <$text>, "<$enumName>")
				}
				*<$v> = (<$enumName>)(<$x>)
				return nil
			case string:
				return <$v>.UnmarshalText([]byte(<$w>))
			default:
				return <$fmt>.Errorf("invalid JSON value %q (%T) to unmarshal into %q", <$t>, <$t>, "<$enumName>")
			}
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
// the given (potentially import qualified) enumName and EnumItem.
func enumItemName(enumName string, spec *compile.EnumItem) (string, error) {
	name, err := goNameAnnotation(spec)
	if err != nil {
		return "", err
	}
	if name == "" {
		name = pascalCase(false, /* all caps */
			strings.Split(spec.ThriftName(), "_")...)
	}
	return enumName + name, err
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
