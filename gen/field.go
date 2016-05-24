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
	"fmt"

	"github.com/thriftrw/thriftrw-go/compile"
)

// fieldGroupGenerator is responsible for generating code for FieldGroups.
type fieldGroupGenerator struct {
	Name   string
	Fields compile.FieldGroup

	// If this field group represents a union of values, exactly one field
	// must be set for it to be valid.
	IsUnion         bool
	AllowEmptyUnion bool
}

func (f fieldGroupGenerator) Generate(g Generator) error {
	if err := f.DefineStruct(g); err != nil {
		return err
	}

	if err := f.ToWire(g); err != nil {
		return err
	}

	if err := f.FromWire(g); err != nil {
		return err
	}

	if err := f.String(g); err != nil {
		return err
	}

	return nil
}

func (f fieldGroupGenerator) DefineStruct(g Generator) error {
	return g.DeclareFromTemplate(
		`type <.Name> struct {
			<range .Fields>
				<if .Required>
					<goCase .Name> <typeReference .Type> <tag .>
				<else>
					<goCase .Name> <typeReferencePtr .Type> <tag .>
				<end>
			<end>
		}`,
		f,
		TemplateFunc("tag", func(f *compile.FieldSpec) string {
			// We want to add omitempty if the field is an optional struct or
			// primitive to redure "null" noise. We won't add omitempty for
			// optional collections because omitempty doesn't differentiate
			// between nil and empty collections.

			if (isStructType(f.Type) || isPrimitiveType(f.Type)) && !f.Required {
				return fmt.Sprintf("`json:\"%s,omitempty\"`", f.Name)
			}

			return fmt.Sprintf("`json:%q`", f.Name)
			// TODO(abg): Take go.tag and js.name annotations into account
		}),
	)
}

func (f fieldGroupGenerator) ToWire(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$v := newVar "v">
		func (<$v> *<.Name>) ToWire() (<$wire>.Value, error) {
    		<$fields := newVar "fields">
    		<$i := newVar "i">
			<$wVal := newVar "w">

			var (
					<$fields> [<len .Fields>]<$wire>.Field
					<$i> int = 0
				<if len .Fields>
					<$wVal> <$wire>.Value
					err error
				<end>
			)

			<$structName := .Name>
			<range .Fields>
				<$f := printf "%s.%s" $v (goCase .Name)>
				<if .Required>
					<if not (isPrimitiveType .Type)>
						if <$f> == nil {
							return <$wVal>, <import "errors">.New(
								"field <goCase .Name> of <$structName> is required")
						}
					<end>
						<$wVal>, err = <toWire .Type $f>
						if err != nil {
							return <$wVal>, err
						}
						<$fields>[<$i>] = <$wire>.Field{
							ID: <.ID>,
							Value: <$wVal>,
						}
						<$i>++
				<else>
					<if .Default>
						if <$f> == nil {
							<$f> = <constantValuePtr .Default .Type>
						}
					<else>
						if <$f> != nil {
					<end>
							<$wVal>, err = <toWirePtr .Type $f>
							if err != nil {
								return <$wVal>, err
							}
							<$fields>[<$i>] = <$wire>.Field{
								ID: <.ID>,
								Value: <$wVal>,
							}
							<$i>++
					<if not .Default>
						}
					<end>
				<end>
			<end>

			<if and .IsUnion (len .Fields)>
				<$fmt := import "fmt">
				<if .AllowEmptyUnion>
					if <$i> > 1 {
						return <$wire>.Value{}, <$fmt>.Errorf(
							"<.Name> should receive at most one field value: received %v values", <$i>)
					}
				<else>
					if <$i> != 1 {
						return <$wire>.Value{}, <$fmt>.Errorf(
							"<.Name> should receive exactly one field value: received %v values", <$i>)
					}
				<end>
			<end>

			return <$wire>.NewValueStruct(
				<$wire>.Struct{Fields: <$fields>[:<$i>]},
			), nil
		}
		`, f, TemplateFunc("constantValuePtr", ConstantValuePtr))
}

func (f fieldGroupGenerator) FromWire(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$v := newVar "v">
		<$w := newVar "w">
		func (<$v> *<.Name>) FromWire(<$w> <$wire>.Value) error {
			<if len .Fields>
				var err error
			<end>
			<$f := newVar "field">

			<$isSet := newNamespace>
			<range .Fields>
				<if .Required>
					<$isSet.NewName (printf "%sIsSet" .Name)> := false
				<end>
			<end>

			<$count := newVar "count">
			<if and .IsUnion (len .Fields)>
				<$count> := 0
			<end>

			<$isUnion := .IsUnion>
			for _, <$f> := range <$w>.GetStruct().Fields {
				switch <$f>.ID {
				<range .Fields>
				case <.ID>:
					if <$f>.Value.Type() == <typeCode .Type> {
						<$lhs := printf "%s.%s" $v (goCase .Name)>
						<$value := printf "%s.Value" $f>
						<if .Required>
							<$lhs>, err = <fromWire .Type $value>
						<else>
							<fromWirePtr .Type $lhs $value>
						<end>
						if err != nil {
							return err
						}
						<if .Required>
							<$isSet.Rotate (printf "%sIsSet" .Name)> = true
						<end>
						<if $isUnion><$count>++<end>
					}
				<end>
				}
			}

			<$structName := .Name>
			<range .Fields>
				<$f := printf "%s.%s" $v (goCase .Name)>
				<if .Default>
					if <$f> == nil {
						<$f> = <constantValuePtr .Default .Type>
					}
				<else>
					<if .Required>
						if !<$isSet.Rotate (printf "%sIsSet" .Name)> {
							return <import "errors">.New(
								"field <goCase .Name> of <$structName> is required")
						}
					<end>
				<end>
			<end>

			<if and .IsUnion (len .Fields)>
				<$fmt := import "fmt">
				<if .AllowEmptyUnion>
					if <$count> > 1 {
						return <$fmt>.Errorf(
							"<.Name> should receive at most one field value: received %v values", <$count>)
					}
				<else>
					if <$count> != 1 {
						return <$fmt>.Errorf(
							"<.Name> should receive exactly one field value: received %v values", <$count>)
					}
				<end>
			<end>
			return nil
		}
		`, f, TemplateFunc("constantValuePtr", ConstantValuePtr))
}

func (f fieldGroupGenerator) String(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$fmt := import "fmt">
		<$strings := import "strings">

		<$v := newVar "v">
		func (<$v> *<.Name>) String() string {
    		<$fields := newVar "fields">
    		<$i := newVar "i">

			var <$fields> [<len .Fields>]string
			<$i> := 0
			<range .Fields>
				<$f := printf "%s.%s" $v (goCase .Name)>

				<if not .Required>
					if <$f> != nil {
						<if isPrimitiveType .Type>
							<$fields>[<$i>] = <$fmt>.Sprintf("<goCase .Name>: %v", *(<$f>))
						<else>
							<$fields>[<$i>] = <$fmt>.Sprintf("<goCase .Name>: %v", <$f>)
						<end>
						<$i>++
					}
				<else>
					<$fields>[<$i>] = <$fmt>.Sprintf("<goCase .Name>: %v", <$f>)
					<$i>++
				<end>
			<end>

			return <$fmt>.Sprintf(
                "<.Name>{%v}", <$strings>.Join(<$fields>[:<$i>], ", "))
		}
		`, f)
}
