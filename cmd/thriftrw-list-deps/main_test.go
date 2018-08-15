package main

import (
	"path/filepath"
	"testing"

	"sort"

	"io/ioutil"

	"os"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestThriftrwListDeps(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "")
	require.NoError(t, err)

	defer os.RemoveAll(tmpDir)

	exampleThrifts := map[string]string{
		"test/a.thrift": "",
		"test/b.thrift": `include "./a.thrift"`,
		"test/c.thrift": `include "./a.thrift"`,
		"test/d.thrift": `include "./b.thrift"
include "./c.thrift"`,
	}

	for name, content := range exampleThrifts {
		path := filepath.Join(tmpDir, name)
		err = os.MkdirAll(filepath.Dir(path), 0755)
		require.NoError(t, err)
		err = ioutil.WriteFile(path, []byte(content), 0644)
		require.NoError(t, err)
	}

	t.Run("no dependencies", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/a.thrift"), tmpDir)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Empty(t, outputLines)
	})
	t.Run("one dependency", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/b.thrift"), tmpDir)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, []string{"test/a.thrift"}, outputLines)
	})
	t.Run("one dependency with relative path", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/b.thrift"), filepath.Join(tmpDir, "test"))
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, []string{"a.thrift"}, outputLines)
	})
	t.Run("transitive and multiple dependencies", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/d.thrift"), tmpDir)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, []string{"test/a.thrift", "test/b.thrift", "test/c.thrift"}, outputLines)
	})
	t.Run("with open error", func(t *testing.T) {
		_, err := listDependentThrifts("/does-not-exist", "")
		require.Error(t, err)
	})
	t.Run("with relative error", func(t *testing.T) {
		_, err := listDependentThrifts(filepath.Join(tmpDir, "test/b.thrift"), "./")
		require.Error(t, err)
	})
}
