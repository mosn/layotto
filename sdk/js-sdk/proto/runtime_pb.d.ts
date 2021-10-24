// package: spec.proto.runtime.v1
// file: runtime.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";
import * as google_protobuf_any_pb from "google-protobuf/google/protobuf/any_pb";

export class GetFileRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): GetFileRequest;
    getName(): string;
    setName(value: string): GetFileRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetFileRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetFileRequest): GetFileRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetFileRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetFileRequest;
    static deserializeBinaryFromReader(message: GetFileRequest, reader: jspb.BinaryReader): GetFileRequest;
}

export namespace GetFileRequest {
    export type AsObject = {
        storeName: string,
        name: string,

        metadataMap: Array<[string, string]>,
    }
}

export class GetFileResponse extends jspb.Message { 
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): GetFileResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetFileResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetFileResponse): GetFileResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetFileResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetFileResponse;
    static deserializeBinaryFromReader(message: GetFileResponse, reader: jspb.BinaryReader): GetFileResponse;
}

export namespace GetFileResponse {
    export type AsObject = {
        data: Uint8Array | string,
    }
}

export class PutFileRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): PutFileRequest;
    getName(): string;
    setName(value: string): PutFileRequest;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): PutFileRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PutFileRequest.AsObject;
    static toObject(includeInstance: boolean, msg: PutFileRequest): PutFileRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PutFileRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PutFileRequest;
    static deserializeBinaryFromReader(message: PutFileRequest, reader: jspb.BinaryReader): PutFileRequest;
}

export namespace PutFileRequest {
    export type AsObject = {
        storeName: string,
        name: string,
        data: Uint8Array | string,

        metadataMap: Array<[string, string]>,
    }
}

export class FileRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): FileRequest;
    getName(): string;
    setName(value: string): FileRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): FileRequest.AsObject;
    static toObject(includeInstance: boolean, msg: FileRequest): FileRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: FileRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): FileRequest;
    static deserializeBinaryFromReader(message: FileRequest, reader: jspb.BinaryReader): FileRequest;
}

export namespace FileRequest {
    export type AsObject = {
        storeName: string,
        name: string,

        metadataMap: Array<[string, string]>,
    }
}

export class ListFileRequest extends jspb.Message { 

    hasRequest(): boolean;
    clearRequest(): void;
    getRequest(): FileRequest | undefined;
    setRequest(value?: FileRequest): ListFileRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListFileRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ListFileRequest): ListFileRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ListFileRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListFileRequest;
    static deserializeBinaryFromReader(message: ListFileRequest, reader: jspb.BinaryReader): ListFileRequest;
}

export namespace ListFileRequest {
    export type AsObject = {
        request?: FileRequest.AsObject,
    }
}

export class ListFileResp extends jspb.Message { 
    clearFileNameList(): void;
    getFileNameList(): Array<string>;
    setFileNameList(value: Array<string>): ListFileResp;
    addFileName(value: string, index?: number): string;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListFileResp.AsObject;
    static toObject(includeInstance: boolean, msg: ListFileResp): ListFileResp.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ListFileResp, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListFileResp;
    static deserializeBinaryFromReader(message: ListFileResp, reader: jspb.BinaryReader): ListFileResp;
}

export namespace ListFileResp {
    export type AsObject = {
        fileNameList: Array<string>,
    }
}

export class DelFileRequest extends jspb.Message { 

    hasRequest(): boolean;
    clearRequest(): void;
    getRequest(): FileRequest | undefined;
    setRequest(value?: FileRequest): DelFileRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DelFileRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DelFileRequest): DelFileRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DelFileRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DelFileRequest;
    static deserializeBinaryFromReader(message: DelFileRequest, reader: jspb.BinaryReader): DelFileRequest;
}

export namespace DelFileRequest {
    export type AsObject = {
        request?: FileRequest.AsObject,
    }
}

export class GetNextIdRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): GetNextIdRequest;
    getKey(): string;
    setKey(value: string): GetNextIdRequest;

    hasOptions(): boolean;
    clearOptions(): void;
    getOptions(): SequencerOptions | undefined;
    setOptions(value?: SequencerOptions): GetNextIdRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetNextIdRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetNextIdRequest): GetNextIdRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetNextIdRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetNextIdRequest;
    static deserializeBinaryFromReader(message: GetNextIdRequest, reader: jspb.BinaryReader): GetNextIdRequest;
}

export namespace GetNextIdRequest {
    export type AsObject = {
        storeName: string,
        key: string,
        options?: SequencerOptions.AsObject,

        metadataMap: Array<[string, string]>,
    }
}

export class SequencerOptions extends jspb.Message { 
    getIncrement(): SequencerOptions.AutoIncrement;
    setIncrement(value: SequencerOptions.AutoIncrement): SequencerOptions;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SequencerOptions.AsObject;
    static toObject(includeInstance: boolean, msg: SequencerOptions): SequencerOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SequencerOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SequencerOptions;
    static deserializeBinaryFromReader(message: SequencerOptions, reader: jspb.BinaryReader): SequencerOptions;
}

export namespace SequencerOptions {
    export type AsObject = {
        increment: SequencerOptions.AutoIncrement,
    }

    export enum AutoIncrement {
    WEAK = 0,
    STRONG = 1,
    }

}

export class GetNextIdResponse extends jspb.Message { 
    getNextId(): string;
    setNextId(value: string): GetNextIdResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetNextIdResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetNextIdResponse): GetNextIdResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetNextIdResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetNextIdResponse;
    static deserializeBinaryFromReader(message: GetNextIdResponse, reader: jspb.BinaryReader): GetNextIdResponse;
}

export namespace GetNextIdResponse {
    export type AsObject = {
        nextId: string,
    }
}

export class TryLockRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): TryLockRequest;
    getResourceId(): string;
    setResourceId(value: string): TryLockRequest;
    getLockOwner(): string;
    setLockOwner(value: string): TryLockRequest;
    getExpire(): number;
    setExpire(value: number): TryLockRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TryLockRequest.AsObject;
    static toObject(includeInstance: boolean, msg: TryLockRequest): TryLockRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TryLockRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TryLockRequest;
    static deserializeBinaryFromReader(message: TryLockRequest, reader: jspb.BinaryReader): TryLockRequest;
}

export namespace TryLockRequest {
    export type AsObject = {
        storeName: string,
        resourceId: string,
        lockOwner: string,
        expire: number,
    }
}

export class TryLockResponse extends jspb.Message { 
    getSuccess(): boolean;
    setSuccess(value: boolean): TryLockResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TryLockResponse.AsObject;
    static toObject(includeInstance: boolean, msg: TryLockResponse): TryLockResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TryLockResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TryLockResponse;
    static deserializeBinaryFromReader(message: TryLockResponse, reader: jspb.BinaryReader): TryLockResponse;
}

export namespace TryLockResponse {
    export type AsObject = {
        success: boolean,
    }
}

export class UnlockRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): UnlockRequest;
    getResourceId(): string;
    setResourceId(value: string): UnlockRequest;
    getLockOwner(): string;
    setLockOwner(value: string): UnlockRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UnlockRequest.AsObject;
    static toObject(includeInstance: boolean, msg: UnlockRequest): UnlockRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UnlockRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UnlockRequest;
    static deserializeBinaryFromReader(message: UnlockRequest, reader: jspb.BinaryReader): UnlockRequest;
}

export namespace UnlockRequest {
    export type AsObject = {
        storeName: string,
        resourceId: string,
        lockOwner: string,
    }
}

export class UnlockResponse extends jspb.Message { 
    getStatus(): UnlockResponse.Status;
    setStatus(value: UnlockResponse.Status): UnlockResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): UnlockResponse.AsObject;
    static toObject(includeInstance: boolean, msg: UnlockResponse): UnlockResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: UnlockResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): UnlockResponse;
    static deserializeBinaryFromReader(message: UnlockResponse, reader: jspb.BinaryReader): UnlockResponse;
}

export namespace UnlockResponse {
    export type AsObject = {
        status: UnlockResponse.Status,
    }

    export enum Status {
    SUCCESS = 0,
    LOCK_UNEXIST = 1,
    LOCK_BELONG_TO_OTHERS = 2,
    INTERNAL_ERROR = 3,
    }

}

export class SayHelloRequest extends jspb.Message { 
    getServiceName(): string;
    setServiceName(value: string): SayHelloRequest;
    getName(): string;
    setName(value: string): SayHelloRequest;

    hasData(): boolean;
    clearData(): void;
    getData(): google_protobuf_any_pb.Any | undefined;
    setData(value?: google_protobuf_any_pb.Any): SayHelloRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SayHelloRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SayHelloRequest): SayHelloRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SayHelloRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SayHelloRequest;
    static deserializeBinaryFromReader(message: SayHelloRequest, reader: jspb.BinaryReader): SayHelloRequest;
}

export namespace SayHelloRequest {
    export type AsObject = {
        serviceName: string,
        name: string,
        data?: google_protobuf_any_pb.Any.AsObject,
    }
}

export class SayHelloResponse extends jspb.Message { 
    getHello(): string;
    setHello(value: string): SayHelloResponse;

    hasData(): boolean;
    clearData(): void;
    getData(): google_protobuf_any_pb.Any | undefined;
    setData(value?: google_protobuf_any_pb.Any): SayHelloResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SayHelloResponse.AsObject;
    static toObject(includeInstance: boolean, msg: SayHelloResponse): SayHelloResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SayHelloResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SayHelloResponse;
    static deserializeBinaryFromReader(message: SayHelloResponse, reader: jspb.BinaryReader): SayHelloResponse;
}

export namespace SayHelloResponse {
    export type AsObject = {
        hello: string,
        data?: google_protobuf_any_pb.Any.AsObject,
    }
}

export class InvokeServiceRequest extends jspb.Message { 
    getId(): string;
    setId(value: string): InvokeServiceRequest;

    hasMessage(): boolean;
    clearMessage(): void;
    getMessage(): CommonInvokeRequest | undefined;
    setMessage(value?: CommonInvokeRequest): InvokeServiceRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvokeServiceRequest.AsObject;
    static toObject(includeInstance: boolean, msg: InvokeServiceRequest): InvokeServiceRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvokeServiceRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvokeServiceRequest;
    static deserializeBinaryFromReader(message: InvokeServiceRequest, reader: jspb.BinaryReader): InvokeServiceRequest;
}

export namespace InvokeServiceRequest {
    export type AsObject = {
        id: string,
        message?: CommonInvokeRequest.AsObject,
    }
}

export class CommonInvokeRequest extends jspb.Message { 
    getMethod(): string;
    setMethod(value: string): CommonInvokeRequest;

    hasData(): boolean;
    clearData(): void;
    getData(): google_protobuf_any_pb.Any | undefined;
    setData(value?: google_protobuf_any_pb.Any): CommonInvokeRequest;
    getContentType(): string;
    setContentType(value: string): CommonInvokeRequest;

    hasHttpExtension(): boolean;
    clearHttpExtension(): void;
    getHttpExtension(): HTTPExtension | undefined;
    setHttpExtension(value?: HTTPExtension): CommonInvokeRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): CommonInvokeRequest.AsObject;
    static toObject(includeInstance: boolean, msg: CommonInvokeRequest): CommonInvokeRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: CommonInvokeRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): CommonInvokeRequest;
    static deserializeBinaryFromReader(message: CommonInvokeRequest, reader: jspb.BinaryReader): CommonInvokeRequest;
}

export namespace CommonInvokeRequest {
    export type AsObject = {
        method: string,
        data?: google_protobuf_any_pb.Any.AsObject,
        contentType: string,
        httpExtension?: HTTPExtension.AsObject,
    }
}

export class HTTPExtension extends jspb.Message { 
    getVerb(): HTTPExtension.Verb;
    setVerb(value: HTTPExtension.Verb): HTTPExtension;
    getQuerystring(): string;
    setQuerystring(value: string): HTTPExtension;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): HTTPExtension.AsObject;
    static toObject(includeInstance: boolean, msg: HTTPExtension): HTTPExtension.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: HTTPExtension, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): HTTPExtension;
    static deserializeBinaryFromReader(message: HTTPExtension, reader: jspb.BinaryReader): HTTPExtension;
}

export namespace HTTPExtension {
    export type AsObject = {
        verb: HTTPExtension.Verb,
        querystring: string,
    }

    export enum Verb {
    NONE = 0,
    GET = 1,
    HEAD = 2,
    POST = 3,
    PUT = 4,
    DELETE = 5,
    CONNECT = 6,
    OPTIONS = 7,
    TRACE = 8,
    }

}

export class InvokeResponse extends jspb.Message { 

    hasData(): boolean;
    clearData(): void;
    getData(): google_protobuf_any_pb.Any | undefined;
    setData(value?: google_protobuf_any_pb.Any): InvokeResponse;
    getContentType(): string;
    setContentType(value: string): InvokeResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvokeResponse.AsObject;
    static toObject(includeInstance: boolean, msg: InvokeResponse): InvokeResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvokeResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvokeResponse;
    static deserializeBinaryFromReader(message: InvokeResponse, reader: jspb.BinaryReader): InvokeResponse;
}

export namespace InvokeResponse {
    export type AsObject = {
        data?: google_protobuf_any_pb.Any.AsObject,
        contentType: string,
    }
}

export class ConfigurationItem extends jspb.Message { 
    getKey(): string;
    setKey(value: string): ConfigurationItem;
    getContent(): string;
    setContent(value: string): ConfigurationItem;
    getGroup(): string;
    setGroup(value: string): ConfigurationItem;
    getLabel(): string;
    setLabel(value: string): ConfigurationItem;

    getTagsMap(): jspb.Map<string, string>;
    clearTagsMap(): void;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ConfigurationItem.AsObject;
    static toObject(includeInstance: boolean, msg: ConfigurationItem): ConfigurationItem.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ConfigurationItem, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ConfigurationItem;
    static deserializeBinaryFromReader(message: ConfigurationItem, reader: jspb.BinaryReader): ConfigurationItem;
}

export namespace ConfigurationItem {
    export type AsObject = {
        key: string,
        content: string,
        group: string,
        label: string,

        tagsMap: Array<[string, string]>,

        metadataMap: Array<[string, string]>,
    }
}

export class GetConfigurationRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): GetConfigurationRequest;
    getAppId(): string;
    setAppId(value: string): GetConfigurationRequest;
    getGroup(): string;
    setGroup(value: string): GetConfigurationRequest;
    getLabel(): string;
    setLabel(value: string): GetConfigurationRequest;
    clearKeysList(): void;
    getKeysList(): Array<string>;
    setKeysList(value: Array<string>): GetConfigurationRequest;
    addKeys(value: string, index?: number): string;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;
    getSubscribeUpdate(): boolean;
    setSubscribeUpdate(value: boolean): GetConfigurationRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetConfigurationRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetConfigurationRequest): GetConfigurationRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetConfigurationRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetConfigurationRequest;
    static deserializeBinaryFromReader(message: GetConfigurationRequest, reader: jspb.BinaryReader): GetConfigurationRequest;
}

export namespace GetConfigurationRequest {
    export type AsObject = {
        storeName: string,
        appId: string,
        group: string,
        label: string,
        keysList: Array<string>,

        metadataMap: Array<[string, string]>,
        subscribeUpdate: boolean,
    }
}

export class GetConfigurationResponse extends jspb.Message { 
    clearItemsList(): void;
    getItemsList(): Array<ConfigurationItem>;
    setItemsList(value: Array<ConfigurationItem>): GetConfigurationResponse;
    addItems(value?: ConfigurationItem, index?: number): ConfigurationItem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetConfigurationResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetConfigurationResponse): GetConfigurationResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetConfigurationResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetConfigurationResponse;
    static deserializeBinaryFromReader(message: GetConfigurationResponse, reader: jspb.BinaryReader): GetConfigurationResponse;
}

export namespace GetConfigurationResponse {
    export type AsObject = {
        itemsList: Array<ConfigurationItem.AsObject>,
    }
}

export class SubscribeConfigurationRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): SubscribeConfigurationRequest;
    getAppId(): string;
    setAppId(value: string): SubscribeConfigurationRequest;
    getGroup(): string;
    setGroup(value: string): SubscribeConfigurationRequest;
    getLabel(): string;
    setLabel(value: string): SubscribeConfigurationRequest;
    clearKeysList(): void;
    getKeysList(): Array<string>;
    setKeysList(value: Array<string>): SubscribeConfigurationRequest;
    addKeys(value: string, index?: number): string;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubscribeConfigurationRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SubscribeConfigurationRequest): SubscribeConfigurationRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubscribeConfigurationRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubscribeConfigurationRequest;
    static deserializeBinaryFromReader(message: SubscribeConfigurationRequest, reader: jspb.BinaryReader): SubscribeConfigurationRequest;
}

export namespace SubscribeConfigurationRequest {
    export type AsObject = {
        storeName: string,
        appId: string,
        group: string,
        label: string,
        keysList: Array<string>,

        metadataMap: Array<[string, string]>,
    }
}

export class SubscribeConfigurationResponse extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): SubscribeConfigurationResponse;
    getAppId(): string;
    setAppId(value: string): SubscribeConfigurationResponse;
    clearItemsList(): void;
    getItemsList(): Array<ConfigurationItem>;
    setItemsList(value: Array<ConfigurationItem>): SubscribeConfigurationResponse;
    addItems(value?: ConfigurationItem, index?: number): ConfigurationItem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SubscribeConfigurationResponse.AsObject;
    static toObject(includeInstance: boolean, msg: SubscribeConfigurationResponse): SubscribeConfigurationResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SubscribeConfigurationResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SubscribeConfigurationResponse;
    static deserializeBinaryFromReader(message: SubscribeConfigurationResponse, reader: jspb.BinaryReader): SubscribeConfigurationResponse;
}

export namespace SubscribeConfigurationResponse {
    export type AsObject = {
        storeName: string,
        appId: string,
        itemsList: Array<ConfigurationItem.AsObject>,
    }
}

export class SaveConfigurationRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): SaveConfigurationRequest;
    getAppId(): string;
    setAppId(value: string): SaveConfigurationRequest;
    clearItemsList(): void;
    getItemsList(): Array<ConfigurationItem>;
    setItemsList(value: Array<ConfigurationItem>): SaveConfigurationRequest;
    addItems(value?: ConfigurationItem, index?: number): ConfigurationItem;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SaveConfigurationRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SaveConfigurationRequest): SaveConfigurationRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SaveConfigurationRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SaveConfigurationRequest;
    static deserializeBinaryFromReader(message: SaveConfigurationRequest, reader: jspb.BinaryReader): SaveConfigurationRequest;
}

export namespace SaveConfigurationRequest {
    export type AsObject = {
        storeName: string,
        appId: string,
        itemsList: Array<ConfigurationItem.AsObject>,

        metadataMap: Array<[string, string]>,
    }
}

export class DeleteConfigurationRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): DeleteConfigurationRequest;
    getAppId(): string;
    setAppId(value: string): DeleteConfigurationRequest;
    getGroup(): string;
    setGroup(value: string): DeleteConfigurationRequest;
    getLabel(): string;
    setLabel(value: string): DeleteConfigurationRequest;
    clearKeysList(): void;
    getKeysList(): Array<string>;
    setKeysList(value: Array<string>): DeleteConfigurationRequest;
    addKeys(value: string, index?: number): string;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteConfigurationRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteConfigurationRequest): DeleteConfigurationRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteConfigurationRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteConfigurationRequest;
    static deserializeBinaryFromReader(message: DeleteConfigurationRequest, reader: jspb.BinaryReader): DeleteConfigurationRequest;
}

export namespace DeleteConfigurationRequest {
    export type AsObject = {
        storeName: string,
        appId: string,
        group: string,
        label: string,
        keysList: Array<string>,

        metadataMap: Array<[string, string]>,
    }
}

export class GetStateRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): GetStateRequest;
    getKey(): string;
    setKey(value: string): GetStateRequest;
    getConsistency(): StateOptions.StateConsistency;
    setConsistency(value: StateOptions.StateConsistency): GetStateRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetStateRequest): GetStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetStateRequest;
    static deserializeBinaryFromReader(message: GetStateRequest, reader: jspb.BinaryReader): GetStateRequest;
}

export namespace GetStateRequest {
    export type AsObject = {
        storeName: string,
        key: string,
        consistency: StateOptions.StateConsistency,

        metadataMap: Array<[string, string]>,
    }
}

export class GetBulkStateRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): GetBulkStateRequest;
    clearKeysList(): void;
    getKeysList(): Array<string>;
    setKeysList(value: Array<string>): GetBulkStateRequest;
    addKeys(value: string, index?: number): string;
    getParallelism(): number;
    setParallelism(value: number): GetBulkStateRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetBulkStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: GetBulkStateRequest): GetBulkStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetBulkStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetBulkStateRequest;
    static deserializeBinaryFromReader(message: GetBulkStateRequest, reader: jspb.BinaryReader): GetBulkStateRequest;
}

export namespace GetBulkStateRequest {
    export type AsObject = {
        storeName: string,
        keysList: Array<string>,
        parallelism: number,

        metadataMap: Array<[string, string]>,
    }
}

export class GetBulkStateResponse extends jspb.Message { 
    clearItemsList(): void;
    getItemsList(): Array<BulkStateItem>;
    setItemsList(value: Array<BulkStateItem>): GetBulkStateResponse;
    addItems(value?: BulkStateItem, index?: number): BulkStateItem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetBulkStateResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetBulkStateResponse): GetBulkStateResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetBulkStateResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetBulkStateResponse;
    static deserializeBinaryFromReader(message: GetBulkStateResponse, reader: jspb.BinaryReader): GetBulkStateResponse;
}

export namespace GetBulkStateResponse {
    export type AsObject = {
        itemsList: Array<BulkStateItem.AsObject>,
    }
}

export class BulkStateItem extends jspb.Message { 
    getKey(): string;
    setKey(value: string): BulkStateItem;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): BulkStateItem;
    getEtag(): string;
    setEtag(value: string): BulkStateItem;
    getError(): string;
    setError(value: string): BulkStateItem;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): BulkStateItem.AsObject;
    static toObject(includeInstance: boolean, msg: BulkStateItem): BulkStateItem.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: BulkStateItem, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): BulkStateItem;
    static deserializeBinaryFromReader(message: BulkStateItem, reader: jspb.BinaryReader): BulkStateItem;
}

export namespace BulkStateItem {
    export type AsObject = {
        key: string,
        data: Uint8Array | string,
        etag: string,
        error: string,

        metadataMap: Array<[string, string]>,
    }
}

export class GetStateResponse extends jspb.Message { 
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): GetStateResponse;
    getEtag(): string;
    setEtag(value: string): GetStateResponse;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): GetStateResponse.AsObject;
    static toObject(includeInstance: boolean, msg: GetStateResponse): GetStateResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: GetStateResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): GetStateResponse;
    static deserializeBinaryFromReader(message: GetStateResponse, reader: jspb.BinaryReader): GetStateResponse;
}

export namespace GetStateResponse {
    export type AsObject = {
        data: Uint8Array | string,
        etag: string,

        metadataMap: Array<[string, string]>,
    }
}

export class DeleteStateRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): DeleteStateRequest;
    getKey(): string;
    setKey(value: string): DeleteStateRequest;

    hasEtag(): boolean;
    clearEtag(): void;
    getEtag(): Etag | undefined;
    setEtag(value?: Etag): DeleteStateRequest;

    hasOptions(): boolean;
    clearOptions(): void;
    getOptions(): StateOptions | undefined;
    setOptions(value?: StateOptions): DeleteStateRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteStateRequest): DeleteStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteStateRequest;
    static deserializeBinaryFromReader(message: DeleteStateRequest, reader: jspb.BinaryReader): DeleteStateRequest;
}

export namespace DeleteStateRequest {
    export type AsObject = {
        storeName: string,
        key: string,
        etag?: Etag.AsObject,
        options?: StateOptions.AsObject,

        metadataMap: Array<[string, string]>,
    }
}

export class DeleteBulkStateRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): DeleteBulkStateRequest;
    clearStatesList(): void;
    getStatesList(): Array<StateItem>;
    setStatesList(value: Array<StateItem>): DeleteBulkStateRequest;
    addStates(value?: StateItem, index?: number): StateItem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): DeleteBulkStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: DeleteBulkStateRequest): DeleteBulkStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: DeleteBulkStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): DeleteBulkStateRequest;
    static deserializeBinaryFromReader(message: DeleteBulkStateRequest, reader: jspb.BinaryReader): DeleteBulkStateRequest;
}

export namespace DeleteBulkStateRequest {
    export type AsObject = {
        storeName: string,
        statesList: Array<StateItem.AsObject>,
    }
}

export class SaveStateRequest extends jspb.Message { 
    getStoreName(): string;
    setStoreName(value: string): SaveStateRequest;
    clearStatesList(): void;
    getStatesList(): Array<StateItem>;
    setStatesList(value: Array<StateItem>): SaveStateRequest;
    addStates(value?: StateItem, index?: number): StateItem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): SaveStateRequest.AsObject;
    static toObject(includeInstance: boolean, msg: SaveStateRequest): SaveStateRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: SaveStateRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): SaveStateRequest;
    static deserializeBinaryFromReader(message: SaveStateRequest, reader: jspb.BinaryReader): SaveStateRequest;
}

export namespace SaveStateRequest {
    export type AsObject = {
        storeName: string,
        statesList: Array<StateItem.AsObject>,
    }
}

export class StateItem extends jspb.Message { 
    getKey(): string;
    setKey(value: string): StateItem;
    getValue(): Uint8Array | string;
    getValue_asU8(): Uint8Array;
    getValue_asB64(): string;
    setValue(value: Uint8Array | string): StateItem;

    hasEtag(): boolean;
    clearEtag(): void;
    getEtag(): Etag | undefined;
    setEtag(value?: Etag): StateItem;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    hasOptions(): boolean;
    clearOptions(): void;
    getOptions(): StateOptions | undefined;
    setOptions(value?: StateOptions): StateItem;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StateItem.AsObject;
    static toObject(includeInstance: boolean, msg: StateItem): StateItem.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StateItem, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StateItem;
    static deserializeBinaryFromReader(message: StateItem, reader: jspb.BinaryReader): StateItem;
}

export namespace StateItem {
    export type AsObject = {
        key: string,
        value: Uint8Array | string,
        etag?: Etag.AsObject,

        metadataMap: Array<[string, string]>,
        options?: StateOptions.AsObject,
    }
}

export class Etag extends jspb.Message { 
    getValue(): string;
    setValue(value: string): Etag;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): Etag.AsObject;
    static toObject(includeInstance: boolean, msg: Etag): Etag.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: Etag, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): Etag;
    static deserializeBinaryFromReader(message: Etag, reader: jspb.BinaryReader): Etag;
}

export namespace Etag {
    export type AsObject = {
        value: string,
    }
}

export class StateOptions extends jspb.Message { 
    getConcurrency(): StateOptions.StateConcurrency;
    setConcurrency(value: StateOptions.StateConcurrency): StateOptions;
    getConsistency(): StateOptions.StateConsistency;
    setConsistency(value: StateOptions.StateConsistency): StateOptions;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): StateOptions.AsObject;
    static toObject(includeInstance: boolean, msg: StateOptions): StateOptions.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: StateOptions, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): StateOptions;
    static deserializeBinaryFromReader(message: StateOptions, reader: jspb.BinaryReader): StateOptions;
}

export namespace StateOptions {
    export type AsObject = {
        concurrency: StateOptions.StateConcurrency,
        consistency: StateOptions.StateConsistency,
    }

    export enum StateConcurrency {
    CONCURRENCY_UNSPECIFIED = 0,
    CONCURRENCY_FIRST_WRITE = 1,
    CONCURRENCY_LAST_WRITE = 2,
    }

    export enum StateConsistency {
    CONSISTENCY_UNSPECIFIED = 0,
    CONSISTENCY_EVENTUAL = 1,
    CONSISTENCY_STRONG = 2,
    }

}

export class TransactionalStateOperation extends jspb.Message { 
    getOperationtype(): string;
    setOperationtype(value: string): TransactionalStateOperation;

    hasRequest(): boolean;
    clearRequest(): void;
    getRequest(): StateItem | undefined;
    setRequest(value?: StateItem): TransactionalStateOperation;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TransactionalStateOperation.AsObject;
    static toObject(includeInstance: boolean, msg: TransactionalStateOperation): TransactionalStateOperation.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TransactionalStateOperation, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TransactionalStateOperation;
    static deserializeBinaryFromReader(message: TransactionalStateOperation, reader: jspb.BinaryReader): TransactionalStateOperation;
}

export namespace TransactionalStateOperation {
    export type AsObject = {
        operationtype: string,
        request?: StateItem.AsObject,
    }
}

export class ExecuteStateTransactionRequest extends jspb.Message { 
    getStorename(): string;
    setStorename(value: string): ExecuteStateTransactionRequest;
    clearOperationsList(): void;
    getOperationsList(): Array<TransactionalStateOperation>;
    setOperationsList(value: Array<TransactionalStateOperation>): ExecuteStateTransactionRequest;
    addOperations(value?: TransactionalStateOperation, index?: number): TransactionalStateOperation;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ExecuteStateTransactionRequest.AsObject;
    static toObject(includeInstance: boolean, msg: ExecuteStateTransactionRequest): ExecuteStateTransactionRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ExecuteStateTransactionRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ExecuteStateTransactionRequest;
    static deserializeBinaryFromReader(message: ExecuteStateTransactionRequest, reader: jspb.BinaryReader): ExecuteStateTransactionRequest;
}

export namespace ExecuteStateTransactionRequest {
    export type AsObject = {
        storename: string,
        operationsList: Array<TransactionalStateOperation.AsObject>,

        metadataMap: Array<[string, string]>,
    }
}

export class PublishEventRequest extends jspb.Message { 
    getPubsubName(): string;
    setPubsubName(value: string): PublishEventRequest;
    getTopic(): string;
    setTopic(value: string): PublishEventRequest;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): PublishEventRequest;
    getDataContentType(): string;
    setDataContentType(value: string): PublishEventRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): PublishEventRequest.AsObject;
    static toObject(includeInstance: boolean, msg: PublishEventRequest): PublishEventRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: PublishEventRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): PublishEventRequest;
    static deserializeBinaryFromReader(message: PublishEventRequest, reader: jspb.BinaryReader): PublishEventRequest;
}

export namespace PublishEventRequest {
    export type AsObject = {
        pubsubName: string,
        topic: string,
        data: Uint8Array | string,
        dataContentType: string,

        metadataMap: Array<[string, string]>,
    }
}

export class InvokeBindingRequest extends jspb.Message { 
    getName(): string;
    setName(value: string): InvokeBindingRequest;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): InvokeBindingRequest;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;
    getOperation(): string;
    setOperation(value: string): InvokeBindingRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvokeBindingRequest.AsObject;
    static toObject(includeInstance: boolean, msg: InvokeBindingRequest): InvokeBindingRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvokeBindingRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvokeBindingRequest;
    static deserializeBinaryFromReader(message: InvokeBindingRequest, reader: jspb.BinaryReader): InvokeBindingRequest;
}

export namespace InvokeBindingRequest {
    export type AsObject = {
        name: string,
        data: Uint8Array | string,

        metadataMap: Array<[string, string]>,
        operation: string,
    }
}

export class InvokeBindingResponse extends jspb.Message { 
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): InvokeBindingResponse;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): InvokeBindingResponse.AsObject;
    static toObject(includeInstance: boolean, msg: InvokeBindingResponse): InvokeBindingResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: InvokeBindingResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): InvokeBindingResponse;
    static deserializeBinaryFromReader(message: InvokeBindingResponse, reader: jspb.BinaryReader): InvokeBindingResponse;
}

export namespace InvokeBindingResponse {
    export type AsObject = {
        data: Uint8Array | string,

        metadataMap: Array<[string, string]>,
    }
}
