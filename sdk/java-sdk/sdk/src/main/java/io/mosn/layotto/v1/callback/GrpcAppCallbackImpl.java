/*
 * Copyright 2021 Layotto Authors
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package io.mosn.layotto.v1.callback;

import com.google.protobuf.Empty;
import io.grpc.stub.StreamObserver;
import io.mosn.layotto.v1.callback.component.pubsub.PubSub;
import io.mosn.layotto.v1.callback.component.pubsub.PubSubRegistry;
import io.mosn.layotto.v1.grpc.PubsubConverter;
import spec.proto.runtime.v1.AppCallbackGrpc;
import spec.proto.runtime.v1.AppCallbackProto;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponse;
import spec.sdk.runtime.v1.domain.pubsub.TopicSubscription;

import java.util.Collection;
import java.util.Set;

public class GrpcAppCallbackImpl extends AppCallbackGrpc.AppCallbackImplBase {

    private final PubSubRegistry pubSubRegistry;

    public GrpcAppCallbackImpl(PubSubRegistry pubSubRegistry) {
        this.pubSubRegistry = pubSubRegistry;
    }

    /**
     * Iterates all pubsub clients registered to pubSubClientRegistry and get all topic subscriptions.
     */
    @Override
    public void listTopicSubscriptions(Empty request,
                                       StreamObserver<AppCallbackProto.ListTopicSubscriptionsResponse> responseObserver) {
        final AppCallbackProto.ListTopicSubscriptionsResponse.Builder builder
                = AppCallbackProto.ListTopicSubscriptionsResponse.newBuilder();
        // get all PubSub callbacks
        Collection<PubSub> pubsubs = pubSubRegistry.getAllPubSubCallbacks();
        if (pubsubs == null) {
            responseObserver.onNext(builder.build());
            responseObserver.onCompleted();
            return;
        }
        // Iterates them and get all topic subscriptions.
        for (PubSub pubSub : pubsubs) {
            final Set<TopicSubscription> topicSubscriptions = pubSub.listTopicSubscriptions();
            if (topicSubscriptions == null || topicSubscriptions.isEmpty()) {
                continue;
            }
            for (TopicSubscription topicSubscription : topicSubscriptions) {
                if (topicSubscription == null) {
                    continue;
                }
                builder.addSubscriptions(PubsubConverter.TopicSubscription2Grpc(topicSubscription));
            }
        }

        // ack
        responseObserver.onNext(builder.build());
        responseObserver.onCompleted();
    }

    /**
     * On message delivery, find pubsub client by pubsub name.
     */
    @Override
    public void onTopicEvent(AppCallbackProto.TopicEventRequest request,
                             StreamObserver<AppCallbackProto.TopicEventResponse> responseObserver) {
        final String pubsubName = request.getPubsubName();
        // dispatch by pub sub name
        final PubSub pubsub = pubSubRegistry.getCallbackByPubSubName(pubsubName);

        // invoke callback
        final TopicEventResponse response = pubsub.onTopicEvent(PubsubConverter.TopicEventRequest2Domain(request));
        responseObserver.onNext(PubsubConverter.TopicEventResponse2Grpc(response));
        responseObserver.onCompleted();
    }
}
