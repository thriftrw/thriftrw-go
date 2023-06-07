// Copyright (c) 2015 Uber Technologies, Inc.
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

package protocol

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/protocol/binary"
	"go.uber.org/thriftrw/protocol/stream"
	"go.uber.org/thriftrw/wire"
)

var (
	_testNonStrictEnvelopeOneWayBytes = []byte{
		// envelope
		0x00, 0x00, 0x00, 0x05, // length:4 = 5
		0x77, 0x72, 0x69, 0x74, 0x65, // 'write'
		0x04,                   // type:1 = OneWay
		0x00, 0x00, 0x00, 0x2a, // seqid:4 = 42
	}

	_testNonStrictEnvelopeExceptionBytes = []byte{
		// envelope
		0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
		0x03,                   // type:1 = Exception
		0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436
	}

	_testStrictEnvelopeCallBytes = []byte{
		// envelope
		0x80, 0x01, 0x00, 0x01, // version|type:4 = 1 | call
		0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
		0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436
	}
)

type testResponse struct {
	t *testing.T
}

var _ stream.Enveloper = (*testResponse)(nil)

func (h testResponse) MethodName() string              { return "" }
func (h testResponse) EnvelopeType() wire.EnvelopeType { return wire.Reply }
func (h testResponse) Encode(w stream.Writer) error {
	h.t.Helper()

	require.NoError(h.t, w.WriteStructBegin(), "failed to write response struct begin")
	require.NoError(h.t, w.WriteStructEnd(), "failed to write response struct end")
	return nil
}

type testRequestBody struct {
	t  *testing.T
	fh *stream.FieldHeader
}

var _ stream.BodyReader = (*testRequestBody)(nil)

func (h *testRequestBody) Decode(sr stream.Reader) error {
	h.t.Helper()

	require.NoError(h.t, sr.ReadStructBegin(), "failed to read struct begin")
	tempFh, ok, err := sr.ReadFieldBegin()
	require.NoError(h.t, err, "failed to read field begin")
	require.True(h.t, ok, "failed to read field begin")
	require.NoError(h.t, sr.Skip(tempFh.Type), "failed to skip field type")
	require.NoError(h.t, sr.ReadFieldEnd(), "failed to read field end")
	require.NoError(h.t, sr.ReadStructEnd(), "failed to read struct end")
	h.fh = &tempFh

	return nil
}

type emptyBody struct {
}

var _ stream.BodyReader = (*emptyBody)(nil)

func (h emptyBody) Decode(stream.Reader) error {
	return nil
}

type encodeDecodeTest struct {
	msg     string
	value   wire.Value
	encoded []byte
}

func checkEncodeDecode(t *testing.T, typ wire.Type, tests []encodeDecodeTest) {
	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			buffer := bytes.Buffer{}

			// encode and match bytes
			err := Binary.Encode(tt.value, &buffer)
			if assert.NoError(t, err, "Encode failed:\n%s", tt.value) {
				assert.Equal(t, tt.encoded, buffer.Bytes())
			}

			// decode and match value
			value, err := Binary.Decode(bytes.NewReader(tt.encoded), typ)
			if assert.NoError(t, err, "Decode failed:\n%s", tt.value) {
				assert.True(
					t, wire.ValuesAreEqual(tt.value, value),
					fmt.Sprintf("\n\t   %v (expected)\n\t!= %v (actual)", tt.value, value),
				)
			}

			// encode the decoded value again
			buffer = bytes.Buffer{}
			err = Binary.Encode(value, &buffer)
			if assert.NoError(t, err, "Encode of decoded value failed:\n%s", tt.value) {
				assert.Equal(t, tt.encoded, buffer.Bytes())
			}
		})
	}
}

func getStreamReader(t *testing.T, encoded []byte) stream.Reader {
	t.Helper()

	return BinaryStreamer.Reader(bytes.NewReader(encoded))
}

type failureTest struct {
	msg     string
	encoded []byte
}

func checkDecodeFailure(t *testing.T, typ wire.Type, tests []failureTest) {
	for _, tt := range tests {
		value, err := Binary.Decode(bytes.NewReader(tt.encoded), typ)
		if err == nil {
			// lazy collections need to be fully evaluated for the failure to
			// propagate
			err = wire.EvaluateValue(value)
		}
		if assert.Error(t, err, "Expected failure parsing %x, got %s", tt, value) {
			assert.True(
				t,
				binary.IsDecodeError(err),
				"Expected decode error while parsing %x, got %s",
				tt,
				err,
			)
		}
	}
}

func checkEOFError(t *testing.T, typ wire.Type, tests []failureTest) {
	for _, tt := range tests {
		value, err := Binary.Decode(bytes.NewReader(tt.encoded), typ)
		if err == nil {
			// lazy collections need to be fully evaluated for the failure to
			// propagate
			err = wire.EvaluateValue(value)
		}
		if assert.Error(t, err, "Expected failure parsing %x, got %s", tt, value) {
			assert.Equal(
				t, io.ErrUnexpectedEOF, err,
				"Expected EOF error while parsing %x, got %s", tt, err,
			)
		}
	}
}

func TestBool(t *testing.T) {
	tests := []encodeDecodeTest{
		{"false", vbool(false), []byte{0x00}},
		{"true", vbool(true), []byte{0x01}},
	}

	checkEncodeDecode(t, wire.TBool, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetBool()
			val, err := reader.ReadBool()
			require.NoError(t, err)
			assert.Equal(t, want, val)
		})
	}
}

func TestBoolDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{"invalid", []byte{0x02}}, // values outside 0 and 1
	}

	checkDecodeFailure(t, wire.TBool, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadBool()
			require.Error(t, err)
		})
	}
}

func TestBoolEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
	}

	checkEOFError(t, wire.TBool, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadBool()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestI8(t *testing.T) {
	tests := []encodeDecodeTest{
		{"0", vi8(0), []byte{0x00}},
		{"1", vi8(1), []byte{0x01}},
		{"-1", vi8(-1), []byte{0xff}},
		{"127", vi8(127), []byte{0x7f}},
		{"-128", vi8(-128), []byte{0x80}},
	}

	checkEncodeDecode(t, wire.TI8, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetI8()
			val, err := reader.ReadInt8()
			require.NoError(t, err)
			assert.Equal(t, want, val)
		})
	}
}

func TestI8EOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
	}

	checkEOFError(t, wire.TI8, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadInt8()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestI16(t *testing.T) {
	tests := []encodeDecodeTest{
		{"1", vi16(1), []byte{0x00, 0x01}},
		{"255", vi16(255), []byte{0x00, 0xff}},
		{"256", vi16(256), []byte{0x01, 0x00}},
		{"257", vi16(257), []byte{0x01, 0x01}},
		{"32767", vi16(32767), []byte{0x7f, 0xff}},
		{"-1", vi16(-1), []byte{0xff, 0xff}},
		{"-2", vi16(-2), []byte{0xff, 0xfe}},
		{"-256", vi16(-256), []byte{0xff, 0x00}},
		{"-255", vi16(-255), []byte{0xff, 0x01}},
		{"-32768", vi16(-32768), []byte{0x80, 0x00}},
	}

	checkEncodeDecode(t, wire.TI16, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetI16()
			val, err := reader.ReadInt16()
			require.NoError(t, err)
			assert.Equal(t, want, val)
		})
	}
}

func TestI16EOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},     // empty
		{"short", []byte{0x00}}, // one byte too short
	}

	checkEOFError(t, wire.TI16, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadInt16()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestI32(t *testing.T) {
	tests := []encodeDecodeTest{
		{"1", vi32(1), []byte{0x00, 0x00, 0x00, 0x01}},
		{"255", vi32(255), []byte{0x00, 0x00, 0x00, 0xff}},
		{"65535", vi32(65535), []byte{0x00, 0x00, 0xff, 0xff}},
		{"16777215", vi32(16777215), []byte{0x00, 0xff, 0xff, 0xff}},
		{"2147483647", vi32(2147483647), []byte{0x7f, 0xff, 0xff, 0xff}},
		{"-1", vi32(-1), []byte{0xff, 0xff, 0xff, 0xff}},
		{"-256", vi32(-256), []byte{0xff, 0xff, 0xff, 0x00}},
		{"-65536", vi32(-65536), []byte{0xff, 0xff, 0x00, 0x00}},
		{"-16777216", vi32(-16777216), []byte{0xff, 0x00, 0x00, 0x00}},
		{"-2147483648", vi32(-2147483648), []byte{0x80, 0x00, 0x00, 0x00}},
	}

	checkEncodeDecode(t, wire.TI32, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetI32()
			val, err := reader.ReadInt32()
			require.NoError(t, err)
			assert.Equal(t, want, val)
		})
	}
}

func TestI32EOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},                 // empty
		{"short", []byte{0x01, 0x02, 0x03}}, // one byte too short
	}

	checkEOFError(t, wire.TI32, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadInt32()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestI64(t *testing.T) {
	tests := []encodeDecodeTest{
		{"1", vi64(1), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{"4294967295", vi64(4294967295), []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}},
		{"1099511627775", vi64(1099511627775), []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"281474976710655", vi64(281474976710655), []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"72057594037927935", vi64(72057594037927935), []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"9223372036854775807", vi64(9223372036854775807), []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"-1", vi64(-1), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{"-4294967296", vi64(-4294967296), []byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00}},
		{"-1099511627776", vi64(-1099511627776), []byte{0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"-281474976710656", vi64(-281474976710656), []byte{0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"-72057594037927936", vi64(-72057594037927936), []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"-9223372036854775808", vi64(-9223372036854775808), []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	checkEncodeDecode(t, wire.TI64, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetI64()
			val, err := reader.ReadInt64()
			require.NoError(t, err)
			assert.Equal(t, want, val)
		})
	}
}

func TestI64EOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
		{"short", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}}, // one byte too short
	}

	checkEOFError(t, wire.TI64, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadInt64()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestDouble(t *testing.T) {
	tests := []encodeDecodeTest{
		{"0.0", vdouble(0.0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"1.0", vdouble(1.0), []byte{0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{"1.0000000001", vdouble(1.0000000001), []byte{0x3f, 0xf0, 0x0, 0x0, 0x0, 0x6, 0xdf, 0x38}},
		{"1.1", vdouble(1.1), []byte{0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{"-1.1", vdouble(-1.1), []byte{0xbf, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{"3.141592653589793", vdouble(3.141592653589793), []byte{0x40, 0x9, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}},
		{"-1.0000000001", vdouble(-1.0000000001), []byte{0xbf, 0xf0, 0x0, 0x0, 0x0, 0x6, 0xdf, 0x38}},
		{"0", vdouble(math.Inf(0)), []byte{0x7f, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		{"-1", vdouble(math.Inf(-1)), []byte{0xff, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
	}

	checkEncodeDecode(t, wire.TDouble, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetDouble()
			val, err := reader.ReadDouble()
			require.NoError(t, err)
			assert.Equal(t, want, val)
		})
	}
}

func TestDoubleEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
		{"short", []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}}, // one byte too short
	}

	checkEOFError(t, wire.TDouble, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadDouble()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestDoubleNaN(t *testing.T) {
	value := vdouble(math.NaN())
	encoded := []byte{0x7f, 0xf8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1}

	buffer := bytes.Buffer{}
	err := Binary.Encode(value, &buffer)
	if assert.NoError(t, err, "Encode failed:\n%s", value) {
		assert.Equal(t, encoded, buffer.Bytes())
	}

	v, err := Binary.Decode(bytes.NewReader(encoded), wire.TDouble)
	if assert.NoError(t, err, "Decode failed:\n%s", value) {
		assert.Equal(t, wire.TDouble, v.Type())
		assert.True(t, math.IsNaN(v.GetDouble()))
	}

	reader := getStreamReader(t, encoded)
	val, err := reader.ReadDouble()
	require.NoError(t, err)
	assert.True(t, math.IsNaN(val))
}

func TestBinary(t *testing.T) {
	tests := []encodeDecodeTest{
		{"empty string", vbinary(""), []byte{0x00, 0x00, 0x00, 0x00}},
		{"hello ", vbinary("hello"), []byte{
			0x00, 0x00, 0x00, 0x05, // len:4 = 5
			0x68, 0x65, 0x6c, 0x6c, 0x6f, // 'h', 'e', 'l', 'l', 'o'
		}},
	}

	checkEncodeDecode(t, wire.TBinary, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			want := tt.value.GetBinary()
			val, err := reader.ReadString()
			require.NoError(t, err)
			assert.Equal(t, string(want), val)
		})
	}
}

func TestBinaryLargeLength(t *testing.T) {
	// 5 MB + 4 bytes for length
	data := make([]byte, 5242880+4)
	data[0], data[1], data[2], data[3] = 0x0, 0x50, 0x0, 0x0 // 5 MB

	value, err := Binary.Decode(bytes.NewReader(data), wire.TBinary)
	require.NoError(t, err, "failed to decode value")

	want := wire.NewValueBinary(data[4:])
	assert.True(t, wire.ValuesAreEqual(want, value), "values did not match")

	reader := getStreamReader(t, data)
	val, err := reader.ReadBinary()
	require.NoError(t, err, "failed to parse binary data")
	assert.Equal(t, data[4:], val)
}

func TestBinaryDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{"negative length", []byte{0xff, 0x30, 0x30, 0x30}}, // negative length
	}

	checkDecodeFailure(t, wire.TBinary, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadBinary()
			require.Error(t, err)
		})
	}
}

func TestBinaryEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},
		{"incomplete length", []byte{0x00}},                 // incomplete length
		{"length mismatch", []byte{0x00, 0x00, 0x00, 0x01}}, // length mismatch
		{"long length", []byte{0x22, 0x6e, 0x6f, 0x74}},     // really long length
	}

	checkEOFError(t, wire.TBinary, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadBinary()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

type limitWriter struct {
	b []byte
}

var errWriteLimitReached = errors.New("write limit reached")

func (lw *limitWriter) Write(p []byte) (int, error) {
	n := copy(lw.b, p)
	lw.b = lw.b[n:]
	if n < len(p) {
		return n, errWriteLimitReached
	}
	return n, nil
}

func TestStringEncodeFailure(t *testing.T) {
	// WriteString("hello") requires 9 bytes
	// - four bytes for length as int32, and
	// - five bytes for "hello"

	testCases := []struct {
		msg     string
		in      string
		len     int
		wantErr error
	}{
		{
			msg:     "int_write_failure",
			in:      "hello",
			len:     3,
			wantErr: errWriteLimitReached,
		},
		{
			msg:     "bytes_write_failure",
			in:      "hello",
			len:     8,
			wantErr: errWriteLimitReached,
		},
		{
			msg: "no_failure",
			in:  "hello",
			len: 9,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.msg, func(t *testing.T) {
			bw := &limitWriter{b: make([]byte, tt.len)}
			sw := binary.NewStreamWriter(bw)
			err := sw.WriteString(tt.in)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestStruct(t *testing.T) {
	tests := []encodeDecodeTest{
		{"empty struct", vstruct(), []byte{0x00}},
		{"simple struct", vstruct(vfield(1, vbool(true))), []byte{
			0x02,       // type:1 = bool
			0x00, 0x01, // id:2 = 1
			0x01, // value = true
			0x00, // stop
		}},
		{
			"complex struct",
			vstruct(
				vfield(1, vi16(42)),
				vfield(2, vlist(wire.TBinary, vbinary("foo"), vbinary("bar"))),
				vfield(3, vset(wire.TBinary, vbinary("baz"), vbinary("qux"))),
			), []byte{
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x2a, // value = 42

				0x0F,       // type:1 = list
				0x00, 0x02, // id:2 = 2

				// <list>
				0x0B,                   // type:1 = binary
				0x00, 0x00, 0x00, 0x02, // size:4 = 2
				// <binary>
				0x00, 0x00, 0x00, 0x03, // len:4 = 3
				0x66, 0x6f, 0x6f, // 'f', 'o', 'o'
				// </binary>
				// <binary>
				0x00, 0x00, 0x00, 0x03, // len:4 = 3
				0x62, 0x61, 0x72, // 'b', 'a', 'r'
				// </binary>
				// </list>

				0x0E,       // type = set
				0x00, 0x03, // id = 3

				// <set>
				0x0B,                   // type:1 = binary
				0x00, 0x00, 0x00, 0x02, // size:4 = 2
				// <binary>
				0x00, 0x00, 0x00, 0x03, // len:4 = 3
				0x62, 0x61, 0x7a, // 'b', 'a', 'z'
				// </binary>
				// <binary>
				0x00, 0x00, 0x00, 0x03, // len:4 = 3
				0x71, 0x75, 0x78, // 'q', 'u', 'x'
				// </binary>
				// </set>

				0x00, // stop
			},
		},
	}

	checkEncodeDecode(t, wire.TStruct, tests)

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			fields := tt.value.GetStruct().Fields
			reader := getStreamReader(t, tt.encoded)

			for i := 0; /* add 1 for the stop field */ i < len(fields)+1; i++ {
				err := reader.ReadStructBegin()
				require.NoError(t, err)

				fh, ok, err := reader.ReadFieldBegin()
				require.NoError(t, err)

				if !ok {
					assert.Equal(t, len(fields), i, "expected to have read all fields before stop-field")
					return
				}

				assert.Equal(t, fields[i].ID, fh.ID)
				assert.Equal(t, fields[i].Value.Type(), fh.Type)

				err = reader.Skip(fields[i].Value.Type())
				require.NoError(t, err)

				err = reader.ReadFieldEnd()
				require.NoError(t, err)

				err = reader.ReadStructEnd()
				require.NoError(t, err)
			}
		})
	}
}

func TestStructBeginAndEndEncode(t *testing.T) {
	var streamBuff bytes.Buffer

	// Encode with Streaming protocol
	w := binary.NewStreamWriter(&streamBuff)
	require.NoError(t, w.WriteStructBegin())
	require.NoError(t, w.WriteStructEnd())
	require.NoError(t, w.Close())

	// Assert that encoded bytes are equivalent
	assert.Equal(t, []byte{0x0}, streamBuff.Bytes())
}

func TestStructEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},
		{"invalid field ID", []byte{0x0B, 0x00}}, // invalid field ID
		{"no value", []byte{0x02, 0x00, 0x01}},   // no value
	}

	checkEOFError(t, wire.TStruct, tests)
}

func TestStructStreamingBegin(t *testing.T) {
	reader := getStreamReader(t, []byte{})
	assert.NoError(t, reader.ReadStructBegin())
}

func TestStructStreamingEnd(t *testing.T) {
	reader := getStreamReader(t, []byte{})
	assert.NoError(t, reader.ReadStructEnd())
}

func TestFieldStreamingBegin(t *testing.T) {
	tests := []struct {
		msg     string
		want    stream.FieldHeader
		encoded []byte
	}{
		{
			msg: "int32, ID:14834",
			want: stream.FieldHeader{
				Type: wire.TI32,
				ID:   int16(14834),
			},
			encoded: []byte{0x08, 0x39, 0xF2},
		},
		{
			msg: "double, ID:30091",
			want: stream.FieldHeader{
				Type: wire.TDouble,
				ID:   int16(30091),
			},
			encoded: []byte{0x04, 0x75, 0x8B},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			fh, _, err := reader.ReadFieldBegin()
			require.NoError(t, err)
			assert.Equal(t, tt.want, fh)
		})
	}
}

func TestFieldStreamingBeginFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},
		{"short ID", []byte{0x02, 0x04}},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, _, err := reader.ReadFieldBegin()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestFieldStreamingEnd(t *testing.T) {
	reader := getStreamReader(t, []byte{})
	assert.NoError(t, reader.ReadFieldEnd())
}

func TestMap(t *testing.T) {
	tests := []encodeDecodeTest{
		{"simple map", vmap(wire.TI64, wire.TBinary), []byte{0x0A, 0x0B, 0x00, 0x00, 0x00, 0x00}},
		{
			"complex map",
			vmap(
				wire.TBinary, wire.TList,
				vitem(vbinary("a"), vlist(wire.TI16, vi16(1))),
				vitem(vbinary("b"), vlist(wire.TI16, vi16(2), vi16(3))),
			), []byte{
				0x0B,                   // ktype = binary
				0x0F,                   // vtype = list
				0x00, 0x00, 0x00, 0x02, // count:4 = 2

				// <item>
				// <key>
				0x00, 0x00, 0x00, 0x01, // len:4 = 1
				0x61, // 'a'
				// </key>
				// <value>
				0x06,                   // type:1 = i16
				0x00, 0x00, 0x00, 0x01, // count:4 = 1
				0x00, 0x01, // 1
				// </value>
				// </item>

				// <item>
				// <key>
				0x00, 0x00, 0x00, 0x01, // len:4 = 1
				0x62, // 'b'
				// </key>
				// <value>
				0x06,                   // type:1 = i16
				0x00, 0x00, 0x00, 0x02, // count:4 = 2
				0x00, 0x02, // 2
				0x00, 0x03, // 3
				// </value>
				// </item>
			},
		},
	}

	checkEncodeDecode(t, wire.TMap, tests)
}

func TestMapBeginEncode(t *testing.T) {
	var (
		streamBuff bytes.Buffer
		err        error
	)

	// Encode with Streaming protocol
	w := binary.NewStreamWriter(&streamBuff)
	err = w.WriteMapBegin(stream.MapHeader{
		KeyType:   wire.TBinary,
		ValueType: wire.TBool,
		Length:    1,
	})
	require.NoError(t, err)
	require.NoError(t, w.Close())

	// Assert that encoded bytes are equivalent
	assert.Equal(t, []byte{0xb, 0x2, 0x0, 0x0, 0x0, 0x1}, streamBuff.Bytes())
}

func TestMapDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{"negative length",
			[]byte{
				0x08, 0x0B, // key: i32, value: binary
				0xff, 0x00, 0x00, 0x30, // negative length
			},
		},
	}

	checkDecodeFailure(t, wire.TMap, tests)
}

func TestMapEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
		{"no values", []byte{0x08, 0x0B, 0x00, 0x00, 0x00, 0x01}}, // no values
	}

	checkEOFError(t, wire.TMap, tests)
}

func TestMapStreamingBegin(t *testing.T) {
	tests := []struct {
		msg     string
		want    stream.MapHeader
		encoded []byte
	}{
		{
			msg: "int64:binary, 0-length",
			want: stream.MapHeader{
				KeyType:   wire.TI64,
				ValueType: wire.TBinary,
			},
			encoded: []byte{0x0A, 0x0B, 0x00, 0x00, 0x00, 0x00},
		},
		{
			msg: "binary:list, 2 length",
			want: stream.MapHeader{
				KeyType:   wire.TBinary,
				ValueType: wire.TList,
				Length:    2,
			},
			encoded: []byte{
				0x0B,                   // ktype = binary
				0x0F,                   // vtype = list
				0x00, 0x00, 0x00, 0x02, // count:4 = 2
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			mh, err := reader.ReadMapBegin()
			require.NoError(t, err)
			assert.Equal(t, tt.want, mh)
		})
	}
}

func TestMapStreamingBeginReadFailure(t *testing.T) {
	negativeLength := []byte{0x0A, 0x0B, 0x80, 0x00, 0x00, 0x22}

	reader := getStreamReader(t, negativeLength)
	_, err := reader.ReadMapBegin()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "got negative length")
}

func TestMapStreamingBeginReadEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},
		{"short value", []byte{0x0A}},
		{"short size", []byte{0x0A, 0x0B, 0x00, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadMapBegin()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestMapStreamingEnd(t *testing.T) {
	reader := getStreamReader(t, []byte{})
	assert.NoError(t, reader.ReadMapEnd())
}

func TestSet(t *testing.T) {
	tests := []encodeDecodeTest{
		{"small set", vset(wire.TBool), []byte{0x02, 0x00, 0x00, 0x00, 0x00}},
		{
			"large set", vset(wire.TBool, vbool(true), vbool(false), vbool(true)),
			[]byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x01},
		},
	}

	checkEncodeDecode(t, wire.TSet, tests)
}

func TestSetBeginEncode(t *testing.T) {
	var (
		streamBuff bytes.Buffer
		err        error
	)

	// Encode with Streaming protocol
	w := binary.NewStreamWriter(&streamBuff)
	err = w.WriteSetBegin(stream.SetHeader{
		Type:   wire.TList,
		Length: 1,
	})
	require.NoError(t, err)
	require.NoError(t, w.Close())

	// Assert that encoded bytes are equivalent
	assert.Equal(t, []byte{0xf, 0x0, 0x0, 0x0, 0x1}, streamBuff.Bytes())
}

func TestSetDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{"negative length",
			[]byte{
				0x08,                   // type: i32
				0xff, 0x00, 0x30, 0x30, // negative length
			},
		},
	}

	checkDecodeFailure(t, wire.TSet, tests)
}

func TestSetEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
		{"no values", []byte{0x08, 0x00, 0x00, 0x00, 0x01}}, // no values
	}

	checkEOFError(t, wire.TSet, tests)
}

func TestSetStreamingBegin(t *testing.T) {
	tests := []struct {
		msg     string
		want    stream.SetHeader
		encoded []byte
	}{
		{
			msg: "bool, 0-length",
			want: stream.SetHeader{
				Type:   wire.TBool,
				Length: 0,
			},
			encoded: []byte{0x02, 0x00, 0x00, 0x00, 0x00},
		},
		{
			msg: "int64, 10-length",
			want: stream.SetHeader{
				Type:   wire.TI64,
				Length: 10,
			},
			encoded: []byte{0x0A, 0x00, 0x00, 0x00, 0x0A},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			sh, err := reader.ReadSetBegin()
			require.NoError(t, err)
			assert.Equal(t, tt.want, sh)
		})
	}
}

func TestSetStreamingBeginReadFailure(t *testing.T) {
	negativeLength := []byte{0x0A, 0x80, 0x00, 0x00, 0x22}

	reader := getStreamReader(t, negativeLength)
	_, err := reader.ReadSetBegin()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "got negative length")
}

func TestSetStreamingBeginReadEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},
		{"short", []byte{0x0A, 0x0B, 0x00, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadSetBegin()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestSetStreamingEnd(t *testing.T) {
	reader := getStreamReader(t, []byte{})
	assert.NoError(t, reader.ReadSetEnd())
}

func TestList(t *testing.T) {
	tests := []encodeDecodeTest{
		{"small list", vlist(wire.TStruct), []byte{0x0C, 0x00, 0x00, 0x00, 0x00}},
		{
			"large list",
			vlist(
				wire.TStruct,
				vstruct(
					vfield(1, vi16(1)),
					vfield(2, vi32(2)),
				),
				vstruct(
					vfield(1, vi16(3)),
					vfield(2, vi32(4)),
				),
			),
			[]byte{
				0x0C,                   // vtype:1 = struct
				0x00, 0x00, 0x00, 0x02, // count:4 = 2

				// <struct>
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x01, // value = 1

				0x08,       // type:1 = i32
				0x00, 0x02, // id:2 = 2
				0x00, 0x00, 0x00, 0x02, // value = 2

				0x00, // stop
				// </struct>

				// <struct>
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x03, // value = 3

				0x08,       // type:1 = i32
				0x00, 0x02, // id:2 = 2
				0x00, 0x00, 0x00, 0x04, // value = 4

				0x00, // stop
				// </struct>
			},
		},
	}

	checkEncodeDecode(t, wire.TList, tests)
}

func TestListBeginEncode(t *testing.T) {
	var (
		streamBuff bytes.Buffer
		err        error
	)

	// Encode with Streaming protocol
	w := binary.NewStreamWriter(&streamBuff)
	err = w.WriteListBegin(stream.ListHeader{
		Type:   wire.TMap,
		Length: 5,
	})
	require.NoError(t, err)
	require.NoError(t, w.Close())

	// Assert that encoded bytes are equivalent
	assert.Equal(t, []byte{0xd, 0x0, 0x0, 0x0, 0x5}, streamBuff.Bytes())
}

func TestListDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{"negative length",
			[]byte{
				0x0B,                   // type: i32
				0xff, 0x00, 0x30, 0x00, // negative length
			},
		},
		{"invalid bool",
			[]byte{
				0x02, // type: bool
				0x00, 0x00, 0x00, 0x01,
				0x10, // invalid bool
			},
		},
	}

	checkDecodeFailure(t, wire.TList, tests)
}

func TestListEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}}, // empty
		{"no values", []byte{0x08, 0x00, 0x00, 0x00, 0x01}}, // no values
	}

	checkEOFError(t, wire.TList, tests)
}

func TestListStreamingBegin(t *testing.T) {
	tests := []struct {
		msg     string
		want    stream.ListHeader
		encoded []byte
	}{
		{
			msg: "struct, 0-length",
			want: stream.ListHeader{
				Type:   wire.TStruct,
				Length: 0,
			},
			encoded: []byte{0x0C, 0x00, 0x00, 0x00, 0x00},
		},
		{
			msg: "map, 14-length",
			want: stream.ListHeader{
				Type:   wire.TMap,
				Length: 14,
			},
			encoded: []byte{0x0D, 0x00, 0x00, 0x00, 0x0E},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			lh, err := reader.ReadListBegin()
			require.NoError(t, err)
			assert.Equal(t, tt.want, lh)
		})
	}
}

func TestListStreamingBeginReadFailure(t *testing.T) {
	negativeLength := []byte{0x0B, 0xFF, 0x00, 0x30, 0x00}

	reader := getStreamReader(t, negativeLength)
	_, err := reader.ReadListBegin()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "got negative length")
}

func TestListStreamingBeginReadEOFFailure(t *testing.T) {
	tests := []failureTest{
		{"empty", []byte{}},
		{"short", []byte{0x0A, 0x0B, 0x00, 0x01}},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			_, err := reader.ReadListBegin()
			require.Error(t, err)
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestListStreamingEnd(t *testing.T) {
	reader := getStreamReader(t, []byte{})
	assert.NoError(t, reader.ReadListEnd())
}

func TestStructOfContainers(t *testing.T) {
	tests := []encodeDecodeTest{
		{
			"struct of containers",
			vstruct(
				vfield(1, vlist(
					wire.TMap,
					vmap(
						wire.TI32, wire.TSet,
						vitem(vi32(1), vset(
							wire.TBinary,
							vbinary("a"), vbinary("b"), vbinary("c"),
						)),
						vitem(vi32(2), vset(wire.TBinary)),
						vitem(vi32(3), vset(
							wire.TBinary,
							vbinary("d"), vbinary("e"), vbinary("f"),
						)),
					),
					vmap(
						wire.TI32, wire.TSet,
						vitem(vi32(4), vset(wire.TBinary, vbinary("g"))),
					),
					vmap(
						wire.TI8, wire.TI16,
						vitem(vi8(2), vi16(4)),
						vitem(vi8(16), vi16(256)),
					),
				)),
				vfield(2, vlist(wire.TI16, vi16(1), vi16(2), vi16(3))),
			),
			[]byte{
				0x0f,       // type:list
				0x00, 0x01, // field ID 1

				0x0d,                   // type: map
				0x00, 0x00, 0x00, 0x03, // length: 3

				// <map-1>
				0x08, 0x0e, // ktype: i32, vtype: set
				0x00, 0x00, 0x00, 0x03, // length: 3

				// 1: {"a", "b", "c"}
				0x00, 0x00, 0x00, 0x01, // 1
				0x0B,                   // type: binary
				0x00, 0x00, 0x00, 0x03, // length: 3
				0x00, 0x00, 0x00, 0x01, 0x61, // 'a'
				0x00, 0x00, 0x00, 0x01, 0x62, // 'b'
				0x00, 0x00, 0x00, 0x01, 0x63, // 'c'

				// 2: {}
				0x00, 0x00, 0x00, 0x02, // 2
				0x0B,                   // type: binary
				0x00, 0x00, 0x00, 0x00, // length: 0

				// 3: {"d", "e", "f"}
				0x00, 0x00, 0x00, 0x03, // 3
				0x0B,                   // type: binary
				0x00, 0x00, 0x00, 0x03, // length: 3
				0x00, 0x00, 0x00, 0x01, 0x64, // 'd'
				0x00, 0x00, 0x00, 0x01, 0x65, // 'e'
				0x00, 0x00, 0x00, 0x01, 0x66, // 'f'

				// </map-1>

				// <map-2>
				0x08, 0x0e, // ktype: i32, vtype: set
				0x00, 0x00, 0x00, 0x01, // length: 1

				// 4: {"g"}
				0x00, 0x00, 0x00, 0x04, // 3
				0x0B,                   // type: binary
				0x00, 0x00, 0x00, 0x01, // length: 1
				0x00, 0x00, 0x00, 0x01, 0x67, // 'g'

				// </map-2>

				// <map-3>
				0x03, 0x06, // ktype: i8, vtype: i16
				0x00, 0x00, 0x00, 0x02, // length: 2
				0x02, 0x00, 0x04, // 2: 4
				0x10, 0x01, 0x00, // 16: 256

				0x0f,       // type: list
				0x00, 0x02, // field ID 2

				0x06,                   // type: i16
				0x00, 0x00, 0x00, 0x03, // length 3
				0x00, 0x01, 0x00, 0x02, 0x00, 0x03, // [1,2,3]

				0x00,
			},
		},
	}

	checkEncodeDecode(t, wire.TStruct, tests)
	for _, tt := range tests {
		t.Run("streaming skip structs", func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			err := reader.Skip(wire.TStruct)
			assert.NoError(t, err)

			_, err = reader.ReadInt8()
			require.Error(t, err, "error expected after reading all encoded data")
			assert.Equal(t, io.ErrUnexpectedEOF, err)
		})
	}
}

func TestSkipStreamingErrors(t *testing.T) {
	tests := []struct {
		msg      string
		skipType wire.Type
		encoded  []byte
	}{
		{
			msg:      "unknown type",
			skipType: wire.Type(1),
			encoded:  []byte{},
		},
		{
			msg:      "binary, empty",
			skipType: wire.TBinary,
			encoded:  []byte{},
		},
		{
			msg:      "binary, negative length",
			skipType: wire.TBinary,
			encoded:  []byte{0xFF, 0x00, 0x00, 0x01},
		},
		{
			msg:      "struct, empty",
			skipType: wire.TStruct,
			encoded:  []byte{},
		},
		{
			msg:      "struct, short ID",
			skipType: wire.TStruct,
			encoded:  []byte{0x0A, 0x00},
		},
		{
			msg:      "struct, short value",
			skipType: wire.TStruct,
			encoded:  []byte{0x08, 0x00, 0x02, 0x00, 0x00},
		},
		{
			msg:      "struct, no stop field",
			skipType: wire.TStruct,
			encoded:  []byte{0x03, 0x00, 0x08, 0x04},
		},
		{
			msg:      "map, empty",
			skipType: wire.TMap,
			encoded:  []byte{},
		},
		{
			msg:      "map, short value",
			skipType: wire.TMap,
			encoded:  []byte{0x04},
		},
		{
			msg:      "map, short size",
			skipType: wire.TMap,
			encoded:  []byte{0x04, 0x06},
		},
		{
			msg:      "map, negativesize",
			skipType: wire.TMap,
			encoded:  []byte{0x04, 0x06, 0xF0, 0xFF, 0x00, 0x00},
		},
		{
			msg:      "map, unknown key type",
			skipType: wire.TMap,
			encoded:  []byte{0x05, 0x06, 0x00, 0xFF, 0x00, 0x00},
		},
		{
			msg:      "map, unknown value type",
			skipType: wire.TMap,
			encoded:  []byte{0x06, 0x07, 0x00, 0xFF, 0x00, 0x00, 0x00, 0x01},
		},
		{
			msg:      "list, empty",
			skipType: wire.TList,
			encoded:  []byte{},
		},
		{
			msg:      "list, unknown type",
			skipType: wire.TList,
			encoded:  []byte{0x07, 0x00, 0x00, 0x00, 0xFF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			reader := getStreamReader(t, tt.encoded)
			err := reader.Skip(tt.skipType)
			require.Error(t, err)
		})
	}
}

func TestBinaryEnvelopeErrors(t *testing.T) {
	tests := []struct {
		encoded []byte
		errMsg  string
	}{
		{
			encoded: []byte{
				0x80, 0x02, 0x00, 0x01, // version|type:4 = 2 | call
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// <struct>
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x64, // value = 100
				0x00, // stop
			},
			errMsg: "cannot decode envelope of version",
		},
		{
			encoded: []byte{
				0x80, 0x02, 0x00, 0x01, // version|type:4 = 2 | call
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// <struct>
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x64, // value = 100
				0x00, // stop
			},
			errMsg: "cannot decode envelope of version",
		},
	}

	for _, tt := range tests {
		reader := bytes.NewReader(tt.encoded)
		_, err := Binary.DecodeEnveloped(reader)
		if !assert.Error(t, err, "%v: should fail", tt.errMsg) {
			continue
		}

		// Also verify that DecodeRequest fails.
		_, env, err := EnvelopeAgnosticBinary.DecodeRequest(wire.Call, reader)
		if !assert.Error(t, err, "%v: should fail to decode as request", tt.errMsg) {
			continue
		}

		if !assert.Equal(t, NoEnvelopeResponder, env, "%v: should fail with noEnvelopeResponder", tt.errMsg) {
			continue
		}

		assert.Contains(t, err.Error(), tt.errMsg, "Unexpected failure")
	}
}

func TestBinaryEnvelopeSuccessful(t *testing.T) {
	tests := []struct {
		msg               string
		encoded           []byte
		want              wire.Envelope
		wantResponderType reflect.Type
		reencode          bool
	}{
		{
			msg: "non-strict envelope, struct",
			encoded: []byte{
				0x00, 0x00, 0x00, 0x05, // length:4 = 5
				0x77, 0x72, 0x69, 0x74, 0x65, // 'write'
				0x04,                   // type:1 = OneWay
				0x00, 0x00, 0x00, 0x2a, // seqid:4 = 42

				// <struct>
				0x0B,       // ttype:1 = BINARY
				0x00, 0x01, // id:2 = 1
				0x00, 0x00, 0x00, 0x05, // length:4 = 5
				0x68, 0x65, 0x6c, 0x6c, 0x6f, // 'hello'
				0x00, // stop
			},
			want: wire.Envelope{
				Name:  "write",
				Type:  wire.OneWay,
				SeqID: 42,
				Value: vstruct(
					vfield(1, vbinary("hello")),
				),
			},
			wantResponderType: reflect.TypeOf((*EnvelopeV0Responder)(nil)),
		},
		{
			msg: "strict envelope, struct",
			encoded: []byte{
				0x80, 0x01, 0x00, 0x01, // version|type:4 = 1 | call
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// <struct>
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x64, // value = 100
				0x00, // stop
			},
			want: wire.Envelope{
				Name:  "abc",
				Type:  wire.Call,
				SeqID: 5436,
				Value: vstruct(
					vfield(1, vi16(100)),
				),
			},
			wantResponderType: reflect.TypeOf((*EnvelopeV1Responder)(nil)),
			reencode:          true,
		},
		{
			msg: "non-strict envelope, struct",
			encoded: []byte{
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x03,                   // type:1 = Exception
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// <struct>
				0x06,       // type:1 = i16
				0x00, 0x01, // id:2 = 1
				0x00, 0x64, // value = 100
				0x00, // stop
			},
			want: wire.Envelope{
				Name:  "abc",
				Type:  wire.Exception,
				SeqID: 5436,
				Value: vstruct(
					vfield(1, vi16(100)),
				),
			},
			wantResponderType: reflect.TypeOf((*EnvelopeV0Responder)(nil)),
		},
	}

	for _, tt := range tests {
		reader := bytes.NewReader(tt.encoded)
		e, err := Binary.DecodeEnveloped(reader)
		if !assert.NoError(t, err, "%v: failed to decode", tt.msg) {
			continue
		}

		if !assert.Equal(t, tt.want, e, "%v: decoded envelope mismatch") {
			continue
		}

		// Also verify whether we can infer the presence of the envelope
		// reliably when reading a request struct.

		r, responder, err := EnvelopeAgnosticBinary.DecodeRequest(tt.want.Type, reader)
		if !assert.NoError(t, err, "%v: failed to decode request with envelope", tt.msg) {
			continue
		}

		if !assert.Equal(t, tt.want.Value, r, "%v: decoded request mismatch", tt.msg) {
			continue
		}

		if !assert.True(t, tt.wantResponderType == reflect.TypeOf(responder), "%v: decoded request should have responder want %v got %T", tt.msg, tt.wantResponderType, responder) {
			continue
		}

		// Verify a round trip, back from encode after decode.

		if !tt.reencode {
			continue
		}

		buf := &bytes.Buffer{}
		if !assert.NoError(t, Binary.EncodeEnveloped(e, buf), "%v: failed to encode", tt.msg) {
			continue
		}

		assert.Equal(t, tt.encoded, buf.Bytes(), "%v: reencoded bytes mismatch")
	}
}

func TestStreamingEnvelopeErrors(t *testing.T) {
	tests := []struct {
		inputEnvType  wire.EnvelopeType
		inputReqBytes []byte
		errMsg        string
	}{
		{
			inputEnvType: wire.OneWay,
			inputReqBytes: []byte{
				// envelope
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x03,                   // type:1 = Exception
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// request
			},
			errMsg: "unexpected envelope type",
		},
		{
			inputEnvType: wire.Call,
			inputReqBytes: []byte{
				// envelope
				0x80, 0x02, 0x00, 0x01, // version|type:4 = 2 | call
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// request
			},
			errMsg: "cannot decode envelope of version",
		},
		{
			inputEnvType: wire.Call,
			inputReqBytes: []byte{
				// envelope
				0x80, 0x02, 0x00, 0x01, // version|type:4 = 2 | call
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// request
			},
			errMsg: "cannot decode envelope of version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.errMsg, func(t *testing.T) {
			responder, err := binary.Default.ReadRequest(context.Background(), tt.inputEnvType, bytes.NewReader(tt.inputReqBytes), &emptyBody{})
			require.Error(t, err, "ReadRequest should fail")
			require.Equal(t, binary.NoEnvelopeResponder, responder, "ReadRequest should fail with noEnvelopeResponder")
			assert.Contains(t, err.Error(), tt.errMsg, "ReadRequest should fail with error message")
		})
	}
}

func TestStreamingEnvelopeSuccessful(t *testing.T) {
	tests := []struct {
		msg               string
		inputReqBytes     []byte
		inputEnvType      wire.EnvelopeType
		wantReqFieldID    int16
		wantResponderType reflect.Type
		wantResBytes      []byte
	}{
		{
			msg:               "no envelope",
			inputEnvType:      wire.OneWay,
			inputReqBytes:     tbinary(vstruct(vfield(1, vbinary("hello")))),
			wantReqFieldID:    1,
			wantResponderType: reflect.TypeOf(binary.NoEnvelopeResponder),
			wantResBytes:      tbinary(vstruct()),
		},
		{
			msg:          "non-strict envelope, envelope type OneWay",
			inputEnvType: wire.OneWay,
			inputReqBytes: append(
				// envelope
				_testNonStrictEnvelopeOneWayBytes,

				// request
				tbinary(vstruct(vfield(2, vbinary("hello"))))...,
			),
			wantReqFieldID:    2,
			wantResponderType: reflect.TypeOf((*binary.EnvelopeV0Responder)(nil)),
			wantResBytes: append(
				// envelope
				_testNonStrictEnvelopeOneWayBytes,

				// response
				tbinary(vstruct())...,
			),
		},
		{
			msg:          "strict envelope, envelope type Call",
			inputEnvType: wire.Call,
			inputReqBytes: append(
				// envelope
				_testStrictEnvelopeCallBytes,

				// request
				tbinary(vstruct(vfield(3, vbinary("hello"))))...,
			),
			wantReqFieldID:    3,
			wantResponderType: reflect.TypeOf((*binary.EnvelopeV1Responder)(nil)),
			wantResBytes: append(
				// envelope
				_testStrictEnvelopeCallBytes,

				// response
				tbinary(vstruct())...,
			),
		},
		{
			msg:          "non-strict envelope, envelope type Exception",
			inputEnvType: wire.Exception,
			inputReqBytes: append(
				// envelope
				_testNonStrictEnvelopeExceptionBytes,

				// request
				tbinary(vstruct(vfield(4, vbinary("hello"))))...,
			),
			wantReqFieldID:    4,
			wantResponderType: reflect.TypeOf((*binary.EnvelopeV0Responder)(nil)),
			wantResBytes: append(
				// envelope
				_testNonStrictEnvelopeExceptionBytes,

				// response
				tbinary(vstruct())...,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.msg, func(t *testing.T) {
			var fh stream.FieldHeader

			// Verify whether we can infer the presence of the envelope
			// reliably when reading a request.
			h := testRequestBody{
				t:  t,
				fh: &fh,
			}

			responder, err := binary.Default.ReadRequest(context.Background(), tt.inputEnvType, bytes.NewReader(tt.inputReqBytes), &h)
			require.NoError(t, err, "failed to read request with envelope")
			require.True(t, tt.wantResponderType == reflect.TypeOf(responder), "read request should have responder want %v got %T", tt.wantResponderType, responder)

			// Verify whether we can read the correct request field ID from the stream.Reader
			// This is a basic test to ensure the stream.Reader is at the correct offset in the request.
			// It is not intended to test the functionality of stream.Reader
			assert.Equal(t, tt.wantReqFieldID, h.fh.ID, "request field ID mismatch")

			// Verify response has the correct envelope
			writer := &bytes.Buffer{}
			err = responder.WriteResponse(tt.inputEnvType, writer, &testResponse{t: t})
			require.NoError(t, err, "failed to write response with envelope")

			// Verify response can be written to with an empty struct
			// This is an action executed on each test to verify stream.Writer is at the correct offset
			require.Equal(t, tt.wantResBytes, writer.Bytes(), "encoded response envelope bytes mismatch")
		})
	}
}

func tbinary(v wire.Value) []byte {
	buf := &bytes.Buffer{}
	Binary.Encode(v, buf)
	return buf.Bytes()
}

func TestReqRes(t *testing.T) {

	complexPayload := vstruct(
		vfield(1, vi16(42)),
		vfield(2, vlist(wire.TBinary, vbinary("foo"), vbinary("bar"))),
		vfield(3, vset(wire.TBinary, vbinary("baz"), vbinary("qux"))),
		vfield(4, vmap(wire.TBinary, wire.TI8, vitem(vbinary("a"), vi8(1)))),
	)

	tests := []struct {
		msg           string
		req           wire.Value
		reqBytes      []byte
		responderType reflect.Type
		res           wire.Value
		resType       wire.EnvelopeType
		resBytes      []byte
	}{
		{
			msg:           "empty req, empty reply, no envelope",
			req:           vstruct(),
			reqBytes:      []byte{0x00},
			responderType: reflect.TypeOf(NoEnvelopeResponder),
			res:           vstruct(),
			resType:       wire.Reply,
			resBytes:      []byte{0x00},
		},
		{
			msg:           "two field req, empty reply, no envelope",
			req:           vstruct(vfield(1, vbool(true))),
			reqBytes:      tbinary(vstruct(vfield(1, vbool(true)))),
			responderType: reflect.TypeOf(NoEnvelopeResponder),
			res:           vstruct(),
			resType:       wire.Reply,
			resBytes:      []byte{0x00},
		},
		{
			msg: "empty reply, no-version non-strict envelope",
			reqBytes: append(
				[]byte{
					0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
					0x01,                   // type:1 = 1 = call
					0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436
				},
				tbinary(vstruct(vfield(1, vi16(100))))...,
			),
			req:           vstruct(vfield(1, vi16(100))),
			responderType: reflect.TypeOf((*EnvelopeV0Responder)(nil)),
			resType:       wire.Exception,
			res:           vstruct(),
			resBytes: []byte{
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x03,                   // type1 = 3 (exception)
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// <struct>
				0x00,
			},
		},
		{
			msg: "empty reply, version 1 strict envelope",
			reqBytes: append(
				[]byte{
					0x80, 0x01, 0x00, 0x01, // version|type:4 = 1 | call
					0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
					0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436
				},
				tbinary(vstruct(
					vfield(1, vi16(100)),
				))...,
			),
			req: vstruct(
				vfield(1, vi16(100)),
			),
			responderType: reflect.TypeOf((*EnvelopeV1Responder)(nil)),
			resType:       wire.Reply,
			res:           vstruct(),
			resBytes: []byte{
				0x80, 0x01, 0x00, 0x02, // version:2 &^ 0x80 = 1, type:2 = 2 (reply)
				0x00, 0x00, 0x00, 0x03, 'a', 'b', 'c', // name~4 = "abc"
				0x00, 0x00, 0x15, 0x3c, // seqID:4 = 5436

				// <struct>
				0x00,
			},
		},
		{
			msg:           "complex request, no envelope, complex response",
			req:           complexPayload,
			reqBytes:      tbinary(complexPayload),
			res:           complexPayload,
			responderType: reflect.TypeOf(NoEnvelopeResponder),
			resType:       wire.Reply,
			resBytes:      tbinary(complexPayload),
		},
	}

	// Verify that all structs can be read as request structs, dispite lacking
	// an envelope.
	for _, tt := range tests {
		reader := bytes.NewReader(tt.reqBytes)
		req, reser, err := EnvelopeAgnosticBinary.DecodeRequest(wire.Call, reader)
		if !assert.NoError(t, err, "%s: failed to decode struct as request without envelope", tt.msg) {
			continue
		}

		if !assert.Equal(t, tt.responderType, reflect.TypeOf(reser), "%s: responder type mismatch", tt.msg) {
			continue
		}

		if !assert.True(t, wire.ValuesAreEqual(tt.req, req), "%s: decoded request mismatch", tt.msg) {
			continue
		}

		writer := &bytes.Buffer{}
		err = reser.EncodeResponse(tt.res, tt.resType, writer)

		if !assert.NoError(t, err, "%s: failed to encode response", tt.msg) {
			continue
		}

		if !assert.Equal(t, tt.resBytes, writer.Bytes(), "%s: response bytes mismatch", tt.msg) {
			continue
		}

	}
}
