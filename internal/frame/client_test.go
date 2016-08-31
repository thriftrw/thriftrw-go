package frame

import (
	"bytes"
	"errors"
	"testing"

	"go.uber.org/thriftrw/internal/iotest"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestClientSendWriteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	r := new(bytes.Buffer)
	w := NewMockWriter(mockCtrl)
	w.EXPECT().Write(gomock.Any()).Return(0, errors.New("great sadness"))

	client := NewClient(w, r)
	_, err := client.Send([]byte{})
	assert.Equal(t, errors.New("great sadness"), err)
}

func TestClientSendReadError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ch, r := iotest.ChunkReader()
	ch <- []byte{0x00, 0x00}
	ch <- nil

	w := NewMockWriter(mockCtrl)
	w.EXPECT().Write(gomock.Any()).Return(10, nil)

	client := NewClient(w, r)
	_, err := client.Send([]byte{})
	assert.Equal(t, iotest.ErrUser, err)
}
