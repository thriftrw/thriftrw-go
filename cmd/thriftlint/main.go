package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
	"go.uber.org/thriftrw/lint"
	"go.uber.org/thriftrw/lint/rule"
)

type catcher struct {
	lints []lint.Lint
}

func (c *catcher) PostLint(l lint.Lint) {
	c.lints = append(c.lints, l)
}

var _rules = []lint.Rule{
	rule.JavaNamespace,
	rule.ErrorField,
	rule.UnannotatedI64,
	rule.UnspecifiedRequiredness,
	rule.TooManyArguments,
}

func main() {
	src, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	prog, err := idl.Parse(src)
	if err != nil {
		log.Fatal(err)
	}

	var (
		c        catcher
		visitors []ast.Visitor
	)
	for _, rule := range _rules {
		visitors = append(visitors, rule.Inspector(&c))
	}

	ast.Walk(ast.MultiVisitor(visitors...), prog)

	lint.SortByLine(c.lints)
	for _, lint := range c.lints {
		prefix := fmt.Sprintf("%d: ", lint.Line)
		indent := len(prefix)
		fmt.Printf("%s%s\n", prefix, indentTail(indent, wrap(lint.Message, 72)))
	}
}

func wrap(s string, width int) string {
	words := strings.Split(s, " ")
	var lines []string
	for len(words) > 0 {
		var (
			line   []string
			length int
		)

		for len(words) > 0 && len(words[0])+length < width {
			line = append(line, words[0])
			length += len(words[0])
			words = words[1:]
		}

		lines = append(lines, strings.Join(line, " "))
	}
	return strings.Join(lines, "\n")
}

func indentTail(spaces int, s string) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(s, "\n")
	for i, line := range lines[1:] {
		lines[i+1] = prefix + line
	}
	return strings.Join(lines, "\n")
}
