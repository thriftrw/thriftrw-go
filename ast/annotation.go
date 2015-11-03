package ast

// Annotation represents a type annotation. Type annotations are key-value
// pairs in the form,
//
// 	(foo = "bar", baz = "qux")
//
// They may be used to customize the generated code.
type Annotation struct {
	Name  string
	Value string
	Line  int
}
