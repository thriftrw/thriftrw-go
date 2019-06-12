// Copyright (c) 2019 Uber Technologies, Inc.
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

	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/ptr"
)

type serviceName string

type generateServiceBuilder struct {
	api.GenerateServiceRequest

	importer thriftPackageImporter

	nextModuleID  api.ModuleID
	nextServiceID api.ServiceID

	// ThriftFile -> Module ID
	moduleIDs map[string]api.ModuleID

	// ThriftFile -> Service name -> Service ID
	serviceIDs map[string]map[serviceName]api.ServiceID

	// To ensure there are no duplicates
	rootServices map[api.ServiceID]struct{}
}

func newGenerateServiceBuilder(i thriftPackageImporter) *generateServiceBuilder {
	return &generateServiceBuilder{
		GenerateServiceRequest: api.GenerateServiceRequest{
			RootServices:  make([]api.ServiceID, 0, 10),
			Services:      make(map[api.ServiceID]*api.Service),
			Modules:       make(map[api.ModuleID]*api.Module),
			ThriftRoot:    i.ThriftRoot,
			PackagePrefix: i.ImportPrefix,
		},
		importer:      i,
		nextModuleID:  1,
		nextServiceID: 1,
		moduleIDs:     make(map[string]api.ModuleID),
		serviceIDs:    make(map[string]map[serviceName]api.ServiceID),
		rootServices:  make(map[api.ServiceID]struct{}),
	}
}

func (g *generateServiceBuilder) Build() *api.GenerateServiceRequest {
	return &g.GenerateServiceRequest
}

// AddRootService adds a service as a root service to this request.
func (g *generateServiceBuilder) AddRootService(spec *compile.ServiceSpec) (api.ServiceID, error) {
	id, err := g.addService(spec)
	if err != nil {
		return id, err
	}

	if _, alreadyAdded := g.rootServices[id]; !alreadyAdded {
		g.RootServices = append(g.RootServices, id)
		g.rootServices[id] = struct{}{}
	}

	return id, err
}

// AddModule adds the module for the given Thrift file to the request.
func (g *generateServiceBuilder) AddModule(thriftPath string) (api.ModuleID, error) {
	if id, ok := g.moduleIDs[thriftPath]; ok {
		return id, nil
	}

	id := g.nextModuleID
	g.nextModuleID++
	g.moduleIDs[thriftPath] = id

	importPath, err := g.importer.Package(thriftPath)
	if err != nil {
		return 0, err
	}

	dir, err := g.importer.RelativePackage(thriftPath)
	if err != nil {
		return 0, err
	}

	g.Modules[id] = &api.Module{
		ImportPath:     importPath,
		Directory:      dir,
		ThriftFilePath: thriftPath,
	}
	return id, nil
}

func (g *generateServiceBuilder) addService(spec *compile.ServiceSpec) (api.ServiceID, error) {
	if moduleServices, ok := g.serviceIDs[spec.ThriftFile()]; ok {
		if id, ok := moduleServices[serviceName(spec.Name)]; ok {
			return id, nil
		}
	} else {
		g.serviceIDs[spec.ThriftFile()] = make(map[serviceName]api.ServiceID)
	}

	var parentID *api.ServiceID
	if spec.Parent != nil {
		parent, err := g.addService(spec.Parent)
		if err != nil {
			return 0, err
		}
		parentID = &parent
	}

	serviceID := g.nextServiceID
	g.nextServiceID++
	g.serviceIDs[spec.ThriftFile()][serviceName(spec.Name)] = serviceID

	// Modules must already be populated.
	moduleID, ok := g.moduleIDs[spec.ThriftFile()]
	if !ok {
		return 0, fmt.Errorf("unable to lookup module ID for Thrift file: %q", spec.ThriftFile())
	}

	functions := make([]*api.Function, 0, len(spec.Functions))
	for _, functionName := range sortStringKeys(spec.Functions) {
		function, err := g.buildFunction(spec.Functions[functionName])
		if err != nil {
			return 0, err
		}
		functions = append(functions, function)
	}

	g.Services[serviceID] = &api.Service{
		ThriftName:  spec.Name,
		Name:        goCase(spec.Name),
		ParentID:    parentID,
		Functions:   functions,
		ModuleID:    moduleID,
		Annotations: spec.Annotations,
	}
	return serviceID, nil
}

func (g *generateServiceBuilder) buildFunction(spec *compile.FunctionSpec) (*api.Function, error) {
	args, err := g.buildFieldGroup(compile.FieldGroup(spec.ArgsSpec))
	if err != nil {
		return nil, err
	}

	function := &api.Function{
		Name:        goCase(spec.Name),
		ThriftName:  spec.Name,
		Arguments:   args,
		Annotations: spec.Annotations,
	}
	if spec.OneWay {
		function.OneWay = ptr.Bool(spec.OneWay)
	}

	if spec.ResultSpec != nil {
		var err error
		result := spec.ResultSpec
		if result.ReturnType != nil {
			function.ReturnType, err = g.buildType(result.ReturnType, true)
			if err != nil {
				return nil, err
			}
		}
		if len(result.Exceptions) > 0 {
			function.Exceptions, err = g.buildFieldGroup(result.Exceptions)
			if err != nil {
				return nil, err
			}
		}
	}

	return function, nil
}

func (g *generateServiceBuilder) buildFieldGroup(fs compile.FieldGroup) ([]*api.Argument, error) {
	args := make([]*api.Argument, 0, len(fs))
	for _, f := range fs {
		t, err := g.buildType(f.Type, f.Required)
		if err != nil {
			return nil, err
		}

		name, err := goName(f)
		if err != nil {
			return nil, err
		}
		args = append(args, &api.Argument{
			Name: name,
			Type: t,
		})
	}
	return args, nil
}

func (g *generateServiceBuilder) buildType(spec compile.TypeSpec, required bool) (*api.Type, error) {
	simpleType := func(t api.SimpleType) *api.SimpleType { return &t }

	// try primitives first since they have to be wrapped inside a pointer if
	// optional.
	var t *api.Type
	switch s := spec.(type) {
	case *compile.BoolSpec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeBool)}
	case *compile.I8Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt8)}
	case *compile.I16Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt16)}
	case *compile.I32Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt32)}
	case *compile.I64Spec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeInt64)}
	case *compile.DoubleSpec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeFloat64)}
	case *compile.StringSpec:
		t = &api.Type{SimpleType: simpleType(api.SimpleTypeString)}
	case *compile.EnumSpec:
		importPath, err := g.importer.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}
		name, err := goName(s)
		if err != nil {
			return nil, err
		}
		t = &api.Type{
			ReferenceType: &api.TypeReference{
				Name:        name,
				ImportPath:  importPath,
				Annotations: s.Annotations,
			},
		}
	}

	if t != nil {
		if !required {
			t = &api.Type{PointerType: t}
		}
		return t, nil
	}

	switch s := spec.(type) {
	case *compile.BinarySpec:
		return &api.Type{SliceType: &api.Type{SimpleType: simpleType(api.SimpleTypeByte)}}, nil

	case *compile.MapSpec:
		k, err := g.buildType(s.KeySpec, true)
		if err != nil {
			return nil, err
		}

		v, err := g.buildType(s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		if !isHashable(s.KeySpec) {
			return &api.Type{KeyValueSliceType: &api.TypePair{Left: k, Right: v}}, nil
		}

		return &api.Type{MapType: &api.TypePair{Left: k, Right: v}}, nil

	case *compile.ListSpec:
		v, err := g.buildType(s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		return &api.Type{SliceType: v}, nil

	case *compile.SetSpec:
		v, err := g.buildType(s.ValueSpec, true)
		if err != nil {
			return nil, err
		}

		if !isHashable(s.ValueSpec) {
			return &api.Type{SliceType: v}, nil
		}

		return &api.Type{MapType: &api.TypePair{
			Left:  v,
			Right: &api.Type{SimpleType: simpleType(api.SimpleTypeStructEmpty)},
		}}, nil

	case *compile.StructSpec:
		importPath, err := g.importer.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		name, err := goName(s)
		if err != nil {
			return nil, err
		}

		return &api.Type{
			PointerType: &api.Type{
				ReferenceType: &api.TypeReference{
					Name:        name,
					ImportPath:  importPath,
					Annotations: s.Annotations,
				},
			},
		}, nil

	case *compile.TypedefSpec:
		importPath, err := g.importer.Package(s.ThriftFile())
		if err != nil {
			return nil, err
		}

		name, err := goName(s)
		if err != nil {
			return nil, err
		}

		t = &api.Type{
			ReferenceType: &api.TypeReference{
				Name:        name,
				ImportPath:  importPath,
				Annotations: s.Annotations,
			},
		}

		if (!required && !isReferenceType(spec)) ||
			isStructType(spec) {
			t = &api.Type{PointerType: t}
		}

		return t, nil
	default:
		panic(fmt.Sprintf("Unknown type (%T) %v", spec, spec))
	}
}
