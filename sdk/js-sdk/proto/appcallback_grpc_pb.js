// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var appcallback_pb = require('./appcallback_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');

function serialize_google_protobuf_Empty(arg) {
  if (!(arg instanceof google_protobuf_empty_pb.Empty)) {
    throw new Error('Expected argument of type google.protobuf.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_google_protobuf_Empty(buffer_arg) {
  return google_protobuf_empty_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_ListTopicSubscriptionsResponse(arg) {
  if (!(arg instanceof appcallback_pb.ListTopicSubscriptionsResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.ListTopicSubscriptionsResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_ListTopicSubscriptionsResponse(buffer_arg) {
  return appcallback_pb.ListTopicSubscriptionsResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_TopicEventRequest(arg) {
  if (!(arg instanceof appcallback_pb.TopicEventRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.TopicEventRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_TopicEventRequest(buffer_arg) {
  return appcallback_pb.TopicEventRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_TopicEventResponse(arg) {
  if (!(arg instanceof appcallback_pb.TopicEventResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.TopicEventResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_TopicEventResponse(buffer_arg) {
  return appcallback_pb.TopicEventResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


// AppCallback V1 allows user application to interact with runtime.
// User application needs to implement AppCallback service if it needs to
// receive message from runtime.
var AppCallbackService = exports.AppCallbackService = {
  // Lists all topics subscribed by this app.
listTopicSubscriptions: {
    path: '/spec.proto.runtime.v1.AppCallback/ListTopicSubscriptions',
    requestStream: false,
    responseStream: false,
    requestType: google_protobuf_empty_pb.Empty,
    responseType: appcallback_pb.ListTopicSubscriptionsResponse,
    requestSerialize: serialize_google_protobuf_Empty,
    requestDeserialize: deserialize_google_protobuf_Empty,
    responseSerialize: serialize_spec_proto_runtime_v1_ListTopicSubscriptionsResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_ListTopicSubscriptionsResponse,
  },
  // Subscribes events from Pubsub
onTopicEvent: {
    path: '/spec.proto.runtime.v1.AppCallback/OnTopicEvent',
    requestStream: false,
    responseStream: false,
    requestType: appcallback_pb.TopicEventRequest,
    responseType: appcallback_pb.TopicEventResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_TopicEventRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_TopicEventRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_TopicEventResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_TopicEventResponse,
  },
};

exports.AppCallbackClient = grpc.makeGenericClientConstructor(AppCallbackService);
//  // Invokes service method with InvokeRequest.
//  rpc OnInvoke (InvokeRequest) returns (InvokeResponse) {}
