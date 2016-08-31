package frame

import (
	"bytes"
	"errors"
	"testing"

	"go.uber.org/thriftrw/internal/iotest"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServeReadError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ch, r := iotest.ChunkReader()
	ch <- []byte{0x00, 0x00}
	ch <- nil

	w := NewMockWriteCloser(mockCtrl)
	w.EXPECT().Close().Return(nil)

	server := NewServer(r, w)
	err := server.Serve(handlerFunc(
		func([]byte) ([]byte, error) {
			return nil, errors.New("unexpected call")
		},
	))

	assert.Equal(t, iotest.ErrUser, err)
}

func TestServeHandleError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	r := bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})

	w := NewMockWriteCloser(mockCtrl)
	w.EXPECT().Close().Return(nil)

	server := NewServer(r, w)
	err := server.Serve(handlerFunc(
		func([]byte) ([]byte, error) {
			return nil, errors.New("great sadness")
		},
	))

	assert.Equal(t, errors.New("great sadness"), err)
}

func TestServeWriteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	r := bytes.NewReader([]byte{0x00, 0x00, 0x00, 0x00})
	w := NewMockWriter(mockCtrl)

	w.EXPECT().Write(gomock.Any()).Return(0, errors.New("great sadness"))

	server := NewServer(r, w)
	err := server.Serve(handlerFunc(
		func([]byte) ([]byte, error) {
			return []byte("hello"), nil
		},
	))

	assert.Equal(t, errors.New("great sadness"), err)
}
