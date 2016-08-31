package envelope

import (
	"errors"
	"io"
	"testing"

	"go.uber.org/thriftrw/internal/envelope/envelopetest"
	"go.uber.org/thriftrw/internal/envelope/exception"
	"go.uber.org/thriftrw/ptr"
	"go.uber.org/thriftrw/wire"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	internalError := exception.ExceptionTypeInternalError
	unknownMethod := exception.ExceptionTypeUnknownMethod

	tests := []struct {
		desc string

		// Either transportError or decode* can be set.
		transportError error
		decodeEnvelope wire.Envelope
		decodeError    error

		wantError error // expected error if any
	}{
		{
			desc: "nothing went wrong",
			decodeEnvelope: wire.Envelope{
				Name:  "hello",
				Type:  wire.Reply,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{}),
			},
		},
		{
			desc:        "decode error",
			decodeError: errors.New("great sadness"),
			wantError:   errors.New("great sadness"),
		},
		{
			desc: "internal error",
			decodeEnvelope: wire.Envelope{
				Name:  "hello",
				Type:  wire.Exception,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("great sadness")},
					{ID: 2, Value: wire.NewValueI32(6)}, // Internal error
				}}),
			},
			wantError: &exception.TApplicationException{
				Message: ptr.String("great sadness"),
				Type:    &internalError,
			},
		},
		{
			desc: "unknown method",
			decodeEnvelope: wire.Envelope{
				Name:  "hello",
				Type:  wire.Exception,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString(`unknown method "hello"`)},
					{ID: 2, Value: wire.NewValueI32(1)}, // Internal error
				}}),
			},
			wantError: &exception.TApplicationException{
				Message: ptr.String(`unknown method "hello"`),
				Type:    &unknownMethod,
			},
		},
		{
			desc: "unknown envelope type",
			decodeEnvelope: wire.Envelope{
				Name:  "hello",
				Type:  wire.EnvelopeType(12),
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{}),
			},
			wantError: errUnknownEnvelopeType(12),
		},
		{
			desc:           "transport error",
			transportError: errors.New("great sadness"),
			wantError:      errors.New("great sadness"),
		},
	}

	for _, tt := range tests {
		proto := NewMockProtocol(mockCtrl)
		proto.EXPECT().EncodeEnveloped(
			wire.Envelope{
				Name:  "hello",
				Type:  wire.Call,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("world")},
				}}),
			},
			gomock.Any(),
		).Do(func(_ wire.Envelope, w io.Writer) {
			_, err := w.Write([]byte{1, 2, 3})
			assert.NoError(t, err, tt.desc)
		}).Return(nil)

		transport := envelopetest.NewMockTransport(mockCtrl)
		if tt.transportError != nil {
			transport.EXPECT().Send([]byte{1, 2, 3}).Return(nil, tt.transportError)
		} else {
			transport.EXPECT().Send([]byte{1, 2, 3}).Return([]byte{4, 5, 6}, nil)
			proto.EXPECT().DecodeEnveloped(gomock.Any()).
				Return(tt.decodeEnvelope, tt.decodeError)
		}

		client := NewClient(proto, transport)
		_, err := client.Send("hello", wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
			{ID: 1, Value: wire.NewValueString("world")},
		}}))
		assert.Equal(t, tt.wantError, err)
	}
}
