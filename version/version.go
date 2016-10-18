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

package version

import "log"

// Version is the current thriftrw version.
const Version = "0.4.0"

// CheckCompatibilityWithGeneratedCodeAt will panics if the thriftrw version
// used to generated code (given by `genCodeVer`) is not compatible with the
// current version of thriftrw. This function is intended to be called from the
// generated code.
// This function is designed to be called during initialization of the
// generated code.
//
// Rational: Let's say you use thriftrw version 1.0 to generate some stubs.
// Later on, you imports the stubs, but also thriftrw in version 1.2. Maybe
// thriftrw 1.2 is not compatible in subtle ways with the generated code from
// version 1.0. This function will make sure to panic during initialization
// preventing potential subtle bugs.
func CheckCompatibilityWithGeneratedCodeAt(genCodeVer string) {
	log.Printf("#### %s - %s", Version, genCodeVer)
}
