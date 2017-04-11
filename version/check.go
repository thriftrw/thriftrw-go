// Copyright (c) 2017 Uber Technologies, Inc.
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

// +build !thriftrw.disableVersionCheck

package version

import (
	"log"

	"go.uber.org/thriftrw/internal/semver"
)

var compatRange = computeCompatibleRange()

func computeCompatibleRange() semver.Range {
	v, err := semver.Parse(Version)
	if err != nil {
		panic(err)
	}

	return semver.CompatibleRange(v)
}

// CheckCompatWithGeneratedCodeAt will panic if the ThriftRW version used to
// generated code (given by `genCodeVersion`) is not compatible with the
// current version of ThriftRW.
// This function is designed to be called during initialization of the
// generated code.
//
// Rationale: It is possible that the old generated code is not compatible with
// the new version of ThriftRW in subtle ways but still compiles successfully.
// This function will ensure that the version mismatch is detected and help
// avoid bugs that could be caused by this discrepancy.
func CheckCompatWithGeneratedCodeAt(genCodeVersion string, fromPkg string) {
	v, err := semver.Parse(genCodeVersion)
	if err != nil {
		panic(err)
	}

	if !compatRange.Contains(v) {
		log.Panicf(`incompatible version from generated package %q, expected >=%s and <%s, got %s`,
			fromPkg, &compatRange.Begin, &compatRange.End, &v)
	}
}
