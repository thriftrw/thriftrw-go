// Copyright (c) 2016 Uber Technologies, Inc.
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

package envelope_test

import (
	"bytes"
	"errors"
	"testing"

	// Use . import so the generated Thrift code can import envelope
	// without causing a circular dependency.
	. "go.uber.org/thriftrw/envelope"

	tv "go.uber.org/thriftrw/gen/testdata/services"
	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
)

type failToWire struct {
	// Embed an Enveloper so we only need to implement the methods we expect.
	Enveloper
}

func (failToWire) ToWire() (wire.Value, error) {
	return wire.Value{}, errors.New("failed")
}

func stringp(s string) *string {
	return &s
}

func TestWrite(t *testing.T) {
	tests := []struct {
		e       Enveloper
		want    []byte
		wantErr bool
	}{
		{
			e:       failToWire{},
			wantErr: true,
		},
		{
			e: tv.KeyValue_GetValue_Helper.Args((*tv.Key)(stringp("foo"))),
			want: []byte{
				0x80, 0x01, 0x00, 0x01, // version|type:4 = 1 | call
				0x00, 0x00, 0x00, 0x08, // name length = 8
				'g', 'e', 't', 'V', 'a', 'l', 'u', 'e', // "getValue"
				0x00, 0x00, 0x04, 0xd2, // seqID:4 = 1234

				// <struct>
				0x0b,       // type:1 = string
				0x00, 0x01, // id:2 = 1
				0x00, 0x00, 0x00, 0x03, // length = 3
				'f', 'o', 'o', // "foo"
				0x00, // stop
			},
		},
	}

	for _, tt := range tests {
		buf := &bytes.Buffer{}
		err := Write(protocol.Binary, buf, 1234, tt.e)
		if tt.wantErr {
			assert.Error(t, err, "%T should fail to be enveloped", tt.e)
			continue
		}
		if !assert.NoError(t, err, "%T should not fail to be enveloped", tt.e) {
			continue
		}

		assert.Equal(t, tt.want, buf.Bytes(), "%T enveloped mismatch")
	}
}

func TestReadReply(t *testing.T) {
	tests := []struct {
		desc      string
		bs        []byte
		want      wire.Value
		wantSeqID int32
		wantErr   string
	}{
		{
			desc:    "Invalid envelope",
			bs:      []byte{0},
			wantErr: "unexpected EOF",
		},
		{
			desc: "Unexpected envelope type",
			bs: []byte{
				0x80, 0x01, 0x00, 0x01, // version|type:4 = 1 | call
				0x00, 0x00, 0x00, 0x03, // name length = 3
				'a', 'b', 'c', // "abc"
				0x00, 0x00, 0x04, 0xd2, // seqID:4 = 1234

				// <struct>
				0x00, // stop
			},
			want:      wire.NewValueStruct(wire.Struct{}),
			wantSeqID: 1234,
			wantErr:   "unknown envelope",
		},
		{
			desc: "Valid reply",
			bs: []byte{
				0x80, 0x01, 0x00, 0x02, // version|type:4 = 2 | reply
				0x00, 0x00, 0x00, 0x03, // name length = 3
				'a', 'b', 'c', // "abc"
				0x00, 0x00, 0x04, 0xd2, // seqID:4 = 1234

				// <struct>
				0x00, // stop
			},
			want:      wire.NewValueStruct(wire.Struct{}),
			wantSeqID: 1234,
		},
		{
			desc: "Invalid exception",
			bs: []byte{
				0x80, 0x01, 0x00, 0x03, // version|type:4 = 3 | exception
				0x00, 0x00, 0x00, 0x03, // name length = 3
				'a', 'b', 'c', // "abc"
				0x00, 0x00, 0x04, 0xd2, // seqID:4 = 1234

				// <struct> (invalid)
				0x08,       // type:1 = i32
				0x00, 0x01, // id:2 = 1
				0x00, 0x00, 0x00, 0x01, // value = 1
				0x00, // stop
			},
			want: wire.NewValueStruct(wire.Struct{
				Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueI32(1)},
				},
			}),
			wantSeqID: 1234,
			// TODO: This should probably fail to decode. Right now, it's being ignored.
			// wantErr:   "failed to decode exception",
			wantErr: "TApplicationException{}",
		},
		{
			desc: "Valid exception",
			bs: []byte{
				0x80, 0x01, 0x00, 0x03, // version|type:4 = 3 | exception
				0x00, 0x00, 0x00, 0x03, // name length = 3
				'a', 'b', 'c', // "abc"
				0x00, 0x00, 0x04, 0xd2, // seqID:4 = 1234

				// <struct>
				0x0b,       // type:1 = string
				0x00, 0x01, // id:2 = 1
				0x00, 0x00, 0x00, 0x06, // length = 3
				'e', 'r', 'r', 'M', 's', 'g', // "errMsg"
				0x08,       // type:1 = i32
				0x00, 0x02, // id:2 = 2
				0x00, 0x00, 0x00, 0x01, // value = 1 (unknown method)
				0x00, // stop
			},
			want: wire.NewValueStruct(wire.Struct{
				Fields: []wire.Field{
					{ID: 1, Value: wire.NewValueString("errMsg")},
					{ID: 2, Value: wire.NewValueI32(1)},
				},
			}),
			wantSeqID: 1234,
			wantErr:   "TApplicationException{Message: errMsg, Type: UNKNOWN_METHOD}",
		},
	}

	for _, tt := range tests {
		result, seqID, err := ReadReply(protocol.Binary, bytes.NewReader(tt.bs))
		if tt.wantErr != "" {
			if assert.Error(t, err, tt.desc) {
				assert.Contains(t, err.Error(), tt.wantErr, "%v: error mismatch", tt.desc)
			}
		} else {
			assert.NoError(t, err, tt.desc)
		}
		assert.Equal(t, tt.want, result, "%v: result mismatch", tt.desc)
		assert.Equal(t, tt.wantSeqID, seqID, "%v: seqID mismatch", tt.desc)
	}
}
