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

package gen

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"go.uber.org/thriftrw/wire"
)

// This file contains helpers for the different test cases in this module.

func singleFieldStruct(id int16, value wire.Value) wire.Value {
	return wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
		{ID: id, Value: value},
	}})
}

func boolp(x bool) *bool         { return &x }
func bytep(x int8) *int8         { return &x }
func int16p(x int16) *int16      { return &x }
func int32p(x int32) *int32      { return &x }
func int64p(x int64) *int64      { return &x }
func doublep(x float64) *float64 { return &x }
func stringp(x string) *string   { return &x }

func hash(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func dirhash(dir string) (string, error) {
	fileHashes := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileHash, err := hash(path)
		if err != nil {
			return fmt.Errorf("failed to hash %q: %v", path, err)
		}

		// We only care about the path relative to the directory being
		// hashed.
		path, err = filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		fileHashes[path] = fileHash
		return nil
	})
	if err != nil {
		return "", err
	}

	fileNames := make([]string, 0, len(fileHashes))
	for name := range fileHashes {
		fileNames = append(fileNames, name)
	}
	sort.Strings(fileNames)

	h := sha1.New()
	for _, name := range fileNames {
		if _, err := fmt.Fprintf(h, "%v\t%v\n", name, fileHashes[name]); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
