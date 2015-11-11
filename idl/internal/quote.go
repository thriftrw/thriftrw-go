package internal

import "strconv"

// UnquoteSingleQuoted unquotes a slice of bytes representing a single quoted
// string.
//
// 	UnquoteSingleQuoted([]byte("'foo'")) == "foo"
func UnquoteSingleQuoted(in []byte) (string, error) {
	out := make([]byte, len(in))

	// s/'/"/g, s/"/'/g
	for idx, c := range in {
		if c == '"' {
			c = '\''
		} else if c == '\'' {
			c = '"'
		}
		out[idx] = c
	}

	str, err := strconv.Unquote(string(out))
	if err != nil {
		return str, err
	}

	// s/'/"/g, s/"/'/g
	out = make([]byte, len(str))
	for idx, c := range []byte(str) {
		if c == '"' {
			c = '\''
		} else if c == '\'' {
			c = '"'
		}
		out[idx] = c
	}

	return string(out), nil
}
