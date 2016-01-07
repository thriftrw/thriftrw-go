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

package gen

import (
	"fmt"
	"go/ast"

	"github.com/uber/thriftrw-go/compile"
)

// TypeDefinition TODO
func (g *Generator) TypeDefinition(spec compile.TypeSpec) {
	switch s := spec.(type) {
	case *compile.EnumSpec:
		g.enum(s)
	case *compile.StructSpec:
		g.structure(s)
	case *compile.TypedefSpec:
		g.typedef(s)
	default:
		panic(fmt.Sprintf("%q is not a defined type", spec.ThriftName()))
	}
}

func isReferenceType(spec compile.TypeSpec) bool {
	if spec == compile.BinarySpec {
		return true
	}

	switch spec.(type) {
	case *compile.MapSpec:
		return true
	case *compile.ListSpec:
		return true
	case *compile.SetSpec:
		return true
	}

	return false
}

func typeReference(spec compile.TypeSpec, ptr bool) (result ast.Expr) {
	if ptr && !isReferenceType(spec) {
		defer func() {
			result = &ast.StarExpr{X: result}
		}()
	}

	switch spec {
	case compile.BoolSpec:
		return ast.NewIdent("bool")
	case compile.I8Spec:
		return ast.NewIdent("int8")
	case compile.I16Spec:
		return ast.NewIdent("int16")
	case compile.I32Spec:
		return ast.NewIdent("int32")
	case compile.I64Spec:
		return ast.NewIdent("int64")
	case compile.DoubleSpec:
		return ast.NewIdent("double64")
	case compile.StringSpec:
		return ast.NewIdent("string")
	case compile.BinarySpec:
		return &ast.ArrayType{Elt: ast.NewIdent("byte")}
	default:
		// Try matching type
	}

	switch s := spec.(type) {
	case *compile.MapSpec:
		// TODO unhashable types
		return &ast.MapType{
			Key:   typeReference(s.KeySpec, false),
			Value: typeReference(s.ValueSpec, false),
		}
	case *compile.ListSpec:
		return &ast.ArrayType{Elt: typeReference(s.ValueSpec, false)}
	case *compile.SetSpec:
		// TODO unhashable types
		return &ast.MapType{
			Key:   typeReference(s.ValueSpec, false),
			Value: &ast.StructType{},
		}
	default:
		return ast.NewIdent(typeDeclName(spec))
	}
}

func typeDeclName(spec compile.TypeSpec) string {
	switch s := spec.(type) {
	case *compile.EnumSpec:
		return capitalize(s.Name)
	case *compile.StructSpec:
		return capitalize(s.Name)
	case *compile.TypedefSpec:
		return capitalize(s.Name)
	default:
		panic(fmt.Sprintf(
			"Type %q doesn't can't have a declaration name", spec.ThriftName(),
		))
	}
}
