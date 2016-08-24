package process

import (
	"bytes"
	"os/exec"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCat(t *testing.T) {
	client, err := NewClient(exec.Command("cat"))
	require.NoError(t, err)
	defer client.Close()

	err = quick.Check(func(give []byte) bool {
		got, err := client.Send(give)
		return assert.NoError(t, err) && assert.Equal(t, got, give)
	}, nil)
	assert.NoError(t, err)
}

func TestSendAfterStop(t *testing.T) {
	client, err := NewClient(exec.Command("cat"))
	require.NoError(t, err)

	_, err = client.Send([]byte("hello"))
	require.NoError(t, err)
	require.NoError(t, client.Close())

	assert.Panics(t, func() {
		client.Send([]byte("hello"))
	})
}

func TestCloseTwice(t *testing.T) {
	client, err := NewClient(exec.Command("cat"))
	require.NoError(t, err)
	require.NoError(t, client.Close())
	require.NoError(t, client.Close())
}

func TestStartErrors(t *testing.T) {
	tests := []struct {
		getCommand      func() *exec.Cmd
		wantMessageLike string
	}{
		{
			getCommand: func() *exec.Cmd {
				return exec.Command("this_command_does_not_exist")
			},
			wantMessageLike: `failed to start "this_command_does_not_exist":`,
		},
		{
			getCommand: func() *exec.Cmd {
				// StdoutPipe will fail if Stdout is already set
				cmd := exec.Command("/bin/cat")
				cmd.Stdout = new(bytes.Buffer)
				return cmd
			},
			wantMessageLike: `failed to create stdout pipe to "/bin/cat":`,
		},
		{
			getCommand: func() *exec.Cmd {
				// StdinPipe will fail if Stdout is already set
				cmd := exec.Command("/bin/cat")
				cmd.Stdin = new(bytes.Buffer)
				return cmd
			},
			wantMessageLike: `failed to create stdin pipe to "/bin/cat":`,
		},
	}

	for _, tt := range tests {
		_, err := NewClient(tt.getCommand())
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), tt.wantMessageLike)
		}
	}
}
