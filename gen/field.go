// Copyright (c) 2017 Uber Technologies, Inc.
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

	"github.com/fatih/structtag"
	"go.uber.org/thriftrw/compile"
)

const (
	// value of this annotation tag will be parsed and included
	// in the generated go struct's tag
	goTagKey = "go.tag"

	// key for tag set on all generated go structs used by encoding/json
	jsonTagKey = "json"
)

var reservedIdentifiers = map[string]struct{}{
	"ToWire":   {},
	"FromWire": {},
	"String":   {},
	"Equals":   {},
}

// fieldGroupGenerator is responsible for generating code for FieldGroups.
type fieldGroupGenerator struct {
	Namespace

	Name   string
	Fields compile.FieldGroup

	// If this field group represents a union of values, exactly one field
	// must be set for it to be valid.
	IsUnion         bool
	AllowEmptyUnion bool

	// This field group represents a Thrift exception.
	IsException bool
}

func (f fieldGroupGenerator) checkReservedIdentifier(name string) error {
	_, match := reservedIdentifiers[name]
	match = match || (f.IsException && name == "Error")
	if match {
		return fmt.Errorf("%q is a reserved ThriftRW identifier", name)
	}
	return nil
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

	if err := f.Equals(g); err != nil {
		return err
	}

	return nil
}

func (f fieldGroupGenerator) DefineStruct(g Generator) error {
	return g.DeclareFromTemplate(
		`type <.Name> struct {
			<range .Fields>
				<if .Required>
					<declFieldName .> <typeReference .Type> <tag .>
				<else>
					<declFieldName .> <typeReferencePtr .Type> <tag .>
				<end>
			<end>
		}`,
		f,
		TemplateFunc("tag", generateTags),
		TemplateFunc("declFieldName", f.declFieldName),
	)
}

// generateTags parses the annotation on the thrift field and creates the resulting go tag
func generateTags(f *compile.FieldSpec) (string, error) {
	tags, err := structtag.Parse("") // no tags
	if err != nil {
		return "", fmt.Errorf("failed to parse tag: %v", err)
	}

	if err := tags.Set(compileJSONTag(f, f.Name)); err != nil {
		return "", fmt.Errorf("failed to set tag: %v", err)
	}

	// process go tags and overwrite json tag if specified in thrift annotation
	if goAnnotation := f.Annotations[goTagKey]; goAnnotation != "" {
		goTags, err := structtag.Parse(goAnnotation)
		if err != nil {
			return "", fmt.Errorf("failed to parse tags %q: %v", goAnnotation, err)
		}

		for _, t := range goTags.Tags() {
			if t.Key == jsonTagKey {
				t = compileJSONTag(f, t.Name, t.Options...)
			}
			if err := tags.Set(t); err != nil {
				return "", fmt.Errorf("failed to set tag: %v", err)
			}
		}
	}

	return fmt.Sprintf("`%s`", tags.String()), nil
}

func compileJSONTag(f *compile.FieldSpec, name string, opts ...string) *structtag.Tag {
	// We want to add omitempty if the field is an optional struct or
	// primitive to reduce "null" noise. We won't add omitempty for
	// optional collections because omitempty doesn't differentiate
	// between nil and empty collections.

	t := &structtag.Tag{
		Key:     jsonTagKey,
		Name:    name,
		Options: opts,
	}

	if name == "-" {
		// If the field name is "-" then it means omit, add no tags
		return t
	}

	if (isStructType(f.Type) || isPrimitiveType(f.Type)) && !f.Required && !t.HasOption("omitempty") {
		t.Options = append(t.Options, "omitempty")
	}

	if f.Required && !t.HasOption("required") {
		t.Options = append(t.Options, "required")
	}

	return t
}

// declFieldName replaces goName during generation of a structure's definition.
// It replicates goName but also register all field names in the
// fieldGroupGenerator namespace, enforcing single field definition when
// generating Go code. TL;DR: will fail during generation, before compilation.
func (f *fieldGroupGenerator) declFieldName(fs *compile.FieldSpec) (string, error) {
	name, fromAnnotation, err := goNameForNamedEntity(fs)
	if err != nil {
		return "", err
	}

	if err = f.checkReservedIdentifier(name); err == nil {
		err = f.Reserve(name)
	}

	if err != nil {
		originalName := (name == fs.ThriftName())
		var note string
		switch {
		case originalName && fromAnnotation:
			note = " (from go.name annotation)"
		case !originalName && fromAnnotation:
			note = fmt.Sprintf(" (from %q go.name annotation)", fs.ThriftName())
		case !originalName && !fromAnnotation:
			note = fmt.Sprintf(" (from %q)", fs.ThriftName())
		}
		return "", fmt.Errorf("could not declare field %q%s: %v", name, note, err)
	}
	return name, nil
}

func (f fieldGroupGenerator) ToWire(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$wire := import "go.uber.org/thriftrw/wire">

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
				<$fname := goName .>
				<$f := printf "%s.%s" $v $fname>
				<if .Required>
					<if not (isPrimitiveType .Type)>
						if <$f> == nil {
							// TODO: Include names of all missing fields in
							// the error message.
							return <$wVal>, <import "errors">.New(
								"field <$fname> of <$structName> is required")
						}
					<end>
						<$wVal>, err = <toWire .Type $f>
						if err != nil {
							// TODO: Nest the error inside a "failed to
							// serialize field X of struct Y" error.
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
						{
					<else>
						if <$f> != nil {
					<end>
							<$wVal>, err = <toWirePtr .Type $f>
							if err != nil {
								// TODO: Nest the error inside a "failed to
								// serialize field X of struct Y" error.
								return <$wVal>, err
							}
							<$fields>[<$i>] = <$wire>.Field{
								ID: <.ID>,
								Value: <$wVal>,
							}
							<$i>++
						}
				<end>
			<end>

			<if and .IsUnion (len .Fields)>
				<$fmt := import "fmt">
				<if .AllowEmptyUnion>
					if <$i> > 1 {
						return <$wire>.Value{}, <$fmt>.Errorf(
							"<.Name> should have at most one field: got %v fields", <$i>)
					}
				<else>
					if <$i> != 1 {
						return <$wire>.Value{}, <$fmt>.Errorf(
							"<.Name> should have exactly one field: got %v fields", <$i>)
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
		<$wire := import "go.uber.org/thriftrw/wire">

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

			for _, <$f> := range <$w>.GetStruct().Fields {
				switch <$f>.ID {
				<range .Fields>
				case <.ID>:
					if <$f>.Value.Type() == <typeCode .Type> {
						<$lhs := printf "%s.%s" $v (goName .)>
						<$value := printf "%s.Value" $f>
						<if .Required>
							<$lhs>, err = <fromWire .Type $value>
						<else>
							<fromWirePtr .Type $lhs $value>
						<end>
						if err != nil {
							return err
							// TODO: Nest the error inside a "failed to read
							// field X of struct Y" error.
						}
						<if .Required>
							<$isSet.Rotate (printf "%sIsSet" .Name)> = true
						<end>
					}
				<end>
				}
			}

			<$structName := .Name>
			<range .Fields>
				<$fname := goName .>
				<$f := printf "%s.%s" $v $fname>
				<if .Default>
					if <$f> == nil {
						<$f> = <constantValuePtr .Default .Type>
					}
				<else>
					<if .Required>
						if !<$isSet.Rotate (printf "%sIsSet" .Name)> {
							return <import "errors">.New(
								"field <$fname> of <$structName> is required")
						}
						// TODO: Include names of all missing fields in the
						// error message.
					<end>
				<end>
			<end>

			<if and .IsUnion (len .Fields)>
				<$fmt := import "fmt">
				<$count := newVar "count">
				<$count> := 0
				<range .Fields>
					if <$v>.<goName .> != nil { <$count>++ }
				<end>
				<if .AllowEmptyUnion>
					if <$count> > 1 {
						return <$fmt>.Errorf(
							"<.Name> should have at most one field: got %v fields", <$count>)
					}
				<else>
					if <$count> != 1 {
						return <$fmt>.Errorf(
							"<.Name> should have exactly one field: got %v fields", <$count>)
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
			if <$v> == nil {
				return "<"<nil>">"
			}

			<$fields := newVar "fields">
			<$i := newVar "i">

			var <$fields> [<len .Fields>]string
			<$i> := 0
			<range .Fields>
				<$fname := goName .>
				<$f := printf "%s.%s" $v $fname>

				<if not .Required>
					if <$f> != nil {
						<if isPrimitiveType .Type>
							<$fields>[<$i>] = <$fmt>.Sprintf("<$fname>: %v", *(<$f>))
						<else>
							<$fields>[<$i>] = <$fmt>.Sprintf("<$fname>: %v", <$f>)
						<end>
						<$i>++
					}
				<else>
					<$fields>[<$i>] = <$fmt>.Sprintf("<$fname>: %v", <$f>)
					<$i>++
				<end>
			<end>

			return <$fmt>.Sprintf(
				"<.Name>{%v}", <$strings>.Join(<$fields>[:<$i>], ", "))
		}
		`, f)
}

func (f fieldGroupGenerator) Equals(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$v := newVar "v">
		<$rhs := newVar "rhs">
		func (<$v> *<.Name>) Equals(<$rhs> *<.Name>) bool {
			<range .Fields>
				<$fname := goName .>
				<$lhsField := printf "%s.%s" $v $fname>
				<$rhsField := printf "%s.%s" $rhs $fname>

				<if .Required>
					if !<equals .Type $lhsField $rhsField> {
						return false
					}
				<else>
					if !<equalsPtr .Type $lhsField $rhsField> {
						return false
					}
				<end>
			<end>
			return true
		}
		`, f)
}
