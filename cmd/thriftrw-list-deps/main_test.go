package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecursiveImports(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		outputBuf := new(bytes.Buffer)
		cwd, err := os.Getwd()
		require.NoError(t, err)
		idlRoot := filepath.Join(cwd, "../../gen/testdata/thrift")
		err = do(idlRoot, filepath.Join(idlRoot, "structs.thrift"), outputBuf)
		require.NoError(t, err)
		outputLines := strings.Split(outputBuf.String(), "\n")
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 2, len(outputLines))
		assert.Equal(t, "enums.thrift", outputLines[0])
		assert.Equal(t, "", outputLines[1])
	})
	t.Run("with relative path", func(t *testing.T) {
		outputBuf := new(bytes.Buffer)
		cwd, err := os.Getwd()
		require.NoError(t, err)
		idlRoot := filepath.Join(cwd, "../../gen/testdata/thrift")
		err = do(filepath.Join(cwd, "../../"), filepath.Join(idlRoot, "structs.thrift"), outputBuf)
		require.NoError(t, err)
		outputLines := strings.Split(outputBuf.String(), "\n")
		t.Logf("output lines: %+v", outputLines)
		assert.Equal(t, 2, len(outputLines))
		assert.Equal(t, "gen/testdata/thrift/enums.thrift", outputLines[0])
		assert.Equal(t, "", outputLines[1])
	})
}

func ExampleUsage() {
	cwd, _ := os.Getwd()
	thriftrwRoot := filepath.Join(cwd, "../..")
	do(thriftrwRoot, filepath.Join(thriftrwRoot, "gen/testdata/thrift/structs.thrift"), os.Stdout)
	// Output:
	// gen/testdata/thrift/enums.thrift
}
