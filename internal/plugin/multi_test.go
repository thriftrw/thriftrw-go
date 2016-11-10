package plugin_test // because import cycle

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	. "go.uber.org/thriftrw/internal/plugin"
	"go.uber.org/thriftrw/internal/plugin/handletest"
	"go.uber.org/thriftrw/plugin/api"
)

func TestMultiHandleClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handles := make([]*handletest.MockHandle, 1000)
	var mh MultiHandle

	for i := range handles {
		handle := handletest.NewMockHandle(mockCtrl)
		handle.EXPECT().Close().Return(nil)
		handles[i] = handle
		mh = append(mh, handle)
	}

	assert.NoError(t, mh.Close())
}

func TestMultiHandleCloseNil(t *testing.T) {
	var mh MultiHandle
	assert.NoError(t, mh.Close())
}

func TestMultiHandleServiceGeneratorNil(t *testing.T) {
	var mh MultiHandle
	mh.ServiceGenerator() // should not panic
}

func TestMultiHandleCloseError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handles := make([]*handletest.MockHandle, 1000)
	var mh MultiHandle

	for i := range handles {
		handle := handletest.NewMockHandle(mockCtrl)
		handles[i] = handle
		mh = append(mh, handle)

		// fail only odd plugins
		if i%2 == 0 {
			handle.EXPECT().Close().Return(nil)
		} else {
			handle.EXPECT().Close().Return(fmt.Errorf("plugin-%d: great sadness", i))
		}
	}

	err := mh.Close()
	require.Error(t, err)

	errMsg := err.Error()
	for i := range handles {
		if i%2 == 0 {
			continue
		}

		assert.Contains(t, errMsg, fmt.Sprintf("plugin-%d: great sadness", i))
	}
}

func TestMultiHandleNoServiceGenerators(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handles := make([]*handletest.MockHandle, 1000)
	var mh MultiHandle
	for i := range handles {
		handle := handletest.NewMockHandle(mockCtrl)
		handles[i] = handle
		mh = append(mh, handle)
		handle.EXPECT().ServiceGenerator().Return(nil)
	}

	assert.Empty(t, mh.ServiceGenerator())
}

func TestMultiHandleServiceGenerator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	handles := make([]*handletest.MockHandle, 1000)

	var mh MultiHandle
	for i := range handles {
		handle := handletest.NewMockHandle(mockCtrl)
		handles[i] = handle
		mh = append(mh, handle)

		// only odd handles have a ServiceGenerator
		if i%2 == 0 {
			handle.EXPECT().ServiceGenerator().Return(nil)
			continue
		}

		handle.EXPECT().ServiceGenerator().Return(
			handletest.NewMockServiceGenerator(mockCtrl))
	}

	assert.NotNil(t, mh.ServiceGenerator())
}

func TestMultiServiceGeneratorGenerate(t *testing.T) {
	type response struct {
		success *api.GenerateServiceResponse
		failure error
	}

	tests := []struct {
		desc string

		// list of responses from different service generators
		responses []response

		// final expected response or errors
		wantResponse    *api.GenerateServiceResponse
		wantErrors      []string
		wantOneOfErrors []string
		// both, wantErrors and wantOneOfErrors may be set. All errors in
		// wantErrors must be present, but only one or more of the errors in
		// wantOneOfErrors must be present.
	}{
		{
			desc: "no conflicts; no errors",
			responses: []response{
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/a.go": {1, 2, 3},
					"foo/b.go": {4, 5, 6},
				}}},
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/c.go": {7, 8, 9},
					"foo/d.go": {1, 2, 3},
				}}},
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
				// no files
				}}},
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/keyvalue/e.go": {4, 5, 6},
				}}},
			},
			wantResponse: &api.GenerateServiceResponse{Files: map[string][]byte{
				"foo/a.go":          {1, 2, 3},
				"foo/b.go":          {4, 5, 6},
				"foo/c.go":          {7, 8, 9},
				"foo/d.go":          {1, 2, 3},
				"foo/keyvalue/e.go": {4, 5, 6},
			}},
		},
		{
			desc: "no conflicts; with errors",
			responses: []response{
				{failure: errors.New("foo: great sadness")},
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/a.go": {1, 2, 3},
				}}},
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/b.go": {4, 5, 6},
				}}},
				{failure: errors.New("bar: great sadness")},
			},
			wantErrors: []string{
				`foo: great sadness`,
				`bar: great sadness`,
			},
		},
		{
			desc: "conflicts",
			responses: []response{
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/a.go": {1, 2, 3},
					"foo/b.go": {4, 5, 6},
				}}},
				{success: &api.GenerateServiceResponse{Files: map[string][]byte{
					"foo/c.go": {7, 8, 9},
					"foo/b.go": {1, 2, 3},
				}}},
			},
			wantErrors: []string{`plugin conflict: cannot write file "foo/b.go" for plugin`},
			wantOneOfErrors: []string{
				`plugin "plugin-1" already wrote to that file`,
				`plugin "plugin-0" already wrote to that file`,
			},
		},
	}

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

	for _, tt := range tests {
		func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			var msg MultiServiceGenerator
			for i, res := range tt.responses {
				handle := handletest.NewMockHandle(mockCtrl)
				handle.EXPECT().Name().Return(fmt.Sprintf("plugin-%d", i)).AnyTimes()

				sg := handletest.NewMockServiceGenerator(mockCtrl)
				msg = append(msg, sg)
				sg.EXPECT().Generate(req).Return(res.success, res.failure)
				sg.EXPECT().Handle().Return(handle).AnyTimes()
			}

			res, err := msg.Generate(req)
			if len(tt.wantErrors) > 0 || len(tt.wantOneOfErrors) > 0 {
				if !assert.Error(t, err, tt.desc) {
					return
				}

				for _, errMsg := range tt.wantErrors {
					assert.Contains(t, err.Error(), errMsg, tt.desc)
				}

				matches := len(tt.wantOneOfErrors) == 0
				for _, errMsg := range tt.wantOneOfErrors {
					if strings.Contains(err.Error(), errMsg) {
						matches = true
						break
					}
				}

				assert.True(t, matches, "expected %v to contain one of %v", err, tt.wantOneOfErrors)
			} else {
				assert.Equal(t, tt.wantResponse, res, tt.desc)
			}
		}()
	}
}

func TestMultiServiceGeneratorGenerateNil(t *testing.T) {
	var msg MultiServiceGenerator
	_, err := msg.Generate(&api.GenerateServiceRequest{})
	assert.NoError(t, err)
}
