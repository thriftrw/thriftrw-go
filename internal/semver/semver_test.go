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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		GreaterVersion string
		LesserVersion  string
	}{
		{"0.0.0", "0.0.0-foo"},
		{"0.0.1", "0.0.0"},
		{"0.10.0", "0.9.0"},
		{"0.99.0", "0.10.0"},
		{"1.0.0", "0.9.9"},
		{"1.0.0", "1.0.0-rc.1"},
		{"1.0.0-alpha.1", "1.0.0-alpha"},
		{"1.0.0-alpha.beta", "1.0.0-alpha.1"},
		{"1.0.0-beta", "1.0.0-alpha.beta"},
		{"1.0.0-beta.11", "1.0.0-beta.2"},
		{"1.0.0-beta.2", "1.0.0-beta"},
		{"1.0.0-rc.1", "1.0.0-beta.11"},
		{"1.0.0-rc.2", "1.0.0-rc.1"},
		{"1.2.3", "1.2.3-4"},
		{"1.2.3", "1.2.3-4-foo"},
		{"1.2.3", "1.2.3-asdf"},
		{"1.2.3-5", "1.2.3-4"},
		{"1.2.3-5-foo", "1.2.3-5"},
		{"1.2.3-5-foo", "1.2.3-5-Foo"},
		{"1.2.3-a.10", "1.2.3-a.5"},
		{"1.2.3-a.b", "1.2.3-a"},
		{"1.2.3-a.b", "1.2.3-a.5"},
		{"1.2.3-a.b.c.10.d.5", "1.2.3-a.b.c.5.d.100"},
		{"2.0.0", "1.2.3"},
		{"3.0.0", "2.7.2+asdf"},
		{"3.0.0+foobar", "2.7.2"},
		{"3.0.0-hello.42+foobar.meta.39", "3.0.0-42.42+barfoo.tame.93"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s > %s", test.GreaterVersion, test.LesserVersion), func(t *testing.T) {
			gv, err := Parse(test.GreaterVersion)
			require.NoError(t, err)
			lv, err := Parse(test.LesserVersion)
			require.NoError(t, err)
			assert.Equal(t, test.GreaterVersion, gv.String(), "GreaterVersion input != parsed output")
			assert.Equal(t, test.LesserVersion, lv.String(), "LesserVersion input != parsed output")
			assert.Equal(t, 1, gv.Compare(&lv), "Greater version must be greater than")
			assert.Equal(t, -1, lv.Compare(&gv), "Lesser version must be less than")
		})
	}
}

func TestCompareEqual(t *testing.T) {
	tests := []struct {
		A string
		B string
	}{
		{"0.0.0", "0.0.0"},
		{"0.0.0-foo", "0.0.0-foo"},
		{"0.0.1", "0.0.1"},
		{"0.1.0", "0.1.0"},
		{"1.0.0", "1.0.0"},
		{"1.0.0-rc.1", "1.0.0-rc.1"},
		{"1.0.0-alpha.beta", "1.0.0-alpha.beta"},
		{"1.2.3-5-foo", "1.2.3-5-foo"},
		{"3.0.0", "3.0.0+asdf"},
		{"3.0.0+foobar", "3.0.0"},
		{"3.0.0-hello.42+foobar.meta.39", "3.0.0-hello.42+barfoo.tame.93"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s == %s", test.A, test.B), func(t *testing.T) {
			av, err := Parse(test.A)
			require.NoError(t, err)
			bv, err := Parse(test.B)
			require.NoError(t, err)
			assert.Equal(t, test.A, av.String(), "a input != parsed output")
			assert.Equal(t, test.B, bv.String(), "b input != parsed output")
			assert.Equal(t, 0, av.Compare(&bv), "a must be equal to b")
			assert.Equal(t, 0, bv.Compare(&av), "b must be equal to a")
		})
	}
}
