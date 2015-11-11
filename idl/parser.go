package idl

import "github.com/uber/thriftrw-go/ast"
import "github.com/uber/thriftrw-go/idl/internal"

// Parse parses a Thrift document.
func Parse(s []byte) (*ast.Program, error) {
	return internal.Parse(s)
}
