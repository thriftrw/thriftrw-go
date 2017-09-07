package internal

import (
	"strings"
	"unicode"
)

const (
	_docstringOpen  = "/**"
	_docstringClose = "*/"
)

// CleanupDocstring takes a docstring in the form,
//
//  /**
//   * foo bar
//   */
//
// And returns,
//
//  foo bar
func CleanupDocstring(s string) string {
	lines := strings.Split(s, "\n")
	dedent(lines, true /* skipFirstIfUnindented */)

	switch len(lines) {
	case 0:
		return s
	case 1:
		// /** foo */
		s := lines[0]
		s = strings.TrimPrefix(s, _docstringOpen)
		s = strings.TrimSuffix(s, _docstringClose)
		return strings.TrimSpace(s)
	}

	if len(lines) > 0 {
		lines[0] = strings.TrimPrefix(lines[0], _docstringOpen)
		lines[len(lines)-1] = strings.TrimSuffix(lines[len(lines)-1], _docstringClose)
	}

	lines = dropLeadingEmptyLines(lines)
	lines = dropTrailingEmptyLines(lines)

	// At this point, we need to strip the leading "*" and " *" from every
	// line and dedent again.
	for i, l := range lines {
		if len(l) == 0 {
			continue
		}

		if l[0] == '*' {
			lines[i] = l[1:]
			continue
		}

		if len(l) > 1 && l[:2] == " *" {
			lines[i] = l[2:]
			continue
		}
	}

	dedent(lines, false /* skipFirstIfUnindented */)
	lines = dropLeadingEmptyLines(lines)
	lines = dropTrailingEmptyLines(lines)

	return strings.Join(lines, "\n")
}

func dropLeadingEmptyLines(lines []string) []string {
	for len(lines) > 0 {
		if len(lines[0]) > 0 {
			break
		}
		lines = lines[1:]
	}
	return lines
}

func dropTrailingEmptyLines(lines []string) []string {
	for i := len(lines) - 1; i >= 0; i-- {
		if len(lines[i]) > 0 {
			break
		}
		lines = lines[:i]
	}
	return lines
}

// skipFirstIfUnindented tells dedent to not consider the first line for the
// prefix if it's not already indented. This is needed because our docstringcs
// will start at "/**" without any leading whitespace.
func dedent(lines []string, skipFirstIfUnindented bool) {
	if len(lines) == 0 {
		return
	}

	var (
		havePrefix bool
		prefix     string
	)

	// Don't consider the first line for the prefix if it hasn't been
	// indented.
	if nonSpace := strings.IndexFunc(lines[0], notIsSpace); nonSpace >= 0 {
		if nonSpace == 0 && skipFirstIfUnindented {
			// The first line starts with a non-space character. Skip this line.
			lines = lines[1:]
		} else {
			havePrefix = true
			prefix = lines[0][:nonSpace]
		}
	}

	for i, s := range lines {
		nonSpace := strings.IndexFunc(s, notIsSpace)
		if nonSpace < 0 {
			// Whitespace-only. Use an empty string.
			lines[i] = ""
			continue
		}

		if !havePrefix {
			prefix = s[:nonSpace]
			havePrefix = true
		}

		// Dedent only if the first non-space character appears at or after
		// the prefix.
		if nonSpace >= len(prefix) {
			lines[i] = s[len(prefix):]
		}
	}
}

func notIsSpace(r rune) bool {
	return !unicode.IsSpace(r)
}
