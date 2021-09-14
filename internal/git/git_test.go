package git

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/utils/merkletrie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/compile"
	"go.uber.org/thriftrw/internal/compare"
)

func writeThrifts(t *testing.T, tmpDir string, contents map[string]string, worktree *git.Worktree, extra string) error {
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
		"test/b.thrift": `include "./a.thrift"`,
		"test/c.thrift": `include "./a.thrift"`,
		"test/d.thrift": `include "./b.thrift"
include "./c.thrift"`,
	}
	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, ""))

	exampleThrifts = map[string]string{
		"v1.thrift": "namespace rb v1\n" +
			"struct AddedRequiredField {\n" +
			"    1: optional string A\n" +
			"    2: optional string B\n" +
			"    3: required string C\n}\n" +
			"service Foo {}",
		"test/b.thrift": `include "./a.thrift"`,
		"test/c.thrift": `include "./a.thrift"`,
		"test/d.thrift": `include "./b.thrift"
include "./c.thrift"`,
	}

	require.NoError(t, writeThrifts(t, tmpDir, exampleThrifts, worktree, "second"))

	return tmpDir, repository
}

func TestOpenRepo(t *testing.T) {
	t.Parallel()
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()
	gitDir, repo := createRepoAndCommit(t, tmpDir)
	changed, err := findChangedThrift(gitDir)
	assert.NoError(t, err)
	assert.Equal(t, []*change{
		{file: "v1.thrift", change: merkletrie.Modify},
	}, changed)

	fs := NewGitFS(gitDir, repo, false)
	fsFrom := NewGitFS(gitDir, repo, true)

	for _, c := range changed {
		toModule, err := compile.Compile(c.file, compile.Filesystem(fs))
		require.NoError(t, err)
		fromModule, err := compile.Compile(c.file, compile.Filesystem(fsFrom))
		require.NoError(t, err)

		err = compare.Modules(fromModule, toModule)
		require.Error(t, err)
		assert.EqualError(t, err,
			"removing method methodA in service Foo is not backwards compatible;"+
				" adding a required field C to AddedRequiredField is not backwards compatible")
	}
}
