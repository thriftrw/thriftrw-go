package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	fromFile := fs.String("from_file", "", "updated file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	args = fs.Args()
	if toFile == nil {
		return fmt.Errorf("must provide an updated Thrift file")
	}
	// TODO: here we'll get into interesting stuff with git, but for now default to
	// using originalFile as an argument.

	if err := compileFiles(*toFile, *fromFile); err != nil {
		return err
	}

	return nil
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


	return checkRequired(toModule, fromModule)
}

func checkRequired(toModule, fromModule *compile.Module) error {
	// How do we find the correct Types to compare.
	for n, spec := range toModule.Types {
		if s, ok := spec.(*compile.StructSpec); ok {
			// Match on Type names. Here we hit a limitation, that if someone
			// renames the struct and then adds a new field, we don't really have
			// a good way of tracking it.
			fromSpec := fromModule.Types[n]
			if fromStructSpec, ok := fromSpec.(*compile.StructSpec); ok {
				err := compare.StructSpecs(fromStructSpec, s)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil

	// for n, spec := range .Types {
	// 	fmt.Println(spec)
	// 	fmt.Printf("n: %s, type is %T", n, spec)
	// 	fmt.Println()
	//
	// 	if s, ok := spec.(*compile.StructSpec); ok {
	// 		fmt.Println(len(s.Fields))
	// 		for _, f := range s.Fields {
	// 			fmt.Println(f)
	// 		}
	// 	}
	// }
	//
	// for _, spec := range originalModule.Types {
	// 	fmt.Println(spec)
	// 	fmt.Printf("type is %T", spec)
	// 	fmt.Println()
	//
	// 	if s, ok := spec.(*compile.StructSpec); ok {
	// 		fmt.Println(len(s.Fields))
	// 		for _, f := range s.Fields {
	// 			fmt.Println(f)
	// 		}
	// 	}
	// }
	//
	// return nil
}

