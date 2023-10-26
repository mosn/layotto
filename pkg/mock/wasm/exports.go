// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/wasm/exports.go

// Package mock_wasm is a generated GoMock package.
package mock_wasm

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1 "mosn.io/proxy-wasm-go-host/proxywasm/v1"
)

// MockExports is a mock of Exports interface.
type MockExports struct {
	ctrl     *gomock.Controller
	recorder *MockExportsMockRecorder
}

// MockExportsMockRecorder is the mock recorder for MockExports.
type MockExportsMockRecorder struct {
	mock *MockExports
}

// NewMockExports creates a new mock instance.
func NewMockExports(ctrl *gomock.Controller) *MockExports {
	mock := &MockExports{ctrl: ctrl}
	mock.recorder = &MockExportsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExports) EXPECT() *MockExportsMockRecorder {
	return m.recorder
}

// ProxyGetID mocks base method.
func (m *MockExports) ProxyGetID() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyGetID")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyGetID indicates an expected call of ProxyGetID.
func (mr *MockExportsMockRecorder) ProxyGetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyGetID", reflect.TypeOf((*MockExports)(nil).ProxyGetID))
}

// ProxyOnConfigure mocks base method.
func (m *MockExports) ProxyOnConfigure(rootContextID, pluginConfigurationSize int32) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnConfigure", rootContextID, pluginConfigurationSize)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnConfigure indicates an expected call of ProxyOnConfigure.
func (mr *MockExportsMockRecorder) ProxyOnConfigure(rootContextID, pluginConfigurationSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnConfigure", reflect.TypeOf((*MockExports)(nil).ProxyOnConfigure), rootContextID, pluginConfigurationSize)
}

// ProxyOnContextCreate mocks base method.
func (m *MockExports) ProxyOnContextCreate(contextID, parentContextID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnContextCreate", contextID, parentContextID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnContextCreate indicates an expected call of ProxyOnContextCreate.
func (mr *MockExportsMockRecorder) ProxyOnContextCreate(contextID, parentContextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnContextCreate", reflect.TypeOf((*MockExports)(nil).ProxyOnContextCreate), contextID, parentContextID)
}

// ProxyOnDelete mocks base method.
func (m *MockExports) ProxyOnDelete(contextID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnDelete", contextID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnDelete indicates an expected call of ProxyOnDelete.
func (mr *MockExportsMockRecorder) ProxyOnDelete(contextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnDelete", reflect.TypeOf((*MockExports)(nil).ProxyOnDelete), contextID)
}

// ProxyOnDone mocks base method.
func (m *MockExports) ProxyOnDone(contextID int32) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnDone", contextID)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnDone indicates an expected call of ProxyOnDone.
func (mr *MockExportsMockRecorder) ProxyOnDone(contextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnDone", reflect.TypeOf((*MockExports)(nil).ProxyOnDone), contextID)
}

// ProxyOnDownstreamConnectionClose mocks base method.
func (m *MockExports) ProxyOnDownstreamConnectionClose(contextID, closeType int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnDownstreamConnectionClose", contextID, closeType)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnDownstreamConnectionClose indicates an expected call of ProxyOnDownstreamConnectionClose.
func (mr *MockExportsMockRecorder) ProxyOnDownstreamConnectionClose(contextID, closeType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnDownstreamConnectionClose", reflect.TypeOf((*MockExports)(nil).ProxyOnDownstreamConnectionClose), contextID, closeType)
}

// ProxyOnDownstreamData mocks base method.
func (m *MockExports) ProxyOnDownstreamData(contextID, dataLength, endOfStream int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnDownstreamData", contextID, dataLength, endOfStream)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnDownstreamData indicates an expected call of ProxyOnDownstreamData.
func (mr *MockExportsMockRecorder) ProxyOnDownstreamData(contextID, dataLength, endOfStream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnDownstreamData", reflect.TypeOf((*MockExports)(nil).ProxyOnDownstreamData), contextID, dataLength, endOfStream)
}

// ProxyOnGrpcCallClose mocks base method.
func (m *MockExports) ProxyOnGrpcCallClose(contextID, calloutID, statusCode int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnGrpcCallClose", contextID, calloutID, statusCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnGrpcCallClose indicates an expected call of ProxyOnGrpcCallClose.
func (mr *MockExportsMockRecorder) ProxyOnGrpcCallClose(contextID, calloutID, statusCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnGrpcCallClose", reflect.TypeOf((*MockExports)(nil).ProxyOnGrpcCallClose), contextID, calloutID, statusCode)
}

// ProxyOnGrpcCallResponseHeaderMetadata mocks base method.
func (m *MockExports) ProxyOnGrpcCallResponseHeaderMetadata(contextID, calloutID, nElements int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnGrpcCallResponseHeaderMetadata", contextID, calloutID, nElements)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnGrpcCallResponseHeaderMetadata indicates an expected call of ProxyOnGrpcCallResponseHeaderMetadata.
func (mr *MockExportsMockRecorder) ProxyOnGrpcCallResponseHeaderMetadata(contextID, calloutID, nElements interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnGrpcCallResponseHeaderMetadata", reflect.TypeOf((*MockExports)(nil).ProxyOnGrpcCallResponseHeaderMetadata), contextID, calloutID, nElements)
}

// ProxyOnGrpcCallResponseMessage mocks base method.
func (m *MockExports) ProxyOnGrpcCallResponseMessage(contextID, calloutID, msgSize int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnGrpcCallResponseMessage", contextID, calloutID, msgSize)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnGrpcCallResponseMessage indicates an expected call of ProxyOnGrpcCallResponseMessage.
func (mr *MockExportsMockRecorder) ProxyOnGrpcCallResponseMessage(contextID, calloutID, msgSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnGrpcCallResponseMessage", reflect.TypeOf((*MockExports)(nil).ProxyOnGrpcCallResponseMessage), contextID, calloutID, msgSize)
}

// ProxyOnGrpcCallResponseTrailerMetadata mocks base method.
func (m *MockExports) ProxyOnGrpcCallResponseTrailerMetadata(contextID, calloutID, nElements int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnGrpcCallResponseTrailerMetadata", contextID, calloutID, nElements)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnGrpcCallResponseTrailerMetadata indicates an expected call of ProxyOnGrpcCallResponseTrailerMetadata.
func (mr *MockExportsMockRecorder) ProxyOnGrpcCallResponseTrailerMetadata(contextID, calloutID, nElements interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnGrpcCallResponseTrailerMetadata", reflect.TypeOf((*MockExports)(nil).ProxyOnGrpcCallResponseTrailerMetadata), contextID, calloutID, nElements)
}

// ProxyOnHttpCallResponse mocks base method.
func (m *MockExports) ProxyOnHttpCallResponse(contextID, token, headers, bodySize, trailers int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnHttpCallResponse", contextID, token, headers, bodySize, trailers)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnHttpCallResponse indicates an expected call of ProxyOnHttpCallResponse.
func (mr *MockExportsMockRecorder) ProxyOnHttpCallResponse(contextID, token, headers, bodySize, trailers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnHttpCallResponse", reflect.TypeOf((*MockExports)(nil).ProxyOnHttpCallResponse), contextID, token, headers, bodySize, trailers)
}

// ProxyOnLog mocks base method.
func (m *MockExports) ProxyOnLog(contextID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnLog", contextID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnLog indicates an expected call of ProxyOnLog.
func (mr *MockExportsMockRecorder) ProxyOnLog(contextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnLog", reflect.TypeOf((*MockExports)(nil).ProxyOnLog), contextID)
}

// ProxyOnMemoryAllocate mocks base method.
func (m *MockExports) ProxyOnMemoryAllocate(size int32) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnMemoryAllocate", size)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnMemoryAllocate indicates an expected call of ProxyOnMemoryAllocate.
func (mr *MockExportsMockRecorder) ProxyOnMemoryAllocate(size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnMemoryAllocate", reflect.TypeOf((*MockExports)(nil).ProxyOnMemoryAllocate), size)
}

// ProxyOnNewConnection mocks base method.
func (m *MockExports) ProxyOnNewConnection(contextID int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnNewConnection", contextID)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnNewConnection indicates an expected call of ProxyOnNewConnection.
func (mr *MockExportsMockRecorder) ProxyOnNewConnection(contextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnNewConnection", reflect.TypeOf((*MockExports)(nil).ProxyOnNewConnection), contextID)
}

// ProxyOnQueueReady mocks base method.
func (m *MockExports) ProxyOnQueueReady(rootContextID, token int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnQueueReady", rootContextID, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnQueueReady indicates an expected call of ProxyOnQueueReady.
func (mr *MockExportsMockRecorder) ProxyOnQueueReady(rootContextID, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnQueueReady", reflect.TypeOf((*MockExports)(nil).ProxyOnQueueReady), rootContextID, token)
}

// ProxyOnRequestBody mocks base method.
func (m *MockExports) ProxyOnRequestBody(contextID, bodyBufferLength, endOfStream int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnRequestBody", contextID, bodyBufferLength, endOfStream)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnRequestBody indicates an expected call of ProxyOnRequestBody.
func (mr *MockExportsMockRecorder) ProxyOnRequestBody(contextID, bodyBufferLength, endOfStream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnRequestBody", reflect.TypeOf((*MockExports)(nil).ProxyOnRequestBody), contextID, bodyBufferLength, endOfStream)
}

// ProxyOnRequestHeaders mocks base method.
func (m *MockExports) ProxyOnRequestHeaders(contextID, headers, endOfStream int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnRequestHeaders", contextID, headers, endOfStream)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnRequestHeaders indicates an expected call of ProxyOnRequestHeaders.
func (mr *MockExportsMockRecorder) ProxyOnRequestHeaders(contextID, headers, endOfStream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnRequestHeaders", reflect.TypeOf((*MockExports)(nil).ProxyOnRequestHeaders), contextID, headers, endOfStream)
}

// ProxyOnRequestMetadata mocks base method.
func (m *MockExports) ProxyOnRequestMetadata(contextID, nElements int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnRequestMetadata", contextID, nElements)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnRequestMetadata indicates an expected call of ProxyOnRequestMetadata.
func (mr *MockExportsMockRecorder) ProxyOnRequestMetadata(contextID, nElements interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnRequestMetadata", reflect.TypeOf((*MockExports)(nil).ProxyOnRequestMetadata), contextID, nElements)
}

// ProxyOnRequestTrailers mocks base method.
func (m *MockExports) ProxyOnRequestTrailers(contextID, trailers int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnRequestTrailers", contextID, trailers)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnRequestTrailers indicates an expected call of ProxyOnRequestTrailers.
func (mr *MockExportsMockRecorder) ProxyOnRequestTrailers(contextID, trailers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnRequestTrailers", reflect.TypeOf((*MockExports)(nil).ProxyOnRequestTrailers), contextID, trailers)
}

// ProxyOnResponseBody mocks base method.
func (m *MockExports) ProxyOnResponseBody(contextID, bodyBufferLength, endOfStream int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnResponseBody", contextID, bodyBufferLength, endOfStream)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnResponseBody indicates an expected call of ProxyOnResponseBody.
func (mr *MockExportsMockRecorder) ProxyOnResponseBody(contextID, bodyBufferLength, endOfStream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnResponseBody", reflect.TypeOf((*MockExports)(nil).ProxyOnResponseBody), contextID, bodyBufferLength, endOfStream)
}

// ProxyOnResponseHeaders mocks base method.
func (m *MockExports) ProxyOnResponseHeaders(contextID, headers, endOfStream int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnResponseHeaders", contextID, headers, endOfStream)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnResponseHeaders indicates an expected call of ProxyOnResponseHeaders.
func (mr *MockExportsMockRecorder) ProxyOnResponseHeaders(contextID, headers, endOfStream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnResponseHeaders", reflect.TypeOf((*MockExports)(nil).ProxyOnResponseHeaders), contextID, headers, endOfStream)
}

// ProxyOnResponseMetadata mocks base method.
func (m *MockExports) ProxyOnResponseMetadata(contextID, nElements int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnResponseMetadata", contextID, nElements)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnResponseMetadata indicates an expected call of ProxyOnResponseMetadata.
func (mr *MockExportsMockRecorder) ProxyOnResponseMetadata(contextID, nElements interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnResponseMetadata", reflect.TypeOf((*MockExports)(nil).ProxyOnResponseMetadata), contextID, nElements)
}

// ProxyOnResponseTrailers mocks base method.
func (m *MockExports) ProxyOnResponseTrailers(contextID, trailers int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnResponseTrailers", contextID, trailers)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnResponseTrailers indicates an expected call of ProxyOnResponseTrailers.
func (mr *MockExportsMockRecorder) ProxyOnResponseTrailers(contextID, trailers interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnResponseTrailers", reflect.TypeOf((*MockExports)(nil).ProxyOnResponseTrailers), contextID, trailers)
}

// ProxyOnTick mocks base method.
func (m *MockExports) ProxyOnTick(rootContextID int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnTick", rootContextID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnTick indicates an expected call of ProxyOnTick.
func (mr *MockExportsMockRecorder) ProxyOnTick(rootContextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnTick", reflect.TypeOf((*MockExports)(nil).ProxyOnTick), rootContextID)
}

// ProxyOnUpstreamConnectionClose mocks base method.
func (m *MockExports) ProxyOnUpstreamConnectionClose(contextID, closeType int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnUpstreamConnectionClose", contextID, closeType)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProxyOnUpstreamConnectionClose indicates an expected call of ProxyOnUpstreamConnectionClose.
func (mr *MockExportsMockRecorder) ProxyOnUpstreamConnectionClose(contextID, closeType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnUpstreamConnectionClose", reflect.TypeOf((*MockExports)(nil).ProxyOnUpstreamConnectionClose), contextID, closeType)
}

// ProxyOnUpstreamData mocks base method.
func (m *MockExports) ProxyOnUpstreamData(contextID, dataLength, endOfStream int32) (v1.Action, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnUpstreamData", contextID, dataLength, endOfStream)
	ret0, _ := ret[0].(v1.Action)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnUpstreamData indicates an expected call of ProxyOnUpstreamData.
func (mr *MockExportsMockRecorder) ProxyOnUpstreamData(contextID, dataLength, endOfStream interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnUpstreamData", reflect.TypeOf((*MockExports)(nil).ProxyOnUpstreamData), contextID, dataLength, endOfStream)
}

// ProxyOnVmStart mocks base method.
func (m *MockExports) ProxyOnVmStart(rootContextID, vmConfigurationSize int32) (int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProxyOnVmStart", rootContextID, vmConfigurationSize)
	ret0, _ := ret[0].(int32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ProxyOnVmStart indicates an expected call of ProxyOnVmStart.
func (mr *MockExportsMockRecorder) ProxyOnVmStart(rootContextID, vmConfigurationSize interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProxyOnVmStart", reflect.TypeOf((*MockExports)(nil).ProxyOnVmStart), rootContextID, vmConfigurationSize)
}
