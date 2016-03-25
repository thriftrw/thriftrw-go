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

package gen

// hasReaders helps generators track whether a reader for the given name was
// already generated.
type hasReaders struct{ have map[string]struct{} }

// HasReader returns true if a reader with the given name was already generated,
// otherwise it returns false for this call but marks the reader as generated
// for all consecutive calls.
func (r *hasReaders) HasReader(name string) bool {
	if r.have == nil {
		r.have = make(map[string]struct{})
	}

	if _, ok := r.have[name]; ok {
		return ok
	}

	r.have[name] = struct{}{}
	return false
}

// hasLazyLists helps generators track whether a lazy list for the given name
// was already generated.
type hasLazyLists struct{ have map[string]struct{} }

// HasLazyList returns true if a lazy list with the given name was already generated,
// otherwise it returns false for this call but marks the lazy list as generated
// for all consecutive calls.
func (r *hasLazyLists) HasLazyList(name string) bool {
	if r.have == nil {
		r.have = make(map[string]struct{})
	}

	if _, ok := r.have[name]; ok {
		return ok
	}

	r.have[name] = struct{}{}
	return false
}
