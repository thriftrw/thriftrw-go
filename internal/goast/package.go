// Copyright (c) 2016 Uber Technologies, Inc.
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

package goast

import (
	"go/build"
	"path/filepath"
)

// DeterminePackageName determines the name of the package at the given import
// path.
func DeterminePackageName(importPath string) string {
	// TODO(abg): This can be a lot faster by using build.FindOnly and parsing one
	// of the .go files in the directory with parser.PackageClauseOnly set. See
	// how goimports does this:
	// https://github.com/golang/tools/blob/0e9f43fcb67267967af8c15d7dc54b373e341d20/imports/fix.go#L284

	pkg, err := build.Import(importPath, "", 0)
	if err != nil {
		return filepath.Base(importPath)
	}
	return pkg.Name
}
