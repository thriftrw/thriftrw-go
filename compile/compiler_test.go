// Copyright (c) 2016 Uber Technologies, Inc.
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

package compile

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/wire"
)

func TestFilesystem(t *testing.T) {
	files := map[string]string{
		"/some/prefix/main.thrift": `
			include "./shared/shared.thrift"

			struct S {
				1: optional shared.UUID uuid;
			}
		`,
		"/some/prefix/shared/shared.thrift": `
			typedef string UUID;
		`,
	}

	fs := dummyFS{"/some/prefix/", files}

	module, err := Compile("main.thrift", Filesystem(fs))
	require.NoError(t, err, "Compile failed")

	sType, err := module.LookupType("S")
	require.NoError(t, err, "Lookup S failed")
	require.NotNil(t, sType, "Type S is nil")
	assert.Equal(t, wire.TStruct, sType.TypeCode(), "Type mismatch")
}

func TestCompile(t *testing.T) {
	module, err := Compile("../gen/testdata/thrift/services.thrift")
	require.NoError(t, err, "Compile failed")

	kvSvc, err := module.LookupService("KeyValue")
	require.NoError(t, err, "Lookup KeyValue failed")
	require.NotNil(t, kvSvc, "KeyValue service is nil")
}

func TestCompileNonStrict(t *testing.T) {
	files := map[string]string{
		"/some/prefix/main.thrift": `
			struct S {
				1: string uuid;
			}
		`,
	}

	fs := dummyFS{"/some/prefix/", files}

	module, err := Compile("main.thrift", Filesystem(fs), NonStrict())
	require.NoError(t, err, "Compile failed")

	sType, err := module.LookupType("S")
	require.NoError(t, err, "Lookup S failed")
	require.NotNil(t, sType, "Type S is nil")
	assert.Equal(t, wire.TStruct, sType.TypeCode(), "Type mismatch")

	uuidField, err := sType.(*StructSpec).Fields.FindByName("uuid")
	require.NoError(t, err, "Failed to find UUID field in struct")
	assert.False(t, uuidField.Required, "Unspecified requiredness should be treated as optional")
}
