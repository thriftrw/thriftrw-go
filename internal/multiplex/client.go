package multiplex

import (
	"github.com/thriftrw/thriftrw-go/internal/envelope"
	"github.com/thriftrw/thriftrw-go/wire"
)

// NewClient builds a new multiplexing client.
//
// name is the name of the service for which requests are being sent through
// this client.
func NewClient(name string, c envelope.Client) envelope.Client {
	return client{name: name, c: c}
}

// client is a multiplexing client.
type client struct {
	c    envelope.Client
	name string
}

func (c client) Send(name string, reqValue wire.Value) (wire.Value, error) {
	return c.c.Send(c.name+":"+name, reqValue)
}
