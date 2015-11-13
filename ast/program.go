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

package ast

import "fmt"

// Program represents the full syntax tree for a single .thrift file.
type Program struct {

	// Headers

	Includes   []*Include
	Namespaces []*Namespace

	// Definitions

	Constants []*Constant
	Typedefs  []*Typedef
	Enums     []*Enum
	Structs   []*Struct
	Services  []*Service
}

// TODO is AddHeader and AddDefinition even needed

// AddHeader adds a new Header to the AST for the Thrift file.
func (p *Program) AddHeader(header Header) {
	switch hdr := header.(type) {
	case *Include:
		p.Includes = append(p.Includes, hdr)
	case *Namespace:
		p.Namespaces = append(p.Namespaces, hdr)
	default:
		panic(fmt.Sprintf("parser: unknown header type %T: %v", hdr, hdr))
	}
}

// AddDefinition adds a new Definition to the AST for the Thrift file.
func (p *Program) AddDefinition(definition Definition) error {
	switch def := definition.(type) {
	case *Constant:
		p.Constants = append(p.Constants, def)
	case *Typedef:
		p.Typedefs = append(p.Typedefs, def)
	case *Enum:
		p.Enums = append(p.Enums, def)
	case *Struct:
		p.Structs = append(p.Structs, def)
	case *Service:
		p.Services = append(p.Services, def)
	default:
		return fmt.Errorf("parser: unknown definition type %T: %v", def, def)
	}
	return nil
}
