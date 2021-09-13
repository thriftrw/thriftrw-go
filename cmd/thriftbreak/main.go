package main

import (
	"flag"
	"log"
	"os"

	"go.uber.org/thriftrw/internal/compare"
	"go.uber.org/thriftrw/internal/git"

)

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("linter error: %v", err)
	}
}

// NoFileError is returned is we fail to provide any thrift files to check
// and no --git_repo option.
type NoFileError struct{}

func (e NoFileError) Error() string {
	return "must provide an updated Thrift file"
}

func run(args []string) error {
	fs := flag.NewFlagSet("thriftcompat", flag.ContinueOnError)
	toFile := fs.String("to_file", "", "updated file")
	fromFile := fs.String("from_file", "", "original file")
	gitRepo := fs.String("git_repo", "", "location of git repository")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if *toFile == "" && *gitRepo == "" {
		return NoFileError{}
	}
	// Here we are going to access git repo to find the previous version of a file.
	if *gitRepo != "" {
		return git.UseGit(*gitRepo)
	}

	return compare.Files(*toFile, *fromFile)
}
