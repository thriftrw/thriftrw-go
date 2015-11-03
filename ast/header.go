package ast

// Header unifies types representing header in the AST.
type Header interface {
	header()
}

func (*Include) header()   {}
func (*Namespace) header() {}

// Include is a request to include another Thrift file.
//
// 	include "shared.thrift"
type Include struct {
	Path string
	Line int
}

// Namespace statements allow users to choose the package name used by the
// generated code in certain languages.
//
// 	namespace py foo.bar
type Namespace struct {
	Scope string
	Name  string
	Line  int
}
