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

package multiplex

import (
	"strings"

	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/wire"
)

// Handler implements a service multiplexer
type Handler struct {
	services map[string]envelope.Handler
}

// NewHandler builds a new handler.
func NewHandler() Handler {
	return Handler{services: make(map[string]envelope.Handler)}
}

// Put adds the given service to the multiplexer.
func (h Handler) Put(name string, service envelope.Handler) {
	h.services[name] = service
}

// Handle handles the given request, dispatching to one of the
// registered services.
func (h Handler) Handle(name string, req wire.Value) (wire.Value, error) {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) < 2 {
		return wire.Value{}, envelope.ErrUnknownMethod(name)
	}

	service, ok := h.services[parts[0]]
	if !ok {
		return wire.Value{}, envelope.ErrUnknownMethod(name)
	}

	name = parts[1]
	return service.Handle(name, req)
}
