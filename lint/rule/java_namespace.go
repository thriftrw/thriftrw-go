package rule

import (
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/lint"
)

// JavaNamespace defines a linter rule that complains if a Thrift file does not
// define a Java namespace.
var JavaNamespace = javaNamespace.MustBuild()

var javaNamespace = &lint.SimpleRule{
	Name:     "Java namespacing",
	Severity: lint.Error,
	Valid: func(w ast.Walker, prog *ast.Program) bool {
		for _, header := range prog.Headers {
			ns, ok := header.(*ast.Namespace)
			if !ok {
				continue
			}

			if ns.Scope == "java" || ns.Scope == "*" {
				return true
			}
		}
		return false
	},
	MessageTemplate: "No namespace specified for Java. Please specify a namespace to ensure " +
		"that your Thrift file is compatible with all Java clients.",
}
