package mock_appcallback

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
	empty "google.golang.org/protobuf/types/known/emptypb"

	dapr_common_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/common/v1"
	dapr_v1pb "mosn.io/layotto/pkg/grpc/dapr/proto/runtime/v1"
)

type MockDaprAppCallbackServer struct {
	ctrl     *gomock.Controller
	recorder *MockDaprAppCallbackServerMockRecorder
}

type MockDaprAppCallbackServerMockRecorder struct {
	mock *MockDaprAppCallbackServer
}

func (m *MockDaprAppCallbackServer) OnInvoke(ctx context.Context, in *dapr_common_v1pb.InvokeRequest) (*dapr_common_v1pb.InvokeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTopicSubscriptions", ctx, in)
	ret0, _ := ret[0].(*dapr_common_v1pb.InvokeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockDaprAppCallbackServer) ListInputBindings(ctx context.Context, in *empty.Empty) (*dapr_v1pb.ListInputBindingsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTopicSubscriptions", ctx, in)
	ret0, _ := ret[0].(*dapr_v1pb.ListInputBindingsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockDaprAppCallbackServer) OnBindingEvent(ctx context.Context, in *dapr_v1pb.BindingEventRequest) (*dapr_v1pb.BindingEventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTopicSubscriptions", ctx, in)
	ret0, _ := ret[0].(*dapr_v1pb.BindingEventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockDaprAppCallbackServer) ListTopicSubscriptions(arg0 context.Context, arg1 *empty.Empty) (*dapr_v1pb.ListTopicSubscriptionsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTopicSubscriptions", arg0, arg1)
	ret0, _ := ret[0].(*dapr_v1pb.ListTopicSubscriptionsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (m *MockDaprAppCallbackServer) OnTopicEvent(ctx context.Context, in *dapr_v1pb.TopicEventRequest) (*dapr_v1pb.TopicEventResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OnTopicEvent", ctx, in)
	ret0, _ := ret[0].(*dapr_v1pb.TopicEventResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockDaprAppCallbackServerMockRecorder) OnInvoke(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnInvoke", reflect.TypeOf((*MockDaprAppCallbackServer)(nil).OnInvoke), arg0, arg1)
}

func (mr *MockDaprAppCallbackServerMockRecorder) ListInputBindings(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListInputBindings", reflect.TypeOf((*MockDaprAppCallbackServer)(nil).ListInputBindings), arg0, arg1)
}

func (mr *MockDaprAppCallbackServerMockRecorder) OnBindingEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnBindingEvent", reflect.TypeOf((*MockDaprAppCallbackServer)(nil).OnBindingEvent), arg0, arg1)
}

func (mr *MockDaprAppCallbackServerMockRecorder) ListTopicSubscriptions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTopicSubscriptions", reflect.TypeOf((*MockDaprAppCallbackServer)(nil).ListTopicSubscriptions), arg0, arg1)
}

func (mr *MockDaprAppCallbackServerMockRecorder) OnTopicEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OnTopicEvent", reflect.TypeOf((*MockDaprAppCallbackServer)(nil).OnTopicEvent), arg0, arg1)
}

func NewMockDaprAppCallbackServer(ctrl *gomock.Controller) *MockDaprAppCallbackServer {
	mock := &MockDaprAppCallbackServer{ctrl: ctrl}
	mock.recorder = &MockDaprAppCallbackServerMockRecorder{mock}
	return mock
}

func (m *MockDaprAppCallbackServer) EXPECT() *MockDaprAppCallbackServerMockRecorder {
	return m.recorder
}
