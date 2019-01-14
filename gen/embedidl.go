// Copyright (c) 2019 Uber Technologies, Inc.
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
	"sort"

	"go.uber.org/thriftrw/compile"
)

// embedIDL generate Go code with a full copy of the IDL embeded.
func embedIDL(g Generator, i thriftPackageImporter, m *compile.Module) error {
	pkg, err := i.Package(m.ThriftPath)
	if err != nil {
		return wrapGenerateError("idl embedding", err)
	}
	packageRelPath, err := i.RelativeThriftFilePath(m.ThriftPath)
	if err != nil {
		return wrapGenerateError("idl embedding", err)
	}

	hash := sha1.Sum(m.Raw)
	var includes []string
	for _, v := range m.Includes {
		importPath, err := i.Package(v.Module.ThriftPath)
		if err != nil {
			return wrapGenerateError("idl embedding", err)
		}
		includes = append(includes, g.Import(importPath))
	}

	sort.Strings(includes)

	data := struct {
		Name     string
		Package  string
		FilePath string
		SHA1     string
		Includes []string
		Raw      []byte
	}{
		Name:     m.Name,
		Package:  pkg,
		FilePath: packageRelPath,
		SHA1:     hex.EncodeToString(hash[:]),
		Includes: includes,
		Raw:      m.Raw,
	}
	err = g.DeclareFromTemplate(`
		<$idl := import "go.uber.org/thriftrw/thriftreflect">

		// ThriftModule represents the IDL file used to generate this package.
		var ThriftModule = &<$idl>.ThriftModule {
			Name: "<.Name>",
			Package: "<.Package>",
			FilePath: <printf "%q" .FilePath>,
			SHA1: "<.SHA1>",
			<if .Includes ->
				Includes: []*<$idl>.ThriftModule {<range .Includes>
						<.>.ThriftModule, <end>
					},
			<end ->
			Raw: rawIDL,
		}
		const rawIDL = <printf "%q" .Raw>
		`, data)
	return wrapGenerateError("idl embedding", err)
}
