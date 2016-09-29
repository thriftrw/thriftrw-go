package multiplex

import (
	"testing"

	"go.uber.org/thriftrw/internal/envelope/envelopetest"
	"go.uber.org/thriftrw/wire"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := envelopetest.NewMockClient(mockCtrl)
	mockClient.EXPECT().Send("Foo:hello", wire.NewValueStruct(wire.Struct{})).
		Return(wire.NewValueStruct(wire.Struct{}), nil)

	client := NewClient("Foo", mockClient)
	_, err := client.Send("hello", wire.NewValueStruct(wire.Struct{}))
	assert.NoError(t, err)
}
