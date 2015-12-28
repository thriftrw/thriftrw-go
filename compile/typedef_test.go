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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uber/thriftrw-go/ast"
	"github.com/uber/thriftrw-go/idl"
)

func parseTypedef(s string) *ast.Typedef {
	prog, err := idl.Parse([]byte(s))
	if err != nil {
		panic(fmt.Sprintf("failure to parse: %v: %s", err, s))
	}

	if len(prog.Definitions) != 1 {
		panic("parseTypedef may be used to parse a single typedef only")
	}

	return prog.Definitions[0].(*ast.Typedef)
}

func TestCompileTypedef(t *testing.T) {
	tests := []struct {
		src   string
		scope Scope
		spec  *TypedefSpec
	}{
		{
			"typedef i64 timestamp",
			nil,
			&TypedefSpec{Name: "timestamp", Target: I64Spec},
		},
		{
			"typedef Bar Foo",
			scope("Bar", &TypedefSpec{Name: "Bar", Target: I32Spec}),
			&TypedefSpec{
				Name:   "Foo",
				Target: &TypedefSpec{Name: "Bar", Target: I32Spec},
			},
		},
	}

	for _, tt := range tests {
		src := parseTypedef(tt.src)
		typedefSpec := compileTypedef(src)

		scp := tt.scope
		if scp == nil {
			scp = scope()
		}

		expected := mustLink(t, tt.spec, scope())
		spec, err := typedefSpec.Link(scp)
		if assert.NoError(t, err) {
			assert.Equal(t, expected, spec)
		}
	}
}
