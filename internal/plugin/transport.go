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

package plugin

import (
	"fmt"
	"io"
	"strings"

	"go.uber.org/thriftrw/internal"
	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/internal/multiplex"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/plugin/api/service/plugin"
	"go.uber.org/thriftrw/plugin/api/service/servicegenerator"
	"go.uber.org/thriftrw/protocol"

	"github.com/uber-go/atomic"
)

var _proto = protocol.Binary

// transportHandle is a Handle to a plugin which is behind an envelope.Transport.
type transportHandle struct {
	Name      string
	Transport envelope.Transport
	Client    api.Plugin
	Running   *atomic.Bool
	Features  map[api.Feature]struct{}
}

// NewTransportHandle builds a new Handle which speaks to the given transport.
//
// If the transport is an io.Closer, it will be closed when the handle is closed.
func NewTransportHandle(name string, t envelope.Transport) (Handle, error) {
	client := plugin.NewClient(multiplex.NewClient(
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

	if handshake.ApiVersion != api.Version {
		return nil, errHandshakeFailed{
			Name:   name,
			Reason: errVersionMismatch{Want: api.Version, Got: handshake.ApiVersion},
		}
	}

	features := make(map[api.Feature]struct{}, len(handshake.Features))
	for _, feature := range handshake.Features {
		features[feature] = struct{}{}
	}

	return &transportHandle{
		Name:      name,
		Transport: t,
		Client:    client,
		Running:   atomic.NewBool(true),
		Features:  features,
	}, nil
}

func (h *transportHandle) Close() error {
	if !h.Running.Swap(false) {
		return nil // already closed
	}

	err := h.Client.Goodbye()
	if closer, ok := h.Transport.(io.Closer); ok {
		err = internal.CombineErrors(err, closer.Close())
	}
	return err
}

func (h *transportHandle) ServiceGenerator() api.ServiceGenerator {
	if !h.Running.Load() {
		panic(fmt.Sprintf("handle for plugin %q has already been closed", h.Name))
	}

	if _, hasFeature := h.Features[api.FeatureServiceGenerator]; !hasFeature {
		return nil
	}

	return &serviceGenerator{
		Name:    h.Name,
		Running: h.Running,
		ServiceGenerator: servicegenerator.NewClient(multiplex.NewClient(
			"ServiceGenerator",
			envelope.NewClient(_proto, h.Transport),
		)),
	}
}

// serviceGenerator is a ServiceGenerator that validates the output of an api.ServiceGenerator.
//
// It also panics if a request is made to it after it has been closed.
type serviceGenerator struct {
	Name             string
	ServiceGenerator api.ServiceGenerator
	Running          *atomic.Bool
}

func (sg *serviceGenerator) Generate(req *api.GenerateServiceRequest) (*api.GenerateServiceResponse, error) {
	if !sg.Running.Load() {
		panic(fmt.Sprintf("handle for plugin %q has already been closed", sg.Name))
	}

	res, err := sg.ServiceGenerator.Generate(req)
	if err != nil {
		return res, err
	}

	for path := range res.Files {
		if strings.Contains(path, "..") {
			return res, fmt.Errorf(
				"plugin %q is attempting to write to a parent directory: "+
					`path %q contains ".."`, sg.Name, path)
		}
	}

	return res, nil
}
