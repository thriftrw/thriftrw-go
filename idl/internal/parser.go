package internal

import (
	"errors"

	"github.com/uber/thriftrw-go/ast"
)

func init() {
	// TODO configure parser here
	yyErrorVerbose = true
}

func Parse(s []byte) (*ast.Program, error) {
	lex := newLexer(s)
	e := yyParse(lex)
	if e == 0 && !lex.parseFailed {
		return lex.program, nil
	}
	return nil, errors.New("failed to parse")
}

//go:generate ragel -Z -G2 -o lex.go lex.rl
//go:generate goimports -w ./lex.go

//go:generate go tool yacc thrift.y
//go:generate goimports -w ./y.go
