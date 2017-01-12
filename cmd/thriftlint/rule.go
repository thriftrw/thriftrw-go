package main

import (
	"fmt"
	"strings"
	"text/template"

	"go.uber.org/thriftrw/ast"
)

var javaNamespace = &SimpleRule{
	Name:     "Java namespacing",
	Severity: Error,
	Inspect: func(w ast.Walker, prog *ast.Program) bool {
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

var errorFieldRule = &SimpleRule{
	Name:     "Error Field",
	Severity: Warning,
	Inspect: func(w ast.Walker, field *ast.Field) bool {
		return strings.ToLower(field.Name) != "error"
	},
	MessageTemplate: `Field "{{.N.Name}}" of "{{.P.Name}}" may be exported as "Error" in Go. ` +
		"This will conflict with the `Error()` method of the `error` interface. " +
		"Please rename this field to something else.",
}

var unannotedI64Rule = &SimpleRule{
	Name:     "Unannoted i64",
	Severity: Warning,
	Inspect: func(w ast.Walker, t ast.BaseType) bool {
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

var unspecifiedRequiredness = &SimpleRule{
	Name:     "Unspecified requiredness",
	Severity: Error,
	Inspect: func(w ast.Walker, field *ast.Field) bool {
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

var tooManyArguments = &SimpleRule{
	Name:     "Too many arguments",
	Severity: Advice,
	Inspect: func(w ast.Walker, f *ast.Function) bool {
		return len(f.Parameters) <= 2
	},
	MessageTemplate: `"{{.N.Name}}" of "{{.P.Name}}" has too many arguments. ` +
		"Consider creating a struct that holds all its arguments and use that as " +
		"the only argument to the function. Note that if this API is in production " +
		"this will be considered a breaking change.",
}
