package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatAnnotations(t *testing.T) {
	var anns []*Annotation
	assert.Equal(t, "", FormatAnnotations(anns))

	anns = append(anns, &Annotation{Name: "foo", Value: "bar"})
	assert.Equal(t, `(foo = "bar")`, FormatAnnotations(anns))

	anns = append(anns, &Annotation{Name: "baz", Value: "qux"})
	assert.Equal(t, `(foo = "bar", baz = "qux")`, FormatAnnotations(anns))
}
