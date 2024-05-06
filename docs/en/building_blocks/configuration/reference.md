# Configuration API

## What is Configuration API
Configuration API provides a set of API for developers to manage the configuration data of their applications. Developers can use this API to perform CRUD operations on configuration data and subscribe to change notifications for the configuration data.  

## How to use Configuration API
You can call the Configuration API through grpc. Configuration API is defined in [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto).  

Before using the component, it is necessary to configure it according to the different configuration centers. For example, if you are using etcd, you can refer to the [ETCD Configuration Documentation](en/component_specs/configuration/etcd.md).  

### Usage Example
An example of calling the Configuration API through the SDK can be found in the following **Quick Start Examples**:  
- [Apollo Configuration Center](en/start/configuration/start-apollo.md)
- [Etcd Configuration Center](en/start/configuration/start.md)
- [Nacos Configuration Center](en/start/confguration/start-nacos.md) 

## Configuration API Introduction

To avoid inconsistencies between the documentation and the code, please refer to [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) for detailed input parameters and return values.

### GetConfiguration
Used to get configuration data  
```protobuf
  // GetConfiguration gets configuration from configuration store.
  rpc GetConfiguration(GetConfigurationRequest) returns (GetConfigurationResponse) {}
```

Input parameter：  
```protobuf
message GetConfigurationRequest {
  // The name of configuration store.
  string store_name = 1;

  // The application id which
  // Only used for admin, Ignored and reset for normal client
  string app_id = 2;

  // The group of keys.
  string group = 3;

  // The label for keys.
  string label = 4;

  // The keys to get.
  repeated string keys = 5;

  // The metadata which will be sent to configuration store components.
  map<string, string> metadata = 6;

  // Subscribes update event for given keys.
  // If true, when any configuration item in this request is updated, app will receive event by OnConfigurationEvent() of app callback
  bool subscribe_update = 7;
}
```

Return value：  
```protobuf
// GetConfigurationResponse is the response conveying the list of configuration values.
message GetConfigurationResponse {
  // The list of items containing configuration values.
  repeated ConfigurationItem items = 1;
}

// ConfigurationItem represents a configuration item with key, content and other information.
message ConfigurationItem {
  // Required. The key of configuration item
  string key = 1;

  // The content of configuration item
  // Empty if the configuration is not set, including the case that the configuration is changed from value-set to value-not-set.
  string content = 2;

  // The group of configuration item.
  string group = 3;

  // The label of configuration item.
  string label = 4;

  // The tag list of configuration item.
  map<string, string> tags = 5;

  // The metadata which will be passed to configuration store component.
  map<string, string> metadata = 6;
}
```

### SaveConfiguration
Used to save configuration data  
```protobuf
// SaveConfiguration saves configuration into configuration store.
  rpc SaveConfiguration(SaveConfigurationRequest) returns (google.protobuf.Empty) {}
```

Input parameter：  
```protobuf
message SaveConfigurationRequest {
  // The name of configuration store.
  string store_name = 1;

  // The application id which
  // Only used for admin, ignored and reset for normal client
  string app_id = 2;

  // The list of configuration items to save.
  // To delete a exist item, set the key (also label) and let content to be empty
  repeated ConfigurationItem items = 3;

  // The metadata which will be sent to configuration store components.
  map<string, string> metadata = 4;
}
```

Return value：  
`google.protobuf.Empty`

### DeleteConfiguration
Used to delete configuration data  
```protobuf
  // DeleteConfiguration deletes configuration from configuration store.
  rpc DeleteConfiguration(DeleteConfigurationRequest) returns (google.protobuf.Empty) {}
```

Input parameter:  
```protobuf
// DeleteConfigurationRequest is the message to delete a list of key-value configuration from specified configuration store.
message DeleteConfigurationRequest {
  // The name of configuration store.
  string store_name = 1;

  // The application id which
  // Only used for admin, Ignored and reset for normal client
  string app_id = 2;

  // The group of keys.
  string group = 3;

  // The label for keys.
  string label = 4;

  // The keys to get.
  repeated string keys = 5;

  // The metadata which will be sent to configuration store components.
  map<string, string> metadata = 6;
}
```

Return value：  
`google.protobuf.Empty`

### SubscribeConfiguration
For subscribing to changes in configuration data  
```protobuf
  // SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
  rpc SubscribeConfiguration(stream SubscribeConfigurationRequest) returns (stream SubscribeConfigurationResponse) {}
```

Input parameter：  
```protobuf
// SubscribeConfigurationRequest is the message to get a list of key-value configuration from specified configuration store.
message SubscribeConfigurationRequest {
  // The name of configuration store.
  string store_name = 1;

  // The application id which
  // Only used for admin, ignored and reset for normal client
  string app_id = 2;

  // The group of keys.
  string group = 3;

  // The label for keys.
  string label = 4;

  // The keys to get.
  repeated string keys = 5;

  // The metadata which will be sent to configuration store components.
  map<string, string> metadata = 6;
}
```

Return value：  
```protobuf
// SubscribeConfigurationResponse is the response conveying the list of configuration values.
message SubscribeConfigurationResponse {
  // The name of configuration store.
  string store_name = 1;

  // The application id.
  // Only used for admin client.
  string app_id = 2;

  // The list of items containing configuration values.
  repeated ConfigurationItem items = 3;
}
```