package io.mosn.layotto.v1.grpc;

import com.google.protobuf.ByteString;
import spec.sdk.runtime.v1.domain.pubsub.PublishEventRequest;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventRequest;
import spec.sdk.runtime.v1.domain.pubsub.TopicEventResponse;
import spec.sdk.runtime.v1.domain.pubsub.TopicSubscription;
import spec.proto.runtime.v1.AppCallbackProto;
import spec.proto.runtime.v1.RuntimeProto;

/**
 * pubsub related converter
 */
public class PubsubConverter {
    public static RuntimeProto.PublishEventRequest PublishEventRequest2Grpc(PublishEventRequest req) {
        if (req == null) {
            return null;
        }
        RuntimeProto.PublishEventRequest result = RuntimeProto.PublishEventRequest.newBuilder()
                .setPubsubName(req.getPubsubName())
                .setTopic(req.getTopic())
                .setData(ByteString.copyFrom(req.getData()))
                .setDataContentType(req.getContentType())
                .putAllMetadata(req.getMetadata())
                .build();

        return result;
    }

    public static AppCallbackProto.TopicSubscription TopicSubscription2Grpc(TopicSubscription sub) {
        if (sub == null) {
            return null;
        }
        AppCallbackProto.TopicSubscription result = AppCallbackProto.TopicSubscription.newBuilder()
                .setPubsubName(sub.getPubsubName())
                .setTopic(sub.getTopic())
                .putAllMetadata(sub.getMetadata())
                .build();

        return result;
    }

    public static AppCallbackProto.TopicEventResponse TopicEventResponse2Grpc(TopicEventResponse resp) {
        if (resp == null) {
            return null;
        }
        int idx = 0;
        if (resp.getStatus() != null) {
            idx = resp.getStatus().getIdx();
        }
        AppCallbackProto.TopicEventResponse result = AppCallbackProto.TopicEventResponse.newBuilder()
                .setStatusValue(idx)
                .build();

        return result;
    }

    public static AppCallbackProto.TopicEventRequest TopicEventRequest2Grpc(TopicEventRequest req) {
        if (req == null) {
            return null;
        }
        AppCallbackProto.TopicEventRequest result = AppCallbackProto.TopicEventRequest.newBuilder()
                .setId(req.getId())
                .setSource(req.getSource())
                .setType(req.getType())
                .setSpecVersion(req.getSpecVersion())
                .setDataContentType(req.getContentType())
                .setData(ByteString.copyFrom(req.getData()))
                .setTopic(req.getTopic())
                .setPubsubName(req.getPubsubName())
                //.putAllMetadata(req.getMetadata())
                .build();

        return result;
    }

    public static TopicEventRequest TopicEventRequest2Domain(AppCallbackProto.TopicEventRequest req) {
        TopicEventRequest result = new TopicEventRequest();
        result.setId(req.getId());
        result.setSource(req.getSource());
        result.setType(req.getType());
        result.setSpecVersion(req.getSpecVersion());
        result.setContentType(req.getDataContentType());
        result.setData(req.getData().toByteArray());
        result.setTopic(req.getTopic());
        result.setPubsubName(req.getPubsubName());
        //result.setMetadata(req.getMetadataMap());
        return result;
    }
}
