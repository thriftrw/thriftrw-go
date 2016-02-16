// Copyright (c) 2015 Uber Technologies, Inc.
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
	"flag"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/thriftrw/thriftrw-go/compile"
	"github.com/thriftrw/thriftrw-go/gen"
)

func constantNames(constants map[string]*compile.Constant) []string {
	names := make([]string, 0, len(constants))
	for name := range constants {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func typeNames(types map[string]compile.TypeSpec) []string {
	names := make([]string, 0, len(types))
	for name := range types {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func main() {
	// TODO proper command line argument parsing

	output := flag.String("o", "", "Output file")
	flag.Parse()
	file := flag.Arg(0)

	outDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if len(*output) > 0 {
		var err error
		if *output, err = filepath.Abs(*output); err != nil {
			log.Fatal(err)
		}
		outDir = filepath.Dir(*output)
	}

	module, err := compile.Compile(file)
	if err != nil {
		log.Fatal(err)
	}

	packageName := filepath.Base(outDir)
	g := gen.NewGenerator(packageName)

	for _, constantName := range constantNames(module.Constants) {
		c := module.Constants[constantName]
		if err := g.Constant(c); err != nil {
			log.Fatal(err)
		}
	}

	for _, typeName := range typeNames(module.Types) {
		t := module.Types[typeName]
		if err := g.TypeDefinition(t); err != nil {
			log.Fatal(err)
		}
	}

	outFile := os.Stdout
	if len(*output) > 0 {
		outFile, err = os.Create(*output)
		defer outFile.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	if err := g.Write(outFile, token.NewFileSet()); err != nil {
		log.Fatal(err)
	}
}
