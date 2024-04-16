// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	http "net/http"
	reflect "reflect"
	time "time"

	core "github.com/totegamma/concurrent/x/core"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockClient) Commit(ctx context.Context, domain, body string) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit", ctx, domain, body)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Commit indicates an expected call of Commit.
func (mr *MockClientMockRecorder) Commit(ctx, domain, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockClient)(nil).Commit), ctx, domain, body)
}

// GetChunks mocks base method.
func (m *MockClient) GetChunks(ctx context.Context, domain string, timelines []string, queryTime time.Time) (map[string]core.Chunk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChunks", ctx, domain, timelines, queryTime)
	ret0, _ := ret[0].(map[string]core.Chunk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChunks indicates an expected call of GetChunks.
func (mr *MockClientMockRecorder) GetChunks(ctx, domain, timelines, queryTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChunks", reflect.TypeOf((*MockClient)(nil).GetChunks), ctx, domain, timelines, queryTime)
}

// GetDomain mocks base method.
func (m *MockClient) GetDomain(ctx context.Context, domain string) (core.Domain, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDomain", ctx, domain)
	ret0, _ := ret[0].(core.Domain)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDomain indicates an expected call of GetDomain.
func (mr *MockClientMockRecorder) GetDomain(ctx, domain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDomain", reflect.TypeOf((*MockClient)(nil).GetDomain), ctx, domain)
}

// GetEntity mocks base method.
func (m *MockClient) GetEntity(ctx context.Context, domain, address string) (core.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", ctx, domain, address)
	ret0, _ := ret[0].(core.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntity indicates an expected call of GetEntity.
func (mr *MockClientMockRecorder) GetEntity(ctx, domain, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockClient)(nil).GetEntity), ctx, domain, address)
}

// GetKey mocks base method.
func (m *MockClient) GetKey(ctx context.Context, domain, id string) ([]core.Key, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetKey", ctx, domain, id)
	ret0, _ := ret[0].([]core.Key)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetKey indicates an expected call of GetKey.
func (mr *MockClientMockRecorder) GetKey(ctx, domain, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetKey", reflect.TypeOf((*MockClient)(nil).GetKey), ctx, domain, id)
}

// GetTimeline mocks base method.
func (m *MockClient) GetTimeline(ctx context.Context, domain, id string) (core.Timeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTimeline", ctx, domain, id)
	ret0, _ := ret[0].(core.Timeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTimeline indicates an expected call of GetTimeline.
func (mr *MockClientMockRecorder) GetTimeline(ctx, domain, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTimeline", reflect.TypeOf((*MockClient)(nil).GetTimeline), ctx, domain, id)
}

// ResolveAddress mocks base method.
func (m *MockClient) ResolveAddress(ctx context.Context, domain, address string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveAddress", ctx, domain, address)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResolveAddress indicates an expected call of ResolveAddress.
func (mr *MockClientMockRecorder) ResolveAddress(ctx, domain, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveAddress", reflect.TypeOf((*MockClient)(nil).ResolveAddress), ctx, domain, address)
}
