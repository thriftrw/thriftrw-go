package rule

import (
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/lint"
)

// TooManyArguments is a linter rule that fails if a service function accepts
// more than two arguments.
var TooManyArguments = tooManyArguments.MustBuild()

var tooManyArguments = &lint.SimpleRule{
	Name:     "Too many arguments",
	Severity: lint.Advice,
	Valid: func(w ast.Walker, f *ast.Function) bool {
		return len(f.Parameters) <= 2
	},
	MessageTemplate: `"{{.N.Name}}" of "{{.P.Name}}" has too many arguments. ` +
		"Consider creating a struct that holds all its arguments and use that as " +
		"the only argument to the function. Note that if this API is in production " +
		"this will be considered a breaking change.",
}
