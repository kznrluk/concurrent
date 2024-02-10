// Code generated by MockGen. DO NOT EDIT.
// Source: manager.go

// Package mock_socket is a generated GoMock package.
package mock_socket

import (
	reflect "reflect"

	websocket "github.com/gorilla/websocket"
	gomock "go.uber.org/mock/gomock"
)

// MockManager is a mock of Manager interface.
type MockManager struct {
	ctrl     *gomock.Controller
	recorder *MockManagerMockRecorder
}

// MockManagerMockRecorder is the mock recorder for MockManager.
type MockManagerMockRecorder struct {
	mock *MockManager
}

// NewMockManager creates a new mock instance.
func NewMockManager(ctrl *gomock.Controller) *MockManager {
	mock := &MockManager{ctrl: ctrl}
	mock.recorder = &MockManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockManager) EXPECT() *MockManagerMockRecorder {
	return m.recorder
}

// GetAllRemoteSubs mocks base method.
func (m *MockManager) GetAllRemoteSubs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllRemoteSubs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetAllRemoteSubs indicates an expected call of GetAllRemoteSubs.
func (mr *MockManagerMockRecorder) GetAllRemoteSubs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllRemoteSubs", reflect.TypeOf((*MockManager)(nil).GetAllRemoteSubs))
}

// Subscribe mocks base method.
func (m *MockManager) Subscribe(conn *websocket.Conn, streams []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Subscribe", conn, streams)
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockManagerMockRecorder) Subscribe(conn, streams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockManager)(nil).Subscribe), conn, streams)
}

// Unsubscribe mocks base method.
func (m *MockManager) Unsubscribe(conn *websocket.Conn) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Unsubscribe", conn)
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockManagerMockRecorder) Unsubscribe(conn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockManager)(nil).Unsubscribe), conn)
}