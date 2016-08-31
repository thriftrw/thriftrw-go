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

package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/thriftrw/compile"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeterminePackagePrefix(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "thriftrw-main-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	genDir, err := filepath.Abs("gen")
	require.NoError(t, err)

	tests := []struct {
		overrideGoPath func(curGoPath string) string

		dir    string
		result string
		errMsg string
	}{
		{
			overrideGoPath: func(gopath string) string { return "" },
			dir:            genDir,
			errMsg:         "$GOPATH is not set",
		},
		{
			dir:    tmpDir,
			errMsg: "not inside $GOPATH",
		},
		{
			dir:    genDir,
			result: "go.uber.org/thriftrw/gen",
		},
		{
			overrideGoPath: func(gopath string) string {
				return tmpDir + ":" + gopath
			},
			dir:    filepath.Join(tmpDir, "src/go.uber.org/thriftrw"),
			result: "go.uber.org/thriftrw",
		},
		{
			overrideGoPath: func(gopath string) string {
				return tmpDir + ":" + gopath
			},
			dir:    genDir,
			result: "go.uber.org/thriftrw/gen",
		},
	}

	realGoPath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", realGoPath)

	for _, tt := range tests {
		require.NoError(t, os.Setenv("GOPATH", realGoPath))

		if tt.overrideGoPath != nil {
			fakeGoPath := tt.overrideGoPath(realGoPath)
			if fakeGoPath == "" {
				require.NoError(t, os.Unsetenv("GOPATH"))
			} else {
				require.NoError(t, os.Setenv("GOPATH", fakeGoPath))
			}
		}

		result, err := determinePackagePrefix(tt.dir)

		if tt.result != "" {
			if assert.NoError(t, err) {
				assert.Equal(t, tt.result, result)
			}
		}

		if tt.errMsg != "" {
			if assert.Error(t, err, "determinePackagePrefix(%q) should fail", tt.dir) {
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		}
	}
}

func TestCommonPrefix(t *testing.T) {
	tests := []struct {
		left     []string
		right    []string
		expected []string
	}{
		{
			[]string{},
			[]string{"foo"},
			[]string{},
		},
		{
			[]string{"foo"},
			[]string{},
			[]string{},
		},
		{
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "bar", "qux"},
			[]string{"foo", "bar"},
		},
		{
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "bar", "baz", "qux"},
			[]string{"foo", "bar", "baz"},
		},
		{
			[]string{"foo", "bar", "baz", "qux"},
			[]string{"foo", "bar", "baz"},
			[]string{"foo", "bar", "baz"},
		},
	}

	for _, tt := range tests {
		assert.Equal(
			t, tt.expected, commonPrefix(tt.left, tt.right),
			"commonPrefix(%v, %v) != %v", tt.left, tt.right, tt.expected)
	}
}

func TestVerifyAncestry(t *testing.T) {
	cyclicFoo := &compile.Module{
		Name:       "foo",
		ThriftPath: "/tmp/service/foo.thrift",
		Includes: map[string]*compile.IncludedModule{
			"bar": {
				Name: "bar",
				Module: &compile.Module{
					Name:       "bar",
					ThriftPath: "/tmp/service/foo/bar.thrift",
				},
			},
		},
	}
	cyclicFoo.Includes["bar"].Module.Includes =
		map[string]*compile.IncludedModule{
			"foo": {
				Name:   "foo",
				Module: cyclicFoo,
			},
		}

	tests := []struct {
		desc   string
		module *compile.Module
		root   string
		errMsg string
	}{
		{
			desc: "success without includes",
			module: &compile.Module{
				Name:       "foo.thrift",
				ThriftPath: "/tmp/service/foo.thrift",
			},
			root: "/tmp/service",
		},
		{
			desc: "success with includes",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/tmp/service/foo/bar.thrift",
						},
					},
					"baz": {
						Name: "baz",
						Module: &compile.Module{
							Name:       "baz",
							ThriftPath: "/tmp/service/baz.thrift",
						},
					},
				},
			},
			root: "/tmp/service",
		},
		{
			desc:   "success with cyclic includes",
			module: cyclicFoo,
			root:   "/tmp/service",
		},
		{
			desc: "fail without includes",
			module: &compile.Module{
				Name:       "foo.thrift",
				ThriftPath: "/tmp/service/foo.thrift",
			},
			root:   "/tmp/anotherService",
			errMsg: `not contained in the "/tmp/anotherService" directory`,
		},
		{
			desc: "fail with includes",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/tmp/service2/bar.thrift",
						},
					},
				},
			},
			root:   "/tmp/service",
			errMsg: `"/tmp/service2/bar.thrift" is not contained in the "/tmp/service" directory`,
		},
	}

	for _, tt := range tests {
		err := verifyAncestry(tt.module, tt.root)
		if tt.errMsg != "" {
			if assert.Error(t, err, tt.desc) {
				assert.Contains(t, err.Error(), tt.errMsg, tt.desc)
			}
		} else {
			assert.NoError(t, err, tt.desc)
		}
	}
}

func TestFindCommonAncestor(t *testing.T) {
	cyclicFoo := &compile.Module{
		Name:       "foo",
		ThriftPath: "/tmp/service/foo.thrift",
		Includes: map[string]*compile.IncludedModule{
			"bar": {
				Name: "bar",
				Module: &compile.Module{
					Name:       "bar",
					ThriftPath: "/tmp/service/foo/bar.thrift",
				},
			},
		},
	}
	cyclicFoo.Includes["bar"].Module.Includes =
		map[string]*compile.IncludedModule{
			"foo": {
				Name:   "foo",
				Module: cyclicFoo,
			},
		}

	tests := []struct {
		desc     string
		module   *compile.Module
		expected string
		errMsg   string
	}{
		{
			desc: "success: no includes",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
			},
			expected: "/tmp/service",
		},
		{
			desc: "success: include sibling",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/tmp/service/bar.thrift",
						},
					},
				},
			},
			expected: "/tmp/service",
		},
		{
			desc: "success: include child",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/tmp/service/common/bar.thrift",
						},
					},
				},
			},
			expected: "/tmp/service",
		},
		{
			desc: "success: include multiple levels",
			module: &compile.Module{
				Name:       "service",
				ThriftPath: "/tmp/service/foo/service.thrift",
				Includes: map[string]*compile.IncludedModule{
					"common": {
						Name: "common",
						Module: &compile.Module{
							Name:       "common",
							ThriftPath: "/tmp/service/shared/types/common.thrift",
						},
					},
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/tmp/service/bar/bar.thrift",
							Includes: map[string]*compile.IncludedModule{
								"common": {
									Name: "common",
									Module: &compile.Module{
										Name:       "common",
										ThriftPath: "/tmp/service/shared/types/common.thrift",
									},
								},
							},
						},
					},
				},
			},
			expected: "/tmp/service",
		},
		{
			desc: "success: include parent",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/tmp/common/bar.thrift",
						},
					},
				},
			},
			expected: "/tmp",
		},
		{
			desc:     "success: include cyclic",
			module:   cyclicFoo,
			expected: "/tmp/service",
		},
		{
			desc: "failure: relative path",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "service/foo.thrift",
			},
			errMsg: `"service/foo.thrift" is not absolute`,
		},
		{
			desc: "failure: relative include",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "common/bar.thrift",
						},
					},
				},
			},
			errMsg: `"common/bar.thrift" is not absolute`,
		},
		{
			desc: "failure: different trees",
			module: &compile.Module{
				Name:       "foo",
				ThriftPath: "/tmp/service/foo.thrift",
				Includes: map[string]*compile.IncludedModule{
					"bar": {
						Name: "bar",
						Module: &compile.Module{
							Name:       "bar",
							ThriftPath: "/home/thriftrw/common/shared.thrift",
						},
					},
				},
			},
			errMsg: `"/home/thriftrw/common/shared.thrift" does not share an ancestor with "/tmp/service/foo.thrift"`,
		},
	}

	for _, tt := range tests {
		got, err := findCommonAncestor(tt.module)
		if tt.errMsg != "" {
			if assert.Error(t, err, "expected failure for %q but got: %v", tt.desc, got) {
				assert.Contains(t, err.Error(), tt.errMsg, tt.desc)
			}
		} else {
			if assert.NoError(t, err, tt.desc) {
				assert.Equal(t, tt.expected, got, tt.desc)
			}
		}
	}
}
