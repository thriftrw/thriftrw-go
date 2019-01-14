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
	"log"
	"os"

	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/internal/frame"
	"go.uber.org/thriftrw/internal/multiplex"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/version"
)

const _fastPathFrameSize = 10 * 1024 * 1024 // 10 MB

var _proto = protocol.Binary

// Plugin defines a ThriftRW plugin.
//
// At minimum, a plugin name must be provided and it MUST match the name of
// the plugin in the executable.
type Plugin struct {
	// Name of the plugin. The name of the executable providing this plugin MUST
	// be thriftrw-plugin-$name.
	Name string

	// Plugins can implement a ServiceGenerator to generate arbitrary code for
	// Thrift services. This may be nil if the plugin does not provide this
	// functionality.
	ServiceGenerator api.ServiceGenerator

	// Reader and Writer may be specified to change the communication channel
	// this plugin uses. By default, plugins listen on stdin and write to
	// stdout.
	Reader io.Reader
	Writer io.Writer
}

// Main serves the given plugin. It is the entry point to the plugin system.
// User-defined plugins should call Main with their main function.
//
// 	func main() {
// 		plugin.Main(myPlugin)
// 	}
func Main(p *Plugin) {
	if p.Name == "" {
		panic("a plugin name must be provided")
	}

	// The plugin communicates with the ThriftRW process over stdout and stdin
	// of this process. Requests and responses are Thrift envelopes with a
	// 4-byte big-endian encoded length prefix. Envelope names contain method
	// names prefixed with the service name and a ":".

	mainHandler := multiplex.NewHandler()

	features := []api.Feature{}

	if p.ServiceGenerator != nil {
		features = append(features, api.FeatureServiceGenerator)
		mainHandler.Put("ServiceGenerator", api.NewServiceGeneratorHandler(p.ServiceGenerator))
	}

	// TODO(abg): Check for other features and register handlers here.
	reader := p.Reader
	if reader == nil {
		reader = os.Stdin
	}
	writer := p.Writer
	if writer == nil {
		writer = os.Stdout
	}

	server := frame.NewServer(reader, writer)
	mainHandler.Put("Plugin", api.NewPluginHandler(pluginHandler{
		server:   server,
		plugin:   p,
		features: features,
	}))

	if err := server.Serve(envelope.NewServer(_proto, mainHandler)); err != nil {
		log.Fatalf("plugin server failed with error: %v", err)
	}
}

// pluginHandler implements the Plugin service.
type pluginHandler struct {
	server   *frame.Server
	plugin   *Plugin
	features []api.Feature
}

func (h pluginHandler) Handshake(request *api.HandshakeRequest) (*api.HandshakeResponse, error) {
	return &api.HandshakeResponse{
		Name:           h.plugin.Name,
		APIVersion:     api.APIVersion,
		Features:       h.features,
		LibraryVersion: ptr.String(version.Version),
	}, nil
}

func (h pluginHandler) Goodbye() error {
	return h.server.Stop()
}
