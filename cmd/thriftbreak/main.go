package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go.uber.org/thriftrw/internal/compare"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("linter error: %v", err)
	}
}

type NoFileError struct {}

func (e NoFileError) Error() string {
	return fmt.Sprint("must provide an updated Thrift file")
}

func run(args []string) error {
	fs := flag.NewFlagSet("thriftcompat", flag.ContinueOnError)
	toFile := fs.String("to_file", "", "updated file")
	fromFile := fs.String("from_file", "", "original file")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *toFile == "" {
		return NoFileError{}
	}
	// TODO: here we'll get into interesting stuff with git, but for now default to
	// using originalFile as an argument.
	return compare.Files(*toFile, *fromFile)
}
