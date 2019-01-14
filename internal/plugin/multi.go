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
	"strings"
	"sync"

	"go.uber.org/thriftrw/internal/concurrent"
	"go.uber.org/thriftrw/plugin/api"
)

// MultiHandle wraps a collection of handles into a single handle.
type MultiHandle []Handle

// Name is the combined name of this plugin
func (mh MultiHandle) Name() string {
	names := make([]string, 0, len(mh))
	for _, h := range mh {
		names = append(names, h.Name())
	}
	return fmt.Sprintf("MultiHandle{%v}", strings.Join(names, ", "))
}

// Close closes all Handles associated with this MultiHandle.
func (mh MultiHandle) Close() error {
	return concurrent.Range(mh, func(_ int, h Handle) error {
		return h.Close()
	})
}

// ServiceGenerator returns a ServiceGenerator which calls into the
// ServiceGenerators of all plugins associated with this MultiHandle and
// consolidates their results.
func (mh MultiHandle) ServiceGenerator() ServiceGenerator {
	msg := make(MultiServiceGenerator, 0, len(mh))
	for _, h := range mh {
		if sg := h.ServiceGenerator(); sg != nil {
			msg = append(msg, sg)
		}
	}
	return msg
}

// MultiServiceGenerator wraps a collection of ServiceGenerators into a single
// ServiceGenerator.
type MultiServiceGenerator []ServiceGenerator

// Handle returns a reference to the Handle that owns this ServiceGenerator.
func (msg MultiServiceGenerator) Handle() Handle {
	mh := make(MultiHandle, len(msg))
	for i, sg := range msg {
		mh[i] = sg.Handle()
	}
	return mh
}

// Generate calls all the service generators associated with this plugin and
// consolidates their output.
//
// Any conflicts in the generated files will result in a failure.
func (msg MultiServiceGenerator) Generate(req *api.GenerateServiceRequest) (*api.GenerateServiceResponse, error) {
	var (
		lock      sync.Mutex
		files     = make(map[string][]byte)
		usedPaths = make(map[string]string) // path -> plugin name
	)

	err := concurrent.Range(msg, func(_ int, sg ServiceGenerator) error {
		res, err := sg.Generate(req)
		if err != nil {
			return err
		}

		lock.Lock()
		defer lock.Unlock()

		pluginName := sg.Handle().Name()
		for path, contents := range res.Files {
			if takenBy, taken := usedPaths[path]; taken {
				return fmt.Errorf("plugin conflict: cannot write file %q for plugin %q: "+
					"plugin %q already wrote to that file", path, pluginName, takenBy)
			}

			usedPaths[path] = pluginName
			files[path] = contents
		}

		return nil
	})

	return &api.GenerateServiceResponse{Files: files}, err
}
