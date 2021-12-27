/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package default_api

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/pkg/log"
	"testing"
)

type mockClient struct {
}

func (m *mockClient) ListTopicSubscriptions(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*runtimev1pb.ListTopicSubscriptionsResponse, error) {
	return nil, fmt.Errorf("mock failure")
}

func (m *mockClient) OnTopicEvent(ctx context.Context, in *runtimev1pb.TopicEventRequest, opts ...grpc.CallOption) (*runtimev1pb.TopicEventResponse, error) {
	panic("implement me")
}

func Test_listTopicSubscriptions(t *testing.T) {
	topics := listTopicSubscriptions(&mockClient{}, log.DefaultLogger)
	assert.True(t, topics != nil && len(topics) == 0)
}
