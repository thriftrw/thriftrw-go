// Copyright (c) 2019 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package envelope

import (
	"bytes"
	"fmt"

	"go.uber.org/thriftrw/internal/envelope/exception"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/wire"
)

// Transport sends and receive binary payloads.
type Transport interface {
	Send([]byte) ([]byte, error)
}

// Client sends Thrift requests and returns their responses.
//
// Client is thread-safe only if the underlying transport is thread-safe.
type Client interface {
	// Send sends a request to the method with the given name and body.
	Send(name string, body wire.Value) (wire.Value, error)
}

// NewClient builds a new client which sends requests over the given
// transport, encoding them using the given protocol.
func NewClient(p protocol.Protocol, t Transport) Client {
	return client{p: p, t: t}
}

type client struct {
	p protocol.Protocol
	t Transport
}

// Send sends the given request envelope over this transport.
func (c client) Send(name string, reqValue wire.Value) (wire.Value, error) {
	reqEnvelope := wire.Envelope{
		Name:  name,
		Type:  wire.Call,
		SeqID: 1, // don't care
		Value: reqValue,
	}

	// TODO(abg): We don't use or support out-of-order requests and responses
	// for plugin communaction so this should be fine for now but we may
	// eventually want to match responses to requests using seqID.

	var buff bytes.Buffer
	if err := c.p.EncodeEnveloped(reqEnvelope, &buff); err != nil {
		return wire.Value{}, err
	}

	resBody, err := c.t.Send(buff.Bytes())
	if err != nil {
		return wire.Value{}, err
	}

	resEnvelope, err := c.p.DecodeEnveloped(bytes.NewReader(resBody))
	if err != nil {
		return wire.Value{}, err
	}

	switch resEnvelope.Type {
	case wire.Exception:
		var exc exception.TApplicationException
		if err := exc.FromWire(resEnvelope.Value); err != nil {
			return wire.Value{}, err
		}
		return wire.Value{}, &exc

	case wire.Reply:
		return resEnvelope.Value, nil

	default:
		return wire.Value{}, errUnknownEnvelopeType(resEnvelope.Type)
	}
}

type errUnknownEnvelopeType wire.EnvelopeType

func (e errUnknownEnvelopeType) Error() string {
	return fmt.Sprintf("unknown envelope type: expected Reply, got %v", wire.EnvelopeType(e))
}
