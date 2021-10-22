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
	"context"
	"fmt"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/compare"
)

// FS holds reference to components needed for git FS.
type FS struct {
	repo *git.Repository
	root string
	tree *object.Tree
}

// NewGitFS creates an implementation of FS to use git to discover
// Thrift changes and previous version of a Thrift file.
func NewGitFS(gitDir string, repo *git.Repository, tree *object.Tree) *FS {
	return &FS{
		repo: repo,
		root: gitDir,
		tree: tree,
	}
}

// Compare takes a path to a git repository and returns errors between HEAD and HEAD~
// for any incompatible Thrift changes between the two shas.
func Compare(path string) (compare.Pass, error) {
	var pass compare.Pass
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return pass, err
	}

	h, err := findChangedThrift(r)
	if err != nil {
		return pass, fmt.Errorf("failed to find changed thrift files: %w", err)
	}
	fs := NewGitFS(path, r, h.to)
	fsFrom := NewGitFS(path, r, h.from)
	var errs error
	for _, c := range h.changes {
		var toModule *compile.Module
		if c.change == merkletrie.Modify {
			toModule, err = compile.Compile(c.file, compile.Filesystem(fs))
			if err != nil {
				return pass, err
			}
		} else if c.change == merkletrie.Delete {
			// something got deleted, so we are creating an empty module here.
			toModule = &compile.Module{
				Name: c.file,
			}
		}

		fromModule, err := compile.Compile(c.file, compile.Filesystem(fsFrom))
		if err != nil {
			return pass, err
		}
		pass.CompareModules(fromModule, toModule)
	}
	// p will have lints as a field which we can sort in cli.

	return pass, errs
}

func (fs FS) Read(path string) ([]byte, error) {
	// filename is going to be the full path. We don't want that.
	filename, err := filepath.Rel(fs.root, path)
	if err != nil {
		return nil, err
	}

	// It's possible that file was deleted and it will not exist.
	f, err := fs.tree.File(filename)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", filename, err)
	}
	s, err := f.Contents()
	if err != nil {
		return nil, err
	}
	body := []byte(s)

	return body, nil
}

// Abs returns absolute path to a file.
func (fs FS) Abs(p string) (string, error) {
	// Sometimes p can be a full path already on includes, and sometimes it can be relative.
	if filepath.IsAbs(p) {
		return p, nil
	}

	return filepath.Join(fs.root, p), nil
}

type treeChanges struct {
	from    *object.Tree
	to      *object.Tree
	changes []*change
}

type change struct {
	file   string
	change merkletrie.Action
}

// findChangedThrift reads a git repo and finds any Thrift files that got changed
// between HEAD and previous commit.
func findChangedThrift(r *git.Repository) (*treeChanges, error) {
	// Get Repo's HEAD
	refHead, err := r.Head() // *plumbing.Reference
	if err != nil {
		return nil, err
	}
	commit, err := r.CommitObject(refHead.Hash()) // *object.Commit
	if err != nil {
		return nil, err
	}
	parentCommit, err := commit.Parent(0) // *object.Commit
	if err != nil {
		return nil, err
	}
	// Get the two commit trees.
	c, err := commit.Tree() // *object.Tree
	if err != nil {
		return nil, err
	}
	pc, err := parentCommit.Tree() // *object.Tree
	if err != nil {
		return nil, err
	}
	// Diff the trees and find what changed.
	objects, _ := object.DiffTreeWithOptions(context.Background(),
		pc, c, &object.DiffTreeOptions{DetectRenames: true}) // *object.Changes
	var changed []*change
	for _, o := range objects {
		a, err := o.Action() // Insert, delete or modify.
		if err != nil {
			return nil, err
		}
		from, _, _ := o.Files()
		// New file was added which doesnt have a name.
		if from == nil {
			continue
		}
		if filepath.Ext(from.Name) == ".thrift" {
			changed = append(changed, &change{
				file:   o.From.Name,
				change: a,
			})
		}
	}

	return &treeChanges{
		from:    pc,
		to:      c,
		changes: changed}, nil
}
