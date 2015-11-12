package ast

import (
	"fmt"
	"strings"
)

// Annotation represents a type annotation. Type annotations are key-value
// pairs in the form,
//
// 	(foo = "bar", baz = "qux")
//
// They may be used to customize the generated code. Annotations are optional
// anywhere in the code where they're accepted and may be skipped completely.
type Annotation struct {
	Name  string
	Value string
	Line  int
}

func (ann *Annotation) String() string {
	return fmt.Sprintf("%s = %q", ann.Name, ann.Value)
}

// FormatAnnotations formats a collection of annotations into a string.
func FormatAnnotations(anns []*Annotation) string {
	if anns == nil || len(anns) == 0 {
		return ""
	}

	as := make([]string, len(anns))
	for idx, ann := range anns {
		as[idx] = ann.String()
	}

	return "(" + strings.Join(as, ", ") + ")"
}
