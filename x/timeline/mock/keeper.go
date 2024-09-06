// Code generated by MockGen. DO NOT EDIT.
// Source: keeper.go
//
// Generated by this command:
//
//	mockgen -source=keeper.go -destination=mock/keeper.go
//

// Package mock_timeline is a generated GoMock package.
package mock_timeline

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockKeeper is a mock of Keeper interface.
type MockKeeper struct {
	ctrl     *gomock.Controller
	recorder *MockKeeperMockRecorder
}

// MockKeeperMockRecorder is the mock recorder for MockKeeper.
type MockKeeperMockRecorder struct {
	mock *MockKeeper
}

// NewMockKeeper creates a new mock instance.
func NewMockKeeper(ctrl *gomock.Controller) *MockKeeper {
	mock := &MockKeeper{ctrl: ctrl}
	mock.recorder = &MockKeeperMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeeper) EXPECT() *MockKeeperMockRecorder {
	return m.recorder
}

// GetCurrentSubs mocks base method.
func (m *MockKeeper) GetCurrentSubs(ctx context.Context) []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentSubs", ctx)
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetCurrentSubs indicates an expected call of GetCurrentSubs.
func (mr *MockKeeperMockRecorder) GetCurrentSubs(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentSubs", reflect.TypeOf((*MockKeeper)(nil).GetCurrentSubs), ctx)
}

// GetRemoteSubs mocks base method.
func (m *MockKeeper) GetRemoteSubs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRemoteSubs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetRemoteSubs indicates an expected call of GetRemoteSubs.
func (mr *MockKeeperMockRecorder) GetRemoteSubs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRemoteSubs", reflect.TypeOf((*MockKeeper)(nil).GetRemoteSubs))
}

// Start mocks base method.
func (m *MockKeeper) Start(ctx context.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start", ctx)
}

// Start indicates an expected call of Start.
func (mr *MockKeeperMockRecorder) Start(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockKeeper)(nil).Start), ctx)
}
