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
