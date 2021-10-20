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
        Collection<PubSub> pubsubs = pubSubRegistry.getAllPubSubCallbacks();
        if (pubsubs == null) {
            responseObserver.onNext(builder.build());
            responseObserver.onCompleted();
            return;
        }
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
