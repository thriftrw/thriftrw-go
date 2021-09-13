package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThriftBreak(t *testing.T) {
	t.Parallel()
	t.Run("wrong flag", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--out_file=tests"})
		require.Error(t, err)
		assert.EqualError(t, err, "flag provided but not defined: -out_file")
	})
	t.Run("wrong file name", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/something.thrift"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no such file")
	})
	t.Run("invalid thrift", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/invalid.thrift"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse file")
	})
	t.Run("invalid thrift for from_file", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/v1.thrift", "--from_file=tests/invalid.thrift"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "could not parse file")
	})
	t.Run("missing to_file", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--from_file=tests/something.thrift"})
		require.Error(t, err)
		assert.EqualError(t, err, "must provide an updated Thrift file")
	})
	t.Run("integration test all errors", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/v2.thrift", "--from_file=tests/v1.thrift"})
		require.Error(t, err)

		assert.EqualError(t, err, "removing method methodA in service Foo is not backwards compatible;"+
			" deleting service Bar is not backwards compatible;"+
			" changing an optional field B in AddedRequiredField to required is not backwards compatible;"+
			" adding a required field C to AddedRequiredField is not backwards compatible",
		)
	})
	t.Run("integration test single method", func(t *testing.T) {
		t.Parallel()
		err := run([]string{"--to_file=tests/v3.thrift", "--from_file=tests/v1.thrift"})
		require.Error(t, err)

		assert.EqualError(t, err, "removing method methodA in service Foo is not backwards compatible")
	})
}

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
func createRepoAndCommit(t *testing.T, tmpDir string) (string, *git.Repository) {
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
		"test/d.thrift":  `service Qux {}`,
		"somefile.go":    `service Quux{}`,
	}
	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, nil, ""))

	exampleThrifts = map[string]string{
		"v1.thrift": "namespace rb v1\n" +
			"struct AddedRequiredField {\n" +
			"    1: optional string A\n" +
			"    2: optional string B\n" +
			"    3: required string C\n}\n" +
			"service Foo {}",
		"test/v2.thrift": `service Foo {}`,
		"test/c.thrift":  `service Bar {}`,
		"somefile.go":    `service Qux{}`, // Change name for Go file.
	}

	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, []string{"test/d.thrift"}, "second"))

	return tmpDir, repository
}

func TestThriftBreakIntegration(t *testing.T) {
	t.Parallel()
	t.Run("integration test git repo", func(t *testing.T) {
		t.Parallel()
		tmpDir, err := ioutil.TempDir("", "")
		require.NoError(t, err)

		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()
		gitDir, _ := createRepoAndCommit(t, tmpDir)
		assert.NoError(t, err)

		err = run([]string{fmt.Sprintf("--git_repo=%s", gitDir)})
		require.Error(t, err, "expected lint errors")
		assert.EqualError(t, err, "deleting service Baz is not backwards compatible;"+
			" deleting service Qux is not backwards compatible;"+
			" deleting service Bar is not backwards compatible;"+
			" removing method methodA in service Foo is not backwards compatible;"+
			" adding a required field C to AddedRequiredField is not backwards compatible")
	})
}
