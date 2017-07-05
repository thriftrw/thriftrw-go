package lint

import "go.uber.org/thriftrw/ast"

// Severity of an issue
type Severity int

// Varying levels of issue severity.
const (
	Advice Severity = iota
	Warning
	Error
)

// Lint is an issue found in a Thrift file.
type Lint struct {
	Severity Severity

	// If non-zero, this is the line in the Thrift file on which this problem
	// was found.
	Line int

	// Message is a plain-text message to the user indicating what the problem
	// was.
	Message string
}

// Trap receives lint messages from different linter rules.
type Trap interface {
	PostLint(Lint)
}

// Rule is a Linter rule.
type Rule interface {
	// Name of the rule
	Name() string

	// Inspector walks the Thrift AST and posts linter messages to the given
	// Trap.
	Inspector(Trap) ast.Visitor
}
