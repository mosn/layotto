// GENERATED CODE -- DO NOT EDIT!

'use strict';
var grpc = require('@grpc/grpc-js');
var runtime_pb = require('./runtime_pb.js');
var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js');
var google_protobuf_any_pb = require('google-protobuf/google/protobuf/any_pb.js');

function serialize_google_protobuf_Empty(arg) {
  if (!(arg instanceof google_protobuf_empty_pb.Empty)) {
    throw new Error('Expected argument of type google.protobuf.Empty');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_google_protobuf_Empty(buffer_arg) {
  return google_protobuf_empty_pb.Empty.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_DelFileRequest(arg) {
  if (!(arg instanceof runtime_pb.DelFileRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.DelFileRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_DelFileRequest(buffer_arg) {
  return runtime_pb.DelFileRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_DeleteBulkStateRequest(arg) {
  if (!(arg instanceof runtime_pb.DeleteBulkStateRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.DeleteBulkStateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_DeleteBulkStateRequest(buffer_arg) {
  return runtime_pb.DeleteBulkStateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_DeleteConfigurationRequest(arg) {
  if (!(arg instanceof runtime_pb.DeleteConfigurationRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.DeleteConfigurationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_DeleteConfigurationRequest(buffer_arg) {
  return runtime_pb.DeleteConfigurationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_DeleteStateRequest(arg) {
  if (!(arg instanceof runtime_pb.DeleteStateRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.DeleteStateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_DeleteStateRequest(buffer_arg) {
  return runtime_pb.DeleteStateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_ExecuteStateTransactionRequest(arg) {
  if (!(arg instanceof runtime_pb.ExecuteStateTransactionRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.ExecuteStateTransactionRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_ExecuteStateTransactionRequest(buffer_arg) {
  return runtime_pb.ExecuteStateTransactionRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetBulkStateRequest(arg) {
  if (!(arg instanceof runtime_pb.GetBulkStateRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetBulkStateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetBulkStateRequest(buffer_arg) {
  return runtime_pb.GetBulkStateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetBulkStateResponse(arg) {
  if (!(arg instanceof runtime_pb.GetBulkStateResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetBulkStateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetBulkStateResponse(buffer_arg) {
  return runtime_pb.GetBulkStateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetConfigurationRequest(arg) {
  if (!(arg instanceof runtime_pb.GetConfigurationRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetConfigurationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetConfigurationRequest(buffer_arg) {
  return runtime_pb.GetConfigurationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetConfigurationResponse(arg) {
  if (!(arg instanceof runtime_pb.GetConfigurationResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetConfigurationResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetConfigurationResponse(buffer_arg) {
  return runtime_pb.GetConfigurationResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetFileRequest(arg) {
  if (!(arg instanceof runtime_pb.GetFileRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetFileRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetFileRequest(buffer_arg) {
  return runtime_pb.GetFileRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetFileResponse(arg) {
  if (!(arg instanceof runtime_pb.GetFileResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetFileResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetFileResponse(buffer_arg) {
  return runtime_pb.GetFileResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetNextIdRequest(arg) {
  if (!(arg instanceof runtime_pb.GetNextIdRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetNextIdRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetNextIdRequest(buffer_arg) {
  return runtime_pb.GetNextIdRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetNextIdResponse(arg) {
  if (!(arg instanceof runtime_pb.GetNextIdResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetNextIdResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetNextIdResponse(buffer_arg) {
  return runtime_pb.GetNextIdResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetStateRequest(arg) {
  if (!(arg instanceof runtime_pb.GetStateRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetStateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetStateRequest(buffer_arg) {
  return runtime_pb.GetStateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_GetStateResponse(arg) {
  if (!(arg instanceof runtime_pb.GetStateResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.GetStateResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_GetStateResponse(buffer_arg) {
  return runtime_pb.GetStateResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_InvokeBindingRequest(arg) {
  if (!(arg instanceof runtime_pb.InvokeBindingRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.InvokeBindingRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_InvokeBindingRequest(buffer_arg) {
  return runtime_pb.InvokeBindingRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_InvokeBindingResponse(arg) {
  if (!(arg instanceof runtime_pb.InvokeBindingResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.InvokeBindingResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_InvokeBindingResponse(buffer_arg) {
  return runtime_pb.InvokeBindingResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_InvokeResponse(arg) {
  if (!(arg instanceof runtime_pb.InvokeResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.InvokeResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_InvokeResponse(buffer_arg) {
  return runtime_pb.InvokeResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_InvokeServiceRequest(arg) {
  if (!(arg instanceof runtime_pb.InvokeServiceRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.InvokeServiceRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_InvokeServiceRequest(buffer_arg) {
  return runtime_pb.InvokeServiceRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_ListFileRequest(arg) {
  if (!(arg instanceof runtime_pb.ListFileRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.ListFileRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_ListFileRequest(buffer_arg) {
  return runtime_pb.ListFileRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_ListFileResp(arg) {
  if (!(arg instanceof runtime_pb.ListFileResp)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.ListFileResp');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_ListFileResp(buffer_arg) {
  return runtime_pb.ListFileResp.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_PublishEventRequest(arg) {
  if (!(arg instanceof runtime_pb.PublishEventRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.PublishEventRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_PublishEventRequest(buffer_arg) {
  return runtime_pb.PublishEventRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_PutFileRequest(arg) {
  if (!(arg instanceof runtime_pb.PutFileRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.PutFileRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_PutFileRequest(buffer_arg) {
  return runtime_pb.PutFileRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_SaveConfigurationRequest(arg) {
  if (!(arg instanceof runtime_pb.SaveConfigurationRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.SaveConfigurationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_SaveConfigurationRequest(buffer_arg) {
  return runtime_pb.SaveConfigurationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_SaveStateRequest(arg) {
  if (!(arg instanceof runtime_pb.SaveStateRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.SaveStateRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_SaveStateRequest(buffer_arg) {
  return runtime_pb.SaveStateRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_SayHelloRequest(arg) {
  if (!(arg instanceof runtime_pb.SayHelloRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.SayHelloRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_SayHelloRequest(buffer_arg) {
  return runtime_pb.SayHelloRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_SayHelloResponse(arg) {
  if (!(arg instanceof runtime_pb.SayHelloResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.SayHelloResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_SayHelloResponse(buffer_arg) {
  return runtime_pb.SayHelloResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_SubscribeConfigurationRequest(arg) {
  if (!(arg instanceof runtime_pb.SubscribeConfigurationRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.SubscribeConfigurationRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_SubscribeConfigurationRequest(buffer_arg) {
  return runtime_pb.SubscribeConfigurationRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_SubscribeConfigurationResponse(arg) {
  if (!(arg instanceof runtime_pb.SubscribeConfigurationResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.SubscribeConfigurationResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_SubscribeConfigurationResponse(buffer_arg) {
  return runtime_pb.SubscribeConfigurationResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_TryLockRequest(arg) {
  if (!(arg instanceof runtime_pb.TryLockRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.TryLockRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_TryLockRequest(buffer_arg) {
  return runtime_pb.TryLockRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_TryLockResponse(arg) {
  if (!(arg instanceof runtime_pb.TryLockResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.TryLockResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_TryLockResponse(buffer_arg) {
  return runtime_pb.TryLockResponse.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_UnlockRequest(arg) {
  if (!(arg instanceof runtime_pb.UnlockRequest)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.UnlockRequest');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_UnlockRequest(buffer_arg) {
  return runtime_pb.UnlockRequest.deserializeBinary(new Uint8Array(buffer_arg));
}

function serialize_spec_proto_runtime_v1_UnlockResponse(arg) {
  if (!(arg instanceof runtime_pb.UnlockResponse)) {
    throw new Error('Expected argument of type spec.proto.runtime.v1.UnlockResponse');
  }
  return Buffer.from(arg.serializeBinary());
}

function deserialize_spec_proto_runtime_v1_UnlockResponse(buffer_arg) {
  return runtime_pb.UnlockResponse.deserializeBinary(new Uint8Array(buffer_arg));
}


var RuntimeService = exports.RuntimeService = {
  // SayHello used for test
sayHello: {
    path: '/spec.proto.runtime.v1.Runtime/SayHello',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.SayHelloRequest,
    responseType: runtime_pb.SayHelloResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_SayHelloRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_SayHelloRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_SayHelloResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_SayHelloResponse,
  },
  // InvokeService do rpc calls
invokeService: {
    path: '/spec.proto.runtime.v1.Runtime/InvokeService',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.InvokeServiceRequest,
    responseType: runtime_pb.InvokeResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_InvokeServiceRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_InvokeServiceRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_InvokeResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_InvokeResponse,
  },
  // GetConfiguration gets configuration from configuration store.
getConfiguration: {
    path: '/spec.proto.runtime.v1.Runtime/GetConfiguration',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.GetConfigurationRequest,
    responseType: runtime_pb.GetConfigurationResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_GetConfigurationRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_GetConfigurationRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_GetConfigurationResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_GetConfigurationResponse,
  },
  // SaveConfiguration saves configuration into configuration store.
saveConfiguration: {
    path: '/spec.proto.runtime.v1.Runtime/SaveConfiguration',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.SaveConfigurationRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_SaveConfigurationRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_SaveConfigurationRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // DeleteConfiguration deletes configuration from configuration store.
deleteConfiguration: {
    path: '/spec.proto.runtime.v1.Runtime/DeleteConfiguration',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.DeleteConfigurationRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_DeleteConfigurationRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_DeleteConfigurationRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
subscribeConfiguration: {
    path: '/spec.proto.runtime.v1.Runtime/SubscribeConfiguration',
    requestStream: true,
    responseStream: true,
    requestType: runtime_pb.SubscribeConfigurationRequest,
    responseType: runtime_pb.SubscribeConfigurationResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_SubscribeConfigurationRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_SubscribeConfigurationRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_SubscribeConfigurationResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_SubscribeConfigurationResponse,
  },
  // Distributed Lock API
// A non-blocking method trying to get a lock with ttl.
tryLock: {
    path: '/spec.proto.runtime.v1.Runtime/TryLock',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.TryLockRequest,
    responseType: runtime_pb.TryLockResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_TryLockRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_TryLockRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_TryLockResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_TryLockResponse,
  },
  unlock: {
    path: '/spec.proto.runtime.v1.Runtime/Unlock',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.UnlockRequest,
    responseType: runtime_pb.UnlockResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_UnlockRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_UnlockRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_UnlockResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_UnlockResponse,
  },
  // Sequencer API
// Get next unique id with some auto-increment guarantee
getNextId: {
    path: '/spec.proto.runtime.v1.Runtime/GetNextId',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.GetNextIdRequest,
    responseType: runtime_pb.GetNextIdResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_GetNextIdRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_GetNextIdRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_GetNextIdResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_GetNextIdResponse,
  },
  //  Below are the APIs compatible with Dapr.
//  We try to keep them same as Dapr's because we want to work with Dapr to build an API spec for cloud native runtime
//  ,like CloudEvent for event data.
//
// Gets the state for a specific key.
getState: {
    path: '/spec.proto.runtime.v1.Runtime/GetState',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.GetStateRequest,
    responseType: runtime_pb.GetStateResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_GetStateRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_GetStateRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_GetStateResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_GetStateResponse,
  },
  // Gets a bulk of state items for a list of keys
getBulkState: {
    path: '/spec.proto.runtime.v1.Runtime/GetBulkState',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.GetBulkStateRequest,
    responseType: runtime_pb.GetBulkStateResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_GetBulkStateRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_GetBulkStateRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_GetBulkStateResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_GetBulkStateResponse,
  },
  // Saves an array of state objects
saveState: {
    path: '/spec.proto.runtime.v1.Runtime/SaveState',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.SaveStateRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_SaveStateRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_SaveStateRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // Deletes the state for a specific key.
deleteState: {
    path: '/spec.proto.runtime.v1.Runtime/DeleteState',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.DeleteStateRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_DeleteStateRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_DeleteStateRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // Deletes a bulk of state items for a list of keys
deleteBulkState: {
    path: '/spec.proto.runtime.v1.Runtime/DeleteBulkState',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.DeleteBulkStateRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_DeleteBulkStateRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_DeleteBulkStateRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // Executes transactions for a specified store
executeStateTransaction: {
    path: '/spec.proto.runtime.v1.Runtime/ExecuteStateTransaction',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.ExecuteStateTransactionRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_ExecuteStateTransactionRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_ExecuteStateTransactionRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // Publishes events to the specific topic
publishEvent: {
    path: '/spec.proto.runtime.v1.Runtime/PublishEvent',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.PublishEventRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_PublishEventRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_PublishEventRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // Get file with stream
getFile: {
    path: '/spec.proto.runtime.v1.Runtime/GetFile',
    requestStream: false,
    responseStream: true,
    requestType: runtime_pb.GetFileRequest,
    responseType: runtime_pb.GetFileResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_GetFileRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_GetFileRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_GetFileResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_GetFileResponse,
  },
  // Put file with stream
putFile: {
    path: '/spec.proto.runtime.v1.Runtime/PutFile',
    requestStream: true,
    responseStream: false,
    requestType: runtime_pb.PutFileRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_PutFileRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_PutFileRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // List all files
listFile: {
    path: '/spec.proto.runtime.v1.Runtime/ListFile',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.ListFileRequest,
    responseType: runtime_pb.ListFileResp,
    requestSerialize: serialize_spec_proto_runtime_v1_ListFileRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_ListFileRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_ListFileResp,
    responseDeserialize: deserialize_spec_proto_runtime_v1_ListFileResp,
  },
  // Delete specific file
delFile: {
    path: '/spec.proto.runtime.v1.Runtime/DelFile',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.DelFileRequest,
    responseType: google_protobuf_empty_pb.Empty,
    requestSerialize: serialize_spec_proto_runtime_v1_DelFileRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_DelFileRequest,
    responseSerialize: serialize_google_protobuf_Empty,
    responseDeserialize: deserialize_google_protobuf_Empty,
  },
  // Invokes binding data to specific output bindings
invokeBinding: {
    path: '/spec.proto.runtime.v1.Runtime/InvokeBinding',
    requestStream: false,
    responseStream: false,
    requestType: runtime_pb.InvokeBindingRequest,
    responseType: runtime_pb.InvokeBindingResponse,
    requestSerialize: serialize_spec_proto_runtime_v1_InvokeBindingRequest,
    requestDeserialize: deserialize_spec_proto_runtime_v1_InvokeBindingRequest,
    responseSerialize: serialize_spec_proto_runtime_v1_InvokeBindingResponse,
    responseDeserialize: deserialize_spec_proto_runtime_v1_InvokeBindingResponse,
  },
};

exports.RuntimeClient = grpc.makeGenericClientConstructor(RuntimeService);
