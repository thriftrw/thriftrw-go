// Copyright (c) 2026 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPtrHelpers(t *testing.T) {
	t.Run("Bool", func(t *testing.T) {
		assert.Equal(t, true, *Bool(true))
		assert.Equal(t, false, *Bool(false))
	})

	t.Run("Int8", func(t *testing.T) {
		val := int8(42)
		assert.Equal(t, val, *Int8(val))
	})

	t.Run("Int16", func(t *testing.T) {
		val := int16(42)
		assert.Equal(t, val, *Int16(val))
	})

	t.Run("Int32", func(t *testing.T) {
		val := int32(42)
		assert.Equal(t, val, *Int32(val))
	})

	t.Run("Int64", func(t *testing.T) {
		val := int64(42)
		assert.Equal(t, val, *Int64(val))
	})

	t.Run("Float64", func(t *testing.T) {
		val := 3.14159
		assert.Equal(t, val, *Float64(val))
	})

	t.Run("String", func(t *testing.T) {
		val := "hello"
		assert.Equal(t, val, *String(val))
	})
}
