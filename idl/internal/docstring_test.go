package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanupDocstring(t *testing.T) {
	tests := []struct {
		give []string
		want []string
	}{
		{},
		{
			give: []string{"/** foo bar */"},
			want: []string{"foo bar"},
		},
		{
			give: []string{
				"/**",
				" * foo",
				" *   bar",
				" * baz",
				" */",
			},
			want: []string{
				"foo",
				"  bar",
				"baz",
			},
		},
		{
			give: []string{
				"	/**",
				"	 * foo",
				"	 * 	bar",
				"	 * baz",
				"	 */",
			},
			want: []string{
				"foo",
				"	bar",
				"baz",
			},
		},
		{
			give: []string{
				"/**",
				"	 * hello",
				"	 * world",
				"	 */",
			},
			want: []string{
				"hello",
				"world",
			},
		},
		{
			give: []string{
				"/**",
				"	 * hello",
				"	 *",
				"	 * world",
				"	 */",
			},
			want: []string{
				"hello",
				"",
				"world",
			},
		},
		{
			give: []string{
				"/**",
				"	 * hello",
				"",
				"	 * world",
				"	 */",
			},
			want: []string{
				"hello",
				"",
				"world",
			},
		},
		{
			give: []string{
				"	/**",
				"	 *foo",
				"	 *	bar",
				"	 *baz",
				"	 */",
			},
			want: []string{
				"foo",
				"	bar",
				"baz",
			},
		},
		{
			give: []string{
				"/**",
				"	    * foo does stuff",
				"	    */",
			},
			want: []string{"foo does stuff"},
		},
	}

	for _, tt := range tests {
		give := strings.Join(tt.give, "\n")
		want := strings.Join(tt.want, "\n")
		got := CleanupDocstring(give)
		assert.Equalf(t, want, got, "failed to cleanup %#v", tt.give)
	}
}
