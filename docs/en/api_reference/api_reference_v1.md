# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [appcallback.proto](#appcallback.proto)
    - [ListTopicSubscriptionsResponse](#spec.proto.runtime.v1.ListTopicSubscriptionsResponse)
    - [TopicEventRequest](#spec.proto.runtime.v1.TopicEventRequest)
    - [TopicEventResponse](#spec.proto.runtime.v1.TopicEventResponse)
    - [TopicSubscription](#spec.proto.runtime.v1.TopicSubscription)
    - [TopicSubscription.MetadataEntry](#spec.proto.runtime.v1.TopicSubscription.MetadataEntry)
  
    - [TopicEventResponse.TopicEventResponseStatus](#spec.proto.runtime.v1.TopicEventResponse.TopicEventResponseStatus)
  
    - [AppCallback](#spec.proto.runtime.v1.AppCallback)
  
- [runtime.proto](#runtime.proto)
    - [BulkStateItem](#spec.proto.runtime.v1.BulkStateItem)
    - [BulkStateItem.MetadataEntry](#spec.proto.runtime.v1.BulkStateItem.MetadataEntry)
    - [CommonInvokeRequest](#spec.proto.runtime.v1.CommonInvokeRequest)
    - [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem)
    - [ConfigurationItem.MetadataEntry](#spec.proto.runtime.v1.ConfigurationItem.MetadataEntry)
    - [ConfigurationItem.TagsEntry](#spec.proto.runtime.v1.ConfigurationItem.TagsEntry)
    - [DelFileRequest](#spec.proto.runtime.v1.DelFileRequest)
    - [DeleteBulkStateRequest](#spec.proto.runtime.v1.DeleteBulkStateRequest)
    - [DeleteConfigurationRequest](#spec.proto.runtime.v1.DeleteConfigurationRequest)
    - [DeleteConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.DeleteConfigurationRequest.MetadataEntry)
    - [DeleteStateRequest](#spec.proto.runtime.v1.DeleteStateRequest)
    - [DeleteStateRequest.MetadataEntry](#spec.proto.runtime.v1.DeleteStateRequest.MetadataEntry)
    - [Etag](#spec.proto.runtime.v1.Etag)
    - [ExecuteStateTransactionRequest](#spec.proto.runtime.v1.ExecuteStateTransactionRequest)
    - [ExecuteStateTransactionRequest.MetadataEntry](#spec.proto.runtime.v1.ExecuteStateTransactionRequest.MetadataEntry)
    - [FileRequest](#spec.proto.runtime.v1.FileRequest)
    - [FileRequest.MetadataEntry](#spec.proto.runtime.v1.FileRequest.MetadataEntry)
    - [GetBulkStateRequest](#spec.proto.runtime.v1.GetBulkStateRequest)
    - [GetBulkStateRequest.MetadataEntry](#spec.proto.runtime.v1.GetBulkStateRequest.MetadataEntry)
    - [GetBulkStateResponse](#spec.proto.runtime.v1.GetBulkStateResponse)
    - [GetConfigurationRequest](#spec.proto.runtime.v1.GetConfigurationRequest)
    - [GetConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.GetConfigurationRequest.MetadataEntry)
    - [GetConfigurationResponse](#spec.proto.runtime.v1.GetConfigurationResponse)
    - [GetFileRequest](#spec.proto.runtime.v1.GetFileRequest)
    - [GetFileRequest.MetadataEntry](#spec.proto.runtime.v1.GetFileRequest.MetadataEntry)
    - [GetFileResponse](#spec.proto.runtime.v1.GetFileResponse)
    - [GetNextIdRequest](#spec.proto.runtime.v1.GetNextIdRequest)
    - [GetNextIdRequest.MetadataEntry](#spec.proto.runtime.v1.GetNextIdRequest.MetadataEntry)
    - [GetNextIdResponse](#spec.proto.runtime.v1.GetNextIdResponse)
    - [GetStateRequest](#spec.proto.runtime.v1.GetStateRequest)
    - [GetStateRequest.MetadataEntry](#spec.proto.runtime.v1.GetStateRequest.MetadataEntry)
    - [GetStateResponse](#spec.proto.runtime.v1.GetStateResponse)
    - [GetStateResponse.MetadataEntry](#spec.proto.runtime.v1.GetStateResponse.MetadataEntry)
    - [HTTPExtension](#spec.proto.runtime.v1.HTTPExtension)
    - [InvokeBindingRequest](#spec.proto.runtime.v1.InvokeBindingRequest)
    - [InvokeBindingRequest.MetadataEntry](#spec.proto.runtime.v1.InvokeBindingRequest.MetadataEntry)
    - [InvokeBindingResponse](#spec.proto.runtime.v1.InvokeBindingResponse)
    - [InvokeBindingResponse.MetadataEntry](#spec.proto.runtime.v1.InvokeBindingResponse.MetadataEntry)
    - [InvokeResponse](#spec.proto.runtime.v1.InvokeResponse)
    - [InvokeServiceRequest](#spec.proto.runtime.v1.InvokeServiceRequest)
    - [ListFileRequest](#spec.proto.runtime.v1.ListFileRequest)
    - [ListFileResp](#spec.proto.runtime.v1.ListFileResp)
    - [PublishEventRequest](#spec.proto.runtime.v1.PublishEventRequest)
    - [PublishEventRequest.MetadataEntry](#spec.proto.runtime.v1.PublishEventRequest.MetadataEntry)
    - [PutFileRequest](#spec.proto.runtime.v1.PutFileRequest)
    - [PutFileRequest.MetadataEntry](#spec.proto.runtime.v1.PutFileRequest.MetadataEntry)
    - [SaveConfigurationRequest](#spec.proto.runtime.v1.SaveConfigurationRequest)
    - [SaveConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.SaveConfigurationRequest.MetadataEntry)
    - [SaveStateRequest](#spec.proto.runtime.v1.SaveStateRequest)
    - [SayHelloRequest](#spec.proto.runtime.v1.SayHelloRequest)
    - [SayHelloResponse](#spec.proto.runtime.v1.SayHelloResponse)
    - [SequencerOptions](#spec.proto.runtime.v1.SequencerOptions)
    - [StateItem](#spec.proto.runtime.v1.StateItem)
    - [StateItem.MetadataEntry](#spec.proto.runtime.v1.StateItem.MetadataEntry)
    - [StateOptions](#spec.proto.runtime.v1.StateOptions)
    - [SubscribeConfigurationRequest](#spec.proto.runtime.v1.SubscribeConfigurationRequest)
    - [SubscribeConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.SubscribeConfigurationRequest.MetadataEntry)
    - [SubscribeConfigurationResponse](#spec.proto.runtime.v1.SubscribeConfigurationResponse)
    - [TransactionalStateOperation](#spec.proto.runtime.v1.TransactionalStateOperation)
    - [TryLockRequest](#spec.proto.runtime.v1.TryLockRequest)
    - [TryLockResponse](#spec.proto.runtime.v1.TryLockResponse)
    - [UnlockRequest](#spec.proto.runtime.v1.UnlockRequest)
    - [UnlockResponse](#spec.proto.runtime.v1.UnlockResponse)
  
    - [HTTPExtension.Verb](#spec.proto.runtime.v1.HTTPExtension.Verb)
    - [SequencerOptions.AutoIncrement](#spec.proto.runtime.v1.SequencerOptions.AutoIncrement)
    - [StateOptions.StateConcurrency](#spec.proto.runtime.v1.StateOptions.StateConcurrency)
    - [StateOptions.StateConsistency](#spec.proto.runtime.v1.StateOptions.StateConsistency)
    - [UnlockResponse.Status](#spec.proto.runtime.v1.UnlockResponse.Status)
  
    - [Runtime](#spec.proto.runtime.v1.Runtime)
  
- [Scalar Value Types](#scalar-value-types)



<a name="appcallback.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## appcallback.proto



<a name="spec.proto.runtime.v1.ListTopicSubscriptionsResponse"></a>

### ListTopicSubscriptionsResponse
ListTopicSubscriptionsResponse is the message including the list of the subscribing topics.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| subscriptions | [TopicSubscription](#spec.proto.runtime.v1.TopicSubscription) | repeated | The list of topics. |






<a name="spec.proto.runtime.v1.TopicEventRequest"></a>

### TopicEventRequest
TopicEventRequest message is compatible with CloudEvent spec v1.0
https://github.com/cloudevents/spec/blob/v1.0/spec.md


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | id identifies the event. Producers MUST ensure that source &#43; id is unique for each distinct event. If a duplicate event is re-sent (e.g. due to a network error) it MAY have the same id. |
| source | [string](#string) |  | source identifies the context in which an event happened. Often this will include information such as the type of the event source, the organization publishing the event or the process that produced the event. The exact syntax and semantics behind the data encoded in the URI is defined by the event producer. |
| type | [string](#string) |  | The type of event related to the originating occurrence. |
| spec_version | [string](#string) |  | The version of the CloudEvents specification. |
| data_content_type | [string](#string) |  | The content type of data value. |
| data | [bytes](#bytes) |  | The content of the event. |
| topic | [string](#string) |  | The pubsub topic which publisher sent to. |
| pubsub_name | [string](#string) |  | The name of the pubsub the publisher sent to. |






<a name="spec.proto.runtime.v1.TopicEventResponse"></a>

### TopicEventResponse
TopicEventResponse is response from app on published message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [TopicEventResponse.TopicEventResponseStatus](#spec.proto.runtime.v1.TopicEventResponse.TopicEventResponseStatus) |  | The list of output bindings. |






<a name="spec.proto.runtime.v1.TopicSubscription"></a>

### TopicSubscription
TopicSubscription represents topic and metadata.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pubsub_name | [string](#string) |  | Required. The name of the pubsub containing the topic below to subscribe to. |
| topic | [string](#string) |  | Required. The name of topic which will be subscribed |
| metadata | [TopicSubscription.MetadataEntry](#spec.proto.runtime.v1.TopicSubscription.MetadataEntry) | repeated | The optional properties used for this topic&#39;s subscription e.g. session id |






<a name="spec.proto.runtime.v1.TopicSubscription.MetadataEntry"></a>

### TopicSubscription.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |





 


<a name="spec.proto.runtime.v1.TopicEventResponse.TopicEventResponseStatus"></a>

### TopicEventResponse.TopicEventResponseStatus
TopicEventResponseStatus allows apps to have finer control over handling of the message.

| Name | Number | Description |
| ---- | ------ | ----------- |
| SUCCESS | 0 | SUCCESS is the default behavior: message is acknowledged and not retried or logged. |
| RETRY | 1 | RETRY status signals runtime to retry the message as part of an expected scenario (no warning is logged). |
| DROP | 2 | DROP status signals runtime to drop the message as part of an unexpected scenario (warning is logged). |


 

 


<a name="spec.proto.runtime.v1.AppCallback"></a>

### AppCallback
AppCallback V1 allows user application to interact with runtime.
User application needs to implement AppCallback service if it needs to
receive message from runtime.

// Invokes service method with InvokeRequest.
 rpc OnInvoke (InvokeRequest) returns (InvokeResponse) {}

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListTopicSubscriptions | [.google.protobuf.Empty](#google.protobuf.Empty) | [ListTopicSubscriptionsResponse](#spec.proto.runtime.v1.ListTopicSubscriptionsResponse) | Lists all topics subscribed by this app. |
| OnTopicEvent | [TopicEventRequest](#spec.proto.runtime.v1.TopicEventRequest) | [TopicEventResponse](#spec.proto.runtime.v1.TopicEventResponse) | Subscribes events from Pubsub |

 



<a name="runtime.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## runtime.proto



<a name="spec.proto.runtime.v1.BulkStateItem"></a>

### BulkStateItem
BulkStateItem is the response item for a bulk get operation.
Return values include the item key, data and etag.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | state item key |
| data | [bytes](#bytes) |  | The byte array data |
| etag | [string](#string) |  | The entity tag which represents the specific version of data. ETag format is defined by the corresponding data store. |
| error | [string](#string) |  | The error that was returned from the state store in case of a failed get operation. |
| metadata | [BulkStateItem.MetadataEntry](#spec.proto.runtime.v1.BulkStateItem.MetadataEntry) | repeated | The metadata which will be sent to app. |






<a name="spec.proto.runtime.v1.BulkStateItem.MetadataEntry"></a>

### BulkStateItem.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.CommonInvokeRequest"></a>

### CommonInvokeRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  |  |
| data | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| content_type | [string](#string) |  |  |
| http_extension | [HTTPExtension](#spec.proto.runtime.v1.HTTPExtension) |  |  |






<a name="spec.proto.runtime.v1.ConfigurationItem"></a>

### ConfigurationItem
ConfigurationItem represents a configuration item with key, content and other information.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Required. The key of configuration item |
| content | [string](#string) |  | The content of configuration item Empty if the configuration is not set, including the case that the configuration is changed from value-set to value-not-set. |
| group | [string](#string) |  | The group of configuration item. |
| label | [string](#string) |  | The label of configuration item. |
| tags | [ConfigurationItem.TagsEntry](#spec.proto.runtime.v1.ConfigurationItem.TagsEntry) | repeated | The tag list of configuration item. |
| metadata | [ConfigurationItem.MetadataEntry](#spec.proto.runtime.v1.ConfigurationItem.MetadataEntry) | repeated | The metadata which will be passed to configuration store component. |






<a name="spec.proto.runtime.v1.ConfigurationItem.MetadataEntry"></a>

### ConfigurationItem.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.ConfigurationItem.TagsEntry"></a>

### ConfigurationItem.TagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.DelFileRequest"></a>

### DelFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [FileRequest](#spec.proto.runtime.v1.FileRequest) |  |  |






<a name="spec.proto.runtime.v1.DeleteBulkStateRequest"></a>

### DeleteBulkStateRequest
DeleteBulkStateRequest is the message to delete a list of key-value states from specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| states | [StateItem](#spec.proto.runtime.v1.StateItem) | repeated | Required. The array of the state key values. |






<a name="spec.proto.runtime.v1.DeleteConfigurationRequest"></a>

### DeleteConfigurationRequest
DeleteConfigurationRequest is the message to delete a list of key-value configuration from specified configuration store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id which Only used for admin, Ignored and reset for normal client |
| group | [string](#string) |  | The group of keys. |
| label | [string](#string) |  | The label for keys. |
| keys | [string](#string) | repeated | The keys to get. |
| metadata | [DeleteConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.DeleteConfigurationRequest.MetadataEntry) | repeated | The metadata which will be sent to configuration store components. |






<a name="spec.proto.runtime.v1.DeleteConfigurationRequest.MetadataEntry"></a>

### DeleteConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.DeleteStateRequest"></a>

### DeleteStateRequest
DeleteStateRequest is the message to delete key-value states in the specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| key | [string](#string) |  | Required. The key of the desired state |
| etag | [Etag](#spec.proto.runtime.v1.Etag) |  | (optional) The entity tag which represents the specific version of data. The exact ETag format is defined by the corresponding data store. |
| options | [StateOptions](#spec.proto.runtime.v1.StateOptions) |  | (optional) State operation options which includes concurrency/ consistency/retry_policy. |
| metadata | [DeleteStateRequest.MetadataEntry](#spec.proto.runtime.v1.DeleteStateRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to state store components. |






<a name="spec.proto.runtime.v1.DeleteStateRequest.MetadataEntry"></a>

### DeleteStateRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.Etag"></a>

### Etag
Etag represents a state item version


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  | value sets the etag value |






<a name="spec.proto.runtime.v1.ExecuteStateTransactionRequest"></a>

### ExecuteStateTransactionRequest
ExecuteStateTransactionRequest is the message to execute multiple operations on a specified store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| storeName | [string](#string) |  | Required. name of state store. |
| operations | [TransactionalStateOperation](#spec.proto.runtime.v1.TransactionalStateOperation) | repeated | Required. transactional operation list. |
| metadata | [ExecuteStateTransactionRequest.MetadataEntry](#spec.proto.runtime.v1.ExecuteStateTransactionRequest.MetadataEntry) | repeated | (optional) The metadata used for transactional operations. |






<a name="spec.proto.runtime.v1.ExecuteStateTransactionRequest.MetadataEntry"></a>

### ExecuteStateTransactionRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.FileRequest"></a>

### FileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  |  |
| name | [string](#string) |  | The name of the directory |
| metadata | [FileRequest.MetadataEntry](#spec.proto.runtime.v1.FileRequest.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.FileRequest.MetadataEntry"></a>

### FileRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetBulkStateRequest"></a>

### GetBulkStateRequest
GetBulkStateRequest is the message to get a list of key-value states from specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| keys | [string](#string) | repeated | Required. The keys to get. |
| parallelism | [int32](#int32) |  | (optional) The number of parallel operations executed on the state store for a get operation. |
| metadata | [GetBulkStateRequest.MetadataEntry](#spec.proto.runtime.v1.GetBulkStateRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to state store components. |






<a name="spec.proto.runtime.v1.GetBulkStateRequest.MetadataEntry"></a>

### GetBulkStateRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetBulkStateResponse"></a>

### GetBulkStateResponse
GetBulkStateResponse is the response conveying the list of state values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [BulkStateItem](#spec.proto.runtime.v1.BulkStateItem) | repeated | The list of items containing the keys to get values for. |






<a name="spec.proto.runtime.v1.GetConfigurationRequest"></a>

### GetConfigurationRequest
GetConfigurationRequest is the message to get a list of key-value configuration from specified configuration store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id which Only used for admin, Ignored and reset for normal client |
| group | [string](#string) |  | The group of keys. |
| label | [string](#string) |  | The label for keys. |
| keys | [string](#string) | repeated | The keys to get. |
| metadata | [GetConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.GetConfigurationRequest.MetadataEntry) | repeated | The metadata which will be sent to configuration store components. |
| subscribe_update | [bool](#bool) |  | Subscribes update event for given keys. If true, when any configuration item in this request is updated, app will receive event by OnConfigurationEvent() of app callback |






<a name="spec.proto.runtime.v1.GetConfigurationRequest.MetadataEntry"></a>

### GetConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetConfigurationResponse"></a>

### GetConfigurationResponse
GetConfigurationResponse is the response conveying the list of configuration values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem) | repeated | The list of items containing configuration values. |






<a name="spec.proto.runtime.v1.GetFileRequest"></a>

### GetFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  |  |
| name | [string](#string) |  | The name of the file or object want to get. |
| metadata | [GetFileRequest.MetadataEntry](#spec.proto.runtime.v1.GetFileRequest.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.GetFileRequest.MetadataEntry"></a>

### GetFileRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetFileResponse"></a>

### GetFileResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  |  |






<a name="spec.proto.runtime.v1.GetNextIdRequest"></a>

### GetNextIdRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. Name of sequencer storage |
| key | [string](#string) |  | Required. key is the identifier of a sequencer namespace,e.g. &#34;order_table&#34;. |
| options | [SequencerOptions](#spec.proto.runtime.v1.SequencerOptions) |  | (optional) SequencerOptions configures requirements for auto-increment guarantee |
| metadata | [GetNextIdRequest.MetadataEntry](#spec.proto.runtime.v1.GetNextIdRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to the component. |






<a name="spec.proto.runtime.v1.GetNextIdRequest.MetadataEntry"></a>

### GetNextIdRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetNextIdResponse"></a>

### GetNextIdResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| next_id | [int64](#int64) |  | The next unique id |






<a name="spec.proto.runtime.v1.GetStateRequest"></a>

### GetStateRequest
GetStateRequest is the message to get key-value states from specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| key | [string](#string) |  | Required. The key of the desired state |
| consistency | [StateOptions.StateConsistency](#spec.proto.runtime.v1.StateOptions.StateConsistency) |  | (optional) read consistency mode |
| metadata | [GetStateRequest.MetadataEntry](#spec.proto.runtime.v1.GetStateRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to state store components. |






<a name="spec.proto.runtime.v1.GetStateRequest.MetadataEntry"></a>

### GetStateRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetStateResponse"></a>

### GetStateResponse
GetStateResponse is the response conveying the state value and etag.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The byte array data |
| etag | [string](#string) |  | The entity tag which represents the specific version of data. ETag format is defined by the corresponding data store. |
| metadata | [GetStateResponse.MetadataEntry](#spec.proto.runtime.v1.GetStateResponse.MetadataEntry) | repeated | The metadata which will be sent to app. |






<a name="spec.proto.runtime.v1.GetStateResponse.MetadataEntry"></a>

### GetStateResponse.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.HTTPExtension"></a>

### HTTPExtension



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| verb | [HTTPExtension.Verb](#spec.proto.runtime.v1.HTTPExtension.Verb) |  |  |
| querystring | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.InvokeBindingRequest"></a>

### InvokeBindingRequest
InvokeBindingRequest is the message to send data to output bindings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the output binding to invoke. |
| data | [bytes](#bytes) |  | The data which will be sent to output binding. |
| metadata | [InvokeBindingRequest.MetadataEntry](#spec.proto.runtime.v1.InvokeBindingRequest.MetadataEntry) | repeated | The metadata passing to output binding components Common metadata property: - ttlInSeconds : the time to live in seconds for the message. If set in the binding definition will cause all messages to have a default time to live. The message ttl overrides any value in the binding definition. |
| operation | [string](#string) |  | The name of the operation type for the binding to invoke |






<a name="spec.proto.runtime.v1.InvokeBindingRequest.MetadataEntry"></a>

### InvokeBindingRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.InvokeBindingResponse"></a>

### InvokeBindingResponse
InvokeBindingResponse is the message returned from an output binding invocation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data which will be sent to output binding. |
| metadata | [InvokeBindingResponse.MetadataEntry](#spec.proto.runtime.v1.InvokeBindingResponse.MetadataEntry) | repeated | The metadata returned from an external system |






<a name="spec.proto.runtime.v1.InvokeBindingResponse.MetadataEntry"></a>

### InvokeBindingResponse.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.InvokeResponse"></a>

### InvokeResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| content_type | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.InvokeServiceRequest"></a>

### InvokeServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  |  |
| message | [CommonInvokeRequest](#spec.proto.runtime.v1.CommonInvokeRequest) |  |  |






<a name="spec.proto.runtime.v1.ListFileRequest"></a>

### ListFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [FileRequest](#spec.proto.runtime.v1.FileRequest) |  |  |






<a name="spec.proto.runtime.v1.ListFileResp"></a>

### ListFileResp



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file_name | [string](#string) | repeated |  |






<a name="spec.proto.runtime.v1.PublishEventRequest"></a>

### PublishEventRequest
PublishEventRequest is the message to publish event data to pubsub topic


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| pubsub_name | [string](#string) |  | The name of the pubsub component |
| topic | [string](#string) |  | The pubsub topic |
| data | [bytes](#bytes) |  | The data which will be published to topic. |
| data_content_type | [string](#string) |  | The content type for the data (optional). |
| metadata | [PublishEventRequest.MetadataEntry](#spec.proto.runtime.v1.PublishEventRequest.MetadataEntry) | repeated | The metadata passing to pub components

metadata property: - key : the key of the message. |






<a name="spec.proto.runtime.v1.PublishEventRequest.MetadataEntry"></a>

### PublishEventRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.PutFileRequest"></a>

### PutFileRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  |  |
| name | [string](#string) |  | The name of the file or object want to put. |
| data | [bytes](#bytes) |  | The data will be store. |
| metadata | [PutFileRequest.MetadataEntry](#spec.proto.runtime.v1.PutFileRequest.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.PutFileRequest.MetadataEntry"></a>

### PutFileRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SaveConfigurationRequest"></a>

### SaveConfigurationRequest
SaveConfigurationRequest is the message to save a list of key-value configuration into specified configuration store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id which Only used for admin, ignored and reset for normal client |
| items | [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem) | repeated | The list of configuration items to save. To delete a exist item, set the key (also label) and let content to be empty |
| metadata | [SaveConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.SaveConfigurationRequest.MetadataEntry) | repeated | The metadata which will be sent to configuration store components. |






<a name="spec.proto.runtime.v1.SaveConfigurationRequest.MetadataEntry"></a>

### SaveConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SaveStateRequest"></a>

### SaveStateRequest
SaveStateRequest is the message to save multiple states into state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| states | [StateItem](#spec.proto.runtime.v1.StateItem) | repeated | Required. The array of the state key values. |






<a name="spec.proto.runtime.v1.SayHelloRequest"></a>

### SayHelloRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_name | [string](#string) |  |  |
| name | [string](#string) |  |  |
| data | [google.protobuf.Any](#google.protobuf.Any) |  | Optional. This field is used to control the packet size during load tests. |






<a name="spec.proto.runtime.v1.SayHelloResponse"></a>

### SayHelloResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hello | [string](#string) |  |  |
| data | [google.protobuf.Any](#google.protobuf.Any) |  |  |






<a name="spec.proto.runtime.v1.SequencerOptions"></a>

### SequencerOptions
SequencerOptions configures requirements for auto-increment guarantee


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| increment | [SequencerOptions.AutoIncrement](#spec.proto.runtime.v1.SequencerOptions.AutoIncrement) |  |  |






<a name="spec.proto.runtime.v1.StateItem"></a>

### StateItem
StateItem represents state key, value, and additional options to save state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Required. The state key |
| value | [bytes](#bytes) |  | Required. The state data for key |
| etag | [Etag](#spec.proto.runtime.v1.Etag) |  | (optional) The entity tag which represents the specific version of data. The exact ETag format is defined by the corresponding data store. Layotto runtime only treats ETags as opaque strings. |
| metadata | [StateItem.MetadataEntry](#spec.proto.runtime.v1.StateItem.MetadataEntry) | repeated | (optional) additional key-value pairs to be passed to the state store. |
| options | [StateOptions](#spec.proto.runtime.v1.StateOptions) |  | (optional) Options for concurrency and consistency to save the state. |






<a name="spec.proto.runtime.v1.StateItem.MetadataEntry"></a>

### StateItem.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.StateOptions"></a>

### StateOptions
StateOptions configures concurrency and consistency for state operations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| concurrency | [StateOptions.StateConcurrency](#spec.proto.runtime.v1.StateOptions.StateConcurrency) |  |  |
| consistency | [StateOptions.StateConsistency](#spec.proto.runtime.v1.StateOptions.StateConsistency) |  |  |






<a name="spec.proto.runtime.v1.SubscribeConfigurationRequest"></a>

### SubscribeConfigurationRequest
SubscribeConfigurationRequest is the message to get a list of key-value configuration from specified configuration store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id which Only used for admin, ignored and reset for normal client |
| group | [string](#string) |  | The group of keys. |
| label | [string](#string) |  | The label for keys. |
| keys | [string](#string) | repeated | The keys to get. |
| metadata | [SubscribeConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.SubscribeConfigurationRequest.MetadataEntry) | repeated | The metadata which will be sent to configuration store components. |






<a name="spec.proto.runtime.v1.SubscribeConfigurationRequest.MetadataEntry"></a>

### SubscribeConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SubscribeConfigurationResponse"></a>

### SubscribeConfigurationResponse
SubscribeConfigurationResponse is the response conveying the list of configuration values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id. Only used for admin client. |
| items | [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem) | repeated | The list of items containing configuration values. |






<a name="spec.proto.runtime.v1.TransactionalStateOperation"></a>

### TransactionalStateOperation
TransactionalStateOperation is the message to execute a specified operation with a key-value pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operationType | [string](#string) |  | Required. The type of operation to be executed |
| request | [StateItem](#spec.proto.runtime.v1.StateItem) |  | Required. State values to be operated on |






<a name="spec.proto.runtime.v1.TryLockRequest"></a>

### TryLockRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The lock store name,e.g. `redis`. |
| resource_id | [string](#string) |  | Required. resource_id is the lock key. e.g. `order_id_111` It stands for &#34;which resource I want to protect&#34; |
| lock_owner | [string](#string) |  | Required. lock_owner indicate the identifier of lock owner. You can generate a uuid as lock_owner.For example,in golang: req.LockOwner = uuid.New().String() This field is per request,not per process,so it is different for each request, which aims to prevent multi-thread in the same process trying the same lock concurrently. The reason why we don&#39;t make it automatically generated is: 1. If it is automatically generated,there must be a &#39;my_lock_owner_id&#39; field in the response. This name is so weird that we think it is inappropriate to put it into the api spec 2. If we change the field &#39;my_lock_owner_id&#39; in the response to &#39;lock_owner&#39;,which means the current lock owner of this lock, we find that in some lock services users can&#39;t get the current lock owner.Actually users don&#39;t need it at all. 3. When reentrant lock is needed,the existing lock_owner is required to identify client and check &#34;whether this client can reenter this lock&#34;. So this field in the request shouldn&#39;t be removed. |
| expire | [int32](#int32) |  | Required. expire is the time before expire.The time unit is second. |






<a name="spec.proto.runtime.v1.TryLockResponse"></a>

### TryLockResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  |  |






<a name="spec.proto.runtime.v1.UnlockRequest"></a>

### UnlockRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  |  |
| resource_id | [string](#string) |  | resource_id is the lock key. |
| lock_owner | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.UnlockResponse"></a>

### UnlockResponse



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [UnlockResponse.Status](#spec.proto.runtime.v1.UnlockResponse.Status) |  |  |





 


<a name="spec.proto.runtime.v1.HTTPExtension.Verb"></a>

### HTTPExtension.Verb


| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 |  |
| GET | 1 |  |
| HEAD | 2 |  |
| POST | 3 |  |
| PUT | 4 |  |
| DELETE | 5 |  |
| CONNECT | 6 |  |
| OPTIONS | 7 |  |
| TRACE | 8 |  |



<a name="spec.proto.runtime.v1.SequencerOptions.AutoIncrement"></a>

### SequencerOptions.AutoIncrement
requirements for auto-increment guarantee

| Name | Number | Description |
| ---- | ------ | ----------- |
| WEAK | 0 | (default) WEAK means a &#34;best effort&#34; incrementing service.But there is no strict guarantee of global monotonically increasing. The next id is &#34;probably&#34; greater than current id. |
| STRONG | 1 | STRONG means a strict guarantee of global monotonically increasing. The next id &#34;must&#34; be greater than current id. |



<a name="spec.proto.runtime.v1.StateOptions.StateConcurrency"></a>

### StateOptions.StateConcurrency
Enum describing the supported concurrency for state.
The API server uses Optimized Concurrency Control (OCC) with ETags.
When an ETag is associated with an save or delete request, the store shall allow the update only if the attached ETag matches with the latest ETag in the database.
But when ETag is missing in the write requests, the state store shall handle the requests in the specified strategy(e.g. a last-write-wins fashion).

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONCURRENCY_UNSPECIFIED | 0 |  |
| CONCURRENCY_FIRST_WRITE | 1 | First write wins |
| CONCURRENCY_LAST_WRITE | 2 | Last write wins |



<a name="spec.proto.runtime.v1.StateOptions.StateConsistency"></a>

### StateOptions.StateConsistency
Enum describing the supported consistency for state.

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONSISTENCY_UNSPECIFIED | 0 |  |
| CONSISTENCY_EVENTUAL | 1 | The API server assumes data stores are eventually consistent by default.A state store should: - For read requests, the state store can return data from any of the replicas - For write request, the state store should asynchronously replicate updates to configured quorum after acknowledging the update request. |
| CONSISTENCY_STRONG | 2 | When a strong consistency hint is attached, a state store should: - For read requests, the state store should return the most up-to-date data consistently across replicas. - For write/delete requests, the state store should synchronisely replicate updated data to configured quorum before completing the write request. |



<a name="spec.proto.runtime.v1.UnlockResponse.Status"></a>

### UnlockResponse.Status


| Name | Number | Description |
| ---- | ------ | ----------- |
| SUCCESS | 0 |  |
| LOCK_UNEXIST | 1 |  |
| LOCK_BELONG_TO_OTHERS | 2 |  |
| INTERNAL_ERROR | 3 |  |


 

 


<a name="spec.proto.runtime.v1.Runtime"></a>

### Runtime


| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SayHello | [SayHelloRequest](#spec.proto.runtime.v1.SayHelloRequest) | [SayHelloResponse](#spec.proto.runtime.v1.SayHelloResponse) | SayHello used for test |
| InvokeService | [InvokeServiceRequest](#spec.proto.runtime.v1.InvokeServiceRequest) | [InvokeResponse](#spec.proto.runtime.v1.InvokeResponse) | InvokeService do rpc calls |
| GetConfiguration | [GetConfigurationRequest](#spec.proto.runtime.v1.GetConfigurationRequest) | [GetConfigurationResponse](#spec.proto.runtime.v1.GetConfigurationResponse) | GetConfiguration gets configuration from configuration store. |
| SaveConfiguration | [SaveConfigurationRequest](#spec.proto.runtime.v1.SaveConfigurationRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | SaveConfiguration saves configuration into configuration store. |
| DeleteConfiguration | [DeleteConfigurationRequest](#spec.proto.runtime.v1.DeleteConfigurationRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | DeleteConfiguration deletes configuration from configuration store. |
| SubscribeConfiguration | [SubscribeConfigurationRequest](#spec.proto.runtime.v1.SubscribeConfigurationRequest) stream | [SubscribeConfigurationResponse](#spec.proto.runtime.v1.SubscribeConfigurationResponse) stream | SubscribeConfiguration gets configuration from configuration store and subscribe the updates. |
| TryLock | [TryLockRequest](#spec.proto.runtime.v1.TryLockRequest) | [TryLockResponse](#spec.proto.runtime.v1.TryLockResponse) | Distributed Lock API A non-blocking method trying to get a lock with ttl. |
| Unlock | [UnlockRequest](#spec.proto.runtime.v1.UnlockRequest) | [UnlockResponse](#spec.proto.runtime.v1.UnlockResponse) |  |
| GetNextId | [GetNextIdRequest](#spec.proto.runtime.v1.GetNextIdRequest) | [GetNextIdResponse](#spec.proto.runtime.v1.GetNextIdResponse) | Sequencer API Get next unique id with some auto-increment guarantee |
| GetState | [GetStateRequest](#spec.proto.runtime.v1.GetStateRequest) | [GetStateResponse](#spec.proto.runtime.v1.GetStateResponse) | Gets the state for a specific key. |
| GetBulkState | [GetBulkStateRequest](#spec.proto.runtime.v1.GetBulkStateRequest) | [GetBulkStateResponse](#spec.proto.runtime.v1.GetBulkStateResponse) | Gets a bulk of state items for a list of keys |
| SaveState | [SaveStateRequest](#spec.proto.runtime.v1.SaveStateRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | Saves an array of state objects |
| DeleteState | [DeleteStateRequest](#spec.proto.runtime.v1.DeleteStateRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | Deletes the state for a specific key. |
| DeleteBulkState | [DeleteBulkStateRequest](#spec.proto.runtime.v1.DeleteBulkStateRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | Deletes a bulk of state items for a list of keys |
| ExecuteStateTransaction | [ExecuteStateTransactionRequest](#spec.proto.runtime.v1.ExecuteStateTransactionRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | Executes transactions for a specified store |
| PublishEvent | [PublishEventRequest](#spec.proto.runtime.v1.PublishEventRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | Publishes events to the specific topic |
| GetFile | [GetFileRequest](#spec.proto.runtime.v1.GetFileRequest) | [GetFileResponse](#spec.proto.runtime.v1.GetFileResponse) stream | Get file with stream |
| PutFile | [PutFileRequest](#spec.proto.runtime.v1.PutFileRequest) stream | [.google.protobuf.Empty](#google.protobuf.Empty) | Put file with stream |
| ListFile | [ListFileRequest](#spec.proto.runtime.v1.ListFileRequest) | [ListFileResp](#spec.proto.runtime.v1.ListFileResp) | List all files |
| DelFile | [DelFileRequest](#spec.proto.runtime.v1.DelFileRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | Delete specific file |
| InvokeBinding | [InvokeBindingRequest](#spec.proto.runtime.v1.InvokeBindingRequest) | [InvokeBindingResponse](#spec.proto.runtime.v1.InvokeBindingResponse) | Invokes binding data to specific output bindings |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

