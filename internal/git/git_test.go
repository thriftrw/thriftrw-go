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

package git

import (
	"testing"

	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/internal/breaktest"
)

func TestOpenRepo(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	breaktest.CreateRepoAndCommit(t, tmpDir)
	changed, err := findChangedThrift(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, []*change{
		{file: "test/c.thrift", change: merkletrie.Modify},
		{file: "test/d.thrift", change: merkletrie.Delete},
		{file: "test/v2.thrift", change: merkletrie.Modify},
		{file: "v1.thrift", change: merkletrie.Modify},
	}, changed)

	pass, err := Compare(tmpDir)
	require.NoError(t, err)
	assert.Equal(t,
		"c.thrift:deleting service Baz is not backwards compatible\n"+
			"d.thrift:deleting service Qux is not backwards compatible\n"+
			"v2.thrift:deleting service Bar is not backwards compatible\n"+
			"v1.thrift:removing method methodA in service Foo is not backwards compatible\n"+
			"v1.thrift:adding a required field C to AddedRequiredField is not backwards compatible\n",
		pass.String())
}
