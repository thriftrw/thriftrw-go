// Copyright (c) 2021 Uber Technologies, Inc.
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
	"fmt"
	"log"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	"go.uber.org/thriftrw/compile"
)

var opts struct {
	RelativeTo string `long:"relative-to" description:"If specified, output paths will be relative to this directory"`
	Args       struct {
		ThriftFile string `positional-arg-name:"file" description:"Path to the Thrift file"`
	} `positional-args:"yes" required:"yes"`
}

// listDependentThrifts lists the Thrift files that the given Thrift file depends on or includes, both directly (through imports) and
// indirectly (through transitive imports).
//
// The returned file paths are absolute, unless relativeTo parameter is given, in which case the paths are relative to
// the relativeTo directory.
func listDependentThrifts(input string, relativeTo string) ([]string, error) {
	var deps []string

	module, err := compile.Compile(input)
	if err != nil {
		return nil, fmt.Errorf("could not compile %q: %v", input, err)
	}

	err = module.Walk(func(mod *compile.Module) error {
		output := mod.ThriftPath

		// Do not return a self-referencing dependency
		if output == input {
			return nil
		}

		if relativeTo != "" {
			output, err = filepath.Rel(relativeTo, output)
			if err != nil {
				return fmt.Errorf("%q depends on %q, which is not relative to %q; ensure that --relative-to"+
					" is an an ancestor of %q: %v", input, output, relativeTo, output, err)
			}
		}

		deps = append(deps, output)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return deps, nil
}

func run() error {
	if _, err := flags.Parse(&opts); err != nil {
		return fmt.Errorf("error parsing arguments: %v", err)
	}

	file := opts.Args.ThriftFile

	paths, err := listDependentThrifts(file, opts.RelativeTo)
	if err != nil {
		return fmt.Errorf("error listing deps of %q: %v", file, err)
	}

	for _, path := range paths {
		fmt.Println(path)
	}

	return nil
}

func main() {
	log.SetFlags(0) // so that the error message isn't noisy
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
