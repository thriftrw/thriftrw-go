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

package thriftbreaktest

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/require"
)

func writeThrifts(t *testing.T, tmpDir string, contents map[string]string, worktree *git.Worktree, toRemove []string,
	extra string) error {
	t.Helper()
	for name, content := range contents {
		path := filepath.Join(tmpDir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, []byte(content), 0600)
		if err != nil {
			return err
		}
	}
	err := worktree.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}
	for _, f := range toRemove {
		_, err := worktree.Remove(f)
		if err != nil {
			return err
		}
	}

	_, err = worktree.Commit("thrift update file"+extra, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "update v1.thrift",
			Email: "thriftforeverornever@uber.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// createRepoAndCommit creates a temporary repository and adds
// a commit of a thrift file for us to look up later.
func CreateRepoAndCommit(t *testing.T, tmpDir string) {
	t.Helper()
	// Create a new repo in temp directory.
	repository, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)
	worktree, err := repository.Worktree()
	require.NoError(t, err)
	exampleThrifts := map[string]string{
		"v1.thrift": "namespace rb v1\n" +
			"struct AddedRequiredField {\n" +
			"    1: optional string A\n" +
			"    2: optional string B\n" +
			"}\n" +
			"\nservice Foo {\n    void methodA()\n}",
		"test/v2.thrift": `service Bar {}`,
		"test/c.thrift":  `service Baz {}`,
		"test/d.thrift": `include "../v1.thrift"
service Qux {}`,
		"somefile.go": `service Quux{}`,
	}
	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, nil, ""))

	// For c.thrift we are also checking to make sure includes work as expected.
	exampleThrifts = map[string]string{
		"v1.thrift": "namespace rb v1\n" +
			"struct AddedRequiredField {\n" +
			"    1: optional string A\n" +
			"    2: optional string B\n" +
			"    3: required string C\n}\n" +
			"service Foo {}",
		"test/v2.thrift": `service Foo {}`,
		"test/c.thrift": `include "../v1.thrift"
service Bar {}`,
		"somefile.go": `service Qux{}`, // Change name for Go file.
	}

	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, []string{"test/d.thrift"}, "second"))
}
