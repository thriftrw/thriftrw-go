package git

import (
	"context"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"go.uber.org/multierr"
	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/compare"
)

// NewGitFS creates an implementation of FS to use git to discover
// Thrift changes and previous version of a Thrift file.
func NewGitFS(gitDir string, repo *git.Repository, from bool) *FS {
	if repo == nil {
		var err error
		repo, err = git.PlainOpenWithOptions(gitDir, &git.PlainOpenOptions{
			DetectDotGit:          true,
			EnableDotGitCommonDir: true,
		})
		if err != nil {
			return nil
		}
	}

	return &FS{
		repo:   repo,
		gitDir: gitDir,
		from:   from,
	}
}

// UseGit takes a path to a git repository and returns errors between HEAD and HEAD~
// for any incompatible Thrift changes between the two shas.
func UseGit(path string) error {
	r, err := git.PlainOpenWithOptions(path, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return err
	}
	fs := NewGitFS(path, r, false)
	fsFrom := NewGitFS(path, r, true)

	changed, err := findChangedThrift(path)
	if err != nil {
		return err
	}
	var errs error
	for _, c := range changed {
		var toModule *compile.Module
		if c.change == merkletrie.Modify {
			toModule, err = compile.Compile(c.file, compile.Filesystem(fs))
			if err != nil {
				return err
			}
		} else if c.change == merkletrie.Delete {
			// something got deleted, so we are creating an empty module here.
			toModule = &compile.Module{
				Name: c.file,
			}
		}

		fromModule, err := compile.Compile(c.file, compile.Filesystem(fsFrom))
		if err != nil {
			return err
		}
		errs = multierr.Append(errs, compare.Modules(fromModule, toModule))
	}

	return errs
}

// FS holds reference to components needed for git FS.
type FS struct {
	repo   *git.Repository
	gitDir string
	from   bool // Whether we are looking for previous version.
}

func (fs FS) findCommit() (*object.Commit, error) {
	// Default is look at HEAD.
	refHead, err := fs.repo.Head()
	if err != nil {
		return nil, err
	}
	hash := refHead.Hash()

	// Return first commit.
	if !fs.from {
		return fs.repo.CommitObject(hash)
	}

	// Otherwise, we need to look at provided sha.
	commitIter, err := fs.repo.Log(&git.LogOptions{From: hash})
	if err != nil {
		return nil, err
	}
	// Skip one.
	_, err = commitIter.Next()
	if err != nil {
		return nil, err
	}

	return commitIter.Next()
}

func (fs FS) Read(path string) ([]byte, error) {
	commit, err := fs.findCommit()

	if err != nil {
		return nil, err
	}

	// filename is going to be the full path. We don't want that.
	filename, err := filepath.Rel(fs.gitDir, path)
	if err != nil {
		return nil, err
	}
	// It's possible that file was deleted and it will not exist.
	f, err := commit.File(filename)
	if err != nil {
		return nil, err
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

	return filepath.Join(fs.gitDir, p), nil
}

type change struct {
	file   string
	change merkletrie.Action
}

// findChangedThrift reads a git repo and finds any Thrift files that got changed
// between HEAD and previous commit.
func findChangedThrift(gitDir string) ([]*change, error) {
	// Init git repo.
	r, err := git.PlainOpenWithOptions(gitDir, &git.PlainOpenOptions{
		DetectDotGit:          true,
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		return nil, err
	}

	// Get Repo's HEAD
	refHead, err := r.Head() // *plumbing.Reference
	if err != nil {
		return nil, err
	}
	// Look at the log. Can't use filters as they need exact path which
	// we don't know yet.
	commitIter, err := r.Log(&git.LogOptions{From: refHead.Hash()})
	if err != nil {
		return nil, err
	}
	commit, err := commitIter.Next() // *object.Commit
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
		if filepath.Ext(from.Name) == ".thrift" {
			changed = append(changed, &change{
				file:   o.From.Name,
				change: a,
			})
		}
	}

	return changed, nil
}
