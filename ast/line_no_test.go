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

func TestLineNumber(t *testing.T) {
	tests := []struct {
		give Node
		want int
	}{
		// The following don't record their line numbers because they're just
		// type aliases. We can't change this anymore because converting them
		// to structs is a breaking change.
		{give: ConstantBoolean(true), want: 0},
		{give: ConstantDouble(42.0), want: 0},
		{give: ConstantInteger(42), want: 0},
		{give: ConstantString("foo"), want: 0},

		// Program is the whole Thrift file. It doesn't have a line number.
		{give: &Program{}, want: 0},

		{give: &Annotation{Line: 1}, want: 1},
		{give: ConstantMap{Line: 2}, want: 2},
		{give: ConstantMapItem{Line: 3}, want: 3},
		{give: ConstantList{Line: 4}, want: 4},
		{give: ConstantReference{Line: 5}, want: 5},
		{give: &Constant{Line: 6}, want: 6},
		{give: &Typedef{Line: 7}, want: 7},
		{give: &Enum{Line: 8}, want: 8},
		{give: &EnumItem{Line: 9}, want: 9},
		{give: &Struct{Line: 10}, want: 10},
		{give: &Service{Line: 11}, want: 11},
		{give: &Function{Line: 12}, want: 12},
		{give: &Field{Line: 13}, want: 13},
		{give: &Include{Line: 14}, want: 14},
		{give: &Namespace{Line: 15}, want: 15},
		{give: BaseType{Line: 16}, want: 16},
		{give: MapType{Line: 17}, want: 17},
		{give: ListType{Line: 18}, want: 18},
		{give: SetType{Line: 19}, want: 19},
		{give: TypeReference{Line: 20}, want: 20},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.give), func(t *testing.T) {
			assert.Equal(t, tt.want, LineNumber(tt.give))
		})
	}
}
