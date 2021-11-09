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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"go.uber.org/thriftrw/internal/git"
)

func main() {
	if err := run(os.Args[1:]); err != nil && err != flag.ErrHelp {
		log.Fatalf("%+v", err)
	}
}

func run(args []string) error {
	flag := flag.NewFlagSet("thriftbreak", flag.ContinueOnError)
	gitRepo := flag.String("C", "", "location of git repository. Defaults to current directory.")
	jsonOut := flag.Bool("json", false, "output as JSON")
	if err := flag.Parse(args); err != nil {
		return err
	}

	if *gitRepo == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("cannot determine current directory: %v", err)
		}
		*gitRepo = cwd
	}

	pass, err := git.Compare(*gitRepo)
	// Errors in compiling phase, but not in backwards compatibility.
	if err != nil {
		return err
	}

	lints := pass.Lints()
	var out string
	for _, l := range lints {
		if *jsonOut {
			bytes, err := json.Marshal(l)
			if err != nil {
				return fmt.Errorf("failed to marshal error: %v", err)
			}
			out = string(bytes)
		} else {
			out = l.String()
		}
		fmt.Println(out)
	}
	if len(lints) > 0 {
		return fmt.Errorf("found %d issues", len(lints))
	}

	return nil
}
