// Copyright (c) 2019 Uber Technologies, Inc.
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
	"io/ioutil"
	"path/filepath"
)

// Option represents a compiler option.
type Option func(*compiler)

// FS is used by the compiler to interact with the filesystem.
type FS interface {
	// Read reads the file named by filename and returns the contents.
	// See: https://golang.org/pkg/io/ioutil/#ReadFile
	Read(filename string) ([]byte, error)
	// Abs returns an absolute representation of path.
	// See: https://golang.org/pkg/path/filepath/#Abs
	Abs(p string) (string, error)
}

type realFS struct{}

func (realFS) Read(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (realFS) Abs(p string) (string, error) {
	return filepath.Abs(p)
}

// Filesystem controls how the Thrift compiler accesses the filesystem.
func Filesystem(fs FS) Option {
	return func(c *compiler) {
		c.fs = fs
	}
}

// NonStrict disables strict validation of the Thrift file. This allows
// struct fields which are not marked as optional or required.
func NonStrict() Option {
	return func(c *compiler) {
		c.nonStrict = true
	}
}
