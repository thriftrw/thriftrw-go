// Copyright (c) 2016 Uber Technologies, Inc.
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
	"bytes"
	"fmt"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/thriftrw/thriftrw-go/compile"
)

// Options controls how code gets generated.
type Options struct {
	// OutputDir is the directory into which all generated code is written.
	//
	// This must be an absolute path.
	OutputDir string

	// PackagePrefix controls the import path prefix for all generated
	// packages.
	PackagePrefix string

	// ThriftRoot is the directory within whose tree all Thrift files consumed
	// are contained. The locations of the Thrift files relative to the
	// ThriftFile determines the module structure in OutputDir.
	//
	// This must be an absolute path.
	ThriftRoot string

	// NoRecurse determines whether code should be generated for included Thrift
	// files as well. If true, code gets generated only for the first module.
	NoRecurse bool
}

// Generate generates code based on the given options.
func Generate(m *compile.Module, o *Options) error {
	if !filepath.IsAbs(o.ThriftRoot) {
		return fmt.Errorf(
			"ThriftRoot must be an absolute path: %q is not absolute",
			o.ThriftRoot)
	}

	if !filepath.IsAbs(o.OutputDir) {
		return fmt.Errorf(
			"OutputDir must be an absolute path: %q is not absolute",
			o.OutputDir)
	}

	importer := thriftPackageImporter{
		ImportPrefix: o.PackagePrefix,
		ThriftRoot:   o.ThriftRoot,
	}

	if o.NoRecurse {
		return generateModule(m, importer, o.OutputDir)
	}

	return m.Walk(func(m *compile.Module) error {
		if err := generateModule(m, importer, o.OutputDir); err != nil {
			return generateError{Name: m.ThriftPath, Reason: err}
		}
		return nil
	})
}

type thriftPackageImporter struct {
	ImportPrefix string
	ThriftRoot   string
}

// RelativePackage returns the import path for the top-level package of the
// given Thrift file relative to the ImportPrefix.
func (i thriftPackageImporter) RelativePackage(file string) (string, error) {
	return filepath.Rel(i.ThriftRoot, strings.TrimSuffix(file, ".thrift"))
}

// Package returns the import path for the top-level package of the given Thrift
// file.
func (i thriftPackageImporter) Package(file string) (string, error) {
	pkg, err := i.RelativePackage(file)
	if err != nil {
		return "", err
	}
	return filepath.Join(i.ImportPrefix, pkg), nil
}

// ServicePackage returns the import path for the package for the Thrift service
// with the given name defined in the given file.
func (i thriftPackageImporter) ServicePackage(file, name string) (string, error) {
	topPackage, err := i.Package(file)
	if err != nil {
		return "", err
	}

	return filepath.Join(topPackage, strings.ToLower(name)), nil
}

// generates code for only the given module, assuming that code for included
// modules has already been generated.
func generateModule(m *compile.Module, i thriftPackageImporter, outDir string) error {
	// packageRelPath is the path relative to outputDir into which we'll be
	// writing the package for this Thrift file. For $thriftRoot/foo/bar.thrift,
	// packageRelPath is foo/bar, and packageDir is $outputDir/foo/bar. All
	// files for bar.thrift will be written to the $outputDir/foo/bar/ tree. The
	// package will be importable via $importPrefix/foo/bar.
	packageRelPath, err := i.RelativePackage(m.ThriftPath)
	if err != nil {
		return err
	}

	// TODO(abg): Prefer top-level package name from `namespace go` directive.
	packageName := filepath.Base(packageRelPath)

	// importPath is the full import path for the top-level package generated
	// for this Thrift file.
	importPath, err := i.Package(m.ThriftPath)
	if err != nil {
		return err
	}

	// packageOutDir is the directory whithin which all files and folders for
	// this Thrift file will be written.
	packageOutDir := filepath.Join(outDir, packageRelPath)

	// Mapping of file names relative to $packageOutDir and their contents.
	files := make(map[string]*bytes.Buffer)

	if len(m.Constants) > 0 {
		g := NewGenerator(i, importPath, packageName)

		for _, constantName := range sortStringKeys(m.Constants) {
			if err := Constant(g, m.Constants[constantName]); err != nil {
				return err
			}
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, token.NewFileSet()); err != nil {
			return fmt.Errorf(
				"could not generate constants for %q: %v", m.ThriftPath, err)
		}

		// TODO(abg): Verify no file collisions
		files["constants.go"] = buff
	}

	if len(m.Types) > 0 {
		g := NewGenerator(i, importPath, packageName)

		for _, typeName := range sortStringKeys(m.Types) {
			if err := TypeDefinition(g, m.Types[typeName]); err != nil {
				return err
			}
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, token.NewFileSet()); err != nil {
			return fmt.Errorf(
				"could not generate types for %q: %v", m.ThriftPath, err)
		}

		// TODO(abg): Verify no file collisions
		files["types.go"] = buff
	}

	// TODO(abg): Services

	for relPath, contents := range files {
		fullPath := filepath.Join(packageOutDir, relPath)
		directory := filepath.Dir(fullPath)

		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("could not create directory %q: %v", directory, err)
		}

		if err := ioutil.WriteFile(fullPath, contents.Bytes(), 0644); err != nil {
			return fmt.Errorf("failed to write %q: %v", fullPath, err)
		}
	}
	return nil
}

// sortStringKeys returns a sorted list of strings given a map[string]*.
func sortStringKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	t := v.Type()
	if t.Kind() != reflect.Map || t.Key().Kind() != reflect.String {
		panic(fmt.Sprintf(
			"sortStringKeys may be called with a map[string]* only"))
	}

	keys := v.MapKeys()
	sortedKeys := make([]string, 0, len(keys))

	for _, k := range keys {
		key := k.Interface().(string)
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)
	return sortedKeys
}
