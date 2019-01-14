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

package compile

import (
	"errors"
	"fmt"

	"go.uber.org/thriftrw/ast"
)

// ConstantValue represents a compiled constant value or a reference to one.
type ConstantValue interface {
	// Link the constant value with the given scope, casting it to the given
	// type if necessary.
	Link(scope Scope, t TypeSpec) (ConstantValue, error)
}

// compileConstantValue compiles a constant value AST into a ConstantValue.
func compileConstantValue(v ast.ConstantValue) ConstantValue {
	if v == nil {
		return nil
	}

	// TODO(abg): Support typedefs

	switch src := v.(type) {
	case ast.ConstantReference:
		return constantReference(src)
	case ast.ConstantMap:
		return compileConstantMap(src)
	case ast.ConstantList:
		return compileConstantList(src)
	case ast.ConstantBoolean:
		return ConstantBool(src)
	case ast.ConstantInteger:
		return ConstantInt(src)
	case ast.ConstantString:
		return ConstantString(src)
	case ast.ConstantDouble:
		return ConstantDouble(src)
	default:
		panic(fmt.Sprintf("unknown constant value of type %T: %v", src, src))
	}
}

type (
	// ConstantBool represents a boolean constant from the Thrift file.
	ConstantBool bool

	// ConstantInt represents an integer constant from the Thrift file.
	ConstantInt int64

	// ConstantString represents a string constant from the Thrift file.
	ConstantString string

	// ConstantDouble represents a double constant from the Thrift file.
	ConstantDouble float64
)

// Link for ConstantBool
func (c ConstantBool) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	if _, ok := RootTypeSpec(t).(*BoolSpec); !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}
	return c, nil
}

// Link for ConstantInt.
func (c ConstantInt) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	rt := RootTypeSpec(t)
	switch spec := rt.(type) {
	case *I8Spec, *I16Spec, *I32Spec, *I64Spec:
		// TODO bounds checks?
		return c, nil
	case *DoubleSpec:
		return ConstantDouble(float64(c)).Link(scope, t)
	case *BoolSpec:
		switch v := int64(c); v {
		case 0, 1:
			return ConstantBool(v == 1).Link(scope, t)
		default:
			return nil, constantValueCastError{
				Value:  c,
				Type:   t,
				Reason: errors.New("the value must be 0 or 1"),
			}
		}
	case *EnumSpec:
		for _, item := range spec.Items {
			if item.Value == int32(c) {
				return EnumItemReference{Enum: spec, Item: &item}, nil
			}
		}

		return nil, constantValueCastError{
			Value: c,
			Type:  t,
			Reason: fmt.Errorf(
				"%v is not a valid value for enum %q", int32(c), spec.ThriftName()),
		}
	}

	return nil, constantValueCastError{Value: c, Type: t}
	// TODO: AST for constants will need to track positions for us to
	// include them in the error messages.
}

// Link for ConstantString.
func (c ConstantString) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	// TODO(abg): Are binary literals a thing?
	if _, ok := RootTypeSpec(t).(*StringSpec); !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}
	return c, nil
}

// Link for ConstantDouble.
func (c ConstantDouble) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	if _, ok := RootTypeSpec(t).(*DoubleSpec); !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}
	return c, nil
}

// ConstantStruct represents a struct literal from the Thrift file.
type ConstantStruct struct {
	Fields map[string]ConstantValue
}

// buildConstantStruct builds a constant struct from a ConstantMap.
func buildConstantStruct(c ConstantMap) (*ConstantStruct, error) {
	fields := make(map[string]ConstantValue, len(c))
	for _, pair := range c {
		s, isString := pair.Key.(ConstantString)
		if !isString {
			return nil, fmt.Errorf(
				"%v is not a string: all keys must be strings", pair.Key)
		}
		fields[string(s)] = pair.Value
	}
	return &ConstantStruct{Fields: fields}, nil
}

// Link for ConstantStruct
func (c *ConstantStruct) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	s, ok := RootTypeSpec(t).(*StructSpec)
	if !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}

	for _, field := range s.Fields {
		f, ok := c.Fields[field.Name]
		if !ok {
			if field.Default == nil {
				if field.Required {
					return nil, constantValueCastError{
						Value:  c,
						Type:   t,
						Reason: fmt.Errorf("%q is a required field", field.Name),
					}
				}
				continue
			}
			f = field.Default
			c.Fields[field.Name] = f
		}

		f, err := f.Link(scope, field.Type)
		if err != nil {
			return nil, constantValueCastError{
				Value: c,
				Type:  t,
				Reason: constantStructFieldCastError{
					FieldName: field.Name,
					Reason:    err,
				},
			}
		}

		c.Fields[field.Name] = f
	}

	return c, nil
}

// ConstantMap represents a map literal from the Thrift file.
type ConstantMap []ConstantValuePair

func compileConstantMap(src ast.ConstantMap) ConstantMap {
	items := make([]ConstantValuePair, len(src.Items))
	for i, item := range src.Items {
		items[i] = ConstantValuePair{
			Key:   compileConstantValue(item.Key),
			Value: compileConstantValue(item.Value),
		}
	}
	return ConstantMap(items)
}

// ConstantValuePair represents a key-value pair of ConstantValues in a constant
// map.
type ConstantValuePair struct {
	Key, Value ConstantValue
}

// Link for ConstantMap.
func (c ConstantMap) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	rt := RootTypeSpec(t)
	if _, isStruct := rt.(*StructSpec); isStruct {
		cs, err := buildConstantStruct(c)
		if err != nil {
			return nil, constantValueCastError{
				Value:  c,
				Type:   t,
				Reason: err,
			}
		}
		return cs.Link(scope, t)
	}

	m, ok := rt.(*MapSpec)
	if !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}

	items := make([]ConstantValuePair, len(c))
	for i, item := range c {
		key, err := item.Key.Link(scope, m.KeySpec)
		if err != nil {
			return nil, err
		}

		value, err := item.Value.Link(scope, m.ValueSpec)
		if err != nil {
			return nil, err
		}

		// TODO(abg): Duplicate key check
		items[i] = ConstantValuePair{Key: key, Value: value}
	}

	return ConstantMap(items), nil
}

// ConstantSet represents a set of constant values from the Thrift file.
type ConstantSet []ConstantValue

// Link for ConstantSet.
func (c ConstantSet) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	s, ok := RootTypeSpec(t).(*SetSpec)
	if !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}

	// TODO(abg): Track whether things are linked so that we don't re-link here
	// TODO(abg): Fail for duplicates
	values := make([]ConstantValue, len(c))
	for i, v := range c {
		value, err := v.Link(scope, s.ValueSpec)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return ConstantSet(values), nil
}

// ConstantList represents a list of constant values from the Thrift file.
type ConstantList []ConstantValue

func compileConstantList(src ast.ConstantList) ConstantList {
	values := make([]ConstantValue, len(src.Items))
	for i, v := range src.Items {
		values[i] = compileConstantValue(v)
	}
	return ConstantList(values)
}

// Link for ConstantList.
func (c ConstantList) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	rt := RootTypeSpec(t)
	if _, isSet := rt.(*SetSpec); isSet {
		return ConstantSet(c).Link(scope, t)
	}

	l, ok := rt.(*ListSpec)
	if !ok {
		return nil, constantValueCastError{Value: c, Type: t}
	}

	values := make([]ConstantValue, len(c))
	for i, v := range c {
		value, err := v.Link(scope, l.ValueSpec)
		if err != nil {
			return nil, err
		}
		values[i] = value
	}

	return ConstantList(values), nil
}

// ConstReference represents a reference to a 'const' declared in the Thrift
// file.
type ConstReference struct {
	Target *Constant
}

// Link for ConstReference.
func (c ConstReference) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	if t == c.Target.Type {
		return c, nil
	}
	return c.Target.Value.Link(scope, t)
}

// EnumItemReference represents a reference to an item of an enum defined in the
// THrift file.
type EnumItemReference struct {
	Enum *EnumSpec
	Item *EnumItem
}

// Link for EnumItemReference.
func (e EnumItemReference) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	if RootTypeSpec(t) != e.Enum {
		return nil, constantValueCastError{Value: e, Type: t}
	}
	return e, nil
}

// constantReference represents a reference to another constant.
//
// This gets resolved to a ConstReference or an EnumItemReference during the
// link stage.
type constantReference ast.ConstantReference

// Link a constantReference.
//
// This resolves the reference to a ConstReference or an EnumItemReference.
func (r constantReference) Link(scope Scope, t TypeSpec) (ConstantValue, error) {
	src := ast.ConstantReference(r)

	c, err := scope.LookupConstant(src.Name)
	if err == nil {
		if err := c.Link(scope); err != nil {
			return nil, err
		}
		return ConstReference{Target: c}.Link(scope, t)
	}

	mname, iname := splitInclude(src.Name)
	if len(mname) == 0 {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	if enum, ok := lookupEnum(scope, mname); ok {
		if item, ok := enum.LookupItem(iname); ok {
			return EnumItemReference{
				Enum: enum,
				Item: item,
			}, nil
		}

		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason: unrecognizedEnumItemError{
				EnumName: mname,
				ItemName: iname,
			},
		}
	}

	includedScope, err := getIncludedScope(scope, mname)
	if err != nil {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	value, err := constantReference{Name: iname}.Link(includedScope, t)
	if err != nil {
		return nil, referenceError{
			Target:    src.Name,
			Line:      src.Line,
			ScopeName: scope.GetName(),
			Reason:    err,
		}
	}

	return value, nil
}

// lookupEnum looks up an enum with the given name in the given scope.
func lookupEnum(scope Scope, name string) (*EnumSpec, bool) {
	t, err := scope.LookupType(name)
	if err != nil {
		return nil, false
	}

	if enum, ok := t.(*EnumSpec); ok {
		return enum, true
	}
	return nil, false
}
