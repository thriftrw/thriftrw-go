// Copyright (c) 2015 Uber Technologies, Inc.
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
	"fmt"

	"github.com/uber/thriftrw-go/ast"
)

// ConstantValue represents a compiled constant value or a reference to one.
type ConstantValue interface {
	Link(scope Scope) (ConstantValue, error)

	// TODO link probably needs a reference to the type the constant is expected
	// to be.

	// TODO all constants must have a type associated with them.
}

// compileConstantValue compiles a constant value AST into a ConstantValue.
func compileConstantValue(v ast.ConstantValue) ConstantValue {
	if v == nil {
		return nil
	}
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
func (c ConstantBool) Link(scope Scope) (ConstantValue, error) {
	return c, nil
}

// Link for ConstantInt.
func (c ConstantInt) Link(scope Scope) (ConstantValue, error) {
	// TODO ConstantInt can resolve to a ConstantBool if it's being treated as
	// a bool.
	return c, nil
}

// Link for ConstantString.
func (c ConstantString) Link(scope Scope) (ConstantValue, error) {
	return c, nil
}

// Link for ConstantDouble.
func (c ConstantDouble) Link(scope Scope) (ConstantValue, error) {
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
func (c ConstantMap) Link(scope Scope) (ConstantValue, error) {
	items := make([]ConstantValuePair, len(c))

	// TODO ConstantMap can resolve into a constant struct if the type it is
	// being cast to is a struct. Otherwise, all keys and values must be the
	// same type.

	for i, item := range c {
		key, err := item.Key.Link(scope)
		if err != nil {
			return nil, err
		}

		value, err := item.Value.Link(scope)
		if err != nil {
			return nil, err
		}

		// TODO(abg): Duplicate key check
		items[i] = ConstantValuePair{Key: key, Value: value}
	}

	return ConstantMap(items), nil
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
func (c ConstantList) Link(scope Scope) (ConstantValue, error) {
	values := make([]ConstantValue, len(c))

	// TODO ConstantList can resolve to a constant set if it's being treated as
	// such.

	for i, v := range c {
		value, err := v.Link(scope)
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
func (c ConstReference) Link(scope Scope) (ConstantValue, error) {
	// The reference has already been verified. Nothing to do.
	return c, nil
}

// EnumItemReference represents a reference to an item of an enum defined in the
// THrift file.
type EnumItemReference struct {
	Enum *EnumSpec
	Item EnumItem
}

// Link for EnumItemReference.
func (e EnumItemReference) Link(scope Scope) (ConstantValue, error) {
	// The reference has already been verified. Nothing to do.
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
func (r constantReference) Link(scope Scope) (ConstantValue, error) {
	src := ast.ConstantReference(r)

	c, err := scope.LookupConstant(src.Name)
	if err == nil {
		err = c.Link(scope)
		return ConstReference{Target: c}, err
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

	value, err := constantReference{Name: iname}.Link(includedScope)
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
