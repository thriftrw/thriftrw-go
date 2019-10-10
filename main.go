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

package main // import "go.uber.org/thriftrw"

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/gen"
	"go.uber.org/thriftrw/internal/plugin"
	"go.uber.org/thriftrw/internal/plugin/builtin/pluginapigen"
	"go.uber.org/thriftrw/version"

	flags "github.com/jessevdk/go-flags"
	"go.uber.org/multierr"
)

type options struct {
	DisplayVersion bool       `long:"version" short:"v" description:"Show the ThriftRW version number"`
	GOpts          genOptions `group:"Generator Options"`
}

type genOptions struct {
	OutputDirectory string `long:"out" short:"o" value-name:"DIR" description:"Directory to which the generated files will be written."`
	PackagePrefix   string `long:"pkg-prefix" value-name:"PREFIX" description:"Prefix for import paths of generated module. By default, this is based on the output directory's location relative to $GOPATH."`
	ThriftRoot      string `long:"thrift-root" value-name:"DIR" description:"Directory whose descendants contain all Thrift files. The structure of the generated Go packages mirrors the paths to the Thrift files relative to this directory. By default, this is the deepest common ancestor directory of the Thrift files."`

	NoRecurse bool         `long:"no-recurse" description:"Don't generate code for included Thrift files."`
	Plugins   plugin.Flags `long:"plugin" short:"p" value-name:"PLUGIN" description:"Code generation plugin for ThriftRW. This option may be provided multiple times to apply multiple plugins."`

	GeneratePluginAPI bool   `long:"generate-plugin-api" hidden:"true" description:"Generates code for the plugin API"`
	NoVersionCheck    bool   `long:"no-version-check" hidden:"true" description:"Does not add library version checks to generated code."`
	NoTypes           bool   `long:"no-types" description:"Do not generate code for types, implies --no-service-helpers."`
	NoConstants       bool   `long:"no-constants" description:"Do not generate code for const declarations."`
	NoServiceHelpers  bool   `long:"no-service-helpers" description:"Do not generate service helpers."`
	NoEmbedIDL        bool   `long:"no-embed-idl" description:"Do not embed IDLs into the generated code."`
	NoZap             bool   `long:"no-zap" description:"Do not generate code for Zap logging."`
	OutputFile        string `long:"output-file" value-name:"FILENAME" description:"Generates a single .go file as an output. Specifying an OutputFile prevents code generation for included Thrift Files."`

	// TODO(abg): Detailed help with examples of --thrift-root, --pkg-prefix,
	// and --plugin

}

func main() {
	if err := do(); err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
	os.Exit(0)
}

func do() (err error) {
	log.SetFlags(0) // don't include timestamps, etc. in the output

	var opts options

	parser := flags.NewParser(&opts, flags.Default & ^flags.PrintErrors)
	parser.Usage = "[OPTIONS] FILE"

	args, err := parser.Parse()
	if ferr, ok := err.(*flags.Error); ok && ferr.Type == flags.ErrHelp {
		parser.WriteHelp(os.Stdout)
		return nil
	} else if err != nil {
		return err
	}

	if opts.DisplayVersion {
		fmt.Printf("thriftrw v%s\n", version.Version)
		return nil
	}

	if len(args) != 1 {
		var buffer bytes.Buffer
		parser.WriteHelp(&buffer)
		return errors.New(buffer.String())
	}

	inputFile := args[0]
	if _, err := os.Stat(inputFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("File %q does not exist: %v", inputFile, err)
		}
		return fmt.Errorf("Could not stat file %q: %v", inputFile, err)
	}
	gopts := opts.GOpts

	if len(gopts.OutputDirectory) == 0 {
		gopts.OutputDirectory = "."
	}
	gopts.OutputDirectory, err = filepath.Abs(gopts.OutputDirectory)
	if err != nil {
		return fmt.Errorf("Unable to resolve absolute path for %q: %v", gopts.OutputDirectory, err)
	}

	if gopts.PackagePrefix == "" {
		gopts.PackagePrefix, err = determinePackagePrefix(gopts.OutputDirectory)
		if err != nil {
			return fmt.Errorf(
				"Could not determine a package prefix automatically: %v\n"+
					"A package prefix is required to use correct import paths in the generated code.\n"+
					"Use the --pkg-prefix option to provide a package prefix manually.", err)
		}
	}

	module, err := compile.Compile(inputFile)
	if err != nil {
		// TODO(abg): For nested compile errors, split causal chain across
		// multiple lines.
		return fmt.Errorf("Failed to compile %q: %+v", inputFile, err)
	}

	if gopts.ThriftRoot == "" {
		gopts.ThriftRoot, err = findCommonAncestor(module)
		if err != nil {
			return fmt.Errorf(
				"Could not find a common parent directory for %q and the Thrift files "+
					"imported by it.\nThis directory is required to generate a consistent "+
					"hierarchy for generated packages.\nUse the --thrift-root option to "+
					"provide this path.\n\t%v", inputFile, err)
		}
	} else {
		gopts.ThriftRoot, err = filepath.Abs(gopts.ThriftRoot)
		if err != nil {
			return fmt.Errorf("Unable to resolve absolute path for %q: %v", gopts.ThriftRoot, err)
		}
		if err := verifyAncestry(module, gopts.ThriftRoot); err != nil {
			return fmt.Errorf(
				"An included Thrift file is not contained in the %q directory tree: %v",
				gopts.ThriftRoot, err)
		}
	}

	if len(gopts.OutputFile) > 0 && filepath.Ext(gopts.OutputFile) != ".go" {
		return fmt.Errorf("output-file value: %q invalid. A {FILENAME}.go name must be provided", gopts.OutputFile)
	}

	pluginHandle, err := gopts.Plugins.Handle()
	if err != nil {
		return fmt.Errorf("Failed to initialize plugins: %+v", err)
	}

	if gopts.GeneratePluginAPI {
		pluginHandle = append(pluginHandle, pluginapigen.Handle)
	}

	defer func() {
		err = multierr.Append(err, pluginHandle.Close())
	}()

	generatorOptions := gen.Options{
		OutputDir:        gopts.OutputDirectory,
		PackagePrefix:    gopts.PackagePrefix,
		ThriftRoot:       gopts.ThriftRoot,
		NoRecurse:        gopts.NoRecurse,
		NoVersionCheck:   gopts.NoVersionCheck,
		Plugin:           pluginHandle,
		NoTypes:          gopts.NoTypes,
		NoConstants:      gopts.NoConstants,
		NoServiceHelpers: gopts.NoServiceHelpers || gopts.NoTypes,
		NoEmbedIDL:       gopts.NoEmbedIDL,
		NoZap:            gopts.NoZap,
		OutputFile:       gopts.OutputFile,
	}
	if err := gen.Generate(module, &generatorOptions); err != nil {
		return fmt.Errorf("Failed to generate code: %+v", err)
	}
	return nil
}

// verifyAncestry verifies that the Thrift file for the given module and the
// Thrift files for all imported modules are contained within the directory
// tree rooted at the given path.
func verifyAncestry(m *compile.Module, root string) error {
	return m.Walk(func(m *compile.Module) error {
		path, err := filepath.Rel(root, m.ThriftPath)
		if err != nil {
			return fmt.Errorf(
				"could not resolve path for %q: %v", m.ThriftPath, err)
		}

		if strings.HasPrefix(path, "..") {
			return fmt.Errorf(
				"%q is not contained in the %q directory tree",
				m.ThriftPath, root)
		}

		return nil
	})
}

// findCommonAncestor finds the deepest common ancestor for the given module
// and all modules imported by it.
func findCommonAncestor(m *compile.Module) (string, error) {
	var result []string
	var lastString string

	err := m.Walk(func(m *compile.Module) error {
		thriftPath := m.ThriftPath
		if !filepath.IsAbs(thriftPath) {
			return fmt.Errorf(
				"ThriftPath must be absolute: %q is not absolute", thriftPath)
		}

		thriftDir := filepath.Dir(thriftPath)

		// Split("/foo/bar", "/") = ["", "foo", "bar"]
		parts := strings.Split(thriftDir, string(filepath.Separator))
		if result == nil {
			result = parts
			lastString = thriftPath
			return nil
		}

		result = commonPrefix(result, parts)
		if len(result) == 1 && result[0] == "" {
			return fmt.Errorf(
				"%q does not share an ancestor with %q",
				thriftPath, lastString)
		}

		lastString = thriftPath
		return nil
	})
	if err != nil {
		return "", err
	}

	return strings.Join(result, string(filepath.Separator)), nil
}

// commonPrefix finds the shortest common prefix for the two lists.
//
// An empty slice may be returned if the two lists don't have a common prefix.
func commonPrefix(l, r []string) []string {
	var i int
	for i = 0; i < len(l) && i < len(r); i++ {
		if l[i] != r[i] {
			break
		}
	}
	return l[:i]
}

// determinePackagePrefix determines the package prefix for Go packages
// generated in this file.
//
// dir must be an absolute path.
func determinePackagePrefix(dir string) (string, error) {
	gopathList := os.Getenv("GOPATH")
	if gopathList == "" {
		return "", errors.New("$GOPATH is not set")
	}

	for _, gopath := range filepath.SplitList(gopathList) {
		packagePath, err := filepath.Rel(filepath.Join(gopath, "src"), dir)
		if err != nil {
			return "", err
		}

		// The match is valid only if it's within the directory tree.
		if !strings.HasPrefix(packagePath, "..") {
			return packagePath, nil
		}
	}

	return "", fmt.Errorf("directory %q is not inside $GOPATH/src", dir)
}
