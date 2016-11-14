package plugin

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/internal/envelope/envelopetest"
	"go.uber.org/thriftrw/internal/frame"
	"go.uber.org/thriftrw/internal/multiplex"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/plugin/plugintest"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/version"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakePluginServer struct {
	errCh  chan error
	server *frame.Server

	ClientTransport  envelope.Transport
	Plugin           *plugintest.MockPlugin
	ServiceGenerator *plugintest.MockServiceGenerator
}

func newFakePluginServer(mockCtrl *gomock.Controller) *fakePluginServer {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()

	server := frame.NewServer(serverReader, serverWriter)
	client := frame.NewClient(clientWriter, clientReader)

	mockPlugin := plugintest.NewMockPlugin(mockCtrl)
	mockServiceGenerator := plugintest.NewMockServiceGenerator(mockCtrl)

	handler := multiplex.NewHandler()
	handler.Put("Plugin", api.NewPluginHandler(mockPlugin))
	handler.Put("ServiceGenerator", api.NewServiceGeneratorHandler(mockServiceGenerator))

	done := make(chan error)
	go func() {
		err := server.Serve(envelope.NewServer(protocol.Binary, handler))
		if err != nil {
			done <- err
		}
		close(done)
	}()

	return &fakePluginServer{
		errCh:            done,
		server:           server,
		ClientTransport:  client,
		Plugin:           mockPlugin,
		ServiceGenerator: mockServiceGenerator,
	}
}

// Handshake with this server without any expectation of failure.
func (s *fakePluginServer) Handshake(t *testing.T, pluginName string, features []api.Feature) Handle {
	s.Plugin.EXPECT().Handshake(&api.HandshakeRequest{}).
		Return(&api.HandshakeResponse{
			Name:           pluginName,
			APIVersion:     api.APIVersion,
			LibraryVersion: ptr.String(version.Version),
			Features:       features,
		}, nil)

	handle, err := NewTransportHandle(pluginName, s.ClientTransport)
	require.NoError(t, err, "handshake with fakePluginServer failed")

	return handle
}

func (s *fakePluginServer) ExpectGoodbye() {
	s.Plugin.EXPECT().Goodbye().Return(nil)
}

func (s *fakePluginServer) Close() error {
	s.server.Stop()
	err := <-s.errCh
	return err
}

// Wraps an envelope.Transport to become an io.Closer.
type transportCloser struct {
	envelope.Transport

	CloseError error
	WasClosed  bool
}

func (t *transportCloser) Close() error {
	t.WasClosed = true
	return t.CloseError
}

func TestTransportHandleTransportError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	transport := envelopetest.NewMockTransport(mockCtrl)
	transport.EXPECT().Send(gomock.Any()).Return(nil, errors.New("great sadness"))

	_, err := NewTransportHandle("foo", transport)
	require.Error(t, err)
	assert.Equal(t, err.Error(), `handshake with plugin "foo" failed: great sadness`)
}

func TestTransportHandleHandshakeError(t *testing.T) {
	tests := []struct {
		desc      string
		name      string
		response  *api.HandshakeResponse
		wantError string
	}{
		{
			desc: "name mismatch",
			name: "foo",
			response: &api.HandshakeResponse{
				Name:           "bar",
				APIVersion:     api.APIVersion,
				LibraryVersion: ptr.String(version.Version),
				Features:       []api.Feature{},
			},
			wantError: `handshake with plugin "foo" failed: ` +
				`plugin name mismatch: expected "foo" but got "bar"`,
		},
		{
			desc: "version mismatch",
			name: "foo",
			response: &api.HandshakeResponse{
				Name:           "foo",
				APIVersion:     42,
				LibraryVersion: ptr.String(version.Version),
				Features:       []api.Feature{},
			},
			wantError: `handshake with plugin "foo" failed: ` +
				fmt.Sprintf("plugin API version mismatch: expected %d but got 42", api.APIVersion),
		},
		{
			desc: "missing version",
			name: "foo",
			response: &api.HandshakeResponse{
				Name:           "foo",
				APIVersion:     api.APIVersion,
				LibraryVersion: nil,
				Features:       []api.Feature{},
			},
			wantError: `handshake with plugin "foo" failed: Version is required`,
		},
		{
			desc: "unparseable version",
			name: "foo",
			response: &api.HandshakeResponse{
				Name:           "foo",
				APIVersion:     api.APIVersion,
				LibraryVersion: ptr.String("hello"),
				Features:       []api.Feature{},
			},
			wantError: `handshake with plugin "foo" failed: ` +
				`cannot parse as semantic version: "hello"`,
		},
		{
			desc: "semver mismatch",
			name: "foo",
			response: &api.HandshakeResponse{
				Name:           "foo",
				APIVersion:     api.APIVersion,
				LibraryVersion: ptr.String("12.3.4"),
				Features:       []api.Feature{},
			},
			wantError: `handshake with plugin "foo" failed: ` +
				"plugin compiled with the wrong version of ThriftRW: " +
				fmt.Sprintf("expected >=%v and <%v but got 12.3.4", &compatRange.Begin, &compatRange.End),
		},
	}

	for _, tt := range tests {
		func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			server := newFakePluginServer(mockCtrl)
			defer server.Close()

			server.Plugin.EXPECT().Handshake(&api.HandshakeRequest{}).Return(tt.response, nil)
			_, err := NewTransportHandle(tt.name, server.ClientTransport)
			if assert.Error(t, err, tt.desc) {
				assert.Equal(t, err.Error(), tt.wantError, tt.desc)
			}
		}()
	}
}

func TestTransportHandleServiceGenerator(t *testing.T) {
	tests := []struct {
		desc                string
		features            []api.Feature
		hasServiceGenerator bool
	}{
		{
			desc:     "no ServiceGenerator",
			features: []api.Feature{},
		},
		{
			desc:                "has ServiceGenerator",
			features:            []api.Feature{api.FeatureServiceGenerator},
			hasServiceGenerator: true,
		},
	}

	for _, tt := range tests {
		func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			server := newFakePluginServer(mockCtrl)
			defer server.Close()

			handle := server.Handshake(t, "foo", tt.features)

			sg := handle.ServiceGenerator()
			if tt.hasServiceGenerator {
				assert.NotNil(t, sg, tt.desc)
			} else {
				assert.Nil(t, sg, tt.desc)
			}

			server.ExpectGoodbye()
			assert.NoError(t, handle.Close(), tt.desc)
		}()
	}
}

func TestTransportHandleDoubleClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := newFakePluginServer(mockCtrl)
	defer server.Close()

	handle := server.Handshake(t, "foo", []api.Feature{})

	server.ExpectGoodbye()
	assert.NoError(t, handle.Close())
	assert.NoError(t, handle.Close())
}

func TestTransportHandleCloseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := newFakePluginServer(mockCtrl)
	defer server.Close()

	handle := server.Handshake(t, "foo", []api.Feature{})

	server.Plugin.EXPECT().Goodbye().Return(errors.New("great sadness"))
	if err := handle.Close(); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "great sadness")
	}
}

func TestTransportHandleCloseCloser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := newFakePluginServer(mockCtrl)
	defer server.Close()

	transport := &transportCloser{Transport: server.ClientTransport}
	server.ClientTransport = transport

	handle := server.Handshake(t, "foo", []api.Feature{})

	server.ExpectGoodbye()
	if err := handle.Close(); assert.NoError(t, err) {
		assert.True(t, transport.WasClosed, "expected Close() to be called")
	}
}

func TestTransportHandleCloserError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := newFakePluginServer(mockCtrl)
	defer server.Close()

	transport := &transportCloser{
		Transport:  server.ClientTransport,
		CloseError: errors.New("great sadness"),
	}
	server.ClientTransport = transport

	handle := server.Handshake(t, "foo", []api.Feature{})

	server.ExpectGoodbye()
	if err := handle.Close(); assert.Error(t, err) {
		assert.True(t, transport.WasClosed, "expected Close() to be called")
		assert.Equal(t, "great sadness", err.Error())
	}
}

func TestTransportHandleCloseServiceGenerator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := newFakePluginServer(mockCtrl)
	defer server.Close()

	handle := server.Handshake(t, "foo", []api.Feature{})

	server.ExpectGoodbye()
	require.NoError(t, handle.Close())

	assert.Panics(t, func() {
		handle.ServiceGenerator()
	})
}

func TestServiceGeneratorClosed(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	server := newFakePluginServer(mockCtrl)
	defer server.Close()

	handle := server.Handshake(t, "foo", []api.Feature{api.FeatureServiceGenerator})
	sg := handle.ServiceGenerator()

	server.ExpectGoodbye()
	require.NoError(t, handle.Close())

	assert.Panics(t, func() {
		sg.Generate(&api.GenerateServiceRequest{})
	})
}

func TestServiceGeneratorGenerate(t *testing.T) {
	tests := []struct {
		desc             string
		generateResponse *api.GenerateServiceResponse
		generateError    error

		wantError string
	}{
		{
			desc: "success",
			generateResponse: &api.GenerateServiceResponse{
				Files: map[string][]byte{"foo/bar.go": []byte("package foo")},
			},
		},
		{
			desc: "parent directory",
			generateResponse: &api.GenerateServiceResponse{
				Files: map[string][]byte{"../foo/bar.go": []byte("package foo")},
			},
			wantError: `plugin "foo" is attempting to write to a parent directory: ` +
				`path "../foo/bar.go" contains ".."`,
		},
		{
			desc:          "call error",
			generateError: errors.New("great sadness"),
			wantError: `plugin "foo" failed to generate service code: ` +
				"TApplicationException{Message: great sadness, Type: INTERNAL_ERROR}",
		},
	}

	for _, tt := range tests {
		func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			server := newFakePluginServer(mockCtrl)
			defer server.Close()

			handle := server.Handshake(t, "foo", []api.Feature{api.FeatureServiceGenerator})
			defer func() {
				server.ExpectGoodbye()
				require.NoError(t, handle.Close(), tt.desc)
			}()

			req := &api.GenerateServiceRequest{
				RootServices: []api.ServiceID{1},
				Services: map[api.ServiceID]*api.Service{
					1: {
						Name:       "KeyValue",
						ThriftName: "KeyValue",
						Functions:  []*api.Function{},
						ModuleID:   api.ModuleID(1),
					},
				},
				Modules: map[api.ModuleID]*api.Module{
					1: {
						ImportPath: "go.uber.org/thriftrw/foo",
						Directory:  "foo",
					},
				},
			}

			server.ServiceGenerator.EXPECT().Generate(req).
				Return(tt.generateResponse, tt.generateError)

			res, err := handle.ServiceGenerator().Generate(req)
			if tt.wantError != "" {
				if assert.Error(t, err, tt.desc) {
					assert.Equal(t, tt.wantError, err.Error(), tt.desc)
				}
			} else {
				assert.NoError(t, err, tt.desc)
				assert.Equal(t, tt.generateResponse, res, tt.desc)
			}
		}()
	}
}
