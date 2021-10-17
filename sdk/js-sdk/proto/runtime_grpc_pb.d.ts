// package: spec.proto.runtime.v1
// file: runtime.proto

/* tslint:disable */
/* eslint-disable */

import * as grpc from "@grpc/grpc-js";
import * as runtime_pb from "./runtime_pb";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";

interface IRuntimeService extends grpc.ServiceDefinition<grpc.UntypedServiceImplementation> {
    sayHello: IRuntimeService_ISayHello;
    invokeService: IRuntimeService_IInvokeService;
    getConfiguration: IRuntimeService_IGetConfiguration;
    saveConfiguration: IRuntimeService_ISaveConfiguration;
    deleteConfiguration: IRuntimeService_IDeleteConfiguration;
    subscribeConfiguration: IRuntimeService_ISubscribeConfiguration;
    tryLock: IRuntimeService_ITryLock;
    unlock: IRuntimeService_IUnlock;
    getNextId: IRuntimeService_IGetNextId;
    getState: IRuntimeService_IGetState;
    getBulkState: IRuntimeService_IGetBulkState;
    saveState: IRuntimeService_ISaveState;
    deleteState: IRuntimeService_IDeleteState;
    deleteBulkState: IRuntimeService_IDeleteBulkState;
    executeStateTransaction: IRuntimeService_IExecuteStateTransaction;
    publishEvent: IRuntimeService_IPublishEvent;
    getFile: IRuntimeService_IGetFile;
    putFile: IRuntimeService_IPutFile;
    listFile: IRuntimeService_IListFile;
    delFile: IRuntimeService_IDelFile;
    invokeBinding: IRuntimeService_IInvokeBinding;
}

interface IRuntimeService_ISayHello extends grpc.MethodDefinition<runtime_pb.SayHelloRequest, runtime_pb.SayHelloResponse> {
    path: "/spec.proto.runtime.v1.Runtime/SayHello";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.SayHelloRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.SayHelloRequest>;
    responseSerialize: grpc.serialize<runtime_pb.SayHelloResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.SayHelloResponse>;
}
interface IRuntimeService_IInvokeService extends grpc.MethodDefinition<runtime_pb.InvokeServiceRequest, runtime_pb.InvokeResponse> {
    path: "/spec.proto.runtime.v1.Runtime/InvokeService";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.InvokeServiceRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.InvokeServiceRequest>;
    responseSerialize: grpc.serialize<runtime_pb.InvokeResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.InvokeResponse>;
}
interface IRuntimeService_IGetConfiguration extends grpc.MethodDefinition<runtime_pb.GetConfigurationRequest, runtime_pb.GetConfigurationResponse> {
    path: "/spec.proto.runtime.v1.Runtime/GetConfiguration";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.GetConfigurationRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.GetConfigurationRequest>;
    responseSerialize: grpc.serialize<runtime_pb.GetConfigurationResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.GetConfigurationResponse>;
}
interface IRuntimeService_ISaveConfiguration extends grpc.MethodDefinition<runtime_pb.SaveConfigurationRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/SaveConfiguration";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.SaveConfigurationRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.SaveConfigurationRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IDeleteConfiguration extends grpc.MethodDefinition<runtime_pb.DeleteConfigurationRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/DeleteConfiguration";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.DeleteConfigurationRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.DeleteConfigurationRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_ISubscribeConfiguration extends grpc.MethodDefinition<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse> {
    path: "/spec.proto.runtime.v1.Runtime/SubscribeConfiguration";
    requestStream: true;
    responseStream: true;
    requestSerialize: grpc.serialize<runtime_pb.SubscribeConfigurationRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.SubscribeConfigurationRequest>;
    responseSerialize: grpc.serialize<runtime_pb.SubscribeConfigurationResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.SubscribeConfigurationResponse>;
}
interface IRuntimeService_ITryLock extends grpc.MethodDefinition<runtime_pb.TryLockRequest, runtime_pb.TryLockResponse> {
    path: "/spec.proto.runtime.v1.Runtime/TryLock";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.TryLockRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.TryLockRequest>;
    responseSerialize: grpc.serialize<runtime_pb.TryLockResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.TryLockResponse>;
}
interface IRuntimeService_IUnlock extends grpc.MethodDefinition<runtime_pb.UnlockRequest, runtime_pb.UnlockResponse> {
    path: "/spec.proto.runtime.v1.Runtime/Unlock";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.UnlockRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.UnlockRequest>;
    responseSerialize: grpc.serialize<runtime_pb.UnlockResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.UnlockResponse>;
}
interface IRuntimeService_IGetNextId extends grpc.MethodDefinition<runtime_pb.GetNextIdRequest, runtime_pb.GetNextIdResponse> {
    path: "/spec.proto.runtime.v1.Runtime/GetNextId";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.GetNextIdRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.GetNextIdRequest>;
    responseSerialize: grpc.serialize<runtime_pb.GetNextIdResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.GetNextIdResponse>;
}
interface IRuntimeService_IGetState extends grpc.MethodDefinition<runtime_pb.GetStateRequest, runtime_pb.GetStateResponse> {
    path: "/spec.proto.runtime.v1.Runtime/GetState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.GetStateRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.GetStateRequest>;
    responseSerialize: grpc.serialize<runtime_pb.GetStateResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.GetStateResponse>;
}
interface IRuntimeService_IGetBulkState extends grpc.MethodDefinition<runtime_pb.GetBulkStateRequest, runtime_pb.GetBulkStateResponse> {
    path: "/spec.proto.runtime.v1.Runtime/GetBulkState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.GetBulkStateRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.GetBulkStateRequest>;
    responseSerialize: grpc.serialize<runtime_pb.GetBulkStateResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.GetBulkStateResponse>;
}
interface IRuntimeService_ISaveState extends grpc.MethodDefinition<runtime_pb.SaveStateRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/SaveState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.SaveStateRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.SaveStateRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IDeleteState extends grpc.MethodDefinition<runtime_pb.DeleteStateRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/DeleteState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.DeleteStateRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.DeleteStateRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IDeleteBulkState extends grpc.MethodDefinition<runtime_pb.DeleteBulkStateRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/DeleteBulkState";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.DeleteBulkStateRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.DeleteBulkStateRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IExecuteStateTransaction extends grpc.MethodDefinition<runtime_pb.ExecuteStateTransactionRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/ExecuteStateTransaction";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.ExecuteStateTransactionRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.ExecuteStateTransactionRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IPublishEvent extends grpc.MethodDefinition<runtime_pb.PublishEventRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/PublishEvent";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.PublishEventRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.PublishEventRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IGetFile extends grpc.MethodDefinition<runtime_pb.GetFileRequest, runtime_pb.GetFileResponse> {
    path: "/spec.proto.runtime.v1.Runtime/GetFile";
    requestStream: false;
    responseStream: true;
    requestSerialize: grpc.serialize<runtime_pb.GetFileRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.GetFileRequest>;
    responseSerialize: grpc.serialize<runtime_pb.GetFileResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.GetFileResponse>;
}
interface IRuntimeService_IPutFile extends grpc.MethodDefinition<runtime_pb.PutFileRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/PutFile";
    requestStream: true;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.PutFileRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.PutFileRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IListFile extends grpc.MethodDefinition<runtime_pb.ListFileRequest, runtime_pb.ListFileResp> {
    path: "/spec.proto.runtime.v1.Runtime/ListFile";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.ListFileRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.ListFileRequest>;
    responseSerialize: grpc.serialize<runtime_pb.ListFileResp>;
    responseDeserialize: grpc.deserialize<runtime_pb.ListFileResp>;
}
interface IRuntimeService_IDelFile extends grpc.MethodDefinition<runtime_pb.DelFileRequest, google_protobuf_empty_pb.Empty> {
    path: "/spec.proto.runtime.v1.Runtime/DelFile";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.DelFileRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.DelFileRequest>;
    responseSerialize: grpc.serialize<google_protobuf_empty_pb.Empty>;
    responseDeserialize: grpc.deserialize<google_protobuf_empty_pb.Empty>;
}
interface IRuntimeService_IInvokeBinding extends grpc.MethodDefinition<runtime_pb.InvokeBindingRequest, runtime_pb.InvokeBindingResponse> {
    path: "/spec.proto.runtime.v1.Runtime/InvokeBinding";
    requestStream: false;
    responseStream: false;
    requestSerialize: grpc.serialize<runtime_pb.InvokeBindingRequest>;
    requestDeserialize: grpc.deserialize<runtime_pb.InvokeBindingRequest>;
    responseSerialize: grpc.serialize<runtime_pb.InvokeBindingResponse>;
    responseDeserialize: grpc.deserialize<runtime_pb.InvokeBindingResponse>;
}

export const RuntimeService: IRuntimeService;

export interface IRuntimeServer extends grpc.UntypedServiceImplementation {
    sayHello: grpc.handleUnaryCall<runtime_pb.SayHelloRequest, runtime_pb.SayHelloResponse>;
    invokeService: grpc.handleUnaryCall<runtime_pb.InvokeServiceRequest, runtime_pb.InvokeResponse>;
    getConfiguration: grpc.handleUnaryCall<runtime_pb.GetConfigurationRequest, runtime_pb.GetConfigurationResponse>;
    saveConfiguration: grpc.handleUnaryCall<runtime_pb.SaveConfigurationRequest, google_protobuf_empty_pb.Empty>;
    deleteConfiguration: grpc.handleUnaryCall<runtime_pb.DeleteConfigurationRequest, google_protobuf_empty_pb.Empty>;
    subscribeConfiguration: grpc.handleBidiStreamingCall<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse>;
    tryLock: grpc.handleUnaryCall<runtime_pb.TryLockRequest, runtime_pb.TryLockResponse>;
    unlock: grpc.handleUnaryCall<runtime_pb.UnlockRequest, runtime_pb.UnlockResponse>;
    getNextId: grpc.handleUnaryCall<runtime_pb.GetNextIdRequest, runtime_pb.GetNextIdResponse>;
    getState: grpc.handleUnaryCall<runtime_pb.GetStateRequest, runtime_pb.GetStateResponse>;
    getBulkState: grpc.handleUnaryCall<runtime_pb.GetBulkStateRequest, runtime_pb.GetBulkStateResponse>;
    saveState: grpc.handleUnaryCall<runtime_pb.SaveStateRequest, google_protobuf_empty_pb.Empty>;
    deleteState: grpc.handleUnaryCall<runtime_pb.DeleteStateRequest, google_protobuf_empty_pb.Empty>;
    deleteBulkState: grpc.handleUnaryCall<runtime_pb.DeleteBulkStateRequest, google_protobuf_empty_pb.Empty>;
    executeStateTransaction: grpc.handleUnaryCall<runtime_pb.ExecuteStateTransactionRequest, google_protobuf_empty_pb.Empty>;
    publishEvent: grpc.handleUnaryCall<runtime_pb.PublishEventRequest, google_protobuf_empty_pb.Empty>;
    getFile: grpc.handleServerStreamingCall<runtime_pb.GetFileRequest, runtime_pb.GetFileResponse>;
    putFile: grpc.handleClientStreamingCall<runtime_pb.PutFileRequest, google_protobuf_empty_pb.Empty>;
    listFile: grpc.handleUnaryCall<runtime_pb.ListFileRequest, runtime_pb.ListFileResp>;
    delFile: grpc.handleUnaryCall<runtime_pb.DelFileRequest, google_protobuf_empty_pb.Empty>;
    invokeBinding: grpc.handleUnaryCall<runtime_pb.InvokeBindingRequest, runtime_pb.InvokeBindingResponse>;
}

export interface IRuntimeClient {
    sayHello(request: runtime_pb.SayHelloRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.SayHelloResponse) => void): grpc.ClientUnaryCall;
    sayHello(request: runtime_pb.SayHelloRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.SayHelloResponse) => void): grpc.ClientUnaryCall;
    sayHello(request: runtime_pb.SayHelloRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.SayHelloResponse) => void): grpc.ClientUnaryCall;
    invokeService(request: runtime_pb.InvokeServiceRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeResponse) => void): grpc.ClientUnaryCall;
    invokeService(request: runtime_pb.InvokeServiceRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeResponse) => void): grpc.ClientUnaryCall;
    invokeService(request: runtime_pb.InvokeServiceRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeResponse) => void): grpc.ClientUnaryCall;
    getConfiguration(request: runtime_pb.GetConfigurationRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetConfigurationResponse) => void): grpc.ClientUnaryCall;
    getConfiguration(request: runtime_pb.GetConfigurationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetConfigurationResponse) => void): grpc.ClientUnaryCall;
    getConfiguration(request: runtime_pb.GetConfigurationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetConfigurationResponse) => void): grpc.ClientUnaryCall;
    saveConfiguration(request: runtime_pb.SaveConfigurationRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    saveConfiguration(request: runtime_pb.SaveConfigurationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    saveConfiguration(request: runtime_pb.SaveConfigurationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteConfiguration(request: runtime_pb.DeleteConfigurationRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteConfiguration(request: runtime_pb.DeleteConfigurationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteConfiguration(request: runtime_pb.DeleteConfigurationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    subscribeConfiguration(): grpc.ClientDuplexStream<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse>;
    subscribeConfiguration(options: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse>;
    subscribeConfiguration(metadata: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse>;
    tryLock(request: runtime_pb.TryLockRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.TryLockResponse) => void): grpc.ClientUnaryCall;
    tryLock(request: runtime_pb.TryLockRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.TryLockResponse) => void): grpc.ClientUnaryCall;
    tryLock(request: runtime_pb.TryLockRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.TryLockResponse) => void): grpc.ClientUnaryCall;
    unlock(request: runtime_pb.UnlockRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.UnlockResponse) => void): grpc.ClientUnaryCall;
    unlock(request: runtime_pb.UnlockRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.UnlockResponse) => void): grpc.ClientUnaryCall;
    unlock(request: runtime_pb.UnlockRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.UnlockResponse) => void): grpc.ClientUnaryCall;
    getNextId(request: runtime_pb.GetNextIdRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetNextIdResponse) => void): grpc.ClientUnaryCall;
    getNextId(request: runtime_pb.GetNextIdRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetNextIdResponse) => void): grpc.ClientUnaryCall;
    getNextId(request: runtime_pb.GetNextIdRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetNextIdResponse) => void): grpc.ClientUnaryCall;
    getState(request: runtime_pb.GetStateRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetStateResponse) => void): grpc.ClientUnaryCall;
    getState(request: runtime_pb.GetStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetStateResponse) => void): grpc.ClientUnaryCall;
    getState(request: runtime_pb.GetStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetStateResponse) => void): grpc.ClientUnaryCall;
    getBulkState(request: runtime_pb.GetBulkStateRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetBulkStateResponse) => void): grpc.ClientUnaryCall;
    getBulkState(request: runtime_pb.GetBulkStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetBulkStateResponse) => void): grpc.ClientUnaryCall;
    getBulkState(request: runtime_pb.GetBulkStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetBulkStateResponse) => void): grpc.ClientUnaryCall;
    saveState(request: runtime_pb.SaveStateRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    saveState(request: runtime_pb.SaveStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    saveState(request: runtime_pb.SaveStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteState(request: runtime_pb.DeleteStateRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteState(request: runtime_pb.DeleteStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteState(request: runtime_pb.DeleteStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteBulkState(request: runtime_pb.DeleteBulkStateRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteBulkState(request: runtime_pb.DeleteBulkStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    deleteBulkState(request: runtime_pb.DeleteBulkStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    executeStateTransaction(request: runtime_pb.ExecuteStateTransactionRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    executeStateTransaction(request: runtime_pb.ExecuteStateTransactionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    executeStateTransaction(request: runtime_pb.ExecuteStateTransactionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    publishEvent(request: runtime_pb.PublishEventRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    publishEvent(request: runtime_pb.PublishEventRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    publishEvent(request: runtime_pb.PublishEventRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    getFile(request: runtime_pb.GetFileRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<runtime_pb.GetFileResponse>;
    getFile(request: runtime_pb.GetFileRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<runtime_pb.GetFileResponse>;
    putFile(callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    putFile(metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    putFile(options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    putFile(metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    listFile(request: runtime_pb.ListFileRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.ListFileResp) => void): grpc.ClientUnaryCall;
    listFile(request: runtime_pb.ListFileRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.ListFileResp) => void): grpc.ClientUnaryCall;
    listFile(request: runtime_pb.ListFileRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.ListFileResp) => void): grpc.ClientUnaryCall;
    delFile(request: runtime_pb.DelFileRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    delFile(request: runtime_pb.DelFileRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    delFile(request: runtime_pb.DelFileRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    invokeBinding(request: runtime_pb.InvokeBindingRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeBindingResponse) => void): grpc.ClientUnaryCall;
    invokeBinding(request: runtime_pb.InvokeBindingRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeBindingResponse) => void): grpc.ClientUnaryCall;
    invokeBinding(request: runtime_pb.InvokeBindingRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeBindingResponse) => void): grpc.ClientUnaryCall;
}

export class RuntimeClient extends grpc.Client implements IRuntimeClient {
    constructor(address: string, credentials: grpc.ChannelCredentials, options?: Partial<grpc.ClientOptions>);
    public sayHello(request: runtime_pb.SayHelloRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.SayHelloResponse) => void): grpc.ClientUnaryCall;
    public sayHello(request: runtime_pb.SayHelloRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.SayHelloResponse) => void): grpc.ClientUnaryCall;
    public sayHello(request: runtime_pb.SayHelloRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.SayHelloResponse) => void): grpc.ClientUnaryCall;
    public invokeService(request: runtime_pb.InvokeServiceRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeResponse) => void): grpc.ClientUnaryCall;
    public invokeService(request: runtime_pb.InvokeServiceRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeResponse) => void): grpc.ClientUnaryCall;
    public invokeService(request: runtime_pb.InvokeServiceRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeResponse) => void): grpc.ClientUnaryCall;
    public getConfiguration(request: runtime_pb.GetConfigurationRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetConfigurationResponse) => void): grpc.ClientUnaryCall;
    public getConfiguration(request: runtime_pb.GetConfigurationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetConfigurationResponse) => void): grpc.ClientUnaryCall;
    public getConfiguration(request: runtime_pb.GetConfigurationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetConfigurationResponse) => void): grpc.ClientUnaryCall;
    public saveConfiguration(request: runtime_pb.SaveConfigurationRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public saveConfiguration(request: runtime_pb.SaveConfigurationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public saveConfiguration(request: runtime_pb.SaveConfigurationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteConfiguration(request: runtime_pb.DeleteConfigurationRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteConfiguration(request: runtime_pb.DeleteConfigurationRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteConfiguration(request: runtime_pb.DeleteConfigurationRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public subscribeConfiguration(options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse>;
    public subscribeConfiguration(metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientDuplexStream<runtime_pb.SubscribeConfigurationRequest, runtime_pb.SubscribeConfigurationResponse>;
    public tryLock(request: runtime_pb.TryLockRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.TryLockResponse) => void): grpc.ClientUnaryCall;
    public tryLock(request: runtime_pb.TryLockRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.TryLockResponse) => void): grpc.ClientUnaryCall;
    public tryLock(request: runtime_pb.TryLockRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.TryLockResponse) => void): grpc.ClientUnaryCall;
    public unlock(request: runtime_pb.UnlockRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.UnlockResponse) => void): grpc.ClientUnaryCall;
    public unlock(request: runtime_pb.UnlockRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.UnlockResponse) => void): grpc.ClientUnaryCall;
    public unlock(request: runtime_pb.UnlockRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.UnlockResponse) => void): grpc.ClientUnaryCall;
    public getNextId(request: runtime_pb.GetNextIdRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetNextIdResponse) => void): grpc.ClientUnaryCall;
    public getNextId(request: runtime_pb.GetNextIdRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetNextIdResponse) => void): grpc.ClientUnaryCall;
    public getNextId(request: runtime_pb.GetNextIdRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetNextIdResponse) => void): grpc.ClientUnaryCall;
    public getState(request: runtime_pb.GetStateRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetStateResponse) => void): grpc.ClientUnaryCall;
    public getState(request: runtime_pb.GetStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetStateResponse) => void): grpc.ClientUnaryCall;
    public getState(request: runtime_pb.GetStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetStateResponse) => void): grpc.ClientUnaryCall;
    public getBulkState(request: runtime_pb.GetBulkStateRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetBulkStateResponse) => void): grpc.ClientUnaryCall;
    public getBulkState(request: runtime_pb.GetBulkStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetBulkStateResponse) => void): grpc.ClientUnaryCall;
    public getBulkState(request: runtime_pb.GetBulkStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.GetBulkStateResponse) => void): grpc.ClientUnaryCall;
    public saveState(request: runtime_pb.SaveStateRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public saveState(request: runtime_pb.SaveStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public saveState(request: runtime_pb.SaveStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteState(request: runtime_pb.DeleteStateRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteState(request: runtime_pb.DeleteStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteState(request: runtime_pb.DeleteStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteBulkState(request: runtime_pb.DeleteBulkStateRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteBulkState(request: runtime_pb.DeleteBulkStateRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public deleteBulkState(request: runtime_pb.DeleteBulkStateRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public executeStateTransaction(request: runtime_pb.ExecuteStateTransactionRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public executeStateTransaction(request: runtime_pb.ExecuteStateTransactionRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public executeStateTransaction(request: runtime_pb.ExecuteStateTransactionRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public publishEvent(request: runtime_pb.PublishEventRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public publishEvent(request: runtime_pb.PublishEventRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public publishEvent(request: runtime_pb.PublishEventRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public getFile(request: runtime_pb.GetFileRequest, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<runtime_pb.GetFileResponse>;
    public getFile(request: runtime_pb.GetFileRequest, metadata?: grpc.Metadata, options?: Partial<grpc.CallOptions>): grpc.ClientReadableStream<runtime_pb.GetFileResponse>;
    public putFile(callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    public putFile(metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    public putFile(options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    public putFile(metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientWritableStream<runtime_pb.PutFileRequest>;
    public listFile(request: runtime_pb.ListFileRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.ListFileResp) => void): grpc.ClientUnaryCall;
    public listFile(request: runtime_pb.ListFileRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.ListFileResp) => void): grpc.ClientUnaryCall;
    public listFile(request: runtime_pb.ListFileRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.ListFileResp) => void): grpc.ClientUnaryCall;
    public delFile(request: runtime_pb.DelFileRequest, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public delFile(request: runtime_pb.DelFileRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public delFile(request: runtime_pb.DelFileRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: google_protobuf_empty_pb.Empty) => void): grpc.ClientUnaryCall;
    public invokeBinding(request: runtime_pb.InvokeBindingRequest, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeBindingResponse) => void): grpc.ClientUnaryCall;
    public invokeBinding(request: runtime_pb.InvokeBindingRequest, metadata: grpc.Metadata, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeBindingResponse) => void): grpc.ClientUnaryCall;
    public invokeBinding(request: runtime_pb.InvokeBindingRequest, metadata: grpc.Metadata, options: Partial<grpc.CallOptions>, callback: (error: grpc.ServiceError | null, response: runtime_pb.InvokeBindingResponse) => void): grpc.ClientUnaryCall;
}
