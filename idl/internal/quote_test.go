package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnquoteSingleQuoted(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{`'foo'`, "foo"},
		{`'a "b" c'`, `a "b" c`},
		{`'a \'b\' c'`, `a 'b' c`},
		{`'a \\"b\\" c'`, `a \"b\" c`},
	}

	for _, tt := range tests {
		got, err := UnquoteSingleQuoted([]byte(tt.in))
		if assert.NoError(t, err, "Failed to unquote: %#v", tt.in) {
			assert.Equal(t, tt.out, got, "Unquote incorrect: %#v", tt.in)
		}
	}
}
