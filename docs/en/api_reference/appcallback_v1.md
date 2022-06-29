

<a name="appcallback.proto"></a>

# appcallback.proto
<a name="top"></a>

This document is automaticallly generated from the [`.proto`](https://github.com/mosn/layotto/tree/main/spec/proto/runtime/v1) files.




<a name="spec.proto.runtime.v1.AppCallback"></a>

## [gRPC Service] AppCallback
AppCallback V1 allows user application to interact with runtime.
User application needs to implement AppCallback service if it needs to
receive message from runtime.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListTopicSubscriptions | [.google.protobuf.Empty](#google.protobuf.Empty) | [ListTopicSubscriptionsResponse](#spec.proto.runtime.v1.ListTopicSubscriptionsResponse) | Lists all topics subscribed by this app. |
| OnTopicEvent | [TopicEventRequest](#spec.proto.runtime.v1.TopicEventRequest) | [TopicEventResponse](#spec.proto.runtime.v1.TopicEventResponse) | Subscribes events from Pubsub |

 <!-- end services -->


<a name="spec.proto.runtime.v1.ListTopicSubscriptionsResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## ListTopicSubscriptionsResponse
ListTopicSubscriptionsResponse is the message including the list of the subscribing topics.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriptions | [TopicSubscription](#spec.proto.runtime.v1.TopicSubscription) | repeated | The list of topics. |






<a name="spec.proto.runtime.v1.TopicEventRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## TopicEventRequest
TopicEventRequest message is compatible with CloudEvent spec v1.0
https://github.com/cloudevents/spec/blob/v1.0/spec.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id identifies the event. Producers MUST ensure that source + id is unique for each distinct event. If a duplicate event is re-sent (e.g. due to a network error) it MAY have the same id. |
| source | [string](#string) |  | source identifies the context in which an event happened. Often this will include information such as the type of the event source, the organization publishing the event or the process that produced the event. The exact syntax and semantics behind the data encoded in the URI is defined by the event producer. |
| type | [string](#string) |  | The type of event related to the originating occurrence. |
| spec_version | [string](#string) |  | The version of the CloudEvents specification. |
| data_content_type | [string](#string) |  | The content type of data value. |
| data | [bytes](#bytes) |  | The content of the event. |
| topic | [string](#string) |  | The pubsub topic which publisher sent to. |
| pubsub_name | [string](#string) |  | The name of the pubsub the publisher sent to. |
| metadata | [TopicEventRequest.MetadataEntry](#spec.proto.runtime.v1.TopicEventRequest.MetadataEntry) | repeated | add a map to pass some extra properties. |






<a name="spec.proto.runtime.v1.TopicEventRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## TopicEventRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.TopicEventResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## TopicEventResponse
TopicEventResponse is response from app on published message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [TopicEventResponse.TopicEventResponseStatus](#spec.proto.runtime.v1.TopicEventResponse.TopicEventResponseStatus) |  | The list of output bindings. |






<a name="spec.proto.runtime.v1.TopicSubscription"></a>
<p align="right"><a href="#top">Top</a></p>

## TopicSubscription
TopicSubscription represents topic and metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pubsub_name | [string](#string) |  | Required. The name of the pubsub containing the topic below to subscribe to. |
| topic | [string](#string) |  | Required. The name of topic which will be subscribed |
| metadata | [TopicSubscription.MetadataEntry](#spec.proto.runtime.v1.TopicSubscription.MetadataEntry) | repeated | The optional properties used for this topic's subscription e.g. session id |






<a name="spec.proto.runtime.v1.TopicSubscription.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## TopicSubscription.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 <!-- end messages -->


<a name="spec.proto.runtime.v1.TopicEventResponse.TopicEventResponseStatus"></a>

## TopicEventResponse.TopicEventResponseStatus
TopicEventResponseStatus allows apps to have finer control over handling of the message.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SUCCESS | 0 | SUCCESS is the default behavior: message is acknowledged and not retried or logged. |
| RETRY | 1 | RETRY status signals runtime to retry the message as part of an expected scenario (no warning is logged). |
| DROP | 2 | DROP status signals runtime to drop the message as part of an unexpected scenario (warning is logged). |


 <!-- end enums -->

 <!-- end HasExtensions -->

