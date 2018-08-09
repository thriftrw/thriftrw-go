package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecursiveImports(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		cwd, err := os.Getwd()
		require.NoError(t, err)
		idlRoot := filepath.Join(cwd, "../../gen/testdata/thrift")
		outputLines, err := listDependentThrifts(filepath.Join(idlRoot, "structs.thrift"), idlRoot)
		require.NoError(t, err)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 1, len(outputLines))
		assert.Equal(t, "enums.thrift", outputLines[0])
	})
	t.Run("with relative path", func(t *testing.T) {
		cwd, err := os.Getwd()
		require.NoError(t, err)
		idlRoot := filepath.Join(cwd, "../../gen/testdata/thrift")
		outputLines, err := listDependentThrifts(filepath.Join(idlRoot, "structs.thrift"), filepath.Join(cwd, "../../"))
		require.NoError(t, err)
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 1, len(outputLines))
		assert.Equal(t, "gen/testdata/thrift/enums.thrift", outputLines[0])
	})
}
