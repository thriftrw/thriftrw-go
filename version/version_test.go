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

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodeCompatRange(t *testing.T) {
	test := []struct {
		libVer     string
		compatVer  []string
		invalidVer []string
	}{
		{
			"1.2.3",
			[]string{"1.2.3", "1.2.4", "1.0.0"},
			[]string{"1.3.0", "0.9.9"},
		},
		{
			"2.0.0",
			[]string{"2.0.0", "2.0.9", "2.1.0-pre"},
			[]string{"2.1.0", "1.9.9", "2.0.0-pre"},
		},
	}
	for _, tt := range test {
		t.Run(tt.libVer, func(t *testing.T) {
			compatRange := computeGenCodeCompabilityRange(tt.libVer)
			for _, compatVer := range tt.compatVer {
				t.Run(compatVer, func(t *testing.T) {
					genv, err := parseSemVer(compatVer)
					require.NoError(t, err)
					assert.True(t, compatRange.IsCompatibleWith(genv))
				})
			}
			for _, invalidVer := range tt.invalidVer {
				t.Run(invalidVer, func(t *testing.T) {
					genv, err := parseSemVer(invalidVer)
					require.NoError(t, err)
					assert.False(t, compatRange.IsCompatibleWith(genv))
				})
			}
		})
	}
}
