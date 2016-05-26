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

package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/thriftrw/thriftrw-go/compile"
	"github.com/thriftrw/thriftrw-go/gen"
)

func main() {
	var outputDir string
	flag.StringVar(
		&outputDir, "out", ".",
		"Directory to which the generated files will be written.")

	var packagePrefix string
	flag.StringVar(
		&packagePrefix, "pkg-prefix", "",
		"Prefix for import paths of generated modules. By default, this is "+
			"based on the output directory's location relative to $GOPATH.")

	var thriftRoot string
	flag.StringVar(
		&thriftRoot, "thrift-root", "",
		"Directory whose descendants contain all the used Thrift files. "+
			"The structure of the generated Go packages mirrors the paths to "+
			"the Thrift files relative to this directory. By default, this is "+
			"the deepest common ancestor of the Thrift files.")

	var yarpc bool
	flag.BoolVar(&yarpc, "yarpc", false, "Generate code for YARPC.")

	var noRecurse bool
	flag.BoolVar(&noRecurse, "no-recurse", false,
		"Disable code generation for included Thrift files.")

	flag.Parse()

	inputFile := flag.Arg(0)
	if len(inputFile) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if _, err := os.Stat(inputFile); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("file %q does not exist: %v", inputFile, err)
		}
		log.Fatalf("error describing file %q: %v", inputFile, err)
	}

	if len(outputDir) == 0 {
		outputDir = "."
	}
	outputDir, err := filepath.Abs(outputDir)
	if err != nil {
		log.Fatalf("could not resolve path %q: %v", outputDir, err)
	}

	if packagePrefix == "" {
		packagePrefix, err = determinePackagePrefix(outputDir)
		if err != nil {
			log.Fatalf("could not determine the package prefix: %v", err)
		}
	}

	module, err := compile.Compile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if thriftRoot == "" {
		thriftRoot, err = findCommonAncestor(module)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		thriftRoot, err = filepath.Abs(thriftRoot)
		if err != nil {
			log.Fatal(err)
		}
		if err := verifyAncestry(module, thriftRoot); err != nil {
			log.Fatal(err)
		}
	}

	opts := gen.Options{
		OutputDir:     outputDir,
		PackagePrefix: packagePrefix,
		ThriftRoot:    thriftRoot,
		NoRecurse:     noRecurse,
		YARPC:         yarpc,
	}

	if err := gen.Generate(module, &opts); err != nil {
		log.Fatal(err)
	}
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
