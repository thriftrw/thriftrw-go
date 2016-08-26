package plugin

import (
	"io"
	"testing"

	"github.com/thriftrw/thriftrw-go/internal/envelope"
	"github.com/thriftrw/thriftrw-go/internal/frame"
	"github.com/thriftrw/thriftrw-go/internal/multiplex"
	"github.com/thriftrw/thriftrw-go/plugin/api"
	"github.com/thriftrw/thriftrw-go/plugin/api/service/plugin"
	"github.com/thriftrw/thriftrw-go/plugin/api/service/servicegenerator"
	"github.com/thriftrw/thriftrw-go/plugin/plugintest"
	"github.com/thriftrw/thriftrw-go/ptr"

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

func TestEmptyPlugin(t *testing.T) {
	transport, done := fakeEnvelopeClient()
	defer done()

	go Main(&Plugin{Name: "hello"})

	client := plugin.NewClient(multiplex.NewClient("Plugin", transport))

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

	pluginClient := plugin.NewClient(multiplex.NewClient("Plugin", transport))
	defer pluginClient.Goodbye()

	handshake, err := pluginClient.Handshake(&api.HandshakeRequest{})
	require.NoError(t, err)
	assert.Equal(t, api.Version, handshake.ApiVersion)
	assert.Equal(t, "hello", handshake.Name)
	assert.Contains(t, handshake.Features, api.FeatureServiceGenerator)

	sgClient := servicegenerator.NewClient(multiplex.NewClient("ServiceGenerator", transport))
	req := &api.GenerateServiceRequest{
		RootServices: []api.ServiceID{1},
		Services: map[api.ServiceID]*api.Service{
			1: {
				Name:      "MyService",
				Package:   "github.com/thriftrw/thriftrw-go/plugin/fake/myservice",
				Directory: "fake/myservice",
				Functions: []*api.Function{},
				ParentID:  (*api.ServiceID)(ptr.Int32(2)),
				ModuleID:  1,
			},
			2: {
				Name:      "BaseService",
				Package:   "github.com/thriftrw/thriftrw-go/plugin/fake/baseservice",
				Directory: "fake/baseservice",
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
				Package:   "github.com/thriftrw/thriftrw-go/plugin/fake",
				Directory: "fake",
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
