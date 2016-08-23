package frame

import (
	"errors"
	"io"
	"testing"
	"testing/quick"

	"github.com/stretchr/testify/assert"
)

func TestClientServer(t *testing.T) {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()

	defer func() {
		assert.NoError(t, serverWriter.Close())
		assert.NoError(t, clientWriter.Close())
		assert.NoError(t, clientReader.Close())
		assert.NoError(t, serverReader.Close())
	}()

	server := NewServer(serverReader, serverWriter)
	client := NewClient(clientWriter, clientReader)

	type entry struct{ receive, reply []byte }
	expect := make(chan entry, 100)

	go func() {
		err := server.Serve(handlerFunc(
			func(got []byte) ([]byte, error) {
				entry := <-expect
				if assert.Equal(t, entry.receive, got) {
					return entry.reply, nil
				}
				return nil, errors.New("unexpected request")
			},
		))
		assert.NoError(t, err)
	}()
	defer server.Stop()

	err := quick.Check(func(send, get []byte) bool {
		expect <- entry{receive: send, reply: get}
		got, err := client.Send(send)
		return assert.NoError(t, err) && assert.Equal(t, get, got)
	}, nil)
	assert.NoError(t, err)
}

func TestClientServerHandleError(t *testing.T) {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()

	defer func() {
		assert.NoError(t, serverWriter.Close())
		assert.NoError(t, clientWriter.Close())
		assert.NoError(t, clientReader.Close())
		assert.NoError(t, serverReader.Close())
	}()

	server := NewServer(serverReader, serverWriter)
	client := NewClient(clientWriter, clientReader)

	go func() {
		err := server.Serve(handlerFunc(
			func([]byte) ([]byte, error) {
				return nil, errors.New("great sadness")
			},
		))
		assert.Error(t, err)
	}()

	_, err := client.Send([]byte("hello"))
	assert.Equal(t, io.EOF, err)
}

type handlerFunc func([]byte) ([]byte, error)

func (f handlerFunc) Handle(b []byte) ([]byte, error) {
	return f(b)
}
