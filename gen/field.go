// Copyright (c) 2021 Uber Technologies, Inc.
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

	omitempty    = "omitempty"
	notOmitempty = "!omitempty"
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

	Name       string
	ThriftName string
	Fields     compile.FieldGroup

	// If this field group represents a union of values, exactly one field
	// must be set for it to be valid.
	IsUnion         bool
	AllowEmptyUnion bool

	// This field group represents a Thrift exception.
	IsException bool

	Doc string
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
	if err := verifyUniqueFieldLabels(f.Fields); err != nil {
		return err
	}

	if err := f.DefineStruct(g); err != nil {
		return err
	}

	if err := f.DefineDefaultConstructor(g); err != nil {
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

	if f.IsException {
		if err := f.ErrorName(g); err != nil {
			return err
		}
	}

	if err := f.Equals(g); err != nil {
		return err
	}

	if !checkNoZap(g) {
		if err := f.Zap(g); err != nil {
			return err
		}
	}

	return f.Accessors(g)
}

func (f fieldGroupGenerator) DefineStruct(g Generator) error {
	return g.DeclareFromTemplate(
		`<formatDoc .Doc>type <.Name> struct {
			<range .Fields>
				<- if .Required ->
					<formatDoc .Doc><declFieldName .> <typeReference .Type> <tag .>
				<- else ->
					<formatDoc .Doc><declFieldName .> <typeReferencePtr .Type> <tag .>
				<- end>
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

	// Default to the field name or label as the name used in the JSON
	// representation.
	if err := tags.Set(compileJSONTag(f, entityLabel(f))); err != nil {
		return "", fmt.Errorf("failed to set tag: %v", err)
	}

	// Process go.tags and overwrite JSON tag if specified in Thrift
	// annotation.
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
	t := &structtag.Tag{
		Key:     jsonTagKey,
		Name:    name,
		Options: opts,
	}

	// If the field name is "-" then it means omit, add no tags
	if name == "-" {
		return t
	}

	// We want to add the "omitempty" JSON tag if the field is an "optional" to
	// reduce "null" noise. The "omitempty" json tag specifies that the field
	// should be omitted from the encoding if the field has an empty value,
	// defined as false, 0, a nil pointer, a nil interface value, and any empty
	// array, slice, map, or string.
	//
	// If the field is marked with "!omitempty", then "omitempty" will not be added.
	// If both "!omitempty" and "omitempty" are present, then "omitempty" is removed.
	if (isReferenceType(f.Type) || isStructType(f.Type) || isPrimitiveType(f.Type)) && !f.Required {
		hasNotOmitempty := t.HasOption(notOmitempty)
		hasOmitempty := t.HasOption(omitempty)
		if !hasNotOmitempty && !hasOmitempty {
			// Add omitempty if it's not there and the user
			// has not opted out with !omitempty.
			t.Options = append(t.Options, omitempty)
		} else if hasNotOmitempty && hasOmitempty {
			// !omitempty takes precedence so remove omitempty
			// if present.
			newOptions := make([]string, 0, len(t.Options))
			for _, option := range t.Options {
				if option == omitempty {
					continue
				}
				newOptions = append(newOptions, option)
			}

			t.Options = newOptions
		}
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

func (f fieldGroupGenerator) DefineDefaultConstructor(g Generator) error {
	var hasDefaults bool
	for _, f := range f.Fields {
		if f.Default != nil {
			hasDefaults = true
			break
		}
	}
	if !hasDefaults {
		return nil
	}
	return g.DeclareFromTemplate(
		`
		// Default_<.Name> constructs a new <.Name> struct,
		// pre-populating any fields with defined default values. 
		func Default_<.Name>() *<.Name> {
			<- $v := newVar "v" ->
			var v <.Name>
			<- range .Fields ->
				<- $fname := goName . ->
				<- if isNotNil .Default>
					<$v>.<$fname> = <constantValuePtr .Default .Type>
				<- end ->
			<end>
			return &v
		}
		`, f, TemplateFunc("constantValuePtr", ConstantValuePtr))
}

func (f fieldGroupGenerator) ToWire(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$wire := import "go.uber.org/thriftrw/wire">

		<$v := newVar "v">
		// ToWire translates a <.Name> struct into a Thrift-level intermediate
		// representation. This intermediate representation may be serialized
		// into bytes using a ThriftRW protocol implementation.
		//
		// An error is returned if the struct or any of its fields failed to
		// validate.
		//
		//   x, err := <$v>.ToWire()
		//   if err != nil {
		//     return err
		//   }
		//
		//   if err := binaryProtocol.Encode(x, writer); err != nil {
		//     return err
		//   }
		func (<$v> *<.Name>) ToWire() (<$wire>.Value, error) {
			return <$wire>.NewValueFieldList((*_fieldList_<.Name>)(<$v>)), nil
		}

		type _fieldList_<.Name> <.Name>

		<$fl := newVar "fl">
		func (<$fl> *_fieldList_<.Name>) ForEach(writeField func(<$wire>.Field) error) error {
			<- $i := newVar "i" ->
			<- $wVal := newVar "w" ->

			<if len .Fields ->
				var (
					<$i> int = 0
					<$v> = (*<.Name>)(<$fl>)
					<$wVal> <$wire>.Value
					err error
				)
			<- end>
			<$structName := .Name>
			<range .Fields>
				<- $fname := goName . ->
				<- $f := printf "%s.%s" $v $fname ->
				<- if .Required ->
					<- if and (not (isPrimitiveType .Type)) (not (isListType .Type)) ->
						if <$f> == nil {
							return <import "errors">.New("field <$fname> of <$structName> is required")
						}
					<- end>
						<$wVal>, err = <toWire .Type $f>
						if err != nil {
							return err
						}
						if err := writeField(<$wire>.Field{ID: <.ID>, Value: <$wVal>}); err != nil {
							return err
						}
						<$i>++
				<- else ->
					<- if isNotNil .Default ->
						<- $fval := printf "%s%s" $v $fname ->
						<$fval> := <$f>
						if <$fval> == nil {
							<$fval> = <constantValuePtr .Default .Type>
						}
						{
							<$wVal>, err = <toWirePtr .Type $fval>
					<- else ->
						if <$f> != nil {
							<$wVal>, err = <toWirePtr .Type $f>
					<- end>
							if err != nil {
								return err
							}
							if err := writeField(<$wire>.Field{ID: <.ID>, Value: <$wVal>}); err != nil {
								return err
							}
							<$i>++
						}
				<- end>
			<end>

			<if and .IsUnion (len .Fields)>
				<$fmt := import "fmt">
				<if .AllowEmptyUnion>
					if <$i> > 1 {
						return <$fmt>.Errorf("<.Name> should have at most one field: got %v fields", <$i>)
					}
				<else>
					if <$i> != 1 {
						return <$fmt>.Errorf("<.Name> should have exactly one field: got %v fields", <$i>)
					}
				<end>
			<end>

			return nil
		}

		func (<$fl> *_fieldList_<.Name>) Close() {}
		`, f, TemplateFunc("constantValuePtr", ConstantValuePtr))
}

func (f fieldGroupGenerator) FromWire(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$wire := import "go.uber.org/thriftrw/wire">

		<$v := newVar "v">
		<$w := newVar "w">
		// FromWire deserializes a <.Name> struct from its Thrift-level
		// representation. The Thrift-level representation may be obtained
		// from a ThriftRW protocol implementation.
		//
		// An error is returned if we were unable to build a <.Name> struct
		// from the provided intermediate representation.
		//
		//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
		//   if err != nil {
		//     return nil, err
		//   }
		//
		//   var <$v> <.Name>
		//   if err := <$v>.FromWire(x); err != nil {
		//     return nil, err
		//   }
		//   return &<$v>, nil
		func (<$v> *<.Name>) FromWire(<$w> <$wire>.Value) error {
			<$f := newVar "field">

			<$isSet := newNamespace>
			<range .Fields>
				<- if .Required ->
					<$isSet.NewName (printf "%sIsSet" .Name)> := false
				<- end>
			<end>

			fields := <$w>.GetFieldList()
			err := fields.ForEach(func(<$f> <$wire>.Field) (err error) {
				switch <$f>.ID {
				<range .Fields ->
				case <.ID>:
					if <$f>.Value.Type() == <typeCode .Type> {
						<- $lhs := printf "%s.%s" $v (goName .) ->
						<- $value := printf "%s.Value" $f ->
						<- if .Required ->
							<$lhs>, err = <fromWire .Type $value>
						<- else ->
							<fromWirePtr .Type $lhs $value>
						<- end>
						if err != nil {
							return err
						}
						<if .Required ->
							<$isSet.Rotate (printf "%sIsSet" .Name)> = true
						<- end>
					}
				<end ->
				}
				return nil
			})
			if err != nil {
				return err
			}
			fields.Close()

			<$structName := .Name>
			<range .Fields>
				<$fname := goName .>
				<$f := printf "%s.%s" $v $fname>
				<if isNotNil .Default>
					if <$f> == nil {
						<$f> = <constantValuePtr .Default .Type>
					}
				<else>
					<if .Required>
						if !<$isSet.Rotate (printf "%sIsSet" .Name)> {
							return <import "errors">.New("field <$fname> of <$structName> is required")
						}
					<end>
				<end>
			<end>

			<if and .IsUnion (len .Fields)>
				<$fmt := import "fmt">
				<$count := newVar "count">
				<$count> := 0
				<range .Fields ->
					if <$v>.<goName .> != nil {
						<$count>++
					}
				<end>
				<- if .AllowEmptyUnion ->
					if <$count> > 1 {
						return <$fmt>.Errorf( "<.Name> should have at most one field: got %v fields", <$count>)
					}
				<- else ->
					if <$count> != 1 {
						return <$fmt>.Errorf( "<.Name> should have exactly one field: got %v fields", <$count>)
					}
				<- end>
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
		// String returns a readable string representation of a <.Name>
		// struct.
		func (<$v> *<.Name>) String() string {
			if <$v> == nil {
				return "<"<nil>">"
			}

			<$fields := newVar "fields">
			<$i := newVar "i">

			var <$fields> [<len .Fields>]string
			<$i> := 0
			<range .Fields>
				<- $fname := goName . ->
				<- $f := printf "%s.%s" $v $fname ->

				<- if not .Required ->
					if <$f> != nil {
						<if isPrimitiveType .Type ->
							<$fields>[<$i>] = <$fmt>.Sprintf("<$fname>: %v", *(<$f>))
						<- else ->
							<$fields>[<$i>] = <$fmt>.Sprintf("<$fname>: %v", <$f>)
						<- end>
						<$i>++
					}
				<- else ->
					<$fields>[<$i>] = <$fmt>.Sprintf("<$fname>: %v", <$f>)
					<$i>++
				<- end>
			<end>

			return <$fmt>.Sprintf("<.Name>{%v}", <$strings>.Join(<$fields>[:<$i>], ", "))
		}
		`, f)
}

func (f fieldGroupGenerator) ErrorName(g Generator) error {
	return g.DeclareFromTemplate(
		`
		// ErrorName is the name of this type as defined in the Thrift
		// file.
		func (*<.Name>) ErrorName() string {
			return "<.ThriftName>"
		}
		`, f)
}

func (f fieldGroupGenerator) Equals(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$v := newVar "v">
		<$rhs := newVar "rhs">
		// Equals returns true if all the fields of this <.Name> match the
		// provided <.Name>.
		//
		// This function performs a deep comparison.
		func (<$v> *<.Name>) Equals(<$rhs> *<.Name>) bool {
			if <$v> == nil {
				return <$rhs> == nil
			} else if <$rhs> == nil {
				return false
			}
			<range .Fields>
				<- $fname := goName . ->
				<- $lhsField := printf "%s.%s" $v $fname ->
				<- $rhsField := printf "%s.%s" $rhs $fname ->

				<- if .Required ->
					if !<equals .Type $lhsField $rhsField> {
						return false
					}
				<- else ->
					if !<equalsPtr .Type $lhsField $rhsField> {
						return false
					}
				<- end>
			<end>
			return true
		}
		`, f)
}

func (f fieldGroupGenerator) Zap(g Generator) error {
	return g.DeclareFromTemplate(
		`
		<$zapcore := import "go.uber.org/zap/zapcore">
		<$v := newVar "v">
		<$enc := newVar "enc">

		// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
		// fast logging of <.Name>.
		func (<$v> *<.Name>) MarshalLogObject(<$enc> <$zapcore>.ObjectEncoder) (err error) {
			if <$v> == nil {
				return nil
			}
			<range .Fields>
				<- if not (zapOptOut .) ->
					<- $fval := printf "%s.%s" $v (goName .) ->
					<- if .Required ->
						<zapEncodeBegin .Type ->
							<$enc>.Add<zapEncoder .Type>("<fieldLabel .>", <zapMarshaler .Type $fval>)
						<- zapEncodeEnd .Type>
					<- else ->
						if <$fval> != nil {
							<zapEncodeBegin .Type ->
								<$enc>.Add<zapEncoder .Type>("<fieldLabel .>", <zapMarshalerPtr .Type $fval>)
							<- zapEncodeEnd .Type>
						}
					<- end>
				<- end>
			<end ->
			return err
		}
		`, f,
		TemplateFunc("zapOptOut", zapOptOut),
		TemplateFunc("fieldLabel", entityLabel),
	)
}

func (f fieldGroupGenerator) Accessors(g Generator) error {
	// Namespace to ensure that field names don't conflict with method names.
	fieldsAndMethods := NewNamespace()

	return g.DeclareFromTemplate(
		`
		<$v := newVar "v">
		<$o := newVar "o">
		<$name := .Name>

		<range .Fields>
			<$fname := goName .>
			<reserveFieldOrMethod $fname>

			<reserveFieldOrMethod (printf "Get%v" $fname)>
			// Get<$fname> returns the value of <$fname> if it is set or its
			// <if isNotNil .Default>default<else>zero<end> value if it is unset.
			func (<$v> *<$name>) Get<$fname>() (<$o> <typeReference .Type>) {
				<- if .Required ->
				  if <$v> != nil {
				    <$o> = <$v>.<$fname>
				  }
				  return
				<- else ->
				  if <$v> != nil && <$v>.<$fname> != nil {
					<- if and (not .Required) (isPrimitiveType .Type) ->
					  return *<$v>.<$fname>
					<- else ->
					  return <$v>.<$fname>
					<- end ->
				  }
				  <if isNotNil .Default><$o> = <constantValue .Default .Type><end>
				  return
				<- end ->
			}

			<if shouldGenerateIsSet .>
				<reserveFieldOrMethod (printf "IsSet%v" $fname)>
				// IsSet<$fname> returns true if <$fname> is not nil.
				func (<$v> *<$name>) IsSet<$fname>() bool {
					return <$v> != nil && <$v>.<$fname> != nil
				}
			<end>
		<end>
		`, f,
		TemplateFunc("constantValue", ConstantValue),
		TemplateFunc("shouldGenerateIsSet", func(f *compile.FieldSpec) bool {
			// Generate IsSet functions for a field only if the field is
			// optional or the field value itself is nillable.
			return !f.Required || isReferenceType(f.Type) || isStructType(f.Type)
		}),
		TemplateFunc("reserveFieldOrMethod", func(name string) (string, error) {
			// we return an empty string for the sake of the templating system
			err := fieldsAndMethods.Reserve(name)
			return "", err
		}),
	)
}

func verifyUniqueFieldLabels(fs compile.FieldGroup) error {
	used := make(map[string]*compile.FieldSpec, len(fs))
	for _, f := range fs {
		label := entityLabel(f)
		if conflict, isUsed := used[label]; isUsed {
			return fmt.Errorf(
				"field %q with label %q conflicts with field %q",
				f.Name, label, conflict.Name)
		}
		used[label] = f
	}
	return nil
}
