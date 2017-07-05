package rule

import (
	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/lint"
)

// UnspecifiedRequiredness is a linter rule that fails if any field of a struct
// or exception does not explicitly state whether it is required or optional.
var UnspecifiedRequiredness = unspecifiedRequiredness.MustBuild()

var unspecifiedRequiredness = &lint.SimpleRule{
	Name:     "Unspecified requiredness",
	Severity: lint.Error,
	Valid: func(w ast.Walker, field *ast.Field) bool {
		// Field requiredness may be absent if the field is a function
		// parameter or exception, or a field inside a union.
		switch p := w.Parent().(type) {
		case *ast.Function:
			return true
		case *ast.Struct:
			if p.Type == ast.UnionType {
				return true
			}
		}

		return field.Requiredness != ast.Unspecified
	},
	MessageTemplate: `"{{.N.Name}}" of "{{.P.Name}}" is not marked as required or optional. ` +
		"Thrift has inconsistent default behavior for structs and " +
		"exception fields which are not explicitly marked either way. " +
		`Please specify whether "{{.N.Name}}" is required or optional.`,
}
