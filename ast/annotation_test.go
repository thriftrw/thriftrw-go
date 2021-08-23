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

package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnotations(t *testing.T) {
	anns := []*Annotation{
		{Name: "annotation", Value: "value"},
	}

	tests := []struct {
		give Node
		want []*Annotation
	}{
		{give: &Annotation{}, want: nil},
		{give: BaseType{Annotations: anns}, want: anns},
		{give: &Constant{}, want: nil},
		{give: ConstantBoolean(true), want: nil},
		{give: ConstantDouble(42.0), want: nil},
		{give: ConstantInteger(42), want: nil},
		{give: ConstantMap{Line: 2, Column: 1}, want: nil},
		{give: ConstantMapItem{Line: 3, Column: 1}, want: nil},
		{give: ConstantList{Line: 4, Column: 1}, want: nil},
		{give: ConstantReference{Line: 5, Column: 1}, want: nil},
		{give: ConstantString("foo"), want: nil},
		{give: &CppInclude{}, want: nil},
		{give: &Enum{Annotations: anns}, want: anns},
		{give: &EnumItem{Annotations: anns}, want: anns},
		{give: &Field{Annotations: anns}, want: anns},
		{give: &Function{Annotations: anns}, want: anns},
		{give: &Include{}, want: nil},
		{give: ListType{Annotations: anns}, want: anns},
		{give: MapType{Annotations: anns}, want: anns},
		{give: &Namespace{}, want: nil},
		{give: &Program{}, want: nil},
		{give: &Service{Annotations: anns}, want: anns},
		{give: SetType{Annotations: anns}, want: anns},
		{give: &Struct{Annotations: anns}, want: anns},
		{give: &Typedef{Annotations: anns}, want: anns},
		{give: TypeReference{}, want: nil},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.give), func(t *testing.T) {
			actual := Annotations(tt.give)
			assert.Equal(t, tt.want, actual)
		})
	}
}

func TestFormatAnnotations(t *testing.T) {
	var anns []*Annotation
	assert.Equal(t, "", FormatAnnotations(anns))

	anns = append(anns, &Annotation{Name: "foo", Value: "bar"})
	assert.Equal(t, `(foo = "bar")`, FormatAnnotations(anns))

	anns = append(anns, &Annotation{Name: "baz", Value: "qux"})
	assert.Equal(t, `(foo = "bar", baz = "qux")`, FormatAnnotations(anns))
}
