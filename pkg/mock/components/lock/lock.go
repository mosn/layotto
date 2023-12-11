// Code generated by MockGen. DO NOT EDIT.
// Source: lock_store.go

// Package mock_lock is a generated GoMock package.
package mock_lock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	lock "mosn.io/layotto/components/lock"
)

// MockLockStore is a mock of LockStore interface.
type MockLockStore struct {
	ctrl     *gomock.Controller
	recorder *MockLockStoreMockRecorder
}

// MockLockStoreMockRecorder is the mock recorder for MockLockStore.
type MockLockStoreMockRecorder struct {
	mock *MockLockStore
}

// NewMockLockStore creates a new mock instance.
func NewMockLockStore(ctrl *gomock.Controller) *MockLockStore {
	mock := &MockLockStore{ctrl: ctrl}
	mock.recorder = &MockLockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLockStore) EXPECT() *MockLockStoreMockRecorder {
	return m.recorder
}

// Features mocks base method.
func (m *MockLockStore) Features() []lock.Feature {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Features")
	ret0, _ := ret[0].([]lock.Feature)
	return ret0
}

// Features indicates an expected call of Features.
func (mr *MockLockStoreMockRecorder) Features() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Features", reflect.TypeOf((*MockLockStore)(nil).Features))
}

// Init mocks base method.
func (m *MockLockStore) Init(metadata lock.Metadata) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Init", metadata)
	ret0, _ := ret[0].(error)
	return ret0
}

// Init indicates an expected call of Init.
func (mr *MockLockStoreMockRecorder) Init(metadata interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Init", reflect.TypeOf((*MockLockStore)(nil).Init), metadata)
}

// LockKeepAlive mocks base method.
func (m *MockLockStore) LockKeepAlive(arg0 context.Context, arg1 *lock.LockKeepAliveRequest) (*lock.LockKeepAliveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LockKeepAlive", arg0, arg1)
	ret0, _ := ret[0].(*lock.LockKeepAliveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LockKeepAlive indicates an expected call of LockKeepAlive.
func (mr *MockLockStoreMockRecorder) LockKeepAlive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LockKeepAlive", reflect.TypeOf((*MockLockStore)(nil).LockKeepAlive), arg0, arg1)
}

// TryLock mocks base method.
func (m *MockLockStore) TryLock(ctx context.Context, req *lock.TryLockRequest) (*lock.TryLockResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryLock", req)
	ret0, _ := ret[0].(*lock.TryLockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryLock indicates an expected call of TryLock.
func (mr *MockLockStoreMockRecorder) TryLock(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryLock", reflect.TypeOf((*MockLockStore)(nil).TryLock), req)
}

// Unlock mocks base method.
func (m *MockLockStore) Unlock(ctx context.Context, req *lock.UnlockRequest) (*lock.UnlockResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unlock", req)
	ret0, _ := ret[0].(*lock.UnlockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Unlock indicates an expected call of Unlock.
func (mr *MockLockStoreMockRecorder) Unlock(req interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockLockStore)(nil).Unlock), req)
}
