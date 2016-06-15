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
	"io/ioutil"
	"os"
	"testing"

	"github.com/thriftrw/thriftrw-go/compile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			PackagePrefix: "github.com/thriftrw/thriftrw-go/gen",
			ThriftRoot:    "testdata",
		},
		{
			OutputDir:     "testdata",
			PackagePrefix: "github.com/thriftrw/thriftrw-go/gen",
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

func TestThriftPackageImporter(t *testing.T) {
	importer := thriftPackageImporter{
		ImportPrefix: "github.com/myteam/myservice",
		ThriftRoot:   "/src/thrift",
	}

	tests := []struct {
		File, ServiceName string // Inputs

		// If non-empty, these are the expected outputs for RelativePackage,
		// Package, and ServicePackage.
		Relative, Package, ServicePackage string
	}{
		{
			File:           "/src/thrift/foo.thrift",
			Relative:       "foo",
			Package:        "github.com/myteam/myservice/foo",
			ServiceName:    "MyService",
			ServicePackage: "github.com/myteam/myservice/foo/service/myservice",
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

		if tt.ServicePackage != "" {
			got, err := importer.ServicePackage(tt.File, tt.ServiceName)
			if assert.NoError(t, err, "ServicePackage(%q, %q)", tt.File, tt.ServiceName) {
				assert.Equal(t, tt.ServicePackage, got,
					"ServicePackage(%q, %q)", tt.File, tt.ServiceName)
			}
		}
	}
}
