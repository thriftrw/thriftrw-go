package envelope

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"go.uber.org/thriftrw/wire"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tests := []struct {
		desc    string
		handler handlerFunc

		// Only one of the following must be set
		giveEnvelope    wire.Envelope
		giveDecodeError error

		wantEnvelope *wire.Envelope // envelope to expect as output
		wantError    error          // whether Handle should fail
	}{
		{
			desc: "nothing went wrong",
			handler: func(name string, body wire.Value) (wire.Value, error) {
				assert.Equal(t, "write", name)
				return wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("world")},
				}}), nil
			},
			giveEnvelope: wire.Envelope{
				Name:  "write",
				Type:  wire.Call,
				SeqID: 42,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("hello")},
				}}),
			},
			wantEnvelope: &wire.Envelope{
				Name:  "write",
				Type:  wire.Reply,
				SeqID: 42,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("world")},
				}}),
			},
		},
		{
			desc:            "decode failed",
			giveDecodeError: errors.New("great sadness"),
			wantError:       errors.New("great sadness"),
		},
		{
			desc: "unhandled internal error",
			giveEnvelope: wire.Envelope{
				Name:  "hello",
				Type:  wire.Call,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{}),
			},
			handler: func(name string, body wire.Value) (wire.Value, error) {
				assert.Equal(t, "hello", name)
				return wire.Value{}, errors.New("great sadness")
			},
			wantEnvelope: &wire.Envelope{
				Name:  "hello",
				Type:  wire.Exception,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("great sadness")},
					{ID: 2, Value: wire.NewValueI32(6)}, // Internal error
				}}),
			},
		},
		{
			desc: "unknown method",
			giveEnvelope: wire.Envelope{
				Name:  "hello",
				Type:  wire.Call,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{}),
			},
			handler: func(string, wire.Value) (wire.Value, error) {
				return wire.Value{}, ErrUnknownMethod("hello")
			},
			wantEnvelope: &wire.Envelope{
				Name:  "hello",
				Type:  wire.Exception,
				SeqID: 1,
				Value: wire.NewValueStruct(wire.Struct{Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString(`unknown method "hello"`)},
					{ID: 2, Value: wire.NewValueI32(1)}, // Internal error
				}}),
			},
		},
	}

	for _, tt := range tests {
		proto := NewMockProtocol(mockCtrl)
		if tt.giveDecodeError != nil {
			proto.EXPECT().DecodeEnveloped(gomock.Any()).
				Return(wire.Envelope{}, tt.giveDecodeError)
		} else {
			proto.EXPECT().DecodeEnveloped(gomock.Any()).
				Return(tt.giveEnvelope, nil)
		}

		if tt.wantEnvelope != nil {
			proto.EXPECT().EncodeEnveloped(matchesEnvelope(*tt.wantEnvelope), gomock.Any()).
				Do(func(_ wire.Envelope, w io.Writer) {
					_, err := w.Write([]byte{1, 2, 3})
					assert.NoError(t, err, tt.desc)
				}).Return(nil)
		}

		server := NewServer(proto, tt.handler)
		_, err := server.Handle([]byte{1, 2, 3})
		if tt.wantError != nil {
			assert.Equal(t, tt.wantError, err, tt.desc)
		}
	}
}

type handlerFunc func(string, wire.Value) (wire.Value, error)

func (f handlerFunc) Handle(name string, body wire.Value) (wire.Value, error) {
	return f(name, body)
}

type matchesEnvelope wire.Envelope

func (m matchesEnvelope) Matches(x interface{}) bool {
	got, ok := x.(wire.Envelope)
	if !ok {
		return false
	}

	want := wire.Envelope(m)
	return want.Name == got.Name &&
		want.SeqID == got.SeqID &&
		want.Type == got.Type &&
		wire.ValuesAreEqual(want.Value, got.Value)
}

// String describes what the matcher matches.
func (m matchesEnvelope) String() string {
	return fmt.Sprintf("matches %v", wire.Envelope(m))
}
