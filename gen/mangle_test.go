package gen

import (
	"testing"

	"go.uber.org/thriftrw/compile"

	"github.com/stretchr/testify/assert"
)

func TestMangleType(t *testing.T) {
	m := newMangler()

	tests := []struct {
		spec compile.TypeSpec
		want string
	}{
		{
			spec: &compile.MapSpec{
				KeySpec: &compile.StringSpec{},
				ValueSpec: &compile.ListSpec{
					ValueSpec: &compile.I32Spec{},
				},
			},
			want: "Map_String_List_I32",
		},
		{
			spec: &compile.SetSpec{ValueSpec: &compile.StringSpec{}},
			want: "Set_String",
		},
		{
			spec: &compile.StructSpec{Name: "foo", File: "bar.thrift"},
			want: "Foo",
		},
		{
			spec: &compile.StructSpec{Name: "foo", File: "baz.thrift"},
			want: "Foo_1",
		},
		{
			spec: &compile.MapSpec{
				KeySpec: &compile.TypedefSpec{Name: "UUID", File: "users.thrift"},
				ValueSpec: &compile.MapSpec{
					KeySpec: &compile.TypedefSpec{Name: "UUID", File: "common.thrift"},
					ValueSpec: &compile.StructSpec{
						Name: "User",
						File: "users.thrift",
					},
				},
			},
			want: "Map_UUID_Map_UUID_1_User",
		},
	}

	for _, tt := range tests {
		t.Run(tt.spec.ThriftName(), func(t *testing.T) {
			assert.Equal(t, tt.want, m.MangleType(tt.spec))
		})
	}
}
