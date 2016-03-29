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

	if o.NoRecurse {
		return generateModule(m, o)
	}

	return m.Walk(func(m *compile.Module) error {
		if err := generateModule(m, o); err != nil {
			return generateError{Name: m.ThriftPath, Reason: err}
		}
		return nil
	})
}

// generates code for only the given module, assuming that code for included
// modules has already been generated.
func generateModule(m *compile.Module, o *Options) error {
	packagePath, err := filepath.Rel(o.ThriftRoot, m.ThriftPath)
	if err != nil {
		return err
	}

	// packagePath is the path relative to outputDir into which we'll be
	// writing the package for this Thrift file. For
	// $thriftRoot/foo/bar.thrift, packagePath is foo/bar, and packageDir is
	// $outputDir/foo/bar. All files for bar.thrift will be written to
	// $outputDir/foo/bar/. The package will be importable via
	// $packagePrefix/foo/bar.
	packagePath = strings.TrimSuffix(packagePath, ".thrift")
	packageDir := filepath.Join(o.OutputDir, packagePath)

	// files contains a collection of files relative to packageDir and their
	// contents.
	files := make(map[string]*bytes.Buffer)

	if len(m.Constants) > 0 {
		g := NewGenerator(o.PackagePrefix, o.ThriftRoot, m.ThriftPath)

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
		g := NewGenerator(o.PackagePrefix, o.ThriftRoot, m.ThriftPath)

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

		files["types.go"] = buff
	}

	// TODO(abg): Services

	for relPath, contents := range files {
		fullPath := filepath.Join(packageDir, relPath)
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
