package envelopetest

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/thriftrw/wire"
)

// AssertEqual asserts that the provided envelopes are equivalent.
func AssertEqual(t assert.TestingT, want, got wire.Envelope, msgAndArgs ...interface{}) (ok bool) {
	ok = true

	// Perform all checks and return false if any of them fails.
	ok = assert.Equal(t, want.Name, got.Name, msgAndArgs...) && ok
	ok = assert.Equal(t, want.Type, got.Type, msgAndArgs...) && ok
	ok = assert.Equal(t, want.SeqID, got.SeqID, msgAndArgs...) && ok
	ok = assert.True(t, wire.ValuesAreEqual(want.Value, got.Value), msgAndArgs...) && ok

	return ok
}
