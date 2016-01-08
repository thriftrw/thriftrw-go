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

import "github.com/uber/thriftrw-go/compile"

func (g *Generator) enum(spec *compile.EnumSpec) {
	// TODO(abg) define an error type in the library for unrecognized enums.
	err := g.DeclareFromTemplate(
		`
		{{ $fmt := import "fmt" }}
		{{ $wire := import "github.com/uber/thriftrw-go/wire" }}

		{{ $enumName := defName . }}
		type {{$enumName}} int32

		const (
		{{ range .Items }}
			{{$enumName}}{{.Name | goCase}} {{$enumName}} = {{.Value}}
		{{ end }}
		)

		func (v {{$enumName}}) ToWire() {{$wire}}.Value {
			return {{$wire}}.NewI32Value(int32(v))
		}

		func (v *{{$enumName}}) FromWire(w {{$wire}}.Value) error {
			switch w.GetI32() {
			{{ range .Items }}
			case {{.Value}}:
				*v = {{$enumName}}{{.Name | goCase}}
			{{ end }}
			default:
				return {{$fmt}}.Errorf("Unknown {{$enumName}}: %d", w.GetI32())
			}
			return nil
		}
		`,
		spec,
	)

	if err != nil {
		panic(err)
	}
}
