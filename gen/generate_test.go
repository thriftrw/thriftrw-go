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
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/plugin"
	"go.uber.org/thriftrw/internal/plugin/handletest"
	"go.uber.org/thriftrw/plugin/api"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testdata(t *testing.T, paths ...string) string {
	// We need an absolute path to CWD.
	cwd, err := os.Getwd()
	require.NoError(t, err, "could not determine CWD")

	args := []string{cwd, "testdata"}
	args = append(args, paths...)
	return filepath.Join(args...)
}

func TestGenerateWithRelativePaths(t *testing.T) {
	outputDir, err := ioutil.TempDir("", "thriftrw-generate-test")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	thriftRoot, err := os.Getwd()
	require.NoError(t, err)

	module, err := compile.Compile("testdata/thrift/structs.thrift")
	require.NoError(t, err)

	opts := []*Options{
		{
			OutputDir:     outputDir,
			PackagePrefix: "go.uber.org/thriftrw/gen",
			ThriftRoot:    "testdata",
		},
		{
			OutputDir:     "testdata",
			PackagePrefix: "go.uber.org/thriftrw/gen",
			ThriftRoot:    thriftRoot,
		},
	}

	for _, opt := range opts {
		err := Generate(module, opt)
		if assert.Error(t, err, "expected code generation with %v to fail", opt) {
			assert.Contains(t, err.Error(), "must be an absolute path")
		}
	}
}

func TestGenerate(t *testing.T) {
	var (
		ts compile.TypeSpec = &compile.TypedefSpec{
			Name:   "Timestamp",
			File:   testdata(t, "thrift/common/bar.thrift"),
			Target: &compile.I64Spec{},
		}
		ts2 compile.TypeSpec = &compile.TypedefSpec{
			Name:   "Timestamp",
			File:   testdata(t, "thrift/foo.thrift"),
			Target: ts,
		}
	)

	ts2, err := ts2.Link(compile.EmptyScope("bar"))
	require.NoError(t, err)

	ts, err = ts.Link(compile.EmptyScope("bar"))
	require.NoError(t, err)

	module := &compile.Module{
		Name:       "foo",
		ThriftPath: testdata(t, "thrift/foo.thrift"),
		Includes: map[string]*compile.IncludedModule{
			"bar": {
				Name: "bar",
				Module: &compile.Module{
					Name:       "bar",
					ThriftPath: testdata(t, "thrift/common/bar.thrift"),
					Types:      map[string]compile.TypeSpec{"Timestamp": ts},
				},
			},
		},
		Types: map[string]compile.TypeSpec{"Timestamp": ts2},
	}

	tests := []struct {
		desc      string
		noRecurse bool
		getPlugin func(*gomock.Controller) plugin.Handle

		wantFiles []string
		wantError string
	}{
		{
			desc:      "nil plugin; no recurse",
			noRecurse: true,
			wantFiles: []string{"foo/types.go"},
		},
		{
			desc: "nil plugin; recurse",
			wantFiles: []string{
				"foo/types.go",
				"common/bar/types.go",
			},
		},
		{
			desc: "no service generator",
			getPlugin: func(mockCtrl *gomock.Controller) plugin.Handle {
				handle := handletest.NewMockHandle(mockCtrl)
				handle.EXPECT().ServiceGenerator().Return(nil)
				return handle
			},
			wantFiles: []string{
				"foo/types.go",
				"common/bar/types.go",
			},
		},
		{
			desc: "empty plugin",
			getPlugin: func(mockCtrl *gomock.Controller) plugin.Handle {
				return plugin.EmptyHandle
			},
			wantFiles: []string{
				"foo/types.go",
				"common/bar/types.go",
			},
		},
		{
			desc: "ServiceGenerator plugin",
			getPlugin: func(mockCtrl *gomock.Controller) plugin.Handle {
				sgen := handletest.NewMockServiceGenerator(mockCtrl)
				sgen.EXPECT().Generate(gomock.Any()).
					Return(&api.GenerateServiceResponse{
						Files: map[string][]byte{
							"foo.txt":    []byte("hello world\n"),
							"bar/baz.go": []byte("package bar\n"),
						},
					}, nil)

				handle := handletest.NewMockHandle(mockCtrl)
				handle.EXPECT().ServiceGenerator().Return(sgen)
				return handle
			},
			wantFiles: []string{
				"foo/types.go",
				"common/bar/types.go",
				"foo.txt",
				"bar/baz.go",
			},
		},
		{
			desc: "ServiceGenerator plugin conflict",
			getPlugin: func(mockCtrl *gomock.Controller) plugin.Handle {
				sgen := handletest.NewMockServiceGenerator(mockCtrl)
				sgen.EXPECT().Generate(gomock.Any()).
					Return(&api.GenerateServiceResponse{
						Files: map[string][]byte{
							"common/bar/types.go": []byte("hulk smash"),
						},
					}, nil)

				handle := handletest.NewMockHandle(mockCtrl)
				handle.EXPECT().ServiceGenerator().Return(sgen)
				return handle
			},
			wantError: `file generation conflict: multiple sources are trying to write to "common/bar/types.go"`,
		},
		{
			desc: "ServiceGenerator plugin error",
			getPlugin: func(mockCtrl *gomock.Controller) plugin.Handle {
				sgen := handletest.NewMockServiceGenerator(mockCtrl)
				sgen.EXPECT().Generate(gomock.Any()).Return(nil, errors.New("great sadness"))

				handle := handletest.NewMockHandle(mockCtrl)
				handle.EXPECT().ServiceGenerator().Return(sgen)
				return handle
			},
			wantError: `great sadness`,
		},
	}

	for _, tt := range tests {
		func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			outputDir, err := ioutil.TempDir(os.TempDir(), "test-generate-recurse")
			require.NoError(t, err)
			defer os.RemoveAll(outputDir)

			var p plugin.Handle
			if tt.getPlugin != nil {
				p = tt.getPlugin(mockCtrl)
			}

			err = Generate(module, &Options{
				OutputDir:     outputDir,
				PackagePrefix: "go.uber.org/thriftrw/gen/testdata",
				ThriftRoot:    testdata(t, "thrift"),
				Plugin:        p,
				NoRecurse:     tt.noRecurse,
			})
			if tt.wantError != "" {
				assert.Contains(t, err.Error(), tt.wantError)
				return
			}

			if assert.NoError(t, err, tt.desc) {
				for _, f := range tt.wantFiles {
					_, err = os.Stat(filepath.Join(outputDir, f))
					assert.NoError(t, err, tt.desc)
				}
			}
		}()
	}
}

func TestThriftPackageImporter(t *testing.T) {
	importer := thriftPackageImporter{
		ImportPrefix: "github.com/myteam/myservice",
		ThriftRoot:   "/src/thrift",
	}

	tests := []struct {
		File, ServiceName string // Inputs

		// If non-empty, these are the expected outputs for RelativePackage,
		// Package
		Relative, Package string
	}{
		{
			File:        "/src/thrift/foo.thrift",
			Relative:    "foo",
			Package:     "github.com/myteam/myservice/foo",
			ServiceName: "MyService",
		},
		{
			File:     "/src/thrift/shared/common.thrift",
			Relative: "shared/common",
			Package:  "github.com/myteam/myservice/shared/common",
		},
	}

	for _, tt := range tests {
		if tt.Relative != "" {
			got, err := importer.RelativePackage(tt.File)
			if assert.NoError(t, err, "RelativePackage(%q)", tt.File) {
				assert.Equal(t, tt.Relative, got, "RelativePackage(%q)", tt.File)
			}
		}

		if tt.Package != "" {
			got, err := importer.Package(tt.File)
			if assert.NoError(t, err, "Package(%q)", tt.File) {
				assert.Equal(t, tt.Package, got, "Package(%q)", tt.File)
			}
		}
	}
}
