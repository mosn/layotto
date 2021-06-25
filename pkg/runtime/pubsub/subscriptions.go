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

package pubsub

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	runtimev1pb "mosn.io/layotto/spec/proto/runtime/v1"
	"mosn.io/pkg/log"
)

func ListTopicSubscriptions(client runtimev1pb.AppCallbackClient, log log.ErrorLogger) []*runtimev1pb.TopicSubscription {
	resp, err := client.ListTopicSubscriptions(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Errorf("[runtime][ListTopicSubscriptions]error after callback: %s", err)
		return make([]*runtimev1pb.TopicSubscription, 0)
	}
	if resp != nil && len(resp.Subscriptions) > 0 {
		return resp.Subscriptions
	}
	return make([]*runtimev1pb.TopicSubscription, 0)
}
