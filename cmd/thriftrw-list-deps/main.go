package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	"go.uber.org/thriftrw/compile"
)

type args struct {
	ThriftFile string `positional-arg-name:"file" description:"Path to the Thrift file"`
}

var opts struct {
	RelativeTo string `long:"relative-to" description:"If specified, output paths will be relative to this directory"`
	Args       args   `positional-args:"yes" required:"yes"`
}

// listDependentThrifts lists the Thrift files that the given Thrift file depends on or includes, both directly (through imports) and
// indirectly (through transitive imports).
//
// The returned file paths are absolute, unless relativeDir parameter is given, in which case the paths are relative to
// the relativeDir.
func listDependentThrifts(inputThriftFile string, relativeDir string) ([]string, error) {
	dependentThriftFiles := make([]string, 0)

	module, err := compile.Compile(inputThriftFile)
	if err != nil {
		return nil, fmt.Errorf("could not compile %q: %v", inputThriftFile, err)
	}

	for _, included := range module.Includes {
		outputFilePath := included.Module.ThriftPath

		if relativeDir != "" {
			outputFilePath, err = filepath.Rel(relativeDir, outputFilePath)
			if err != nil {
				return nil, fmt.Errorf("cannot make %q relativeDir to %q: %v", outputFilePath, relativeDir, err)
			}
		}

		dependentThriftFiles = append(dependentThriftFiles, outputFilePath)
	}

	return dependentThriftFiles, nil
}

func run() error {
	if _, err := flags.Parse(&opts); err != nil {
		return fmt.Errorf("error parsing arguments: %v", err)
	}

	inputThriftFile := opts.Args.ThriftFile

	paths, err := listDependentThrifts(inputThriftFile, opts.RelativeTo)
	if err != nil {
		return fmt.Errorf("error listing deps of %q: %v", inputThriftFile, err)
	}

	for _, path := range paths {
		fmt.Println(path)
	}

	return nil
}

func main() {
	log.SetFlags(0) // so that the error message isn't noisy
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}
