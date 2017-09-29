package gen

import (
	"testing"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/ptr"

	"github.com/stretchr/testify/assert"
)

func TestAddRootService(t *testing.T) {
	tests := []struct {
		desc string
		spec *compile.ServiceSpec
		want *api.GenerateServiceRequest
	}{
		{
			desc: "empty service",
			spec: &compile.ServiceSpec{
				Name: "EmptyService",
				File: "idl/empty.thrift",
			},
			want: &api.GenerateServiceRequest{
				RootServices: []api.ServiceID{1},
				Services: map[api.ServiceID]*api.Service{
					1: {
						Name:       "EmptyService",
						ThriftName: "EmptyService",
						Functions:  []*api.Function{}, // must be non-nil
						ModuleID:   1,
					},
				},
				Modules: map[api.ModuleID]*api.Module{
					1: {
						ImportPath: "go.uber.org/thriftrw/gen/testdata/empty",
						Directory:  "empty",
					},
				},
			},
		},
		{
			desc: "Non standard names",
			spec: &compile.ServiceSpec{
				Name: "non_standard_service_name",
				File: "idl/service.thrift",
			},
			want: &api.GenerateServiceRequest{
				RootServices: []api.ServiceID{1},
				Services: map[api.ServiceID]*api.Service{
					1: {
						Name:       "NonStandardServiceName",
						ThriftName: "non_standard_service_name",
						Functions:  []*api.Function{}, // must be non-nil
						ModuleID:   1,
					},
				},
				Modules: map[api.ModuleID]*api.Module{
					1: {
						ImportPath: "go.uber.org/thriftrw/gen/testdata/service",
						Directory:  "service",
					},
				},
			},
		},
		{
			desc: "service with a parent",
			spec: &compile.ServiceSpec{
				Name: "KeyValue",
				File: "idl/kv.thrift",
				Parent: &compile.ServiceSpec{
					Name: "AbstractService",
					File: "idl/common/abstract.thrift",
				},
			},
			want: &api.GenerateServiceRequest{
				RootServices: []api.ServiceID{2},
				Services: map[api.ServiceID]*api.Service{
					1: {
						Name:       "AbstractService",
						ThriftName: "AbstractService",
						Functions:  []*api.Function{}, // must be non-nil
						ModuleID:   1,
					},
					2: {
						Name:       "KeyValue",
						ThriftName: "KeyValue",
						ParentID:   (*api.ServiceID)(ptr.Int32(1)),
						Functions:  []*api.Function{}, // must be non-nil
						ModuleID:   2,
					},
				},
				Modules: map[api.ModuleID]*api.Module{
					1: {
						ImportPath: "go.uber.org/thriftrw/gen/testdata/common/abstract",
						Directory:  "common/abstract",
					},
					2: {
						ImportPath: "go.uber.org/thriftrw/gen/testdata/kv",
						Directory:  "kv",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		spec := tt.spec
		err := spec.Link(compile.EmptyScope("foo"))
		if !assert.NoError(t, err, "%v: invalid test: scope must link", tt.desc) {
			continue
		}

		importer := thriftPackageImporter{
			ImportPrefix: "go.uber.org/thriftrw/gen/testdata",
			ThriftRoot:   "idl",
		}

		g := newGenerateServiceBuilder(importer)
		if _, err := g.AddRootService(spec); assert.NoError(t, err, tt.desc) {
			assert.Equal(t, tt.want, g.Build(), tt.desc)
		}
	}
}

func TestBuildFunction(t *testing.T) {
	tests := []struct {
		desc string
		spec *compile.FunctionSpec
		want *api.Function
	}{
		{
			desc: "returns and throws",
			spec: &compile.FunctionSpec{
				Name: "getValue",
				ArgsSpec: compile.ArgsSpec{
					{
						ID:   1,
						Name: "key",
						Type: &compile.StringSpec{},
					},
				},
				ResultSpec: &compile.ResultSpec{
					ReturnType: &compile.BinarySpec{},
					Exceptions: compile.FieldGroup{
						{
							ID:   1,
							Name: "doesNotExist",
							Type: &compile.StructSpec{
								Name: "KeyDoesNotExist",
								File: "idl/keyvalue.thrift",
								Type: ast.ExceptionType,
								Fields: compile.FieldGroup{
									{
										ID:   1,
										Name: "message",
										Type: &compile.StringSpec{},
									},
								},
							},
						},
					},
				},
			},
			want: &api.Function{
				Name:       "GetValue",
				ThriftName: "getValue",
				Arguments: []*api.Argument{
					{
						Name: "Key",
						Type: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeString)}},
					},
				},
				ReturnType: &api.Type{SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)}},
				Exceptions: []*api.Argument{
					{
						Name: "DoesNotExist",
						Type: &api.Type{
							PointerType: &api.Type{
								ReferenceType: &api.TypeReference{
									Name:       "KeyDoesNotExist",
									ImportPath: "go.uber.org/thriftrw/gen/testdata/keyvalue",
								},
							},
						},
					},
				},
			},
		},
		{
			desc: "no return, no throw",
			spec: &compile.FunctionSpec{
				Name: "setValue",
				ArgsSpec: compile.ArgsSpec{
					{
						ID:   1,
						Name: "key",
						Type: &compile.StringSpec{},
					},
					{
						ID:   2,
						Name: "value",
						Type: &compile.BinarySpec{},
					},
				},
				ResultSpec: &compile.ResultSpec{},
			},
			want: &api.Function{
				Name:       "SetValue",
				ThriftName: "setValue",
				Arguments: []*api.Argument{
					{
						Name: "Key",
						Type: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeString)}},
					},
					{
						Name: "Value",
						Type: &api.Type{SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)}},
					},
				},
			},
		},
		{
			desc: "oneway",
			spec: &compile.FunctionSpec{
				Name:   "clearCache",
				OneWay: true,
				ArgsSpec: compile.ArgsSpec{
					{
						ID:   1,
						Name: "delayMS",
						Type: &compile.I64Spec{},
					},
				},
			},
			want: &api.Function{
				Name:       "ClearCache",
				ThriftName: "clearCache",
				OneWay:     ptr.Bool(true),
				Arguments: []*api.Argument{
					{
						Name: "DelayMS",
						Type: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeInt64)}},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		spec := tt.spec
		err := spec.Link(compile.EmptyScope("foo"))
		if !assert.NoError(t, err, "%v: invalid test: scope must link", tt.desc) {
			continue
		}

		importer := thriftPackageImporter{
			ImportPrefix: "go.uber.org/thriftrw/gen/testdata",
			ThriftRoot:   "idl",
		}

		g := newGenerateServiceBuilder(importer)
		got, err := g.buildFunction(spec)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, tt.want, got, tt.desc)
		}
	}
}

func TestBuildType(t *testing.T) {
	tests := []struct {
		desc     string
		spec     compile.TypeSpec
		required bool

		want *api.Type
	}{
		// required primitives
		{
			desc:     "bool",
			spec:     &compile.BoolSpec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeBool)},
		},
		{
			desc:     "int8",
			spec:     &compile.I8Spec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeInt8)},
		},
		{
			desc:     "int16",
			spec:     &compile.I16Spec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeInt16)},
		},
		{
			desc:     "int32",
			spec:     &compile.I32Spec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeInt32)},
		},
		{
			desc:     "int64",
			spec:     &compile.I64Spec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeInt64)},
		},
		{
			desc:     "float64",
			spec:     &compile.DoubleSpec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeFloat64)},
		},
		{
			desc:     "string",
			spec:     &compile.StringSpec{},
			required: true,
			want:     &api.Type{SimpleType: simpleType(api.SimpleTypeString)},
		},
		{
			desc:     "[]byte",
			spec:     &compile.BinarySpec{},
			required: true,
			want:     &api.Type{SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)}},
		},

		// optional primitives
		{
			desc: "*bool",
			spec: &compile.BoolSpec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeBool)}},
		},
		{
			desc: "*int8",
			spec: &compile.I8Spec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeInt8)}},
		},
		{
			desc: "*int16",
			spec: &compile.I16Spec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeInt16)}},
		},
		{
			desc: "*int32",
			spec: &compile.I32Spec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeInt32)}},
		},
		{
			desc: "*int64",
			spec: &compile.I64Spec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeInt64)}},
		},
		{
			desc: "*float64",
			spec: &compile.DoubleSpec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeFloat64)}},
		},
		{
			desc: "*string",
			spec: &compile.StringSpec{},
			want: &api.Type{PointerType: &api.Type{SimpleType: simpleType(api.SimpleTypeString)}},
		},
		{
			desc: "[]byte",
			spec: &compile.BinarySpec{},
			want: &api.Type{SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)}},
		},

		// containers
		{
			// hashable map key
			desc: "map[string]int32",
			spec: &compile.MapSpec{
				KeySpec:   &compile.StringSpec{},
				ValueSpec: &compile.I32Spec{},
			},
			want: &api.Type{MapType: &api.TypePair{
				Left:  &api.Type{SimpleType: simpleType(api.SimpleTypeString)},
				Right: &api.Type{SimpleType: simpleType(api.SimpleTypeInt32)},
			}},
		},
		{
			// unhashable map key
			desc: "[]struct{Key []byte; Value int32}",
			spec: &compile.MapSpec{
				KeySpec:   &compile.BinarySpec{},
				ValueSpec: &compile.I32Spec{},
			},
			want: &api.Type{KeyValueSliceType: &api.TypePair{
				Left: &api.Type{
					SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)},
				},
				Right: &api.Type{SimpleType: simpleType(api.SimpleTypeInt32)},
			}},
		},
		{
			// hashable set item
			desc: "map[float64]struct{}",
			spec: &compile.SetSpec{ValueSpec: &compile.DoubleSpec{}},
			want: &api.Type{MapType: &api.TypePair{
				Left:  &api.Type{SimpleType: simpleType(api.SimpleTypeFloat64)},
				Right: &api.Type{SimpleType: simpleType(api.SimpleTypeStructEmpty)},
			}},
		},
		{
			// unhashable set item
			desc: "[]*foo.Foo",
			spec: &compile.SetSpec{
				ValueSpec: &compile.StructSpec{
					Name: "Foo",
					File: "idl/foo.thrift",
					Type: ast.StructType,
					Fields: compile.FieldGroup{
						{
							ID:       1,
							Name:     "value",
							Type:     &compile.StringSpec{},
							Required: true,
						},
					},
				},
			},
			want: &api.Type{
				SliceType: &api.Type{
					PointerType: &api.Type{
						ReferenceType: &api.TypeReference{
							Name:       "Foo",
							ImportPath: "go.uber.org/thriftrw/gen/testdata/foo",
						},
					},
				},
			},
		},
		{
			// list
			desc: "[]map[string][]byte",
			spec: &compile.ListSpec{
				ValueSpec: &compile.MapSpec{
					KeySpec:   &compile.StringSpec{},
					ValueSpec: &compile.BinarySpec{},
				},
			},
			want: &api.Type{
				SliceType: &api.Type{
					MapType: &api.TypePair{
						Left: &api.Type{SimpleType: simpleType(api.SimpleTypeString)},
						Right: &api.Type{
							SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)},
						},
					},
				},
			},
		},
		{
			// required enum
			desc: "required enum",
			spec: &compile.EnumSpec{
				Name: "Foo",
				File: "idl/bar.thrift",
				Items: []compile.EnumItem{
					{Name: "A", Value: 0},
					{Name: "B", Value: 2},
				},
			},
			required: true,
			want: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/gen/testdata/bar",
				},
			},
		},
		{
			// optional enum
			desc: "optional enum",
			spec: &compile.EnumSpec{
				Name: "Foo",
				File: "idl/bar.thrift",
				Items: []compile.EnumItem{
					{Name: "A", Value: 0},
					{Name: "B", Value: 2},
				},
			},
			want: &api.Type{
				PointerType: &api.Type{
					ReferenceType: &api.TypeReference{
						Name:       "Foo",
						ImportPath: "go.uber.org/thriftrw/gen/testdata/bar",
					},
				},
			},
		},
		{
			desc: "enum with annotations",
			spec: &compile.EnumSpec{
				Name: "Foo",
				File: "idl/bar.thrift",
				Items: []compile.EnumItem{
					{Name: "A", Value: 0},
					{Name: "B", Value: 2},
				},
				Annotations: compile.Annotations{
					"foo": "bar",
					"baz": "",
				},
			},
			required: true,
			want: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/gen/testdata/bar",
					Annotations: map[string]string{
						"foo": "bar",
						"baz": "",
					},
				},
			},
		},
		{
			// struct
			desc: "struct",
			spec: &compile.StructSpec{
				Name: "Foo",
				File: "idl/foo.thrift",
				Type: ast.StructType,
				Fields: compile.FieldGroup{
					{
						ID:       1,
						Name:     "value",
						Type:     &compile.StringSpec{},
						Required: true,
					},
				},
			},
			want: &api.Type{
				PointerType: &api.Type{
					ReferenceType: &api.TypeReference{
						Name:       "Foo",
						ImportPath: "go.uber.org/thriftrw/gen/testdata/foo",
					},
				},
			},
		},
		{
			desc: "struct with annotations",
			spec: &compile.StructSpec{
				Name: "Foo",
				File: "idl/foo.thrift",
				Type: ast.StructType,
				Fields: compile.FieldGroup{
					{
						ID:       1,
						Name:     "value",
						Type:     &compile.StringSpec{},
						Required: true,
					},
				},
				Annotations: compile.Annotations{
					"validate":  "true",
					"obfuscate": "",
				},
			},
			want: &api.Type{
				PointerType: &api.Type{
					ReferenceType: &api.TypeReference{
						Name:       "Foo",
						ImportPath: "go.uber.org/thriftrw/gen/testdata/foo",
						Annotations: map[string]string{
							"validate":  "true",
							"obfuscate": "",
						},
					},
				},
			},
		},
		{
			desc: "required typedef with a primitive",
			spec: &compile.TypedefSpec{
				Name:   "Foo",
				File:   "idl/foo/bar.thrift",
				Target: &compile.I64Spec{},
			},
			required: true,
			want: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/gen/testdata/foo/bar",
				},
			},
		},
		{
			desc: "optional typedef with a primitive",
			spec: &compile.TypedefSpec{
				Name:   "Foo",
				File:   "idl/foo/bar.thrift",
				Target: &compile.I64Spec{},
			},
			want: &api.Type{
				PointerType: &api.Type{
					ReferenceType: &api.TypeReference{
						Name:       "Foo",
						ImportPath: "go.uber.org/thriftrw/gen/testdata/foo/bar",
					},
				},
			},
		},
		{
			desc: "required typedef with non-primitive",
			spec: &compile.TypedefSpec{
				Name:   "Foo",
				File:   "idl/foo/bar.thrift",
				Target: &compile.BinarySpec{},
			},
			required: true,
			want: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/gen/testdata/foo/bar",
				},
			},
		},
		{
			desc: "optional typedef with non-primitive",
			spec: &compile.TypedefSpec{
				Name:   "Foo",
				File:   "idl/foo/bar.thrift",
				Target: &compile.ListSpec{ValueSpec: &compile.StringSpec{}},
			},
			want: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Foo",
					ImportPath: "go.uber.org/thriftrw/gen/testdata/foo/bar",
				},
			},
		},
		{
			desc: "typedef with annotations",
			spec: &compile.TypedefSpec{
				Name:   "Timestamp",
				File:   "idl/common.thrift",
				Target: &compile.ListSpec{ValueSpec: &compile.StringSpec{}},
				Annotations: compile.Annotations{
					"format":   "ISO8601",
					"validate": "true",
				},
			},
			want: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:       "Timestamp",
					ImportPath: "go.uber.org/thriftrw/gen/testdata/common",
					Annotations: map[string]string{
						"format":   "ISO8601",
						"validate": "true",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		spec, err := tt.spec.Link(compile.EmptyScope("foo"))
		if !assert.NoError(t, err, "%v: invalid test: scope must link", tt.desc) {
			continue
		}

		importer := thriftPackageImporter{
			ImportPrefix: "go.uber.org/thriftrw/gen/testdata",
			ThriftRoot:   "idl",
		}

		g := newGenerateServiceBuilder(importer)
		got, err := g.buildType(spec, tt.required)
		if assert.NoError(t, err, tt.desc) {
			assert.Equal(t, tt.want, got, tt.desc)
		}
	}
}

func simpleType(s api.SimpleType) *api.SimpleType {
	return &s
}
