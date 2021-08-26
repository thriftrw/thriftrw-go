package compare

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/compile"
)

func TestErrorRequiredCase(t *testing.T) {
	type test struct {
		desc string
		fromStruct *compile.StructSpec
		toStruct *compile.StructSpec
		wantError string
	}
	tests := []test{
		{
			desc: "changed an optional field to required",
			fromStruct: &compile.StructSpec{
				Name: "structA",
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    false,
						Name: "fieldA",
					},
				},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    true,
						Name: "fieldA",
					},
				},
			},
			wantError: "changing an optional field fieldA in structA to required is not backwards compatible",
		},
		{
			desc: "found a new required field",
			fromStruct: &compile.StructSpec{
				Fields:      compile.FieldGroup{},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    true,
						Name: "fieldA",
					},
				},
			},
			wantError: "adding a required field fieldA to structA is not backwards compatible",
		},
		{
			desc: "found a new required and changed optional field",
			fromStruct: &compile.StructSpec{
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    false,
						Name: "fieldA",
					},

				},
			},
			toStruct: &compile.StructSpec{
				Name: "structA",
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    true,
						Name: "fieldA",
					},
					&compile.FieldSpec{
						Required:    true,
						Name: "fieldB",
					},
				},
			},
			wantError: "changing an optional field fieldA in structA to" +
				" required is not backwards compatible; changing an optional" +
				" field fieldB in structA to required is not backwards" +
				" compatible",
		},

	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := StructSpecs(tt.fromStruct, tt.toStruct)
			require.Error(t, err, "expected error")
			assert.EqualError(t, err, tt.wantError, "wrong error message")
		})
	}
}

func TestRequiredCaseOk(t *testing.T) {
	type test struct {
		desc string
		fromStruct *compile.StructSpec
		toStruct *compile.StructSpec
	}
	tests := []test{
		{
			desc: "adding an optional field",
			fromStruct: &compile.StructSpec{
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    false,
						Name: "fieldA",
					},
				},
			},
			toStruct: &compile.StructSpec{
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    false,
						Name: "fieldA",
					},
					&compile.FieldSpec{
						Required:    false,
						Name: "fieldA",
					},
				},
			},
		},
		{
			desc: "removing a field from a struct",
			fromStruct: &compile.StructSpec{
				Fields:      compile.FieldGroup{
					&compile.FieldSpec{
						Required:    false,
					},
				},
			},
			toStruct: &compile.StructSpec{
				Fields:      compile.FieldGroup{
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := StructSpecs(tt.fromStruct, tt.toStruct)
			require.NoError(t, err, "do not expect an error")
		})
	}
}
