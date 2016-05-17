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

// Constant TODO
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
		return g.TextTemplate(
			`
			<$valueType := .Spec.ValueSpec>
			<typeReference .Spec>{
				<range .Value>
					<constantValue . $valueType>,
				<end>
			}`, struct {
				Spec  compile.TypeSpec
				Value compile.ConstantList
			}{Spec: t, Value: v},
			TemplateFunc("constantValue", ConstantValue))
	case compile.ConstantMap:
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
				Spec  compile.TypeSpec
				Value compile.ConstantMap
			}{Spec: t, Value: v},
			TemplateFunc("constantValue", ConstantValue))
	case compile.ConstantSet:
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
				Spec  compile.TypeSpec
				Value compile.ConstantSet
			}{Spec: t, Value: v},
			TemplateFunc("constantValue", ConstantValue))
	case compile.ConstantString:
		return strconv.Quote(string(v)), nil
	case *compile.ConstantStruct:
		return g.TextTemplate(
			`
			<$fields := .Value.Fields>
			&<typeName .Spec>{
				<range .Spec.Fields>
					<$value := index $fields .Name>
					<if $value>
						<goCase .Name>: <constantValue $value .Type>,
					<end>
				<end>
			}`, struct {
				Spec  compile.TypeSpec
				Value *compile.ConstantStruct
			}{Spec: t, Value: v},
			TemplateFunc("constantValue", ConstantValue))
	case compile.EnumItemReference:
		return g.TextTemplate(`<typeName .Enum><goCase .Item.Name>`, v)
	case compile.ConstReference:
		return g.LookupConstantName(v.Target)
	default:
		panic(fmt.Sprintf("Unknown constant value %v (%T)", c, c))
	}
}
