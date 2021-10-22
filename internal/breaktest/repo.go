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

package breaktest

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

type writeThrift struct {
	tmpDir   string
	contents map[string]string
	worktree *git.Worktree
	toRemove []string
}

func newWriteThrift(tmpDir string, contents map[string]string, worktree *git.Worktree, toRemove []string) *writeThrift {
	return &writeThrift{
		tmpDir:   tmpDir,
		contents: contents,
		worktree: worktree,
		toRemove: toRemove,
	}
}

// commit commits all changes staged before it is called.
func (w *writeThrift) commit(extra string) error {
	err := w.worktree.AddWithOptions(&git.AddOptions{All: true})
	if err != nil {
		return err
	}
	for _, f := range w.toRemove {
		if _, err := w.worktree.Remove(f); err != nil {
			return err
		}
	}

	_, err = w.worktree.Commit("thrift update file"+extra, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "update v1.thrift",
			Email: "thriftforeverornever@example.com",
			When:  time.Now(),
		},
	})

	return err
}

func (w *writeThrift) writeThrifts(extraMsg string) error {
	for name, content := range w.contents {
		path := filepath.Join(w.tmpDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, []byte(content), 0600); err != nil {
			return err
		}
	}

	return w.commit(extraMsg)
}

// CreateRepoAndCommit creates a temporary repository and adds
// a commit of a thrift files for us to look up later.
// TODO(GO-891): finish implementation of this integration.
func CreateRepoAndCommit(t *testing.T, tmpDir string, from map[string]string,
	to map[string]string, remove []string) *git.Repository {
	t.Helper()
	// Create a new repo in temp directory.
	repository, err := git.PlainInit(tmpDir, false)
	require.NoError(t, err)
	worktree, err := repository.Worktree()
	require.NoError(t, err)

	w := newWriteThrift(tmpDir, from, worktree, nil)
	require.NoError(t, w.writeThrifts(""))
	w = newWriteThrift(tmpDir, to, worktree, remove)
	require.NoError(t, w.writeThrifts("second"))

	return repository
}
