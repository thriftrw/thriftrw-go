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
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/compile"
)

// Header for generated code as per https://golang.org/s/generatedcode
var generatedByRegex = regexp.MustCompile(`^// Code generated .* DO NOT EDIT\.$`)

func TestCodeIsUpToDate(t *testing.T) {
	// This test just verifies that the generated code in internal/tests/ is up to
	// date. If this test failed, run 'make' in the internal/tests/ directory and
	// commit the changes.
	var (
		// Set of files that are passed a --no-zap flag in code generation
		noZapFiles = map[string]struct{}{
			"nozap": {},
		}
		// Set of files that are passed a --no-error flag in code generation
		noErrorFiles = map[string]struct{}{
			"noerror": {},
		}
	)

	thriftRoot, err := filepath.Abs("internal/tests/thrift")
	require.NoError(t, err, "could not resolve absolute path to internal/tests/thrift")

	thriftFiles, err := filepath.Glob(thriftRoot + "/*.thrift")
	require.NoError(t, err)

	outputDir, err := ioutil.TempDir("", "thriftrw-golden-test")
	require.NoError(t, err)
	defer os.RemoveAll(outputDir)

	for _, thriftFile := range thriftFiles {
		pkgRelPath := strings.TrimSuffix(filepath.Base(thriftFile), ".thrift")
		currentPackageDir := filepath.Join("internal/tests", pkgRelPath)
		newPackageDir := filepath.Join(outputDir, pkgRelPath)

		currentHash, err := dirhash(currentPackageDir)
		require.NoError(t, err, "could not hash %q", currentPackageDir)

		module, err := compile.Compile(thriftFile)
		require.NoError(t, err, "failed to compile %q", thriftFile)

		_, nozap := noZapFiles[pkgRelPath]
		_, noerror := noErrorFiles[pkgRelPath]
		err = Generate(module, &Options{
			OutputDir:     outputDir,
			PackagePrefix: "go.uber.org/thriftrw/gen/internal/tests",
			ThriftRoot:    thriftRoot,
			NoRecurse:     true,
			NoZap:         nozap,
			NoError:       noerror,
		})
		require.NoError(t, err, "failed to generate code for %q", thriftFile)

		// All generated Go files must have a line that matches
		// generatedByRegex.
		err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() || !strings.HasSuffix(path, ".go") {
				return nil
			}

			f, err := os.Open(path)
			if !assert.NoError(t, err, "failed to open %q", path) {
				return err
			}
			defer f.Close()

			scanner := bufio.NewScanner(f)
			if scanner.Scan() {
				// Check the first line only if the file is non-empty.
				line := scanner.Text()
				assert.Regexp(t, generatedByRegex, line,
					"first line of %q does not have the correct header", path)
			}

			err = scanner.Err()
			assert.NoError(t, err, "failed to scan %q", path)
			return err
		})
		require.NoError(t, err)

		newHash, err := dirhash(newPackageDir)
		require.NoError(t, err, "could not hash %q", newPackageDir)

		if newHash != currentHash {
			// TODO(abg): Diff the two directories?
			t.Fatalf(
				"Generated code for %q is out of date. "+
					"Please run 'make' in gen/internal/tests.", thriftFile)
		}
	}
}
