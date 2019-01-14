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
	"os"
	"os/exec"
	"sync"

	"go.uber.org/thriftrw/internal/concurrent"
	"go.uber.org/thriftrw/internal/process"

	"github.com/anmitsu/go-shlex"
	"go.uber.org/multierr"
)

const _pluginExecPrefix = "thriftrw-plugin-"

// Flag defines an external plugin specification passed over the command line.
//
// Plugin specifications received over the command line are simply plugin
// names followed by arguments for the plugin.
//
// An executable with the name thriftrw-plugin-$name is expected on the $PATH.
// Remaining arguments are passed to the program. For example,
//
// 	-p "foo -a --bc"
//
// Will pass the arguments "-a --bc" to the executable "thriftrw-plugin-foo".
type Flag struct {
	Name    string    // Name of the plugin
	Command *exec.Cmd // Command specification
}

// Handle gets a Handle to this plugin specification.
//
// The returned handle MUST be closed by the caller if error was nil.
func (f *Flag) Handle() (Handle, error) {
	transport, err := process.NewClient(f.Command)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %q: %v", f.Name, err)
	}

	handle, err := NewTransportHandle(f.Name, transport)
	if err != nil {
		return nil, multierr.Combine(
			fmt.Errorf("failed to open plugin %q: %v", f.Name, err),
			transport.Close(),
		)
	}

	return handle, nil
}

// UnmarshalFlag parses a string specification of a plugin.
func (f *Flag) UnmarshalFlag(value string) error {
	tokens, err := shlex.Split(value, true /* posix */)
	if err != nil {
		return fmt.Errorf("invalid plugin %q: %v", value, err)
	}

	if len(tokens) < 1 {
		return fmt.Errorf("invalid plugin %q: please provide a name", value)
	}

	f.Name = tokens[0]
	exe := _pluginExecPrefix + f.Name
	path, err := exec.LookPath(exe)
	if err != nil {
		return fmt.Errorf("invalid plugin %q: could not find executable %q: %v", value, exe, err)
	}

	cmd := exec.Command(path, tokens[1:]...)
	cmd.Stderr = os.Stderr // connect stderr so that plugins can log

	f.Command = cmd
	return nil
}

// Flags is a collection of ThriftRW external plugin specifications.
type Flags []Flag

// Handle gets a MultiHandle to all the plugins in this list or nil if the
// list is empty.
//
// The returned handle MUST be closed by the caller if error was nil.
func (fs Flags) Handle() (MultiHandle, error) {
	var (
		lock  sync.Mutex
		multi MultiHandle
	)

	err := concurrent.Range(fs, func(_ int, f Flag) error {
		h, err := f.Handle()
		if err != nil {
			return err
		}

		lock.Lock()
		defer lock.Unlock()

		multi = append(multi, h)
		return nil
	})

	if err == nil {
		return multi, nil
	}

	return nil, multierr.Append(err, multi.Close())
}
