package multiplex

import (
	"testing"

	"go.uber.org/thriftrw/internal/envelope"
	"go.uber.org/thriftrw/internal/envelope/envelopetest"
	"go.uber.org/thriftrw/wire"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandlerErrors(t *testing.T) {
	tests := []struct {
		name      string
		expect    func(meta, kv *envelopetest.MockHandler)
		wantError error
	}{
		{
			name: "Meta:health",
			expect: func(meta, kv *envelopetest.MockHandler) {
				meta.EXPECT().Handle("health", gomock.Any()).
					Return(wire.NewValueStruct(wire.Struct{}), nil)
			},
		},
		{
			name: "KeyValue:setValue",
			expect: func(meta, kv *envelopetest.MockHandler) {
				kv.EXPECT().Handle("setValue", gomock.Any()).
					Return(wire.NewValueStruct(wire.Struct{}), nil)
			},
		},
		{
			name:      "keyValue:setValue",
			wantError: envelope.ErrUnknownMethod("keyValue:setValue"),
		},
		{
			name:      "setValue",
			wantError: envelope.ErrUnknownMethod("setValue"),
		},
	}

	for _, tt := range tests {
		func() {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			meta := envelopetest.NewMockHandler(mockCtrl)
			kv := envelopetest.NewMockHandler(mockCtrl)

			handler := NewHandler()
			handler.Put("Meta", meta)
			handler.Put("KeyValue", kv)

			if tt.expect != nil {
				tt.expect(meta, kv)
			}

			_, err := handler.Handle(tt.name, wire.NewValueStruct(wire.Struct{}))
			assert.Equal(t, tt.wantError, err)
		}()
	}
}
