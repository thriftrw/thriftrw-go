// Copyright (c) 2015 Uber Technologies, Inc.
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

package curry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOne(t *testing.T) {
	tests := []func(*bool){
		func(called *bool) {
			One(func(i int) {
				*called = true
				assert.Equal(t, i, 42)
			}, 42).(func())()
		},
		func(called *bool) {
			One(func(s string, i int) {
				*called = true
				assert.Equal(t, s, "hello")
				assert.Equal(t, i, 42)
			}, "hello").(func(int))(42)
		},
	}

	for _, tt := range tests {
		called := false
		tt(&called)
		assert.True(t, called)
	}
}

func TestOneFailures(t *testing.T) {
	tests := []func(){
		func() {
			One(func(i int) {}, nil)
		},
		func() {
			One(nil, 42)
		},
		func() {
			One(func() {}, 42)
		},
	}

	for _, tt := range tests {
		assert.Panics(t, tt)
	}
}
