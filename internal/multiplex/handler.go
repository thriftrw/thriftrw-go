package multiplex

import (
	"strings"

	"github.com/thriftrw/thriftrw-go/internal/envelope"
	"github.com/thriftrw/thriftrw-go/wire"
)

// Handler implements a service multiplexer
type Handler struct {
	services map[string]envelope.Handler
}

// NewHandler builds a new handler.
func NewHandler() Handler {
	return Handler{services: make(map[string]envelope.Handler)}
}

// Put adds the given service to the multiplexer.
func (h Handler) Put(name string, service envelope.Handler) {
	h.services[name] = service
}

// Handle handles the given request, dispatching to one of the
// registered services.
func (h Handler) Handle(name string, req wire.Value) (wire.Value, error) {
	parts := strings.SplitN(name, ":", 2)
	if len(parts) < 2 {
		return wire.Value{}, envelope.ErrUnknownMethod(name)
	}

	service, ok := h.services[parts[0]]
	if !ok {
		return wire.Value{}, envelope.ErrUnknownMethod(name)
	}

	name = parts[1]
	return service.Handle(name, req)
}
