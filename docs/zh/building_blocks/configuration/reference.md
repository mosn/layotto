# Configuration API

## 什么是 Configuration API
Configuration API 为开发者提供了一组用于管理应用程序的配置数据的 API，开发者可以通过使用该 API 对配置数据进行增删改查以及订阅相关配置数据的更新。  

## 如何使用 Configuration API
Configuration API 定义在 [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto) 中。应用可以通过grpc调用对应的 Configuration API，实现对应用程序的配置数据的操作。  

在使用前需要对组件进行配置，根据不同的配置中心进行不同的配置，比如使用 etcd 可以参考 [ETCD配置文档](zh/component_specs/configuration/etcd.md)。  

### 使用示例
通过 SDK 调用 Configuration API 的示例可以参考以下**快速开始示例**：
- [Apollo 配置中心](zh/start/configuration/start-apollo.md)
- [Etcd 配置中心](zh/start/configuration/start.md)
- [Nacos 配置中心](zh/start/confguration/start-nacos.md) 

## Configuration API 介绍

为避免文档和代码不一致，详细入参和返回值请参考 [runtime.proto](https://github.com/mosn/layotto/blob/main/spec/proto/runtime/v1/runtime.proto)。

### GetConfiguration
用于获取配置数据  
```protobuf
  // GetConfiguration gets configuration from configuration store.
  rpc GetConfiguration(GetConfigurationRequest) returns (GetConfigurationResponse) {}
```

入参：  
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

返回参数：  
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
用于保存配置数据  
```protobuf
// SaveConfiguration saves configuration into configuration store.
  rpc SaveConfiguration(SaveConfigurationRequest) returns (google.protobuf.Empty) {}
```

入参：  
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

返回参数：  
`google.protobuf.Empty`

### DeleteConfiguration
用于删除配置数据  
```protobuf
  // DeleteConfiguration deletes configuration from configuration store.
  rpc DeleteConfiguration(DeleteConfigurationRequest) returns (google.protobuf.Empty) {}
```

入参：  
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

返回参数：  
`google.protobuf.Empty`

### SubscribeConfiguration
用于订阅配置数据的变更
```protobuf
  // SubscribeConfiguration gets configuration from configuration store and subscribe the updates.
  rpc SubscribeConfiguration(stream SubscribeConfigurationRequest) returns (stream SubscribeConfigurationResponse) {}
```

入参：  
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

返回参数：  
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