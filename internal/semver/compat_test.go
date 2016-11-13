// Copyright (c) 2016 Uber Technologies, Inc.
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

package semver

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompatibleRange(t *testing.T) {
	test := []struct {
		library      string
		compatible   []string
		incompatible []string
	}{
		{
			library:      "1.2.3",
			compatible:   []string{"1.2.3", "1.2.4", "1.0.0"},
			incompatible: []string{"1.3.0", "0.9.9"},
		},
		{
			library:      "2.0.0",
			compatible:   []string{"2.0.0", "2.0.9", "2.1.0-pre"},
			incompatible: []string{"2.1.0", "1.9.9", "2.0.0-pre"},
		},
	}
	for _, tt := range test {
		t.Run(tt.library, func(t *testing.T) {
			libVersion, err := Parse(tt.library)
			require.NoError(t, err)

			compatRange := CompatibleRange(libVersion)

			for _, v := range tt.compatible {
				t.Run(v, func(t *testing.T) {
					version, err := Parse(v)
					require.NoError(t, err)
					assert.True(t, compatRange.Contains(version))
				})
			}

			for _, v := range tt.incompatible {
				t.Run(v, func(t *testing.T) {
					version, err := Parse(v)
					require.NoError(t, err)
					assert.False(t, compatRange.Contains(version))
				})
			}
		})
	}
}
