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

package gen

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T) {
	tests := [][]struct{ Path, Name string }{
		{
			{
				Path: "go.uber.org/thriftrw/wire",
				Name: "wire",
			},
			{
				Path: "go.uber.org/thriftrw/another/wire",
				Name: "wire2",
			},
		},
		{
			{
				Path: "foo/bar",
				Name: "bar",
			},
			{
				Path: "baz/bar-go",
				Name: "bar2",
			},
		},
		{
			{
				Path: "github.com/yarpc/yarpc-go",
				Name: "yarpc",
			},
			{
				Path: "go.uber.org/thriftrw/yarpc",
				Name: "yarpc2",
			},
		},
	}

	for _, tt := range tests {
		imp := newImporter(NewNamespace())
		for _, e := range tt {
			assert.Equal(t, e.Name, imp.Import(e.Path))
		}

		for _, e := range tt {
			i, ok := imp.imports[e.Path]
			if !assert.True(t, ok, "could not find %q", e.Path) {
				continue
			}
			if e.Name == filepath.Base(e.Path) {
				continue
			}
			if assert.NotNil(t, i.Name, "expected non-nil name for %q", e.Path) {
				assert.Equal(t, e.Name, i.Name.Name, "name for %q did not match", e.Path)
			}
		}
	}
}
