package internal

import (
	"bytes"
	"fmt"
)

// parseError is an error type to keep track of any parse errors and the
// positions they occur at.
type parseError struct {
	Errors map[int]string
}

func newParseError() parseError {
	return parseError{Errors: make(map[int]string)}
}

func (pe parseError) add(line int, msg string) {
	pe.Errors[line] = msg
}

func (pe parseError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("parse error\n")
	for line, msg := range pe.Errors {
		buffer.WriteString(fmt.Sprintf("  line %d: %s\n", line, msg))
	}
	return buffer.String()
}
