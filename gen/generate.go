// Copyright (c) 2018 Uber Technologies, Inc.
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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/plugin"

	"go.uber.org/multierr"
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

	// If true, we will not generate versioncheck.go files.
	NoVersionCheck bool

	// Code generation plugin
	Plugin plugin.Handle

	// Do not generate types.go
	NoTypes bool

	// Do not generate constants.go
	NoConstants bool

	// Do not generate service helpers
	NoServiceHelpers bool

	// Do not embed IDLs in generated code
	NoEmbedIDL bool
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

	// Mapping of filenames relative to OutputDir to their contents.
	files := make(map[string][]byte)
	genBuilder := newGenerateServiceBuilder(importer)

	generate := func(m *compile.Module) error {
		moduleFiles, err := generateModule(m, importer, genBuilder, o)
		if err != nil {
			return generateError{Name: m.ThriftPath, Reason: err}
		}
		if err := mergeFiles(files, moduleFiles); err != nil {
			return generateError{Name: m.ThriftPath, Reason: err}
		}
		return nil
	}

	// Note that we call generate directly on only those modules that we need
	// to generate code for. If the user used --no-recurse, we're not going to
	// generate code for included modules.
	if o.NoRecurse {
		if err := generate(m); err != nil {
			return err
		}
	} else {
		if err := m.Walk(generate); err != nil {
			return err
		}
	}

	plug := o.Plugin
	if plug == nil {
		plug = plugin.EmptyHandle
	}

	if sgen := plug.ServiceGenerator(); sgen != nil {
		res, err := sgen.Generate(genBuilder.Build())
		if err != nil {
			return err
		}

		if err := mergeFiles(files, res.Files); err != nil {
			return err
		}
	}

	for relPath, contents := range files {
		fullPath := filepath.Join(o.OutputDir, relPath)
		directory := filepath.Dir(fullPath)

		if err := os.MkdirAll(directory, 0755); err != nil {
			return fmt.Errorf("could not create directory %q: %v", directory, err)
		}

		if err := ioutil.WriteFile(fullPath, contents, 0644); err != nil {
			return fmt.Errorf("failed to write %q: %v", fullPath, err)
		}
	}

	return nil
}

// TODO(abg): Make some sort of public interface out of the Importer

type thriftPackageImporter struct {
	ImportPrefix string
	ThriftRoot   string
}

// RelativePackage returns the import path for the top-level package of the
// given Thrift file relative to the ImportPrefix.
func (i thriftPackageImporter) RelativePackage(file string) (string, error) {
	return filepath.Rel(i.ThriftRoot, strings.TrimSuffix(file, ".thrift"))
}

func (i thriftPackageImporter) RelativeThriftFilePath(file string) (string, error) {
	return filepath.Rel(i.ThriftRoot, file)
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

func mergeFiles(dest, src map[string][]byte) error {
	var errors []error
	for path, contents := range src {
		if _, ok := dest[path]; ok {
			errors = append(errors, fmt.Errorf("file generation conflict: "+
				"multiple sources are trying to write to %q", path))
		}
		dest[path] = contents
	}
	return multierr.Combine(errors...)
}

// generateModule returns a mapping from filename to file contents of files that
// should be generated relative to o.OutputDir.
func generateModule(m *compile.Module, i thriftPackageImporter, builder *generateServiceBuilder, o *Options) (map[string][]byte, error) {
	// packageRelPath is the path relative to outputDir into which we'll be
	// writing the package for this Thrift file. For $thriftRoot/foo/bar.thrift,
	// packageRelPath is foo/bar, and packageDir is $outputDir/foo/bar. All
	// files for bar.thrift will be written to the $outputDir/foo/bar/ tree. The
	// package will be importable via $importPrefix/foo/bar.
	packageRelPath, err := i.RelativePackage(m.ThriftPath)
	if err != nil {
		return nil, err
	}

	// TODO(abg): Prefer top-level package name from `namespace go` directive.
	packageName := filepath.Base(packageRelPath)

	// importPath is the full import path for the top-level package generated
	// for this Thrift file.
	importPath, err := i.Package(m.ThriftPath)
	if err != nil {
		return nil, err
	}

	// Mapping of file names relative to packageRelPath to their contents.
	// Note that we need to return a mapping relative to o.OutputDir so we
	// will prepend $packageRelPath/ to all these paths.
	files := make(map[string][]byte)

	g := NewGenerator(i, importPath, packageName)

	if len(m.Constants) > 0 {
		for _, constantName := range sortStringKeys(m.Constants) {
			if err := Constant(g, m.Constants[constantName]); err != nil {
				return nil, err
			}
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, nil /* fset */); err != nil {
			return nil, fmt.Errorf(
				"could not generate constants for %q: %v", m.ThriftPath, err)
		}

		// TODO(abg): Verify no file collisions
		if !o.NoConstants {
			files["constants.go"] = buff.Bytes()
		}
	}

	if len(m.Types) > 0 {
		for _, typeName := range sortStringKeys(m.Types) {
			if err := TypeDefinition(g, m.Types[typeName]); err != nil {
				return nil, err
			}
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, nil /* fset */); err != nil {
			return nil, fmt.Errorf(
				"could not generate types for %q: %v", m.ThriftPath, err)
		}

		// TODO(abg): Verify no file collisions
		if !o.NoTypes {
			files["types.go"] = buff.Bytes()
		}
	}

	if !o.NoEmbedIDL {
		if err := embedIDL(g, i, m); err != nil {
			return nil, err
		}

		buff := new(bytes.Buffer)
		if err := g.Write(buff, nil /* fset */); err != nil {
			return nil, fmt.Errorf(
				"could not generate idl.go for %q: %v", m.ThriftPath, err)
		}

		files["idl.go"] = buff.Bytes()
	}

	// Services must be generated last because names of user-defined types take
	// precedence over the names we pick for the service types.
	if len(m.Services) > 0 {
		for _, serviceName := range sortStringKeys(m.Services) {
			service := m.Services[serviceName]

			// generateModule gets called only for those modules for which we
			// need to generate code. With --no-recurse, generateModule is
			// called only on the root file specified by the user and not its
			// included modules. Only services defined in these files are
			// considered root services; plugins will generate code only for
			// root services, even though they have information about the
			// whole service tree.
			if _, err := builder.AddRootService(service); err != nil {
				return nil, err
			}

			serviceFiles, err := Service(g, service)
			if err != nil {
				return nil, fmt.Errorf(
					"could not generate code for service %q: %v",
					serviceName, err)
			}

			if !o.NoServiceHelpers {
				for name, buff := range serviceFiles {
					files[name] = buff.Bytes()
				}
			}
		}
	}

	newFiles := make(map[string][]byte, len(files))
	for path, contents := range files {
		newFiles[filepath.Join(packageRelPath, path)] = contents
	}
	return newFiles, nil
}
