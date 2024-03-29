// Code generated by MockGen. DO NOT EDIT.
// Source: go.uber.org/thriftrw/ast (interfaces: Visitor)

// Package ast is a generated GoMock package.
package ast

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockVisitor is a mock of Visitor interface.
type MockVisitor struct {
	ctrl     *gomock.Controller
	recorder *MockVisitorMockRecorder
}

// MockVisitorMockRecorder is the mock recorder for MockVisitor.
type MockVisitorMockRecorder struct {
	mock *MockVisitor
}

// NewMockVisitor creates a new mock instance.
func NewMockVisitor(ctrl *gomock.Controller) *MockVisitor {
	mock := &MockVisitor{ctrl: ctrl}
	mock.recorder = &MockVisitorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVisitor) EXPECT() *MockVisitorMockRecorder {
	return m.recorder
}

// Visit mocks base method.
func (m *MockVisitor) Visit(arg0 Walker, arg1 Node) Visitor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Visit", arg0, arg1)
	ret0, _ := ret[0].(Visitor)
	return ret0
}

// Visit indicates an expected call of Visit.
func (mr *MockVisitorMockRecorder) Visit(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Visit", reflect.TypeOf((*MockVisitor)(nil).Visit), arg0, arg1)
}
