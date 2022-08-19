

<a name="runtime.proto"></a>

# runtime.proto
<a name="top"></a>

This document is automaticallly generated from the [`.proto`](https://github.com/mosn/layotto/tree/main/spec/proto/runtime/v1) files.




<a name="spec.proto.runtime.v1.Runtime"></a>

## [gRPC Service] Runtime
Runtime encapsulates variours Runtime APIs(such as Configuration API, Pub/Sub API, etc)

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| SayHello | [SayHelloRequest](#spec.proto.runtime.v1.SayHelloRequest) | [SayHelloResponse](#spec.proto.runtime.v1.SayHelloResponse) | SayHello used for test |
| InvokeService | [InvokeServiceRequest](#spec.proto.runtime.v1.InvokeServiceRequest) | [InvokeResponse](#spec.proto.runtime.v1.InvokeResponse) | InvokeService do rpc calls |
| GetConfiguration | [GetConfigurationRequest](#spec.proto.runtime.v1.GetConfigurationRequest) | [GetConfigurationResponse](#spec.proto.runtime.v1.GetConfigurationResponse) | GetConfiguration gets configuration from configuration store. |
| SaveConfiguration | [SaveConfigurationRequest](#spec.proto.runtime.v1.SaveConfigurationRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | SaveConfiguration saves configuration into configuration store. |
| DeleteConfiguration | [DeleteConfigurationRequest](#spec.proto.runtime.v1.DeleteConfigurationRequest) | [.google.protobuf.Empty](#google.protobuf.Empty) | DeleteConfiguration deletes configuration from configuration store. |
| SubscribeConfiguration | [SubscribeConfigurationRequest](#spec.proto.runtime.v1.SubscribeConfigurationRequest) stream | [SubscribeConfigurationResponse](#spec.proto.runtime.v1.SubscribeConfigurationResponse) stream | SubscribeConfiguration gets configuration from configuration store and subscribe the updates. |
| TryLock | [TryLockRequest](#spec.proto.runtime.v1.TryLockRequest) | [TryLockResponse](#spec.proto.runtime.v1.TryLockResponse) | Distributed Lock API A non-blocking method trying to get a lock with ttl. |
| Unlock | [UnlockRequest](#spec.proto.runtime.v1.UnlockRequest) | [UnlockResponse](#spec.proto.runtime.v1.UnlockResponse) | A method trying to unlock. |
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
| GetFileMeta | [GetFileMetaRequest](#spec.proto.runtime.v1.GetFileMetaRequest) | [GetFileMetaResponse](#spec.proto.runtime.v1.GetFileMetaResponse) | Get file meta data, if file not exist,return code.NotFound error |
| InvokeBinding | [InvokeBindingRequest](#spec.proto.runtime.v1.InvokeBindingRequest) | [InvokeBindingResponse](#spec.proto.runtime.v1.InvokeBindingResponse) | Invokes binding data to specific output bindings |
| GetSecret | [GetSecretRequest](#spec.proto.runtime.v1.GetSecretRequest) | [GetSecretResponse](#spec.proto.runtime.v1.GetSecretResponse) | Gets secrets from secret stores. |
| GetBulkSecret | [GetBulkSecretRequest](#spec.proto.runtime.v1.GetBulkSecretRequest) | [GetBulkSecretResponse](#spec.proto.runtime.v1.GetBulkSecretResponse) | Gets a bulk of secrets |

 <!-- end services -->


<a name="spec.proto.runtime.v1.BulkStateItem"></a>
<p align="right"><a href="#top">Top</a></p>

## BulkStateItem
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
<p align="right"><a href="#top">Top</a></p>

## BulkStateItem.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.CommonInvokeRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## CommonInvokeRequest
Common invoke request message which includes invoke method and data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| method | [string](#string) |  | The method of requset |
| data | [google.protobuf.Any](#google.protobuf.Any) |  | The request data |
| content_type | [string](#string) |  | The content type of request data |
| http_extension | [HTTPExtension](#spec.proto.runtime.v1.HTTPExtension) |  | The extra information of http |






<a name="spec.proto.runtime.v1.ConfigurationItem"></a>
<p align="right"><a href="#top">Top</a></p>

## ConfigurationItem
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
<p align="right"><a href="#top">Top</a></p>

## ConfigurationItem.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.ConfigurationItem.TagsEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## ConfigurationItem.TagsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.DelFileRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## DelFileRequest
Delete file request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [FileRequest](#spec.proto.runtime.v1.FileRequest) |  | File request |






<a name="spec.proto.runtime.v1.DeleteBulkStateRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteBulkStateRequest
DeleteBulkStateRequest is the message to delete a list of key-value states from specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| states | [StateItem](#spec.proto.runtime.v1.StateItem) | repeated | Required. The array of the state key values. |






<a name="spec.proto.runtime.v1.DeleteConfigurationRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteConfigurationRequest
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
<p align="right"><a href="#top">Top</a></p>

## DeleteConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.DeleteStateRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteStateRequest
DeleteStateRequest is the message to delete key-value states in the specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| key | [string](#string) |  | Required. The key of the desired state |
| etag | [Etag](#spec.proto.runtime.v1.Etag) |  | (optional) The entity tag which represents the specific version of data. The exact ETag format is defined by the corresponding data store. |
| options | [StateOptions](#spec.proto.runtime.v1.StateOptions) |  | (optional) State operation options which includes concurrency/ consistency/retry_policy. |
| metadata | [DeleteStateRequest.MetadataEntry](#spec.proto.runtime.v1.DeleteStateRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to state store components. |






<a name="spec.proto.runtime.v1.DeleteStateRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## DeleteStateRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.Etag"></a>
<p align="right"><a href="#top">Top</a></p>

## Etag
Etag represents a state item version


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) |  | value sets the etag value |






<a name="spec.proto.runtime.v1.ExecuteStateTransactionRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## ExecuteStateTransactionRequest
ExecuteStateTransactionRequest is the message to execute multiple operations on a specified store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| storeName | [string](#string) |  | Required. name of state store. |
| operations | [TransactionalStateOperation](#spec.proto.runtime.v1.TransactionalStateOperation) | repeated | Required. transactional operation list. |
| metadata | [ExecuteStateTransactionRequest.MetadataEntry](#spec.proto.runtime.v1.ExecuteStateTransactionRequest.MetadataEntry) | repeated | (optional) The metadata used for transactional operations. |






<a name="spec.proto.runtime.v1.ExecuteStateTransactionRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## ExecuteStateTransactionRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.FileInfo"></a>
<p align="right"><a href="#top">Top</a></p>

## FileInfo
File info message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| file_name | [string](#string) |  | The name of file |
| size | [int64](#int64) |  | The size of file |
| last_modified | [string](#string) |  | The modified time of file |
| metadata | [FileInfo.MetadataEntry](#spec.proto.runtime.v1.FileInfo.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.FileInfo.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## FileInfo.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.FileMeta"></a>
<p align="right"><a href="#top">Top</a></p>

## FileMeta
A map that store FileMetaValue


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| metadata | [FileMeta.MetadataEntry](#spec.proto.runtime.v1.FileMeta.MetadataEntry) | repeated | A data structure to store metadata |






<a name="spec.proto.runtime.v1.FileMeta.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## FileMeta.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [FileMetaValue](#spec.proto.runtime.v1.FileMetaValue) |  |  |






<a name="spec.proto.runtime.v1.FileMetaValue"></a>
<p align="right"><a href="#top">Top</a></p>

## FileMetaValue
FileMeta value


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| value | [string](#string) | repeated | File meta value |






<a name="spec.proto.runtime.v1.FileRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## FileRequest
File request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of store |
| name | [string](#string) |  | The name of the directory |
| metadata | [FileRequest.MetadataEntry](#spec.proto.runtime.v1.FileRequest.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.FileRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## FileRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetBulkSecretRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkSecretRequest
GetBulkSecretRequest is the message to get the secrets from secret store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of secret store. |
| metadata | [GetBulkSecretRequest.MetadataEntry](#spec.proto.runtime.v1.GetBulkSecretRequest.MetadataEntry) | repeated | The metadata which will be sent to secret store components. |






<a name="spec.proto.runtime.v1.GetBulkSecretRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkSecretRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetBulkSecretResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkSecretResponse
GetBulkSecretResponse is the response message to convey the requested secrets.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [GetBulkSecretResponse.DataEntry](#spec.proto.runtime.v1.GetBulkSecretResponse.DataEntry) | repeated | data hold the secret values. Some secret store, such as kubernetes secret store, can save multiple secrets for single secret key. |






<a name="spec.proto.runtime.v1.GetBulkSecretResponse.DataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkSecretResponse.DataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [SecretResponse](#spec.proto.runtime.v1.SecretResponse) |  |  |






<a name="spec.proto.runtime.v1.GetBulkStateRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkStateRequest
GetBulkStateRequest is the message to get a list of key-value states from specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| keys | [string](#string) | repeated | Required. The keys to get. |
| parallelism | [int32](#int32) |  | (optional) The number of parallel operations executed on the state store for a get operation. |
| metadata | [GetBulkStateRequest.MetadataEntry](#spec.proto.runtime.v1.GetBulkStateRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to state store components. |






<a name="spec.proto.runtime.v1.GetBulkStateRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkStateRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetBulkStateResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetBulkStateResponse
GetBulkStateResponse is the response conveying the list of state values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [BulkStateItem](#spec.proto.runtime.v1.BulkStateItem) | repeated | The list of items containing the keys to get values for. |






<a name="spec.proto.runtime.v1.GetConfigurationRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetConfigurationRequest
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
<p align="right"><a href="#top">Top</a></p>

## GetConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetConfigurationResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetConfigurationResponse
GetConfigurationResponse is the response conveying the list of configuration values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| items | [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem) | repeated | The list of items containing configuration values. |






<a name="spec.proto.runtime.v1.GetFileMetaRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetFileMetaRequest
Get fileMeta request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [FileRequest](#spec.proto.runtime.v1.FileRequest) |  | File meta request |






<a name="spec.proto.runtime.v1.GetFileMetaResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetFileMetaResponse
Get fileMeta response message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| size | [int64](#int64) |  | The size of file |
| last_modified | [string](#string) |  | The modified time of file |
| response | [FileMeta](#spec.proto.runtime.v1.FileMeta) |  | File meta response |






<a name="spec.proto.runtime.v1.GetFileRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetFileRequest
Get file request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of store |
| name | [string](#string) |  | The name of the file or object want to get. |
| metadata | [GetFileRequest.MetadataEntry](#spec.proto.runtime.v1.GetFileRequest.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.GetFileRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetFileRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetFileResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetFileResponse
Get file response message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data of file |






<a name="spec.proto.runtime.v1.GetNextIdRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetNextIdRequest
Get next id request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. Name of sequencer storage |
| key | [string](#string) |  | Required. key is the identifier of a sequencer namespace,e.g. "order_table". |
| options | [SequencerOptions](#spec.proto.runtime.v1.SequencerOptions) |  | (optional) SequencerOptions configures requirements for auto-increment guarantee |
| metadata | [GetNextIdRequest.MetadataEntry](#spec.proto.runtime.v1.GetNextIdRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to the component. |






<a name="spec.proto.runtime.v1.GetNextIdRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetNextIdRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetNextIdResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetNextIdResponse
Get next id response message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| next_id | [int64](#int64) |  | The next unique id Fixed int64 overflow problems on JavaScript https://github.com/improbable-eng/ts-protoc-gen#gotchas |






<a name="spec.proto.runtime.v1.GetSecretRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetSecretRequest
GetSecretRequest is the message to get secret from secret store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of secret store. |
| key | [string](#string) |  | The name of secret key. |
| metadata | [GetSecretRequest.MetadataEntry](#spec.proto.runtime.v1.GetSecretRequest.MetadataEntry) | repeated | The metadata which will be sent to secret store components. Contains version, status, and so on... |






<a name="spec.proto.runtime.v1.GetSecretRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetSecretRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetSecretResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetSecretResponse
GetSecretResponse is the response message to convey the requested secret.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [GetSecretResponse.DataEntry](#spec.proto.runtime.v1.GetSecretResponse.DataEntry) | repeated | data is the secret value. Some secret store, such as kubernetes secret store, can save multiple secrets for single secret key. |






<a name="spec.proto.runtime.v1.GetSecretResponse.DataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetSecretResponse.DataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetStateRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## GetStateRequest
GetStateRequest is the message to get key-value states from specific state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| key | [string](#string) |  | Required. The key of the desired state |
| consistency | [StateOptions.StateConsistency](#spec.proto.runtime.v1.StateOptions.StateConsistency) |  | (optional) read consistency mode |
| metadata | [GetStateRequest.MetadataEntry](#spec.proto.runtime.v1.GetStateRequest.MetadataEntry) | repeated | (optional) The metadata which will be sent to state store components. |






<a name="spec.proto.runtime.v1.GetStateRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetStateRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.GetStateResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## GetStateResponse
GetStateResponse is the response conveying the state value and etag.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The byte array data |
| etag | [string](#string) |  | The entity tag which represents the specific version of data. ETag format is defined by the corresponding data store. |
| metadata | [GetStateResponse.MetadataEntry](#spec.proto.runtime.v1.GetStateResponse.MetadataEntry) | repeated | The metadata which will be sent to app. |






<a name="spec.proto.runtime.v1.GetStateResponse.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## GetStateResponse.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.HTTPExtension"></a>
<p align="right"><a href="#top">Top</a></p>

## HTTPExtension
Http extension message is about invoke http information


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| verb | [HTTPExtension.Verb](#spec.proto.runtime.v1.HTTPExtension.Verb) |  | The method of http reuest |
| querystring | [string](#string) |  | The query information of http |






<a name="spec.proto.runtime.v1.InvokeBindingRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## InvokeBindingRequest
InvokeBindingRequest is the message to send data to output bindings


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | The name of the output binding to invoke. |
| data | [bytes](#bytes) |  | The data which will be sent to output binding. |
| metadata | [InvokeBindingRequest.MetadataEntry](#spec.proto.runtime.v1.InvokeBindingRequest.MetadataEntry) | repeated | The metadata passing to output binding components Common metadata property: - ttlInSeconds : the time to live in seconds for the message. If set in the binding definition will cause all messages to have a default time to live. The message ttl overrides any value in the binding definition. |
| operation | [string](#string) |  | The name of the operation type for the binding to invoke |






<a name="spec.proto.runtime.v1.InvokeBindingRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## InvokeBindingRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.InvokeBindingResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## InvokeBindingResponse
InvokeBindingResponse is the message returned from an output binding invocation


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [bytes](#bytes) |  | The data which will be sent to output binding. |
| metadata | [InvokeBindingResponse.MetadataEntry](#spec.proto.runtime.v1.InvokeBindingResponse.MetadataEntry) | repeated | The metadata returned from an external system |






<a name="spec.proto.runtime.v1.InvokeBindingResponse.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## InvokeBindingResponse.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.InvokeResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## InvokeResponse
Invoke service response message is result of invoke service queset


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | [google.protobuf.Any](#google.protobuf.Any) |  | The response data |
| content_type | [string](#string) |  | The content type of response data |






<a name="spec.proto.runtime.v1.InvokeServiceRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## InvokeServiceRequest
Invoke service request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | [string](#string) |  | The identify of InvokeServiceRequest |
| message | [CommonInvokeRequest](#spec.proto.runtime.v1.CommonInvokeRequest) |  | InvokeServiceRequest message |






<a name="spec.proto.runtime.v1.ListFileRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## ListFileRequest
List file request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| request | [FileRequest](#spec.proto.runtime.v1.FileRequest) |  | File request |
| page_size | [int32](#int32) |  | Page size |
| marker | [string](#string) |  | Marker |






<a name="spec.proto.runtime.v1.ListFileResp"></a>
<p align="right"><a href="#top">Top</a></p>

## ListFileResp
List file response message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| files | [FileInfo](#spec.proto.runtime.v1.FileInfo) | repeated | File info |
| marker | [string](#string) |  | Marker |
| is_truncated | [bool](#bool) |  | Is truncated |






<a name="spec.proto.runtime.v1.PublishEventRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## PublishEventRequest
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
<p align="right"><a href="#top">Top</a></p>

## PublishEventRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.PutFileRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## PutFileRequest
Put file request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of store |
| name | [string](#string) |  | The name of the file or object want to put. |
| data | [bytes](#bytes) |  | The data will be store. |
| metadata | [PutFileRequest.MetadataEntry](#spec.proto.runtime.v1.PutFileRequest.MetadataEntry) | repeated | The metadata for user extension. |






<a name="spec.proto.runtime.v1.PutFileRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## PutFileRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SaveConfigurationRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## SaveConfigurationRequest
SaveConfigurationRequest is the message to save a list of key-value configuration into specified configuration store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id which Only used for admin, ignored and reset for normal client |
| items | [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem) | repeated | The list of configuration items to save. To delete a exist item, set the key (also label) and let content to be empty |
| metadata | [SaveConfigurationRequest.MetadataEntry](#spec.proto.runtime.v1.SaveConfigurationRequest.MetadataEntry) | repeated | The metadata which will be sent to configuration store components. |






<a name="spec.proto.runtime.v1.SaveConfigurationRequest.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## SaveConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SaveStateRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## SaveStateRequest
SaveStateRequest is the message to save multiple states into state store.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The name of state store. |
| states | [StateItem](#spec.proto.runtime.v1.StateItem) | repeated | Required. The array of the state key values. |






<a name="spec.proto.runtime.v1.SayHelloRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## SayHelloRequest
Hello request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service_name | [string](#string) |  | The name of service |
| name | [string](#string) |  | Reuqest name |
| data | [google.protobuf.Any](#google.protobuf.Any) |  | Optional. This field is used to control the packet size during load tests. |






<a name="spec.proto.runtime.v1.SayHelloResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## SayHelloResponse
Hello response message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| hello | [string](#string) |  | Hello |
| data | [google.protobuf.Any](#google.protobuf.Any) |  | Hello message of data |






<a name="spec.proto.runtime.v1.SecretResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## SecretResponse
SecretResponse is a map of decrypted string/string values


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| secrets | [SecretResponse.SecretsEntry](#spec.proto.runtime.v1.SecretResponse.SecretsEntry) | repeated | The data struct of secrets |






<a name="spec.proto.runtime.v1.SecretResponse.SecretsEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## SecretResponse.SecretsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SequencerOptions"></a>
<p align="right"><a href="#top">Top</a></p>

## SequencerOptions
SequencerOptions configures requirements for auto-increment guarantee


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| increment | [SequencerOptions.AutoIncrement](#spec.proto.runtime.v1.SequencerOptions.AutoIncrement) |  | Default STRONG auto-increment |






<a name="spec.proto.runtime.v1.StateItem"></a>
<p align="right"><a href="#top">Top</a></p>

## StateItem
StateItem represents state key, value, and additional options to save state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  | Required. The state key |
| value | [bytes](#bytes) |  | Required. The state data for key |
| etag | [Etag](#spec.proto.runtime.v1.Etag) |  | (optional) The entity tag which represents the specific version of data. The exact ETag format is defined by the corresponding data store. Layotto runtime only treats ETags as opaque strings. |
| metadata | [StateItem.MetadataEntry](#spec.proto.runtime.v1.StateItem.MetadataEntry) | repeated | (optional) additional key-value pairs to be passed to the state store. |
| options | [StateOptions](#spec.proto.runtime.v1.StateOptions) |  | (optional) Options for concurrency and consistency to save the state. |






<a name="spec.proto.runtime.v1.StateItem.MetadataEntry"></a>
<p align="right"><a href="#top">Top</a></p>

## StateItem.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.StateOptions"></a>
<p align="right"><a href="#top">Top</a></p>

## StateOptions
StateOptions configures concurrency and consistency for state operations


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| concurrency | [StateOptions.StateConcurrency](#spec.proto.runtime.v1.StateOptions.StateConcurrency) |  | The state operation of concurrency |
| consistency | [StateOptions.StateConsistency](#spec.proto.runtime.v1.StateOptions.StateConsistency) |  | The state operation of consistency |






<a name="spec.proto.runtime.v1.SubscribeConfigurationRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## SubscribeConfigurationRequest
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
<p align="right"><a href="#top">Top</a></p>

## SubscribeConfigurationRequest.MetadataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [string](#string) |  |  |






<a name="spec.proto.runtime.v1.SubscribeConfigurationResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## SubscribeConfigurationResponse
SubscribeConfigurationResponse is the response conveying the list of configuration values.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of configuration store. |
| app_id | [string](#string) |  | The application id. Only used for admin client. |
| items | [ConfigurationItem](#spec.proto.runtime.v1.ConfigurationItem) | repeated | The list of items containing configuration values. |






<a name="spec.proto.runtime.v1.TransactionalStateOperation"></a>
<p align="right"><a href="#top">Top</a></p>

## TransactionalStateOperation
TransactionalStateOperation is the message to execute a specified operation with a key-value pair.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| operationType | [string](#string) |  | Required. The type of operation to be executed. Legal values include: "upsert" represents an update or create operation "delete" represents a delete operation |
| request | [StateItem](#spec.proto.runtime.v1.StateItem) |  | Required. State values to be operated on |






<a name="spec.proto.runtime.v1.TryLockRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## TryLockRequest
Lock request message is distributed lock API which is not blocking method tring to get a lock with ttl


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | Required. The lock store name,e.g. `redis`. |
| resource_id | [string](#string) |  | Required. resource_id is the lock key. e.g. `order_id_111` It stands for "which resource I want to protect" |
| lock_owner | [string](#string) |  | Required. lock_owner indicate the identifier of lock owner. You can generate a uuid as lock_owner.For example,in golang: req.LockOwner = uuid.New().String() This field is per request,not per process,so it is different for each request, which aims to prevent multi-thread in the same process trying the same lock concurrently. The reason why we don't make it automatically generated is: 1. If it is automatically generated,there must be a 'my_lock_owner_id' field in the response. This name is so weird that we think it is inappropriate to put it into the api spec 2. If we change the field 'my_lock_owner_id' in the response to 'lock_owner',which means the current lock owner of this lock, we find that in some lock services users can't get the current lock owner.Actually users don't need it at all. 3. When reentrant lock is needed,the existing lock_owner is required to identify client and check "whether this client can reenter this lock". So this field in the request shouldn't be removed. |
| expire | [int32](#int32) |  | Required. expire is the time before expire.The time unit is second. |






<a name="spec.proto.runtime.v1.TryLockResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## TryLockResponse
Lock response message returns is the lock obtained.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| success | [bool](#bool) |  | Is lock success |






<a name="spec.proto.runtime.v1.UnlockRequest"></a>
<p align="right"><a href="#top">Top</a></p>

## UnlockRequest
UnLock request message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| store_name | [string](#string) |  | The name of store |
| resource_id | [string](#string) |  | resource_id is the lock key. |
| lock_owner | [string](#string) |  | The owner of the lock |






<a name="spec.proto.runtime.v1.UnlockResponse"></a>
<p align="right"><a href="#top">Top</a></p>

## UnlockResponse
UnLock response message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| status | [UnlockResponse.Status](#spec.proto.runtime.v1.UnlockResponse.Status) |  | The status of unlock |





 <!-- end messages -->


<a name="spec.proto.runtime.v1.HTTPExtension.Verb"></a>

## HTTPExtension.Verb
The enum of http reuest method

| Name | Number | Description |
| ---- | ------ | ----------- |
| NONE | 0 | NONE |
| GET | 1 | GET method |
| HEAD | 2 | HEAD method |
| POST | 3 | POST method |
| PUT | 4 | PUT method |
| DELETE | 5 | DELETE method |
| CONNECT | 6 | CONNECT method |
| OPTIONS | 7 | CONNECT method |
| TRACE | 8 | CONNECT method |
| PATCH | 9 | PATCH method |



<a name="spec.proto.runtime.v1.SequencerOptions.AutoIncrement"></a>

## SequencerOptions.AutoIncrement
requirements for auto-increment guarantee

| Name | Number | Description |
| ---- | ------ | ----------- |
| WEAK | 0 | (default) WEAK means a "best effort" incrementing service.But there is no strict guarantee of global monotonically increasing. The next id is "probably" greater than current id. |
| STRONG | 1 | STRONG means a strict guarantee of global monotonically increasing. The next id "must" be greater than current id. |



<a name="spec.proto.runtime.v1.StateOptions.StateConcurrency"></a>

## StateOptions.StateConcurrency
Enum describing the supported concurrency for state.
The API server uses Optimized Concurrency Control (OCC) with ETags.
When an ETag is associated with an save or delete request, the store shall allow the update only if the attached ETag matches with the latest ETag in the database.
But when ETag is missing in the write requests, the state store shall handle the requests in the specified strategy(e.g. a last-write-wins fashion).

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONCURRENCY_UNSPECIFIED | 0 | Concurrency state is unspecified |
| CONCURRENCY_FIRST_WRITE | 1 | First write wins |
| CONCURRENCY_LAST_WRITE | 2 | Last write wins |



<a name="spec.proto.runtime.v1.StateOptions.StateConsistency"></a>

## StateOptions.StateConsistency
Enum describing the supported consistency for state.

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONSISTENCY_UNSPECIFIED | 0 | Consistency state is unspecified |
| CONSISTENCY_EVENTUAL | 1 | The API server assumes data stores are eventually consistent by default.A state store should: - For read requests, the state store can return data from any of the replicas - For write request, the state store should asynchronously replicate updates to configured quorum after acknowledging the update request. |
| CONSISTENCY_STRONG | 2 | When a strong consistency hint is attached, a state store should: - For read requests, the state store should return the most up-to-date data consistently across replicas. - For write/delete requests, the state store should synchronisely replicate updated data to configured quorum before completing the write request. |



<a name="spec.proto.runtime.v1.UnlockResponse.Status"></a>

## UnlockResponse.Status
The enum of unlock status

| Name | Number | Description |
| ---- | ------ | ----------- |
| SUCCESS | 0 | Unlock is success |
| LOCK_UNEXIST | 1 | The lock is not exist |
| LOCK_BELONG_TO_OTHERS | 2 | The lock is belong to others |
| INTERNAL_ERROR | 3 | Internal error |


 <!-- end enums -->

 <!-- end HasExtensions -->

