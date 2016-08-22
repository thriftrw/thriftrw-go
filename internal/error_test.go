package internal

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiError(t *testing.T) {
	tests := []struct {
		give        []error
		want        error
		wantMessage string
	}{
		{
			give: []error{},
			want: nil,
		},
		{
			give:        []error{errors.New("great sadness")},
			want:        errors.New("great sadness"),
			wantMessage: "great sadness",
		},
		{
			give: []error{
				errors.New("foo"),
				errors.New("bar"),
			},
			want: multiError{
				errors.New("foo"),
				errors.New("bar"),
			},
			wantMessage: "The following errors occurred:\n" +
				" -  foo\n" +
				" -  bar",
		},
		{
			give: []error{
				errors.New("great sadness"),
				errors.New("multi\n  line\nerror message"),
				errors.New("single line error message"),
			},
			want: multiError{
				errors.New("great sadness"),
				errors.New("multi\n  line\nerror message"),
				errors.New("single line error message"),
			},
			wantMessage: "The following errors occurred:\n" +
				" -  great sadness\n" +
				" -  multi\n" +
				"      line\n" +
				"    error message\n" +
				" -  single line error message",
		},
		{
			give: []error{
				errors.New("foo"),
				multiError{
					errors.New("bar"),
					errors.New("baz"),
				},
				errors.New("qux"),
			},
			want: multiError{
				errors.New("foo"),
				errors.New("bar"),
				errors.New("baz"),
				errors.New("qux"),
			},
			wantMessage: "The following errors occurred:\n" +
				" -  foo\n" +
				" -  bar\n" +
				" -  baz\n" +
				" -  qux",
		},
	}

	for _, tt := range tests {
		err := MultiError(tt.give)
		if assert.Equal(t, tt.want, err) && tt.wantMessage != "" {
			assert.Equal(t, tt.wantMessage, err.Error())
		}
	}
}
