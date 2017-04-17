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
	"go.uber.org/thriftrw/thriftreflect"
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

func TestIDLEmbedding(t *testing.T) {
	for _, tt := range []struct {
		N  string
		TM *thriftreflect.ThriftModule
	}{
		{
			"enums",
			te.ThriftModule,
		},
		{
			"structs",
			ts.ThriftModule,
		},
	} {
		tm := tt.TM
		assert.Equal(t, tt.N, tm.Name)
		assert.True(t, strings.HasSuffix(tm.Package, "thriftrw/gen/testdata/"+tt.N))
		assert.Equal(t, tt.N+".thrift", tm.FilePath)

		rawIDL, sha1, err := loadIDL(tt.N + ".thrift")
		if assert.NoError(t, err) {
			assert.Equal(t, rawIDL, tm.Raw)
			assert.Equal(t, sha1, tm.SHA1)
		}
	}
}

func TestIDLEmbeddingInclude(t *testing.T) {
	tm := ts.ThriftModule
	if assert.Equal(t, 1, len(tm.Includes)) {
		assert.Equal(t, te.ThriftModule, tm.Includes[0])
	}
}
