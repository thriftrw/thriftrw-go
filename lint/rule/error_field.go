package rule

import (
	"strings"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/lint"
)

// ErrorField defines a linter rule that fails if a Thrift exception is found
// with the name `error`.
var ErrorField = errorField.MustBuild()

var errorField = &lint.SimpleRule{
	Name:     "Error Field",
	Severity: lint.Warning,
	Valid: func(w ast.Walker, field *ast.Field) bool {
		s, ok := w.Parent().(*ast.Struct)
		if !ok || s.Type != ast.ExceptionType {
			return true
		}

		return strings.ToLower(field.Name) != "error"
	},
	MessageTemplate: `Field "{{.N.Name}}" of "{{.P.Name}}" may be exported as "Error" in Go. ` +
		"This will conflict with the `Error()` method of the `error` interface. " +
		"Please rename this field to something else.",
}
