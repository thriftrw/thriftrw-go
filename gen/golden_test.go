package gen

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func hash(t *testing.T, name string) string {
	f, err := os.Open(name)
	require.NoError(t, err)
	defer f.Close()

	h := sha1.New()
	_, err = io.Copy(h, f)
	require.NoError(t, err)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func run(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func TestCodeIsUpToDate(t *testing.T) {
	// This test just verifies that the generated code in testdata/ is up to
	// date. If this test failed, run 'make generate' in the testdata/ directory
	// and commit the changes.

	files, err := filepath.Glob("testdata/*.thrift")
	require.NoError(t, err)

	tmpDir, err := ioutil.TempDir("", "thriftrw-golden-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	outputDir := filepath.Join(tmpDir, "testdata")
	require.NoError(t, os.Mkdir(outputDir, 755))

	for _, file := range files {
		currentPath := file + ".go"
		currentHash := hash(t, currentPath)

		newPath := filepath.Join(outputDir, filepath.Base(file)+".go")
		require.NoError(t, run("../thriftrw-go", "-o", newPath, file))

		if hash(t, newPath) != currentHash {
			run("diff", "-u", currentPath, newPath)
			t.Fatalf(
				"Generated code for %s is out of date. "+
					"Please run 'make generate' in gen/testdata.",
				file,
			)
		}
	}
}
