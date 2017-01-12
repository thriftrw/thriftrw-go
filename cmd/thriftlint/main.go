package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"go.uber.org/thriftrw/ast"
	"go.uber.org/thriftrw/idl"
)

type catcher struct {
	lints []Lint
}

func (c *catcher) PostLint(l Lint) {
	c.lints = append(c.lints, l)
}

type byLine []Lint

func (ll byLine) Len() int {
	return len(ll)
}

func (ll byLine) Less(i, j int) bool {
	return ll[i].Line < ll[j].Line
}

func (ll byLine) Swap(i, j int) {
	ll[i], ll[j] = ll[j], ll[i]
}

var _simpleRules = []*SimpleRule{
	javaNamespace,
	errorFieldRule,
	unannotedI64Rule,
	unspecifiedRequiredness,
	tooManyArguments,
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

	var rules []*Rule
	for _, r := range _simpleRules {
		rules = append(rules, r.mustBuild())
	}

	var (
		c        catcher
		visitors []ast.Visitor
	)
	for _, rule := range rules {
		visitors = append(visitors, rule.Inspector(&c))
	}

	ast.Walk(ast.MultiVisitor(visitors...), prog)

	sort.Sort(byLine(c.lints))
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
