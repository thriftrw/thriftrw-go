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
	if err != nil {
		t.Skipf("error creating temporary directory: %v", err)
	}

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
		if err != nil {
			t.Skipf("error mkdir %s: %v", path, err)
		}
		err = ioutil.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Skipf("error writing %s: %v", path, err)
		}
	}

	t.Run("no dependencies", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/a.thrift"), tmpDir)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 0, len(outputLines))
	})
	t.Run("one dependency", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/b.thrift"), tmpDir)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 1, len(outputLines))
		assert.Equal(t, "test/a.thrift", outputLines[0])
	})
	t.Run("one dependency with relative path", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/b.thrift"), filepath.Join(tmpDir, "test"))
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 1, len(outputLines))
		assert.Equal(t, "a.thrift", outputLines[0])
	})
	t.Run("transitive and multiple dependencies", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(tmpDir, "test/d.thrift"), tmpDir)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 3, len(outputLines))
		assert.Equal(t, "test/a.thrift", outputLines[0])
		assert.Equal(t, "test/b.thrift", outputLines[1])
		assert.Equal(t, "test/c.thrift", outputLines[2])
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
