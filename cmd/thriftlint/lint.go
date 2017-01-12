package main

import (
	"bytes"
	"fmt"
	"reflect"
	"text/template"

	"go.uber.org/thriftrw/ast"
)

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

// Catcher receives lint messages from different linter rules.
type Catcher interface {
	PostLint(Lint)
}

// Rule is a Linter rule.
type Rule struct {
	// Name of the rule
	Name string

	// Inspector walks the Thrift AST and posts linter messages to the given
	// catcher.
	Inspector func(Catcher) ast.Visitor
}

// SimpleRule provides a convenient API to write linter rules that don't need
// full control over the AST traversal and that post at most one message per
// node.
//
// SimpleRule linters call Inspect on matching node types and for any call
// that fails, the MessageTemplate is rendered using that Node.
type SimpleRule struct {
	Name string

	// Severity of messages posted by this rule.
	Severity Severity

	// A function in the form,
	//
	// 	func(w ast.Walker, n N) bool
	//
	// Where N is ast.Node or any type that implements ast.Node.
	Inspect interface{}

	// If Inspect fails for a Node, this template will be rendered (using
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

// Builds the SimpleRule into a Rule.
func (r *SimpleRule) build() (*Rule, error) {
	tmpl, err := template.New(r.Name).Funcs(r.TemplateFuncs).Parse(r.MessageTemplate)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to parse message template for rule %q: %v", r.Name, err)
	}

	inspect := reflect.ValueOf(r.Inspect)
	t := inspect.Type()
	switch {
	case t.Kind() != reflect.Func:
		return nil, fmt.Errorf("Inspect must be a function, found %v", t)
	case t.NumIn() != 2:
		return nil, fmt.Errorf("Inspect must accept two arguments, found %v", t.NumIn())
	case t.In(0) != _typeOfWalker:
		return nil, fmt.Errorf(
			"Inspect must accept an ast.Walker as its first argument, found %v", t.In(0))
	case t.In(1) != _typeOfNode && !t.In(1).Implements(_typeOfNode):
		return nil, fmt.Errorf(
			"Inspect must accept an ast.Node or a type that implements ast.Node "+
				"as its second argument, found %v", t.In(1))
	case t.NumOut() != 1:
		return nil, fmt.Errorf("Inspect must return one result, found %v", t.NumOut())
	case t.Out(0) != _typeOfBool:
		return nil, fmt.Errorf("Inspect must return a bool, found %v", t.Out(0))
	}

	return &Rule{
		Name: r.Name,
		Inspector: func(c Catcher) ast.Visitor {
			return &simpleRuleVisitor{
				Rule:            r,
				Catcher:         c,
				Inspect:         inspect,
				NodeType:        t.In(1),
				MessageTemplate: tmpl,
			}
		},
	}, nil
}

func (r *SimpleRule) mustBuild() *Rule {
	rule, err := r.build()
	if err != nil {
		panic(fmt.Errorf("failed to build simple rule %q: %v", r.Name, err))
	}
	return rule
}

type simpleRuleVisitor struct {
	Rule            *SimpleRule
	Catcher         Catcher
	Inspect         reflect.Value
	NodeType        reflect.Type
	MessageTemplate *template.Template
}

func (v *simpleRuleVisitor) Visit(w ast.Walker, n ast.Node) ast.Visitor {
	nV := reflect.ValueOf(n)
	if !nV.Type().AssignableTo(v.NodeType) {
		return v
	}

	result := v.Inspect.Call([]reflect.Value{reflect.ValueOf(w), nV})
	if result[0].Bool() {
		return v
	}

	var buff bytes.Buffer
	err := v.MessageTemplate.Execute(&buff, struct {
		N ast.Node
		P ast.Node
	}{N: n, P: w.Parent()})
	if err != nil {
		panic(fmt.Sprintf(
			"failed to execute message template for rule %q: %v", v.Rule.Name, err))
	}

	lint := Lint{
		Severity: v.Rule.Severity,
		Line:     getLineNumber(nV),
		Message:  buff.String(),
	}
	v.Catcher.PostLint(lint)
	return v
}

func getLineNumber(node reflect.Value) int {
	// TODO: Don't use reflection for this. We can introduce a HasLine
	// subinterface for Nodes which know their line number.

	for node.Kind() != reflect.Struct {
		if node.Kind() == reflect.Ptr {
			node = node.Elem()
		} else {
			return 0
		}
	}

	if line := node.FieldByName("Line"); line.IsValid() {
		return int(line.Int())
	}

	return 0
}
