// Code generated by MockGen. DO NOT EDIT.

// Copyright (c) 2024 Uber Technologies, Inc.
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
// Source: go.uber.org/thriftrw/protocol/stream (interfaces: Protocol,Writer,Reader,BodyReader,Enveloper,RequestReader,ResponseWriter)

// Package streamtest is a generated GoMock package.
package streamtest

import (
	context "context"
	io "io"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	stream "go.uber.org/thriftrw/protocol/stream"
	wire "go.uber.org/thriftrw/wire"
)

// MockProtocol is a mock of Protocol interface.
type MockProtocol struct {
	ctrl     *gomock.Controller
	recorder *MockProtocolMockRecorder
}

// MockProtocolMockRecorder is the mock recorder for MockProtocol.
type MockProtocolMockRecorder struct {
	mock *MockProtocol
}

// NewMockProtocol creates a new mock instance.
func NewMockProtocol(ctrl *gomock.Controller) *MockProtocol {
	mock := &MockProtocol{ctrl: ctrl}
	mock.recorder = &MockProtocolMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProtocol) EXPECT() *MockProtocolMockRecorder {
	return m.recorder
}

// Reader mocks base method.
func (m *MockProtocol) Reader(arg0 io.Reader) stream.Reader {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reader", arg0)
	ret0, _ := ret[0].(stream.Reader)
	return ret0
}

// Reader indicates an expected call of Reader.
func (mr *MockProtocolMockRecorder) Reader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reader", reflect.TypeOf((*MockProtocol)(nil).Reader), arg0)
}

// Writer mocks base method.
func (m *MockProtocol) Writer(arg0 io.Writer) stream.Writer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Writer", arg0)
	ret0, _ := ret[0].(stream.Writer)
	return ret0
}

// Writer indicates an expected call of Writer.
func (mr *MockProtocolMockRecorder) Writer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Writer", reflect.TypeOf((*MockProtocol)(nil).Writer), arg0)
}

// MockWriter is a mock of Writer interface.
type MockWriter struct {
	ctrl     *gomock.Controller
	recorder *MockWriterMockRecorder
}

// MockWriterMockRecorder is the mock recorder for MockWriter.
type MockWriterMockRecorder struct {
	mock *MockWriter
}

// NewMockWriter creates a new mock instance.
func NewMockWriter(ctrl *gomock.Controller) *MockWriter {
	mock := &MockWriter{ctrl: ctrl}
	mock.recorder = &MockWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWriter) EXPECT() *MockWriterMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockWriter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockWriter)(nil).Close))
}

// WriteBinary mocks base method.
func (m *MockWriter) WriteBinary(arg0 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteBinary", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteBinary indicates an expected call of WriteBinary.
func (mr *MockWriterMockRecorder) WriteBinary(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBinary", reflect.TypeOf((*MockWriter)(nil).WriteBinary), arg0)
}

// WriteBool mocks base method.
func (m *MockWriter) WriteBool(arg0 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteBool", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteBool indicates an expected call of WriteBool.
func (mr *MockWriterMockRecorder) WriteBool(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteBool", reflect.TypeOf((*MockWriter)(nil).WriteBool), arg0)
}

// WriteDouble mocks base method.
func (m *MockWriter) WriteDouble(arg0 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteDouble", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteDouble indicates an expected call of WriteDouble.
func (mr *MockWriterMockRecorder) WriteDouble(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteDouble", reflect.TypeOf((*MockWriter)(nil).WriteDouble), arg0)
}

// WriteEnvelopeBegin mocks base method.
func (m *MockWriter) WriteEnvelopeBegin(arg0 stream.EnvelopeHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteEnvelopeBegin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteEnvelopeBegin indicates an expected call of WriteEnvelopeBegin.
func (mr *MockWriterMockRecorder) WriteEnvelopeBegin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteEnvelopeBegin", reflect.TypeOf((*MockWriter)(nil).WriteEnvelopeBegin), arg0)
}

// WriteEnvelopeEnd mocks base method.
func (m *MockWriter) WriteEnvelopeEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteEnvelopeEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteEnvelopeEnd indicates an expected call of WriteEnvelopeEnd.
func (mr *MockWriterMockRecorder) WriteEnvelopeEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteEnvelopeEnd", reflect.TypeOf((*MockWriter)(nil).WriteEnvelopeEnd))
}

// WriteFieldBegin mocks base method.
func (m *MockWriter) WriteFieldBegin(arg0 stream.FieldHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFieldBegin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFieldBegin indicates an expected call of WriteFieldBegin.
func (mr *MockWriterMockRecorder) WriteFieldBegin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFieldBegin", reflect.TypeOf((*MockWriter)(nil).WriteFieldBegin), arg0)
}

// WriteFieldEnd mocks base method.
func (m *MockWriter) WriteFieldEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFieldEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFieldEnd indicates an expected call of WriteFieldEnd.
func (mr *MockWriterMockRecorder) WriteFieldEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFieldEnd", reflect.TypeOf((*MockWriter)(nil).WriteFieldEnd))
}

// WriteInt16 mocks base method.
func (m *MockWriter) WriteInt16(arg0 int16) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteInt16", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteInt16 indicates an expected call of WriteInt16.
func (mr *MockWriterMockRecorder) WriteInt16(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt16", reflect.TypeOf((*MockWriter)(nil).WriteInt16), arg0)
}

// WriteInt32 mocks base method.
func (m *MockWriter) WriteInt32(arg0 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteInt32", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteInt32 indicates an expected call of WriteInt32.
func (mr *MockWriterMockRecorder) WriteInt32(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt32", reflect.TypeOf((*MockWriter)(nil).WriteInt32), arg0)
}

// WriteInt64 mocks base method.
func (m *MockWriter) WriteInt64(arg0 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteInt64", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteInt64 indicates an expected call of WriteInt64.
func (mr *MockWriterMockRecorder) WriteInt64(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt64", reflect.TypeOf((*MockWriter)(nil).WriteInt64), arg0)
}

// WriteInt8 mocks base method.
func (m *MockWriter) WriteInt8(arg0 int8) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteInt8", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteInt8 indicates an expected call of WriteInt8.
func (mr *MockWriterMockRecorder) WriteInt8(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteInt8", reflect.TypeOf((*MockWriter)(nil).WriteInt8), arg0)
}

// WriteListBegin mocks base method.
func (m *MockWriter) WriteListBegin(arg0 stream.ListHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteListBegin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteListBegin indicates an expected call of WriteListBegin.
func (mr *MockWriterMockRecorder) WriteListBegin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteListBegin", reflect.TypeOf((*MockWriter)(nil).WriteListBegin), arg0)
}

// WriteListEnd mocks base method.
func (m *MockWriter) WriteListEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteListEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteListEnd indicates an expected call of WriteListEnd.
func (mr *MockWriterMockRecorder) WriteListEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteListEnd", reflect.TypeOf((*MockWriter)(nil).WriteListEnd))
}

// WriteMapBegin mocks base method.
func (m *MockWriter) WriteMapBegin(arg0 stream.MapHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMapBegin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMapBegin indicates an expected call of WriteMapBegin.
func (mr *MockWriterMockRecorder) WriteMapBegin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMapBegin", reflect.TypeOf((*MockWriter)(nil).WriteMapBegin), arg0)
}

// WriteMapEnd mocks base method.
func (m *MockWriter) WriteMapEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteMapEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMapEnd indicates an expected call of WriteMapEnd.
func (mr *MockWriterMockRecorder) WriteMapEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMapEnd", reflect.TypeOf((*MockWriter)(nil).WriteMapEnd))
}

// WriteSetBegin mocks base method.
func (m *MockWriter) WriteSetBegin(arg0 stream.SetHeader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSetBegin", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSetBegin indicates an expected call of WriteSetBegin.
func (mr *MockWriterMockRecorder) WriteSetBegin(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSetBegin", reflect.TypeOf((*MockWriter)(nil).WriteSetBegin), arg0)
}

// WriteSetEnd mocks base method.
func (m *MockWriter) WriteSetEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteSetEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteSetEnd indicates an expected call of WriteSetEnd.
func (mr *MockWriterMockRecorder) WriteSetEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteSetEnd", reflect.TypeOf((*MockWriter)(nil).WriteSetEnd))
}

// WriteString mocks base method.
func (m *MockWriter) WriteString(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteString", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteString indicates an expected call of WriteString.
func (mr *MockWriterMockRecorder) WriteString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteString", reflect.TypeOf((*MockWriter)(nil).WriteString), arg0)
}

// WriteStructBegin mocks base method.
func (m *MockWriter) WriteStructBegin() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteStructBegin")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteStructBegin indicates an expected call of WriteStructBegin.
func (mr *MockWriterMockRecorder) WriteStructBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteStructBegin", reflect.TypeOf((*MockWriter)(nil).WriteStructBegin))
}

// WriteStructEnd mocks base method.
func (m *MockWriter) WriteStructEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteStructEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteStructEnd indicates an expected call of WriteStructEnd.
func (mr *MockWriterMockRecorder) WriteStructEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteStructEnd", reflect.TypeOf((*MockWriter)(nil).WriteStructEnd))
}

// MockReader is a mock of Reader interface.
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader.
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance.
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockReader) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockReaderMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockReader)(nil).Close))
}

// ReadBinary mocks base method.
func (m *MockReader) ReadBinary() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadBinary")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadBinary indicates an expected call of ReadBinary.
func (mr *MockReaderMockRecorder) ReadBinary() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadBinary", reflect.TypeOf((*MockReader)(nil).ReadBinary))
}

// ReadBool mocks base method.
func (m *MockReader) ReadBool() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadBool")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadBool indicates an expected call of ReadBool.
func (mr *MockReaderMockRecorder) ReadBool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadBool", reflect.TypeOf((*MockReader)(nil).ReadBool))
}

// ReadDouble mocks base method.
func (m *MockReader) ReadDouble() (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadDouble")
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadDouble indicates an expected call of ReadDouble.
func (mr *MockReaderMockRecorder) ReadDouble() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadDouble", reflect.TypeOf((*MockReader)(nil).ReadDouble))
}

// ReadEnvelopeBegin mocks base method.
func (m *MockReader) ReadEnvelopeBegin() (stream.EnvelopeHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadEnvelopeBegin")
	ret0, _ := ret[0].(stream.EnvelopeHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadEnvelopeBegin indicates an expected call of ReadEnvelopeBegin.
func (mr *MockReaderMockRecorder) ReadEnvelopeBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEnvelopeBegin", reflect.TypeOf((*MockReader)(nil).ReadEnvelopeBegin))
}

// ReadEnvelopeEnd mocks base method.
func (m *MockReader) ReadEnvelopeEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadEnvelopeEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadEnvelopeEnd indicates an expected call of ReadEnvelopeEnd.
func (mr *MockReaderMockRecorder) ReadEnvelopeEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadEnvelopeEnd", reflect.TypeOf((*MockReader)(nil).ReadEnvelopeEnd))
}

// ReadFieldBegin mocks base method.
func (m *MockReader) ReadFieldBegin() (stream.FieldHeader, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFieldBegin")
	ret0, _ := ret[0].(stream.FieldHeader)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadFieldBegin indicates an expected call of ReadFieldBegin.
func (mr *MockReaderMockRecorder) ReadFieldBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFieldBegin", reflect.TypeOf((*MockReader)(nil).ReadFieldBegin))
}

// ReadFieldEnd mocks base method.
func (m *MockReader) ReadFieldEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadFieldEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadFieldEnd indicates an expected call of ReadFieldEnd.
func (mr *MockReaderMockRecorder) ReadFieldEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadFieldEnd", reflect.TypeOf((*MockReader)(nil).ReadFieldEnd))
}

// ReadInt16 mocks base method.
func (m *MockReader) ReadInt16() (int16, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadInt16")
	ret0, _ := ret[0].(int16)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadInt16 indicates an expected call of ReadInt16.
func (mr *MockReaderMockRecorder) ReadInt16() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadInt16", reflect.TypeOf((*MockReader)(nil).ReadInt16))
}

// ReadInt32 mocks base method.
func (m *MockReader) ReadInt32() (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadInt32")
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadInt32 indicates an expected call of ReadInt32.
func (mr *MockReaderMockRecorder) ReadInt32() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadInt32", reflect.TypeOf((*MockReader)(nil).ReadInt32))
}

// ReadInt64 mocks base method.
func (m *MockReader) ReadInt64() (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadInt64")
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadInt64 indicates an expected call of ReadInt64.
func (mr *MockReaderMockRecorder) ReadInt64() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadInt64", reflect.TypeOf((*MockReader)(nil).ReadInt64))
}

// ReadInt8 mocks base method.
func (m *MockReader) ReadInt8() (int8, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadInt8")
	ret0, _ := ret[0].(int8)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadInt8 indicates an expected call of ReadInt8.
func (mr *MockReaderMockRecorder) ReadInt8() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadInt8", reflect.TypeOf((*MockReader)(nil).ReadInt8))
}

// ReadListBegin mocks base method.
func (m *MockReader) ReadListBegin() (stream.ListHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadListBegin")
	ret0, _ := ret[0].(stream.ListHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadListBegin indicates an expected call of ReadListBegin.
func (mr *MockReaderMockRecorder) ReadListBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadListBegin", reflect.TypeOf((*MockReader)(nil).ReadListBegin))
}

// ReadListEnd mocks base method.
func (m *MockReader) ReadListEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadListEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadListEnd indicates an expected call of ReadListEnd.
func (mr *MockReaderMockRecorder) ReadListEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadListEnd", reflect.TypeOf((*MockReader)(nil).ReadListEnd))
}

// ReadMapBegin mocks base method.
func (m *MockReader) ReadMapBegin() (stream.MapHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMapBegin")
	ret0, _ := ret[0].(stream.MapHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadMapBegin indicates an expected call of ReadMapBegin.
func (mr *MockReaderMockRecorder) ReadMapBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMapBegin", reflect.TypeOf((*MockReader)(nil).ReadMapBegin))
}

// ReadMapEnd mocks base method.
func (m *MockReader) ReadMapEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadMapEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadMapEnd indicates an expected call of ReadMapEnd.
func (mr *MockReaderMockRecorder) ReadMapEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadMapEnd", reflect.TypeOf((*MockReader)(nil).ReadMapEnd))
}

// ReadSetBegin mocks base method.
func (m *MockReader) ReadSetBegin() (stream.SetHeader, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSetBegin")
	ret0, _ := ret[0].(stream.SetHeader)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadSetBegin indicates an expected call of ReadSetBegin.
func (mr *MockReaderMockRecorder) ReadSetBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSetBegin", reflect.TypeOf((*MockReader)(nil).ReadSetBegin))
}

// ReadSetEnd mocks base method.
func (m *MockReader) ReadSetEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadSetEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadSetEnd indicates an expected call of ReadSetEnd.
func (mr *MockReaderMockRecorder) ReadSetEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadSetEnd", reflect.TypeOf((*MockReader)(nil).ReadSetEnd))
}

// ReadString mocks base method.
func (m *MockReader) ReadString() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadString")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadString indicates an expected call of ReadString.
func (mr *MockReaderMockRecorder) ReadString() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadString", reflect.TypeOf((*MockReader)(nil).ReadString))
}

// ReadStructBegin mocks base method.
func (m *MockReader) ReadStructBegin() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadStructBegin")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadStructBegin indicates an expected call of ReadStructBegin.
func (mr *MockReaderMockRecorder) ReadStructBegin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadStructBegin", reflect.TypeOf((*MockReader)(nil).ReadStructBegin))
}

// ReadStructEnd mocks base method.
func (m *MockReader) ReadStructEnd() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadStructEnd")
	ret0, _ := ret[0].(error)
	return ret0
}

// ReadStructEnd indicates an expected call of ReadStructEnd.
func (mr *MockReaderMockRecorder) ReadStructEnd() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadStructEnd", reflect.TypeOf((*MockReader)(nil).ReadStructEnd))
}

// Skip mocks base method.
func (m *MockReader) Skip(arg0 wire.Type) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Skip", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Skip indicates an expected call of Skip.
func (mr *MockReaderMockRecorder) Skip(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Skip", reflect.TypeOf((*MockReader)(nil).Skip), arg0)
}

// MockBodyReader is a mock of BodyReader interface.
type MockBodyReader struct {
	ctrl     *gomock.Controller
	recorder *MockBodyReaderMockRecorder
}

// MockBodyReaderMockRecorder is the mock recorder for MockBodyReader.
type MockBodyReaderMockRecorder struct {
	mock *MockBodyReader
}

// NewMockBodyReader creates a new mock instance.
func NewMockBodyReader(ctrl *gomock.Controller) *MockBodyReader {
	mock := &MockBodyReader{ctrl: ctrl}
	mock.recorder = &MockBodyReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBodyReader) EXPECT() *MockBodyReaderMockRecorder {
	return m.recorder
}

// Decode mocks base method.
func (m *MockBodyReader) Decode(arg0 stream.Reader) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Decode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Decode indicates an expected call of Decode.
func (mr *MockBodyReaderMockRecorder) Decode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Decode", reflect.TypeOf((*MockBodyReader)(nil).Decode), arg0)
}

// MockEnveloper is a mock of Enveloper interface.
type MockEnveloper struct {
	ctrl     *gomock.Controller
	recorder *MockEnveloperMockRecorder
}

// MockEnveloperMockRecorder is the mock recorder for MockEnveloper.
type MockEnveloperMockRecorder struct {
	mock *MockEnveloper
}

// NewMockEnveloper creates a new mock instance.
func NewMockEnveloper(ctrl *gomock.Controller) *MockEnveloper {
	mock := &MockEnveloper{ctrl: ctrl}
	mock.recorder = &MockEnveloperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEnveloper) EXPECT() *MockEnveloperMockRecorder {
	return m.recorder
}

// Encode mocks base method.
func (m *MockEnveloper) Encode(arg0 stream.Writer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Encode", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Encode indicates an expected call of Encode.
func (mr *MockEnveloperMockRecorder) Encode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Encode", reflect.TypeOf((*MockEnveloper)(nil).Encode), arg0)
}

// EnvelopeType mocks base method.
func (m *MockEnveloper) EnvelopeType() wire.EnvelopeType {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnvelopeType")
	ret0, _ := ret[0].(wire.EnvelopeType)
	return ret0
}

// EnvelopeType indicates an expected call of EnvelopeType.
func (mr *MockEnveloperMockRecorder) EnvelopeType() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnvelopeType", reflect.TypeOf((*MockEnveloper)(nil).EnvelopeType))
}

// MethodName mocks base method.
func (m *MockEnveloper) MethodName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MethodName")
	ret0, _ := ret[0].(string)
	return ret0
}

// MethodName indicates an expected call of MethodName.
func (mr *MockEnveloperMockRecorder) MethodName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MethodName", reflect.TypeOf((*MockEnveloper)(nil).MethodName))
}

// MockRequestReader is a mock of RequestReader interface.
type MockRequestReader struct {
	ctrl     *gomock.Controller
	recorder *MockRequestReaderMockRecorder
}

// MockRequestReaderMockRecorder is the mock recorder for MockRequestReader.
type MockRequestReaderMockRecorder struct {
	mock *MockRequestReader
}

// NewMockRequestReader creates a new mock instance.
func NewMockRequestReader(ctrl *gomock.Controller) *MockRequestReader {
	mock := &MockRequestReader{ctrl: ctrl}
	mock.recorder = &MockRequestReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequestReader) EXPECT() *MockRequestReaderMockRecorder {
	return m.recorder
}

// ReadRequest mocks base method.
func (m *MockRequestReader) ReadRequest(arg0 context.Context, arg1 wire.EnvelopeType, arg2 io.Reader, arg3 stream.BodyReader) (stream.ResponseWriter, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadRequest", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(stream.ResponseWriter)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadRequest indicates an expected call of ReadRequest.
func (mr *MockRequestReaderMockRecorder) ReadRequest(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadRequest", reflect.TypeOf((*MockRequestReader)(nil).ReadRequest), arg0, arg1, arg2, arg3)
}

// MockResponseWriter is a mock of ResponseWriter interface.
type MockResponseWriter struct {
	ctrl     *gomock.Controller
	recorder *MockResponseWriterMockRecorder
}

// MockResponseWriterMockRecorder is the mock recorder for MockResponseWriter.
type MockResponseWriterMockRecorder struct {
	mock *MockResponseWriter
}

// NewMockResponseWriter creates a new mock instance.
func NewMockResponseWriter(ctrl *gomock.Controller) *MockResponseWriter {
	mock := &MockResponseWriter{ctrl: ctrl}
	mock.recorder = &MockResponseWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResponseWriter) EXPECT() *MockResponseWriterMockRecorder {
	return m.recorder
}

// WriteResponse mocks base method.
func (m *MockResponseWriter) WriteResponse(arg0 wire.EnvelopeType, arg1 io.Writer, arg2 stream.Enveloper) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteResponse", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteResponse indicates an expected call of WriteResponse.
func (mr *MockResponseWriterMockRecorder) WriteResponse(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteResponse", reflect.TypeOf((*MockResponseWriter)(nil).WriteResponse), arg0, arg1, arg2)
}
