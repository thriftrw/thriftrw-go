package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnquoteSingleQuoted(t *testing.T) {
	cases := []struct {
		in  string
		out string
	}{
		{`'foo'`, "foo"},
		{`'a "b" c'`, `a "b" c`},
		{`'a \'b\' c'`, `a 'b' c`},
		{`'a \\"b\\" c'`, `a \"b\" c`},
	}

	for _, c := range cases {
		got, err := UnquoteSingleQuoted([]byte(c.in))
		if assert.NoError(t, err, "Failed to unquote: %#v", c.in) {
			assert.Equal(t, c.out, got, "Unquote incorrect: %#v", c.in)
		}
	}
}
