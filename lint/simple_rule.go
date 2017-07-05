package lint

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"go.uber.org/thriftrw/ast"
)

// SimpleRule provides a convenient API to write linter rules that don't need
// full control over the AST traversal and that post at most one message per
// node.
//
// SimpleRule linters call Valid on matching node types and for any call that
// fails, the MessageTemplate is rendered using that Node.
type SimpleRule struct {
	Name string

	// Severity of messages posted by this rule.
	Severity Severity

	// A function in the form,
	//
	// 	func(w ast.Walker, n N) (isValid bool)
	//
	// Where N is ast.Node or any type that implements ast.Node.
	//
	// This function will be called on nodes of the Thrift AST and must return
	// whether that node is valid. For nodes that are marked invalid, a message
	// of the associated Severity will be posted to the linter.
	Valid interface{}

	// If Valid fails for a Node, this template will be rendered (using
	// text/template) with the following context:
	//
	// 	.N: The Node for which Inspect failed
	// 	.P: Parent node of .Node (if any)
	//
	// The system will panic if MessageTemplate fails to execute for a Node.
	MessageTemplate string

	// Functions to make available in MessageTemplate
	TemplateFuncs template.FuncMap
}

var (
	_typeOfWalker = reflect.TypeOf((*ast.Walker)(nil)).Elem()
	_typeOfNode   = reflect.TypeOf((*ast.Node)(nil)).Elem()
	_typeOfBool   = reflect.TypeOf(true)
)

// Build converts a SimpleRule into a Rule.
func (s *SimpleRule) Build() (Rule, error) {
	tmpl, err := template.New(s.Name).Funcs(s.TemplateFuncs).Parse(s.MessageTemplate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse message template for rule %q: %v", s.Name, err)
	}

	valid := reflect.ValueOf(s.Valid)
	t := valid.Type()
	switch {
	case t.Kind() != reflect.Func:
		return nil, fmt.Errorf("Valid must be a function, found %v", t)
	case t.NumIn() != 2:
		return nil, fmt.Errorf("Valid must accept two arguments, found %v", t.NumIn())
	case t.In(0) != _typeOfWalker:
		return nil, fmt.Errorf(
			"Valid must accept an ast.Walker as its first argument, found %v", t.In(0))
	case t.In(1) != _typeOfNode && !t.In(1).Implements(_typeOfNode):
		return nil, fmt.Errorf(
			"Valid must accept an ast.Node or a type that implements ast.Node "+
				"as its second argument, found %v", t.In(1))
	case t.NumOut() != 1:
		return nil, fmt.Errorf("Valid must return one result, found %v", t.NumOut())
	case t.Out(0) != _typeOfBool:
		return nil, fmt.Errorf("Valid must return a bool, found %v", t.Out(0))
	}

	return simpleRule{
		SimpleRule: s,
		valid:      valid,
		nodeType:   t.In(1),
		template:   tmpl,
	}, nil
}

// MustBuild converts this SimpleRule into a Rule or panics if the operation
// fails.
func (s *SimpleRule) MustBuild() Rule {
	r, err := s.Build()
	if err != nil {
		panic(err)
	}
	return r
}

type simpleRule struct {
	*SimpleRule

	valid    reflect.Value
	nodeType reflect.Type
	template *template.Template

	// Set only if we're being used as a visitor.
	trap Trap
}

var (
	_ Rule        = simpleRule{}
	_ ast.Visitor = simpleRule{}
)

func (s simpleRule) Name() string {
	return s.SimpleRule.Name
}

func (s simpleRule) Inspector(trap Trap) ast.Visitor {
	// Pass-by-value. We can just use this object as the visitor.
	s.trap = trap
	return s
}

func (s simpleRule) Visit(w ast.Walker, n ast.Node) ast.Visitor {
	nV := reflect.ValueOf(n)
	if !nV.Type().AssignableTo(s.nodeType) {
		return s
	}

	result := s.valid.Call([]reflect.Value{reflect.ValueOf(w), nV})
	if result[0].Bool() {
		return s
	}

	var buff bytes.Buffer
	err := s.template.Execute(&buff, struct {
		N ast.Node
		P ast.Node
	}{N: n, P: w.Parent()})
	if err != nil {
		panic(fmt.Sprintf(
			"failed to execute message template for rule %q: %v", s.Name(), err))
	}

	lint := Lint{
		Severity: s.SimpleRule.Severity,
		Line:     ast.LineNumber(n),
		Message:  buff.String(),
	}
	s.trap.PostLint(lint)
	return s
}
