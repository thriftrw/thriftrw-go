package rule

import (
	"fmt"
	"text/template"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/lint"
)

// UnannotatedI64 is a liner rule that complains if it finds usage of the i64
// type without the `js.type` annotation.
var UnannotatedI64 = unannotatedI64.MustBuild()

var unannotatedI64 = &lint.SimpleRule{
	Name:     "Unannoted i64",
	Severity: lint.Warning,
	Valid: func(w ast.Walker, t ast.BaseType) bool {
		if t.ID != ast.I64TypeID {
			return true
		}

		for _, ann := range t.Annotations {
			if ann.Name == "js.type" {
				return true
			}
		}

		return false
	},
	MessageTemplate: `{{parentRelationship .P}} "i64" but it does not have a "js.type" annotation. ` +
		"JavaScript does not support 64-bit integers. Instead it uses Buffer, " +
		"Date, or Long (https://npmjs.com/package/long). We recommend that you " +
		`provide a "js.type" annotation for all uses of "i64" to let ` +
		"thriftrw-node know which of these options you want to use. Valid " +
		`values are "Date", "Long", "Buffer". See ` +
		"https://github.com/thriftrw/thriftrw-node#i64 for more information.",
	TemplateFuncs: template.FuncMap{
		"parentRelationship": func(n ast.Node) string {
			switch p := n.(type) {
			case *ast.Typedef:
				return fmt.Sprintf("%q is an alias for", p.Name)
			case *ast.Function:
				return fmt.Sprintf("%q returns an", p.Name)
			case *ast.Field:
				return fmt.Sprintf("%q has type", p.Name)
			case *ast.Constant:
				return fmt.Sprintf("%q is a constant with type", p.Name)
			default:
				panic(fmt.Sprintf("unknown parent for i64 type reference: %v", n))
			}
		},
	},
}
