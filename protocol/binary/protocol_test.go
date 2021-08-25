package binary_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/protocol/binary"
	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

func TestReadRequestSeekMiddle(t *testing.T) {
	// ReadRequest looks ahead into the buffer and seeks back if it sees an
	// envelope. This test ensures that when we seek back, we seek back to
	// the start of the envelope, and not to the start of the entire file.

	buff := append(
		[]byte("hello"),
		encodeEnvelope(t, "Cache::clear", wire.Call, 42)...,
	)

	r := bytes.NewReader(buff)

	// Move to somewhere near the middle of the buffer so that when
	// ReadRequest seeks back, it has to come back to this exact position.
	hello := make([]byte, 5)
	_, err := r.Read(hello)
	require.NoError(t, err)
	require.Equal(t, "hello", string(hello))

	_, err = binary.Default.ReadRequest(
		context.Background(),
		wire.Call,
		r,
		new(emptyStructReader),
	)
	require.NoError(t, err)
}

// Serializes an envelope to bytes.
func encodeEnvelope(t testing.TB,
	name string,
	typ wire.EnvelopeType,
	seqID int32,
	fields ...wire.Field,
) []byte {
	var buff bytes.Buffer
	err := binary.Default.EncodeEnveloped(wire.Envelope{
		Name:  name,
		Type:  typ,
		SeqID: seqID,
		Value: wire.NewValueStruct(wire.Struct{
			Fields: fields,
		}),
	}, &buff)
	require.NoError(t, err)
	return buff.Bytes()
}

// Measures the cost of looking ahead in Binary.ReadRequest to determine if the
// request is enveloped.
func BenchmarkReadRequestLookahead(b *testing.B) {
	// Build an enveloped payload with an empty struct to only measure the
	// cost of looking ahead.
	bs := encodeEnvelope(b, "Cache::clear", wire.Call, 42) // empty struct

	// bufferedReader is an interface implemented by bytes.Reader. It lets
	// us use the same benchmarking logic for readers that support seek and
	// those that do not.
	type bufferReader interface {
		io.Reader

		// Reset resets the position of this buffered reader,
		// specifying that it should read from the provided buffer.
		Reset([]byte)
	}

	benchmarks := []struct {
		desc   string
		reader bufferReader
	}{
		{
			// bytes.Reader supports seek This will use the
			// optimized path.
			desc:   "with seek",
			reader: bytes.NewReader(nil),
		},
		{
			// sliceReader does not support seek. This will
			// construct a multi-reader with the buffered text.
			desc:   "no seek",
			reader: new(sliceReader),
		},
	}

	for _, bb := range benchmarks {
		b.Run(bb.desc, func(b *testing.B) {
			var body emptyStructReader

			ctx := context.Background()
			for i := 0; i < b.N; i++ {
				bb.reader.Reset(bs)

				_, err := binary.Default.ReadRequest(
					ctx,
					wire.Call,
					bb.reader,
					&body,
				)
				require.NoError(b, err)
			}
		})
	}
}

type emptyStructReader struct{}

func (*emptyStructReader) Decode(r stream.Reader) error {
	if err := r.ReadStructBegin(); err != nil {
		return err
	}
	return r.ReadStructEnd()
}

// sliceReader is an io.Reader that reads from a slice incrementally. It does
// not implement io.Seeker.
type sliceReader struct {
	bs  []byte
	idx int
}

// Reset resets the position of the slice reader and specifies that it should
// read from the given slice.
func (r *sliceReader) Reset(bs []byte) {
	r.bs = bs
	r.idx = 0
}

func (r *sliceReader) Read(bs []byte) (int, error) {
	n := copy(bs, r.bs[r.idx:])
	r.idx += n
	return n, nil
}
