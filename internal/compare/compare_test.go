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

package compare

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/compile"
)

func TestErrorRequiredCase(t *testing.T) {
	t.Parallel()
	type test struct {
		desc       string
		fromStruct *compile.StructSpec
		toStruct   *compile.StructSpec
		wantError  string
	}
	tests := []test{
		{
			desc: "changed an optional field to required",
			fromStruct: &compile.StructSpec{
				Name: "structA",
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: false,
						Name:     "fieldA",
					},
				},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: true,
						Name:     "fieldA",
					},
				},
			},
			wantError: `foo.thrift:changing an optional field "fieldA" in "structA" to required`,
		},
		{
			desc: "found a new required field",
			fromStruct: &compile.StructSpec{
				Fields: compile.FieldGroup{},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: true,
						Name:     "fieldA",
					},
				},
			},
			wantError: `foo.thrift:adding a required field "fieldA" to "structA"`,
		},
		{
			desc: "found a new required and changed optional field",
			fromStruct: &compile.StructSpec{
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: false,
						Name:     "fieldA",
					},
				},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: true,
						Name:     "fieldA",
					},
					&compile.FieldSpec{
						Required: true,
						Name:     "fieldB",
					},
				},
			},
			wantError: `foo.thrift:changing an optional field "fieldA" in "structA" to required` +
				"\n" + `foo.thrift:changing an optional field "fieldB" in "structA" to required`,
		},
		{
			desc: "found a type changed",
			fromStruct: &compile.StructSpec{
				Name: "structA",
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Name: "fieldA",
						Type: &compile.BinarySpec{},
					},
				},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Name: "fieldA",
						Type: &compile.BoolSpec{},
					},
					&compile.FieldSpec{
						Name: "fieldB",
						Type: &compile.BinarySpec{},
					},
				},
			},
			wantError: `foo.thrift:changing type of field "fieldA" in struct "structA" from "binary" to "bool"`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			pass := Pass{}
			pass.structSpecs(tt.fromStruct, tt.toStruct, "foo.thrift")
			want := fmt.Sprintf("%s\n", tt.wantError)
			assert.Equal(t, want, pass.String(), "wrong lint diagnostics")
		})
	}
}

func TestRequiredCaseOk(t *testing.T) {
	t.Parallel()
	type test struct {
		desc       string
		fromStruct *compile.StructSpec
		toStruct   *compile.StructSpec
	}
	tests := []test{
		{
			desc: "adding an optional field",
			fromStruct: &compile.StructSpec{
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: false,
						Name:     "fieldA",
					},
				},
			},
			toStruct: &compile.StructSpec{
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: false,
						Name:     "fieldA",
					},
					&compile.FieldSpec{
						Required: false,
						Name:     "fieldA",
					},
				},
			},
		},
		{
			desc: "removing a field from a struct",
			fromStruct: &compile.StructSpec{
				Fields: compile.FieldGroup{
					&compile.FieldSpec{
						Required: false,
					},
				},
			},
			toStruct: &compile.StructSpec{
				Fields: compile.FieldGroup{},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			pass := Pass{}
			pass.structSpecs(tt.fromStruct, tt.toStruct, "foo.thrift")
			assert.Equal(t, "", pass.String(), "wrong lint diagnostics")
		})
	}
}

func TestServicesError(t *testing.T) {
	t.Parallel()
	type test struct {
		desc      string
		from      *compile.ServiceSpec
		to        *compile.ServiceSpec
		wantError string
	}
	tests := []test{
		{
			desc: "removing service",
			from: &compile.ServiceSpec{
				Name: "foo",
				File: "foo.thrift",
			},
			to:        nil,
			wantError: `foo.thrift:deleting service "foo"`,
		},
		{
			desc: "removing a method",
			from: &compile.ServiceSpec{
				Name:      "foo",
				File:      "foo.thrift",
				Functions: map[string]*compile.FunctionSpec{"bar": {}},
			},
			to: &compile.ServiceSpec{
				Name: "foo",
				File: "foo.thrift",
			},
			wantError: `foo.thrift:removing method "bar" in service "foo"`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			pass := Pass{}
			pass.service(tt.from, tt.to)
			want := fmt.Sprintf("%s\n", tt.wantError)
			assert.Equal(t, want, pass.String(), "wrong error message")
		})
	}
}

func TestRelativePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		desc string
		path string
		want string
	}{
		{
			desc: "good case",
			path: "/home/foo/bar/buz/buz.go",
			want: "buz/buz.go",
		},
		{
			desc: "fallback case where we just return file name",
			path: "foo/buz.go",
			want: "buz.go",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			p := Pass{
				GitDir: "/home/foo/bar",
			}
			assert.Equal(t, tt.want, p.getRelativePath(tt.path))
		})
	}

}
