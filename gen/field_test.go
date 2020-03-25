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
	"strings"
	"testing"

	"go.uber.org/thriftrw/compile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFieldLabelConflict(t *testing.T) {
	tests := []struct {
		desc    string
		spec    compile.FieldGroup
		wantErr string
	}{
		{
			desc: "conflict with field name",
			spec: compile.FieldGroup{
				{
					Name:        "foo",
					Annotations: compile.Annotations{"go.label": "bar"},
				},
				{
					Name: "bar",
				},
			},
			wantErr: `field "bar" with label "bar" conflicts with field "foo"`,
		},
		{
			desc: "conflict with another label",
			spec: compile.FieldGroup{
				{
					Name:        "foo",
					Annotations: compile.Annotations{"go.label": "baz"},
				},
				{
					Name:        "bar",
					Annotations: compile.Annotations{"go.label": "baz"},
				},
			},
			wantErr: `field "bar" with label "baz" conflicts with field "foo"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			fg := fieldGroupGenerator{
				Namespace: NewNamespace(),
				Fields:    tt.spec,
			}
			err := fg.Generate(nil /* generator */)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestCompileJSONTag(t *testing.T) {
	tests := []struct {
		desc          string
		required      bool
		options       []string
		wantOmitempty bool
	}{
		{
			desc:          "optional fields",
			required:      false,
			wantOmitempty: true,
		},
		{
			desc:          "required fields",
			required:      true,
			wantOmitempty: false,
		},
		{
			desc:          "optional with !omitempty",
			required:      false,
			options:       []string{notOmitempty},
			wantOmitempty: false,
		},
		{
			desc: "required with !omitempty",
			required: true,
			options: []string{notOmitEmpty},
			wantOmitEmpty: false,
		},
	}

	fieldName := "numbers"

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			goTag := fmt.Sprintf(
				`json:"%s,%s"`,
				fieldName,
				strings.Join(tt.options, ","),
			)

			fieldSpec := &compile.FieldSpec{
				ID:   0,
				Name: fieldName,
				Type: &compile.ListSpec{
					ValueSpec: &compile.I64Spec{},
				},
				Required: tt.required,
				Annotations: compile.Annotations{
					goTagKey: goTag,
				},
			}

			result := compileJSONTag(fieldSpec, fieldName, tt.options...)
			require.Equal(t, tt.wantOmitempty, result.HasOption(omitempty))
		})
	}
}
