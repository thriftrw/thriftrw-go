package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.uber.org/multierr"
	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/compare"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("linter error: %v", err)
	}
}

func run(args []string) error {

	fs := flag.NewFlagSet("thriftcompat", flag.ContinueOnError)
	toFile := fs.String("to_file", "", "updated file")
	fromFile := fs.String("from_file", "", "original file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	args = fs.Args()
	if *toFile == "" {
		return fmt.Errorf("must provide an updated Thrift file")
	}
	// TODO: here we'll get into interesting stuff with git, but for now default to
	// using originalFile as an argument.

	return compileFiles(*toFile, *fromFile)
}

func compileFiles(toFile, fromFile string) error {
	toModule, err := compile.Compile(toFile)
	if err != nil {
		return err
	}
	fromModule, err := compile.Compile(fromFile)
	if err != nil {
		return err
	}

	err = checkRemovedMethods(toModule, fromModule)

	return multierr.Combine(err, checkRequiredFields(toModule, fromModule))
}

func checkRemovedMethods(toModule, fromModule *compile.Module) error {
	return compare.Services(toModule, fromModule)
}

func checkRequiredFields(toModule, fromModule *compile.Module) error {
	for n, spec := range toModule.Types {
		fromSpec, ok := fromModule.Types[n]
		if !ok {
			// This is a new Type, which is backwards compatible.
			continue
		}
		if s, ok := spec.(*compile.StructSpec); ok {
			// Match on Type names. Here we hit a limitation, that if someone
			// renames the struct and then adds a new field, we don't really have
			// a good way of tracking it.
			if fromStructSpec, ok := fromSpec.(*compile.StructSpec); ok {
				err := compare.StructSpecs(fromStructSpec, s)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

