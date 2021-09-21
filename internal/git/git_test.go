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
	"go.uber.org/thriftrw/thriftbreaktest"
)

// func writeThrifts(t *testing.T, tmpDir string, contents map[string]string, worktree *git.Worktree, extra string) error {
// 	t.Helper()
// 	for name, content := range contents {
// 		path := filepath.Join(tmpDir, name)
// 		err := os.MkdirAll(filepath.Dir(path), 0755)
// 		if err != nil {
// 			return err
// 		}
// 		err = ioutil.WriteFile(path, []byte(content), 0600)
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	err := worktree.AddWithOptions(&git.AddOptions{All: true})
// 	if err != nil {
// 		return err
// 	}
//
// 	_, err = worktree.Commit("thrift update file"+extra, &git.CommitOptions{
// 		Author: &object.Signature{
// 			Name:  "update v1.thrift",
// 			Email: "thriftforeverornever@uber.com",
// 			When:  time.Now(),
// 		},
// 	})
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// createRepoAndCommit creates a temporary repository and adds
// a commit of a thrift file for us to look up later.
// func createRepoAndCommit(t *testing.T, tmpDir string) (string, *git.Repository) {
// 	t.Helper()
// 	// Create a new repo in temp directory.
// 	repository, err := git.PlainInit(tmpDir, false)
// 	require.NoError(t, err)
// 	worktree, err := repository.Worktree()
// 	require.NoError(t, err)
// 	exampleThrifts := map[string]string{
// 		"v1.thrift": "namespace rb v1\n" +
// 			"struct AddedRequiredField {\n" +
// 			"    1: optional string A\n" +
// 			"    2: optional string B\n" +
// 			"}\n" +
// 			"\nservice Foo {\n    void methodA()\n}",
// 		"test/b.thrift": `include "./a.thrift"`,
// 		"test/c.thrift": `include "./a.thrift"`,
// 		"test/d.thrift": `include "./b.thrift"
// include "./c.thrift"`,
// 	}
// 	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, ""))
//
// 	exampleThrifts = map[string]string{
// 		"v1.thrift": "namespace rb v1\n" +
// 			"struct AddedRequiredField {\n" +
// 			"    1: optional string A\n" +
// 			"    2: optional string B\n" +
// 			"    3: required string C\n}\n" +
// 			"service Foo {}",
// 		"test/b.thrift": `include "./a.thrift"`,
// 		"test/c.thrift": `include "./a.thrift"`,
// 		"test/d.thrift": `include "./b.thrift"
// include "./c.thrift"`,
// 	}
//
// 	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, "second"))
//
// 	return tmpDir, repository
// }

func TestOpenRepo(t *testing.T) {
	t.Parallel()
	tmpDir := t.TempDir()
	thriftbreaktest.CreateRepoAndCommit(t, tmpDir)
	changed, err := findChangedThrift(tmpDir)
	assert.NoError(t, err)
	assert.Equal(t, []*change{
		{file: "test/c.thrift", change: merkletrie.Modify},
		{file: "test/d.thrift", change: merkletrie.Delete},
		{file: "test/v2.thrift", change: merkletrie.Modify},
		{file: "v1.thrift", change: merkletrie.Modify},
	}, changed)

	err = Compare(tmpDir)
	require.Error(t, err)
	assert.EqualError(t, err,
		"deleting service Baz is not backwards compatible;" +
		" deleting service Qux is not backwards compatible;" +
		" deleting service Bar is not backwards compatible;" +
		" removing method methodA in service Foo is not backwards compatible;"+
		" adding a required field C to AddedRequiredField is not backwards compatible")
}
