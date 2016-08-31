package frame

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"go.uber.org/thriftrw/internal/iotest"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -destination mock_reader_test.go -package=frame io Reader,ReadCloser

// changes the fast path threshold and returns a function that may be called to
// reset the old value.
func setFastPathThreshold(newSize int64) func() {
	oldSize := _fastPathFrameSize
	_fastPathFrameSize = newSize
	return func() {
		_fastPathFrameSize = oldSize
	}
}

func TestReader(t *testing.T) {
	type wantRead struct {
		frame []byte
		err   error
	}

	tests := []struct {
		desc       string
		giveReader func() io.Reader
		wantReads  []wantRead // reads to perform and what to expect

		// if non-zero, the _fastPathFrameSize will be set to this value for the test
		fastPathThreshold int64
	}{
		{
			desc: "error while reading length",
			giveReader: func() io.Reader {
				ch, reader := iotest.ChunkReader()

				// one successful read
				ch <- []byte{0x00, 0x00, 0x00, 0x01, 0x01}

				// error half way through reading the length
				ch <- []byte{0x00, 0x00}
				ch <- nil

				return reader
			},
			wantReads: []wantRead{
				{frame: []byte{0x01}},
				{err: iotest.ErrUser},
			},
		},
		{
			desc: "fast path, no errors",
			giveReader: func() io.Reader {
				return bytes.NewReader([]byte{
					0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x01, 0x01,
					0x00, 0x00, 0x00, 0x02, 0x01, 0x02,
					0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03,
				})
			},
			wantReads: []wantRead{
				{frame: []byte{}},
				{frame: []byte{0x01}},
				{frame: []byte{0x01, 0x02}},
				{frame: []byte{0x01, 0x02, 0x03}},
			},
		},
		{
			desc: "fast path, error while reading body",
			giveReader: func() io.Reader {
				ch, reader := iotest.ChunkReader()

				// one successful read
				ch <- []byte{0x00, 0x00, 0x00, 0x01, 0x01}

				// error while reading the body
				ch <- []byte{0x00, 0x00, 0x00, 0x10}
				ch <- []byte{0x01, 0x02, 0x03, 0x04, 0x05}
				ch <- nil

				return reader
			},
			wantReads: []wantRead{
				{frame: []byte{0x01}},
				{err: iotest.ErrUser},
			},
		},
		{
			desc: "slow path, no errors",
			giveReader: func() io.Reader {
				return bytes.NewReader([]byte{
					0x00, 0x00, 0x00, 0x03, 0x01, 0x02, 0x03,
					0x00, 0x00, 0x00, 0x04, 0x01, 0x02, 0x03, 0x04,
					0x00, 0x00, 0x00, 0x05, 0x01, 0x02, 0x03, 0x04, 0x05,
				})
			},
			wantReads: []wantRead{
				{frame: []byte{0x01, 0x02, 0x03}},
				{frame: []byte{0x01, 0x02, 0x03, 0x04}},
				{frame: []byte{0x01, 0x02, 0x03, 0x04, 0x05}},
			},
			fastPathThreshold: 3,
		},
		{
			desc: "slow path, error while reading body",
			giveReader: func() io.Reader {
				ch, reader := iotest.ChunkReader()
				ch <- []byte{0x1f, 0x40, 0x00, 0x00} // 500 MB
				ch <- []byte{0x00, 0x00, 0x00, 0x00}
				ch <- nil
				return reader
			},
			wantReads: []wantRead{
				{err: iotest.ErrUser},
			},
		},
		{
			desc: "slow path, body too short",
			giveReader: func() io.Reader {
				ch, reader := iotest.ChunkReader()
				ch <- []byte{0x1f, 0x40, 0x00, 0x00} // 500 MB
				ch <- []byte{0x00}
				close(ch)
				return reader
			},
			wantReads: []wantRead{
				{err: io.EOF},
			},
		},
	}

	for _, tt := range tests {
		func() {
			if tt.fastPathThreshold != 0 {
				reset := setFastPathThreshold(tt.fastPathThreshold)
				defer reset()
			}

			r := NewReader(tt.giveReader())

			for _, want := range tt.wantReads {
				frame, err := r.Read()
				if want.frame != nil && assert.NoError(t, err, tt.desc) {
					assert.Equal(t, want.frame, frame, tt.desc)
				} else if assert.Error(t, err, tt.desc) {
					assert.Equal(t, want.err, err, tt.desc)
				}
			}

			assert.NoError(t, r.Close(), tt.desc)
		}()
	}
}

func TestReaderClose(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		give func() io.Reader
		want error
	}{
		{
			give: func() io.Reader {
				return NewMockReader(mockCtrl)
			},
		},
		{
			give: func() io.Reader {
				w := NewMockReadCloser(mockCtrl)
				w.EXPECT().Close().Return(nil)
				return w
			},
		},
		{
			give: func() io.Reader {
				w := NewMockReadCloser(mockCtrl)
				w.EXPECT().Close().Return(errors.New("great sadness"))
				return w
			},
			want: errors.New("great sadness"),
		},
	}

	for _, tt := range tests {
		err := NewReader(tt.give()).Close()
		assert.Equal(t, tt.want, err)
	}
}
