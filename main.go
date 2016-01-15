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
	"go/token"
	"log"
	"os"
	"sort"

	"github.com/uber/thriftrw-go/compile"
	"github.com/uber/thriftrw-go/gen"
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
	file := os.Args[1]
	module, err := compile.Compile(file)
	if err != nil {
		log.Fatal(err)
	}

	g := gen.NewGenerator()

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

	if err := g.Write(os.Stdout, token.NewFileSet()); err != nil {
		log.Fatal(err)
	}
}
