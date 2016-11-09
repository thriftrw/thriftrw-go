package plugin

import (
	"io"
	"testing"

	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/internal/frame"
	"go.uber.org/thriftrw/internal/multiplex"
	"go.uber.org/thriftrw/plugin/api"
	"go.uber.org/thriftrw/plugin/plugintest"
	"go.uber.org/thriftrw/ptr"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeStreams is a helper for tests to control the output and input used by
// plugin.Main while testing.
//
// It returns a Writer to write to the plugin's stdin, a Reader to read from it,
// and a function that should be called after the test is finished to restore
// the old values.
//
// 	in, out, done := fakeStreams()
// 	defer done()
func fakeStreams() (stdin io.Writer, stdout io.Reader, done func()) {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	oldIn := _in
	oldOut := _out
	_in = stdinReader
	_out = stdoutWriter

	return stdinWriter, stdoutReader, func() {
		_in = oldIn
		_out = oldOut
	}
}

func fakeEnvelopeClient() (envelope.Client, func()) {
	in, out, done := fakeStreams()
	return envelope.NewClient(_proto, frame.NewClient(in, out)), done
}

func TestEmptyPluginName(t *testing.T) {
	assert.Panics(t, func() { Main(&Plugin{}) })
}

func TestEmptyPlugin(t *testing.T) {
	transport, done := fakeEnvelopeClient()
	defer done()

	go Main(&Plugin{Name: "hello"})

	client := api.NewPluginClient(multiplex.NewClient("Plugin", transport))

	response, err := client.Handshake(&api.HandshakeRequest{})
	require.NoError(t, err)
	assert.Equal(t, api.Version, response.ApiVersion)
	assert.Equal(t, "hello", response.Name)
	assert.Empty(t, response.Features)

	assert.NoError(t, client.Goodbye())
}

func TestServiceGenerator(t *testing.T) {
	transport, done := fakeEnvelopeClient()
	defer done()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	serviceGenerator := plugintest.NewMockServiceGenerator(mockCtrl)

	go Main(&Plugin{
		Name:             "hello",
		ServiceGenerator: serviceGenerator,
	})

	pluginClient := api.NewPluginClient(multiplex.NewClient("Plugin", transport))
	defer pluginClient.Goodbye()

	handshake, err := pluginClient.Handshake(&api.HandshakeRequest{})
	require.NoError(t, err)
	assert.Equal(t, api.Version, handshake.ApiVersion)
	assert.Equal(t, "hello", handshake.Name)
	assert.Contains(t, handshake.Features, api.FeatureServiceGenerator)

	sgClient := api.NewServiceGeneratorClient(multiplex.NewClient("ServiceGenerator", transport))
	req := &api.GenerateServiceRequest{
		RootServices: []api.ServiceID{1},
		Services: map[api.ServiceID]*api.Service{
			1: {
				Name:       "MyService",
				ThriftName: "MyService",
				Functions:  []*api.Function{},
				ParentID:   (*api.ServiceID)(ptr.Int32(2)),
				ModuleID:   1,
			},
			2: {
				Name:       "BaseService",
				ThriftName: "BaseService",
				Functions: []*api.Function{
					{
						Name:       "Healthy",
						ThriftName: "healthy",
						Arguments:  []*api.Argument{},
					},
				},
				ModuleID: 1,
			},
		},
		Modules: map[api.ModuleID]*api.Module{
			1: {
				ImportPath: "go.uber.org/thriftrw/plugin/fake",
				Directory:  "fake",
			},
		},
	}

	res := &api.GenerateServiceResponse{
		Files: map[string][]byte{
			"fake/myservice/foo.go":   {1, 2, 3},
			"fake/baseservice/bar.go": {4, 5, 6},
			"fake/baz.go":             {7, 8, 9},
		},
	}

	serviceGenerator.EXPECT().Generate(req).Return(res, nil)
	gotRes, err := sgClient.Generate(req)
	if assert.NoError(t, err) {
		assert.Equal(t, res, gotRes)
	}
}
