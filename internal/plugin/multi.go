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
	"sync"

	"go.uber.org/thriftrw/internal/concurrent"
	"go.uber.org/thriftrw/plugin/api"
)

// MultiHandle wraps a collection of handles into a single handle.
type MultiHandle map[string]Handle

// Close closes all Handles associated with this MultiHandle.
func (mh MultiHandle) Close() error {
	return concurrent.Range(mh, func(name string, h Handle) error {
		if err := h.Close(); err != nil {
			return fmt.Errorf("plugin %q failed to close: %v", name, err)
		}
		return nil
	})
}

// ServiceGenerator returns a ServiceGenerator which calls into the
// ServiceGenerators of all plugins associated with this MultiHandle and
// consolidates their results.
func (mh MultiHandle) ServiceGenerator() api.ServiceGenerator {
	msg := make(MultiServiceGenerator, len(mh))
	for name, h := range mh {
		sg := h.ServiceGenerator()
		if sg != nil {
			msg[name] = sg
		}
	}
	if len(msg) == 0 {
		return nil
	}
	return msg
}

// MultiServiceGenerator wraps a collection of ServiceGenerators into a single
// ServiceGenerator.
type MultiServiceGenerator map[string]api.ServiceGenerator

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

	err := concurrent.Range(msg, func(name string, sg api.ServiceGenerator) error {
		res, err := sg.Generate(req)
		if err != nil {
			return fmt.Errorf("plugin %q failed to generate service code: %v", name, err)
		}

		lock.Lock()
		defer lock.Unlock()

		for path, contents := range res.Files {
			if takenBy, taken := usedPaths[path]; taken {
				return fmt.Errorf("plugin conflict: cannot write file %q for plugin %q: "+
					"plugin %q already wrote to that file", path, name, takenBy)
			}

			usedPaths[path] = name
			files[path] = contents
		}

		return nil
	})

	return &api.GenerateServiceResponse{Files: files}, err
}
