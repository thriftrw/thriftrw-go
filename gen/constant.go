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
	"strconv"

	"github.com/thriftrw/thriftrw-go/compile"
)

// Constant generates code for `const` expressions in Thrift files.
func Constant(g Generator, c *compile.Constant) error {
	err := g.DeclareFromTemplate(
		`<if canBeConstant .Type>const<else>var<end> <goCase .Name> <typeReference .Type> = <constantValue .Value .Type>`,
		c,
		TemplateFunc("constantValue", ConstantValue),
		TemplateFunc("canBeConstant", canBeConstant),
	)
	return wrapGenerateError(c.Name, err)
}

// ConstantValue generates an expression containing the given constant value of
// the given type.
//
// The constant must already have been linked to the given type.
func ConstantValue(g Generator, c compile.ConstantValue, t compile.TypeSpec) (string, error) {
	switch v := c.(type) {
	case compile.ConstantBool:
		if v {
			return "true", nil
		}
		return "false", nil
	case compile.ConstantDouble:
		return fmt.Sprint(float64(v)), nil
	case compile.ConstantInt:
		return fmt.Sprint(int(v)), nil
	case compile.ConstantList:
		l, ok := t.(*compile.ListSpec)
		if !ok {
			return "", fmt.Errorf(
				"cannot cast list %v to %v", v, t.ThriftName())
		}
		return constantList(g, v, l)
	case compile.ConstantMap:
		m, ok := t.(*compile.MapSpec)
		if !ok {
			return "", fmt.Errorf(
				"cannot cast map %v to %v", v, t.ThriftName())
		}
		return constantMap(g, v, m)
	case compile.ConstantSet:
		s, ok := t.(*compile.SetSpec)
		if !ok {
			return "", fmt.Errorf(
				"cannot cast set %v to %v", v, t.ThriftName())
		}
		return constantSet(g, v, s)
	case compile.ConstantString:
		return strconv.Quote(string(v)), nil
	case *compile.ConstantStruct:
		s, ok := t.(*compile.StructSpec)
		if !ok {
			return "", fmt.Errorf(
				"cannot cast struct %v to %v", v, t.ThriftName())
		}
		return constantStruct(g, v, s)
	case compile.EnumItemReference:
		return g.TextTemplate(`<typeName .Enum><goCase .Item.Name>`, v)
	case compile.ConstReference:
		return g.LookupConstantName(v.Target)
	default:
		panic(fmt.Sprintf("Unknown constant value %v (%T)", c, c))
	}
}

func constantList(g Generator, v compile.ConstantList, l *compile.ListSpec) (string, error) {
	return g.TextTemplate(
		`
		<$valueType := .Spec.ValueSpec>
		<typeReference .Spec>{
			<range .Value>
				<constantValue . $valueType>,
			<end>
		}`, struct {
			Spec  *compile.ListSpec
			Value compile.ConstantList
		}{Spec: l, Value: v},
		TemplateFunc("constantValue", ConstantValue))
}

func constantMap(g Generator, v compile.ConstantMap, m *compile.MapSpec) (string, error) {
	return g.TextTemplate(
		`
		<$keyType := .Spec.KeySpec>
		<$valueType := .Spec.ValueSpec>
		<typeReference .Spec>{
			<range .Value>
				<if isHashable $keyType>
					<constantValue .Key $keyType>:
						<constantValue .Value $valueType>,
				<else>
					{
						Key: <constantValue .Key $keyType>,
						Value: <constantValue .Value $valueType>,
					},
				<end>
			<end>
		}`, struct {
			Spec  *compile.MapSpec
			Value compile.ConstantMap
		}{Spec: m, Value: v},
		TemplateFunc("constantValue", ConstantValue))
}

func constantSet(g Generator, v compile.ConstantSet, s *compile.SetSpec) (string, error) {
	return g.TextTemplate(
		`
		<$valueType := .Spec.ValueSpec>
		<typeReference .Spec>{
			<range .Value>
				<if isHashable $valueType>
					<constantValue . $valueType>: struct{}{},
				<else>
					<constantValue . $valueType>,
				<end>
			<end>
		}`, struct {
			Spec  *compile.SetSpec
			Value compile.ConstantSet
		}{Spec: s, Value: v},
		TemplateFunc("constantValue", ConstantValue))
}

func constantStruct(g Generator, v *compile.ConstantStruct, s *compile.StructSpec) (string, error) {
	return g.TextTemplate(
		`
		<$fields := .Spec.Fields>
		&<typeName .Spec>{
			<range $name, $value := .Value.Fields>
				<$field := $fields.FindByName $name>
				<if and (not $field.Required) (isPrimitiveType $field.Type)>
					<goCase $field.Name>: <primitiveValueRef $value $field.Type>,
				<else>
					<goCase $field.Name>: <constantValue $value $field.Type>,
				<end>
			<end>
		}`, struct {
			Spec  *compile.StructSpec
			Value *compile.ConstantStruct
		}{Spec: s, Value: v},
		TemplateFunc("constantValue", ConstantValue),
		TemplateFunc("primitiveValueRef", primitiveValueRef),
	)
}

// helper to generate pointers to primitives
func primitiveValueRef(g Generator, c compile.ConstantValue, t compile.TypeSpec) (string, error) {
	var (
		f   string
		err error
	)

	switch t {
	case compile.BoolSpec:
		f = "_boolptr"
		err = g.EnsureDeclared(`func _boolptr(v bool) *bool { return &v }`, nil)
	case compile.I8Spec:
		f = "_i8ptr"
		err = g.EnsureDeclared(`func _i8ptr(v int8) *int8 { return &v }`, nil)
	case compile.I16Spec:
		f = "_i16ptr"
		err = g.EnsureDeclared(`func _i16ptr(v int16) *int16 { return &v }`, nil)
	case compile.I32Spec:
		f = "_i32ptr"
		err = g.EnsureDeclared(`func _i32ptr(v int32) *int32 { return &v }`, nil)
	case compile.I64Spec:
		f = "_i64ptr"
		err = g.EnsureDeclared(`func _i64ptr(v int64) *int64 { return &v }`, nil)
	case compile.DoubleSpec:
		f = "_doubleptr"
		err = g.EnsureDeclared(`func _doubleptr(v float64) *float64 { return &v }`, nil)
	case compile.StringSpec:
		f = "_stringptr"
		err = g.EnsureDeclared(`func _stringptr(v string) *string { return &v }`, nil)
	}

	switch t.(type) {
	case *compile.EnumSpec, *compile.TypedefSpec:
		f = fmt.Sprintf("_%v_ptr", t.ThriftName())
		err = g.EnsureDeclared(
			`func _<.ThriftName>_ptr(v <typeReference .>) *<typeReference .> {
				return &v
			}`, t)
	}

	if err != nil {
		return "", err
	}

	if f == "" {
		panic(fmt.Sprintf(
			"primitiveValueRef called with %v which is not a primitive",
			t.ThriftName()))
	}

	s, err := ConstantValue(g, c, t)
	s = fmt.Sprintf("%v(%v)", f, s)
	return s, err
}
