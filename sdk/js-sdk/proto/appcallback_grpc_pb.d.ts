// package: spec.proto.runtime.v1
// file: appcallback.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as appcallback_pb from "./appcallback_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

interface IAppCallbackService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    listTopicSubscriptions: IAppCallbackService_IListTopicSubscriptions;
    onTopicEvent: IAppCallbackService_IOnTopicEvent;
}

interface IAppCallbackService_IListTopicSubscriptions extends grpc.MethodDefinition<google_protobuf_empty_pb.Empty, appcallback_pb.ListTopicSubscriptionsResponse> {
    path: "/spec.proto.runtime.v1.AppCallback/ListTopicSubscriptions";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    requestDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
    responseSerialize: grpc.serialize<appcallback_pb.ListTopicSubscriptionsResponse>;
    responseDeserialize: grpc.deserialize<appcallback_pb.ListTopicSubscriptionsResponse>;
}
interface IAppCallbackService_IOnTopicEvent extends grpc.MethodDefinition<appcallback_pb.TopicEventRequest, appcallback_pb.TopicEventResponse> {
    path: "/spec.proto.runtime.v1.AppCallback/OnTopicEvent";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<appcallback_pb.TopicEventRequest>;
    requestDeserialize: grpc.deserialize<appcallback_pb.TopicEventRequest>;
    responseSerialize: grpc.serialize<appcallback_pb.TopicEventResponse>;
    responseDeserialize: grpc.deserialize<appcallback_pb.TopicEventResponse>;
}

export const AppCallbackService: IAppCallbackService;

export interface IAppCallbackServer extends grpc.UntypedServiceImplementation {
    listTopicSubscriptions: grpc.handleUnaryCall<google_protobuf_empty_pb.Empty, appcallback_pb.ListTopicSubscriptionsResponse>;
    onTopicEvent: grpc.handleUnaryCall<appcallback_pb.TopicEventRequest, appcallback_pb.TopicEventResponse>;
}

export interface IAppCallbackClient {
    listTopicSubscriptions(request: google_protobuf_empty_pb.Empty, callback: (error: grpc.ServiceError | null, response: appcallback_pb.ListTopicSubscriptionsResponse) => void): grpc.ClientUnaryCall;
    listTopicSubscriptions(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: appcallback_pb.ListTopicSubscriptionsResponse) => void): grpc.ClientUnaryCall;
    listTopicSubscriptions(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: appcallback_pb.ListTopicSubscriptionsResponse) => void): grpc.ClientUnaryCall;
    onTopicEvent(request: appcallback_pb.TopicEventRequest, callback: (error: grpc.ServiceError | null, response: appcallback_pb.TopicEventResponse) => void): grpc.ClientUnaryCall;
    onTopicEvent(request: appcallback_pb.TopicEventRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: appcallback_pb.TopicEventResponse) => void): grpc.ClientUnaryCall;
    onTopicEvent(request: appcallback_pb.TopicEventRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: appcallback_pb.TopicEventResponse) => void): grpc.ClientUnaryCall;
}

export class AppCallbackClient extends grpc.Client implements IAppCallbackClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public listTopicSubscriptions(request: google_protobuf_empty_pb.Empty, callback: (error: grpc.ServiceError | null, response: appcallback_pb.ListTopicSubscriptionsResponse) => void): grpc.ClientUnaryCall;
    public listTopicSubscriptions(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: appcallback_pb.ListTopicSubscriptionsResponse) => void): grpc.ClientUnaryCall;
    public listTopicSubscriptions(request: google_protobuf_empty_pb.Empty, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: appcallback_pb.ListTopicSubscriptionsResponse) => void): grpc.ClientUnaryCall;
    public onTopicEvent(request: appcallback_pb.TopicEventRequest, callback: (error: grpc.ServiceError | null, response: appcallback_pb.TopicEventResponse) => void): grpc.ClientUnaryCall;
    public onTopicEvent(request: appcallback_pb.TopicEventRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: appcallback_pb.TopicEventResponse) => void): grpc.ClientUnaryCall;
    public onTopicEvent(request: appcallback_pb.TopicEventRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: appcallback_pb.TopicEventResponse) => void): grpc.ClientUnaryCall;
}
