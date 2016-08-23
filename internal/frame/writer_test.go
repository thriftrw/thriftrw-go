package frame

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination mock_writer_test.go -package=frame io Writer,WriteCloser

func TestWriterWrite(t *testing.T) {
	tests := []struct {
		giveChunks [][]byte
		wantBody   []byte
	}{
		{
			giveChunks: [][]byte{
				{},
				{0x01},
				{0x01, 0x02},
				{0x01, 0x02, 0x03},
			},
			wantBody: []byte{
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x01, 0x01,
				0x00, 0x00, 0x00, 0x02, 0x01, 0x02,
				0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03,
			},
		},
	}

	for _, tt := range tests {
		var buff bytes.Buffer
		w := NewWriter(&buff)
		for _, chunk := range tt.giveChunks {
			assert.NoError(t, w.Write(chunk))
		}
		assert.Equal(t, tt.wantBody, buff.Bytes())
	}
}

func TestWriterWriteError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		expect  func(*MockWriter)
		wantErr error
	}{
		{
			expect: func(w *MockWriter) {
				w.EXPECT().Write([]byte{0x00, 0x00, 0x00, 0x01}).
					Return(0, errors.New("failed to write length"))
			},
			wantErr: errors.New("failed to write length"),
		},
		{
			expect: func(w *MockWriter) {
				gomock.InOrder(
					w.EXPECT().Write([]byte{0x00, 0x00, 0x00, 0x01}).Return(4, nil),
					w.EXPECT().Write([]byte{0x00}).Return(0, errors.New("great sadness")),
				)
			},
			wantErr: errors.New("great sadness"),
		},
	}

	for _, tt := range tests {
		w := NewMockWriter(mockCtrl)
		tt.expect(w)

		err := NewWriter(w).Write([]byte{0x00})
		if assert.Error(t, err) {
			assert.Equal(t, tt.wantErr, err)
		}
	}
}

func TestWriterClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		give func() io.Writer
		want error
	}{
		{
			give: func() io.Writer {
				return NewMockWriter(mockCtrl)
			},
		},
		{
			give: func() io.Writer {
				w := NewMockWriteCloser(mockCtrl)
				w.EXPECT().Close().Return(nil)
				return w
			},
		},
		{
			give: func() io.Writer {
				w := NewMockWriteCloser(mockCtrl)
				w.EXPECT().Close().Return(errors.New("great sadness"))
				return w
			},
			want: errors.New("great sadness"),
		},
	}

	for _, tt := range tests {
		err := NewWriter(tt.give()).Close()
		assert.Equal(t, tt.want, err)
	}
}
