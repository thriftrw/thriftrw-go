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

package compile

import (
	"fmt"
	"strings"
)

// dummyFS provides a fake in-memory filesystem for use with Compile.
type dummyFS struct {
	// Absolute path to the fake current working directory
	CWD string

	// Map of absolute paths of files to their contents.
	Files map[string]string
}

func (fs dummyFS) Abs(p string) (string, error) {
	if strings.HasPrefix(p, "/") {
		return p, nil
	}
	return fs.CWD + p, nil
}

func (fs dummyFS) Read(path string) ([]byte, error) {
	if contents, ok := fs.Files[path]; ok {
		return []byte(contents), nil
	}

	return nil, fmt.Errorf("file not found: %v", path)
}
