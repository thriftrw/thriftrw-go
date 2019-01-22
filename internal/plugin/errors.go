// Copyright (c) 2019 Uber Technologies, Inc.
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

package plugin

import (
	"fmt"

	"go.uber.org/thriftrw/internal/semver"
)

type errHandshakeFailed struct {
	Name   string
	Reason error
}

func (e errHandshakeFailed) Error() string {
	return fmt.Sprintf("handshake with plugin %q failed: %v", e.Name, e.Reason)
}

type errNameMismatch struct {
	Want, Got string
}

func (e errNameMismatch) Error() string {
	return fmt.Sprintf("plugin name mismatch: expected %q but got %q", e.Want, e.Got)
}

type errAPIVersionMismatch struct {
	Want, Got int32
}

func (e errAPIVersionMismatch) Error() string {
	return fmt.Sprintf("plugin API version mismatch: expected %v but got %v", e.Want, e.Got)
}

type errVersionMismatch struct {
	Want semver.Range
	Got  string
}

func (e errVersionMismatch) Error() string {
	return fmt.Sprintf(
		"plugin compiled with the wrong version of ThriftRW: "+
			"expected >=%v and <%v but got %v", &e.Want.Begin, &e.Want.End, e.Got)
}
