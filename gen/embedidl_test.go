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
	"encoding/hex"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	te "go.uber.org/thriftrw/gen/testdata/enums"
	ts "go.uber.org/thriftrw/gen/testdata/structs"
)

func loadIDL(filename string) (string, string, error) {
	f, err := os.Open("./testdata/thrift/" + filename)
	if err != nil {
		return "", "", err
	}
	rawIDL, err := ioutil.ReadAll(f)
	if err != nil {
		return "", "", err
	}
	hash := sha1.Sum(rawIDL)
	return string(rawIDL), hex.EncodeToString(hash[:]), nil
}

func TestIDLEmbeddingEnums(t *testing.T) {
	tm := te.ThriftModule
	assert.Equal(t, "enums", tm.Name)
	assert.True(t, strings.HasSuffix(tm.Package, "thriftrw/gen/testdata/enums"))
	assert.Equal(t, "enums.thrift", tm.FilePath)

	rawIDL, sha1, err := loadIDL("enums.thrift")
	if assert.NoError(t, err) {
		assert.Equal(t, rawIDL, tm.Raw)
		assert.Equal(t, sha1, tm.SHA1)
	}
}

func TestIDLEmbeddingStructs(t *testing.T) {
	tm := ts.ThriftModule
	assert.Equal(t, "structs", tm.Name)
	assert.True(t, strings.HasSuffix(tm.Package, "thriftrw/gen/testdata/structs"))
	assert.Equal(t, "structs.thrift", tm.FilePath)

	rawIDL, sha1, err := loadIDL("structs.thrift")
	if assert.NoError(t, err) {
		assert.Equal(t, rawIDL, tm.Raw)
		assert.Equal(t, sha1, tm.SHA1)
	}

	if assert.Equal(t, 1, len(tm.Includes)) {
		assert.Equal(t, te.ThriftModule, tm.Includes[0])
	}
}
