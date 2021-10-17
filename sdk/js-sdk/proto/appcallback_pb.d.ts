// package: spec.proto.runtime.v1
// file: appcallback.proto

/* tslint:disable */
/* eslint-disable */

import * as jspb from "google-protobuf";
import * as google_protobuf_empty_pb from "google-protobuf/google/protobuf/empty_pb";

export class TopicEventRequest extends jspb.Message { 
    getId(): string;
    setId(value: string): TopicEventRequest;
    getSource(): string;
    setSource(value: string): TopicEventRequest;
    getType(): string;
    setType(value: string): TopicEventRequest;
    getSpecVersion(): string;
    setSpecVersion(value: string): TopicEventRequest;
    getDataContentType(): string;
    setDataContentType(value: string): TopicEventRequest;
    getData(): Uint8Array | string;
    getData_asU8(): Uint8Array;
    getData_asB64(): string;
    setData(value: Uint8Array | string): TopicEventRequest;
    getTopic(): string;
    setTopic(value: string): TopicEventRequest;
    getPubsubName(): string;
    setPubsubName(value: string): TopicEventRequest;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TopicEventRequest.AsObject;
    static toObject(includeInstance: boolean, msg: TopicEventRequest): TopicEventRequest.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TopicEventRequest, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TopicEventRequest;
    static deserializeBinaryFromReader(message: TopicEventRequest, reader: jspb.BinaryReader): TopicEventRequest;
}

export namespace TopicEventRequest {
    export type AsObject = {
        id: string,
        source: string,
        type: string,
        specVersion: string,
        dataContentType: string,
        data: Uint8Array | string,
        topic: string,
        pubsubName: string,
    }
}

export class TopicEventResponse extends jspb.Message { 
    getStatus(): TopicEventResponse.TopicEventResponseStatus;
    setStatus(value: TopicEventResponse.TopicEventResponseStatus): TopicEventResponse;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TopicEventResponse.AsObject;
    static toObject(includeInstance: boolean, msg: TopicEventResponse): TopicEventResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TopicEventResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TopicEventResponse;
    static deserializeBinaryFromReader(message: TopicEventResponse, reader: jspb.BinaryReader): TopicEventResponse;
}

export namespace TopicEventResponse {
    export type AsObject = {
        status: TopicEventResponse.TopicEventResponseStatus,
    }

    export enum TopicEventResponseStatus {
    SUCCESS = 0,
    RETRY = 1,
    DROP = 2,
    }

}

export class ListTopicSubscriptionsResponse extends jspb.Message { 
    clearSubscriptionsList(): void;
    getSubscriptionsList(): Array<TopicSubscription>;
    setSubscriptionsList(value: Array<TopicSubscription>): ListTopicSubscriptionsResponse;
    addSubscriptions(value?: TopicSubscription, index?: number): TopicSubscription;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): ListTopicSubscriptionsResponse.AsObject;
    static toObject(includeInstance: boolean, msg: ListTopicSubscriptionsResponse): ListTopicSubscriptionsResponse.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: ListTopicSubscriptionsResponse, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): ListTopicSubscriptionsResponse;
    static deserializeBinaryFromReader(message: ListTopicSubscriptionsResponse, reader: jspb.BinaryReader): ListTopicSubscriptionsResponse;
}

export namespace ListTopicSubscriptionsResponse {
    export type AsObject = {
        subscriptionsList: Array<TopicSubscription.AsObject>,
    }
}

export class TopicSubscription extends jspb.Message { 
    getPubsubName(): string;
    setPubsubName(value: string): TopicSubscription;
    getTopic(): string;
    setTopic(value: string): TopicSubscription;

    getMetadataMap(): jspb.Map<string, string>;
    clearMetadataMap(): void;

    serializeBinary(): Uint8Array;
    toObject(includeInstance?: boolean): TopicSubscription.AsObject;
    static toObject(includeInstance: boolean, msg: TopicSubscription): TopicSubscription.AsObject;
    static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
    static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
    static serializeBinaryToWriter(message: TopicSubscription, writer: jspb.BinaryWriter): void;
    static deserializeBinary(bytes: Uint8Array): TopicSubscription;
    static deserializeBinaryFromReader(message: TopicSubscription, reader: jspb.BinaryReader): TopicSubscription;
}

export namespace TopicSubscription {
    export type AsObject = {
        pubsubName: string,
        topic: string,

        metadataMap: Array<[string, string]>,
    }
}
