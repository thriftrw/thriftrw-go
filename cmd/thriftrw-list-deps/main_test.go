package main

import (
	"os"
	"path/filepath"
	"testing"

	"sort"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecursiveImports(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)
	idlRoot := filepath.Join(cwd, "../../gen/testdata/thrift")

	t.Run("default options", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(idlRoot, "structs.thrift"), idlRoot)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 1, len(outputLines))
		assert.Equal(t, "enums.thrift", outputLines[0])
	})
	t.Run("with relative path", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(idlRoot, "structs.thrift"), filepath.Join(cwd, "../../"))
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 1, len(outputLines))
		assert.Equal(t, "gen/testdata/thrift/enums.thrift", outputLines[0])
	})
	t.Run("with transitive dependency", func(t *testing.T) {
		outputLines, err := listDependentThrifts(filepath.Join(idlRoot, "unions.thrift"), idlRoot)
		require.NoError(t, err)
		sort.Strings(outputLines)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 3, len(outputLines))
		assert.Equal(t, "enums.thrift", outputLines[0])
		assert.Equal(t, "structs.thrift", outputLines[1])
		assert.Equal(t, "typedefs.thrift", outputLines[2])
	})
	t.Run("with open error", func(t *testing.T) {
		_, err := listDependentThrifts("/does-not-exist", "")
		require.Error(t, err)
	})
	t.Run("with relative error", func(t *testing.T) {
		_, err := listDependentThrifts(filepath.Join(idlRoot, "unions.thrift"), "./")
		require.Error(t, err)
	})
}
