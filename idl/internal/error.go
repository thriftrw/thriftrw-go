// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package internal

import (
	"bytes"
	"fmt"
)

// parseError is an error type to keep track of any parse errors and the
// positions they occur at.
type parseError struct {
	Errors map[int][]string
}

func newParseError() parseError {
	return parseError{Errors: make(map[int][]string)}
}

func (pe parseError) add(line int, msg string) {
	pe.Errors[line] = append(pe.Errors[line], msg)
}

func (pe parseError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString("parse error\n")
	for line, msgs := range pe.Errors {
		for _, msg := range msgs {
			buffer.WriteString(fmt.Sprintf("  line %d: %s\n", line, msg))
		}
	}
	return buffer.String()
}
