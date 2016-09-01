package plugin

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testdata(t *testing.T, paths ...string) string {
	cwd, err := os.Getwd()
	require.NoError(t, err, "could not determine CWD")

	args := []string{cwd, "testdata"}
	args = append(args, paths...)
	return filepath.Join(args...)
}

func prependToPath(p string) func() {
	oldPath := os.Getenv("PATH")
	newPath := p + string(os.PathListSeparator) + oldPath

	os.Setenv("PATH", newPath)
	return func() { os.Setenv("PATH", oldPath) }
}

func TestUnmarshalFlag(t *testing.T) {
	done := prependToPath(testdata(t))
	defer done()

	tests := []struct {
		giveValue string

		wantName string
		wantPath string
		wantArgs []string

		wantErrorPrefix string
	}{
		{
			giveValue:       "",
			wantErrorPrefix: `invalid plugin "": please provide a name`,
		},
		{
			giveValue:       `foo '`,
			wantErrorPrefix: `invalid plugin "foo '": No closing quotation`,
		},
		{
			giveValue:       "unknown-plugin",
			wantErrorPrefix: `invalid plugin "unknown-plugin": could not find executable "thriftrw-plugin-unknown-plugin":`,
		},
		{
			giveValue: "empty",
			wantName:  "empty",
			wantPath:  testdata(t, "thriftrw-plugin-empty"),
		},
		{
			giveValue: "handshake-failed",
			wantName:  "handshake-failed",
			wantPath:  testdata(t, "thriftrw-plugin-handshake-failed"),
		},
		{
			giveValue: "handshake-failed --bar baz",
			wantName:  "handshake-failed",
			wantPath:  testdata(t, "thriftrw-plugin-handshake-failed"),
			wantArgs:  []string{"--bar", "baz"},
		},
		{
			giveValue: `handshake-failed bar\ baz`,
			wantName:  "handshake-failed",
			wantPath:  testdata(t, "thriftrw-plugin-handshake-failed"),
			wantArgs:  []string{"bar baz"},
		},
		{
			giveValue: `handshake-failed bar="baz qux"`,
			wantName:  "handshake-failed",
			wantPath:  testdata(t, "thriftrw-plugin-handshake-failed"),
			wantArgs:  []string{"bar=baz qux"},
		},
	}

	for _, tt := range tests {
		var f Flag
		err := f.UnmarshalFlag(tt.giveValue)

		if tt.wantErrorPrefix != "" {
			if assert.Error(t, err, "expected error for %q", tt.giveValue) {
				assert.True(t,
					strings.HasPrefix(err.Error(), tt.wantErrorPrefix),
					"expected prefix %q on error message %q for %q",
					tt.wantErrorPrefix, err.Error(), tt.giveValue,
				)
			}
			continue
		}

		if !assert.NoError(t, err, "expected no error for %q", tt.giveValue) {
			continue
		}

		assert.Equal(t, tt.wantName, f.Name, "name for %q does not match", tt.giveValue)
		assert.Equal(t, tt.wantPath, f.Command.Path, "path for %q does not match", tt.giveValue)

		// We check args[1:] because args[0] is always path.
		if len(tt.wantArgs) == 0 {
			assert.Empty(t, f.Command.Args[1:], "args for %q do not match", tt.giveValue)
		} else {
			assert.Equal(t, tt.wantArgs, f.Command.Args[1:], "args for %q do not match", tt.giveValue)
		}
	}
}

func TestFlagHandle(t *testing.T) {
	tests := []struct {
		desc string
		name string
		path string
		args []string

		wantErrors []string
	}{
		{
			desc: "working plugin",
			name: "empty",
			path: testdata(t, "thriftrw-plugin-empty"),
		},
		{
			desc: "handshake failed",
			name: "handshake-failed",
			path: testdata(t, "thriftrw-plugin-handshake-failed"),
			wantErrors: []string{
				`failed to open plugin "handshake-failed":`,
				`handshake with plugin "handshake-failed" failed:`,
			},
		},
		{
			desc:       "non existent path",
			name:       "bar",
			path:       testdata(t, "thriftrw-plugin-bar"),
			wantErrors: []string{`failed to open plugin "bar":`, "no such file or directory"},
		},
	}

	for _, tt := range tests {
		f := Flag{
			Name:    tt.name,
			Command: exec.Command(tt.path, tt.args...),
		}
		h, err := f.Handle()

		if len(tt.wantErrors) > 0 {
			if !assert.Error(t, err, "%v: expected error", tt.desc) {
				continue
			}

			for _, msg := range tt.wantErrors {
				assert.Contains(t, err.Error(), msg, "%v: error message mismatch", tt.desc)
			}
		} else {
			if assert.NoError(t, err, "%v: expected no error", tt.desc) {
				assert.NoError(t, h.Close(), "%v: failed to close", tt.desc)
			}
		}
	}
}

func TestFlagsHandle(t *testing.T) {
	type plug struct {
		name string
		path string
		args []string
	}

	tests := []struct {
		desc  string
		plugs []plug

		wantErrors []string
	}{
		{
			desc: "no plugins",
		},
		{
			desc:  "empty plugins",
			plugs: []plug{},
		},
		{
			desc: "all success",
			plugs: []plug{
				{
					name: "empty",
					path: testdata(t, "thriftrw-plugin-empty"),
				},
				{
					name: "another-empty",
					path: testdata(t, "thriftrw-plugin-another-empty"),
				},
			},
		},
		{
			desc: "all fail",
			plugs: []plug{
				{
					name: "handshake-failed",
					path: testdata(t, "thriftrw-plugin-handshake-failed"),
				},
				{
					name: "another-handshake-failed",
					path: testdata(t, "thriftrw-plugin-another-handshake-failed"),
				},
			},
			wantErrors: []string{
				`failed to open plugin "handshake-failed": handshake with plugin "handshake-failed" failed:`,
				`failed to open plugin "another-handshake-failed": handshake with plugin "another-handshake-failed" failed:`,
			},
		},
		{
			desc: "partial failure",
			plugs: []plug{
				{
					name: "handshake-failed",
					path: testdata(t, "thriftrw-plugin-handshake-failed"),
				},
				{
					name: "empty",
					path: testdata(t, "thriftrw-plugin-empty"),
				},
			},
			wantErrors: []string{
				`failed to open plugin "handshake-failed": handshake with plugin "handshake-failed" failed:`,
			},
		},
	}

	for _, tt := range tests {
		var flags Flags
		for _, p := range tt.plugs {
			flags = append(flags, Flag{
				Name:    p.name,
				Command: exec.Command(p.path, p.args...),
			})
		}

		h, err := flags.Handle()

		if len(tt.wantErrors) > 0 {
			if !assert.Error(t, err, "%v: expected error", tt.desc) {
				continue
			}

			for _, msg := range tt.wantErrors {
				assert.Contains(t, err.Error(), msg, "%v: error message mismatch", tt.desc)
			}

			continue
		}

		if assert.NoError(t, err, "%v: expected no error", tt.desc) {
			assert.NoError(t, h.Close(), "%v: failed to close", tt.desc)
		}
	}
}
