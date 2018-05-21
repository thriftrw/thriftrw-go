package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jessevdk/go-flags"

	"go.uber.org/thriftrw/compile"
)

type args struct {
	ThriftFile string `positional-arg-name:"thrift-file"`
}

var opts struct {
	Relative string `long:"relative" description:"Output paths will be relative to this directory"`
	Args     args   `positional-args:"yes" required:"yes"`
}

func do(relative string, inputThriftFile string, output io.Writer) error {
	module, err := compile.Compile(inputThriftFile)
	if err != nil {
		return err
	}

	for _, included := range module.Includes {
		outputFilePath := included.Module.ThriftPath
		if relative != "" {
			outputFilePath, err = filepath.Rel(relative, outputFilePath)
			if err != nil {
				return err
			}
		}
		fmt.Fprintf(output, "%s\n", outputFilePath)
	}

	return nil
}

func main() {
	if _, err := flags.Parse(&opts); err != nil {
		log.Fatalf("error parsing arguments: %s", err.Error())
		os.Exit(1)
	}

	inputThriftFile := opts.Args.ThriftFile

	if err := do(opts.Relative, inputThriftFile, os.Stdout); err != nil {
		log.Fatalf("error listing deps of %s: %s", inputThriftFile, err.Error())
		os.Exit(1)
	}
}
