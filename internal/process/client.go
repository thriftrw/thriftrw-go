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

package process

import (
	"fmt"
	"io"
	"os/exec"

	"go.uber.org/thriftrw/internal/frame"

	"go.uber.org/atomic"
	"go.uber.org/multierr"
)

// Client sends framed requests and receives framed responses from an external
// process.
type Client struct {
	running *atomic.Bool
	cmd     *exec.Cmd
	stdout  io.ReadCloser
	stdin   io.WriteCloser
	client  *frame.Client
}

// NewClient starts up the given external process and communicates with it over
// stdin and stdout using framed requests and responses.
//
// The Cmd MUST NOT have Stdout or Stdin set.
func NewClient(cmd *exec.Cmd) (*Client, error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe to %q: %v", cmd.Path, err)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe to %q: %v", cmd.Path, err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start %q: %v", cmd.Path, err)
	}

	return &Client{
		stdout:  stdout,
		stdin:   stdin,
		running: atomic.NewBool(true),
		client:  frame.NewClient(stdin, stdout),
		cmd:     cmd,
	}, nil
}

// Send sends the given frame to the external process and returns the response.
//
// Panics if Close was already called.
func (c *Client) Send(data []byte) ([]byte, error) {
	if !c.running.Load() {
		panic(fmt.Sprintf("process.Client for %q has been closed", c.cmd.Path))
	}

	return c.client.Send(data)
}

// Close detaches from the external process and waits for it to exit.
func (c *Client) Close() error {
	if !c.running.Swap(false) {
		return nil // already stopped
	}

	var errors []error
	if err := c.stdout.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to detach stdout from %q: %v", c.cmd.Path, err))
	}
	if err := c.stdin.Close(); err != nil {
		errors = append(errors, fmt.Errorf("failed to detach stdin from %q: %v", c.cmd.Path, err))
	}
	if err := c.cmd.Wait(); err != nil {
		errors = append(errors, fmt.Errorf("%q failed with: %v", c.cmd.Path, err))
	}
	return multierr.Combine(errors...)
}
