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
	"fmt"
	"io"
	"math"
	"reflect"
	"testing"

	"go.uber.org/thriftrw/internal/envelope/envelopetest"
	"go.uber.org/thriftrw/protocol/binary"
	"go.uber.org/thriftrw/wire"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type encodeDecodeTest struct {
	value   wire.Value
	encoded []byte
}

func checkEncodeDecode(t *testing.T, typ wire.Type, tests []encodeDecodeTest) {
	for _, tt := range tests {
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
	}
}

type failureTest []byte

func checkDecodeFailure(t *testing.T, typ wire.Type, tests []failureTest) {
	for _, tt := range tests {
		value, err := Binary.Decode(bytes.NewReader(tt), typ)
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
		value, err := Binary.Decode(bytes.NewReader(tt), typ)
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
		{vbool(false), []byte{0x00}},
		{vbool(true), []byte{0x01}},
	}

	checkEncodeDecode(t, wire.TBool, tests)
}

func TestBoolDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{0x02}, // values outside 0 and 1
	}

	checkDecodeFailure(t, wire.TBool, tests)
}

func TestBoolEOFFailure(t *testing.T) {
	tests := []failureTest{
		{}, // empty
	}

	checkEOFError(t, wire.TBool, tests)
}

func TestI8(t *testing.T) {
	tests := []encodeDecodeTest{
		{vi8(0), []byte{0x00}},
		{vi8(1), []byte{0x01}},
		{vi8(-1), []byte{0xff}},
		{vi8(127), []byte{0x7f}},
		{vi8(-128), []byte{0x80}},
	}

	checkEncodeDecode(t, wire.TI8, tests)
}

func TestI8EOFFailure(t *testing.T) {
	tests := []failureTest{
		{}, // empty
	}

	checkEOFError(t, wire.TI8, tests)
}

func TestI16(t *testing.T) {
	tests := []encodeDecodeTest{
		{vi16(1), []byte{0x00, 0x01}},
		{vi16(255), []byte{0x00, 0xff}},
		{vi16(256), []byte{0x01, 0x00}},
		{vi16(257), []byte{0x01, 0x01}},
		{vi16(32767), []byte{0x7f, 0xff}},
		{vi16(-1), []byte{0xff, 0xff}},
		{vi16(-2), []byte{0xff, 0xfe}},
		{vi16(-256), []byte{0xff, 0x00}},
		{vi16(-255), []byte{0xff, 0x01}},
		{vi16(-32768), []byte{0x80, 0x00}},
	}

	checkEncodeDecode(t, wire.TI16, tests)
}

func TestI16EOFFailure(t *testing.T) {
	tests := []failureTest{
		{},     // empty
		{0x00}, // one byte too short
	}

	checkEOFError(t, wire.TI16, tests)
}

func TestI32(t *testing.T) {
	tests := []encodeDecodeTest{
		{vi32(1), []byte{0x00, 0x00, 0x00, 0x01}},
		{vi32(255), []byte{0x00, 0x00, 0x00, 0xff}},
		{vi32(65535), []byte{0x00, 0x00, 0xff, 0xff}},
		{vi32(16777215), []byte{0x00, 0xff, 0xff, 0xff}},
		{vi32(2147483647), []byte{0x7f, 0xff, 0xff, 0xff}},
		{vi32(-1), []byte{0xff, 0xff, 0xff, 0xff}},
		{vi32(-256), []byte{0xff, 0xff, 0xff, 0x00}},
		{vi32(-65536), []byte{0xff, 0xff, 0x00, 0x00}},
		{vi32(-16777216), []byte{0xff, 0x00, 0x00, 0x00}},
		{vi32(-2147483648), []byte{0x80, 0x00, 0x00, 0x00}},
	}

	checkEncodeDecode(t, wire.TI32, tests)
}

func TestI32EOFFailure(t *testing.T) {
	tests := []failureTest{
		{},                 // empty
		{0x01, 0x02, 0x03}, // one byte too short
	}

	checkEOFError(t, wire.TI32, tests)
}

func TestI64(t *testing.T) {
	tests := []encodeDecodeTest{
		{vi64(1), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}},
		{vi64(4294967295), []byte{0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff}},
		{vi64(1099511627775), []byte{0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{vi64(281474976710655), []byte{0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{vi64(72057594037927935), []byte{0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{vi64(9223372036854775807), []byte{0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{vi64(-1), []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}},
		{vi64(-4294967296), []byte{0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00}},
		{vi64(-1099511627776), []byte{0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{vi64(-281474976710656), []byte{0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{vi64(-72057594037927936), []byte{0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{vi64(-9223372036854775808), []byte{0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
	}

	checkEncodeDecode(t, wire.TI64, tests)
}

func TestI64EOFFailure(t *testing.T) {
	tests := []failureTest{
		{}, // empty
		{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, // one byte too short
	}

	checkEOFError(t, wire.TI64, tests)
}

func TestDouble(t *testing.T) {
	tests := []encodeDecodeTest{
		{vdouble(0.0), []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{vdouble(1.0), []byte{0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{vdouble(1.0000000001), []byte{0x3f, 0xf0, 0x0, 0x0, 0x0, 0x6, 0xdf, 0x38}},
		{vdouble(1.1), []byte{0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{vdouble(-1.1), []byte{0xbf, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}},
		{vdouble(3.141592653589793), []byte{0x40, 0x9, 0x21, 0xfb, 0x54, 0x44, 0x2d, 0x18}},
		{vdouble(-1.0000000001), []byte{0xbf, 0xf0, 0x0, 0x0, 0x0, 0x6, 0xdf, 0x38}},
		{vdouble(math.Inf(0)), []byte{0x7f, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
		{vdouble(math.Inf(-1)), []byte{0xff, 0xf0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}},
	}

	checkEncodeDecode(t, wire.TDouble, tests)
}

func TestDoubleEOFFailure(t *testing.T) {
	tests := []failureTest{
		{}, // empty
		{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}, // one byte too short
	}

	checkEOFError(t, wire.TDouble, tests)
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
}

func TestBinary(t *testing.T) {
	tests := []encodeDecodeTest{
		{vbinary(""), []byte{0x00, 0x00, 0x00, 0x00}},
		{vbinary("hello"), []byte{
			0x00, 0x00, 0x00, 0x05, // len:4 = 5
			0x68, 0x65, 0x6c, 0x6c, 0x6f, // 'h', 'e', 'l', 'l', 'o'
		}},
	}

	checkEncodeDecode(t, wire.TBinary, tests)
}

func TestBinaryLargeLength(t *testing.T) {
	// 5 MB + 4 bytes for length
	data := make([]byte, 5242880+4)
	data[0], data[1], data[2], data[3] = 0x0, 0x50, 0x0, 0x0 // 5 MB

	value, err := Binary.Decode(bytes.NewReader(data), wire.TBinary)
	require.NoError(t, err, "failed to decode value")

	want := wire.NewValueBinary(data[4:])
	assert.True(t, wire.ValuesAreEqual(want, value), "values did not match")
}

func TestBinaryDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{0xff, 0x30, 0x30, 0x30}, // negative length
	}

	checkDecodeFailure(t, wire.TBinary, tests)
}

func TestBinaryEOFFailure(t *testing.T) {
	tests := []failureTest{
		{},
		{0x00},                   // incomplete length
		{0x00, 0x00, 0x00, 0x01}, // length mismatch
		{0x22, 0x6e, 0x6f, 0x74}, // really long length
	}

	checkEOFError(t, wire.TBinary, tests)
}

func TestStruct(t *testing.T) {
	tests := []encodeDecodeTest{
		{vstruct(), []byte{0x00}},
		{vstruct(vfield(1, vbool(true))), []byte{
			0x02,       // type:1 = bool
			0x00, 0x01, // id:2 = 1
			0x01, // value = true
			0x00, // stop
		}},
		{
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
}

func TestStructEOFFailure(t *testing.T) {
	tests := []failureTest{
		{},
		{0x0B, 0x00},       // invalid field ID
		{0x02, 0x00, 0x01}, // no value
	}

	checkEOFError(t, wire.TStruct, tests)
}

func TestMap(t *testing.T) {
	tests := []encodeDecodeTest{
		{vmap(wire.TI64, wire.TBinary), []byte{0x0A, 0x0B, 0x00, 0x00, 0x00, 0x00}},
		{
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

func TestMapDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{
			0x08, 0x0B, // key: i32, value: binary
			0xff, 0x00, 0x00, 0x30, // negative length
		},
	}

	checkDecodeFailure(t, wire.TMap, tests)
}

func TestMapEOFFailure(t *testing.T) {
	tests := []failureTest{
		{},                                   // empty
		{0x08, 0x0B, 0x00, 0x00, 0x00, 0x01}, // no values
	}

	checkEOFError(t, wire.TMap, tests)
}

func TestSet(t *testing.T) {
	tests := []encodeDecodeTest{
		{vset(wire.TBool), []byte{0x02, 0x00, 0x00, 0x00, 0x00}},
		{
			vset(wire.TBool, vbool(true), vbool(false), vbool(true)),
			[]byte{0x02, 0x00, 0x00, 0x00, 0x03, 0x01, 0x00, 0x01},
		},
	}

	checkEncodeDecode(t, wire.TSet, tests)
}

func TestSetDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{
			0x08,                   // type: i32
			0xff, 0x00, 0x30, 0x30, // negative length
		},
	}

	checkDecodeFailure(t, wire.TSet, tests)
}

func TestSetEOFFailure(t *testing.T) {
	tests := []failureTest{
		{},                             // empty
		{0x08, 0x00, 0x00, 0x00, 0x01}, // no values
	}

	checkEOFError(t, wire.TSet, tests)
}

func TestList(t *testing.T) {
	tests := []encodeDecodeTest{
		{vlist(wire.TStruct), []byte{0x0C, 0x00, 0x00, 0x00, 0x00}},
		{
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

func TestListDecodeFailure(t *testing.T) {
	tests := []failureTest{
		{
			0x0B,                   // type: i32
			0xff, 0x00, 0x30, 0x00, // negative length
		},
		{
			0x02, // type: bool
			0x00, 0x00, 0x00, 0x01,
			0x10, // invalid bool
		},
	}

	checkDecodeFailure(t, wire.TList, tests)
}

func TestListEOFFailure(t *testing.T) {
	tests := []failureTest{
		{},                             // empty
		{0x08, 0x00, 0x00, 0x00, 0x01}, // no values
	}

	checkEOFError(t, wire.TList, tests)
}

func TestStructOfContainers(t *testing.T) {
	tests := []encodeDecodeTest{
		{
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
				)),
				vfield(2, vlist(wire.TI16, vi16(1), vi16(2), vi16(3))),
			),
			[]byte{
				0x0f,       // type:list
				0x00, 0x01, // field ID 1

				0x0d,                   // type: map
				0x00, 0x00, 0x00, 0x02, // length: 2

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

		if !envelopetest.AssertEqual(t, tt.want, e, "%v: decoded envelope mismatch", tt) {
			continue
		}

		// Also verify whether we can infer the presence of the envelope
		// reliably when reading a request struct.

		r, responder, err := EnvelopeAgnosticBinary.DecodeRequest(tt.want.Type, reader)
		if !assert.NoError(t, err, "%v: failed to decode request with envelope", tt.msg) {
			continue
		}

		if !assert.True(t, wire.ValuesAreEqual(tt.want.Value, r), "%v: decoded request mismatch", tt.msg) {
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
