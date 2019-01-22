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
	"io"
	"strings"

	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/internal/multiplex"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/protocol"

	"go.uber.org/atomic"
	"go.uber.org/multierr"
)

var _proto = protocol.Binary

// transportHandle is a Handle to a plugin which is behind an envelope.Transport.
type transportHandle struct {
	name string

	Transport envelope.Transport
	Client    api.Plugin
	Running   *atomic.Bool
	Features  map[api.Feature]struct{}
}

// NewTransportHandle builds a new Handle which speaks to the given transport.
//
// If the transport is an io.Closer, it will be closed when the handle is closed.
func NewTransportHandle(name string, t envelope.Transport) (Handle, error) {
	client := api.NewPluginClient(multiplex.NewClient(
		"Plugin",
		envelope.NewClient(_proto, t),
	))

	handshake, err := client.Handshake(&api.HandshakeRequest{})
	if err != nil {
		return nil, errHandshakeFailed{Name: name, Reason: err}
	}

	if handshake.Name != name {
		return nil, errHandshakeFailed{
			Name:   name,
			Reason: errNameMismatch{Want: name, Got: handshake.Name},
		}
	}

	if handshake.APIVersion != api.APIVersion {
		return nil, errHandshakeFailed{
			Name:   name,
			Reason: errAPIVersionMismatch{Want: api.APIVersion, Got: handshake.APIVersion},
		}
	}

	features := make(map[api.Feature]struct{}, len(handshake.Features))
	for _, feature := range handshake.Features {
		features[feature] = struct{}{}
	}

	return &transportHandle{
		name:      name,
		Transport: t,
		Client:    client,
		Running:   atomic.NewBool(true),
		Features:  features,
	}, nil
}

func (h *transportHandle) Name() string {
	return h.name
}

func (h *transportHandle) Close() error {
	if !h.Running.Swap(false) {
		return nil // already closed
	}

	err := h.Client.Goodbye()
	if closer, ok := h.Transport.(io.Closer); ok {
		err = multierr.Append(err, closer.Close())
	}
	return err
}

func (h *transportHandle) ServiceGenerator() ServiceGenerator {
	if !h.Running.Load() {
		panic(fmt.Sprintf("handle for plugin %q has already been closed", h.name))
	}

	if _, hasFeature := h.Features[api.FeatureServiceGenerator]; !hasFeature {
		return nil
	}

	return &serviceGenerator{
		handle:  h,
		Running: h.Running,
		ServiceGenerator: api.NewServiceGeneratorClient(multiplex.NewClient(
			"ServiceGenerator",
			envelope.NewClient(_proto, h.Transport),
		)),
	}
}

// serviceGenerator is a ServiceGenerator that validates the output of an ServiceGenerator.
//
// It also panics if a request is made to it after it has been closed.
type serviceGenerator struct {
	handle *transportHandle

	ServiceGenerator api.ServiceGenerator
	Running          *atomic.Bool
}

func (sg *serviceGenerator) Handle() Handle {
	return sg.handle
}

func (sg *serviceGenerator) Generate(req *api.GenerateServiceRequest) (*api.GenerateServiceResponse, error) {
	name := sg.handle.name
	if !sg.Running.Load() {
		panic(fmt.Sprintf("handle for plugin %q has already been closed", name))
	}

	res, err := sg.ServiceGenerator.Generate(req)
	if err != nil {
		return res, fmt.Errorf("plugin %q failed to generate service code: %v", name, err)
	}

	for path := range res.Files {
		if strings.Contains(path, "..") {
			return res, fmt.Errorf(
				"plugin %q is attempting to write to a parent directory: "+
					`path %q contains ".."`, name, path)
		}
	}

	return res, nil
}
