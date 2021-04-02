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
	"fmt"
	"testing"

	// Use . import so the generated Thrift code can import envelope
	// without causing a circular dependency.
	. "go.uber.org/thriftrw/envelope"

	"go.uber.org/thriftrw/protocol"
	"go.uber.org/thriftrw/wire"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEnveloper struct {
	Name  string
	Type  wire.EnvelopeType
	Value wire.Value
	Err   error
}

var _ Enveloper = fakeEnveloper{}

func (e fakeEnveloper) MethodName() string { return e.Name }

func (e fakeEnveloper) EnvelopeType() wire.EnvelopeType { return e.Type }

func (e fakeEnveloper) ToWire() (wire.Value, error) {
	return e.Value, e.Err
}

func TestWrite(t *testing.T) {
	// This enveloper represents the arguments of a method,
	//   getValue(1: string key)
	enveloper := fakeEnveloper{
		Name: "getValue",
		Type: wire.Call,
		Value: wire.NewValueStruct(wire.Struct{
			Fields: []wire.Field{
				{
					ID:    1,
					Value: wire.NewValueString("foo"),
				},
			},
		}),
	}

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var buff bytes.Buffer
		require.NoError(t, Write(protocol.Binary, &buff, 1234, enveloper))
		assert.Equal(t,
			[]byte{
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
			}, buff.Bytes())
	})

	t.Run("failure", func(t *testing.T) {
		errEnveloper := enveloper
		errEnveloper.Err = fmt.Errorf("great sadness")

		var buff bytes.Buffer
		require.Error(t, Write(protocol.Binary, &buff, 1234, errEnveloper))
	})
}

func TestReadReply(t *testing.T) {
	tests := []struct {
		desc string
		bs   []byte

		wantValue   wire.Value
		wantNoValue bool //

		wantSeqID int32
		wantErr   string
	}{
		{
			desc:        "Invalid envelope",
			bs:          []byte{0},
			wantNoValue: true,
			wantErr:     "unexpected EOF",
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
			wantValue: wire.NewValueStruct(wire.Struct{}),
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
			wantValue: wire.NewValueStruct(wire.Struct{}),
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
			wantSeqID:   1234,
			wantNoValue: true,
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
			wantValue: wire.NewValueStruct(wire.Struct{
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
		t.Run(tt.desc, func(t *testing.T) {
			result, seqID, err := ReadReply(protocol.Binary, bytes.NewReader(tt.bs))
			if tt.wantErr != "" {
				if assert.Error(t, err, tt.desc) {
					assert.Contains(t, err.Error(), tt.wantErr, "%v: error mismatch", tt.desc)
				}
			} else {
				assert.NoError(t, err, tt.desc)
			}
			if !tt.wantNoValue {
				assert.True(t, wire.ValuesAreEqual(tt.wantValue, result), "%v: result mismatch", tt.desc)
			}
			assert.Equal(t, tt.wantSeqID, seqID, "%v: seqID mismatch", tt.desc)
		})
	}
}
