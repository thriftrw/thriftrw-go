// Copyright (c) 2016 Uber Technologies, Inc.
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

import "strings"

type multiError []error

func (errs multiError) Error() string {
	msg := "The following errors occurred:"
	for _, err := range errs {
		msg += "\n -  " + indentTail(4, err.Error())
	}
	return msg
}

// indentTail prepends the given number of spaces to all lines following the
// first line of the given string.
func indentTail(spaces int, s string) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(s, "\n")
	for i, line := range lines[1:] {
		lines[i+1] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// MultiError combines a list of errors into one.
//
// Returns nil if the error list is empty.
func MultiError(errors []error) error {
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	}

	newErrors := make(multiError, 0, len(errors))
	for _, err := range errors {
		switch e := err.(type) {
		case multiError:
			newErrors = append(newErrors, e...)
		default:
			newErrors = append(newErrors, e)
		}
	}

	return multiError(newErrors)
}
