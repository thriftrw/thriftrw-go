package internal

import "strings"

type multiError []error

func (errs multiError) Error() string {
	msg := "The following errors occurred:"
	for _, err := range errs {
		msg += "\n -  " + indentTail(4, err.Error())
	}
	return msg
}

// indentTail prepends the given number of spaces to all lines following the
// first line of the given string.
func indentTail(spaces int, s string) string {
	prefix := strings.Repeat(" ", spaces)
	lines := strings.Split(s, "\n")
	for i, line := range lines[1:] {
		lines[i+1] = prefix + line
	}
	return strings.Join(lines, "\n")
}

// MultiError combines a list of errors into one.
//
// Returns nil if the error list is empty.
func MultiError(errors []error) error {
	switch len(errors) {
	case 0:
		return nil
	case 1:
		return errors[0]
	}

	newErrors := make(multiError, 0, len(errors))
	for _, err := range errors {
		switch e := err.(type) {
		case multiError:
			newErrors = append(newErrors, e...)
		default:
			newErrors = append(newErrors, e)
		}
	}

	return multiError(newErrors)
}
