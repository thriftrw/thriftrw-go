// Copyright (c) 2015 Uber Technologies, Inc.
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
	"github.com/thriftrw/thriftrw-go/compile"
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

func generate(t *testing.T, inPath, outPath string) {
	module, err := compile.Compile(inPath)
	require.NoError(t, err, "failed to compile %v", inPath)

	outFile, err := os.Create(outPath)
	require.NoError(t, err, "could not create %v", outPath)
	defer outFile.Close()

	opts := Options{PackageName: "testdata", Output: outFile}
	require.NoError(
		t, Generate(module, &opts), "could not generate code for %v", inPath)
}

func TestCodeIsUpToDate(t *testing.T) {
	// This test just verifies that the generated code in testdata/ is up to
	// date. If this test failed, run 'make generate' in the testdata/ directory
	// and commit the changes.

	files, err := filepath.Glob("testdata/*.thrift")
	require.NoError(t, err)

	outputDir, err := ioutil.TempDir("", "thriftrw-golden-test")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	for _, file := range files {
		currentPath := file + ".go"
		currentHash := hash(t, currentPath)
		newPath := filepath.Join(outputDir, filepath.Base(file)+".go")
		generate(t, file, newPath)

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
