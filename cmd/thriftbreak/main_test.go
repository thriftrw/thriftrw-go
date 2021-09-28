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
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/internal/breaktest"
)

func TestThriftBreakIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	breaktest.CreateRepoAndCommit(t, tmpDir)

	f, err := ioutil.TempFile(tmpDir, "stdout")
	require.NoError(t, err, "create temporary file")
	defer func(oldStdout *os.File) {
		assert.NoError(t, f.Close())
		os.Stdout = oldStdout
	}(os.Stdout)
	os.Stdout = f

	err = run([]string{fmt.Sprintf("-C=%s", tmpDir)})
	require.Error(t, err, "expected no errors")
	assert.EqualError(t, err, "found 5 issues")

	stderr, err := ioutil.ReadFile(f.Name())
	require.NoError(t, err)

	out := string(stderr)

	assert.Equal(t,
		`c.thrift:deleting service "Baz"`+"\n"+
			`d.thrift:deleting service "Qux"`+"\n"+
			`v2.thrift:deleting service "Bar"`+"\n"+
			`v1.thrift:removing method "methodA" in service "Foo"`+"\n"+
			`v1.thrift:adding a required field "C" to "AddedRequiredField"`+"\n",
		out)
}
