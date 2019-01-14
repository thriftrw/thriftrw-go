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

import "strconv"

// UnquoteSingleQuoted unquotes a slice of bytes representing a single quoted
// string.
//
// 	UnquoteSingleQuoted([]byte("'foo'")) == "foo"
func UnquoteSingleQuoted(in []byte) (string, error) {
	out := string(swapQuotes(in))
	str, err := strconv.Unquote(out)
	if err != nil {
		return str, err
	}

	// s/'/"/g, s/"/'/g
	out = string(swapQuotes([]byte(str)))
	return out, nil
}

// swapQuotes replaces all single quotes with double quotes and all double
// quotes with single quotes.
func swapQuotes(in []byte) []byte {
	// s/'/"/g, s/"/'/g
	out := make([]byte, len(in))
	for i, c := range in {
		if c == '"' {
			c = '\''
		} else if c == '\'' {
			c = '"'
		}
		out[i] = c
	}
	return out
}
