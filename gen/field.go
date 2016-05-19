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
	// TODO(abg): Default values
	return g.DeclareFromTemplate(
		`
		<$wire := import "github.com/thriftrw/thriftrw-go/wire">

		<$v := newVar "v">
		func (<$v> *<.Name>) ToWire() <$wire>.Value {
    		<$fields := newVar "fields">
    		<$i := newVar "i">
			// TODO check if required fields that are reference types are nil

			var <$fields> [<len .Fields>]<$wire>.Field
			<$i> := 0

			<range .Fields>
				<$f := printf "%s.%s" $v (goCase .Name)>
				<if .Required>
					<$wVal := toWire .Type $f>
					<$fields>[<$i>] = <$wire>.Field{ID: <.ID>, Value: <$wVal>}
					<$i>++
				<else>
					<if .Default>
						if <$f> == nil {
							<$f> = <constantValuePtr .Default .Type>
						}
						<$fields>[<$i>] = <$wire>.Field{
							ID: <.ID>,
							Value: <toWirePtr .Type $f>,
						}
						<$i>++
					<else>
						if <$f> != nil {
							<$fields>[<$i>] = <$wire>.Field{
								ID: <.ID>,
								Value: <toWirePtr .Type $f>,
							}
							<$i>++
						}
					<end>
				<end>
			<end>

			return <$wire>.NewValueStruct(
				<$wire>.Struct{Fields: <$fields>[:<$i>]},
			)
		}
		`, f, TemplateFunc("constantValuePtr", ConstantValuePtr))
}

func (f fieldGroupGenerator) FromWire(g Generator) error {
	// TODO(abg): Default values
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
					}
				<end>
				}
			}

			<range .Fields>
				<$f := printf "%s.%s" $v (goCase .Name)>
				<if .Default>
					if <$f> == nil {
						<$f> = <constantValuePtr .Default .Type>
					}
				<end>
			<end>

			// TODO(abg): Check that all required fields were set.
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
