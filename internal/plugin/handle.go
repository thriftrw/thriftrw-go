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
	"io"

	"go.uber.org/thriftrw/plugin/api"
)

// Handle is a handle to ThriftRW plugin.
type Handle interface {
	io.Closer

	// Name is the name of this plugin.
	Name() string

	// ServiceGenerator returns a ServiceGenerator for this plugin or nil if
	// this plugin does not implement that feature.
	//
	// Note that the ServiceGenerator is valid only as long as Close is not
	// called on the Handle.
	ServiceGenerator() ServiceGenerator
}

// ServiceGenerator generates files for Thrift services.
type ServiceGenerator interface {
	api.ServiceGenerator

	// Handle returns the Handle that owns this ServiceGenerator.
	Handle() Handle
}
